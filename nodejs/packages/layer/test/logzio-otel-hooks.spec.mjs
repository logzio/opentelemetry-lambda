import assert from 'assert';
import { createRequire } from 'module';
import {
  context,
  propagation,
  trace,
  ROOT_CONTEXT,
  TraceFlags,
} from '@opentelemetry/api';
import { W3CTraceContextPropagator } from '@opentelemetry/core';
import { AsyncLocalStorageContextManager } from '@opentelemetry/context-async-hooks';

const require = createRequire(import.meta.url);

const TRACE_ID = '0af7651916cd43dd8448eb211c80319c';
const SPAN_ID = 'b7ad6b7169203331';
const TRACEPARENT = `00-${TRACE_ID}-${SPAN_ID}-01`;

function fakeSpan() {
  const span = {
    attributes: {},
    links: [],
    name: undefined,
    setAttribute(k, v) {
      span.attributes[k] = v;
    },
    setAttributes(a) {
      Object.assign(span.attributes, a);
    },
    updateName(n) {
      span.name = n;
    },
    addLinks(l) {
      span.links.push(...l);
    },
    spanContext: () => ({
      traceId: TRACE_ID,
      spanId: SPAN_ID,
      traceFlags: TraceFlags.SAMPLED,
    }),
  };
  return span;
}

describe('Logz.io hooks', () => {
  let lambdaCfg;
  let awsCfg;

  before(() => {
    require('../logzio/logzio-otel-hooks.cjs');
    lambdaCfg = globalThis.configureLambdaInstrumentation({});
    awsCfg = globalThis.configureAwsInstrumentation({});
  });

  // Other specs reset the global API state, so set it up for every test.
  beforeEach(() => {
    propagation.disable();
    propagation.setGlobalPropagator(new W3CTraceContextPropagator());
    context.disable();
    context.setGlobalContextManager(
      new AsyncLocalStorageContextManager().enable(),
    );
  });

  describe('requestHook', () => {
    it('names API Gateway spans by route and sets both HTTP attribute names', () => {
      const span = fakeSpan();
      lambdaCfg.requestHook(span, {
        event: { httpMethod: 'POST', resource: '/orders' },
      });
      assert.strictEqual(span.name, 'POST /orders');
      assert.strictEqual(span.attributes['http.route'], '/orders');
      assert.strictEqual(span.attributes['http.method'], 'POST');
      assert.strictEqual(span.attributes['http.request.method'], 'POST');
    });

    it('names HTTP API (v2) spans from routeKey', () => {
      const span = fakeSpan();
      lambdaCfg.requestHook(span, {
        event: {
          routeKey: 'GET /orders/{id}',
          requestContext: { http: { method: 'GET' } },
        },
      });
      assert.strictEqual(span.name, 'GET /orders/{id}');
    });

    it('tags Kinesis-triggered invocations with the stream name', () => {
      const span = fakeSpan();
      lambdaCfg.requestHook(span, {
        event: {
          Records: [
            {
              eventSource: 'aws:kinesis',
              eventSourceARN:
                'arn:aws:kinesis:us-east-1:111111111111:stream/telemetry',
              kinesis: { data: Buffer.from('{}').toString('base64') },
            },
          ],
        },
      });
      assert.strictEqual(span.attributes['messaging.system'], 'aws_kinesis');
      assert.strictEqual(
        span.attributes['messaging.destination.name'],
        'telemetry',
      );
      assert.strictEqual(span.attributes['faas.trigger'], 'pubsub');
    });

    it('tags SQS-triggered invocations with the queue name', () => {
      const span = fakeSpan();
      lambdaCfg.requestHook(span, {
        event: {
          Records: [
            {
              eventSource: 'aws:sqs',
              eventSourceARN: 'arn:aws:sqs:us-east-1:111111111111:jobs',
              body: 'not json',
            },
          ],
        },
      });
      assert.strictEqual(span.attributes['messaging.system'], 'aws_sqs');
      assert.strictEqual(span.attributes['messaging.destination.name'], 'jobs');
    });

    it('links every record trace in a batch', () => {
      const span = fakeSpan();
      const body = JSON.stringify({
        _logzio_trace: { traceparent: TRACEPARENT },
      });
      lambdaCfg.requestHook(span, {
        event: { Records: [{ body }, { body }] },
      });
      assert.strictEqual(span.links.length, 2);
      assert.strictEqual(span.attributes['faas.batch.size'], 2);
    });

    it('never sets faas.coldstart (the aws-lambda instrumentation owns it)', () => {
      const span = fakeSpan();
      lambdaCfg.requestHook(span, { event: {} });
      assert.ok(!('faas.coldstart' in span.attributes));
    });
  });

  describe('eventContextExtractor', () => {
    it('continues the trace from API Gateway headers, case-insensitively', () => {
      const ctx = lambdaCfg.eventContextExtractor({
        headers: { TraceParent: TRACEPARENT },
      });
      assert.strictEqual(trace.getSpanContext(ctx)?.traceId, TRACE_ID);
    });

    it('continues the trace from an EventBridge detail field', () => {
      const ctx = lambdaCfg.eventContextExtractor({
        detail: { _logzio_trace: { traceparent: TRACEPARENT } },
      });
      assert.strictEqual(trace.getSpanContext(ctx)?.traceId, TRACE_ID);
    });

    it('continues the trace from a Kinesis record payload', () => {
      const data = Buffer.from(
        JSON.stringify({ _logzio_trace: { traceparent: TRACEPARENT } }),
      ).toString('base64');
      const ctx = lambdaCfg.eventContextExtractor({
        Records: [{ kinesis: { data } }],
      });
      assert.strictEqual(trace.getSpanContext(ctx)?.traceId, TRACE_ID);
    });

    it('returns the root context when nothing carries a trace', () => {
      assert.strictEqual(
        lambdaCfg.eventContextExtractor({ Records: [{ body: 'x' }] }),
        ROOT_CONTEXT,
      );
    });
  });

  describe('preRequestHook', () => {
    const run = (commandName, commandInput) => {
      const span = fakeSpan();
      context.with(trace.setSpan(ROOT_CONTEXT, span), () =>
        awsCfg.preRequestHook(span, {
          moduleVersion: '3.0.0',
          request: { commandName, commandInput },
        }),
      );
      return commandInput;
    };

    it('adds trace context to EventBridge JSON details', () => {
      const input = run('PutEvents', { Entries: [{ Detail: '{"a":1}' }] });
      const detail = JSON.parse(input.Entries[0].Detail);
      assert.strictEqual(detail.a, 1);
      assert.ok(detail._logzio_trace.traceparent.includes(TRACE_ID));
    });

    it('adds trace context to Kinesis JSON records', () => {
      const input = run('PutRecord', { Data: Buffer.from('{"a":1}') });
      const data = JSON.parse(Buffer.from(input.Data).toString());
      assert.ok(data._logzio_trace.traceparent.includes(TRACE_ID));
    });

    it('leaves non-JSON payloads untouched', () => {
      const input = run('PutRecord', { Data: 'plain text' });
      assert.strictEqual(input.Data, 'plain text');
    });
  });

  describe('logs', () => {
    it('emits console output as log records in the active context and flushes them with spans', async () => {
      const emitted = [];
      let flushed = 0;
      globalThis.configureLoggerProvider({
        getLogger: () => ({ emit: r => emitted.push(r) }),
        forceFlush: async () => {
          flushed++;
        },
      });
      const tracerCfg = globalThis.configureTracer({ spanProcessors: [] });
      const original = console.info;
      console.info = () => {};
      try {
        const span = fakeSpan();
        lambdaCfg.requestHook(span, { event: {} }); // wraps console on the first invocation
        context.with(trace.setSpan(ROOT_CONTEXT, span), () =>
          console.info('order %s created', 'o-1'),
        );
      } finally {
        console.info = original;
      }
      assert.strictEqual(emitted.length, 1);
      assert.strictEqual(emitted[0].body, 'order o-1 created');
      assert.strictEqual(emitted[0].severityText, 'INFO');
      assert.strictEqual(
        trace.getSpanContext(emitted[0].context)?.traceId,
        TRACE_ID,
      );
      await tracerCfg.spanProcessors[0].forceFlush();
      assert.strictEqual(flushed, 1);
    });
  });
});

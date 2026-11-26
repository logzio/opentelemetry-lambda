// Logz.io OpenTelemetry hooks for AWS Lambda (Node.js).
// Loaded with NODE_OPTIONS=--require /opt/nodejs/logzio-otel-hooks.cjs, before the OpenTelemetry layer starts.
//
// 1. Names API Gateway spans by route and adds HTTP attributes APM reads.
// 2. Carries trace context across services AWS does not propagate for us:
//    - sending side: adds a `_logzio_trace` field to JSON payloads sent with
//      EventBridge PutEvents, Kinesis PutRecord(s) and IoT Publish. Non-JSON payloads are left untouched.
//      Set LOGZIO_PROPAGATION=off to disable.
//    - receiving side: the invocation continues the first record's trace, and links to every
//      record's trace in a batch (Kinesis, SQS, EventBridge).
// 3. Cold starts: the first invocation of a new Lambda instance gets a `cold start` child span covering
//    the time from process start to the handler call, and `faas.init_duration_ms` on the invocation span.
// 4. Event sources: invocations triggered by Kinesis, SQS, SNS or DynamoDB Streams get messaging.system and
//    messaging.destination.name (stream, queue or topic), so APM shows the trigger as a dependency.
// 5. Logs: console output during an invocation is also sent as OpenTelemetry logs carrying the trace and
//    span IDs, and flushed together with the spans when the invocation ends.
// 6. Limits: the invocation span gets the function's timeout (`faas.timeout_ms`) and the peak memory the
//    instance has used so far (`faas.memory_used_mb`). Configured memory is already on the resource
//    (`faas.max_memory`). Aggregate them over spans (for example max by function) to compare with the limits.
const {
  propagation,
  context,
  trace,
  ROOT_CONTEXT,
} = require('@opentelemetry/api');
const { performance } = require('perf_hooks');
const { format } = require('util');

const FIELD = '_logzio_trace';
const ENABLED = (process.env.LOGZIO_PROPAGATION || 'on') !== 'off';
const getter = { keys: c => Object.keys(c || {}), get: (c, k) => c?.[k] };
const fromCarrier = c =>
  c && c.traceparent ? propagation.extract(ROOT_CONTEXT, c, getter) : null;

// ---------- receiving side ----------
function parseJson(text) {
  try {
    const v = JSON.parse(text);
    return v && typeof v === 'object' && !Array.isArray(v) ? v : null;
  } catch {
    return null;
  }
}
function recordCarrier(r) {
  if (r?.messageAttributes?.traceparent)
    return { traceparent: r.messageAttributes.traceparent.stringValue }; // SQS (SDK-injected)
  if (r?.kinesis?.data)
    return parseJson(Buffer.from(r.kinesis.data, 'base64').toString())?.[FIELD]; // Kinesis
  if (typeof r?.body === 'string') {
    const b = parseJson(r.body);
    return b?.[FIELD] || b?.detail?.[FIELD];
  } // SQS JSON body / EventBridge→SQS
  return null;
}
function eventContextExtractor(event) {
  if (event?.headers) {
    // API Gateway / function URL
    const h = Object.fromEntries(
      Object.entries(event.headers).map(([k, v]) => [k.toLowerCase(), v]),
    );
    return propagation.extract(ROOT_CONTEXT, h, getter);
  }
  if (event?.detail?.[FIELD])
    return fromCarrier(event.detail[FIELD]) || ROOT_CONTEXT; // EventBridge
  for (const r of event?.Records || []) {
    const c = fromCarrier(recordCarrier(r));
    if (c) return c;
  } // first record with context
  return ROOT_CONTEXT;
}

// ---------- cold start ----------
// Provisioned-concurrency instances are initialized ahead of traffic, so no request waits for their init.
let firstInvocation =
  process.env.AWS_LAMBDA_INITIALIZATION_TYPE !== 'provisioned-concurrency';
function recordColdStart(span) {
  if (!firstInvocation) return;
  firstInvocation = false;
  const start = performance.timeOrigin; // Node process start (epoch ms)
  const end = Date.now(); // handler is about to run
  const ms = Math.round(end - start);
  if (!(ms > 0 && ms < 15 * 60 * 1000)) return;
  span.setAttribute('faas.init_duration_ms', ms); // faas.coldstart is set by the aws-lambda instrumentation
  const init = trace.getTracer('logzio-lambda').startSpan(
    'cold start',
    {
      startTime: start,
      attributes: { 'faas.init_duration_ms': ms },
    },
    trace.setSpan(context.active(), span),
  );
  init.end(end);
}

// ---------- event source ----------
// Names the stream, queue or topic that triggered the invocation, so APM shows it as a dependency
// of the consumer (the trigger itself sends no spans).
const SOURCES = {
  'aws:kinesis': ['aws_kinesis', arn => arn.split(':stream/')[1]],
  'aws:sqs': ['aws_sqs', arn => arn.split(':').pop()],
  'aws:sns': ['aws_sns', arn => arn.split(':').pop()],
  'aws:dynamodb': [
    'aws_dynamodb_streams',
    arn => arn.split(':table/')[1]?.split('/')[0],
  ],
};
function recordEventSource(span, event) {
  const r = event?.Records?.[0];
  const src = SOURCES[r?.eventSource || r?.EventSource];
  if (!src) return;
  const arn = r.eventSourceARN || r.EventSubscriptionArn || '';
  const name = arn && src[1](arn);
  span.setAttributes({
    'faas.trigger': 'pubsub',
    'messaging.system': src[0],
    'messaging.operation.type': 'process',
    ...(name ? { 'messaging.destination.name': name } : {}),
  });
}

// ---------- logs ----------
// The Lambda runtime replaces console methods at startup, so they are wrapped on the first invocation.
// Each call still prints as before; it is also emitted as a log record in the active trace context.
const LEVELS = {
  trace: ['TRACE', 1],
  debug: ['DEBUG', 5],
  log: ['INFO', 9],
  info: ['INFO', 9],
  warn: ['WARN', 13],
  error: ['ERROR', 17],
};
let loggerProvider;
let otelLogger;
let consoleWrapped = false;
function wrapConsole() {
  if (consoleWrapped || !loggerProvider) return;
  consoleWrapped = true;
  otelLogger = loggerProvider.getLogger('logzio-lambda-console');
  for (const [method, [severityText, severityNumber]] of Object.entries(
    LEVELS,
  )) {
    const original = console[method];
    if (typeof original !== 'function') continue;
    console[method] = function (...args) {
      original.apply(this, args);
      try {
        otelLogger.emit({
          severityText,
          severityNumber,
          body: format(...args),
          context: context.active(),
        });
      } catch {
        /* logging must never break the function */
      }
    };
  }
}
// A span processor whose only job is to flush logs whenever spans are flushed (end of every invocation).
const flushLogsWithSpans = {
  onStart() {},
  onEnd() {},
  forceFlush: () => loggerProvider?.forceFlush() ?? Promise.resolve(),
  shutdown: () => loggerProvider?.shutdown() ?? Promise.resolve(),
};

// ---------- limits ----------
// Remaining time at the very start of an invocation is the configured timeout (Lambda timeouts are whole seconds).
function recordTimeout(span, event, lambdaContext) {
  const remaining = lambdaContext?.getRemainingTimeInMillis?.();
  if (typeof remaining === 'number' && remaining > 0)
    span.setAttribute('faas.timeout_ms', Math.ceil(remaining / 1000) * 1000);
}
// Peak resident memory of the instance so far (maxRSS is in kilobytes). Close to Lambda's "Max Memory Used",
// which also counts the extension, so slightly lower.
function recordMemory(span) {
  span.setAttribute(
    'faas.memory_used_mb',
    Math.round(process.resourceUsage().maxRSS / 1024),
  );
}

// ---------- span naming + batch links ----------
function requestHook(span, { event, context: lambdaContext }) {
  for (const step of [
    recordColdStart,
    recordEventSource,
    wrapConsole,
    recordTimeout,
  ]) {
    try {
      step(span, event, lambdaContext);
    } catch {
      /* never break the invocation */
    }
  }
  const method = event?.httpMethod || event?.requestContext?.http?.method;
  const route = event?.resource || event?.routeKey?.split(' ')[1];
  if (method && route) {
    span.updateName(`${method} ${route}`);
    // Both old and current names: APM's status panel reads the old ones until APPZ-3884 ships.
    span.setAttributes({
      'http.request.method': method,
      'http.method': method,
      'http.route': route,
      'faas.trigger': 'http',
    });
  }
  const records = event?.Records;
  if (
    Array.isArray(records) &&
    records.length > 1 &&
    typeof span.addLinks === 'function'
  ) {
    const links = [];
    for (const r of records.slice(0, 128)) {
      const c = fromCarrier(recordCarrier(r));
      const sc = c && trace.getSpanContext(c);
      if (sc) links.push({ context: sc });
    }
    if (links.length) span.addLinks(links);
    span.setAttribute('faas.batch.size', records.length);
  }
}
function responseHook(span, { res }) {
  try {
    recordMemory(span);
  } catch {
    /* never break the invocation */
  }
  if (res && typeof res.statusCode === 'number')
    span.setAttributes({
      'http.response.status_code': res.statusCode,
      'http.status_code': res.statusCode,
    });
}

// ---------- sending side ----------
function carrierFor(span) {
  const c = {};
  propagation.inject(trace.setSpan(context.active(), span), c);
  return c.traceparent ? c : null;
}
function addToJsonString(text, carrier) {
  const obj = parseJson(text);
  if (!obj || obj[FIELD]) return null;
  obj[FIELD] = carrier;
  return JSON.stringify(obj);
}
function addToBytes(data, carrier) {
  const text =
    typeof data === 'string' ? data : Buffer.from(data).toString('utf8');
  const out = addToJsonString(text, carrier);
  return out === null
    ? null
    : typeof data === 'string'
      ? out
      : Buffer.from(out);
}
function preRequestHook(span, info) {
  if (!ENABLED) return;
  try {
    const request = info?.request || info; // newer instrumentation passes { moduleVersion, request }
    const input = request?.commandInput;
    const cmd = request?.commandName;
    const carrier = input && carrierFor(span);
    if (!carrier) return;
    if (cmd === 'PutEvents' && Array.isArray(input.Entries)) {
      // EventBridge
      for (const e of input.Entries) {
        const d =
          typeof e.Detail === 'string' && addToJsonString(e.Detail, carrier);
        if (d) e.Detail = d;
      }
    } else if (cmd === 'PutRecord' && input.Data) {
      // Kinesis / Firehose
      const d = addToBytes(input.Data, carrier);
      if (d) input.Data = d;
    } else if (cmd === 'PutRecords' && Array.isArray(input.Records)) {
      for (const r of input.Records) {
        const d = r.Data && addToBytes(r.Data, carrier);
        if (d) r.Data = d;
      }
    } else if (cmd === 'Publish' && input.payload && input.topic) {
      // IoT Core
      const d = addToBytes(input.payload, carrier);
      if (d) input.payload = d;
    }
  } catch {
    /* never break the customer's call */
  }
}

globalThis.configureLambdaInstrumentation = cfg => ({
  ...cfg,
  eventContextExtractor,
  requestHook,
  responseHook,
});
globalThis.configureAwsInstrumentation = cfg => ({ ...cfg, preRequestHook });
globalThis.configureTracer = cfg => ({
  ...cfg,
  spanProcessors: [...(cfg.spanProcessors || []), flushLogsWithSpans],
});
// Called by the layer with its logger provider (only when OTEL_LOGS_EXPORTER is not "none").
globalThis.configureLoggerProvider = provider => {
  loggerProvider = provider;
};

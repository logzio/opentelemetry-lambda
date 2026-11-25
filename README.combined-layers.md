## Combined Layers (New)

**Simplified Deployment**: We now offer combined layers that bundle both the language-specific instrumentation and the collector into a single layer. This approach:
- Reduces the number of layers from 2 to 1
- Simplifies configuration and deployment
- Maintains all the functionality of the separate layers
- Is available for Python, Node.js, Java, Ruby, and Go

### What's included in combined layers:
- **Language-specific OpenTelemetry instrumentation** - Automatically instruments your Lambda function and popular libraries
- **OpenTelemetry Collector** - Built-in collector that exports telemetry data to your configured backend
- **Auto-instrumentation** - Automatic instrumentation for AWS SDK and popular libraries in each language
- **Optimized packaging** - Reduced cold start impact with optimized layer packaging

### Benefits:
- **Single layer deployment** - No need to manage separate collector and instrumentation layers
- **Simplified configuration** - Fewer environment variables and layer configurations
- **Reduced cold start impact** - Optimized packaging reduces overhead
- **Production-ready** - Includes all necessary components for complete observability

### Common Environment Variables

Most combined layers support these common environment variables:

**Required:**
- `AWS_LAMBDA_EXEC_WRAPPER` – set to `/opt/otel-handler` (or language-specific handler)
- `LOGZIO_TRACES_TOKEN` – account token for traces
- `LOGZIO_METRICS_TOKEN` – account token for metrics  
- `LOGZIO_LOGS_TOKEN` – account token for logs
- `LOGZIO_REGION` – Logz.io region code (for example, `us`, `eu`)

**Optional:**
- `OTEL_SERVICE_NAME` – explicit service name
- `OTEL_RESOURCE_ATTRIBUTES` – comma-separated resource attributes (for example, `service.name=my-func,env_id=${LOGZIO_ENV_ID},deployment.environment=${ENVIRONMENT}`)
- `LOGZIO_ENV_ID` – environment identifier you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `env_id=prod`)
- `ENVIRONMENT` – logical environment name you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `deployment.environment=prod`)
- `OPENTELEMETRY_COLLECTOR_CONFIG_URI` – custom collector config URI/file path; defaults to `/opt/collector-config/config.yaml`
- `OPENTELEMETRY_EXTENSION_LOG_LEVEL` – extension log level (`debug`, `info`, `warn`, `error`)

### Language-Specific Details

#### Java Combined Layer
- **Multiple handler types available:**
  - `/opt/otel-handler` - for regular handlers (implementing RequestHandler)
  - `/opt/otel-sqs-handler` - for SQS-triggered functions
  - `/opt/otel-proxy-handler` - for API Gateway proxied handlers
  - `/opt/otel-stream-handler` - for streaming handlers
- **Fast startup mode:** Set `OTEL_JAVA_AGENT_FAST_STARTUP_ENABLED=true` to enable optimized startup (disables JIT tiered compilation level 2)
- **Agent and wrapper variants:** Both Java agent and wrapper approaches are available in the combined layer

#### Node.js Combined Layer
- **Based on OpenTelemetry Node.js layer 0.24.0.** Produces AWS SDK spans (`DynamoDB.PutItem`, `SQS.SendMessage`, `EventBridge.PutEvents`, ...) with `rpc.*`, `db.system.name` and `messaging.*` attributes, including SDK v3 releases that moved to `@smithy/core`.
- **Logz.io hooks** (`/opt/nodejs/logzio-otel-hooks.cjs`, loaded automatically by `/opt/otel-handler`):
  - API Gateway spans are named by route (`GET /orders/{id}`) and get `http.route`, `http.request.method`, `faas.trigger` and the response status code (both current and legacy attribute names).
  - Trace context crosses services AWS does not propagate: a `_logzio_trace` field is added to JSON payloads sent with EventBridge `PutEvents`, Kinesis `PutRecord`/`PutRecords` and IoT `Publish` (non-JSON payloads are left untouched). The receiving function continues that trace, and links every record of a Kinesis/SQS batch. Set `LOGZIO_PROPAGATION=off` to disable.
- **Defaults set by the wrapper:** `OTEL_NODE_ENABLED_INSTRUMENTATIONS` includes `undici` (Node's built-in `fetch`); SDK metrics and logs exporters are off (function logs are shipped by the collector).
- **Gateway mode (APM service pages):** set `OPENTELEMETRY_COLLECTOR_CONFIG_URI=/opt/collector-config/config.gateway.yaml`, `LOGZIO_GATEWAY_ENDPOINT` and `LOGZIO_GATEWAY_TOKEN`. Traces and logs go to a Logz.io gateway collector, which computes APM span metrics and forwards to Logz.io.
- **ESM and CommonJS support:** Works with both module systems
- **Instrumentation control:**
  - `OTEL_NODE_ENABLED_INSTRUMENTATIONS` - comma-separated list to enable only specific instrumentations
  - `OTEL_NODE_DISABLED_INSTRUMENTATIONS` - comma-separated list to disable specific instrumentations
- **Popular libraries included:** AWS SDK v3, Express, HTTP, MongoDB, Redis, and many more

#### Ruby Combined Layer
- **Ruby version support:** Compatible with Ruby 3.2.0, 3.3.0, and 3.4.0
- **Popular gems included:** AWS SDK, Rails, Sinatra, and other popular Ruby libraries
- **Additional configuration:** `OTEL_RUBY_INSTRUMENTATION_NET_HTTP_ENABLED` - toggle net/http instrumentation (true/false)

#### Go Combined Layer
- **Collector-only layer:** Since Go uses manual instrumentation, this provides only the collector component
- **Manual instrumentation required:** You must instrument your Go code using the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/aws/aws-lambda-go/otellambda)
- **No AWS_LAMBDA_EXEC_WRAPPER needed:** Go layer doesn't require the wrapper environment variable

#### Python Combined Layer
- **Auto-instrumentation:** Automatically instruments Lambda functions and popular Python libraries like boto3, requests, urllib3
- **Trace context propagation:** Automatically propagates trace context through AWS services

### Build Scripts
Each language provides a `build-combined.sh` script for creating combined layers:
- `python/src/build-combined.sh`
- `java/build-combined.sh` 
- `nodejs/packages/layer/build-combined.sh`
- `ruby/build-combined.sh`
- `go/build-combined.sh`

For detailed build instructions and sample applications, see the individual language README files below.

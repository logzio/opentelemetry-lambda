# OpenTelemetry Lambda

![GitHub Java Workflow Status](https://img.shields.io/github/actions/workflow/status/open-telemetry/opentelemetry-lambda/ci-java.yml?branch%3Amain&label=CI%20%28Java%29&style=for-the-badge)
![GitHub Collector Workflow Status](https://img.shields.io/github/actions/workflow/status/open-telemetry/opentelemetry-lambda/ci-collector.yml?branch%3Amain&label=CI%20%28Collector%29&style=for-the-badge)
![GitHub NodeJS Workflow Status](https://img.shields.io/github/actions/workflow/status/open-telemetry/opentelemetry-lambda/ci-nodejs.yml?branch%3Amain&label=CI%20%28NodeJS%29&style=for-the-badge)
![GitHub Terraform Lint Workflow Status](https://img.shields.io/github/actions/workflow/status/open-telemetry/opentelemetry-lambda/ci-terraform.yml?branch%3Amain&label=CI%20%28Terraform%20Lint%29&style=for-the-badge)
![GitHub Python Pull Request Workflow Status](https://img.shields.io/github/actions/workflow/status/open-telemetry/opentelemetry-lambda/ci-python.yml?branch%3Amain&label=Pull%20Request%20%28Python%29&style=for-the-badge)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/open-telemetry/opentelemetry-lambda/badge?style=for-the-badge)](https://scorecard.dev/viewer/?uri=github.com/open-telemetry/opentelemetry-lambda)

## OpenTelemetry Lambda Layers

The OpenTelemetry Lambda Layers provide the OpenTelemetry (OTel) code to export telemetry asynchronously from AWS Lambda functions. It does this by embedding a stripped-down version of [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) inside an [AWS Lambda Extension Layer](https://aws.amazon.com/blogs/compute/introducing-aws-lambda-extensions-in-preview/). This allows Lambda functions to use OpenTelemetry to send traces and metrics to any configured backend.

There are 2 types of lambda layers
1. Collector Layer - Embeds a stripped down version of the OpenTelemetry Collector
2. Language Specific Layer - Includes language specific nuances to allow lambda functions to automatically consume context from upstream callers, create spans, and automatically instrument the AWS SDK

These 2 layers are meant to be used in conjunction to instrument your lambda functions. The reason that the collector is not embedded in specific language layers is to give users flexibility

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

## Collector Layer
* ### [Collector Lambda Layer](collector/README.md)

## Extension Layer Language Support
* ### [Python Lambda Layer](python/README.md) - *Combined layer available*
* ### [Java Lambda Layer](java/README.md) - *Combined layer available*
* ### [NodeJS Lambda Layer](nodejs/README.md) - *Combined layer available*
* ### [Ruby Lambda Layer](ruby/README.md) - *Combined layer available*

## Additional language tooling not currently supported
* ### [Go Lambda Library](go/README.md) - *Combined layer available (collector only)*
* ### [.NET Lambda Layer](dotnet/README.md) 

## Latest Layer Versions
| Name         | ARN                                                                                                                    | Version |
|--------------|:-----------------------------------------------------------------------------------------------------------------------|:--------|
| collector    | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-collector-<amd64\|arm64>-<version>:1` | ![Collector](https://api.globadge.com/v1/badgen/http/jq/e3309d56-dfd6-4dae-ac00-4498070d84f0) |
| nodejs       | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-nodejs-<version>:1` | ![NodeJS](https://api.globadge.com/v1/badgen/http/jq/91b0f102-25fc-425f-8de9-f05491b9f757) |
| python       | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-python-<version>:1` | ![Python](https://api.globadge.com/v1/badgen/http/jq/ab030ce1-ee7d-4c14-b643-eb20ec050e0b) |
| java-agent   | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-javaagent-<version>:1` | ![Java Agent](https://api.globadge.com/v1/badgen/http/jq/301ad852-ccb4-4bb4-997e-60282ad11f71) |
| java-wrapper | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-javawrapper-<version>:1` | ![Java Wrapper](https://api.globadge.com/v1/badgen/http/jq/e10281c6-3d0e-42e4-990b-7a725301bef4) |
| ruby         | `arn:aws:lambda:<region>:184161586896:layer:opentelemetry-ruby-dev-<version>:1` | ![Ruby](https://api.globadge.com/v1/badgen/http/jq/4d9b9e93-7d6b-4dcf-836e-1878de566fdb) |

## FAQ

* **What exporters/receivers/processors are included from the OpenTelemetry Collector?**
    > You can check out [the stripped-down collector's imports](https://github.com/open-telemetry/opentelemetry-lambda/blob/main/collector/lambdacomponents/default.go#L18) in this repository for a full list of currently included components.

    > Self-built binaries of the collector have **experimental** support for a custom set of connectors/exporters/receivers/processors. For more information, see [(Experimental) Customized collector build](./collector/README.md#experimental-customized-collector-build)
* **Is the Lambda layer provided or do I need to build it and distribute it myself?**
    > This repository provides pre-built Lambda layers, their ARNs are available in the [Releases](https://github.com/open-telemetry/opentelemetry-lambda/releases). You can also build the layers manually and publish them in your AWS account. This repo has files to facilitate doing that. More information is provided in [the Collector folder's README](collector/README.md).

## Design Proposal

To get a better understanding of the proposed design for the OpenTelemetry Lambda extension, you can see the [Design Proposal here.](docs/design_proposal.md)

## Features

The following is a list of features provided by the OpenTelemetry layers.

### OpenTelemetry collector

The layer includes the OpenTelemetry Collector as a Lambda extension.

### Custom context propagation carrier extraction

Context can be propagated through various mechanisms (e.g. http headers (APIGW), message attributes (SQS), ...). In some cases, it may be required to pass a custom context propagation extractor in Lambda through configuration, this feature allows this through Lambda instrumentation configuration.

### X-Ray Env Var Span Link

This links a context extracted from the Lambda runtime environment to the instrumentation-generated span rather than disabling that context extraction entirely.

### Semantic conventions

The Lambda language implementation follows the semantic conventions specified in the OpenTelemetry Specification.

### Auto instrumentation

The Lambda layer includes support for automatically instrumentation code via the use of instrumentation libraries.

### Flush TracerProvider

The Lambda instrumentation will flush the `TracerProvider` at the end of an invocation.

### Flush MeterProvider

The Lambda instrumentation will flush the `MeterProvider` at the end of an invocation.

### Support matrix

The table below captures the state of various features and their levels of support different runtimes.

| Feature                    | Node | Python | Java | .NET | Go   | Ruby |
| -------------------------- | :--: | :----: | :--: | :--: | :--: | :--: |
| OpenTelemetry collector    |  +   |  +     |  +   |  +   |  +   |  +   |
| Custom context propagation |  +   |  -     |  -   |  -   | N/A  |  +   |
| X-Ray Env Var Span Link    |  -   |  -     |  -   |  -   | N/A  |  -   |
| Semantic Conventions^      |      |  +     |  +   |  +   | N/A  |  +   |
| - Trace General^<sup>[1]</sup>           |  +   |        |  +   |  +   | N/A  |   +  |
| - Trace Incoming^<sup>[2]</sup>          |  -   |        |  -   |  +   | N/A  |   -  |
| - Trace Outgoing^<sup>[3]</sup>          |  +   |        |  -   |  +   | N/A  |   +  |
| - Metrics^<sup>[4]</sup>                 |  -   |        |  -   |  -   | N/A  |   -  |
| Auto instrumentation       |  +   |   +    |  +   |  -   | N/A  |   +  |
| Flush TracerProvider       |  +   |   +    |      |  +   |  +   |   +  |
| Flush MeterProvider        |  +   |   +    |      |      |      |   -  |

#### Legend

* `+` is supported
* `-` not supported
* `^` subject to change depending on spec updates
* `N/A` not applicable to the particular language
* blank cell means the status of the feature is not known.

The following are runtimes which are no longer or not yet supported by this repository:

* Node.js 12, Node.js 16 - not [officially supported](https://github.com/open-telemetry/opentelemetry-js#supported-runtimes) by OpenTelemetry JS

[1]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/faas/faas-spans.md#general-attributes
[2]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/faas/faas-spans.md#incoming-invocations
[3]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/faas/faas-spans.md#outgoing-invocations
[4]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/faas/faas-metrics.md

## Contributing

See the [Contributing Guide](CONTRIBUTING.md) for details.

### Maintainers

- [Serkan Özal](https://github.com/serkan-ozal), Catchpoint
- [Tyler Benson](https://github.com/tylerbenson), ServiceNow

For more information about the maintainer role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#maintainer).

### Approvers

- [Ivan Santos](https://github.com/pragmaticivan)
- [Warre Pessers](https://github.com/wpessers)

For more information about the approver role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#approver).

### Emeritus Maintainers

- [Alex Boten](https://github.com/codeboten)
- [Anthony Mirabella](https://github.com/Aneurysm9)
- [Raphael Philipe Mendes da Silva](https://github.com/rapphil)

For more information about the emeritus role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#emeritus-maintainerapprovertriager).

### Emeritus Approvers

- [Lei Wang](https://github.com/wangzlei)
- [Nathaniel Ruiz Nowell](https://github.com/NathanielRN)
- [Tristan Sloughter](https://github.com/tsloughter)

- Maintainers ([@open-telemetry/lambda-extension-maintainers](https://github.com/orgs/open-telemetry/teams/lambda-extension-maintainers)):

  - [Raphael Philipe Mendes da Silva](https://github.com/rapphil), AWS
  - [Serkan Özal](https://github.com/serkan-ozal), Catchpoint
  - [Tyler Benson](https://github.com/tylerbenson), Lightstep

- Emeritus Maintainers:

  - [Alex Boten](https://github.com/codeboten)
  - [Anthony Mirabella](https://github.com/Aneurysm9)

Learn more about roles in the [community repository](https://github.com/open-telemetry/community/blob/main/community-membership.md).

# Configuration Example

Replace `<<LOGZIO_TRACING_SHIPPING_TOKEN>>`, `<<LOGZIO_SPM_SHIPPING_TOKEN>>`, `<<LOGZIO_ACCOUNT_REGION_CODE>>`, and `<<LOGZIO_LISTENER_HOST>>` with your Logz.io account's information.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
      http:
        endpoint: "0.0.0.0:4318"

connectors:
  spanmetrics:
    aggregation_temporality: AGGREGATION_TEMPORALITY_CUMULATIVE
    dimensions:
      - name: rpc.grpc.status_code
      - name: http.method
      - name: http.status_code
      - name: cloud.provider
      - name: cloud.region
      - name: db.system
      - name: messaging.system
      - default: DEV
        name: env_id
    dimensions_cache_size: 100000
    histogram:
      explicit:
        buckets:
          - 2ms
          - 8ms
          - 50ms
          - 100ms
          - 200ms
          - 500ms
          - 1s
          - 5s
          - 10s
    metrics_expiration: 5m
    resource_metrics_key_attributes:
      - service.name
      - telemetry.sdk.language
      - telemetry.sdk.name

exporters:
  logzio/traces:
    account_token: <<LOGZIO_TRACING_SHIPPING_TOKEN>>
    region: <<LOGZIO_ACCOUNT_REGION_CODE>>
  prometheusremotewrite/spm:
    endpoint: https://<<LOGZIO_LISTENER_HOST>>:8053
    add_metric_suffixes: false
    headers:
      Authorization: Bearer <<LOGZIO_SPM_SHIPPING_TOKEN>>

processors:
  batch:
  tail_sampling:
    policies:
      - name: policy-errors
        type: status_code
        status_code: {status_codes: [ERROR]}
      - name: policy-slow
        type: latency
        latency: {threshold_ms: 1000}
      - name: policy-random-ok
        type: probabilistic
        probabilistic: {sampling_percentage: 10}
  metricstransform/metrics-rename:
    transforms:
    - include: ^duration(.*)$$
      action: update
      match_type: regexp
      new_name: latency.$${1} 
    - action: update
      include: calls
      new_name: calls_total
  metricstransform/labels-rename:
    transforms:
    - action: update
      include: ^latency
      match_type: regexp
      operations:
      - action: update_label
        label: span.name
        new_label: operation
    - action: update
      include: ^calls
      match_type: regexp
      operations:
      - action: update_label
        label: span.name
        new_label: operation  

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [tail_sampling, batch]
      exporters: [logzio/traces]
    traces/spm:
      receivers: [otlp]
      processors: [batch]
      exporters: [spanmetrics]
    metrics/spanmetrics:
      receivers: [spanmetrics]
      processors: [metricstransform/metrics-rename, metricstransform/labels-rename, batch]
      exporters: [prometheusremotewrite/spm]
  telemetry: 
    logs:
      level: "info"
```
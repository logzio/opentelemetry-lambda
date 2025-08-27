# OpenTelemetry Lambda Python

Scripts and files used to build AWS Lambda Layers for running OpenTelemetry on AWS Lambda for Python.

## Combined OpenTelemetry Lambda Layer

**New**: We now offer a simplified deployment option with a combined layer that bundles both the OpenTelemetry Python instrumentation and the collector into a single layer. This reduces the number of layers you need to manage and simplifies your Lambda function configuration.

### What's included in the combined layer:
- **Python OpenTelemetry instrumentation** - Automatically instruments your Lambda function and common Python libraries
- **OpenTelemetry Collector** - Built-in collector that exports telemetry data to your configured backend
- **Auto-instrumentation for popular libraries** - Automatic instrumentation for libraries like boto3, requests, urllib3, and more
- **Trace context propagation** - Automatically propagates trace context through AWS services

### Benefits:
- **Single layer deployment** - No need to manage separate collector and instrumentation layers
- **Simplified configuration** - Fewer environment variables and layer configurations
- **Reduced cold start impact** - Optimized packaging reduces overhead
- **Production-ready** - Includes all necessary components for complete observability

### Usage:
To use the combined layer, simply add it to your Lambda function and set the `AWS_LAMBDA_EXEC_WRAPPER` environment variable:
```
AWS_LAMBDA_EXEC_WRAPPER=/opt/otel-instrument
```

The layer handles both instrumentation and telemetry export automatically. For detailed build instructions, see the build script at `python/src/build-combined.sh` in this repository.

### Environment variables

Required:
- `AWS_LAMBDA_EXEC_WRAPPER` – set to `/opt/otel-instrument`
- `LOGZIO_TRACES_TOKEN` – account token for traces
- `LOGZIO_METRICS_TOKEN` – account token for metrics
- `LOGZIO_LOGS_TOKEN` – account token for logs
- `LOGZIO_REGION` – Logz.io region code (for example, `us`, `eu`)

Optional:
- `OTEL_SERVICE_NAME` – explicit service name
- `OTEL_RESOURCE_ATTRIBUTES` – comma-separated resource attributes (for example, `service.name=my-func,env_id=${LOGZIO_ENV_ID},deployment.environment=${ENVIRONMENT}`)
- `LOGZIO_ENV_ID` – environment identifier you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `env_id=prod`)
- `ENVIRONMENT` – logical environment name you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `deployment.environment=prod`)

## Sample App 

1. Install
   * [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
   * [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
   * [Go](https://go.dev/doc/install)
   * [Docker](https://docs.docker.com/get-docker)
2. Run aws configure to [set aws credential(with administrator permissions)](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install-mac.html#serverless-sam-cli-install-mac-iam-permissions) and default region.
3. Download a local copy of this repository from Github.
4. Navigate to the path `cd python/src`
5. If you just want to create a zip file with the OpenTelemetry Python AWS Lambda layer, then use the `-b true` option: `bash run.sh -n <LAYER_NAME_HERE> -b true`
6. If you want to create the layer and automatically publish it, use no options: `bash run.sh`

# OpenTelemetry Lambda Ruby

Scripts and files used to build AWS Lambda Layers for running OpenTelemetry on AWS Lambda for Ruby.

## Combined OpenTelemetry Lambda Layer

**New**: We now offer a simplified deployment option with a combined layer that bundles both the OpenTelemetry Ruby instrumentation and the collector into a single layer. This reduces the number of layers you need to manage and simplifies your Lambda function configuration.

### What's included in the combined layer:
- **Ruby OpenTelemetry instrumentation** - Automatically instruments your Lambda function
- **OpenTelemetry Collector** - Built-in collector that exports telemetry data to your configured backend
- **Auto-instrumentation for popular gems** - Includes instrumentation for AWS SDK, Rails, Sinatra, and many other popular Ruby libraries
- **Support for Ruby 3.2.0, 3.3.0, and 3.4.0** - Compatible with recent Ruby versions

### Benefits:
- **Single layer deployment** - No need to manage separate collector and instrumentation layers
- **Simplified configuration** - Fewer environment variables and layer configurations
- **Reduced cold start impact** - Optimized packaging with stripped binaries and pruned gem files
- **Production-ready** - Includes all necessary components for complete observability

### Usage:
To use the combined layer, add it to your Lambda function and set the `AWS_LAMBDA_EXEC_WRAPPER` environment variable:
```
AWS_LAMBDA_EXEC_WRAPPER=/opt/otel-handler
```

The layer handles both instrumentation and telemetry export automatically. For detailed build instructions, see the build script at `ruby/build-combined.sh` in this repository.

### Environment variables

Required:
- `AWS_LAMBDA_EXEC_WRAPPER` – set to `/opt/otel-handler`
- `LOGZIO_TRACES_TOKEN` – account token for traces
- `LOGZIO_METRICS_TOKEN` – account token for metrics
- `LOGZIO_LOGS_TOKEN` – account token for logs
- `LOGZIO_REGION` – Logz.io region code (for example, `us`, `eu`)

Optional:
- `OTEL_SERVICE_NAME` – explicit service name
- `OTEL_RESOURCE_ATTRIBUTES` – comma-separated resource attributes (for example, `service.name=my-func,env_id=${LOGZIO_ENV_ID},deployment.environment=${ENVIRONMENT}`)
- `OTEL_RUBY_INSTRUMENTATION_NET_HTTP_ENABLED` – toggle net/http instrumentation (true/false)
- `LOGZIO_ENV_ID` – environment identifier you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `env_id=prod`)
- `ENVIRONMENT` – logical environment name you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `deployment.environment=prod`)

## Requirement
* Ruby 3.2.0/3.3.0/3.4.0
* [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
* [Go](https://go.dev/doc/install)
* [Docker](https://docs.docker.com/get-docker)

**Building Lambda Ruby Layer With OpenTelemetry Ruby Dependencies**

1. Run build script

```bash
./build.sh
```

Layer is stored in `src/build` folder

**Default GEM_PATH**

The [default GEM_PATH](https://docs.aws.amazon.com/lambda/latest/dg/ruby-package.html#ruby-package-dependencies-layers) for aws lambda ruby is `/opt/ruby/gems/<ruby_vesion>` after lambda function loads this layer.

**Define AWS_LAMBDA_EXEC_WRAPPER**

Point `AWS_LAMBDA_EXEC_WRAPPER` to `/opt/otel-handler` to take advantage of layer wrapper that load all opentelemetry ruby components
e.g.
```
AWS_LAMBDA_EXEC_WRAPPER: /opt/otel-handler
```

#### There are two ways to define the AWS_LAMBDA_EXEC_WRAPPER that point to either binary executable or script (normally bash).

Method 1: define the AWS_LAMBDA_EXEC_WRAPPER in function from template.yml
```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: 'AWS::Serverless-2016-10-31'
Description: OpenTelemetry Ruby Lambda layer for Ruby
Parameters:
  LayerName:
    ...
Resources:
  OTelLayer:
    ...
  api:
    ...
  function:
    Type: AWS::Serverless::Function
    Properties:
      ...
      Environment:
        Variables:
          AWS_LAMBDA_EXEC_WRAPPER: /opt/otel-handler # this is an example of the path

```

Method 2: directly update the environmental variable in lambda console: Configuration -> Environemntal variables

For more information about aws lambda wrapper and wrapper layer, check [aws lambda runtime-wrapper](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-modify.html#runtime-wrapper). We provide a sample wrapper file in `src/layer/otel-handler` as reference.

### Sample App

1. Make sure the requirements are met (e.g. sam, aws, docker, ruby version.). Current sample app only support testing Ruby 3.2.0. If you wish to play with other ruby version, please modify ruby version from Runtime in sample-apps/template.yml and src/otel/layer/Makefile.

2. Navigate to the path `cd ruby/src` to build layer

```bash
sam build -u -t template.yml
```

3. Navigate to the path `cd ruby/sample-apps`
4. Build the layer and function based on template.yml. You will see .aws-sam folder after executed the command
```bash
sam build -u -t template.yml
# for different arch, define it in properties from template.yml
   # Architectures:
    #   - arm64
```
5. Test with local simulation
```bash
sam local start-api --skip-pull-image
```

6. curl the lambda function
```bash
curl http://127.0.0.1:3000
# you should expect: Hello 1.4.1
```
In this sample-apps, we use `src/layer/otel-handler` as default `AWS_LAMBDA_EXEC_WRAPPER`; to change it, please edit in `sample-apps/template.yml`

In `ruby/sample-apps/template.yml`, the OTelLayer -> Properties -> ContentUri is pointing to `ruby/src/layer/`. This is for local testing purpose. If you wish to deploy (e.g. `sam deploy`), please point it to correct location or zip file.

### Test with Jaeger Endpoint

Assume you have a lambda function with current [released layer](https://github.com/open-telemetry/opentelemetry-lambda/releases/tag/layer-ruby%2F0.1.0), and you want to test it out that send trace to jaeger endpoint, below should be your environmental variable.
```
AWS_LAMBDA_EXEC_WRAPPER=/opt/otel-handler
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://<jaeger_endpoint:port_number>/v1/traces
```
Try with `jaeger-all-in-one` at [Jaeger](https://www.jaegertracing.io/docs/1.57/getting-started/)



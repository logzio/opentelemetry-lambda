# OpenTelemetry Lambda NodeJS

Layer for running NodeJS applications on AWS Lambda with OpenTelemetry. Adding the layer and pointing to it with
the `AWS_LAMBDA_EXEC_WRAPPER` environment variable will initialize OpenTelemetry, enabling tracing with no code change.

## Combined OpenTelemetry Lambda Layer

**New**: We now offer a simplified deployment option with a combined layer that bundles both the OpenTelemetry Node.js instrumentation and the collector into a single layer. This reduces the number of layers you need to manage and simplifies your Lambda function configuration.

### What's included in the combined layer:
- **Node.js OpenTelemetry instrumentation** - Automatically instruments your Lambda function
- **OpenTelemetry Collector** - Built-in collector that exports telemetry data to your configured backend
- **Auto-instrumentation for popular libraries** - Includes AWS SDK v3 and a subset of popular Node.js libraries
- **ESM and CommonJS support** - Works with both module systems

### Benefits:
- **Single layer deployment** - No need to manage separate collector and instrumentation layers
- **Simplified configuration** - Fewer environment variables and layer configurations
- **Reduced cold start impact** - Optimized packaging reduces overhead
- **Production-ready** - Includes all necessary components for complete observability

### Usage:
To use the combined layer, add it to your Lambda function and set `AWS_LAMBDA_EXEC_WRAPPER` to `/opt/otel-handler`.

For detailed build instructions, see the build script at `nodejs/packages/layer/build-combined.sh` in this repository.

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
- `LOGZIO_ENV_ID` – environment identifier you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `env_id=prod`)
- `ENVIRONMENT` – logical environment name you can include in `OTEL_RESOURCE_ATTRIBUTES` (for example, `deployment.environment=prod`)

## Configuring auto instrumentation

[AWS SDK v3 instrumentation](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/packages/instrumentation-aws-sdk)
is included and loaded automatically by default.
A subset of instrumentations from the [OTEL auto-instrumentations-node metapackage](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/packages/auto-instrumentations-node)
are also included.

Following instrumentations from the metapackage are included:
- `amqplib`
- `bunyan`
- `cassandra-driver`
- `connect`
- `dataloader`
- `dns` *- default*
- `express` *- default*
- `fs`
- `graphql` *- default*
- `grpc` *- default*
- `hapi` *- default*
- `http` *- default*
- `ioredis` *- default*
- `kafkajs`
- `knex`
- `koa` *- default*
- `memcached`
- `mongodb` *- default*
- `mongoose`
- `mysql` *- default*
- `mysql2`
- `nestjs-core`
- `net` *- default*
- `pg` *- default*
- `pino`
- `redis` *- default*
- `restify`
- `socket.io`
- `undici`
- `winston`

Instrumentations annotated with "*- default*" are loaded by default.

To only load specific instrumentations, specify the `OTEL_NODE_ENABLED_INSTRUMENTATIONS` environment variable in the lambda configuration.
This disables all the defaults, and only enables the ones you specify. Selectively disabling instrumentations from the defaults is also possible with the `OTEL_NODE_DISABLED_INSTRUMENTATIONS` environment variable.

The environment variables should be set to a comma-separated list of the instrumentation package names without the
`@opentelemetry/instrumentation-` prefix.

For example, to enable only `@opentelemetry/instrumentation-http` and `@opentelemetry/instrumentation-undici`:
```shell
OTEL_NODE_ENABLED_INSTRUMENTATIONS="http,undici"
```
To disable only `@opentelemetry/instrumentation-net`:
```shell
OTEL_NODE_DISABLED_INSTRUMENTATIONS="net"
```

## Building

To build the layer and sample applications in this `nodejs` folder:

First install dependencies:

```
npm install
```

Then build the project:

```
npm run build
```

You'll find the generated layer zip file at `./packages/layer/build/layer.zip`.

## Sample applications

Sample applications are provided to show usage of the above layer.

- Application using AWS SDK - shows using the wrapper with an application using AWS SDK without code change.
  - [WIP] [Using OTel Public Layer](./sample-apps/aws-sdk)

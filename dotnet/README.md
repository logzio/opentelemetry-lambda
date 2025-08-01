# OpenTelemetry Lambda .NET

NuGet package for running .NET applications on AWS Lambda with OpenTelemetry.

## Provided SDK

[OpenTelemetry Lambda SDK for .NET](https://github.com/open-telemetry/opentelemetry-dotnet-contrib/tree/main/src/OpenTelemetry.Instrumentation.AWSLambda) includes tracing APIs to instrument Lambda handlers and is provided on [NuGet](https://www.nuget.org/packages/OpenTelemetry.Instrumentation.AWSLambda). Follow the instructions in the [user guide](https://aws-otel.github.io/docs/getting-started/lambda/lambda-dotnet#instrumentation) to manually instrument the Lambda handler.
For other instrumentations, such as HTTP, you'll need to include the corresponding library instrumentation from the [instrumentation project](https://github.com/open-telemetry/opentelemetry-dotnet) and modify your code to initialize it in your function.

## Provided Layer

[OpenTelemetry Lambda Layer for Collector](https://aws-otel.github.io/docs/getting-started/lambda/lambda-dotnet#lambda-layer) includes OpenTelemetry Collector for Lambda components. Follow the [user guide](https://aws-otel.github.io/docs/getting-started/lambda/lambda-dotnet#enable-tracing) to apply this layer to your Lambda handler that's already been instrumented with the OpenTelemetry Lambda .NET SDK to enable end-to-end tracing.

## Sample application

The [sample application](https://github.com/open-telemetry/opentelemetry-lambda/blob/main/dotnet/sample-apps/aws-sdk/wrapper/SampleApps/AwsSdkSample/Function.cs) shows the manual instrumentations of OpenTelemetry Lambda .NET SDK on a Lambda handler that triggers a downstream request to AWS S3.

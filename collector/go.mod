module github.com/open-telemetry/opentelemetry-lambda/collector

go 1.24.4

replace github.com/open-telemetry/opentelemetry-lambda/collector/lambdacomponents => ./lambdacomponents

replace github.com/open-telemetry/opentelemetry-lambda/collector/lambdalifecycle => ./lambdalifecycle

replace github.com/open-telemetry/opentelemetry-lambda/collector/processor/coldstartprocessor => ./processor/coldstartprocessor

replace github.com/open-telemetry/opentelemetry-lambda/collector/processor/decoupleprocessor => ./processor/decoupleprocessor

replace github.com/open-telemetry/opentelemetry-lambda/collector/receiver/telemetryapireceiver => ./receiver/telemetryapireceiver

// fixes ambiguous import error: found package cloud.google.com/go/compute/metadata in multiple modules:
//        cloud.google.com/go
//        cloud.google.com/go/compute
// Force cloud.google.com/go to be at least v0.107.0, so that the metadata is not present.
replace cloud.google.com/go => cloud.google.com/go v0.107.0

require (
	github.com/golang-collections/go-datastructures v0.0.0-20150211160725-59788d5eb259
	github.com/google/go-cmp v0.7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/s3provider v0.130.0
	github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/secretsmanagerprovider v0.130.0
	github.com/open-telemetry/opentelemetry-lambda/collector/lambdacomponents v0.98.0
	github.com/open-telemetry/opentelemetry-lambda/collector/lambdalifecycle v0.0.0-20250728180610-2945d3555cc8
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/component v1.36.1
	go.opentelemetry.io/collector/confmap v1.36.1
	go.opentelemetry.io/collector/confmap/provider/envprovider v1.36.1
	go.opentelemetry.io/collector/confmap/provider/fileprovider v1.36.1
	go.opentelemetry.io/collector/confmap/provider/httpprovider v1.36.1
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.36.1
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.36.1
	go.opentelemetry.io/collector/otelcol v0.130.1
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.130.0 // indirect
	github.com/prometheus/otlptranslator v0.0.0-20250725141939-ab8d56d297e3 // indirect
)

require (
	cloud.google.com/go/auth v0.16.3 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.18.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.10.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.2 // indirect
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5 // indirect
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/antchfx/xmlquery v1.4.4 // indirect
	github.com/antchfx/xpath v1.3.4 // indirect
	github.com/aws/aws-sdk-go-v2 v1.37.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.30.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.17.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.85.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.26.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.31.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.35.0 // indirect
	github.com/aws/smithy-go v1.22.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/elastic/go-grok v0.3.1 // indirect
	github.com/elastic/lunes v0.1.0 // indirect
	github.com/expr-lang/expr v1.17.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/foxboron/go-tpm-keyfiles v0.0.0-20250520203025-c3c3a4ec1653 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect; indirect	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/knadh/koanf/v2 v2.2.2 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/sigv4authextension v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/pdatautil v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/sampling v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.130.0 // indirect
	github.com/open-telemetry/opentelemetry-lambda/collector/processor/coldstartprocessor v0.98.0 // indirect
	github.com/open-telemetry/opentelemetry-lambda/collector/processor/decoupleprocessor v0.0.0-20250728180610-2945d3555cc8 // indirect
	github.com/open-telemetry/opentelemetry-lambda/collector/receiver/telemetryapireceiver v0.98.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/prometheus/prometheus v0.305.0 // indirect
	github.com/prometheus/sigv4 v0.2.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/shirou/gopsutil/v4 v4.25.6 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/tg123/go-htpasswd v1.2.4 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/tinylru v1.2.1 // indirect
	github.com/tidwall/wal v1.1.8 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/ua-parser/uap-go v0.0.0-20250326155420-f7f5a2f9f5bc // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector v0.130.1 // indirect
	go.opentelemetry.io/collector/client v1.36.1 // indirect
	go.opentelemetry.io/collector/component/componentstatus v0.130.1 // indirect
	go.opentelemetry.io/collector/component/componenttest v0.130.1 // indirect
	go.opentelemetry.io/collector/config/configauth v0.130.1 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.36.1 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.130.1 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.130.1 // indirect
	go.opentelemetry.io/collector/config/configmiddleware v0.130.1 // indirect
	go.opentelemetry.io/collector/config/confignet v1.36.1 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.36.1 // indirect
	go.opentelemetry.io/collector/config/configoptional v0.130.1 // indirect
	go.opentelemetry.io/collector/config/configretry v1.36.1 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.130.1 // indirect
	go.opentelemetry.io/collector/config/configtls v1.36.1 // indirect
	go.opentelemetry.io/collector/confmap/xconfmap v0.130.1 // indirect
	go.opentelemetry.io/collector/connector v0.130.1 // indirect
	go.opentelemetry.io/collector/connector/connectortest v0.130.1 // indirect
	go.opentelemetry.io/collector/connector/xconnector v0.130.1 // indirect
	go.opentelemetry.io/collector/consumer v1.36.1 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.130.1 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/xconsumererror v0.130.1 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.130.1 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/debugexporter v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/exportertest v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/otlpexporter v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.130.1 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.130.1 // indirect
	go.opentelemetry.io/collector/extension v1.36.1 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v1.36.1 // indirect
	go.opentelemetry.io/collector/extension/extensioncapabilities v0.130.1 // indirect
	go.opentelemetry.io/collector/extension/extensionmiddleware v0.130.1 // indirect
	go.opentelemetry.io/collector/extension/extensiontest v0.130.1 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.130.1 // indirect
	go.opentelemetry.io/collector/featuregate v1.36.1 // indirect
	go.opentelemetry.io/collector/internal/fanoutconsumer v0.130.1 // indirect
	go.opentelemetry.io/collector/internal/memorylimiter v0.130.1 // indirect
	go.opentelemetry.io/collector/internal/sharedcomponent v0.130.1 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.130.1 // indirect
	go.opentelemetry.io/collector/pdata v1.36.1 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.130.1 // indirect
	go.opentelemetry.io/collector/pdata/testdata v0.130.1 // indirect
	go.opentelemetry.io/collector/pdata/xpdata v0.130.1 // indirect
	go.opentelemetry.io/collector/pipeline v0.130.1 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.130.1 // indirect
	go.opentelemetry.io/collector/processor v1.36.1 // indirect
	go.opentelemetry.io/collector/processor/batchprocessor v0.130.1 // indirect
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.130.1 // indirect
	go.opentelemetry.io/collector/processor/processorhelper v0.130.1 // indirect
	go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper v0.130.1 // indirect
	go.opentelemetry.io/collector/processor/processortest v0.130.1 // indirect
	go.opentelemetry.io/collector/processor/xprocessor v0.130.1 // indirect
	go.opentelemetry.io/collector/receiver v1.36.1 // indirect
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.130.1 // indirect
	go.opentelemetry.io/collector/receiver/receiverhelper v0.130.1 // indirect
	go.opentelemetry.io/collector/receiver/receivertest v0.130.1 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.130.1 // indirect
	go.opentelemetry.io/collector/semconv v0.128.0 // indirect
	go.opentelemetry.io/collector/service v0.130.1 // indirect
	go.opentelemetry.io/collector/service/hostcapabilities v0.130.1 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.12.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.62.0 // indirect
	go.opentelemetry.io/contrib/otelconf v0.17.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.37.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.59.1 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.37.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	google.golang.org/api v0.243.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/grpc v1.74.2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.33.3 // indirect
	k8s.io/client-go v0.33.3 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20250604170112-4c0f3b243397 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

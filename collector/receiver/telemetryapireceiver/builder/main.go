package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor/batchprocessor"

	// Import the telemetryapireceiver receiver factory
	"github.com/open-telemetry/opentelemetry-lambda/collector/receiver/telemetryapireceiver"
)

func main() {
	// The extensionID is normally discovered at runtime. For this builder,
	// we can pass a placeholder. The receiver's code doesn't use it
	// directly, the API client does.
	extensionID := os.Getenv("AWS_LAMBDA_EXTENSION_IDENTIFIER")
	if extensionID == "" {
		extensionID = "test-extension-id"
	}

	// Create factory maps for each component type
	receivers, err := otelcol.MakeFactoryMap(
		telemetryapireceiver.NewFactory(extensionID),
	)
	if err != nil {
		log.Fatalf("failed to create receiver factories: %v", err)
	}

	processors, err := otelcol.MakeFactoryMap(
		batchprocessor.NewFactory(),
	)
	if err != nil {
		log.Fatalf("failed to create processor factories: %v", err)
	}

	exporters, err := otelcol.MakeFactoryMap(
		debugexporter.NewFactory(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter factories: %v", err)
	}

	// Create the factories struct
	factories := otelcol.Factories{
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Extensions: make(map[component.Type]extension.Factory),
		Connectors: make(map[component.Type]connector.Factory),
	}

	// Create a new collector with our custom set of components
	info := component.BuildInfo{
		Command:     "otelcol-custom",
		Description: "Custom OTel Collector for testing",
		Version:     "1.0.0",
	}

	// Configure the config provider settings
	configProviderSettings := otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs: []string{"config.yaml"},
			ProviderFactories: []confmap.ProviderFactory{
				fileprovider.NewFactory(),
				envprovider.NewFactory(),
				yamlprovider.NewFactory(),
			},
		},
	}

	app, err := otelcol.NewCollector(otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: func() (otelcol.Factories, error) {
			return factories, nil
		},
		ConfigProviderSettings: configProviderSettings,
	})
	if err != nil {
		log.Fatalf("failed to create collector: %v", err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("collector server run finished with error: %v", err)
	}
}

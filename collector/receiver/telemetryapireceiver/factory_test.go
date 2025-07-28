package telemetryapireceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestFactory(t *testing.T) {
	factory := NewFactory("test-id")
	require.Equal(t, component.MustNewType(typeStr), factory.Type())

	// Test default config creation
	expectedCfg := &Config{
		extensionID: "test-id",
		Port:        defaultPort,
		Types:       []string{platform, function, extension},
		MaxItems:    defaultMaxItems,
		MaxBytes:    defaultMaxBytes,
		TimeoutMS:   defaultTimeoutMS,
	}
	require.Equal(t, expectedCfg, factory.CreateDefaultConfig())

	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))
	nopConsumer := consumertest.NewNop()

	// Test logs receiver creation
	_, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		factory.CreateDefaultConfig(),
		nopConsumer,
	)
	require.NoError(t, err)

	// Test traces receiver creation
	_, err = factory.CreateTraces(
		context.Background(),
		nopSettings,
		factory.CreateDefaultConfig(),
		nopConsumer,
	)
	require.NoError(t, err)

	// Test metrics receiver creation
	_, err = factory.CreateMetrics(
		context.Background(),
		nopSettings,
		factory.CreateDefaultConfig(),
		nopConsumer,
	)
	require.NoError(t, err)
}

// TestFactory_Constants tests factory constants
func TestFactory_Constants(t *testing.T) {
	assert.Equal(t, "telemetryapireceiver", typeStr)
	assert.Equal(t, component.StabilityLevelBeta, stability)
	assert.Equal(t, "platform", platform)
	assert.Equal(t, "function", function)
	assert.Equal(t, "extension", extension)
	assert.Equal(t, 4325, defaultPort)
	assert.Equal(t, 1000, defaultMaxItems)
	assert.Equal(t, 262144, defaultMaxBytes)
	assert.Equal(t, 1000, defaultTimeoutMS)
}

// TestFactory_InvalidConfig tests factory behavior with invalid configurations
func TestFactory_InvalidConfig(t *testing.T) {
	factory := NewFactory("test-id")
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))
	nopConsumer := consumertest.NewNop()

	tests := []struct {
		name   string
		config component.Config
	}{
		{
			name:   "non-telemetryapi config",
			config: &invalidConfig{},
		},
		{
			name:   "nil config",
			config: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logs receiver creation with invalid config
			_, err := factory.CreateLogs(
				context.Background(),
				nopSettings,
				tt.config,
				nopConsumer,
			)
			assert.Error(t, err)

			// Test traces receiver creation with invalid config
			_, err = factory.CreateTraces(
				context.Background(),
				nopSettings,
				tt.config,
				nopConsumer,
			)
			assert.Error(t, err)

			// Test metrics receiver creation with invalid config
			_, err = factory.CreateMetrics(
				context.Background(),
				nopSettings,
				tt.config,
				nopConsumer,
			)
			assert.Error(t, err)
		})
	}
}

// TestFactory_SharedComponent tests that the same receiver is shared across signal types
func TestFactory_SharedComponent(t *testing.T) {
	factory := NewFactory("test-id")
	config := factory.CreateDefaultConfig()
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	// Create receivers for different signal types with the same config
	logsReceiver, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		config,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	tracesReceiver, err := factory.CreateTraces(
		context.Background(),
		nopSettings,
		config,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	metricsReceiver, err := factory.CreateMetrics(
		context.Background(),
		nopSettings,
		config,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	// All receivers should be the same shared instance
	assert.Equal(t, logsReceiver, tracesReceiver)
	assert.Equal(t, tracesReceiver, metricsReceiver)
}

// TestFactory_DifferentConfigs tests that different configs create different receivers
func TestFactory_DifferentConfigs(t *testing.T) {
	factory := NewFactory("test-id")
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	config1 := &Config{
		extensionID: "test-id",
		Port:        4325,
		Types:       []string{platform},
		MaxItems:    1000,
		MaxBytes:    262144,
		TimeoutMS:   1000,
	}

	config2 := &Config{
		extensionID: "test-id",
		Port:        4326, // Different port
		Types:       []string{function},
		MaxItems:    2000,
		MaxBytes:    524288,
		TimeoutMS:   2000,
	}

	// Create receivers with different configs
	receiver1, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		config1,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	receiver2, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		config2,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	// Should be different instances
	assert.NotEqual(t, receiver1, receiver2)
}

// TestFactory_WithNilConsumer tests factory behavior with nil consumers
func TestFactory_WithNilConsumer(t *testing.T) {
	factory := NewFactory("test-id")
	config := factory.CreateDefaultConfig()
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	// Test with nil consumers - should not error but handle gracefully
	_, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		config,
		nil,
	)
	// This should succeed as consumer registration happens after creation
	assert.NoError(t, err)

	_, err = factory.CreateTraces(
		context.Background(),
		nopSettings,
		config,
		nil,
	)
	assert.NoError(t, err)

	_, err = factory.CreateMetrics(
		context.Background(),
		nopSettings,
		config,
		nil,
	)
	assert.NoError(t, err)
}

// TestFactory_ReceiverCreationError tests factory behavior when receiver creation fails
func TestFactory_ReceiverCreationError(t *testing.T) {
	factory := NewFactory("test-id")
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	// Create a config that will cause receiver creation to fail
	// This simulates internal errors during receiver creation
	config := &Config{
		extensionID: "",
		Port:        -1, // Invalid port
		Types:       []string{},
		MaxItems:    0,
		MaxBytes:    0,
		TimeoutMS:   0,
	}

	// The current implementation doesn't validate these fields during creation,
	// it only validates during Start(). So creation should succeed here.
	_, err := factory.CreateLogs(
		context.Background(),
		nopSettings,
		config,
		consumertest.NewNop(),
	)
	// Creation should succeed - validation happens later in Start()
	assert.NoError(t, err)
}

// TestGetOrCreateReceiver tests the shared component creation logic
func TestGetOrCreateReceiver(t *testing.T) {
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	// Test with valid config
	config := &Config{
		extensionID: "test-id",
		Port:        4325,
		Types:       []string{platform},
		MaxItems:    1000,
		MaxBytes:    262144,
		TimeoutMS:   1000,
	}

	shared1, err := getOrCreateReceiver(nopSettings, config)
	require.NoError(t, err)
	require.NotNil(t, shared1)
	require.NotNil(t, shared1.Unwrap())

	// Second call with same config should return the same instance
	shared2, err := getOrCreateReceiver(nopSettings, config)
	require.NoError(t, err)
	assert.Equal(t, shared1, shared2)

	// Test with invalid config type
	_, err = getOrCreateReceiver(nopSettings, &invalidConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config was not a Telemetry API receiver config")
}

// TestFactory_ExtensionIDPropagation tests that extension ID is properly set
func TestFactory_ExtensionIDPropagation(t *testing.T) {
	extensionID := "custom-extension-id"
	factory := NewFactory(extensionID)

	config := factory.CreateDefaultConfig().(*Config)
	assert.Equal(t, extensionID, config.extensionID)
}

// TestFactory_ConfigTypes tests different configuration combinations
func TestFactory_ConfigTypes(t *testing.T) {
	factory := NewFactory("test-id")
	nopSettings := receivertest.NewNopSettings(component.MustNewType(typeStr))

	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "platform only",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{platform},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
		},
		{
			name: "function only",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{function},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
		},
		{
			name: "extension only",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{extension},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
		},
		{
			name: "multiple types",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{platform, function, extension},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
		},
		{
			name: "empty types",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all signal type receivers can be created with different config types
			_, err := factory.CreateLogs(context.Background(), nopSettings, tt.config, consumertest.NewNop())
			assert.NoError(t, err)

			_, err = factory.CreateTraces(context.Background(), nopSettings, tt.config, consumertest.NewNop())
			assert.NoError(t, err)

			_, err = factory.CreateMetrics(context.Background(), nopSettings, tt.config, consumertest.NewNop())
			assert.NoError(t, err)
		})
	}
}

// invalidConfig is a helper struct for testing invalid configuration
type invalidConfig struct{}

func (c *invalidConfig) Validate() error {
	return nil
}

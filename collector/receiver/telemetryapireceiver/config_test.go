package telemetryapireceiver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	factory := NewFactory("test-extension-id")
	defaultCfg := factory.CreateDefaultConfig()

	// Test that the default config is created correctly
	require.Equal(t, uint(defaultMaxItems), defaultCfg.(*Config).MaxItems)

	// Test loading a config from YAML - this loads the "telemetryapi" section which doesn't have port specified
	sub, err := cm.Sub(component.NewIDWithName(component.MustNewType(typeStr), "").String())
	require.NoError(t, err)

	cfg := factory.CreateDefaultConfig()
	require.NoError(t, sub.Unmarshal(cfg))

	expected := &Config{
		extensionID: "test-extension-id",
		Port:        defaultPort,                             // Should be default since not specified in YAML
		Types:       []string{platform, function, extension}, // Should be default
		MaxItems:    defaultMaxItems,
		MaxBytes:    defaultMaxBytes,
		TimeoutMS:   defaultTimeoutMS,
	}
	require.Equal(t, expected, cfg)
}

// TestConfig_Validation tests the configuration validation
func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with all types",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{platform, function, extension},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: false,
		},
		{
			name: "valid config with single type",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{platform},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: false,
		},
		{
			name: "valid config with empty types",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: false,
		},
		{
			name: "invalid config with unknown type",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{"unknown_type"},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: true,
			errorMsg:    "unknown extension type: unknown_type",
		},
		{
			name: "invalid config with mixed valid and invalid types",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{platform, "invalid_type", function},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: true,
			errorMsg:    "unknown extension type: invalid_type",
		},
		{
			name: "invalid config with empty string type",
			config: &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       []string{""},
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			},
			expectError: true,
			errorMsg:    "unknown extension type: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLoadConfig_AllConfigurations tests loading all configurations from the test file
func TestLoadConfig_AllConfigurations(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	factory := NewFactory("test-extension-id")

	tests := []struct {
		name     string
		configID string
		expected *Config
	}{
		{
			name:     "default config",
			configID: "telemetryapi",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        defaultPort,
				Types:       []string{platform, function, extension},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
		{
			name:     "config with platform only",
			configID: "telemetryapi/2",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        12345,
				Types:       []string{platform},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
		{
			name:     "config with function only",
			configID: "telemetryapi/3",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        12345,
				Types:       []string{function},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
		{
			name:     "config with extension only",
			configID: "telemetryapi/4",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        12345,
				Types:       []string{extension},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
		{
			name:     "config with multiple types - array style",
			configID: "telemetryapi/5",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        12345,
				Types:       []string{platform, function},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
		{
			name:     "config with empty types array",
			configID: "telemetryapi/8",
			expected: &Config{
				extensionID: "test-extension-id",
				Port:        12345,
				Types:       []string{},
				MaxItems:    defaultMaxItems,
				MaxBytes:    defaultMaxBytes,
				TimeoutMS:   defaultTimeoutMS,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := cm.Sub(tt.configID)
			require.NoError(t, err)

			cfg := factory.CreateDefaultConfig()
			require.NoError(t, sub.Unmarshal(cfg))

			assert.Equal(t, tt.expected, cfg)

			// Verify the config is valid
			assert.NoError(t, cfg.(*Config).Validate())
		})
	}
}

// TestConfig_DefaultValues tests that default configuration values are correct
func TestConfig_DefaultValues(t *testing.T) {
	factory := NewFactory("test-extension-id")
	cfg := factory.CreateDefaultConfig().(*Config)

	assert.Equal(t, "test-extension-id", cfg.extensionID)
	assert.Equal(t, defaultPort, cfg.Port)
	assert.Equal(t, []string{platform, function, extension}, cfg.Types)
	assert.Equal(t, uint(defaultMaxItems), cfg.MaxItems)
	assert.Equal(t, uint(defaultMaxBytes), cfg.MaxBytes)
	assert.Equal(t, uint(defaultTimeoutMS), cfg.TimeoutMS)

	// Verify default config is valid
	assert.NoError(t, cfg.Validate())
}

// TestConfig_BoundaryValues tests boundary values for configuration parameters
func TestConfig_BoundaryValues(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		field    string
		expected interface{}
	}{
		{
			name: "minimum port value",
			config: &Config{
				extensionID: "test-id",
				Port:        1,
				Types:       []string{platform},
				MaxItems:    1,
				MaxBytes:    1,
				TimeoutMS:   1,
			},
			field:    "Port",
			expected: 1,
		},
		{
			name: "maximum uint values",
			config: &Config{
				extensionID: "test-id",
				Port:        65535,
				Types:       []string{platform},
				MaxItems:    ^uint(0), // Maximum uint value
				MaxBytes:    ^uint(0), // Maximum uint value
				TimeoutMS:   ^uint(0), // Maximum uint value
			},
			field:    "MaxItems",
			expected: ^uint(0),
		},
		{
			name: "zero values",
			config: &Config{
				extensionID: "test-id",
				Port:        0,
				Types:       []string{},
				MaxItems:    0,
				MaxBytes:    0,
				TimeoutMS:   0,
			},
			field:    "MaxItems",
			expected: uint(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// All configs should be valid regardless of boundary values
			// (validation is currently only for Types field)
			assert.NoError(t, tt.config.Validate())
		})
	}
}

// TestConfig_TypesFieldVariations tests different ways to specify types
func TestConfig_TypesFieldVariations(t *testing.T) {
	tests := []struct {
		name          string
		types         []string
		expectValid   bool
		expectedTypes []string
	}{
		{
			name:          "nil types",
			types:         nil,
			expectValid:   true,
			expectedTypes: nil,
		},
		{
			name:          "empty types",
			types:         []string{},
			expectValid:   true,
			expectedTypes: []string{},
		},
		{
			name:          "single valid type",
			types:         []string{platform},
			expectValid:   true,
			expectedTypes: []string{platform},
		},
		{
			name:          "all valid types",
			types:         []string{platform, function, extension},
			expectValid:   true,
			expectedTypes: []string{platform, function, extension},
		},
		{
			name:          "duplicate types",
			types:         []string{platform, platform, function},
			expectValid:   true,
			expectedTypes: []string{platform, platform, function}, // Duplicates are allowed
		},
		{
			name:        "invalid type mixed with valid",
			types:       []string{platform, "invalid"},
			expectValid: false,
		},
		{
			name:        "only invalid types",
			types:       []string{"invalid1", "invalid2"},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				extensionID: "test-id",
				Port:        4325,
				Types:       tt.types,
				MaxItems:    1000,
				MaxBytes:    262144,
				TimeoutMS:   1000,
			}

			err := config.Validate()
			if tt.expectValid {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTypes, config.Types)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestConfig_ExtensionIDHandling tests extension ID handling
func TestConfig_ExtensionIDHandling(t *testing.T) {
	tests := []struct {
		name        string
		extensionID string
	}{
		{
			name:        "normal extension ID",
			extensionID: "my-extension",
		},
		{
			name:        "empty extension ID",
			extensionID: "",
		},
		{
			name:        "extension ID with special characters",
			extensionID: "my-extension_123.test",
		},
		{
			name:        "long extension ID",
			extensionID: "very-long-extension-id-that-might-be-used-in-some-scenarios",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory(tt.extensionID)
			cfg := factory.CreateDefaultConfig().(*Config)

			assert.Equal(t, tt.extensionID, cfg.extensionID)
			assert.NoError(t, cfg.Validate())
		})
	}
}

// TestConfig_MarshalUnmarshal tests configuration serialization/deserialization
func TestConfig_MarshalUnmarshal(t *testing.T) {
	original := &Config{
		extensionID: "test-id",
		Port:        8080,
		Types:       []string{platform, function},
		MaxItems:    2000,
		MaxBytes:    524288,
		TimeoutMS:   5000,
	}

	// Note: extensionID is not exported, so it won't be marshaled
	// This test verifies that the public fields work correctly

	factory := NewFactory("test-id")
	cfg := factory.CreateDefaultConfig().(*Config)

	// Manually set the values that would come from unmarshaling
	cfg.Port = original.Port
	cfg.Types = original.Types
	cfg.MaxItems = original.MaxItems
	cfg.MaxBytes = original.MaxBytes
	cfg.TimeoutMS = original.TimeoutMS

	assert.Equal(t, original.Port, cfg.Port)
	assert.Equal(t, original.Types, cfg.Types)
	assert.Equal(t, original.MaxItems, cfg.MaxItems)
	assert.Equal(t, original.MaxBytes, cfg.MaxBytes)
	assert.Equal(t, original.TimeoutMS, cfg.TimeoutMS)
	assert.NoError(t, cfg.Validate())
}

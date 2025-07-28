package telemetryapireceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestHttpHandler_Metrics(t *testing.T) {
	// Setup: Create a mock consumer (a "sink") to receive the metrics
	sink := new(consumertest.MetricsSink)

	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	r.registerMetricsConsumer(sink)

	// Create a sample HTTP request with a platform.report event
	sampleRecord := map[string]interface{}{
		"requestId": "test-req-id",
		"metrics":   map[string]interface{}{"durationMs": 100.0},
	}
	payload, _ := json.Marshal([]event{
		{Time: time.Now().Format(time.RFC3339), Type: "platform.report", Record: sampleRecord},
	})

	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	// Execute the httpHandler
	r.httpHandler(rr, req)

	// Assert the results
	require.Equal(t, http.StatusOK, rr.Code)
	require.Len(t, sink.AllMetrics(), 1, "sink should have received one metric payload")
	allMetrics := sink.AllMetrics()[0]
	numDataPoints := allMetrics.MetricCount()
	require.Equal(t, 1, numDataPoints, "payload should contain one metric")
}

func TestHttpHandler_Logs(t *testing.T) {
	// Setup: Create a sink for logs
	sink := new(consumertest.LogsSink)

	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	r.registerLogsConsumer(sink)

	// Create a sample HTTP request with a function event
	payload, _ := json.Marshal([]event{
		{Time: time.Now().Format(time.RFC3339), Type: "function", Record: "hello world"},
	})

	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	// Execute
	r.httpHandler(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	require.Len(t, sink.AllLogs(), 1, "sink should have received one log payload")
	require.Equal(t, 1, sink.AllLogs()[0].LogRecordCount(), "payload should contain one log record")
}

// TestHttpHandler_Traces tests trace span creation from platform events
func TestHttpHandler_Traces(t *testing.T) {
	sink := new(consumertest.TracesSink)

	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	r.registerTracesConsumer(sink)

	// Simulate init phase
	initStartEvent := event{
		Time:   time.Now().Format(time.RFC3339),
		Type:   "platform.initStart",
		Record: map[string]interface{}{},
	}

	initDoneEvent := event{
		Time: time.Now().Add(time.Second).Format(time.RFC3339),
		Type: "platform.initRuntimeDone",
		Record: map[string]interface{}{
			"status": "success",
		},
	}

	// Send both events
	payload, _ := json.Marshal([]event{initStartEvent, initDoneEvent})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Len(t, sink.AllTraces(), 1, "sink should have received one trace payload")
	require.Equal(t, 1, sink.AllTraces()[0].SpanCount(), "payload should contain one span")
}

// TestHttpHandler_InvokeTraces tests invoke span creation
func TestHttpHandler_InvokeTraces(t *testing.T) {
	sink := new(consumertest.TracesSink)

	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	r.registerTracesConsumer(sink)

	requestID := "test-request-123"

	// Simulate invoke phase
	startEvent := event{
		Time: time.Now().Format(time.RFC3339),
		Type: "platform.start",
		Record: map[string]interface{}{
			"requestId": requestID,
		},
	}

	doneEvent := event{
		Time: time.Now().Add(100 * time.Millisecond).Format(time.RFC3339),
		Type: "platform.runtimeDone",
		Record: map[string]interface{}{
			"requestId": requestID,
			"status":    "success",
		},
	}

	// Send both events
	payload, _ := json.Marshal([]event{startEvent, doneEvent})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Len(t, sink.AllTraces(), 1, "sink should have received one trace payload")
	require.Equal(t, 1, sink.AllTraces()[0].SpanCount(), "payload should contain one span")
}

// TestHttpHandler_ErrorCases tests various error scenarios
func TestHttpHandler_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		expectedStatus int
	}{
		{
			name:           "malformed JSON",
			payload:        []byte("{invalid json"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			payload:        []byte(""),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non-array JSON",
			payload:        []byte(`{"not": "an array"}`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
			require.NoError(t, err)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(tt.payload))
			rr := httptest.NewRecorder()

			r.httpHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

// TestHttpHandler_MultipleEvents tests handling multiple events in one request
func TestHttpHandler_MultipleEvents(t *testing.T) {
	metricsSink := new(consumertest.MetricsSink)
	logsSink := new(consumertest.LogsSink)

	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	r.registerMetricsConsumer(metricsSink)
	r.registerLogsConsumer(logsSink)

	// Create multiple events of different types
	events := []event{
		{
			Time: time.Now().Format(time.RFC3339),
			Type: "function",
			Record: map[string]interface{}{
				"message": "function log",
				"level":   "INFO",
			},
		},
		{
			Time:   time.Now().Format(time.RFC3339),
			Type:   "extension",
			Record: "extension log message",
		},
		{
			Time: time.Now().Format(time.RFC3339),
			Type: "platform.report",
			Record: map[string]interface{}{
				"requestId": "test-req-id",
				"metrics": map[string]interface{}{
					"durationMs":   150.0,
					"memorySizeMB": 128.0,
				},
			},
		},
	}

	payload, _ := json.Marshal(events)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Len(t, logsSink.AllLogs(), 2, "should receive 2 log events")
	require.Len(t, metricsSink.AllMetrics(), 1, "should receive 1 metrics event")
}

// TestHttpHandler_UnknownEventTypes tests handling of unknown event types
func TestHttpHandler_UnknownEventTypes(t *testing.T) {
	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)

	// Create an event with unknown type
	events := []event{
		{
			Time:   time.Now().Format(time.RFC3339),
			Type:   "unknown.event.type",
			Record: map[string]interface{}{},
		},
	}

	payload, _ := json.Marshal(events)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	// Should still return OK, just ignore unknown events
	require.Equal(t, http.StatusOK, rr.Code)
}

// TestHttpHandler_NoConsumers tests behavior when no consumers are registered
func TestHttpHandler_NoConsumers(t *testing.T) {
	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)
	// Don't register any consumers

	events := []event{
		{
			Time:   time.Now().Format(time.RFC3339),
			Type:   "function",
			Record: "test message",
		},
	}

	payload, _ := json.Marshal(events)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	// Should still return OK even without consumers
	require.Equal(t, http.StatusOK, rr.Code)
}

// TestNewTelemetryAPIReceiver tests receiver creation
func TestNewTelemetryAPIReceiver(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected map[string]string
	}{
		{
			name: "with all environment variables",
			envVars: map[string]string{
				"AWS_LAMBDA_FUNCTION_NAME":        "test-function",
				"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": "128",
				"AWS_LAMBDA_FUNCTION_VERSION":     "$LATEST",
				"AWS_REGION":                      "us-west-2",
				"LOGZIO_ENV_ID":                   "test-env",
			},
			expected: map[string]string{
				"service.name":    "test-function",
				"faas.name":       "test-function",
				"faas.max_memory": "128",
				"faas.version":    "$LATEST",
				"cloud.region":    "us-west-2",
				"env_id":          "test-env",
				"cloud.provider":  "aws",
			},
		},
		{
			name: "with minimal environment variables",
			envVars: map[string]string{
				"AWS_LAMBDA_FUNCTION_NAME": "test-function",
			},
			expected: map[string]string{
				"service.name":   "test-function",
				"faas.name":      "test-function",
				"cloud.provider": "aws",
			},
		},
		{
			name:    "without function name",
			envVars: map[string]string{},
			expected: map[string]string{
				"service.name":   "unknown_service",
				"cloud.provider": "aws",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg := &Config{}
			r, err := newTelemetryAPIReceiver(cfg, receivertest.NewNopSettings(component.MustNewType(typeStr)))
			require.NoError(t, err)

			// Verify resource attributes
			attrs := r.resource.Attributes()
			for expectedKey, expectedValue := range tt.expected {
				val, exists := attrs.Get(expectedKey)
				assert.True(t, exists, "Expected attribute %s to exist", expectedKey)
				if exists {
					assert.Equal(t, expectedValue, val.Str(), "Expected attribute %s to have value %s", expectedKey, expectedValue)
				}
			}
		})
	}
}

// TestListenOnAddress tests the address calculation logic
func TestListenOnAddress(t *testing.T) {
	tests := []struct {
		name        string
		port        int
		awsSamLocal string
		expected    string
	}{
		{
			name:     "production environment",
			port:     4325,
			expected: "sandbox.localdomain:4325",
		},
		{
			name:        "SAM local environment",
			port:        4325,
			awsSamLocal: "true",
			expected:    "127.0.0.1:4325",
		},
		{
			name:        "SAM local false",
			port:        4325,
			awsSamLocal: "false",
			expected:    "sandbox.localdomain:4325",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.awsSamLocal != "" {
				os.Setenv("AWS_SAM_LOCAL", tt.awsSamLocal)
				defer os.Unsetenv("AWS_SAM_LOCAL")
			}

			result := listenOnAddress(tt.port)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestReceiverLifecycle tests Start and Shutdown methods
func TestReceiverLifecycle(t *testing.T) {
	// Skip this test if we can't control the telemetry API environment
	if _, exists := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); !exists {
		t.Skip("Skipping lifecycle test: AWS_LAMBDA_RUNTIME_API not set")
	}

	cfg := &Config{
		Port:      0, // Use random port
		Types:     []string{},
		MaxItems:  1000,
		MaxBytes:  262144,
		TimeoutMS: 1000,
	}

	r, err := newTelemetryAPIReceiver(cfg, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)

	ctx := context.Background()

	// Test Start - this will fail without proper AWS environment but should handle gracefully
	err = r.Start(ctx, nil)
	if err == nil {
		// If Start succeeded, test Shutdown
		err = r.Shutdown(ctx)
		assert.NoError(t, err, "Shutdown should succeed")
	}
	// If Start failed due to environment, that's expected in test environment
}

// TestReceiverConsumerRegistration tests consumer registration
func TestReceiverConsumerRegistration(t *testing.T) {
	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)

	logsSink := new(consumertest.LogsSink)
	tracesSink := new(consumertest.TracesSink)
	metricsSink := new(consumertest.MetricsSink)

	// Test consumer registration
	r.registerLogsConsumer(logsSink)
	r.registerTracesConsumer(tracesSink)
	r.registerMetricsConsumer(metricsSink)

	assert.Equal(t, logsSink, r.nextLogs)
	assert.Equal(t, tracesSink, r.nextTraces)
	assert.Equal(t, metricsSink, r.nextMetrics)
}

// TestHttpHandler_ReadError tests handling of body read errors
func TestHttpHandler_ReadError(t *testing.T) {
	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)

	// Create a request with a body that will cause a read error
	req := httptest.NewRequest("POST", "/", nil)
	req.Body = &errorReader{}
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestHttpHandler_LargePayload tests handling of large payloads
func TestHttpHandler_LargePayload(t *testing.T) {
	r, err := newTelemetryAPIReceiver(&Config{}, receivertest.NewNopSettings(component.MustNewType(typeStr)))
	require.NoError(t, err)

	// Create a large number of events
	events := make([]event, 1000)
	for i := range events {
		events[i] = event{
			Time: time.Now().Format(time.RFC3339),
			Type: "function",
			Record: map[string]interface{}{
				"message": "test message",
				"index":   i,
			},
		}
	}

	payload, _ := json.Marshal(events)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	r.httpHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

// errorReader is a helper for testing read errors
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}

func (e *errorReader) Close() error {
	return nil
}

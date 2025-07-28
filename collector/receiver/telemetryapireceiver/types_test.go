package telemetryapireceiver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEvent_GetTime tests the getTime method of the event struct
func TestEvent_GetTime(t *testing.T) {
	tests := []struct {
		name        string
		timeStr     string
		expectError bool
		expectedLen int // Length to check if it's close to current time
	}{
		{
			name:        "valid RFC3339 timestamp",
			timeStr:     "2022-10-12T00:03:50.000Z",
			expectError: false,
		},
		{
			name:        "valid RFC3339 timestamp with timezone",
			timeStr:     "2022-10-12T00:03:50.000-07:00",
			expectError: false,
		},
		{
			name:        "valid RFC3339 timestamp with microseconds",
			timeStr:     "2022-10-12T00:03:50.123456Z",
			expectError: false,
		},
		{
			name:        "invalid timestamp format",
			timeStr:     "invalid-timestamp",
			expectError: true,
		},
		{
			name:        "empty timestamp",
			timeStr:     "",
			expectError: true,
		},
		{
			name:        "partial timestamp",
			timeStr:     "2022-10-12",
			expectError: true,
		},
		{
			name:        "wrong format",
			timeStr:     "2022/10/12 00:03:50",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := event{
				Time:   tt.timeStr,
				Type:   "test",
				Record: "test record",
			}

			result := event.getTime()

			if tt.expectError {
				// For invalid timestamps, should return current time as fallback
				// Check that it's within the last few seconds (reasonable for test execution time)
				now := time.Now()
				diff := now.Sub(result)
				assert.True(t, diff >= 0 && diff < 5*time.Second,
					"Expected fallback time to be close to current time, got diff: %v", diff)
			} else {
				// For valid timestamps, should parse correctly
				expected, err := time.Parse(time.RFC3339, tt.timeStr)
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			}
		})
	}
}

// TestEvent_Structure tests the event structure
func TestEvent_Structure(t *testing.T) {
	tests := []struct {
		name   string
		event  event
		checks func(t *testing.T, e event)
	}{
		{
			name: "string record",
			event: event{
				Time:   "2022-10-12T00:03:50.000Z",
				Type:   "function",
				Record: "plain string record",
			},
			checks: func(t *testing.T, e event) {
				assert.Equal(t, "2022-10-12T00:03:50.000Z", e.Time)
				assert.Equal(t, "function", e.Type)
				assert.Equal(t, "plain string record", e.Record)
			},
		},
		{
			name: "map record",
			event: event{
				Time: "2022-10-12T00:03:50.000Z",
				Type: "platform.report",
				Record: map[string]interface{}{
					"requestId": "test-123",
					"metrics": map[string]interface{}{
						"durationMs": 150.5,
					},
				},
			},
			checks: func(t *testing.T, e event) {
				assert.Equal(t, "2022-10-12T00:03:50.000Z", e.Time)
				assert.Equal(t, "platform.report", e.Type)

				record, ok := e.Record.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "test-123", record["requestId"])

				metrics, ok := record["metrics"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, 150.5, metrics["durationMs"])
			},
		},
		{
			name: "complex nested record",
			event: event{
				Time: "2022-10-12T00:03:50.000Z",
				Type: "function",
				Record: map[string]interface{}{
					"timestamp": "2022-10-12T00:03:50.000Z",
					"level":     "INFO",
					"requestId": "test-req-id",
					"message":   "Complex log message",
					"trace_id":  "80e1afed08e019fc1110464cfa66635c",
					"span_id":   "7a085853722dc6d2",
					"attributes": map[string]interface{}{
						"key1": "value1",
						"key2": 42,
						"key3": true,
					},
				},
			},
			checks: func(t *testing.T, e event) {
				record, ok := e.Record.(map[string]interface{})
				require.True(t, ok)

				assert.Equal(t, "INFO", record["level"])
				assert.Equal(t, "test-req-id", record["requestId"])
				assert.Equal(t, "Complex log message", record["message"])
				assert.Equal(t, "80e1afed08e019fc1110464cfa66635c", record["trace_id"])
				assert.Equal(t, "7a085853722dc6d2", record["span_id"])

				attrs, ok := record["attributes"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "value1", attrs["key1"])
				assert.Equal(t, 42, attrs["key2"])
				assert.Equal(t, true, attrs["key3"])
			},
		},
		{
			name: "nil record",
			event: event{
				Time:   "2022-10-12T00:03:50.000Z",
				Type:   "test",
				Record: nil,
			},
			checks: func(t *testing.T, e event) {
				assert.Nil(t, e.Record)
			},
		},
		{
			name: "numeric record",
			event: event{
				Time:   "2022-10-12T00:03:50.000Z",
				Type:   "test",
				Record: 42.5,
			},
			checks: func(t *testing.T, e event) {
				assert.Equal(t, 42.5, e.Record)
			},
		},
		{
			name: "boolean record",
			event: event{
				Time:   "2022-10-12T00:03:50.000Z",
				Type:   "test",
				Record: true,
			},
			checks: func(t *testing.T, e event) {
				assert.Equal(t, true, e.Record)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.checks(t, tt.event)
		})
	}
}

// TestEvent_RecordTypes tests different record types that might be encountered
func TestEvent_RecordTypes(t *testing.T) {
	baseEvent := event{
		Time: "2022-10-12T00:03:50.000Z",
		Type: "test",
	}

	// Test various record types
	records := []interface{}{
		"string record",
		123,
		123.456,
		true,
		false,
		nil,
		[]string{"array", "of", "strings"},
		[]interface{}{"mixed", 123, true},
		map[string]interface{}{"key": "value"},
		map[string]interface{}{
			"nested": map[string]interface{}{
				"deep": "value",
			},
		},
	}

	for _, record := range records {
		t.Run(fmt.Sprintf("record_type_%T", record), func(t *testing.T) {
			event := baseEvent
			event.Record = record

			// Verify the record is stored correctly
			assert.Equal(t, record, event.Record)

			// Verify getTime works regardless of record type
			parsedTime := event.getTime()
			expectedTime, _ := time.Parse(time.RFC3339, event.Time)
			assert.Equal(t, expectedTime, parsedTime)
		})
	}
}

// TestEvent_RealWorldExamples tests with real-world example data
func TestEvent_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name  string
		event event
	}{
		{
			name: "AWS Lambda platform.initStart",
			event: event{
				Time: "2023-10-12T12:34:56.789Z",
				Type: "platform.initStart",
				Record: map[string]interface{}{
					"initializationType": "on-demand",
				},
			},
		},
		{
			name: "AWS Lambda platform.initRuntimeDone",
			event: event{
				Time: "2023-10-12T12:34:57.123Z",
				Type: "platform.initRuntimeDone",
				Record: map[string]interface{}{
					"initializationType": "on-demand",
					"status":             "success",
				},
			},
		},
		{
			name: "AWS Lambda platform.start",
			event: event{
				Time: "2023-10-12T12:35:00.000Z",
				Type: "platform.start",
				Record: map[string]interface{}{
					"requestId": "c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
					"version":   "$LATEST",
				},
			},
		},
		{
			name: "AWS Lambda platform.runtimeDone",
			event: event{
				Time: "2023-10-12T12:35:00.456Z",
				Type: "platform.runtimeDone",
				Record: map[string]interface{}{
					"requestId": "c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
					"status":    "success",
				},
			},
		},
		{
			name: "AWS Lambda platform.report",
			event: event{
				Time: "2023-10-12T12:35:00.500Z",
				Type: "platform.report",
				Record: map[string]interface{}{
					"requestId": "c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
					"metrics": map[string]interface{}{
						"durationMs":       456.78,
						"billedDurationMs": 500.0,
						"memorySizeMB":     128.0,
						"maxMemoryUsedMB":  64.0,
						"initDurationMs":   234.56,
					},
				},
			},
		},
		{
			name: "Function log with structured data",
			event: event{
				Time: "2023-10-12T12:35:00.200Z",
				Type: "function",
				Record: map[string]interface{}{
					"timestamp": "2023-10-12T12:35:00.200Z",
					"level":     "INFO",
					"requestId": "c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
					"message":   "Processing request",
					"trace_id":  "1-5e1b4151-c4b5ff3f-1b2a3c4d5e6f7890",
					"span_id":   "1234567890abcdef",
				},
			},
		},
		{
			name: "Function log plain text",
			event: event{
				Time:   "2023-10-12T12:35:00.300Z",
				Type:   "function",
				Record: "2023-10-12T12:35:00.300Z\tINFO\tSimple log message\n",
			},
		},
		{
			name: "Extension log",
			event: event{
				Time: "2023-10-12T12:35:00.100Z",
				Type: "extension",
				Record: map[string]interface{}{
					"timestamp": "2023-10-12T12:35:00.100Z",
					"level":     "DEBUG",
					"message":   "Extension initialized",
					"extension": "my-extension",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the event structure is valid
			assert.NotEmpty(t, tt.event.Time)
			assert.NotEmpty(t, tt.event.Type)
			assert.NotNil(t, tt.event.Record)

			// Verify time parsing works
			parsedTime := tt.event.getTime()
			assert.False(t, parsedTime.IsZero())

			// Verify it's a valid timestamp
			expectedTime, err := time.Parse(time.RFC3339, tt.event.Time)
			assert.NoError(t, err)
			assert.Equal(t, expectedTime, parsedTime)
		})
	}
}

// TestEvent_EdgeCases tests edge cases for event handling
func TestEvent_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		event event
	}{
		{
			name: "very long timestamp",
			event: event{
				Time:   "2023-10-12T12:35:00.123456789012345Z", // More precision than typically supported
				Type:   "test",
				Record: "test",
			},
		},
		{
			name: "minimum valid timestamp",
			event: event{
				Time:   "0001-01-01T00:00:00Z",
				Type:   "test",
				Record: "test",
			},
		},
		{
			name: "empty type",
			event: event{
				Time:   "2023-10-12T12:35:00.000Z",
				Type:   "",
				Record: "test",
			},
		},
		{
			name: "very long type",
			event: event{
				Time:   "2023-10-12T12:35:00.000Z",
				Type:   "very.long.event.type.name.that.might.be.used.in.some.scenarios.with.deep.nesting",
				Record: "test",
			},
		},
		{
			name: "unicode in fields",
			event: event{
				Time: "2023-10-12T12:35:00.000Z",
				Type: "测试.event",
				Record: map[string]interface{}{
					"message": "Unicode message: 你好世界 🌍",
					"emoji":   "🚀🎉",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The event should handle these edge cases gracefully
			result := tt.event.getTime()

			// Either parse successfully or fallback to current time
			if _, err := time.Parse(time.RFC3339, tt.event.Time); err != nil {
				// Should fallback to current time
				now := time.Now()
				diff := now.Sub(result)
				assert.True(t, diff >= 0 && diff < 5*time.Second)
			} else {
				// Should parse correctly
				expected, _ := time.Parse(time.RFC3339, tt.event.Time)
				assert.Equal(t, expected, result)
			}
		})
	}
}

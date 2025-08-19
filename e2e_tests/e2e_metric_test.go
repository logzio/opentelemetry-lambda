//go:build e2e

package e2e

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EMetrics(t *testing.T) {
	skipIfEnvVarsMissing(t, t.Name())
	e2eLogger.Infof("Starting E2E Metrics Test for environment: %s", e2eTestEnvironmentLabel)

	expectedFaasName := os.Getenv("EXPECTED_LAMBDA_FUNCTION_NAME")
	require.NotEmpty(t, expectedFaasName, "EXPECTED_LAMBDA_FUNCTION_NAME environment variable must be set")

	expectedServiceName := os.Getenv("EXPECTED_SERVICE_NAME")
	require.NotEmpty(t, expectedServiceName, "EXPECTED_SERVICE_NAME environment variable must be set")

	e2eLogger.Infof("Expecting metrics with common labels - faas.name: %s, service_name: %s, environment: %s", expectedFaasName, expectedServiceName, e2eTestEnvironmentLabel)

	// Note: The environment label will be dynamic (python-e2e-{GITHUB_RUN_ID}), but we'll still validate it in assertions
	query := `{faas_name="one-layer-e2e-test-python", service_name="logzio-e2e-python-service"}`
	e2eLogger.Infof("Querying for any metrics matching: %s", query)

	metricResponse, err := fetchLogzMetricsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, query)

	if err != nil {
		if errors.Is(err, ErrNoDataFoundAfterRetries) {
			t.Fatalf("Failed to find metrics after all retries for query '%s': %v", query, err)
		} else {
			t.Fatalf("Error fetching metrics for query '%s': %v", query, err)
		}
	}
	require.NotNil(t, metricResponse, "Metric response should not be nil if error is nil")
	require.Equal(t, "success", metricResponse.Status, "Metric API status should be success")
	require.GreaterOrEqual(t, len(metricResponse.Data.Result), 1, "Should find at least one metric series matching the core labels. Query: %s", query)

	e2eLogger.Info("Validating labels on the first found metric series...")
	firstSeries := metricResponse.Data.Result[0]
	metricLabels := firstSeries.Metric
	e2eLogger.Infof("Found metric '%s' with labels: %+v", metricLabels["__name__"], metricLabels)

	assert.Equal(t, e2eTestEnvironmentLabel, metricLabels["environment"], "Label 'environment' mismatch")
	assert.Equal(t, expectedFaasName, metricLabels["faas_name"], "Label 'faas_name' mismatch")
	assert.Equal(t, expectedServiceName, metricLabels["service_name"], "Label 'service_name' mismatch")
	assert.Equal(t, "aws_lambda", metricLabels["cloud_platform"], "Label 'cloud_platform' should be 'aws_lambda'")
	assert.Equal(t, "aws", metricLabels["cloud_provider"], "Label 'cloud_provider' should be 'aws'")
	assert.NotEmpty(t, metricLabels["cloud_region"], "Label 'cloud_region' should be present")

	if metricName, ok := metricLabels["__name__"]; ok && (metricName == "aws_lambda_duration_milliseconds" || metricName == "aws_lambda_maxMemoryUsed_megabytes" || metricName == "aws_lambda_invocations" || metricName == "aws_lambda_errors") {
		assert.NotEmpty(t, metricLabels["faas_execution"], "Label 'faas_execution' (Lambda Request ID) should be present for AWS platform metrics")
	}

	foundDurationMetric := false
	for _, series := range metricResponse.Data.Result {
		if series.Metric["__name__"] == "aws_lambda_duration_milliseconds" {
			foundDurationMetric = true
			e2eLogger.Info("Confirmed 'aws_lambda_duration_milliseconds' is among the found metrics with correct labels.")
			break
		}
	}
	assert.True(t, foundDurationMetric, "Expected 'aws_lambda_duration_milliseconds' to be one of the metrics reported with the correct labels.")
	e2eLogger.Info("E2E Metrics Test: Core label validation successful.")
}

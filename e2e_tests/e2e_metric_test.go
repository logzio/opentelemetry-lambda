//go:build e2e

package e2e

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EMetrics(t *testing.T) {
	skipIfEnvVarsMissing(t, t.Name())
	e2eLogger.Infof("Starting E2E Metrics Test for environment: %s", e2eTestEnvironmentLabel)

	expectedServiceName := os.Getenv("EXPECTED_SERVICE_NAME")
	require.NotEmpty(t, expectedServiceName, "EXPECTED_SERVICE_NAME environment variable must be set")

	// We'll validate two representative metrics visible in Logz.io Grafana
	metricsToCheck := []string{
		"aws_lambda_billedDurationMs_milliseconds",
		"http_client_duration_milliseconds_count",
	}

	for _, metricName := range metricsToCheck {
		promql := fmt.Sprintf(`%s{job="%s"}`, metricName, expectedServiceName)
		e2eLogger.Infof("Querying metrics: %s", promql)

		metricResponse, err := fetchLogzMetricsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, promql)
		if err != nil {
			if errors.Is(err, ErrNoDataFoundAfterRetries) {
				t.Fatalf("Failed to find metrics after all retries for query '%s': %v", promql, err)
			} else {
				t.Fatalf("Error fetching metrics for query '%s': %v", promql, err)
			}
		}
		require.NotNil(t, metricResponse, "Metric response should not be nil if error is nil")
		require.Equal(t, "success", metricResponse.Status, "Metric API status should be success")
		require.GreaterOrEqual(t, len(metricResponse.Data.Result), 1, "Should find at least one series for %s with job=%s", metricName, expectedServiceName)

		first := metricResponse.Data.Result[0]
		labels := first.Metric
		assert.Equal(t, expectedServiceName, labels["job"], "metric %s should have job=%s", metricName, expectedServiceName)

		if metricName == "http_client_duration_milliseconds_count" {
			// Optional helpful context if present
			if v := labels["http_host"]; v != "" {
				e2eLogger.Infof("http_host=%s", v)
			}
			if v := labels["http_method"]; v != "" {
				e2eLogger.Infof("http_method=%s", v)
			}
			if v := labels["http_status_code"]; v != "" {
				e2eLogger.Infof("http_status_code=%s", v)
			}
		}
	}

	// Also verify the labels endpoint returns common labels (e.g., job)
	start := time.Now().Add(-30 * time.Minute)
	end := time.Now()
	labelsResp, err := fetchLogzMetricsLabelsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, start, end, []string{fmt.Sprintf(`{job="%s"}`, expectedServiceName)}, 200)
	if err != nil {
		if errors.Is(err, ErrNoDataFoundAfterRetries) {
			t.Fatalf("Failed to retrieve labels after all retries: %v", err)
		} else {
			t.Fatalf("Error retrieving labels: %v", err)
		}
	}
	require.NotNil(t, labelsResp)
	require.Equal(t, "success", labelsResp.Status)
	foundJob := false
	for _, l := range labelsResp.Data {
		if l == "job" {
			foundJob = true
			break
		}
	}
	require.True(t, foundJob, "labels endpoint should include 'job'")

	e2eLogger.Info("E2E Metrics Test: Specific metric validation successful.")
}

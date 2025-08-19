//go:build e2e

package e2e

import (
	"errors"
	"fmt"
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

	e2eLogger.Infof("Validating presence of key metrics for job: %s (function: %s)", expectedServiceName, expectedFaasName)

	// Helper to run a PromQL query and assert results
	runQuery := func(t *testing.T, promql string) *logzioPrometheusResponse {
		e2eLogger.Infof("Querying metrics: %s", promql)
		metricResponse, err := fetchLogzMetricsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, promql)
		if err != nil {
			if errors.Is(err, ErrNoDataFoundAfterRetries) {
				t.Fatalf("No metrics found after retries for query '%s': %v", promql, err)
			} else {
				t.Fatalf("Error fetching metrics for query '%s': %v", promql, err)
			}
		}
		require.NotNil(t, metricResponse)
		require.Equal(t, "success", metricResponse.Status)
		require.GreaterOrEqual(t, len(metricResponse.Data.Result), 1, "Expected at least one series for query: %s", promql)
		return metricResponse
	}

	// 1) AWS Lambda platform metrics (names as seen in Grafana)
	awsLambdaMetrics := []string{
		"aws_lambda_billedDurationMs_milliseconds",
		"aws_lambda_durationMs_milliseconds",
		"aws_lambda_initDurationMs_milliseconds",
		"aws_lambda_maxMemoryUsedMB_bytes",
		"aws_lambda_memorySizeMB_bytes",
	}
	for _, m := range awsLambdaMetrics {
		promql := fmt.Sprintf(`%s{job="%s"}`, m, expectedServiceName)
		t.Run(m, func(t *testing.T) {
			resp := runQuery(t, promql)
			first := resp.Data.Result[0].Metric
			// Basic label sanity
			assert.Equal(t, expectedServiceName, first["job"], "job label should match service name")
		})
	}

	// 2) HTTP client duration metrics (count/sum/bucket) for httpbin GETs
	httpMetrics := []string{
		"http_client_duration_milliseconds_count",
		"http_client_duration_milliseconds_sum",
		"http_client_duration_milliseconds_bucket",
	}
	for _, m := range httpMetrics {
		// Filter by job, host and method as in typical dashboards
		promql := fmt.Sprintf(`%s{job="%s",http_host="httpbin.org",http_method="GET"}`, m, expectedServiceName)
		t.Run(m, func(t *testing.T) {
			resp := runQuery(t, promql)
			first := resp.Data.Result[0].Metric
			assert.Equal(t, expectedServiceName, first["job"], "job label should match service name")
			if host, ok := first["http_host"]; ok {
				assert.Equal(t, "httpbin.org", host)
			}
			if method, ok := first["http_method"]; ok {
				assert.Equal(t, "GET", method)
			}
		})
	}

	// 3) Optional: verify that platform metrics include a Lambda invocation label when present
	promqlExec := fmt.Sprintf(`aws_lambda_durationMs_milliseconds{job="%s"}`, expectedServiceName)
	resp := runQuery(t, promqlExec)
	first := resp.Data.Result[0].Metric
	if _, ok := first["faas_invocation_id"]; ok {
		assert.NotEmpty(t, first["faas_invocation_id"], "faas_invocation_id should not be empty when present")
	}
}

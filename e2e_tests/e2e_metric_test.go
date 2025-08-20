//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2EMetrics(t *testing.T) {
	skipIfEnvVarsMissing(t, t.Name())
	e2eLogger.Infof("Starting E2E Metrics Test for environment: %s", e2eTestEnvironmentLabel)

	expectedServiceName := os.Getenv("EXPECTED_SERVICE_NAME")
	require.NotEmpty(t, expectedServiceName, "EXPECTED_SERVICE_NAME environment variable must be set")

	// Metrics to validate (as seen in Grafana)
	metricsToCheck := []string{
		"aws_lambda_billedDurationMs_milliseconds",
		"http_client_duration_milliseconds_count",
	}

	// Candidate label selectors; we'll try them in order
	labelSelectors := []string{
		`{job="%s"}`,
		`{service_name="%s"}`,
	}

	for _, metricName := range metricsToCheck {
		found := false

		for _, selectorFmt := range labelSelectors {
			selector := fmt.Sprintf(selectorFmt, expectedServiceName)

			// 1) Instant query
			promql := fmt.Sprintf(`%s%s`, metricName, selector)
			e2eLogger.Infof("Querying metrics (instant): %s", promql)
			metricResponse, err := fetchLogzMetricsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, promql)
			if err == nil && metricResponse != nil && metricResponse.Status == "success" && len(metricResponse.Data.Result) > 0 {
				first := metricResponse.Data.Result[0]
				labels := first.Metric
				// Validate one of the labels matches expected service
				if labels["job"] == expectedServiceName || labels["service_name"] == expectedServiceName {
					found = true
					e2eLogger.Infof("Found series for %s with labels: %+v", metricName, labels)
					break
				}
			}

			// 2) Range query fallback (helps if no sample at the exact instant)
			var rangeQuery string
			if strings.Contains(metricName, "_count") {
				rangeQuery = fmt.Sprintf(`increase(%s%s[30m])`, metricName, selector)
			} else {
				rangeQuery = fmt.Sprintf(`max_over_time(%s%s[30m])`, metricName, selector)
			}
			e2eLogger.Infof("Querying metrics (range fallback): %s", rangeQuery)
			metricResponse, err = fetchLogzMetricsAPI(t, logzioMetricsQueryAPIKey, logzioMetricsQueryBaseURL, rangeQuery)
			if err == nil && metricResponse != nil && metricResponse.Status == "success" && len(metricResponse.Data.Result) > 0 {
				first := metricResponse.Data.Result[0]
				labels := first.Metric
				if labels["job"] == expectedServiceName || labels["service_name"] == expectedServiceName {
					found = true
					e2eLogger.Infof("Found series for %s (range) with labels: %+v", metricName, labels)
					break
				}
			}
		}

		if !found {
			t.Fatalf("Failed to find any series for metric %s with expected service label (job or service_name=%s)", metricName, expectedServiceName)
		}
	}

	e2eLogger.Info("E2E Metrics Test: Specific metric validation successful.")
}

//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2ELogs(t *testing.T) {
	skipIfEnvVarsMissing(t, t.Name())
	e2eLogger.Infof("Starting E2E Log Test for environment label: %s", e2eTestEnvironmentLabel)

	expectedServiceName := os.Getenv("EXPECTED_SERVICE_NAME")
	require.NotEmpty(t, expectedServiceName, "EXPECTED_SERVICE_NAME environment variable must be set for log tests")
	expectedFaasName := os.Getenv("EXPECTED_LAMBDA_FUNCTION_NAME")
	require.NotEmpty(t, expectedFaasName, "EXPECTED_LAMBDA_FUNCTION_NAME must be set for log tests")

	// Query for logs from our function - start with basic search
	baseQuery := fmt.Sprintf(`faas.name:"%q"`, expectedFaasName)

	logChecks := []struct {
		name        string
		mustContain string
		assertion   func(t *testing.T, hits []map[string]interface{})
	}{
		{
			name:        "telemetry_api_subscription",
			mustContain: `"Successfully subscribed to Telemetry API"`,
			assertion: func(t *testing.T, hits []map[string]interface{}) {
				assert.GreaterOrEqual(t, len(hits), 1, "Should find telemetry API subscription log")
				hit := hits[0]
				assert.Equal(t, expectedFaasName, hit["faas.name"])
			},
		},
		{
			name:        "function_invocation_log",
			mustContain: `"📍 Lambda invocation started"`,
			assertion: func(t *testing.T, hits []map[string]interface{}) {
				assert.GreaterOrEqual(t, len(hits), 1, "Should find function invocation start log")
				hit := hits[0]
				assert.Equal(t, expectedFaasName, hit["faas.name"])
			},
		},
	}

	allChecksPassed := true

	for _, check := range logChecks {
		t.Run(check.name, func(t *testing.T) {
			query := fmt.Sprintf(`%s AND %s`, baseQuery, check.mustContain)
			e2eLogger.Infof("Querying for logs: %s", query)

			logResponse, err := fetchLogzSearchAPI(t, logzioLogsQueryAPIKey, logzioAPIURL, query, "logs")
			if err != nil {
				e2eLogger.Errorf("Failed to fetch logs for check '%s' after all retries: %v", check.name, err)
				allChecksPassed = false
				t.Fail()
				return
			}

			require.NotNil(t, logResponse, "Log response should not be nil if error is nil for check '%s'", check.name)

			var sources []map[string]interface{}
			for _, hit := range logResponse.Hits.Hits {
				sources = append(sources, hit.Source)
				if len(sources) <= 2 {
					logSample, _ := json.Marshal(hit.Source)
					e2eLogger.Debugf("Sample log for check '%s': %s", check.name, string(logSample))
				}
			}

			if check.assertion != nil {
				check.assertion(t, sources)
			}
		})
	}

	require.True(t, allChecksPassed, "One or more E2E log checks failed.")
	e2eLogger.Info("E2E Log Test Completed Successfully.")
}

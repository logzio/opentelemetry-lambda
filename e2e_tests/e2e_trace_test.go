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

func TestE2ETraces(t *testing.T) {
	skipIfEnvVarsMissing(t, t.Name())
	e2eLogger.Infof("Starting E2E Trace Test for environment: %s", e2eTestEnvironmentLabel)

	tracesQueryKey := logzioTracesQueryAPIKey
	expectedFaasName := os.Getenv("EXPECTED_LAMBDA_FUNCTION_NAME")
	require.NotEmpty(t, expectedFaasName, "EXPECTED_LAMBDA_FUNCTION_NAME must be set")
	expectedServiceName := os.Getenv("EXPECTED_SERVICE_NAME")
	require.NotEmpty(t, expectedServiceName, "EXPECTED_SERVICE_NAME must be set")

	e2eLogger.Infof("Expecting traces for service: %s, function: %s, environment: %s", expectedServiceName, expectedFaasName, e2eTestEnvironmentLabel)

	// Simple query for any traces from our service and function
	query := fmt.Sprintf(`type:jaegerSpan AND process\.serviceName:"%s" AND process\.tag\.faas@name:"%s"`, expectedServiceName, expectedFaasName)
	e2eLogger.Infof("Querying for traces: %s", query)

	traceResponse, err := fetchLogzSearchAPI(t, tracesQueryKey, logzioAPIURL, query, "traces")
	require.NoError(t, err, "Failed to find any matching traces after all retries.")
	require.NotNil(t, traceResponse, "Trace response should not be nil if no error was returned")
	require.GreaterOrEqual(t, traceResponse.getTotalHits(), 1, "Should find at least one trace matching the query.")

	e2eLogger.Info("✅ Found traces! Validating content of the first trace...")

	hit := traceResponse.Hits.Hits[0].Source
	logSample, _ := json.Marshal(hit)
	e2eLogger.Debugf("Sample trace for validation: %s", string(logSample))

	// Basic content checks
	assert.Equal(t, expectedServiceName, getNestedValue(hit, "process", "serviceName"))
	assert.Equal(t, expectedFaasName, getNestedValue(hit, "process", "tag", "faas@name"))

	e2eLogger.Info("E2E Trace Test Completed Successfully.")
}

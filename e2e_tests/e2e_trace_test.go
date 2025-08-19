//go:build e2e

package e2e

import (
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

	baseQuery := fmt.Sprintf(`type:jaegerSpan AND process.serviceName:"%s" AND process.tag.faas@name:"%s"`, expectedServiceName, expectedFaasName)

	// 1) Fetch platform/server span for the Lambda handler
	serverQuery := baseQuery + " AND JaegerTag.span@kind:server"
	e2eLogger.Infof("Querying for server span: %s", serverQuery)
	serverResp, err := fetchLogzSearchAPI(t, tracesQueryKey, logzioAPIURL, serverQuery, "traces")
	require.NoError(t, err, "Failed to find server span after all retries.")
	require.NotNil(t, serverResp)
	require.GreaterOrEqual(t, serverResp.getTotalHits(), 1, "Should find at least one server span.")

	serverHit := serverResp.Hits.Hits[0].Source

	// Extract traceID from common variants
	candKeys := []string{"traceID", "traceId", "trace_id"}
	var traceID string
	for _, k := range candKeys {
		if v, ok := serverHit[k]; ok {
			if s, ok2 := v.(string); ok2 && s != "" {
				traceID = s
				break
			}
		}
	}
	if traceID == "" {
		// Log available keys to aid debugging if missing
		keys := make([]string, 0, len(serverHit))
		for k := range serverHit {
			keys = append(keys, k)
		}
		e2eLogger.Warnf("traceID not found on server span. Available keys: %v", keys)
	}
	require.NotEmpty(t, traceID, "traceID should be present on server span")
	e2eLogger.Infof("Found server span with traceID: %s", traceID)

	// Basic content checks for server span
	assert.Equal(t, expectedServiceName, getNestedValue(serverHit, "process", "serviceName"))
	assert.Equal(t, expectedFaasName, getNestedValue(serverHit, "process", "tag", "faas@name"))

	// 2) Fetch custom/client spans within the same trace
	clientQuery := fmt.Sprintf(`type:jaegerSpan AND traceID:"%s" AND JaegerTag.span@kind:client`, traceID)
	e2eLogger.Infof("Querying for client spans in same trace: %s", clientQuery)
	clientResp, err := fetchLogzSearchAPI(t, tracesQueryKey, logzioAPIURL, clientQuery, "traces")
	require.NoError(t, err, "Failed to find client spans for the trace after all retries.")
	require.NotNil(t, clientResp)
	require.GreaterOrEqual(t, clientResp.getTotalHits(), 1, "Should find at least one client span in the same trace.")

	clientHit := clientResp.Hits.Hits[0].Source
	// Optional light checks for client spans
	if m := getNestedValue(clientHit, "JaegerTag.http@method"); m != nil {
		e2eLogger.Infof("Client span HTTP method: %v", m)
	}
	if sc := getNestedValue(clientHit, "JaegerTag.http@status_code"); sc != nil {
		e2eLogger.Infof("Client span HTTP status: %v", sc)
	}

	e2eLogger.Info("E2E Trace Test Completed Successfully.")
}

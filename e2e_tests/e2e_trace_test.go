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

	// Verify at least one platform/server span exists
	serverQuery := baseQuery + " AND JaegerTag.span@kind:server"
	e2eLogger.Infof("Querying for server span: %s", serverQuery)
	serverResp, err := fetchLogzSearchAPI(t, tracesQueryKey, logzioAPIURL, serverQuery, "traces")
	require.NoError(t, err, "Failed to find server span after all retries.")
	require.NotNil(t, serverResp)
	require.GreaterOrEqual(t, serverResp.getTotalHits(), 1, "Should find at least one server span.")
	serverHit := serverResp.Hits.Hits[0].Source
	assert.Equal(t, expectedServiceName, getNestedValue(serverHit, "process", "serviceName"))
	assert.Equal(t, expectedFaasName, getNestedValue(serverHit, "process", "tag", "faas@name"))

	// Verify at least one custom/client span exists
	clientQuery := baseQuery + " AND JaegerTag.span@kind:client"
	e2eLogger.Infof("Querying for client spans: %s", clientQuery)
	clientResp, err := fetchLogzSearchAPI(t, tracesQueryKey, logzioAPIURL, clientQuery, "traces")
	require.NoError(t, err, "Failed to find client spans after all retries.")
	require.NotNil(t, clientResp)
	require.GreaterOrEqual(t, clientResp.getTotalHits(), 1, "Should find at least one client span.")

	clientHit := clientResp.Hits.Hits[0].Source
	if m := getNestedValue(clientHit, "JaegerTag.http@method"); m != nil {
		e2eLogger.Infof("Client span HTTP method: %v", m)
	}
	if sc := getNestedValue(clientHit, "JaegerTag.http@status_code"); sc != nil {
		e2eLogger.Infof("Client span HTTP status: %v", sc)
	}

	e2eLogger.Info("E2E Trace Test Completed Successfully.")
}

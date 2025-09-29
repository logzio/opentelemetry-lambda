//go:build e2e

package e2e

import (
	"testing"
	"time"
)

func TestE2ERunner(t *testing.T) {
	e2eLogger.Info("E2E Test Runner: Waiting 60 seconds for initial Lambda execution and data ingestion before starting tests...")
	time.Sleep(60 * time.Second)

	initTimeTracking()
	e2eLogger.Infof("E2E Test Runner starting with a total budget of %d seconds.", totalBudgetSeconds)
	e2eLogger.Info("Tests will run in order: Metrics -> Logs -> Traces.")

	t.Run("E2EMetricsTest", func(t *testing.T) {
		e2eLogger.Info("=== Starting E2E Metrics Test ===")
		startTime := time.Now()
		TestE2EMetrics(t)
		duration := time.Since(startTime)
		recordTimeSpent("metrics", duration)
		e2eLogger.Infof("=== E2E Metrics Test completed in %.1f seconds ===", duration.Seconds())
	})

	if t.Failed() {
		e2eLogger.Error("Metrics test or previous setup failed. Subsequent tests might be affected or also fail.")
	}

	t.Run("E2ELogsTest", func(t *testing.T) {
		e2eLogger.Info("=== Starting E2E Logs Test ===")
		startTime := time.Now()
		TestE2ELogs(t)
		duration := time.Since(startTime)
		recordTimeSpent("logs", duration)
		e2eLogger.Infof("=== E2E Logs Test completed in %.1f seconds ===", duration.Seconds())
	})

	if t.Failed() {
		e2eLogger.Error("Logs test or previous setup/tests failed. Subsequent tests might be affected or also fail.")
	}

	t.Run("E2ETracesTest", func(t *testing.T) {
		e2eLogger.Info("=== Starting E2E Traces Test ===")
		startTime := time.Now()
		TestE2ETraces(t)
		duration := time.Since(startTime)
		recordTimeSpent("traces", duration)
		e2eLogger.Infof("=== E2E Traces Test completed in %.1f seconds ===", duration.Seconds())
	})

	totalElapsed := time.Since(testStartTime)
	e2eLogger.Infof("E2E Test Runner finished all tests in %.1f seconds. Remaining budget: %ds", totalElapsed.Seconds(), getRemainingBudgetSeconds())

	if t.Failed() {
		e2eLogger.Error("One or more E2E tests failed.")
	} else {
		e2eLogger.Info("All E2E tests passed successfully!")
	}
}

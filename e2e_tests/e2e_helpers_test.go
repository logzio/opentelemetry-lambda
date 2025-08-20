//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var e2eLogger = logrus.WithField("test_type", "e2e")

var (
	logzioLogsQueryAPIKey     = os.Getenv("LOGZIO_API_KEY")
	logzioAPIURL              = os.Getenv("LOGZIO_API_URL")
	e2eTestEnvironmentLabel   = os.Getenv("E2E_TEST_ENVIRONMENT_LABEL")
	logzioMetricsQueryAPIKey  = os.Getenv("LOGZIO_API_METRICS_KEY")
	logzioMetricsQueryBaseURL = os.Getenv("LOGZIO_METRICS_QUERY_URL")
	logzioTracesQueryAPIKey   = os.Getenv("LOGZIO_API_TRACES_KEY")
)

var (
	totalBudgetSeconds = 400
	testStartTime      time.Time
	timeSpentMetrics   time.Duration
	timeSpentLogs      time.Duration
	timeSpentTraces    time.Duration
)

func initTimeTracking() {
	testStartTime = time.Now()
	timeSpentMetrics = 0
	timeSpentLogs = 0
	timeSpentTraces = 0
}

func getRemainingBudgetSeconds() int {
	elapsed := time.Since(testStartTime)
	remaining := time.Duration(totalBudgetSeconds)*time.Second - elapsed
	return max(0, int(remaining.Seconds()))
}

func getDynamicRetryConfig(testType string) (maxRetries int, retryDelay time.Duration) {
	defaultMaxRetries := 30
	defaultRetryDelay := 10 * time.Second

	remainingBudget := getRemainingBudgetSeconds()
	retryDelay = defaultRetryDelay

	var allocatedBudgetPortion float64
	switch testType {
	case "metrics":
		allocatedBudgetPortion = 0.1
	case "logs":
		allocatedBudgetPortion = 0.6
	case "traces":
		allocatedBudgetPortion = 0.3
	default:
		allocatedBudgetPortion = 0.2
	}

	var effectiveBudget int
	if timeSpentMetrics == 0 && timeSpentLogs == 0 && timeSpentTraces == 0 {
		effectiveBudget = int(float64(totalBudgetSeconds) * allocatedBudgetPortion)
	} else {
		effectiveBudget = int(float64(remainingBudget) * allocatedBudgetPortion)
	}

	effectiveBudget = max(effectiveBudget, int(defaultRetryDelay.Seconds())*2+1)

	maxRetries = effectiveBudget / int(defaultRetryDelay.Seconds())
	maxRetries = max(2, min(maxRetries, defaultMaxRetries))

	e2eLogger.Infof("Time budget for %s: %d attempts (delay %s). Total remaining: %ds. Effective budget for this test: %ds", testType, maxRetries, retryDelay, remainingBudget, effectiveBudget)
	return maxRetries, retryDelay
}

func recordTimeSpent(testType string, duration time.Duration) {
	switch testType {
	case "metrics":
		timeSpentMetrics += duration
	case "logs":
		timeSpentLogs += duration
	case "traces":
		timeSpentTraces += duration
	}
	total := timeSpentMetrics + timeSpentLogs + timeSpentTraces
	e2eLogger.Infof("Time spent - Metrics: %.1fs, Logs: %.1fs, Traces: %.1fs, Total: %.1fs/%ds", timeSpentMetrics.Seconds(), timeSpentLogs.Seconds(), timeSpentTraces.Seconds(), total.Seconds(), totalBudgetSeconds)
}

const (
	apiTimeout     = 45 * time.Second
	searchLookback = "30m"
)

var ErrNoDataFoundAfterRetries = errors.New("no data found after all retries")

func skipIfEnvVarsMissing(t *testing.T, testName string) {
	baseRequired := []string{"E2E_TEST_ENVIRONMENT_LABEL"}
	specificRequiredMissing := false

	if logzioAPIURL == "" {
		e2eLogger.Errorf("Skipping E2E test %s: Missing base required environment variable LOGZIO_API_URL.", testName)
		t.Skipf("Skipping E2E test %s: Missing base required environment variable LOGZIO_API_URL.", testName)
		return
	}

	if strings.Contains(testName, "Logs") || strings.Contains(testName, "E2ELogsTest") {
		if logzioLogsQueryAPIKey == "" {
			e2eLogger.Errorf("Skipping E2E Log test %s: Missing LOGZIO_API_KEY.", testName)
			t.Skipf("Skipping E2E Log test %s: Missing LOGZIO_API_KEY.", testName)
			specificRequiredMissing = true
		}
	}
	if strings.Contains(testName, "Metrics") || strings.Contains(testName, "E2EMetricsTest") {
		if logzioMetricsQueryAPIKey == "" {
			e2eLogger.Errorf("Skipping E2E Metrics test %s: Missing LOGZIO_API_METRICS_KEY.", testName)
			t.Skipf("Skipping E2E Metrics test %s: Missing LOGZIO_API_METRICS_KEY.", testName)
			specificRequiredMissing = true
		}
		if logzioMetricsQueryBaseURL == "" {
			e2eLogger.Errorf("Skipping E2E Metrics test %s: Missing LOGZIO_METRICS_QUERY_URL.", testName)
			t.Skipf("Skipping E2E Metrics test %s: Missing LOGZIO_METRICS_QUERY_URL.", testName)
			specificRequiredMissing = true
		}
	}
	if strings.Contains(testName, "Traces") || strings.Contains(testName, "E2ETracesTest") {
		if logzioTracesQueryAPIKey == "" {
			e2eLogger.Errorf("Skipping E2E Traces test %s: Missing required environment variable LOGZIO_API_TRACES_KEY.", testName)
			t.Skipf("Skipping E2E Traces test %s: Missing required environment variable LOGZIO_API_TRACES_KEY.", testName)
			specificRequiredMissing = true
		}
	}

	if specificRequiredMissing {
		return
	}

	for _, v := range baseRequired {
		if os.Getenv(v) == "" {
			e2eLogger.Errorf("Skipping E2E test %s: Missing base required environment variable %s.", testName, v)
			t.Skipf("Skipping E2E test %s: Missing base required environment variable %s.", testName, v)
			return
		}
	}
}

type logzioSearchQueryBody struct {
	Query          map[string]interface{} `json:"query"`
	From           int                    `json:"from"`
	Size           int                    `json:"size"`
	Sort           []interface{}          `json:"sort"`
	Source         interface{}            `json:"_source"`
	PostFilter     interface{}            `json:"post_filter,omitempty"`
	DocvalueFields []string               `json:"docvalue_fields"`
	Version        bool                   `json:"version"`
	StoredFields   []string               `json:"stored_fields"`
	Highlight      map[string]interface{} `json:"highlight"`
	Aggregations   map[string]interface{} `json:"aggregations,omitempty"`
}

type logzioSearchResponse struct {
	Hits struct {
		Total json.RawMessage `json:"total"`
		Hits  []struct {
			Source map[string]interface{} `json:"_source"`
			Sort   []interface{}          `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
	Error *struct {
		Reason string `json:"reason"`
	} `json:"error,omitempty"`
}

func (r *logzioSearchResponse) getTotalHits() int {
	if len(r.Hits.Total) == 0 {
		return 0
	}
	var totalInt int
	if err := json.Unmarshal(r.Hits.Total, &totalInt); err == nil {
		return totalInt
	}
	var totalObj struct {
		Value int `json:"value"`
	}
	if err := json.Unmarshal(r.Hits.Total, &totalObj); err == nil {
		return totalObj.Value
	}
	e2eLogger.Warnf("Could not determine total hits from raw message: %s", string(r.Hits.Total))
	return 0
}

func fetchLogzSearchAPI(t *testing.T, apiKey, queryBaseAPIURL, luceneQuery string, testType string) (*logzioSearchResponse, error) {
	maxRetries, retryDelay := getDynamicRetryConfig(testType)
	return fetchLogzSearchAPIWithRetries(t, apiKey, queryBaseAPIURL, luceneQuery, maxRetries, retryDelay)
}

func fetchLogzSearchAPIWithRetries(t *testing.T, apiKey, queryBaseAPIURL, luceneQuery string, maxRetries int, retryDelay time.Duration) (*logzioSearchResponse, error) {
	searchAPIEndpoint := fmt.Sprintf("%s/v1/search", strings.TrimSuffix(queryBaseAPIURL, "/"))

	// Build request body per Logz.io Search API example
	queryBodyMap := logzioSearchQueryBody{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"query_string": map[string]interface{}{
							"query":                  luceneQuery,
							"allow_leading_wildcard": false,
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": "now-30m",
								"lte": "now",
							},
						},
					},
				},
			},
		},
		From:           0,
		Size:           100,
		Sort:           []interface{}{map[string]interface{}{}},
		Source:         true,
		PostFilter:     nil,
		DocvalueFields: []string{"@timestamp"},
		Version:        true,
		StoredFields:   []string{"*"},
		Highlight:      map[string]interface{}{},
		Aggregations: map[string]interface{}{
			"byType": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "type",
					"size":  5,
				},
			},
		},
	}
	queryBytes, err := json.Marshal(queryBodyMap)
	require.NoError(t, err)

	// Debug: Log the actual JSON query being sent
	e2eLogger.Debugf("Logz.io search query JSON: %s", string(queryBytes))

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		e2eLogger.Infof("Attempt %d/%d to fetch Logz.io search results (Query: %s)...", i+1, maxRetries, luceneQuery)
		req, err := http.NewRequest("POST", searchAPIEndpoint, bytes.NewBuffer(queryBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-API-TOKEN", apiKey)
		client := &http.Client{Timeout: apiTimeout}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("API request failed on attempt %d: %w", i+1, err)
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		respBodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = fmt.Errorf("failed to read API response body on attempt %d: %w", i+1, readErr)
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API returned status %d on attempt %d: %s", resp.StatusCode, i+1, string(respBodyBytes))
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			e2eLogger.Debugf("Failed request body was: %s", string(queryBytes))
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		var logResponse logzioSearchResponse
		unmarshalErr := json.Unmarshal(respBodyBytes, &logResponse)
		if unmarshalErr != nil {
			lastErr = fmt.Errorf("failed to unmarshal API response on attempt %d: %w. Body: %s", i+1, unmarshalErr, string(respBodyBytes))
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if logResponse.Error != nil {
			lastErr = fmt.Errorf("Logz.io API error in response on attempt %d: %s", i+1, logResponse.Error.Reason)
			if strings.Contains(logResponse.Error.Reason, "parse_exception") || strings.Contains(logResponse.Error.Reason, "query_shard_exception") {
				e2eLogger.Errorf("Non-retryable API error encountered: %v", lastErr)
				return nil, lastErr
			}
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if logResponse.getTotalHits() > 0 {
			e2eLogger.Infof("Attempt %d successful. Found %d total hits.", i+1, logResponse.getTotalHits())
			return &logResponse, nil
		}
		lastErr = fmt.Errorf("attempt %d/%d: no data found for query '%s'", i+1, maxRetries, luceneQuery)
		e2eLogger.Infof("%s. Retrying in %s...", lastErr.Error(), retryDelay)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	e2eLogger.Warnf("No data found for query '%s' after %d retries.", luceneQuery, maxRetries)
	return nil, ErrNoDataFoundAfterRetries
}

type logzioPrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

func fetchLogzMetricsAPI(t *testing.T, apiKey, metricsAPIBaseURL, promqlQuery string) (*logzioPrometheusResponse, error) {
	maxRetries, retryDelay := getDynamicRetryConfig("metrics")
	return fetchLogzMetricsAPIWithRetries(t, apiKey, metricsAPIBaseURL, promqlQuery, maxRetries, retryDelay)
}

func fetchLogzMetricsAPIWithRetries(t *testing.T, apiKey, metricsAPIBaseURL, promqlQuery string, maxRetries int, retryDelay time.Duration) (*logzioPrometheusResponse, error) {
	queryAPIEndpoint := fmt.Sprintf("%s/v1/metrics/prometheus/api/v1/query?query=%s", strings.TrimSuffix(metricsAPIBaseURL, "/"), url.QueryEscape(promqlQuery))
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		e2eLogger.Infof("Attempt %d/%d to fetch Logz.io metrics (Query: %s)...", i+1, maxRetries, promqlQuery)
		req, err := http.NewRequest("GET", queryAPIEndpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("metrics API request creation failed: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-API-TOKEN", apiKey)

		client := &http.Client{Timeout: apiTimeout}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("metrics API request failed on attempt %d: %w", i+1, err)
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		respBodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = fmt.Errorf("failed to read metrics API response body on attempt %d: %w", i+1, readErr)
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("metrics API returned status %d on attempt %d: %s", resp.StatusCode, i+1, string(respBodyBytes))
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		var metricResponse logzioPrometheusResponse
		unmarshalErr := json.Unmarshal(respBodyBytes, &metricResponse)
		if unmarshalErr != nil {
			lastErr = fmt.Errorf("failed to unmarshal metrics API response on attempt %d: %w. Body: %s", i+1, unmarshalErr, string(respBodyBytes))
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if metricResponse.Status != "success" {
			lastErr = fmt.Errorf("Logz.io Metrics API returned status '%s' on attempt %d, ErrorType: '%s', Error: '%s'", metricResponse.Status, i+1, metricResponse.ErrorType, metricResponse.Error)
			e2eLogger.Warnf("%v. Retrying in %s...", lastErr, retryDelay)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
			continue
		}
		if len(metricResponse.Data.Result) > 0 {
			e2eLogger.Infof("Attempt %d successful. Found %d metric series.", i+1, len(metricResponse.Data.Result))
			return &metricResponse, nil
		}
		lastErr = fmt.Errorf("attempt %d/%d: no data found for query '%s'", i+1, maxRetries, promqlQuery)
		e2eLogger.Infof("%s. Retrying in %s...", lastErr.Error(), retryDelay)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	e2eLogger.Warnf("No data found for query '%s' after %d retries.", promqlQuery, maxRetries)
	return nil, ErrNoDataFoundAfterRetries
}

func getNestedValue(data map[string]interface{}, path ...string) interface{} {
	var current interface{} = data
	for _, key := range path {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[key]
		if !ok {
			return nil
		}
	}
	return current
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

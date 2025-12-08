//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL        = "http://localhost:8080"
	apiPrefix      = "/api/v1"
	requestTimeout = 10 * time.Second
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: requestTimeout,
	}
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	// Verify environment variables are set
	requiredEnvVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Warning: Environment variable %s is not set\n", envVar)
		}
	}

	// Wait for server to be ready
	fmt.Println("Waiting for server to be ready...")
	if err := waitForServer(baseURL+apiPrefix+"/crypto/market", 30*time.Second); err != nil {
		fmt.Printf("Warning: Server may not be ready: %v\n", err)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// waitForServer waits for the server to be ready
func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := httpClient.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("server did not become ready within %v", timeout)
}

// TestGetMarketData tests GET /api/v1/crypto/market endpoint
func TestGetMarketData(t *testing.T) {
	url := baseURL + apiPrefix + "/crypto/market"

	resp, err := httpClient.Get(url)
	require.NoError(t, err, "Failed to make request to %s", url)
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK")

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Expected JSON content type")

	// Parse response
	var marketResponse model.MarketResponse
	err = json.NewDecoder(resp.Body).Decode(&marketResponse)
	require.NoError(t, err, "Failed to decode response body")

	// Verify response structure
	assert.NotEmpty(t, marketResponse.Data, "Data field should not be empty")
	assert.NotZero(t, marketResponse.Timestamp, "Timestamp field should not be zero")

	// Verify we have at least Bitcoin and Ethereum
	assert.GreaterOrEqual(t, len(marketResponse.Data), 2, "Should have at least 2 cryptocurrencies")

	// Verify each crypto data has required fields
	for _, crypto := range marketResponse.Data {
		assert.NotEmpty(t, crypto.ID, "ID should not be empty")
		assert.NotEmpty(t, crypto.Name, "Name should not be empty")
		assert.NotEmpty(t, crypto.Symbol, "Symbol should not be empty")
		assert.NotEmpty(t, crypto.Pair, "Pair should not be empty")
		assert.NotEmpty(t, crypto.Icon, "Icon should not be empty")
		assert.NotEmpty(t, crypto.IconColor, "IconColor should not be empty")
		assert.Greater(t, crypto.CurrentPrice, 0.0, "CurrentPrice should be greater than 0")
		// Note: ChangePercent can be negative, so we don't check > 0
		assert.NotNil(t, crypto.ChartData, "ChartData should not be nil")

		t.Logf("Crypto: %s, Price: %.2f JPY, Change: %.2f%%", crypto.Name, crypto.CurrentPrice, crypto.ChangePercent)
	}
}

// TestGetCryptoByID tests GET /api/v1/crypto/:id endpoint
func TestGetCryptoByID(t *testing.T) {
	testCases := []struct {
		name           string
		cryptoID       string
		expectedStatus int
		shouldHaveData bool
	}{
		{
			name:           "Get Bitcoin data",
			cryptoID:       "bitcoin",
			expectedStatus: http.StatusOK,
			shouldHaveData: true,
		},
		{
			name:           "Get Ethereum data",
			cryptoID:       "ethereum",
			expectedStatus: http.StatusOK,
			shouldHaveData: true,
		},
		{
			name:           "Get non-existent crypto",
			cryptoID:       "nonexistent",
			expectedStatus: http.StatusNotFound,
			shouldHaveData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s/crypto/%s", baseURL, apiPrefix, tc.cryptoID)

			resp, err := httpClient.Get(url)
			require.NoError(t, err, "Failed to make request to %s", url)
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected status %d", tc.expectedStatus)

			if tc.shouldHaveData {
				// Parse response
				var cryptoData model.CryptoData
				err = json.NewDecoder(resp.Body).Decode(&cryptoData)
				require.NoError(t, err, "Failed to decode response body")

				// Verify required fields
				assert.Equal(t, tc.cryptoID, cryptoData.ID, "ID should match")
				assert.NotEmpty(t, cryptoData.Name, "Name should not be empty")
				assert.NotEmpty(t, cryptoData.Symbol, "Symbol should not be empty")
				assert.NotEmpty(t, cryptoData.Pair, "Pair should not be empty")
				assert.NotEmpty(t, cryptoData.Icon, "Icon should not be empty")
				assert.NotEmpty(t, cryptoData.IconColor, "IconColor should not be empty")
				assert.Greater(t, cryptoData.CurrentPrice, 0.0, "CurrentPrice should be greater than 0")
				assert.NotNil(t, cryptoData.ChartData, "ChartData should not be nil")

				t.Logf("Crypto: %s, Price: %.2f JPY, Change: %.2f%%", cryptoData.Name, cryptoData.CurrentPrice, cryptoData.ChangePercent)
			} else {
				// Parse error response
				var errorResponse model.ErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				require.NoError(t, err, "Failed to decode error response body")

				assert.NotEmpty(t, errorResponse.Error, "Error field should not be empty")
				assert.NotEmpty(t, errorResponse.Message, "Message field should not be empty")

				t.Logf("Error: %s, Message: %s", errorResponse.Error, errorResponse.Message)
			}
		})
	}
}

// TestGetChartData tests GET /api/v1/crypto/:id/chart endpoint
func TestGetChartData(t *testing.T) {
	testCases := []struct {
		name           string
		cryptoID       string
		period         string
		expectedStatus int
		shouldHaveData bool
	}{
		{
			name:           "Get Bitcoin chart data (default period)",
			cryptoID:       "bitcoin",
			period:         "",
			expectedStatus: http.StatusOK,
			shouldHaveData: true,
		},
		{
			name:           "Get Ethereum chart data (7d)",
			cryptoID:       "ethereum",
			period:         "7d",
			expectedStatus: http.StatusOK,
			shouldHaveData: true,
		},
		{
			name:           "Get Bitcoin chart data (30d)",
			cryptoID:       "bitcoin",
			period:         "30d",
			expectedStatus: http.StatusOK,
			shouldHaveData: true,
		},
		{
			name:           "Get non-existent crypto chart",
			cryptoID:       "nonexistent",
			period:         "7d",
			expectedStatus: http.StatusNotFound,
			shouldHaveData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s/crypto/%s/chart", baseURL, apiPrefix, tc.cryptoID)
			if tc.period != "" {
				url += fmt.Sprintf("?period=%s", tc.period)
			}

			resp, err := httpClient.Get(url)
			require.NoError(t, err, "Failed to make request to %s", url)
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected status %d", tc.expectedStatus)

			if tc.shouldHaveData {
				// Parse response
				var chartResponse model.ChartResponse
				err = json.NewDecoder(resp.Body).Decode(&chartResponse)
				require.NoError(t, err, "Failed to decode response body")

				// Verify response structure
				assert.NotNil(t, chartResponse.Data, "Data field should not be nil")
				assert.NotEmpty(t, chartResponse.Period, "Period field should not be empty")

				// Verify chart data points
				if len(chartResponse.Data) > 0 {
					for _, point := range chartResponse.Data {
						assert.NotEmpty(t, point.Day, "Day should not be empty")
						assert.Greater(t, point.Price, 0.0, "Price should be greater than 0")
					}
					t.Logf("Chart data points: %d, Period: %s", len(chartResponse.Data), chartResponse.Period)
				} else {
					t.Logf("Warning: No chart data points returned (may be expected if no historical data)")
				}
			} else {
				// Parse error response
				var errorResponse model.ErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				require.NoError(t, err, "Failed to decode error response body")

				assert.NotEmpty(t, errorResponse.Error, "Error field should not be empty")
				assert.NotEmpty(t, errorResponse.Message, "Message field should not be empty")

				t.Logf("Error: %s, Message: %s", errorResponse.Error, errorResponse.Message)
			}
		})
	}
}

// TestCORSHeaders tests that CORS headers are properly set
func TestCORSHeaders(t *testing.T) {
	url := baseURL + apiPrefix + "/crypto/market"

	// Create OPTIONS request
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	require.NoError(t, err, "Failed to create OPTIONS request")

	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to make OPTIONS request")
	defer resp.Body.Close()

	// Verify CORS headers
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.NotEmpty(t, allowOrigin, "Access-Control-Allow-Origin header should be set")

	t.Logf("CORS Allow-Origin: %s", allowOrigin)
}

// TestServerHealth tests basic server health
func TestServerHealth(t *testing.T) {
	url := baseURL + apiPrefix + "/crypto/market"

	// Make multiple requests to ensure server is stable
	for i := 0; i < 5; i++ {
		resp, err := httpClient.Get(url)
		require.NoError(t, err, "Failed to make request %d", i+1)
		resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Request %d should return 200", i+1)
		time.Sleep(100 * time.Millisecond)
	}

	t.Log("Server health check passed: 5 consecutive requests successful")
}

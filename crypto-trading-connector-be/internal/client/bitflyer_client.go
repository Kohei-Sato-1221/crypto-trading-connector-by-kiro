package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
)

// BitFlyerClient implements CryptoExchangeClient for bitFlyer exchange
type BitFlyerClient struct {
	baseURL   string
	apiKey    string
	apiSecret string
	client    *http.Client
}

// NewBitFlyerClient creates a new bitFlyer API client
func NewBitFlyerClient(baseURL string) *BitFlyerClient {
	return &BitFlyerClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewBitFlyerClientWithAuth creates a new bitFlyer API client with authentication
func NewBitFlyerClientWithAuth(baseURL, apiKey, apiSecret string) *BitFlyerClient {
	return &BitFlyerClient{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTicker retrieves ticker information for a specific product
func (c *BitFlyerClient) GetTicker(productCode string) (*model.TickerResponse, error) {
	url := fmt.Sprintf("%s/v1/ticker?product_code=%s", c.baseURL, productCode)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call bitFlyer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bitFlyer API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ticker model.TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return nil, fmt.Errorf("failed to decode ticker response: %w", err)
	}

	return &ticker, nil
}

// GetBalance retrieves JPY balance from bitFlyer API
func (c *BitFlyerClient) GetBalance() (float64, error) {
	path := "/v1/me/getbalance"
	method := "GET"
	body := ""

	req, err := c.createAuthenticatedRequest(method, path, body)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call bitFlyer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bitFlyer API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var balances []model.BitFlyerBalance
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return 0, fmt.Errorf("failed to decode balance response: %w", err)
	}

	// Find JPY balance
	for _, balance := range balances {
		if balance.CurrencyCode == "JPY" {
			return balance.Available, nil
		}
	}

	return 0, fmt.Errorf("JPY balance not found")
}

// SendOrder sends an order to bitFlyer API
func (c *BitFlyerClient) SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
	path := "/v1/me/sendchildorder"
	method := "POST"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	httpReq, err := c.createAuthenticatedRequest(method, path, string(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call bitFlyer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bitFlyer API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var orderResp model.BitFlyerOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &orderResp, nil
}

// createAuthenticatedRequest creates an HTTP request with bitFlyer API authentication headers
func (c *BitFlyerClient) createAuthenticatedRequest(method, path, body string) (*http.Request, error) {
	url := c.baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Create request
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add authentication headers if credentials are provided
	if c.apiKey != "" && c.apiSecret != "" {
		// Create signature: HMAC-SHA256(timestamp + method + path + body, secret)
		text := timestamp + method + path + body
		signature := c.createSignature(text)

		req.Header.Set("ACCESS-KEY", c.apiKey)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-SIGN", signature)
	}

	return req, nil
}

// createSignature creates HMAC-SHA256 signature for bitFlyer API authentication
func (c *BitFlyerClient) createSignature(text string) string {
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}

// RoundPrice rounds the price according to the product code
// BTC_JPY: rounds to nearest million (1,000,000)
// ETH_JPY: rounds to nearest ten thousand (10,000)
func (c *BitFlyerClient) RoundPrice(price float64, productCode string) float64 {
	switch productCode {
	case "BTC_JPY":
		// Round to nearest million
		return float64(int(price/1000000)) * 1000000
	case "ETH_JPY":
		// Round to nearest ten thousand
		return float64(int(price/10000)) * 10000
	default:
		// Default: round to nearest integer
		return float64(int(price))
	}
}

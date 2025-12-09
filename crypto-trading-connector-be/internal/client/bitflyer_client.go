package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
)

// BitFlyerClient defines the interface for bitFlyer API operations
type BitFlyerClient interface {
	GetTicker(productCode string) (*model.TickerResponse, error)
	GetBalance() (float64, error)
	SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error)
}

// BitFlyerClientImpl implements BitFlyerClient
type BitFlyerClientImpl struct {
	baseURL   string
	apiKey    string
	apiSecret string
	client    *http.Client
}

// NewBitFlyerClient creates a new bitFlyer API client
func NewBitFlyerClient(baseURL string) *BitFlyerClientImpl {
	return &BitFlyerClientImpl{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewBitFlyerClientWithAuth creates a new bitFlyer API client with authentication
func NewBitFlyerClientWithAuth(baseURL, apiKey, apiSecret string) *BitFlyerClientImpl {
	return &BitFlyerClientImpl{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTicker retrieves ticker information for a specific product
func (c *BitFlyerClientImpl) GetTicker(productCode string) (*model.TickerResponse, error) {
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
func (c *BitFlyerClientImpl) GetBalance() (float64, error) {
	url := fmt.Sprintf("%s/v1/me/getbalance", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers (simplified - in production use proper HMAC)
	if c.apiKey != "" {
		req.Header.Set("ACCESS-KEY", c.apiKey)
		req.Header.Set("ACCESS-SIGN", c.apiSecret)
		req.Header.Set("ACCESS-TIMESTAMP", fmt.Sprintf("%d", time.Now().Unix()))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call bitFlyer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bitFlyer API returned status %d: %s", resp.StatusCode, string(body))
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
func (c *BitFlyerClientImpl) SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
	url := fmt.Sprintf("%s/v1/me/sendchildorder", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// Add authentication headers (simplified - in production use proper HMAC)
	if c.apiKey != "" {
		httpReq.Header.Set("ACCESS-KEY", c.apiKey)
		httpReq.Header.Set("ACCESS-SIGN", c.apiSecret)
		httpReq.Header.Set("ACCESS-TIMESTAMP", fmt.Sprintf("%d", time.Now().Unix()))
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

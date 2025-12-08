package client

import (
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
}

// BitFlyerClientImpl implements BitFlyerClient
type BitFlyerClientImpl struct {
	baseURL string
	client  *http.Client
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

package client

import "github.com/crypto-trading-connector/backend/internal/model"

// CryptoExchangeClient defines the common interface for all cryptocurrency exchange APIs
// This interface allows the application to support multiple exchanges (bitFlyer, Coinbase, Binance, etc.)
type CryptoExchangeClient interface {
	// GetTicker retrieves current ticker information for a specific trading pair
	// productCode: Trading pair identifier (e.g., "BTC_JPY", "ETH_JPY")
	GetTicker(productCode string) (*model.TickerResponse, error)

	// GetBalance retrieves the available balance in the base currency (JPY)
	GetBalance() (float64, error)

	// SendOrder submits a new order to the exchange
	SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error)
}

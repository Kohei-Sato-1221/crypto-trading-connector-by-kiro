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

	// GetChildOrders retrieves current orders (child orders) from the exchange
	// productCode: Trading pair identifier (e.g., "BTC_JPY", "ETH_JPY"), empty for all pairs
	// childOrderState: Order state filter (e.g., "ACTIVE"), empty for all states
	GetChildOrders(productCode string, childOrderState string) ([]model.BitflyerChildOrder, error)
}

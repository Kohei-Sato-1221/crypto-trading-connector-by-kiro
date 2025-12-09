# Exchange Client Architecture

## Overview

This document describes the architecture for supporting multiple cryptocurrency exchanges in the application.

## Design Pattern

The application uses the **Strategy Pattern** to support multiple cryptocurrency exchanges. This allows easy addition of new exchanges without modifying existing code.

## Architecture

```
┌─────────────────────────────────────┐
│   CryptoExchangeClient Interface    │
│  (Common interface for all exchanges)│
└─────────────────────────────────────┘
                 ▲
                 │ implements
    ┌────────────┼────────────┐
    │            │            │
┌───┴───┐   ┌───┴───┐   ┌───┴───┐
│bitFlyer│   │Coinbase│   │Binance│
│ Client │   │ Client │   │ Client │
└────────┘   └────────┘   └────────┘
```

## Interface Definition

```go
// CryptoExchangeClient defines the common interface for all cryptocurrency exchange APIs
type CryptoExchangeClient interface {
    // GetTicker retrieves current ticker information for a specific trading pair
    GetTicker(productCode string) (*model.TickerResponse, error)

    // GetBalance retrieves the available balance in the base currency (JPY)
    GetBalance() (float64, error)

    // SendOrder submits a new order to the exchange
    SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error)
}
```

## Current Implementation

### bitFlyer Client

**File**: `internal/client/bitflyer_client.go`

**Features**:
- HMAC-SHA256 authentication
- Support for public and private APIs
- Ticker retrieval
- Balance retrieval
- Order submission

**Usage**:
```go
// Without authentication (public API only)
client := client.NewBitFlyerClient("https://api.bitflyer.com")

// With authentication (private API)
client := client.NewBitFlyerClientWithAuth(
    "https://api.bitflyer.com",
    "your-api-key",
    "your-api-secret"
)
```

## Adding New Exchanges

To add support for a new exchange (e.g., Coinbase):

1. **Create a new client file**: `internal/client/coinbase_client.go`

```go
package client

import "github.com/crypto-trading-connector/backend/internal/model"

// CoinbaseClient implements CryptoExchangeClient for Coinbase exchange
type CoinbaseClient struct {
    baseURL   string
    apiKey    string
    apiSecret string
    client    *http.Client
}

// NewCoinbaseClient creates a new Coinbase API client
func NewCoinbaseClient(baseURL string) *CoinbaseClient {
    return &CoinbaseClient{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// Implement CryptoExchangeClient interface methods
func (c *CoinbaseClient) GetTicker(productCode string) (*model.TickerResponse, error) {
    // Coinbase-specific implementation
}

func (c *CoinbaseClient) GetBalance() (float64, error) {
    // Coinbase-specific implementation
}

func (c *CoinbaseClient) SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
    // Coinbase-specific implementation
}
```

2. **Update main.go** to support exchange selection:

```go
// Get exchange type from environment
exchangeType := getEnv("EXCHANGE_TYPE", "bitflyer")

var exchangeClient client.CryptoExchangeClient

switch exchangeType {
case "bitflyer":
    exchangeClient = client.NewBitFlyerClientWithAuth(...)
case "coinbase":
    exchangeClient = client.NewCoinbaseClientWithAuth(...)
default:
    log.Fatalf("Unsupported exchange: %s", exchangeType)
}
```

3. **Add environment variables** to `.env`:

```bash
EXCHANGE_TYPE=bitflyer  # or coinbase, binance, etc.
```

## Benefits

1. **Extensibility**: Easy to add new exchanges without modifying existing code
2. **Testability**: Mock implementations for testing
3. **Maintainability**: Clear separation of concerns
4. **Flexibility**: Switch between exchanges via configuration

## Testing

All exchange clients should implement the `CryptoExchangeClient` interface and can be tested using the `MockBitFlyerClient` pattern:

```go
mockClient := &client.MockBitFlyerClient{
    GetBalanceFunc: func() (float64, error) {
        return 1000000.0, nil
    },
}

service := service.NewOrderService(mockClient, orderRepo)
```

## Future Enhancements

1. **Exchange-specific features**: Add optional interfaces for exchange-specific functionality
2. **Rate limiting**: Implement per-exchange rate limiting
3. **Failover**: Support multiple exchanges with automatic failover
4. **Aggregation**: Aggregate prices from multiple exchanges

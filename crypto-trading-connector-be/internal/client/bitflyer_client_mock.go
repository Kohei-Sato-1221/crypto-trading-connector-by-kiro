package client

import "github.com/crypto-trading-connector/backend/internal/model"

// MockBitFlyerClient is a mock implementation of BitFlyerClient for testing
type MockBitFlyerClient struct {
	GetTickerFunc  func(productCode string) (*model.TickerResponse, error)
	GetBalanceFunc func() (float64, error)
	SendOrderFunc  func(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error)
}

// GetTicker calls the mock function if set, otherwise returns default values
func (m *MockBitFlyerClient) GetTicker(productCode string) (*model.TickerResponse, error) {
	if m.GetTickerFunc != nil {
		return m.GetTickerFunc(productCode)
	}
	return &model.TickerResponse{
		ProductCode: productCode,
		Ltp:         10000000.0,
	}, nil
}

// GetBalance calls the mock function if set, otherwise returns default balance
func (m *MockBitFlyerClient) GetBalance() (float64, error) {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc()
	}
	return 1000000.0, nil // Default: 1,000,000 JPY
}

// SendOrder calls the mock function if set, otherwise returns default response
func (m *MockBitFlyerClient) SendOrder(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
	if m.SendOrderFunc != nil {
		return m.SendOrderFunc(req)
	}
	return &model.BitFlyerOrderResponse{
		ChildOrderAcceptanceID: "TEST_ORDER_123",
	}, nil
}

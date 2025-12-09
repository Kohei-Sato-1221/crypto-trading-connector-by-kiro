package service

import (
	"errors"
	"testing"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/model"
)

// MockOrderRepository is a mock implementation of OrderRepository for testing
type MockOrderRepository struct {
	SaveOrderFunc     func(order *model.BuyOrder) error
	GetOrderByIDFunc  func(orderID string) (*model.BuyOrder, error)
}

func (m *MockOrderRepository) SaveOrder(order *model.BuyOrder) error {
	if m.SaveOrderFunc != nil {
		return m.SaveOrderFunc(order)
	}
	return nil
}

func (m *MockOrderRepository) GetOrderByID(orderID string) (*model.BuyOrder, error) {
	if m.GetOrderByIDFunc != nil {
		return m.GetOrderByIDFunc(orderID)
	}
	return nil, errors.New("not implemented")
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return 2000000.0, nil // 2,000,000 JPY
		},
		SendOrderFunc: func(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
			return &model.BitFlyerOrderResponse{
				ChildOrderAcceptanceID: "ORDER_123",
			}, nil
		},
	}

	mockRepo := &MockOrderRepository{
		SaveOrderFunc: func(order *model.BuyOrder) error {
			return nil
		},
	}

	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0.001,
	}

	order, err := service.CreateOrder(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("expected order, got nil")
	}
	// OrderId is UUID type, just check it's not empty
	if order.OrderId.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected valid order ID, got empty UUID")
	}
	if order.EstimatedTotal != 14000 {
		t.Errorf("expected estimated total 14000, got %f", order.EstimatedTotal)
	}
	if order.Status != generated.Pending {
		t.Errorf("expected status pending, got %s", order.Status)
	}
}

func TestOrderService_CreateOrder_InsufficientBalance(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return 10000.0, nil // Only 10,000 JPY
		},
	}

	mockRepo := &MockOrderRepository{}

	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0.001,
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error for insufficient balance, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_CreateOrder_InvalidPrice(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{}
	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     0, // Invalid price
		Amount:    0.001,
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error for invalid price, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_CreateOrder_InvalidAmount(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{}
	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0, // Invalid amount
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error for invalid amount, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_CreateOrder_AmountBelowMinimum(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{}
	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0.0001, // Below minimum 0.001
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error for amount below minimum, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_CreateOrder_BitFlyerAPIError(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return 2000000.0, nil
		},
		SendOrderFunc: func(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
			return nil, errors.New("bitFlyer API error")
		},
	}

	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0.001,
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error from bitFlyer API, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_CreateOrder_BalanceFetchError(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return 0, errors.New("balance fetch error")
		},
	}

	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	req := &generated.CreateOrderRequest{
		Pair:      generated.CreateOrderRequestPairBTCJPY,
		OrderType: generated.CreateOrderRequestOrderTypeLimit,
		Price:     14000000,
		Amount:    0.001,
	}

	order, err := service.CreateOrder(req)

	if err == nil {
		t.Error("expected error from balance fetch, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}

func TestOrderService_GetBalance_Success(t *testing.T) {
	expectedBalance := 1540200.0
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return expectedBalance, nil
		},
	}

	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	balance, err := service.GetBalance()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if balance == nil {
		t.Fatal("expected balance, got nil")
	}
	if balance.AvailableBalance != expectedBalance {
		t.Errorf("expected balance %f, got %f", expectedBalance, balance.AvailableBalance)
	}
	if balance.Currency != generated.JPY {
		t.Errorf("expected currency JPY, got %s", balance.Currency)
	}
}

func TestOrderService_GetBalance_Error(t *testing.T) {
	mockClient := &client.MockBitFlyerClient{
		GetBalanceFunc: func() (float64, error) {
			return 0, errors.New("balance fetch error")
		},
	}

	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockClient, mockRepo)

	balance, err := service.GetBalance()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if balance != nil {
		t.Errorf("expected nil balance, got %v", balance)
	}
}

func TestOrderService_ValidateOrderRequest(t *testing.T) {
	service := &OrderServiceImpl{}

	tests := []struct {
		name    string
		req     *generated.CreateOrderRequest
		wantErr bool
	}{
		{
			name: "valid BTC order",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     14000000,
				Amount:    0.001,
			},
			wantErr: false,
		},
		{
			name: "valid ETH order",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairETHJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     480000,
				Amount:    0.01,
			},
			wantErr: false,
		},
		{
			name: "zero price",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     0,
				Amount:    0.001,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     -100,
				Amount:    0.001,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     14000000,
				Amount:    0,
			},
			wantErr: true,
		},
		{
			name: "BTC amount below minimum",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     14000000,
				Amount:    0.0001,
			},
			wantErr: true,
		},
		{
			name: "ETH amount below minimum",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairETHJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     480000,
				Amount:    0.001,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateOrderRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOrderRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

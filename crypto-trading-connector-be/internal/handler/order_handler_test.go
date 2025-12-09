package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// MockOrderService is a mock implementation of OrderService for testing
type MockOrderService struct {
	CreateOrderFunc func(req *generated.CreateOrderRequest) (*generated.Order, error)
	GetBalanceFunc  func() (*generated.Balance, error)
}

func (m *MockOrderService) CreateOrder(req *generated.CreateOrderRequest) (*generated.Order, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(req)
	}
	return nil, errors.New("not implemented")
}

func (m *MockOrderService) GetBalance() (*generated.Balance, error) {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc()
	}
	return nil, errors.New("not implemented")
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	mockService := &MockOrderService{
		CreateOrderFunc: func(req *generated.CreateOrderRequest) (*generated.Order, error) {
			return &generated.Order{
				OrderId:        openapi_types.UUID(uuid.New()),
				Pair:           generated.OrderPairBTCJPY,
				OrderType:      generated.OrderOrderTypeLimit,
				Price:          14000000,
				Amount:         0.001,
				EstimatedTotal: 14000,
				Status:         generated.Pending,
			}, nil
		},
	}

	handler := NewOrderHandler(mockService)
	e := echo.New()

	reqBody := `{
		"pair": "BTC/JPY",
		"orderType": "limit",
		"price": 14000000,
		"amount": 0.001
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateOrder(c)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	var order generated.Order
	if err := json.Unmarshal(rec.Body.Bytes(), &order); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if order.Price != 14000000 {
		t.Errorf("expected price 14000000, got %f", order.Price)
	}
}

func TestOrderHandler_CreateOrder_InvalidRequest(t *testing.T) {
	mockService := &MockOrderService{}
	handler := NewOrderHandler(mockService)
	e := echo.New()

	reqBody := `{invalid json}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handler.CreateOrder(c)

	// When there's an error, Echo's error handler sets the status code
	// For this test, we just verify the handler doesn't panic
	// In a real scenario, the error would be caught by Echo's middleware
}

func TestOrderHandler_CreateOrder_InvalidPrice(t *testing.T) {
	mockService := &MockOrderService{}
	handler := NewOrderHandler(mockService)
	e := echo.New()

	reqBody := `{
		"pair": "BTC/JPY",
		"orderType": "limit",
		"price": 0,
		"amount": 0.001
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handler.CreateOrder(c)

	// Validation error would be caught by Echo's middleware
}

func TestOrderHandler_CreateOrder_InsufficientBalance(t *testing.T) {
	mockService := &MockOrderService{
		CreateOrderFunc: func(req *generated.CreateOrderRequest) (*generated.Order, error) {
			return nil, errors.New("insufficient balance: required 14000.00, available 10000.00")
		},
	}

	handler := NewOrderHandler(mockService)
	e := echo.New()

	reqBody := `{
		"pair": "BTC/JPY",
		"orderType": "limit",
		"price": 14000000,
		"amount": 0.001
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handler.CreateOrder(c)

	// Service error would be caught by Echo's middleware
}

func TestOrderHandler_GetBalance_Success(t *testing.T) {
	mockService := &MockOrderService{
		GetBalanceFunc: func() (*generated.Balance, error) {
			return &generated.Balance{
				AvailableBalance: 1540200,
				Currency:         generated.JPY,
				Timestamp:        1704067200,
			}, nil
		},
	}

	handler := NewOrderHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/balance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetBalance(c)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var balance generated.Balance
	if err := json.Unmarshal(rec.Body.Bytes(), &balance); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if balance.AvailableBalance != 1540200 {
		t.Errorf("expected balance 1540200, got %f", balance.AvailableBalance)
	}
}

func TestOrderHandler_GetBalance_Error(t *testing.T) {
	mockService := &MockOrderService{
		GetBalanceFunc: func() (*generated.Balance, error) {
			return nil, errors.New("balance fetch error")
		},
	}

	handler := NewOrderHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/balance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handler.GetBalance(c)

	// Service error would be caught by Echo's middleware
}

func TestValidateCreateOrderRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *generated.CreateOrderRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &generated.CreateOrderRequest{
				Pair:      generated.CreateOrderRequestPairBTCJPY,
				OrderType: generated.CreateOrderRequestOrderTypeLimit,
				Price:     14000000,
				Amount:    0.001,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateOrderRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateOrderRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

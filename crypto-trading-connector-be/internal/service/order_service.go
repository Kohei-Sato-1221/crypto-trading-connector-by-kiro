package service

import (
	"fmt"
	"strings"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/model"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(req *generated.CreateOrderRequest) (*generated.Order, error)
	GetBalance() (*generated.Balance, error)
}

// OrderServiceImpl implements OrderService
type OrderServiceImpl struct {
	exchangeClient client.CryptoExchangeClient
	orderRepo      repository.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(exchangeClient client.CryptoExchangeClient, orderRepo repository.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{
		exchangeClient: exchangeClient,
		orderRepo:      orderRepo,
	}
}

// CreateOrder creates a new order
func (s *OrderServiceImpl) CreateOrder(req *generated.CreateOrderRequest) (*generated.Order, error) {
	// Validate input
	if err := s.validateOrderRequest(req); err != nil {
		return nil, err
	}

	// Get balance from exchange
	balance, err := s.exchangeClient.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Calculate estimated total
	estimatedTotal := req.Price * req.Amount

	// Check if balance is sufficient
	if estimatedTotal > balance {
		return nil, fmt.Errorf("insufficient balance: required %.2f, available %.2f", estimatedTotal, balance)
	}

	// Convert pair format (BTC/JPY -> BTC_JPY)
	productCode := strings.ReplaceAll(string(req.Pair), "/", "_")

	// Send order to exchange
	exchangeReq := &model.BitFlyerOrderRequest{
		ProductCode:    productCode,
		ChildOrderType: "LIMIT",
		Side:           "BUY",
		Price:          req.Price,
		Size:           req.Amount,
		TimeInForce:    "GTC", // Good Till Cancelled
	}

	exchangeResp, err := s.exchangeClient.SendOrder(exchangeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send order to exchange: %w", err)
	}

	// Save order to database
	buyOrder := &model.BuyOrder{
		OrderID:     exchangeResp.ChildOrderAcceptanceID,
		ProductCode: productCode,
		Side:        "BUY",
		Price:       req.Price,
		Size:        req.Amount,
		Exchange:    "bitflyer",
		Filled:      0, // 0: unfilled
		Strategy:    99, // 99: not recorded
		Remarks:     nil,
	}

	if err := s.orderRepo.SaveOrder(buyOrder); err != nil {
		// Log error but don't fail - order was already sent to exchange
		fmt.Printf("Warning: failed to save order to database: %v\n", err)
	}

	// Create response
	// Parse UUID from string
	orderUUID, err := uuid.Parse(exchangeResp.ChildOrderAcceptanceID)
	if err != nil {
		// If not a valid UUID, generate a new one
		// In production, exchanges may return acceptance IDs that are not UUIDs
		orderUUID = uuid.New()
	}

	order := &generated.Order{
		OrderId:        openapi_types.UUID(orderUUID),
		Pair:           generated.OrderPair(req.Pair),
		OrderType:      generated.OrderOrderType(req.OrderType),
		Price:          req.Price,
		Amount:         req.Amount,
		EstimatedTotal: estimatedTotal,
		Status:         generated.Pending,
	}

	return order, nil
}

// GetBalance retrieves the current balance
func (s *OrderServiceImpl) GetBalance() (*generated.Balance, error) {
	balance, err := s.exchangeClient.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &generated.Balance{
		AvailableBalance: balance,
		Currency:         generated.JPY,
		Timestamp:        0, // TODO: Add proper timestamp
	}, nil
}

// validateOrderRequest validates the order request
func (s *OrderServiceImpl) validateOrderRequest(req *generated.CreateOrderRequest) error {
	// Validate price
	if req.Price <= 0 {
		return fmt.Errorf("invalid price: must be greater than 0")
	}

	// Validate amount
	if req.Amount <= 0 {
		return fmt.Errorf("invalid amount: must be greater than 0")
	}

	// Validate minimum amount based on pair
	switch req.Pair {
	case generated.CreateOrderRequestPairBTCJPY:
		if req.Amount < 0.001 {
			return fmt.Errorf("invalid amount for BTC/JPY: minimum is 0.001")
		}
	case generated.CreateOrderRequestPairETHJPY:
		if req.Amount < 0.01 {
			return fmt.Errorf("invalid amount for ETH/JPY: minimum is 0.01")
		}
	default:
		return fmt.Errorf("unsupported pair: %s", req.Pair)
	}

	// Validate order type
	if req.OrderType != generated.CreateOrderRequestOrderTypeLimit {
		return fmt.Errorf("unsupported order type: %s", req.OrderType)
	}

	return nil
}

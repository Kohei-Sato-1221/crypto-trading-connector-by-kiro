package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

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
	GetCurrentOrders(pair string, limit int) (*model.CurrentOrdersResponse, error)
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
		Status:      "UNFILLED", // UNFILLED status
		Strategy:    99,         // 99: not recorded
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

// GetCurrentOrders retrieves current orders from the exchange
func (s *OrderServiceImpl) GetCurrentOrders(pair string, limit int) (*model.CurrentOrdersResponse, error) {
	// Default limit to 10 if not specified or invalid
	if limit <= 0 {
		limit = 10
	}

	// Convert pair format if specified (BTC/JPY -> BTC_JPY)
	var productCode string
	if pair != "" {
		productCode = strings.ReplaceAll(pair, "/", "_")
	}

	// Get active orders from exchange
	orders, err := s.exchangeClient.GetChildOrders(productCode, "ACTIVE")
	if err != nil {
		return nil, fmt.Errorf("failed to get child orders from exchange: %w", err)
	}

	// Convert to CurrentOrder format and separate by type
	var buyOrders []model.CurrentOrder
	var sellOrders []model.CurrentOrder

	for _, order := range orders {
		currentOrder := s.convertToCurrentOrder(&order)

		if order.Side == "BUY" {
			buyOrders = append(buyOrders, currentOrder)
		} else if order.Side == "SELL" {
			sellOrders = append(sellOrders, currentOrder)
		}
	}

	// Sort by creation date descending (newest first)
	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].CreatedAt > buyOrders[j].CreatedAt
	})
	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].CreatedAt > sellOrders[j].CreatedAt
	})

	// Limit results
	if len(buyOrders) > limit {
		buyOrders = buyOrders[:limit]
	}
	if len(sellOrders) > limit {
		sellOrders = sellOrders[:limit]
	}

	// Create response
	response := &model.CurrentOrdersResponse{
		BuyOrders:  buyOrders,
		SellOrders: sellOrders,
		Timestamp:  time.Now().Unix(),
	}

	// Add pair filter if specified
	if pair != "" {
		response.Pair = &pair
	}

	return response, nil
}

// convertToCurrentOrder converts BitflyerChildOrder to CurrentOrder
func (s *OrderServiceImpl) convertToCurrentOrder(order *model.BitflyerChildOrder) model.CurrentOrder {
	// Convert product code format (BTC_JPY -> BTC/JPY)
	pair := strings.ReplaceAll(order.ProductCode, "_", "/")

	// Convert side to lowercase
	orderType := strings.ToLower(order.Side)

	// Parse and format creation date to ISO 8601
	createdAt := order.ChildOrderDate
	if parsedTime, err := time.Parse("2006-01-02T15:04:05.999", order.ChildOrderDate); err == nil {
		createdAt = parsedTime.Format(time.RFC3339)
	}

	return model.CurrentOrder{
		ID:        order.ChildOrderID,
		Type:      orderType,
		Pair:      pair,
		Price:     order.Price,
		Amount:    order.Size,
		CreatedAt: createdAt,
	}
}

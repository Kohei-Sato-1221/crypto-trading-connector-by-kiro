package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/labstack/echo/v4"
)

// OrderHandler handles HTTP requests for order endpoints
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req generated.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return handleError(c, http.StatusBadRequest, generated.BADREQUEST, "Invalid request body")
	}

	// Validate request
	if err := validateCreateOrderRequest(&req); err != nil {
		return handleError(c, http.StatusBadRequest, generated.INVALIDREQUEST, err.Error())
	}

	// Create order
	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		return h.handleOrderError(c, err)
	}

	return c.JSON(http.StatusCreated, order)
}

// GetBalance handles GET /api/v1/balance
func (h *OrderHandler) GetBalance(c echo.Context) error {
	balance, err := h.orderService.GetBalance()
	if err != nil {
		return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, "Failed to get balance")
	}

	return c.JSON(http.StatusOK, balance)
}

// GetCurrentOrders handles GET /api/v1/orders/current
func (h *OrderHandler) GetCurrentOrders(c echo.Context) error {
	// Parse query parameters
	pair := c.QueryParam("pair")
	limitStr := c.QueryParam("limit")

	// Default limit to 10
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Validate pair parameter if provided
	if pair != "" && pair != "BTC_JPY" && pair != "ETH_JPY" {
		return handleError(c, http.StatusBadRequest, generated.INVALIDREQUEST, "Invalid pair parameter. Supported pairs: BTC_JPY, ETH_JPY")
	}

	// Get current orders from service
	orders, err := h.orderService.GetCurrentOrders(pair, limit)
	if err != nil {
		return h.handleCurrentOrdersError(c, err)
	}

	return c.JSON(http.StatusOK, orders)
}

// handleCurrentOrdersError handles errors from GetCurrentOrders
func (h *OrderHandler) handleCurrentOrdersError(c echo.Context, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "authentication") || strings.Contains(errMsg, "unauthorized") {
		return handleError(c, http.StatusUnauthorized, generated.UNAUTHORIZED, "Authentication required")
	}
	if strings.Contains(errMsg, "invalid") {
		return handleError(c, http.StatusBadRequest, generated.INVALIDREQUEST, errMsg)
	}

	// Default to internal server error
	return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, "Failed to retrieve current orders")
}

// validateCreateOrderRequest validates the create order request
func validateCreateOrderRequest(req *generated.CreateOrderRequest) error {
	if req.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "price must be greater than 0")
	}
	if req.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "amount must be greater than 0")
	}
	if req.Pair != generated.CreateOrderRequestPairBTCJPY && req.Pair != generated.CreateOrderRequestPairETHJPY {
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported pair")
	}
	if req.OrderType != generated.CreateOrderRequestOrderTypeLimit {
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported order type")
	}
	return nil
}

// handleOrderError handles errors from order service
func (h *OrderHandler) handleOrderError(c echo.Context, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "insufficient balance") {
		return handleError(c, http.StatusPaymentRequired, generated.INSUFFICIENTBALANCE, errMsg)
	}
	if strings.Contains(errMsg, "invalid price") {
		return handleError(c, http.StatusBadRequest, generated.INVALIDPRICE, errMsg)
	}
	if strings.Contains(errMsg, "invalid amount") {
		return handleError(c, http.StatusBadRequest, generated.INVALIDAMOUNT, errMsg)
	}
	if strings.Contains(errMsg, "unsupported pair") {
		return handleError(c, http.StatusBadRequest, generated.UNSUPPORTEDPAIR, errMsg)
	}

	// Default to internal server error
	return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, "Failed to create order")
}

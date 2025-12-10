package handler

import (
	"net/http"
	"strconv"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/labstack/echo/v4"
)

// TradeHistoryHandler handles HTTP requests for trade history endpoints
type TradeHistoryHandler struct {
	tradeHistoryService service.TradeHistoryService
}

// NewTradeHistoryHandler creates a new trade history handler
func NewTradeHistoryHandler(tradeHistoryService service.TradeHistoryService) *TradeHistoryHandler {
	return &TradeHistoryHandler{
		tradeHistoryService: tradeHistoryService,
	}
}

// GetTradeStatistics handles GET /api/v1/trade-history/statistics
func (h *TradeHistoryHandler) GetTradeStatistics(c echo.Context) error {
	// Get query parameters directly
	assetFilter := c.QueryParam("asset_filter")
	if assetFilter == "" {
		assetFilter = "all"
	}

	timeFilter := c.QueryParam("time_filter")
	if timeFilter == "" {
		timeFilter = "all"
	}

	// Validate filters
	if err := h.validateFilters(assetFilter, timeFilter); err != nil {
		return handleError(c, http.StatusBadRequest, generated.INVALIDFILTER, err.Error())
	}

	// Get statistics from service
	statistics, err := h.tradeHistoryService.GetTradeStatistics(assetFilter, timeFilter)
	if err != nil {
		return h.handleTradeHistoryError(c, err)
	}

	return c.JSON(http.StatusOK, statistics)
}

// GetTradeTransactions handles GET /api/v1/trade-history/transactions
func (h *TradeHistoryHandler) GetTradeTransactions(c echo.Context) error {
	// Get query parameters directly
	assetFilter := c.QueryParam("asset_filter")
	if assetFilter == "" {
		assetFilter = "all"
	}

	timeFilter := c.QueryParam("time_filter")
	if timeFilter == "" {
		timeFilter = "all"
	}

	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Validate filters
	if err := h.validateFilters(assetFilter, timeFilter); err != nil {
		return handleError(c, http.StatusBadRequest, generated.INVALIDFILTER, err.Error())
	}

	// Validate pagination
	if err := h.validatePagination(page, limit); err != nil {
		return handleError(c, http.StatusBadRequest, generated.INVALIDPAGINATION, err.Error())
	}

	// Get transactions from service
	response, err := h.tradeHistoryService.GetTradeTransactions(assetFilter, timeFilter, page, limit)
	if err != nil {
		return h.handleTradeHistoryError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// validateFilters validates asset and time filter parameters
func (h *TradeHistoryHandler) validateFilters(assetFilter, timeFilter string) error {
	// Validate asset filter
	validAssetFilters := map[string]bool{
		"all": true,
		"BTC": true,
		"ETH": true,
	}
	if !validAssetFilters[assetFilter] {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset filter: "+assetFilter+". Valid values are: all, BTC, ETH")
	}

	// Validate time filter
	validTimeFilters := map[string]bool{
		"all":   true,
		"7days": true,
	}
	if !validTimeFilters[timeFilter] {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid time filter: "+timeFilter+". Valid values are: all, 7days")
	}

	return nil
}

// validatePagination validates pagination parameters
func (h *TradeHistoryHandler) validatePagination(page, limit int) error {
	if page < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page: "+strconv.Itoa(page)+". Page must be >= 1")
	}

	if limit < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid limit: "+strconv.Itoa(limit)+". Limit must be >= 1")
	}

	if limit > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid limit: "+strconv.Itoa(limit)+". Limit must be <= 100")
	}

	return nil
}

// handleTradeHistoryError handles service errors and converts them to appropriate HTTP responses
func (h *TradeHistoryHandler) handleTradeHistoryError(c echo.Context, err error) error {
	// Check for specific error types
	errMsg := err.Error()

	// Filter validation errors
	if contains(errMsg, "invalid asset filter") || contains(errMsg, "invalid time filter") {
		return handleError(c, http.StatusBadRequest, generated.INVALIDFILTER, errMsg)
	}

	// Pagination validation errors
	if contains(errMsg, "invalid page") || contains(errMsg, "invalid limit") {
		return handleError(c, http.StatusBadRequest, generated.INVALIDPAGINATION, errMsg)
	}

	// Database or other internal errors
	return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, "Internal server error")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 indexOf(s, substr) >= 0)))
}

// indexOf returns the index of the first occurrence of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}


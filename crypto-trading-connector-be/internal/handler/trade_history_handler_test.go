package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTradeHistoryService is a mock implementation of TradeHistoryService
type MockTradeHistoryService struct {
	mock.Mock
}

func (m *MockTradeHistoryService) GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error) {
	args := m.Called(assetFilter, timeFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.TradeStatistics), args.Error(1)
}

func (m *MockTradeHistoryService) GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error) {
	args := m.Called(assetFilter, timeFilter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.TransactionLogResponse), args.Error(1)
}

/**
 * **Feature: trade-history-page, Property 7: ページネーション時のフィルター維持**
 * 
 * 任意のページネーション操作について、フィルター条件（資産・時間）が
 * 正しく維持されなければならない
 */

func TestTradeHistoryHandler_PaginationFilterMaintenance_Property(t *testing.T) {
	tests := []struct {
		name        string
		assetFilter string
		timeFilter  string
		page        int
		limit       int
	}{
		{
			name:        "BTC filter with pagination page 1",
			assetFilter: "BTC",
			timeFilter:  "all",
			page:        1,
			limit:       10,
		},
		{
			name:        "ETH filter with pagination page 2",
			assetFilter: "ETH",
			timeFilter:  "7days",
			page:        2,
			limit:       5,
		},
		{
			name:        "all assets filter with pagination page 3",
			assetFilter: "all",
			timeFilter:  "all",
			page:        3,
			limit:       20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTradeHistoryService)
			handler := NewTradeHistoryHandler(mockService)

			// Create mock response
			mockResponse := &generated.TransactionLogResponse{
				Transactions: []generated.Transaction{
					{
						Id:             "1",
						Cryptocurrency: generated.Bitcoin,
						Timestamp:      time.Now(),
						Profit:         25000.0,
						OrderType:      generated.Sell,
						OrderId:        "SELL001",
						BuyPrice:       5800000.0,
						SellPrice:      6050000.0,
						Amount:         0.1,
						BuyOrderId:     "BUY001",
					},
				},
				Pagination: generated.Pagination{
					CurrentPage: tt.page,
					TotalPages:  5,
					TotalCount:  50,
					HasNext:     tt.page < 5,
				},
			}

			// Expect service to be called with exact filter parameters
			mockService.On("GetTradeTransactions", tt.assetFilter, tt.timeFilter, tt.page, tt.limit).Return(mockResponse, nil)

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions", nil)
			
			// Set query parameters
			q := req.URL.Query()
			q.Add("asset_filter", tt.assetFilter)
			q.Add("time_filter", tt.timeFilter)
			q.Add("page", strconv.Itoa(tt.page))
			q.Add("limit", strconv.Itoa(tt.limit))
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute handler
			err := handler.GetTradeTransactions(c)
			require.NoError(t, err)

			// Verify response
			assert.Equal(t, http.StatusOK, rec.Code)

			var response generated.TransactionLogResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify that filters were maintained in the service call
			// This is verified by the mock expectations
			assert.Equal(t, tt.page, response.Pagination.CurrentPage, "Current page should match requested page")
			assert.Equal(t, tt.page < 5, response.Pagination.HasNext, "HasNext should be calculated correctly")

			mockService.AssertExpectations(t)
		})
	}
}

func TestTradeHistoryHandler_GetTradeStatistics_Success(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	// Mock response
	mockStats := &generated.TradeStatistics{
		TotalProfit:      125000.0,
		ProfitPercentage: 12.5,
		ExecutionCount:   5,
		Period:           generated.TradeStatisticsPeriodAll,
	}

	mockService.On("GetTradeStatistics", "all", "all").Return(mockStats, nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute handler
	err := handler.GetTradeStatistics(c)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, mockStats.TotalProfit, response.TotalProfit)
	assert.Equal(t, mockStats.ExecutionCount, response.ExecutionCount)

	mockService.AssertExpectations(t)
}

func TestTradeHistoryHandler_GetTradeStatistics_WithFilters(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	// Mock response
	mockStats := &generated.TradeStatistics{
		TotalProfit:      50000.0,
		ProfitPercentage: 8.5,
		ExecutionCount:   2,
		Period:           generated.TradeStatisticsPeriodN7days,
	}

	mockService.On("GetTradeStatistics", "BTC", "7days").Return(mockStats, nil)

	// Create request with filters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics?asset_filter=BTC&time_filter=7days", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute handler
	err := handler.GetTradeStatistics(c)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, mockStats.TotalProfit, response.TotalProfit)
	assert.Equal(t, mockStats.Period, response.Period)

	mockService.AssertExpectations(t)
}

func TestTradeHistoryHandler_GetTradeTransactions_Success(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	// Mock response
	mockResponse := &generated.TransactionLogResponse{
		Transactions: []generated.Transaction{
			{
				Id:             "1",
				Cryptocurrency: generated.Bitcoin,
				Timestamp:      time.Now(),
				Profit:         25000.0,
				OrderType:      generated.Sell,
				OrderId:        "SELL001",
				BuyPrice:       5800000.0,
				SellPrice:      6050000.0,
				Amount:         0.1,
				BuyOrderId:     "BUY001",
			},
		},
		Pagination: generated.Pagination{
			CurrentPage: 1,
			TotalPages:  1,
			TotalCount:  1,
			HasNext:     false,
		},
	}

	mockService.On("GetTradeTransactions", "all", "all", 1, 10).Return(mockResponse, nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute handler
	err := handler.GetTradeTransactions(c)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TransactionLogResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Transactions, 1)
	assert.Equal(t, "1", response.Transactions[0].Id)

	mockService.AssertExpectations(t)
}

func TestTradeHistoryHandler_InvalidFilters(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	tests := []struct {
		name        string
		assetFilter string
		timeFilter  string
		expectError bool
	}{
		{
			name:        "invalid asset filter",
			assetFilter: "INVALID",
			timeFilter:  "all",
			expectError: true,
		},
		{
			name:        "invalid time filter",
			assetFilter: "BTC",
			timeFilter:  "INVALID",
			expectError: true,
		},
		{
			name:        "valid filters",
			assetFilter: "ETH",
			timeFilter:  "7days",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				mockService.On("GetTradeStatistics", tt.assetFilter, tt.timeFilter).Return(&generated.TradeStatistics{}, nil)
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics", nil)
			
			// Set query parameters
			q := req.URL.Query()
			q.Add("asset_filter", tt.assetFilter)
			q.Add("time_filter", tt.timeFilter)
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute handler
			err := handler.GetTradeStatistics(c)

			if tt.expectError {
				require.NoError(t, err) // Handler should not return error, but HTTP error response
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var errorResponse generated.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Equal(t, generated.INVALIDFILTER, errorResponse.Error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestTradeHistoryHandler_InvalidPagination(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	tests := []struct {
		name        string
		page        string
		limit       string
		expectError bool
	}{
		{
			name:        "invalid page (zero)",
			page:        "0",
			limit:       "10",
			expectError: true,
		},
		{
			name:        "invalid limit (zero)",
			page:        "1",
			limit:       "0",
			expectError: true,
		},
		{
			name:        "invalid limit (too large)",
			page:        "1",
			limit:       "101",
			expectError: true,
		},
		{
			name:        "valid pagination",
			page:        "2",
			limit:       "20",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				mockService.On("GetTradeTransactions", "all", "all", 2, 20).Return(&generated.TransactionLogResponse{}, nil)
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions", nil)
			
			// Set query parameters
			q := req.URL.Query()
			q.Add("page", tt.page)
			q.Add("limit", tt.limit)
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute handler
			err := handler.GetTradeTransactions(c)

			if tt.expectError {
				require.NoError(t, err) // Handler should not return error, but HTTP error response
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var errorResponse generated.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Equal(t, generated.INVALIDPAGINATION, errorResponse.Error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestTradeHistoryHandler_ServiceError(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	// Mock service error
	mockService.On("GetTradeStatistics", "all", "all").Return(nil, errors.New("database connection failed"))

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute handler
	err := handler.GetTradeStatistics(c)
	require.NoError(t, err) // Handler should not return error, but HTTP error response

	// Verify error response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errorResponse generated.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, generated.INTERNALSERVERERROR, errorResponse.Error)

	mockService.AssertExpectations(t)
}

func TestTradeHistoryHandler_DefaultValues(t *testing.T) {
	mockService := new(MockTradeHistoryService)
	handler := NewTradeHistoryHandler(mockService)

	// Mock service call with default values
	mockService.On("GetTradeTransactions", "all", "all", 1, 10).Return(&generated.TransactionLogResponse{
		Transactions: []generated.Transaction{},
		Pagination: generated.Pagination{
			CurrentPage: 1,
			TotalPages:  0,
			TotalCount:  0,
			HasNext:     false,
		},
	}, nil)

	// Create request without query parameters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute handler
	err := handler.GetTradeTransactions(c)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify that service was called with default values
	mockService.AssertExpectations(t)
}
package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/handler"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TradeHistoryIntegrationTest tests the complete flow from HTTP handler to database
func TestTradeHistoryIntegration_GetStatistics_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query and result
	rows := sqlmock.NewRows([]string{"execution_count", "total_profit"}).
		AddRow(5, 12500.75)
	
	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(SUM\(\(s\.price - b\.price\) \* s\.size\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED'`).
		WillReturnRows(rows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 12500.8, response.TotalProfit) // Rounded to 1 decimal
	assert.Equal(t, 5, response.ExecutionCount)
	assert.Equal(t, generated.TradeStatisticsPeriodAll, response.Period)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetStatistics_WithFilters(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query with BTC filter and 7days filter
	rows := sqlmock.NewRows([]string{"execution_count", "total_profit"}).
		AddRow(3, 7500.25)
	
	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(SUM\(\(s\.price - b\.price\) \* s\.size\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED' AND s\.product_code = \? AND s\.timestamp >= DATE_SUB\(NOW\(\), INTERVAL 7 DAY\)`).
		WithArgs("BTC_JPY").
		WillReturnRows(rows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo with query parameters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics?asset_filter=BTC&time_filter=7days", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 7500.3, response.TotalProfit) // Rounded to 1 decimal
	assert.Equal(t, 3, response.ExecutionCount)
	assert.Equal(t, generated.TradeStatisticsPeriodN7days, response.Period)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetTransactions_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query for transactions
	timestamp := time.Now()
	transactionRows := sqlmock.NewRows([]string{
		"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "timestamp", "profit",
	}).
		AddRow("tx1", "sell1", "buy1", "BTC_JPY", 5000000.0, 4900000.0, 0.01, timestamp, 1000.0).
		AddRow("tx2", "sell2", "buy2", "ETH_JPY", 300000.0, 290000.0, 0.1, timestamp, 1000.0)

	mock.ExpectQuery(`SELECT s\.id, s\.order_id, b\.order_id as buy_order_id, s\.product_code, s\.price as sell_price, b\.price as buy_price, s\.size, s\.timestamp, \(s\.price - b\.price\) \* s\.size as profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED' ORDER BY s\.timestamp DESC LIMIT \? OFFSET \?`).
		WithArgs(10, 0).
		WillReturnRows(transactionRows)

	// Setup expected database query for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED'`).
		WillReturnRows(countRows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeTransactions(c)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TransactionLogResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Transactions, 2)
	assert.Equal(t, "tx1", response.Transactions[0].Id)
	assert.Equal(t, generated.Bitcoin, response.Transactions[0].Cryptocurrency)
	assert.Equal(t, 1000.0, response.Transactions[0].Profit)
	
	assert.Equal(t, "tx2", response.Transactions[1].Id)
	assert.Equal(t, generated.Ethereum, response.Transactions[1].Cryptocurrency)
	assert.Equal(t, 1000.0, response.Transactions[1].Profit)

	// Verify pagination
	assert.Equal(t, 1, response.Pagination.CurrentPage)
	assert.Equal(t, 1, response.Pagination.TotalPages)
	assert.Equal(t, 2, response.Pagination.TotalCount)
	assert.False(t, response.Pagination.HasNext)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetTransactions_WithPagination(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query for transactions (page 2, limit 5)
	timestamp := time.Now()
	transactionRows := sqlmock.NewRows([]string{
		"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "timestamp", "profit",
	}).
		AddRow("tx6", "sell6", "buy6", "BTC_JPY", 5000000.0, 4900000.0, 0.01, timestamp, 1000.0)

	mock.ExpectQuery(`SELECT s\.id, s\.order_id, b\.order_id as buy_order_id, s\.product_code, s\.price as sell_price, b\.price as buy_price, s\.size, s\.timestamp, \(s\.price - b\.price\) \* s\.size as profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED' ORDER BY s\.timestamp DESC LIMIT \? OFFSET \?`).
		WithArgs(5, 5). // page 2, limit 5 -> offset 5
		WillReturnRows(transactionRows)

	// Setup expected database query for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED'`).
		WillReturnRows(countRows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo with pagination parameters
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/transactions?page=2&limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeTransactions(c)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TransactionLogResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Transactions, 1)
	assert.Equal(t, "tx6", response.Transactions[0].Id)

	// Verify pagination
	assert.Equal(t, 2, response.Pagination.CurrentPage)
	assert.Equal(t, 2, response.Pagination.TotalPages) // 10 total / 5 per page = 2 pages
	assert.Equal(t, 10, response.Pagination.TotalCount)
	assert.False(t, response.Pagination.HasNext) // page 2 of 2

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_InvalidFilters_Error(t *testing.T) {
	// Setup mock database (no expectations since validation should fail before DB call)
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Test cases for invalid filters
	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid asset filter",
			queryParams:    "asset_filter=XRP",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_FILTER",
		},
		{
			name:           "invalid time filter",
			queryParams:    "time_filter=30days",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_FILTER",
		},
		{
			name:           "invalid pagination - zero page",
			queryParams:    "page=0",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_PAGINATION",
		},
		{
			name:           "invalid pagination - large limit",
			queryParams:    "limit=200",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_PAGINATION",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Echo
			e := echo.New()
			
			var endpoint string
			if tc.name == "invalid pagination - zero page" || tc.name == "invalid pagination - large limit" {
				endpoint = "/api/v1/trade-history/transactions?" + tc.queryParams
			} else {
				endpoint = "/api/v1/trade-history/statistics?" + tc.queryParams
			}
			
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			var err error
			if tc.name == "invalid pagination - zero page" || tc.name == "invalid pagination - large limit" {
				err = h.GetTradeTransactions(c)
			} else {
				err = h.GetTradeStatistics(c)
			}

			// Verify error response
			if err != nil {
				// Echo error handling
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tc.expectedStatus, he.Code)
				}
			} else {
				// JSON error response
				assert.Equal(t, tc.expectedStatus, rec.Code)
				
				var errorResponse map[string]interface{}
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Contains(t, fmt.Sprintf("%v", errorResponse["error"]), tc.expectedError)
			}
		})
	}

	// Verify no database calls were made
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_DatabaseError_Handling(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database error
	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(SUM\(\(s\.price - b\.price\) \* s\.size\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND b\.status = 'FILLED'`).
		WillReturnError(sql.ErrConnDone)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade-history/statistics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Verify error handling
	if err != nil {
		// Echo error handling
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(t, http.StatusInternalServerError, he.Code)
		}
	} else {
		// JSON error response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
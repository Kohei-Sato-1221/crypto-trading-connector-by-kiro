package integration

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/handler"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTradeHistoryIntegration_GetStatistics_Success(t *testing.T) {
	// Setup database mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query
	rows := sqlmock.NewRows([]string{"execution_count", "total_profit"}).
		AddRow(5, 12500.75)

	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(ROUND\(SUM\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989\), 2\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED'`).
		WillReturnRows(rows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/trade-history/statistics?asset_filter=all&time_filter=all", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 12500.8, response.TotalProfit)
	assert.Equal(t, 5, response.ExecutionCount)
	assert.Equal(t, generated.TradeStatisticsPeriodAll, response.Period)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetStatistics_WithFilters(t *testing.T) {
	// Setup database mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query with filters
	rows := sqlmock.NewRows([]string{"execution_count", "total_profit"}).
		AddRow(3, 7500.25)

	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(ROUND\(SUM\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989\), 2\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' AND s\.product_code = \? AND s\.updatetime >= DATE_SUB\(NOW\(\), INTERVAL 7 DAY\)`).
		WithArgs("BTC_JPY").
		WillReturnRows(rows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/trade-history/statistics?asset_filter=BTC&time_filter=7days", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TradeStatistics
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 7500.3, response.TotalProfit)
	assert.Equal(t, 3, response.ExecutionCount)
	assert.Equal(t, generated.TradeStatisticsPeriodN7days, response.Period)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetTransactions_Success(t *testing.T) {
	// Setup database mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query for transactions
	timestamp := time.Now()
	transactionRows := sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
		AddRow("tx1", "sell1", "buy1", "BTC_JPY", 6000000.0, 5800000.0, 0.1, timestamp, 2000.0).
		AddRow("tx2", "sell2", "buy2", "ETH_JPY", 300000.0, 290000.0, 0.1, timestamp, 1000.0)

	mock.ExpectQuery(`SELECT s\.id, s\.order_id, b\.order_id as buy_order_id, s\.product_code, s\.price as sell_price, b\.price as buy_price, s\.size, s\.updatetime, ROUND\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989, 2\) as profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' ORDER BY s\.updatetime DESC LIMIT \? OFFSET \?`).
		WithArgs(10, 0).
		WillReturnRows(transactionRows)

	// Setup expected database query for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED'`).
		WillReturnRows(countRows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/trade-history/transactions?asset_filter=all&time_filter=all&page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeTransactions(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response generated.TransactionLogResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Transactions, 2)
	assert.Equal(t, "tx1", response.Transactions[0].Id)
	assert.Equal(t, generated.Bitcoin, response.Transactions[0].Cryptocurrency)
	assert.Equal(t, 2000.0, response.Transactions[0].Profit)

	// Verify pagination
	assert.Equal(t, 1, response.Pagination.CurrentPage)
	assert.Equal(t, 1, response.Pagination.TotalPages)
	assert.Equal(t, 2, response.Pagination.TotalCount)
	assert.False(t, response.Pagination.HasNext)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_GetTransactions_Pagination(t *testing.T) {
	// Setup database mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database query for transactions (page 2)
	timestamp := time.Now()
	transactionRows := sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
		AddRow("tx6", "sell6", "buy6", "BTC_JPY", 5000000.0, 4900000.0, 0.01, timestamp, 1000.0)

	mock.ExpectQuery(`SELECT s\.id, s\.order_id, b\.order_id as buy_order_id, s\.product_code, s\.price as sell_price, b\.price as buy_price, s\.size, s\.updatetime, ROUND\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989, 2\) as profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED' ORDER BY s\.updatetime DESC LIMIT \? OFFSET \?`).
		WithArgs(5, 5). // page 2, limit 5 -> offset 5
		WillReturnRows(transactionRows)

	// Setup expected database query for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED'`).
		WillReturnRows(countRows)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/trade-history/transactions?asset_filter=all&time_filter=all&page=2&limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeTransactions(c)

	// Assert
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
	assert.False(t, response.Pagination.HasNext) // page 2 of 2, so no next

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_InvalidFilters(t *testing.T) {
	// Setup database mock (no expectations since validation should fail before DB call)
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "invalid asset filter",
			url:            "/trade-history/statistics?asset_filter=INVALID&time_filter=all",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid time filter",
			url:            "/trade-history/statistics?asset_filter=all&time_filter=INVALID",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Echo
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err = h.GetTradeStatistics(c)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}

	// Verify all expectations were met (should be none)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryIntegration_DatabaseError(t *testing.T) {
	// Setup database mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected database error
	mock.ExpectQuery(`SELECT COUNT\(\*\) as execution_count, COALESCE\(ROUND\(SUM\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989\), 2\), 0\) as total_profit FROM sell_orders s INNER JOIN buy_orders b ON s\.parentid = b\.order_id WHERE s\.status = 'FILLED'`).
		WillReturnError(sql.ErrConnDone)

	// Setup components
	repo := repository.NewMySQLTradeHistoryRepository(db)
	svc := service.NewTradeHistoryService(repo)
	h := handler.NewTradeHistoryHandler(svc)

	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/trade-history/statistics?asset_filter=all&time_filter=all", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = h.GetTradeStatistics(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

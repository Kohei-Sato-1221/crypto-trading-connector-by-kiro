package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresTradeHistoryRepository_GetTradeStatistics_BTCFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresTradeHistoryRepository(db)

	// With BTC filter, the query includes " AND s.product_code = $1"
	statsQuery := `
		SELECT
			COUNT(*) as execution_count,
			COALESCE(ROUND(CAST(SUM(((s.price * s.size) - (b.price * b.size)) * 0.9989) AS numeric), 2), 0) as total_profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	 AND s.product_code = $1`

	mock.ExpectQuery(regexp.QuoteMeta(statsQuery)).
		WithArgs("BTC_JPY").
		WillReturnRows(sqlmock.NewRows([]string{"execution_count", "total_profit"}).
			AddRow(10, 250000.0))

	result, err := repo.GetTradeStatistics("BTC", "all")

	require.NoError(t, err)
	assert.Equal(t, 10, result.ExecutionCount)
	assert.Equal(t, 250000.0, result.TotalProfit)
	assert.Equal(t, generated.TradeStatisticsPeriodAll, result.Period)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTradeHistoryRepository_GetTradeStatistics_AllFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresTradeHistoryRepository(db)

	// With "all" filter, no product_code condition is added
	statsQuery := `
		SELECT
			COUNT(*) as execution_count,
			COALESCE(ROUND(CAST(SUM(((s.price * s.size) - (b.price * b.size)) * 0.9989) AS numeric), 2), 0) as total_profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	`

	mock.ExpectQuery(regexp.QuoteMeta(statsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"execution_count", "total_profit"}).
			AddRow(25, 500000.0))

	result, err := repo.GetTradeStatistics("all", "all")

	require.NoError(t, err)
	assert.Equal(t, 25, result.ExecutionCount)
	assert.Equal(t, 500000.0, result.TotalProfit)
	assert.Equal(t, generated.TradeStatisticsPeriodAll, result.Period)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTradeHistoryRepository_GetTradeStatistics_7DaysFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresTradeHistoryRepository(db)

	statsQuery := `
		SELECT
			COUNT(*) as execution_count,
			COALESCE(ROUND(CAST(SUM(((s.price * s.size) - (b.price * b.size)) * 0.9989) AS numeric), 2), 0) as total_profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	 AND s.updatetime >= NOW() - INTERVAL '7 days'`

	mock.ExpectQuery(regexp.QuoteMeta(statsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"execution_count", "total_profit"}).
			AddRow(5, 120000.0))

	result, err := repo.GetTradeStatistics("all", "7days")

	require.NoError(t, err)
	assert.Equal(t, 5, result.ExecutionCount)
	assert.Equal(t, 120000.0, result.TotalProfit)
	assert.Equal(t, generated.TradeStatisticsPeriodN7days, result.Period)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTradeHistoryRepository_GetTradeTransactions_BTCFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresTradeHistoryRepository(db)

	// Transaction query with BTC filter
	txQuery := `
		SELECT
			s.id,
			s.order_id,
			b.order_id as buy_order_id,
			s.product_code,
			s.price as sell_price,
			b.price as buy_price,
			s.size,
			s.updatetime,
			ROUND(CAST(((s.price * s.size) - (b.price * b.size)) * 0.9989 AS numeric), 2) as profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	 AND s.product_code = $1 ORDER BY s.updatetime DESC LIMIT $2 OFFSET $3`

	mock.ExpectQuery(regexp.QuoteMeta(txQuery)).
		WithArgs("BTC_JPY", 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
			AddRow("1", "SELL001", "BUY001", "BTC_JPY", 6000000.0, 5800000.0, 0.1, time.Now(), 19977.8))

	// Count query with BTC filter
	countQuery := `
		SELECT COUNT(*)
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	 AND s.product_code = $1`

	mock.ExpectQuery(regexp.QuoteMeta(countQuery)).
		WithArgs("BTC_JPY").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	result, err := repo.GetTradeTransactions("BTC", "all", 1, 10)

	require.NoError(t, err)
	assert.Len(t, result.Transactions, 1)
	assert.Equal(t, generated.Bitcoin, result.Transactions[0].Cryptocurrency)
	assert.Equal(t, 1, result.Pagination.TotalCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTradeHistoryRepository_GetTradeTransactions_AllFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresTradeHistoryRepository(db)

	// Transaction query without asset filter
	txQuery := `
		SELECT
			s.id,
			s.order_id,
			b.order_id as buy_order_id,
			s.product_code,
			s.price as sell_price,
			b.price as buy_price,
			s.size,
			s.updatetime,
			ROUND(CAST(((s.price * s.size) - (b.price * b.size)) * 0.9989 AS numeric), 2) as profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	 ORDER BY s.updatetime DESC LIMIT $1 OFFSET $2`

	mock.ExpectQuery(regexp.QuoteMeta(txQuery)).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
			AddRow("1", "SELL001", "BUY001", "BTC_JPY", 6000000.0, 5800000.0, 0.1, time.Now(), 19977.8).
			AddRow("2", "SELL002", "BUY002", "ETH_JPY", 400000.0, 390000.0, 1.0, time.Now(), 9989.0))

	// Count query without asset filter
	countQuery := `
		SELECT COUNT(*)
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	`

	mock.ExpectQuery(regexp.QuoteMeta(countQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	result, err := repo.GetTradeTransactions("all", "all", 1, 10)

	require.NoError(t, err)
	assert.Len(t, result.Transactions, 2)
	assert.Equal(t, generated.Bitcoin, result.Transactions[0].Cryptocurrency)
	assert.Equal(t, generated.Ethereum, result.Transactions[1].Cryptocurrency)
	assert.Equal(t, 2, result.Pagination.TotalCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

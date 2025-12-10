package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * **Feature: trade-history-page, Property 3: 資産フィルタリングの正確性**
 *
 * 任意の資産フィルター（all/BTC/ETH）について、返される取引データは
 * 指定された資産のみを含まなければならない
 */

func TestTradeHistoryRepository_AssetFiltering_Property(t *testing.T) {
	tests := []struct {
		name        string
		assetFilter string
		productCode string
		shouldMatch bool
	}{
		{
			name:        "BTC filter should return only BTC transactions",
			assetFilter: "BTC",
			productCode: "BTC_JPY",
			shouldMatch: true,
		},
		{
			name:        "BTC filter should not return ETH transactions",
			assetFilter: "BTC",
			productCode: "ETH_JPY",
			shouldMatch: false,
		},
		{
			name:        "ETH filter should return only ETH transactions",
			assetFilter: "ETH",
			productCode: "ETH_JPY",
			shouldMatch: true,
		},
		{
			name:        "ETH filter should not return BTC transactions",
			assetFilter: "ETH",
			productCode: "BTC_JPY",
			shouldMatch: false,
		},
		{
			name:        "all filter should return BTC transactions",
			assetFilter: "all",
			productCode: "BTC_JPY",
			shouldMatch: true,
		},
		{
			name:        "all filter should return ETH transactions",
			assetFilter: "all",
			productCode: "ETH_JPY",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewMySQLTradeHistoryRepository(db)

			// Mock the query based on asset filter
			if tt.assetFilter == "all" {
				// For "all" filter, no WHERE clause for product_code
				mock.ExpectQuery(`SELECT.*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'.*ORDER BY s.updatetime DESC LIMIT.*OFFSET`).
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
						AddRow("1", "SELL001", "BUY001", tt.productCode, 6000000.0, 5800000.0, 0.1, time.Now(), 20000.0))
			} else {
				// For specific asset filter, expect WHERE clause with product_code
				expectedProductCode := getProductCodeFromAsset(tt.assetFilter)
				mock.ExpectQuery(`SELECT.*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED' AND s.product_code = \?.*ORDER BY s.updatetime DESC LIMIT.*OFFSET`).
					WithArgs(expectedProductCode, 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
						AddRow("1", "SELL001", "BUY001", tt.productCode, 6000000.0, 5800000.0, 0.1, time.Now(), 20000.0))
			}

			// Mock total count query
			if tt.assetFilter == "all" {
				mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			} else {
				expectedProductCode := getProductCodeFromAsset(tt.assetFilter)
				mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED' AND s.product_code = \?`).
					WithArgs(expectedProductCode).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			}

			result, err := repo.GetTradeTransactions(tt.assetFilter, "all", 1, 10)

			if tt.shouldMatch {
				require.NoError(t, err)
				assert.Len(t, result.Transactions, 1)
				// Check that the transaction has the expected product code through the cryptocurrency field
				if tt.productCode == "BTC_JPY" {
					assert.Equal(t, generated.Bitcoin, result.Transactions[0].Cryptocurrency)
				} else if tt.productCode == "ETH_JPY" {
					assert.Equal(t, generated.Ethereum, result.Transactions[0].Cryptocurrency)
				}
			} else {
				// For non-matching cases, we expect no results but no error
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/**
 * **Feature: trade-history-page, Property 4: 時間フィルタリングの正確性**
 *
 * 任意の時間フィルター（all/7days）について、返される取引データは
 * 指定された期間内のもののみを含まなければならない
 */

func TestTradeHistoryRepository_TimeFiltering_Property(t *testing.T) {
	tests := []struct {
		name        string
		timeFilter  string
		timestamp   time.Time
		shouldMatch bool
	}{
		{
			name:        "7days filter should return recent transactions",
			timeFilter:  "7days",
			timestamp:   time.Now().AddDate(0, 0, -3), // 3 days ago
			shouldMatch: true,
		},
		{
			name:        "7days filter should not return old transactions",
			timeFilter:  "7days",
			timestamp:   time.Now().AddDate(0, 0, -10), // 10 days ago
			shouldMatch: false,
		},
		{
			name:        "all filter should return recent transactions",
			timeFilter:  "all",
			timestamp:   time.Now().AddDate(0, 0, -3), // 3 days ago
			shouldMatch: true,
		},
		{
			name:        "all filter should return old transactions",
			timeFilter:  "all",
			timestamp:   time.Now().AddDate(0, 0, -10), // 10 days ago
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewMySQLTradeHistoryRepository(db)

			// Mock the query based on time filter
			if tt.timeFilter == "all" {
				// For "all" filter, no WHERE clause for timestamp
				mock.ExpectQuery(`SELECT.*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'.*ORDER BY s.updatetime DESC LIMIT.*OFFSET`).
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
						AddRow("1", "SELL001", "BUY001", "BTC_JPY", 6000000.0, 5800000.0, 0.1, tt.timestamp, 20000.0))
			} else {
				// For "7days" filter, expect WHERE clause with timestamp condition
				mock.ExpectQuery(`SELECT.*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED' AND s.updatetime >= DATE_SUB\(NOW\(\), INTERVAL 7 DAY\).*ORDER BY s.updatetime DESC LIMIT.*OFFSET`).
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
						AddRow("1", "SELL001", "BUY001", "BTC_JPY", 6000000.0, 5800000.0, 0.1, tt.timestamp, 20000.0))
			}

			// Mock total count query
			if tt.timeFilter == "all" {
				mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			} else {
				mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED' AND s.updatetime >= DATE_SUB\(NOW\(\), INTERVAL 7 DAY\)`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			}

			result, err := repo.GetTradeTransactions("all", tt.timeFilter, 1, 10)

			if tt.shouldMatch {
				require.NoError(t, err)
				assert.Len(t, result.Transactions, 1)
			} else {
				// For non-matching cases, we expect no results but no error
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/**
 * **Feature: trade-history-page, Property 5: データ整合性の保証**
 *
 * 任意の取引について、profit = (sell_price * size - buy_price * size) * 0.9989
 * の計算が正しく行われなければならない
 */

func TestTradeHistoryRepository_DataIntegrity_Property(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTradeHistoryRepository(db)

	// Mock the query
	mock.ExpectQuery(`SELECT.*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'.*ORDER BY s.updatetime DESC LIMIT.*OFFSET`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "buy_order_id", "product_code", "sell_price", "buy_price", "size", "updatetime", "profit"}).
			AddRow("1", "SELL001", "BUY001", "BTC_JPY", 6000000.0, 5800000.0, 0.1, time.Now(), 19977.8))

	// Mock total count query
	mock.ExpectQuery(`SELECT COUNT\(\*\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s.status = 'FILLED'`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	result, err := repo.GetTradeTransactions("all", "all", 1, 10)

	require.NoError(t, err)
	assert.Len(t, result.Transactions, 1)

	// Verify profit calculation: (6000000 * 0.1 - 5800000 * 0.1) * 0.9989 = 20000 * 0.9989 = 19977.8
	expectedProfit := 19977.8
	assert.Equal(t, expectedProfit, result.Transactions[0].Profit)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeHistoryRepository_GetTradeStatistics(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTradeHistoryRepository(db)

	// Mock the statistics query
	mock.ExpectQuery(`SELECT.*COUNT\(\*\) as execution_count.*ROUND\(SUM\(\(\(s\.price \* s\.size\) - \(b\.price \* b\.size\)\) \* 0\.9989\), 2\).*FROM sell_orders s.*INNER JOIN buy_orders b.*WHERE s\.status = 'FILLED'`).
		WillReturnRows(sqlmock.NewRows([]string{"execution_count", "total_profit"}).
			AddRow(5, 125000.0))

	result, err := repo.GetTradeStatistics("all", "all")

	require.NoError(t, err)
	assert.Equal(t, 5, result.ExecutionCount)
	assert.Equal(t, 125000.0, result.TotalProfit)
	assert.Equal(t, generated.TradeStatisticsPeriodAll, result.Period)

	assert.NoError(t, mock.ExpectationsWereMet())
}

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
)

// PostgresTradeHistoryRepository implements TradeHistoryRepository with PostgreSQL.
type PostgresTradeHistoryRepository struct {
	db *sql.DB
}

// NewPostgresTradeHistoryRepository creates a new PostgreSQL trade history repository.
func NewPostgresTradeHistoryRepository(db *sql.DB) *PostgresTradeHistoryRepository {
	return &PostgresTradeHistoryRepository{
		db: db,
	}
}

// GetTradeStatistics retrieves aggregated trade statistics from PostgreSQL.
// It calculates total profit and execution count with optional asset and time filters.
func (r *PostgresTradeHistoryRepository) GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error) {
	baseQuery := `
		SELECT
			COUNT(*) as execution_count,
			COALESCE(ROUND(CAST(SUM(((s.price * s.size) - (b.price * b.size)) * 0.9989) AS numeric), 2), 0) as total_profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	`

	args := []any{}
	paramCount := 1

	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		baseQuery += fmt.Sprintf(" AND s.product_code = $%d", paramCount)
		args = append(args, productCode)
		paramCount++
	}

	// Add time filter
	if timeFilter == "7days" {
		baseQuery += " AND s.updatetime >= NOW() - INTERVAL '7 days'"
	}

	var executionCount int
	var totalProfit float64

	err := r.db.QueryRow(baseQuery, args...).Scan(&executionCount, &totalProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade statistics: %w", err)
	}

	// Calculate profit percentage based on average profit per trade
	profitPercentage := 0.0
	if executionCount > 0 && totalProfit != 0 {
		avgProfitPerTrade := totalProfit / float64(executionCount)
		// Calculate percentage based on average trade size (assuming ~100,000 JPY per trade)
		profitPercentage = (avgProfitPerTrade / 100000) * 100
	}

	period := generated.TradeStatisticsPeriodAll
	if timeFilter == "7days" {
		period = generated.TradeStatisticsPeriodN7days
	}

	return &generated.TradeStatistics{
		TotalProfit:      roundToOneDecimal(totalProfit),
		ProfitPercentage: roundToOneDecimal(profitPercentage),
		ExecutionCount:   executionCount,
		Period:           period,
	}, nil
}

// GetTradeTransactions retrieves paginated trade transactions from PostgreSQL.
// It returns a list of transactions with pagination metadata, supporting asset and time filters.
func (r *PostgresTradeHistoryRepository) GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error) {
	offset := (page - 1) * limit

	query := `
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
	`

	args := []any{}
	paramCount := 1

	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		query += fmt.Sprintf(" AND s.product_code = $%d", paramCount)
		args = append(args, productCode)
		paramCount++
	}

	// Add time filter
	if timeFilter == "7days" {
		query += " AND s.updatetime >= NOW() - INTERVAL '7 days'"
	}

	query += fmt.Sprintf(" ORDER BY s.updatetime DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query trade transactions: %w", err)
	}
	defer rows.Close()

	var transactions []generated.Transaction

	for rows.Next() {
		var (
			id          string
			sellOrderID string
			buyOrderID  string
			productCode string
			sellPrice   float64
			buyPrice    float64
			size        float64
			timestamp   time.Time
			profit      float64
		)

		if err := rows.Scan(&id, &sellOrderID, &buyOrderID, &productCode, &sellPrice, &buyPrice, &size, &timestamp, &profit); err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		cryptocurrency := getCryptocurrencyFromProductCode(productCode)

		transaction := generated.Transaction{
			Id:             id,
			Cryptocurrency: cryptocurrency,
			Timestamp:      timestamp,
			Profit:         roundToOneDecimal(profit),
			OrderType:      generated.Sell,
			OrderId:        sellOrderID,
			BuyPrice:       buyPrice,
			SellPrice:      sellPrice,
			Amount:         size,
			BuyOrderId:     buyOrderID,
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	// Get total count for pagination
	totalCount, err := r.GetTotalTransactionCount(assetFilter, timeFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get total transaction count: %w", err)
	}

	totalPages := (totalCount + limit - 1) / limit
	hasNext := page < totalPages

	pagination := generated.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		HasNext:     hasNext,
	}

	return &generated.TransactionLogResponse{
		Transactions: transactions,
		Pagination:   pagination,
	}, nil
}

// GetTotalTransactionCount gets the total count of transactions for pagination from PostgreSQL.
// It supports optional asset and time filters.
func (r *PostgresTradeHistoryRepository) GetTotalTransactionCount(assetFilter, timeFilter string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED'
	`

	args := []any{}
	paramCount := 1

	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		query += fmt.Sprintf(" AND s.product_code = $%d", paramCount)
		args = append(args, productCode)
	}

	// Add time filter
	if timeFilter == "7days" {
		query += " AND s.updatetime >= NOW() - INTERVAL '7 days'"
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}

	return count, nil
}

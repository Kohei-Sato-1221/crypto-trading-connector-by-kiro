package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
)

// TradeHistoryRepository defines the interface for trade history data access
type TradeHistoryRepository interface {
	GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error)
	GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error)
	GetTotalTransactionCount(assetFilter, timeFilter string) (int, error)
}

// MySQLTradeHistoryRepository implements TradeHistoryRepository with MySQL
type MySQLTradeHistoryRepository struct {
	db *sql.DB
}

// NewMySQLTradeHistoryRepository creates a new MySQL trade history repository
func NewMySQLTradeHistoryRepository(db *sql.DB) *MySQLTradeHistoryRepository {
	return &MySQLTradeHistoryRepository{
		db: db,
	}
}

// GetTradeStatistics retrieves aggregated trade statistics
func (r *MySQLTradeHistoryRepository) GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error) {
	query := `
		SELECT 
			COUNT(*) as execution_count,
			COALESCE(SUM((s.price - b.price) * s.size), 0) as total_profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED' AND b.status = 'FILLED'
	`
	
	args := []interface{}{}
	
	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		query += " AND s.product_code = ?"
		args = append(args, productCode)
	}
	
	// Add time filter
	if timeFilter == "7days" {
		query += " AND s.timestamp >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
	}

	var executionCount int
	var totalProfit float64
	
	err := r.db.QueryRow(query, args...).Scan(&executionCount, &totalProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade statistics: %w", err)
	}

	// Calculate profit percentage (mock calculation based on total profit)
	profitPercentage := 0.0
	if totalProfit > 0 {
		profitPercentage = (totalProfit / 10000000) * 100 // Assuming 10M JPY base
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

// GetTradeTransactions retrieves paginated trade transactions
func (r *MySQLTradeHistoryRepository) GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error) {
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
			s.timestamp,
			(s.price - b.price) * s.size as profit
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED' AND b.status = 'FILLED'
	`
	
	args := []interface{}{}
	
	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		query += " AND s.product_code = ?"
		args = append(args, productCode)
	}
	
	// Add time filter
	if timeFilter == "7days" {
		query += " AND s.timestamp >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
	}
	
	query += " ORDER BY s.timestamp DESC LIMIT ? OFFSET ?"
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

// GetTotalTransactionCount gets the total count of transactions for pagination
func (r *MySQLTradeHistoryRepository) GetTotalTransactionCount(assetFilter, timeFilter string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM sell_orders s
		INNER JOIN buy_orders b ON s.parentid = b.order_id
		WHERE s.status = 'FILLED' AND b.status = 'FILLED'
	`
	
	args := []interface{}{}
	
	// Add asset filter
	if assetFilter != "all" {
		productCode := getProductCodeFromAsset(assetFilter)
		query += " AND s.product_code = ?"
		args = append(args, productCode)
	}
	
	// Add time filter
	if timeFilter == "7days" {
		query += " AND s.timestamp >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}

	return count, nil
}

// Helper functions

// getProductCodeFromAsset converts asset filter to product code
func getProductCodeFromAsset(assetFilter string) string {
	switch assetFilter {
	case "BTC":
		return "BTC_JPY"
	case "ETH":
		return "ETH_JPY"
	default:
		return ""
	}
}

// getCryptocurrencyFromProductCode converts product code to cryptocurrency name
func getCryptocurrencyFromProductCode(productCode string) generated.TransactionCryptocurrency {
	if strings.Contains(productCode, "BTC") {
		return generated.Bitcoin
	}
	if strings.Contains(productCode, "ETH") {
		return generated.Ethereum
	}
	return generated.Bitcoin // Default fallback
}

// roundToOneDecimal rounds a float64 to 1 decimal place
func roundToOneDecimal(value float64) float64 {
	if value >= 0 {
		return float64(int(value*10+0.5)) / 10
	}
	return float64(int(value*10-0.5)) / 10
}
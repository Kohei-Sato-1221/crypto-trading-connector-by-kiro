package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
)

// CryptoRepository defines the interface for cryptocurrency data operations
type CryptoRepository interface {
	GetDailyAveragePrices(productCode string, days int) ([]model.ChartDataPoint, error)
}

// MySQLCryptoRepository implements CryptoRepository with MySQL
type MySQLCryptoRepository struct {
	db *sql.DB
}

// NewMySQLCryptoRepository creates a new MySQL repository
func NewMySQLCryptoRepository(db *sql.DB) *MySQLCryptoRepository {
	return &MySQLCryptoRepository{
		db: db,
	}
}

// GetDailyAveragePrices retrieves daily average prices for the specified product and number of days
func (r *MySQLCryptoRepository) GetDailyAveragePrices(productCode string, days int) ([]model.ChartDataPoint, error) {
	query := `
		SELECT 
			DATE(datetime) as date,
			AVG(price) as avg_price
		FROM price_histories
		WHERE product_code = ?
			AND datetime >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(datetime)
		ORDER BY date ASC
	`

	rows, err := r.db.Query(query, productCode, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query price histories: %w", err)
	}
	defer rows.Close()

	var chartData []model.ChartDataPoint
	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	for rows.Next() {
		var date time.Time
		var avgPrice float64

		if err := rows.Scan(&date, &avgPrice); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get day name from date
		dayName := dayNames[date.Weekday()]

		chartData = append(chartData, model.ChartDataPoint{
			Day:   dayName,
			Price: avgPrice,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return chartData, nil
}

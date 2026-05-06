package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
)

// PostgresCryptoRepository implements CryptoRepository with PostgreSQL
type PostgresCryptoRepository struct {
	db *sql.DB
}

// NewPostgresCryptoRepository creates a new PostgreSQL repository
func NewPostgresCryptoRepository(db *sql.DB) *PostgresCryptoRepository {
	return &PostgresCryptoRepository{db: db}
}

// GetDailyAveragePrices retrieves daily average prices for the specified product and number of days
func (r *PostgresCryptoRepository) GetDailyAveragePrices(productCode string, days int) ([]generated.ChartDataPoint, error) {
	query := `
		SELECT
			DATE(datetime) as date,
			AVG(price) as avg_price
		FROM price_histories
		WHERE product_code = $1
			AND datetime >= CURRENT_DATE - CAST($2 AS INTEGER) * INTERVAL '1 day'
		GROUP BY DATE(datetime)
		ORDER BY date ASC
	`

	rows, err := r.db.Query(query, productCode, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query price histories: %w", err)
	}
	defer rows.Close()

	var chartData []generated.ChartDataPoint
	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	for rows.Next() {
		var date time.Time
		var avgPrice float64
		if err := rows.Scan(&date, &avgPrice); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		dayName := dayNames[date.Weekday()]
		chartData = append(chartData, generated.ChartDataPoint{
			Day:   dayName,
			Price: avgPrice,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return chartData, nil
}

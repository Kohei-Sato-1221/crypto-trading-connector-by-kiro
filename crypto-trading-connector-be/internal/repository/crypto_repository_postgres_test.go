package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresCryptoRepository_GetDailyAveragePrices(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresCryptoRepository(db)

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

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("BTC_JPY", 7).
		WillReturnRows(sqlmock.NewRows([]string{"date", "avg_price"}).
			AddRow(yesterday, 5800000.0).
			AddRow(now, 6000000.0))

	result, err := repo.GetDailyAveragePrices("BTC_JPY", 7)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 5800000.0, result[0].Price)
	assert.Equal(t, 6000000.0, result[1].Price)
	// Day names should be set based on the weekday
	assert.NotEmpty(t, result[0].Day)
	assert.NotEmpty(t, result[1].Day)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresCryptoRepository_GetDailyAveragePrices_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresCryptoRepository(db)

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

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("ETH_JPY", 30).
		WillReturnRows(sqlmock.NewRows([]string{"date", "avg_price"}))

	result, err := repo.GetDailyAveragePrices("ETH_JPY", 30)

	require.NoError(t, err)
	assert.Empty(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

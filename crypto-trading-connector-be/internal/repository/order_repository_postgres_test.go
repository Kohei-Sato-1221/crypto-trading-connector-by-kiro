package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crypto-trading-connector/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresOrderRepository_SaveOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	remarks := "test order"
	order := &model.BuyOrder{
		OrderID:     "ORDER001",
		ProductCode: "BTC_JPY",
		Side:        "BUY",
		Price:       5800000.0,
		Size:        0.01,
		Exchange:    "bitflyer",
		Status:      "ACTIVE",
		Strategy:    1,
		Remarks:     &remarks,
	}

	insertQuery := `
		INSERT INTO buy_orders (
			order_id, product_code, side, price, size,
			exchange, status, strategy, remarks, timestamp, updatetime
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(
			order.OrderID,
			order.ProductCode,
			order.Side,
			order.Price,
			order.Size,
			order.Exchange,
			order.Status,
			order.Strategy,
			order.Remarks,
			sqlmock.AnyArg(), // timestamp
			sqlmock.AnyArg(), // updatetime
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveOrder(order)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOrderRepository_SaveOrder_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	order := &model.BuyOrder{
		OrderID:     "ORDER002",
		ProductCode: "ETH_JPY",
		Side:        "BUY",
		Price:       400000.0,
		Size:        0.1,
		Exchange:    "bitflyer",
		Status:      "ACTIVE",
		Strategy:    1,
		Remarks:     nil,
	}

	insertQuery := `
		INSERT INTO buy_orders (
			order_id, product_code, side, price, size,
			exchange, status, strategy, remarks, timestamp, updatetime
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(
			order.OrderID,
			order.ProductCode,
			order.Side,
			order.Price,
			order.Size,
			order.Exchange,
			order.Status,
			order.Strategy,
			order.Remarks,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(assert.AnError)

	err = repo.SaveOrder(order)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save order")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOrderRepository_GetOrderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	selectQuery := `
		SELECT id, order_id, product_code, side, price, size,
		       exchange, status, strategy, remarks, timestamp, updatetime
		FROM buy_orders
		WHERE order_id = $1
	`

	now := time.Now().Format(time.RFC3339)
	remarks := "test"

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
		WithArgs("ORDER001").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "order_id", "product_code", "side", "price", "size",
			"exchange", "status", "strategy", "remarks", "timestamp", "updatetime",
		}).AddRow(
			1, "ORDER001", "BTC_JPY", "BUY", 5800000.0, 0.01,
			"bitflyer", "ACTIVE", 1, &remarks, now, now,
		))

	order, err := repo.GetOrderByID("ORDER001")

	require.NoError(t, err)
	assert.Equal(t, "ORDER001", order.OrderID)
	assert.Equal(t, "BTC_JPY", order.ProductCode)
	assert.Equal(t, "BUY", order.Side)
	assert.Equal(t, 5800000.0, order.Price)
	assert.Equal(t, 0.01, order.Size)
	assert.Equal(t, "bitflyer", order.Exchange)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOrderRepository_GetOrderByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	selectQuery := `
		SELECT id, order_id, product_code, side, price, size,
		       exchange, status, strategy, remarks, timestamp, updatetime
		FROM buy_orders
		WHERE order_id = $1
	`

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
		WithArgs("NONEXISTENT").
		WillReturnError(sql.ErrNoRows)

	order, err := repo.GetOrderByID("NONEXISTENT")

	require.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "order not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

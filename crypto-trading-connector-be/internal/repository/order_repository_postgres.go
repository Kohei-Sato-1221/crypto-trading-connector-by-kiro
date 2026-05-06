package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
)

// PostgresOrderRepository implements OrderRepository using PostgreSQL.
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new PostgreSQL-backed order repository.
func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

// SaveOrder saves a buy order to the PostgreSQL database.
func (r *PostgresOrderRepository) SaveOrder(order *model.BuyOrder) error {
	query := `
		INSERT INTO buy_orders (
			order_id, product_code, side, price, size,
			exchange, status, strategy, remarks, timestamp, updatetime
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		order.OrderID,
		order.ProductCode,
		order.Side,
		order.Price,
		order.Size,
		order.Exchange,
		order.Status,
		order.Strategy,
		order.Remarks,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

// GetOrderByID retrieves an order by its order ID from the PostgreSQL database.
func (r *PostgresOrderRepository) GetOrderByID(orderID string) (*model.BuyOrder, error) {
	query := `
		SELECT id, order_id, product_code, side, price, size,
		       exchange, status, strategy, remarks, timestamp, updatetime
		FROM buy_orders
		WHERE order_id = $1
	`

	var order model.BuyOrder
	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&order.OrderID,
		&order.ProductCode,
		&order.Side,
		&order.Price,
		&order.Size,
		&order.Exchange,
		&order.Status,
		&order.Strategy,
		&order.Remarks,
		&order.Timestamp,
		&order.Updatetime,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/model"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	SaveOrder(order *model.BuyOrder) error
	GetOrderByID(orderID string) (*model.BuyOrder, error)
}

// OrderRepositoryImpl implements OrderRepository
type OrderRepositoryImpl struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{
		db: db,
	}
}

// SaveOrder saves a buy order to the database
func (r *OrderRepositoryImpl) SaveOrder(order *model.BuyOrder) error {
	query := `
		INSERT INTO buy_orders (
			order_id, product_code, side, price, size, 
			exchange, filled, strategy, remarks, timestamp, updatetime
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		order.Filled,
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

// GetOrderByID retrieves an order by its order ID
func (r *OrderRepositoryImpl) GetOrderByID(orderID string) (*model.BuyOrder, error) {
	query := `
		SELECT id, order_id, product_code, side, price, size, 
		       exchange, filled, strategy, remarks, timestamp, updatetime
		FROM buy_orders
		WHERE order_id = ?
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
		&order.Filled,
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

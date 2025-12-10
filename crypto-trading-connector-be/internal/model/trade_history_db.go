package model

import "time"

// SellOrder represents a record from sell_orders table
type SellOrder struct {
	ID          int       `db:"id"`
	ParentID    string    `db:"parentid"`
	OrderID     string    `db:"order_id"`
	ProductCode string    `db:"product_code"`
	Side        string    `db:"side"`
	Price       float64   `db:"price"`
	Size        float64   `db:"size"`
	Exchange    string    `db:"exchange"`
	Status      string    `db:"status"`
	Remarks     *string   `db:"remarks"`
	Timestamp   time.Time `db:"timestamp"`
	Updatetime  time.Time `db:"updatetime"`
}

// TradeHistoryQuery represents the result of joining buy and sell orders for profit calculation
type TradeHistoryQuery struct {
	SellOrderID     string    `db:"sell_order_id"`
	BuyOrderID      string    `db:"buy_order_id"`
	ProductCode     string    `db:"product_code"`
	SellPrice       float64   `db:"sell_price"`
	BuyPrice        float64   `db:"buy_price"`
	Size            float64   `db:"size"`
	SellTimestamp   time.Time `db:"sell_timestamp"`
	BuyTimestamp    time.Time `db:"buy_timestamp"`
	Profit          float64   `db:"profit"`
	Cryptocurrency  string    `db:"cryptocurrency"`
}

// TradeHistoryFilter represents filter parameters for trade history queries
type TradeHistoryFilter struct {
	AssetFilter string `json:"asset_filter" form:"asset_filter"`
	TimeFilter  string `json:"time_filter" form:"time_filter"`
	Page        int    `json:"page" form:"page"`
	Limit       int    `json:"limit" form:"limit"`
}
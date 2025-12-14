package model

// PriceHistory represents a record from price_histories table
// This is a database-specific model not defined in OpenAPI
type PriceHistory struct {
	ID            int
	Datetime      string
	ProductCode   string
	Price         float64
	PriceRatio24h *float64
}

// TickerResponse represents bitFlyer ticker API response
// This is a bitFlyer-specific model not defined in OpenAPI
type TickerResponse struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"` // Last Traded Price
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// BuyOrder represents a record from buy_orders table
type BuyOrder struct {
	ID          int     `db:"id"`
	OrderID     string  `db:"order_id"`
	ProductCode string  `db:"product_code"`
	Side        string  `db:"side"`
	Price       float64 `db:"price"`
	Size        float64 `db:"size"`
	Exchange    string  `db:"exchange"`
	Status      string  `db:"status"`
	Strategy    int     `db:"strategy"`
	Remarks     *string `db:"remarks"`
	Timestamp   string  `db:"timestamp"`
	Updatetime  string  `db:"updatetime"`
}

// BitFlyerBalance represents balance response from bitFlyer API
type BitFlyerBalance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

// BitFlyerOrderRequest represents order request to bitFlyer API
type BitFlyerOrderRequest struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"` // LIMIT or MARKET
	Side           string  `json:"side"`             // BUY or SELL
	Price          float64 `json:"price"`
	Size           float64 `json:"size"`
	TimeInForce    string  `json:"time_in_force,omitempty"` // GTC, IOC, FOK
}

// BitFlyerOrderResponse represents order response from bitFlyer API
type BitFlyerOrderResponse struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

// BitflyerChildOrder represents a child order response from bitFlyer API
type BitflyerChildOrder struct {
	ID                     int64   `json:"id"`
	ChildOrderID           string  `json:"child_order_id"`
	ProductCode            string  `json:"product_code"`
	Side                   string  `json:"side"`
	ChildOrderType         string  `json:"child_order_type"`
	Price                  float64 `json:"price"`
	AveragePrice           float64 `json:"average_price"`
	Size                   float64 `json:"size"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
}

// CurrentOrder represents a current order for API response
type CurrentOrder struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"` // "buy" or "sell"
	Pair      string  `json:"pair"` // "BTC/JPY" or "ETH/JPY"
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"createdAt"` // ISO 8601 format
}

// CurrentOrdersResponse represents the API response for current orders
type CurrentOrdersResponse struct {
	BuyOrders  []CurrentOrder `json:"buyOrders"`
	SellOrders []CurrentOrder `json:"sellOrders"`
	Timestamp  int64          `json:"timestamp"`
	Pair       *string        `json:"pair,omitempty"` // Only present if pair filter was specified
}

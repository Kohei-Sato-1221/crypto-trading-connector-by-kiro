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

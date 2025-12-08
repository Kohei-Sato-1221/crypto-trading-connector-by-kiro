package model

// ChartDataPoint represents a single data point in a price chart
type ChartDataPoint struct {
	Day   string  `json:"day"`
	Price float64 `json:"price"`
}

// CryptoData represents cryptocurrency market data
type CryptoData struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Symbol        string           `json:"symbol"`
	Pair          string           `json:"pair"`
	Icon          string           `json:"icon"`
	IconColor     string           `json:"iconColor"`
	CurrentPrice  float64          `json:"currentPrice"`
	ChangePercent float64          `json:"changePercent"`
	ChartData     []ChartDataPoint `json:"chartData"`
}

// MarketResponse represents the response for market data endpoint
type MarketResponse struct {
	Data      []CryptoData `json:"data"`
	Timestamp int64        `json:"timestamp"`
}

// ChartResponse represents the response for chart data endpoint
type ChartResponse struct {
	Data   []ChartDataPoint `json:"data"`
	Period string           `json:"period"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// PriceHistory represents a record from price_histories table
type PriceHistory struct {
	ID            int
	Datetime      string
	ProductCode   string
	Price         float64
	PriceRatio24h *float64
}

// TickerResponse represents bitFlyer ticker API response
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

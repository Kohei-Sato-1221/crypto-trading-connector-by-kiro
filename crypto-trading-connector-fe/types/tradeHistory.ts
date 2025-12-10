/**
 * Trade statistics interface (API response format)
 */
export interface TradeStatistics {
  total_profit: number      // Total profit in JPY, rounded to 1 decimal place
  profit_percentage: number // Profit percentage, rounded to 1 decimal place
  execution_count: number   // Number of completed sell orders
  period: 'all' | '7days'  // Time period
}

/**
 * Trade statistics interface (UI format)
 */
export interface TradeStatisticsUI {
  totalProfit: number      // Total profit in JPY, rounded to 1 decimal place
  profitPercentage: number // Profit percentage, rounded to 1 decimal place
  executionCount: number   // Number of completed sell orders
  period: string          // Time period ('all' or '7days')
}

/**
 * Individual transaction interface (API response format)
 */
export interface Transaction {
  id: string
  cryptocurrency: 'Bitcoin' | 'Ethereum'
  timestamp: string        // ISO date string
  profit: number           // Profit in JPY, rounded to 1 decimal place
  order_type: 'sell'       // Order type
  order_id: string         // e.g., '#BF-88219'
  buy_price: number        // Buy price in JPY
  sell_price: number       // Sell price in JPY
  amount: number           // Amount of cryptocurrency
  buy_order_id: string     // Corresponding buy order ID
}

/**
 * Individual transaction interface (UI format)
 */
export interface TransactionUI {
  id: string
  cryptocurrency: string    // e.g., 'Bitcoin', 'Ethereum'
  timestamp: Date
  profit: number           // Profit in JPY, rounded to 1 decimal place
  orderType: string        // 'sell'
  orderId: string          // e.g., '#BF-88219'
  buyPrice: number         // Buy price in JPY
  sellPrice: number        // Sell price in JPY
  amount: number           // Amount of cryptocurrency
  buyOrderId: string       // Corresponding buy order ID
}

/**
 * Pagination interface (API response format)
 */
export interface Pagination {
  current_page: number
  total_pages: number
  total_count: number
  has_next: boolean
}

/**
 * Transaction log response interface (API response format)
 */
export interface TransactionLogResponse {
  transactions: Transaction[]
  pagination: Pagination
}

/**
 * Transaction log response interface (UI format)
 */
export interface TransactionLogResponseUI {
  transactions: TransactionUI[]
  hasMore: boolean
  total: number
}

/**
 * Filter state interface
 */
export interface FilterState {
  timeFilter: 'all' | '7days'
  assetFilter: 'all' | 'BTC' | 'ETH'
}

/**
 * API request parameters interface
 */
export interface ApiRequestParams {
  asset_filter?: 'all' | 'BTC' | 'ETH'
  time_filter?: 'all' | '7days'
  page?: number
  limit?: number
}
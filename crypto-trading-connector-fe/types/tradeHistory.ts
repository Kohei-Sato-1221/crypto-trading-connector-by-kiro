/**
 * Trade statistics interface
 */
export interface TradeStatistics {
  totalProfit: number      // Total profit in JPY, rounded to 1 decimal place
  profitPercentage: number // Profit percentage, rounded to 1 decimal place
  executionCount: number   // Number of completed sell orders
  period: string          // Time period ('all' or '7days')
}

/**
 * Individual transaction interface
 */
export interface Transaction {
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
 * Transaction log response interface
 */
export interface TransactionLogResponse {
  transactions: Transaction[]
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
 * Pagination interface
 */
export interface Pagination {
  page: number
  limit: number
}
/**
 * Crypto Trading App Type Definitions
 */

/**
 * Represents a single data point in a price chart
 */
export interface ChartDataPoint {
  day: string
  price: number
}

/**
 * Represents cryptocurrency market data
 */
export interface CryptoData {
  id: string
  name: string
  symbol: string
  pair: string
  icon: string
  iconColor: string
  currentPrice: number
  changePercent: number
  chartData: ChartDataPoint[]
}

/**
 * Time filter options for chart display
 */
export type TimeFilter = '24h' | '7d' | '30d' | '1y' | 'all'

/**
 * Navigation item structure
 */
export interface NavigationItem {
  name: string
  path: string
  icon: string
  label: string
}

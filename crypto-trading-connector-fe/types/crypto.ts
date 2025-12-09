/**
 * Crypto Trading App Type Definitions
 * 
 * Re-exports types from OpenAPI-generated types for convenience
 */

import type { components } from './api'

/**
 * Represents a single data point in a price chart
 * Generated from OpenAPI spec
 */
export type ChartDataPoint = components['schemas']['ChartDataPoint']

/**
 * Represents cryptocurrency market data
 * Generated from OpenAPI spec
 */
export type CryptoData = components['schemas']['CryptoData']

/**
 * Market response with timestamp
 * Generated from OpenAPI spec
 */
export type MarketResponse = components['schemas']['MarketResponse']

/**
 * Chart response with period
 * Generated from OpenAPI spec
 */
export type ChartResponse = components['schemas']['ChartResponse']

/**
 * Error response
 * Generated from OpenAPI spec
 */
export type ErrorResponse = components['schemas']['ErrorResponse']

/**
 * Time filter options for chart display
 * Frontend-specific type
 */
export type TimeFilter = '24h' | '7d' | '30d' | '1y' | 'all'

/**
 * Navigation item structure
 * Frontend-specific type
 */
export interface NavigationItem {
  name: string
  path: string
  icon: string
  label: string
}

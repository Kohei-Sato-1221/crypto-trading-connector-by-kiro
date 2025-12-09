import { computed, type Ref } from 'vue'

/**
 * Composable for price rounding and discount calculations
 * Handles currency-specific rounding rules
 */
export const usePriceRounding = () => {
  /**
   * Round price based on currency pair
   * BTC/JPY: Round to nearest million (1,000,000)
   * ETH/JPY: Round to nearest ten thousand (10,000)
   * 
   * @param price - Current price to round
   * @param pair - Currency pair (BTC/JPY or ETH/JPY)
   * @returns Rounded price
   */
  const roundPrice = (price: number, pair: string): number => {
    if (pair === 'BTC/JPY') {
      // Round to nearest million
      return Math.round(price / 1000000) * 1000000
    } else if (pair === 'ETH/JPY') {
      // Round to nearest ten thousand
      return Math.round(price / 10000) * 10000
    }
    return price
  }

  /**
   * Calculate discounted price
   * 
   * @param price - Original price
   * @param percentage - Discount percentage (e.g., 99 for 99%)
   * @returns Discounted price (floored to integer)
   */
  const calculateDiscountPrice = (price: number, percentage: number): number => {
    return Math.floor(price * (percentage / 100))
  }

  return {
    roundPrice,
    calculateDiscountPrice
  }
}

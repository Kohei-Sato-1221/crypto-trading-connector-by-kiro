import { ref, computed, watch, type Ref } from 'vue'
import { usePriceRounding } from './usePriceRounding'

/**
 * Composable for order form logic
 * Handles price, amount, validation, and calculations
 */
export const useOrderForm = (
  pair: Ref<string>,
  currentPrice: Ref<number>,
  availableBalance: Ref<number>
) => {
  const { roundPrice, calculateDiscountPrice } = usePriceRounding()

  // Form state
  const price = ref<number>(0)
  const amount = ref<number>(0)

  // Computed values
  const estimatedTotal = computed(() => {
    return Math.floor(price.value * amount.value)
  })

  const minAmount = computed(() => {
    return pair.value === 'BTC/JPY' ? 0.001 : 0.01
  })

  const step = computed(() => {
    return pair.value === 'BTC/JPY' ? 0.001 : 0.01
  })

  const isBalanceSufficient = computed(() => {
    return estimatedTotal.value <= availableBalance.value
  })

  const isValidOrder = computed(() => {
    return (
      price.value > 0 &&
      amount.value > 0 &&
      amount.value >= minAmount.value &&
      isBalanceSufficient.value
    )
  })

  /**
   * Initialize form with default values
   */
  const initializeDefaults = () => {
    price.value = roundPrice(currentPrice.value, pair.value)
    amount.value = minAmount.value
  }

  /**
   * Increment amount by step value
   */
  const incrementAmount = () => {
    const newAmount = amount.value + step.value
    // Round to 3 decimal places to avoid floating point issues
    amount.value = parseFloat(newAmount.toFixed(3))
  }

  /**
   * Decrement amount by step value
   * Ensures amount doesn't go below minimum
   */
  const decrementAmount = () => {
    const newAmount = amount.value - step.value
    // Ensure we don't go below minimum
    amount.value = Math.max(minAmount.value, parseFloat(newAmount.toFixed(3)))
  }

  /**
   * Set price to a discounted percentage of current price
   * 
   * @param percentage - Discount percentage (e.g., 99 for 99%)
   */
  const setDiscountPrice = (percentage: number) => {
    price.value = calculateDiscountPrice(currentPrice.value, percentage)
  }

  /**
   * Validate order before submission
   * 
   * @returns Object with validation result and error message
   */
  const validateOrder = (): { valid: boolean; error?: string } => {
    if (price.value <= 0) {
      return { valid: false, error: '有効な価格を入力してください' }
    }
    if (amount.value <= 0) {
      return { valid: false, error: '有効な数量を入力してください' }
    }
    if (amount.value < minAmount.value) {
      return { valid: false, error: `最小数量は${minAmount.value}です` }
    }
    if (!isBalanceSufficient.value) {
      return { valid: false, error: '利用可能残高が不足しています' }
    }
    return { valid: true }
  }

  /**
   * Reset form to default values
   */
  const resetForm = () => {
    initializeDefaults()
  }

  // Watch for pair changes and reset amount to new minimum
  watch(pair, () => {
    // When pair changes, update amount to the new pair's minimum
    amount.value = minAmount.value
    // Also update price if currentPrice is available
    if (currentPrice.value > 0) {
      price.value = roundPrice(currentPrice.value, pair.value)
    }
  })

  // Watch for current price changes and update price
  // This handles both initial load and pair changes (when new price data arrives)
  watch(currentPrice, (newPrice) => {
    if (newPrice > 0) {
      price.value = roundPrice(newPrice, pair.value)
    }
  })

  return {
    // State
    price,
    amount,
    
    // Computed
    estimatedTotal,
    minAmount,
    step,
    isBalanceSufficient,
    isValidOrder,
    
    // Methods
    initializeDefaults,
    incrementAmount,
    decrementAmount,
    setDiscountPrice,
    validateOrder,
    resetForm
  }
}

import { ref, onMounted, onUnmounted, watch, readonly, type Ref } from 'vue'
import type { components } from '~/types/api'

// Current Order interface (matches API schema)
interface CurrentOrder {
  id: string
  type: 'buy' | 'sell'
  pair: 'BTC/JPY' | 'ETH/JPY'
  price: number
  amount: number
  createdAt: string // ISO 8601 format
}

// API response type
type CurrentOrdersResponse = components['schemas']['CurrentOrdersResponse']



export const useCurrentOrders = (selectedPair?: Ref<'BTC/JPY' | 'ETH/JPY'>) => {
  // Reactive state
  const buyOrders = ref<CurrentOrder[]>([])
  const sellOrders = ref<CurrentOrder[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  // Track if this is the initial load
  const isInitialLoad = ref(true)

  // Auto-refresh timer
  let refreshTimer: NodeJS.Timeout | null = null
  
  // Retry state for exponential backoff
  let retryCount = 0
  const maxRetries = 5
  const baseDelay = 1000 // 1 second

  // Get API composable
  const { get } = useApi()

  // Calculate exponential backoff delay
  const getRetryDelay = (attempt: number): number => {
    return Math.min(baseDelay * Math.pow(2, attempt), 30000) // Max 30 seconds
  }



  // Convert display pair format to API format
  const convertToApiPairFormat = (displayPair: 'BTC/JPY' | 'ETH/JPY'): 'BTC_JPY' | 'ETH_JPY' => {
    return displayPair.replace('/', '_') as 'BTC_JPY' | 'ETH_JPY'
  }

  // Fetch current orders from API
  const fetchCurrentOrders = async (showLoading: boolean = false): Promise<void> => {
    try {
      // Only show loading for initial load or when explicitly requested
      if (showLoading || isInitialLoad.value) {
        loading.value = true
      }
      
      // Clear error on successful retry
      if (error.value && retryCount === 0) {
        error.value = null
      }

      // Build API parameters
      const params: Record<string, string | number> = {
        limit: 10 // Always limit to 10 orders per type
      }
      
      // Add pair filter if specified
      if (selectedPair?.value) {
        params.pair = convertToApiPairFormat(selectedPair.value)
      }

      // Make API call
      const response = await get<CurrentOrdersResponse>('/orders/current', params)

      // Transform API response to match our interface
      const transformedBuyOrders: CurrentOrder[] = response.buyOrders.map(order => ({
        id: order.id,
        type: order.type,
        pair: order.pair,
        price: order.price,
        amount: order.amount,
        createdAt: order.createdAt
      }))

      const transformedSellOrders: CurrentOrder[] = response.sellOrders.map(order => ({
        id: order.id,
        type: order.type,
        pair: order.pair,
        price: order.price,
        amount: order.amount,
        createdAt: order.createdAt
      }))

      // Update state
      buyOrders.value = transformedBuyOrders
      sellOrders.value = transformedSellOrders

      // Reset retry count on success
      retryCount = 0

      console.log('Current orders fetched from API:', {
        pair: selectedPair?.value || 'all',
        buyOrders: buyOrders.value.length,
        sellOrders: sellOrders.value.length,
        isInitialLoad: isInitialLoad.value,
        timestamp: response.timestamp
      })
      
      // Mark initial load as complete
      if (isInitialLoad.value) {
        isInitialLoad.value = false
      }
    } catch (err) {
      console.error('Error fetching current orders:', err)
      
      // Handle different types of errors with specific messages
      let errorMessage = 'Failed to fetch current orders'
      let shouldRetry = false
      
      if (err instanceof Error) {
        const message = err.message.toLowerCase()
        
        if (message.includes('network error') || message.includes('fetch')) {
          errorMessage = 'Network connection failed'
          shouldRetry = true
        } else if (message.includes('timeout')) {
          errorMessage = 'Request timed out'
          shouldRetry = true
        } else if (message.includes('401') || message.includes('unauthorized')) {
          errorMessage = 'Authentication failed. Please check your credentials.'
          shouldRetry = false
        } else if (message.includes('403') || message.includes('forbidden')) {
          errorMessage = 'Access denied. Please check your permissions.'
          shouldRetry = false
        } else if (message.includes('500') || message.includes('internal server error')) {
          errorMessage = 'Server error. Please try again later.'
          shouldRetry = true
        } else if (message.includes('400') || message.includes('bad request')) {
          errorMessage = 'Invalid request. Please check your settings.'
          shouldRetry = false
        } else {
          errorMessage = err.message
          // For unknown errors, try once more
          shouldRetry = retryCount === 0
        }
      }
      
      // Check if this error should trigger retry
      if (shouldRetry && retryCount < maxRetries) {
        // Exponential backoff retry
        retryCount++
        const delay = getRetryDelay(retryCount - 1)
        
        console.log(`Retrying in ${delay}ms (attempt ${retryCount}/${maxRetries}): ${errorMessage}`)
        
        setTimeout(() => {
          fetchCurrentOrders(false) // Don't show loading for retries
        }, delay)
        
        // Show retry message to user
        if (retryCount === 1) {
          error.value = 'Connection issues, retrying...'
        } else {
          error.value = `Retrying... (${retryCount}/${maxRetries})`
        }
      } else {
        // Set final error message
        retryCount = 0
        if (retryCount >= maxRetries) {
          error.value = 'Unable to connect after multiple attempts. Please check your connection and try again.'
        } else {
          error.value = errorMessage
        }
      }
    } finally {
      // Only hide loading if we're not going to retry
      if (retryCount === 0 || retryCount >= maxRetries) {
        loading.value = false
      }
    }
  }

  // Start auto-refresh with 10-second interval
  const startAutoRefresh = (): void => {
    // Clear existing timer if any
    if (refreshTimer) {
      clearInterval(refreshTimer)
    }

    // Set up 10-second interval
    refreshTimer = setInterval(() => {
      // Only refresh if document is visible (performance optimization)
      if (!document.hidden) {
        // Background update - don't show loading
        fetchCurrentOrders(false)
      }
    }, 10000) // 10 seconds

    console.log('Auto-refresh started (10 second interval)')
  }

  // Stop auto-refresh
  const stopAutoRefresh = (): void => {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
      console.log('Auto-refresh stopped')
    }
  }

  // Handle visibility change (pause/resume on tab switch)
  const handleVisibilityChange = (): void => {
    if (document.hidden) {
      console.log('Tab hidden - pausing auto-refresh')
      stopAutoRefresh()
    } else {
      console.log('Tab visible - resuming auto-refresh')
      startAutoRefresh()
      // Fetch immediately when returning to tab (background update)
      fetchCurrentOrders(false)
    }
  }

  // Store reference to the handler for cleanup
  let visibilityHandler: (() => void) | null = null

  // Watch for pair changes
  if (selectedPair) {
    watch(selectedPair, async () => {
      // Reset initial load flag for pair changes to show loading
      isInitialLoad.value = true
      await fetchCurrentOrders(true)
    })
  }

  // Lifecycle management
  onMounted(async () => {
    // Initial fetch with loading
    await fetchCurrentOrders(true)
    
    // Start auto-refresh
    startAutoRefresh()
    
    // Store handler reference and listen for visibility changes
    visibilityHandler = handleVisibilityChange
    document.addEventListener('visibilitychange', visibilityHandler)
  })

  onUnmounted(() => {
    // Clean up timer and event listeners
    stopAutoRefresh()
    if (visibilityHandler) {
      document.removeEventListener('visibilitychange', visibilityHandler)
      visibilityHandler = null
    }
  })

  // Manual retry function that resets retry count
  const retryFetch = async (): Promise<void> => {
    retryCount = 0
    error.value = null
    await fetchCurrentOrders(true)
  }

  return {
    buyOrders: readonly(buyOrders),
    sellOrders: readonly(sellOrders),
    loading: readonly(loading),
    error: readonly(error),
    fetchCurrentOrders: retryFetch, // Export the retry function instead
    startAutoRefresh,
    stopAutoRefresh
  }
}
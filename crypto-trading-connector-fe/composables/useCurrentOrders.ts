import { ref, onMounted, onUnmounted, watch, type Ref } from 'vue'

// Current Order interface
interface CurrentOrder {
  id: string
  type: 'buy' | 'sell'
  pair: 'BTC/JPY' | 'ETH/JPY'
  price: number
  amount: number
  createdAt: string // ISO 8601 format
}

// Mock data for development - separate data for each pair
const generateMockOrders = (selectedPair: 'BTC/JPY' | 'ETH/JPY'): { buyOrders: CurrentOrder[], sellOrders: CurrentOrder[] } => {
  const now = new Date()
  
  // BTC/JPY Mock Orders
  const btcBuyOrders: CurrentOrder[] = [
    {
      id: 'btc-buy-1',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 14000000,
      amount: 0.001,
      createdAt: new Date(now.getTime() - 2 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-2',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13950000,
      amount: 0.002,
      createdAt: new Date(now.getTime() - 5 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-3',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13900000,
      amount: 0.0015,
      createdAt: new Date(now.getTime() - 8 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-4',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13850000,
      amount: 0.003,
      createdAt: new Date(now.getTime() - 12 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-5',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13800000,
      amount: 0.0025,
      createdAt: new Date(now.getTime() - 18 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-6',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13750000,
      amount: 0.004,
      createdAt: new Date(now.getTime() - 25 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-7',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13700000,
      amount: 0.0035,
      createdAt: new Date(now.getTime() - 32 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-8',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13650000,
      amount: 0.0028,
      createdAt: new Date(now.getTime() - 40 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-9',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13600000,
      amount: 0.0045,
      createdAt: new Date(now.getTime() - 48 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-buy-10',
      type: 'buy',
      pair: 'BTC/JPY',
      price: 13550000,
      amount: 0.005,
      createdAt: new Date(now.getTime() - 55 * 60 * 1000).toISOString()
    }
  ]

  const btcSellOrders: CurrentOrder[] = [
    {
      id: 'btc-sell-1',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14100000,
      amount: 0.001,
      createdAt: new Date(now.getTime() - 3 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-2',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14150000,
      amount: 0.0018,
      createdAt: new Date(now.getTime() - 7 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-3',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14200000,
      amount: 0.0022,
      createdAt: new Date(now.getTime() - 10 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-4',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14250000,
      amount: 0.0035,
      createdAt: new Date(now.getTime() - 15 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-5',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14300000,
      amount: 0.0028,
      createdAt: new Date(now.getTime() - 22 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-6',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14350000,
      amount: 0.004,
      createdAt: new Date(now.getTime() - 28 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-7',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14400000,
      amount: 0.0032,
      createdAt: new Date(now.getTime() - 35 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-8',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14450000,
      amount: 0.0038,
      createdAt: new Date(now.getTime() - 42 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-9',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14500000,
      amount: 0.0025,
      createdAt: new Date(now.getTime() - 50 * 60 * 1000).toISOString()
    },
    {
      id: 'btc-sell-10',
      type: 'sell',
      pair: 'BTC/JPY',
      price: 14550000,
      amount: 0.0045,
      createdAt: new Date(now.getTime() - 58 * 60 * 1000).toISOString()
    }
  ]

  // ETH/JPY Mock Orders
  const ethBuyOrders: CurrentOrder[] = [
    {
      id: 'eth-buy-1',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 480000,
      amount: 0.05,
      createdAt: new Date(now.getTime() - 2 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-2',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 475000,
      amount: 0.08,
      createdAt: new Date(now.getTime() - 5 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-3',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 470000,
      amount: 0.12,
      createdAt: new Date(now.getTime() - 8 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-4',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 465000,
      amount: 0.15,
      createdAt: new Date(now.getTime() - 12 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-5',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 460000,
      amount: 0.18,
      createdAt: new Date(now.getTime() - 18 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-6',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 455000,
      amount: 0.22,
      createdAt: new Date(now.getTime() - 25 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-7',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 450000,
      amount: 0.25,
      createdAt: new Date(now.getTime() - 32 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-8',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 445000,
      amount: 0.28,
      createdAt: new Date(now.getTime() - 40 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-9',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 440000,
      amount: 0.32,
      createdAt: new Date(now.getTime() - 48 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-buy-10',
      type: 'buy',
      pair: 'ETH/JPY',
      price: 435000,
      amount: 0.35,
      createdAt: new Date(now.getTime() - 55 * 60 * 1000).toISOString()
    }
  ]

  const ethSellOrders: CurrentOrder[] = [
    {
      id: 'eth-sell-1',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 490000,
      amount: 0.03,
      createdAt: new Date(now.getTime() - 3 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-2',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 495000,
      amount: 0.06,
      createdAt: new Date(now.getTime() - 7 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-3',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 500000,
      amount: 0.045,
      createdAt: new Date(now.getTime() - 10 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-4',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 505000,
      amount: 0.075,
      createdAt: new Date(now.getTime() - 15 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-5',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 510000,
      amount: 0.09,
      createdAt: new Date(now.getTime() - 22 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-6',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 515000,
      amount: 0.12,
      createdAt: new Date(now.getTime() - 28 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-7',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 520000,
      amount: 0.085,
      createdAt: new Date(now.getTime() - 35 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-8',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 525000,
      amount: 0.11,
      createdAt: new Date(now.getTime() - 42 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-9',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 530000,
      amount: 0.095,
      createdAt: new Date(now.getTime() - 50 * 60 * 1000).toISOString()
    },
    {
      id: 'eth-sell-10',
      type: 'sell',
      pair: 'ETH/JPY',
      price: 535000,
      amount: 0.13,
      createdAt: new Date(now.getTime() - 58 * 60 * 1000).toISOString()
    }
  ]

  // Select orders based on the selected pair
  const buyOrders = selectedPair === 'BTC/JPY' ? btcBuyOrders : ethBuyOrders
  const sellOrders = selectedPair === 'BTC/JPY' ? btcSellOrders : ethSellOrders

  // Sort by creation date descending (newest first)
  buyOrders.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
  sellOrders.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

  // Limit to 10 orders each
  return {
    buyOrders: buyOrders.slice(0, 10),
    sellOrders: sellOrders.slice(0, 10)
  }
}

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

  // Fetch current orders (mock implementation)
  const fetchCurrentOrders = async (showLoading: boolean = false): Promise<void> => {
    try {
      // Only show loading for initial load or when explicitly requested
      if (showLoading || isInitialLoad.value) {
        loading.value = true
      }
      error.value = null

      // Simulate API delay (shorter for background updates)
      const delay = isInitialLoad.value ? 500 : 200
      await new Promise(resolve => setTimeout(resolve, delay))

      // Generate mock data based on selected pair
      const currentPair = selectedPair?.value || 'BTC/JPY'
      const mockData = generateMockOrders(currentPair)
      buyOrders.value = mockData.buyOrders
      sellOrders.value = mockData.sellOrders

      console.log('Current orders fetched:', {
        pair: currentPair,
        buyOrders: buyOrders.value.length,
        sellOrders: sellOrders.value.length,
        isInitialLoad: isInitialLoad.value
      })
      
      // Mark initial load as complete
      if (isInitialLoad.value) {
        isInitialLoad.value = false
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch current orders'
      console.error('Error fetching current orders:', err)
    } finally {
      loading.value = false
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

  return {
    buyOrders: readonly(buyOrders),
    sellOrders: readonly(sellOrders),
    loading: readonly(loading),
    error: readonly(error),
    fetchCurrentOrders,
    startAutoRefresh,
    stopAutoRefresh
  }
}
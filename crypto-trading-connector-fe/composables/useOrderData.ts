import { ref, type Ref } from 'vue'
import type { CryptoData, ChartDataPoint } from '~/types/crypto'

/**
 * Mock data for order page
 */
const getMockOrderData = (pair: string) => {
  const isBTC = pair === 'BTC/JPY'
  
  return {
    currentPrice: isBTC ? 14062621 : 485318,
    priceChange: isBTC ? 2.5 : -1.2,
    chartData: generateMockChartData(isBTC ? 14000000 : 480000, 7),
    availableBalance: 1540200
  }
}

/**
 * Generate mock chart data
 */
const generateMockChartData = (basePrice: number, days: number): ChartDataPoint[] => {
  const data: ChartDataPoint[] = []
  const now = new Date()
  
  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(now)
    date.setDate(date.getDate() - i)
    
    // Add some random variation
    const variation = (Math.random() - 0.5) * 0.1 // ±5% variation
    const price = basePrice * (1 + variation)
    
    data.push({
      timestamp: date.toISOString(),
      price: Math.floor(price)
    })
  }
  
  return data
}

/**
 * Composable for fetching order page data
 * Currently uses mock data, will be replaced with API calls in integration phase
 */
export const useOrderData = (pair: Ref<string>) => {
  const currentPrice = ref<number>(0)
  const priceChange = ref<number>(0)
  const chartData = ref<ChartDataPoint[]>([])
  const availableBalance = ref<number>(0)
  const loading = ref<boolean>(false)
  const error = ref<Error | null>(null)

  /**
   * Fetch price data for selected pair
   */
  const fetchPriceData = async () => {
    loading.value = true
    error.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300))
      
      const mockData = getMockOrderData(pair.value)
      currentPrice.value = mockData.currentPrice
      priceChange.value = mockData.priceChange
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch price data:', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch chart data for selected time filter
   * 
   * @param timeFilter - Time range filter (1H, 24H, 7D, 30D, 1Y)
   */
  const fetchChartData = async (timeFilter: string) => {
    loading.value = true
    error.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300))
      
      // Convert time filter to days
      const days = timeFilterToDays(timeFilter)
      const mockData = getMockOrderData(pair.value)
      const basePrice = pair.value === 'BTC/JPY' ? 14000000 : 480000
      
      chartData.value = generateMockChartData(basePrice, days)
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch chart data:', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch available balance
   */
  const fetchBalance = async () => {
    loading.value = true
    error.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 200))
      
      const mockData = getMockOrderData(pair.value)
      availableBalance.value = mockData.availableBalance
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch balance:', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch all data (price, chart, balance)
   */
  const fetchAllData = async (timeFilter: string = '7D') => {
    await Promise.all([
      fetchPriceData(),
      fetchChartData(timeFilter),
      fetchBalance()
    ])
  }

  /**
   * Convert time filter to number of days
   */
  const timeFilterToDays = (filter: string): number => {
    switch (filter) {
      case '1H':
        return 1
      case '24H':
        return 1
      case '7D':
        return 7
      case '30D':
        return 30
      case '1Y':
        return 365
      default:
        return 7
    }
  }

  return {
    // State
    currentPrice,
    priceChange,
    chartData,
    availableBalance,
    loading,
    error,
    
    // Methods
    fetchPriceData,
    fetchChartData,
    fetchBalance,
    fetchAllData
  }
}

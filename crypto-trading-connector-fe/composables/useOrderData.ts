import { ref, type Ref } from 'vue'
import type { ChartDataPoint } from '~/types/crypto'
import type { components } from '~/types/api'
import { getErrorMessage } from '~/utils/errorHandler'

type Balance = components['schemas']['Balance']
type CreateOrderRequest = components['schemas']['CreateOrderRequest']
type Order = components['schemas']['Order']
type ErrorResponse = components['schemas']['ErrorResponse']

/**
 * Composable for fetching order page data from API
 */
export const useOrderData = (pair: Ref<string>) => {
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl || 'http://localhost:8080/api/v1'
  
  const currentPrice = ref<number>(0)
  const priceChange = ref<number>(0)
  const chartData = ref<ChartDataPoint[]>([])
  const availableBalance = ref<number>(0)
  const loading = ref<boolean>(false)
  const error = ref<Error | null>(null)

  /**
   * Convert pair format from BTC/JPY to bitcoin
   */
  const pairToId = (pairValue: string): string => {
    if (pairValue === 'BTC/JPY') return 'bitcoin'
    if (pairValue === 'ETH/JPY') return 'ethereum'
    return 'bitcoin'
  }

  /**
   * Fetch price data for selected pair from API
   */
  const fetchPriceData = async () => {
    loading.value = true
    error.value = null

    try {
      const cryptoId = pairToId(pair.value)
      const response = await $fetch(`${apiBaseUrl}/crypto/${cryptoId}`)
      
      if (response && typeof response === 'object' && 'currentPrice' in response) {
        currentPrice.value = (response as any).currentPrice
        priceChange.value = (response as any).changePercent || 0
      }
    } catch (e: any) {
      const errorMessage = getErrorMessage(e)
      error.value = new Error(errorMessage)
      console.error('Failed to fetch price data:', errorMessage, e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch chart data for selected time filter from API
   * 
   * @param timeFilter - Time range filter (1H, 24H, 7D, 30D, 1Y)
   */
  const fetchChartData = async (timeFilter: string) => {
    loading.value = true
    error.value = null

    try {
      const cryptoId = pairToId(pair.value)
      const period = timeFilterToPeriod(timeFilter)
      const response = await $fetch(`${apiBaseUrl}/crypto/${cryptoId}/chart?period=${period}`)
      
      if (response && typeof response === 'object' && 'data' in response) {
        const data = (response as any).data
        if (Array.isArray(data)) {
          chartData.value = data.map((item: any) => ({
            timestamp: item.day || '',
            price: item.price || 0
          }))
        }
      }
    } catch (e: any) {
      const errorMessage = getErrorMessage(e)
      error.value = new Error(errorMessage)
      console.error('Failed to fetch chart data:', errorMessage, e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch available balance from API
   */
  const fetchBalance = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await $fetch<Balance>(`${apiBaseUrl}/balance`)
      
      if (response && 'availableBalance' in response) {
        availableBalance.value = response.availableBalance
      }
    } catch (e: any) {
      const errorMessage = getErrorMessage(e)
      error.value = new Error(errorMessage)
      console.error('Failed to fetch balance:', errorMessage, e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Submit order to API
   */
  const submitOrder = async (orderData: CreateOrderRequest): Promise<Order | null> => {
    loading.value = true
    error.value = null

    try {
      const response = await $fetch<Order>(`${apiBaseUrl}/orders`, {
        method: 'POST',
        body: orderData
      })
      
      return response
    } catch (e: any) {
      const errorMessage = getErrorMessage(e)
      error.value = new Error(errorMessage)
      console.error('Failed to submit order:', errorMessage, e)
      
      // Create a more detailed error for the UI
      const enhancedError = new Error(errorMessage)
      ;(enhancedError as any).data = e.data
      throw enhancedError
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
   * Convert time filter to API period format
   */
  const timeFilterToPeriod = (filter: string): string => {
    switch (filter) {
      case '1H':
        return '24h'
      case '24H':
        return '24h'
      case '7D':
        return '7d'
      case '30D':
        return '30d'
      case '1Y':
        return '1y'
      default:
        return '7d'
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
    fetchAllData,
    submitOrder
  }
}

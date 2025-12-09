import { ref, type Ref } from 'vue'
import type { ChartDataPoint } from '~/types/crypto'
import type { components } from '~/types/api'
import { getErrorMessage } from '~/utils/errorHandler'

type Balance = components['schemas']['Balance']
type CreateOrderRequest = components['schemas']['CreateOrderRequest']
type Order = components['schemas']['Order']

/**
 * Composable for fetching order page data from API
 * Uses the same endpoint as market page for consistency (DRY principle)
 */
export const useOrderData = (pair: Ref<string>) => {
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl || 'http://localhost:8080'
  
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

  /**
   * Fetch crypto data (price + chart) for selected pair from API
   * Note: The /crypto/:id endpoint returns both price and chart data,
   * so we don't need separate calls. This matches how market page works.
   * 
   * @param timeFilter - Time range filter (1H, 24H, 7D, 30D, 1Y)
   */
  const fetchCryptoData = async (timeFilter: string = '7D') => {
    loading.value = true
    error.value = null

    try {
      const cryptoId = pairToId(pair.value)
      const period = timeFilterToPeriod(timeFilter)
      const response = await $fetch(`${apiBaseUrl}/api/v1/crypto/${cryptoId}?period=${period}`)
      
      if (response && typeof response === 'object') {
        const data = response as any
        
        // Extract price data
        if ('currentPrice' in data) {
          currentPrice.value = data.currentPrice
        }
        if ('changePercent' in data) {
          priceChange.value = data.changePercent || 0
        }
        
        // Extract chart data (already in correct format)
        if ('chartData' in data && Array.isArray(data.chartData)) {
          chartData.value = data.chartData
        }
      }
    } catch (e: any) {
      const errorMessage = getErrorMessage(e)
      error.value = new Error(errorMessage)
      console.error('Failed to fetch crypto data:', errorMessage, e)
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
      const response = await $fetch<Balance>(`${apiBaseUrl}/api/v1/balance`)
      
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
      const response = await $fetch<Order>(`${apiBaseUrl}/api/v1/orders`, {
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
   * Fetch all data (crypto data + balance)
   * 
   * @param timeFilter - Time range filter (1H, 24H, 7D, 30D, 1Y)
   */
  const fetchAllData = async (timeFilter: string = '7D') => {
    await Promise.all([
      fetchCryptoData(timeFilter),
      fetchBalance()
    ])
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
    fetchCryptoData,
    fetchBalance,
    fetchAllData,
    submitOrder
  }
}

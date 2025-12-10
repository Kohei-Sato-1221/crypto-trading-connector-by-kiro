import { ref } from 'vue'
import type { CryptoData } from '~/types/crypto'
import { getMockCryptoData } from '~/utils/mockData'

export interface CryptoDataWithError extends CryptoData {
  hasError?: boolean
  errorMessage?: string
}

/**
 * Composable for fetching cryptocurrency data
 * Supports both mock data and API integration
 */
export const useCryptoData = () => {
  const cryptoData = ref<CryptoDataWithError[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const config = useRuntimeConfig()
  const { get } = useApi()

  // Use mock data if USE_MOCK_DATA is true, otherwise use API
  const useMockData = config.public.useMockData ?? false

  // Crypto IDs to fetch
  const cryptoIds = ['bitcoin', 'ethereum']

  /**
   * Fetch single cryptocurrency data from API
   */
  const fetchSingleCrypto = async (id: string): Promise<CryptoDataWithError> => {
    try {
      const response = await get<CryptoData>(`/crypto/${id}`)
      return { ...response, hasError: false }
    } catch (e) {
      console.error(`Failed to fetch ${id} data from API:`, e)
      // Return error placeholder
      return {
        id,
        name: id.charAt(0).toUpperCase() + id.slice(1),
        symbol: id === 'bitcoin' ? 'BTC' : 'ETH',
        pair: id === 'bitcoin' ? 'BTC/JPY' : 'ETH/JPY',
        icon: id === 'bitcoin' ? '₿' : 'Ξ',
        iconColor: id === 'bitcoin' ? '#f7931a' : '#627eea',
        currentPrice: 0,
        changePercent: 0,
        chartData: [],
        hasError: true,
        errorMessage: 'Failed to load data'
      }
    }
  }

  /**
   * Fetch cryptocurrency data from API (individual requests)
   */
  const fetchFromAPI = async (): Promise<CryptoDataWithError[]> => {
    const promises = cryptoIds.map(id => fetchSingleCrypto(id))
    return await Promise.all(promises)
  }

  /**
   * Fetch cryptocurrency data
   * Uses mock data or API based on configuration
   */
  const fetchCryptoData = async (): Promise<CryptoDataWithError[]> => {
    loading.value = true
    error.value = null

    try {
      let data: CryptoDataWithError[]
      
      if (useMockData) {
        // Use mock data
        data = getMockCryptoData().map(crypto => ({ ...crypto, hasError: false }))
      } else {
        // Use API (fetch individually)
        data = await fetchFromAPI()
      }
      
      cryptoData.value = data
      return data
    } catch (e) {
      error.value = e as Error
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Refresh cryptocurrency data
   */
  const refresh = async () => {
    return await fetchCryptoData()
  }

  return {
    cryptoData,
    loading,
    error,
    fetchCryptoData,
    refresh,
    useMockData
  }
}

import { ref } from 'vue'
import type { CryptoData } from '~/types/crypto'
import { getMockCryptoData } from '~/utils/mockData'

/**
 * Composable for fetching cryptocurrency data
 * Designed to support both mock data and future API integration
 */
export const useCryptoData = () => {
  const cryptoData = ref<CryptoData[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)

  /**
   * Fetch cryptocurrency data
   * Currently uses mock data, but designed for future API integration
   */
  const fetchCryptoData = async (): Promise<CryptoData[]> => {
    loading.value = true
    error.value = null

    try {
      // TODO: Replace with API call when backend is ready
      // const response = await $fetch('/api/crypto/market')
      // return response.data
      
      // For now, use mock data
      const data = getMockCryptoData()
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
    refresh
  }
}

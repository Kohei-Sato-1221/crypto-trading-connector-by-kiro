import { ref, computed } from 'vue'
import type { TradeStatisticsUI, TransactionUI, FilterState } from '~/types/tradeHistory'
import { getMockTradeStatistics, getFilteredTransactions, calculateFilteredStatistics } from '~/utils/tradeHistoryMockData'

/**
 * Composable for managing trade history data
 */
export const useTradeHistory = () => {
  const config = useRuntimeConfig()
  const useMockData = config.public.useMockData
  const { fetchTradeStatistics, fetchTradeTransactions } = useTradeHistoryApi()

  // State
  const statistics = ref<TradeStatisticsUI>(getMockTradeStatistics())
  const transactions = ref<TransactionUI[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  
  // Individual component error states
  const statisticsError = ref<Error | null>(null)
  const transactionsError = ref<Error | null>(null)
  
  const hasMore = ref(true)
  const currentPage = ref(1)
  const pageSize = 10
  const totalTransactions = ref(0)

  // Filter state
  const filters = ref<FilterState>({
    timeFilter: 'all',
    assetFilter: 'all'
  })

  // Computed
  const canLoadMore = computed(() => {
    if (useMockData) {
      const allFiltered = getFilteredTransactions(filters.value.assetFilter, filters.value.timeFilter)
      return transactions.value.length < allFiltered.length
    }
    return hasMore.value
  })

  /**
   * Update statistics based on current filters
   */
  const updateStatistics = async () => {
    if (useMockData) {
      statistics.value = calculateFilteredStatistics(
        filters.value.assetFilter,
        filters.value.timeFilter
      )
      statisticsError.value = null
      return
    }

    try {
      loading.value = true
      statisticsError.value = null
      
      const stats = await fetchTradeStatistics(
        filters.value.assetFilter,
        filters.value.timeFilter
      )
      statistics.value = stats
    } catch (e) {
      statisticsError.value = e as Error
      console.error('Failed to fetch statistics:', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Load transactions based on current filters and pagination
   */
  const loadTransactions = async (reset: boolean = false) => {
    if (useMockData) {
      const allFiltered = getFilteredTransactions(filters.value.assetFilter, filters.value.timeFilter)
      if (reset) {
        transactions.value = allFiltered.slice(0, pageSize)
        currentPage.value = 1
      } else {
        const startIndex = currentPage.value * pageSize
        const endIndex = startIndex + pageSize
        const newTransactions = allFiltered.slice(startIndex, endIndex)
        transactions.value = [...transactions.value, ...newTransactions]
        currentPage.value++
      }
      return
    }

    try {
      loading.value = true
      transactionsError.value = null

      const page = reset ? 1 : currentPage.value + 1
      
      const response = await fetchTradeTransactions(
        filters.value.assetFilter,
        filters.value.timeFilter,
        page,
        pageSize
      )

      if (reset) {
        transactions.value = response.transactions
        currentPage.value = 1
      } else {
        transactions.value = [...transactions.value, ...response.transactions]
        currentPage.value = page
      }

      hasMore.value = response.hasMore
      totalTransactions.value = response.total

    } catch (e) {
      transactionsError.value = e as Error
      console.error('Failed to fetch transactions:', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Set time filter
   */
  const setTimeFilter = async (timeFilter: 'all' | '7days') => {
    filters.value.timeFilter = timeFilter
    await Promise.all([
      updateStatistics(),
      loadTransactions(true) // Reset transactions
    ])
  }

  /**
   * Set asset filter
   */
  const setAssetFilter = async (assetFilter: 'all' | 'BTC' | 'ETH') => {
    filters.value.assetFilter = assetFilter
    await Promise.all([
      updateStatistics(),
      loadTransactions(true) // Reset transactions
    ])
  }

  /**
   * Load more transactions (pagination)
   */
  const loadMoreTransactions = async () => {
    if (canLoadMore.value && !loading.value) {
      await loadTransactions(false) // Append transactions
    }
  }

  /**
   * Refresh all data
   */
  const refresh = async () => {
    await Promise.all([
      updateStatistics(),
      loadTransactions(true) // Reset transactions
    ])
  }

  /**
   * Initialize data
   */
  const initialize = async () => {
    await refresh()
  }

  return {
    // State
    statistics: readonly(statistics),
    transactions: readonly(transactions),
    loading: readonly(loading),
    error: readonly(error),
    
    // Individual component errors
    statisticsError: readonly(statisticsError),
    transactionsError: readonly(transactionsError),
    
    filters: readonly(filters),
    canLoadMore,
    totalTransactions: readonly(totalTransactions),

    // Actions
    setTimeFilter,
    setAssetFilter,
    loadMoreTransactions,
    refresh,
    initialize
  }
}
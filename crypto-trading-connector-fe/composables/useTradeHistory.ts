import { ref, computed } from 'vue'
import type { TradeStatistics, Transaction, FilterState, TransactionLogResponse } from '~/types/tradeHistory'
import { getMockTradeStatistics, getFilteredTransactions, calculateFilteredStatistics } from '~/utils/tradeHistoryMockData'

/**
 * Composable for managing trade history data
 */
export const useTradeHistory = () => {
  // State
  const statistics = ref<TradeStatistics>(getMockTradeStatistics())
  const transactions = ref<Transaction[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const hasMore = ref(true)
  const currentPage = ref(1)
  const pageSize = 10

  // Filter state
  const filters = ref<FilterState>({
    timeFilter: 'all',
    assetFilter: 'all'
  })

  // Computed
  const filteredTransactions = computed(() => {
    return getFilteredTransactions(filters.value.assetFilter, filters.value.timeFilter)
  })

  const displayedTransactions = computed(() => {
    const allFiltered = filteredTransactions.value
    return allFiltered.slice(0, currentPage.value * pageSize)
  })

  const canLoadMore = computed(() => {
    return displayedTransactions.value.length < filteredTransactions.value.length
  })

  /**
   * Update statistics based on current filters
   */
  const updateStatistics = () => {
    statistics.value = calculateFilteredStatistics(
      filters.value.assetFilter,
      filters.value.timeFilter
    )
  }

  /**
   * Set time filter
   */
  const setTimeFilter = (timeFilter: 'all' | '7days') => {
    filters.value.timeFilter = timeFilter
    currentPage.value = 1
    updateStatistics()
  }

  /**
   * Set asset filter
   */
  const setAssetFilter = (assetFilter: 'all' | 'BTC' | 'ETH') => {
    filters.value.assetFilter = assetFilter
    currentPage.value = 1
    updateStatistics()
  }

  /**
   * Load more transactions (pagination)
   */
  const loadMoreTransactions = () => {
    if (canLoadMore.value) {
      currentPage.value++
    }
  }

  /**
   * Refresh all data
   */
  const refresh = async () => {
    loading.value = true
    error.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500))
      
      // Reset pagination
      currentPage.value = 1
      
      // Update statistics
      updateStatistics()
      
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
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
    transactions: displayedTransactions,
    loading: readonly(loading),
    error: readonly(error),
    filters: readonly(filters),
    canLoadMore,

    // Actions
    setTimeFilter,
    setAssetFilter,
    loadMoreTransactions,
    refresh,
    initialize
  }
}
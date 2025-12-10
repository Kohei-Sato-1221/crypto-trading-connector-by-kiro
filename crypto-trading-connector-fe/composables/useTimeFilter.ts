import { ref } from 'vue'
import type { TimeFilter } from '~/types/crypto'

// Global state for time filter (shared across components)
const globalSelectedFilter = ref<TimeFilter>('7d')

/**
 * Composable for managing time filter state
 * Used for chart display period selection
 * Uses global state to ensure consistency across components
 */
export const useTimeFilter = () => {
  /**
   * Set the selected time filter
   */
  const setFilter = (filter: TimeFilter) => {
    globalSelectedFilter.value = filter
  }

  /**
   * Check if a filter is currently selected
   */
  const isSelected = (filter: TimeFilter): boolean => {
    return globalSelectedFilter.value === filter
  }

  /**
   * Reset filter to default value (for testing)
   */
  const resetFilter = () => {
    globalSelectedFilter.value = '7d'
  }

  return {
    selectedFilter: globalSelectedFilter,
    setFilter,
    isSelected,
    resetFilter
  }
}

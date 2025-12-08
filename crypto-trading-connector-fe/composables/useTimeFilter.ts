import { ref } from 'vue'
import type { TimeFilter } from '~/types/crypto'

/**
 * Composable for managing time filter state
 * Used for chart display period selection
 */
export const useTimeFilter = () => {
  const selectedFilter = ref<TimeFilter>('24h')

  /**
   * Set the selected time filter
   */
  const setFilter = (filter: TimeFilter) => {
    selectedFilter.value = filter
  }

  /**
   * Check if a filter is currently selected
   */
  const isSelected = (filter: TimeFilter): boolean => {
    return selectedFilter.value === filter
  }

  return {
    selectedFilter,
    setFilter,
    isSelected
  }
}

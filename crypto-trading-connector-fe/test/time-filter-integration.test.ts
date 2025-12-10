import { describe, it, expect, beforeEach } from 'vitest'
import { useTimeFilter } from '../composables/useTimeFilter'

describe('Time Filter Integration Tests', () => {
  beforeEach(() => {
    // Reset filter to default before each test
    const { resetFilter } = useTimeFilter()
    resetFilter()
  })

  it('should share state between multiple instances', () => {
    const instance1 = useTimeFilter()
    const instance2 = useTimeFilter()

    // Both instances should start with the same default value
    expect(instance1.selectedFilter.value).toBe('7d')
    expect(instance2.selectedFilter.value).toBe('7d')

    // Change filter in instance1
    instance1.setFilter('30d')

    // Both instances should reflect the change
    expect(instance1.selectedFilter.value).toBe('30d')
    expect(instance2.selectedFilter.value).toBe('30d')

    // Change filter in instance2
    instance2.setFilter('1y')

    // Both instances should reflect the change
    expect(instance1.selectedFilter.value).toBe('1y')
    expect(instance2.selectedFilter.value).toBe('1y')
  })

  it('should correctly identify selected filter across instances', () => {
    const instance1 = useTimeFilter()
    const instance2 = useTimeFilter()

    // Set filter in instance1
    instance1.setFilter('30d')

    // Both instances should identify the same filter as selected
    expect(instance1.isSelected('30d')).toBe(true)
    expect(instance2.isSelected('30d')).toBe(true)
    expect(instance1.isSelected('7d')).toBe(false)
    expect(instance2.isSelected('7d')).toBe(false)
  })

  it('should reset to default value', () => {
    const instance = useTimeFilter()

    // Change to a different filter
    instance.setFilter('1y')
    expect(instance.selectedFilter.value).toBe('1y')

    // Reset to default
    instance.resetFilter()
    expect(instance.selectedFilter.value).toBe('7d')
    expect(instance.isSelected('7d')).toBe(true)
  })

  it('should handle all valid time filter values', () => {
    const instance = useTimeFilter()
    const validFilters = ['7d', '30d', '1y', 'all'] as const

    for (const filter of validFilters) {
      instance.setFilter(filter)
      expect(instance.selectedFilter.value).toBe(filter)
      expect(instance.isSelected(filter)).toBe(true)
      
      // Other filters should not be selected
      for (const otherFilter of validFilters) {
        if (otherFilter !== filter) {
          expect(instance.isSelected(otherFilter)).toBe(false)
        }
      }
    }
  })
})
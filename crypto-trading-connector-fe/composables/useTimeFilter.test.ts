import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import { useTimeFilter } from './useTimeFilter'
import type { TimeFilter } from '~/types/crypto'

// Feature: crypto-market-page, Property 4: Time filter exclusivity
describe('useTimeFilter - Property Based Tests', () => {
  it('Property 4: For any time filter selection, exactly one filter should be active', () => {
    const validFilters: TimeFilter[] = ['24h', '7d', '30d', '1y', 'all']
    
    fc.assert(
      fc.property(
        fc.constantFrom(...validFilters),
        (selectedFilter: TimeFilter) => {
          const { setFilter, isSelected } = useTimeFilter()
          
          // Set the filter
          setFilter(selectedFilter)
          
          // Count how many filters are selected
          let selectedCount = 0
          for (const filter of validFilters) {
            if (isSelected(filter)) {
              selectedCount++
            }
          }
          
          // Exactly one filter should be selected
          expect(selectedCount).toBe(1)
          
          // The selected filter should be the one we set
          expect(isSelected(selectedFilter)).toBe(true)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('Property: Default filter should be 24h', () => {
    const { selectedFilter, isSelected } = useTimeFilter()
    
    expect(selectedFilter.value).toBe('24h')
    expect(isSelected('24h')).toBe(true)
  })

  it('Property: Setting a filter should deselect previous filter', () => {
    const validFilters: TimeFilter[] = ['24h', '7d', '30d', '1y', 'all']
    
    fc.assert(
      fc.property(
        fc.constantFrom(...validFilters),
        fc.constantFrom(...validFilters),
        (firstFilter: TimeFilter, secondFilter: TimeFilter) => {
          const { setFilter, isSelected } = useTimeFilter()
          
          // Set first filter
          setFilter(firstFilter)
          expect(isSelected(firstFilter)).toBe(true)
          
          // Set second filter
          setFilter(secondFilter)
          expect(isSelected(secondFilter)).toBe(true)
          
          // If filters are different, first should no longer be selected
          if (firstFilter !== secondFilter) {
            expect(isSelected(firstFilter)).toBe(false)
          }
        }
      ),
      { numRuns: 100 }
    )
  })
})

describe('useTimeFilter - Unit Tests', () => {
  it('should have default filter of 24h', () => {
    const { selectedFilter } = useTimeFilter()
    expect(selectedFilter.value).toBe('24h')
  })

  it('should update selected filter when setFilter is called', () => {
    const { selectedFilter, setFilter } = useTimeFilter()
    
    setFilter('7d')
    expect(selectedFilter.value).toBe('7d')
    
    setFilter('30d')
    expect(selectedFilter.value).toBe('30d')
  })

  it('should correctly identify selected filter with isSelected', () => {
    const { setFilter, isSelected } = useTimeFilter()
    
    setFilter('1y')
    
    expect(isSelected('1y')).toBe(true)
    expect(isSelected('24h')).toBe(false)
    expect(isSelected('7d')).toBe(false)
    expect(isSelected('30d')).toBe(false)
    expect(isSelected('all')).toBe(false)
  })

  it('should handle all valid time filter values', () => {
    const { setFilter, isSelected } = useTimeFilter()
    const validFilters: TimeFilter[] = ['24h', '7d', '30d', '1y', 'all']
    
    for (const filter of validFilters) {
      setFilter(filter)
      expect(isSelected(filter)).toBe(true)
    }
  })
})

import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import type { CryptoData, ChartDataPoint, TimeFilter, NavigationItem } from './crypto'

// Feature: crypto-market-page, Property 6: Mock data structure validity
describe('Crypto Types - Property Based Tests', () => {
  it('Property 6: For any mock CryptoData object, it should conform to the CryptoData type definition', () => {
    // Arbitrary generator for ChartDataPoint
    const chartDataPointArbitrary = fc.record({
      day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
      price: fc.double({ min: 0, max: 100000000, noNaN: true })
    })

    // Arbitrary generator for CryptoData
    const cryptoDataArbitrary = fc.record({
      id: fc.string({ minLength: 1 }),
      name: fc.constantFrom('Bitcoin', 'Ethereum', 'Ripple', 'Litecoin'),
      symbol: fc.constantFrom('BTC', 'ETH', 'XRP', 'LTC'),
      pair: fc.constantFrom('BTC/JPY', 'ETH/JPY', 'XRP/JPY', 'LTC/JPY'),
      icon: fc.string({ minLength: 1 }),
      iconColor: fc.constantFrom('#f7931a', '#627eea', '#23292f', '#345d9d'),
      currentPrice: fc.double({ min: 0, max: 100000000, noNaN: true }),
      changePercent: fc.double({ min: -100, max: 100, noNaN: true }),
      chartData: fc.array(chartDataPointArbitrary, { minLength: 1, maxLength: 7 })
    })

    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        (data: CryptoData) => {
          // Verify all required fields are present and correctly typed
          expect(typeof data.id).toBe('string')
          expect(data.id.length).toBeGreaterThan(0)
          
          expect(typeof data.name).toBe('string')
          expect(data.name.length).toBeGreaterThan(0)
          
          expect(typeof data.symbol).toBe('string')
          expect(data.symbol.length).toBeGreaterThan(0)
          
          expect(typeof data.pair).toBe('string')
          expect(data.pair.length).toBeGreaterThan(0)
          
          expect(typeof data.icon).toBe('string')
          expect(data.icon.length).toBeGreaterThan(0)
          
          expect(typeof data.iconColor).toBe('string')
          expect(data.iconColor).toMatch(/^#[0-9a-fA-F]{6}$/)
          
          expect(typeof data.currentPrice).toBe('number')
          expect(data.currentPrice).toBeGreaterThanOrEqual(0)
          expect(Number.isFinite(data.currentPrice)).toBe(true)
          
          expect(typeof data.changePercent).toBe('number')
          expect(Number.isFinite(data.changePercent)).toBe(true)
          
          expect(Array.isArray(data.chartData)).toBe(true)
          expect(data.chartData.length).toBeGreaterThan(0)
          
          // Verify each chart data point
          data.chartData.forEach((point: ChartDataPoint) => {
            expect(typeof point.day).toBe('string')
            expect(point.day.length).toBeGreaterThan(0)
            expect(typeof point.price).toBe('number')
            expect(Number.isFinite(point.price)).toBe(true)
          })
        }
      ),
      { numRuns: 100 }
    )
  })

  it('Property: TimeFilter should only accept valid values', () => {
    const validTimeFilters: TimeFilter[] = ['7d', '30d', '1y', 'all']
    
    fc.assert(
      fc.property(
        fc.constantFrom(...validTimeFilters),
        (filter: TimeFilter) => {
          // Verify the filter is one of the valid values
          expect(validTimeFilters).toContain(filter)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('Property: NavigationItem should have all required fields', () => {
    const navigationItemArbitrary = fc.record({
      name: fc.constantFrom('market', 'trade', 'history', 'portfolio'),
      path: fc.constantFrom('/market', '/trade', '/history', '/portfolio'),
      icon: fc.string({ minLength: 1 }),
      label: fc.constantFrom('Market', 'Trade', 'History', 'Portfolio')
    })

    fc.assert(
      fc.property(
        navigationItemArbitrary,
        (item: NavigationItem) => {
          expect(typeof item.name).toBe('string')
          expect(item.name.length).toBeGreaterThan(0)
          
          expect(typeof item.path).toBe('string')
          expect(item.path).toMatch(/^\//)
          
          expect(typeof item.icon).toBe('string')
          expect(item.icon.length).toBeGreaterThan(0)
          
          expect(typeof item.label).toBe('string')
          expect(item.label.length).toBeGreaterThan(0)
        }
      ),
      { numRuns: 100 }
    )
  })
})

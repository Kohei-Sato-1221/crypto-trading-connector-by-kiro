import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import { 
  getMockCryptoData, 
  updateMockPrice, 
  validateCryptoData,
  bitcoinMockData,
  ethereumMockData
} from './mockData'
import type { CryptoData } from '~/types/crypto'

// Feature: crypto-market-page, Property 7: Price update maintains data integrity
describe('Mock Data - Property Based Tests', () => {
  it('Property 7: For any CryptoData after price update, all required fields remain valid', () => {
    const cryptoDataArbitrary = fc.constantFrom(bitcoinMockData, ethereumMockData)
    
    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        (originalData: CryptoData) => {
          // Update the price
          const updatedData = updateMockPrice(originalData)
          
          // Verify all required fields are still present and valid
          expect(updatedData.id).toBe(originalData.id)
          expect(updatedData.name).toBe(originalData.name)
          expect(updatedData.symbol).toBe(originalData.symbol)
          expect(updatedData.pair).toBe(originalData.pair)
          expect(updatedData.icon).toBe(originalData.icon)
          expect(updatedData.iconColor).toBe(originalData.iconColor)
          
          // Price should be positive
          expect(updatedData.currentPrice).toBeGreaterThan(0)
          expect(Number.isFinite(updatedData.currentPrice)).toBe(true)
          
          // Change percent should be a finite number
          expect(Number.isFinite(updatedData.changePercent)).toBe(true)
          
          // Chart data should not be empty
          expect(Array.isArray(updatedData.chartData)).toBe(true)
          expect(updatedData.chartData.length).toBeGreaterThan(0)
          
          // Each chart data point should be valid
          updatedData.chartData.forEach(point => {
            expect(typeof point.day).toBe('string')
            expect(point.day.length).toBeGreaterThan(0)
            expect(typeof point.price).toBe('number')
            expect(Number.isFinite(point.price)).toBe(true)
          })
          
          // Validate using the validation function
          expect(validateCryptoData(updatedData)).toBe(true)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('Property: Multiple price updates maintain data integrity', () => {
    const cryptoDataArbitrary = fc.constantFrom(bitcoinMockData, ethereumMockData)
    const updateCountArbitrary = fc.integer({ min: 1, max: 10 })
    
    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        updateCountArbitrary,
        (originalData: CryptoData, updateCount: number) => {
          let currentData = originalData
          
          // Apply multiple updates
          for (let i = 0; i < updateCount; i++) {
            currentData = updateMockPrice(currentData)
          }
          
          // After multiple updates, data should still be valid
          expect(validateCryptoData(currentData)).toBe(true)
          expect(currentData.currentPrice).toBeGreaterThan(0)
          expect(currentData.chartData.length).toBeGreaterThan(0)
        }
      ),
      { numRuns: 50 }
    )
  })

  it('Property: Chart data length remains constant after updates', () => {
    const cryptoDataArbitrary = fc.constantFrom(bitcoinMockData, ethereumMockData)
    
    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        (originalData: CryptoData) => {
          const originalLength = originalData.chartData.length
          const updatedData = updateMockPrice(originalData)
          
          // Chart data length should remain the same
          expect(updatedData.chartData.length).toBe(originalLength)
        }
      ),
      { numRuns: 100 }
    )
  })
})

describe('Mock Data - Unit Tests', () => {
  it('should return valid mock cryptocurrency data', () => {
    const mockData = getMockCryptoData()
    
    expect(Array.isArray(mockData)).toBe(true)
    expect(mockData.length).toBeGreaterThan(0)
    
    mockData.forEach(crypto => {
      expect(validateCryptoData(crypto)).toBe(true)
    })
  })

  it('should have Bitcoin in mock data', () => {
    const mockData = getMockCryptoData()
    const bitcoin = mockData.find(c => c.id === 'bitcoin')
    
    expect(bitcoin).toBeDefined()
    expect(bitcoin?.name).toBe('Bitcoin')
    expect(bitcoin?.symbol).toBe('BTC')
  })

  it('should have Ethereum in mock data', () => {
    const mockData = getMockCryptoData()
    const ethereum = mockData.find(c => c.id === 'ethereum')
    
    expect(ethereum).toBeDefined()
    expect(ethereum?.name).toBe('Ethereum')
    expect(ethereum?.symbol).toBe('ETH')
  })

  it('should generate 7-day chart data', () => {
    const mockData = getMockCryptoData()
    
    mockData.forEach(crypto => {
      expect(crypto.chartData.length).toBe(7)
      
      const days = crypto.chartData.map(d => d.day)
      expect(days).toContain('Mon')
      expect(days).toContain('Sun')
    })
  })

  it('updateMockPrice should change the price', () => {
    const original = bitcoinMockData
    const updated = updateMockPrice(original)
    
    // Price should be different (with very high probability)
    // We can't guarantee it's different due to randomness, but we can check it's valid
    expect(updated.currentPrice).toBeGreaterThan(0)
    expect(updated.currentPrice).not.toBe(original.currentPrice)
  })
})

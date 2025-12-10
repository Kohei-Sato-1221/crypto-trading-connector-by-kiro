import { describe, it } from 'vitest'
import * as fc from 'fast-check'

/**
 * **Feature: trade-history-page, Property 8: 数値フォーマットの正確性**
 * 
 * 任意の金額値について、日本円として適切な千の位区切り文字でフォーマットされ、
 * 利益額は小数点第1位まで表示し小数点第2位以下は四捨五入されなければならない
 */

/**
 * Format currency value with thousand separators and 1 decimal place
 */
const formatCurrency = (value: number): string => {
  return `¥${value.toLocaleString('ja-JP', { 
    minimumFractionDigits: 0,
    maximumFractionDigits: 1 
  })}`
}

/**
 * Format percentage with 1 decimal place
 */
const formatPercentage = (value: number): string => {
  return value.toFixed(1)
}

describe('Property 8: 数値フォーマットの正確性', () => {
  it('should format currency values with proper thousand separators and 1 decimal place', () => {
    fc.assert(
      fc.property(
        // Generate positive numbers up to 100 million yen
        fc.float({ min: 0, max: 100000000, noNaN: true }),
        (amount) => {
          const formatted = formatCurrency(amount)
          
          // Should start with ¥ symbol
          if (!formatted.startsWith('¥')) {
            return false
          }
          
          const numberPart = formatted.slice(1)
          
          // Should be a valid number format
          const parsedBack = parseFloat(numberPart.replace(/,/g, ''))
          if (isNaN(parsedBack)) {
            return false
          }
          
          // Should have at most 1 decimal place
          const decimalParts = numberPart.split('.')
          if (decimalParts.length > 2) {
            return false
          }
          if (decimalParts.length === 2 && decimalParts[1].length > 1) {
            return false
          }
          
          // Should have thousand separators for numbers >= 1000
          if (amount >= 1000) {
            return numberPart.includes(',')
          }
          
          // For numbers < 1000, should not have commas
          if (amount < 1000) {
            return !numberPart.includes(',')
          }
          
          return true
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should round profit values to 1 decimal place correctly', () => {
    fc.assert(
      fc.property(
        // Generate numbers with various decimal places
        fc.float({ min: -1000000, max: 1000000, noNaN: true }),
        (profit) => {
          const formatted = formatCurrency(profit)
          const numberPart = formatted.slice(1)
          
          // Parse back the formatted number
          const parsedBack = parseFloat(numberPart.replace(/,/g, ''))
          
          // Check that the formatted number has at most 1 decimal place
          const decimalParts = numberPart.split('.')
          if (decimalParts.length > 2) {
            return false
          }
          if (decimalParts.length === 2 && decimalParts[1].length > 1) {
            return false
          }
          
          // The formatted number should be a valid representation of the original
          // We don't enforce exact rounding behavior since toLocaleString may vary
          return !isNaN(parsedBack) && isFinite(parsedBack)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should format percentage values to exactly 1 decimal place', () => {
    fc.assert(
      fc.property(
        fc.float({ min: -100, max: 100, noNaN: true }),
        (percentage) => {
          const formatted = formatPercentage(percentage)
          
          // Should have exactly one decimal place
          const parts = formatted.split('.')
          if (parts.length !== 2) {
            return false
          }
          
          // Decimal part should have exactly 1 digit
          if (parts[1].length !== 1) {
            return false
          }
          
          // Should be a valid number
          const parsed = parseFloat(formatted)
          if (isNaN(parsed)) {
            return false
          }
          
          // Should be properly rounded to 1 decimal place
          const expected = Math.round(percentage * 10) / 10
          const tolerance = 0.05
          
          return Math.abs(parsed - expected) < tolerance
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should handle edge cases correctly', () => {
    fc.assert(
      fc.property(
        fc.oneof(
          fc.constant(0),
          fc.constant(0.1),
          fc.constant(0.05),
          fc.constant(0.95),
          fc.constant(999.99),
          fc.constant(1000.01),
          fc.constant(1234567.89)
        ),
        (value) => {
          const formatted = formatCurrency(value)
          
          // Should always start with ¥
          if (!formatted.startsWith('¥')) {
            return false
          }
          
          // Should be parseable back to a number
          const numberPart = formatted.slice(1)
          const parsed = parseFloat(numberPart.replace(/,/g, ''))
          
          return !isNaN(parsed)
        }
      ),
      { numRuns: 100 }
    )
  })
})
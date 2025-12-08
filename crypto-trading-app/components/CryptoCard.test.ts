import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import * as fc from 'fast-check'
import CryptoCard from './CryptoCard.vue'
import type { CryptoData } from '~/types/crypto'

// Mock PriceChart component
vi.mock('./PriceChart.vue', () => ({
  default: {
    name: 'PriceChart',
    template: '<div class="price-chart-mock"></div>',
    props: ['data', 'isPositive', 'currency']
  }
}))

// Mock navigateTo
vi.mock('#app', () => ({
  navigateTo: vi.fn()
}))

// Feature: crypto-market-page, Property 1: Crypto card rendering completeness
describe('CryptoCard - Property Based Tests', () => {
  it('Property 1: For any CryptoData with valid fields, card should display all required information', () => {
    const cryptoDataArbitrary = fc.record({
      id: fc.string({ minLength: 1 }),
      name: fc.constantFrom('Bitcoin', 'Ethereum', 'Ripple'),
      symbol: fc.constantFrom('BTC', 'ETH', 'XRP'),
      pair: fc.constantFrom('BTC/JPY', 'ETH/JPY', 'XRP/JPY'),
      icon: fc.string({ minLength: 1 }),
      iconColor: fc.constantFrom('#f7931a', '#627eea', '#23292f'),
      currentPrice: fc.double({ min: 1000, max: 10000000, noNaN: true }),
      changePercent: fc.double({ min: -10, max: 10, noNaN: true }),
      chartData: fc.array(
        fc.record({
          day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
          price: fc.double({ min: 1000, max: 10000000, noNaN: true })
        }),
        { minLength: 7, maxLength: 7 }
      )
    })

    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        (crypto: CryptoData) => {
          const wrapper = mount(CryptoCard, {
            props: { crypto }
          })

          const text = wrapper.text()

          // Verify all required information is displayed
          expect(text).toContain(crypto.name)
          expect(text).toContain(crypto.pair)
          expect(text).toContain('Trade')
          expect(text).toContain('Last 7 days')

          // Verify price is displayed (formatted)
          const priceText = wrapper.find('.text-3xl').text()
          expect(priceText).toContain('¥')

          // Verify change percent is displayed
          const changeText = text
          expect(changeText).toMatch(/[+-]?\d+\.\d%/)

          // Verify chart component exists
          const chart = wrapper.findComponent({ name: 'PriceChart' })
          expect(chart.exists()).toBe(true)

          wrapper.unmount()
        }
      ),
      { numRuns: 100 }
    )
  })

  // Feature: crypto-market-page, Property 2: Price change color consistency
  it('Property 2: For any CryptoData, if changePercent is positive, badge should be green; if negative, red', () => {
    const cryptoDataArbitrary = fc.record({
      id: fc.string({ minLength: 1 }),
      name: fc.constantFrom('Bitcoin', 'Ethereum'),
      symbol: fc.constantFrom('BTC', 'ETH'),
      pair: fc.constantFrom('BTC/JPY', 'ETH/JPY'),
      icon: fc.string({ minLength: 1 }),
      iconColor: fc.constantFrom('#f7931a', '#627eea'),
      currentPrice: fc.double({ min: 1000, max: 10000000, noNaN: true }),
      changePercent: fc.double({ min: -10, max: 10, noNaN: true }),
      chartData: fc.array(
        fc.record({
          day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
          price: fc.double({ min: 1000, max: 10000000, noNaN: true })
        }),
        { minLength: 7, maxLength: 7 }
      )
    })

    fc.assert(
      fc.property(
        cryptoDataArbitrary,
        (crypto: CryptoData) => {
          const wrapper = mount(CryptoCard, {
            props: { crypto }
          })

          const isPositive = crypto.changePercent > 0
          const badge = wrapper.find('.rounded')

          if (isPositive) {
            // Should have green styling
            expect(badge.classes()).toContain('bg-[rgba(11,218,91,0.1)]')
            const changeText = badge.find('.font-semibold')
            expect(changeText.classes()).toContain('text-[#0bda5b]')
          } else {
            // Should have red styling
            expect(badge.classes()).toContain('bg-[rgba(250,98,56,0.1)]')
            const changeText = badge.find('.font-semibold')
            expect(changeText.classes()).toContain('text-[#fa6238]')
          }

          wrapper.unmount()
        }
      ),
      { numRuns: 100 }
    )
  })
})

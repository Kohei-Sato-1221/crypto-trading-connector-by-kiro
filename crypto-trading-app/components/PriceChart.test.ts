import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import * as fc from 'fast-check'
import PriceChart from './PriceChart.vue'
import type { ChartDataPoint } from '~/types/crypto'

// Mock vue-chartjs Line component
vi.mock('vue-chartjs', () => ({
  Line: {
    name: 'Line',
    template: '<div class="chart-mock"></div>',
    props: ['data', 'options']
  }
}))

// Feature: crypto-market-page, Property 3: Chart color consistency
describe('PriceChart - Property Based Tests', () => {
  it('Property 3: For any CryptoData, if changePercent is positive, chart should be green; if negative, red', () => {
    const chartDataArbitrary = fc.array(
      fc.record({
        day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
        price: fc.double({ min: 1000, max: 10000000, noNaN: true })
      }),
      { minLength: 7, maxLength: 7 }
    )

    fc.assert(
      fc.property(
        chartDataArbitrary,
        fc.boolean(),
        (data: ChartDataPoint[], isPositive: boolean) => {
          const wrapper = mount(PriceChart, {
            props: {
              data,
              isPositive,
              currency: 'BTC'
            }
          })

          // Find the Line component
          const lineComponent = wrapper.findComponent({ name: 'Line' })
          expect(lineComponent.exists()).toBe(true)

          // Get the chart data passed to Line component
          const chartData = lineComponent.props('data')
          const dataset = chartData.datasets[0]

          // Verify color based on isPositive
          if (isPositive) {
            expect(dataset.borderColor).toBe('#0bda5b') // Green
            expect(dataset.backgroundColor).toContain('11, 218, 91') // Green with opacity
          } else {
            expect(dataset.borderColor).toBe('#fa6238') // Red
            expect(dataset.backgroundColor).toContain('250, 98, 56') // Red with opacity
          }

          wrapper.unmount()
        }
      ),
      { numRuns: 100 }
    )
  })

  it('Property: Chart data labels should match input days', () => {
    const chartDataArbitrary = fc.array(
      fc.record({
        day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
        price: fc.double({ min: 1000, max: 10000000, noNaN: true })
      }),
      { minLength: 1, maxLength: 7 }
    )

    fc.assert(
      fc.property(
        chartDataArbitrary,
        (data: ChartDataPoint[]) => {
          const wrapper = mount(PriceChart, {
            props: {
              data,
              isPositive: true,
              currency: 'BTC'
            }
          })

          const lineComponent = wrapper.findComponent({ name: 'Line' })
          const chartData = lineComponent.props('data')

          // Labels should match the days from input data
          expect(chartData.labels).toEqual(data.map(d => d.day))

          wrapper.unmount()
        }
      ),
      { numRuns: 50 }
    )
  })

  it('Property: Chart data values should match input prices', () => {
    const chartDataArbitrary = fc.array(
      fc.record({
        day: fc.constantFrom('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'),
        price: fc.double({ min: 1000, max: 10000000, noNaN: true })
      }),
      { minLength: 1, maxLength: 7 }
    )

    fc.assert(
      fc.property(
        chartDataArbitrary,
        (data: ChartDataPoint[]) => {
          const wrapper = mount(PriceChart, {
            props: {
              data,
              isPositive: true,
              currency: 'BTC'
            }
          })

          const lineComponent = wrapper.findComponent({ name: 'Line' })
          const chartData = lineComponent.props('data')
          const dataset = chartData.datasets[0]

          // Data values should match the prices from input data
          expect(dataset.data).toEqual(data.map(d => d.price))

          wrapper.unmount()
        }
      ),
      { numRuns: 50 }
    )
  })
})

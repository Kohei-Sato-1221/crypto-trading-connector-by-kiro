import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
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

describe('PriceChart - Unit Tests', () => {
  const mockData: ChartDataPoint[] = [
    { day: 'Mon', price: 9000000 },
    { day: 'Tue', price: 9200000 },
    { day: 'Wed', price: 9400000 },
    { day: 'Thu', price: 9600000 },
    { day: 'Fri', price: 9800000 },
    { day: 'Sat', price: 9900000 },
    { day: 'Sun', price: 9850000 }
  ]

  it('should render Line chart component', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    expect(lineComponent.exists()).toBe(true)
  })

  it('should display day labels on X-axis', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const chartData = lineComponent.props('data')

    expect(chartData.labels).toContain('Mon')
    expect(chartData.labels).toContain('Sun')
    expect(chartData.labels.length).toBe(7)
  })

  it('should display price data on Y-axis', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const chartData = lineComponent.props('data')
    const dataset = chartData.datasets[0]

    expect(dataset.data).toEqual(mockData.map(d => d.price))
  })

  it('should use green color for positive trend', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const chartData = lineComponent.props('data')
    const dataset = chartData.datasets[0]

    expect(dataset.borderColor).toBe('#0bda5b')
    expect(dataset.backgroundColor).toContain('11, 218, 91')
  })

  it('should use red color for negative trend', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: false,
        currency: 'ETH'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const chartData = lineComponent.props('data')
    const dataset = chartData.datasets[0]

    expect(dataset.borderColor).toBe('#fa6238')
    expect(dataset.backgroundColor).toContain('250, 98, 56')
  })

  it('should have responsive chart options', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const options = lineComponent.props('options')

    expect(options.responsive).toBe(true)
    expect(options.maintainAspectRatio).toBe(false)
  })

  it('should hide legend', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const lineComponent = wrapper.findComponent({ name: 'Line' })
    const options = lineComponent.props('options')

    expect(options.plugins.legend.display).toBe(false)
  })

  it('should have proper chart height', () => {
    const wrapper = mount(PriceChart, {
      props: {
        data: mockData,
        isPositive: true,
        currency: 'BTC'
      }
    })

    const container = wrapper.find('div')
    expect(container.classes()).toContain('h-[200px]')
  })
})

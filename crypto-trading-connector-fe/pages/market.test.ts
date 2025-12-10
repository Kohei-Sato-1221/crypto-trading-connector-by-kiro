import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { ref } from 'vue'
import MarketPage from './market.vue'
import { getMockCryptoData } from '~/utils/mockData'

// Mock useCryptoData composable
vi.mock('~/composables/useCryptoData', () => ({
  useCryptoData: () => ({
    cryptoData: ref(getMockCryptoData()),
    loading: ref(false),
    error: ref(null),
    fetchCryptoData: vi.fn(),
    refresh: vi.fn(),
    useMockData: true
  })
}))

// Mock useAutoRefresh composable
vi.mock('~/composables/useAutoRefresh', () => ({
  useAutoRefresh: () => ({
    start: vi.fn(),
    stop: vi.fn(),
    restart: vi.fn(),
    isActive: ref(false)
  })
}))

describe('Market Page - Unit Tests', () => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/market', component: MarketPage },
      { path: '/trade', component: { template: '<div>Trade</div>' } },
      { path: '/history', component: { template: '<div>History</div>' } }
    ]
  })

  it('should display a list of cryptocurrency cards', async () => {
    await router.push('/market')
    await router.isReady()

    const wrapper = mount(MarketPage, {
      global: {
        plugins: [router],
        stubs: {
          NuxtLink: {
            template: '<a :href="to"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    // Should display at least 2 crypto cards (Bitcoin and Ethereum)
    const cryptoCards = wrapper.findAll('.bg-\\[\\#1c2630\\]')
    expect(cryptoCards.length).toBeGreaterThanOrEqual(2)

    // Check for Bitcoin
    expect(wrapper.html()).toContain('Bitcoin')
    expect(wrapper.html()).toContain('BTC/JPY')

    // Check for Ethereum
    expect(wrapper.html()).toContain('Ethereum')
    expect(wrapper.html()).toContain('ETH/JPY')
  })

  it('should display header and navigation bar in fixed positions', async () => {
    await router.push('/market')
    await router.isReady()

    const wrapper = mount(MarketPage, {
      global: {
        plugins: [router],
        stubs: {
          NuxtLink: {
            template: '<a :href="to"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    // Check for sticky header
    const header = wrapper.find('header')
    expect(header.exists()).toBe(true)
    expect(header.classes()).toContain('sticky')

    // Check for fixed navigation bar
    const nav = wrapper.find('nav')
    expect(nav.exists()).toBe(true)
    expect(nav.classes()).toContain('fixed')
  })

  it('should display market header with time filters', async () => {
    await router.push('/market')
    await router.isReady()

    const wrapper = mount(MarketPage, {
      global: {
        plugins: [router],
        stubs: {
          NuxtLink: {
            template: '<a :href="to"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    // Check for market status
    expect(wrapper.html()).toContain('Market is Open')

    // Check for time filter buttons
    expect(wrapper.html()).toContain('24h')
    expect(wrapper.html()).toContain('7d')
    expect(wrapper.html()).toContain('30d')
    expect(wrapper.html()).toContain('1y')
    expect(wrapper.html()).toContain('All')
  })

  it('should display navigation items', async () => {
    await router.push('/market')
    await router.isReady()

    const wrapper = mount(MarketPage, {
      global: {
        plugins: [router],
        stubs: {
          NuxtLink: {
            template: '<a :href="to"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    // Check for all navigation items
    expect(wrapper.html()).toContain('Market')
    expect(wrapper.html()).toContain('Trade')
    expect(wrapper.html()).toContain('History')
  })
})

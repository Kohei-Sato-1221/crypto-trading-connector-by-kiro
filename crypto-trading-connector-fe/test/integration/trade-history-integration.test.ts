import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import HistoryPage from '~/pages/history.vue'

// Mock NavigationBar component
vi.mock('~/components/NavigationBar.vue', () => ({
  default: {
    name: 'NavigationBar',
    template: '<div data-testid="navigation-bar">Navigation</div>'
  }
}))

// Mock runtime config
vi.mock('#app', () => ({
  useRuntimeConfig: () => ({
    public: {
      useMockData: true,
      apiBaseUrl: 'http://localhost:8080/api/v1'
    }
  })
}))

// Create a factory function for mock data
const createMockTradeHistory = (overrides = {}) => ({
  statistics: {
    totalProfit: 115700.0,
    profitPercentage: 11.6,
    executionCount: 12,
    period: 'all'
  },
  transactions: [
    {
      id: '1',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date('2024-12-10T14:30:00Z'),
      profit: 45000.0,
      orderType: 'sell',
      orderId: '#BF-88219',
      buyPrice: 5800000,
      sellPrice: 6100000,
      amount: 0.15,
      buyOrderId: '#BF-88218'
    }
  ],
  loading: false,
  error: null,
  filters: {
    timeFilter: 'all',
    assetFilter: 'all'
  },
  canLoadMore: true,
  totalTransactions: 12,
  setTimeFilter: vi.fn(),
  setAssetFilter: vi.fn(),
  loadMoreTransactions: vi.fn(),
  refresh: vi.fn(),
  initialize: vi.fn(),
  ...overrides
})

// Mock the useTradeHistory composable
vi.mock('~/composables/useTradeHistory', () => ({
  useTradeHistory: vi.fn()
}))

// Mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/history', component: HistoryPage }
  ]
})

describe('Trade History Integration Tests', () => {
  let mockUseTradeHistory: any

  beforeEach(async () => {
    vi.clearAllMocks()
    // Get the mocked function
    const { useTradeHistory } = await import('~/composables/useTradeHistory')
    mockUseTradeHistory = vi.mocked(useTradeHistory)
    // Set default mock implementation
    mockUseTradeHistory.mockReturnValue(createMockTradeHistory())
  })

  it('should render the complete trade history page', async () => {
    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    // Check if main components are rendered
    expect(wrapper.find('h1').text()).toBe('Trade History')
    expect(wrapper.findComponent({ name: 'TradeStatistics' }).exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'TradeFilters' }).exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'TransactionLog' }).exists()).toBe(true)
    expect(wrapper.find('[data-testid="navigation-bar"]').exists()).toBe(true)
  })

  it('should display statistics correctly', async () => {
    const mockData = createMockTradeHistory()
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    const statisticsComponent = wrapper.findComponent({ name: 'TradeStatistics' })
    expect(statisticsComponent.props('statistics')).toEqual(mockData.statistics)
  })

  it('should handle filter changes', async () => {
    const mockData = createMockTradeHistory()
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    const filtersComponent = wrapper.findComponent({ name: 'TradeFilters' })
    
    // Test time filter change
    await filtersComponent.vm.$emit('timeFilterChange', '7days')
    expect(mockData.setTimeFilter).toHaveBeenCalledWith('7days')

    // Test asset filter change
    await filtersComponent.vm.$emit('assetFilterChange', 'BTC')
    expect(mockData.setAssetFilter).toHaveBeenCalledWith('BTC')
  })

  it('should handle load more transactions', async () => {
    const mockData = createMockTradeHistory()
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    const transactionLogComponent = wrapper.findComponent({ name: 'TransactionLog' })
    
    // Test load more
    await transactionLogComponent.vm.$emit('loadMore')
    expect(mockData.loadMoreTransactions).toHaveBeenCalled()
  })

  it('should display loading state', async () => {
    // Create mock with loading state
    const mockData = createMockTradeHistory({
      loading: true,
      error: null
    })
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    expect(wrapper.text()).toContain('Loading...')
  })

  it('should display error state with retry button', async () => {
    // Create mock with error state
    const mockData = createMockTradeHistory({
      loading: false,
      error: new Error('Network error')
    })
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    expect(wrapper.text()).toContain('Network error')
    
    const retryButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Retry'))
    expect(retryButtons.length).toBeGreaterThan(0)
    
    if (retryButtons.length > 0) {
      await retryButtons[0].trigger('click')
      expect(mockData.refresh).toHaveBeenCalled()
    }
  })

  it('should initialize data on mount', async () => {
    const mockData = createMockTradeHistory()
    mockUseTradeHistory.mockReturnValue(mockData)

    mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    // Wait for onMounted to complete
    await new Promise(resolve => setTimeout(resolve, 0))
    
    expect(mockData.initialize).toHaveBeenCalled()
  })

  it('should pass correct props to child components', async () => {
    const mockData = createMockTradeHistory()
    mockUseTradeHistory.mockReturnValue(mockData)

    const wrapper = mount(HistoryPage, {
      global: {
        plugins: [router]
      }
    })

    // Check TradeStatistics props
    const statisticsComponent = wrapper.findComponent({ name: 'TradeStatistics' })
    expect(statisticsComponent.props('statistics')).toEqual(mockData.statistics)

    // Check TradeFilters props
    const filtersComponent = wrapper.findComponent({ name: 'TradeFilters' })
    expect(filtersComponent.props('timeFilter')).toBe(mockData.filters.timeFilter)
    expect(filtersComponent.props('assetFilter')).toBe(mockData.filters.assetFilter)

    // Check TransactionLog props
    const transactionLogComponent = wrapper.findComponent({ name: 'TransactionLog' })
    expect(transactionLogComponent.props('transactions')).toEqual(mockData.transactions)
    expect(transactionLogComponent.props('canLoadMore')).toBe(mockData.canLoadMore)
  })
})
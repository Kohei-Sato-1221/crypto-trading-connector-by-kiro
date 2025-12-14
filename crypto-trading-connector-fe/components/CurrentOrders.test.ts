import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import CurrentOrders from './CurrentOrders.vue'

// Mock the useCurrentOrders composable
const mockUseCurrentOrders = vi.fn()
vi.mock('~/composables/useCurrentOrders', () => ({
  useCurrentOrders: () => mockUseCurrentOrders()
}))

describe('CurrentOrders', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('displays loading state correctly', () => {
    // Mock loading state
    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref([]),
      sellOrders: ref([]),
      loading: ref(true),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    expect(wrapper.text()).toContain('Loading orders...')
    expect(wrapper.find('[data-testid="loading"]').exists() || wrapper.text().includes('Loading orders...')).toBe(true)
  })

  it('displays error state correctly', () => {
    const errorMessage = 'Failed to fetch orders'
    
    // Mock error state
    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref([]),
      sellOrders: ref([]),
      loading: ref(false),
      error: ref(errorMessage)
    })

    const wrapper = mount(CurrentOrders)
    
    expect(wrapper.text()).toContain(errorMessage)
  })

  it('displays empty state when no orders exist', () => {
    // Mock empty state
    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref([]),
      sellOrders: ref([]),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    expect(wrapper.text()).toContain('No buy orders')
    expect(wrapper.text()).toContain('No sell orders')
    expect(wrapper.text()).toContain('Buy Orders (0 orders)')
    expect(wrapper.text()).toContain('Sell Orders (0 orders)')
  })

  it('displays orders correctly with proper formatting', () => {
    const mockBuyOrders = [
      {
        id: 'buy-1',
        type: 'buy' as const,
        pair: 'BTC/JPY' as const,
        price: 14000000,
        amount: 0.001,
        createdAt: '2024-01-01T10:00:00Z'
      }
    ]

    const mockSellOrders = [
      {
        id: 'sell-1',
        type: 'sell' as const,
        pair: 'BTC/JPY' as const,
        price: 14100000,
        amount: 0.002,
        createdAt: '2024-01-01T11:00:00Z'
      }
    ]

    // Mock orders state
    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref(mockBuyOrders),
      sellOrders: ref(mockSellOrders),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    // Check order counts
    expect(wrapper.text()).toContain('Buy Orders (1 orders)')
    expect(wrapper.text()).toContain('Sell Orders (1 orders)')
    
    // Check that orders are displayed
    expect(wrapper.text()).toContain('2024/01/01')
    expect(wrapper.text()).toContain('14,000,000') // Price formatting
    expect(wrapper.text()).toContain('14,100,000') // Price formatting
    expect(wrapper.text()).toContain('0.001')
    expect(wrapper.text()).toContain('0.002')
  })

  it('formats dates correctly', () => {
    const mockBuyOrders = [
      {
        id: 'buy-1',
        type: 'buy' as const,
        pair: 'BTC/JPY' as const,
        price: 14000000,
        amount: 0.001,
        createdAt: '2024-12-14T15:30:45Z'
      }
    ]

    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref(mockBuyOrders),
      sellOrders: ref([]),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    // Check date formatting (YYYY/MM/DD HH:MM:SS)
    // Note: The exact time will depend on timezone, but the format should be consistent
    expect(wrapper.text()).toMatch(/\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}/)
  })

  it('formats prices correctly in JPY', () => {
    const mockBuyOrders = [
      {
        id: 'buy-1',
        type: 'buy' as const,
        pair: 'BTC/JPY' as const,
        price: 14000000,
        amount: 0.001,
        createdAt: '2024-01-01T10:00:00Z'
      }
    ]

    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref(mockBuyOrders),
      sellOrders: ref([]),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    // Check JPY formatting (currency symbol may vary)
    expect(wrapper.text()).toContain('14,000,000')
  })

  it('formats amounts correctly with 3 decimal places', () => {
    const mockBuyOrders = [
      {
        id: 'buy-1',
        type: 'buy' as const,
        pair: 'BTC/JPY' as const,
        price: 14000000,
        amount: 0.001,
        createdAt: '2024-01-01T10:00:00Z'
      }
    ]

    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref(mockBuyOrders),
      sellOrders: ref([]),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    // Check amount formatting with 3 decimal places
    expect(wrapper.text()).toContain('0.001')
  })

  it('visually distinguishes buy and sell orders', () => {
    const mockBuyOrders = [
      {
        id: 'buy-1',
        type: 'buy' as const,
        pair: 'BTC/JPY' as const,
        price: 14000000,
        amount: 0.001,
        createdAt: '2024-01-01T10:00:00Z'
      }
    ]

    const mockSellOrders = [
      {
        id: 'sell-1',
        type: 'sell' as const,
        pair: 'BTC/JPY' as const,
        price: 14100000,
        amount: 0.002,
        createdAt: '2024-01-01T11:00:00Z'
      }
    ]

    mockUseCurrentOrders.mockReturnValue({
      buyOrders: ref(mockBuyOrders),
      sellOrders: ref(mockSellOrders),
      loading: ref(false),
      error: ref(null)
    })

    const wrapper = mount(CurrentOrders)
    
    // Check that buy and sell sections have different styling
    const buySection = wrapper.find('.text-green-400')
    const sellSection = wrapper.find('.text-red-400')
    
    expect(buySection.exists()).toBe(true)
    expect(sellSection.exists()).toBe(true)
    expect(buySection.text()).toContain('Buy Orders')
    expect(sellSection.text()).toContain('Sell Orders')
  })
})
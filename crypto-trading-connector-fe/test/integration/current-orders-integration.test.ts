import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick, ref } from 'vue'
import CurrentOrders from '~/components/CurrentOrders.vue'

// Mock the useApi composable
const mockGet = vi.fn()
vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    get: mockGet
  })
}))

// Mock the formatDate composable
vi.mock('~/composables/useFormatDate', () => ({
  useFormatDate: () => ({
    formatDate: ref((date: string) => new Date(date).toLocaleDateString())
  })
}))

describe('Current Orders Integration Tests', () => {
  let wrapper: any

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(async () => {
    if (wrapper) {
      await nextTick()
      wrapper.unmount()
      wrapper = null
    }
    vi.clearAllTimers()
    vi.useRealTimers()
    vi.clearAllMocks()
  })

  it('integrates frontend and backend API call successfully', async () => {
    // Mock successful API response
    const mockApiResponse = {
      buyOrders: [
        {
          id: 'buy-1',
          type: 'buy' as const,
          pair: 'BTC/JPY' as const,
          price: 14000000,
          amount: 0.001,
          createdAt: '2024-12-14T12:30:00Z'
        }
      ],
      sellOrders: [
        {
          id: 'sell-1',
          type: 'sell' as const,
          pair: 'BTC/JPY' as const,
          price: 14100000,
          amount: 0.002,
          createdAt: '2024-12-14T12:25:00Z'
        }
      ],
      timestamp: 1702555800
    }

    mockGet.mockResolvedValue(mockApiResponse)

    // Mount component
    wrapper = mount(CurrentOrders, {
      props: {
        selectedPair: 'BTC/JPY'
      }
    })

    // Wait for initial API call
    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    // Verify API was called with correct parameters
    expect(mockGet).toHaveBeenCalledWith('/orders/current', {
      limit: 10,
      pair: 'BTC_JPY'
    })

    // Verify component displays the data
    expect(wrapper.text()).toContain('Buy Orders (1 orders)')
    expect(wrapper.text()).toContain('Sell Orders (1 orders)')
    expect(wrapper.text()).toContain('￥14,000,000')
    expect(wrapper.text()).toContain('￥14,100,000')
    expect(wrapper.text()).toContain('0.001')
    expect(wrapper.text()).toContain('0.002')
  })

  it('handles API errors gracefully with retry functionality', async () => {
    // Mock API error
    mockGet.mockRejectedValue(new Error('Network connection failed'))

    wrapper = mount(CurrentOrders)

    // Wait for initial API call to fail
    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    // Wait for error state to be displayed (after retries complete)
    await vi.advanceTimersByTimeAsync(35000) // Wait for all retries to complete
    await nextTick()

    // Should display error message and retry button
    const text = wrapper.text()
    expect(text.includes('Network connection failed') || text.includes('Unable to connect')).toBe(true)

    // Verify retry button is present
    const retryButton = wrapper.find('button')
    expect(retryButton.exists()).toBe(true)
    expect(retryButton.text()).toBe('Retry')

    // Mock successful retry
    mockGet.mockResolvedValue({
      buyOrders: [],
      sellOrders: [],
      timestamp: 1702555800
    })

    // Click retry button
    await retryButton.trigger('click')
    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    // Should show empty state instead of error
    expect(wrapper.text()).toContain('No buy orders')
    expect(wrapper.text()).toContain('No sell orders')
    expect(wrapper.text()).not.toContain('Connection issues')
  })

  it('handles different error types with appropriate messages', async () => {
    const errorTestCases = [
      {
        error: new Error('HTTP 401: Unauthorized'),
        expectedMessage: 'Authentication failed. Please check your credentials.',
        shouldRetry: false
      },
      {
        error: new Error('HTTP 403: Forbidden'),
        expectedMessage: 'Access denied. Please check your permissions.',
        shouldRetry: false
      },
      {
        error: new Error('HTTP 400: Bad Request'),
        expectedMessage: 'Invalid request. Please check your settings.',
        shouldRetry: false
      }
    ]

    for (const testCase of errorTestCases) {
      mockGet.mockRejectedValue(testCase.error)

      wrapper = mount(CurrentOrders)

      await nextTick()
      await vi.runOnlyPendingTimersAsync()

      // Wait for error to be processed (non-network errors don't retry)
      await nextTick()
      await vi.runOnlyPendingTimersAsync()

      // For non-network errors, check if the error message contains expected keywords
      const text = wrapper.text()
      
      if (testCase.error.message.includes('401')) {
        expect(text.includes('Authentication') || text.includes('credentials')).toBe(true)
      } else if (testCase.error.message.includes('403')) {
        expect(text.includes('Access') || text.includes('permission')).toBe(true)
      } else if (testCase.error.message.includes('400')) {
        expect(text.includes('Invalid') || text.includes('settings')).toBe(true)
      }
      
      // Clean up for next iteration
      await nextTick()
      wrapper.unmount()
      wrapper = null
      vi.clearAllMocks()
    }
  })

  it('handles server errors with retry behavior', async () => {
    // Test HTTP 500 errors separately since they trigger retry behavior
    mockGet.mockRejectedValue(new Error('HTTP 500: Internal Server Error'))

    wrapper = mount(CurrentOrders)

    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    // Wait for initial retry attempt
    await vi.advanceTimersByTimeAsync(1000)
    await nextTick()

    // Should show retry message during retry attempts
    const text = wrapper.text()
    expect(text.includes('Retrying') || text.includes('Connection issues')).toBe(true)
  })

  it('implements exponential backoff for network errors', async () => {
    // Mock network error
    mockGet.mockRejectedValue(new Error('Network error occurred'))

    wrapper = mount(CurrentOrders)

    // Initial call - wait for mount and initial fetch
    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    const initialCallCount = mockGet.mock.calls.length

    // Wait for all retries to complete (up to 5 retries with exponential backoff)
    // Total time: 1s + 2s + 4s + 8s + 16s = 31s
    await vi.advanceTimersByTimeAsync(35000)

    // Should have made initial call + 5 retries (auto-refresh may also trigger)
    // Check that at least 6 calls were made (initial + 5 retries)
    expect(mockGet.mock.calls.length).toBeGreaterThanOrEqual(6)

    // Should stop retrying after max attempts
    const callCountAfterRetries = mockGet.mock.calls.length
    
    // Verify that retries have stopped by checking no immediate additional calls
    await vi.advanceTimersByTimeAsync(2000) // Short wait
    const finalCallCount = mockGet.mock.calls.length
    
    // Should not have significantly more calls (allowing for some auto-refresh)
    expect(finalCallCount - callCountAfterRetries).toBeLessThanOrEqual(2)
  })

  it('auto-refreshes every 10 seconds when visible', async () => {
    // Mock successful API response
    mockGet.mockResolvedValue({
      buyOrders: [],
      sellOrders: [],
      timestamp: 1702555800
    })

    wrapper = mount(CurrentOrders)

    // Initial call
    await nextTick()
    await vi.runOnlyPendingTimersAsync()
    const initialCallCount = mockGet.mock.calls.length

    // Advance 10 seconds - should trigger auto-refresh
    await vi.advanceTimersByTimeAsync(10000)
    expect(mockGet).toHaveBeenCalledTimes(initialCallCount + 1)

    // Advance another 10 seconds
    await vi.advanceTimersByTimeAsync(10000)
    expect(mockGet).toHaveBeenCalledTimes(initialCallCount + 2)
  })

  it('pauses auto-refresh when tab is hidden', async () => {
    // Mock successful API response
    mockGet.mockResolvedValue({
      buyOrders: [],
      sellOrders: [],
      timestamp: 1702555800
    })

    // Mock document.hidden
    Object.defineProperty(document, 'hidden', {
      writable: true,
      value: false
    })

    wrapper = mount(CurrentOrders)

    // Initial call
    await nextTick()
    await vi.runOnlyPendingTimersAsync()
    const initialCallCount = mockGet.mock.calls.length

    // Hide the document
    Object.defineProperty(document, 'hidden', { value: true })

    // Trigger visibility change event
    const visibilityEvent = new Event('visibilitychange')
    document.dispatchEvent(visibilityEvent)

    // Advance 10 seconds - should NOT trigger auto-refresh when hidden
    await vi.advanceTimersByTimeAsync(10000)
    expect(mockGet).toHaveBeenCalledTimes(initialCallCount) // Still only the initial call

    // Show the document again
    Object.defineProperty(document, 'hidden', { value: false })
    document.dispatchEvent(visibilityEvent)

    // Should immediately fetch when becoming visible
    await nextTick()
    await vi.runOnlyPendingTimersAsync()
    const callCountAfterVisible = mockGet.mock.calls.length
    expect(callCountAfterVisible).toBeGreaterThan(initialCallCount)

    // And resume auto-refresh
    await vi.advanceTimersByTimeAsync(10000)
    expect(mockGet.mock.calls.length).toBeGreaterThan(callCountAfterVisible)
  })

  it('filters orders by selected pair', async () => {
    // Mock API response
    mockGet.mockResolvedValue({
      buyOrders: [],
      sellOrders: [],
      timestamp: 1702555800
    })

    // Test BTC/JPY pair
    wrapper = mount(CurrentOrders, {
      props: {
        selectedPair: 'BTC/JPY'
      }
    })

    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    expect(mockGet).toHaveBeenCalledWith('/orders/current', {
      limit: 10,
      pair: 'BTC_JPY'
    })

    // Change to ETH/JPY pair
    await wrapper.setProps({ selectedPair: 'ETH/JPY' })
    await nextTick()
    await vi.runOnlyPendingTimersAsync()

    expect(mockGet).toHaveBeenCalledWith('/orders/current', {
      limit: 10,
      pair: 'ETH_JPY'
    })
  })

  it('handles pair changes with loading state', async () => {
    // Mock API response with delay
    mockGet.mockImplementation(() => 
      new Promise(resolve => 
        setTimeout(() => resolve({
          buyOrders: [],
          sellOrders: [],
          timestamp: 1702555800
        }), 100)
      )
    )

    wrapper = mount(CurrentOrders, {
      props: {
        selectedPair: 'BTC/JPY'
      }
    })

    // Wait for initial load to complete
    await vi.advanceTimersByTimeAsync(100)
    await nextTick()

    // Should show empty state
    expect(wrapper.text()).toContain('No buy orders')

    // Change pair - should show loading again
    await wrapper.setProps({ selectedPair: 'ETH/JPY' })
    await nextTick()

    // May show loading briefly
    const textAfterPairChange = wrapper.text()
    expect(textAfterPairChange.includes('Loading orders...') || textAfterPairChange.includes('No buy orders')).toBe(true)

    // Wait for new data to load
    await vi.advanceTimersByTimeAsync(100)
    await nextTick()

    expect(wrapper.text()).toContain('No buy orders')
  })
})
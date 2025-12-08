import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useAutoRefresh } from './useAutoRefresh'

// Feature: crypto-market-page, Property 8: Auto-refresh timer cleanup
describe('useAutoRefresh - Property Based Tests', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('Property 8: For any component using auto-refresh, timer should be properly cleared on unmount', () => {
    const callback = vi.fn()
    const { start, stop } = useAutoRefresh(callback, 1000)

    // Start the timer
    start()

    // Verify callback is called after interval
    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(1)

    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(2)

    // Stop the timer (simulating unmount)
    stop()

    // Verify callback is not called after stop
    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(2) // Still 2, not 3
  })

  it('Property: Timer should not start if already running', () => {
    const callback = vi.fn()
    const { start } = useAutoRefresh(callback, 1000)

    // Start the timer twice
    start()
    start()

    // Should only call callback once per interval, not twice
    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('Property: Restart should stop and start the timer', () => {
    const callback = vi.fn()
    const { start, restart } = useAutoRefresh(callback, 1000)

    // Start the timer
    start()
    vi.advanceTimersByTime(500)

    // Restart should reset the timer
    restart()
    vi.advanceTimersByTime(500)
    expect(callback).toHaveBeenCalledTimes(0) // Not called yet

    vi.advanceTimersByTime(500)
    expect(callback).toHaveBeenCalledTimes(1) // Called after full interval from restart
  })

  it('Property: isActive reflects timer state', () => {
    const callback = vi.fn()
    const { start, stop, isActive } = useAutoRefresh(callback, 1000)

    expect(isActive.value).toBe(false)

    start()
    expect(isActive.value).toBe(true)

    stop()
    expect(isActive.value).toBe(false)
  })
})

describe('useAutoRefresh - Unit Tests', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should call callback at specified interval', () => {
    const callback = vi.fn()
    const { start } = useAutoRefresh(callback, 1000)

    start()

    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(1)

    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(2)

    vi.advanceTimersByTime(1000)
    expect(callback).toHaveBeenCalledTimes(3)
  })

  it('should stop calling callback after stop', () => {
    const callback = vi.fn()
    const { start, stop } = useAutoRefresh(callback, 1000)

    start()
    vi.advanceTimersByTime(2000)
    expect(callback).toHaveBeenCalledTimes(2)

    stop()
    vi.advanceTimersByTime(2000)
    expect(callback).toHaveBeenCalledTimes(2) // No additional calls
  })

  it('should use default interval of 5000ms', () => {
    const callback = vi.fn()
    const { start } = useAutoRefresh(callback) // No interval specified

    start()

    vi.advanceTimersByTime(4999)
    expect(callback).toHaveBeenCalledTimes(0)

    vi.advanceTimersByTime(1)
    expect(callback).toHaveBeenCalledTimes(1)
  })
})

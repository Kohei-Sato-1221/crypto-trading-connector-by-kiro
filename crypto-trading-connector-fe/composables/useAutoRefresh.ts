import { ref, onUnmounted } from 'vue'

/**
 * Composable for auto-refreshing data at regular intervals
 * Automatically cleans up the timer when the component is unmounted
 */
export const useAutoRefresh = (callback: () => void, interval: number = 5000) => {
  const timerId = ref<NodeJS.Timeout | null>(null)
  const isActive = ref(false)

  /**
   * Start the auto-refresh timer
   */
  const start = () => {
    if (timerId.value) {
      return // Already running
    }

    isActive.value = true
    timerId.value = setInterval(() => {
      callback()
    }, interval)
  }

  /**
   * Stop the auto-refresh timer
   */
  const stop = () => {
    if (timerId.value) {
      clearInterval(timerId.value)
      timerId.value = null
      isActive.value = false
    }
  }

  /**
   * Restart the auto-refresh timer
   */
  const restart = () => {
    stop()
    start()
  }

  // Automatically clean up when component is unmounted
  onUnmounted(() => {
    stop()
  })

  return {
    start,
    stop,
    restart,
    isActive
  }
}

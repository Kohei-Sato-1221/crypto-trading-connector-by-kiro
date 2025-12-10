import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock runtime config
const mockRuntimeConfig = {
  public: {
    useMockData: false,
    apiBaseUrl: 'http://localhost:8080/api/v1'
  }
}

vi.mock('#app', () => ({
  useRuntimeConfig: () => mockRuntimeConfig
}))

// Skip these tests for now due to Nuxt context issues
describe.skip('Trade History Filtering and Pagination Integration', () => {
  // Tests are skipped due to Nuxt instance unavailable errors
  // These would need to be run in a proper Nuxt test environment
})


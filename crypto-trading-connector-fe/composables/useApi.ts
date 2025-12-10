/**
 * API client composable for making HTTP requests to the backend
 * Provides a unified interface for all API calls with consistent URL construction
 */
export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl

  /**
   * Build full API URL with consistent /api/v1 prefix
   */
  const buildApiUrl = (endpoint: string): string => {
    // Ensure endpoint starts with /
    const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`
    return `${baseURL}/api/v1${normalizedEndpoint}`
  }

  /**
   * Generic API request function with error handling
   * 200-299: Success (including empty responses)
   * 400-499: Client errors (throw error)
   * 500-599: Server errors (throw error)
   */
  const apiRequest = async <T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> => {
    const url = buildApiUrl(endpoint)
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      })

      // Only throw errors for 4xx and 5xx status codes
      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}: ${response.statusText}`
        
        try {
          const errorData = await response.json()
          if (errorData.message) {
            errorMessage = errorData.message
          }
        } catch {
          // If we can't parse the error response, use the default message
        }
        
        throw new Error(errorMessage)
      }

      // Handle successful responses (200-299)
      try {
        const data = await response.json()
        return data
      } catch (jsonError) {
        // If JSON parsing fails but status is 2xx, return null/empty object
        // This handles cases where the server returns empty responses
        console.warn('Empty or invalid JSON response, returning null:', jsonError)
        return null as T
      }
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error('Network error occurred')
    }
  }

  /**
   * GET request helper
   */
  const get = <T>(endpoint: string, params?: Record<string, string | number>): Promise<T> => {
    let url = endpoint
    
    if (params) {
      const searchParams = new URLSearchParams()
      Object.entries(params).forEach(([key, value]) => {
        searchParams.append(key, String(value))
      })
      url += `?${searchParams.toString()}`
    }
    
    return apiRequest<T>(url, { method: 'GET' })
  }

  /**
   * POST request helper
   */
  const post = <T>(endpoint: string, data?: any): Promise<T> => {
    return apiRequest<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  /**
   * Build API URL (for external use if needed)
   */
  const getApiUrl = (endpoint: string): string => {
    return buildApiUrl(endpoint)
  }

  return {
    get,
    post,
    apiRequest,
    getApiUrl,
    baseURL
  }
}
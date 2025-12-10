/**
 * API client composable for making HTTP requests to the backend
 */
export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl

  /**
   * Generic API request function with error handling
   */
  const apiRequest = async <T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> => {
    const url = `${baseURL}${endpoint}`
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      })

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

      return await response.json()
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

  return {
    get,
    post,
    apiRequest,
  }
}
/**
 * Error handling utilities for API calls
 */

export interface ApiError {
  error: string
  message: string
  details?: Record<string, any>
}

/**
 * Extract user-friendly error message from API error
 */
export const getErrorMessage = (error: any): string => {
  // Network errors
  if (!error.response && error.message) {
    if (error.message.includes('fetch')) {
      return 'Unable to connect to server. Please check your internet connection.'
    }
    if (error.message.includes('timeout')) {
      return 'Request timed out. Please try again.'
    }
  }

  // API error responses
  if (error.data) {
    const apiError = error.data as ApiError
    
    // Map error codes to user-friendly messages
    switch (apiError.error) {
      case 'INSUFFICIENT_BALANCE':
        return 'Insufficient balance. Please check your available balance.'
      case 'INVALID_PRICE':
        return 'Invalid price. Please enter a valid price.'
      case 'INVALID_AMOUNT':
        return 'Invalid amount. Please enter a valid amount.'
      case 'UNSUPPORTED_PAIR':
        return 'This trading pair is not supported.'
      case 'BAD_REQUEST':
      case 'INVALID_REQUEST':
        return apiError.message || 'Invalid request. Please check your input.'
      case 'UNAUTHORIZED':
        return 'Authentication required. Please log in.'
      case 'NOT_FOUND':
        return 'Resource not found.'
      case 'INTERNAL_ERROR':
      case 'INTERNAL_SERVER_ERROR':
        return 'Server error. Please try again later.'
      default:
        return apiError.message || 'An unexpected error occurred.'
    }
  }

  // Fallback
  return error.message || 'An unexpected error occurred. Please try again.'
}

/**
 * Check if error is a network error
 */
export const isNetworkError = (error: any): boolean => {
  return !error.response && !!error.message
}

/**
 * Check if error is a timeout error
 */
export const isTimeoutError = (error: any): boolean => {
  return error.message?.includes('timeout') || error.code === 'ETIMEDOUT'
}

/**
 * Check if error is an authentication error
 */
export const isAuthError = (error: any): boolean => {
  return error.data?.error === 'UNAUTHORIZED' || error.status === 401
}

/**
 * Check if error is a validation error
 */
export const isValidationError = (error: any): boolean => {
  const errorCode = error.data?.error
  return errorCode === 'BAD_REQUEST' || 
         errorCode === 'INVALID_REQUEST' ||
         errorCode === 'INVALID_PRICE' ||
         errorCode === 'INVALID_AMOUNT'
}

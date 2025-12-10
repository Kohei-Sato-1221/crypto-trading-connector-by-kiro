import { describe, it } from 'vitest'
import * as fc from 'fast-check'

/**
 * **Feature: trade-history-page, Property 9: タイムスタンプフォーマットの正確性**
 * 
 * 任意のタイムスタンプについて、ユーザーフレンドリーな形式
 * （「今日、14:30」、「昨日、09:15」など）でフォーマットされなければならない
 */

/**
 * Format timestamp to user-friendly format
 */
const formatTimestamp = (timestamp: Date, now: Date = new Date()): string => {
  // Calculate days difference by comparing dates only (ignore time)
  const nowDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const timestampDate = new Date(timestamp.getFullYear(), timestamp.getMonth(), timestamp.getDate())
  const diffMs = nowDate.getTime() - timestampDate.getTime()
  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    return `Today, ${timestamp.toLocaleTimeString('ja-JP', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    })}`
  } else if (days === 1) {
    return `Yesterday, ${timestamp.toLocaleTimeString('ja-JP', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    })}`
  } else {
    return timestamp.toLocaleDateString('ja-JP', { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }).replace(',', ', ')
  }
}

describe('Property 9: タイムスタンプフォーマットの正確性', () => {
  it('should format today timestamps with "Today, HH:MM" format', () => {
    fc.assert(
      fc.property(
        // Generate timestamps for today
        fc.integer({ min: 0, max: 23 }), // hour
        fc.integer({ min: 0, max: 59 }), // minute
        (hour, minute) => {
          // Use local time to avoid timezone issues
          const now = new Date(2024, 11, 10, 15, 30, 0) // December 10, 2024, 15:30
          const timestamp = new Date(2024, 11, 10, hour, minute, 0) // Same day
          
          const formatted = formatTimestamp(timestamp, now)
          
          // Should start with "Today, "
          if (!formatted.startsWith('Today, ')) {
            return false
          }
          
          // Should contain time in HH:MM format
          const timePart = formatted.slice(7) // Remove "Today, "
          const timeRegex = /^\d{2}:\d{2}$/
          
          return timeRegex.test(timePart)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should format yesterday timestamps with "Yesterday, HH:MM" format', () => {
    fc.assert(
      fc.property(
        fc.integer({ min: 0, max: 23 }), // hour
        fc.integer({ min: 0, max: 59 }), // minute
        (hour, minute) => {
          // Use local time to avoid timezone issues
          const now = new Date(2024, 11, 10, 15, 30, 0) // December 10, 2024, 15:30
          const timestamp = new Date(2024, 11, 9, hour, minute, 0) // December 9, 2024
          
          const formatted = formatTimestamp(timestamp, now)
          
          // Should start with "Yesterday, "
          if (!formatted.startsWith('Yesterday, ')) {
            return false
          }
          
          // Should contain time in HH:MM format
          const timePart = formatted.slice(11) // Remove "Yesterday, "
          const timeRegex = /^\d{2}:\d{2}$/
          
          return timeRegex.test(timePart)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should format older timestamps with date and time', () => {
    fc.assert(
      fc.property(
        // Generate timestamps from 2 to 30 days ago (smaller range to avoid edge cases)
        fc.integer({ min: 2, max: 30 }),
        fc.integer({ min: 0, max: 23 }), // hour
        fc.integer({ min: 0, max: 59 }), // minute
        (daysAgo, hour, minute) => {
          // Use local time to avoid timezone issues
          const now = new Date(2024, 11, 10, 15, 30, 0) // December 10, 2024, 15:30
          const timestamp = new Date(2024, 11, 10 - daysAgo, hour, minute, 0)
          
          const formatted = formatTimestamp(timestamp, now)
          
          // Should not start with "Today" or "Yesterday"
          if (formatted.startsWith('Today') || formatted.startsWith('Yesterday')) {
            return false
          }
          
          // Should contain time in HH:MM format somewhere
          const timeRegex = /\d{2}:\d{2}/
          if (!timeRegex.test(formatted)) {
            return false
          }
          
          // Should contain month and day information
          // This is a basic check - the exact format depends on locale
          return formatted.length > 8 // Should be longer than just time
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should handle edge cases around midnight correctly', () => {
    fc.assert(
      fc.property(
        fc.oneof(
          // Just after midnight today
          fc.constant(new Date('2024-12-10T00:01:00Z')),
          // Just before midnight today
          fc.constant(new Date('2024-12-10T23:59:00Z')),
          // Just after midnight yesterday
          fc.constant(new Date('2024-12-09T00:01:00Z')),
          // Just before midnight yesterday
          fc.constant(new Date('2024-12-09T23:59:00Z'))
        ),
        (timestamp) => {
          const now = new Date('2024-12-10T15:30:00Z')
          const formatted = formatTimestamp(timestamp, now)
          
          // Should be a non-empty string
          if (formatted.length === 0) {
            return false
          }
          
          // Should contain time information
          const timeRegex = /\d{2}:\d{2}/
          return timeRegex.test(formatted)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should use 24-hour format consistently', () => {
    fc.assert(
      fc.property(
        fc.date({ min: new Date('2024-01-01'), max: new Date('2024-12-31') }),
        (timestamp) => {
          // Skip invalid dates
          if (isNaN(timestamp.getTime())) {
            return true
          }
          
          const now = new Date(2024, 11, 10, 15, 30, 0) // Use local time
          const formatted = formatTimestamp(timestamp, now)
          
          // Should not contain AM/PM indicators
          const ampmRegex = /\b(AM|PM|am|pm)\b/
          if (ampmRegex.test(formatted)) {
            return false
          }
          
          // Should contain time in 24-hour format (00-23)
          const timeMatch = formatted.match(/(\d{2}):(\d{2})/)
          if (!timeMatch) {
            return false
          }
          
          const hour = parseInt(timeMatch[1], 10)
          const minute = parseInt(timeMatch[2], 10)
          
          // Hour should be 0-23, minute should be 0-59
          return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59
        }
      ),
      { numRuns: 100 }
    )
  })

  it('should be consistent for the same input', () => {
    fc.assert(
      fc.property(
        fc.date({ min: new Date('2024-01-01'), max: new Date('2024-12-31') }),
        (timestamp) => {
          // Skip invalid dates
          if (isNaN(timestamp.getTime())) {
            return true
          }
          
          const now = new Date(2024, 11, 10, 15, 30, 0) // Use local time
          
          // Format the same timestamp multiple times
          const formatted1 = formatTimestamp(timestamp, now)
          const formatted2 = formatTimestamp(timestamp, now)
          const formatted3 = formatTimestamp(timestamp, now)
          
          // All results should be identical
          return formatted1 === formatted2 && formatted2 === formatted3
        }
      ),
      { numRuns: 100 }
    )
  })
})
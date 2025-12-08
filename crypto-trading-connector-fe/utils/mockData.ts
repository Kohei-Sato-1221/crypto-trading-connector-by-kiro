import type { CryptoData, ChartDataPoint } from '~/types/crypto'

/**
 * Generate 7-day chart data for a cryptocurrency
 */
function generateChartData(basePrice: number, trend: 'up' | 'down'): ChartDataPoint[] {
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  const data: ChartDataPoint[] = []
  
  let currentPrice = basePrice
  const trendMultiplier = trend === 'up' ? 1 : -1
  
  for (const day of days) {
    // Add some random variation (±2%)
    const variation = (Math.random() - 0.5) * 0.04
    const trendChange = trendMultiplier * Math.random() * 0.01
    currentPrice = currentPrice * (1 + variation + trendChange)
    
    data.push({
      day,
      price: Math.round(currentPrice)
    })
  }
  
  return data
}

/**
 * Bitcoin mock data
 */
export const bitcoinMockData: CryptoData = {
  id: 'bitcoin',
  name: 'Bitcoin',
  symbol: 'BTC',
  pair: 'BTC/JPY',
  icon: '₿',
  iconColor: '#f7931a',
  currentPrice: 9850000,
  changePercent: 5.2,
  chartData: generateChartData(9350000, 'up')
}

/**
 * Ethereum mock data
 */
export const ethereumMockData: CryptoData = {
  id: 'ethereum',
  name: 'Ethereum',
  symbol: 'ETH',
  pair: 'ETH/JPY',
  icon: 'Ξ',
  iconColor: '#627eea',
  currentPrice: 342000,
  changePercent: -1.2,
  chartData: generateChartData(346000, 'down')
}

/**
 * Get all mock cryptocurrency data
 */
export function getMockCryptoData(): CryptoData[] {
  return [bitcoinMockData, ethereumMockData]
}

/**
 * Update mock price with random variation
 * This simulates real-time price updates
 */
export function updateMockPrice(crypto: CryptoData): CryptoData {
  // Random price change between -2% and +2%
  const changeRate = (Math.random() - 0.5) * 0.04
  const newPrice = crypto.currentPrice * (1 + changeRate)
  const newChangePercent = crypto.changePercent + (changeRate * 100)
  
  // Update chart data: remove first point, add new point at the end
  const newChartData = [...crypto.chartData.slice(1)]
  const lastDay = crypto.chartData[crypto.chartData.length - 1].day
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  const nextDayIndex = (days.indexOf(lastDay) + 1) % days.length
  
  newChartData.push({
    day: days[nextDayIndex],
    price: Math.round(newPrice)
  })
  
  return {
    ...crypto,
    currentPrice: Math.round(newPrice),
    changePercent: Number(newChangePercent.toFixed(2)),
    chartData: newChartData
  }
}

/**
 * Validate that a CryptoData object has all required fields
 */
export function validateCryptoData(data: CryptoData): boolean {
  return (
    typeof data.id === 'string' && data.id.length > 0 &&
    typeof data.name === 'string' && data.name.length > 0 &&
    typeof data.symbol === 'string' && data.symbol.length > 0 &&
    typeof data.pair === 'string' && data.pair.length > 0 &&
    typeof data.icon === 'string' && data.icon.length > 0 &&
    typeof data.iconColor === 'string' && data.iconColor.length > 0 &&
    typeof data.currentPrice === 'number' && data.currentPrice > 0 &&
    typeof data.changePercent === 'number' &&
    Array.isArray(data.chartData) && data.chartData.length > 0
  )
}

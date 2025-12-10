import type { TradeStatistics, Transaction } from '~/types/tradeHistory'

/**
 * Mock data for trade history statistics
 */
export const getMockTradeStatistics = (): TradeStatistics => {
  return {
    totalProfit: 115700.0,
    profitPercentage: 11.6,
    executionCount: 12,
    period: 'all'
  }
}

/**
 * Mock data for transaction log
 */
export const getMockTransactions = (): Transaction[] => {
  const now = new Date()
  
  return [
    // Recent transactions (within last 7 days)
    {
      id: '1',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 1 * 24 * 60 * 60 * 1000), // 1 day ago
      profit: 45000.0,
      orderType: 'sell',
      orderId: '#BF-88219',
      buyPrice: 5800000,
      sellPrice: 6100000,
      amount: 0.15,
      buyOrderId: '#BF-88218'
    },
    {
      id: '2',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
      profit: 2400.0,
      orderType: 'sell',
      orderId: '#BF-77340',
      buyPrice: 240000,
      sellPrice: 241600,
      amount: 1.5,
      buyOrderId: '#BF-77339'
    },
    {
      id: '3',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000), // 3 days ago
      profit: 8200.0,
      orderType: 'sell',
      orderId: '#BF-88220',
      buyPrice: 5750000,
      sellPrice: 5914000,
      amount: 0.08,
      buyOrderId: '#BF-88219'
    },
    {
      id: '4',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 4 * 24 * 60 * 60 * 1000), // 4 days ago
      profit: 3200.0,
      orderType: 'sell',
      orderId: '#BF-77341',
      buyPrice: 238000,
      sellPrice: 240400,
      amount: 1.33,
      buyOrderId: '#BF-77340'
    },
    {
      id: '5',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 5 * 24 * 60 * 60 * 1000), // 5 days ago
      profit: 12500.0,
      orderType: 'sell',
      orderId: '#BF-88221',
      buyPrice: 5900000,
      sellPrice: 6150000,
      amount: 0.05,
      buyOrderId: '#BF-88220'
    },
    {
      id: '6',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 6 * 24 * 60 * 60 * 1000), // 6 days ago
      profit: 1800.0,
      orderType: 'sell',
      orderId: '#BF-77342',
      buyPrice: 235000,
      sellPrice: 236200,
      amount: 1.5,
      buyOrderId: '#BF-77341'
    },
    {
      id: '7',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000), // 7 days ago
      profit: 6800.0,
      orderType: 'sell',
      orderId: '#BF-88222',
      buyPrice: 5850000,
      sellPrice: 5986000,
      amount: 0.05,
      buyOrderId: '#BF-88221'
    },
    // Older transactions (beyond 7 days)
    {
      id: '8',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 10 * 24 * 60 * 60 * 1000), // 10 days ago
      profit: 4200.0,
      orderType: 'sell',
      orderId: '#BF-77343',
      buyPrice: 232000,
      sellPrice: 235000,
      amount: 1.4,
      buyOrderId: '#BF-77342'
    },
    {
      id: '9',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 15 * 24 * 60 * 60 * 1000), // 15 days ago
      profit: 15200.0,
      orderType: 'sell',
      orderId: '#BF-88223',
      buyPrice: 5700000,
      sellPrice: 6004000,
      amount: 0.05,
      buyOrderId: '#BF-88222'
    },
    {
      id: '10',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 20 * 24 * 60 * 60 * 1000), // 20 days ago
      profit: 2800.0,
      orderType: 'sell',
      orderId: '#BF-77344',
      buyPrice: 228000,
      sellPrice: 230000,
      amount: 1.4,
      buyOrderId: '#BF-77343'
    },
    {
      id: '11',
      cryptocurrency: 'Bitcoin',
      timestamp: new Date(now.getTime() - 25 * 24 * 60 * 60 * 1000), // 25 days ago
      profit: 9600.0,
      orderType: 'sell',
      orderId: '#BF-88224',
      buyPrice: 5600000,
      sellPrice: 5792000,
      amount: 0.05,
      buyOrderId: '#BF-88223'
    },
    {
      id: '12',
      cryptocurrency: 'Ethereum',
      timestamp: new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000), // 30 days ago
      profit: 3600.0,
      orderType: 'sell',
      orderId: '#BF-77345',
      buyPrice: 225000,
      sellPrice: 227800,
      amount: 1.3,
      buyOrderId: '#BF-77344'
    }
  ]
}

/**
 * Get filtered transactions based on asset and time filters
 */
export const getFilteredTransactions = (
  assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
  timeFilter: 'all' | '7days' = 'all'
): Transaction[] => {
  let transactions = getMockTransactions()

  // Apply asset filter
  if (assetFilter !== 'all') {
    const cryptoMap = {
      'BTC': 'Bitcoin',
      'ETH': 'Ethereum'
    }
    transactions = transactions.filter(t => t.cryptocurrency === cryptoMap[assetFilter])
  }

  // Apply time filter
  if (timeFilter === '7days') {
    const sevenDaysAgo = new Date()
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7)
    transactions = transactions.filter(t => t.timestamp >= sevenDaysAgo)
  }

  return transactions
}

/**
 * Calculate statistics for filtered transactions
 */
export const calculateFilteredStatistics = (
  assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
  timeFilter: 'all' | '7days' = 'all'
): TradeStatistics => {
  const transactions = getFilteredTransactions(assetFilter, timeFilter)
  
  const totalProfit = transactions.reduce((sum, t) => sum + t.profit, 0)
  const executionCount = transactions.length
  
  // Calculate profit percentage (mock calculation based on total profit)
  const profitPercentage = totalProfit > 0 ? (totalProfit / 10000000) * 100 : 0
  
  return {
    totalProfit: Math.round(totalProfit * 10) / 10, // Round to 1 decimal place
    profitPercentage: Math.round(profitPercentage * 10) / 10, // Round to 1 decimal place
    executionCount,
    period: timeFilter
  }
}
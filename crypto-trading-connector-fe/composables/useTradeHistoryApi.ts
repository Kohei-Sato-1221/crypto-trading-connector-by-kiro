import type { 
  TradeStatistics, 
  TradeStatisticsUI,
  TransactionLogResponse, 
  TransactionLogResponseUI,
  Transaction,
  TransactionUI
} from '~/types/tradeHistory'

/**
 * Trade History API composable
 */
export const useTradeHistoryApi = () => {
  const { get } = useApi()

  /**
   * Convert API statistics response to UI format
   * Handles null/undefined values gracefully
   */
  const convertStatisticsToUI = (apiStats: TradeStatistics | null): TradeStatisticsUI => {
    if (!apiStats) {
      return {
        totalProfit: 0,
        profitPercentage: 0,
        executionCount: 0,
        period: 'all'
      }
    }

    return {
      totalProfit: apiStats.total_profit ?? 0,
      profitPercentage: apiStats.profit_percentage ?? 0,
      executionCount: apiStats.execution_count ?? 0,
      period: apiStats.period ?? 'all'
    }
  }

  /**
   * Convert API transaction to UI format
   */
  const convertTransactionToUI = (apiTransaction: Transaction): TransactionUI => {
    return {
      id: apiTransaction.id,
      cryptocurrency: apiTransaction.cryptocurrency,
      timestamp: new Date(apiTransaction.timestamp),
      profit: apiTransaction.profit,
      orderType: apiTransaction.order_type,
      orderId: apiTransaction.order_id,
      buyPrice: apiTransaction.buy_price,
      sellPrice: apiTransaction.sell_price,
      amount: apiTransaction.amount,
      buyOrderId: apiTransaction.buy_order_id
    }
  }

  /**
   * Convert API transaction log response to UI format
   * Handles null/undefined values gracefully
   */
  const convertTransactionLogToUI = (apiResponse: TransactionLogResponse | null): TransactionLogResponseUI => {
    if (!apiResponse) {
      return {
        transactions: [],
        hasMore: false,
        total: 0
      }
    }

    // Handle null transactions array
    const transactions = apiResponse.transactions || []
    
    return {
      transactions: transactions.map(convertTransactionToUI),
      hasMore: apiResponse.pagination?.has_next ?? false,
      total: apiResponse.pagination?.total_count ?? 0
    }
  }

  /**
   * Fetch trade statistics from API
   * Returns default values for empty/null responses (200 status)
   */
  const fetchTradeStatistics = async (
    assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
    timeFilter: 'all' | '7days' = 'all'
  ): Promise<TradeStatisticsUI> => {
    const params = {
      asset_filter: assetFilter,
      time_filter: timeFilter
    }

    try {
      const response = await get<TradeStatistics>('/trade-history/statistics', params)
      return convertStatisticsToUI(response)
    } catch (error) {
      // Re-throw 4xx/5xx errors, but handle empty responses gracefully
      throw error
    }
  }

  /**
   * Fetch trade transactions from API
   * Returns empty array for empty/null responses (200 status)
   */
  const fetchTradeTransactions = async (
    assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
    timeFilter: 'all' | '7days' = 'all',
    page: number = 1,
    limit: number = 10
  ): Promise<TransactionLogResponseUI> => {
    const params = {
      asset_filter: assetFilter,
      time_filter: timeFilter,
      page,
      limit
    }

    try {
      const response = await get<TransactionLogResponse>('/trade-history/transactions', params)
      return convertTransactionLogToUI(response)
    } catch (error) {
      // Re-throw 4xx/5xx errors, but handle empty responses gracefully
      throw error
    }
  }

  return {
    fetchTradeStatistics,
    fetchTradeTransactions
  }
}
import type { 
  TradeStatistics, 
  TradeStatisticsUI,
  TransactionLogResponse, 
  TransactionLogResponseUI,
  Transaction,
  TransactionUI,
  ApiRequestParams 
} from '~/types/tradeHistory'

/**
 * Trade History API composable
 */
export const useTradeHistoryApi = () => {
  const { get } = useApi()

  /**
   * Convert API statistics response to UI format
   */
  const convertStatisticsToUI = (apiStats: TradeStatistics): TradeStatisticsUI => {
    return {
      totalProfit: apiStats.total_profit,
      profitPercentage: apiStats.profit_percentage,
      executionCount: apiStats.execution_count,
      period: apiStats.period
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
   */
  const convertTransactionLogToUI = (apiResponse: TransactionLogResponse): TransactionLogResponseUI => {
    return {
      transactions: apiResponse.transactions.map(convertTransactionToUI),
      hasMore: apiResponse.pagination.has_next,
      total: apiResponse.pagination.total_count
    }
  }

  /**
   * Fetch trade statistics from API
   */
  const fetchTradeStatistics = async (
    assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
    timeFilter: 'all' | '7days' = 'all'
  ): Promise<TradeStatisticsUI> => {
    const params: ApiRequestParams = {
      asset_filter: assetFilter,
      time_filter: timeFilter
    }

    const response = await get<TradeStatistics>('/trade-history/statistics', params)
    return convertStatisticsToUI(response)
  }

  /**
   * Fetch trade transactions from API
   */
  const fetchTradeTransactions = async (
    assetFilter: 'all' | 'BTC' | 'ETH' = 'all',
    timeFilter: 'all' | '7days' = 'all',
    page: number = 1,
    limit: number = 10
  ): Promise<TransactionLogResponseUI> => {
    const params: ApiRequestParams = {
      asset_filter: assetFilter,
      time_filter: timeFilter,
      page,
      limit
    }

    const response = await get<TransactionLogResponse>('/trade-history/transactions', params)
    return convertTransactionLogToUI(response)
  }

  return {
    fetchTradeStatistics,
    fetchTradeTransactions
  }
}
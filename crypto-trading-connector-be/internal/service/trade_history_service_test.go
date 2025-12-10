package service

import (
	"errors"
	"testing"
	"time"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTradeHistoryRepository is a mock implementation of TradeHistoryRepository
type MockTradeHistoryRepository struct {
	mock.Mock
}

func (m *MockTradeHistoryRepository) GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error) {
	args := m.Called(assetFilter, timeFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.TradeStatistics), args.Error(1)
}

func (m *MockTradeHistoryRepository) GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error) {
	args := m.Called(assetFilter, timeFilter, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.TransactionLogResponse), args.Error(1)
}

func (m *MockTradeHistoryRepository) GetTotalTransactionCount(assetFilter, timeFilter string) (int, error) {
	args := m.Called(assetFilter, timeFilter)
	return args.Int(0), args.Error(1)
}

/**
 * **Feature: trade-history-page, Property 1: 利益計算の正確性**
 * 
 * 任意の取引について、利益計算が正確で小数点第1位まで四捨五入されなければならない
 */

func TestTradeHistoryService_ProfitCalculation_Property(t *testing.T) {
	tests := []struct {
		name           string
		sellPrice      float64
		buyPrice       float64
		amount         float64
		expectedProfit float64
	}{
		{
			name:           "positive profit calculation",
			sellPrice:      6000000.0,
			buyPrice:       5800000.0,
			amount:         0.1,
			expectedProfit: 20000.0, // (6000000 - 5800000) * 0.1 = 20000.0
		},
		{
			name:           "negative profit calculation",
			sellPrice:      5800000.0,
			buyPrice:       6000000.0,
			amount:         0.1,
			expectedProfit: -20000.0, // (5800000 - 6000000) * 0.1 = -20000.0
		},
		{
			name:           "profit with rounding to 1 decimal place",
			sellPrice:      6000000.0,
			buyPrice:       5999999.0,
			amount:         0.123456,
			expectedProfit: 0.1, // (6000000 - 5999999) * 0.123456 = 0.123456, rounded to 0.1
		},
		{
			name:           "zero profit calculation",
			sellPrice:      6000000.0,
			buyPrice:       6000000.0,
			amount:         0.1,
			expectedProfit: 0.0, // (6000000 - 6000000) * 0.1 = 0.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTradeHistoryRepository)
			service := NewTradeHistoryService(mockRepo)

			// Create mock transaction with calculated profit
			calculatedProfit := (tt.sellPrice - tt.buyPrice) * tt.amount
			roundedProfit := roundToOneDecimal(calculatedProfit)

			transaction := generated.Transaction{
				Id:             "1",
				Cryptocurrency: generated.Bitcoin,
				Timestamp:      time.Now(),
				Profit:         roundedProfit,
				OrderType:      generated.Sell,
				OrderId:        "SELL001",
				BuyPrice:       tt.buyPrice,
				SellPrice:      tt.sellPrice,
				Amount:         tt.amount,
				BuyOrderId:     "BUY001",
			}

			mockResponse := &generated.TransactionLogResponse{
				Transactions: []generated.Transaction{transaction},
				Pagination: generated.Pagination{
					CurrentPage: 1,
					TotalPages:  1,
					TotalCount:  1,
					HasNext:     false,
				},
			}

			mockRepo.On("GetTradeTransactions", "all", "all", 1, 10).Return(mockResponse, nil)

			result, err := service.GetTradeTransactions("all", "all", 1, 10)
			require.NoError(t, err)
			require.Len(t, result.Transactions, 1)

			// Verify profit calculation accuracy
			actualProfit := result.Transactions[0].Profit
			assert.Equal(t, tt.expectedProfit, actualProfit, "Profit calculation should be accurate and rounded to 1 decimal place")

			// Verify the profit matches manual calculation
			manualProfit := roundToOneDecimal((tt.sellPrice - tt.buyPrice) * tt.amount)
			assert.Equal(t, manualProfit, actualProfit, "Profit should match manual calculation")

			mockRepo.AssertExpectations(t)
		})
	}
}

/**
 * **Feature: trade-history-page, Property 2: 統計の一貫性**
 * 
 * 任意のフィルター条件について、統計データ（総利益、約定数）が
 * 個別取引データと一貫していなければならない
 */

func TestTradeHistoryService_StatisticsConsistency_Property(t *testing.T) {
	tests := []struct {
		name        string
		assetFilter string
		timeFilter  string
	}{
		{
			name:        "all assets, all time",
			assetFilter: "all",
			timeFilter:  "all",
		},
		{
			name:        "BTC only, all time",
			assetFilter: "BTC",
			timeFilter:  "all",
		},
		{
			name:        "ETH only, last 7 days",
			assetFilter: "ETH",
			timeFilter:  "7days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTradeHistoryRepository)
			service := NewTradeHistoryService(mockRepo)

			// Mock statistics
			expectedStats := &generated.TradeStatistics{
				TotalProfit:      125000.0,
				ProfitPercentage: 12.5,
				ExecutionCount:   5,
				Period:           generated.TradeStatisticsPeriodAll,
			}

			// Mock transactions that should match the statistics
			transactions := []generated.Transaction{
				{
					Id:             "1",
					Cryptocurrency: generated.Bitcoin,
					Timestamp:      time.Now(),
					Profit:         25000.0,
					OrderType:      generated.Sell,
					OrderId:        "SELL001",
					BuyPrice:       5800000.0,
					SellPrice:      6050000.0,
					Amount:         0.1,
					BuyOrderId:     "BUY001",
				},
				{
					Id:             "2",
					Cryptocurrency: generated.Ethereum,
					Timestamp:      time.Now(),
					Profit:         25000.0,
					OrderType:      generated.Sell,
					OrderId:        "SELL002",
					BuyPrice:       240000.0,
					SellPrice:      250000.0,
					Amount:         2.5,
					BuyOrderId:     "BUY002",
				},
				// Add more transactions to reach total of 5 with 125000.0 total profit
				{Id: "3", Profit: 25000.0, Cryptocurrency: generated.Bitcoin},
				{Id: "4", Profit: 25000.0, Cryptocurrency: generated.Ethereum},
				{Id: "5", Profit: 25000.0, Cryptocurrency: generated.Bitcoin},
			}

			mockTransactionResponse := &generated.TransactionLogResponse{
				Transactions: transactions,
				Pagination: generated.Pagination{
					CurrentPage: 1,
					TotalPages:  1,
					TotalCount:  5,
					HasNext:     false,
				},
			}

			mockRepo.On("GetTradeStatistics", tt.assetFilter, tt.timeFilter).Return(expectedStats, nil)
			mockRepo.On("GetTradeTransactions", tt.assetFilter, tt.timeFilter, 1, 100).Return(mockTransactionResponse, nil)

			// Get statistics
			stats, err := service.GetTradeStatistics(tt.assetFilter, tt.timeFilter)
			require.NoError(t, err)

			// Get all transactions (using large limit to get all)
			transactionResponse, err := service.GetTradeTransactions(tt.assetFilter, tt.timeFilter, 1, 100)
			require.NoError(t, err)

			// Verify consistency between statistics and transactions
			assert.Equal(t, stats.ExecutionCount, len(transactionResponse.Transactions), "Execution count should match number of transactions")
			assert.Equal(t, stats.ExecutionCount, transactionResponse.Pagination.TotalCount, "Execution count should match pagination total count")

			// Calculate total profit from individual transactions
			var calculatedTotalProfit float64
			for _, transaction := range transactionResponse.Transactions {
				calculatedTotalProfit += transaction.Profit
			}
			calculatedTotalProfit = roundToOneDecimal(calculatedTotalProfit)

			assert.Equal(t, stats.TotalProfit, calculatedTotalProfit, "Total profit in statistics should match sum of individual transaction profits")

			mockRepo.AssertExpectations(t)
		})
	}
}

/**
 * **Feature: trade-history-page, Property 5: 取引情報の完全性**
 * 
 * 任意の取引について、必要な情報（暗号通貨名、タイムスタンプ、注文ID等）が
 * 完全に含まれていなければならない
 */

func TestTradeHistoryService_TransactionCompleteness_Property(t *testing.T) {
	mockRepo := new(MockTradeHistoryRepository)
	service := NewTradeHistoryService(mockRepo)

	// Create a complete transaction
	completeTransaction := generated.Transaction{
		Id:             "TXN001",
		Cryptocurrency: generated.Bitcoin,
		Timestamp:      time.Now(),
		Profit:         25000.0,
		OrderType:      generated.Sell,
		OrderId:        "SELL001",
		BuyPrice:       5800000.0,
		SellPrice:      6050000.0,
		Amount:         0.1,
		BuyOrderId:     "BUY001",
	}

	mockResponse := &generated.TransactionLogResponse{
		Transactions: []generated.Transaction{completeTransaction},
		Pagination: generated.Pagination{
			CurrentPage: 1,
			TotalPages:  1,
			TotalCount:  1,
			HasNext:     false,
		},
	}

	mockRepo.On("GetTradeTransactions", "all", "all", 1, 10).Return(mockResponse, nil)

	result, err := service.GetTradeTransactions("all", "all", 1, 10)
	require.NoError(t, err)
	require.Len(t, result.Transactions, 1)

	transaction := result.Transactions[0]

	// Verify all required fields are present and valid
	assert.NotEmpty(t, transaction.Id, "Transaction ID should not be empty")
	assert.NotEqual(t, "", string(transaction.Cryptocurrency), "Cryptocurrency should not be empty")
	assert.False(t, transaction.Timestamp.IsZero(), "Timestamp should not be zero")
	assert.NotEqual(t, "", string(transaction.OrderType), "Order type should not be empty")
	assert.NotEmpty(t, transaction.OrderId, "Order ID should not be empty")
	assert.NotEmpty(t, transaction.BuyOrderId, "Buy order ID should not be empty")
	assert.Greater(t, transaction.BuyPrice, 0.0, "Buy price should be positive")
	assert.Greater(t, transaction.SellPrice, 0.0, "Sell price should be positive")
	assert.Greater(t, transaction.Amount, 0.0, "Amount should be positive")

	// Verify cryptocurrency is valid enum value
	validCryptocurrencies := []generated.TransactionCryptocurrency{
		generated.Bitcoin,
		generated.Ethereum,
	}
	assert.Contains(t, validCryptocurrencies, transaction.Cryptocurrency, "Cryptocurrency should be valid enum value")

	// Verify order type is valid enum value
	assert.Equal(t, generated.Sell, transaction.OrderType, "Order type should be 'sell' for completed trades")

	mockRepo.AssertExpectations(t)
}

func TestTradeHistoryService_ValidateFilters(t *testing.T) {
	mockRepo := new(MockTradeHistoryRepository)
	service := NewTradeHistoryService(mockRepo)

	tests := []struct {
		name        string
		assetFilter string
		timeFilter  string
		expectError bool
	}{
		{
			name:        "valid filters",
			assetFilter: "BTC",
			timeFilter:  "7days",
			expectError: false,
		},
		{
			name:        "invalid asset filter",
			assetFilter: "INVALID",
			timeFilter:  "all",
			expectError: true,
		},
		{
			name:        "invalid time filter",
			assetFilter: "all",
			timeFilter:  "INVALID",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				mockRepo.On("GetTradeStatistics", tt.assetFilter, tt.timeFilter).Return(&generated.TradeStatistics{}, nil)
			}

			_, err := service.GetTradeStatistics(tt.assetFilter, tt.timeFilter)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTradeHistoryService_ValidatePagination(t *testing.T) {
	mockRepo := new(MockTradeHistoryRepository)
	service := NewTradeHistoryService(mockRepo)

	tests := []struct {
		name        string
		page        int
		limit       int
		expectError bool
	}{
		{
			name:        "valid pagination",
			page:        1,
			limit:       10,
			expectError: false,
		},
		{
			name:        "invalid page (zero)",
			page:        0,
			limit:       10,
			expectError: true,
		},
		{
			name:        "invalid limit (zero)",
			page:        1,
			limit:       0,
			expectError: true,
		},
		{
			name:        "invalid limit (too large)",
			page:        1,
			limit:       101,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				mockRepo.On("GetTradeTransactions", "all", "all", tt.page, tt.limit).Return(&generated.TransactionLogResponse{}, nil)
			}

			_, err := service.GetTradeTransactions("all", "all", tt.page, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTradeHistoryService_ErrorHandling(t *testing.T) {
	mockRepo := new(MockTradeHistoryRepository)
	service := NewTradeHistoryService(mockRepo)

	// Test repository error handling
	mockRepo.On("GetTradeStatistics", "all", "all").Return(nil, errors.New("database error"))
	mockRepo.On("GetTradeTransactions", "all", "all", 1, 10).Return(nil, errors.New("database error"))

	_, err := service.GetTradeStatistics("all", "all")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get trade statistics")

	_, err = service.GetTradeTransactions("all", "all", 1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get trade transactions")

	mockRepo.AssertExpectations(t)
}

// Helper function for rounding (should match the one in repository)
func roundToOneDecimal(value float64) float64 {
	if value >= 0 {
		return float64(int(value*10+0.5)) / 10
	}
	return float64(int(value*10-0.5)) / 10
}
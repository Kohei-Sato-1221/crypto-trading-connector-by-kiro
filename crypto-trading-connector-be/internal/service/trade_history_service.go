package service

import (
	"fmt"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/repository"
)

// TradeHistoryService defines the interface for trade history business logic
type TradeHistoryService interface {
	GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error)
	GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error)
}

// TradeHistoryServiceImpl implements TradeHistoryService
type TradeHistoryServiceImpl struct {
	tradeHistoryRepo repository.TradeHistoryRepository
}

// NewTradeHistoryService creates a new trade history service
func NewTradeHistoryService(tradeHistoryRepo repository.TradeHistoryRepository) *TradeHistoryServiceImpl {
	return &TradeHistoryServiceImpl{
		tradeHistoryRepo: tradeHistoryRepo,
	}
}

// GetTradeStatistics retrieves aggregated trade statistics with filtering
func (s *TradeHistoryServiceImpl) GetTradeStatistics(assetFilter, timeFilter string) (*generated.TradeStatistics, error) {
	// Validate filters
	if err := s.validateFilters(assetFilter, timeFilter); err != nil {
		return nil, err
	}

	// Get statistics from repository
	statistics, err := s.tradeHistoryRepo.GetTradeStatistics(assetFilter, timeFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade statistics: %w", err)
	}

	return statistics, nil
}

// GetTradeTransactions retrieves paginated trade transactions with filtering
func (s *TradeHistoryServiceImpl) GetTradeTransactions(assetFilter, timeFilter string, page, limit int) (*generated.TransactionLogResponse, error) {
	// Validate filters
	if err := s.validateFilters(assetFilter, timeFilter); err != nil {
		return nil, err
	}

	// Validate pagination parameters
	if err := s.validatePagination(page, limit); err != nil {
		return nil, err
	}

	// Get transactions from repository
	response, err := s.tradeHistoryRepo.GetTradeTransactions(assetFilter, timeFilter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade transactions: %w", err)
	}

	return response, nil
}

// validateFilters validates asset and time filter parameters
func (s *TradeHistoryServiceImpl) validateFilters(assetFilter, timeFilter string) error {
	// Validate asset filter
	validAssetFilters := map[string]bool{
		"all": true,
		"BTC": true,
		"ETH": true,
	}
	if !validAssetFilters[assetFilter] {
		return fmt.Errorf("invalid asset filter: %s. Valid values are: all, BTC, ETH", assetFilter)
	}

	// Validate time filter
	validTimeFilters := map[string]bool{
		"all":   true,
		"7days": true,
	}
	if !validTimeFilters[timeFilter] {
		return fmt.Errorf("invalid time filter: %s. Valid values are: all, 7days", timeFilter)
	}

	return nil
}

// validatePagination validates pagination parameters
func (s *TradeHistoryServiceImpl) validatePagination(page, limit int) error {
	if page < 1 {
		return fmt.Errorf("invalid page: %d. Page must be >= 1", page)
	}

	if limit < 1 {
		return fmt.Errorf("invalid limit: %d. Limit must be >= 1", limit)
	}

	if limit > 100 {
		return fmt.Errorf("invalid limit: %d. Limit must be <= 100", limit)
	}

	return nil
}
package service

import (
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/model"
	"github.com/crypto-trading-connector/backend/internal/repository"
)

// CryptoService defines the interface for cryptocurrency business logic
type CryptoService interface {
	GetMarketData() (*model.MarketResponse, error)
	GetCryptoByID(id string) (*model.CryptoData, error)
	GetChartData(id string, period string) (*model.ChartResponse, error)
}

// CryptoServiceImpl implements CryptoService
type CryptoServiceImpl struct {
	repo           repository.CryptoRepository
	bitflyerClient client.BitFlyerClient
}

// NewCryptoService creates a new crypto service
func NewCryptoService(repo repository.CryptoRepository, bitflyerClient client.BitFlyerClient) *CryptoServiceImpl {
	return &CryptoServiceImpl{
		repo:           repo,
		bitflyerClient: bitflyerClient,
	}
}

// cryptoConfig holds configuration for each cryptocurrency
type cryptoConfig struct {
	ID          string
	Name        string
	Symbol      string
	Pair        string
	Icon        string
	IconColor   string
	ProductCode string
}

var cryptoConfigs = []cryptoConfig{
	{
		ID:          "bitcoin",
		Name:        "Bitcoin",
		Symbol:      "BTC",
		Pair:        "BTC/JPY",
		Icon:        "₿",
		IconColor:   "#f7931a",
		ProductCode: "BTC_JPY",
	},
	{
		ID:          "ethereum",
		Name:        "Ethereum",
		Symbol:      "ETH",
		Pair:        "ETH/JPY",
		Icon:        "Ξ",
		IconColor:   "#627eea",
		ProductCode: "ETH_JPY",
	},
}

// GetMarketData retrieves market data for all cryptocurrencies
func (s *CryptoServiceImpl) GetMarketData() (*model.MarketResponse, error) {
	var cryptoDataList []model.CryptoData

	for _, config := range cryptoConfigs {
		// Get current price from bitFlyer API
		ticker, err := s.bitflyerClient.GetTicker(config.ProductCode)
		if err != nil {
			return nil, fmt.Errorf("failed to get ticker for %s: %w", config.ProductCode, err)
		}

		// Get chart data from database (last 7 days)
		chartData, err := s.repo.GetDailyAveragePrices(config.ProductCode, 7)
		if err != nil {
			return nil, fmt.Errorf("failed to get chart data for %s: %w", config.ProductCode, err)
		}

		// Calculate change percent
		changePercent := calculateChangePercent(chartData, ticker.Ltp)

		cryptoData := model.CryptoData{
			ID:            config.ID,
			Name:          config.Name,
			Symbol:        config.Symbol,
			Pair:          config.Pair,
			Icon:          config.Icon,
			IconColor:     config.IconColor,
			CurrentPrice:  ticker.Ltp,
			ChangePercent: changePercent,
			ChartData:     chartData,
		}

		cryptoDataList = append(cryptoDataList, cryptoData)
	}

	return &model.MarketResponse{
		Data:      cryptoDataList,
		Timestamp: time.Now().Unix(),
	}, nil
}

// GetCryptoByID retrieves data for a specific cryptocurrency
func (s *CryptoServiceImpl) GetCryptoByID(id string) (*model.CryptoData, error) {
	// Find config for the requested ID
	var config *cryptoConfig
	for _, c := range cryptoConfigs {
		if c.ID == id {
			config = &c
			break
		}
	}

	if config == nil {
		return nil, fmt.Errorf("cryptocurrency not found: %s", id)
	}

	// Get current price from bitFlyer API
	ticker, err := s.bitflyerClient.GetTicker(config.ProductCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker for %s: %w", config.ProductCode, err)
	}

	// Get chart data from database (last 7 days)
	chartData, err := s.repo.GetDailyAveragePrices(config.ProductCode, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data for %s: %w", config.ProductCode, err)
	}

	// Calculate change percent
	changePercent := calculateChangePercent(chartData, ticker.Ltp)

	return &model.CryptoData{
		ID:            config.ID,
		Name:          config.Name,
		Symbol:        config.Symbol,
		Pair:          config.Pair,
		Icon:          config.Icon,
		IconColor:     config.IconColor,
		CurrentPrice:  ticker.Ltp,
		ChangePercent: changePercent,
		ChartData:     chartData,
	}, nil
}

// GetChartData retrieves chart data for a specific cryptocurrency
func (s *CryptoServiceImpl) GetChartData(id string, period string) (*model.ChartResponse, error) {
	// Find config for the requested ID
	var config *cryptoConfig
	for _, c := range cryptoConfigs {
		if c.ID == id {
			config = &c
			break
		}
	}

	if config == nil {
		return nil, fmt.Errorf("cryptocurrency not found: %s", id)
	}

	// Default period is 7d
	if period == "" {
		period = "7d"
	}

	// Convert period to days
	days := periodToDays(period)

	// Get chart data from database
	chartData, err := s.repo.GetDailyAveragePrices(config.ProductCode, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data for %s: %w", config.ProductCode, err)
	}

	return &model.ChartResponse{
		Data:   chartData,
		Period: period,
	}, nil
}

// calculateChangePercent calculates the percentage change from the first chart data point to current price
func calculateChangePercent(chartData []model.ChartDataPoint, currentPrice float64) float64 {
	if len(chartData) == 0 {
		return 0
	}

	firstPrice := chartData[0].Price
	if firstPrice == 0 {
		return 0
	}

	return ((currentPrice - firstPrice) / firstPrice) * 100
}

// periodToDays converts period string to number of days
func periodToDays(period string) int {
	switch period {
	case "24h":
		return 1
	case "7d":
		return 7
	case "30d":
		return 30
	case "1y":
		return 365
	case "all":
		return 3650 // ~10 years
	default:
		return 7
	}
}

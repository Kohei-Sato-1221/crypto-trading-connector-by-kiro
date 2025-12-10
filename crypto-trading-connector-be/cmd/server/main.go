package main

import (
	"fmt"
	"log"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/handler"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/crypto-trading-connector/backend/pkg/database"
	"github.com/crypto-trading-connector/backend/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Initialize exchange client (bitFlyer)
	bitflyerAPIURL := utils.GetEnv("BITFLYER_API_URL", "https://api.bitflyer.com")
	bitflyerAPIKey := utils.GetEnv("BITFLYER_API_KEY", "xxxx")
	bitflyerAPISecret := utils.GetEnv("BITFLYER_API_SECRET", "xxxx")
	
	var exchangeClient client.CryptoExchangeClient
	if bitflyerAPIKey != "" && bitflyerAPISecret != "" {
		exchangeClient = client.NewBitFlyerClientWithAuth(bitflyerAPIURL, bitflyerAPIKey, bitflyerAPISecret)
		log.Println("Exchange client (bitFlyer) initialized with authentication")
	} else {
		exchangeClient = client.NewBitFlyerClient(bitflyerAPIURL)
		log.Println("Exchange client (bitFlyer) initialized without authentication (public API only)")
	}

	// Initialize repositories
	cryptoRepo := repository.NewMySQLCryptoRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	tradeHistoryRepo := repository.NewMySQLTradeHistoryRepository(db)

	// Initialize services
	cryptoService := service.NewCryptoService(cryptoRepo, exchangeClient)
	orderService := service.NewOrderService(exchangeClient, orderRepo)
	tradeHistoryService := service.NewTradeHistoryService(tradeHistoryRepo)

	// Initialize handlers
	cryptoHandler := handler.NewCryptoHandler(cryptoService)
	orderHandler := handler.NewOrderHandler(orderService)
	tradeHistoryHandler := handler.NewTradeHistoryHandler(tradeHistoryService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware - Allow all origins in development
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	api := e.Group("/api/v1")
	{
		// Crypto routes
		crypto := api.Group("/crypto")
		{
			crypto.GET("/market", cryptoHandler.GetMarketData)
			crypto.GET("/:id", cryptoHandler.GetCryptoByID)
			crypto.GET("/:id/chart", cryptoHandler.GetChartData)
		}

		// Order routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/balance", orderHandler.GetBalance)

		// Trade History routes
		tradeHistory := api.Group("/trade-history")
		{
			tradeHistory.GET("/test", func(c echo.Context) error {
				return c.JSON(200, map[string]string{"message": "Trade history routing works"})
			})
			tradeHistory.GET("/statistics", tradeHistoryHandler.GetTradeStatistics)
			tradeHistory.GET("/transactions", tradeHistoryHandler.GetTradeTransactions)
		}

		// Debug route to test if routing works
		api.GET("/debug", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"message": "Debug endpoint works"})
		})
	}

	// Log registered routes
	log.Println("Registered routes:")
	for _, route := range e.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	// Start server on 0.0.0.0 to allow access from other devices
	port := utils.GetEnv("SERVER_PORT", "8080")
	log.Printf("Starting server on 0.0.0.0:%s", port)
	log.Printf("Local: http://localhost:%s", port)
	if err := e.Start(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

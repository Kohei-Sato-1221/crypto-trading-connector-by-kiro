package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/handler"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/crypto-trading-connector/backend/pkg/database"
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

	// Initialize bitFlyer client
	bitflyerAPIURL := getEnv("BITFLYER_API_URL", "https://api.bitflyer.com")
	bitflyerAPIKey := getEnv("BITFLYER_API_KEY", "")
	bitflyerAPISecret := getEnv("BITFLYER_API_SECRET", "")
	
	var bitflyerClient client.BitFlyerClient
	if bitflyerAPIKey != "" && bitflyerAPISecret != "" {
		bitflyerClient = client.NewBitFlyerClientWithAuth(bitflyerAPIURL, bitflyerAPIKey, bitflyerAPISecret)
		log.Println("BitFlyer client initialized with authentication")
	} else {
		bitflyerClient = client.NewBitFlyerClient(bitflyerAPIURL)
		log.Println("BitFlyer client initialized without authentication (public API only)")
	}

	// Initialize repositories
	cryptoRepo := repository.NewMySQLCryptoRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	cryptoService := service.NewCryptoService(cryptoRepo, bitflyerClient)
	orderService := service.NewOrderService(bitflyerClient, orderRepo)

	// Initialize handlers
	cryptoHandler := handler.NewCryptoHandler(cryptoService)
	orderHandler := handler.NewOrderHandler(orderService)

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
	}

	// Start server on 0.0.0.0 to allow access from other devices
	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Starting server on 0.0.0.0:%s", port)
	log.Printf("Local: http://localhost:%s", port)
	if err := e.Start(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

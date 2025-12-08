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
	bitflyerClient := client.NewBitFlyerClient(bitflyerAPIURL)

	// Initialize repository
	cryptoRepo := repository.NewMySQLCryptoRepository(db)

	// Initialize service
	cryptoService := service.NewCryptoService(cryptoRepo, bitflyerClient)

	// Initialize handler
	cryptoHandler := handler.NewCryptoHandler(cryptoService)

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
		crypto := api.Group("/crypto")
		{
			crypto.GET("/market", cryptoHandler.GetMarketData)
			crypto.GET("/:id", cryptoHandler.GetCryptoByID)
			crypto.GET("/:id/chart", cryptoHandler.GetChartData)
		}
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

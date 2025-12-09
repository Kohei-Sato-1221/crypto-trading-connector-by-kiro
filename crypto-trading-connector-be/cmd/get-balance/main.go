package main

import (
	"fmt"
	"log"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get bitFlyer API credentials from environment
	apiURL := utils.GetEnv("BITFLYER_API_URL", "https://api.bitflyer.com")
	apiKey := utils.GetEnv("BITFLYER_API_KEY", "xxxx")
	apiSecret := utils.GetEnv("BITFLYER_API_SECRET", "xxxx")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Error: BITFLYER_API_KEY and BITFLYER_API_SECRET must be set in .env file")
	}

	// Initialize bitFlyer client with authentication
	bitflyerClient := client.NewBitFlyerClientWithAuth(apiURL, apiKey, apiSecret)

	// Get balance
	fmt.Println("Fetching balance from bitFlyer API...")
	balance, err := bitflyerClient.GetBalance()
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	// Display result
	fmt.Printf("\n✅ Success!\n")
	fmt.Printf("Available JPY Balance: ¥%.2f\n", balance)
}
package main

import (
	"fmt"
	"log"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/repository"
	"github.com/crypto-trading-connector/backend/internal/service"
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
	apiKey := utils.GetEnv("BITFLYER_API_KEY", "")
	apiSecret := utils.GetEnv("BITFLYER_API_SECRET", "")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Error: BITFLYER_API_KEY and BITFLYER_API_SECRET must be set in .env file")
	}

	// Initialize bitFlyer client with authentication
	bitflyerClient := client.NewBitFlyerClientWithAuth(apiURL, apiKey, apiSecret)

	// Initialize order repository (nil for this CLI tool since we don't need DB operations)
	orderRepo := repository.NewOrderRepository(nil)

	// Initialize order service
	orderService := service.NewOrderService(bitflyerClient, orderRepo)

	fmt.Println("🔍 Fetching current orders from bitFlyer API...")
	fmt.Println("=" + string(make([]byte, 60)) + "=")

	// Test different scenarios
	scenarios := []struct {
		name  string
		pair  string
		limit int
	}{
		{"All pairs, default limit", "", 10},
		{"BTC_JPY only, limit 5", "BTC_JPY", 5},
		{"ETH_JPY only, limit 5", "ETH_JPY", 5},
	}

	for i, scenario := range scenarios {
		fmt.Printf("\n📊 Scenario %d: %s\n", i+1, scenario.name)
		fmt.Printf("   Pair: %s (empty = all pairs)\n", scenario.pair)
		fmt.Printf("   Limit: %d\n", scenario.limit)

		// Get current orders
		response, err := orderService.GetCurrentOrders(scenario.pair, scenario.limit)
		if err != nil {
			log.Printf("❌ Failed to get current orders: %v\n", err)
			continue
		}

		// Display results
		fmt.Printf("   ✅ Success! Retrieved at: %d\n", response.Timestamp)

		if len(response.BuyOrders) > 0 {
			fmt.Printf("   📈 Buy Orders (%d):\n", len(response.BuyOrders))
			for j, order := range response.BuyOrders {
				fmt.Printf("      %d. ID: %s | Pair: %s | Price: ¥%.0f | Amount: %f | Created: %s\n",
					j+1, order.ID, order.Pair, order.Price, order.Amount, order.CreatedAt)
			}
		} else {
			fmt.Printf("   📈 Buy Orders: None\n")
		}

		if len(response.SellOrders) > 0 {
			fmt.Printf("   📉 Sell Orders (%d):\n", len(response.SellOrders))
			for j, order := range response.SellOrders {
				fmt.Printf("      %d. ID: %s | Pair: %s | Price: ¥%.0f | Amount: %f | Created: %s\n",
					j+1, order.ID, order.Pair, order.Price, order.Amount, order.CreatedAt)
			}
		} else {
			fmt.Printf("   📉 Sell Orders: None\n")
		}

		totalOrders := len(response.BuyOrders) + len(response.SellOrders)
		fmt.Printf("   📊 Total Orders: %d\n", totalOrders)
	}

	fmt.Println("\n" + string(make([]byte, 60)) + "=")
	fmt.Println("✨ Current orders fetch completed!")
	fmt.Println("")
	fmt.Println("💡 You can also test the API endpoint directly:")
	fmt.Println("   curl -X GET \"http://localhost:8090/api/v1/orders/current\"")
	fmt.Println("   curl -X GET \"http://localhost:8090/api/v1/orders/current?pair=BTC_JPY&limit=5\"")
}

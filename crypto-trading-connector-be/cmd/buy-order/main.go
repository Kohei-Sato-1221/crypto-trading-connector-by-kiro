package main

import (
	"fmt"
	"log"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/model"
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

	fmt.Println("🚀 Starting buy order process...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Place orders for both BTC and ETH
	pairs := []struct {
		productCode string
		symbol      string
		minSize     float64
	}{
		{"BTC_JPY", "BTC", 0.001},
		{"ETH_JPY", "ETH", 0.01},
	}

	for _, pair := range pairs {
		fmt.Printf("\n📊 Processing %s order...\n", pair.symbol)

		// Get current price
		ticker, err := bitflyerClient.GetTicker(pair.productCode)
		if err != nil {
			log.Printf("❌ Failed to get %s ticker: %v", pair.symbol, err)
			continue
		}

		currentPrice := ticker.Ltp
		fmt.Printf("   Current Price: ¥%.0f\n", currentPrice)

		// Calculate 97% of current price and round to integer (bitFlyer requires integer prices)
		orderPrice := float64(int(currentPrice * 0.97))
		fmt.Printf("   Order Price (97%%): ¥%.0f\n", orderPrice)

		// Use minimum order size
		orderSize := pair.minSize
		fmt.Printf("   Order Size: %f %s\n", orderSize, pair.symbol)

		estimatedTotal := orderPrice * orderSize
		fmt.Printf("   Estimated Total: ¥%.2f\n", estimatedTotal)

		// Create order request
		orderReq := &model.BitFlyerOrderRequest{
			ProductCode:    pair.productCode,
			ChildOrderType: "LIMIT",
			Side:           "BUY",
			Price:          orderPrice,
			Size:           orderSize,
			TimeInForce:    "GTC", // Good Till Cancelled
		}

		// Submit order
		fmt.Printf("   Submitting order...\n")
		orderResp, err := bitflyerClient.SendOrder(orderReq)
		if err != nil {
			log.Printf("❌ Failed to submit %s order: %v\n", pair.symbol, err)
			continue
		}

		// Display success
		fmt.Printf("   ✅ Order submitted successfully!\n")
		fmt.Printf("   Order ID: %s\n", orderResp.ChildOrderAcceptanceID)
	}

	fmt.Println("\n" + string(make([]byte, 50)) + "=")
	fmt.Println("✨ Buy order process completed!")
}

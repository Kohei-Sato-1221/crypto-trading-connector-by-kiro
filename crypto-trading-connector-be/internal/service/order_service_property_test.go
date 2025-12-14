package service

import (
	"testing"
	"time"

	"github.com/crypto-trading-connector/backend/internal/client"
	"github.com/crypto-trading-connector/backend/internal/model"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestOrderService_OrderSortConsistency_Property tests Property 2: Order sort consistency
// **Feature: current-orders-component, Property 2: Order sort consistency**
// For any order collection, orders are sorted by creation date in descending order (newest first)
// **Validates: Requirements 5.4**
func TestOrderService_OrderSortConsistency_Property(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("orders should be sorted by creation date descending", prop.ForAll(
		func(orders []model.BitflyerChildOrder) bool {
			// Create mock client that returns the generated orders
			mockClient := &client.MockBitFlyerClient{
				GetChildOrdersFunc: func(productCode string, childOrderState string) ([]model.BitflyerChildOrder, error) {
					return orders, nil
				},
			}

			mockRepo := &MockOrderRepository{}
			service := NewOrderService(mockClient, mockRepo)

			// Get current orders
			response, err := service.GetCurrentOrders("", 100)
			if err != nil {
				return false
			}

			// Check buy orders are sorted by creation date descending
			for i := 1; i < len(response.BuyOrders); i++ {
				if response.BuyOrders[i-1].CreatedAt < response.BuyOrders[i].CreatedAt {
					return false
				}
			}

			// Check sell orders are sorted by creation date descending
			for i := 1; i < len(response.SellOrders); i++ {
				if response.SellOrders[i-1].CreatedAt < response.SellOrders[i].CreatedAt {
					return false
				}
			}

			return true
		},
		genBitflyerChildOrders(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genBitflyerChildOrders generates random BitflyerChildOrder slices for property testing
func genBitflyerChildOrders() gopter.Gen {
	return gen.SliceOfN(20, genBitflyerChildOrder()) // Generate up to 20 orders
}

// genBitflyerChildOrder generates a single random BitflyerChildOrder
func genBitflyerChildOrder() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000),               // ID
		genNonEmptyAlphaString(),              // ChildOrderID
		gen.OneConstOf("BTC_JPY", "ETH_JPY"),  // ProductCode
		gen.OneConstOf("BUY", "SELL"),         // Side
		gen.OneConstOf("LIMIT", "MARKET"),     // ChildOrderType
		gen.Float64Range(1000, 20000000),      // Price
		gen.Float64Range(0.001, 10),           // Size
		gen.OneConstOf("ACTIVE", "COMPLETED"), // ChildOrderState
		genRandomTimestamp(),                  // ChildOrderDate
		genNonEmptyAlphaString(),              // ChildOrderAcceptanceID
		gen.Float64Range(0.001, 10),           // OutstandingSize
	).Map(func(values []interface{}) model.BitflyerChildOrder {
		return model.BitflyerChildOrder{
			ID:                     values[0].(int64),
			ChildOrderID:           values[1].(string),
			ProductCode:            values[2].(string),
			Side:                   values[3].(string),
			ChildOrderType:         values[4].(string),
			Price:                  values[5].(float64),
			Size:                   values[6].(float64),
			ChildOrderState:        values[7].(string),
			ChildOrderDate:         values[8].(string),
			ChildOrderAcceptanceID: values[9].(string),
			OutstandingSize:        values[10].(float64),
			CancelSize:             0,
			ExecutedSize:           0,
			TotalCommission:        0,
			AveragePrice:           0,
			ExpireDate:             "",
		}
	})
}

// TestOrderService_APIDataStructure_Property tests Property 11: API Data Structure
// **Feature: current-orders-component, Property 11: API Data Structure**
// For any API response, returned orders include order ID, type, pair, price, amount, creation date
// **Validates: Requirements 5.2**
func TestOrderService_APIDataStructure_Property(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("API response should contain all required fields", prop.ForAll(
		func(orders []model.BitflyerChildOrder) bool {
			// Create mock client that returns the generated orders
			mockClient := &client.MockBitFlyerClient{
				GetChildOrdersFunc: func(productCode string, childOrderState string) ([]model.BitflyerChildOrder, error) {
					return orders, nil
				},
			}

			mockRepo := &MockOrderRepository{}
			service := NewOrderService(mockClient, mockRepo)

			// Get current orders
			response, err := service.GetCurrentOrders("", 100)
			if err != nil {
				return false
			}

			// Check all buy orders have required fields
			for _, order := range response.BuyOrders {
				if order.ID == "" || order.Type == "" || order.Pair == "" ||
					order.Price <= 0 || order.Amount <= 0 || order.CreatedAt == "" {
					return false
				}
				// Verify type is valid
				if order.Type != "buy" && order.Type != "sell" {
					return false
				}
				// Verify pair format
				if order.Pair != "BTC/JPY" && order.Pair != "ETH/JPY" {
					return false
				}
			}

			// Check all sell orders have required fields
			for _, order := range response.SellOrders {
				if order.ID == "" || order.Type == "" || order.Pair == "" ||
					order.Price <= 0 || order.Amount <= 0 || order.CreatedAt == "" {
					return false
				}
				// Verify type is valid
				if order.Type != "buy" && order.Type != "sell" {
					return false
				}
				// Verify pair format
				if order.Pair != "BTC/JPY" && order.Pair != "ETH/JPY" {
					return false
				}
			}

			// Check response structure
			if response.Timestamp <= 0 {
				return false
			}

			return true
		},
		genBitflyerChildOrders(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genNonEmptyAlphaString generates non-empty alphabetic strings
func genNonEmptyAlphaString() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0
	})
}

// genRandomTimestamp generates random timestamps in the format expected by BitFlyer
func genRandomTimestamp() gopter.Gen {
	return gen.Int64Range(1640995200, 1735689600).Map(func(timestamp int64) string {
		// Convert Unix timestamp to BitFlyer format: 2006-01-02T15:04:05.000
		t := time.Unix(timestamp, 0)
		return t.Format("2006-01-02T15:04:05.000")
	})
}

// TestOrderService_AuthenticatedDataRetrieval_Property tests Property 12: Authenticated data retrieval
// **Feature: current-orders-component, Property 12: Authenticated data retrieval**
// For any authenticated API request, system returns only authenticated user's unfilled orders
// **Validates: Requirements 5.1**
func TestOrderService_AuthenticatedDataRetrieval_Property(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("authenticated requests should call exchange with correct parameters", prop.ForAll(
		func(pair string, limit int) bool {
			// Ensure valid inputs
			if limit <= 0 {
				limit = 10
			}

			var expectedProductCode string
			if pair == "BTC/JPY" {
				expectedProductCode = "BTC_JPY"
			} else if pair == "ETH/JPY" {
				expectedProductCode = "ETH_JPY"
			} else {
				expectedProductCode = ""
			}

			// Track if the exchange client was called with correct parameters
			exchangeCalled := false
			correctParameters := false

			mockClient := &client.MockBitFlyerClient{
				GetChildOrdersFunc: func(productCode string, childOrderState string) ([]model.BitflyerChildOrder, error) {
					exchangeCalled = true
					// Verify correct parameters are passed to exchange
					if productCode == expectedProductCode && childOrderState == "ACTIVE" {
						correctParameters = true
					}
					return []model.BitflyerChildOrder{}, nil
				},
			}

			mockRepo := &MockOrderRepository{}
			service := NewOrderService(mockClient, mockRepo)

			// Call GetCurrentOrders
			_, err := service.GetCurrentOrders(pair, limit)
			if err != nil {
				return false
			}

			// Verify exchange was called with correct authentication parameters
			return exchangeCalled && correctParameters
		},
		gen.OneConstOf("", "BTC/JPY", "ETH/JPY"), // pair parameter
		gen.IntRange(1, 50),                      // limit parameter
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

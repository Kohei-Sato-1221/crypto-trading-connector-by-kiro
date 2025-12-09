package client

import (
	"errors"
	"testing"

	"github.com/crypto-trading-connector/backend/internal/model"
)

func TestMockBitFlyerClient_GetBalance(t *testing.T) {
	t.Run("returns default balance when no mock function set", func(t *testing.T) {
		mock := &MockBitFlyerClient{}
		balance, err := mock.GetBalance()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if balance != 1000000.0 {
			t.Errorf("expected balance 1000000.0, got %f", balance)
		}
	})

	t.Run("calls custom mock function when set", func(t *testing.T) {
		expectedBalance := 2500000.0
		mock := &MockBitFlyerClient{
			GetBalanceFunc: func() (float64, error) {
				return expectedBalance, nil
			},
		}

		balance, err := mock.GetBalance()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if balance != expectedBalance {
			t.Errorf("expected balance %f, got %f", expectedBalance, balance)
		}
	})

	t.Run("returns error when mock function returns error", func(t *testing.T) {
		expectedError := errors.New("balance fetch failed")
		mock := &MockBitFlyerClient{
			GetBalanceFunc: func() (float64, error) {
				return 0, expectedError
			},
		}

		balance, err := mock.GetBalance()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if balance != 0 {
			t.Errorf("expected balance 0, got %f", balance)
		}
	})
}

func TestMockBitFlyerClient_SendOrder(t *testing.T) {
	t.Run("returns default response when no mock function set", func(t *testing.T) {
		mock := &MockBitFlyerClient{}
		req := &model.BitFlyerOrderRequest{
			ProductCode:    "BTC_JPY",
			ChildOrderType: "LIMIT",
			Side:           "BUY",
			Price:          10000000,
			Size:           0.001,
		}

		resp, err := mock.SendOrder(req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.ChildOrderAcceptanceID != "TEST_ORDER_123" {
			t.Errorf("expected order ID TEST_ORDER_123, got %s", resp.ChildOrderAcceptanceID)
		}
	})

	t.Run("calls custom mock function when set", func(t *testing.T) {
		expectedOrderID := "CUSTOM_ORDER_456"
		mock := &MockBitFlyerClient{
			SendOrderFunc: func(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
				return &model.BitFlyerOrderResponse{
					ChildOrderAcceptanceID: expectedOrderID,
				}, nil
			},
		}

		req := &model.BitFlyerOrderRequest{
			ProductCode:    "BTC_JPY",
			ChildOrderType: "LIMIT",
			Side:           "BUY",
			Price:          10000000,
			Size:           0.001,
		}

		resp, err := mock.SendOrder(req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.ChildOrderAcceptanceID != expectedOrderID {
			t.Errorf("expected order ID %s, got %s", expectedOrderID, resp.ChildOrderAcceptanceID)
		}
	})

	t.Run("returns error when mock function returns error", func(t *testing.T) {
		expectedError := errors.New("order submission failed")
		mock := &MockBitFlyerClient{
			SendOrderFunc: func(req *model.BitFlyerOrderRequest) (*model.BitFlyerOrderResponse, error) {
				return nil, expectedError
			},
		}

		req := &model.BitFlyerOrderRequest{
			ProductCode:    "BTC_JPY",
			ChildOrderType: "LIMIT",
			Side:           "BUY",
			Price:          10000000,
			Size:           0.001,
		}

		resp, err := mock.SendOrder(req)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})
}

func TestMockBitFlyerClient_GetTicker(t *testing.T) {
	t.Run("returns default ticker when no mock function set", func(t *testing.T) {
		mock := &MockBitFlyerClient{}
		ticker, err := mock.GetTicker("BTC_JPY")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if ticker.ProductCode != "BTC_JPY" {
			t.Errorf("expected product code BTC_JPY, got %s", ticker.ProductCode)
		}
		if ticker.Ltp != 10000000.0 {
			t.Errorf("expected LTP 10000000.0, got %f", ticker.Ltp)
		}
	})
}

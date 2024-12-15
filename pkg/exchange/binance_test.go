package exchange

import (
	"context"
	"os"
	"testing"

	"spot_grid_bot/pkg/types"
)

func TestNewBinanceClient(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		wantErr   bool
	}{
		{
			name:      "Valid credentials",
			apiKey:    "test_api_key",
			apiSecret: "test_api_secret",
			wantErr:   false,
		},
		{
			name:      "Empty credentials",
			apiKey:    "",
			apiSecret: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewBinanceClient(tt.apiKey, tt.apiSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBinanceClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("Expected non-nil client")
			}
		})
	}
}

func TestGetSymbolPrice(t *testing.T) {
	// Skip if no API credentials available
	apiKey := os.Getenv("BINANCE_TEST_API_KEY")
	apiSecret := os.Getenv("BINANCE_TEST_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: no API credentials")
	}

	client, err := NewBinanceClient(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "Valid symbol BTCUSDT",
			symbol:  "BTCUSDT",
			wantErr: false,
		},
		{
			name:    "Invalid symbol",
			symbol:  "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := client.GetSymbolPrice(context.Background(), tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSymbolPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && price <= 0 {
				t.Errorf("Expected positive price, got %v", price)
			}
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	// Skip if no API credentials available
	apiKey := os.Getenv("BINANCE_TEST_API_KEY")
	apiSecret := os.Getenv("BINANCE_TEST_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: no API credentials")
	}

	client, err := NewBinanceClient(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		order   types.Order
		wantErr bool
	}{
		{
			name: "Valid limit buy order",
			order: types.Order{
				Symbol:      "BTCUSDT",
				Side:        "BUY",
				Type:        "LIMIT",
				Quantity:    0.001,
				Price:       20000.0,
				TimeInForce: "GTC",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID, err := client.PlaceOrder(context.Background(), tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && orderID == "" {
				t.Error("Expected non-empty orderID")
			}

			// If order was placed successfully, try to cancel it
			if orderID != "" {
				err = client.CancelOrder(context.Background(), tt.order.Symbol, orderID)
				if err != nil {
					t.Errorf("Failed to cancel order: %v", err)
				}
			}
		})
	}
}

func TestGetBalance(t *testing.T) {
	// Skip if no API credentials available
	apiKey := os.Getenv("BINANCE_TEST_API_KEY")
	apiSecret := os.Getenv("BINANCE_TEST_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: no API credentials")
	}

	client, err := NewBinanceClient(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		asset   string
		wantErr bool
	}{
		{
			name:    "Get USDT balance",
			asset:   "USDT",
			wantErr: false,
		},
		{
			name:    "Get BTC balance",
			asset:   "BTC",
			wantErr: false,
		},
		{
			name:    "Invalid asset",
			asset:   "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := client.GetBalance(context.Background(), tt.asset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && balance < 0 {
				t.Errorf("Expected non-negative balance, got %v", balance)
			}
		})
	}
}

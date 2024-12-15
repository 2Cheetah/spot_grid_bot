package bot

import (
	"context"
	"testing"

	"spot_grid_bot/pkg/types"
)

type mockExchange struct {
	currentPrice float64
	orders       map[string]mockOrder
	orderCounter int // Added to generate unique order IDs
}

type mockOrder struct {
	symbol   string
	side     string
	price    float64
	quantity float64
	orderID  string
}

func (m *mockExchange) GetSymbolPrice(ctx context.Context, symbol string) (float64, error) {
	return m.currentPrice, nil
}

func (m *mockExchange) PlaceOrder(ctx context.Context, order types.Order) (string, error) {
	m.orderCounter++
	orderID := "test_order_" + string(rune(m.orderCounter+'0'))

	m.orders[orderID] = mockOrder{
		symbol:   order.Symbol,
		side:     order.Side,
		price:    order.Price,
		quantity: order.Quantity,
		orderID:  orderID,
	}
	return orderID, nil
}

func (m *mockExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	delete(m.orders, orderID)
	return nil
}

func (m *mockExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	return 1000.0, nil // Mock balance for testing
}

func TestNewGridBot(t *testing.T) {
	tests := []struct {
		name    string
		config  GridBotConfig
		wantErr bool
	}{
		{
			name: "Valid configuration",
			config: GridBotConfig{
				Symbol:     "BTCUSDT",
				LowerPrice: 25000.0,
				UpperPrice: 35000.0,
				GridNum:    5,
				Investment: 1000.0,
			},
			wantErr: false,
		},
		{
			name: "Invalid price range",
			config: GridBotConfig{
				Symbol:     "BTCUSDT",
				LowerPrice: 35000.0, // Lower price > Upper price
				UpperPrice: 25000.0,
				GridNum:    5,
				Investment: 1000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid grid number",
			config: GridBotConfig{
				Symbol:     "BTCUSDT",
				LowerPrice: 25000.0,
				UpperPrice: 35000.0,
				GridNum:    1, // Must be at least 2
				Investment: 1000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid investment amount",
			config: GridBotConfig{
				Symbol:     "BTCUSDT",
				LowerPrice: 25000.0,
				UpperPrice: 35000.0,
				GridNum:    5,
				Investment: 0.0, // Must be positive
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exchange := &mockExchange{
				currentPrice: 30000.0,
				orders:       make(map[string]mockOrder),
			}

			bot, err := NewGridBot(exchange, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGridBot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && bot == nil {
				t.Error("Expected non-nil bot")
			}
		})
	}
}

func TestGridBotInitialOrders(t *testing.T) {
	config := GridBotConfig{
		Symbol:     "BTCUSDT",
		LowerPrice: 25000.0,
		UpperPrice: 35000.0,
		GridNum:    5,
		Investment: 1000.0,
	}

	exchange := &mockExchange{
		currentPrice: 30000.0,
		orders:       make(map[string]mockOrder),
		orderCounter: 0,
	}

	bot, err := NewGridBot(exchange, config)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	ctx := context.Background()
	if err := bot.Start(ctx); err != nil {
		t.Fatalf("Failed to start bot: %v", err)
	}

	// Check if initial orders were placed
	if len(exchange.orders) == 0 {
		t.Errorf("No orders were placed, expected multiple orders")
		return
	}

	// Verify order distribution
	var buyOrders, sellOrders int
	for _, order := range exchange.orders {
		t.Logf("Order: %+v", order) // Add logging for debugging
		if order.side == "BUY" {
			buyOrders++
		} else if order.side == "SELL" {
			sellOrders++
		}
	}

	t.Logf("Buy orders: %d, Sell orders: %d", buyOrders, sellOrders)

	if buyOrders == 0 {
		t.Error("Expected at least one buy order")
	}
	if sellOrders == 0 {
		t.Error("Expected at least one sell order")
	}

	// Test stopping the bot
	if err := bot.Stop(ctx); err != nil {
		t.Errorf("Failed to stop bot: %v", err)
	}

	// Verify all orders were cancelled
	if len(exchange.orders) != 0 {
		t.Errorf("Expected all orders to be cancelled, but %d orders remain", len(exchange.orders))
	}
}

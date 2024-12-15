package bot

import (
	"context"
	"fmt"
	"log"
	"sync"

	"spot_grid_bot/pkg/grid"
	"spot_grid_bot/pkg/types"
)

// GridBotConfig holds the configuration for the grid trading bot
type GridBotConfig struct {
	Symbol     string  // Trading pair symbol (e.g., "BTCUSDT")
	LowerPrice float64 // Lower price bound of the grid
	UpperPrice float64 // Upper price bound of the grid
	GridNum    int     // Number of grid levels
	Investment float64 // Total investment amount in quote currency
}

// Exchange defines the interface for interacting with the exchange
type Exchange interface {
	GetSymbolPrice(ctx context.Context, symbol string) (float64, error)
	PlaceOrder(ctx context.Context, order types.Order) (string, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error
	GetBalance(ctx context.Context, asset string) (float64, error)
}

// GridBot implements a grid trading strategy
type GridBot struct {
	exchange Exchange
	config   GridBotConfig
	levels   []float64
	orders   map[string]types.Order
	mu       sync.RWMutex
	running  bool
}

// NewGridBot creates a new grid trading bot
func NewGridBot(exchange Exchange, config GridBotConfig) (*GridBot, error) {
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// Calculate grid levels
	levels := grid.CalculateGridLevels(config.LowerPrice, config.UpperPrice, config.GridNum)
	if levels == nil {
		return nil, fmt.Errorf("failed to calculate grid levels")
	}

	return &GridBot{
		exchange: exchange,
		config:   config,
		levels:   levels,
		orders:   make(map[string]types.Order),
	}, nil
}

// validateConfig validates the bot configuration
func validateConfig(config GridBotConfig) error {
	if config.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if config.Investment <= 0 {
		return fmt.Errorf("investment must be positive")
	}
	if err := grid.ValidateGridParams(config.LowerPrice, config.UpperPrice, config.GridNum); err != nil {
		return fmt.Errorf("invalid grid parameters: %w", err)
	}
	return nil
}

// Start initializes the grid and starts the trading bot
func (b *GridBot) Start(ctx context.Context) error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return fmt.Errorf("bot is already running")
	}
	b.running = true
	b.mu.Unlock()

	// Get current price
	currentPrice, err := b.exchange.GetSymbolPrice(ctx, b.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}

	// Calculate order quantities
	quantityPerGrid := b.config.Investment / float64(b.config.GridNum*2) // Split investment across grids

	// Place initial orders
	for _, level := range b.levels {
		// Skip levels too close to current price
		if level == currentPrice {
			continue
		}

		var order types.Order
		if level < currentPrice {
			// Place buy order
			order = types.Order{
				Symbol:      b.config.Symbol,
				Side:        "BUY",
				Type:        "LIMIT",
				Quantity:    quantityPerGrid / level, // Convert quote currency to base currency
				Price:       level,
				TimeInForce: "GTC",
			}
		} else {
			// Place sell order
			order = types.Order{
				Symbol:      b.config.Symbol,
				Side:        "SELL",
				Type:        "LIMIT",
				Quantity:    quantityPerGrid / level,
				Price:       level,
				TimeInForce: "GTC",
			}
		}

		orderID, err := b.exchange.PlaceOrder(ctx, order)
		if err != nil {
			log.Printf("Failed to place order at level %v: %v", level, err)
			continue
		}

		b.mu.Lock()
		b.orders[orderID] = order
		b.mu.Unlock()

		log.Printf("Placed %s order at price %.2f, quantity %.8f", order.Side, order.Price, order.Quantity)
	}

	return nil
}

// Stop cancels all open orders and stops the bot
func (b *GridBot) Stop(ctx context.Context) error {
	b.mu.Lock()
	if !b.running {
		b.mu.Unlock()
		return fmt.Errorf("bot is not running")
	}
	b.running = false
	b.mu.Unlock()

	// Cancel all open orders
	var lastError error
	b.mu.RLock()
	for orderID, order := range b.orders {
		if err := b.exchange.CancelOrder(ctx, order.Symbol, orderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", orderID, err)
			lastError = err
		}
	}
	b.mu.RUnlock()

	return lastError
}

// GetStatus returns the current status of the grid bot
func (b *GridBot) GetStatus() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return map[string]interface{}{
		"running":    b.running,
		"symbol":     b.config.Symbol,
		"lowerPrice": b.config.LowerPrice,
		"upperPrice": b.config.UpperPrice,
		"gridNum":    b.config.GridNum,
		"investment": b.config.Investment,
		"openOrders": len(b.orders),
	}
}

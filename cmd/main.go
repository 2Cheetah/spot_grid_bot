package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"spot_grid_bot/pkg/bot"
	"spot_grid_bot/pkg/exchange"
)

func main() {
	// Parse command line flags
	symbol := flag.String("symbol", "BTCUSDT", "Trading pair symbol")
	lowerPrice := flag.Float64("lower", 0, "Lower price bound")
	upperPrice := flag.Float64("upper", 0, "Upper price bound")
	gridNum := flag.Int("grids", 5, "Number of grid levels")
	investment := flag.Float64("investment", 0, "Total investment amount in quote currency")
	flag.Parse()

	// Validate required flags
	if *lowerPrice == 0 || *upperPrice == 0 || *investment == 0 {
		log.Fatal("Lower price, upper price, and investment amount are required")
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("BINANCE_TEST_API_KEY")
	apiSecret := os.Getenv("BINANCE_TEST_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		log.Fatal("BINANCE_TEST_API_KEY and BINANCE_TEST_API_SECRET environment variables are required")
	}

	// Initialize Binance client
	client, err := exchange.NewBinanceClient(apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Failed to create Binance client: %v", err)
	}

	// Create bot configuration
	config := bot.GridBotConfig{
		Symbol:     *symbol,
		LowerPrice: *lowerPrice,
		UpperPrice: *upperPrice,
		GridNum:    *gridNum,
		Investment: *investment,
	}

	// Create grid bot
	gridBot, err := bot.NewGridBot(client, config)
	if err != nil {
		log.Fatalf("Failed to create grid bot: %v", err)
	}

	// Create context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Start the bot
	log.Printf("Starting grid bot for %s...", *symbol)
	log.Printf("Grid configuration: Lower: %.2f, Upper: %.2f, Grids: %d, Investment: %.2f",
		*lowerPrice, *upperPrice, *gridNum, *investment)

	if err := gridBot.Start(ctx); err != nil {
		log.Fatalf("Failed to start grid bot: %v", err)
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Stop the bot
	if err := gridBot.Stop(context.Background()); err != nil {
		log.Printf("Error stopping bot: %v", err)
	}

	log.Println("Bot stopped successfully")
}

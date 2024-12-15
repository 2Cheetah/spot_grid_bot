# Spot Grid Trading Bot

A cryptocurrency grid trading bot for Binance testnet, implemented in Go. The bot creates a grid of buy and sell orders within a specified price range to profit from price oscillations.

## Features

- Grid trading strategy implementation
- Binance testnet support
- Configurable grid parameters
- Real-time price monitoring
- Automatic order management
- Clean shutdown with order cancellation

## Prerequisites

- Go 1.16 or higher
- Binance testnet account
- API key and secret from Binance testnet

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/spot_grid_bot.git
cd spot_grid_bot
```

2. Install dependencies:
```bash
go mod download
```

## Configuration

1. Create a Binance testnet account at https://testnet.binance.vision/
2. Generate API key and secret
3. Set environment variables:
```bash
export BINANCE_TEST_API_KEY="your_api_key"
export BINANCE_TEST_API_SECRET="your_api_secret"
```

## Usage

Run the bot with command line flags:

```bash
go run cmd/main.go \
  -symbol BTCUSDT \
  -lower 25000 \
  -upper 35000 \
  -grids 5 \
  -investment 1000
```

Parameters:
- `-symbol`: Trading pair (default: BTCUSDT)
- `-lower`: Lower price bound of the grid
- `-upper`: Upper price bound of the grid
- `-grids`: Number of grid levels (minimum: 2)
- `-investment`: Total investment amount in quote currency

## Architecture

The project is organized into several packages:

- `pkg/grid`: Grid calculation logic
- `pkg/exchange`: Binance API client wrapper
- `pkg/bot`: Grid trading bot implementation
- `pkg/types`: Common type definitions
- `cmd`: Main application entry point

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

## Example

Here's how the grid bot works:

1. If you set a grid between $25,000 and $35,000 with 5 levels and $1000 investment:
   - The bot will create grid levels at: $25,000, $27,500, $30,000, $32,500, and $35,000
   - Investment will be split across grid levels
   - Buy orders will be placed below current price
   - Sell orders will be placed above current price

2. When orders are filled:
   - If a buy order is filled, a sell order is placed above
   - If a sell order is filled, a buy order is placed below
   - This creates a continuous trading cycle

## Safety Notes

- This bot is for educational purposes
- Test thoroughly on testnet before using real funds
- Start with small amounts to understand the behavior
- Monitor the bot's performance regularly

## License

MIT License

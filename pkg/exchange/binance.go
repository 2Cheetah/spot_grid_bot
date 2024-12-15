package exchange

import (
	"context"
	"fmt"
	"strconv"

	"spot_grid_bot/pkg/types"

	"github.com/adshao/go-binance/v2"
)

// BinanceClient wraps the Binance API client with testnet support
type BinanceClient struct {
	client *binance.Client
}

// NewBinanceClient creates a new Binance client configured for testnet
func NewBinanceClient(apiKey, apiSecret string) (*BinanceClient, error) {
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("API key and secret are required")
	}

	// Use testnet
	binance.UseTestnet = true
	client := binance.NewClient(apiKey, apiSecret)

	return &BinanceClient{
		client: client,
	}, nil
}

// GetSymbolPrice gets the current price for a symbol
func (c *BinanceClient) GetSymbolPrice(ctx context.Context, symbol string) (float64, error) {
	prices, err := c.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get price: %w", err)
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("no price found for symbol %s", symbol)
	}

	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// PlaceOrder places a new order
func (c *BinanceClient) PlaceOrder(ctx context.Context, order types.Order) (string, error) {
	service := c.client.NewCreateOrderService().
		Symbol(order.Symbol).
		Side(binance.SideType(order.Side)).
		Type(binance.OrderType(order.Type))

	// Convert quantity to string with appropriate precision
	quantityStr := strconv.FormatFloat(order.Quantity, 'f', 8, 64)
	service.Quantity(quantityStr)

	if order.Type == "LIMIT" {
		priceStr := strconv.FormatFloat(order.Price, 'f', 8, 64)
		service.TimeInForce(binance.TimeInForceType(order.TimeInForce)).
			Price(priceStr)
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to place order: %w", err)
	}

	return strconv.FormatInt(resp.OrderID, 10), nil
}

// CancelOrder cancels an existing order
func (c *BinanceClient) CancelOrder(ctx context.Context, symbol, orderID string) error {
	// Convert string orderID to int64
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	_, err = c.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderIDInt).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// GetBalance gets the balance for a specific asset
func (c *BinanceClient) GetBalance(ctx context.Context, asset string) (float64, error) {
	account, err := c.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}

	for _, balance := range account.Balances {
		if balance.Asset == asset {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse balance: %w", err)
			}
			return free, nil
		}
	}

	return 0, fmt.Errorf("asset %s not found", asset)
}

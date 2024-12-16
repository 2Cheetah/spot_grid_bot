package main

import (
	"os"

	binance_connector "github.com/binance/binance-connector-go"
)

func main() {
	apiKey := os.Getenv("BINANCE_TEST_API_KEY")
	secretKey := os.Getenv("BINANCE_TEST_API_SECRET")
	baseURL := "https://testnet.binance.vision"

	// Initialise the client
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// // Get balance of an asset
	// balance, err := client.NewGetAccountService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // fmt.Println(binance_connector.PrettyPrint(balance.Balances))
	// for _, asset := range balance.Balances {
	// 	if asset.Asset == "USDT" {
	// 		fmt.Println(binance_connector.PrettyPrint(asset))
	// 	}
	// }

	// // Get BNBUSDT ticker
	// ticker, err := client.NewTickerService().Symbol("BNBUSDT").Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(binance_connector.PrettyPrint(ticker))

	// // Create new order
	// newOrder, err := client.NewCreateOrderService().Symbol("BNBUSDT").
	// 	Side("BUY").Type("MARKET").Quantity(1).
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(binance_connector.PrettyPrint(newOrder))
}

package types

// Order represents a trading order
type Order struct {
	Symbol      string
	Side        string // BUY or SELL
	Type        string // LIMIT or MARKET
	Quantity    float64
	Price       float64
	TimeInForce string // GTC, IOC, FOK
}

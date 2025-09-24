package dto

type Transaction struct {
	ID         string
	CustomerID string
	FromAcc    string
	ToAcc      string
	Amount     float64
	Currency   string
	Metadata   map[string]any
}

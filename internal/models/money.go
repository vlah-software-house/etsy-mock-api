package models

type Money struct {
	Amount       int    `json:"amount"`
	Divisor      int    `json:"divisor"`
	CurrencyCode string `json:"currency_code"`
}

func USD(cents int) Money {
	return Money{Amount: cents, Divisor: 100, CurrencyCode: "USD"}
}

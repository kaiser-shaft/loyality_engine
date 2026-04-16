package models

import "errors"

var (
	ErrInvalidItem     = errors.New("invalid item data: price and quantity must be positive")
	ErrEmptyCart       = errors.New("items list is empty")
	ErrInternalStorage = errors.New("internal storage error")
)

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Quantity int     `json:"quantity"`
}

type AppliedDiscount struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Calculation struct {
	OriginalTotal float64           `json:"original_total"`
	Discounts     []AppliedDiscount `json:"discounts"`
	FinalTotal    float64           `json:"final_total"`
}

package finance

import (
	"fmt"
)

// Money represents a monetary value in a specific currency.
// It uses integer math (minor units) to avoid floating point errors.
type Money struct {
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"` // ISO 4217 code
	Scale       int    `json:"scale"`    // e.g. 2 for USD/EUR, 8 for BTC
}

// NewMoney creates a new Money instance.
func NewMoney(amount int64, currency string) Money {
	// Default scale lookup could go here, for now assuming 2 for fiat
	scale := 2
	if currency == "BTC" || currency == "ETH" {
		scale = 8
	}
	return Money{
		AmountMinor: amount,
		Currency:    currency,
		Scale:       scale,
	}
}

// Add adds two Money amounts. Returns error on currency mismatch.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	if m.Scale != other.Scale {
		return Money{}, fmt.Errorf("scale mismatch: %d vs %d", m.Scale, other.Scale)
	}
	return Money{
		AmountMinor: m.AmountMinor + other.AmountMinor,
		Currency:    m.Currency,
		Scale:       m.Scale,
	}, nil
}

// Sub subtracts other Money from m. Returns error on currency mismatch.
func (m Money) Sub(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{
		AmountMinor: m.AmountMinor - other.AmountMinor,
		Currency:    m.Currency,
		Scale:       m.Scale,
	}, nil
}

// IsZero returns true if the amount is 0.
func (m Money) IsZero() bool {
	return m.AmountMinor == 0
}

// IsPositive returns true if the amount is > 0.
func (m Money) IsPositive() bool {
	return m.AmountMinor > 0
}

// IsNegative returns true if the amount is < 0.
func (m Money) IsNegative() bool {
	return m.AmountMinor < 0
}

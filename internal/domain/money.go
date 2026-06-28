package domain

import (
	"errors"
	"strings"
)

// Money holds amount in minor units scaled ×10000 (4 decimal places).
// ponytail: int64 scale — upgrade to big.Int only if overflow becomes measurable.
type Money struct {
	amount   int64  // e.g. 1.0000 USD = 10000
	currency string // ISO 4217
}

var (
	ErrInvalidCurrency  = errors.New("currency must be a non-empty ISO 4217 code")
	ErrCurrencyMismatch = errors.New("currency mismatch: cannot add different currencies")
)

func NewMoney(amount int64, currency string) (Money, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Money{}, ErrInvalidCurrency
	}
	return Money{amount: amount, currency: currency}, nil
}

func MustMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }

func (m Money) IsZero() bool { return m.amount == 0 }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

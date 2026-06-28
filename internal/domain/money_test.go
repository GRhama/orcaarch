package domain

import "testing"

func TestNewMoney_valid(t *testing.T) {
	m, err := NewMoney(10000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if m.Amount() != 10000 || m.Currency() != "USD" {
		t.Fatalf("got %d %s", m.Amount(), m.Currency())
	}
}

func TestNewMoney_zero_amount(t *testing.T) {
	m, err := NewMoney(0, "BRL")
	if err != nil {
		t.Fatal(err)
	}
	if !m.IsZero() {
		t.Fatal("expected zero")
	}
}

func TestNewMoney_empty_currency(t *testing.T) {
	_, err := NewMoney(100, "")
	if err != ErrInvalidCurrency {
		t.Fatalf("expected ErrInvalidCurrency, got %v", err)
	}
}

func TestNewMoney_normalizes_currency(t *testing.T) {
	m, err := NewMoney(5000, " usd ")
	if err != nil {
		t.Fatal(err)
	}
	if m.Currency() != "USD" {
		t.Fatalf("expected USD, got %s", m.Currency())
	}
}

func TestMoney_Add_SameCurrency(t *testing.T) {
	a := MustMoney(10000, "USD")
	b := MustMoney(5000, "USD")
	got, err := a.Add(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Amount() != 15000 || got.Currency() != "USD" {
		t.Fatalf("got %d %s", got.Amount(), got.Currency())
	}
}

func TestMoney_Add_CurrencyMismatch(t *testing.T) {
	a := MustMoney(10000, "USD")
	b := MustMoney(5000, "EUR")
	_, err := a.Add(b)
	if err != ErrCurrencyMismatch {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}

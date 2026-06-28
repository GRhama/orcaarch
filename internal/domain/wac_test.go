package domain

import (
	"errors"
	"testing"
)

// OldStock=0, WAC = NewLandedCost / NewTons
func TestCalculateWAC_FirstShipment(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    0,
		OldWACPerMilliTon:    MustMoney(0, "USD"),
		NewShipmentMilliTons: 2000,
		NewLandedCostUSD:     MustMoney(10_000_000, "USD"), // 1000.0000 USD
	}
	got, err := CalculateWAC(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 10_000_000 / 2000 = 5000
	if got.Amount() != 5000 {
		t.Errorf("amount: got %d, want 5000", got.Amount())
	}
	if got.Currency() != "USD" {
		t.Errorf("currency: got %s, want USD", got.Currency())
	}
}

// Blended WAC between existing stock and new shipment at different cost.
func TestCalculateWAC_SubsequentShipment(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    3000,
		OldWACPerMilliTon:    MustMoney(4000, "USD"), // 0.4000 USD/milli-ton
		NewShipmentMilliTons: 2000,
		NewLandedCostUSD:     MustMoney(10_000_000, "USD"), // 1000.0000 USD (= 5000/milli-ton)
	}
	// numerator = 3000×4000 + 10_000_000 = 12_000_000 + 10_000_000 = 22_000_000
	// denominator = 3000 + 2000 = 5000
	// newWAC = 22_000_000 / 5000 = 4400
	got, err := CalculateWAC(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount() != 4400 {
		t.Errorf("amount: got %d, want 4400", got.Amount())
	}
}

func TestCalculateWAC_ZeroNewShipmentWeight(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    1000,
		OldWACPerMilliTon:    MustMoney(5000, "USD"),
		NewShipmentMilliTons: 0,
		NewLandedCostUSD:     MustMoney(5_000_000, "USD"),
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACInvalidNewShipmentWeight) {
		t.Errorf("got %v, want ErrWACInvalidNewShipmentWeight", err)
	}
}

func TestCalculateWAC_NegativeNewShipmentWeight(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    1000,
		OldWACPerMilliTon:    MustMoney(5000, "USD"),
		NewShipmentMilliTons: -1,
		NewLandedCostUSD:     MustMoney(5_000_000, "USD"),
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACInvalidNewShipmentWeight) {
		t.Errorf("got %v, want ErrWACInvalidNewShipmentWeight", err)
	}
}

func TestCalculateWAC_NegativeOldStock(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    -1,
		OldWACPerMilliTon:    MustMoney(5000, "USD"),
		NewShipmentMilliTons: 1000,
		NewLandedCostUSD:     MustMoney(5_000_000, "USD"),
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACNegativeOldStock) {
		t.Errorf("got %v, want ErrWACNegativeOldStock", err)
	}
}

func TestCalculateWAC_CurrencyNotUSD_NewLandedCost(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    0,
		OldWACPerMilliTon:    MustMoney(0, "USD"),
		NewShipmentMilliTons: 1000,
		NewLandedCostUSD:     MustMoney(5_000_000, "BRL"),
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACCurrencyNotUSD) {
		t.Errorf("got %v, want ErrWACCurrencyNotUSD", err)
	}
}

func TestCalculateWAC_CurrencyNotUSD_NewLandedCostEmpty(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    0,
		OldWACPerMilliTon:    MustMoney(0, "USD"),
		NewShipmentMilliTons: 1000,
		NewLandedCostUSD:     Money{}, // currency="" must be rejected
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACCurrencyNotUSD) {
		t.Errorf("got %v, want ErrWACCurrencyNotUSD", err)
	}
}

// Integer division truncates toward zero; no rounding.
func TestCalculateWAC_IntegerTruncation(t *testing.T) {
	// numerator = 1×10000 + 10002 = 20002; denominator = 1+2 = 3; 20002/3 = 6667 remainder 1 (truncated)
	input := WACInput{
		OldStockMilliTons:    1,
		OldWACPerMilliTon:    MustMoney(10000, "USD"),
		NewShipmentMilliTons: 2,
		NewLandedCostUSD:     MustMoney(10002, "USD"),
	}
	got, err := CalculateWAC(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount() != 6667 {
		t.Errorf("amount: got %d, want 6667 (truncated)", got.Amount())
	}
}

// Money{} zero-value has currency="" — must be rejected, not treated as USD.
func TestCalculateWAC_CurrencyNotUSD_OldWAC(t *testing.T) {
	input := WACInput{
		OldStockMilliTons:    1000,
		OldWACPerMilliTon:    Money{},
		NewShipmentMilliTons: 1000,
		NewLandedCostUSD:     MustMoney(5_000_000, "USD"),
	}
	_, err := CalculateWAC(input)
	if !errors.Is(err, ErrWACCurrencyNotUSD) {
		t.Errorf("got %v, want ErrWACCurrencyNotUSD", err)
	}
}

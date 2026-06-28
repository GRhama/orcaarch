package domain

import "testing"

func TestCalculateLandedCost_AllComponents(t *testing.T) {
	input := LandedCostInput{
		FactoryCostUSD:   MustMoney(100_0000, "USD"), // 100.0000
		ActualFreightUSD: MustMoney(20_0000, "USD"),  // 20.0000
		CustomsDutiesUSD: MustMoney(15_0000, "USD"),  // 15.0000
		AncillaryFeesUSD: MustMoney(5_0000, "USD"),   // 5.0000
		InsuranceCostUSD: MustMoney(2_0000, "USD"),   // 2.0000
	}
	got, err := CalculateLandedCost(input)
	if err != nil {
		t.Fatal(err)
	}
	// 100 + 20 + 15 + 5 + 2 = 142.0000 -> 1_420_000
	if got.Amount() != 142_0000 || got.Currency() != "USD" {
		t.Fatalf("got %d %s", got.Amount(), got.Currency())
	}
}

func TestCalculateLandedCost_ZeroInsurance(t *testing.T) {
	input := LandedCostInput{
		FactoryCostUSD:   MustMoney(50_0000, "USD"),
		ActualFreightUSD: MustMoney(10_0000, "USD"),
		CustomsDutiesUSD: MustMoney(5_0000, "USD"),
		AncillaryFeesUSD: MustMoney(2_0000, "USD"),
		InsuranceCostUSD: MustMoney(0, "USD"),
	}
	got, err := CalculateLandedCost(input)
	if err != nil {
		t.Fatal(err)
	}
	// 50 + 10 + 5 + 2 = 67.0000 -> 670_000
	if got.Amount() != 67_0000 {
		t.Fatalf("got %d", got.Amount())
	}
}

func TestCalculateLandedCost_ZeroFreight(t *testing.T) {
	input := LandedCostInput{
		FactoryCostUSD:   MustMoney(80_0000, "USD"),
		ActualFreightUSD: MustMoney(0, "USD"),
		CustomsDutiesUSD: MustMoney(8_0000, "USD"),
		AncillaryFeesUSD: MustMoney(1_0000, "USD"),
		InsuranceCostUSD: MustMoney(0, "USD"),
	}
	got, err := CalculateLandedCost(input)
	if err != nil {
		t.Fatal(err)
	}
	// 80 + 0 + 8 + 1 + 0 = 89.0000 -> 890_000
	if got.Amount() != 89_0000 {
		t.Fatalf("got %d", got.Amount())
	}
}

func TestCalculateLandedCost_CurrencyMismatch(t *testing.T) {
	input := LandedCostInput{
		FactoryCostUSD:   MustMoney(100_0000, "USD"),
		ActualFreightUSD: MustMoney(20_0000, "EUR"), // wrong currency
		CustomsDutiesUSD: MustMoney(5_0000, "USD"),
		AncillaryFeesUSD: MustMoney(2_0000, "USD"),
		InsuranceCostUSD: MustMoney(0, "USD"),
	}
	_, err := CalculateLandedCost(input)
	if err != ErrLandedCostCurrencyNotUSD {
		t.Fatalf("expected ErrLandedCostCurrencyNotUSD, got %v", err)
	}
}

func TestCalculateLandedCost_AllZero(t *testing.T) {
	input := LandedCostInput{
		FactoryCostUSD:   MustMoney(0, "USD"),
		ActualFreightUSD: MustMoney(0, "USD"),
		CustomsDutiesUSD: MustMoney(0, "USD"),
		AncillaryFeesUSD: MustMoney(0, "USD"),
		InsuranceCostUSD: MustMoney(0, "USD"),
	}
	got, err := CalculateLandedCost(input)
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsZero() || got.Currency() != "USD" {
		t.Fatalf("got %d %s", got.Amount(), got.Currency())
	}
}

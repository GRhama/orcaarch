package domain

import "errors"

var ErrLandedCostCurrencyNotUSD = errors.New("all landed cost components must be in USD")

// LandedCostInput holds all cost components, each pre-converted to USD by the caller.
type LandedCostInput struct {
	FactoryCostUSD   Money
	ActualFreightUSD Money
	CustomsDutiesUSD Money
	AncillaryFeesUSD Money
	InsuranceCostUSD Money // zero if not applicable
}

// CalculateLandedCost returns the sum of all cost components.
// All components must carry currency "USD"; returns ErrLandedCostCurrencyNotUSD otherwise.
func CalculateLandedCost(input LandedCostInput) (Money, error) {
	components := []Money{
		input.FactoryCostUSD,
		input.ActualFreightUSD,
		input.CustomsDutiesUSD,
		input.AncillaryFeesUSD,
		input.InsuranceCostUSD,
	}
	for _, c := range components {
		if c.Currency() != "USD" {
			return Money{}, ErrLandedCostCurrencyNotUSD
		}
	}
	total := MustMoney(0, "USD")
	for _, c := range components {
		var err error
		total, err = total.Add(c)
		if err != nil {
			return Money{}, err
		}
	}
	return total, nil
}

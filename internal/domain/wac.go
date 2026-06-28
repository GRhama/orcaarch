package domain

import "errors"

var (
	ErrWACInvalidNewShipmentWeight = errors.New("new_shipment_milli_tons must be positive")
	ErrWACNegativeOldStock         = errors.New("old_stock_milli_tons must be non-negative")
	ErrWACCurrencyNotUSD           = errors.New("WAC money components must be in USD")
)

// WACInput holds the inputs required to compute a new Weighted Average Cost.
// OldStockMilliTons may be zero (first receipt). NewShipmentMilliTons must be > 0.
type WACInput struct {
	OldStockMilliTons    int64 // >= 0
	OldWACPerMilliTon    Money // USD/milli-ton ×10000; zero-value Money{} rejected
	NewShipmentMilliTons int64 // > 0
	NewLandedCostUSD     Money // USD
}

// CalculateWAC returns the new blended WAC (USD/milli-ton ×10000).
// Formula: ((OldStockMilliTons × OldWAC.Amount()) + NewLandedCost.Amount())
//
//	/ (OldStockMilliTons + NewShipmentMilliTons)
func CalculateWAC(input WACInput) (Money, error) {
	if input.NewShipmentMilliTons <= 0 {
		return Money{}, ErrWACInvalidNewShipmentWeight
	}
	if input.OldStockMilliTons < 0 {
		return Money{}, ErrWACNegativeOldStock
	}
	if input.OldWACPerMilliTon.Currency() != "USD" {
		return Money{}, ErrWACCurrencyNotUSD
	}
	if input.NewLandedCostUSD.Currency() != "USD" {
		return Money{}, ErrWACCurrencyNotUSD
	}

	// ponytail: int64 overflow possible for extreme old stock × WAC values; safe for realistic logistics volumes.
	numerator := input.OldStockMilliTons*input.OldWACPerMilliTon.Amount() + input.NewLandedCostUSD.Amount()
	denominator := input.OldStockMilliTons + input.NewShipmentMilliTons

	return NewMoney(numerator/denominator, "USD")
}

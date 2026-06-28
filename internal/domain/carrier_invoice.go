package domain

import (
	"errors"
	"time"
)

// CarrierInvoice is the real charge from the carrier.
// ActualFreightCurrency.Currency = original invoice currency (may differ from USD).
// CustomsDutiesLocal.Currency    = local country currency.
// InsuranceCostUSD.Currency      = "USD".
// AncillaryFeesUSD.Currency      = "USD".
type CarrierInvoice struct {
	InvoiceID             string
	InvoiceDate           time.Time
	TrackingNumber        string
	ActualFreightCurrency Money
	CustomsDutiesLocal    Money
	InsuranceCostUSD      Money
	AncillaryFeesUSD      Money
}

var (
	ErrEmptyInvoiceID              = errors.New("invoice_id must not be empty")
	ErrInsuranceMustBeUSD          = errors.New("insurance_cost must be in USD")
	ErrAncillaryFeesMustBeUSD      = errors.New("ancillary_fees must be in USD")
	ErrEmptyActualFreightCurrency  = errors.New("actual_freight_currency must have a currency set")
	ErrEmptyCustomsDutiesCurrency  = errors.New("customs_duties_local must have a currency set")
)

func (c CarrierInvoice) Validate() error {
	if c.InvoiceID == "" {
		return ErrEmptyInvoiceID
	}
	if c.TrackingNumber == "" {
		return ErrEmptyTrackingNumber
	}
	if c.ActualFreightCurrency.Currency() == "" {
		return ErrEmptyActualFreightCurrency
	}
	if c.CustomsDutiesLocal.Currency() == "" {
		return ErrEmptyCustomsDutiesCurrency
	}
	if c.InsuranceCostUSD.Currency() != "USD" {
		return ErrInsuranceMustBeUSD
	}
	if c.AncillaryFeesUSD.Currency() != "USD" {
		return ErrAncillaryFeesMustBeUSD
	}
	return nil
}

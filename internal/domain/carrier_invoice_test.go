package domain

import (
	"testing"
	"time"
)

func validInvoice() CarrierInvoice {
	return CarrierInvoice{
		InvoiceID:             "INV-001",
		InvoiceDate:           time.Now(),
		TrackingNumber:        "TRK-ABC123",
		ActualFreightCurrency: MustMoney(30000000, "EUR"), // EUR 3000.0000
		CustomsDutiesLocal:    MustMoney(5000000, "BRL"),  // BRL 500.0000
		InsuranceCostUSD:      MustMoney(1000000, "USD"),  // USD 100.0000
		AncillaryFeesUSD:      MustMoney(500000, "USD"),   // USD 50.0000
	}
}

func TestCarrierInvoice_valid(t *testing.T) {
	if err := validInvoice().Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestCarrierInvoice_empty_invoice_id(t *testing.T) {
	c := validInvoice()
	c.InvoiceID = ""
	if err := c.Validate(); err != ErrEmptyInvoiceID {
		t.Fatalf("expected ErrEmptyInvoiceID, got %v", err)
	}
}

func TestCarrierInvoice_empty_tracking_number(t *testing.T) {
	c := validInvoice()
	c.TrackingNumber = ""
	if err := c.Validate(); err != ErrEmptyTrackingNumber {
		t.Fatalf("expected ErrEmptyTrackingNumber, got %v", err)
	}
}

func TestCarrierInvoice_insurance_not_usd(t *testing.T) {
	c := validInvoice()
	c.InsuranceCostUSD = MustMoney(1000000, "EUR")
	if err := c.Validate(); err != ErrInsuranceMustBeUSD {
		t.Fatalf("expected ErrInsuranceMustBeUSD, got %v", err)
	}
}

func TestCarrierInvoice_ancillary_not_usd(t *testing.T) {
	c := validInvoice()
	c.AncillaryFeesUSD = MustMoney(500000, "BRL")
	if err := c.Validate(); err != ErrAncillaryFeesMustBeUSD {
		t.Fatalf("expected ErrAncillaryFeesMustBeUSD, got %v", err)
	}
}

func TestCarrierInvoice_empty_actual_freight_currency(t *testing.T) {
	c := validInvoice()
	c.ActualFreightCurrency = Money{} // zero value: currency == ""
	if err := c.Validate(); err != ErrEmptyActualFreightCurrency {
		t.Fatalf("expected ErrEmptyActualFreightCurrency, got %v", err)
	}
}

func TestCarrierInvoice_empty_customs_duties_currency(t *testing.T) {
	c := validInvoice()
	c.CustomsDutiesLocal = Money{} // zero value: currency == ""
	if err := c.Validate(); err != ErrEmptyCustomsDutiesCurrency {
		t.Fatalf("expected ErrEmptyCustomsDutiesCurrency, got %v", err)
	}
}

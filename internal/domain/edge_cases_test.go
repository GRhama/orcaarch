package domain

import "testing"

func TestProvisionalLandedCostInput_PreservesFreight(t *testing.T) {
	freight := MustMoney(5_000_0000, "USD")
	erp := ERPShipment{TrackingNumber: "TRK-001", BookingID: "BK-001", TotalWeightMilliTons: 1000, EstimatedFreightUSD: freight}
	input := ProvisionalLandedCostInput(erp)
	if input.ActualFreightUSD.Amount() != freight.Amount() || input.ActualFreightUSD.Currency() != "USD" {
		t.Errorf("ActualFreightUSD = %v, want %v", input.ActualFreightUSD, freight)
	}
}

func TestProvisionalLandedCostInput_OtherComponentsAreUSD(t *testing.T) {
	erp := ERPShipment{TrackingNumber: "TRK-001", BookingID: "BK-001", TotalWeightMilliTons: 1000, EstimatedFreightUSD: MustMoney(1_0000, "USD")}
	input := ProvisionalLandedCostInput(erp)
	for _, m := range []Money{input.FactoryCostUSD, input.CustomsDutiesUSD, input.AncillaryFeesUSD, input.InsuranceCostUSD} {
		if m.Currency() != "USD" {
			t.Errorf("expected USD, got %q", m.Currency())
		}
	}
}

func TestProvisionalLandedCostInput_PassesIntoCalculateLandedCost(t *testing.T) {
	erp := ERPShipment{TrackingNumber: "TRK-001", BookingID: "BK-001", TotalWeightMilliTons: 1000, EstimatedFreightUSD: MustMoney(3_000_0000, "USD")}
	_, err := CalculateLandedCost(ProvisionalLandedCostInput(erp))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProvisionalLandedCostInput_ZeroFreight(t *testing.T) {
	erp := ERPShipment{TrackingNumber: "TRK-001", BookingID: "BK-001", TotalWeightMilliTons: 1000, EstimatedFreightUSD: MustMoney(0, "USD")}
	result, err := CalculateLandedCost(ProvisionalLandedCostInput(erp))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Amount() != 0 || result.Currency() != "USD" {
		t.Errorf("want Money(0,USD), got %v", result)
	}
}

func TestNewQuarantine_PreservesInvoiceAndReason(t *testing.T) {
	invoice := CarrierInvoice{
		InvoiceID:             "INV-001",
		TrackingNumber:        "TRK-999",
		ActualFreightCurrency: MustMoney(2_000_0000, "USD"),
		CustomsDutiesLocal:    MustMoney(100_0000, "BRL"),
		InsuranceCostUSD:      MustMoney(0, "USD"),
		AncillaryFeesUSD:      MustMoney(0, "USD"),
	}
	q := NewQuarantine(invoice)
	if q.Invoice.InvoiceID != invoice.InvoiceID {
		t.Errorf("Invoice not preserved: got %v", q.Invoice.InvoiceID)
	}
	if q.Reason != ReasonMissingERPShipment {
		t.Errorf("Reason = %q, want %q", q.Reason, ReasonMissingERPShipment)
	}
}

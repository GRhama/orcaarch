package domain

import (
	"math"
	"testing"
)

// helpers — minimal structs, Reconcile does not call Validate
func erpUSD(tracking string, amountScaled int64) ERPShipment {
	return ERPShipment{
		BookingID:            "BK-" + tracking,
		TotalWeightMilliTons: 1000,
		TrackingNumber:       tracking,
		EstimatedFreightUSD:  MustMoney(amountScaled, "USD"),
	}
}

func carrierUSD(tracking string, amountScaled int64) CarrierInvoice {
	return CarrierInvoice{
		InvoiceID:             "INV-" + tracking,
		TrackingNumber:        tracking,
		ActualFreightCurrency: MustMoney(amountScaled, "USD"),
		InsuranceCostUSD:      MustMoney(0, "USD"),
		AncillaryFeesUSD:      MustMoney(0, "USD"),
		CustomsDutiesLocal:    MustMoney(0, "USD"),
	}
}

func carrierFX(tracking string, amountScaled int64, currency string) CarrierInvoice {
	return CarrierInvoice{
		InvoiceID:             "INV-" + tracking,
		TrackingNumber:        tracking,
		ActualFreightCurrency: MustMoney(amountScaled, currency),
		InsuranceCostUSD:      MustMoney(0, "USD"),
		AncillaryFeesUSD:      MustMoney(0, "USD"),
		CustomsDutiesLocal:    MustMoney(0, "USD"),
	}
}

func findResult(results []ReconciliationResult, tracking string) (ReconciliationResult, bool) {
	for _, r := range results {
		if r.TrackingNumber == tracking {
			return r, true
		}
	}
	return ReconciliationResult{}, false
}

// 1 — diff < tolerance → MATCHED
func TestReconcile_Matched(t *testing.T) {
	// 100 USD estimated, 100.50 USD actual → 50 bp diff, tolerance 100 bp
	erps := []ERPShipment{erpUSD("TRK-1", 1_000_000)}
	carriers := []CarrierInvoice{carrierUSD("TRK-1", 1_005_000)}

	results := Reconcile(erps, carriers, 100)

	if len(results) != 1 {
		t.Fatalf("want 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Status != StatusMatched {
		t.Errorf("want MATCHED, got %s", r.Status)
	}
	if r.Reason != ReasonWithinTolerance {
		t.Errorf("want %s, got %s", ReasonWithinTolerance, r.Reason)
	}
	if r.DifferenceBasisPoints != 50 {
		t.Errorf("want 50 bp, got %d", r.DifferenceBasisPoints)
	}
	if r.DifferenceAmount.Amount() != 5_000 {
		t.Errorf("want DifferenceAmount 5000, got %d", r.DifferenceAmount.Amount())
	}
	if r.DifferenceAmount.Currency() != "USD" {
		t.Errorf("want DifferenceAmount currency USD, got %s", r.DifferenceAmount.Currency())
	}
	if r.Version != 1 {
		t.Errorf("want version 1, got %d", r.Version)
	}
}

// 2 — diffBP == toleranceBP → MATCHED (inclusive boundary, BR-002)
func TestReconcile_MatchedAtBoundary(t *testing.T) {
	// 100 USD estimated, 101 USD actual → 100 bp diff, tolerance 100 bp
	erps := []ERPShipment{erpUSD("TRK-2", 1_000_000)}
	carriers := []CarrierInvoice{carrierUSD("TRK-2", 1_010_000)}

	results := Reconcile(erps, carriers, 100)

	if results[0].Status != StatusMatched {
		t.Errorf("want MATCHED at boundary, got %s", results[0].Status)
	}
	if results[0].DifferenceBasisPoints != 100 {
		t.Errorf("want 100 bp, got %d", results[0].DifferenceBasisPoints)
	}
}

// 3 — diff > tolerance → DISCREPANCY (BR-003)
func TestReconcile_Discrepancy(t *testing.T) {
	// 100 USD estimated, 102 USD actual → 200 bp diff, tolerance 100 bp
	erps := []ERPShipment{erpUSD("TRK-3", 1_000_000)}
	carriers := []CarrierInvoice{carrierUSD("TRK-3", 1_020_000)}

	results := Reconcile(erps, carriers, 100)

	r := results[0]
	if r.Status != StatusDiscrepancy {
		t.Errorf("want DISCREPANCY, got %s", r.Status)
	}
	if r.Reason != ReasonAboveTolerance {
		t.Errorf("want %s, got %s", ReasonAboveTolerance, r.Reason)
	}
	if r.DifferenceBasisPoints != 200 {
		t.Errorf("want 200 bp, got %d", r.DifferenceBasisPoints)
	}
}

// 4 — no carrier → UNRECONCILED_ERP (BR-004)
func TestReconcile_UnreconciledERP(t *testing.T) {
	erps := []ERPShipment{erpUSD("TRK-4", 1_000_000)}

	results := Reconcile(erps, nil, 100)

	r := results[0]
	if r.Status != StatusUnreconciledERP {
		t.Errorf("want UNRECONCILED_ERP, got %s", r.Status)
	}
	if r.Reason != ReasonMissingCarrierInvoice {
		t.Errorf("want %s, got %s", ReasonMissingCarrierInvoice, r.Reason)
	}
	if r.EstimatedFreight.Amount() != 1_000_000 {
		t.Errorf("estimated must carry original value")
	}
	if r.DifferenceBasisPoints != 0 {
		t.Errorf("diff bp must be 0 for orphan, got %d", r.DifferenceBasisPoints)
	}
}

// 5 — no ERP → UNRECONCILED_CARRIER (BR-005)
func TestReconcile_UnreconciledCarrier(t *testing.T) {
	carriers := []CarrierInvoice{carrierUSD("TRK-5", 1_000_000)}

	results := Reconcile(nil, carriers, 100)

	r := results[0]
	if r.Status != StatusUnreconciledCarrier {
		t.Errorf("want UNRECONCILED_CARRIER, got %s", r.Status)
	}
	if r.Reason != ReasonMissingERPShipment {
		t.Errorf("want %s, got %s", ReasonMissingERPShipment, r.Reason)
	}
	if r.ActualFreight.Amount() != 1_000_000 {
		t.Errorf("actual must carry original value")
	}
}

// 6 — ActualFreight non-USD → DISCREPANCY + currency_not_normalized
func TestReconcile_CurrencyNotNormalized(t *testing.T) {
	erps := []ERPShipment{erpUSD("TRK-6", 1_000_000)}
	carriers := []CarrierInvoice{carrierFX("TRK-6", 5_000_000, "EUR")}

	results := Reconcile(erps, carriers, 100)

	r := results[0]
	if r.Status != StatusDiscrepancy {
		t.Errorf("want DISCREPANCY, got %s", r.Status)
	}
	if r.Reason != ReasonCurrencyNotNormalized {
		t.Errorf("want %s, got %s", ReasonCurrencyNotNormalized, r.Reason)
	}
	if r.DifferenceBasisPoints != 0 {
		t.Errorf("diff bp must be 0 when non-USD, got %d", r.DifferenceBasisPoints)
	}
	if r.ActualFreight.Currency() != "EUR" {
		t.Errorf("ActualFreight currency must be preserved")
	}
}

// 7 — zero estimated, non-zero actual → MaxInt64 → DISCREPANCY
func TestReconcile_ZeroEstimated_NonZeroActual(t *testing.T) {
	erps := []ERPShipment{erpUSD("TRK-7", 0)}
	carriers := []CarrierInvoice{carrierUSD("TRK-7", 1_000_000)}

	results := Reconcile(erps, carriers, 100)

	r := results[0]
	if r.Status != StatusDiscrepancy {
		t.Errorf("want DISCREPANCY for zero estimated, got %s", r.Status)
	}
	if r.DifferenceBasisPoints != math.MaxInt64 {
		t.Errorf("want MaxInt64 bp, got %d", r.DifferenceBasisPoints)
	}
}

// 8 — both zero → 0 bp → MATCHED
func TestReconcile_BothZero(t *testing.T) {
	erps := []ERPShipment{erpUSD("TRK-8", 0)}
	carriers := []CarrierInvoice{carrierUSD("TRK-8", 0)}

	results := Reconcile(erps, carriers, 100)

	if results[0].Status != StatusMatched {
		t.Errorf("want MATCHED for 0/0, got %s", results[0].Status)
	}
	if results[0].DifferenceBasisPoints != 0 {
		t.Errorf("want 0 bp, got %d", results[0].DifferenceBasisPoints)
	}
}

// 10 — actual < estimated → abs() fires → correct bp → DISCREPANCY
func TestReconcile_Discrepancy_ActualBelowEstimated(t *testing.T) {
	// estimated 102 USD (1_020_000), actual 100 USD (1_000_000)
	// diff = 20_000; diffBP = 20_000 * 10_000 / 1_020_000 = 196 bp (integer division)
	erps := []ERPShipment{erpUSD("TRK-10", 1_020_000)}
	carriers := []CarrierInvoice{carrierUSD("TRK-10", 1_000_000)}

	results := Reconcile(erps, carriers, 100)

	r := results[0]
	if r.Status != StatusDiscrepancy {
		t.Errorf("want DISCREPANCY, got %s", r.Status)
	}
	if r.Reason != ReasonAboveTolerance {
		t.Errorf("want %s, got %s", ReasonAboveTolerance, r.Reason)
	}
	if r.DifferenceBasisPoints != 196 {
		t.Errorf("want 196 bp, got %d", r.DifferenceBasisPoints)
	}
	if r.DifferenceAmount.Amount() != 20_000 {
		t.Errorf("want DifferenceAmount 20000, got %d", r.DifferenceAmount.Amount())
	}
}

// 9 — mixed batch: MATCHED + DISCREPANCY + UNRECONCILED_ERP + UNRECONCILED_CARRIER
func TestReconcile_MixedBatch(t *testing.T) {
	erps := []ERPShipment{
		erpUSD("MATCH-1", 1_000_000),
		erpUSD("DISC-1", 1_000_000),
		erpUSD("ERP-ONLY", 1_000_000),
	}
	carriers := []CarrierInvoice{
		carrierUSD("MATCH-1", 1_005_000),    // 50 bp → MATCHED
		carrierUSD("DISC-1", 1_020_000),     // 200 bp → DISCREPANCY
		carrierUSD("CARRIER-ONLY", 500_000), // no ERP → UNRECONCILED_CARRIER
	}

	results := Reconcile(erps, carriers, 100)

	if len(results) != 4 {
		t.Fatalf("want 4 results, got %d", len(results))
	}

	check := func(tracking string, wantStatus ReconciliationStatus) {
		r, ok := findResult(results, tracking)
		if !ok {
			t.Errorf("no result for %s", tracking)
			return
		}
		if r.Status != wantStatus {
			t.Errorf("%s: want %s, got %s", tracking, wantStatus, r.Status)
		}
	}

	check("MATCH-1", StatusMatched)
	check("DISC-1", StatusDiscrepancy)
	check("ERP-ONLY", StatusUnreconciledERP)
	check("CARRIER-ONLY", StatusUnreconciledCarrier)
}

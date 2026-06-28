package domain

import "testing"

func TestReprocess_ChangedInvoice(t *testing.T) {
	erp := ERPShipment{
		TrackingNumber:      "TRK-001",
		EstimatedFreightUSD: MustMoney(100_0000, "USD"),
	}
	origResult := ReconciliationResult{
		TrackingNumber:   "TRK-001",
		Status:           StatusMatched,
		Reason:           ReasonWithinTolerance,
		EstimatedFreight: MustMoney(100_0000, "USD"),
		ActualFreight:    MustMoney(100_0000, "USD"),
		Version:          1,
	}
	origEntries, _ := GenerateLedgerEntries(origResult)

	newInvoice := CarrierInvoice{
		TrackingNumber:        "TRK-001",
		ActualFreightCurrency: MustMoney(150_0000, "USD"),
	}

	out, err := Reprocess(ReprocessingInput{
		Original:    origResult,
		OrigEntries: origEntries,
		NewInvoice:  newInvoice,
		ERP:         erp,
		ToleranceBP: 100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Changed {
		t.Error("expected Changed=true")
	}
	if out.NewResult.Version != origResult.Version+1 {
		t.Errorf("Version: want %d, got %d", origResult.Version+1, out.NewResult.Version)
	}
	if len(out.Reversals) == 0 {
		t.Error("expected non-empty Reversals")
	}
	if len(out.NewEntries) == 0 {
		t.Error("expected non-empty NewEntries")
	}
}

func TestReprocess_IdempotentSameInvoice(t *testing.T) {
	erp := ERPShipment{
		TrackingNumber:      "TRK-002",
		EstimatedFreightUSD: MustMoney(100_0000, "USD"),
	}
	origResult := ReconciliationResult{
		TrackingNumber:   "TRK-002",
		Status:           StatusMatched,
		Reason:           ReasonWithinTolerance,
		EstimatedFreight: MustMoney(100_0000, "USD"),
		ActualFreight:    MustMoney(100_0000, "USD"),
		Version:          1,
	}
	origEntries, _ := GenerateLedgerEntries(origResult)

	sameInvoice := CarrierInvoice{
		TrackingNumber:        "TRK-002",
		ActualFreightCurrency: MustMoney(100_0000, "USD"),
	}

	out, err := Reprocess(ReprocessingInput{
		Original:    origResult,
		OrigEntries: origEntries,
		NewInvoice:  sameInvoice,
		ERP:         erp,
		ToleranceBP: 100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Changed {
		t.Error("expected Changed=false")
	}
	if len(out.Reversals) != 0 {
		t.Errorf("expected empty Reversals, got %d", len(out.Reversals))
	}
	if len(out.NewEntries) != 0 {
		t.Errorf("expected empty NewEntries, got %d", len(out.NewEntries))
	}
	if out.NewResult.Version != origResult.Version {
		t.Errorf("Version: want %d, got %d", origResult.Version, out.NewResult.Version)
	}
}

func TestReprocess_ReversalExact(t *testing.T) {
	erp := ERPShipment{
		TrackingNumber:      "TRK-003",
		EstimatedFreightUSD: MustMoney(200_0000, "USD"),
	}
	origResult := ReconciliationResult{
		TrackingNumber:   "TRK-003",
		Status:           StatusMatched,
		Reason:           ReasonWithinTolerance,
		EstimatedFreight: MustMoney(200_0000, "USD"),
		ActualFreight:    MustMoney(200_0000, "USD"),
		Version:          2,
	}
	origEntries, _ := GenerateLedgerEntries(origResult)

	newInvoice := CarrierInvoice{
		TrackingNumber:        "TRK-003",
		ActualFreightCurrency: MustMoney(210_0000, "USD"),
	}

	out, err := Reprocess(ReprocessingInput{
		Original:    origResult,
		OrigEntries: origEntries,
		NewInvoice:  newInvoice,
		ERP:         erp,
		ToleranceBP: 100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Reversals) != len(origEntries) {
		t.Fatalf("Reversals len: want %d, got %d", len(origEntries), len(out.Reversals))
	}
	for i, rev := range out.Reversals {
		orig := origEntries[i]
		if rev.AccountDebit != orig.AccountCredit {
			t.Errorf("[%d] reversal debit %q != orig credit %q", i, rev.AccountDebit, orig.AccountCredit)
		}
		if rev.AccountCredit != orig.AccountDebit {
			t.Errorf("[%d] reversal credit %q != orig debit %q", i, rev.AccountCredit, orig.AccountDebit)
		}
		if rev.Amount.Amount() != orig.Amount.Amount() || rev.Amount.Currency() != orig.Amount.Currency() {
			t.Errorf("[%d] amount mismatch: want %v %v, got %v %v",
				i, orig.Amount.Amount(), orig.Amount.Currency(),
				rev.Amount.Amount(), rev.Amount.Currency())
		}
	}
}

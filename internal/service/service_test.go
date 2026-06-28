package service_test

import (
	"testing"

	"orcaarch/internal/domain"
	"orcaarch/internal/service"
)

func TestProcess_AllStatuses(t *testing.T) {
	res, err := service.Process(100)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	statuses := make(map[domain.ReconciliationStatus]bool)
	for _, r := range res.Reconciliations {
		statuses[r.Status] = true
	}
	for _, want := range []domain.ReconciliationStatus{
		domain.StatusMatched,
		domain.StatusDiscrepancy,
		domain.StatusUnreconciledERP,
		domain.StatusUnreconciledCarrier,
	} {
		if !statuses[want] {
			t.Errorf("status %q missing from reconciliations", want)
		}
	}
}

func TestProcess_Reprocessing(t *testing.T) {
	res, err := service.Process(100)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if res.Reprocessed == nil {
		t.Fatal("expected non-nil Reprocessed")
	}
	if !res.Reprocessed.Changed {
		t.Error("Reprocessed.Changed must be true")
	}
	if res.Reprocessed.NewResult.Version <= 1 {
		t.Errorf("expected version > 1, got %d", res.Reprocessed.NewResult.Version)
	}
	if len(res.Reprocessed.Reversals) == 0 {
		t.Error("expected non-empty Reversals")
	}
}

func TestProcess_DoubleEntry(t *testing.T) {
	res, err := service.Process(100)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	// Each LedgerEntry debits and credits the same Amount; sum(debits)==sum(credits) per TN.
	type sums struct{ debit, credit int64 }
	byTN := make(map[string]*sums)
	for _, e := range res.LedgerEntries {
		if byTN[e.TrackingNumber] == nil {
			byTN[e.TrackingNumber] = &sums{}
		}
		byTN[e.TrackingNumber].debit += e.Amount.Amount()
		byTN[e.TrackingNumber].credit += e.Amount.Amount()
	}
	for tn, s := range byTN {
		if s.debit != s.credit {
			t.Errorf("TN %s: debit %d != credit %d", tn, s.debit, s.credit)
		}
	}
}

func TestProcess_ZeroN(t *testing.T) {
	res, err := service.Process(0)
	if err != nil {
		t.Fatalf("Process(0): %v", err)
	}
	if len(res.Reconciliations) != 0 {
		t.Errorf("expected 0 reconciliations for n=0, got %d", len(res.Reconciliations))
	}
	if res.Reprocessed != nil {
		t.Error("expected nil Reprocessed for n=0")
	}
}

func TestProcess_LargeN(t *testing.T) {
	_, err := service.Process(10_000)
	if err != nil {
		t.Fatalf("Process(10000): %v", err)
	}
}

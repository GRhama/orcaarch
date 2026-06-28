package service

import (
	"testing"

	"orcaarch/internal/domain"
)

func BenchmarkProcess50k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Process(50_000); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcess100k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Process(100_000); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcess250k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Process(250_000); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSmokeProcess100k(t *testing.T) {
	result, err := Process(100_000)
	if err != nil {
		t.Fatal(err)
	}

	const n = 100_000
	if got := len(result.Reconciliations); got != n {
		t.Errorf("Reconciliations: want %d, got %d", n, got)
	}
	if got := len(result.Quarantine); got != n/4 {
		t.Errorf("Quarantine: want %d, got %d", n/4, got)
	}
	if got := len(result.LedgerEntries); got != n+2 {
		t.Errorf("LedgerEntries: want %d, got %d", n+2, got)
	}

	statuses := make(map[domain.ReconciliationStatus]bool)
	for _, r := range result.Reconciliations {
		statuses[r.Status] = true
	}
	for _, s := range []domain.ReconciliationStatus{
		domain.StatusMatched,
		domain.StatusDiscrepancy,
		domain.StatusUnreconciledERP,
		domain.StatusUnreconciledCarrier,
	} {
		if !statuses[s] {
			t.Errorf("missing reconciliation status: %v", s)
		}
	}
}

package mock

import (
	"testing"

	"orcaarch/internal/domain"
)

func TestGenerate_AllFourStates(t *testing.T) {
	erps, carriers := Generate(1000)
	results := domain.Reconcile(erps, carriers, 500) // 500 bp = 5% tolerance

	counts := map[domain.ReconciliationStatus]int{}
	for _, r := range results {
		counts[r.Status]++
	}

	want := []domain.ReconciliationStatus{
		domain.StatusMatched,
		domain.StatusDiscrepancy,
		domain.StatusUnreconciledERP,
		domain.StatusUnreconciledCarrier,
	}
	for _, s := range want {
		if counts[s] == 0 {
			t.Errorf("status %v not present in results", s)
		}
	}
}

func TestGenerate_NoPanic(t *testing.T) {
	for _, n := range []int{0, 1, 10000} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Generate(%d) panicked: %v", n, r)
				}
			}()
			Generate(n)
		}()
	}
}

func TestGenerate_UniqueTrackingNumbers(t *testing.T) {
	erps, carriers := Generate(1000)

	seen := map[string]bool{}
	for _, e := range erps {
		if seen[e.TrackingNumber] {
			t.Errorf("duplicate tracking number in ERPs: %s", e.TrackingNumber)
		}
		seen[e.TrackingNumber] = true
	}

	seen = map[string]bool{}
	for _, c := range carriers {
		if seen[c.TrackingNumber] {
			t.Errorf("duplicate tracking number in Carriers: %s", c.TrackingNumber)
		}
		seen[c.TrackingNumber] = true
	}
}

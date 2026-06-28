package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"orcaarch/internal/api"
	"orcaarch/internal/domain"
	"orcaarch/internal/report"
	"orcaarch/internal/service"
)

func makeResult() service.ProcessingResult {
	return service.ProcessingResult{
		Reconciliations: []domain.ReconciliationResult{
			{
				TrackingNumber:        "TN-001",
				Status:                domain.StatusMatched,
				EstimatedFreight:      domain.MustMoney(100000, "USD"),
				ActualFreight:         domain.MustMoney(100000, "USD"),
				DifferenceAmount:      domain.MustMoney(0, "USD"),
				DifferenceBasisPoints: 0,
				Version:               1,
			},
		},
		LedgerEntries: []domain.LedgerEntry{
			{
				TrackingNumber: "TN-001",
				AccountDebit:   "Estoque de Mercadorias",
				AccountCredit:  "Fretes a Pagar",
				Amount:         domain.MustMoney(100000, "USD"),
			},
		},
		Quarantine: []domain.Quarantine{
			{
				Invoice: domain.CarrierInvoice{
					InvoiceID:             "INV-004",
					TrackingNumber:        "TN-004",
					ActualFreightCurrency: domain.MustMoney(75000, "USD"),
					CustomsDutiesLocal:    domain.MustMoney(0, "USD"),
					InsuranceCostUSD:      domain.MustMoney(0, "USD"),
					AncillaryFeesUSD:      domain.MustMoney(0, "USD"),
				},
				Reason: domain.ReasonMissingERPShipment,
			},
		},
		LandedCosts: map[string]domain.Money{},
		WACBySKU:    map[string]domain.Money{},
	}
}

func TestInventoryEndpoint(t *testing.T) {
	srv := api.NewServer(makeResult())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/inventory", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json, got %q", ct)
	}
	var rows []report.ReconciliationRow
	if err := json.NewDecoder(rec.Body).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].TrackingNumber != "TN-001" {
		t.Fatalf("unexpected inventory rows: %+v", rows)
	}
}

func TestLedgerEndpoint(t *testing.T) {
	srv := api.NewServer(makeResult())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/ledger", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var rows []report.LedgerRow
	if err := json.NewDecoder(rec.Body).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Amount != 100000 {
		t.Fatalf("unexpected ledger rows: %+v", rows)
	}
}

func TestRiskEndpoint(t *testing.T) {
	srv := api.NewServer(makeResult())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/risk", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var rows []report.RiskRow
	if err := json.NewDecoder(rec.Body).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].TrackingNumber != "TN-004" || rows[0].Amount != 75000 {
		t.Fatalf("unexpected risk rows: %+v", rows)
	}
}

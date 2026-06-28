package export_test

import (
	"bytes"
	"encoding/csv"
	"testing"

	"orcaarch/internal/export"
	"orcaarch/internal/report"
)

func TestWriteReconciliationCSV(t *testing.T) {
	rows := []report.ReconciliationRow{
		{TrackingNumber: "TN-001", Status: "MATCHED", EstimatedAmount: 100000, EstimatedCurrency: "USD", ActualAmount: 100000, ActualCurrency: "USD", DifferenceAmount: 0, DifferenceCurrency: "USD", DifferenceBP: 0, Version: 1},
		{TrackingNumber: "TN-002", Status: "DISCREPANCY", EstimatedAmount: 90000, EstimatedCurrency: "USD", ActualAmount: 108000, ActualCurrency: "USD", DifferenceAmount: 18000, DifferenceCurrency: "USD", DifferenceBP: 2000, Version: 1},
	}
	var buf bytes.Buffer
	if err := export.WriteReconciliationCSV(&buf, rows); err != nil {
		t.Fatal(err)
	}
	records, err := csv.NewReader(&buf).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("want 3 records (header+2), got %d", len(records))
	}
	if records[0][0] != "tracking_number" {
		t.Fatalf("want header[0]=tracking_number, got %q", records[0][0])
	}
	if records[1][0] != "TN-001" || records[2][0] != "TN-002" {
		t.Fatalf("unexpected data rows: %v %v", records[1], records[2])
	}
}

func TestWriteLedgerCSV(t *testing.T) {
	rows := []report.LedgerRow{
		{TrackingNumber: "TN-001", AccountDebit: "Estoque de Mercadorias", AccountCredit: "Fretes a Pagar", Amount: 100000, Currency: "USD"},
		{TrackingNumber: "TN-002", AccountDebit: "Despesas Logísticas Suspensas", AccountCredit: "Contas a Pagar", Amount: 85000, Currency: "USD"},
	}
	var buf bytes.Buffer
	if err := export.WriteLedgerCSV(&buf, rows); err != nil {
		t.Fatal(err)
	}
	records, err := csv.NewReader(&buf).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("want 3 records (header+2), got %d", len(records))
	}
	if records[0][0] != "tracking_number" {
		t.Fatalf("want header[0]=tracking_number, got %q", records[0][0])
	}
	if records[1][3] != "100000" {
		t.Fatalf("want amount=100000, got %q", records[1][3])
	}
}

func TestWriteRiskCSV(t *testing.T) {
	rows := []report.RiskRow{
		{TrackingNumber: "TN-004", Reason: "missing_erp_shipment", Amount: 75000, Currency: "USD"},
	}
	var buf bytes.Buffer
	if err := export.WriteRiskCSV(&buf, rows); err != nil {
		t.Fatal(err)
	}
	records, err := csv.NewReader(&buf).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("want 2 records (header+1), got %d", len(records))
	}
	if records[1][0] != "TN-004" || records[1][2] != "75000" {
		t.Fatalf("unexpected risk row: %v", records[1])
	}
}

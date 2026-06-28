package domain

import (
	"errors"
	"testing"
)

func makeResult(status ReconciliationStatus, estimated, actual Money) ReconciliationResult {
	return ReconciliationResult{
		TrackingNumber:   "TRK-001",
		Status:           status,
		EstimatedFreight: estimated,
		ActualFreight:    actual,
		Version:          1,
	}
}

func TestGenerateLedgerEntries_Matched(t *testing.T) {
	actual := MustMoney(50000, "USD")
	entries, err := GenerateLedgerEntries(makeResult(StatusMatched, MustMoney(51000, "USD"), actual))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.AccountDebit != AccountEstoqueMercadorias {
		t.Errorf("debit: want %q, got %q", AccountEstoqueMercadorias, e.AccountDebit)
	}
	if e.AccountCredit != AccountFretesAPagar {
		t.Errorf("credit: want %q, got %q", AccountFretesAPagar, e.AccountCredit)
	}
	if e.Amount != actual {
		t.Errorf("amount: want %v, got %v", actual, e.Amount)
	}
}

func TestGenerateLedgerEntries_Discrepancy(t *testing.T) {
	actual := MustMoney(80000, "USD")
	entries, err := GenerateLedgerEntries(makeResult(StatusDiscrepancy, MustMoney(50000, "USD"), actual))
	if err != nil {
		t.Fatal(err)
	}
	e := entries[0]
	if e.AccountDebit != AccountEstoqueMercadorias || e.AccountCredit != AccountFretesAPagar {
		t.Errorf("unexpected accounts: debit=%q credit=%q", e.AccountDebit, e.AccountCredit)
	}
	if e.Amount != actual {
		t.Errorf("amount: want %v, got %v", actual, e.Amount)
	}
}

func TestGenerateLedgerEntries_UnreconciledERP(t *testing.T) {
	estimated := MustMoney(30000, "USD")
	entries, err := GenerateLedgerEntries(makeResult(StatusUnreconciledERP, estimated, Money{}))
	if err != nil {
		t.Fatal(err)
	}
	e := entries[0]
	if e.AccountDebit != AccountEstoqueMercadorias {
		t.Errorf("debit: want %q, got %q", AccountEstoqueMercadorias, e.AccountDebit)
	}
	if e.AccountCredit != AccountProvisoesFretePagar {
		t.Errorf("credit: want %q, got %q", AccountProvisoesFretePagar, e.AccountCredit)
	}
	if e.Amount != estimated {
		t.Errorf("amount: want %v, got %v", estimated, e.Amount)
	}
}

func TestGenerateLedgerEntries_UnreconciledCarrier(t *testing.T) {
	actual := MustMoney(45000, "USD")
	entries, err := GenerateLedgerEntries(makeResult(StatusUnreconciledCarrier, Money{}, actual))
	if err != nil {
		t.Fatal(err)
	}
	e := entries[0]
	if e.AccountDebit != AccountDespesasLogisticasSuspensas {
		t.Errorf("debit: want %q, got %q", AccountDespesasLogisticasSuspensas, e.AccountDebit)
	}
	if e.AccountCredit != AccountContasAPagar {
		t.Errorf("credit: want %q, got %q", AccountContasAPagar, e.AccountCredit)
	}
	if e.Amount != actual {
		t.Errorf("amount: want %v, got %v", actual, e.Amount)
	}
}

func TestGenerateLedgerEntries_UnknownStatus(t *testing.T) {
	_, err := GenerateLedgerEntries(ReconciliationResult{Status: "INVALID"})
	if !errors.Is(err, ErrUnknownReconciliationStatus) {
		t.Errorf("want ErrUnknownReconciliationStatus, got %v", err)
	}
}

func TestGenerateLedgerEntries_Discrepancy_NonUSD(t *testing.T) {
	actual := MustMoney(80000, "BRL")
	r := ReconciliationResult{
		TrackingNumber:   "TRK-002",
		Status:           StatusDiscrepancy,
		EstimatedFreight: MustMoney(50000, "USD"),
		ActualFreight:    actual,
		Reason:           ReasonCurrencyNotNormalized,
		Version:          1,
	}
	entries, err := GenerateLedgerEntries(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.AccountDebit != AccountEstoqueMercadorias {
		t.Errorf("debit: want %q, got %q", AccountEstoqueMercadorias, e.AccountDebit)
	}
	if e.AccountCredit != AccountFretesAPagar {
		t.Errorf("credit: want %q, got %q", AccountFretesAPagar, e.AccountCredit)
	}
	if e.Amount != actual {
		t.Errorf("amount: want %v, got %v", actual, e.Amount)
	}
	if e.Amount.Currency() != "BRL" {
		t.Errorf("currency preserved: want BRL, got %s", e.Amount.Currency())
	}
}

func TestGenerateReversal(t *testing.T) {
	amount := MustMoney(20000, "USD")
	original := []LedgerEntry{{
		TrackingNumber: "TRK-001",
		AccountDebit:   AccountEstoqueMercadorias,
		AccountCredit:  AccountFretesAPagar,
		Amount:         amount,
	}}
	reversed := GenerateReversal(original)
	if len(reversed) != 1 {
		t.Fatalf("want 1 reversed entry, got %d", len(reversed))
	}
	r := reversed[0]
	if r.AccountDebit != AccountFretesAPagar {
		t.Errorf("reversal debit: want %q, got %q", AccountFretesAPagar, r.AccountDebit)
	}
	if r.AccountCredit != AccountEstoqueMercadorias {
		t.Errorf("reversal credit: want %q, got %q", AccountEstoqueMercadorias, r.AccountCredit)
	}
	if r.Amount != amount {
		t.Errorf("reversal amount changed: want %v, got %v", amount, r.Amount)
	}
}

func TestGenerateReversal_Empty(t *testing.T) {
	reversed := GenerateReversal([]LedgerEntry{})
	if len(reversed) != 0 {
		t.Errorf("want empty slice, got %d entries", len(reversed))
	}
}

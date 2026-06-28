package domain

import "errors"

const (
	AccountEstoqueMercadorias          = "Estoque de Mercadorias"
	AccountFretesAPagar                = "Fretes a Pagar"
	AccountProvisoesFretePagar         = "Provisões de Frete a Pagar"
	AccountDespesasLogisticasSuspensas = "Despesas Logísticas Suspensas"
	AccountContasAPagar                = "Contas a Pagar"
)

var ErrUnknownReconciliationStatus = errors.New("unknown reconciliation status")

type LedgerEntry struct {
	TrackingNumber string
	AccountDebit   string
	AccountCredit  string
	Amount         Money
}

// GenerateLedgerEntries produces one double-entry record per ReconciliationResult.
// Debit == Credit by construction (single Amount, two accounts).
func GenerateLedgerEntries(r ReconciliationResult) ([]LedgerEntry, error) {
	var debit, credit string
	var amount Money

	switch r.Status {
	case StatusMatched, StatusDiscrepancy:
		// ponytail: DISCREPANCY may carry non-USD ActualFreight when Reconcile sets ReasonCurrencyNotNormalized; no FX conversion here, currency passes through as-is.
		debit, credit = AccountEstoqueMercadorias, AccountFretesAPagar
		amount = r.ActualFreight
	case StatusUnreconciledERP:
		debit, credit = AccountEstoqueMercadorias, AccountProvisoesFretePagar
		amount = r.EstimatedFreight
	case StatusUnreconciledCarrier:
		debit, credit = AccountDespesasLogisticasSuspensas, AccountContasAPagar
		amount = r.ActualFreight
	default:
		return nil, ErrUnknownReconciliationStatus
	}

	return []LedgerEntry{{
		TrackingNumber: r.TrackingNumber,
		AccountDebit:   debit,
		AccountCredit:  credit,
		Amount:         amount,
	}}, nil
}

// GenerateReversal inverts each entry (swap debit↔credit), preserving amount.
func GenerateReversal(entries []LedgerEntry) []LedgerEntry {
	reversed := make([]LedgerEntry, len(entries))
	for i, e := range entries {
		reversed[i] = LedgerEntry{
			TrackingNumber: e.TrackingNumber,
			AccountDebit:   e.AccountCredit,
			AccountCredit:  e.AccountDebit,
			Amount:         e.Amount,
		}
	}
	return reversed
}

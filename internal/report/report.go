package report

import "orcaarch/internal/domain"

type ReconciliationRow struct {
	TrackingNumber     string `json:"tracking_number"`
	Status             string `json:"status"`
	EstimatedAmount    int64  `json:"estimated_amount"`
	EstimatedCurrency  string `json:"estimated_currency"`
	ActualAmount       int64  `json:"actual_amount"`
	ActualCurrency     string `json:"actual_currency"`
	DifferenceAmount   int64  `json:"difference_amount"`
	DifferenceCurrency string `json:"difference_currency"`
	DifferenceBP       int64  `json:"difference_bp"`
	Version            int    `json:"version"`
}

type LedgerRow struct {
	TrackingNumber string `json:"tracking_number"`
	AccountDebit   string `json:"account_debit"`
	AccountCredit  string `json:"account_credit"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
}

type RiskRow struct {
	TrackingNumber string `json:"tracking_number"`
	Reason         string `json:"reason"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
}

func ToReconciliationRows(rs []domain.ReconciliationResult) []ReconciliationRow {
	rows := make([]ReconciliationRow, len(rs))
	for i, r := range rs {
		rows[i] = ReconciliationRow{
			TrackingNumber:     r.TrackingNumber,
			Status:             string(r.Status),
			EstimatedAmount:    r.EstimatedFreight.Amount(),
			EstimatedCurrency:  r.EstimatedFreight.Currency(),
			ActualAmount:       r.ActualFreight.Amount(),
			ActualCurrency:     r.ActualFreight.Currency(),
			DifferenceAmount:   r.DifferenceAmount.Amount(),
			DifferenceCurrency: r.DifferenceAmount.Currency(),
			DifferenceBP:       r.DifferenceBasisPoints,
			Version:            r.Version,
		}
	}
	return rows
}

func ToLedgerRows(entries []domain.LedgerEntry) []LedgerRow {
	rows := make([]LedgerRow, len(entries))
	for i, e := range entries {
		rows[i] = LedgerRow{
			TrackingNumber: e.TrackingNumber,
			AccountDebit:   e.AccountDebit,
			AccountCredit:  e.AccountCredit,
			Amount:         e.Amount.Amount(),
			Currency:       e.Amount.Currency(),
		}
	}
	return rows
}

func ToRiskRows(qs []domain.Quarantine) []RiskRow {
	rows := make([]RiskRow, len(qs))
	for i, q := range qs {
		rows[i] = RiskRow{
			TrackingNumber: q.Invoice.TrackingNumber,
			Reason:         q.Reason,
			Amount:         q.Invoice.ActualFreightCurrency.Amount(),
			Currency:       q.Invoice.ActualFreightCurrency.Currency(),
		}
	}
	return rows
}

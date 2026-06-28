package export

import (
	"encoding/csv"
	"io"
	"strconv"

	"orcaarch/internal/report"
)

func WriteReconciliationCSV(w io.Writer, rows []report.ReconciliationRow) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"tracking_number", "status",
		"estimated_amount", "estimated_currency",
		"actual_amount", "actual_currency",
		"difference_amount", "difference_currency",
		"difference_bp", "version",
	}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.Write([]string{
			r.TrackingNumber, r.Status,
			i64(r.EstimatedAmount), r.EstimatedCurrency,
			i64(r.ActualAmount), r.ActualCurrency,
			i64(r.DifferenceAmount), r.DifferenceCurrency,
			i64(r.DifferenceBP), strconv.Itoa(r.Version),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteLedgerCSV(w io.Writer, rows []report.LedgerRow) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"tracking_number", "account_debit", "account_credit", "amount", "currency",
	}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.Write([]string{
			r.TrackingNumber, r.AccountDebit, r.AccountCredit,
			i64(r.Amount), r.Currency,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteRiskCSV(w io.Writer, rows []report.RiskRow) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"tracking_number", "reason", "amount", "currency",
	}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.Write([]string{
			r.TrackingNumber, r.Reason, i64(r.Amount), r.Currency,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func i64(n int64) string { return strconv.FormatInt(n, 10) }

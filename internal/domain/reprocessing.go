package domain

// ReprocessingInput holds all state needed to retract and recompute a reconciliation.
type ReprocessingInput struct {
	Original    ReconciliationResult
	OrigEntries []LedgerEntry
	NewInvoice  CarrierInvoice
	ERP         ERPShipment
	ToleranceBP int64
}

// ReprocessingOutput carries reversals, updated result, and new ledger entries.
// ponytail: WAC recalc excluded — CalculateWAC requires stock state (OldStockMilliTons,
// OldWACPerMilliTon) absent from ReconciliationResult and CarrierInvoice; caller chains
// CalculateWAC after Reprocess using its own stock context.
type ReprocessingOutput struct {
	Changed    bool
	Reversals  []LedgerEntry
	NewResult  ReconciliationResult
	NewEntries []LedgerEntry
}

// Reprocess detects a carrier invoice change, generates exact reversals of OrigEntries,
// re-reconciles with the new invoice, increments Version, and produces new ledger entries.
// Idempotent: identical NewInvoice returns Changed=false with no new entries.
// ERP.TrackingNumber and NewInvoice.TrackingNumber must match; divergent tracking numbers
// produce undefined results because Reconcile may return multiple results in unspecified order.
func Reprocess(in ReprocessingInput) (ReprocessingOutput, error) {
	orig := in.Original.ActualFreight
	novo := in.NewInvoice.ActualFreightCurrency
	if orig.Amount() == novo.Amount() && orig.Currency() == novo.Currency() {
		return ReprocessingOutput{Changed: false, NewResult: in.Original}, nil
	}

	reversals := GenerateReversal(in.OrigEntries)

	results := Reconcile([]ERPShipment{in.ERP}, []CarrierInvoice{in.NewInvoice}, in.ToleranceBP)
	newResult := results[0]
	newResult.Version = in.Original.Version + 1

	newEntries, err := GenerateLedgerEntries(newResult)
	if err != nil {
		return ReprocessingOutput{}, err
	}

	return ReprocessingOutput{
		Changed:    true,
		Reversals:  reversals,
		NewResult:  newResult,
		NewEntries: newEntries,
	}, nil
}

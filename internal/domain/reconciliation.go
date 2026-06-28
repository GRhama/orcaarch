package domain

import "math"

type ReconciliationStatus string

const (
	StatusMatched             ReconciliationStatus = "MATCHED"
	StatusDiscrepancy         ReconciliationStatus = "DISCREPANCY"
	StatusUnreconciledERP     ReconciliationStatus = "UNRECONCILED_ERP"
	StatusUnreconciledCarrier ReconciliationStatus = "UNRECONCILED_CARRIER"
)

const (
	ReasonWithinTolerance       = "within_tolerance"
	ReasonAboveTolerance        = "above_tolerance"
	ReasonMissingCarrierInvoice = "missing_carrier_invoice"
	ReasonMissingERPShipment    = "missing_erp_shipment"
	ReasonCurrencyNotNormalized = "currency_not_normalized"
)

type ReconciliationResult struct {
	TrackingNumber        string
	Status                ReconciliationStatus
	Reason                string
	EstimatedFreight      Money
	ActualFreight         Money
	DifferenceAmount      Money // Money{} when not computable (orphan or non-USD actual)
	DifferenceBasisPoints int64 // 100 = 1%; 0 when not computable
	Version               int
}

// Reconcile matches ERPShipments to CarrierInvoices by TrackingNumber and
// classifies each pair. toleranceBasisPoints: 100 = 1%, 50 = 0.5%.
// The order of results is undefined; callers must not rely on slice ordering.
func Reconcile(erps []ERPShipment, carriers []CarrierInvoice, toleranceBasisPoints int64) []ReconciliationResult {
	erpMap := make(map[string]ERPShipment, len(erps))
	for _, e := range erps {
		erpMap[e.TrackingNumber] = e
	}

	carrierMap := make(map[string]CarrierInvoice, len(carriers))
	for _, c := range carriers {
		carrierMap[c.TrackingNumber] = c
	}

	results := make([]ReconciliationResult, 0, len(erps)+len(carriers))

	for _, erp := range erpMap {
		carrier, ok := carrierMap[erp.TrackingNumber]
		if !ok {
			results = append(results, ReconciliationResult{
				TrackingNumber:   erp.TrackingNumber,
				Status:           StatusUnreconciledERP,
				Reason:           ReasonMissingCarrierInvoice,
				EstimatedFreight: erp.EstimatedFreightUSD,
				Version:          1,
			})
			continue
		}

		if carrier.ActualFreightCurrency.Currency() != "USD" {
			// ponytail: no FX conversion — non-USD actual cannot be compared; caller must normalize first
			results = append(results, ReconciliationResult{
				TrackingNumber:   erp.TrackingNumber,
				Status:           StatusDiscrepancy,
				Reason:           ReasonCurrencyNotNormalized,
				EstimatedFreight: erp.EstimatedFreightUSD,
				ActualFreight:    carrier.ActualFreightCurrency,
				Version:          1,
			})
			continue
		}

		diffBP := absDiffBP(erp.EstimatedFreightUSD, carrier.ActualFreightCurrency)
		diffAmt := absDiffMoney(erp.EstimatedFreightUSD, carrier.ActualFreightCurrency)

		status := StatusDiscrepancy
		reason := ReasonAboveTolerance
		if diffBP <= toleranceBasisPoints {
			status = StatusMatched
			reason = ReasonWithinTolerance
		}

		results = append(results, ReconciliationResult{
			TrackingNumber:        erp.TrackingNumber,
			Status:                status,
			Reason:                reason,
			EstimatedFreight:      erp.EstimatedFreightUSD,
			ActualFreight:         carrier.ActualFreightCurrency,
			DifferenceAmount:      diffAmt,
			DifferenceBasisPoints: diffBP,
			Version:               1,
		})
	}

	for _, carrier := range carrierMap {
		if _, ok := erpMap[carrier.TrackingNumber]; !ok {
			results = append(results, ReconciliationResult{
				TrackingNumber: carrier.TrackingNumber,
				Status:         StatusUnreconciledCarrier,
				Reason:         ReasonMissingERPShipment,
				ActualFreight:  carrier.ActualFreightCurrency,
				Version:        1,
			})
		}
	}

	return results
}

// absDiffBP returns |estimated - actual| in basis points.
// Returns math.MaxInt64 when estimated is zero and actual is non-zero (infinite divergence).
func absDiffBP(estimated, actual Money) int64 {
	e := estimated.Amount()
	a := actual.Amount()
	if e == 0 && a == 0 {
		return 0
	}
	if e == 0 {
		return math.MaxInt64
	}
	diff := e - a
	if diff < 0 {
		diff = -diff
	}
	return diff * 10_000 / e
}

func absDiffMoney(estimated, actual Money) Money {
	diff := estimated.Amount() - actual.Amount()
	if diff < 0 {
		diff = -diff
	}
	return MustMoney(diff, estimated.Currency())
}

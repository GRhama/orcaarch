package domain

// Quarantine holds a CarrierInvoice with no matching ERP shipment.
// WAC must not be updated for quarantined entries; stock integrity depends on this invariant.
type Quarantine struct {
	Invoice CarrierInvoice
	Reason  string
}

// NewQuarantine creates a Quarantine for a carrier invoice that has no ERP counterpart.
func NewQuarantine(invoice CarrierInvoice) Quarantine {
	return Quarantine{Invoice: invoice, Reason: ReasonMissingERPShipment}
}

// ProvisionalLandedCostInput builds a LandedCostInput for UNRECONCILED_ERP shipments.
// Uses EstimatedFreightUSD as ActualFreightUSD; all other components zero in USD.
// Caller must ensure erp passes Validate(); non-USD EstimatedFreightUSD propagates to CalculateLandedCost as ErrLandedCostCurrencyNotUSD.
// ponytail: factory cost, duties, ancillary fees unknown without carrier invoice — caller enriches when available.
func ProvisionalLandedCostInput(erp ERPShipment) LandedCostInput {
	zero := MustMoney(0, "USD")
	return LandedCostInput{
		FactoryCostUSD:   zero,
		ActualFreightUSD: erp.EstimatedFreightUSD,
		CustomsDutiesUSD: zero,
		AncillaryFeesUSD: zero,
		InsuranceCostUSD: zero,
	}
}

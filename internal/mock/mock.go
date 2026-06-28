package mock

import (
	"fmt"
	"math/rand"
	"time"

	"orcaarch/internal/domain"
)

// Generate returns n ERPShipments and CarrierInvoices distributed across all 4 reconciliation states.
// Deterministic: seed fixed at 42.
func Generate(n int) ([]domain.ERPShipment, []domain.CarrierInvoice) {
	rng := rand.New(rand.NewSource(42))
	base := n / 4
	// buckets: [matched, discrepancy, erp_only, carrier_only]
	// remainder goes to matched
	sizes := [4]int{n - 3*base, base, base, base}

	var erps []domain.ERPShipment
	var carriers []domain.CarrierInvoice

	tn := 0 // global counter for unique tracking numbers

	// MATCHED: same TN both sides, carrier freight == ERP freight (0% diff)
	for i := 0; i < sizes[0]; i++ {
		freight := randFreight(rng)
		e, c := matchedPair(tn, freight, freight)
		erps = append(erps, e)
		carriers = append(carriers, c)
		tn++
	}

	// DISCREPANCY: same TN both sides, carrier freight = ERP * 1.20 (20% > any sane tolerance)
	for i := 0; i < sizes[1]; i++ {
		freight := randFreight(rng)
		higher := freight * 120 / 100
		e, c := matchedPair(tn, freight, higher)
		erps = append(erps, e)
		carriers = append(carriers, c)
		tn++
	}

	// UNRECONCILED_ERP: TN only in ERP
	for i := 0; i < sizes[2]; i++ {
		erps = append(erps, randERP(rng, tn))
		tn++
	}

	// UNRECONCILED_CARRIER: TN only in Carrier
	for i := 0; i < sizes[3]; i++ {
		carriers = append(carriers, randCarrier(rng, tn))
		tn++
	}

	return erps, carriers
}

func trackingNum(i int) string { return fmt.Sprintf("TN%08d", i) }

// randFreight returns a random freight amount in Money minor units (scale ×10000).
// Range: $1 000–$50 000 USD → 10_000_000–500_000_000.
func randFreight(rng *rand.Rand) int64 {
	return int64(rng.Intn(490_000_000)+10_000_000)
}

func randWeight(rng *rand.Rand) int64 { return int64(rng.Intn(99_000)+1_000) }

func baseDate() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }

func matchedPair(i int, erpFreight, carrierFreight int64) (domain.ERPShipment, domain.CarrierInvoice) {
	tn := trackingNum(i)
	e := domain.ERPShipment{
		BookingID:            fmt.Sprintf("BK%08d", i),
		TradeDate:            baseDate(),
		SKUID:                fmt.Sprintf("SKU%04d", i%100),
		TotalWeightMilliTons: 10_000,
		OriginCountry:        "CN",
		DestinationCountry:   "BR",
		EstimatedFreightUSD:  domain.MustMoney(erpFreight, "USD"),
		TrackingNumber:       tn,
	}
	c := domain.CarrierInvoice{
		InvoiceID:             fmt.Sprintf("INV%08d", i),
		InvoiceDate:           baseDate(),
		TrackingNumber:        tn,
		ActualFreightCurrency: domain.MustMoney(carrierFreight, "USD"),
		CustomsDutiesLocal:    domain.MustMoney(0, "USD"),
		InsuranceCostUSD:      domain.MustMoney(0, "USD"),
		AncillaryFeesUSD:      domain.MustMoney(0, "USD"),
	}
	return e, c
}

func randERP(rng *rand.Rand, i int) domain.ERPShipment {
	return domain.ERPShipment{
		BookingID:            fmt.Sprintf("BK%08d", i),
		TradeDate:            baseDate(),
		SKUID:                fmt.Sprintf("SKU%04d", i%100),
		TotalWeightMilliTons: randWeight(rng),
		OriginCountry:        "CN",
		DestinationCountry:   "BR",
		EstimatedFreightUSD:  domain.MustMoney(randFreight(rng), "USD"),
		TrackingNumber:       trackingNum(i),
	}
}

func randCarrier(rng *rand.Rand, i int) domain.CarrierInvoice {
	return domain.CarrierInvoice{
		InvoiceID:             fmt.Sprintf("INV%08d", i),
		InvoiceDate:           baseDate(),
		TrackingNumber:        trackingNum(i),
		ActualFreightCurrency: domain.MustMoney(randFreight(rng), "USD"),
		CustomsDutiesLocal:    domain.MustMoney(0, "USD"),
		InsuranceCostUSD:      domain.MustMoney(0, "USD"),
		AncillaryFeesUSD:      domain.MustMoney(0, "USD"),
	}
}

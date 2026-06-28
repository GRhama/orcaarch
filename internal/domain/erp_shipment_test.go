package domain

import (
	"testing"
	"time"
)

func validERP() ERPShipment {
	return ERPShipment{
		BookingID:            "BK-001",
		TradeDate:            time.Now(),
		SKUID:                "SKU-XYZ",
		TotalWeightMilliTons: 1500, // 1.500 tons
		OriginCountry:        "BR",
		DestinationCountry:   "US",
		EstimatedFreightUSD:  MustMoney(25000000, "USD"), // USD 2500.0000
		TrackingNumber:       "TRK-ABC123",
	}
}

func TestERPShipment_valid(t *testing.T) {
	if err := validERP().Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestERPShipment_empty_tracking_number(t *testing.T) {
	s := validERP()
	s.TrackingNumber = ""
	if err := s.Validate(); err != ErrEmptyTrackingNumber {
		t.Fatalf("expected ErrEmptyTrackingNumber, got %v", err)
	}
}

func TestERPShipment_empty_booking_id(t *testing.T) {
	s := validERP()
	s.BookingID = ""
	if err := s.Validate(); err != ErrEmptyBookingID {
		t.Fatalf("expected ErrEmptyBookingID, got %v", err)
	}
}

func TestERPShipment_zero_weight(t *testing.T) {
	s := validERP()
	s.TotalWeightMilliTons = 0
	if err := s.Validate(); err != ErrZeroWeight {
		t.Fatalf("expected ErrZeroWeight, got %v", err)
	}
}

func TestERPShipment_negative_weight(t *testing.T) {
	s := validERP()
	s.TotalWeightMilliTons = -1
	if err := s.Validate(); err != ErrZeroWeight {
		t.Fatalf("expected ErrZeroWeight, got %v", err)
	}
}

func TestERPShipment_estimated_freight_not_usd(t *testing.T) {
	s := validERP()
	s.EstimatedFreightUSD = MustMoney(25000000, "EUR")
	if err := s.Validate(); err != ErrEstimatedFreightMustBeUSD {
		t.Fatalf("expected ErrEstimatedFreightMustBeUSD, got %v", err)
	}
}

func TestERPShipment_estimated_freight_empty_currency(t *testing.T) {
	s := validERP()
	s.EstimatedFreightUSD = Money{} // zero value: currency == ""
	if err := s.Validate(); err != ErrEstimatedFreightMustBeUSD {
		t.Fatalf("expected ErrEstimatedFreightMustBeUSD, got %v", err)
	}
}

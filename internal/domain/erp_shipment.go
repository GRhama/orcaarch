package domain

import (
	"errors"
	"time"
)

// ERPShipment is the internally authorized dispatch record from the ERP log.
// TotalWeightMilliTons: 1 unit = 0.001 metric tons. 1500 = 1.500 tons.
type ERPShipment struct {
	BookingID           string
	TradeDate           time.Time
	SKUID               string
	TotalWeightMilliTons int64
	OriginCountry       string
	DestinationCountry  string
	EstimatedFreightUSD Money
	TrackingNumber      string
}

var (
	ErrEmptyTrackingNumber       = errors.New("tracking_number must not be empty")
	ErrEmptyBookingID            = errors.New("booking_id must not be empty")
	ErrZeroWeight                = errors.New("total_weight_milli_tons must be positive")
	ErrEstimatedFreightMustBeUSD = errors.New("estimated_freight_usd must be in USD")
)

func (s ERPShipment) Validate() error {
	if s.TrackingNumber == "" {
		return ErrEmptyTrackingNumber
	}
	if s.BookingID == "" {
		return ErrEmptyBookingID
	}
	if s.TotalWeightMilliTons <= 0 {
		return ErrZeroWeight
	}
	if s.EstimatedFreightUSD.Currency() != "USD" {
		return ErrEstimatedFreightMustBeUSD
	}
	return nil
}

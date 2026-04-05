package domain

import "errors"

// Tariff defines pricing for energy consumption at a charger or location.
type Tariff struct {
	ID          string
	PricePerKWh float64
	Currency    string
}

// Validate checks the tariff fields for consistency.
func (t *Tariff) Validate() error {
	if t.ID == "" {
		return errors.New("tariff id is required")
	}
	if t.PricePerKWh < 0 {
		return errors.New("tariff price_per_kwh must not be negative")
	}
	if t.Currency == "" {
		return errors.New("tariff currency is required")
	}
	return nil
}

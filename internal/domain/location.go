package domain

import "errors"

// Location represents a physical site containing one or more chargers.
type Location struct {
	ID        string
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
}

// Validate checks the location fields for consistency.
func (l *Location) Validate() error {
	if l.ID == "" {
		return errors.New("location id is required")
	}
	if l.Name == "" {
		return errors.New("location name is required")
	}
	if l.Address == "" {
		return errors.New("location address is required")
	}
	if l.Latitude < -90 || l.Latitude > 90 {
		return errors.New("location latitude must be between -90 and 90")
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("location longitude must be between -180 and 180")
	}
	return nil
}

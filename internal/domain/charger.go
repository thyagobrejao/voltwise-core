// Package domain defines the core business entities and rules for the VoltWise
// EV charging platform. These types are shared across all services (cloud, OCPP,
// agent) and contain no infrastructure dependencies.
package domain

import (
	"errors"
	"fmt"
)

// ChargerStatus represents the operational state of a charging station.
type ChargerStatus string

const (
	ChargerStatusAvailable ChargerStatus = "available"
	ChargerStatusCharging  ChargerStatus = "charging"
	ChargerStatusOffline   ChargerStatus = "offline"
	ChargerStatusFault     ChargerStatus = "fault"
)

// ChargerProtocol represents the communication protocol used by the charger.
type ChargerProtocol string

const (
	ChargerProtocolOCPP16 ChargerProtocol = "ocpp1.6"
)

// Charger represents a physical EV charging station.
type Charger struct {
	ID         string
	Name       string
	LocationID string
	Status     ChargerStatus
	Protocol   ChargerProtocol
}

// Validate checks the charger fields for consistency.
func (c *Charger) Validate() error {
	if c.ID == "" {
		return errors.New("charger id is required")
	}
	if c.Name == "" {
		return errors.New("charger name is required")
	}
	if c.LocationID == "" {
		return errors.New("charger location_id is required")
	}
	if !c.Status.IsValid() {
		return fmt.Errorf("invalid charger status: %q", c.Status)
	}
	if !c.Protocol.IsValid() {
		return fmt.Errorf("invalid charger protocol: %q", c.Protocol)
	}
	return nil
}

// IsValid returns true if the status is a known charger status.
func (s ChargerStatus) IsValid() bool {
	switch s {
	case ChargerStatusAvailable, ChargerStatusCharging, ChargerStatusOffline, ChargerStatusFault:
		return true
	}
	return false
}

// IsValid returns true if the protocol is a known charger protocol.
func (p ChargerProtocol) IsValid() bool {
	switch p {
	case ChargerProtocolOCPP16:
		return true
	}
	return false
}

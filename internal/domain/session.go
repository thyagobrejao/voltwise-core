package domain

import (
	"errors"
	"fmt"
	"time"
)

// SessionStatus represents the lifecycle state of a charging session.
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusFailed    SessionStatus = "failed"
)

// ChargingSession represents a single charging event on a charger.
type ChargingSession struct {
	ID        string
	ChargerID string
	UserID    string
	StartTime time.Time
	EndTime   *time.Time // nil while session is active
	EnergyKWh float64
	Status    SessionStatus
}

// Validate checks the session fields for consistency.
func (s *ChargingSession) Validate() error {
	if s.ID == "" {
		return errors.New("session id is required")
	}
	if s.ChargerID == "" {
		return errors.New("session charger_id is required")
	}
	if s.UserID == "" {
		return errors.New("session user_id is required")
	}
	if s.StartTime.IsZero() {
		return errors.New("session start_time is required")
	}
	if s.EndTime != nil && s.EndTime.Before(s.StartTime) {
		return errors.New("session end_time must be after start_time")
	}
	if s.EnergyKWh < 0 {
		return errors.New("session energy_kwh must not be negative")
	}
	if !s.Status.IsValid() {
		return fmt.Errorf("invalid session status: %q", s.Status)
	}
	return nil
}

// IsActive returns true if the session is currently in progress.
func (s *ChargingSession) IsActive() bool {
	return s.Status == SessionStatusActive
}

// IsValid returns true if the status is a known session status.
func (s SessionStatus) IsValid() bool {
	switch s {
	case SessionStatusActive, SessionStatusCompleted, SessionStatusFailed:
		return true
	}
	return false
}

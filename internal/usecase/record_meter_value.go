package usecase

import (
	"context"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
	"github.com/thyagobrejao/voltwise-core/internal/ports"
)

// RecordMeterValue updates the cumulative energy reading on the active session
// for a charger. This is typically called when the charger reports a MeterValues
// message via OCPP.
type RecordMeterValue struct {
	Sessions ports.SessionRepository
}

// MeterValueInput carries the meter reading to record.
type MeterValueInput struct {
	ChargerID string
	EnergyKWh float64 // cumulative energy delivered so far
}

// Execute updates the energy consumption on the active session.
func (uc *RecordMeterValue) Execute(ctx context.Context, in MeterValueInput) (*domain.ChargingSession, error) {
	session, err := uc.Sessions.FindActiveByChargerID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrNoActiveSession
	}

	session.EnergyKWh = in.EnergyKWh

	if err := session.Validate(); err != nil {
		return nil, err
	}
	if err := uc.Sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

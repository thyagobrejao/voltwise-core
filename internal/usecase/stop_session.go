package usecase

import (
	"context"
	"time"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
	"github.com/thyagobrejao/voltwise-core/internal/ports"
)

// StopChargingSession finalises an active charging session on a charger and
// transitions the charger back to "available".
type StopChargingSession struct {
	Chargers ports.ChargerRepository
	Sessions ports.SessionRepository
}

// StopSessionInput carries the parameters needed to stop a session.
type StopSessionInput struct {
	ChargerID string
}

// Execute stops the active session for the given charger.
func (uc *StopChargingSession) Execute(ctx context.Context, in StopSessionInput) (*domain.ChargingSession, error) {
	charger, err := uc.Chargers.FindByID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}

	session, err := uc.Sessions.FindActiveByChargerID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrNoActiveSession
	}

	// Complete the session.
	now := time.Now().UTC()
	session.EndTime = &now
	session.Status = domain.SessionStatusCompleted

	if err := session.Validate(); err != nil {
		return nil, err
	}
	if err := uc.Sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	// Charger becomes available again.
	charger.Status = domain.ChargerStatusAvailable
	if err := uc.Chargers.Save(ctx, charger); err != nil {
		return nil, err
	}

	return session, nil
}

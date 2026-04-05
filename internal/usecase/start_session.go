// Package usecase implements application-level business operations. Each use case
// orchestrates domain entities and repository ports without depending on any
// concrete infrastructure.
package usecase

import (
	"context"
	"time"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
	"github.com/thyagobrejao/voltwise-core/internal/ports"
)

// StartChargingSession initiates a new charging session on a charger.
// It enforces the rule that a charger can only have one active session at a time
// and transitions the charger status to "charging".
type StartChargingSession struct {
	Chargers ports.ChargerRepository
	Sessions ports.SessionRepository
	Users    ports.UserRepository
}

// StartSessionInput carries the parameters needed to begin a session.
type StartSessionInput struct {
	SessionID string
	ChargerID string
	UserID    string
}

// Execute starts a charging session after validating preconditions.
func (uc *StartChargingSession) Execute(ctx context.Context, in StartSessionInput) (*domain.ChargingSession, error) {
	// Verify the charger exists.
	charger, err := uc.Chargers.FindByID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}

	// Verify the user exists.
	if _, err := uc.Users.FindByID(ctx, in.UserID); err != nil {
		return nil, err
	}

	// Business rule: a charger can only have one active session.
	existing, err := uc.Sessions.FindActiveByChargerID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrSessionAlreadyActive
	}

	// Create the session.
	session := &domain.ChargingSession{
		ID:        in.SessionID,
		ChargerID: in.ChargerID,
		UserID:    in.UserID,
		StartTime: time.Now().UTC(),
		Status:    domain.SessionStatusActive,
	}
	if err := session.Validate(); err != nil {
		return nil, err
	}
	if err := uc.Sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	// Reflect session state on the charger.
	charger.Status = domain.ChargerStatusCharging
	if err := uc.Chargers.Save(ctx, charger); err != nil {
		return nil, err
	}

	return session, nil
}

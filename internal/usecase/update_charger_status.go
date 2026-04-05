package usecase

import (
	"context"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
	"github.com/thyagobrejao/voltwise-core/internal/ports"
)

// UpdateChargerStatus changes the operational status of a charger.
// It guards against invalid transitions—for example, a charger that is actively
// charging should not be set to "available" without stopping the session first.
type UpdateChargerStatus struct {
	Chargers ports.ChargerRepository
	Sessions ports.SessionRepository
}

// UpdateStatusInput carries the parameters for a status change.
type UpdateStatusInput struct {
	ChargerID string
	Status    domain.ChargerStatus
}

// Execute updates the charger status after validating the transition.
func (uc *UpdateChargerStatus) Execute(ctx context.Context, in UpdateStatusInput) (*domain.Charger, error) {
	if !in.Status.IsValid() {
		return nil, domain.ErrInvalidChargerStatus
	}

	charger, err := uc.Chargers.FindByID(ctx, in.ChargerID)
	if err != nil {
		return nil, err
	}

	// Don't allow switching away from "charging" if an active session exists;
	// the session must be stopped explicitly first.
	if charger.Status == domain.ChargerStatusCharging && in.Status == domain.ChargerStatusAvailable {
		active, err := uc.Sessions.FindActiveByChargerID(ctx, in.ChargerID)
		if err != nil {
			return nil, err
		}
		if active != nil {
			return nil, domain.ErrSessionAlreadyActive
		}
	}

	charger.Status = in.Status
	if err := uc.Chargers.Save(ctx, charger); err != nil {
		return nil, err
	}

	return charger, nil
}

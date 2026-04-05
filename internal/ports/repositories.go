// Package ports defines the interfaces (driven side) that use cases depend on.
// Infrastructure adapters (database, message broker, etc.) implement these
// interfaces, keeping the domain and use-case layers free of external concerns.
package ports

import (
	"context"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
)

// ChargerRepository provides access to charger persistence.
type ChargerRepository interface {
	// FindByID retrieves a charger by its unique identifier.
	FindByID(ctx context.Context, id string) (*domain.Charger, error)

	// Save persists a charger (create or update).
	Save(ctx context.Context, charger *domain.Charger) error
}

// SessionRepository provides access to charging session persistence.
type SessionRepository interface {
	// FindByID retrieves a session by its unique identifier.
	FindByID(ctx context.Context, id string) (*domain.ChargingSession, error)

	// FindActiveByChargerID returns the active session for a charger, or nil if none.
	FindActiveByChargerID(ctx context.Context, chargerID string) (*domain.ChargingSession, error)

	// Save persists a session (create or update).
	Save(ctx context.Context, session *domain.ChargingSession) error
}

// UserRepository provides access to user persistence.
type UserRepository interface {
	// FindByID retrieves a user by its unique identifier.
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

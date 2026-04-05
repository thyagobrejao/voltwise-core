package domain

import "errors"

// Domain-level sentinel errors used across use cases and services.
// These allow callers to match on specific error conditions without
// coupling to implementation details.
var (
	ErrChargerNotFound        = errors.New("charger not found")
	ErrSessionNotFound        = errors.New("session not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrSessionAlreadyActive   = errors.New("charger already has an active session")
	ErrNoActiveSession        = errors.New("no active session for charger")
	ErrInvalidChargerStatus   = errors.New("invalid charger status transition")
)

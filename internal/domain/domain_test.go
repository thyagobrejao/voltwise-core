package domain_test

import (
	"testing"
	"time"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
)

// --- Charger validation ---

func TestCharger_Validate_Valid(t *testing.T) {
	c := &domain.Charger{
		ID:         "chg-1",
		Name:       "Charger A",
		LocationID: "loc-1",
		Status:     domain.ChargerStatusAvailable,
		Protocol:   domain.ChargerProtocolOCPP16,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid charger, got error: %v", err)
	}
}

func TestCharger_Validate_MissingID(t *testing.T) {
	c := &domain.Charger{Name: "Charger A", LocationID: "loc-1", Status: domain.ChargerStatusAvailable, Protocol: domain.ChargerProtocolOCPP16}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestCharger_Validate_InvalidStatus(t *testing.T) {
	c := &domain.Charger{ID: "chg-1", Name: "Charger A", LocationID: "loc-1", Status: "unknown", Protocol: domain.ChargerProtocolOCPP16}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid status")
	}
}

// --- ChargingSession validation ---

func TestSession_Validate_Valid(t *testing.T) {
	s := &domain.ChargingSession{
		ID:        "ses-1",
		ChargerID: "chg-1",
		UserID:    "usr-1",
		StartTime: time.Now(),
		Status:    domain.SessionStatusActive,
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("expected valid session, got error: %v", err)
	}
}

func TestSession_Validate_EndBeforeStart(t *testing.T) {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(-time.Hour)
	s := &domain.ChargingSession{
		ID:        "ses-1",
		ChargerID: "chg-1",
		UserID:    "usr-1",
		StartTime: start,
		EndTime:   &end,
		Status:    domain.SessionStatusCompleted,
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error when end_time is before start_time")
	}
}

func TestSession_Validate_NegativeEnergy(t *testing.T) {
	s := &domain.ChargingSession{
		ID:        "ses-1",
		ChargerID: "chg-1",
		UserID:    "usr-1",
		StartTime: time.Now(),
		EnergyKWh: -5,
		Status:    domain.SessionStatusActive,
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for negative energy")
	}
}

// --- User validation ---

func TestUser_Validate_Valid(t *testing.T) {
	u := &domain.User{ID: "usr-1", Name: "Alice", Email: "alice@example.com"}
	if err := u.Validate(); err != nil {
		t.Fatalf("expected valid user, got error: %v", err)
	}
}

func TestUser_Validate_InvalidEmail(t *testing.T) {
	u := &domain.User{ID: "usr-1", Name: "Alice", Email: "not-an-email"}
	if err := u.Validate(); err == nil {
		t.Fatal("expected error for invalid email")
	}
}

// --- Location validation ---

func TestLocation_Validate_Valid(t *testing.T) {
	l := &domain.Location{ID: "loc-1", Name: "HQ", Address: "123 Main St", Latitude: 40.71, Longitude: -74.00}
	if err := l.Validate(); err != nil {
		t.Fatalf("expected valid location, got error: %v", err)
	}
}

func TestLocation_Validate_LatitudeOutOfRange(t *testing.T) {
	l := &domain.Location{ID: "loc-1", Name: "HQ", Address: "123 Main St", Latitude: 91, Longitude: 0}
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for latitude > 90")
	}
}

// --- Tariff validation ---

func TestTariff_Validate_Valid(t *testing.T) {
	tr := &domain.Tariff{ID: "tar-1", PricePerKWh: 0.30, Currency: "USD"}
	if err := tr.Validate(); err != nil {
		t.Fatalf("expected valid tariff, got error: %v", err)
	}
}

func TestTariff_Validate_NegativePrice(t *testing.T) {
	tr := &domain.Tariff{ID: "tar-1", PricePerKWh: -1, Currency: "USD"}
	if err := tr.Validate(); err == nil {
		t.Fatal("expected error for negative price")
	}
}

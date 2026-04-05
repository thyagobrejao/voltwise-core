package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/thyagobrejao/voltwise-core/internal/domain"
	"github.com/thyagobrejao/voltwise-core/internal/usecase"
)

// --- In-memory test doubles ---

// memChargerRepo is a minimal in-memory ChargerRepository for tests.
type memChargerRepo struct {
	chargers map[string]*domain.Charger
}

func newMemChargerRepo(chargers ...*domain.Charger) *memChargerRepo {
	m := &memChargerRepo{chargers: make(map[string]*domain.Charger)}
	for _, c := range chargers {
		m.chargers[c.ID] = c
	}
	return m
}

func (r *memChargerRepo) FindByID(_ context.Context, id string) (*domain.Charger, error) {
	c, ok := r.chargers[id]
	if !ok {
		return nil, domain.ErrChargerNotFound
	}
	return c, nil
}

func (r *memChargerRepo) Save(_ context.Context, charger *domain.Charger) error {
	r.chargers[charger.ID] = charger
	return nil
}

// memSessionRepo is a minimal in-memory SessionRepository for tests.
type memSessionRepo struct {
	sessions map[string]*domain.ChargingSession
}

func newMemSessionRepo(sessions ...*domain.ChargingSession) *memSessionRepo {
	m := &memSessionRepo{sessions: make(map[string]*domain.ChargingSession)}
	for _, s := range sessions {
		m.sessions[s.ID] = s
	}
	return m
}

func (r *memSessionRepo) FindByID(_ context.Context, id string) (*domain.ChargingSession, error) {
	s, ok := r.sessions[id]
	if !ok {
		return nil, domain.ErrSessionNotFound
	}
	return s, nil
}

func (r *memSessionRepo) FindActiveByChargerID(_ context.Context, chargerID string) (*domain.ChargingSession, error) {
	for _, s := range r.sessions {
		if s.ChargerID == chargerID && s.Status == domain.SessionStatusActive {
			return s, nil
		}
	}
	return nil, nil
}

func (r *memSessionRepo) Save(_ context.Context, session *domain.ChargingSession) error {
	r.sessions[session.ID] = session
	return nil
}

// memUserRepo is a minimal in-memory UserRepository for tests.
type memUserRepo struct {
	users map[string]*domain.User
}

func newMemUserRepo(users ...*domain.User) *memUserRepo {
	m := &memUserRepo{users: make(map[string]*domain.User)}
	for _, u := range users {
		m.users[u.ID] = u
	}
	return m
}

func (r *memUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

// --- StartChargingSession tests ---

func TestStartSession_Success(t *testing.T) {
	charger := &domain.Charger{
		ID: "chg-1", Name: "Charger A", LocationID: "loc-1",
		Status: domain.ChargerStatusAvailable, Protocol: domain.ChargerProtocolOCPP16,
	}
	user := &domain.User{ID: "usr-1", Name: "Alice", Email: "alice@example.com"}

	uc := &usecase.StartChargingSession{
		Chargers: newMemChargerRepo(charger),
		Sessions: newMemSessionRepo(),
		Users:    newMemUserRepo(user),
	}

	session, err := uc.Execute(context.Background(), usecase.StartSessionInput{
		SessionID: "ses-1",
		ChargerID: "chg-1",
		UserID:    "usr-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.Status != domain.SessionStatusActive {
		t.Errorf("expected session status %q, got %q", domain.SessionStatusActive, session.Status)
	}
	if charger.Status != domain.ChargerStatusCharging {
		t.Errorf("expected charger status %q, got %q", domain.ChargerStatusCharging, charger.Status)
	}
}

func TestStartSession_DuplicateActiveSession(t *testing.T) {
	charger := &domain.Charger{
		ID: "chg-1", Name: "Charger A", LocationID: "loc-1",
		Status: domain.ChargerStatusCharging, Protocol: domain.ChargerProtocolOCPP16,
	}
	user := &domain.User{ID: "usr-1", Name: "Alice", Email: "alice@example.com"}
	existing := &domain.ChargingSession{
		ID: "ses-0", ChargerID: "chg-1", UserID: "usr-1",
		StartTime: time.Now(), Status: domain.SessionStatusActive,
	}

	uc := &usecase.StartChargingSession{
		Chargers: newMemChargerRepo(charger),
		Sessions: newMemSessionRepo(existing),
		Users:    newMemUserRepo(user),
	}

	_, err := uc.Execute(context.Background(), usecase.StartSessionInput{
		SessionID: "ses-1",
		ChargerID: "chg-1",
		UserID:    "usr-1",
	})
	if !errors.Is(err, domain.ErrSessionAlreadyActive) {
		t.Fatalf("expected ErrSessionAlreadyActive, got: %v", err)
	}
}

func TestStartSession_ChargerNotFound(t *testing.T) {
	uc := &usecase.StartChargingSession{
		Chargers: newMemChargerRepo(),
		Sessions: newMemSessionRepo(),
		Users:    newMemUserRepo(),
	}

	_, err := uc.Execute(context.Background(), usecase.StartSessionInput{
		SessionID: "ses-1",
		ChargerID: "chg-missing",
		UserID:    "usr-1",
	})
	if !errors.Is(err, domain.ErrChargerNotFound) {
		t.Fatalf("expected ErrChargerNotFound, got: %v", err)
	}
}

// --- StopChargingSession tests ---

func TestStopSession_Success(t *testing.T) {
	charger := &domain.Charger{
		ID: "chg-1", Name: "Charger A", LocationID: "loc-1",
		Status: domain.ChargerStatusCharging, Protocol: domain.ChargerProtocolOCPP16,
	}
	active := &domain.ChargingSession{
		ID: "ses-1", ChargerID: "chg-1", UserID: "usr-1",
		StartTime: time.Now().Add(-10 * time.Minute),
		Status:    domain.SessionStatusActive,
	}

	uc := &usecase.StopChargingSession{
		Chargers: newMemChargerRepo(charger),
		Sessions: newMemSessionRepo(active),
	}

	session, err := uc.Execute(context.Background(), usecase.StopSessionInput{
		ChargerID: "chg-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.Status != domain.SessionStatusCompleted {
		t.Errorf("expected session status %q, got %q", domain.SessionStatusCompleted, session.Status)
	}
	if session.EndTime == nil {
		t.Error("expected end_time to be set")
	}
	if charger.Status != domain.ChargerStatusAvailable {
		t.Errorf("expected charger status %q, got %q", domain.ChargerStatusAvailable, charger.Status)
	}
}

func TestStopSession_NoActiveSession(t *testing.T) {
	charger := &domain.Charger{
		ID: "chg-1", Name: "Charger A", LocationID: "loc-1",
		Status: domain.ChargerStatusAvailable, Protocol: domain.ChargerProtocolOCPP16,
	}

	uc := &usecase.StopChargingSession{
		Chargers: newMemChargerRepo(charger),
		Sessions: newMemSessionRepo(),
	}

	_, err := uc.Execute(context.Background(), usecase.StopSessionInput{
		ChargerID: "chg-1",
	})
	if !errors.Is(err, domain.ErrNoActiveSession) {
		t.Fatalf("expected ErrNoActiveSession, got: %v", err)
	}
}

// --- RecordMeterValue tests ---

func TestRecordMeterValue_Success(t *testing.T) {
	active := &domain.ChargingSession{
		ID: "ses-1", ChargerID: "chg-1", UserID: "usr-1",
		StartTime: time.Now().Add(-5 * time.Minute),
		Status:    domain.SessionStatusActive,
	}

	uc := &usecase.RecordMeterValue{
		Sessions: newMemSessionRepo(active),
	}

	session, err := uc.Execute(context.Background(), usecase.MeterValueInput{
		ChargerID: "chg-1",
		EnergyKWh: 12.5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.EnergyKWh != 12.5 {
		t.Errorf("expected energy_kwh 12.5, got %f", session.EnergyKWh)
	}
}

func TestRecordMeterValue_NoActiveSession(t *testing.T) {
	uc := &usecase.RecordMeterValue{
		Sessions: newMemSessionRepo(),
	}

	_, err := uc.Execute(context.Background(), usecase.MeterValueInput{
		ChargerID: "chg-1",
		EnergyKWh: 5.0,
	})
	if !errors.Is(err, domain.ErrNoActiveSession) {
		t.Fatalf("expected ErrNoActiveSession, got: %v", err)
	}
}

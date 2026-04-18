package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service is the application service for session management.
type Service struct {
	store          Store
	sessionTTL     time.Duration
	maxConcurrency int // per-event cap; 0 means unlimited
}

func NewService(store Store, sessionTTL time.Duration, maxConcurrency int) *Service {
	return &Service{
		store:          store,
		sessionTTL:     sessionTTL,
		maxConcurrency: maxConcurrency,
	}
}

// CreateSession creates a new seat-map session for the user.
// Idempotent: returns existing session if one already exists for this user+event.
func (s *Service) CreateSession(ctx context.Context, eventID, userID string) (string, error) {
	// Check if user already has a session
	existing, err := s.store.FindByUserEvent(ctx, userID, eventID)
	if err != nil {
		return "", err
	}
	if existing != "" {
		return existing, nil
	}

	// Check concurrency cap
	if s.maxConcurrency > 0 {
		count, err := s.store.GetActiveCount(ctx, eventID)
		if err != nil {
			return "", err
		}
		if count >= s.maxConcurrency {
			return "", fmt.Errorf("seat map is full, please wait")
		}
	}

	sess := &Session{
		ID:        uuid.NewString(),
		EventID:   eventID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := s.store.Create(ctx, sess, s.sessionTTL); err != nil {
		return "", err
	}
	return sess.ID, nil
}

// RefreshSession extends the session TTL (heartbeat).
func (s *Service) RefreshSession(ctx context.Context, sessionID, eventID string) error {
	return s.store.Refresh(ctx, sessionID, eventID, s.sessionTTL)
}

// EndSession explicitly ends a session.
func (s *Service) EndSession(ctx context.Context, sessionID, eventID string) error {
	return s.store.End(ctx, sessionID, eventID)
}

// ValidateSession checks that a session exists and belongs to the given user.
func (s *Service) ValidateSession(ctx context.Context, sessionID, userID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	// We need to scan for the session across events since we don't have eventID here.
	// The caller (ReserveTicket) doesn't know the eventID at the handler level,
	// but the session key includes it. We use the user-event mapping in reverse:
	// the sessionID is stored as value, so we can't look up by sessionID alone easily.
	// Solution: accept that ValidateSession is called with eventID derived from the ticket.
	return fmt.Errorf("use ValidateSessionForEvent instead")
}

// ValidateSessionForEvent checks that a session exists for the user on the given event.
func (s *Service) ValidateSessionForEvent(ctx context.Context, sessionID, eventID, userID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	sess, err := s.store.Find(ctx, sessionID, eventID)
	if err != nil {
		return fmt.Errorf("invalid or expired session")
	}
	if sess.UserID != userID {
		return fmt.Errorf("session does not belong to this user")
	}
	return nil
}

// HeartbeatIntervalMs returns the recommended heartbeat interval for clients.
func (s *Service) HeartbeatIntervalMs() int {
	return int(s.sessionTTL.Milliseconds()) / 3
}

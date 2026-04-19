package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const estimatedSessionDuration = 60 // seconds, used for wait time estimation

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
// If the concurrency cap is reached, the user is enqueued and a QueuedError is returned.
func (s *Service) CreateSession(ctx context.Context, eventID, userID string) (string, error) {
	// Check if user already has a session
	existing, err := s.store.FindByUserEvent(ctx, userID, eventID)
	if err != nil {
		return "", err
	}
	if existing != "" {
		return existing, nil
	}

	// Check if user has been admitted from the queue
	admitted, err := s.store.IsAdmitted(ctx, eventID, userID)
	if err != nil {
		return "", err
	}
	if admitted {
		if err := s.store.ClearAdmission(ctx, eventID, userID); err != nil {
			return "", err
		}
		return s.createSessionInternal(ctx, eventID, userID)
	}

	// Check concurrency cap
	if s.maxConcurrency > 0 {
		count, err := s.store.GetActiveCount(ctx, eventID)
		if err != nil {
			return "", err
		}
		if count >= s.maxConcurrency {
			return "", s.enqueueUser(ctx, eventID, userID)
		}
	}

	return s.createSessionInternal(ctx, eventID, userID)
}

func (s *Service) createSessionInternal(ctx context.Context, eventID, userID string) (string, error) {
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

func (s *Service) enqueueUser(ctx context.Context, eventID, userID string) *QueuedError {
	s.store.EnqueueUser(ctx, eventID, userID)
	pos, _ := s.store.GetQueuePosition(ctx, eventID, userID)
	length, _ := s.store.GetQueueLength(ctx, eventID)
	return &QueuedError{
		Position:      pos + 1, // 1-based
		EstimatedWait: (pos + 1) * estimatedSessionDuration / max(s.maxConcurrency, 1),
		QueueLength:   length,
	}
}

// RefreshSession extends the session TTL (heartbeat).
func (s *Service) RefreshSession(ctx context.Context, sessionID, eventID string) error {
	return s.store.Refresh(ctx, sessionID, eventID, s.sessionTTL)
}

// EndSession explicitly ends a session, then admits the next queued user.
func (s *Service) EndSession(ctx context.Context, sessionID, eventID string) error {
	if err := s.store.End(ctx, sessionID, eventID); err != nil {
		return err
	}
	s.AdmitNext(ctx, eventID)
	return nil
}

// ValidateSession checks that a session exists and belongs to the given user.
func (s *Service) ValidateSession(ctx context.Context, sessionID, userID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
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

// GetQueueStatus returns the queue status for a user.
func (s *Service) GetQueueStatus(ctx context.Context, eventID, userID string) (*QueueStatus, error) {
	// Check if admitted
	admitted, err := s.store.IsAdmitted(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if admitted {
		return &QueueStatus{Status: "admitted"}, nil
	}

	// Check queue position
	pos, err := s.store.GetQueuePosition(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if pos < 0 {
		return &QueueStatus{Status: "none"}, nil
	}

	length, _ := s.store.GetQueueLength(ctx, eventID)
	return &QueueStatus{
		Status:        "waiting",
		Position:      pos + 1,
		EstimatedWait: (pos + 1) * estimatedSessionDuration / max(s.maxConcurrency, 1),
		QueueLength:   length,
	}, nil
}

// AdmitNext admits the next user from the queue if capacity is available.
func (s *Service) AdmitNext(ctx context.Context, eventID string) {
	if s.maxConcurrency <= 0 {
		return
	}
	s.store.AdmitNextUser(ctx, eventID, s.maxConcurrency)
}

// LeaveQueue removes a user from the waiting queue.
func (s *Service) LeaveQueue(ctx context.Context, eventID, userID string) error {
	return s.store.DequeueUser(ctx, eventID, userID)
}

// HeartbeatIntervalMs returns the recommended heartbeat interval for clients.
func (s *Service) HeartbeatIntervalMs() int {
	return int(s.sessionTTL.Milliseconds()) / 3
}

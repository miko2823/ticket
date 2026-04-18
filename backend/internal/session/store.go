package session

import (
	"context"
	"time"
)

// Store defines the port for session persistence (Redis adapter).
type Store interface {
	Create(ctx context.Context, session *Session, ttl time.Duration) error
	Refresh(ctx context.Context, sessionID, eventID string, ttl time.Duration) error
	End(ctx context.Context, sessionID, eventID string) error
	Find(ctx context.Context, sessionID, eventID string) (*Session, error)
	FindByUserEvent(ctx context.Context, userID, eventID string) (string, error)
	GetActiveCount(ctx context.Context, eventID string) (int, error)
}

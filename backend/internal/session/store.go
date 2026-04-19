package session

import (
	"context"
	"time"
)

// Store defines the port for session persistence (Redis adapter).
type Store interface {
	// Session operations
	Create(ctx context.Context, session *Session, ttl time.Duration) error
	Refresh(ctx context.Context, sessionID, eventID string, ttl time.Duration) error
	End(ctx context.Context, sessionID, eventID string) error
	Find(ctx context.Context, sessionID, eventID string) (*Session, error)
	FindByUserEvent(ctx context.Context, userID, eventID string) (string, error)
	GetActiveCount(ctx context.Context, eventID string) (int, error)

	// Queue operations
	EnqueueUser(ctx context.Context, eventID, userID string) error
	DequeueUser(ctx context.Context, eventID, userID string) error
	GetQueuePosition(ctx context.Context, eventID, userID string) (int, error) // -1 if not in queue
	GetQueueLength(ctx context.Context, eventID string) (int, error)
	AdmitNextUser(ctx context.Context, eventID string, maxConcurrency int) (string, error) // returns admitted userID or ""
	IsAdmitted(ctx context.Context, eventID, userID string) (bool, error)
	ClearAdmission(ctx context.Context, eventID, userID string) error
}

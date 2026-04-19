package session

import "time"

// Session represents an active user session on the seat map page.
type Session struct {
	ID        string
	EventID   string
	UserID    string
	CreatedAt time.Time
}

// QueuedError is returned when a user is placed in the waiting room queue.
type QueuedError struct {
	Position      int // 1-based
	EstimatedWait int // seconds
	QueueLength   int
}

func (e *QueuedError) Error() string { return "queued" }

// QueueStatus represents the current queue state for a user.
type QueueStatus struct {
	Status        string `json:"status"` // "waiting", "admitted", "none"
	Position      int    `json:"position,omitempty"`
	EstimatedWait int    `json:"estimated_wait_seconds,omitempty"`
	QueueLength   int    `json:"queue_length,omitempty"`
}

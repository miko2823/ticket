package session

import "time"

// Session represents an active user session on the seat map page.
type Session struct {
	ID        string
	EventID   string
	UserID    string
	CreatedAt time.Time
}

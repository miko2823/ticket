package event

import (
	"context"
	"time"
)

// Repository defines the port for Event persistence.
type Repository interface {
	FindAll(ctx context.Context) ([]Event, error)
	FindByID(ctx context.Context, id string) (*Event, error)
	FindTicketsByEventID(ctx context.Context, eventID string) ([]Ticket, error)
	FindTicketByID(ctx context.Context, id string) (*Ticket, error)
	ReserveTicket(ctx context.Context, ticketID string, currentVersion int, reservedUntil time.Time) error
	ReleaseExpiredReservations(ctx context.Context, now time.Time) error
}

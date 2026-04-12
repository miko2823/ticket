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
	ReserveTicket(ctx context.Context, ticketID string, currentVersion int, userID string, reservedUntil time.Time) error
	UpdateTicketStatus(ctx context.Context, ticketID string, status TicketStatus) error
	ReleaseExpiredReservations(ctx context.Context, now time.Time) error
	FindSeatLayoutByEventID(ctx context.Context, eventID string) (*SeatLayout, error)
}

package booking

import "time"

// Booking is the aggregate root for the Booking bounded context.
type Booking struct {
	ID        string
	UserID    string
	TicketID  string
	Status    BookingStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BookingStatus is a value object for booking state.
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusFailed    BookingStatus = "failed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// CanBeCancelled returns whether this booking is in a cancellable state.
func (b *Booking) CanBeCancelled() bool {
	return b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed
}

// Domain Events

// BookingConfirmed is published when a booking is successfully confirmed.
type BookingConfirmed struct {
	BookingID string
	UserID    string
	TicketID  string
	At        time.Time
}

// TicketReserved is published when a ticket is reserved for a booking.
type TicketReserved struct {
	TicketID  string
	BookingID string
	At        time.Time
}

// FailedBooking represents a booking that failed after retries (dead letter).
type FailedBooking struct {
	ID        string
	BookingID string
	Reason    string
	FailedAt  time.Time
}

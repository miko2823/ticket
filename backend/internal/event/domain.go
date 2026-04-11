package event

import (
	"time"
)

// Event is the aggregate root for the Event bounded context.
type Event struct {
	ID                string
	Name              string
	Venue             string
	StartsAt          time.Time
	TicketingStartsAt time.Time
	TicketingEndsAt   time.Time
	CreatedAt         time.Time
}

// IsTicketingOpen returns true if the current time is within the ticketing window.
func (e *Event) IsTicketingOpen(now time.Time) bool {
	return !now.Before(e.TicketingStartsAt) && now.Before(e.TicketingEndsAt)
}

// Ticket is an entity belonging to the Event aggregate.
type Ticket struct {
	ID            string
	EventID       string
	SeatLabel     SeatLabel
	PriceJPY      Price
	Status        TicketStatus
	Version       int
	ReservedBy    *string
	ReservedUntil *time.Time
	CreatedAt     time.Time
}

// IsReservedBy returns true if the ticket is reserved by the given user.
func (t *Ticket) IsReservedBy(userID string) bool {
	return t.Status == TicketStatusReserved && t.ReservedBy != nil && *t.ReservedBy == userID
}

// IsAvailable returns true if the ticket can be reserved.
// A ticket is available if its status is "available", or if it was
// reserved but the reservation has expired.
func (t *Ticket) IsAvailable(now time.Time) bool {
	if t.Status == TicketStatusAvailable {
		return true
	}
	if t.Status == TicketStatusReserved && t.ReservedUntil != nil && now.After(*t.ReservedUntil) {
		return true
	}
	return false
}

const ReservationDuration = 5 * time.Minute

// SeatLabel is a value object representing a seat identifier (e.g. "A-12").
type SeatLabel string

// Price is a value object representing a price in JPY.
type Price int

// TicketStatus is a value object for ticket state.
type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusReserved  TicketStatus = "reserved"
	TicketStatusSold      TicketStatus = "sold"
	TicketStatusCancelled TicketStatus = "cancelled"
)

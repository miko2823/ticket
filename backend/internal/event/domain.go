package event

import (
	"time"
)

// Event is the aggregate root for the Event bounded context.
type Event struct {
	ID        string
	Name      string
	Venue     string
	StartsAt  time.Time
	CreatedAt time.Time
}

// Ticket is an entity belonging to the Event aggregate.
type Ticket struct {
	ID        string
	EventID   string
	SeatLabel SeatLabel
	PriceJPY  Price
	Status    TicketStatus
	Version   int
	CreatedAt time.Time
}

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

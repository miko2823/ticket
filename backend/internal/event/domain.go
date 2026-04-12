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

// SeatLayout is a value object holding the visual seat map configuration.
type SeatLayout struct {
	Canvas   CanvasSize      `json:"canvas"`
	Stage    StageConfig     `json:"stage"`
	Sections []SectionConfig `json:"sections"`
	Seats    []SeatPosition  `json:"seats"`
}

// CanvasSize defines the coordinate space for the seat map.
type CanvasSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// StageConfig defines the stage area drawn on the seat map.
type StageConfig struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Label  string `json:"label"`
}

// SectionConfig defines a named section with a display color.
type SectionConfig struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

// SeatPosition defines where a seat is drawn on the canvas.
type SeatPosition struct {
	Label   string `json:"label"`
	Section string `json:"section"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	R       int    `json:"r"`
}

// SeatMapSeat is a seat with both visual position and live ticket data merged.
type SeatMapSeat struct {
	TicketID string `json:"ticket_id"`
	Label    string `json:"label"`
	Section  string `json:"section"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	R        int    `json:"r"`
	PriceJPY int    `json:"price_jpy"`
	Status   string `json:"status"`
}

// SeatMapLayout holds the static parts of the seat map response.
type SeatMapLayout struct {
	Canvas   CanvasSize      `json:"canvas"`
	Stage    StageConfig     `json:"stage"`
	Sections []SectionConfig `json:"sections"`
}

// SeatMapResponse is the merged API response for the seat map endpoint.
type SeatMapResponse struct {
	EventID string        `json:"event_id"`
	Layout  SeatMapLayout `json:"layout"`
	Seats   []SeatMapSeat `json:"seats"`
}

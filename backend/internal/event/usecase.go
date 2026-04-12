package event

import (
	"context"
	"fmt"
	"time"
)

// UseCase is the application service for the Event bounded context.
type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) ListEvents(ctx context.Context) ([]Event, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *UseCase) GetEvent(ctx context.Context, id string) (*Event, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) GetTickets(ctx context.Context, eventID string) ([]Ticket, error) {
	return uc.repo.FindTicketsByEventID(ctx, eventID)
}

func (uc *UseCase) GetTicket(ctx context.Context, id string) (*Ticket, error) {
	return uc.repo.FindTicketByID(ctx, id)
}

// ReserveTicket checks availability and ticketing window, then reserves a ticket for 5 minutes.
func (uc *UseCase) ReserveTicket(ctx context.Context, ticketID, userID string) (*Ticket, error) {
	ticket, err := uc.repo.FindTicketByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found")
	}

	now := time.Now()

	// Check ticketing window
	ev, err := uc.repo.FindByID(ctx, ticket.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}
	if !ev.IsTicketingOpen(now) {
		return nil, fmt.Errorf("ticketing is not open")
	}

	if !ticket.IsAvailable(now) {
		return nil, fmt.Errorf("ticket not available")
	}

	reservedUntil := now.Add(ReservationDuration)
	if err := uc.repo.ReserveTicket(ctx, ticketID, ticket.Version, userID, reservedUntil); err != nil {
		return nil, err
	}

	ticket.Status = TicketStatusReserved
	ticket.ReservedBy = &userID
	ticket.ReservedUntil = &reservedUntil
	return ticket, nil
}

// ReleaseTicket releases a reservation if the user owns it.
func (uc *UseCase) ReleaseTicket(ctx context.Context, ticketID, userID string) error {
	ticket, err := uc.repo.FindTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found")
	}

	if ticket.Status != TicketStatusReserved {
		return fmt.Errorf("ticket is not reserved")
	}
	if !ticket.IsReservedBy(userID) {
		return fmt.Errorf("ticket is not reserved by you")
	}

	return uc.repo.UpdateTicketStatus(ctx, ticketID, TicketStatusAvailable)
}

// GetSeatMap returns the seat layout merged with live ticket availability.
func (uc *UseCase) GetSeatMap(ctx context.Context, eventID string) (*SeatMapResponse, error) {
	layout, err := uc.repo.FindSeatLayoutByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}
	if layout == nil {
		return nil, fmt.Errorf("no seat layout for this event")
	}

	tickets, err := uc.repo.FindTicketsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Build lookup: seat_label -> ticket
	ticketMap := make(map[string]*Ticket, len(tickets))
	for i := range tickets {
		ticketMap[string(tickets[i].SeatLabel)] = &tickets[i]
	}

	// Merge layout seats with live ticket data
	seats := make([]SeatMapSeat, 0, len(layout.Seats))
	for _, sp := range layout.Seats {
		s := SeatMapSeat{
			Label:   sp.Label,
			Section: sp.Section,
			X:       sp.X,
			Y:       sp.Y,
			R:       sp.R,
		}
		if t, ok := ticketMap[sp.Label]; ok {
			s.TicketID = t.ID
			s.PriceJPY = int(t.PriceJPY)
			s.Status = string(t.Status)
		}
		seats = append(seats, s)
	}

	return &SeatMapResponse{
		EventID: eventID,
		Layout: SeatMapLayout{
			Canvas:   layout.Canvas,
			Stage:    layout.Stage,
			Sections: layout.Sections,
		},
		Seats: seats,
	}, nil
}

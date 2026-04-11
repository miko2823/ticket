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

package event

import "context"

// UseCase is the application service for the Event bounded context.
type UseCase struct {
	service *Service
}

func NewUseCase(service *Service) *UseCase {
	return &UseCase{service: service}
}

// ListEvents returns all events.
func (uc *UseCase) ListEvents(ctx context.Context) ([]Event, error) {
	// TODO: implement
	return nil, nil
}

// GetEvent returns a single event by ID.
func (uc *UseCase) GetEvent(ctx context.Context, id string) (*Event, error) {
	// TODO: implement
	return nil, nil
}

// GetTickets returns all tickets for an event.
func (uc *UseCase) GetTickets(ctx context.Context, eventID string) ([]Ticket, error) {
	// TODO: implement
	return nil, nil
}

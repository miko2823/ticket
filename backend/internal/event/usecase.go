package event

import "context"

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

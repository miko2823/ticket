package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/event"
)

// Repository is the PostgreSQL adapter for the Event repository port.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindAll(ctx context.Context) ([]event.Event, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*event.Event, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) FindTicketsByEventID(ctx context.Context, eventID string) ([]event.Ticket, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) FindTicketByID(ctx context.Context, id string) (*event.Ticket, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) ReserveTicket(ctx context.Context, ticketID string, currentVersion int) error {
	// TODO: implement with optimistic locking
	return nil
}

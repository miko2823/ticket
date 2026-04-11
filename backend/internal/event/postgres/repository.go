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
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, venue, starts_at, created_at FROM events ORDER BY starts_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Venue, &e.StartsAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id string) (*event.Event, error) {
	var e event.Event
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, venue, starts_at, created_at FROM events WHERE id = $1", id).
		Scan(&e.ID, &e.Name, &e.Venue, &e.StartsAt, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) FindTicketsByEventID(ctx context.Context, eventID string) ([]event.Ticket, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, event_id, seat_label, price_jpy, status, version, created_at FROM tickets WHERE event_id = $1 ORDER BY seat_label",
		eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []event.Ticket
	for rows.Next() {
		var t event.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.SeatLabel, &t.PriceJPY, &t.Status, &t.Version, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

func (r *Repository) FindTicketByID(ctx context.Context, id string) (*event.Ticket, error) {
	var t event.Ticket
	err := r.pool.QueryRow(ctx,
		"SELECT id, event_id, seat_label, price_jpy, status, version, created_at FROM tickets WHERE id = $1", id).
		Scan(&t.ID, &t.EventID, &t.SeatLabel, &t.PriceJPY, &t.Status, &t.Version, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ReserveTicket(ctx context.Context, ticketID string, currentVersion int) error {
	// TODO: implement with optimistic locking
	return nil
}

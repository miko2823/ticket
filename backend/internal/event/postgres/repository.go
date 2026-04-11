package postgres

import (
	"context"
	"fmt"
	"time"

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
		`SELECT id, event_id, seat_label, price_jpy,
			CASE WHEN status = 'reserved' AND reserved_until < now() THEN 'available' ELSE status END,
			version, reserved_until, created_at
		FROM tickets WHERE event_id = $1 ORDER BY seat_label`,
		eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []event.Ticket
	for rows.Next() {
		var t event.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.SeatLabel, &t.PriceJPY, &t.Status, &t.Version, &t.ReservedUntil, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

func (r *Repository) FindTicketByID(ctx context.Context, id string) (*event.Ticket, error) {
	var t event.Ticket
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, seat_label, price_jpy,
			CASE WHEN status = 'reserved' AND reserved_until < now() THEN 'available' ELSE status END,
			version, reserved_until, created_at
		FROM tickets WHERE id = $1`, id).
		Scan(&t.ID, &t.EventID, &t.SeatLabel, &t.PriceJPY, &t.Status, &t.Version, &t.ReservedUntil, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ReserveTicket atomically reserves a ticket using optimistic locking.
// Returns an error if the ticket was already taken (version mismatch or not available).
func (r *Repository) ReserveTicket(ctx context.Context, ticketID string, currentVersion int, reservedUntil time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE tickets
		SET status = 'reserved', version = version + 1, reserved_until = $3
		WHERE id = $1
			AND version = $2
			AND (status = 'available' OR (status = 'reserved' AND reserved_until < now()))`,
		ticketID, currentVersion, reservedUntil)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ticket not available")
	}
	return nil
}

func (r *Repository) UpdateTicketStatus(ctx context.Context, ticketID string, status event.TicketStatus) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE tickets SET status = $2, reserved_until = NULL WHERE id = $1",
		ticketID, status)
	return err
}

// ReleaseExpiredReservations resets expired reserved tickets back to available.
func (r *Repository) ReleaseExpiredReservations(ctx context.Context, now time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tickets SET status = 'available', reserved_until = NULL
		WHERE status = 'reserved' AND reserved_until < $1`, now)
	return err
}

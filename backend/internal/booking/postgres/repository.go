package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/booking"
)

// Repository is the PostgreSQL adapter for the Booking repository port.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, b *booking.Booking) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO bookings (user_id, ticket_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		b.UserID, b.TicketID, b.Status).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *Repository) FindByID(ctx context.Context, id string) (*booking.Booking, error) {
	var b booking.Booking
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, ticket_id, status, created_at, updated_at FROM bookings WHERE id = $1", id).
		Scan(&b.ID, &b.UserID, &b.TicketID, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]booking.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, ticket_id, status, created_at, updated_at FROM bookings WHERE user_id = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.TicketID, &b.Status, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status booking.BookingStatus) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE bookings SET status = $2, updated_at = now() WHERE id = $1",
		id, status)
	return err
}

func (r *Repository) CreateFailedBooking(ctx context.Context, fb *booking.FailedBooking) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO failed_bookings (booking_id, reason)
		VALUES ($1, $2) RETURNING id, failed_at`,
		fb.BookingID, fb.Reason).
		Scan(&fb.ID, &fb.FailedAt)
}

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
	// TODO: implement
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*booking.Booking, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]booking.Booking, error) {
	// TODO: implement
	return nil, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status booking.BookingStatus) error {
	// TODO: implement
	return nil
}

func (r *Repository) CreateFailedBooking(ctx context.Context, fb *booking.FailedBooking) error {
	// TODO: implement
	return nil
}

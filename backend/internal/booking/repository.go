package booking

import "context"

// Repository defines the port for Booking persistence.
type Repository interface {
	Create(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id string) (*Booking, error)
	FindByUserID(ctx context.Context, userID string) ([]Booking, error)
	UpdateStatus(ctx context.Context, id string, status BookingStatus) error
	CreateFailedBooking(ctx context.Context, fb *FailedBooking) error
}

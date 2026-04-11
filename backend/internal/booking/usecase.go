package booking

import "context"

// UseCase is the application service for the Booking bounded context.
type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// CreateBooking reserves a ticket and enqueues async confirmation.
func (uc *UseCase) CreateBooking(ctx context.Context, userID, ticketID string) (*Booking, error) {
	// TODO: implement
	return nil, nil
}

// GetBookingStatus returns the current status of a booking.
func (uc *UseCase) GetBookingStatus(ctx context.Context, id string) (*Booking, error) {
	// TODO: implement
	return nil, nil
}

// GetUserBookings returns all bookings for a user.
func (uc *UseCase) GetUserBookings(ctx context.Context, userID string) ([]Booking, error) {
	// TODO: implement
	return nil, nil
}

// CancelBooking cancels a booking if allowed.
func (uc *UseCase) CancelBooking(ctx context.Context, id, userID string) error {
	// TODO: implement
	return nil
}

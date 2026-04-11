package booking

import (
	"context"
	"fmt"
)

// UseCase is the application service for the Booking bounded context.
type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// CreateBooking creates a pending booking for a reserved ticket.
// In production this would enqueue a Cloud Tasks job for async confirmation.
// For now it creates the booking as pending (simulating the async pattern).
func (uc *UseCase) CreateBooking(ctx context.Context, userID, ticketID string) (*Booking, error) {
	b := &Booking{
		UserID:   userID,
		TicketID: ticketID,
		Status:   BookingStatusPending,
	}
	if err := uc.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	// TODO: enqueue Cloud Tasks job for async confirmation
	// For now, confirm immediately (simulating payment success)
	if err := uc.repo.UpdateStatus(ctx, b.ID, BookingStatusConfirmed); err != nil {
		return nil, err
	}
	b.Status = BookingStatusConfirmed

	return b, nil
}

// GetBookingStatus returns the current status of a booking.
func (uc *UseCase) GetBookingStatus(ctx context.Context, id string) (*Booking, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetUserBookings returns all bookings for a user.
func (uc *UseCase) GetUserBookings(ctx context.Context, userID string) ([]Booking, error) {
	return uc.repo.FindByUserID(ctx, userID)
}

// CancelBooking cancels a booking if allowed.
func (uc *UseCase) CancelBooking(ctx context.Context, id, userID string) error {
	b, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("booking not found")
	}
	if b.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	if !b.CanBeCancelled() {
		return fmt.Errorf("booking cannot be cancelled")
	}
	return uc.repo.UpdateStatus(ctx, id, BookingStatusCancelled)
}

package booking

import (
	"context"
	"fmt"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/event"
)

// UseCase is the application service for the Booking bounded context.
type UseCase struct {
	repo      Repository
	eventRepo event.Repository
}

func NewUseCase(repo Repository, eventRepo event.Repository) *UseCase {
	return &UseCase{repo: repo, eventRepo: eventRepo}
}

// CreateBooking creates a booking and marks the ticket as sold.
// In production this would enqueue a Cloud Tasks job for async confirmation.
// For now it confirms immediately (simulating payment success).
func (uc *UseCase) CreateBooking(ctx context.Context, userID, ticketID string) (*Booking, error) {
	// Verify the ticket is still reserved
	ticket, err := uc.eventRepo.FindTicketByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found")
	}
	if ticket.Status != event.TicketStatusReserved {
		return nil, fmt.Errorf("ticket is not reserved")
	}
	if !ticket.IsReservedBy(userID) {
		return nil, fmt.Errorf("ticket is reserved by another user")
	}

	b := &Booking{
		UserID:   userID,
		TicketID: ticketID,
		Status:   BookingStatusPending,
	}
	if err := uc.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	// TODO: enqueue Cloud Tasks job for async confirmation
	// For now, confirm immediately and mark ticket as sold
	if err := uc.repo.UpdateStatus(ctx, b.ID, BookingStatusConfirmed); err != nil {
		return nil, err
	}
	if err := uc.eventRepo.UpdateTicketStatus(ctx, ticketID, event.TicketStatusSold); err != nil {
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

// CancelBooking cancels a booking and releases the ticket.
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
	if err := uc.repo.UpdateStatus(ctx, id, BookingStatusCancelled); err != nil {
		return err
	}
	return uc.eventRepo.UpdateTicketStatus(ctx, b.TicketID, event.TicketStatusAvailable)
}

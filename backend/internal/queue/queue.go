package queue

import "context"

// Producer enqueues async tasks (Cloud Tasks).
type Producer interface {
	EnqueueBookingConfirmation(ctx context.Context, bookingID string) error
}

// Consumer processes async tasks.
type Consumer interface {
	HandleBookingConfirmation(ctx context.Context, bookingID string) error
}

// TODO: Cloud Tasks implementation

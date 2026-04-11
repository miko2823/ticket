package booking

import (
	"net/http"
)

// Handler is the HTTP delivery layer for the Booking bounded context.
type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// TODO: implement — return 202 Accepted with booking_id
}

func (h *Handler) GetBookingStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

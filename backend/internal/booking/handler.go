package booking

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is the HTTP delivery layer for the Booking bounded context.
type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// RegisterRoutes registers booking routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/bookings", h.CreateBooking)
	r.Get("/api/v1/bookings/{id}/status", h.GetBookingStatus)
	r.Get("/api/v1/bookings/me", h.GetMyBookings)
	r.Delete("/api/v1/bookings/{id}", h.CancelBooking)
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

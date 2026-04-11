package booking

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers booking routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/bookings", h.CreateBooking)
	r.Get("/api/v1/bookings/{id}/status", h.GetBookingStatus)
	r.Get("/api/v1/bookings/me", h.GetMyBookings)
	r.Delete("/api/v1/bookings/{id}", h.CancelBooking)
}

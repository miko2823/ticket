package booking

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CaptchaMiddleware returns reCAPTCHA middleware for a given action.
type CaptchaMiddleware = func(action string) func(http.Handler) http.Handler

// RegisterRoutes registers booking routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router, captcha CaptchaMiddleware) {
	r.With(captcha("create_booking")).Post("/api/v1/bookings", h.CreateBooking)
	r.Get("/api/v1/bookings/{id}/status", h.GetBookingStatus)
	r.Get("/api/v1/bookings/me", h.GetMyBookings)
	r.Delete("/api/v1/bookings/{id}", h.CancelBooking)
}

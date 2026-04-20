package session

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CaptchaMiddleware returns reCAPTCHA middleware for a given action.
type CaptchaMiddleware = func(action string) func(http.Handler) http.Handler

// RegisterRoutes registers session routes (all require auth).
func (h *Handler) RegisterRoutes(r chi.Router, captcha CaptchaMiddleware) {
	r.With(captcha("create_session")).Post("/api/v1/events/{id}/session", h.CreateSession)
	r.Put("/api/v1/events/{id}/session/{sessionId}", h.RefreshSession)
	r.Delete("/api/v1/events/{id}/session/{sessionId}", h.EndSession)
	r.Get("/api/v1/events/{id}/queue", h.GetQueueStatus)
	r.Delete("/api/v1/events/{id}/queue", h.LeaveQueue)
}

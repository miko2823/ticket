package session

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers session routes (all require auth).
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/events/{id}/session", h.CreateSession)
	r.Put("/api/v1/events/{id}/session/{sessionId}", h.RefreshSession)
	r.Delete("/api/v1/events/{id}/session/{sessionId}", h.EndSession)
}

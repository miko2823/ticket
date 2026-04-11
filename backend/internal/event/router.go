package event

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers event routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/events", h.ListEvents)
	r.Get("/api/v1/events/{id}", h.GetEvent)
	r.Get("/api/v1/events/{id}/tickets", h.GetTickets)
}

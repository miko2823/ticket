package event

import "github.com/go-chi/chi/v5"

// RegisterPublicRoutes registers event routes that don't require auth.
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/api/v1/events", h.ListEvents)
	r.Get("/api/v1/events/{id}", h.GetEvent)
	r.Get("/api/v1/events/{id}/tickets", h.GetTickets)
	r.Get("/api/v1/events/{id}/seatmap", h.GetSeatMap)
	r.Get("/api/v1/tickets/{ticketId}", h.GetTicket)
}

// RegisterProtectedRoutes registers event routes that require auth.
func (h *Handler) RegisterProtectedRoutes(r chi.Router) {
	r.Post("/api/v1/tickets/{ticketId}/reserve", h.ReserveTicket)
	r.Delete("/api/v1/tickets/{ticketId}/reserve", h.ReleaseTicket)
}

package event

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is the HTTP delivery layer for the Event bounded context.
type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// RegisterRoutes registers event routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/events", h.ListEvents)
	r.Get("/api/v1/events/{id}", h.GetEvent)
	r.Get("/api/v1/events/{id}/tickets", h.GetTickets)
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) GetTickets(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

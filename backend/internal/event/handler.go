package event

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/auth"
	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

// Handler is the HTTP delivery layer for the Event bounded context.
type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

type eventResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Venue             string    `json:"venue"`
	StartsAt          time.Time `json:"starts_at"`
	TicketingStartsAt time.Time `json:"ticketing_starts_at"`
	TicketingEndsAt   time.Time `json:"ticketing_ends_at"`
}

type ticketResponse struct {
	ID            string  `json:"id"`
	EventID       string  `json:"event_id"`
	SeatLabel     string  `json:"seat_label"`
	PriceJPY      int     `json:"price_jpy"`
	Status        string  `json:"status"`
	ReservedUntil *string `json:"reserved_until,omitempty"`
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.useCase.ListEvents(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list events")
		return
	}

	out := make([]eventResponse, len(events))
	for i, e := range events {
		out[i] = toEventResponse(&e)
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, err := h.useCase.GetEvent(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}

	response.JSON(w, http.StatusOK, toEventResponse(e))
}

func (h *Handler) GetTickets(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	tickets, err := h.useCase.GetTickets(r.Context(), eventID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list tickets")
		return
	}

	out := make([]ticketResponse, len(tickets))
	for i, t := range tickets {
		out[i] = ticketResponse{
			ID:        t.ID,
			EventID:   t.EventID,
			SeatLabel: string(t.SeatLabel),
			PriceJPY:  int(t.PriceJPY),
			Status:    string(t.Status),
		}
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	t, err := h.useCase.GetTicket(r.Context(), ticketID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "ticket not found")
		return
	}
	response.JSON(w, http.StatusOK, toTicketResponse(t))
}

func (h *Handler) ReserveTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	userID := auth.UserIDFromContext(r.Context())
	ticket, err := h.useCase.ReserveTicket(r.Context(), ticketID, userID)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, toTicketResponse(ticket))
}

func (h *Handler) ReleaseTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	userID := auth.UserIDFromContext(r.Context())
	if err := h.useCase.ReleaseTicket(r.Context(), ticketID, userID); err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "released"})
}

func (h *Handler) GetSeatMap(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	seatMap, err := h.useCase.GetSeatMap(r.Context(), eventID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, seatMap)
}

func toEventResponse(e *Event) eventResponse {
	return eventResponse{
		ID:                e.ID,
		Name:              e.Name,
		Venue:             e.Venue,
		StartsAt:          e.StartsAt,
		TicketingStartsAt: e.TicketingStartsAt,
		TicketingEndsAt:   e.TicketingEndsAt,
	}
}

func toTicketResponse(t *Ticket) ticketResponse {
	resp := ticketResponse{
		ID:        t.ID,
		EventID:   t.EventID,
		SeatLabel: string(t.SeatLabel),
		PriceJPY:  int(t.PriceJPY),
		Status:    string(t.Status),
	}
	if t.ReservedUntil != nil {
		s := t.ReservedUntil.Format(time.RFC3339)
		resp.ReservedUntil = &s
	}
	return resp
}

package booking

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/auth"
	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

// Handler is the HTTP delivery layer for the Booking bounded context.
type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

type createBookingRequest struct {
	TicketID string `json:"ticket_id"`
}

type bookingResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TicketID  string    `json:"ticket_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TicketID == "" {
		response.Error(w, http.StatusBadRequest, "ticket_id is required")
		return
	}

	b, err := h.useCase.CreateBooking(r.Context(), userID, req.TicketID)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.JSON(w, http.StatusAccepted, bookingResponse{
		ID:        b.ID,
		UserID:    b.UserID,
		TicketID:  b.TicketID,
		Status:    string(b.Status),
		CreatedAt: b.CreatedAt,
	})
}

func (h *Handler) GetBookingStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b, err := h.useCase.GetBookingStatus(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "booking not found")
		return
	}

	response.JSON(w, http.StatusOK, bookingResponse{
		ID:       b.ID,
		UserID:   b.UserID,
		TicketID: b.TicketID,
		Status:   string(b.Status),
	})
}

func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	bookings, err := h.useCase.GetUserBookings(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list bookings")
		return
	}

	out := make([]bookingResponse, len(bookings))
	for i, b := range bookings {
		out[i] = bookingResponse{
			ID:        b.ID,
			UserID:    b.UserID,
			TicketID:  b.TicketID,
			Status:    string(b.Status),
			CreatedAt: b.CreatedAt,
		}
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := auth.UserIDFromContext(r.Context())

	if err := h.useCase.CancelBooking(r.Context(), id, userID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

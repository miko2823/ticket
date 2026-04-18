package session

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/auth"
	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

// Handler is the HTTP delivery layer for session management.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	userID := auth.UserIDFromContext(r.Context())

	sessionID, err := h.service.CreateSession(r.Context(), eventID, userID)
	if err != nil {
		if err.Error() == "seat map is full, please wait" {
			response.Error(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"session_id":            sessionID,
		"heartbeat_interval_ms": h.service.HeartbeatIntervalMs(),
	})
}

func (h *Handler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	sessionID := chi.URLParam(r, "sessionId")

	if err := h.service.RefreshSession(r.Context(), sessionID, eventID); err != nil {
		response.Error(w, http.StatusNotFound, "session not found or expired")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
}

func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	sessionID := chi.URLParam(r, "sessionId")

	if err := h.service.EndSession(r.Context(), sessionID, eventID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to end session")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

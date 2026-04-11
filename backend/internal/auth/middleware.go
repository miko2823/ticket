package auth

import (
	"context"
	"net/http"
	"strings"

	firebaseAuth "firebase.google.com/go/v4/auth"

	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

type contextKey string

const userIDKey contextKey = "userID"

// Middleware verifies Firebase ID tokens and injects userID into context.
type Middleware struct {
	client *firebaseAuth.Client
}

func NewMiddleware(client *firebaseAuth.Client) *Middleware {
	return &Middleware{client: client}
}

// Authenticate returns an HTTP middleware that validates Firebase tokens.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		token, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			response.Error(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		verified, err := m.client.VerifyIDToken(r.Context(), token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, verified.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts the userID set by the Authenticate middleware.
func UserIDFromContext(ctx context.Context) string {
	uid, _ := ctx.Value(userIDKey).(string)
	return uid
}

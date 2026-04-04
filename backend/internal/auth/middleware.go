package auth

import "net/http"

// Middleware verifies Firebase ID tokens and injects userID into context.
type Middleware struct {
	// TODO: Firebase Auth client
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

// Authenticate returns an HTTP middleware that validates Firebase tokens.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract token from Authorization header
		// TODO: verify with Firebase Admin SDK
		// TODO: inject userID into context
		next.ServeHTTP(w, r)
	})
}

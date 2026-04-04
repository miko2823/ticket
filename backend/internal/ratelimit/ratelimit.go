package ratelimit

import "net/http"

// Middleware provides token-bucket rate limiting per user (Firebase UID) and per IP.
// NOTE: This uses an in-memory store suitable for single-instance demo.
// In production with multiple Cloud Run instances, replace with Redis.
type Middleware struct {
	// TODO: token bucket state
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

// LimitByUser returns middleware that rate-limits by Firebase UID (10 req/min).
func (m *Middleware) LimitByUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement token bucket per user
		next.ServeHTTP(w, r)
	})
}

// LimitByIP returns middleware that rate-limits by IP (100 req/min).
func (m *Middleware) LimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement token bucket per IP
		next.ServeHTTP(w, r)
	})
}

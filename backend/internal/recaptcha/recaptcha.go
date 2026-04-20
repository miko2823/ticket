package recaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Verifier verifies reCAPTCHA v3 tokens with Google's API.
type Verifier struct {
	secretKey    string
	minScore     float64
	httpClient   *http.Client
}

func NewVerifier(secretKey string, minScore float64) *Verifier {
	return &Verifier{
		secretKey: secretKey,
		minScore:  minScore,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type verifyResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ErrorCodes  []string `json:"error-codes"`
}

// Verify checks a reCAPTCHA token against Google's API and validates the action.
func (v *Verifier) Verify(ctx context.Context, token, expectedAction string) error {
	if token == "" {
		return fmt.Errorf("missing recaptcha token")
	}

	resp, err := v.httpClient.PostForm(verifyURL, url.Values{
		"secret":   {v.secretKey},
		"response": {token},
	})
	if err != nil {
		return fmt.Errorf("recaptcha verification request failed: %w", err)
	}
	defer resp.Body.Close()

	var result verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode recaptcha response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("recaptcha verification failed: %v", result.ErrorCodes)
	}

	if result.Action != expectedAction {
		return fmt.Errorf("recaptcha action mismatch: got %q, want %q", result.Action, expectedAction)
	}

	if result.Score < v.minScore {
		return fmt.Errorf("recaptcha score too low: %.1f", result.Score)
	}

	return nil
}

// RequireToken returns a middleware that verifies the X-Recaptcha-Token header.
// The expected reCAPTCHA action is derived from the action parameter.
func (v *Verifier) RequireToken(action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Recaptcha-Token")
			if err := v.Verify(r.Context(), token, action); err != nil {
				response.Error(w, http.StatusForbidden, "bot detection: "+err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

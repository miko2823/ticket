package circuitbreaker

// Breaker wraps external calls with circuit breaker logic using sony/gobreaker.
// States: Closed → Open (after 5 failures in 10s) → Half-Open (after 30s)
type Breaker struct {
	// TODO: gobreaker.CircuitBreaker
}

func NewBreaker(name string) *Breaker {
	// TODO: configure gobreaker settings
	return &Breaker{}
}

// Execute runs fn through the circuit breaker.
func (b *Breaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	// TODO: implement
	return fn()
}

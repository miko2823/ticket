package booking

// Service is the domain service for the Booking bounded context.
// Domain service methods will be added here when cross-aggregate
// logic is needed (e.g. double-booking prevention checks).
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

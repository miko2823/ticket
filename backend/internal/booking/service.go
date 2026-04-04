package booking

// Service is the domain service for the Booking bounded context.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// TODO: domain service methods (double-booking prevention logic, etc.)

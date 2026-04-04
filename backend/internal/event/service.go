package event

// Service is the domain service for the Event bounded context.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// TODO: domain service methods

package event

// Service is the domain service for the Event bounded context.
// Currently empty — domain service methods will be added here when
// cross-aggregate or complex domain logic is needed (e.g. validation
// that spans multiple entities).
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

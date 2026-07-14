package dashboard

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAdminDashboard(ctx context.Context) (*AdminDashboard, error) {
	return s.repo.GetAdminDashboard(ctx)
}

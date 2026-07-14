package dashboard

import "context"

type Repository interface {
	GetAdminDashboard(ctx context.Context) (*AdminDashboard, error)
}

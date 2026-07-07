package role

type Service struct {
	roleRepo Repository
}

func NewService(roleRepo Repository) *Service {
	return &Service{roleRepo: roleRepo}
}
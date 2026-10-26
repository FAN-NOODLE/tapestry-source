package activity

// NewActivity Service初始化
func NewActivity(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

package enterprise

// New 实例化存储
func New(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

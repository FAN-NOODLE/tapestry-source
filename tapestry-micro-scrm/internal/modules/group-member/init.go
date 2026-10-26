package groupmember

// New 初始化存储
func New(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

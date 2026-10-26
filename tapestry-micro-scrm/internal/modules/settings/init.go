package settings

//New 实例化
func New(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

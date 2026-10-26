package shop_member

func New(repository Repository) *Service {
	return &Service{
		Repo: repository,
	}
}

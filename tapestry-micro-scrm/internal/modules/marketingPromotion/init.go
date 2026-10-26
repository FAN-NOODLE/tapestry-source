package marketingPromotion

func New(repository Repository) *Service {
	return &Service{
		Repo: repository,
	}
}

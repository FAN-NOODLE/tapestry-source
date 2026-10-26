package sms

func New(repository Repository) *Service {
	return &Service{
		Repo: repository,
	}
}

func InsideNew(repository Repository, marketingInterface MarketingInterface) *Service {
	return &Service{
		Repo:      repository,
		Marketing: marketingInterface,
	}
}

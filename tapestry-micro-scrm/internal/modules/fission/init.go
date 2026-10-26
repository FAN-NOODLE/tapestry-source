package fission

func NewService(repo Repository, contact ContactInterface, user UserInterface, bossApi BossApiInterface, enterprise EnterpriseInterface) *Service {
	return &Service{
		repo:       repo,
		Contact:    contact,
		User:       user,
		BossApi:    bossApi,
		Enterprise: enterprise,
	}
}

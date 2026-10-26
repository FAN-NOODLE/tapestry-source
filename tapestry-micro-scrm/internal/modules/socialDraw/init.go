package socialDraw

func NewGroupChat(repo Repository, ck Ckrepository) *Service {
	return &Service{
		repo:   repo,
		ckrepo: ck,
	}
}

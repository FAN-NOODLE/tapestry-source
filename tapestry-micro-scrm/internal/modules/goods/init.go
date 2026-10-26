package goods

// NewGoods Service初始化
func NewGoods(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

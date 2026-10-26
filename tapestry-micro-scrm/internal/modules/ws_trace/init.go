package ws_trace

// Service 用户服务
type Service struct {
	repo Repository
}

// NewService 初始化服务
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

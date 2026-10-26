package workexternal

// New 初始化存储
func New(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

// NewExternalTrace 初始化外部联系人统计服务
func NewExternalTrace(repository Repository) *ExternalTraceService {
	return &ExternalTraceService{
		repo: repository,
	}
}

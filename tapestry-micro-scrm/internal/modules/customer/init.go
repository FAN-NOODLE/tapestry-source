package customer

import (
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	shopMemberMysql "micro-scrm/internal/modules/shop_member/respository/mysql"
	tagMySQL "micro-scrm/internal/modules/tag/respository/mysql"
)

// New 实例化存储
func New(repository Repository, workRepo WorkUserRepository, esRepo EsRepository) *Service {
	return &Service{
		repo:            repository,
		workRepo:        workRepo,
		esRepo:          esRepo,
		enterpriseRepo:  enterpriseMysql.New(),
		tagRepo:         tagMySQL.NewTagRepo(),
		memberLevelRepo: shopMemberMysql.New(),
	}
}

func NewDelMessage(repository Repository) *MessageService {
	return &MessageService{
		repo: repository,
	}
}

// NewExternalTag 获取标签服务对象
func NewExternalTag(repository Repository) *ExternalTagService {
	return &ExternalTagService{
		repo: repository,
	}
}

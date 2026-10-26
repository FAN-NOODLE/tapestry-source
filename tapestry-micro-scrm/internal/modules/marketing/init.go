package marketing

import (
	customerMySQL "micro-scrm/internal/modules/customer/respository/mysql"
)

// New 实例化存储
func New(repository Repository, esRepo EsInterface, arg ...Option) *Service {
	service := &Service{
		repo:         repository,
		customerRepo: customerMySQL.NewCustomer(),
		esRepo:       esRepo,
	}

	for _, opt := range arg {
		opt(service)
	}

	return service
}

// Option 选项操作
type Option func(s *Service)

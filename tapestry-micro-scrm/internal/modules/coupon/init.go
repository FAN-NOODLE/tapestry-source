package coupon

import (
	customerMySQL "micro-scrm/internal/modules/customer/respository/mysql"
	enterpriseMySQL "micro-scrm/internal/modules/enterprise/respository/mysql"
)

// New 实例化存储
func New(repository Repository) *Service {
	return &Service{
		repo:           repository,
		customerRepo:   customerMySQL.NewCustomer(),
		enterpriseRepo: enterpriseMySQL.New(),
	}
}

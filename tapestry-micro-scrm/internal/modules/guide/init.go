package guide

import (
	customerMySQL "micro-scrm/internal/modules/customer/respository/mysql"
)

// New 实例化存储
func New(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

// NewInside New 实例化存储
func NewInside(repository Repository, memberInterface MemberInterface) *Service {
	return &Service{
		repo:         repository,
		Member:       memberInterface,
		customerRepo: customerMySQL.NewCustomer(),
	}
}

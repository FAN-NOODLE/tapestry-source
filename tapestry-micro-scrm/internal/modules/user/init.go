package user

import (
	"micro-scrm/internal/modules/customer"
	customerMySQL "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
)

// NewUserService 初始化服务
func NewUserService(repo Repository, eventRepo EventRepository, listRepo customer.EsRepository, enterpriseRepo enterprise.Repository) *Service {
	return &Service{
		repo:           repo,
		enterpriseRepo: enterpriseRepo,
		eventRepo:      eventRepo,
		workRepo:       customerMySQL.NewWorkUser(),
		esRepo:         listRepo,
	}
}

// NewOrderService 初始化订单服务
func NewOrderService() *OrderService {
	return &OrderService{
		customerRepo: customerMySQL.NewCustomer(),
	}
}

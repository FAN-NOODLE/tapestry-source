package sidebarGroupChat

import (
	customerMySQL "micro-scrm/internal/modules/customer/respository/mysql"
	enterpriseMySQL "micro-scrm/internal/modules/enterprise/respository/mysql"
)

func New(repository Repository) *SidebarGroupChatService {
	return &SidebarGroupChatService{
		Repo:           repository,
		enterpriseRepo: enterpriseMySQL.New(),
		customerRepo:   customerMySQL.NewCustomer(),
	}
}

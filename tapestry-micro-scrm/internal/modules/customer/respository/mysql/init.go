package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func NewCustomer() *CustomerRepo {
	return &CustomerRepo{repos.NewSQLRepo("t_customer", db.DB())}
}

func NewDelMessage() customer.Repository {
	return &CustomerRepo{repos.NewSQLRepo("t_customers_del_message", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &customer.DelMessage{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &customer.MessageListResponse{}
	}))}
}

func NewWorkUser() customer.WorkUserRepository {
	return &WorkRepo{repos.NewSQLRepo("t_work_external_user", db.DB())}
}

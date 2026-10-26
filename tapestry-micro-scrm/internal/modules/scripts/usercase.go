package scripts

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/customer"
)

type Service struct {
	repo            Repository
	CustomerService *customer.Service
}

// UserCategoryTable 表格执行
func (s Service) SpiltTable(cc ctx.Context, enterpriseID int64, domain string) error {
	return s.repo.SpiltTable(cc, enterpriseID)
}

//分库service方法
func (s Service) SpiltDatabase(cc ctx.Context, enterpriseID int64, domain string) error {
	return s.repo.SpiltDatabase(cc, enterpriseID)
}

func (s Service) SyncRemarkData(cc ctx.Context) error {
	// 查询有备注数据的数据
	list, _ := s.repo.GetHasRemarkData(cc)
	remarkMap := make(map[int64]string)
	for _, v := range list {
		if _, ok := remarkMap[v.WorkExternalUserID]; ok {
			remarkMap[v.WorkExternalUserID] = remarkMap[v.WorkExternalUserID] + "," + v.Description
		} else {
			remarkMap[v.WorkExternalUserID] = v.Description
		}
	}
	for _, v := range list {
		if _, ok := remarkMap[v.WorkExternalUserID]; !ok {
			continue
		}
		enterpriseUser, err := s.CustomerService.QueryEnterpriseUser(cc, v.EnterpriseID, &customer.SourceUserBase{
			ID:           v.WorkExternalUserID,
			CustomerType: customer.CustomerTypeExternal,
		})
		if err != nil || enterpriseUser == nil {
			continue
		}

		// 对euid对应的enterprise_user进行更新
		_ = s.repo.UpdateRemarkData(cc, v.EnterpriseID, enterpriseUser.ID, remarkMap[v.WorkExternalUserID])
	}
	return nil
}

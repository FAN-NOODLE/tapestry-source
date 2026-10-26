package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/departments"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type DepartmentApi struct {
}

// DepartmentDetail 组织详情
func (d DepartmentApi) DepartmentDetail(c ctx.Context, req *departments.DepartmentDetailRequest) interface{} {
	res, err := microscrm.GetDepartmentDetail(c, req)
	if err != nil {
		return nil
	}
	return res
}

// DepartmentList 组织列表
func (d DepartmentApi) DepartmentList(c ctx.Context, req *departments.DepartmentListRequest) interface{} {
	res, err := microscrm.GetDepartmentList(c, req)
	if err != nil {
		return nil
	}
	return res
}

// UpdateRecord 组织关系更新记录
func (d DepartmentApi) UpdateRecord(c ctx.Context, req *departments.UpdateRecordRequest) interface{} {
	res, err := microscrm.GetUpdateRecord(c, req)
	if err != nil {
		return nil
	}
	return res
}

// DepartmentMomentList 选择部门
func (d DepartmentApi) DepartmentMomentList(c ctx.Context, req *departments.MomentListRequest) interface{} {
	res, err := microscrm.DepartmentMomentList(c, req)
	if err != nil {
		return nil
	}
	return res
}

// UserLostConfig 导购丢失客户推送消息配置
func (d DepartmentApi) UserLostConfig(c ctx.Context, req *departments.LostConfigRequest) interface{} {
	res, err := microscrm.UserLostConfig(c, req)
	if err != nil {
		return nil
	}
	return res
}

// UserLostInfo 导购丢失客户推送消息配置
func (d DepartmentApi) UserLostInfo(c ctx.Context, req *departments.LostInfoRequest) interface{} {
	res, err := microscrm.UserLostInfo(c, req)
	if err != nil {
		return nil
	}
	return res
}

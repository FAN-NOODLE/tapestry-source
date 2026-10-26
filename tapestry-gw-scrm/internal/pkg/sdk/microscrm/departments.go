package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// GetDepartmentList 同步通讯录
func GetDepartmentList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.DepartmentListUrl), request)
}

//GetDepartmentDetail 获取组织详情
func GetDepartmentDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.DepartmentDetailUrl), request)
}

//GetUpdateRecord 获取更新时间
func GetUpdateRecord(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.UpdateRecordUrl), request)
}

// DepartmentMomentList 选择部门
func DepartmentMomentList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.DepartmentMomentListUrl), request)
}

// UserLostConfig 导购丢失客户推送消息配置
func UserLostConfig(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.UserLostConfig), request)
}

// UserLostInfo 导购丢失客户推送消息配置
func UserLostInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.UserLostInfo), request)
}

// GetSaStoreLists 获取coach sa的门店列表，奥莱还是零售
func GetSaStoreLists(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/department/sa-store-lists"), req)
}

package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// SyncMemberContact 同步通讯录
func SyncMemberContact(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.SyncContactUrl), request)
}

// GetMemberList 获取用户列表
func GetMemberList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MemberListUrl), request)
}

// GetMemberDetail 获取用户详情
func GetMemberDetail(c ctx.Context, request interface{}) (interface{}, error) {

	return sdk.SendPost(c, GetUrl(urlEnum.MemberDetailUrl), request)
}

// MomentList 选择员工
func MomentList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MomentListUrl), request)
}
func MomentBaList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MomentBaListUrl), request)
}

func GetMembersByDepartmentIds(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.SimpleMembersOfDepartmentsUrl), request)
}

// DownloadWhiteLists 下载中心的白名单
func DownloadWhiteLists(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/member/download-white-list"), request)
}

func MemberListByThirdIds(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/member/list-by-third-ids"), request)
}

// GetMemberListByStoreCodes resolves store codes to the active employees in those stores.
func GetMemberListByStoreCodes(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/member/get-member-list-by-store-codes"), request)
}

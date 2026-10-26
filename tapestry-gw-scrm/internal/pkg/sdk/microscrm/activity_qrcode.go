package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type ActivityQRCodeApi struct {
}

func NewActivityQRCodeApi() ActivityQRCodeApi {
	return ActivityQRCodeApi{}
}

// ActivitySave  活动二维码保存
func (a ActivityQRCodeApi) ActivitySave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/activity-save"), request)
}

// ActivityDetail  活动二维码详情
func (a ActivityQRCodeApi) ActivityDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/activity-detail"), request)
}

// ActivityDelete   活动二维码删除
func (a ActivityQRCodeApi) ActivityDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/activity-delete"), request)
}

// ActivityList   活动二维码列表
func (a ActivityQRCodeApi) ActivityList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/activity-list"), request)
}

// ActivityCategoryList   投放渠道列表
func (a ActivityQRCodeApi) ActivityCategoryList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/activity-category-list"), request)
}

// WelcomeSave  投放渠道保存
func (a ActivityQRCodeApi) WelcomeSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/welcome-save"), request)
}

// WelcomeDetail   投放渠道欢迎语详情
func (a ActivityQRCodeApi) WelcomeDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/welcome-detail"), request)
}

// WelcomeDelete   投放渠道欢迎语删除
func (a ActivityQRCodeApi) WelcomeDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/welcome-delete"), request)
}

// WelcomeList   投放渠道欢迎语列表
func (a ActivityQRCodeApi) WelcomeList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/welcome-list"), request)
}

// Statistics    管理后台统计
func (a ActivityQRCodeApi) Statistics(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/activity-qrcode/statistics"), request)
}

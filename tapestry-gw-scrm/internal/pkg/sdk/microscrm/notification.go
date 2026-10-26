package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type NotificationApi struct {
}

func NewNotificationApi() NotificationApi {
	return NotificationApi{}
}

// Save 新建|编辑传讯
func (a NotificationApi) Save(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/save"), request)
}

// List 传讯列表
func (a NotificationApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/list"), request)
}

func (a NotificationApi) CrowdSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/crowd/save"), request)
}
func (a NotificationApi) CrowdDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/crowd/detail"), request)
}

func (a NotificationApi) CrowdDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/crowd/delete"), request)
}
func (a NotificationApi) CrowdList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/crowd/list"), request)
}

// Info 传讯内容详情
func (a NotificationApi) Info(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/info"), request)
}

func (a NotificationApi) Departments(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/departments"), request)
}

func (a NotificationApi) InfoMergeSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/merge/save"), request)
}
func (a NotificationApi) InfoMergeRelease(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/merge/release"), request)
}
func (a NotificationApi) InfoMergeDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/merge/detail"), request)
}
func (a NotificationApi) SavePdf(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/merge/save-pdf"), request)
}

// Copy 传讯内容的复制
func (a NotificationApi) Copy(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/copy"), request)
}
func (a NotificationApi) InfoMergeList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/merge/list"), request)
}

func (a NotificationApi) Recall(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/recall"), request)
}

func (a NotificationApi) Check(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/check"), request)
}

func (a NotificationApi) Manager(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/manager"), request)
}

func (a NotificationApi) NoticeConfigSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/notice-config/save"), request)
}

func (a NotificationApi) NoticeConfigDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/notice-config/detail"), request)
}

func (a NotificationApi) Cancel(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/notification/cancel"), request)
}

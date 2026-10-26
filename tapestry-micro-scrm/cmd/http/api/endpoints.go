package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	activityHttp "micro-scrm/internal/modules/activity/delivery/http"
	activityCampaign "micro-scrm/internal/modules/activityCampaign/delivery/http"
	baKanBan "micro-scrm/internal/modules/ba_kanban/delivery/http"
	callback_http "micro-scrm/internal/modules/callback/delivery/http"
	chat_http "micro-scrm/internal/modules/chat/delivery/http"
	couponHttp "micro-scrm/internal/modules/coupon/delivery/http"
	customerHttp "micro-scrm/internal/modules/customer/delivery/http"
	data_compass "micro-scrm/internal/modules/dataComPass/delivery/http"
	departmentshttp "micro-scrm/internal/modules/department/delivery/http"
	drainagehttp "micro-scrm/internal/modules/drainage/delivery/http"
	enterpriseHttp "micro-scrm/internal/modules/enterprise/delivery/http"
	fissionHttp "micro-scrm/internal/modules/fission/delivery/http"
	goodsHttp "micro-scrm/internal/modules/goods/delivery/http"
	groupmember_http "micro-scrm/internal/modules/group-member/delivery/http"
	groupMsgHttp "micro-scrm/internal/modules/groupmsg/delivery/http"
	marketinghttp "micro-scrm/internal/modules/marketing/delivery/http"
	marketingPromotionHttp "micro-scrm/internal/modules/marketingPromotion/delivery/http"
	membershttp "micro-scrm/internal/modules/members/delivery/http"
	miniProgramHttp "micro-scrm/internal/modules/miniProgram/delivery/http"
	momenthttp "micro-scrm/internal/modules/moment/delivery/http"
	overtimehttp "micro-scrm/internal/modules/overtime/delivery/http"
	plan_task_http "micro-scrm/internal/modules/plan_task/delivery/http"
	remindconfigthttp "micro-scrm/internal/modules/remindconfig/delivery/http"
	sensitivehttp "micro-scrm/internal/modules/sensitive/delivery/http"
	settingshttp "micro-scrm/internal/modules/settings/delivery/http"
	shop_member_http "micro-scrm/internal/modules/shop_member/delivery/http"
	sidebarGroupChathttp "micro-scrm/internal/modules/sidebarGroupChat/delivery/http"
	smshttp "micro-scrm/internal/modules/sms/delivery/http"
	socialDraw_http "micro-scrm/internal/modules/socialDraw/delivery/http"
	sourceHttp "micro-scrm/internal/modules/source/delivery/http"
	sourceGroupHttp "micro-scrm/internal/modules/sourceGroup/delivery/http"
	"micro-scrm/internal/modules/tag/delivery/http"
	transferCustomerHttp "micro-scrm/internal/modules/transfer_customer/delivery/http"
	userhttp "micro-scrm/internal/modules/user/delivery/http"
	version_http "micro-scrm/internal/modules/version/delivery/http"
	work_external_http "micro-scrm/internal/modules/work-external/delivery/http"
	work_trace "micro-scrm/internal/modules/ws_trace/delivery/http"
)

func endpoints() []*swagger.Endpoint {
	var eps []*swagger.Endpoint

	//EDIT THIS, ADD YOUR ENDPOINT
	eps = append(eps, http.Endpoints()...)
	eps = append(eps, callback_http.Endpoints()...)
	eps = append(eps, departmentshttp.Endpoints()...)
	eps = append(eps, membershttp.Endpoints()...)
	eps = append(eps, momenthttp.Endpoints()...)
	eps = append(eps, customerHttp.Endpoints()...)
	eps = append(eps, smshttp.Endpoints()...)
	eps = append(eps, userhttp.Endpoints()...)
	eps = append(eps, marketinghttp.Endpoints()...)
	eps = append(eps, miniProgramHttp.Endpoints()...)
	eps = append(eps, sourceHttp.Endpoints()...)
	eps = append(eps, chat_http.Endpoints()...)
	eps = append(eps, sensitivehttp.Endpoints()...)
	eps = append(eps, remindconfigthttp.Endpoints()...)
	eps = append(eps, settingshttp.Endpoints()...)
	eps = append(eps, overtimehttp.Endpoints()...)
	eps = append(eps, work_external_http.Endpoints()...)
	eps = append(eps, groupmember_http.Endpoints()...)
	eps = append(eps, sourceGroupHttp.Endpoints()...)
	eps = append(eps, drainagehttp.Endpoints()...)
	eps = append(eps, sidebarGroupChathttp.Endpoints()...)
	eps = append(eps, enterpriseHttp.Endpoints()...)

	eps = append(eps, couponHttp.Endpoints()...)
	eps = append(eps, marketingPromotionHttp.Endpoints()...)
	eps = append(eps, activityHttp.Endpoints()...)
	eps = append(eps, goodsHttp.Endpoints()...)
	eps = append(eps, shop_member_http.Endpoints()...)
	eps = append(eps, fissionHttp.Endpoints()...)
	eps = append(eps, groupMsgHttp.Endpoints()...)
	eps = append(eps, version_http.Endpoints()...)
	eps = append(eps, plan_task_http.Endpoints()...)
	eps = append(eps, baKanBan.Endpoints()...)

	eps = append(eps, transferCustomerHttp.Endpoints()...)
	eps = append(eps, activityCampaign.Endpoints()...)
	eps = append(eps, data_compass.Endpoints()...)
	eps = append(eps, socialDraw_http.Endpoints()...)
	eps = append(eps, work_trace.Endpoints()...)
	return eps
}

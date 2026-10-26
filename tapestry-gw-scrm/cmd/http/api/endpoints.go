package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	activityCampaign "gw-scrm/internal/modules/activity_campaign/delivery/http"
	auth_template_http "gw-scrm/internal/modules/auth_template/delivery/http"
	baKanBan "gw-scrm/internal/modules/ba_kanban/delivery/http"
	chatHttp "gw-scrm/internal/modules/chat/delivery/http"
	config_guide_http "gw-scrm/internal/modules/config_guide/delivery/http"
	customer_http "gw-scrm/internal/modules/customer/delivery/http"
	dataCompasshttp "gw-scrm/internal/modules/data_compass/delivery/http"
	departments_http "gw-scrm/internal/modules/departments/delivery/http"
	downloadcenter_http "gw-scrm/internal/modules/download/delivery/http"
	drainage_http "gw-scrm/internal/modules/drainage/delivery/http"
	file_http "gw-scrm/internal/modules/file/delivery/http"
	fission_http "gw-scrm/internal/modules/fission/delivery/http"
	groupmsg_http "gw-scrm/internal/modules/groupmsg/delivery/http"
	index_http "gw-scrm/internal/modules/index/delivery/http"
	legal_right_http "gw-scrm/internal/modules/legal_right/delivery/http"
	level_direct_http "gw-scrm/internal/modules/level_direct/delivery/http"
	marketinghttp "gw-scrm/internal/modules/marketing/delivery/http"
	member_level_http "gw-scrm/internal/modules/member_level/delivery/http"
	members_http "gw-scrm/internal/modules/members/delivery/http"
	miniProgramHttp "gw-scrm/internal/modules/miniProgram/delivery/http"
	moment_http "gw-scrm/internal/modules/moment/delivery/http"
	open_http "gw-scrm/internal/modules/open/delivery/http"
	overtimehttp "gw-scrm/internal/modules/overtime/delivery/http"
	plan_task_http "gw-scrm/internal/modules/plan_task/delivery/http"
	remindconfighttp "gw-scrm/internal/modules/remindconfig/delivery/http"
	right_rule_http "gw-scrm/internal/modules/right_rule/delivery/http"
	sensitivehttp "gw-scrm/internal/modules/sensitive/delivery/http"
	settings_http "gw-scrm/internal/modules/settings/delivery/http"
	sms_http "gw-scrm/internal/modules/sms/delivery/http"
	socialDrawHttp "gw-scrm/internal/modules/socialDraw/delivery/http"
	sourceHttp "gw-scrm/internal/modules/source/delivery/http"
	sourceGroupHttp "gw-scrm/internal/modules/sourceGroup/delivery/http"
	tag_http "gw-scrm/internal/modules/tag/delivery/http"
	tag_csv_upload "gw-scrm/internal/modules/tag_csv_upload/delivery/http"
	transferCustomerHttp "gw-scrm/internal/modules/transfer_customer/delivery/http"
	user_http "gw-scrm/internal/modules/user/delivery/http"
	versionHttp "gw-scrm/internal/modules/version/delivery/http"
	wxchannelhttp "gw-scrm/internal/modules/wxChannel/delivery/http"
	wxcouponhttp "gw-scrm/internal/modules/wxCoupon/delivery/http"
)

func endpoints() []*swagger.Endpoint {
	var eps []*swagger.Endpoint

	// EDIT THIS, ADD YOUR ENDPOINT
	eps = append(eps, user_http.Endpoints()...)
	eps = append(eps, members_http.Endpoints()...)
	eps = append(eps, departments_http.Endpoints()...)
	eps = append(eps, tag_http.Endpoints()...)
	eps = append(eps, settings_http.Endpoints()...)
	eps = append(eps, marketinghttp.Endpoints()...)
	eps = append(eps, customer_http.Endpoints()...)
	eps = append(eps, sms_http.Endpoints()...)
	eps = append(eps, file_http.Endpoints()...)
	eps = append(eps, moment_http.Endpoints()...)
	eps = append(eps, index_http.Endpoints()...)
	eps = append(eps, sensitivehttp.Endpoints()...)
	eps = append(eps, remindconfighttp.Endpoints()...)
	eps = append(eps, chatHttp.Endpoints()...)
	eps = append(eps, sourceHttp.Endpoints()...)
	eps = append(eps, sourceGroupHttp.Endpoints()...)
	eps = append(eps, miniProgramHttp.Endpoints()...)
	eps = append(eps, drainage_http.Endpoints()...)
	eps = append(eps, overtimehttp.Endpoints()...)
	eps = append(eps, wxchannelhttp.Endpoints()...)
	eps = append(eps, wxcouponhttp.Endpoints()...)
	eps = append(eps, config_guide_http.Endpoints()...)
	eps = append(eps, level_direct_http.Endpoints()...)
	eps = append(eps, member_level_http.Endpoints()...)
	eps = append(eps, right_rule_http.Endpoints()...)
	eps = append(eps, legal_right_http.Endpoints()...)
	eps = append(eps, auth_template_http.Endpoints()...)
	eps = append(eps, fission_http.Endpoints()...)
	eps = append(eps, groupmsg_http.Endpoints()...)
	eps = append(eps, versionHttp.Endpoints()...)
	eps = append(eps, plan_task_http.Endpoints()...)
	eps = append(eps, baKanBan.Endpoints()...)
	eps = append(eps, transferCustomerHttp.Endpoints()...)
	eps = append(eps, activityCampaign.Endpoints()...)
	eps = append(eps, dataCompasshttp.Endpoints()...)
	eps = append(eps, tag_csv_upload.Endpoints()...)
	eps = append(eps, socialDrawHttp.Endpoints()...)
	eps = append(eps, open_http.Endpoints()...)
	eps = append(eps, downloadcenter_http.Endpoints()...)
	return eps
}

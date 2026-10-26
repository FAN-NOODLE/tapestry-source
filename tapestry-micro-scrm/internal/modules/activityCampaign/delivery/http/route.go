package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/activityCampaign"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/activities/h5/info", "H5活动任务详情",
			endpoint.Handler(activitiesH5Info),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.UsersignInRequest{}, "请求参数", true),
		),
		endpoint.Post("/activities/h5/user/signin/list", "H5活动用户签到列表",
			endpoint.Handler(usersignInH5List),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.UsersignInRequest{}, "请求参数", true),
		),
		endpoint.Post("/activities/h5/user/date/signin", "H5活动用户每日签到",
			endpoint.Handler(usersignInH5),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.UsersignInRequest{}, "请求参数", true),
		),
		endpoint.Post("/activities/h5/user/date/signin/info", "H5活动用户每日签到详情",
			endpoint.Handler(usersignInInfoH5),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.UsersignInRequest{}, "请求参数", true),
		),
		endpoint.Post("/activities/sign/list", "用户签到列表",
			endpoint.Handler(signInList),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.SignReceiveInfo{}, "请求参数", true),
		),
		endpoint.Post("/activities/sign/detail", "用户签到详情",
			endpoint.Handler(signInDetail),
			endpoint.Tags("活动签到"),
			endpoint.Body(activityCampaign.SignDetailInfo{}, "请求参数", true),
		),
	}
}

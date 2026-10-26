package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	activityHttp "gw-work-scrm/internal/modules/activity/delivery/http"
	auth "gw-work-scrm/internal/modules/auth/delivery/http"
	couponhttp "gw-work-scrm/internal/modules/coupon/delivery/http"
	customerhttp "gw-work-scrm/internal/modules/customer/delivery/http"
	fissionActivityhttp "gw-work-scrm/internal/modules/fissionActivity/delivery/http"
	goodsHttp "gw-work-scrm/internal/modules/goods/delivery/http"
	groupChathttp "gw-work-scrm/internal/modules/groupChat/delivery/http"
	posterhttp "gw-work-scrm/internal/modules/poster/delivery/http"
	sidebarGroupChathttp "gw-work-scrm/internal/modules/sidebarGroupChat/delivery/http"
	sourcehttp "gw-work-scrm/internal/modules/source/delivery/http"
	taghttp "gw-work-scrm/internal/modules/tag/delivery/http"
)

func endpoints() []*swagger.Endpoint {
	var eps []*swagger.Endpoint

	// EDIT THIS, ADD YOUR ENDPOINT
	eps = append(eps, auth.Endpoints()...)
	eps = append(eps, customerhttp.Endpoints()...)
	eps = append(eps, taghttp.Endpoints()...)
	eps = append(eps, sidebarGroupChathttp.Endpoints()...)
	eps = append(eps, sourcehttp.Endpoints()...)
	eps = append(eps, couponhttp.Endpoints()...)
	eps = append(eps, goodsHttp.Endpoints()...)
	eps = append(eps, activityHttp.Endpoints()...)
	eps = append(eps, posterhttp.Endpoints()...)
	eps = append(eps, fissionActivityhttp.Endpoints()...)
	eps = append(eps, groupChathttp.Endpoints()...)
	return eps
}

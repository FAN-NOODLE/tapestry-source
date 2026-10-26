package wxchannelhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/wx-channel/wx-channel-list", "微信场景值列表",
			endpoint.Handler(wxChannelList),
			endpoint.Tags("微信场景值"),
			endpoint.Body("", "请求参数", true),
			endpoint.Response(http.StatusOK, "", "wxChannel", "接口调用成功")),
	}
}

package settings_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/settings"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/settings", "企微配置",
			endpoint.Handler(getSettings),
			endpoint.Tags("设置"),
			endpoint.Body(settings.AuthRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: settings.SettingResponse{},
			}, "GetSettings", "接口调用成功")),

		endpoint.Post("/auth-contact-update", "外部联系人secret更新",
			endpoint.Handler(authContactUpdate),
			endpoint.Tags("设置"),
			endpoint.Body(settings.ContactRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: "",
			}, "GetSettings", "接口调用成功")),
		endpoint.Post("/set-wxqy-token", "客户联系token更新",
			endpoint.Handler(setWxqyToken),
			endpoint.Tags("设置"),
			endpoint.Body(settings.ContactRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: "",
			}, "GetSettings", "接口调用成功")),
		endpoint.Post("/work-plugin-update", "自建应用更新",
			endpoint.Handler(workPluginUpdate),
			endpoint.Tags("设置"),
			endpoint.Body(settings.PluginRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: "",
			}, "GetSettings", "接口调用成功")),
		endpoint.Post("/upload-check-file", "上传校验文件",
			endpoint.Handler(uploadCheckFile),
			endpoint.Tags("设置"),
			endpoint.Body(settings.PluginRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: "",
			}, "uploadCheckFile", "接口调用成功")),
		endpoint.Post("/session-archiving-update", "会话存档自建应用更新",
			endpoint.Handler(sessionArchivingUpdate),
			endpoint.Tags("设置"),
			endpoint.Body(settings.PluginRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: "",
			}, "GetSettings", "接口调用成功")),
		endpoint.Post("/chat-settings", "企微会话存档配置",
			endpoint.Handler(getChatSettings),
			endpoint.Tags("企微会话存档配置"),
			endpoint.Body(settings.AuthRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: settings.ChatSettingsResponse{},
			}, "getChatSettings", "接口调用成功")),
	}
}

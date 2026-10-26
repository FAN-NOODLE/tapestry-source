package drainage_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/drainage"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/drainage-contact/channels-create", "添加渠道",
			endpoint.Handler(channelsCreate),
			endpoint.Tags("添加渠道"),
			endpoint.Response(http.StatusOK, "", "channelCreate", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-list", "渠道列表",
			endpoint.Handler(channelsList),
			endpoint.Tags("渠道列表"),
			endpoint.Response(http.StatusOK, "", "channelList", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-detail", "渠道详情",
			endpoint.Handler(channelsDetail),
			endpoint.Tags("渠道列表"),
			endpoint.Response(http.StatusOK, "", "channelsDetail", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-edit", "渠道详情",
			endpoint.Handler(channelsEdit),
			endpoint.Tags("渠道编辑"),
			endpoint.Response(http.StatusOK, "", "channelsEdit", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-count-export", "渠道统计导出",
			endpoint.Handler(channelsCountExport),
			endpoint.Tags("渠道统计导出"),
			endpoint.Response(http.StatusOK, "", "channelsCountExport", "接口调用成功")),

		endpoint.Post("/drainage/contact/list", "联系我列表",
			endpoint.Handler(contactList),
			endpoint.Tags("联系我列表"),
			endpoint.Response(http.StatusOK, "", "contactList", "接口调用成功")),
		endpoint.Post("/drainage/contact/create", "联系我新增",
			endpoint.Handler(contactCreate),
			endpoint.Tags("联系我新增"),
			endpoint.Response(http.StatusOK, "", "contactCreate", "接口调用成功")),
		endpoint.Post("/drainage/contact/update", "联系我修改",
			endpoint.Handler(contactUpdate),
			endpoint.Tags("联系我修改"),
			endpoint.Response(http.StatusOK, "", "contactUpdate", "接口调用成功")),
		endpoint.Post("/drainage/contact/delete", "联系我删除",
			endpoint.Handler(contactDelete),
			endpoint.Tags("联系我删除"),
			endpoint.Response(http.StatusOK, "", "contactDelete", "接口调用成功")),
		endpoint.Post("/drainage/contact/details", "联系我详情",
			endpoint.Handler(contactDetails),
			endpoint.Tags("联系我详情"),
			endpoint.Response(http.StatusOK, "", "contactDelete", "接口调用成功")),
		endpoint.Post("/drainage/contact/download-qr-code", "下载联系我二维码",
			endpoint.Handler(downloadQrCode),
			endpoint.Tags("联系我详情"),
			endpoint.Response(http.StatusOK, "", "contactDelete", "接口调用成功")),

		endpoint.Post("/drainage/welcome-create", "欢迎语新建",
			endpoint.Handler(welcomeCreate),
			endpoint.Body(drainage.WelcomeCreateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeCreateResponse{},
			}, "welcomeCreate", "接口调用成功")),

		endpoint.Post("/drainage/welcome-detail", "欢迎语详情",
			endpoint.Handler(welcomeDetail),
			endpoint.Body(drainage.WelcomeDetailRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeDetailResponse{},
			}, "welcomeDetail", "接口调用成功")),

		endpoint.Post("/drainage/welcome-default", "欢迎语默认详情",
			endpoint.Handler(welcomeDefault),
			endpoint.Body(drainage.WelcomeDetailRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeDetailResponse{},
			}, "welcomeDefault", "接口调用成功")),

		endpoint.Post("/drainage/welcome-edit", "欢迎语编辑",
			endpoint.Handler(welcomeEdit),
			endpoint.Body(drainage.WelcomeEditRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeEditResponse{},
			}, "welcomeEdit", "接口调用成功")),

		endpoint.Post("/drainage/welcome-delete", "欢迎语删除",
			endpoint.Handler(welcomeDelete),
			endpoint.Body(drainage.WelcomeDeleteRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeDeleteResponse{},
			}, "welcomeDelete", "接口调用成功")),

		endpoint.Post("/drainage/welcome-list", "欢迎语列表",
			endpoint.Handler(welcomeList),
			endpoint.Body(drainage.WelcomeListRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeListResponse{},
			}, "welcomeList", "接口调用成功")),

		endpoint.Post("/drainage/welcome-channel-list", "欢迎语选择渠道列表",
			endpoint.Handler(welcomeChannelList),
			endpoint.Body(drainage.WelcomeValidateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeChannelsList{},
			}, "welcomeValidate", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-data", "渠道数据统计(按日期/按渠道)",
			endpoint.Handler(channelsData),
			endpoint.Body(drainage.ChannelsDataRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "channelsData", "接口调用成功")),

		endpoint.Post("/drainage-contact/channels-data-sum", "渠道数据统计总览",
			endpoint.Handler(channelsDataSum),
			endpoint.Body(drainage.ChannelsDataRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "channelsData", "接口调用成功")),

		endpoint.Post("/drainage/welcome-all", "获取所有的欢迎语",
			endpoint.Handler(WelcomeAll),
			endpoint.Body(drainage.WelcomeAllRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "welcome-all", "接口调用成功")),

		endpoint.Post("/drainage/welcome-del", "获取所有的欢迎语",
			endpoint.Handler(WelcomeDel),
			endpoint.Body(drainage.WelcomeBase{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "welcome-del", "接口调用成功")),

		endpoint.Post("/drainage/welcome-new-edit", "新建-编辑欢迎语",
			endpoint.Handler(WelcomeNewEdit),
			endpoint.Body(drainage.WelcomeNewEditRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "welcome-new-edit", "接口调用成功")),

		endpoint.Post("/drainage/welcome-detailnew", "欢迎语详情",
			endpoint.Handler(WelcomeDetailNew),
			endpoint.Body(drainage.WelcomeBase{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "welcome-detailnew", "接口调用成功")),

		endpoint.Post("/drainage/get-detail", "获取所有卡等信息",
			endpoint.Handler(KdTypeAll),
			endpoint.Body(drainage.WelcomeAllRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, "", "welcome-all", "接口调用成功")),

		endpoint.Post("/drainage/cache-media", "缓存media-id",
			endpoint.Handler(cacheMedia),
			endpoint.Body(drainage.CacheMediaRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.CacheMediaResponse{},
			}, "cacheMedia", "接口调用成功")),

		endpoint.Post("/drainage/channels-simple-list", "channel-simple-list",
			endpoint.Handler(channelSimpleList),
			endpoint.Response(http.StatusOK, "", "channelSimpleList", "渠道列表")),

		endpoint.Post("/drainage/all-type", "获取所有外部条件",
			endpoint.Handler(externalWelcomeTypeList),
			endpoint.Tags("获取所有外部条件"),
			endpoint.Response(http.StatusOK, "", "all-type", "接口调用成功")),
	}
}

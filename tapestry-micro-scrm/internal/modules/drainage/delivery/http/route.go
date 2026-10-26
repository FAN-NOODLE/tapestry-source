package drainagehttp

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/tag"
	tagMysql "micro-scrm/internal/modules/tag/respository/mysql"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {

	//topicService = moment.NewTopic(mysql.NewTopic())
	//originalMemberService = members.New(mysql.NewOriginalMember())

	drainageService = drainage.NewDrainage(mysql.NewDrainage())
	contactService = (*drainage.ContactService)(drainage.NewDrainage(mysql.NewDrainage()))
	tagService = tag.NewTag(tagMysql.NewTagRepo())
	enterpriseService = enterprise.New(enterpriseMysql.New())
	externalDataHandle := NewMemberLevelHandle()

	return []*swagger.Endpoint{

		endpoint.Post("/drainage/welcome-create", "欢迎语新建",
			endpoint.Handler(welcomeCreate),
			endpoint.Body(drainage.WelcomeCreateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeCreateResponse{},
			}, "welcomeCreate", "接口调用成功")),

		endpoint.Post("/drainage/contact/create", "联系我新增",
			endpoint.Handler(contactCreate),
			endpoint.Body(drainage.ContactCreateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeCreateResponse{},
			}, "contactCreate", "接口调用成功")),

		endpoint.Post("/drainage/contact/update", "联系我编辑",
			endpoint.Handler(contactUpdate),
			endpoint.Body(drainage.ContactCreateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeCreateResponse{},
			}, "contactUpdate", "接口调用成功")),

		endpoint.Post("/drainage/contact/details", "联系我详情",
			endpoint.Handler(contactDetails),
			endpoint.Body(drainage.ContactDetailsRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.Contacts{},
			}, "contactDetails", "接口调用成功")),

		endpoint.Post("/drainage/contact/delete", "联系我删除",
			endpoint.Handler(contactDelete),
			endpoint.Body(drainage.ContactDeleteRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.Contacts{},
			}, "contactDelete", "接口调用成功")),

		endpoint.Post("/drainage/contact/list", "联系我列表",
			endpoint.Handler(contactList),
			endpoint.Body(drainage.Contacts{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.Contacts{},
			}, "contactDelete", "接口调用成功")),

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

		endpoint.Post("/drainage/welcome-send", "回调发送欢迎语",
			endpoint.Handler(welcomeSend),
			endpoint.Body(drainage.WelcomeSendRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeSendResponse{},
			}, "welcomeSend", "接口调用成功")),

		endpoint.Post("/drainage/welcome-channel-list", "欢迎语渠道列表",
			endpoint.Handler(welcomeChannelList),
			endpoint.Body(drainage.WelcomeValidateRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.WelcomeChannelsList{},
			}, "welcomeValidate", "接口调用成功")),

		endpoint.Post("/drainage/channels-list", "渠道列表",
			endpoint.Handler(channelsList),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsList", "接口调用成功")),

		endpoint.Post("/drainage/channels-simple-list", "渠道总列表",
			endpoint.Handler(channelSimpleList),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsSimpleList", "接口调用成功")),

		endpoint.Post("/drainage/channels-create", "渠道新增",
			endpoint.Handler(channelsCreate),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsCreate", "接口调用成功")),

		endpoint.Post("/drainage/channels-detail", "渠道详情",
			endpoint.Handler(channelsDetail),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsDetail", "接口调用成功")),

		endpoint.Post("/drainage/channels-edit", "渠道编辑",
			endpoint.Handler(channelsEdit),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsCreate", "接口调用成功")),

		endpoint.Post("/drainage/channels-data", "渠道数据统计",
			endpoint.Handler(channelsData),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsCreate", "接口调用成功")),

		endpoint.Post("/drainage/cache-media", "缓存media-id",
			endpoint.Handler(cacheMedia),
			endpoint.Body(drainage.CacheMediaRequest{}, "", true),
			endpoint.Tags("引流"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: drainage.CacheMediaResponse{},
			}, "cacheMedia", "接口调用成功")),

		endpoint.Post("/drainage/contact/download-qr-code", "下载联系我二维码",
			endpoint.Handler(contactQrCode),
			endpoint.Tags("渠道"),
			endpoint.Response(http.StatusOK, "", "channelsCreate", "接口调用成功")),

		endpoint.Post("/drainage/welcome-new-edit", "新建欢迎语",
			endpoint.Handler(externalDataHandle.welcomeNew),
			endpoint.Tags("新建欢迎语"),
			endpoint.Response(http.StatusOK, "", "welcome-new", "接口调用成功")),

		endpoint.Post("/drainage/welcome-del", "删除欢迎语",
			endpoint.Handler(externalDataHandle.welcomeDel),
			endpoint.Tags("删除欢迎语"),
			endpoint.Response(http.StatusOK, "", "welcome-del", "接口调用成功")),

		endpoint.Post("/drainage/welcome-detailnew", "欢迎语详情",
			endpoint.Handler(externalDataHandle.welcomeDetailNew),
			endpoint.Tags("欢迎语详情"),
			endpoint.Response(http.StatusOK, "", "welcome-detailnew", "接口调用成功")),

		endpoint.Post("/drainage/welcome-all", "获取外部条件欢迎语",
			endpoint.Handler(externalDataHandle.welcomeList),
			endpoint.Tags("获取外部条件欢迎语"),
			endpoint.Response(http.StatusOK, "", "welcome-new", "接口调用成功")),

		endpoint.Post("/drainage/get-detail", "根据typeID获取外部条件",
			endpoint.Handler(externalDataHandle.externalTypeById),
			endpoint.Tags("根据typeID获取外部条件"),
			endpoint.Response(http.StatusOK, "", "type-detail-by-id", "接口调用成功")),

		endpoint.Post("/drainage/all-type", "获取所有外部条件",
			endpoint.Handler(externalDataHandle.externalWelcomeTypeList),
			endpoint.Tags("获取所有外部条件"),
			endpoint.Response(http.StatusOK, "", "all-type", "接口调用成功")),

		endpoint.Post("/drainage/update-all-type", "更新所有外部条件",
			endpoint.Handler(externalDataHandle.externalDataList),
			endpoint.Tags("更新所有外部条件"),
			endpoint.Response(http.StatusOK, "", "all-type", "接口调用成功")),

		endpoint.Post("/syn-contact", "同步渠道和联系我二维码",
			endpoint.Handler(synContact),
			endpoint.Tags("同步渠道和联系我二维码"),
			endpoint.Response(http.StatusOK, "", "synContact", "接口调用成功")),

		endpoint.Post(
			"statistics/all-day",
			"渠道统计按天展示",
			endpoint.Handler(statisticsAllByDay),
			endpoint.Tags("渠道统计"),
			endpoint.Response(http.StatusOK, "", "statistics", "接口调用成功"),
		),

		endpoint.Post(
			"statistics/all-sum",
			"渠道统计总览",
			endpoint.Handler(statisticsAllBySum),
			endpoint.Tags("渠道统计"),
			endpoint.Response(http.StatusOK, "", "statistics", "接口调用成功"),
		),

		endpoint.Post(
			"statistics/channels-sum",
			"渠道统计渠道列表总览",
			endpoint.Handler(statisticsAllByChannel),
			endpoint.Tags("渠道统计"),
			endpoint.Response(http.StatusOK, "", "statistics", "接口调用成功"),
		),

		endpoint.Post(
			"statistics/channels-day",
			"渠道统计渠道按天展示",
			endpoint.Handler(statisticsItemByDay),
			endpoint.Tags("渠道统计"),
			endpoint.Response(http.StatusOK, "", "statistics", "接口调用成功"),
		),
	}
}

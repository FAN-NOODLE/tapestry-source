package sourceHttp

import (
	"net/http"

	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"

	"micro-scrm/internal/modules/chat"
	chatRepository "micro-scrm/internal/modules/chat/repository/mysql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/drainage"
	drainageMysql "micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/source"
	"micro-scrm/internal/modules/source/respository/mysql"
	"micro-scrm/internal/modules/ws_trace"
	wsTraceMysql "micro-scrm/internal/modules/ws_trace/respository/mysql"
)

func Endpoints() []*swagger.Endpoint {

	s = source.New(mysql.NewSource())
	enterpriseService = enterprise.New(enterpriseMysql.New())
	source.DrainageService = drainage.NewDrainage(drainageMysql.NewDrainage())
	wsTraceService = ws_trace.NewService(wsTraceMysql.NewWsTrace())
	groupChatService = chat.NewGroupChat(chatRepository.NewWorkGroupChat())
	customerService = customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	return []*swagger.Endpoint{

		// 文本素材
		endpoint.Post("/source/text-list", "文本素材列表",
			endpoint.Handler(TextList),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-text", "接口调用成功")),

		endpoint.Post("/source/text-save", "文本素材保存",
			endpoint.Handler(TextSave),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-text", "接口调用成功")),

		endpoint.Post("/source/text-detail", "文本素材详情",
			endpoint.Handler(TextDetail),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-text", "接口调用成功")),

		// 图片素材
		endpoint.Post("/source/image-list", "图片素材列表",
			endpoint.Handler(SourcesImageList),
			endpoint.Tags("图片素材列表"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/image-add", "图片素材添加",
			endpoint.Handler(SourceImageAdd),
			endpoint.Tags("图片素材ADD"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 文件素材
		endpoint.Post("/source/file-list", "文件素材列表",
			endpoint.Handler(SourcesFileList),
			endpoint.Tags("素材分组列表"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source/file-add", "文件素材添加",
			endpoint.Handler(SourceFileAdd),
			endpoint.Tags("文件素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 链接素材
		endpoint.Post("/source/link-list", "链接素材列表",
			endpoint.Handler(LinkList),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-link", "接口调用成功")),

		endpoint.Post("/source/link-save", "链接素材保存",
			endpoint.Handler(LinkSave),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-link", "接口调用成功")),

		endpoint.Post("/source/link-detail", "链接素材列表",
			endpoint.Handler(LinkDetail),
			endpoint.Tags("素材中心"),
			endpoint.Response(http.StatusOK, "", "source-link", "接口调用成功")),

		// 小程序素材
		endpoint.Post("/source/mini-program-add", "小程序素材添加",
			endpoint.Handler(sourceMiniProgramAdd),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/mini-program-edit", "小程序素材添加",
			endpoint.Handler(sourceMiniProgramUpdate),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/mini-program-list", "小程序素材列表",
			endpoint.Handler(sourceMiniProgramList),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/mini-program-detail", "小程序素材详情",
			endpoint.Handler(sourceMiniProgramDetail),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 视频素材
		endpoint.Post("/source/video-list", "视频素材列表",
			endpoint.Handler(sourceVideoList),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/video-add", "视频素材添加",
			endpoint.Handler(sourceVideoAdd),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/video-edit", "视频素材编辑",
			endpoint.Handler(sourceVideoUpdate),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/video-detail", "视频素材详情",
			endpoint.Handler(sourceVideoDetail),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/delete", "素材删除",
			endpoint.Handler(SourceDelete),
			endpoint.Tags("素材删除"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/move-group", "移动分组",
			endpoint.Handler(moveGroup),
			endpoint.Tags("素材移动"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),
		// 生成临时素材
		endpoint.Post("/source/make-media", "生成临时素材",
			endpoint.Handler(makeSourceMediaId),
			endpoint.Tags("生成临时素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),
		endpoint.Post("/source/get/group_name", "获取分组名称",
			endpoint.Handler(getSourceName),
			endpoint.Tags("获取分组名称"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),
	}
}

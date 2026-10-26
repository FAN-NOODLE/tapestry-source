package sourceHttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-scrm/internal/config"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
	"net/http"
)

var sourceSDK *microScrm.SourceApi

func Endpoints() []*swagger.Endpoint {
	sourceSDK = microScrm.NewSource(config.CfgData.Service.MicroScrm.Domain)

	return []*swagger.Endpoint{
		//endpoint.Post("/source/image-list", "图片素材列表",
		//	endpoint.Handler(SourcesImageList),
		//	endpoint.Tags("图片素材列表"),
		//	endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/image-add", "图片素材添加",
			endpoint.Handler(SourceImageAdd),
			endpoint.Tags("图片素材ADD"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		//endpoint.Post("/source/file-list", "文本素材列表",
		//	endpoint.Handler(SourcesFileList),
		//	endpoint.Tags("文本素材"),
		//	endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source/file-add", "文件素材添加",
			endpoint.Handler(SourceFileAdd),
			endpoint.Tags("文件素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 小程序素材
		endpoint.Post("/source/mini-program-save", "小程序素材添加",
			endpoint.Handler(sourceMiniProgramSave),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		//// 小程序素材
		//endpoint.Post("/source/mini-program-list", "小程序素材列表",
		//	endpoint.Handler(sourceMiniProgramList),
		//	endpoint.Tags("小程序素材"),
		//	endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 小程序素材
		endpoint.Post("/source/mini-program-detail", "小程序素材详情",
			endpoint.Handler(sourceMiniProgramDetail),
			endpoint.Tags("小程序素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 视频素材
		endpoint.Post("/source/video-save", "视频素材添加",
			endpoint.Handler(sourceVideoSave),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		//// 视频素材
		//endpoint.Post("/source/video-list", "视频素材列表",
		//	endpoint.Handler(sourceVideoList),
		//	endpoint.Tags("视频素材"),
		//	endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 视频素材
		endpoint.Post("/source/video-detail", "视频素材详情",
			endpoint.Handler(sourceVideoDetail),
			endpoint.Tags("视频素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 文本素材保存
		endpoint.Post("/source/text-save", "文本素材保存",
			endpoint.Handler(sourceTextSave),
			endpoint.Tags("文本素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 链接素材保存
		endpoint.Post("/source/link-save", "链接素材保存",
			endpoint.Handler(sourceLinkSave),
			endpoint.Tags("链接素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		////文本素材列表
		//endpoint.Post("/source/text-list", "文本素材列表",
		//	endpoint.Handler(sourceTextList),
		//	endpoint.Tags("文本素材"),
		//	endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		////链接素材列表
		//endpoint.Post("/source/link-list", "链接素材列表",
		//	endpoint.Handler(sourceLinkList),
		//	endpoint.Tags("链接素材"),
		//	endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 文本素材列表
		endpoint.Post("/source/text-detail", "文本素材详情",
			endpoint.Handler(sourceTextDetail),
			endpoint.Tags("文本素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 链接素材列表
		endpoint.Post("/source/link-detail", "链接素材详情",
			endpoint.Handler(sourceLinkDetail),
			endpoint.Tags("链接素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 删除素材
		endpoint.Post("/source/delete", "删除素材",
			endpoint.Handler(sourceDelete),
			endpoint.Tags("删除素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		// 移动素材分组
		endpoint.Post("/source/move-group", "移动素材分组",
			endpoint.Handler(moveGroup),
			endpoint.Tags("链接素材"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),

		endpoint.Post("/source/get/group_name", "获取分组名称",
			endpoint.Handler(getGroupName),
			endpoint.Tags("获取分组名称"),
			endpoint.Response(http.StatusOK, "", "source", "接口调用成功")),
	}
}

package sourceGroupHttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-scrm/internal/config"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
	"net/http"
)

var sourceGroupSDK *microScrm.SourceGroupApi

func Endpoints() []*swagger.Endpoint {
	sourceGroupSDK = microScrm.NewSourceGroup(config.CfgData.Service.MicroScrm.Domain)

	return []*swagger.Endpoint{
		endpoint.Post("/source-group/add", "素材分组创建",
			endpoint.Handler(add),
			endpoint.Tags("素材分组创建"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/edit", "素材分组编辑",
			endpoint.Handler(edit),
			endpoint.Tags("素材分组编辑"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/delete", "素材分组删除",
			endpoint.Handler(del),
			endpoint.Tags("素材分组删除"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/list", "素材分组列表",
			endpoint.Handler(list),
			endpoint.Tags("素材分组列表"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/move", "素材分组移动",
			endpoint.Handler(move),
			endpoint.Tags("素材分组移动"),
			endpoint.Response(http.StatusOK, "", "source-move", "接口调用成功")),

		endpoint.Post("/source-group/update-ignore-index", "素材分组首页不可见",
			endpoint.Handler(updateIgnoreIndex),
			endpoint.Tags("素材分组首页不可见"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),
	}
}

package sourceGroupHttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/source"
	sourceMysql "micro-scrm/internal/modules/source/respository/mysql"
	"micro-scrm/internal/modules/sourceGroup"
	"micro-scrm/internal/modules/sourceGroup/respository/mysql"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	s = sourceGroup.New(mysql.NewSourceGroup())
	sou = source.New(sourceMysql.NewSource())
	return []*swagger.Endpoint{
		endpoint.Post("/source-group/add", "素材分组创建",
			endpoint.Handler(addSource),
			endpoint.Tags("素材分组创建"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/edit", "素材分组编辑",
			endpoint.Handler(editSource),
			endpoint.Tags("素材分组编辑"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/update-ignore-index", "素材分组首页不可见",
			endpoint.Handler(updateIgnoreIndex),
			endpoint.Tags("素材分组首页不可见"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/delete", "素材分组删除",
			endpoint.Handler(removeSource),
			endpoint.Tags("素材分组删除"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/list", "素材分组列表",
			endpoint.Handler(SourceTree),
			endpoint.Tags("素材分组列表"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/source-group/move", "分组位置移动",
			endpoint.Handler(SourceMove),
			endpoint.Tags("素材分组移动"),
			endpoint.Response(http.StatusOK, "", "source-move", "接口调用成功")),
	}
}

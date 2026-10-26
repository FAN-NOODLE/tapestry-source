package tag_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/user/create-group", "创建用户分组",
			endpoint.Handler(userCreateGroup),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, "", "userCreateGroup", "接口调用成功")),

		endpoint.Post("/user/group-user-import", "客户分组客户模板导入",
			endpoint.Handler(userGroupImport),
			endpoint.Tags("客户分组客户模板导入"),
			endpoint.Response(http.StatusOK,"",
			"userGroupUserImport", "接口调用成功")),

		endpoint.Post("/user/group-user-export", "客户分组客户导入的记录模板导出",
			endpoint.Handler(userGroupExport),
			endpoint.Tags("客户分组客户导入的记录模板导出"),
			endpoint.Response(http.StatusOK,"",
				"userGroupUserImport", "接口调用成功")),
		endpoint.Post("/user/group-list", "用户分组列表",
			endpoint.Handler(userGroupList),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, "",
				"userGroupList", "接口调用成功")),

		endpoint.Post("/user/group-detail", "用户分组详情",
			endpoint.Handler(userGroupDetail),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, "",
				"userGroupDetail", "接口调用成功")),

		endpoint.Post("/user/group-edit", "用户分组编辑",
			endpoint.Handler(userGroupEdit),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, "",
				"userGroupEdit", "接口调用成功")),

		endpoint.Post("/user/group-delete", "用户分组删除",
			endpoint.Handler(userGroupDelete),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, "",
				"userGroupDelete", "接口调用成功")),

		endpoint.Post("/tag/list", "标签列表",
			endpoint.Handler(tagList),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagList", "接口调用成功")),

		endpoint.Post("/tag/search", "标签/标签分组搜索",
			endpoint.Handler(tagSearch),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagList", "接口调用成功")),

		endpoint.Post("/tag/create-group", "创建标签组",
			endpoint.Handler(createGroup),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"createGroup", "接口调用成功")),

		endpoint.Post("/tag/create-tag", "新增手动标签",
			endpoint.Handler(createTag),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"createTag", "接口调用成功")),

		endpoint.Post("/tag/update", "标签自动编辑",
			endpoint.Handler(updateAutoTag),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"updateAutoTag", "接口调用成功")),

		endpoint.Post("/tag/goods-spu", "商品spu列表",
			endpoint.Handler(goodsSpu),
			endpoint.Body("", "", true),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"goodsSpu", "接口调用成功")),

		endpoint.Post("/tag/coupon-list", "微信优惠券列表",
			endpoint.Handler(wxCouponList),
			endpoint.Body("", "", true),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"wxCouponList", "接口调用成功")),

		endpoint.Post("/tag/precalculate-usercount", "预计算人数",
			endpoint.Handler(precalculateUserCount),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"precalculateUserCount", "接口调用成功")),

		endpoint.Post("/tag/tag-update", "修改子标签",
			endpoint.Handler(tagUpdate),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagUpdate", "接口调用成功")),

		endpoint.Post("/tag/tag-info", "标签组详情",
			endpoint.Handler(tagInfo),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/rank", "标签组排序",
			endpoint.Handler(tagRand),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagRand", "接口调用成功")),

		endpoint.Post("/tag/tag-group-delete", "手动标签-删除标签组",
			endpoint.Handler(tagGroupDelete),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagGroupDelete", "接口调用成功")),

		endpoint.Post("/tag/update-time", "刷新标签更新时间",
			endpoint.Handler(tagUpdateTime),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagUpdateTime", "接口调用成功")),

		endpoint.Post("/tag/tag-delete", "删除子标签",
			endpoint.Handler(tagDelete),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagUpdateTime", "接口调用成功")),

		endpoint.Post("/tag/sync-qw", "标签同步企业微信",
			endpoint.Handler(tagSyncQw),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagSyncQw", "接口调用成功")),

		endpoint.Post("/tag/test", "标签同步企业微信",
			endpoint.Handler(test),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagSyncQw", "接口调用成功")),

		endpoint.Post("/tag/export", "标签导出",
			endpoint.Handler(tagExport),
			endpoint.Tags("标签管理"),
			endpoint.Response(http.StatusOK, "",
				"tagExport", "接口调用成功")),

		endpoint.Post("/tag/mini/tag-group-detail-list", "标签/小程序/获取标签组详情列表",
			endpoint.Handler(getTagGroupDetailList),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"getTagGroupDetailList", "接口调用成功")),




		endpoint.Post("/tag/detail/list", "标签详情客户列表",
			endpoint.Handler(tagDetailList),
			endpoint.Tags("标签详情客户列表"),
			endpoint.Response(http.StatusOK, "",
				"tagDetailList", "接口调用成功")),
		endpoint.Post("/tag/detail/export", "导出标签详情客户名单",
			endpoint.Handler(tagDetailExport),
			endpoint.Tags("导出标签详情客户名单"),
			endpoint.Response(http.StatusOK, "",
				"tagDetailExport", "接口调用成功")),

		endpoint.Post("/tag/auth-create", "是否可以编辑标签",
			endpoint.Handler(tagAuthCreate),
			endpoint.Tags("是否可以编辑标签"),
			endpoint.Response(http.StatusOK, "",
				"tagAuthCreate", "接口调用成功")),


	}
}

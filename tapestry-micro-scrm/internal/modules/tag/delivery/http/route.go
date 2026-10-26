package http

import (
	"net/http"

	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketing"
	marketingMysql "micro-scrm/internal/modules/marketing/respository/mysql"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/modules/tag/respository/mysql"
)

// Endpoints 路由配置
func Endpoints() []*swagger.Endpoint {

	tagService = tag.NewTag(mysql.NewTagRepo())
	tagGroupService = tag.NewTagGroup(mysql.NewTagRepo())
	marketService = marketing.New(marketingMysql.New(), userEs.NewEsRepo())
	tagService.Market = marketService
	enterpriseService = enterprise.New(enterpriseMysql.New())
	userService = user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())

	// work系列标签服务
	workExternalTagService = tag.NewWorkExternalTag(mysql.NewWorkTagRepo())
	workExternalTagGroupService = tag.NewWorkExternalTagGroup(mysql.NewWorkTagGroupRepo())
	userGroupService = tag.NewUserGroup(mysql.NewUserGroupRepo())

	return []*swagger.Endpoint{
		endpoint.Post("/tag/list", "标签列表",
			endpoint.Handler(tagList),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"tagList", "接口调用成功")),

		endpoint.Post("/tag/create-group", "创建标签组",
			endpoint.Handler(createGroup),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"createGroup", "接口调用成功")),

		endpoint.Post("/user/create-group", "创建客户分组",
			endpoint.Handler(userCreateGroup),
			endpoint.Body(entity.CreateUserGroupBody{}, "", true),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.CreateUserGroupResponse{},
			}, "userCreateGroup", "接口调用成功")),

		endpoint.Post("/user/group-export", "客户分组客户导入的记录模板导出",
			endpoint.Handler(userGroupExport),
			endpoint.Body(entity.CreateUserGroupBody{}, "", true),
			endpoint.Tags("客户分组客户导入的记录模板导出"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.CreateUserGroupResponse{},
			}, "userGroupExport", "接口调用成功")),

		endpoint.Post("/user/group/fileDetection", "客户分组导入客户模板文件判断",
			endpoint.Handler(userGroupFileDetection),
			endpoint.Body(entity.CreateUserGroupBody{}, "", true),
			endpoint.Tags("客户分组导入客户模板文件判断"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.CreateUserGroupResponse{},
			}, "userGroupFileDetection", "接口调用成功")),


		endpoint.Post("/user/group/saveInfo", "客户分组导入客户模板文件数据保存",
			endpoint.Handler(userGroupFileSaveInfo),
			endpoint.Body(entity.CreateUserGroupBody{}, "", true),
			endpoint.Tags("客户分组导入客户模板文件数据保存"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.CreateUserGroupResponse{},
			}, "userGroupFileSaveInfo", "接口调用成功")),

		endpoint.Post("/user/group/update/is_extract", "客户分组导入客户模板文件数据完成状态更改",
			endpoint.Handler(userGroupFileUpdate),
			endpoint.Body(entity.CreateUserGroupBody{}, "", true),
			endpoint.Tags("客户分组导入客户模板文件数据完成状态更改"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.CreateUserGroupResponse{},
			}, "userGroupFileUpdate", "接口调用成功")),


		endpoint.Post("/user/group-edit", "客户分组编辑",
			endpoint.Handler(userGroupEdit),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.UserGroupList{},
			}, "userGroupEdit", "接口调用成功")),

		endpoint.Post("/user/group-delete", "客户分组删除",
			endpoint.Handler(userGroupDelete),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.UserGroupList{},
			}, "userGroupDelete", "接口调用成功")),

		endpoint.Post("/user/group-detail", "客户分组详情",
			endpoint.Handler(userGroupDetail),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.UserGroupList{},
			}, "userGroupDetail", "接口调用成功")),

		endpoint.Post("/user/group-list", "客户分组列表",
			endpoint.Handler(userGroupList),
			endpoint.Tags("客户分组"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: entity.UserGroupList{},
			},
				"userGroupList", "接口调用成功")),

		endpoint.Post("/tag/create-tag", "新增手动标签",
			endpoint.Handler(createManualTag),
			endpoint.Body(entity.CreateTagRequest{}, "", true),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.CreateTagResponse{},
				"createGroup", "接口调用成功")),

		endpoint.Post("/tag/update", "更新标签组配置信息",
			endpoint.Handler(updateAutoTag),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.AutoTagUpdateNumberRange{},
				"updateAutoTag", "接口调用成功")),

		endpoint.Post("/tag/update-sync", "更新标签组同步企业微信标识",
			endpoint.Handler(updateTagGroupSync),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"updateTagSync", "接口调用成功")),

		endpoint.Post("/user/group-precalculate", "客户分组预计算人数",
			endpoint.Handler(userGroupPrecalculate),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"userGroupPrecalculate", "接口调用成功")),

		endpoint.Post("/tag/tag-info", "标签组详情",
			endpoint.Handler(tagGroupInfo),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, "",
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/group-sort", "标签组排序",
			endpoint.Handler(sortTagGroup),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.SortTagGroupRequest{},
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/search", "标签搜索",
			endpoint.Handler(searchTag),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.SearchTagRequest{},
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/tag-group-delete", "删除手动标签组",
			endpoint.Handler(deleteManualTagGroup),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagGroupRequest{},
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/manual/update", "更新手动标签",
			endpoint.Handler(updateManualTag),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.UpdateManualTagRequest{},
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/manual/delete", "删除手动标签",
			endpoint.Handler(deleteManualTag),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagRequest{},
				"tagInfo", "接口调用成功")),

		endpoint.Post("/tag/sync", "同步cdp标签和用户分组数据",
			endpoint.Handler(TagSync),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagRequest{},
				"TagSync", "接口调用成功")),

		endpoint.Post("/tag/tag-list", "获取标签/用户分组列表",
			endpoint.Handler(getTagList),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagRequest{},
				"getTagList", "接口调用成功")),

		endpoint.Post("/tag/customer-manual-list", "获取标签/客户手动标签",
			endpoint.Handler(customerManualList),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagRequest{},
				"getTagList", "接口调用成功")),

		endpoint.Post("/tag/mini/tag-group-detail-list", "标签/小程序/获取标签组详情列表",
			endpoint.Handler(getTagGroupDetailList),
			endpoint.Tags("标签"),
			endpoint.Response(http.StatusOK, entity.DeleteManualTagGroupRequest{},
				"getTagGroupDetailList", "接口调用成功")),
	}
}

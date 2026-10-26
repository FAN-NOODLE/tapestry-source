package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	URL_TAG_DETAIL_INFO = "%s/micro-scrm/v1/tag/detail/info"
	URL_TAG_DETAIL_LIST = "%s/micro-scrm/v1/tag/detail/list"
)

var client = sdk.NewHttpClient(sdk.ParseJson{})

// TagGroupDetail 标签详情
func TagGroupDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TagGroupDetail), request)
}

// TagList 获取标签列表
func TagList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TagList), request)
}

// TagSearch 标签/标签分组搜索
func TagSearch(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagSearch), request)
	return response, err
}

// CreateGroup 创建手动标签组
func CreateGroup(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CreateGroupUrl), request)
	return response, err
}

// CreateTag 新增手动标签
func CreateTag(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CreateTagUrl), request)
	return response, err
}

// UpdateAutoTag 更新自动标签设置
func UpdateAutoTag(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.UpdateAutoTagUrl), request)
	return response, err
}

// UserCreateGroup 创建用户分组
func UserCreateGroup(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CreateUserGroupUrl), request)
	return response, err
}

// UserGroupList 获取客户分组列表
func UserGroupList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.GroupUserListUrl), request)
	return response, err
}

// GroupUserCount 客户分组预计算人数
func GroupUserCount(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.GroupUserUserCount), request)
	return response, err
}

// 编辑子标签
func TagUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagUpdateUrl), request)
	return response, err
}

// 编辑子标签
func TagInfo(c ctx.Context, request interface{}) (*TagInfoData, error) {
	var response TagInfoResp
	err := client.Post(c, GetUrl(urlEnum.TagInfoUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

// TagRand 标签组排序
func TagRand(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagRandUrl), request)
	return response, err
}

// TagGroupDelete 标签组删除
func TagGroupDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagGroupDeleteUrl), request)
	return response, err
}

// TagUpdateTime 刷新标签更新时间
func TagUpdateTime(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagUpdateTimeUrl), request)
	return response, err
}

// TagDelete 标签组删除
func TagDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TagDeleteUrl), request)
	return response, err
}

// TagSyncQw 标签同步企业微信
func TagSyncQw(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TagSyncQw), req)
}

// TagExport 标签导出
func TagExport(c ctx.Context, req interface{}) (*TagListData, error) {
	var response TagListResponse
	err := client.Post(c, GetUrl(urlEnum.TagList), req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

// UserGroupDetail 用户分组详情
func UserGroupDetail(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.UserGroupDetailUrl), request)
	return response, err
}

// UserGroupEdit 用户分组编辑
func UserGroupEdit(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.UserGroupEditUrl), request)
	return response, err
}

// UserGroupDelete 用户分组编辑
func UserGroupDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.UserGroupDeleteUrl), request)
	return response, err
}

// UserGroupExportData 用户分组的导入名单记录
func UserGroupExportData(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.UserGroupExportDataUrl), request)
	return response, err
}

func TagDetailList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(URL_TAG_DETAIL_LIST), request)
	return response, err
}

func TagAuthCreate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/tag/detail/auth"), request)
	return response, err
}

func ImpressionTagList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/tag/impression/tag/list"), request)
	return response, err
}

func ImpressionTagGroupList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/tag/impression/group/list"), request)
	return response, err
}

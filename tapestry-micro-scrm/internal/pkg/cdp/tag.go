package cdp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

var httpClient = NewHTTPClient()

// TagList 获取标签列表
func TagList(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.UserTagList), request)
}

// CreateManualTagGroup 创建手动标签组
func CreateManualTagGroup(c ctx.Context, request CreateManualTagGroupRequest) (*CreateManualTagGroupResponse, error) {
	//return sendPost(c, GetURL(urlEnum.ManualGroupCreate), request)
	var response CreateManualTagGroupResponse
	err := httpClient.Post(c, GetURL(urlEnum.ManualGroupCreate), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateUserGroup 创建用户分组
func CreateUserGroup(c ctx.Context, request interface{}) (*CreateUserGroupResponse, error) {
	//return sendPost(c, GetURL(urlEnum.CreateUserGroup), request)
	var response CreateUserGroupResponse
	err := httpClient.Post(c, GetURL(urlEnum.CreateUserGroup), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// EditUserGroup 编辑用户分组
func EditUserGroup(c ctx.Context, request interface{}) (*CreateUserGroupResponse, error) {
	var response CreateUserGroupResponse
	err := httpClient.Post(c, GetURL(urlEnum.EditUserGroup), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteUserGroup 删除用户分组
func DeleteUserGroup(c ctx.Context, request interface{}) (*CreateUserGroupResponse, error) {
	var response CreateUserGroupResponse
	err := httpClient.Post(c, GetURL(urlEnum.DeleteUserGroup), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteUserGroup 删除用户分组
func UserGroupDetail(c ctx.Context, request *UserGroupDetailRequest) (*UserGroupDetailResponse, error) {
	var response UserGroupDetailResponse
	err := httpClient.Post(c, GetURL(urlEnum.UserGroupDetail), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUserGroupList 获取用户分组列表
func GetUserGroupList(c ctx.Context, request interface{}) (*UserGroupListResponse, error) {
	var response UserGroupListResponse
	err := httpClient.Post(c, GetURL(urlEnum.UserGroupList), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateManualTag 创建手动标签
func CreateManualTag(c ctx.Context, request interface{}) (*CreateManualTagResponse, error) {
	var response CreateManualTagResponse
	err := httpClient.Post(c, GetURL(urlEnum.ManualTagCreate), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// AutoTagUpdate 创建手动标签
func AutoTagUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.AutoTagUpdate), request)
}

// UpdateTagGroupSync 创建手动标签
func UpdateTagGroupSync(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.UpdateTagGroupSync), request)
}

// TagGroupInfo 标签组详情
func TagGroupInfo(c ctx.Context, request *TagGroupInfoRequest) (*TagGroupInfoResponse, error) {
	var response TagGroupInfoResponse
	err := httpClient.Post(c, GetURL(urlEnum.TagGroupInfo), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SortTagGroup 标签组列表排序
func SortTagGroup(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.SortTagGroup), request)
}

// SearchTag 标签搜索
func SearchTag(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.SearchTag), request)
}

// DeleteManualTagGroup 删除手动标签组
func DeleteManualTagGroup(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.DeleteManualTagGroup), request)
}

// UserGroupPrecalculateCount 客户分组预计算人数
func UserGroupPrecalculateCount(c ctx.Context, request interface{}) (*UserGroupPrecalculateCountResponse, error) {
	var response UserGroupPrecalculateCountResponse
	err := httpClient.Post(c, GetURL(urlEnum.UserGroupPrecalculateCount), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdateManualTag 更新手动标签信息
func UpdateManualTag(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.UpdateManualTag), request)
}

// DeleteManualTag 删除手动标签
func DeleteManualTag(c ctx.Context, request interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.DeleteManualTag), request)
}

// TagEnterpriseInit 企业初始化
func TagEnterpriseInit(c ctx.Context, req interface{}) (interface{}, error) {
	return sendPost(c, GetURL(urlEnum.TagEnterpriseInit), req)
}

// GetCdpMarketingStrategyDetails 营销策略详情数据
func GetCdpMarketingStrategyDetails(c ctx.Context, request interface{}, enterpriseId int64) (*StrategyDataDetails, error) {
	var response StrategyDataDetails
	err := httpClient.Post(c, GetURL(urlEnum.GetMarketingStrategyDetailsUrl, enterpriseId), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func GetUserGroupsByUid(c ctx.Context, request interface{}) (*UserGroupsByUidResult, error) {
	var response UserGroupsByUidResult
	err := httpClient.Post(c, GetURL(urlEnum.GetUserGroupsByUid), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

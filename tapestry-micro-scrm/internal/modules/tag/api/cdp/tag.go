package cdp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/cdp"
	"time"
)

// TagGroupListResponse 标签组列表返回结构体
type TagGroupListResponse struct {
	ResponseBase
	Data entity.TagGroupList `json:"data"`
}

// GetTagGroupList 获取标签分组列表
func GetTagGroupList(c ctx.Context) (*entity.TagGroupList, error) {
	sdkResponse, err := cdp.TagList(c, nil)
	if err != nil {
		return nil, err
	}

	responseStr := json.MarshalToString(sdkResponse)
	var response entity.TagGroupList
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Info(c).Msgf("标签分组列表数据解析错误: %s", err.Error())
	}

	return &response, nil
}

// SearchTag 标签搜索
func SearchTag(c ctx.Context, request *entity.SearchTagRequest) (*entity.TagGroupList, error) {
	sdkResponse, err := cdp.SearchTag(c, request)
	if err != nil {
		return nil, err
	}

	responseStr := json.MarshalToString(sdkResponse)
	var response entity.TagGroupList
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Info(c).Msgf("标签搜搜数据解析错误: %s, response:%s", err.Error(), responseStr)
	}

	return &response, nil
}

// CreateTagGroupRequest 创建标签组请求参数结构体
type CreateTagGroupRequest struct {
	Name string `json:"name"`
}

// CreateTagGroup 创建标签组
func CreateTagGroup(c ctx.Context, request entity.CreateGroupRequest) (int64, error) {
	var sdkRequest = cdp.CreateManualTagGroupRequest{
		Name:     request.GroupName,
		TagNames: request.TagNames,
	}
	sdkResponse, err := cdp.CreateManualTagGroup(c, sdkRequest)
	if err != nil {
		return 0, err
	}

	return sdkResponse.Data.TagGroupID, nil
}

// CreateUserGroup 创建用户分组
func CreateUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	// sdkResponse, err := cdp.CreateUserGroup(c, param)
	// if err != nil {
	//	return nil, err
	// }
	response := entity.CreateUserGroupResponse{
		GroupID: time.Now().Unix(),
	}
	return &response, nil
}

// CreateUserGroup 创建用户分组
func EditUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	sdkResponse, err := cdp.EditUserGroup(c, param)
	if err != nil {
		return nil, err
	}
	response := entity.CreateUserGroupResponse{
		GroupID: sdkResponse.Data.GroupID,
	}
	return &response, nil
}

// DeleteUserGroup 删除用户分组
func DeleteUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	sdkResponse, err := cdp.DeleteUserGroup(c, param)
	if err != nil {
		return nil, err
	}
	response := entity.CreateUserGroupResponse{
		GroupID: sdkResponse.Data.GroupID,
	}
	return &response, nil
}

// DeleteUserGroup 删除用户分组
func UserGroupDetail(c ctx.Context, ID int64) (*entity.UserGroupDetail, error) {
	var request = cdp.UserGroupDetailRequest{
		ID: ID,
	}
	sdkResponse, err := cdp.UserGroupDetail(c, &request)
	if err != nil {
		return nil, err
	}
	response := UserGroupDetailTranslator(sdkResponse)
	return response, nil
}

// GroupUserListResponse 用户分组列表返回结构体
type GroupUserListResponse struct {
	ResponseBase
	Data entity.UserGroupList `json:"data"`
}

// GroupUserList 获取客户分组列表
func GroupUserList(c ctx.Context) (*entity.UserGroupList, error) {
	sdkResponse, err := cdp.GetUserGroupList(c, nil)
	if err != nil {
		return nil, err
	}
	userGroupList := GroupUserListResponseTranslator(sdkResponse)
	return userGroupList, nil
}

// CreateManualTag cdp创建手动标签实现
func CreateManualTag(c ctx.Context, groupID int64, name string) (int64, error) {
	var request = cdp.CreateManualTagRequest{
		TagGroupID: groupID,
		TagName:    name,
	}
	sdkResponse, err := cdp.CreateManualTag(c, request)
	if err != nil {
		return 0, err
	}
	return sdkResponse.Data.TagID, nil
}

// UpdateAutoTag 更新自动标签配置
func UpdateAutoTag(c ctx.Context, request entity.AutoTagUpdateRequest) error {
	_, err := cdp.AutoTagUpdate(c, request)
	if err != nil {
		return err
	}
	return nil

	/*
		responseStr := sendPost(c, URLEnum.AutoTagUpdate, request)
		var response CommonResponse
		helper.Must(json.Unmarshal([]byte(responseStr), &response))
		err := CheckResponse(response)
		if err != nil {
			return err
		}
		return nil

	*/
}

// TagGroupInfo 标签组详情接口
func TagGroupInfo(c ctx.Context, request *entity.TagGroupInfoRequest) (*entity.TagGroup, error) {
	var sdkRequest = cdp.TagGroupInfoRequest{
		TagGroupID: request.TagGroupID,
		Order:      request.Order,
		Page:       request.Page,
		Size:       request.Size,
	}
	sdkResponse, err := cdp.TagGroupInfo(c, &sdkRequest)
	if err != nil {
		return nil, err
	}
	tagGroup := TagGroupInfoResponseTranslator(sdkResponse)
	return tagGroup, nil
}

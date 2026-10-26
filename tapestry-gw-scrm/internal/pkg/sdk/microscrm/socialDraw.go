package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
	"mime/multipart"
)

// WorkGroupChatList 获取群列表
func WorkGroupChatList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatInfo 获取群列表
func WorkGroupChatInfo(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatInfoUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatSetting 社群设置
func WorkGroupChatSetting(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatSettingUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatTypeList 获取群列表
func WorkGroupChatTypeList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatTypesUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatTypeManage 获取群列表
func WorkGroupChatTypeManage(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatTypeManageUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatFilters 社群自定义下拉筛选
func WorkGroupChatFilters(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatFiltersUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatMembers 社群成员明细
func WorkGroupChatMembers(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkGroupChatMembersUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// WorkGroupChatDataGrid 社群详情统计数据
func WorkGroupChatDataGrid(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.WorkGroupChatDataGridUrl), request)
}

// WorkGroupUploadSetting 上传社群批量设置数据
func WorkGroupUploadSetting(c ctx.Context, params map[string]string, upFileName string, files ...*multipart.FileHeader) (interface{}, error) {
	return sdk.SendFormData(c, GetUrl(urlEnum.WorkGroupChatUploadSettingUrl), params, upFileName, files...)
}

// WorkGroupBatchSetting 批量设置社群信息
func WorkGroupBatchSetting(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.WorkGroupChatBatchSettingUrl), request)
}

// WorkGroupDataUpdateTime 社群数据最新更新时间
func WorkGroupDataUpdateTime(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.WorkGroupChatDataUpdateTimeUrl), request)
}

// WorkGroupBatchSync 批量更新社群信息
func WorkGroupBatchSync(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.WorkGroupChatBatchSyncUrl), request)
}

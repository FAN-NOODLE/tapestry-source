package gomicrouser

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// MemberLevelList 会员等级列表
func MemberLevelList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.MemberLevelListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// MemberLevelDetail 会员等级详情
func MemberLevelDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.MemberLevelDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// MemberLevelAdd 会员等级添加
func MemberLevelSet(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.MemberLevelSetUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// UpdateStatus 会员等级状态修改
func UpdateStatus(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.UpdateStatusUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ContentSort 等级内容排序
func ContentSort(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ContentSortUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelUpList 等级跃迁列表
func LevelUpList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelUpListURL), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// AddLevelUp 新增等级跃迁
func AddLevelUp(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.AddLevelUpURL), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// EditLevelUp 编辑等级跃迁
func EditLevelUp(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.EditLevelUpURL), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// DeleteLevelUp 删除等级跃迁
func DeleteLevelUp(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.DeleteLevelUpURL), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelUpInfo 等级跃迁详情
func LevelUpInfo(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelUpInfoURL), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SelectLevelList 下拉框会员等级列表
func SelectLevelList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SelectLevelListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

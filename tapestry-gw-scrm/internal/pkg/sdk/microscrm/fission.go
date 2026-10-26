package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

//ExchangeRecord 获取兑换记录
func ExportExchangeRecord(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := client.Post(c, GetUrl(urlEnum.ExportExchangeRecordUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ActivityList 裂变活动列表
func ActivityList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionActivityListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ActivityAdd 裂变活动新增
func ActivityAdd(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionActivityAddUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ActivityUpdate 裂变活动编辑
func ActivityUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionActivityUpdateUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ActivityDelete 裂变活动删除
func ActivityDelete(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionActivityDeleteUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// UpdateStatus 裂变活动更新状态
func UpdateStatus(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionUpdateStatusUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ActivityInfo 裂变活动详情
func ActivityInfo(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.FissionActivityInfoUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ChatGroupList 选择社群列表
func ChatGroupList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ChatGroupListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// MutualActivityList 互斥活动选择列表
func MutualActivityList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.MutualActivityListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//ActivityDetail 活动详情
func ActivityDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ActivityDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//ActivityDetailList 活动详情列表
func ActivityDetailList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ActivityDetailListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//InviteDetail 发起人邀请明细
func InviteDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.InviteDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//ExchangeRecord发起人兑奖明细
func ExchangeRecord(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ExchangeRecordUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//PrizeExchangeRecord  奖品兑换记录
func PrizeExchangeRecord(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ExportExchangeRecordUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//	OpenExchangeWarn奖品兑换提醒
func OpenExchangeWarn(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.OpenExchangeWarnUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//ExchangeTotal 奖品兑换记录总数
func ExchangeTotal(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ExchangeTotalUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//ActivityListSimple活动列表 全部有效活动
func ActivityListSimple(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ActivityListSimpleUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func GetExchangeConfigMember(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ExchangeConfigMemberUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func GetStyleConfig(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.GetStyleConfigUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func SetStyleConfig(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SetStyleConfigUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func AddStyleConfig(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.AddStyleConfigUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func DelStyleConfig(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.DelStyleConfigUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

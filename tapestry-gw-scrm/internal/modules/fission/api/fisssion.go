package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type FissionApi struct{}

// ActivityList 裂变活动列表
func (f *FissionApi) ActivityList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ActivityAdd 裂变活动新增
func (f *FissionApi) ActivityAdd(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityAdd(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ActivityUpdate 裂变活动编辑
func (f *FissionApi) ActivityUpdate(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityUpdate(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ActivityDelete 裂变活动删除
func (f *FissionApi) ActivityDelete(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateStatus 裂变活动更新状态
func (f *FissionApi) UpdateStatus(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.UpdateStatus(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ActivityInfo 裂变活动详情
func (f *FissionApi) ActivityInfo(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityInfo(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ChatGroupList 选择社群列表
func (f *FissionApi) ChatGroupList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatGroupList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MutualActivityList 互斥活动选择列表
func (f *FissionApi) MutualActivityList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.MutualActivityList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//ActivityDetail裂变活动明细
func (f *FissionApi) ActivityDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ActivityDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//ActivityDetailList 活动明细 列表
func (f *FissionApi) ActivityDetailList(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.ActivityDetailList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// InviteDetail发起人邀请明细
func (f *FissionApi) InviteDetail(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.InviteDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//ExchangeRecord 发起人兑换记录
func (f *FissionApi) ExchangeRecord(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.ExchangeRecord(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//prizeExchangeRecord 奖品兑换记录
func (f *FissionApi) PrizeExchangeRecord(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.PrizeExchangeRecord(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//OpenExchangeWarn 奖品兑换提醒
func (f *FissionApi) OpenExchangeWarn(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.OpenExchangeWarn(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//ExchangeTotal 奖品兑换记录总数
func (f *FissionApi) ExchangeTotal(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.ExchangeTotal(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//ActivityListSimple 活动列表 全部有效活动
func (f *FissionApi) ActivityListSimple(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.ActivityListSimple(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//GetExchangeConfigMember获取兑换提醒的额外员工
func (f *FissionApi) GetExchangeConfigMember(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.GetExchangeConfigMember(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//GetStyleConfig
func (f *FissionApi) GetStyleConfig(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.GetStyleConfig(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//SetStyleConfig
func (f *FissionApi) SetStyleConfig(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.SetStyleConfig(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//AddStyleConfig
func (f *FissionApi) AddStyleConfig(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.AddStyleConfig(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//DelStyleConfig
func (f *FissionApi) DelStyleConfig(c ctx.Context, req interface{}) (interface{}, error) {

	res, err := microscrm.DelStyleConfig(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/customer"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type CustomerApi struct {
}

// 重复客户
func (m CustomerApi) RepeatCustomer(c ctx.Context, req *customer.RepeatListRequest) interface{} {
	res, err := microscrm.RepeatCustomer(c, req)
	if err != nil {
		return nil
	}
	return res
}

// 重复客户所属员工
func (m CustomerApi) RepeatFindMember(c ctx.Context, req *customer.ExternalRelationByUserIdRequest) interface{} {
	res, err := microscrm.RepeatFindMember(c, req)
	if err != nil {
		return nil
	}
	return res
}

// 更新设置删除消息
func (m CustomerApi) UpdateDelMsgDetail(c ctx.Context, req *customer.UpdateMessageRequest) interface{} {
	res, err := microscrm.UpdateMessage(c, req)
	if err != nil {
		return nil
	}
	return res
}

// 消息列表
func (m CustomerApi) MessageList(c ctx.Context, req *customer.MessageListRequest) interface{} {
	res, err := microscrm.MessageList(c, req)
	if err != nil {
		return nil
	}
	return res
}

// 流失客户列表
func (m CustomerApi) LossCustomersList(c ctx.Context, req *customer.LossCustomersListRequest) interface{} {
	res, err := microscrm.LossCustomersList(c, req)
	if err != nil {
		return nil
	}
	return res
}

func (m CustomerApi) CustomerList(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.CustomerList(c, req)
}

func (m CustomerApi) EditCustomerRemark(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.EditCustomerRemark(c, req)
}

func (m CustomerApi) TagList(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.CustomerTagList(c, req)
}

func (m CustomerApi) ListConfig(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.ListConfig(c, req)
}

func (m CustomerApi) OrderGoods(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.OrderGoods(c, req)
}

func (m CustomerApi) BelongToChats(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.BelongToChats(c, req)
}

func (m CustomerApi) MemberLevelList(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.MemberLevelList(c, req)
}
func (m CustomerApi) SyncEterpriseUserToES(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.SyncEterpriseUserToES(c, req)
}

package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk"
)

const (
	ALLOT_EXTERNAL_LIST   = "%s/micro-scrm/v1/customer/inherit/allot/external-list"
	ALLOT_GROUP_CHAT_LIST = "%s/micro-scrm/v1/customer/inherit/allot/group-chat-list"
	//TAKE_OVER_LIST        = "%s/micro-scrm/v1/customer/inherit/takeover/list"
)

// 更新删除消息配置
func UpdateMessage(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.UpdateMessageUrl), request)
}

// MessageList
func MessageList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MessageListUrl), request)
}

// RepeatCustomer
func RepeatCustomer(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.RepeatCustomerUrl), request)
}

// RepeatFindMember
func RepeatFindMember(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.RepeatCustomerMemberUrl), request)
}

// LossCustomersList
func LossCustomersList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.LossCustomersListUrl), request)
}

func CustomerList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CustomerListUrl), request)
}

// EditCustomerRemark 编辑顾客备注
func EditCustomerRemark(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.EditCustomerRemarkUrl), request)
}

// CustomerTagList 顾客筛选使用标签列表
func CustomerTagList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CustomerTagListUrl), request)
}

// ListConfig 列表配置
func ListConfig(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CustomerListConfigUrl), request)
}

// OrderGoods 订单商品列表
func OrderGoods(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.OrderGoodsUrl), request)
}

// MemberLevelList 会员等级列表
func MemberLevelList(c ctx.Context, request interface{}) (interface{}, error) {
	//return sdk.SendPost(c, GetUrl(urlEnum.MemberLevelListURL), request)
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"work"})
	return sdk.SendPostByWork(c, ClientUrl+"/service/v1/data-lake/query-tier-list", request)
}

// BelongToChats 客户群列表
func BelongToChats(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.BelongToChatsUrl), req)
}

// 全量更新es索引
func SyncEterpriseUserToES(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.SyncEterpriseUserToESUrl), req)
}

func GetExternalAllotList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(ALLOT_EXTERNAL_LIST), req)
}

func GetGroupChatAllotList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(ALLOT_GROUP_CHAT_LIST), req)
}

func GetList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/list"), req)
}

func TransferCustomer(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/transfer_customer"), req)
}

func TransferCustomerGroup(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/transfer_customer_group"), req)
}

func GetTakeOverList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/takeover/list"), req)
}
func GetTakeOverExternalList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/takeover/external-list"), req)
}
func GetTakeOverGroupList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/takeover/group-list"), req)
}
func SyncCustomerData(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/inherit/sync-customer-data"), req)
}

func CustomerUserList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/user-list"), req)
}

func CustomerUserDetail(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/user-detail"), req)
}

func SowTemplateList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/list"), req)
}

func SowTemplateSave(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/save"), req)
}

func SowTemplateDetail(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/detail"), req)
}

func SowTemplateUpdateStatus(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/update-status"), req)
}

func SowTemplateDelete(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/delete"), req)
}

func SowAllTaskList(c ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/customer/sow-template/all-task-list"), req)
}

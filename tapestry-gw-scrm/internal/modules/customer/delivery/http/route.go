package customer_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/customer"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
	"net/http"
)

var customerSDK *microScrm.CustomerApi

func Endpoints() []*swagger.Endpoint {
	customerSDK = microScrm.NewCustomer(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post("/customer/repeat-find-member", "重复g客户所属员工列表",
			endpoint.Handler(repeatFindMember),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.ExternalRelationByUserIdRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: customer.ExternalRelationByUserIdRequest{},
			}, "memberDetail", "接口调用成功")),
		endpoint.Post("/customer/repeat-customer", "重复客户",
			endpoint.Handler(repeatCustomer),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.RequestCustomer{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []customer.RequestCustomer{},
			}, "memberList", "接口调用成功")),
		endpoint.Post("/customer/message-list", "消息通知列表",
			endpoint.Handler(customerDelMsgList),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.MessageListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []customer.MessageListRequest{},
			}, "memberList", "接口调用成功")),

		endpoint.Post("/customer/update-message", "消息通知更新",
			endpoint.Handler(updateDelMsgDetail),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.UpdateMessageRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []customer.UpdateMessageRequest{},
			}, "memberList", "接口调用成功")),

		endpoint.Post("/customer/loss-customers-list", "流失客户列表",
			endpoint.Handler(lossCustomersList),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.LossCustomersListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []customer.LossCustomersListRequest{},
			}, "memberList", "接口调用成功")),

		endpoint.Post("/customer/footprint-list", "客户足迹",
			endpoint.Handler(footprintList),
			endpoint.Tags("客户足迹"),
			endpoint.Response(http.StatusOK, &customer.ResponseFootprintList{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/order-list", "客户订单",
			endpoint.Handler(orderList),
			endpoint.Tags("客户订单"),
			endpoint.Response(http.StatusOK, &customer.ResponseOrderList{}, "customer", "接口调用成功")),
		endpoint.Post("/customer/order-export", "客户订单导出",
			endpoint.Handler(orderExport),
			endpoint.Tags("客户订单"),
			endpoint.Response(http.StatusOK, &customer.ResponseOrderList{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/customer-detail", "客户详情",
			endpoint.Handler(customerDetail),
			endpoint.Tags("客户详情"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/belong-to-members", "所属员工列表",
			endpoint.Handler(belongToMembers),
			endpoint.Tags("所属员工列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/belong-to-chats", "所属群列表",
			endpoint.Handler(belongToChats),
			endpoint.Tags("所属群列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/consignee-address-list", "收货地址列表",
			endpoint.Handler(consigneeAddressList),
			endpoint.Tags("收货地址列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/list", "客户列表",
			endpoint.Handler(customerList),
			endpoint.Tags("客户列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/edit-remark", "编辑备注",
			endpoint.Handler(editCustomerRemark),
			endpoint.Tags("客户列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/tag-list", "标签列表",
			endpoint.Handler(tagList),
			endpoint.Tags("标签列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/list-config", "客户列表配置",
			endpoint.Handler(listConfig),
			endpoint.Tags("客户列表配置"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),

		endpoint.Post("/customer/order-goods", "订单商品列表",
			endpoint.Handler(orderGoods),
			endpoint.Tags("订单商品列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),
		endpoint.Post("/customer/member-level-list", "会员等级列表",
			endpoint.Handler(memberLevelList),
			endpoint.Tags("会员等级列表"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),
		endpoint.Post("/customer/all-eterpriseuser-fores", "全量更新es索引",
			endpoint.Handler(SyncEterpriseUserToES),
			endpoint.Tags("全量更新es索引"),
			endpoint.Response(http.StatusOK, "", "customer", "接口调用成功")),
	}

}

package customerHttp

import (
	"micro-scrm/internal/modules/guide"
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/drainage"
	drainageMysql "micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	guideMysql "micro-scrm/internal/modules/guide/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	s = customer.New(mysql.NewCustomer(), mysql.NewWorkUser(), userEs.NewEsRepo())
	m = customer.NewDelMessage(mysql.NewDelMessage())
	enterpriseService = enterprise.New(enterpriseMysql.New())
	memberService = members.New(membersMysql.New())
	customer.DrainageContact = drainage.NewContact(drainageMysql.NewDrainage())
	customer.GuideService = guide.New(guideMysql.NewGuide())

	return []*swagger.Endpoint{
		endpoint.Post("/customer/footprint-list", "客户足迹",
			endpoint.Handler(footprintList),
			endpoint.Tags("客户足迹"),
			endpoint.Response(http.StatusOK, &customer.ResponseFootprintList{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/order-list", "客户订单",
			endpoint.Handler(orderList),
			endpoint.Tags("客户订单"),
			endpoint.Response(http.StatusOK, &customer.ResponseOrderList{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/order-detail", "客户订单详情",
			endpoint.Handler(orderDetail),
			endpoint.Tags("客户订单详情"),
			endpoint.Response(http.StatusOK, &customer.ResponseOrderDetail{}, "customer", "接口调用成功")),

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

		endpoint.Post("/customer/loss-customers-list", "流失客户列表",
			endpoint.Handler(lossCustomersList),
			endpoint.Body(customer.LossCustomersListRequest{}, "", true),
			endpoint.Tags("客户列表"),
			endpoint.Response(http.StatusOK, customer.LossCustomersListResponse{}, "LossCustomersList", "接口调用成功")),

		endpoint.Post("/customer/set-message", "设置删除消息",
			endpoint.Handler(addCustomerDelMessage),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.AddMessageRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: customer.DelMessage{},
			}, "customerDelMessage", "接口调用成功")),
		endpoint.Post("/customer/message-detail", "设置删除消息",
			endpoint.Handler(customerDelMsgDetail),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.MessageDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: customer.DelMessage{},
			}, "customerDelMessage", "接口调用成功")),
		endpoint.Post("/customer/message-list", "删除消息列表",
			endpoint.Handler(customerDelMsgList),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.MessageDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: customer.DelMessage{},
			}, "customerDelMessage", "接口调用成功")),
		endpoint.Post("/customer/update-message", "更新消息",
			endpoint.Handler(updateDelMsgDetail),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.UpdateMessageRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []struct{}{},
			}, "customerDelMessage", "接口调用成功")),
		endpoint.Post("/customer/send-message", "发送消息",
			endpoint.Handler(sendMessage),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.SendMessageRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []struct{}{},
			}, "SendMessage", "接口调用成功")),
		endpoint.Post("/customer/repeat-customer", "重复客户",
			endpoint.Handler(repeatCustomer),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.SendMessageRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []struct{}{},
			}, "repeatCustomer", "接口调用成功")),
		endpoint.Post("/customer/repeat-find-member", "重复客户查询员工数据",
			endpoint.Handler(FindExternalRelationByUserId),
			endpoint.Tags("客户统计"),
			endpoint.Body(customer.ExternalRelationByUserIdRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []struct{}{},
			}, "repeatCustomer", "接口调用成功")),

		endpoint.Post("/customer/order-goods", "订单商品列表",
			endpoint.Handler(orderGoods),
			endpoint.Tags("订单商品列表"),
			endpoint.Response(http.StatusOK, &customer.ResponseOrderList{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/portrait-detail", "客户画像详情",
			endpoint.Handler(portraitDetail),
			endpoint.Tags("客户画像详情"),
			endpoint.Body(customer.PortraitRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, &customer.PortraitResponse{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/remark-save", "客户备注信息保存",
			endpoint.Handler(remarkSave),
			endpoint.Tags("客户备注信息"),
			endpoint.Body(customer.RemarkInfo{}, "请求参数", true),
			endpoint.Response(http.StatusOK, &customer.ReturnUID{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/set-external-tags", "设置外部联系人企业标签",
			endpoint.Handler(SetExternalTags),
			endpoint.Tags("设置外部联系人企业标签"),
			endpoint.Body(customer.SetExternalTags{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: struct{}{},
			}, "SetExternalTags", "接口调用成功")),
		endpoint.Post("/customer/all-eterpriseuser-fores", "同步EnterpriseUser至es",
			endpoint.Handler(SyncEterpriseUserToES),
			endpoint.Tags("同步EnterpriseUser至es"),
			endpoint.Body(customer.PortraitRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, &customer.PortraitResponse{}, "customer", "接口调用成功")),

		endpoint.Post("/customer/lost-option", "丢失客户记录,额外数据", endpoint.Handler(lostOption)),
		endpoint.Post("/customer/lost-record", "丢失客户记录", endpoint.Handler(lostRecord)),
	}
}

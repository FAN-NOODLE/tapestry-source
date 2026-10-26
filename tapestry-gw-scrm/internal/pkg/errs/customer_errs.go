package scrm_errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// FootprintListErr 客户足迹
	FootprintListErr = errorx.NewErrorMsg(120001, "客户足迹")
	// OrderListErr 客户订单列表
	OrderListErr = errorx.NewErrorMsg(120002, "客户订单列表")
	// BelongToMembersErr 所属员工列表
	BelongToMembersErr = errorx.NewErrorMsg(120003, "所属员工列表")
	// ConsigneeAddressListErr 收货地址列表
	ConsigneeAddressListErr = errorx.NewErrorMsg(120004, "收货地址列表")
	// CustomerDetailErr 客户详情
	CustomerDetailErr = errorx.NewErrorMsg(120005, "客户详情")
	// BelongToChatErr 所属群列表
	BelongToChatErr = errorx.NewErrorMsg(120005, "所属群列表")
)

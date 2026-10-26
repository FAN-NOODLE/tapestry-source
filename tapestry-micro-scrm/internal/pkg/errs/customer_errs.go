package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// FootprintListErr 客户足迹
	FootprintListErr = errorx.NewErrorMsg(160001, "获取客户足迹失败")
	// OrderListErr 客户订单列表
	OrderListErr = errorx.NewErrorMsg(160002, "客户订单列表")
	// BelongToMembersErr 所属员工列表
	BelongToMembersErr = errorx.NewErrorMsg(160003, "所属员工列表")
	// ConsigneeAddressListErr 收货地址列表
	ConsigneeAddressListErr = errorx.NewErrorMsg(160004, "收货地址列表")
	// CustomerDetailErr 客户详情
	CustomerDetailErr  = errorx.NewErrorMsg(160005, "客户详情")
	CustomerCouponsErr = errorx.NewErrorMsg(160005, "客户卡券列表")
	// BelongToChatErr 所属群列表
	BelongToChatErr = errorx.NewErrorMsg(160006, "所属群列表")
	// ShopUserAddressIDErr
	ShopUserAddressIDErr = errorx.NewErrorMsg(160007, "用户收货地址id错误")
	GetUserInfoErr       = errorx.NewErrorMsg(160008, "获取用户信息失败")
)

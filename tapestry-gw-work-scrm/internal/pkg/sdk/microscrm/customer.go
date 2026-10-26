package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// CustomerApi 员工 SDK
type CustomerApi struct {
	BaseSDK
}

// NewCustomer 返回API实例
func NewCustomer(host string) *CustomerApi {
	urlHost = host
	return &CustomerApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// FootPrintList 足迹列表
func (c CustomerApi) FootPrintList(ctx ctx.Context, req CustomerFootPrintRequest) (interface{}, error) {
	path := "/micro-scrm/v1/customer/footprint-list"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// Detail 客户画像详情
func (c CustomerApi) Detail(ctx ctx.Context, req *CustomerDetailRequest) (*PortraitResponse, error) {
	path := "/micro-scrm/v1/customer/portrait-detail"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &PortraitResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		if result.Code == 999 {
			return nil, errorx.NewErrorMsg(999, "很抱歉，客户数据暂未同步完成，请稍后再试，谢谢！")
		}
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data.(*PortraitResponse), err
}

func (c CustomerApi) CouponsList(ctx ctx.Context, req *CustomerCouponsRequest) (interface{}, error) {
	path := "/micro-scrm/v1/customer/customer-coupons"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// BelongToMembers 所属员工列表
func (c CustomerApi) BelongToMembers(ctx ctx.Context, req *CustomerRequest) (interface{}, error) {
	path := "/micro-scrm/v1/customer/belong-to-members"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// BelongToChat 所属员工列表
func (c CustomerApi) BelongToChat(ctx ctx.Context, req *CustomerRequest) (*GroupChatList, error) {
	path := "/micro-scrm/v1/customer/belong-to-chats"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return &result.Data, err
}

// RemarkSave 客户备注信息保存
func (c *CustomerApi) RemarkSave(ctx ctx.Context, customer *RemarkInfo) (interface{}, error) {
	path := "/micro-scrm/v1/customer/remark-save"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// orderList 客户订单列表
func (c *CustomerApi) OrderList(ctx ctx.Context, customer *RequestOrderList) (interface{}, error) {
	path := "/micro-scrm/v1/customer/order-list"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// orderDetail 客户订单详情
func (c *CustomerApi) OrderDetail(ctx ctx.Context, customer *RequestOrderDetail) (*ResponseOrderDetail, error) {
	path := "/micro-scrm/v1/customer/order-detail"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseOrderDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*ResponseOrderDetail), err
}

// SetExternalTags 设置外部联系人企业标签
func (c *CustomerApi) SetExternalTags(ctx ctx.Context, customer *SetExternalTags) (interface{}, error) {
	path := "/micro-scrm/v1/customer/set-external-tags"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// CustomerCouponList 客户订单详情
func (c *CustomerApi) CustomerCouponList(ctx ctx.Context, customer *CustomerRequest) (interface{}, error) {
	path := "/micro-scrm/v1/coupon/customer"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// CustomerAddressList 客户地址列表
func (c *CustomerApi) CustomerAddressList(ctx ctx.Context, customer *CustomerRequest) (interface{}, error) {
	path := "/micro-scrm/v1/customer/consignee-address-list"
	url := urlHost + path
	res, err := c.Post(ctx, url, customer)
	if err != nil {
		return nil, err
	}
	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// UserDetail 客户画像详情
func (c CustomerApi) UserDetail(ctx ctx.Context, req *CustomerDetailRequest) (*PortraitDetailResponse, error) {
	path := "/micro-scrm/v1/customer/portrait-detail"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &PortraitDetailResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		if result.Code == 999 {
			return nil, errorx.NewErrorMsg(999, "很抱歉，客户数据暂未获取，请稍后再试，谢谢！")
		}
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data.(*PortraitDetailResponse), err
}

// GetUserTags 客户标签
func (c CustomerApi) GetUserTags(ctx ctx.Context, req *UserTagsRequest) ([]UserTagsInfo, error) {
	path := "/micro-scrm/v1/tag/detail/user-tags"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}
	result := &UserTagsRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	// 按照GroupID从小到大排序
	//sort.Slice(result.Data, func(i, j int) bool {
	//	return result.Data[i].GroupID < result.Data[j].GroupID
	//})
	//
	//for _, v := range result.Data {
	//	sort.Slice(v.TagsList, func(i, j int) bool {
	//		return v.TagsList[i].TagID < v.TagsList[j].TagID
	//	})
	//}

	return result.Data, err
}

// GetUserImpressionTags 客户印象标签
func (c CustomerApi) GetUserImpressionTags(ctx ctx.Context, req *UserTagsRequest) ([]UserImpressionTagsInfo, error) {
	path := "/micro-scrm/v1/tag/impression/user-tags-list"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}
	result := &UserImpressionTagsRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	// 按照GroupID从小到大排序
	//sort.Slice(result.Data, func(i, j int) bool {
	//	return result.Data[i].GroupID < result.Data[j].GroupID
	//})
	//
	//for _, v := range result.Data {
	//	sort.Slice(v.ImpressionTagsList, func(i, j int) bool {
	//		return v.ImpressionTagsList[i].TagID < v.ImpressionTagsList[j].TagID
	//	})
	//}

	return result.Data, err
}

// UpdateRemarkForRelation 修改备注
func (c CustomerApi) UpdateRemarkForRelation(ctx ctx.Context, req *UpdateRemarkRequest) (interface{}, error) {
	path := "/micro-scrm/v1/customer/update-remark"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

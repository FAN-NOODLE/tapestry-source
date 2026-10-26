package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// CustomerApi 客户 SDK
type CustomerApi struct {
	BaseSDK
}

func NewCustomer(host string) *CustomerApi {
	urlHost = host
	return &CustomerApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// FootprintList 客户足迹
func (api *CustomerApi) FootprintList(c ctx.Context, customer *RequestCustomer) (int64, *ResponseFootprintList, error) {
	path := "/micro-scrm/v1/customer/footprint-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &ResponseFootprintList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	if result.Data.List == nil || result.Data.Total == 0 {
		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*ResponseFootprintList), err
}

// OrderList 客户订单
func (api *CustomerApi) OrderList(c ctx.Context, customer *RequestCustomer) (int64, *ResponseOrderList, error) {
	path := "/micro-scrm/v1/customer/order-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &ResponseOrderList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data.List == nil || result.Data.Total == 0 {

		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*ResponseOrderList), err
}

// BelongToMembers 所属员工列表
func (api *CustomerApi) BelongToMembers(c ctx.Context, customer *RequestCustomer) (int64, *ResponseMemberList, error) {
	path := "/micro-scrm/v1/customer/belong-to-members"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &ResponseMemberList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data.List == nil || result.Data.Total == 0 {
		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*ResponseMemberList), err
}

// BelongToChats 所属群列表
func (api *CustomerApi) BelongToChats(c ctx.Context, customer *RequestCustomerGroup) (int64, *BelongToGroupChatList, error) {
	path := "/micro-scrm/v1/customer/belong-to-chats"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &BelongToGroupChatList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data.List == nil || result.Data.Total == 0 {
		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*BelongToGroupChatList), err
}

// ConsigneeAddressList 收货人列表
func (api *CustomerApi) ConsigneeAddressList(c ctx.Context, customer *RequestCustomer) (int64, *ResponseAddressList, error) {
	path := "/micro-scrm/v1/customer/consignee-address-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &ResponseAddressList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data.List == nil || result.Data.Total == 0 {

		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*ResponseAddressList), err
}

// CustomerDetail 客户详情
func (api *CustomerApi) CustomerDetail(c ctx.Context, customer *RequestCustomer) (*ResponseCustomerDetail, error) {
	path := "/micro-scrm/v1/customer/customer-detail"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseCustomerDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*ResponseCustomerDetail), err
}

// RemarkSave 客户备注信息保存
func (api *CustomerApi) RemarkSave(c ctx.Context, customer *RemarkInfo) (*ReturnUID, error) {
	path := "/micro-scrm/v1/customer/remark-save"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ReturnUID{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*ReturnUID), err
}

// orderDetail 客户订单详情
func (api *CustomerApi) orderDetail(c ctx.Context, customer *RequestOrderDetail) (*ResponseOrderDetail, error) {
	path := "/micro-scrm/v1/customer/order-detail"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
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
func (api *CustomerApi) SetExternalTags(c ctx.Context, customer *SetExternalTags) (interface{}, error) {
	path := "/micro-scrm/v1/customer/set-external-tags"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ReturnID{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*ReturnID), err
}

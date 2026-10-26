package workwechat

import workSDK "micro-scrm/internal/pkg/sdk/work"

// OnTransfer_AllotCustomer 对外提供BA详情
func (w *WorkWeChat) OnTransfer_AllotCustomer(request workSDK.TransferCustomerAllotRequest) (*workSDK.TransferCustomerAllotResponse, error) {

	data, err := w.obj.OnTransfer_AllotCustomer(w.ctx, request)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// OnTransfer_QueryCustomerStatus 查询客户接替状态
func (w *WorkWeChat) OnTransfer_QueryCustomerStatus(req workSDK.TransferResultRequst) (result *workSDK.TransferResultRespone, err error) {

	result, err = w.obj.OnTransfer_QueryCustomerStatus(w.ctx, req)
	if err != nil {
		return
	}

	return
}

// GetUnassignedList 查询客户接替状态
func (w *WorkWeChat) GetUnassignedList(req workSDK.GetUnassignedListRequest) (result *workSDK.GetUnassignedListResponse, err error) {

	result, err = w.obj.GetUnassignedList(w.ctx, req)
	if err != nil {
		return
	}

	return
}

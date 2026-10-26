package work

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
)

const (
	// URL_OnTransfer_AllotCustomer 分配在职成员的客户
	URL_OnTransfer_AllotCustomer = "cgi-bin/externalcontact/transfer_customer"
	// TRANSFERRESULTURL 查询客户接替状态
	TRANSFERRESULTURL      = "cgi-bin/externalcontact/transfer_result"
	// ResignedTransferResult 查询客户接替状态
	ResignedTransferResult = "cgi-bin/externalcontact/resigned/transfer_result"
	// ResignedTransferCustomer 分配离职成员的客户
	ResignedTransferCustomer = "cgi-bin/externalcontact/resigned/transfer_customer"
	// ResignedTransferCustomerGroup 分配离职成员的客户群
	ResignedTransferCustomerGroup = "cgi-bin/externalcontact/groupchat/transfer"

	GET_UNASSIGNED_LIST = "cgi-bin/externalcontact/get_unassigned_list"
)

// OnTransfer_AllotCustomer 分配在职成员的客户
func (w *QyWorkImpl) OnTransfer_AllotCustomer(ctx ctx.Context, req TransferCustomerAllotRequest) (*TransferCustomerAllotResponse, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = URL_OnTransfer_AllotCustomer
	req.ReqMethod = "POST"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}

	response := &TransferCustomerAllotResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// OnTransfer_QueryCustomerStatus 查询客户接替状态
func (w *QyWorkImpl) OnTransfer_QueryCustomerStatus(ctx ctx.Context, req TransferResultRequst) (*TransferResultRespone, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = TRANSFERRESULTURL
	req.ReqMethod = POST
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	response := new(TransferResultRespone)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return response, wErr
	}
	return response, nil
}

// ResignedTransferCustomer 分配离职成员的客户
func (w *QyWorkImpl) ResignedTransferCustomer(ctx ctx.Context, req ResignedTransferCustomerReq) (*ResignedTransferCustomerResp, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ResignedTransferCustomer
	req.ReqMethod = POST
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	var response ResignedTransferCustomerResp
	if err = workRequest(ctx, url, req, &response); err != nil {
		return &response, err
	}
	return &response, nil
}

// ResignedTransferCustomerGroup 分配离职成员的客户群
func (w *QyWorkImpl) ResignedTransferCustomerGroup(ctx ctx.Context, req ResignedTransferCustomerGroupReq) (*ResignedTransferCustomerGroupReps, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ResignedTransferCustomerGroup
	req.ReqMethod = POST
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	var response ResignedTransferCustomerGroupReps
	if err = workRequest(ctx, url, req, &response); err != nil {
		return &response, err
	}
	return &response, nil
}

// GetUnassignedList 获取待分配的离职成员列表
func (w *QyWorkImpl) GetUnassignedList(ctx ctx.Context, req GetUnassignedListRequest) (*GetUnassignedListResponse, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = GET_UNASSIGNED_LIST
	req.ReqMethod = "POST"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}

	response := &GetUnassignedListResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (w *QyWorkImpl) GetResignedTransferResult(ctx ctx.Context, req GetResignedTransferResultRequest) (*ResignedTransferResultResponse,error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ResignedTransferResult
	req.ReqMethod = POST
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	var response ResignedTransferResultResponse
	if err = workRequest(ctx, url, req, &response); err != nil {
		return &response, err
	}
	return &response, nil
}

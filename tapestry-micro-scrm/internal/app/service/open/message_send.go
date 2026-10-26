package openSvc

import (
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	openReq "micro-scrm/internal/app/http/request/open"
	"micro-scrm/internal/app/service"
)

type MessageSendService struct {
	service.Base
	ExternalPayOrderDao *microScrmDao.ExternalPayOrderDao
}

func NewMessageSendService() *MessageSendService {
	return &MessageSendService{
		ExternalPayOrderDao: microScrmDao.NewExternalPayOrderDao(),
	}
}

func (s MessageSendService) GetOrderMemberUserId(c app.Ctx, params openReq.SendWaitPaymentRequest) (string, error) {

	order, err := s.ExternalPayOrderDao.GetOrderMemberUserId(c, params.OrderNo, params.EnterpriseID)
	if err != nil {
		return "", err
	}

	memberUserID := order.PayeeUserid

	return memberUserID, nil
}

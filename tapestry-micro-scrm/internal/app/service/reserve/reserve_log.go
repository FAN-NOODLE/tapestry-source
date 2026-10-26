package reserveSvc

import (
	"encoding/json"
	"micro-scrm/internal/app"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
)

type ReserveLog struct {
	service.Base
	LogDao *reserveDao.ReservesLog
}

func NewReserveLog() *ReserveLog {
	return &ReserveLog{
		LogDao: reserveDao.NewReserveLog(),
	}
}

// CreateLog 操作记录记入日志,eventName 操作行文名称
func (s *ReserveLog) CreateLog(c app.Ctx, enterpriseId, orderId int64, eventName string, operaType int) error {
	orderInfo, err := NewOrderService().GetOrderInfo(c, enterpriseId, orderId)
	if err != nil {
		return err
	}
	operateName := orderInfo.MemberName
	if operaType == reserveRequest.ReserveOperaType_Customer {
		operateName = orderInfo.UserName
	}
	var logInfo reserveRequest.LogInfo
	logInfo.MemberName = orderInfo.MemberName
	logInfo.Event = eventName
	logInfo.OpreateName = operateName
	marshal, err := json.Marshal(logInfo)
	if err != nil {
		return err
	}
	orderMarshal, err := json.Marshal(orderInfo)
	if err != nil {
		return err
	}
	var LogInfo = microScrmModel.ReserveLogModel{
		BaseModel:      models.CreateBase(app.UIDContext(c)),
		LogInfo:        string(marshal),
		ReserveOrderID: orderInfo.ID,
		OrderInfo:      string(orderMarshal),
	}
	return s.LogDao.Create(c, enterpriseId, LogInfo)
}

func (s *ReserveLog) ListLog(c app.Ctx, enterpriseId, reserveOrderID int64, logTypes []int) ([]microScrmModel.ReserveLogModel, error) {

	list, err := s.LogDao.List(c, enterpriseId, reserveOrderID, logTypes)
	if err != nil {
		return nil, err
	}

	return list, nil

}

// CreateScoreLog  用户评分记入日志,eventName 操作行文名称
func (s *ReserveLog) CreateScoreLog(c app.Ctx, request reserveRequest.UserScoreReq, eventName string, operaType int) error {
	orderInfo, err := NewServiceScoreService().ScoreInfo(c, request.EnterpriseId, request.UserOrderBase)
	if err != nil {
		return err
	}
	operateName := orderInfo.MemberName
	if operaType == 2 {
		operateName = orderInfo.UserName
	}
	var logInfo reserveRequest.LogInfo
	logInfo.MemberName = orderInfo.MemberName
	logInfo.Event = eventName
	logInfo.OpreateName = operateName
	marshal, err := json.Marshal(logInfo)
	if err != nil {
		return err
	}
	orderMarshal, err := json.Marshal(orderInfo)
	if err != nil {
		return err
	}
	var LogInfo = microScrmModel.ReserveLogModel{
		BaseModel:      models.CreateBase(app.UIDContext(c)),
		LogInfo:        string(marshal),
		ReserveOrderID: orderInfo.ReserveOrderID,
		OrderInfo:      string(orderMarshal),
		LogType:        2,
	}
	return s.LogDao.Create(c, request.EnterpriseId, LogInfo)
}

func (s *ReserveLog) CreateShareScoreLog(c app.Ctx, enterpriseId, orderId int64, eventName string, operaType int) error {
	orderInfo, err := NewOrderService().GetOrderInfo(c, enterpriseId, orderId)
	if err != nil {
		return err
	}
	operateName := orderInfo.MemberName
	if operaType == 2 {
		operateName = orderInfo.UserName
	}
	var logInfo reserveRequest.LogInfo
	logInfo.MemberName = orderInfo.MemberName
	logInfo.Event = eventName
	logInfo.OpreateName = operateName
	marshal, err := json.Marshal(logInfo)
	if err != nil {
		return err
	}
	orderMarshal, err := json.Marshal(microScrmModel.ReserveServiceScoreModel{
		BaseModel:      models.BaseModel{},
		ReserveNumber:  "",
		ReserveOrderID: 0,
		UserName:       "",
		UserUnionid:    "",
		UserExternalID: "",
		MemberID:       0,
		MemberName:     "",
		MemberUserID:   "",
		StoreID:        0,
		StoreName:      "",
		ServiceScore:   0,
		StoreScore:     0,
		Remark:         "",
		FlowScore:      0,
		ReserveTime:    0,
	})
	if err != nil {
		return err
	}
	var LogInfo = microScrmModel.ReserveLogModel{
		BaseModel:      models.CreateBase(app.UIDContext(c)),
		LogInfo:        string(marshal),
		ReserveOrderID: orderInfo.ID,
		OrderInfo:      string(orderMarshal),
		LogType:        2,
	}
	return s.LogDao.Create(c, enterpriseId, LogInfo)
}

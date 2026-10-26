package reserveSvc

import (
	"micro-scrm/internal/app"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
)

type ServiceScoreService struct {
	service.Base
	ServiceScoreDao *reserveDao.ServiceScoreDao
}

func NewServiceScoreService() *ServiceScoreService {
	return &ServiceScoreService{
		ServiceScoreDao: reserveDao.NewServiceScoreDao(),
	}
}

func (s *ServiceScoreService) GetList(c app.Ctx, enterpriseId int64, serviceScoreListReq reserveRequest.ServiceScoreListReq) ([]microScrmModel.ReserveServiceScoreModel, int64, error) {
	return s.ServiceScoreDao.GetList(c, enterpriseId, serviceScoreListReq)
}

// UserMakeScore 用户评分
func (s *ServiceScoreService) UserMakeScore(c app.Ctx, enterpriseId int64, request reserveRequest.UserScoreReq, orderInfo microScrmModel.ReserveOrderModel) (int64, error) {

	id := app.UIDContext(c)
	var ScoreModel = microScrmModel.ReserveServiceScoreModel{
		BaseModel:      models.CreateBase(id),
		ReserveNumber:  request.ReserveNumber,
		UserName:       orderInfo.UserName,
		UserUnionid:    orderInfo.UserUnionid,
		UserExternalID: orderInfo.UserExternalID,
		MemberID:       orderInfo.MemberID,
		MemberUserID:   orderInfo.MemberUserID,
		MemberName:     orderInfo.MemberName,
		Remark:         request.Remark,
		ReserveOrderID: orderInfo.ID,
		StoreID:        orderInfo.StoreID,
		StoreName:      orderInfo.StoreName,
		StoreScore:     request.StoreScore,
		ServiceScore:   request.ServiceScore,
		FlowScore:      request.FlowScore,
		ReserveTime:    orderInfo.ReserveTime,
	}
	err := s.ServiceScoreDao.UserMakeScore(c, enterpriseId, ScoreModel)
	return id, err
}

func (s *ServiceScoreService) ScoreInfo(c app.Ctx, enterpriseId int64, request reserveRequest.UserOrderBase) (microScrmModel.ReserveServiceScoreModel, error) {
	return s.ServiceScoreDao.GetScoreInfo(c, enterpriseId, request)
}

// UserEditScore coach 用户评论预约单信息
func (s *ServiceScoreService) UserEditScore(c app.Ctx, enterpriseId int64, request reserveRequest.UserScoreReq) error {
	return s.ServiceScoreDao.UserEditScore(c, enterpriseId, request)
}

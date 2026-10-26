package customerSvc

import (
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerModels "micro-scrm/internal/app/models/customer"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
)

type CustomerInfoAdd struct {
	service.Base
	customerInfoAddDao *customerDao.CustomerInfoAdd
}

func NewCustomerInfoAdd() *CustomerInfoAdd {
	return &CustomerInfoAdd{
		customerInfoAddDao: customerDao.NewCustomerInfoAdd(),
	}
}

//更新用户附加信息
func (s CustomerInfoAdd) SaveInfoAdd(c app.Ctx, request customerReq.CustomerInfoAddReq) error {
	condParam := microScrmModel.EnterpriseUserModel{ID: request.EnterpriseUserID}
	enterpriseUserInfo, err := microScrmDao.NewEnterpriseUser().GetInfo(c, condParam, request.EnterpriseID)
	if err != nil {
		return err
	}
	var unionidMap map[string]string
	helper.Unmarshal(enterpriseUserInfo.WxUnion, &unionidMap)
	info, err := s.customerInfoAddDao.GetInfo(c, request.EnterpriseID, unionidMap["default"])
	if err != nil {
		return err
	}
	if info.ID > 0 {
		err := s.customerInfoAddDao.Save(c, info.ID, request)
		if err != nil {
			return err
		}
	} else {
		var addInfo customerModels.CustomerInfoAddModel
		addInfo.Address = request.Address
		addInfo.ShoeSize = request.ShoeSize
		addInfo.Month = request.Month
		addInfo.Day = request.Day
		addInfo.OperateName = request.OperateName
		addInfo.OperateUserId = request.OperateUserId
		addInfo.EnterpriseID = request.EnterpriseID
		addInfo.ExternalUnionId = unionidMap["default"]
		addInfo.ID = uid.ID(c)
		err := s.customerInfoAddDao.Create(c, addInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

//获取用户附加信息
func (s CustomerInfoAdd) GetInfoAdd(c app.Ctx, request customerReq.GetCustomerInfoAddReq) (customerReq.UserInfoAdd, error) {
	var defInfo customerReq.UserInfoAdd
	info, err := s.customerInfoAddDao.GetInfo(c, request.EnterpriseID, request.UnionID)
	if err != nil {
		return defInfo, err
	}
	var res customerReq.UserInfoAdd
	res.Month = info.Month
	res.Day = info.Day
	res.ShoeSize = info.ShoeSize
	res.Address = info.Address
	return res, err
}

//获取前一天更新的数据
func (s CustomerInfoAdd) GetListByTime(c app.Ctx, enterpriseId int64) ([]customerModels.CustomerInfoAddModel, error) {
	var defInfo []customerModels.CustomerInfoAddModel
	todayTime, _ := helper.PeriodOfTimestamp()
	info, err := s.customerInfoAddDao.GetListByTime(c, enterpriseId, todayTime-86400)
	if err != nil {
		return defInfo, err
	}

	return info, err
}

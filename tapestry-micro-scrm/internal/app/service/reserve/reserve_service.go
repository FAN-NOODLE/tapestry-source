package reserveSvc

import (
	"encoding/json"
	"errors"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
)

type ReserveServiceService struct {
	service.Base
	ReserveServiceDao *reserveDao.ReserveServiceDao
}

func NewReserveServiceService() *ReserveServiceService {
	return &ReserveServiceService{
		ReserveServiceDao: reserveDao.NewReserveServiceDao(),
	}
}

func (s *ReserveServiceService) Create(c app.Ctx, enterpriseId int64, createReq reserveRequest.ReserveServiceCreateReq) (err error) {

	timeSlotByte, err := json.Marshal(createReq.TimeSlot)
	if err != nil {
		return err
	}
	timeSlot := string(timeSlotByte)
	var InfoModel = microScrmModel.ReserveServiceModel{
		BaseModel:   models.CreateBase(app.UIDContext(c)),
		ReserveName: createReq.ReserveName,
		ReserveImg:  createReq.ReserveImg,
		ReserveTime: createReq.ReserveTime,
		EndTime:     createReq.EndTime,
		ServiceType: createReq.ServiceType,
		StartTime:   createReq.StartTime,
		ServiceRole: createReq.ServiceRole,
		TimeSlot:    timeSlot,
	}
	// 这里需要将StoreIds,MemberIds放到关联表里
	serviceId, err := s.ReserveServiceDao.Create(c, enterpriseId, InfoModel)
	if err != nil {
		return err
	}
	err = s.ReserveServiceDao.CreateServiceRelation(c, enterpriseId, serviceId, createReq.RoleRules)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReserveServiceService) Update(c app.Ctx, enterpriseId int64, updateReq reserveRequest.ReserveServiceUpdateReq) (err error) {
	// 处理邀请函模板引用时间
	reserveServiceInfo, _ := s.ReserveServiceDao.GetInfoById(c, enterpriseId, updateReq.ID)
	// 长期有效变指定时间
	if reserveServiceInfo.ServiceType == 1 && updateReq.ServiceType == 0 {
		updateDataMap := map[string]interface{}{
			"end_time":    updateReq.EndTime,
			"start_time":  updateReq.StartTime,
			"status_type": 1,
		}
		err = reserveDao.NewInviteTemplateDao().UpdateForReserveServiceId(c, enterpriseId, updateReq.ID, updateDataMap)
		if err != nil {
			return err
		}
	}
	if reserveServiceInfo.ServiceType == 0 && updateReq.ServiceType == 0 {
		// 有效期时间缩小则变化
		if updateReq.StartTime > reserveServiceInfo.StartTime || updateReq.EndTime < reserveServiceInfo.EndTime {
			updateDataMap := map[string]interface{}{
				"end_time":   updateReq.EndTime,
				"start_time": updateReq.StartTime,
			}
			err = reserveDao.NewInviteTemplateDao().UpdateForReserveServiceId(c, enterpriseId, updateReq.ID, updateDataMap)
			if err != nil {
				return err
			}
		}
	}

	err = s.ReserveServiceDao.Update(c, enterpriseId, updateReq)
	if err != nil {
		return err
	}
	err = s.ReserveServiceDao.DeleteServiceRelation(c, enterpriseId, updateReq.ID)
	if err != nil {
		return err
	}
	err = s.ReserveServiceDao.CreateServiceRelation(c, enterpriseId, updateReq.ID, updateReq.RoleRules)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReserveServiceService) Delete(c app.Ctx, enterpriseId int64, deleteReq reserveRequest.ReserveServiceDeleteReq) (err error) {

	err = s.ReserveServiceDao.Delete(c, enterpriseId, deleteReq)
	if err != nil {
		return err
	}
	err = s.ReserveServiceDao.DeleteServiceRelation(c, enterpriseId, deleteReq.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetList 后台管理的服务配置列表
func (s *ReserveServiceService) GetList(c app.Ctx, enterpriseId int64, listReq reserveRequest.ReserveServiceListReq) ([]reserveRequest.ReserveServiceListRes, int64, error) {
	result, total, err := s.ReserveServiceDao.GetList(c, enterpriseId, listReq)
	if err != nil {
		return nil, 0, err
	}
	roleRulesMap := map[int64][]reserveRequest.RoleRules{}
	var res []reserveRequest.ReserveServiceListRes
	for _, model := range result {
		info := reserveRequest.ReserveServiceListRes{
			ReserveServiceModel: model,
			RoleRules:           []reserveRequest.RoleRules{},
			TimeSlot:            []reserveRequest.TimeSlot{},
		}
		if model.ServiceRole == 1 {
			roles, err := s.ReserveServiceDao.GetServiceRole(c, enterpriseId, model.ID)
			if err != nil {
				roleRulesMap[model.ID] = []reserveRequest.RoleRules{}
				continue
			}
			for _, role := range roles {
				item := reserveRequest.RoleRules{
					ID:                 role.MemberID,
					ThirdId:            role.StoreID,
					StoreType:          role.StoreType,
					MemSubscribeStatus: 1,
					UserId:             role.UserId,
				}
				if role.IsStore == 0 { // 不是门店，name为导购名称
					item.IsDP = false
					item.Name = role.MemberName
				} else {
					item.Name = role.StoreName
					item.IsDP = true
					item.ID = role.DepartmentId
				}
				roleRulesMap[model.ID] = append(roleRulesMap[model.ID], item)
			}
			info.RoleRules = roleRulesMap[model.ID]
		}
		timeSlot := make([]reserveRequest.TimeSlot, 0)
		if model.TimeSlot != "" {
			err := json.Unmarshal([]byte(model.TimeSlot), &timeSlot)
			if err != nil {
				continue
			}
		}
		info.TimeSlot = timeSlot
		res = append(res, info)
	}
	return res, total, err
}

func (s *ReserveServiceService) GetInfoById(c app.Ctx, enterpriseId, id int64) (microScrmModel.ReserveServiceModel, error) {
	return s.ReserveServiceDao.GetInfoById(c, enterpriseId, id)
}

func (s *ReserveServiceService) CheckDelete(c app.Ctx, request reserveRequest.ReserveServiceDeleteReq) (int, error) {
	info, err := reserveDao.NewReservesOrder().GetOrderInfoByServiceId(c, request.EnterpriseId, request.ID)
	if err != nil {
		return 2, err
	}
	if info.ID == 0 {
		return 2, nil
	}
	return 1, err
}

// GetListByMember 导购获取服务配置列表
func (s *ReserveServiceService) GetListByMember(c app.Ctx, enterpriseId int64, mobile string) ([]reserveRequest.ReserveServiceListByMemberRes, error) {
	var result []reserveRequest.ReserveServiceListByMemberRes
	memberId := int64(0)
	var storeIds []int64
	memberInfos, err := membersDao.NewMembers().GetMemberInfoAndAllDepartmentsByMobile(c, enterpriseId, mobile)
	if err != nil {
		return nil, errors.New("查询导购信息异常")
	}
	memberId = memberInfos[0].MembersId
	for _, v := range memberInfos {
		storeIds = append(storeIds, v.DepartmentThirdId)
	}
	info, err := s.ReserveServiceDao.GetListByMember(c, enterpriseId, memberId, storeIds)
	if err != nil {
		return result, err
	}
	for _, model := range info {
		value := reserveRequest.ReserveServiceListByMemberRes{
			ReserveServiceModel: model,
		}
		timeSlot := make([]reserveRequest.TimeSlot, 0)
		if model.TimeSlot != "" {
			err := json.Unmarshal([]byte(model.TimeSlot), &timeSlot)
			if err != nil {
				continue
			}
		}
		value.TimeSlot = timeSlot
		result = append(result, value)
	}

	return result, err
}

// ServiceInfo 获取服务详情接口
func (s *ReserveServiceService) ServiceInfo(c app.Ctx, request reserveRequest.ReserveServiceDeleteReq) (reserveRequest.ReserveServiceListRes, error) {
	var result reserveRequest.ReserveServiceListRes
	info, err := s.ReserveServiceDao.GetInfoById(c, request.EnterpriseId, request.ID)
	if err != nil {
		return result, err
	}
	result.ReserveServiceModel = info
	result.RoleRules = []reserveRequest.RoleRules{}
	if info.ServiceRole == 1 {
		roles, err := s.ReserveServiceDao.GetServiceRole(c, request.EnterpriseId, request.ID)
		if err != nil {
			return result, err
		}
		for _, role := range roles {
			item := reserveRequest.RoleRules{
				ID:                 role.MemberID,
				ThirdId:            role.StoreID,
				StoreType:          role.StoreType,
				MemSubscribeStatus: 1,
				UserId:             role.UserId,
			}
			if role.IsStore == 0 { // 不是门店，name为导购名称
				item.IsDP = false
				item.Name = role.MemberName
			} else {
				item.Name = role.StoreName
				item.IsDP = true
				item.ID = role.DepartmentId
			}
			result.RoleRules = append(result.RoleRules, item)
		}
	}
	timeSlot := make([]reserveRequest.TimeSlot, 0)
	if info.TimeSlot != "" {
		err = json.Unmarshal([]byte(info.TimeSlot), &timeSlot)
		result.TimeSlot = timeSlot
	}
	return result, err
}

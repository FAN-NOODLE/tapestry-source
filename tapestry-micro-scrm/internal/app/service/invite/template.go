package inviteSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"math"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	inviteDao "micro-scrm/internal/app/dao/micro_scrm"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"
)

type InviteTemplateService struct {
	service.Base
	inviteTemplateDao *inviteDao.InviteTemplateDao
	inviteRecordDao   *inviteDao.InviteRecordDao
	reserveServiceDao *inviteDao.ReserveServiceDao
}

func NewInviteTemplateService() *InviteTemplateService {
	return &InviteTemplateService{
		inviteTemplateDao: inviteDao.NewInviteTemplateDao(),
		inviteRecordDao:   inviteDao.NewInviteRecordDao(),
		reserveServiceDao: inviteDao.NewReserveServiceDao(),
	}
}

func (s *InviteTemplateService) Create(c app.Ctx, enterpriseId int64, createReq inviteRequest.TemplateSaveReq) (createId int64, err error) {

	createId = app.UIDContext(c)
	if createReq.StatusType == microScrmModel.StatusTypeLongTerm {
		createReq.StartTime = 0
		createReq.EndTime = math.MaxInt32
	}
	var InfoModel = microScrmModel.InviteTemplateModel{
		BaseModel:        models.CreateBase(createId),
		ReserveServiceID: createReq.ReserveServiceId,
		TitleName:        createReq.TitleName,
		StatusType:       createReq.StatusType,
		StartTime:        createReq.StartTime,
		EndTime:          createReq.EndTime,
		BackgroundURL:    helper.MarshalToString(createReq.BackgroundUrl),
		Content:          createReq.Content,
		ShareTitle:       createReq.ShareTitle,
		ShareURL:         helper.MarshalToString(createReq.ShareUrl),
		IsUp:             microScrmModel.IsUp,
		UserIsShow:       createReq.UserIsShow,
		MemberIsShow:     createReq.MemberIsShow,
	}

	err = s.inviteTemplateDao.Create(c, enterpriseId, InfoModel)
	if err != nil {
		return 0, err
	}

	return createId, nil
}

func (s *InviteTemplateService) Update(c app.Ctx, enterpriseId int64, updateReq inviteRequest.TemplateSaveReq) (err error) {

	nowTime := time.Now().Unix()
	if updateReq.StatusType == microScrmModel.StatusTypeLongTerm {
		updateReq.StartTime = 0
		updateReq.EndTime = math.MaxInt32
	}
	updateDataMap := map[string]interface{}{

		"reserve_service_id": updateReq.ReserveServiceId,
		"title_name":         updateReq.TitleName,
		"status_type":        updateReq.StatusType,
		"start_time":         updateReq.StartTime,
		"end_time":           updateReq.EndTime,
		"background_url":     helper.MarshalToString(updateReq.BackgroundUrl),
		"content":            updateReq.Content,
		"share_title":        updateReq.ShareTitle,
		"share_url":          helper.MarshalToString(updateReq.ShareUrl),
		"updated_at":         nowTime,
		"user_is_show":       updateReq.UserIsShow,
		"member_is_show":     updateReq.MemberIsShow,
	}

	err = s.inviteTemplateDao.Update(c, enterpriseId, updateReq.ID, updateDataMap)
	if err != nil {
		return err
	}

	return nil
}

//func (s *InviteTemplateService) handleDepartment(c app.Ctx, enterpriseId, id int64, thirdIds []int) (inviteRelationData []microScrmModel.InviteRelationModel) {
//	// 部门信息
//	departmentInfoList, _ := departmentsDao.NewDepartments().GetDepartmentsInfoThirdIds(c, enterpriseId, thirdIds)
//	departmentInfoListMap := make(map[int]string, 0)
//	for _, item := range departmentInfoList {
//		departmentInfoListMap[item.ThirdID] = item.Name
//	}
//
//	for _, thirdId := range thirdIds {
//		if _, ok := departmentInfoListMap[thirdId]; ok {
//
//			thirdName := departmentInfoListMap[thirdId]
//			// 添加关联关系
//			relationData := microScrmModel.InviteRelationModel{
//				BaseModel:        models.CreateBase(app.UIDContext(c)),
//				InviteTemplateID: id,
//				ThirdID:          thirdId,
//				ThirdName:        thirdName,
//			}
//			inviteRelationData = append(inviteRelationData, relationData)
//		}
//	}
//	return
//}

func (s *InviteTemplateService) List(c app.Ctx, enterpriseId int64, reqList inviteRequest.TemplateListReq) (lists []microScrmModel.InviteTemplateInfoList, total int64, err error) {

	list, total, err := s.inviteTemplateDao.List(c, enterpriseId, reqList)

	if err != nil {
		return nil, 0, err
	}

	// 查询邀请函对应的服务权限
	var inviteIds []int64
	var serviceIds []int64
	for _, item := range list {
		inviteIds = append(inviteIds, item.ID)
		serviceIds = append(serviceIds, item.ReserveServiceID)
	}
	serviceListMap, err1 := s.reserveServiceDao.GetListsByServiceIds(c, enterpriseId, serviceIds)
	if err != nil {
		return nil, 0, err1
	}

	// 查询关联的订单数
	orderCountNum, err := s.inviteRecordDao.GetRecordNumByInviteIds(c, enterpriseId, inviteIds)
	orderCountNumMap := make(map[int64]int, 0)
	for _, val := range orderCountNum {
		orderCountNumMap[val.InviteTemplateId] = val.Num
	}

	nowTime64 := time.Now().Unix()
	strInt64 := strconv.FormatInt(nowTime64, 10)
	nowTime, _ := strconv.Atoi(strInt64)
	for _, item := range list {
		var tempRes microScrmModel.InviteTemplateInfoList

		tempRes.ID = item.ID
		tempRes.UpdatedAt = item.UpdatedAt
		tempRes.ReserveServiceID = item.ReserveServiceID
		tempRes.ReserveServiceName = item.ReserveServiceName
		tempRes.ReserveServiceRole = item.ReserveServiceRole
		tempRes.TitleName = item.TitleName
		tempRes.BackgroundURL = item.BackgroundURL
		tempRes.Content = item.Content
		tempRes.ShareURL = item.ShareURL
		tempRes.ThirdIds = item.ThirdIds
		tempRes.ThirdNames = item.ThirdNames
		tempRes.StartTime = item.StartTime
		tempRes.EndTime = item.EndTime
		tempRes.StatusType = item.StatusType
		tempRes.IsUp = item.IsUp
		tempRes.UserIsShow = item.UserIsShow
		tempRes.MemberIsShow = item.MemberIsShow

		// 权限部门 没有就是全部
		if _, ok := serviceListMap[item.ReserveServiceID]; ok {
			var temp = serviceListMap[item.ReserveServiceID]
			for _, v := range temp {
				tempRes.RoleRules = append(tempRes.RoleRules,
					microScrmModel.RoleRule{
						ID:                 v.ID,
						ThirdId:            v.ThirdId,
						Name:               v.Name,
						IsDP:               v.IsDP,
						StoreType:          v.StoreType,
						MemSubscribeStatus: v.MemSubscribeStatus,
						UserId:             v.UserId,
					},
				)
				tempRes.ThirdNames += v.Name + ","
			}
		} else {
			tempRes.ThirdNames = "全部"
		}

		if _, ok := orderCountNumMap[item.ID]; ok {
			tempRes.InviteOrderNum = orderCountNumMap[item.ID]
		}

		if item.StatusType == microScrmModel.StatusTypeLongTerm {
			tempRes.StatusTimeStr = "长期有效"
		}

		if nowTime < item.EndTime {
			tempRes.StatusStr = "有效"
			tempRes.StatusTimeStr = helper.GetFormatDateDay(int64(item.StartTime)) + "~" + helper.GetFormatDateDay(int64(item.EndTime))
		} else {
			tempRes.StatusStr = "失效"
			tempRes.StatusTimeStr = helper.GetFormatDateDay(int64(item.StartTime)) + "~" + helper.GetFormatDateDay(int64(item.EndTime))
		}

		if item.IsUp == microScrmModel.IsUp {
			tempRes.IsUpStr = "上架"
		}
		if item.IsUp == microScrmModel.IsDown {
			tempRes.IsUpStr = "下架"
		}

		lists = append(lists, tempRes)
	}

	return

}

func (s *InviteTemplateService) Info(c app.Ctx, enterpriseId int64, reqInfo inviteRequest.TemplateInfoReq) (info microScrmModel.InviteTemplateInfo, err error) {

	info, err = s.inviteTemplateDao.Info(c, enterpriseId, reqInfo)

	if err != nil {
		return microScrmModel.InviteTemplateInfo{}, err
	}

	return

}

func (s *InviteTemplateService) Delete(c app.Ctx, enterpriseId int64, req inviteRequest.TemplateDeleteReq) (err error) {

	err = s.inviteTemplateDao.Delete(c, enterpriseId, req)

	if err != nil {
		return err
	}

	return

}

func (s *InviteTemplateService) ConfigList(c app.Ctx, enterpriseId int64, reqList inviteRequest.TemplateConfigListReq) (lists []microScrmModel.InviteTemplateListInfo, total int64, err error) {

	// 导购部门信息
	department, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByMobile(c, enterpriseId, reqList.Mobile)

	reqList.ThirdId = department.DepartmentThirdId
	reqList.MemberId = department.MembersId

	// 查询导购有权限的邀请函模板
	list, total, err := s.inviteTemplateDao.ConfigList(c, enterpriseId, reqList)
	if err != nil {
		return nil, 0, err
	}

	if total > 0 {

		for _, item := range list {
			var tempRes microScrmModel.InviteTemplateListInfo

			tempRes.ID = item.ID
			tempRes.CreatedAt = item.CreatedAt
			tempRes.UpdatedAt = item.UpdatedAt
			tempRes.ReserveServiceID = item.ReserveServiceID
			tempRes.ReserveServiceName = item.ReserveServiceName
			tempRes.ReserveServiceTime = item.ReserveServiceTime
			tempRes.TitleName = item.TitleName
			tempRes.BackgroundURL = item.BackgroundURL
			tempRes.Content = item.Content
			tempRes.ShareTitle = item.ShareTitle
			tempRes.ShareURL = item.ShareURL
			tempRes.ThirdIds = item.ThirdIds
			tempRes.ThirdNames = item.ThirdNames
			tempRes.StartTime = item.StartTime
			tempRes.EndTime = item.EndTime
			tempRes.StatusType = item.StatusType
			tempRes.IsUp = item.IsUp
			tempRes.UserIsShow = item.UserIsShow
			tempRes.MemberIsShow = item.MemberIsShow

			if item.ThirdNames == "" {
				tempRes.ThirdNames = "全部"
			}

			lists = append(lists, tempRes)
		}
	}

	return

}

func (s *InviteTemplateService) ConfigDetail(c app.Ctx, enterpriseId int64, reqDetail inviteRequest.TemplateConfigDetailReq) (data microScrmModel.InviteTemplateDetailInfo, err error) {

	return s.inviteTemplateDao.ConfigDetail(c, enterpriseId, reqDetail)

}

func (s *InviteTemplateService) CreateRecord(c app.Ctx, enterpriseId int64, createReq inviteRequest.RecordInfoReq) (info microScrmModel.InviteRecordModel, err error) {

	// 获取导购详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		Mobile: createReq.Mobile,
	})

	// 获取发送用户信息
	userInfo, err := inviteDao.NewWxUserCorpExternal().FindOne(c, enterpriseId, createReq.UserOpenid)
	if err != nil {
		return microScrmModel.InviteRecordModel{}, err
	}
	createId := app.UIDContext(c)
	var InfoModel = microScrmModel.InviteRecordModel{
		BaseModel:        models.CreateBase(createId),
		InviteTemplateID: createReq.InviteTemplateId,
		MemberID:         memberInfo.ID,
		MemberUserid:     memberInfo.Userid,
		UserName:         userInfo.Nick,
		UserUnionid:      userInfo.Unionid,
		UserOpenid:       createReq.UserOpenid,
		StyleType:        createReq.StyleType,
		BackgroundUrl:    createReq.BackgroundUrl,
		SendContent:      createReq.SendContent,
	}

	info, err = s.inviteRecordDao.Create(c, enterpriseId, InfoModel)
	if err != nil {
		return microScrmModel.InviteRecordModel{}, err
	}

	return
}

func (s *InviteTemplateService) GetRecordInfo(c app.Ctx, enterpriseId, id int64) (data inviteRequest.RecordDetailInfo, err error) {

	return s.inviteRecordDao.GetRecordInfo(c, enterpriseId, id)

}

func (s *InviteTemplateService) GetRecordOrderList(c app.Ctx, enterpriseId int64, reqList inviteRequest.RecordOrderListReq) (list []*inviteRequest.RecordOrder, total int64, err error) {

	list, total, err = s.inviteRecordDao.GetRecordOrderList(c, enterpriseId, reqList)
	if err != nil {
		return nil, 0, err
	}

	for _, item := range list {
		currentTime := item.ReserveTime
		// 预约时间加 默认1小时
		//addHourStr := fmt.Sprintf("%d", item.ReserveServiceTime) + "h"
		addHourStr := fmt.Sprintf("%d", 1) + "h"
		h, _ := time.ParseDuration(addHourStr)

		t := int64(currentTime) //外部传入的时间戳(秒为单位)，必须为int64类型
		timeStr := time.Unix(t, 0).Add(h).Format("15:04")

		if item.ReserveTime > 0 {
			item.ReserveTimeStr = time.Unix(int64(item.ReserveTime), 0).Format("2006-01-02 15:04") + "~" + timeStr
		}
		item.CreateTimeStr = time.Unix(int64(item.CreatedAt), 0).Format("2006-01-02 15:04:05")
	}

	return
}

func (s *InviteTemplateService) RecordDelete(c ctx.Context, enterpriseId, id int64) error {
	return s.inviteRecordDao.RecordDelete(c, enterpriseId, id)
}

func (s *InviteTemplateService) UpdateStatus(c app.Ctx, enterpriseId int64, req inviteRequest.TemplateUpdateStatusReq) (err error) {

	err = s.inviteTemplateDao.UpdateStatus(c, enterpriseId, req)

	if err != nil {
		return err
	}

	return

}

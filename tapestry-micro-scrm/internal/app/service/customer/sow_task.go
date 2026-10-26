package customerSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	customerReq "micro-scrm/internal/app/http/request/customer"
	planReq "micro-scrm/internal/app/http/request/plan"
	customerResp "micro-scrm/internal/app/http/response/customer"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/customer"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"time"
)

type SowTask struct {
	service.Base
	sowTaskDao *customerDao.SowTask
}

func NewSowTask() *SowTask {
	return &SowTask{
		sowTaskDao: customerDao.NewSowTask(),
	}
}

func (s *SowTask) List(c app.Ctx, data customerReq.SowTaskList) ([]customerResp.SowTaskList, int64, error) {
	var resp []customerResp.SowTaskList
	// 查询员工企微id
	memberInfo, err := s.GetUserIdByMobile(c, data.EnterpriseID, data.MemberMobile)
	if err != nil {
		return resp, 0, err
	}

	// 获取当前导购的外部联系人
	memberUserRelationList, err := microScrmDao.NewUserRelation().GetMany(c, map[string]interface{}{
		"enterprise_id": data.EnterpriseID,
		"member_id":     memberInfo.ID,
	})
	if err != nil {
		return resp, 0, err
	}

	// 当前导购的所有外部联系人id以及专属导购id
	var workExternalUserIds []int64
	var workExternalUserYesKsaIds []int64
	var workExternalUserNoKsaIds []int64
	for _, user := range memberUserRelationList {
		workExternalUserIds = append(workExternalUserIds, user.WorkExternalUserID)
		if user.IsKeySa == 1 {
			workExternalUserYesKsaIds = append(workExternalUserYesKsaIds, user.WorkExternalUserID)
		} else {
			workExternalUserNoKsaIds = append(workExternalUserNoKsaIds, user.WorkExternalUserID)
		}
	}

	// 外部联系人id集合
	var corpExternalIds []int64

	// 顾问筛选(1管理客户,2非管理客户)
	if data.SelfManager > 0 {
		tmpCorpExternalIds := workExternalUserNoKsaIds
		if data.SelfManager == 1 {
			tmpCorpExternalIds = workExternalUserYesKsaIds
		}
		if len(tmpCorpExternalIds) == 0 {
			return resp, 0, err
		}
		corpExternalIds = tmpCorpExternalIds
	}

	// 企微标签筛选
	if len(data.WorkExternalTagIds) > 0 {
		var relationIds []int64
		for _, user := range memberUserRelationList {
			relationIds = append(relationIds, user.ID)
		}
		// 获取外部联系人用户标签关系
		workExternalUserTagList, err := tagsDao.NewWorkExternalUserTag().GetMany(c, map[string]interface{}{
			"enterprise_id":   data.EnterpriseID,
			"external_tag_id": data.WorkExternalTagIds,
			"relation_id":     relationIds,
		})
		if err != nil {
			return resp, 0, err
		}
		if len(workExternalUserTagList) == 0 {
			return resp, 0, err
		}

		var tmpCorpExternalIds []int64
		for _, userTag := range workExternalUserTagList {
			for _, userRelation := range memberUserRelationList {
				if userTag.RelationID == userRelation.ID {
					tmpCorpExternalIds = append(tmpCorpExternalIds, userRelation.WorkExternalUserID)
				}
			}
		}
		if len(corpExternalIds) > 0 && len(tmpCorpExternalIds) > 0 {
			corpExternalIds = helper.ArrayIntersect[int64](corpExternalIds, tmpCorpExternalIds)
		} else {
			corpExternalIds = tmpCorpExternalIds
		}
	}

	// 印象标签筛选
	if len(data.ImpressionTagIds) > 0 {
		impressionRelationList, err := tagsDao.NewImpressionRelation().GetMany(c, map[string]interface{}{
			"enterprise_id": data.EnterpriseID,
			"user_id":       memberInfo.Userid,
			"tag_id":        data.ImpressionTagIds,
		})
		if err != nil {
			return resp, 0, err
		}

		var enterpriseUserIds []int64
		for _, user := range impressionRelationList {
			enterpriseUserIds = append(enterpriseUserIds, user.ExternalUserID)
		}
		corpExternalList, err := microScrmDao.NewWxUserCorpExternalCk().GetMany(c, data.EnterpriseID,
			map[string]interface{}{
				"euid": enterpriseUserIds,
			})
		var tmpCorpExternalIds []int64
		for _, user := range corpExternalList {
			tmpCorpExternalIds = append(tmpCorpExternalIds, user.ID)
		}

		if len(tmpCorpExternalIds) == 0 {
			return resp, 0, err
		}
		if len(corpExternalIds) > 0 && len(tmpCorpExternalIds) > 0 {
			corpExternalIds = helper.ArrayIntersect[int64](corpExternalIds, tmpCorpExternalIds)
		} else {
			corpExternalIds = tmpCorpExternalIds
		}
	}

	// 是否筛选顾客姓名,会员生日月份时间,最近一次订单消费时间
	if data.Username != "" || (data.BirthdayStart > 0 && data.BirthdayEnd > 0) || (data.LastPurchaseStartTime > 0 && data.LastPurchaseEndTime > 0) {
		var corpExternalList []microScrmModel.WxUserCorpExternalModel
		if data.Username != "" {
			corpExternalList, err = microScrmDao.NewWxUserCorpExternalCk().GetMany(c, data.EnterpriseID, "nick like ?", "%"+data.Username+"%")
		}

		if data.BirthdayStart > 0 && data.BirthdayEnd > 0 {
			var birthdayArr []string
			for i := data.BirthdayStart; i <= data.BirthdayEnd; i++ {
				birthdayArr = append(birthdayArr, helper.EncodeToString(strconv.Itoa(i)))
			}
			corpExternalList, err = microScrmDao.NewWxUserCorpExternalCk().GetMany(c, data.EnterpriseID, "birth_month in (?)", birthdayArr)
		}

		if data.LastPurchaseStartTime > 0 && data.LastPurchaseEndTime > 0 {
			corpExternalList, err = microScrmDao.NewWxUserCorpExternalCk().GetMany(c, data.EnterpriseID, "last_purchase_date >= ? and last_purchase_date < ?", data.LastPurchaseStartTime, data.LastPurchaseEndTime)
		}

		if err != nil {
			return resp, 0, err
		}
		var tmpCorpExternalIds []int64
		for _, user := range corpExternalList {
			tmpCorpExternalIds = append(tmpCorpExternalIds, user.ID)
		}
		if len(tmpCorpExternalIds) == 0 {
			return resp, 0, err
		}
		if len(corpExternalIds) > 0 && len(tmpCorpExternalIds) > 0 {
			corpExternalIds = helper.ArrayIntersect[int64](corpExternalIds, tmpCorpExternalIds)
		} else {
			corpExternalIds = tmpCorpExternalIds
		}
	}

	// 会员等级 clt_type
	if len(data.CardType) > 0 || len(data.CltType) > 0 {
		corpExternalList, err := microScrmDao.NewWxUserCorpExternal().GetManyByCardLevelOrCltType(c,
			data.EnterpriseID,
			data.CardType,
			data.CltType,
			workExternalUserIds,
		)

		if err != nil {
			return resp, 0, err
		}
		var tmpCorpExternalIds []int64
		for _, user := range corpExternalList {
			tmpCorpExternalIds = append(tmpCorpExternalIds, user.ID)
		}
		if len(tmpCorpExternalIds) == 0 {
			return resp, 0, err
		}
		if len(corpExternalIds) > 0 && len(tmpCorpExternalIds) > 0 {
			corpExternalIds = helper.ArrayIntersect[int64](corpExternalIds, tmpCorpExternalIds)
		} else {
			corpExternalIds = append(corpExternalIds, tmpCorpExternalIds...)
		}
	}

	// 状态筛选
	status := []int{data.Status}
	if data.Status == 0 {
		status = []int{customerModels.SowTaskStatusUnfinished, customerModels.SowTaskStatusFinished}
	}

	query := map[string]interface{}{
		"enterprise_id": data.EnterpriseID,
		"member_id":     memberInfo.ID,
		"status":        status,
	}

	if data.WorkExternalUserID > 0 {
		corpExternalIds = append(corpExternalIds, data.WorkExternalUserID)
	}
	if len(corpExternalIds) > 0 {
		query["work_external_user_id"] = corpExternalIds
	}

	option := make(map[string]interface{})
	// 时间筛选
	if data.Time > 0 {
		zeroTime := helper.GetZeroDateByParams(data.Time)
		option["start_time"] = zeroTime
		option["end_time"] = zeroTime + 86400 - 1
	}

	list, total, err := s.sowTaskDao.GetList(c, data.Page, data.Size, query, option)
	if err != nil {
		return nil, 0, err
	}

	if len(list) == 0 {
		return resp, 0, nil
	}

	// 清空corpExternalIds集合重新赋值
	corpExternalIds = []int64{}
	for _, val := range list {
		corpExternalIds = append(corpExternalIds, val.WorkExternalUserID)
	}

	// 顾客备注
	userRelationMapList := make(map[int64]workExternalModel.WorkExternalUserRelationModel)
	userRelationRes, _ := microScrmDao.NewUserRelation().GetMany(c, "enterprise_id = ? and member_id = ? and work_external_user_id in ?", data.EnterpriseID, memberInfo.ID, corpExternalIds)
	for _, re := range userRelationRes {
		userRelationMapList[re.WorkExternalUserID] = re
	}

	// 顾客信息
	userCorpExternalMapList := make(map[int64]workExternalModel.WxUserCorpExternalModel)
	corpExternalRes, _ := microScrmDao.NewWxUserCorpExternal().GetMany(c, data.EnterpriseID, "id in ?", corpExternalIds)
	for _, re := range corpExternalRes {
		userCorpExternalMapList[re.ID] = re
	}

	for _, item := range list {
		tmpResp := customerResp.SowTaskList{
			Id:                 item.ID,
			WorkExternalUserID: item.WorkExternalUserID,
			MemberID:           item.MemberID,
			MemberUserID:       item.MemberUserID,
			Content:            item.Content,
			Status:             item.Status,
			RemindTime:         item.RemindTime,
			Username:           userCorpExternalMapList[item.WorkExternalUserID].Nick,
			Nickname:           userRelationMapList[item.WorkExternalUserID].Remark,
			Grade:              userCorpExternalMapList[item.WorkExternalUserID].CardLevel,
			Birthday:           userCorpExternalMapList[item.WorkExternalUserID].Birthday,
			OpenId:             userCorpExternalMapList[item.WorkExternalUserID].Openid,
		}
		if tmpResp.Birthday != "" && userCorpExternalMapList[item.WorkExternalUserID].BirthYear != "0" {
			year, _ := helper.DecodeString(userCorpExternalMapList[item.WorkExternalUserID].BirthYear)
			month, _ := helper.DecodeString(userCorpExternalMapList[item.WorkExternalUserID].BirthMonth)
			day, _ := helper.DecodeString(userCorpExternalMapList[item.WorkExternalUserID].BirthDay)
			tmpResp.Birthday = fmt.Sprintf("%s-%s-%s", year, month, day)
		}
		for _, relation := range userRelationMapList {
			if item.MemberID == relation.MemberID && item.WorkExternalUserID == relation.WorkExternalUserID {
				tmpResp.IsKsa = relation.IsKeySa
			}
		}
		resp = append(resp, tmpResp)
	}

	return resp, total, err
}

func (s *SowTask) Create(c app.Ctx, data customerReq.SowTaskCreate) (int64, error) {
	// 查询员工企微id
	membersUser, err := s.GetUserIdByMobile(c, data.EnterpriseID, data.MemberMobile)
	if err != nil {
		return 0, err
	}

	// 获取员工和外部联系人关系
	relation, err := microScrmDao.NewUserRelation().GetUserSARelation(c, data.EnterpriseID, membersUser.ID, data.WorkExternalUserID)
	if err != nil {
		return 0, err
	}
	if relation.Id == 0 {
		return 0, errors.New("该用户非您好友,暂无法设置播种")
	}

	sowTaskModel := customerModels.SowTaskModel{
		MicroBaseModel:     models.CreateMicroBase(data.EnterpriseID, app.UIDContext(c)),
		WorkExternalUserID: data.WorkExternalUserID,
		MemberUserID:       membersUser.Userid,
		MemberID:           membersUser.ID,
		Content:            data.Content,
		RemindTime:         data.RemindTime,
	}
	err = s.sowTaskDao.CreateData(c, sowTaskModel)
	return sowTaskModel.ID, err
}

func (s *SowTask) Update(c app.Ctx, data customerReq.SowTaskUpdate) error {
	err := s.sowTaskDao.UpdateData(c,
		map[string]interface{}{"id": data.Id},
		map[string]interface{}{
			"content":     data.Content,
			"remind_time": data.RemindTime,
		})

	return err
}

func (s *SowTask) Finish(c app.Ctx, data customerReq.SowTaskFinished) error {
	res, err := s.sowTaskDao.GetOne(c, map[string]interface{}{"id": data.Id, "enterprise_id": data.EnterpriseID})
	if err != nil || res.ID == 0 {
		return errors.New("数据不存在")
	}

	err = s.sowTaskDao.UpdateData(c,
		map[string]interface{}{"id": data.Id},
		map[string]interface{}{"status": customerModels.SowTaskStatusFinished})

	return err
}

func (s *SowTask) GetMonthLists(c app.Ctx, req planReq.MonthListsReq) ([]planReq.MonthListsData, error) {
	monthMap := make(map[int64]planReq.MonthListsStatus)
	var monthSplice []planReq.MonthListsData
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: req.EnterpriseId},
		Mobile:         req.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return monthSplice, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	lists, err := customerDao.NewSowTask().GetMonthLists(c, req, memberInfo.Userid)
	if err != nil {
		return monthSplice, err
	}
	for _, v := range lists {
		zeroTime := helper.GetZeroDateByParams(v.RemindTime)
		if item, ok := monthMap[zeroTime]; !ok {
			if v.Status == 2 {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: true, Afoot: false}
			} else {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: false, Afoot: true}
			}
		} else {
			if v.Status != 2 && item.Finished == true {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: false, Afoot: true}
			}
		}
	}

	for k, value := range monthMap {
		info := planReq.MonthListsData{
			Time:             k,
			MonthListsStatus: planReq.MonthListsStatus{Finished: value.Finished, Afoot: value.Afoot},
		}
		monthSplice = append(monthSplice, info)
	}
	return monthSplice, err
}

func (s *SowTask) Delete(c app.Ctx, data customerReq.SowTaskFinished) error {
	membersInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: data.EnterpriseID},
		Mobile:         data.MemberMobile,
	})
	if err != nil || membersInfo.ID == 0 {
		return errors.New("员工信息未查到")
	}

	res, err := s.sowTaskDao.GetOne(c, map[string]interface{}{
		"id":                    data.Id,
		"enterprise_id":         data.EnterpriseID,
		"member_id":             membersInfo.ID,
		"work_external_user_id": data.WorkExternalUserID})
	if err != nil || res.ID == 0 {
		return errors.New("数据不存在")
	}

	if res.Status == customerModels.SowTaskStatusFinished {
		return errors.New("已完成任务,无法删除")
	}

	err = s.sowTaskDao.DeleteData(c, map[string]interface{}{
		"id":            data.Id,
		"enterprise_id": data.EnterpriseID,
	})

	return err
}

func (s *SowTask) AllTaskList(c app.Ctx, params customerReq.SowAllTaskListReq) (respList []customerResp.SowInfoResp, total int64, err error) {
	respList = make([]customerResp.SowInfoResp, 0)
	if params.TemplateName != "" {
		params.TemplateIds, err = customerDao.NewSowTemplate().SearchIds(c, "enterprise_id = ? AND name like ?", params.EnterpriseID, "%"+strings.Trim(params.TemplateName, " ")+"%")
		if err != nil || len(params.TemplateIds) == 0 {
			return
		}
	}
	if params.SaName != "" {
		params.MemberIds, err = membersDao.NewMembers().GetMemberIdsByName(c, params.EnterpriseID, params.SaName)
		if err != nil || len(params.MemberIds) == 0 {
			return
		}
	}
	if params.UnionCode != "" || params.Username != "" {
		params.ExternalIds, err = microScrmDao.NewWxUserCorpExternal().GetIdsBySearch(c, customerReq.SearchUserIdsReq{
			EnterpriseID: params.EnterpriseID,
			Name:         params.Username,
			UnionCode:    params.UnionCode,
		})
		if err != nil || len(params.ExternalIds) == 0 {
			return
		}

	}
	dbList, total, err := s.sowTaskDao.GetBoardList(c, params)
	if err != nil {
		return
	}
	if len(dbList) > 0 {
		var memberIds []int64
		var externalIds []int64
		var templateIds []int64
		for _, val := range dbList {
			memberIds = append(memberIds, val.MemberID)
			externalIds = append(externalIds, val.WorkExternalUserID)
			templateIds = append(templateIds, val.TemplateId)
		}
		memberDepartmentList, _ := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, params.EnterpriseID, memberIds)
		// 整理导购id和导购信息的key=>value
		var memberMap = make(map[int64]customerResp.SowInfoResp)
		for _, info := range memberDepartmentList {
			dpName := ""
			if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
				dpName = info.MemberDepartments[0].Department.Name
			}
			memberMap[info.ID] = customerResp.SowInfoResp{
				SaName:   info.Name,
				SaDpName: dpName,
			}
		}

		externalList, _ := microScrmDao.NewWxUserCorpExternal().GetInfoByIds(c, params.EnterpriseID, externalIds)
		// 整理顾客id和顾客信息的key=>value
		var externalMap = make(map[int64]customerResp.SowInfoResp)
		for _, info := range externalList {
			externalMap[info.ID] = customerResp.SowInfoResp{
				Username:  info.Name,
				UnionCode: info.UnionCode,
			}
		}
		templateList, _ := customerDao.NewSowTemplate().GetUnScopedMany(c, "enterprise_id = ? AND id in ?", params.EnterpriseID, templateIds)
		// 整理顾客id和顾客信息的key=>value
		var templateMap = make(map[int64]string)
		for _, info := range templateList {
			templateMap[info.ID] = info.Name
		}
		for _, val := range dbList {
			finishedAt := int64(0)
			finishedAtStr := "-"
			if val.Status == customerModels.SowTaskStatusFinished {
				finishedAt = val.UpdatedAt
				finishedAtStr = time.Unix(finishedAt, 0).Format("2006-01-02 15:04:05")
			}
			tmpInfo := customerResp.SowInfoResp{
				Id:            val.ID,
				TemplateId:    val.TemplateId,
				TemplateIdStr: helper.RuleThree(val.TemplateId > 0, strconv.Itoa(int(val.TemplateId)), "--"), // 默认自定义模板为--
				TemplateName:  "--",                                                                          // 默认自定义的模板名称为--
				Status:        int(val.Status),
				CreatedAt:     val.CreatedAt,
				RemindTime:    val.RemindTime,
				Content:       val.Content,
				CreatedAtStr:  time.Unix(val.CreatedAt, 0).Format("2006-01-02 15:04:05"),
				RemindTimeStr: time.Unix(val.RemindTime, 0).Format("2006-01-02 15:04:05"),
				FinishedAt:    finishedAt,
				FinishedAtStr: finishedAtStr,
				TemplateType:  helper.RuleThree(val.TemplateId > 0, "播种模版", "默认播种"),
				StatusStr:     helper.RuleThree(val.Status == 2, "已完成", "未完成"),
			}
			if item, ok := templateMap[val.TemplateId]; ok {
				tmpInfo.TemplateName = item
			}
			if item, ok := memberMap[val.MemberID]; ok {
				tmpInfo.SaName = item.SaName
				tmpInfo.SaDpName = item.SaDpName
			}
			if item, ok := externalMap[val.WorkExternalUserID]; ok {
				tmpInfo.Username = item.Username
				tmpInfo.UnionCode = item.UnionCode
			}
			respList = append(respList, tmpInfo)
		}
	}

	return
}

func (s *SowTask) TemplateList(c app.Ctx, enterpriseId int64) ([]customerModels.SowTemplateModel, error) {
	list, err := customerDao.NewSowTemplate().GetTemplateList(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"status":        customerModels.SowTemplateIsUp,
	})
	return list, err
}

func (s *SowTask) NewCreate(c app.Ctx, data customerReq.SowTaskNewCreate) ([]customerResp.SowTaskSimple, error) {
	sowTaskSimple := make([]customerResp.SowTaskSimple, 0)
	// 查询员工企微id
	membersUser, err := s.GetUserIdByMobile(c, data.EnterpriseID, data.MemberMobile)
	if err != nil {
		return sowTaskSimple, err
	}

	// 获取员工和外部联系人关系
	relation, err := microScrmDao.NewUserRelation().GetUserSARelation(c, data.EnterpriseID, membersUser.ID, data.WorkExternalUserID)
	if err != nil {
		return sowTaskSimple, err
	}
	if relation.Id == 0 {
		return sowTaskSimple, errors.New("该用户非您好友,暂无法设置播种")
	}
	if len(data.TaskContent) == 0 {
		return sowTaskSimple, errors.New("播种内容为空")
	}
	sowTaskModels := make([]customerModels.SowTaskModel, 0)

	for _, base := range data.TaskContent {
		id := app.UID()
		sowTaskModel := customerModels.SowTaskModel{
			MicroBaseModel:     models.CreateMicroBase(data.EnterpriseID, id),
			WorkExternalUserID: data.WorkExternalUserID,
			MemberUserID:       membersUser.Userid,
			MemberID:           membersUser.ID,
			Content:            base.Content,
			RemindTime:         base.RemindTime,
			TemplateId:         data.TemplateId,
		}
		sowTaskSimple = append(sowTaskSimple, customerResp.SowTaskSimple{
			Id:      id,
			Content: base.Content,
		})
		sowTaskModels = append(sowTaskModels, sowTaskModel)
	}
	err = s.sowTaskDao.CreateManyData(c, sowTaskModels)

	if data.TemplateId > 0 {
		err = customerDao.NewSowTemplate().Increment(c, map[string]interface{}{
			"id": data.TemplateId,
		}, "use_num", 1)
	}

	return sowTaskSimple, err
}

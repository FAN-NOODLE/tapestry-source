package storeVisitSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	"micro-scrm/internal/app/http/request"
	emailReq "micro-scrm/internal/app/http/request/email"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	"micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	membersModel "micro-scrm/internal/app/models/members"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	emailSvc "micro-scrm/internal/app/service/email"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type MainSvc struct {
	service.Base
	visitMainDao   *storeVisitDao.MainDao
	visitMemberDao *storeVisitDao.MemberRelationDao
}

func NewMainSvc() *MainSvc {
	return &MainSvc{
		visitMainDao:   storeVisitDao.NewMainDao(),
		visitMemberDao: storeVisitDao.NewMemberRelationDao(),
	}
}

// GetStoreManager 根据门店id获取门店管理员
func (s *MainSvc) GetStoreManager(c ctx.Context, req storeVisitReq.StoreManagerReq) ([]membersModel.StoreManager, error) {
	list, err := membersDao.NewMembers().GetStoreVisitManagerByThirdIds(c, req.EnterpriseId, req.ThirdIds)
	for i, manager := range list {
		positionAbbr := ""
		if tmpP, ok := membersDao.StoreManagerPositionMap[manager.Position]; ok {
			positionAbbr = tmpP
		}
		list[i].PositionAbbr = positionAbbr
	}
	return list, err
}

// GetStoreMM 获取门店MM管理员
func (s *MainSvc) GetStoreMM(c ctx.Context, req storeVisitReq.StoreManagerReq) ([]storeVisitResp.StoreMMResp, error) {
	var list = make([]storeVisitResp.StoreMMResp, 0)
	mmMapByThirdIds, err := s.GetMMMapByThirdIds(c, req.EnterpriseId, req.ThirdIds)
	if err != nil {
		return list, err
	}
	for _, thirdId := range req.ThirdIds {
		if mmMemberIds, ok := mmMapByThirdIds[thirdId]; ok {
			list = append(list, storeVisitResp.StoreMMResp{
				ThirdId:     thirdId,
				MMMemberIds: mmMemberIds,
			})
		}
	}
	return list, nil
}

// SavePlan 巡店预约
func (s *MainSvc) SavePlan(c ctx.Context, req storeVisitReq.SavePlanReq) ([]storeVisitResp.StoreVisitInfo, error) {
	var err error
	var visitList = make([]storeVisitResp.StoreVisitInfo, 0)
	nowTime := time.Now().Unix()
	var thirdIds []int
	for _, storeInfo := range req.StoreList {
		thirdIds = append(thirdIds, storeInfo.ThirdId)
	}
	// 获取可以管理这些门店的 MM 用户
	mmMapByThirdIds, err := s.GetMMMapByThirdIds(c, req.EnterpriseId, thirdIds)
	if err != nil {
		return nil, err
	}
	var mainList []storeVisitModel.StoreVisitMainModel
	var memberList []storeVisitModel.StoreVisitMemberRelationModel
	var msgList []request.SendCardMessageToH5
	var emailList []storeVisitResp.SendEmail
	var timeSpanMap = map[string]string{
		"AM": "上午",
		"PM": "下午",
	}
	for _, storeInfo := range req.StoreList {
		info := storeVisitModel.StoreVisitMainModel{
			Date:           storeInfo.Date,
			TimeSpan:       storeInfo.TimeSpan,
			UnixTime:       storeInfo.UnixTime,
			ThirdId:        storeInfo.ThirdId,
			SendMsgToStore: req.SendMsgToStore,
			VisitorIsShow:  storeVisitModel.VisitorDefaultShow,
			StoreIsShow:    storeVisitModel.StoreDefaultShow,
			Remark:         req.Remark,
			Status:         storeVisitModel.VisitStatusPlan,
		}
		if req.Language == "en" {
			//默认都存中文
			if val, ok := timeSpanMap[storeInfo.TimeSpan]; ok {
				info.TimeSpan = val
			}
		}
		id := app.UID()
		info.BaseModel = models.BaseModel{
			ID:        id,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		}
		mainList = append(mainList, info)
		//该巡店计划相关的人
		var relMemberIds []int64
		//需要企微通知的人
		var msgMemberIds []int64
		//需要邮件通知的人
		var emails []string
		//创建者
		memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
			BaseModel:    models.CreateBase(app.UID()),
			VisitID:      id,
			MemberID:     req.MemberId,
			Type:         storeVisitModel.VisitCreator,
			ReportIsShow: storeVisitModel.VisitReportIsShow,
		})
		relMemberIds = append(relMemberIds, req.MemberId)
		//店长
		for _, managerId := range storeInfo.StoreManager {
			memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
				BaseModel: models.CreateBase(app.UID()),
				VisitID:   id,
				MemberID:  managerId,
				Type:      storeVisitModel.VisitStoreManager,
			})
			msgMemberIds = append(msgMemberIds, managerId)
			relMemberIds = append(relMemberIds, managerId)
		}
		//参与者
		//创建人也是参与人
		memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
			BaseModel: models.CreateBase(app.UID()),
			VisitID:   id,
			MemberID:  req.MemberId,
			Type:      storeVisitModel.VisitVisitor,
		})
		//管理门店的MM也是参与人
		if mmMemberIds, ok := mmMapByThirdIds[storeInfo.ThirdId]; ok {
			if len(mmMemberIds) > 0 {
				for _, mmMemberId := range mmMemberIds {
					if !helper.InArray(mmMemberId, relMemberIds) {
						memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
							BaseModel: models.CreateBase(app.UID()),
							VisitID:   id,
							MemberID:  mmMemberId,
							Type:      storeVisitModel.VisitVisitor,
							IsMM:      storeVisitModel.IsMM,
						})
					}
					msgMemberIds = append(msgMemberIds, mmMemberId)
				}
				msgMemberIds = helper.ArrayUniq(msgMemberIds)
			}
		}
		for _, visitor := range req.VisitorList {
			if !helper.InArray(visitor.MemberId, relMemberIds) {
				memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
					BaseModel: models.CreateBase(app.UID()),
					VisitID:   id,
					MemberID:  visitor.MemberId,
					Type:      storeVisitModel.VisitVisitor,
					Email:     visitor.Email,
				})
				if visitor.MemberId > 0 {
					msgMemberIds = append(msgMemberIds, visitor.MemberId)
				} else if visitor.Email != "" {
					emails = append(emails, visitor.Email)
				}
			}
		}
		if len(msgMemberIds) > 0 {
			userIds, _ := membersDao.NewMembers().GetUserIdsById(c, req.EnterpriseId, msgMemberIds)
			if len(userIds) > 0 {
				msgList = append(msgList, request.SendCardMessageToH5{
					EnterpriseBaseReq: request.EnterpriseBaseReq{
						EnterpriseId: req.EnterpriseId,
					},
					Title:   helper.RuleThree(req.Language == "en", "Store Visit Appointment", "巡店预约提醒"),
					Uri:     fmt.Sprintf("store-visit-report?id=%d", id),
					Party:   nil,
					UserIds: userIds,
					Content: []map[string]string{
						{
							"keyname": helper.RuleThree(req.Language == "en", "Time", "巡店时间"),
							"value":   fmt.Sprintf("%s %s", storeInfo.Date, helper.FormatStoreVisitTimeSpan(req.Language, storeInfo.TimeSpan)),
						},
						{
							"keyname": helper.RuleThree(req.Language == "en", "Store", "预约门店"),
							"value":   storeInfo.StoreName,
						},
						{
							"keyname": helper.RuleThree(req.Language == "en", "Description", "巡店内容"),
							"value":   helper.TruncateString(req.Remark, 25),
						},
					},
					AppSign:  qyWork.StoreVisit,
					Language: req.Language,
				})
			}
		}
		if len(emails) > 0 {
			emailList = append(emailList, storeVisitResp.SendEmail{
				Emails:    emails,
				StoreName: storeInfo.StoreName,
				TimeStr:   fmt.Sprintf("%s %s", storeInfo.Date, helper.FormatStoreVisitTimeSpan(req.Language, storeInfo.TimeSpan)),
			})
		}
		visitList = append(visitList, storeVisitResp.StoreVisitInfo{
			ID:      id,
			ThirdId: info.ThirdId,
		})
	}
	if len(mainList) == 0 || len(memberList) == 0 {
		return nil, errors.New("巡店信息不能为空")
	}
	//写入巡店预约
	err = s.visitMainDao.CreateInBatchesData(c, req.EnterpriseId, mainList)
	if err != nil {
		return nil, err
	}
	//写入预约相关人
	err = s.visitMemberDao.CreateInBatchesData(c, req.EnterpriseId, memberList)
	if err != nil {
		return nil, err
	}

	//发送巡店通知
	for _, msg := range msgList {
		msgParams := msg
		go func() {
			_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, msgParams)
			if err != nil {
				app.Logger().Error(c).Msgf("巡店预约通知消息发送失败：%+v", err)
			}
		}()
	}
	//发送邮件通知
	go func() {
		for _, emailInfo := range emailList {
			err = emailSvc.NewSendEmail().SendEmailByDiyGoFunc(emailReq.SendMailReq{
				Title: fmt.Sprintf(helper.RuleThree(req.Language == "en", "Store Visit Appointment at", "巡店预约提醒")+"【%s】", emailInfo.StoreName),
				Content: fmt.Sprintf(helper.RuleThree(req.Language == "en",
					"<p>Hi，<p><p></p><p>You have received a store visit appointment reminder. Please check the details in the StoreVisit on WeCom app:</p><p>Date & Time:%s</p><p>Store:%s</p>",
					"<p>您好，<p><p></p><p>您收到一条巡店预约提醒，您可以在企微StoreVisit应用中查看详细信息；</p><p>预约时间：%s</p><p>预约门店： %s</p>"),
					emailInfo.TimeStr, emailInfo.StoreName),
				ToAddress: emailInfo.Emails,
			}, config.VchangyiEmailConfig{
				Host:     "smtp.exmail.qq.com",
				Port:     465,
				Username: "store_visit@vchangyi.com",
				Password: "REDACTED_FOR_SOURCE_SHARING",
			})
		}
	}()

	return visitList, err
}

// GetMMMapByThirdIds 获取可以管理这些门店的MM用户
// 返回 map[int][]int64，key 是 thirdId，value 是能管理该门店的 MM 的 member_id 列表
func (s *MainSvc) GetMMMapByThirdIds(c ctx.Context, enterpriseId int64, thirdIds []int) (map[int][]int64, error) {
	// 1. 根据 thirdIds 获取对应的门店 area 记录
	areaList, err := dashboardDao.NewArea().GetMany(c, enterpriseId, "type = ? AND third_id in ?", dashboardModel.AreaStore, thirdIds)
	if err != nil {
		return nil, err
	}

	// 建立 thirdId -> areaId 的映射
	thirdIdToAreaId := make(map[int]int64)
	for _, area := range areaList {
		thirdIdToAreaId[area.ThirdID] = area.ID
	}

	// 2. 获取所有 MM/CL 级别的管理者
	mmList, err := dashboardDao.NewManager().GetMany(c, enterpriseId, "position = ?", "MM/CL(Multi Manager/City Manager)")
	if err != nil {
		return nil, err
	}

	// 3. 建立 areaId -> thirdId 的反向映射
	areaIdToThirdId := make(map[int64]int)
	for thirdId, areaId := range thirdIdToAreaId {
		areaIdToThirdId[areaId] = thirdId
	}

	// 4. 初始化结果 map
	result := make(map[int][]int64)

	// 5. 遍历每个 MM，检查其管理范围内是否包含这些门店
	for _, mm := range mmList {
		// 解析 MM 管理的 area_ids
		mmAreaIds := helper.StringToIntSlice(mm.AreaIds, ",")
		if len(mmAreaIds) == 0 {
			continue
		}

		// 检查门店的 area_id 是否在 MM 的管理范围内
		// 判断逻辑：MM 管理的 area_id 在门店的 path 中即可管理该门店
		for _, area := range areaList {
			// 拆分门店的 path，获取所有父级 area_id
			pathIds := strings.Split(area.Path, ",")
			isManaged := false
			for _, mmAreaId := range mmAreaIds {
				// 情况1：门店直接等于管理的 area
				if area.ID == int64(mmAreaId) {
					isManaged = true
					break
				}
				// 情况2：门店的 path 中包含该 area_id
				for _, pathId := range pathIds {
					if pathId == strconv.Itoa(mmAreaId) {
						isManaged = true
						break
					}
				}
				if isManaged {
					break
				}
			}

			if isManaged {
				thirdId := area.ThirdID
				result[thirdId] = append(result[thirdId], mm.MemberId)
			}
		}
	}

	return result, nil
}

func (s *MainSvc) StoreVisitList(c ctx.Context, req storeVisitReq.PlanListReq) ([]storeVisitResp.StoreVisitInfo, int64, error) {
	var result = make([]storeVisitResp.StoreVisitInfo, 0)
	if req.DataType == storeVisitReq.DataTypeAll {
		managerInfo, _ := storeVisitDao.NewDeptMemberDao().GetOne(c, req.EnterpriseId, "member_id = ? AND is_store_visit_manager = 1", req.MemberId)
		if managerInfo.ID == 0 {
			//不是管理员 没有获取所有数据的权限
			return result, 0, nil
		}
	}
	visitList, total, err := s.visitMainDao.GetStoreVisitList(c, req)
	if err != nil || len(visitList) == 0 {
		return result, 0, err
	}
	var visitIds []int64
	var thirdIds []int
	for _, visit := range visitList {
		visitIds = append(visitIds, visit.ID)
		thirdIds = append(thirdIds, visit.ThirdId)
	}
	storeList, err := microWorkDao.NewStoreDepartment().GetStoreName(c, req.EnterpriseId, thirdIds)
	if err != nil {
		return nil, 0, err
	}
	var storeMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, store := range storeList {
		storeMap[int(store.WechatDepartmentID)] = store
	}
	visitMemberList, err := s.visitMemberDao.GetMany(c, req.EnterpriseId, "visit_id in ?", visitIds)
	if err != nil {
		return nil, 0, err
	}
	var memberIds []int64
	var mmMap = make(map[int64]int)
	for _, row := range visitMemberList {
		memberIds = append(memberIds, row.MemberID)
		mmMap[row.MemberID] = row.IsMM
	}
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, req.EnterpriseId, helper.ArrayUniq(memberIds))
	if err != nil {
		return nil, 0, err
	}
	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]storeVisitResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		positionAbbr := ""
		if tmpP, ok := membersDao.StoreManagerPositionMap[info.Position]; ok {
			positionAbbr = tmpP
		}
		memberMap[info.ID] = storeVisitResp.MemberInfo{
			MemberId:     info.ID,
			Name:         info.Name,
			Email:        info.Email,
			DeptName:     dpName,
			UserId:       info.Userid,
			PositionAbbr: positionAbbr,
		}

	}
	storeDeptMemberList, _ := storeVisitDao.NewDeptMemberDao().GetMany(c, req.EnterpriseId, "member_id in ? ", helper.ArrayUniq(memberIds))
	var storeDeptMemberMap = make(map[int64]string)
	for _, deptMember := range storeDeptMemberList {
		storeDeptMemberMap[deptMember.MemberID] = deptMember.DeptName
	}
	for _, row := range visitList {
		tmp := storeVisitResp.StoreVisitInfo{
			ID:            row.ID,
			Date:          row.Date,
			TimeSpan:      helper.FormatStoreVisitTimeSpan(req.Language, row.TimeSpan),
			ThirdId:       row.ThirdId,
			Remark:        row.Remark,
			Status:        row.Status,
			StoreManager:  []storeVisitResp.MemberInfo{},
			VisitorMember: []storeVisitResp.MemberInfo{},
			CopyMember:    []storeVisitResp.MemberInfo{},
			UpdatedAt:     row.UpdatedAt,
		}
		for _, visitMember := range visitMemberList {
			if row.ID != visitMember.VisitID {
				continue
			}
			var memberInfo storeVisitResp.MemberInfo
			if visitMember.MemberID > 0 {
				if member, exist := memberMap[visitMember.MemberID]; exist {
					memberInfo = member
					if isMM, ok := mmMap[visitMember.MemberID]; ok {
						memberInfo.IsMM = isMM
					}
				} else {
					continue
				}
			} else if visitMember.Email != "" {
				memberInfo = storeVisitResp.MemberInfo{Email: visitMember.Email}
			} else {
				continue
			}
			if deptMemberName, ok := storeDeptMemberMap[visitMember.MemberID]; ok {
				memberInfo.DeptName = deptMemberName
			}
			switch visitMember.Type {
			case storeVisitModel.VisitCreator:
				tmp.Creator = memberInfo
			case storeVisitModel.VisitStoreManager:
				tmp.StoreManager = append(tmp.StoreManager, memberInfo)
			case storeVisitModel.VisitVisitor:
				tmp.VisitorMember = append(tmp.VisitorMember, memberInfo)
			case storeVisitModel.VisitCopyMember:
				tmp.CopyMember = append(tmp.CopyMember, memberInfo)
			}
		}
		//门店信息
		if store, ok := storeMap[row.ThirdId]; ok {
			tmp.StoreCode = store.ThirdStoreCode
			tmp.StoreName = store.StoreName
			if req.Language == "en" && store.InventoryName != "" {
				tmp.StoreName = store.InventoryName
			}
		}
		result = append(result, tmp)
	}

	return result, total, nil
}

func (s *MainSvc) StoreVisitInfo(c ctx.Context, req storeVisitReq.VisitInfoReq) (storeVisitResp.StoreVisitInfo, error) {
	result := storeVisitResp.StoreVisitInfo{
		StoreManager:  []storeVisitResp.MemberInfo{},
		VisitorMember: []storeVisitResp.MemberInfo{},
		CopyMember:    []storeVisitResp.MemberInfo{},
	}
	visitInfo, err := s.visitMainDao.GetOne(c, req.EnterpriseId, "id = ?", req.ID)
	if err != nil {
		return result, err
	}
	helper.CoverStruct(&result, &visitInfo)
	result.TimeSpan = helper.FormatStoreVisitTimeSpan(req.Language, result.TimeSpan)
	result.ID = visitInfo.ID
	result.UpdatedAt = visitInfo.UpdatedAt
	memberList, err := s.visitMemberDao.GetMany(c, req.EnterpriseId, "visit_id = ?", req.ID)
	if err != nil {
		return storeVisitResp.StoreVisitInfo{}, err
	}
	var memberIds []int64
	for _, row := range memberList {
		memberIds = append(memberIds, row.MemberID)
	}
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, req.EnterpriseId, memberIds)
	if err != nil {
		return result, err
	}
	storeDeptMemberList, _ := storeVisitDao.NewDeptMemberDao().GetMany(c, req.EnterpriseId, "member_id in ? ", helper.ArrayUniq(memberIds))
	var storeDeptMemberMap = make(map[int64]string)
	for _, deptMember := range storeDeptMemberList {
		storeDeptMemberMap[deptMember.MemberID] = deptMember.DeptName
	}
	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]storeVisitResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := storeDeptMemberMap[info.ID]
		if dpName == "" && len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		positionAbbr := ""
		if tmpP, ok := membersDao.StoreManagerPositionMap[info.Position]; ok {
			positionAbbr = tmpP
		}
		memberMap[info.ID] = storeVisitResp.MemberInfo{
			MemberId:     info.ID,
			Name:         info.Name,
			Email:        info.Email,
			DeptName:     dpName,
			UserId:       info.Userid,
			PositionAbbr: positionAbbr,
		}
	}

	for _, row := range memberList {
		var tmpMember storeVisitResp.MemberInfo
		if row.MemberID > 0 {
			if member, exist := memberMap[row.MemberID]; exist {
				tmpMember = member
				tmpMember.IsMM = row.IsMM
			}
		} else {
			tmpMember.Email = row.Email
		}
		// 添加成员信息
		switch row.Type {
		case storeVisitModel.VisitCreator:
			result.Creator = tmpMember
		case storeVisitModel.VisitStoreManager:
			result.StoreManager = append(result.StoreManager, tmpMember)
		case storeVisitModel.VisitVisitor:
			result.VisitorMember = append(result.VisitorMember, tmpMember)
		case storeVisitModel.VisitCopyMember:
			result.CopyMember = append(result.CopyMember, tmpMember)
		}
	}
	storeList, _ := microWorkDao.NewStoreDepartment().GetStoreName(c, req.EnterpriseId, []int{result.ThirdId})
	if len(storeList) > 0 {
		result.StoreName = storeList[0].StoreName
		if req.Language == "en" && storeList[0].InventoryName != "" {
			result.StoreName = storeList[0].InventoryName
		}
		result.StoreCode = storeList[0].ThirdStoreCode
	}
	sort.Slice(result.VisitorMember, func(i, j int) bool {
		return result.VisitorMember[i].IsMM > result.VisitorMember[j].IsMM
	})
	return result, nil
}

// AddCopyMembers   同步给指定人员
func (s *MainSvc) AddCopyMembers(c ctx.Context, req storeVisitReq.AddCopyMemberReq) error {
	visitMemberList, err := s.visitMemberDao.GetMany(c, req.EnterpriseId, "visit_id = ?", req.ID)
	if err != nil {
		return err
	}
	var existMemberIds []int64
	var existEmails []string
	var hasPermission bool
	for _, model := range visitMemberList {
		if model.MemberID > 0 {
			existMemberIds = append(existMemberIds, model.MemberID)
		} else if model.Email != "" {
			existEmails = append(existEmails, model.Email)
		}
		if model.Type == storeVisitModel.VisitStoreManager && req.MemberId == model.MemberID {
			hasPermission = true
		}
	}
	if !hasPermission {
		return errors.New("您不是预约时的门店店长，不能进行该操作")
	}
	var addMemberIds []int64
	var addEmails []string
	for _, visitor := range req.CopyMemberList {
		if visitor.MemberId > 0 && !helper.InArray(visitor.MemberId, existMemberIds) {
			addMemberIds = append(addMemberIds, visitor.MemberId)
		} else if visitor.Email != "" && !helper.InArray(visitor.Email, existEmails) {
			addEmails = append(addEmails, visitor.Email)
		}
	}
	var addSlice []storeVisitModel.StoreVisitMemberRelationModel
	for _, id := range addMemberIds {
		addSlice = append(addSlice, storeVisitModel.StoreVisitMemberRelationModel{
			BaseModel:    models.CreateBase(app.UID()),
			VisitID:      req.ID,
			MemberID:     id,
			Type:         storeVisitModel.VisitCopyMember,
			ReportIsShow: storeVisitModel.VisitReportIsShow,
		})
	}
	for _, email := range addEmails {
		addSlice = append(addSlice, storeVisitModel.StoreVisitMemberRelationModel{
			BaseModel:    models.CreateBase(app.UID()),
			VisitID:      req.ID,
			Email:        email,
			Type:         storeVisitModel.VisitCopyMember,
			ReportIsShow: storeVisitModel.VisitReportIsShow,
		})
	}
	return s.visitMemberDao.CreateInBatchesData(c, req.EnterpriseId, addSlice)
}

func (s *MainSvc) StartVisit(c ctx.Context, req storeVisitReq.StartVisitReq) (int64, error) {
	var err error
	if req.ID > 0 {
		//预约巡店开始巡店
		err = s.CheckCreator(c, req.VisitInfoReq)
		if err != nil {
			return 0, err
		}
		saveMap := map[string]interface{}{
			"status":     storeVisitModel.VisitVisiting,
			"updated_at": time.Now().Unix(),
		}
		if req.ThirdId > 0 {
			saveMap["third_id"] = req.ThirdId
		}
		err = s.visitMainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), saveMap)
	} else if req.ThirdId > 0 {
		// 获取可以管理这些门店的 MM 用户
		mmMapByThirdIds, _ := s.GetMMMapByThirdIds(c, req.EnterpriseId, []int{req.ThirdId})
		// 直接开始巡店
		req.ID = app.UID()
		var timeSpanMap = map[string]string{
			"AM": "上午",
			"PM": "下午",
		}
		timeSpan := req.TimeSpan
		if req.Language == "en" {
			//默认都存中文
			if val, ok := timeSpanMap[req.TimeSpan]; ok {
				timeSpan = val
			}
		}
		err = s.visitMainDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitMainModel{
			BaseModel:     models.CreateBase(req.ID),
			Date:          req.Date,
			TimeSpan:      timeSpan,
			VisitorIsShow: storeVisitModel.VisitorDefaultShow,
			UnixTime:      req.UnixTime,
			ThirdId:       req.ThirdId,
			Status:        storeVisitModel.VisitVisiting,
			StoreIsShow:   storeVisitModel.StoreDefaultShow,
		})

		if err == nil {
			var memberList []storeVisitModel.StoreVisitMemberRelationModel
			managerList, _ := membersDao.NewMembers().GetStoreVisitManagerByThirdIds(c, req.EnterpriseId, []int{req.ThirdId})
			if len(managerList) > 0 {
				for _, manager := range managerList {
					memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
						BaseModel: models.CreateBase(app.UID()),
						VisitID:   req.ID,
						MemberID:  manager.MemberId,
						Type:      storeVisitModel.VisitStoreManager,
					})
				}
			}
			memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
				BaseModel: models.CreateBase(app.UID()),
				VisitID:   req.ID,
				MemberID:  req.MemberId,
				Type:      storeVisitModel.VisitCreator,
			})
			memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
				BaseModel: models.CreateBase(app.UID()),
				VisitID:   req.ID,
				MemberID:  req.MemberId,
				Type:      storeVisitModel.VisitVisitor,
			})
			//添加MM作为巡店参与者
			if mmMemberIds, ok := mmMapByThirdIds[req.ThirdId]; ok {
				for _, memberId := range mmMemberIds {
					memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
						BaseModel: models.CreateBase(app.UID()),
						VisitID:   req.ID,
						MemberID:  memberId,
						Type:      storeVisitModel.VisitVisitor,
						IsMM:      storeVisitModel.IsMM,
					})
				}
			}
			err = s.visitMemberDao.CreateInBatchesData(c, req.EnterpriseId, memberList)
		}
	}
	return req.ID, err
}

func (s *MainSvc) VisitCheckIn(c ctx.Context, req storeVisitReq.StartVisitReq) error {
	err := s.CheckCreator(c, req.VisitInfoReq)
	if err != nil {
		return err
	}
	var now = time.Now().Unix()
	return s.visitMainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
		"to_store_time": now,
		"status":        storeVisitModel.VisitVisiting,
		"updated_at":    now,
	})
}

func (s *MainSvc) SendCompleteReport(c ctx.Context, req storeVisitReq.VisitInfoReq) error {
	//提交巡店记录
	err := s.CheckCreator(c, req)
	if err != nil {
		return err
	}
	visitInfo, _ := s.visitMainDao.GetOne(c, req.EnterpriseId, "id = ?", req.ID)
	//需要企微通知的人
	var msgMemberIds []int64
	//需要邮件通知的人
	var emails []string
	memberList, _ := s.visitMemberDao.GetMany(c, req.EnterpriseId, "visit_id = ?", req.ID)
	for _, member := range memberList {
		if member.Type == storeVisitModel.VisitCreator {
			//给创建者不发送通知
			continue
		}
		if member.MemberID > 0 {
			msgMemberIds = append(msgMemberIds, member.MemberID)
		} else if member.Email != "" {
			emails = append(emails, member.Email)
		}
	}
	storeInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": visitInfo.ThirdId,
	})
	if req.Language == "en" && storeInfo.InventoryName != "" {
		storeInfo.StoreName = storeInfo.InventoryName
	}
	err = s.SendCompleteMsg(c, req.BaseReq, msgMemberIds, emails, visitInfo, storeInfo)
	if err != nil {
		return err
	}
	err = s.visitMainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
		"is_report": storeVisitModel.VisitSendReport,
	})
	return err
}

func (s *MainSvc) CompleteVisit(c ctx.Context, req storeVisitReq.VisitInfoReq) error {
	//提交巡店记录
	err := s.CheckCreator(c, req)
	if err != nil {
		return err
	}
	err = s.visitMainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
		"status":     storeVisitModel.VisitComplete,
		"updated_at": time.Now().Unix(),
	})
	return err
}

func (s *MainSvc) SaveVisitInfo(c ctx.Context, req storeVisitReq.SaveVisitInfoReq) error {
	//更新巡店记录内容
	err := s.CheckCreator(c, req.VisitInfoReq)
	if err != nil {
		return err
	}
	updateDate := map[string]interface{}{}
	if req.ResetOverview > 0 {
		updateDate["overview"] = ""
	} else if req.ResetContent > 0 {
		updateDate["content"] = req.Content
	} else {
		if req.Overview != "" {
			updateDate["overview"] = req.Overview
		}
		if req.Content != "" {
			updateDate["content"] = req.Content
		}
	}
	if len(updateDate) > 0 {
		err = s.visitMainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), updateDate)
	}
	return err
}

func (s *MainSvc) SaveReportPermission(c ctx.Context, req storeVisitReq.SaveReportPermissionReq) error {
	err := s.CheckCreator(c, req.VisitInfoReq)
	if err != nil {
		return err
	}
	nowTime := time.Now().Unix()
	if len(req.DelVisitorList) > 0 {
		var delMemberIds []int64
		var delEmails []string
		for _, del := range req.DelVisitorList {
			if del.MemberId > 0 {
				delMemberIds = append(delMemberIds, del.MemberId)
			}
			if del.Email != "" {
				delEmails = append(delEmails, del.Email)
			}
		}
		if len(delMemberIds) > 0 {
			err = s.visitMemberDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
				"member_id": delMemberIds,
				"type":      []int{storeVisitModel.VisitVisitor, storeVisitModel.VisitStoreManager},
			}, map[string]interface{}{
				"deleted_at": nowTime,
			})
		}
		if len(delEmails) > 0 {
			err = s.visitMemberDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
				"email": delEmails,
			}, map[string]interface{}{
				"deleted_at": nowTime,
			})
		}
	}
	if len(req.AddVisitorList) > 0 || len(req.AddStoreManager) > 0 {
		// 获取可以管理这些门店的 MM 用户
		visitInfo, err := s.visitMainDao.GetOne(c, req.EnterpriseId, "id = ?", req.ID)
		if err != nil {
			return err
		}
		mmMapByThirdIds, err := s.GetMMMapByThirdIds(c, req.EnterpriseId, []int{visitInfo.ThirdId})
		storeInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
			"enterprise_id":        req.EnterpriseId,
			"wechat_department_id": visitInfo.ThirdId,
		})
		if req.Language == "en" && storeInfo.InventoryName != "" {
			storeInfo.StoreName = storeInfo.InventoryName
		}
		//需要企微通知的人
		var msgMemberIds []int64
		//需要邮件通知的人
		var emails []string
		var memberList []storeVisitModel.StoreVisitMemberRelationModel
		if len(req.AddVisitorList) > 0 {
			for _, visitor := range req.AddVisitorList {
				addTmp := storeVisitModel.StoreVisitMemberRelationModel{
					BaseModel: models.BaseModel{
						ID:        app.UID(),
						CreatedAt: nowTime,
						UpdatedAt: nowTime,
					},
					VisitID:  req.ID,
					MemberID: visitor.MemberId,
					Type:     storeVisitModel.VisitVisitor,
					Email:    visitor.Email,
				}
				if mmMemberIds, ok := mmMapByThirdIds[visitInfo.ThirdId]; ok {
					if visitor.MemberId > 0 && helper.InArray(visitor.MemberId, mmMemberIds) {
						addTmp.IsMM = storeVisitModel.IsMM
					}
				}
				memberList = append(memberList, addTmp)
				if visitor.MemberId > 0 {
					msgMemberIds = append(msgMemberIds, visitor.MemberId)
				} else if visitor.Email != "" {
					emails = append(emails, visitor.Email)
				}
			}
		}
		if len(req.AddStoreManager) > 0 {
			for _, memberId := range req.AddStoreManager {
				memberList = append(memberList, storeVisitModel.StoreVisitMemberRelationModel{
					BaseModel: models.BaseModel{
						ID:        app.UID(),
						CreatedAt: nowTime,
						UpdatedAt: nowTime,
					},
					VisitID:  req.ID,
					MemberID: memberId,
					Type:     storeVisitModel.VisitStoreManager,
				})
				msgMemberIds = append(msgMemberIds, memberId)
			}
		}

		err = s.visitMemberDao.CreateInBatchesData(c, req.EnterpriseId, memberList)
		if err != nil {
			return err
		}
		if visitInfo.IsReport == storeVisitModel.VisitSendReport {
			err = s.SendCompleteMsg(c, req.BaseReq, msgMemberIds, emails, visitInfo, storeInfo)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (s *MainSvc) SendCompleteMsg(c ctx.Context, req storeVisitReq.BaseReq, msgMemberIds []int64, emails []string, visitInfo storeVisitModel.StoreVisitMainModel, storeInfo microWorkModel.StoreDepartmentModel) error {
	if len(msgMemberIds) > 0 {
		userIds, _ := membersDao.NewMembers().GetUserIdsById(c, req.EnterpriseId, msgMemberIds)
		if len(userIds) > 0 {
			go func() {
				_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
					EnterpriseBaseReq: request.EnterpriseBaseReq{
						EnterpriseId: req.EnterpriseId,
					},
					Title:   helper.RuleThree(req.Language == "en", "Store Visit Completion", "巡店完成通知"),
					Uri:     fmt.Sprintf("store-visit-report?id=%d", visitInfo.ID),
					Party:   nil,
					UserIds: userIds,
					Content: []map[string]string{
						{
							"keyname": helper.RuleThree(req.Language == "en", "Time", "巡店时间"),
							"value":   fmt.Sprintf("%s %s", visitInfo.Date, helper.FormatStoreVisitTimeSpan(req.Language, visitInfo.TimeSpan)),
						},
						{
							"keyname": helper.RuleThree(req.Language == "en", "Store", "巡店门店"),
							"value":   helper.RuleThree(req.Language == "en" && storeInfo.InventoryName != "", storeInfo.InventoryName, storeInfo.StoreName),
						},
					},
					AppSign:  qyWork.StoreVisit,
					Language: req.Language,
				})
				if err != nil {
					app.Logger().Error(c).Msgf("巡店完成通知消息发送失败：%+v", err)
				}
			}()
		}
	}
	if len(emails) > 0 {
		mailReq := emailReq.SendMailReq{
			Title: fmt.Sprintf(helper.RuleThree(req.Language == "en", "Visit Completed", "巡店完成通知")+"【%s】", storeInfo.StoreName),
			Content: fmt.Sprintf(helper.RuleThree(req.Language == "en",
				"<p>Hi,<p><p></p><p>The store visit at 【%s】has been completed. You can review the report via the StoreVisit on WeCom app.</p><p>Visit Time:%s %s</p><p>Store: %s</p> ",
				"<p>您好，<p><p></p><p>【%s】巡店已完成，您可以在企微StoreVisit应用中查看巡店报告；</p><p>巡店时间：%s %s</p><p>巡店门店： %s</p> "),
				storeInfo.StoreName, visitInfo.Date, helper.FormatStoreVisitTimeSpan(req.Language, visitInfo.TimeSpan), storeInfo.StoreName),
			ToAddress: emails,
		}
		if visitInfo.PdfUrl != "" {
			mailReq.FileList = append(mailReq.FileList, emailReq.FileInfo{
				Url:  visitInfo.PdfUrl,
				Name: fmt.Sprintf("%s_StoreVisitReport_%s.pdf", strings.Replace(storeInfo.StoreName, "/", "_", -1), time.Now().Format("20060102")),
			})
		}
		go func() {
			_ = emailSvc.NewSendEmail().SendEmailByDiyGoFunc(mailReq, config.VchangyiEmailConfig{
				Host:     "smtp.exmail.qq.com",
				Port:     465,
				Username: "store_visit@vchangyi.com",
				Password: "REDACTED_FOR_SOURCE_SHARING",
			})
		}()

	}
	return nil
}
func (s *MainSvc) DashboardTree(c ctx.Context, req storeVisitReq.DashboardTreeReq) ([]*storeVisitResp.DashboardTree, error) {
	//获取所有门店数据
	allStore, err := microWorkDao.NewStoreDepartment().GetAll(c, req.EnterpriseId)
	for i, storeInfo := range allStore {
		if req.Language == "en" && storeInfo.InventoryName != "" {
			allStore[i].StoreName = storeInfo.InventoryName
		}
	}
	if err != nil {
		return nil, err
	}
	//整理渠道 - 省 - 市 数据
	tree := buildAdministrativeDivisionTree(allStore)

	//获取渠道下的大区
	areaList, err := dashboardDao.NewArea().GetAreaListByTypeAsc(c, req.EnterpriseId, "is_virtual=0")
	if err != nil {
		return nil, err
	}
	var storeReginMap = make(map[int]string)
	var areaMap = make(map[int64]string)
	var sourceMap = make(map[int][]dashboardModel.DashboardAreaModel)

	for _, area := range areaList {
		if area.Type == dashboardModel.AreaStore {
			// 组装节点的path
			var tmpPathSlice []int64
			for _, tmp := range strings.Split(area.Path, ",") {
				tmpPathSlice = append(tmpPathSlice, int64(helper.StringToInt(tmp)))
			}
			if len(tmpPathSlice) != 3 {
				continue
			}
			storeReginMap[area.ThirdID] = areaMap[tmpPathSlice[1]]
		} else if area.Type == dashboardModel.AreaRegion {
			var tmp []dashboardModel.DashboardAreaModel
			if _, ok := sourceMap[area.ParentID]; ok {
				tmp = sourceMap[area.ParentID]
			}
			tmp = append(tmp, area)
			sourceMap[area.ParentID] = tmp
		}
		areaMap[area.ID] = area.Name
	}
	//获取所有巡店数据
	var dataListMap = make(map[int]int64, 0)
	req.VisitIds, _ = s.visitMemberDao.GetVisitIdsByDataType(c, req)
	if len(req.VisitIds) > 0 {
		dataList, err := s.visitMainDao.TotalCountByStore(c, req)
		if err != nil {
			return nil, err
		}
		for _, data := range dataList {
			dataListMap[data.ThirdId] = data.TotalCount
		}
	}

	for _, divisionTree := range tree {
		// 渠道层展示区域
		for _, area := range areaList {
			if area.Type == dashboardModel.AreaSource && strings.Contains(strings.ToLower(area.Name), strings.ToLower(divisionTree.Name)) {
				//匹配到渠道
				if regionList, ok := sourceMap[int(area.ID)]; ok {
					//这是渠道下的大区 regionList
					for _, region := range regionList {
						tmpReginData := storeVisitResp.DashboardData{Name: region.Name}
						for thirdId, regionName := range storeReginMap {
							if regionName == region.Name {
								if totalCount, ok := dataListMap[thirdId]; ok {
									tmpReginData.TotalCount += totalCount
								}
							}
						}
						divisionTree.DataList = append(divisionTree.DataList, tmpReginData)
					}
				}
			}
		}
		for _, province := range divisionTree.Children {
			for _, thirdId := range province.ThirdIds {
				if totalCount, ok := dataListMap[thirdId]; ok {
					province.TotalCount += totalCount
				}
			}
			for _, city := range province.Children {
				for _, thirdId := range city.ThirdIds {
					if totalCount, ok := dataListMap[thirdId]; ok {
						city.TotalCount += totalCount
						for _, store := range city.Children {
							if store.ThirdId == thirdId {
								store.TotalCount = totalCount
							}
						}
					}

				}
			}
		}
		sortProvinceNodes(divisionTree.Children)
	}
	return tree, nil
}

// 先定义排序函数
func sortProvinceNodes(nodes []*storeVisitResp.ProvinceNode) {
	// 定义标准省级行政区划顺序
	var provinceOrderByPinyin = map[string]int{
		"安徽省":      1,  // Ān huī (A)
		"北京市":      2,  // Běi jīng (B)
		"重庆市":      3,  // Chóng qìng (C)
		"福建省":      4,  // Fú jiàn (F)
		"甘肃省":      5,  // Gān sù (G)
		"广东省":      6,  // Guǎng dōng (G)
		"广西壮族自治区":  7,  // Guǎng xī (G) - 注意"广"是多音字(guǎng/ān)，此处读guǎng
		"贵州省":      8,  // Guì zhōu (G)
		"海南省":      9,  // Hǎi nán (H)
		"河北省":      10, // Hé běi (H)
		"河南省":      11, // Hé nán (H)
		"黑龙江省":     12, // Hēi lóng jiāng (H)
		"湖北省":      13, // Hú běi (H)
		"湖南省":      14, // Hú nán (H)
		"吉林省":      15, // Jí lín (J)
		"江苏省":      16, // Jiāng sū (J)
		"江西省":      17, // Jiāng xī (J)
		"辽宁省":      18, // Liáo níng (L)
		"内蒙古自治区":   19, // Nèi měng gǔ (N) - "蒙"读měng，不是mēng/méng
		"宁夏回族自治区":  20, // Níng xià (N)
		"青海省":      21, // Qīng hǎi (Q)
		"山东省":      22, // Shān dōng (S)
		"山西省":      23, // Shān xī (S)
		"陕西省":      24, // Shǎn xī (S) - 注意"陕"读shǎn，不是shān
		"上海市":      25, // Shàng hǎi (S)
		"四川省":      26, // Sì chuān (S)
		"台湾省":      27, // Tái wān (T)
		"天津市":      28, // Tiān jīn (T)
		"西藏自治区":    29, // Xī zàng (X) - "藏"读zàng，不是cáng
		"香港特别行政区":  30, // Xiāng gǎng (X)
		"新疆维吾尔自治区": 31, // Xīn jiāng (X)
		"云南省":      32, // Yún nán (Y)
		"浙江省":      33, // Zhè jiāng (Z)
		"澳门特别行政区":  34, // Ào mén (A) - 注意"澳"读ào，不是yù
	}
	// 排序省级节点
	sort.Slice(nodes, func(i, j int) bool {
		// 先按TotalCount降序
		if nodes[i].TotalCount != nodes[j].TotalCount {
			return nodes[i].TotalCount > nodes[j].TotalCount
		}

		// TotalCount相同时，按标准行政区划顺序排序
		orderI, existsI := provinceOrderByPinyin[nodes[i].Name]
		orderJ, existsJ := provinceOrderByPinyin[nodes[j].Name]

		// 如果两个都在标准列表中，按标准顺序
		if existsI && existsJ {
			return orderI < orderJ
		}

		// 如果只有一个在标准列表中，标准列表中的排在前面
		if existsI {
			return true
		}
		if existsJ {
			return false
		}

		// 如果都不在标准列表中，保持原有顺序或按拼音排序
		return helper.GetPinyin(nodes[i].Name) < helper.GetPinyin(nodes[j].Name)
	})

	// 对每个省级节点的子节点（城市）进行排序
	for _, province := range nodes {
		sortCityNodes(province.Children)
	}
}

func sortCityNodes(nodes []*storeVisitResp.CityNode) {
	// 排序城市节点
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].TotalCount != nodes[j].TotalCount {
			return nodes[i].TotalCount > nodes[j].TotalCount
		}
		return helper.GetPinyin(nodes[i].Name) < helper.GetPinyin(nodes[j].Name)
	})
}
func buildAdministrativeDivisionTree(stores []microWorkModel.StoreDepartmentModel) []*storeVisitResp.DashboardTree {
	// 第一层：按 store_type 分组
	typeMap := make(map[int32]*storeVisitResp.DashboardTree)

	for _, store := range stores {
		if store.StoreType == 0 || store.CityZh == "" || store.ProvinceZh == "" {
			continue
		}
		// 获取或创建 store_type 节点
		typeNode, ok := typeMap[store.StoreType]
		if !ok {
			typeNode = &storeVisitResp.DashboardTree{
				Name:     helper.RuleThree(store.StoreType == microWorkModel.StoreTypeOle, microWorkModel.StoreTypeOleStr, microWorkModel.StoreTypeRetailStr),
				Level:    storeVisitResp.LevelSource,
				Children: make([]*storeVisitResp.ProvinceNode, 0),
			}
			typeMap[store.StoreType] = typeNode
		}

		// 在 typeNode 下查找或创建 province 节点
		var provinceNode *storeVisitResp.ProvinceNode
		for _, p := range typeNode.Children {
			if p.Name == store.ProvinceZh {
				provinceNode = p
				break
			}
		}

		if provinceNode == nil {
			provinceNode = &storeVisitResp.ProvinceNode{
				Name:     store.ProvinceZh,
				Level:    storeVisitResp.LevelProvince,
				Children: make([]*storeVisitResp.CityNode, 0),
				ThirdIds: make([]int, 0),
			}
			typeNode.Children = append(typeNode.Children, provinceNode)
		}

		// 添加省份级的 wechat_department_id
		if store.WechatDepartmentID != 0 {
			provinceNode.ThirdIds = helper.AppendIfNotExists(provinceNode.ThirdIds, int(store.WechatDepartmentID))
		}

		// 在 provinceNode 下查找或创建 city 节点
		var cityNode *storeVisitResp.CityNode
		for i, c := range provinceNode.Children {
			if c.Name == store.CityZh {
				cityNode = provinceNode.Children[i] // 获取指针
				break
			}
		}

		if cityNode == nil {
			provinceNode.Children = append(provinceNode.Children, &storeVisitResp.CityNode{
				Name:     store.CityZh,
				Level:    storeVisitResp.LevelCity,
				ThirdIds: make([]int, 0),
			})
			cityNode = provinceNode.Children[len(provinceNode.Children)-1]
		}
		// 添加城市级的 wechat_department_id
		if store.WechatDepartmentID != 0 {
			cityNode.ThirdIds = helper.AppendIfNotExists(cityNode.ThirdIds, int(store.WechatDepartmentID))
		}

		// 在 cityNode 下查找或创建 store 节点
		var storeNode *storeVisitResp.StoreNode
		for i, c := range cityNode.Children {
			if c.ThirdId == int(store.WechatDepartmentID) {
				storeNode = cityNode.Children[i] // 获取指针
				break
			}
		}

		if storeNode == nil {
			cityNode.Children = append(cityNode.Children, &storeVisitResp.StoreNode{
				Name:      store.StoreName,
				StoreCode: store.ThirdStoreCode,
				Level:     storeVisitResp.LevelStore,
				ThirdId:   0,
			})
			storeNode = cityNode.Children[len(cityNode.Children)-1]
		}

		// 添加门店级的 wechat_department_id
		if store.WechatDepartmentID != 0 {
			storeNode.ThirdId = int(store.WechatDepartmentID)
		}
	}

	// 将 map 转换为 slice 并返回
	result := make([]*storeVisitResp.DashboardTree, 0, len(typeMap))
	for _, node := range typeMap {
		result = append(result, node)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name > result[j].Name
	})
	return result
}

func (s *MainSvc) ExportDashboardList(c ctx.Context, req storeVisitReq.DashboardExportEmail) error {
	var dataList = make([]storeVisitResp.DashboardExportInfo, 0)
	//提交巡店记录
	req.VisitIds, _ = s.visitMemberDao.GetVisitIdsByDataType(c, req.DashboardTreeReq)
	if len(req.VisitIds) > 0 {
		var storeConnMap = map[string]interface{}{"enterprise_id": req.EnterpriseId}
		if req.SearchName != "" {
			switch req.SearchLevel {
			case storeVisitResp.LevelSource:
				storeConnMap["store_type"] = helper.RuleThree(req.SearchName == microWorkModel.StoreTypeOleStr, microWorkModel.StoreTypeOle, microWorkModel.StoreTypeRetail)
			case storeVisitResp.LevelProvince:
				storeConnMap["province_zh"] = req.SearchName
			case storeVisitResp.LevelCity:
				storeConnMap["city_zh"] = req.SearchName
			}
		}

		storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, storeConnMap)
		if err != nil {
			return err
		}
		var thirdIds []int
		for i, store := range storeList {
			thirdIds = append(thirdIds, int(store.WechatDepartmentID))
			if req.Language == "en" && store.InventoryName != "" {
				storeList[i].StoreName = store.InventoryName
			}
		}
		visitList, err := s.visitMainDao.TotalCountByStoreAndCreator(c, req, thirdIds)
		if err != nil {
			return err
		}
		var memberIds []int64
		for _, row := range visitList {
			memberIds = helper.AppendIfNotExists(memberIds, row.MemberId)
		}
		var storeMap = make(map[int]microWorkModel.StoreDepartmentModel)
		for _, row := range storeList {
			storeMap[int(row.WechatDepartmentID)] = row
		}
		areaList, err := dashboardDao.NewArea().GetAreaListByTypeAsc(c, req.EnterpriseId, "is_virtual=0")
		if err != nil {
			return err
		}
		var areaMap = make(map[int64]string)
		var thirdIdAreaMap = make(map[int]storeVisitResp.DashboardExportInfo)
		for _, row := range areaList {
			if row.Type < dashboardModel.AreaStore {
				areaMap[row.ID] = row.Name
			} else if row.Type == dashboardModel.AreaStore {
				// 组装节点的path
				var tmpPathSlice []int64
				for _, tmp := range strings.Split(row.Path, ",") {
					tmpPathSlice = append(tmpPathSlice, int64(helper.StringToInt(tmp)))
				}
				if len(tmpPathSlice) == 3 && helper.InArray(row.ThirdID, thirdIds) {
					tmp := storeVisitResp.DashboardExportInfo{
						Source:    areaMap[tmpPathSlice[0]],
						Region:    areaMap[tmpPathSlice[1]],
						SubRegion: areaMap[tmpPathSlice[2]],
					}
					if val, ok := storeMap[row.ThirdID]; ok {
						tmp.Province = val.ProvinceZh
						tmp.City = val.CityZh
						tmp.StoreCode = val.ThirdStoreCode
						tmp.StoreName = val.StoreName
					}
					thirdIdAreaMap[row.ThirdID] = tmp
				}
			}
		}
		deptList, err := storeVisitDao.NewDeptMemberDao().GetMany(c, req.EnterpriseId, "member_id in ?", memberIds)
		if err != nil {
			return err
		}
		var deptMap = make(map[int64]string)
		for _, dept := range deptList {
			deptMap[dept.MemberID] = dept.DeptName
		}
		mergedVisits := make(map[storeVisitResp.DashboardExportKey]int64)

		for _, row := range visitList {
			// 获取当前记录的 deptName
			deptName := deptMap[row.MemberId]
			// 创建组合 key
			key := storeVisitResp.DashboardExportKey{
				ThirdId:  row.ThirdId,
				DeptName: deptName,
			}
			// 累加 total
			mergedVisits[key] += row.TotalCount
		}

		// 将合并后的数据转换回列表形式
		for key, total := range mergedVisits {
			if val, ok := thirdIdAreaMap[key.ThirdId]; ok {
				val.TotalCount = total
				val.VisitDeptName = key.DeptName
				dataList = append(dataList, val)
			}
		}
	}

	fileName := fmt.Sprintf("StoreVisitReport_%s~%s__%s", helper.GetFormatDateDay(req.StartTime), helper.GetFormatDateDay(req.EndTime), time.Now().Format("20060102150405")) + ".xlsx"
	defer os.Remove(fileName)
	err := dashboardDataToExcel(req, fileName, dataList)
	if err != nil {
		return err
	}
	mailReq := emailReq.SendMailReq{
		Title:     fileName,
		ToAddress: []string{req.Email},
	}
	mailReq.FileList = append(mailReq.FileList, emailReq.FileInfo{
		LocalPath: fileName,
	})
	_ = emailSvc.NewSendEmail().SendEmailByDiy(mailReq, config.VchangyiEmailConfig{
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		Username: "store_visit@vchangyi.com",
		Password: "REDACTED_FOR_SOURCE_SHARING",
	})
	return nil
}

func dashboardDataToExcel(req storeVisitReq.DashboardExportEmail, fileName string, data []storeVisitResp.DashboardExportInfo) error {
	excelFile := excelize.NewFile()
	sheetName := "StoreVisitReport"
	// 创建sheet
	_ = excelFile.SetSheetName("Sheet1", sheetName)
	err := excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}
	row := make([]any, 0)
	row = []any{
		"Region1", "Region2", "Region3", "Province", "City", "StoreID", "StoreName", helper.RuleThree(req.Language == "en", "VisitDepartment", "巡店部门"), helper.RuleThree(req.Language == "en", "VisitCount", "巡店数"),
	}
	// 创建map来存储解析后的数据
	cityMap := make(map[string]string)
	if req.Language == "en" {
		// 读取省、市JSON文件
		cityJson, _ := os.ReadFile("./storage/store_visit/city.json")
		// 解析JSON到map
		_ = json.Unmarshal(cityJson, &cityMap)
	}
	err = stream.SetRow("A1", row)
	startRow := int64(2)
	for _, item := range data {
		if req.Language == "en" {
			if enName, ok := cityMap[item.Province]; ok {
				item.Province = enName
			}
			if enName, ok := cityMap[item.City]; ok {
				item.City = enName
			}
		}
		rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
		row = []any{item.Source, item.Region, item.SubRegion, item.Province, item.City, item.StoreCode, item.StoreName, item.VisitDeptName, item.TotalCount}
		err = stream.SetRow("A"+rowStartStr, row)
		atomic.AddInt64(&startRow, 1)
	}
	err = stream.Flush()
	if err != nil {
		return err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return err
	}
	return nil
}

func (s *MainSvc) SendPdfToEmail(c ctx.Context, req storeVisitReq.SendPdfToEmailReq) error {
	visitInfo, err := s.visitMainDao.GetOne(c, req.EnterpriseId, "id = ?", req.ID)
	if err != nil || visitInfo.ID == 0 {
		return nil
	}
	if visitInfo.PdfUrl == "" {
		return errors.New("巡店报告还未生成，请稍后重试")
	}
	storeInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": visitInfo.ThirdId,
	})
	if req.Language == "en" && storeInfo.InventoryName != "" {
		storeInfo.StoreName = storeInfo.InventoryName
	}
	title := fmt.Sprintf("【%s】"+helper.RuleThree(req.Language == "en", "Store Visit Report", "巡店报告"), storeInfo.StoreName)
	mailReq := emailReq.SendMailReq{
		Title:     title,
		ToAddress: []string{req.Email},
	}
	if visitInfo.PdfUrl == "" {
		return errors.New("巡店报告PDF文件生成中，请稍后再试")
	}
	mailReq.FileList = append(mailReq.FileList, emailReq.FileInfo{
		Name: fmt.Sprintf("%s_StoreVisitReport_%s.pdf", strings.Replace(storeInfo.StoreName, "/", "_", -1), time.Now().Format("20060102")),
		Url:  visitInfo.PdfUrl,
	})
	_ = emailSvc.NewSendEmail().SendEmailByDiy(mailReq, config.VchangyiEmailConfig{
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		Username: "store_visit@vchangyi.com",
		Password: "REDACTED_FOR_SOURCE_SHARING",
	})
	return nil
}

func (s *MainSvc) CheckCreator(c ctx.Context, req storeVisitReq.VisitInfoReq) error {
	//提交巡店记录
	visitMember, err := s.visitMemberDao.GetOne(c, req.EnterpriseId, "visit_id = ? AND member_id = ? AND type = ?",
		req.ID, req.MemberId, storeVisitModel.VisitCreator)
	if err != nil || visitMember.ID == 0 {
		return errors.New("您不是该巡店记录的创建者")
	}
	return nil
}

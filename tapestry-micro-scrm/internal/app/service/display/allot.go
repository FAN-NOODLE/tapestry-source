package display

import (
	"errors"
	"fmt"
	"math/rand"
	"micro-scrm/internal/app"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	display2 "micro-scrm/internal/app/dao/display"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	"micro-scrm/internal/app/http/request/display"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const enableStoreTaskHorizontalAuth = false

type AllotSvc struct {
	service.Base
	allotDao *display2.AllotDao
}

func NewAllotSvc() *AllotSvc {
	return &AllotSvc{
		allotDao: display2.NewAllotDao(),
	}
}

func (s AllotSvc) GetTaskList(c ctx.Context, request display.MemberTaskListReq) ([]display.MemberTaskListRes, error) {
	result := make([]display.MemberTaskListRes, 0)
	storeIds := make([]int, 0)
	memberInfo, err := membersDao.NewMembers().GetMemberInfo(c, map[string]interface{}{
		"id": request.MemberId,
	})
	if err != nil {
		return result, err
	}
	if len(memberInfo.MemberDepartments) == 0 {
		return result, errors.New("没有门店任务")
	}
	for _, department := range memberInfo.MemberDepartments {
		storeIds = append(storeIds, department.Department.ThirdID)
	}
	return s.allotDao.GetListPagination(c, request, storeIds)
}

// CreateAllotInfo 管理后台创建任务，自动创建allot_info 信息
func (s AllotSvc) CreateAllotInfo(c ctx.Context, request display.AddTaskReq, taskId int64, fileData []display.ReviewImportData) error {
	storeBaseInfo := make([]display.StoreBaseInfo, 0)
	storeDepartments, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"deleted_at":    0,
		"closed_at":     0,
	})
	if err != nil {
		return err
	}
	storeDepartmentMap := make(map[int]display.StoreBaseInfo, 0)
	for _, department := range storeDepartments {
		storeId := int(department.WechatDepartmentID)
		item := display.StoreBaseInfo{
			StoreCode: department.ThirdStoreCode,
			StoreId:   storeId,
			StoreName: department.StoreName,
			StoreType: int(department.StoreType),
			StoreCity: department.City,
		}
		storeDepartmentMap[storeId] = item
		if request.StoreType == microScrmModel.VMStoreTypeAll {
			storeBaseInfo = append(storeBaseInfo, item)
		}
	}
	if request.StoreType == microScrmModel.VMStoreTypeLimit {
		for _, info := range request.StoreList {
			value := display.StoreBaseInfo{
				StoreCode: info.DpName,
				StoreId:   info.DpThirdID,
				StoreName: info.DpStoreName,
				StoreType: info.StoreType,
				StoreCity: "",
			}
			if item, ok := storeDepartmentMap[info.ThirdID]; ok {
				value = item
			}
			storeBaseInfo = append(storeBaseInfo, value)
		}
	}
	if request.ReviewType >= microScrmModel.VMTaskReviewTypeRand {
		specialList, err := s.allotDao.GetManagerList(c, request.EnterpriseId, map[string]interface{}{
			"is_special": 1,
		})
		if err != nil {
			return err
		}
		userIds := make([]string, 0)
		for _, model := range specialList {
			userIds = append(userIds, model.UserId)
		}
		memberInfoList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByUserId(c, request.EnterpriseId, userIds)
		thirdIds := make([]int64, 0)
		for _, v := range memberInfoList {
			thirdIds = append(thirdIds, v.DepartmentThirdId)
		}
		storeDepartmentList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
			"enterprise_id":        request.EnterpriseId,
			"wechat_department_id": thirdIds,
		})
		if err != nil {
			return err
		}
		thirdIdMap := make(map[uint]int32)
		for _, sd := range storeDepartmentList {
			thirdIdMap[sd.WechatDepartmentID] = sd.StoreType
		}
		specialMap := make(map[string]display.SpecialInfo, 0)
		specialTypeMap := make(map[int32][]display.SpecialInfo, 0)
		for _, v := range memberInfoList {
			storeType := thirdIdMap[uint(v.DepartmentThirdId)]
			info := display.SpecialInfo{
				ID:        v.MembersId,
				UserId:    v.MembersUserid,
				StoreCode: v.DepartmentName,
				StoreType: storeType,
				ThirdId:   v.DepartmentThirdId,
			}
			specialMap[v.MembersUserid] = info
			specialTypeMap[storeType] = append(specialTypeMap[storeType], info)
		}

		if request.ReviewType == microScrmModel.VMTaskReviewTypeRand { // 随机平均分配
			for k, specials := range specialTypeMap { // 如果不在奥莱和零售的门店，则special平均分配在不同的类型中
				if !helper.InArray(k, []int32{microWorkModel.StoreTypeOle, microWorkModel.StoreTypeRetail}) {
					for i, special := range specials {
						if i%2 == 0 {
							specialTypeMap[microWorkModel.StoreTypeOle] = append(specialTypeMap[microWorkModel.StoreTypeOle], special)
						} else {
							specialTypeMap[microWorkModel.StoreTypeRetail] = append(specialTypeMap[microWorkModel.StoreTypeRetail], special)
						}
					}
				}
			}
			storeBaseInfoMap := make(map[int32][]display.StoreBaseInfo)
			storeBaseInfoCodeMap := make(map[string]display.StoreBaseInfo)
			newStoreBaseInfo := make([]display.StoreBaseInfo, 0)
			newStoreBaseInfoCodeMap := make(map[string]display.StoreBaseInfo)
			resultStoreBaseInfo := make([]display.StoreBaseInfo, 0)
			for _, v := range storeBaseInfo {
				storeType := int32(v.StoreType)
				storeBaseInfoMap[storeType] = append(storeBaseInfoMap[storeType], v)
				storeBaseInfoCodeMap[v.StoreCode] = v
			}
			// 对每个 storeType 的门店列表随机打乱
			for storeType := range storeBaseInfoMap {
				rand.Shuffle(len(storeBaseInfoMap[storeType]), func(i, j int) {
					storeBaseInfoMap[storeType][i], storeBaseInfoMap[storeType][j] = storeBaseInfoMap[storeType][j], storeBaseInfoMap[storeType][i]
				})
			}

			for k, specialList := range specialTypeMap {
				storeBase := s.getSpecialInfo(specialList, storeBaseInfoMap[k])
				newStoreBaseInfo = append(newStoreBaseInfo, storeBase...)

			}
			for _, info := range newStoreBaseInfo {
				newStoreBaseInfoCodeMap[info.StoreCode] = info
			}
			for k, info := range storeBaseInfoCodeMap {
				if item, ok := newStoreBaseInfoCodeMap[k]; ok {
					resultStoreBaseInfo = append(resultStoreBaseInfo, item)
				} else {
					resultStoreBaseInfo = append(resultStoreBaseInfo, info)
				}
			}
			storeBaseInfo = resultStoreBaseInfo
		}
		if request.ReviewType == microScrmModel.VMTaskReviewTypeFile { // 读取匹配文件内容，然后分配
			// 读取文件内容
			if len(fileData) > 0 {
				storeBaseInfoMap := make(map[string]display.StoreBaseInfo, 0)
				for _, v := range storeBaseInfo {
					storeBaseInfoMap[v.StoreCode] = v
				}
				for _, data := range fileData {
					if !helper.InArray(data.UserId, userIds) {
						continue
					}
					if specialItem, exists := specialMap[data.UserId]; !exists {
						continue
					} else {
						if specialItem.StoreCode == data.StoreCode {
							continue
						}
					}
					if item, ok := storeBaseInfoMap[data.StoreCode]; ok {
						item.SpecialUserId = data.UserId
						storeBaseInfoMap[data.StoreCode] = item
					}
				}
				newStoreBaseInfo := make([]display.StoreBaseInfo, 0)
				for _, info := range storeBaseInfoMap {
					newStoreBaseInfo = append(newStoreBaseInfo, info)
				}
				storeBaseInfo = newStoreBaseInfo
			}
		}
	}
	err = display2.NewAllotDao().CreateTaskByStoreIds(c, request.EnterpriseId, taskId, request.DisplayName, storeBaseInfo)
	return err
}

// SaveTaskInfo VM任务导购的存草稿和提交
func (s AllotSvc) SaveTaskInfo(c ctx.Context, params display.StoreTaskSaveInfoReq) error {
	allotInfo, err := s.allotDao.GetOne(c, params.EnterpriseId, map[string]interface{}{
		"id": params.AllotID,
	})
	if err != nil {
		return err
	}
	if allotInfo.ID == 0 || !helper.InArray(allotInfo.FinishStatus, []int32{microScrmModel.VMFinishStatusIng, microScrmModel.VMFinishStatusFail}) {
		return errors.New("该任务已被同事抢先提交了，请刷新查看哦~")
	}
	taskInfo, err := display2.NewTaskDao().GetOne(c, params.EnterpriseId, map[string]interface{}{
		"id": params.ID,
	})
	if err != nil {
		return err
	}
	if taskInfo.Status == microScrmModel.DisplayStatusEnd && allotInfo.FinishStatus == microScrmModel.VMFinishStatusIng { // 任务已结束，未提交的门店VM不能提交
		return errors.New("任务已结束")
	}
	if params.SubmitType == microScrmModel.VMSubmitCheck {
		key := fmt.Sprintf("save_display_vm:%d", allotInfo.ID)
		defer s.Redis().Del(c, key)
		if !s.Redis().SetNX(c, key, 1, 5*time.Second).Val() { // 提交审核加下锁，防止覆盖
			return errors.New("该任务已提交了")
		}
	}
	memberInfo, err := membersDao.NewMembers().GetMemberInfo(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"id":            params.MemberId,
	})
	if err != nil {
		return err
	}
	if len(memberInfo.MemberDepartments) == 0 {
		return errors.New("用户数据异常")
	}
	canSave := false
	for _, department := range memberInfo.MemberDepartments {
		if allotInfo.StoreID == department.Department.ThirdID {
			canSave = true
		}
	}
	if !canSave {
		return errors.New("用户数据异常")
	}
	imageList, err := s.allotDao.GetImageList(c, params.EnterpriseId, map[string]interface{}{
		"allot_id": params.AllotID,
	})
	if err != nil {
		return err
	}
	imageListIds := make([]int64, 0)
	for _, item := range imageList {
		imageListIds = append(imageListIds, item.ID)
	}
	saveImageList := params.ImageList
	oldImageListIds := make([]int64, 0)
	saveImageListMap := make(map[int64]display.ImageInfo, 0)
	for _, model := range saveImageList {
		index := model.ID
		if model.ID > 0 {
			oldImageListIds = append(oldImageListIds, model.ID)
		}
		if model.ID == 0 {
			index = app.UID()
		}
		saveImageListMap[index] = model
	}
	diffIds := helper.ArrayDiff(imageListIds, oldImageListIds)
	if len(diffIds) > 0 { // 差异照片ID删除照片
		err = s.allotDao.DeleteImageList(c, params.EnterpriseId, diffIds)
		for _, id := range diffIds {
			delete(saveImageListMap, id)
		}
	}
	err = s.allotDao.SaveAllotData(c, params.EnterpriseId, params.ID, params.AllotID, saveImageListMap, oldImageListIds)
	// 保存操作记录
	nowTime := time.Now().Unix()
	submitDataLog := microScrmModel.DisplayTaskAllotSubmitModel{
		BaseModel: models.BaseModel{
			ID:        app.UID(),
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		DisplayTaskID: params.ID,
		AllotID:       params.AllotID,
		SubmitType:    int32(params.SubmitType),
		MemberID:      params.MemberId,
		MemberName:    params.MemberName,
		MemberUserID:  params.MemberUserId,
	}
	err = s.allotDao.SaveAllotSubmitDataLog(c, params.EnterpriseId, submitDataLog)
	updateMap := make(map[string]interface{}, 0)
	updateMap = map[string]interface{}{
		"owner_member_id": params.OwnerMemberId,
		"owner_user_id":   params.OwnerMemberUserId,
		"owner_name":      params.OwnerMemberName,
	}
	sendUser := params.OwnerMemberUserId
	if params.SubmitType == microScrmModel.VMSubmitCheck {
		updateMap["finish_status"] = microScrmModel.VMFinishStatusCheck
		updateMap["submit_time"] = nowTime
		if allotInfo.SpecialUserId != "" && allotInfo.SpecialCheck == microScrmModel.VMSpecialCheckNo {
			sendUser = allotInfo.SpecialUserId
			updateMap["special_check"] = microScrmModel.VMSpecialCheckWait
		}
	}
	err = s.allotDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.AllotID,
	}, updateMap)
	if params.SubmitType == microScrmModel.VMSubmitCheck {
		// VM发送审核通知
		uri := fmt.Sprintf("display-task-audit?id=%d&store_id=%d", allotInfo.DisplayTaskID, allotInfo.StoreID)
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "陈列任务名称",
			"value":   allotInfo.DisplayName,
		})
		content = append(content, map[string]string{
			"keyname": "提交门店",
			"value":   allotInfo.StoreName,
		})
		submitTime := helper.GetFormatDateByParams(nowTime)
		content = append(content, map[string]string{
			"keyname": "提交时间",
			"value":   submitTime,
		})
		go func() {
			_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: params.EnterpriseId,
				},
				Title:   "VM陈列任务待审核提醒",
				Uri:     uri,
				Party:   nil,
				UserIds: []string{sendUser},
				Content: content,
				AppSign: qyWork.VM,
			})
			if err != nil {
				app.Logger().Error(c).Msgf("VM通知消息发送失败：%+v", err)
			}
		}()
	}
	return err
}

// CreateLookLog 导购查看VM备注的信息上报
func (s AllotSvc) CreateLookLog(c ctx.Context, request display.LookReq) error {
	data := microScrmModel.DisplayTaskLookLogModel{
		BaseModel: models.BaseModel{
			ID:        app.UID(),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		AllotDataID: request.ImageID,
		MemberID:    request.MemberId,
		AllotID:     request.AllotID,
		CommentID:   request.CommentId,
	}
	return s.allotDao.CreateLookLog(c, request.EnterpriseId, data)
}

func (s AllotSvc) AuthInfo(c ctx.Context, request display.AuthReq) display.AuthInfo {
	var result display.AuthInfo
	info, err := s.allotDao.AuthInfo(c, request.EnterpriseId, request.MemberId)
	if err != nil || info.ID == 0 {
		return result
	}

	result.IsVm = info.IsCheck
	result.IsManager = info.IsManager
	result.IsSpecial = info.IsSpecial
	return result
}

// VmList 获取VM列表
func (s AllotSvc) VmList(c ctx.Context, req request.EnterpriseBaseReq) ([]display.VmListRes, error) {
	vmListInfo := make([]display.VmListRes, 0)
	memberList, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	if err != nil {
		return vmListInfo, err
	}
	memberMap := make(map[int64]membersModel.MemberModel, 0)
	for _, v := range memberList {
		memberMap[v.ID] = v
	}
	list, err := s.allotDao.VmList(c, req.EnterpriseId, map[string]interface{}{
		"is_check":     1,
		"is_corporate": 1,
		"deleted_at":   0,
	})
	if err != nil {
		return vmListInfo, err
	}
	for _, item := range list {
		if info, ok := memberMap[item]; ok {
			vmInfo := display.VmListRes{
				MemberId:   info.ID,
				UserId:     info.Userid,
				MemberName: info.Name,
			}
			vmListInfo = append(vmListInfo, vmInfo)
		}
	}
	return vmListInfo, err
}

// StoreList 他店陈列数据展示店铺
func (s AllotSvc) StoreList(c ctx.Context, req display.StoreSearchReq) ([]display.VmStore, error) {
	result := make([]display.VmStore, 0)
	list := make([]microScrmModel.DisplayTaskAllotModel, 0)
	var err error
	queryStr := " finish_status= ? and display_task_id= ? and store_id!= ? "
	if req.Name != "" {
		queryStr += " and ( store_name like ? or store_city like ? ) "
		list, err = s.allotDao.GetMany(c, req.EnterpriseId, queryStr, microScrmModel.VMFinishStatusSuccess, req.ID, req.StoreId, "%"+req.Name+"%", "%"+req.Name+"%")
		if err != nil {
			return nil, err
		}
	} else {
		list, err = s.allotDao.GetMany(c, req.EnterpriseId, queryStr, microScrmModel.VMFinishStatusSuccess, req.ID, req.StoreId)
		if err != nil {
			return nil, err
		}
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, req.EnterpriseId)
	storeInfoMap := make(map[string][]display.VmStoreBase)
	for _, item := range list {
		storeBaseInfo := display.VmStoreBase{
			StoreId:   item.StoreID,
			StoreName: item.StoreName,
		}
		if v, ok := storeNameMap[storeBaseInfo.StoreId]; ok {
			storeBaseInfo.StoreName = v
		}
		storeInfoMap[item.StoreCity] = append(storeInfoMap[item.StoreCity], storeBaseInfo)
	}
	for k, info := range storeInfoMap {
		storeInfo := display.VmStore{
			StoreCity: k,
			StoreInfo: info,
		}
		result = append(result, storeInfo)
	}
	// 门店城市首字母排序
	sort.Slice(result, func(i, j int) bool {
		return helper.GetPinYinFirstLetter(helper.CovertPinYin(result[i].StoreCity)) < helper.GetPinYinFirstLetter(helper.CovertPinYin(result[j].StoreCity))
	})
	return result, err
}

// GetListTips 获取列表的任务状态下数量
func (s AllotSvc) GetListTips(c ctx.Context, req display.ListTipsReq) (display.ListTips, error) {
	var res display.ListTips
	storeIds := make([]int, 0)
	//whereMap := make(map[string]interface{}, 0)
	query := ""
	args := make([]interface{}, 0)
	if req.IsVm == 0 { // 导购的身份查询信息
		memberInfo, err := membersDao.NewMembers().GetMemberInfo(c, map[string]interface{}{
			"id": req.MemberId,
		})
		if err != nil {
			return res, err
		}
		if len(memberInfo.MemberDepartments) == 0 {
			return res, errors.New("没有门店任务")
		}
		for _, department := range memberInfo.MemberDepartments {
			storeIds = append(storeIds, department.Department.ThirdID)
		}
		//whereMap = map[string]interface{}{
		//	"store_id": storeIds,
		//}
		query = "store_id in ?"
		args = append(args, storeIds)
	}
	if req.IsVm == 1 { // VM的身份查询
		if !s.CheckIsVm(c, req.EnterpriseId, req.MemberId) {
			return res, nil
		}
		//whereMap["owner_member_id"] = req.MemberId
		if s.CheckSpecialInfo(c, req.EnterpriseId, req.MemberId) {
			query = "special_user_id=? and display_task_allot.special_check>0"
			args = append(args, req.MemberUserId)
		} else {
			query = "owner_member_id=? and (display_task_allot.special_check=2 or special_user_id='')"
			args = append(args, req.MemberId)
		}
	}
	if req.ID > 0 {
		//whereMap["display_task_id"] = req.ID
		query += " and display_task_id = ? "
		args = append(args, req.ID)
	}
	many, err := display2.NewTaskDao().GetMany(c, req.EnterpriseId, map[string]interface{}{
		"show_type":  microScrmModel.ShowStatusUp,
		"status":     microScrmModel.DisplayStatusEnd,
		"deleted_at": 0,
	})
	displayEndTaskIds := make([]int64, 0)
	if len(many) > 0 {
		for _, value := range many {
			displayEndTaskIds = append(displayEndTaskIds, value.ID)
		}
	}
	countList, err := s.allotDao.GetCheckTaskCount(c, req.EnterpriseId, query, args...)
	info := make(map[int]int64, 0)
	var allNum, vmAllNum int64
	for _, value := range countList {
		if value.FinishStatus == microScrmModel.VMFinishStatusIng && helper.InArray(value.DisplayTaskId, displayEndTaskIds) {
			continue
		}
		info[value.FinishStatus] += value.Num
		if helper.InArray(value.FinishStatus, []int{microScrmModel.VMFinishStatusCheck, microScrmModel.VMFinishStatusDoubleCheck, microScrmModel.VMFinishStatusFail, microScrmModel.VMFinishStatusSuccess}) {
			vmAllNum += value.Num
		}
		allNum += value.Num
	}
	res.WaitNum = info[microScrmModel.VMFinishStatusIng]
	res.CheckNum = info[microScrmModel.VMFinishStatusCheck]
	res.DoubleCheckNum = info[microScrmModel.VMFinishStatusDoubleCheck]
	res.FailNum = info[microScrmModel.VMFinishStatusFail]
	res.SuccessNum = info[microScrmModel.VMFinishStatusSuccess]
	res.AllNum = allNum
	if req.IsVm == 1 {
		res.AllNum = vmAllNum
	}
	return res, err
}

func (s AllotSvc) GetCheckList(c ctx.Context, req display.MemberTaskListReq) ([]display.CheckListRes, error) {
	result := make([]display.CheckListRes, 0)
	isSpecial := s.CheckSpecialInfo(c, req.EnterpriseId, req.MemberId)
	list, err := s.allotDao.GetCheckList(c, req, isSpecial)
	if err != nil {
		return result, err
	}

	taskId := make([]int64, 0)
	for _, model := range list {
		taskId = append(taskId, model.Id)
	}
	var argValue interface{}
	query := "display_task_id in ? "
	if isSpecial {
		query += " and special_user_id=? and display_task_allot.special_check>0"
		argValue = req.MemberUserId
	} else {
		query += " and owner_member_id=? and (display_task_allot.special_check=2 or display_task_allot.special_user_id='')"
		argValue = req.MemberId
	}
	countList, err := s.allotDao.GetCheckTaskCount(c, req.EnterpriseId, query, taskId, argValue)
	if err != nil {
		return result, err
	}
	countMap := make(map[int64][]map[int]int64, 0)
	for _, value := range countList {
		info := make(map[int]int64, 0)
		info[value.FinishStatus] = value.Num
		countMap[value.DisplayTaskId] = append(countMap[value.DisplayTaskId], info)
	}
	for _, model := range list {
		var tipsInfo display.ListTips
		tipsSlice := countMap[model.Id]
		for _, m := range tipsSlice {
			for k, v := range m {
				switch k {
				case microScrmModel.VMFinishStatusIng:
					tipsInfo.WaitNum = v
				case microScrmModel.VMFinishStatusCheck:
					tipsInfo.CheckNum = v
				case microScrmModel.VMFinishStatusDoubleCheck:
					tipsInfo.DoubleCheckNum = v
				case microScrmModel.VMFinishStatusFail:
					tipsInfo.FailNum = v
				case microScrmModel.VMFinishStatusSuccess:
					tipsInfo.SuccessNum = v
				}
			}
		}
		taskInfo := display.TaskInfoRes{
			ID:          model.Id,
			DisplayName: model.DisplayName,
			Content:     model.Content,
			Status:      int32(model.Status),
			StartTime:   model.StartTime,
			EndTime:     model.EndTime,
		}
		helper.Unmarshal(model.ClassUrl, &taskInfo.ClassURL)
		item := display.CheckListRes{
			TaskInfoRes: taskInfo,
			ListTips:    tipsInfo,
		}
		result = append(result, item)
	}
	return result, err
}

func (s AllotSvc) GetCheckStoreList(c ctx.Context, req display.StoreListInfoReq) ([]display.StoreListInfo, error) {
	var list []microScrmModel.DisplayTaskAllotModel
	result := make([]display.StoreListInfo, 0)
	var err error
	if !s.CheckIsVm(c, req.EnterpriseId, req.MemberId) {
		return result, err
	}
	finishStatus := 0
	isSpecial := s.CheckSpecialInfo(c, req.EnterpriseId, req.MemberId)
	queryStr := "display_task_id=? and finish_status> ? "
	if req.FinishStatus >= microScrmModel.VMFinishStatusCheck {
		finishStatus = req.FinishStatus
		queryStr = "display_task_id=? and finish_status= ? "
	}
	var argValue interface{}
	if isSpecial {
		queryStr += " and special_user_id=? and special_check>0"
		argValue = req.MemberUserId
	} else {
		queryStr += " and owner_member_id=? and (special_check=2 or special_user_id='')"
		argValue = req.MemberId
	}
	if req.Name != "" {
		queryStr += " and store_name like ? "
		list, err = s.allotDao.GetMany(c, req.EnterpriseId, queryStr, req.ID, finishStatus, argValue, "%"+req.Name+"%")
	} else {
		list, err = s.allotDao.GetMany(c, req.EnterpriseId, queryStr, req.ID, finishStatus, argValue)
	}
	if err != nil {
		return result, err
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, req.EnterpriseId)
	for _, model := range list {
		info := display.StoreListInfo{
			ID:           model.DisplayTaskID,
			AllotID:      model.ID,
			StoreName:    model.StoreName,
			StoreCode:    model.StoreCode,
			StoreId:      model.StoreID,
			FinishStatus: model.FinishStatus,
			Grade:        model.Grade,
		}
		if v, ok := storeNameMap[info.StoreId]; ok {
			info.StoreName = v
		}
		result = append(result, info)
	}
	return result, err
}

func (s AllotSvc) GetStoreTaskInfo(c ctx.Context, req display.MemberStoreTaskInfoReq) (display.StoreTaskInfo, error) {
	var result display.StoreTaskInfo
	storeIds := make([]int, 0)
	isOwnStore := false
	memberInfo, err := membersDao.NewMembers().GetMemberInfo(c, map[string]interface{}{
		"id": req.MemberId,
	})
	if err != nil {
		return result, err
	}
	if req.StoreId == 0 {
		if len(memberInfo.MemberDepartments) == 0 {
			return result, errors.New("没有门店任务")
		}
		for _, department := range memberInfo.MemberDepartments {
			storeIds = append(storeIds, department.Department.ThirdID)
		}
	} else {
		for _, department := range memberInfo.MemberDepartments {
			if department.Department.ThirdID == req.StoreId {
				isOwnStore = true
				break
			}
		}
		storeIds = append(storeIds, req.StoreId)
	}
	info, err := s.allotDao.GetVmStoreInfo(c, req, storeIds)
	if err != nil {
		return result, err
	}
	member, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"userid":        info.SpecialUserId,
	})
	if err != nil {
		return result, err
	}
	if enableStoreTaskHorizontalAuth &&
		!req.IsAdmin &&
		req.StoreId != 0 &&
		!canAccessStoreTask(req, info, isOwnStore, s.CheckIsVm(c, req.EnterpriseId, req.MemberId), member.ID) {
		return result, errors.New("无权访问")
	}
	imageLookList, err := s.allotDao.GetLooKLog(c, req.EnterpriseId, info.ID, req.MemberId)
	if err != nil {
		return result, err
	}
	imageLookMap := make(map[int64]int64, 0)
	for _, res := range imageLookList {
		imageLookMap[res.AllotDataId] = res.CommentId
	}
	reviewRemark := info.ReviewRemark
	// special身份
	specialName := ""
	isSpecialCheck := false
	if req.MemberId == member.ID && info.SpecialCheck == microScrmModel.VMSpecialCheckWait { // special审核，未审核，留言数据未special的留言备注数据
		isSpecialCheck = true
		reviewRemark = info.SpecialRemark
	}
	specialName = member.Name
	// corporate人员审核
	isCorporateCheck := false
	if req.MemberId == info.OwnerMemberID && (info.SpecialUserId == "" || info.SpecialCheck == microScrmModel.VMSpecialCheckYes) {
		isCorporateCheck = true
	}

	imageList := make([]display.ImageInfo, 0)
	for _, model := range info.ImageList {
		isLook := 0
		imageStatus := model.Status
		if model.IsDraft == 1 && isSpecialCheck { //如果是special人员查看，并且是待审核的状态时，有草稿数据要查看草稿箱的状态数据
			imageStatus = model.DraftStatus
		}
		item := display.ImageInfo{
			ID:         model.ID,
			Identifier: model.Identifier,
			Url:        model.ImageURL,
			Status:     int(imageStatus),
			IsDraft:    int(model.IsDraft),
		}
		var num int
		var vmNum int
		commentList := make([]display.CommentInfo, 0)
		for _, remarkModel := range model.Remark {
			isSelf := 0
			if remarkModel.MemberID == req.MemberId {
				isSelf = 1
			}
			if remarkModel.IsVm == 1 && remarkModel.ID > imageLookMap[model.ID] { // 如果门店任务未通过，VM端有新的备注消息,导购需要查看
				isLook = 1 // 1为展示气泡，0为未不展示气泡
				vmNum++
			}
			commentInfo := display.CommentInfo{
				ID:           remarkModel.ID,
				MemberId:     remarkModel.MemberID,
				MemberUserId: remarkModel.MemberUserID,
				MemberName:   remarkModel.MemberName,
				Avatar:       remarkModel.Avatar,
				StoreCode:    remarkModel.StoreCode,
				CommentTime:  remarkModel.CommentTime,
				IsVm:         remarkModel.IsVm,
				IsSelf:       isSelf,
				ReplyName:    remarkModel.ReplyName,
			}
			helper.Unmarshal(remarkModel.Content, &commentInfo.Content)
			commentList = append(commentList, commentInfo)
			num++
		}
		item.Comment = commentList
		item.Num = num
		item.VmNum = vmNum
		item.IsLook = isLook
		imageList = append(imageList, item)
	}
	// 未通过的VM任务列表，查询正常的VM人员信息，如果VM已离职，则需要重新选择VM
	ownerMemberId := info.OwnerMemberID
	ownerUserId := info.OwnerUserID
	ownerName := info.OwnerName
	if info.FinishStatus == microScrmModel.VMFinishStatusFail {
		if !s.CheckIsVm(c, req.EnterpriseId, req.MemberId) {
			ownerMemberId, ownerUserId, ownerName = 0, "", ""
		}
	}
	storeInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]interface{}{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": info.StoreID,
		"deleted_at":           0,
	})
	storeName := info.StoreName
	if storeInfo.ID > 0 {
		storeName = storeInfo.StoreName
	}

	result = display.StoreTaskInfo{
		ID:                info.DisplayTaskID,
		DisplayName:       info.DisplayName,
		AllotID:           info.ID,
		OwnerMemberId:     ownerMemberId,
		OwnerMemberUserId: ownerUserId,
		OwnerMemberName:   ownerName,
		FinishStatus:      info.FinishStatus,
		Grade:             info.Grade,
		ReviewRemark:      reviewRemark,
		StoreCode:         info.StoreCode,
		StoreName:         storeName,
		SubmitTime:        info.SubmitTime,
		ReviewTime:        info.ReviewTime,
		ImageList:         imageList,
		SpecialUserId:     info.SpecialUserId,
		SpecialName:       specialName,
		IsSpecialCheck:    isSpecialCheck,
		IsCorporateCheck:  isCorporateCheck,
	}
	return result, err
}

func canAccessStoreTask(req display.MemberStoreTaskInfoReq, info microScrmModel.DisplayTaskAllotModel, isOwnStore, isVM bool, specialMemberID int64) bool {
	if req.StoreId == 0 {
		return true
	}
	if isOwnStore {
		return true
	}
	if !isVM {
		return false
	}
	if specialMemberID == req.MemberId && info.SpecialCheck > microScrmModel.VMSpecialCheckNo {
		return true
	}
	return info.OwnerMemberID == req.MemberId &&
		(info.SpecialUserId == "" || info.SpecialCheck == microScrmModel.VMSpecialCheckYes)
}

func (s AllotSvc) VmCheckStoreTask(c ctx.Context, req display.StoreTaskReviewReq) error {
	allotInfo, err := s.allotDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.AllotID,
	})
	if err != nil {
		return err
	}
	if allotInfo.ID == 0 || helper.InArray(allotInfo.FinishStatus, []int32{microScrmModel.VMFinishStatusFail, microScrmModel.VMFinishStatusSuccess}) {
		return errors.New("门店VM数据状态异常,已被审核了")
	}
	isSpecialCheck := false
	if allotInfo.SpecialUserId != "" {
		if allotInfo.SpecialCheck == microScrmModel.VMSpecialCheckWait {
			isSpecialCheck = true
		}
	}
	if isSpecialCheck && allotInfo.SpecialUserId != req.MemberUserId {
		return errors.New("该门店的Special人员不是您")
	}
	if allotInfo.OwnerMemberID != req.MemberId && !isSpecialCheck {
		return errors.New("该门店的所属VM不是您")
	}

	err = s.checkAllotInfo(c, req.MemberId, req.MemberName, req.MemberUserId, allotInfo, req.StoreTaskReviewBasReq, isSpecialCheck)
	return err
}

// VmSendNotice vm点击发送提醒
func (s AllotSvc) VmSendNotice(c ctx.Context, req display.SendNoticeReq) error {
	var err error
	if !s.CheckIsVm(c, req.EnterpriseId, req.MemberId) {
		return err
	}
	title := "VM陈列任务未通过提醒"
	content := s.noticeContent(c, req.EnterpriseId, req.ID)
	if req.AllotID == 0 { // 一键提醒
		whereMap := map[string]interface{}{
			"display_task_id": req.ID,
			"finish_status":   microScrmModel.VMFinishStatusFail,
			"deleted_at":      0,
		}
		specialInfo, err := display2.NewTaskDao().GetSpecialInfo(c, req.EnterpriseId, map[string]interface{}{
			"member_id":  req.MemberId,
			"is_special": 1,
		})
		if err != nil {
			return err
		}
		if specialInfo.ID > 0 {
			whereMap["special_user_id"] = specialInfo.UserId
			whereMap["special_check"] = microScrmModel.VMSpecialCheckYes
		} else {
			whereMap["owner_member_id"] = req.MemberId
		}
		many, err := s.allotDao.GetMany(c, req.EnterpriseId, whereMap)
		if err != nil {
			return err
		}
		for _, model := range many {
			storeId := strconv.Itoa(model.StoreID)
			uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", req.ID, model.StoreID)
			_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: req.EnterpriseId,
				},
				Title:   title,
				Uri:     uri,
				Party:   []string{storeId},
				UserIds: nil,
				Content: content,
				AppSign: qyWork.VM,
			})
		}
	} else { // 单个提醒
		info, err := s.allotDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
			"display_task_id": req.ID,
			"id":              req.AllotID,
			"finish_status":   microScrmModel.VMFinishStatusFail,
		})
		if err != nil {
			return err
		}
		if info.ID == 0 {
			return errors.New("门店任务不存在")
		}
		storeId := strconv.Itoa(info.StoreID)
		uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", req.ID, info.StoreID)
		_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			Title:   title,
			Uri:     uri,
			Party:   []string{storeId},
			UserIds: nil,
			Content: content,
			AppSign: qyWork.VM,
		})
	}
	return err
}

func (s AllotSvc) GetAllotList(c ctx.Context, req display.SearchAllotReq) (int64, []display.AllotListRes, error) {
	result := make([]display.AllotListRes, 0)
	total, list, err := s.allotDao.GetAllotList(c, req)
	if err != nil {
		return total, result, err
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, req.EnterpriseId)
	for _, item := range list {
		reviewName := ""
		if len(item.Review) > 0 {
			reviewName = item.Review[0].ReviewName
		}
		submitName := ""
		var submitMemberId int64
		if item.FinishStatus >= microScrmModel.VMFinishStatusCheck {
			submitName = item.Submit[0].MemberName
			submitMemberId = item.Submit[0].MemberID
		}
		info := display.AllotListRes{
			AllotId:         item.ID,
			DisplayTaskId:   item.DisplayTaskID,
			DisplayName:     item.DisplayName,
			StoreId:         item.StoreID,
			StoreCode:       item.StoreCode,
			StoreName:       item.StoreName,
			StoreType:       item.StoreType,
			StoreCity:       item.StoreCity,
			FinishStatus:    item.FinishStatus,
			Grade:           item.Grade,
			OwnerMemberId:   item.OwnerMemberID,
			OwnerMemberName: item.OwnerName,
			SubmitTime:      item.SubmitTime,
			ReviewTime:      item.ReviewTime,
			SubmitName:      submitName,
			SubmitMemberId:  submitMemberId,
			ReviewName:      reviewName,
		}
		if v, ok := storeNameMap[info.StoreId]; ok {
			info.StoreName = v
		}
		result = append(result, info)
	}
	return total, result, err
}

func (s AllotSvc) GetAllotNum(c ctx.Context, req display.SearchAllotBaseReq) (display.ListTips, error) {
	var listTips display.ListTips
	list, err := s.allotDao.GetAllotListCount(c, req)
	if err != nil {
		return listTips, err
	}
	tipsNumMap := make(map[int]int64, 0)
	var allNum int64
	for _, info := range list {
		tipsNumMap[info.FinishStatus] += info.Num
		allNum += info.Num
	}
	listTips.AllNum = allNum
	listTips.WaitNum = tipsNumMap[microScrmModel.VMFinishStatusIng]
	listTips.CheckNum = tipsNumMap[microScrmModel.VMFinishStatusCheck]
	listTips.FailNum = tipsNumMap[microScrmModel.VMFinishStatusFail]
	listTips.SuccessNum = tipsNumMap[microScrmModel.VMFinishStatusSuccess]
	listTips.DoubleCheckNum = tipsNumMap[microScrmModel.VMFinishStatusDoubleCheck]
	return listTips, err
}

// GetAllotListExport 陈列数据的导出
func (s AllotSvc) GetAllotListExport(c ctx.Context, req display.SearchAllotReq) (int64, []display.AllotListExportRes, error) {
	result := make([]display.AllotListExportRes, 0)
	total, list, err := s.allotDao.GetAllotList(c, req)
	if err != nil {
		return total, result, err
	}
	finishStatusStrMap := make(map[int32]string, 0)
	finishStatusStrMap[microScrmModel.VMFinishStatusIng] = "未提交"
	finishStatusStrMap[microScrmModel.VMFinishStatusCheck] = "待审核"
	finishStatusStrMap[microScrmModel.VMFinishStatusFail] = "未通过"
	finishStatusStrMap[microScrmModel.VMFinishStatusSuccess] = "已通过"
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, req.EnterpriseId)
	for _, item := range list {
		reviewName := ""
		if len(item.Review) > 0 {
			reviewName = item.Review[0].ReviewName
		}
		submitName := ""
		var submitMemberId int64
		if item.FinishStatus >= microScrmModel.VMFinishStatusCheck {
			submitName = item.Submit[0].MemberName
			submitMemberId = item.Submit[0].MemberID
		}
		submitTimeStr := "--"
		gradeStr := "--"
		storeTypeStr := "--"
		reviewTimeStr := "--"
		if item.FinishStatus == microScrmModel.VMFinishStatusSuccess {
			gradeStr = strconv.Itoa(int(item.Grade))
		}
		if item.StoreType == microWorkDao.STORE_TYPE_RETAIL {
			storeTypeStr = "零售"
		}
		if item.StoreType == microWorkDao.STORE_TYPE_OUTLETS {
			storeTypeStr = "奥莱"
		}
		if item.SubmitTime > 0 {
			submitTimeStr = helper.GetFormatDateByParams(item.SubmitTime)
		}
		if item.ReviewTime > 0 {
			reviewTimeStr = helper.GetFormatDateByParams(item.ReviewTime)
		}
		info := display.AllotListExportRes{
			AllotListRes: display.AllotListRes{
				AllotId:         item.ID,
				DisplayTaskId:   item.DisplayTaskID,
				DisplayName:     item.DisplayName,
				StoreId:         item.StoreID,
				StoreCode:       item.StoreCode,
				StoreName:       item.StoreName,
				StoreType:       item.StoreType,
				StoreCity:       item.StoreCity,
				FinishStatus:    item.FinishStatus,
				Grade:           item.Grade,
				OwnerMemberId:   item.OwnerMemberID,
				OwnerMemberName: item.OwnerName,
				SubmitTime:      item.SubmitTime,
				ReviewTime:      item.ReviewTime,
				SubmitName:      submitName,
				SubmitMemberId:  submitMemberId,
				ReviewName:      reviewName,
			},
			StoreTypeStr:    storeTypeStr,
			GradeStr:        gradeStr,
			SubmitTimeStr:   submitTimeStr,
			ReviewTimeStr:   reviewTimeStr,
			FinishStatusStr: finishStatusStrMap[item.FinishStatus],
		}
		if v, ok := storeNameMap[info.StoreId]; ok {
			info.StoreName = v
		}
		result = append(result, info)
	}
	return total, result, err
}

// SetNoticeSetting 自动提醒设置
func (s AllotSvc) SetNoticeSetting(c ctx.Context, req display.NoticeSettingReq) error {
	info, err := microBossDao.NewEnterpriseConfig().GetInfo(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"deleted_at":    0,
		"key":           "vm-notice-setting",
	})
	if err != nil {
		return err
	}
	if info.ID == 0 {
		nowTime := time.Now().Unix()
		data := microBossModel.EnterpriseConfigModel{
			ID:           uint64(app.UID()),
			EnterpriseID: uint64(req.EnterpriseId),
			Type:         0,
			Key:          "vm-notice-setting",
			Value:        req.Value,
			Desc:         "vm陈列消息提醒",
			CreatedAt:    uint(nowTime),
			UpdatedAt:    uint(nowTime),
			DeletedAt:    0,
		}
		err = microBossDao.NewEnterpriseConfig().Create(c, data)
	}
	if info.ID > 0 {
		err = microBossDao.NewEnterpriseConfig().Update(c, map[string]interface{}{
			"value": req.Value,
		}, map[string]interface{}{
			"id":         info.ID,
			"deleted_at": 0,
			"key":        "vm-notice-setting",
		})
	}
	return err
}

// AllotCheck 管理后台的门店VM任务审核
func (s AllotSvc) AllotCheck(c ctx.Context, req display.StoreTaskReviewWebReq) error {
	allotInfo, err := s.allotDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.AllotID,
	})
	if err != nil {
		return err
	}
	if allotInfo.ID == 0 || helper.InArray(allotInfo.FinishStatus, []int32{microScrmModel.VMFinishStatusFail, microScrmModel.VMFinishStatusSuccess}) {
		return errors.New("门店VM数据状态异常,已被审核了")
	}
	info, err := microBossDao.NewEnterpriseAdmin().GetInfo(c, microBossModel.EnterpriseAdminModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID: req.AdminId,
			},
			EnterpriseID: req.EnterpriseId,
		},
	})
	if err != nil {
		return err
	}
	adminName := info.Username
	err = s.checkAllotInfo(c, req.AdminId, adminName, "", allotInfo, req.StoreTaskReviewBasReq, false)
	return err
}

// 统一封装处理H5端审核消息和管理后台审核消息的逻辑
func (s AllotSvc) checkAllotInfo(c ctx.Context, reviewId int64, reviewName, reviewUserId string, allotInfo microScrmModel.DisplayTaskAllotModel, req display.StoreTaskReviewBasReq, isSpecialCheck bool) error {
	failImageId := make([]int64, 0)
	for _, info := range req.ImageList {
		if info.Status == microScrmModel.VmImageFail { // 图片未审核通过
			failImageId = append(failImageId, info.ID)
		}
	}
	err := s.allotDao.UpdateImageData(c, req.EnterpriseId, failImageId)
	if err != nil {
		return err
	}
	err = s.allotDao.UpdateCommentInfo(c, req.EnterpriseId, req.ID, req.AllotID, req.ImageList)
	if err != nil {
		return err
	}
	// 保存操作记录
	nowTime := time.Now().Unix()
	reviewDataLog := microScrmModel.DisplayTaskAllotReviewModel{
		BaseModel: models.BaseModel{
			ID:        app.UID(),
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		DisplayTaskID: req.ID,
		AllotID:       req.AllotID,
		ReviewID:      reviewId,
		ReviewName:    reviewName,
		ReviewUserID:  reviewUserId,
		ReviewType:    req.ReviewType,
	}
	err = s.allotDao.SaveAllotReviewDataLog(c, req.EnterpriseId, reviewDataLog)
	updateMap := make(map[string]interface{}, 0)
	updateMap["review_time"] = nowTime
	updateMap["review_remark"] = req.ReviewRemark
	updateMap["grade"] = req.Grade
	title := ""
	specialTitle := ""
	if isSpecialCheck {
		title = "VM陈列任务待复审提醒"
		updateMap["finish_status"] = microScrmModel.VMFinishStatusDoubleCheck
		err = s.allotDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": req.AllotID,
		}, updateMap)
		uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", req.ID, allotInfo.StoreID)
		memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
			"enterprise_id": req.EnterpriseId,
			"userid":        allotInfo.SpecialUserId,
		})
		if err != nil {
			return err
		}
		content := s.noticeSpecialContent(c, allotInfo, memberInfo.Name, false)
		_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			Title:   title,
			Uri:     uri,
			Party:   nil,
			UserIds: []string{allotInfo.OwnerUserID},
			Content: content,
			AppSign: qyWork.VM,
		})
		if err != nil {
			app.Logger().Error(c).Msgf("VM待复审通知消息发送失败：%+v", err)
			return err
		}
		return nil
	}
	if req.ReviewType == microScrmModel.VMReviewFail {
		title = "VM陈列任务未通过提醒"
		specialTitle = "您审核的VM任务未通过"
		updateMap["finish_status"] = microScrmModel.VMFinishStatusFail
	}
	if req.ReviewType == microScrmModel.VMReviewSuccess {
		title = "VM任务通过啦!快来查看你的分数吧~"
		specialTitle = "您审核的VM任务通过啦!"
		updateMap["finish_status"] = microScrmModel.VMFinishStatusSuccess
	}
	if reviewUserId == "" && allotInfo.SpecialUserId != "" && allotInfo.SpecialCheck == microScrmModel.VMSpecialCheckWait { // 如果管理后台直接审核,special审核人员还未审核的时候，需要将special审核人员清空
		updateMap["special_user_id"] = ""
	}
	err = s.allotDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": req.AllotID,
	}, updateMap)
	go func() {
		uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", req.ID, allotInfo.StoreID)
		content := s.noticeContent(c, req.EnterpriseId, req.ID)
		_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			Title:   title,
			Uri:     uri,
			Party:   []string{strconv.Itoa(allotInfo.StoreID)},
			UserIds: nil,
			Content: content,
			AppSign: qyWork.VM,
		})
		if err != nil {
			app.Logger().Error(c).Msgf("审核通知消息发送失败：%+v", err)
		}

		if allotInfo.SpecialUserId != "" && allotInfo.SpecialCheck == microScrmModel.VMSpecialCheckYes {
			uri = fmt.Sprintf("display-task-preview?id=%d&store_id=%d&is_vm=1", req.ID, allotInfo.StoreID)
			content = s.noticeSpecialContent(c, allotInfo, "", true)
			_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: req.EnterpriseId,
				},
				Title:   specialTitle,
				Uri:     uri,
				Party:   nil,
				UserIds: []string{allotInfo.SpecialUserId},
				Content: content,
				AppSign: qyWork.VM,
			})
			if err != nil {
				app.Logger().Error(c).Msgf("special审核通知消息发送失败：%+v", err)
			}
		}
	}()
	go func() {
		if req.ReviewType == microScrmModel.VMReviewSuccess {
			subInfo, total, err := s.allotDao.GetSubmitInfo(c, req.EnterpriseId, allotInfo.ID)
			if err != nil {
				return
			}
			resSubmit := "否"
			if total > 1 {
				resSubmit = "是"
			}
			_, err = microwork.NewShenCe().CompleteDisplayTask(c, req.EnterpriseId, subInfo.MemberUserID, allotInfo.DisplayName, resSubmit)
			if err != nil {
				return
			}
		}
	}()

	return err
}

// GetNoticeSetting 获取自动提醒的配置
func (s AllotSvc) GetNoticeSetting(c ctx.Context, req request.EnterpriseBaseReq) (uint, error) {
	info, err := microBossDao.NewEnterpriseConfig().GetInfo(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"deleted_at":    0,
		"key":           "vm-notice-setting",
	})
	if err != nil {
		return 0, err
	}
	return info.Value, err
}

// AllotSendNotice 管理后台一键提醒
func (s AllotSvc) AllotSendNotice(c ctx.Context, req display.SearchAllotBaseReq) error {
	list, err := s.allotDao.GetAllotSendNoticeList(c, req)
	if err != nil {
		return err
	}
	waitMap := make(map[int64][]display.AllotSendNoticeInfo, 0)
	failMap := make(map[int64][]display.AllotSendNoticeInfo, 0)
	for _, info := range list {
		if info.FinishStatus == microScrmModel.VMFinishStatusIng {
			waitMap[info.DisplayTaskId] = append(waitMap[info.DisplayTaskId], info)
		}
		if info.FinishStatus == microScrmModel.VMFinishStatusFail {
			failMap[info.DisplayTaskId] = append(failMap[info.DisplayTaskId], info)
		}
	}
	go func() {
		if len(waitMap) == 0 {
			return
		}
		title := "VM陈列任务待提交提醒"
		for k, item := range waitMap {
			content := s.noticeContent(c, req.EnterpriseId, k)
			for _, info := range item {
				storeId := strconv.Itoa(info.StoreId)
				uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", info.DisplayTaskId, info.StoreId)
				_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
					EnterpriseBaseReq: request.EnterpriseBaseReq{
						EnterpriseId: req.EnterpriseId,
					},
					Title:   title,
					Uri:     uri,
					Party:   []string{storeId},
					UserIds: nil,
					Content: content,
					AppSign: qyWork.VM,
				})
			}

		}
	}()

	go func() {
		if len(failMap) == 0 {
			return
		}
		title := "VM陈列任务未通过提醒"
		for k, item := range failMap {
			content := s.noticeContent(c, req.EnterpriseId, k)
			for _, info := range item {
				storeId := strconv.Itoa(info.StoreId)
				uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", info.DisplayTaskId, info.StoreId)
				_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
					EnterpriseBaseReq: request.EnterpriseBaseReq{
						EnterpriseId: req.EnterpriseId,
					},
					Title:   title,
					Uri:     uri,
					Party:   []string{storeId},
					UserIds: nil,
					Content: content,
					AppSign: qyWork.VM,
				})
			}
		}
	}()

	return err
}

// 获取整理消息通知的模板内容
func (s AllotSvc) noticeContent(c ctx.Context, enterpriseId, displayId int64) []map[string]string {
	content := make([]map[string]string, 0)
	taskInfo, _ := display2.NewTaskDao().GetOne(c, enterpriseId, map[string]interface{}{
		"id":         displayId,
		"deleted_at": 0,
	})
	content = append(content, map[string]string{
		"keyname": "陈列任务名称",
		"value":   taskInfo.DisplayName,
	})
	startTime := time.Unix(int64(taskInfo.StartTime), 0).Format("2006/01/02")
	endTime := time.Unix(int64(taskInfo.EndTime), 0).Format("2006/01/02")
	content = append(content, map[string]string{
		"keyname": "任务执行时间",
		"value":   startTime + " - " + endTime,
	})
	return content
}

func (s AllotSvc) noticeSpecialContent(c ctx.Context, allotInfo microScrmModel.DisplayTaskAllotModel, memberName string, isFinish bool) []map[string]string {
	content := make([]map[string]string, 0)

	content = append(content, map[string]string{
		"keyname": "陈列任务名称",
		"value":   allotInfo.DisplayName,
	})

	content = append(content, map[string]string{
		"keyname": "提交门店",
		"value":   allotInfo.StoreName,
	})

	if !isFinish {
		content = append(content, map[string]string{
			"keyname": "初审人员",
			"value":   memberName,
		})
	}
	submitTime := helper.GetFormatDateByParams(allotInfo.SubmitTime)
	content = append(content, map[string]string{
		"keyname": "提交时间",
		"value":   submitTime,
	})
	return content
}

func (s AllotSvc) AllotImageList(c ctx.Context, req display.SaveImageReq) ([]string, error) {
	return s.allotDao.AllotImageList(c, req)
}

func (s AllotSvc) AllotNoticeNum(c ctx.Context, req display.SearchAllotBaseReq) (int64, error) {
	// 如果vm任务完成情况筛选不是未提交和未通过，则提醒数为0
	if req.FinishStatus >= 0 && !helper.InArray(req.FinishStatus, []int32{microScrmModel.VMFinishStatusIng, microScrmModel.VMFinishStatusFail}) {
		return 0, nil
	}
	return s.allotDao.GetAllotNoticeNum(c, req)
}

// CheckIsVm  CheckIsVm 判断导购是否还是VM身份
func (s AllotSvc) CheckIsVm(c ctx.Context, enterpriseId, memberId int64) bool {
	list, err := display2.NewAllotDao().VmList(c, enterpriseId, map[string]interface{}{
		"is_check":   1,
		"deleted_at": 0,
	})
	if err != nil {
		return false
	}
	return helper.InArray(memberId, list)
}

func (s AllotSvc) CosAvatar(c ctx.Context, filePath string) (string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)
	md5Name := helper.Md5V(fileName) + helper.Int64ToString(time.Now().UnixNano())
	// 压缩文件地址
	fileLocalPath := fmt.Sprintf("storage/upload/%s%s", md5Name, fileExt)
	err := helper.DownloadFileToLocal(filePath, fileLocalPath)
	defer func() { _ = os.Remove(filePath) }()
	if err != nil {
		return "", err
	}
	cosUrl, err := helper.UploadFileToCos(c, fileLocalPath, false)
	if err != nil {
		return "", err
	}

	return cosUrl, nil
}

func (s AllotSvc) SaveStoreImage(c ctx.Context, req display.SaveStoreImageReq) (int64, error) {
	id := app.UID()
	nowTime := time.Now().Unix()
	imageData := microScrmModel.DisplayTaskAllotDataModel{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		DisplayTaskID: req.DisplayTaskId,
		AllotID:       req.AllotId,
		Identifier:    req.Identifier,
		ImageURL:      req.Url,
		Status:        0,
	}
	err := s.allotDao.SaveAllotDataInfo(c, req.EnterpriseId, imageData)
	return id, err
}

func (s AllotSvc) DelImage(c ctx.Context, req display.DelImageReq) error {
	return s.allotDao.DeleteImageList(c, req.EnterpriseId, []int64{req.Id})
}

func (s AllotSvc) SaveStoreComment(c ctx.Context, req display.SaveStoreCommentReq) (int64, error) {
	var err error
	id := req.ID
	if id == 0 {
		id = app.UID()
		nowTime := time.Now().Unix()
		commentData := microScrmModel.DisplayTaskAllotRemarkModel{
			BaseModel: models.BaseModel{
				ID:        id,
				CreatedAt: nowTime,
				UpdatedAt: nowTime,
			},
			DisplayTaskID: req.DisplayTaskId,
			AllotID:       req.AllotId,
			AllotDataID:   req.AllotDataId,
			MemberID:      req.MemberId,
			MemberUserID:  req.MemberUserId,
			MemberName:    req.MemberName,
			Avatar:        req.Avatar,
			Content:       helper.MarshalToString(req.Content),
			StoreCode:     req.StoreCode,
			CommentTime:   req.CommentTime,
			IsVm:          req.IsVm,
			ReplyName:     req.ReplyName,
		}
		err = s.allotDao.SaveAllotComment(c, req.EnterpriseId, commentData)
	} else {
		content := helper.MarshalToString(req.Content)
		err = s.allotDao.UpdateAllotCommentInfo(c, req.EnterpriseId, req.ID, content)
	}
	return id, err
}

func (s AllotSvc) DelComment(c ctx.Context, req display.DelImageReq) error {
	return s.allotDao.DelAllotComment(c, req.EnterpriseId, req.Id)
}

func (s AllotSvc) VmDashboard(c ctx.Context, req display.DashBoard) (display.DashBoardInfo, error) {
	var res display.DashBoardInfo
	ids, err := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{
			EnterpriseId: req.EnterpriseId,
		},
		AreaIds: req.AreaIds,
	})
	if err != nil || len(ids) == 0 {
		return res, err
	}

	countList, err := s.allotDao.GetVmDashboardCount(c, req.EnterpriseId, map[string]interface{}{
		"display_task_id": req.Ids,
		"store_id":        ids,
		"deleted_at":      0,
	})
	if err != nil {
		return res, err
	}
	waitMap := make(map[int]int, 0)
	checkMap := make(map[int]int, 0)
	doubleCheckMap := make(map[int]int, 0)
	failMap := make(map[int]int, 0)
	successMap := make(map[int]int, 0)
	for _, status := range countList {
		switch status.FinishStatus {
		case microScrmModel.VMFinishStatusIng:
			if _, ok := waitMap[status.StoreId]; !ok {
				waitMap[status.StoreId] = 1
			}
		case microScrmModel.VMFinishStatusCheck:
			if _, ok := checkMap[status.StoreId]; !ok {
				checkMap[status.StoreId] = 1
			}
		case microScrmModel.VMFinishStatusFail:
			if _, ok := failMap[status.StoreId]; !ok {
				failMap[status.StoreId] = 1
			}
		case microScrmModel.VMFinishStatusSuccess:
			if _, ok := successMap[status.StoreId]; !ok {
				successMap[status.StoreId] = 1
			}
		case microScrmModel.VMFinishStatusDoubleCheck:
			if _, ok := doubleCheckMap[status.StoreId]; !ok {
				doubleCheckMap[status.StoreId] = 1
			}
		}
	}
	res.WaitNum = int64(len(waitMap))
	res.CheckNum = int64(len(checkMap))
	res.FailNum = int64(len(failMap))
	res.SuccessNum = int64(len(successMap))
	res.DoubleCheckNum = int64(len(doubleCheckMap))
	res.AllNum = res.WaitNum + res.CheckNum + res.FailNum + res.SuccessNum + res.DoubleCheckNum
	res.WaitRate = "0%"
	res.CheckRate = "0%"
	res.FailRate = "0%"
	res.SuccessRate = "0%"
	res.DoubleCheckRate = "0%"
	if res.AllNum > 0 {
		allNum := decimal.NewFromInt(res.AllNum)
		waitRate, _ := decimal.NewFromInt(res.WaitNum).DivRound(allNum, 4).Mul(decimal.NewFromInt(100)).Float64()
		res.WaitRate = fmt.Sprintf("%v%%", waitRate)
		checkRate, _ := decimal.NewFromInt(res.CheckNum).DivRound(allNum, 4).Mul(decimal.NewFromInt(100)).Float64()
		res.CheckRate = fmt.Sprintf("%v%%", checkRate)
		failRate, _ := decimal.NewFromInt(res.FailNum).DivRound(allNum, 4).Mul(decimal.NewFromInt(100)).Float64()
		res.FailRate = fmt.Sprintf("%v%%", failRate)
		successRate, _ := decimal.NewFromInt(res.SuccessNum).DivRound(allNum, 4).Mul(decimal.NewFromInt(100)).Float64()
		res.SuccessRate = fmt.Sprintf("%v%%", successRate)
		doubleCheckRate, _ := decimal.NewFromInt(res.DoubleCheckNum).DivRound(allNum, 4).Mul(decimal.NewFromInt(100)).Float64()
		res.DoubleCheckRate = fmt.Sprintf("%v%%", doubleCheckRate)
	}
	return res, err
}

func (s AllotSvc) FailList(c ctx.Context, req display.MemberTaskReq) (display.VmFailList, error) {
	var res display.VmFailList
	list := make([]display.VmFailInfo, 0)
	res.List = list
	ids, err := dashboardDao.NewManager().GetManagerAuthAreaIds(c, dashboardReq.ManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{
			EnterpriseId: req.EnterpriseId,
		},
		MemberId: req.MemberId,
	})
	if err != nil || len(ids) == 0 {
		return res, err
	}
	storeIds, err := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{
			EnterpriseId: req.EnterpriseId,
		},
		AreaIds: ids,
	})
	if err != nil || len(storeIds) == 0 {
		return res, err
	}
	task, err := display2.NewTaskDao().GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		return res, err
	}
	res.ID = task.ID
	res.DisplayName = task.DisplayName
	many, err := s.allotDao.GetMany(c, req.EnterpriseId, map[string]interface{}{
		"display_task_id": task.ID,
		"store_id":        storeIds,
		"finish_status":   microScrmModel.VMFinishStatusFail,
	})
	if err != nil {
		return res, err
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, req.EnterpriseId)
	for _, v := range many {
		info := display.VmFailInfo{
			AllotId:    v.ID,
			StoreName:  v.StoreName,
			StoreId:    v.StoreID,
			ReviewTime: v.ReviewTime,
		}
		if v, ok := storeNameMap[info.StoreId]; ok {
			info.StoreName = v
		}
		list = append(list, info)
	}
	res.List = list
	return res, err
}

func (s AllotSvc) GetDashBoardCount(c ctx.Context, enterpriseId int64, taskIds []int64, finishStatus int) ([]display.StoreCount, error) {
	storeCount, err := s.allotDao.GetStoreCount(c, enterpriseId, finishStatus, taskIds)
	if err != nil {
		return storeCount, err
	}
	return storeCount, err
}

// 获取随机的special的人员,special所在门店不能审核自己门店
func (s AllotSvc) getSpecialInfo(specialInfo []display.SpecialInfo, storeBaseList []display.StoreBaseInfo) []display.StoreBaseInfo {
	var result []display.StoreBaseInfo
	for _, storeBase := range storeBaseList {
		sort.Slice(specialInfo, func(i, j int) bool {
			return specialInfo[i].Num < specialInfo[j].Num
		})
		for k, info := range specialInfo {
			if storeBase.StoreCode == info.StoreCode {
				continue
			}
			value := storeBase
			value.SpecialUserId = info.UserId
			result = append(result, value)
			specialInfo[k].Num++
			break
		}
	}
	return result
}

func (s AllotSvc) SpecialSaveStoreTask(c ctx.Context, req display.StoreTaskReviewReq) error {
	allotInfo, err := s.allotDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.AllotID,
	})
	if err != nil {
		return err
	}
	if allotInfo.ID == 0 || allotInfo.FinishStatus >= microScrmModel.VMFinishStatusFail {
		return errors.New("门店VM数据状态异常,已被审核了")
	}
	isSpecialCheck := false
	if allotInfo.SpecialUserId != "" {
		if allotInfo.SpecialCheck == microScrmModel.VMSpecialCheckWait {
			isSpecialCheck = true
		}
	}
	if !isSpecialCheck {
		return errors.New("该门店不能被初审")
	}
	if allotInfo.SpecialUserId != req.MemberUserId {
		return errors.New("该门店的Special人员不是您")
	}

	err = s.specialCheckAllotInfo(c, req.MemberId, req.MemberName, req.MemberUserId, allotInfo, req.StoreTaskReviewBasReq)
	return err
}

func (s AllotSvc) specialCheckAllotInfo(c ctx.Context, reviewId int64, reviewName, reviewUserId string, allotInfo microScrmModel.DisplayTaskAllotModel, req display.StoreTaskReviewBasReq) error {
	failImageId := make([]int64, 0)
	updateImageData := make(map[string]interface{})
	if req.ReviewType == microScrmModel.VMSpecialDraft { // special保存草稿
		successImageId := make([]int64, 0)
		for _, info := range req.ImageList {
			if info.Status == microScrmModel.VmImageFail { // 图片未审核通过
				failImageId = append(failImageId, info.ID)
			}
			if info.Status == microScrmModel.VmImageSuccess {
				successImageId = append(successImageId, info.ID)
			}
		}
		updateImageData = map[string]interface{}{
			"draft_status": microScrmModel.VmImageFail,
			"is_draft":     1,
		}
		successUpdateImageData := map[string]interface{}{
			"draft_status": microScrmModel.VmImageSuccess,
			"is_draft":     1,
		}
		_ = s.allotDao.SpecialUpdateImageData(c, req.EnterpriseId, successImageId, successUpdateImageData)
	} else {
		for _, info := range req.ImageList {
			if info.Status == microScrmModel.VmImageFail { // 图片未审核通过
				failImageId = append(failImageId, info.ID)
			}
		}
		updateImageData = map[string]interface{}{
			"status":   microScrmModel.VmImageFail,
			"is_draft": 0,
		}
	}
	err := s.allotDao.SpecialUpdateImageData(c, req.EnterpriseId, failImageId, updateImageData)
	if err != nil {
		return err
	}
	err = s.allotDao.UpdateCommentInfo(c, req.EnterpriseId, req.ID, req.AllotID, req.ImageList)
	if err != nil {
		return err
	}
	updateMap := make(map[string]interface{})
	if req.ReviewType == microScrmModel.VMSpecialDraft { // 保存草稿
		updateMap["special_remark"] = req.ReviewRemark
		err = s.allotDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": req.AllotID,
		}, updateMap)
		return err
	}
	// 保存操作记录
	nowTime := time.Now().Unix()
	reviewDataLog := microScrmModel.DisplayTaskAllotReviewModel{
		BaseModel: models.BaseModel{
			ID:        app.UID(),
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		DisplayTaskID: req.ID,
		AllotID:       req.AllotID,
		ReviewID:      reviewId,
		ReviewName:    reviewName,
		ReviewUserID:  reviewUserId,
		ReviewType:    req.ReviewType,
	}
	err = s.allotDao.SaveAllotReviewDataLog(c, req.EnterpriseId, reviewDataLog)
	updateMap["review_time"] = nowTime
	updateMap["review_remark"] = req.ReviewRemark
	updateMap["grade"] = req.Grade
	updateMap["finish_status"] = microScrmModel.VMFinishStatusDoubleCheck
	updateMap["special_check"] = microScrmModel.VMSpecialCheckYes
	err = s.allotDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": req.AllotID,
	}, updateMap)
	title := "VM陈列任务待复审提醒"
	uri := fmt.Sprintf("display-task-audit?id=%d&store_id=%d", req.ID, allotInfo.StoreID)
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"userid":        allotInfo.SpecialUserId,
	})
	if err != nil {
		return err
	}
	content := s.noticeSpecialContent(c, allotInfo, memberInfo.Name, false)
	_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
		EnterpriseBaseReq: request.EnterpriseBaseReq{
			EnterpriseId: req.EnterpriseId,
		},
		Title:   title,
		Uri:     uri,
		Party:   nil,
		UserIds: []string{allotInfo.OwnerUserID},
		Content: content,
		AppSign: qyWork.VM,
	})
	if err != nil {
		app.Logger().Error(c).Msgf("VM待复审通知消息发送失败：%+v", err)
		return err
	}
	return nil
}

func (s AllotSvc) ReviewExport(c ctx.Context, req display.TaskBaseReq) (string, error) {
	task, err := display2.NewTaskDao().GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		return "", err
	}
	if task.ReviewType == 0 {
		return "", err
	}
	reviewName := helper.RuleThree(task.ReviewType == microScrmModel.VMTaskReviewTypeRand, "随机平均分配", "定向分配")
	fileName := task.DisplayName + "_" + reviewName + "_" + time.Now().Format("20060102150405") + ".xlsx"
	allotList, err := s.allotDao.GetMany(c, req.EnterpriseId, "display_task_id=? and special_user_id <> ? and deleted_at=0", req.ID, "")
	if err != nil {
		return fileName, err
	}
	specialUserIds := make([]string, 0)
	for _, v := range allotList {
		specialUserIds = append(specialUserIds, v.SpecialUserId)
	}
	reviewLog, err := s.allotDao.GetAllotReviewDataLog(c, req.EnterpriseId, map[string]interface{}{
		"display_task_id": req.ID,
		"review_user_id":  specialUserIds,
	})
	if err != nil {
		return fileName, err
	}
	reviewMap := make(map[string]string, 0)
	for _, v := range reviewLog {
		reviewMap[v.ReviewUserID] = helper.GetFormatDateByParams(v.CreatedAt)
	}
	memberList, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"userid":        specialUserIds,
		"deleted_at":    0,
	})
	if err != nil {
		return fileName, err
	}
	memberMap := make(map[string]string, 0)
	for _, v := range memberList {
		memberMap[v.Userid] = v.Name
	}
	excelFile := excelize.NewFile()
	sheetName := reviewName
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	defer excelFile.Close()
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return fileName, err
	}
	err = excelFile.MergeCell(sheetName, "A1", "E1")
	if err != nil {
		return fileName, err
	}
	row := make([]any, 0)
	row = []any{
		"员工姓名", "*导购工号", "*审核门店(code)", "审核状态", "审核时间",
	}
	err = stream.SetRow("A1", []any{
		"数据统计截止时间：" + helper.GetFormatDateByParams(time.Now().Unix()),
	})
	err = stream.SetRow("A2", row)
	startRow := int64(3)
	for _, item := range allotList {
		rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
		memberName := ""
		if value, ok := memberMap[item.SpecialUserId]; ok {
			memberName = value
		}
		statusStr := "未审核"
		reviewTime := "--"
		if item.SpecialCheck == microScrmModel.VMSpecialCheckYes {
			statusStr = "已审核"
			reviewTime = reviewMap[item.SpecialUserId]
		}
		row = []any{memberName, item.SpecialUserId, item.StoreCode, statusStr, reviewTime}
		err = stream.SetRow("A"+rowStartStr, row)
		atomic.AddInt64(&startRow, 1)
	}
	if err != nil {
		return fileName, err
	}
	err = stream.Flush()
	if err != nil {
		return fileName, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return fileName, err
	}
	return fileName, err
}

// CheckSpecialInfo 判断是否是special人员
func (s AllotSvc) CheckSpecialInfo(c ctx.Context, enterpriseId, memberId int64) bool {
	info, err := display2.NewTaskDao().GetSpecialInfo(c, enterpriseId, map[string]interface{}{
		"member_id":  memberId,
		"is_special": 1,
	})
	if err != nil || info.ID == 0 {
		return false
	}
	return true
}

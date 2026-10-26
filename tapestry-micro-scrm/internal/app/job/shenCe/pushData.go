package shenCeJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	sourceDao "micro-scrm/internal/app/dao/source"
	appJob "micro-scrm/internal/app/job"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PushData struct {
	name string

	shenCeClient microwork.ShenCe
	startTime    int64
	endTime      int64
}

func NewPushData() *PushData {
	var startTime, endTime int64 = helper.GetBetweenByTime(1)
	return &PushData{
		name:         "神策数据推送",
		shenCeClient: microwork.NewShenCe(),
		startTime:    startTime,
		endTime:      endTime,
	}
}

func (job *PushData) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	wg := threading.NewRoutineGroup()

	for _, enterpriseId := range enterpriseIds {
		id := enterpriseId
		wg.RunSafe(func() {
			err = job.handle(c, id)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *PushData) handle(c ctx.Context, enterpriseId int64) error {

	// 外部联系人
	if err := job.userCorpExternal(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("外部联系人数据推送 error : %v")
	}

	// 计划任务
	if err := job.task(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("计划任务数据推送 error : %v")
	}

	// 素材库
	if err := job.source(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("素材库数据推送 error : %v")
	}

	// 企微对外收款记录
	if err := job.qyWechatOrder(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("企微对外收款记录数据推送 error : %v")
	}

	// 邀约模板
	if err := job.invite(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("邀约模板数据推送 error : %v")
	}

	// 员工信息
	if err := job.members(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("员工信息数据推送 error : %v")
	}

	if err := job.onjobInherit(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("在职分配数据推送 error : %v")
	}
	return nil
}

// 外部联系人数据推送
func (job *PushData) userCorpExternal(c ctx.Context, enterpriseId int64) error {
	wxUserCorpExternals, err := microScrmDao.NewWxUserCorpExternalCk().GetMany(c, enterpriseId, "unionid <> '' and updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}

	for _, user := range wxUserCorpExternals {
		properties := map[string]any{
			"unionid":       user.Unionid,
			"nick":          user.Nick,
			"sex":           user.Sex,
			"avatar":        user.Avatar,
			"union_code":    user.UnionCode,
			"clt_type":      user.CltType,
			"customer_code": user.CustomerCode,
			"card_level":    user.CardLevel,
			"created_at":    helper.GetFormatDateByParams(user.CreatedAt),
		}
		_, err = job.shenCeClient.External(c, enterpriseId, user.Openid, properties)
	}

	return nil
}

// 计划任务数据推送
func (job *PushData) task(c ctx.Context, enterpriseId int64) error {
	// 任务信息维度
	planTasks, err := microScrmDao.NewPlanTaskDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, plan := range planTasks {
		properties := map[string]any{
			"name":       plan.Name,
			"remark":     plan.Remark,
			"start_time": helper.GetFormatDateByParams(int64(plan.Plan.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(plan.Plan.EndTime)),
		}
		_, _ = job.shenCeClient.Task(c, enterpriseId, helper.Int64ToString(plan.ID), properties)

		if plan.Status == 2 {
			// 任务过期
			_, _ = job.shenCeClient.TaskExpire(c, enterpriseId, helper.Int64ToString(plan.ID), plan.ID)
		}
	}

	// 完成任务
	// taskAllotLogs, err := microScrmDao.NewPlanTaskAllotLogDao().GetMany(c, enterpriseId, "end_time >= ? and end_time <= ? and is_finished = 1", job.startTime, job.endTime)
	// if err != nil {
	// 	return err
	// }
	// for _, task := range taskAllotLogs {
	// 	_, _ = job.shenCeClient.TaskComplete(c, enterpriseId, task.UserId, task.TaskId)
	// }

	// 收到任务
	taskReceives, err := microScrmDao.NewPlanTaskAllotLogDao().GetTaskIdGroupNum(c, enterpriseId, job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, receive := range taskReceives {
		// 获取子任务名称
		planTaskRes, err := microScrmDao.NewPlanTaskDao().GetOne(c, enterpriseId, map[string]any{"id": receive.TaskId})
		if err != nil {
			return err
		}
		_, _ = job.shenCeClient.TaskReceive(c, enterpriseId, helper.Int64ToString(receive.PlanId), map[string]any{
			"task_id":        receive.TaskId,
			"task_name":      planTaskRes.Name,
			"task_coverage":  receive.GroupMsgCount,
			"task_starttime": helper.GetFormatDateByParams(receive.CreatedAt),
		})
	}

	return nil
}

// 素材库数据推送
func (job *PushData) source(c ctx.Context, enterpriseId int64) error {
	// 视频
	sourceVideo, err := sourceDao.NewVideoDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, video := range sourceVideo {
		properties := map[string]any{
			"name":       video.Name,
			"type":       source.TypeVideo,
			"content":    video.URL,
			"status":     video.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(video.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(video.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(video.ID), properties)
	}

	// 文件
	sourceFile, err := sourceDao.NewFileDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, file := range sourceFile {
		properties := map[string]any{
			"name":       file.Name,
			"type":       source.TypeFile,
			"content":    file.URL,
			"status":     file.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(file.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(file.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(file.ID), properties)
	}

	// 图片
	sourceImage, err := sourceDao.NewImageDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, image := range sourceImage {
		properties := map[string]any{
			"name":       image.Name,
			"type":       source.TypeImg,
			"content":    image.URL,
			"status":     image.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(image.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(image.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(image.ID), properties)
	}

	// 小程序
	sourceMiniProgram, err := sourceDao.NewMiniProgramDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, mini := range sourceMiniProgram {
		properties := map[string]any{
			"name":       mini.Title,
			"type":       source.TypeMiniProgram,
			"content":    mini.CoverURL,
			"status":     mini.Source.Status,
			"start_time": mini.Source.StartTime,
			"end_time":   mini.Source.EndTime,
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(mini.ID), properties)
	}

	// 链接
	sourceLink, err := sourceDao.NewLinkDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, link := range sourceLink {
		properties := map[string]any{
			"name":       link.Title,
			"type":       source.TypeLink,
			"content":    link.Content,
			"status":     link.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(link.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(link.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(link.ID), properties)
	}

	// 文本
	sourceText, err := sourceDao.NewTextDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, text := range sourceText {
		properties := map[string]any{
			"name":       "",
			"type":       source.TypeText,
			"content":    text.Content,
			"status":     text.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(text.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(text.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(text.ID), properties)
	}

	// 朋友圈
	sourceFriend, err := sourceDao.NewFriendDao().GetManyPreload(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, friend := range sourceFriend {
		properties := map[string]any{
			"name":       friend.Name,
			"type":       source.TypeFriend,
			"content":    helper.RuleThree(friend.TextContent != "", friend.TextContent, friend.OriginalURL),
			"status":     friend.Source.Status,
			"start_time": helper.GetFormatDateByParams(int64(friend.Source.StartTime)),
			"end_time":   helper.GetFormatDateByParams(int64(friend.Source.EndTime)),
		}
		_, _ = job.shenCeClient.Material(c, enterpriseId, helper.Int64ToString(friend.ID), properties)
	}

	return nil
}

// 邀约数据推送
func (job *PushData) invite(c ctx.Context, enterpriseId int64) error {
	// 邀约模板
	inviteTemplates, err := microScrmDao.NewInviteTemplateDao().GetMany(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, invite := range inviteTemplates {
		properties := map[string]any{
			"title_name":         invite.TitleName,
			"reserve_service_id": invite.ReserveServiceID,
			"start_time":         helper.GetFormatDateByParams(int64(invite.StartTime)),
			"end_time":           helper.GetFormatDateByParams(int64(invite.EndTime)),
			"share_title":        invite.ShareTitle,
			"content":            invite.Content,
		}
		_, err = job.shenCeClient.InviteTemplate(c, enterpriseId, helper.Int64ToString(invite.ID), properties)
	}

	// 邀约服务
	reserveServices, err := microScrmDao.NewReserveServiceDao().GetMany(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, reserve := range reserveServices {
		properties := map[string]any{
			"reserve_name": reserve.ReserveName,
			"service_type": reserve.ServiceType,
		}
		if reserve.ServiceType == 0 {
			properties["start_time"] = helper.GetFormatDateByParams(int64(reserve.StartTime))
			properties["end_time"] = helper.GetFormatDateByParams(int64(reserve.EndTime))
		}
		_, err = job.shenCeClient.ReserveService(c, enterpriseId, helper.Int64ToString(reserve.ID), properties)
	}

	// 邀约成功
	reserveSuccess, err := microScrmDao.NewReservesOrder().GetMany(c, enterpriseId, "updated_at >= ? and updated_at <= ? and status = 1", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, reserve := range reserveSuccess {
		serviceDate := time.Unix(int64(reserve.ReserveTime), 0).Format("2006/01/02")
		serviceTimeslot := time.Unix(int64(reserve.ReserveTime), 0).Format("15:04:05")
		// 获取当前员工的门店编码
		storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
			"enterprise_id":        enterpriseId,
			"wechat_department_id": reserve.StoreID,
		})

		properties := map[string]any{
			"appointment_id":   reserve.ReserveNumber,
			"customer_id":      reserve.UserExternalID,
			"service_type":     reserve.ReserveService.ReserveName,
			"service_date":     serviceDate,
			"service_timeslot": serviceTimeslot,
			"store_name":       storeDepartmentInfo.StoreName,
			"store_code":       storeDepartmentInfo.ThirdStoreCode,
			"customer_type":    reserve.CustomerStatus,
			"is_instore":       reserve.IsStore,
			"order_amount":     reserve.ReserveOrderPrice,
			"order_date":       helper.GetFormatDateByParams(int64(reserve.ReserveOrderTime)),
		}
		_, err = job.shenCeClient.InviteSuccess(c, enterpriseId, reserve.MemberUserID, properties)
	}

	return nil
}

// 企微对外收款记录
func (job *PushData) qyWechatOrder(c ctx.Context, enterpriseId int64) error {
	payOrders, err := microScrmDao.NewExternalPayOrderDao().GetMany(c, enterpriseId, "updated_at >= ? and updated_at <= ?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, payOrder := range payOrders {
		properties := map[string]any{
			"transaction_id":  payOrder.TransactionID,
			"pay_time":        helper.GetFormatDateByParams(int64(payOrder.PayTime)),
			"external_userid": payOrder.ExternalUserid,
			"total_fee":       payOrder.TotalFee,
			"payee_userid":    payOrder.PayeeUserid,
			"remark":          payOrder.Remark,
			"delivery_type":   payOrder.DeliveryType,
			"tickets":         payOrder.Tickets,
		}

		_, err = job.shenCeClient.ExternalPay(c, enterpriseId, payOrder.PayeeUserid, properties)
	}
	return nil
}

// 员工信息
func (job *PushData) members(c ctx.Context, enterpriseId int64) error {
	// 获取需要推送给神策的当前企业的员工
	members, err := membersDao.NewMembers().GetMemberForShenCe(c, enterpriseId, job.startTime, job.endTime)
	if err != nil {
		return err
	}

	shenCeClient := microwork.NewShenCe()
	areaList, err := dashboardDao.NewArea().GetAll(c, enterpriseId)
	if err != nil {
		return err
	}
	areaMap := make(map[int64]dashboardModel.DashboardAreaModel, 0)
	thirdAreaMap := make(map[int][]string, 0)
	for _, model := range areaList {
		if model.ThirdID > 0 {
			thirdAreaMap[model.ThirdID] = strings.Split(model.Path, ",")
		} else {
			areaMap[model.ID] = model
		}
	}
	thirdAreaInfoMap := make(map[int][]dashboardModel.DashboardAreaModel, 0)
	for k, pathSlice := range thirdAreaMap {
		areaSlice := make([]dashboardModel.DashboardAreaModel, 0)
		for _, v := range pathSlice {
			index, _ := strconv.ParseInt(v, 10, 64)
			areaSlice = append(areaSlice, areaMap[index])
		}
		sort.Slice(areaSlice, func(i, j int) bool {
			return areaSlice[i].Type < areaSlice[j].Type
		})
		thirdAreaInfoMap[k] = areaSlice

	}
	var memberIds []int64
	for _, member := range members {
		memberIds = append(memberIds, member.ID)
	}
	deptMemberList, err := departmentsDao.NewDepartments().GetDepartmentsMemberByMemberIds(c, enterpriseId, memberIds)
	if err != nil {
		return err
	}
	var thirdIds []int
	var memberDeptMap = make(map[int64]departmentsModels.DepartmentsMembersInfo)
	for _, info := range deptMemberList {
		memberDeptMap[info.MemberId] = info
		thirdIds = append(thirdIds, info.ThirdId)
	}
	if err != nil {
		return err
	}
	storeDeptList, _ := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, enterpriseId, thirdIds)

	var storeDeptMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, info := range storeDeptList {
		storeDeptMap[int(info.WechatDepartmentID)] = info
	}
	for _, member := range members {
		var storeName, storeCode, wecomDepartment string
		var thirdId int
		// 获取当前员工的部门关系
		if memberDept, ok := memberDeptMap[member.ID]; ok {
			thirdId = memberDept.ThirdId
			wecomDepartment = memberDept.DepartmentName
			if storeDept, ok1 := storeDeptMap[memberDept.ThirdId]; ok1 {
				storeName = storeDept.StoreName
				storeCode = storeDept.ThirdStoreCode
			}
		}

		typeStr := "sa"
		if member.Type == 4 {
			typeStr = "sm"
		}
		areaInfo := thirdAreaInfoMap[thirdId]
		var storeChannel, storeRegion, storeSubRegion, svRole string
		for _, model := range areaInfo {
			switch model.Type {
			case dashboardModel.AreaSource:
				storeChannel = model.Name
			case dashboardModel.AreaRegion:
				storeRegion = model.Name
			case dashboardModel.AreaArea:
				storeSubRegion = model.Name
			default:
			}
		}
		if member.SVManager > 0 {
			svRole = "Corporate User"
		} else if helper.InArray(member.Position, membersDao.StoreManagerPosition) {
			svRole = "Store Manager"
		} else {
			svRole = "SA"
		}
		_, _ = shenCeClient.ProfileSet(c, microwork.ProfileSetReq{
			EnterpriseId: enterpriseId,
			DistinctId:   member.Userid,
			IsLoginId:    false,
			Profiles: map[string]any{
				"userid":           member.Userid,
				"mobile":           member.Mobile,
				"position":         member.Position,
				"type":             typeStr,
				"store_name":       storeName,
				"store_code":       storeCode,
				"wecom_department": wecomDepartment,
				"gender":           member.Gender,
				"avatar":           member.Avatar,
				"name":             member.Name,
				"store_channel":    storeChannel,
				"store_region":     storeRegion,
				"store_subregion":  storeSubRegion,
				"sv_department":    member.SVDepartment,
				"sv_role":          svRole,
			},
		})
	}

	return err
}

// 在职分配上报
func (job *PushData) onjobInherit(c ctx.Context, enterpriseId int64) error {
	// 获取当前企业的所有员工
	queryStr := fmt.Sprintf("enterprise_id = %d and updated_at >= %d and updated_at <= %d", enterpriseId, job.startTime, job.endTime)

	list, err := microScrmDao.NewCustomerOnInheritRecord().GetMany(c, queryStr)
	if err != nil {
		return err
	}
	for _, value := range list {
		assignmentId := helper.Int64ToString(value.InheritID)
		properties := map[string]any{
			"assignment_id":   assignmentId,
			"customer_id":     value.ExternalUserid,
			"handover_userid": value.HandoverUserid,
			"takeover_userid": value.TakeoverUserid,
		}
		_, _ = job.shenCeClient.AssignCustomer(c, enterpriseId, assignmentId, properties)
		if value.Status == 1 {
			properties2 := map[string]any{
				"assignment_id":   assignmentId,
				"customer_id":     value.ExternalUserid,
				"takeover_userid": value.TakeoverUserid,
			}
			_, _ = job.shenCeClient.CustomerAssigned(c, enterpriseId, assignmentId, properties2)
		}

	}

	return err
}

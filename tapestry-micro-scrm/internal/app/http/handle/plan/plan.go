package plan

import (
	"errors"
	"fmt"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	"micro-scrm/internal/app/http/handle"
	plan_task2 "micro-scrm/internal/app/http/middleware/plan_task"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	planReq "micro-scrm/internal/app/http/request/plan"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/plan_task"
	"micro-scrm/internal/pkg/helper"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type PlansHandle struct {
	handle.BaseHandle
}

func NewPlansHandle() PlansHandle {
	return PlansHandle{}
}

func (p PlansHandle) FileDetection(c *gin.Context) {
	var request planReq.FileDetection
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	detection, err := planService.NewPlanTaskImport().FileDetection(p.WarpContext(c), request.EnterpriseID, request.FileName)
	if err != nil {
		p.ResError(c, err)
		return
	}
	var result plan_task.CountInfo
	result.Count = detection
	p.ResSuccess(c, result)
}

func (p PlansHandle) SaveInfo(c *gin.Context) {
	var request planReq.ImportFileReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	var res planReq.ImportFileRes
	var data microScrmModel.TaskImportFileModel
	data.FileName = request.FileName
	data.FileURL = request.FileUrl
	data.TaskType = request.TaskType
	ID, err := planService.NewPlanTaskImport().SaveImportFile(p.WarpContext(c), request.EnterpriseID, data)
	if err != nil {
		p.ResError(c, err)
		return
	}
	res.ID = ID
	p.ResSuccess(c, res)
}

func (p PlansHandle) ConsumerList(c *gin.Context) {
	var request planReq.ConsumerList
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	var result planReq.ConsumerListRes
	if request.Total != 0 {
		result.Total = request.Total
	} else {
		total, err := planService.NewPlanTaskImportDataLog().GetImportDataLogCount(p.WarpContext(c), request.EnterpriseID, request.ConsumerFileID)
		if err != nil {
			p.ResError(c, err)
			return
		}
		result.Total = total
	}

	list, errList := planService.NewPlanTaskImportDataLog().GetImportDataLogList(p.WarpContext(c), request.EnterpriseID, request.ConsumerFileID, request.Page, request.Size)
	if errList != nil {
		p.ResError(c, errList)
		return
	}
	result.List = list
	p.ResSuccess(c, result)
}

// MonthLists 获取任务中心的日历接口
func (p PlansHandle) MonthLists(c *gin.Context) {
	var request planReq.MonthListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	if request.EndTime <= request.StartTime {
		p.ResError(c, errors.New("时间区间错误"))
		return
	}
	var res planReq.MonthListsRes
	var result []planReq.MonthListsData
	var err error
	switch request.Type {
	case 1:
		result, err = planService.NewPlanTask().GetPlanOrderMonthLists(p.WarpContext(c), request)
	case 2:
		result, err = customerSvc.NewSowTask().GetMonthLists(p.WarpContext(c), request)
	case 3:
		result, err = planService.NewPlanTask().GetOrderMonthLists(p.WarpContext(c), request)
	}
	if err != nil {
		p.ResError(c, err)
		return
	}
	res.List = result
	p.ResSuccess(c, res)
	return
}

// Lists 任务中心的主任务列表
func (p PlansHandle) Lists(c *gin.Context) {
	var request planReq.ListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	if request.EndTime <= request.StartTime {
		p.ResError(c, errors.New("时间区间错误"))
		return
	}
	var res planReq.PlanListRes
	result, err := planService.NewPlanTask().GetPlanLists(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	res.Lists = result
	p.ResSuccess(c, res)
}

// Detail 获取主任务详情
func (p PlansHandle) Detail(c *gin.Context) {
	var request planReq.PlanInfo
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTask().GetPlanDetail(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, result)
}

func (p PlansHandle) Tips(c *gin.Context) {
	var request planReq.TipsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	// 获取任务中心的指定时间未完成主任务数
	result, err := planService.NewPlanTask().GetPlanTipsCount(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}

	p.ResSuccess(c, result)
}

// LeaveGroupChat 计划任务的退群顾客列表
func (p PlansHandle) LeaveGroupChat(c *gin.Context) {
	var request planReq.LeaveGroupChatReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	// 获取任务中心的指定时间未完成主任务数
	total, result, err := planService.NewPlanTask().LeaveGroupChat(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}

	p.ResList(c, total, request.Page, request.Size, result)
}

// SaCustomerLists 计划任务的sacode&&unionid|customerid的顾客列表
func (p PlansHandle) SaCustomerLists(c *gin.Context) {
	var request planReq.CustomerListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	// 获取任务中心的指定时间未完成主任务数
	total, result, err := planService.NewPlanTask().SaCustomerLists(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResList(c, total, request.Page, request.Size, result)
}

// ConfirmCall 确认电话沟通
func (p PlansHandle) ConfirmCall(c *gin.Context) {
	var request planReq.ConfirmReq
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		p.ResError(c, err)
		return
	}
	request.MemberId = c.GetInt64(plan_task2.MemberId)
	request.UserId = c.GetString(plan_task2.UserId)
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(p.WarpContext(c), map[string]interface{}{"enterprise_id": request.EnterpriseId, "key": "store_name"})
	if err != nil {
		p.ResError(c, err)
		return
	}
	if res.ID == 0 {
		p.ResError(c, errors.New("企业不存在"))
		return
	}
	isCoach := false
	if res.Name == "coach" { // coach 为进行中的订单
		isCoach = true
	}
	canAddPoints, err := planService.NewPlanTaskConditionData().ConfirmCall(p.WarpContext(c), request, isCoach)
	if err != nil {
		p.ResError(c, err)
		return
	}
	if canAddPoints {
		// 积分队列埋点
		md5Str := helper.Md5V(request.UnionId + strconv.FormatInt(request.TaskAllotLogId, 10) + request.CustomerId)
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			MemberId:     request.MemberId,
			Type:         excitationModels.RuleCategoryId9,
			Value:        1,
			UniqueHash:   md5Str + fmt.Sprintf("-%d-%d", request.MemberId, excitationModels.RuleCategoryId9),
			EnterpriseId: request.EnterpriseId,
		})
	}
	p.ResOk(c)
}

// ConfirmContact 确认发消息按钮的点击
func (p PlansHandle) ConfirmContact(c *gin.Context) {
	var request planReq.ConfirmReq
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		p.ResError(c, err)
		return
	}
	request.MemberId = c.GetInt64(plan_task2.MemberId)
	request.UserId = c.GetString(plan_task2.UserId)
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(p.WarpContext(c), map[string]interface{}{"enterprise_id": request.EnterpriseId, "key": "store_name"})
	if err != nil {
		p.ResError(c, err)
		return
	}
	if res.ID == 0 {
		p.ResError(c, errors.New("企业不存在"))
		return
	}
	isCoach := false
	if res.Name == "coach" { // coach 为进行中的订单
		isCoach = true
	}
	canAddPoints, err := planService.NewPlanTaskConditionData().ConfirmContact(p.WarpContext(c), request, isCoach)
	if err != nil {
		p.ResError(c, err)
		return
	}
	if canAddPoints {
		// 积分队列埋点
		md5Str := helper.Md5V(request.UnionId + strconv.FormatInt(request.TaskAllotLogId, 10) + request.CustomerId)
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			MemberId:     request.MemberId,
			Type:         excitationModels.RuleCategoryId9,
			Value:        1,
			UniqueHash:   md5Str + fmt.Sprintf("-%d-%d", request.MemberId, excitationModels.RuleCategoryId9),
			EnterpriseId: request.EnterpriseId,
		})
	}

	p.ResOk(c)
}

// ConfirmFinish 导购点击标记完成按钮
func (p PlansHandle) ConfirmFinish(c *gin.Context) {
	var request planReq.ConfirmReq
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		p.ResError(c, err)
		return
	}
	cc := p.WarpContext(c)

	request.MemberId = c.GetInt64(plan_task2.MemberId)
	request.UserId = c.GetString(plan_task2.UserId)
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(cc, map[string]interface{}{"enterprise_id": request.EnterpriseId, "key": "store_name"})
	if err != nil {
		p.ResError(c, err)
		return
	}
	if res.ID == 0 {
		p.ResError(c, errors.New("企业不存在"))
		return
	}
	isCoach := false
	if res.Name == "coach" { // coach 为进行中的订单
		isCoach = true
	}
	canAddPoints, err := planService.NewPlanTaskConditionData().ConfirmFinish(cc, request, isCoach)
	if err != nil {
		p.ResError(c, err)
		return
	}

	if canAddPoints {
		// 积分队列埋点
		md5Str := helper.Md5V(request.UnionId + strconv.FormatInt(request.TaskAllotLogId, 10) + request.CustomerId)
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			MemberId:     request.MemberId,
			Type:         excitationModels.RuleCategoryId9,
			Value:        1,
			UniqueHash:   md5Str + fmt.Sprintf("-%d-%d", request.MemberId, excitationModels.RuleCategoryId9),
			EnterpriseId: request.EnterpriseId,
		})
	}
	// // 神策数据上报-完成任务
	// _, _ = microwork.NewShenCe().TaskComplete(cc, request.EnterpriseId, request.UserId, request.TaskAllotLogId)

	p.ResOk(c)
}

// SetRemark 导购编辑联系说明
func (p PlansHandle) SetRemark(c *gin.Context) {
	var request planReq.SetRemark
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		p.ResError(c, err)
		return
	}
	request.MemberId = c.GetInt64(plan_task2.MemberId)
	request.UserId = c.GetString(plan_task2.UserId)
	err := planService.NewPlanTaskConditionData().SetRemark(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

// SetForward 点击转发给顾客的按钮上报
func (p PlansHandle) SetForward(c *gin.Context) {
	var request planReq.SetForward
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(p.WarpContext(c), map[string]interface{}{"enterprise_id": request.EnterpriseId, "key": "store_name"})
	if err != nil {
		p.ResError(c, err)
		return
	}
	if res.ID == 0 {
		p.ResError(c, errors.New("企业不存在"))
		return
	}
	isCoach := false
	if res.Name == "coach" { // coach 为进行中的订单
		isCoach = true
	}
	err = planService.NewPlanTaskConditionData().SetForward(p.WarpContext(c), request, isCoach)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

// TaskDetailExportBase 根据筛选查出要查询的基本数据
func (p PlansHandle) TaskDetailExportBase(c *gin.Context) {
	var request planReq.TaskDetailExportReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTaskConditionData().TaskDetailExportBase(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, result)
}

// TaskDetailExportList 根据查询的结果查询要导出的数据
func (p PlansHandle) TaskDetailExportList(c *gin.Context) {
	var request planReq.TaskDetailExportListReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTaskConditionData().TaskDetailExportList(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, result)
}

// ConfirmFriend 添加好友按钮的上报
func (p PlansHandle) ConfirmFriend(c *gin.Context) {
	var request planReq.ConfirmReq
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		p.ResError(c, err)
		return
	}
	request.MemberId = c.GetInt64(plan_task2.MemberId)
	request.UserId = c.GetString(plan_task2.UserId)
	err := planService.NewPlanTaskConditionData().ConfirmFriend(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

// NewStat  主任务的任务数据新接口
func (p PlansHandle) NewStat(c *gin.Context) {
	var request planReq.NewStatReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTask().GetPlanStat(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, result)
}

// NewStatDetail 主任务的数据明细接口
func (p PlansHandle) NewStatDetail(c *gin.Context) {
	var request planReq.NewStatDetailReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	var result planReq.NewStatDetailRes
	total, lists, err := planService.NewPlanTask().GetPlanStatDetail(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	result.Count = total
	result.Data = lists
	p.ResSuccess(c, result)
}

// NewCloseTouchTips  主任务的关闭未触达
func (p PlansHandle) NewCloseTouchTips(c *gin.Context) {
	var request planReq.PlanInfo
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	err := planService.NewPlanTask().CloseTouchTips(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

// NewCloseTaskTouchTips 子任务的未触达提示关闭
func (p PlansHandle) NewCloseTaskTouchTips(c *gin.Context) {
	var request planReq.PlanInfo
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	err := planService.NewPlanTask().CloseTaskTouchTips(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

func (p PlansHandle) TaskDetailExportListNew(c *gin.Context) {
	var request planReq.TaskDetailExportListNewReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	var total int64
	result, total, err := planService.NewPlanTaskConditionData().TaskDetailExportBaseNew(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResList(c, total, request.Page, request.Size, result)
}

// SearchCustomerLists coach小程序端，升级计划任务的顾客列表
func (p PlansHandle) SearchCustomerLists(c *gin.Context) {
	var request planReq.CustomerSearchListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(p.WarpContext(c), request.EnterpriseId, "store_name")
	if err != nil {
		p.ResError(c, err)
	}
	// 获取任务中心的顾客数据
	total, result, err := planService.NewPlanTask().SearchCustomerLists(p.WarpContext(c), request, resConfig.Name)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResList(c, total, request.Page, request.Size, result)
}

// SearchExternalLists coach小程序端，升级计划任务的顾客列表
func (p PlansHandle) SearchExternalLists(c *gin.Context) {
	var request planReq.ExternalSearchListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(p.WarpContext(c), request.EnterpriseId, "store_name")
	if err != nil {
		p.ResError(c, err)
	}
	// 获取任务中心的顾客数据
	total, result, err := planService.NewPlanTask().SearchExternalLists(p.WarpContext(c), request, resConfig.Name)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResList(c, total, request.Page, request.Size, result)
}

// SearchExternalAllOpenId  coach升级的计划任务顾客列表,全选获取openid
func (p PlansHandle) SearchExternalAllOpenId(c *gin.Context) {
	var request planReq.ExternalSearchListsReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	// 获取任务中心的顾客数据
	unionIds, err := planService.NewPlanTask().SearchExternalAllOpenId(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, unionIds)
}

func (p PlansHandle) CustomerMove(c *gin.Context) {
	var request planReq.CustomerMoveReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	// 顾客列表名单移动
	err := planService.NewPlanTask().CustomerMove(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

func (p PlansHandle) SendFinish(c *gin.Context) {
	var request planReq.SendFinishReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	err := planService.NewPlanTask().SendFinish(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResOk(c)
}

func (p PlansHandle) IsShow(c *gin.Context) {
	var request planReq.NewStatReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTask().ShowTips(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, gin.H{
		"show": result,
	})
}

func (p PlansHandle) TaskPlanNameList(c *gin.Context) {
	var request planReq.NewStatReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	result, err := planService.NewPlanTask().TaskPlanNameList(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResSuccess(c, result)
}

func (p PlansHandle) CustomerTaskList(c *gin.Context) {
	var request planReq.CustomerTaskReq
	if err := c.ShouldBind(&request); err != nil {
		p.ResError(c, err)
		return
	}
	taskList, total, err := planService.NewPlanTask().GetCustomerTask(p.WarpContext(c), request)
	if err != nil {
		p.ResError(c, err)
		return
	}
	p.ResList(c, total, request.Page, request.Size, taskList)
}

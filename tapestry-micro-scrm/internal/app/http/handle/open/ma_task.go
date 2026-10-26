package open

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	planService "micro-scrm/internal/app/service/plan"
)

type MaTaskHandle struct {
	handle.BaseHandle
}

func NewMaTaskHandle() MaTaskHandle {
	return MaTaskHandle{}
}

// PlanTemplateLists MA三方获取任务模板列表
func (h MaTaskHandle) PlanTemplateLists(c *gin.Context) {
	var request openReq.PlanTemplateListsReq
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	lists, err := planService.NewPlanTemplateSvc().MaTemplateLists(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, lists)
}

// AddPlan MA三方创建任务
func (h MaTaskHandle) AddPlan(c *gin.Context) {
	var request openReq.AddPlanReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := c.GetInt64(middleware.EnterpriseId)
	ctx := h.WarpContext(c)
	planId, _, err := planService.NewPlanTask().AddMaPlan(ctx, enterpriseId, request)
	// returnTaskId := taskId
	// if request.PlanType == 1 { // 持续任务返回的任务id为主任务id,定时任务返回子任务id
	// 	returnTaskId = planId
	// }
	err = planService.NewPlanTaskConditionData().CreateShardTable(ctx, enterpriseId, planId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"task_id": planId})
}

// CreatePlanByContent creates a CGE MA main task.
func (h MaTaskHandle) CreatePlanByContent(c *gin.Context) {
	var request openReq.CreatePlanByContentReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := c.GetInt64(middleware.EnterpriseId)
	ctx := h.WarpContext(c)
	planId, err := planService.NewPlanTask().CreateMaPlanByContent(ctx, enterpriseId, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = planService.NewPlanTaskConditionData().CreateShardTable(ctx, enterpriseId, planId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"plan_id": planId})
}

// NewTaskCrowd MA的人群包
func (h MaTaskHandle) NewTaskCrowd(c *gin.Context) {
	var request openReq.TaskCrowdReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := c.GetInt64(middleware.EnterpriseId)
	err := planService.NewPlanTask().MaTaskCrowd(h.WarpContext(c), enterpriseId, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h MaTaskHandle) CreateTask(c *gin.Context) {
	var request openReq.CreateTaskReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := c.GetInt64(middleware.EnterpriseId)
	taskId, err := planService.NewPlanTask().CreateMaTask(h.WarpContext(c), enterpriseId, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"task_id": taskId})
}

func (h MaTaskHandle) PushTaskCrowd(c *gin.Context) {
	var request openReq.PushTaskCrowdReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := c.GetInt64(middleware.EnterpriseId)
	if err := planService.NewPlanTask().PushMaTaskCrowd(h.WarpContext(c), enterpriseId, request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

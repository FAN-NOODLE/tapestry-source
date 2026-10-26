package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/transfer_customer"
)

type tasksHandler struct {
	Name string
}

func newTasksHandler() *tasksHandler {
	return &tasksHandler{
		Name: "管理后台Handler",
	}
}

// 任务报表
func (handler tasksHandler) report(c *gin.Context) {
	var param transfer_customer.TasksReportRequest
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if result, resErr := transfer_customer.NewTasksCase().Report(ctx.Wrap(c), param); resErr == nil {
		ginplus.ResSuccess(c, result)
	} else {
		ginplus.ResError(c, resErr)
	}
}

// 根据任务的ID，获取任务的详情
func (handler tasksHandler) get(c *gin.Context) {
	var param transfer_customer.TasksGetRequest
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if result, resErr := transfer_customer.NewTasksCase().Get(ctx.Wrap(c), param); resErr == nil {
		ginplus.ResSuccess(c, result)
	} else {
		ginplus.ResError(c, resErr)
	}
}

// 根据条件查询任务的详情
func (handler tasksHandler) getList(c *gin.Context) {
	var param transfer_customer.TasksGetListRequest
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if result, resErr := transfer_customer.NewTasksCase().GetList(ctx.Wrap(c), param); resErr == nil {
		ginplus.ResSuccess(c, result)
	} else {
		ginplus.ResError(c, resErr)
	}
}

// 新增任务
func (handler tasksHandler) add(c *gin.Context) {
	var param transfer_customer.TasksAddRequest
	paramErr := c.ShouldBind(&param)

	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if addErr := transfer_customer.NewTasksCase().Add(ctx.Wrap(c), param); addErr == nil {
		ginplus.ResSuccess(c, "新增成功")
	} else {
		ginplus.ResError(c, addErr)
	}
}

// 更新任务
func (handler tasksHandler) update(c *gin.Context) {
	var param transfer_customer.TasksUpdateRequest
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if updateErr := transfer_customer.NewTasksCase().Update(ctx.Wrap(c), param); updateErr == nil {
		ginplus.ResSuccess(c, "更新成功")
	} else {
		ginplus.ResError(c, updateErr)
	}
}

// 删除任务
func (handler tasksHandler) delete(c *gin.Context) {
	var param transfer_customer.TasksDeleteRequest
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	if deleteErr := transfer_customer.NewTasksCase().Delete(ctx.Wrap(c), param); deleteErr == nil {
		ginplus.ResSuccess(c, "删除成功")
	} else {
		ginplus.ResError(c, deleteErr)
	}
}

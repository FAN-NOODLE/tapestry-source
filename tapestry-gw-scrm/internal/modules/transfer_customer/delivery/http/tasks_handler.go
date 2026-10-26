package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type tasksHandler struct {
	Name string
	sdk  *microscrm.TransferCustomerApi
}

func newTasksHandler() *tasksHandler {
	return &tasksHandler{
		Name: "管理后台网关Handler",
		sdk:  microscrm.NewTransferCustomerApi(),
	}
}

// 任务报表
func (handler tasksHandler) report(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId

	response, err := handler.sdk.TasksReport(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 根据任务的ID，获取任务的详情
func (handler tasksHandler) get(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId

	response, err := handler.sdk.TasksGet(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 根据条件查询任务的详情
func (handler tasksHandler) getList(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	response, err := handler.sdk.TasksList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 新增任务
func (handler tasksHandler) add(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	response, err := handler.sdk.TasksAdd(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 更新任务
func (handler tasksHandler) update(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	response, err := handler.sdk.TasksUpdate(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 删除任务
func (handler tasksHandler) delete(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	response, err := handler.sdk.TasksDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

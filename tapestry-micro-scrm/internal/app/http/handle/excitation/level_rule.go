package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

type LevelRuleHandle struct {
	handle.BaseHandle
	levelRuleSvc *excitationSvc.LevelRuleService
}

func NewLevelRuleHandle() LevelRuleHandle {
	return LevelRuleHandle{
		levelRuleSvc: excitationSvc.NewLevelRuleService(),
	}
}

// Save 新增或修改等级规则
func (h LevelRuleHandle) Save(c *gin.Context) {

	var reqSave excitationReq.LevelRuleSaveArrReq
	if err := c.ShouldBind(&reqSave); err != nil {
		h.ResError(c, err)
		return
	}

	if err := h.levelRuleSvc.CreateOrUpdate(h.WarpContext(c), reqSave); err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)

	return
}

// DetailList 等级规则详情
func (h LevelRuleHandle) DetailList(c *gin.Context) {

	var req excitationReq.LevelRuleDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.levelRuleSvc.GetDetailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

// List 等级规则列表
func (h LevelRuleHandle) List(c *gin.Context) {

	var req excitationReq.LevelRuleDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.levelRuleSvc.GetList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

// IncentiveHandle 成长激励维度
type IncentiveHandle struct {
	handle.BaseHandle
	incentiveConfigSvc   *excitationSvc.IncentiveConfigService
	incentiveRelationSvc *excitationSvc.IncentiveRelationService
	integralRuleSvc      *excitationSvc.IntegralRuleService
}

func NewIncentiveHandle() IncentiveHandle {
	return IncentiveHandle{
		incentiveConfigSvc:   excitationSvc.NewIncentiveConfigService(),
		incentiveRelationSvc: excitationSvc.NewIncentiveRelationService(),
		integralRuleSvc:      excitationSvc.NewIntegralRuleService(),
	}
}

// IncentiveConfig 积分规则维度配置
func (h IncentiveHandle) IncentiveConfig(c *gin.Context) {
	var request excitationReq.IncentiveConfigListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// 初始化维度配置
	errs := h.incentiveConfigSvc.IncentiveConfigInit(h.WarpContext(c), request)
	if errs != nil {
		h.ResError(c, errs)
		return
	}

	res, _, err := h.incentiveConfigSvc.GetConfigList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// RelationSave 新增或修改指标配置
func (h IncentiveHandle) RelationSave(c *gin.Context) {

	var reqSave excitationReq.IncentiveRelationSaveReq
	if err := c.ShouldBind(&reqSave); err != nil {
		h.ResError(c, err)
		return
	}

	if reqSave.ID > 0 {
		if err := h.incentiveRelationSvc.RelationUpdate(h.WarpContext(c), reqSave); err != nil {
			h.ResError(c, err)
			return
		}
	} else {
		if err := h.incentiveRelationSvc.RelationCreate(h.WarpContext(c), reqSave); err != nil {
			h.ResError(c, err)
			return
		}
	}
	h.ResOk(c)

	return
}

// RelationDetail 关联图谱详情
func (h IncentiveHandle) RelationDetail(c *gin.Context) {
	var request excitationReq.IncentiveRelationDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.incentiveConfigSvc.RelationDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

func (h IncentiveHandle) RelationDetailList(c *gin.Context) {
	var request excitationReq.IncentiveRelationDetailListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.incentiveConfigSvc.RelationDetailList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

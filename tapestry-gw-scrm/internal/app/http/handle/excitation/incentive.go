package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type IncentiveHandle struct {
	handle.BaseHandle
	api microscrm.IncentiveApi
}

func NewIncentiveHandle() IncentiveHandle {
	return IncentiveHandle{
		api: microscrm.NewIncentiveApi(),
	}
}

// IncentiveConfig 积分规则维度配置
func (h IncentiveHandle) IncentiveConfig(c *gin.Context) {
	var request excitationReq.IncentiveConfigListReq

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID

	res, err := h.api.GetConfig(h.WarpContext(c), request)
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

	enterpriseID, _ := common.GetEnterprise(c)
	reqSave.EnterpriseId = enterpriseID
	res, err := h.api.RelationSave(h.WarpContext(c), reqSave)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)

	return
}

func (h IncentiveHandle) RelationDetail(c *gin.Context) {

	var reqSave excitationReq.IncentiveRelationDetailReq
	if err := c.ShouldBind(&reqSave); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	reqSave.EnterpriseId = enterpriseID
	res, err := h.api.RelationDetail(h.WarpContext(c), reqSave)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)

	return
}

func (h IncentiveHandle) RelationDetailList(c *gin.Context) {

	var req excitationReq.IncentiveRelationDetailListReq
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.RelationDetailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)

	return
}

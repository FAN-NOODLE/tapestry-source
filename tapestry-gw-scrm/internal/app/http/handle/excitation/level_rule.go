package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type LevelRuleHandle struct {
	handle.BaseHandle
	api microscrm.LevelRuleApi
}

func NewLevelRuleHandle() LevelRuleHandle {
	return LevelRuleHandle{
		api: microscrm.NewLevelRuleApi(),
	}
}

// Save 新增或修改等级规则
func (h LevelRuleHandle) Save(c *gin.Context) {

	var reqSave excitationReq.LevelRuleSaveArrReq
	if err := c.ShouldBind(&reqSave); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	reqSave.EnterpriseId = enterpriseID

	res, err := h.api.Save(h.WarpContext(c), reqSave)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// DetailList 等级规则详情
func (h LevelRuleHandle) DetailList(c *gin.Context) {

	var req excitationReq.LevelRuleDetailListReq

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID
	res, err := h.api.GetDetailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// List 等级规则列表
func (h LevelRuleHandle) List(c *gin.Context) {

	var req excitationReq.LevelRuleDetailListReq

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID
	res, err := h.api.GetList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	"time"
)

type Crm struct {
	handle.BaseHandle
}

func NewCrm() Crm {
	return Crm{}
}

// CustomerClear 顾客清理
func (h Crm) CustomerClear(c *gin.Context) {
	cc := h.WarpContext(c)

	var request openReq.CustomerClearReq

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)

	externalRes, err := microScrmDao.NewWxUserCorpExternal().GetOne(cc, map[string]any{"unionid": request.UnionId}, request.EnterpriseId)
	if err != nil {
		h.ResError(c, errors.New("查询用户信息异常"))
		return
	}

	if externalRes.ID == 0 {
		h.ResError(c, errors.New("未查询到用户信息"))
		return
	}

	err = microScrmDao.NewWxUserCorpExternal().UpdateData(cc, request.EnterpriseId, map[string]any{"id": externalRes.ID}, map[string]interface{}{
		"sex":                      0,
		"mobile":                   "",
		"birthday":                 "",
		"union_code":               "",
		"point":                    "",
		"grade":                    "",
		"birth_year":               "",
		"birth_month":              "",
		"birth_day":                "",
		"card_level":               "",
		"clt_type":                 "",
		"amount_total":             0,
		"customer_code":            "",
		"customer_id":              "",
		"last_purchase_date":       0,
		"last_purchase_store_code": "",
		"upgrade_need_amount":      0,
		"customer_type":            0,
		"customer_type_updated_at": 0,
		"vip_type":                 0,
		"customer_deleted_at":      time.Now().Unix(),
		"updated_at":               time.Now().Unix(),
	})
	if err != nil {
		h.ResError(c, errors.New("清理顾客信息异常"))
		return
	}

	h.ResSuccess(c, externalRes.ID)

	return
}

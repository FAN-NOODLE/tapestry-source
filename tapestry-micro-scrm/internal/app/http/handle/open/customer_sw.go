package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type CustomerSw struct {
	handle.BaseHandle
}

func NewCustomerSw() CustomerSw {
	return CustomerSw{}
}

// BatchUpdateCouponStatus   顾客优惠券状态批量更新
func (h CustomerSw) BatchUpdateCouponStatus(c *gin.Context) {
	var request openReq.BatchUpdateCouponSwReq

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if len(request.DataList) > 10000 {
		h.ResError(c, errors.New("入参数量过大"))
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	config, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseId, "store_name")

	if config.Name != "sw" {
		h.ResError(c, errors.New("账号无权限"))
		return
	}
	err := customerSvc.NewExternalSearchSw().BatchSaveCouponExtraData(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, nil)
	return
}

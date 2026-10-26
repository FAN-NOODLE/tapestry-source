package customer

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type InfoAddHandle struct {
	handle.BaseHandle
}

func NewInfoAddHandle() InfoAddHandle {
	return InfoAddHandle{}
}

func (h InfoAddHandle) SaveInfoAdd(c *gin.Context) {
	var request customerReq.CustomerInfoAddReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if err := customerSvc.NewCustomerInfoAdd().SaveInfoAdd(ctx.Wrap(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

func (h InfoAddHandle) GetInfoAdd(c *gin.Context) {
	var request customerReq.GetCustomerInfoAddReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := customerSvc.NewCustomerInfoAdd().GetInfoAdd(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

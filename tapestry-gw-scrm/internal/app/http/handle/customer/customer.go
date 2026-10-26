package customer

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	customerReq "gw-scrm/internal/app/http/request/customer"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type CustomerHandle struct {
	handle.BaseHandle
}

func NewCustomerHandle() CustomerHandle {
	return CustomerHandle{}
}

func (h CustomerHandle) UserList(c *gin.Context) {
	var request customerReq.UserList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.CustomerUserList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

func (h CustomerHandle) UserDetail(c *gin.Context) {
	var request struct {
		UserID       int64 `json:"user_id" binding:"required" required:"客户ID必填"`
		EnterpriseID int64 `json:"enterprise_id"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.CustomerUserDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

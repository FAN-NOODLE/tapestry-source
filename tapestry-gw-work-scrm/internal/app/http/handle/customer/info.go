package customer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	customerReq "gw-work-scrm/internal/app/http/request/customer"
	"gw-work-scrm/internal/pkg/redis"
	microScrmCustomer "gw-work-scrm/internal/pkg/sdk/microscrm/customer"
)

type InfoHandle struct {
	handle.BaseHandle
	sdkCustomerInfo microScrmCustomer.CustomerInfoApi
}

func NewInfoHandle() InfoHandle {
	return InfoHandle{
		sdkCustomerInfo: microScrmCustomer.NewCustomerInfoApi(),
	}
}

func (h InfoHandle) SaveInfoAdd(c *gin.Context) {
	var request customerReq.SaveInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	request.EnterpriseID = h.EnterpriseID(c)

	memberInfo := h.MemberInfo(c)
	request.OperateUserId = memberInfo.Userid
	request.OperateName = memberInfo.Name

	// 删除客户画像缓存
	cacheKey := fmt.Sprintf("USER:CRM:%d:info:%s", request.EnterpriseID, request.ExternalUserID)
	redis.Client().Del(c, cacheKey)

	result, err := h.sdkCustomerInfo.SaveInfoAdd(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return
}

func (h InfoHandle) SimpleInfo(c *gin.Context) {
	var request customerReq.SimpleInfoRequest
	request.EnterpriseID = h.EnterpriseID(c)

	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.sdkCustomerInfo.SimpleInfo(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, info)
}

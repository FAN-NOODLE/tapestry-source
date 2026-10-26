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

type TagsHandle struct {
	handle.BaseHandle
	sdkCustomerTag microScrmCustomer.CustomerTagsApi
}

func NewTagsHandle() TagsHandle {
	return TagsHandle{
		sdkCustomerTag: microScrmCustomer.NewCustomerTagsApi(),
	}
}

func (h TagsHandle) Operate(c *gin.Context) {
	var request customerReq.CreateTags
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	request.EnterpriseID = h.EnterpriseID(c)

	memberInfo := h.MemberInfo(c)
	request.OperateUserId = memberInfo.Userid
	request.OperateName = memberInfo.Name
	request.OperateMobile = memberInfo.Mobile

	// 删除客户画像缓存
	cacheKey := fmt.Sprintf("USER:CRM:%d:info:%s", request.EnterpriseID, request.ExternalUserID)
	redis.Client().Del(c, cacheKey)

	result, err := h.sdkCustomerTag.Operate(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return
}

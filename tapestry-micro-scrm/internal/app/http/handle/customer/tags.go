package customer

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
	membersSvc "micro-scrm/internal/app/service/member"
)

type TagsHandle struct {
	handle.BaseHandle
}

func NewTagsHandle() TagsHandle {
	return TagsHandle{}
}

func (h TagsHandle) Operate(c *gin.Context) {
	var request customerReq.CreateTags
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.OperateUserId == "" {
		// 获取member信息
		memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(ctx.Wrap(c), request.EnterpriseID, request.OperateMobile)
		if err1 != nil {
			h.ResError(c, err1)
			return
		}

		request.OperateUserId = memberInfo.Userid
		request.OperateName = memberInfo.Name
		request.OperateMobile = memberInfo.Mobile
	}

	if err := customerSvc.NewImpressionTag().Update(ctx.Wrap(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

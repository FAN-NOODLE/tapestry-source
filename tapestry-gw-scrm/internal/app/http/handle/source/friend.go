package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	sourceRequest "gw-scrm/internal/app/http/request/source"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type FriendHandle struct {
	handle.BaseHandle
	sdk microscrm.SourceFriendApi
}

func NewFriendHandle() FriendHandle {
	return FriendHandle{
		sdk: microscrm.NewSourceFriendApi(),
	}
}

func (h FriendHandle) Create(c *gin.Context) {
	var request sourceRequest.FriendCreateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := h.sdk.Create(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

func (h FriendHandle) Detail(c *gin.Context) {
	var request sourceRequest.FriendDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := h.sdk.Create(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

func (h FriendHandle) Update(c *gin.Context) {
	var request sourceRequest.FriendUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := h.sdk.Update(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

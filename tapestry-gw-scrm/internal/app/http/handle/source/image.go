package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	sourceRequest "gw-scrm/internal/app/http/request/source"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type ImageHandle struct {
	handle.BaseHandle
	imageSdk microscrm.SourceImageApi
}

func NewImageHandle() ImageHandle {
	return ImageHandle{
		imageSdk: microscrm.NewSourceImageApi(),
	}
}

func (s ImageHandle) Detail(c *gin.Context) {
	var request sourceRequest.ImageDetailRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := s.imageSdk.Detail(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s ImageHandle) Update(c *gin.Context) {
	var request sourceRequest.ImageUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.imageSdk.Update(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

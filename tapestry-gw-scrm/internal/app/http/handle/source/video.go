package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	sourceRequest "gw-scrm/internal/app/http/request/source"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type VideoHandle struct {
	handle.BaseHandle
	videoSdk microscrm.SourceVideoApi
}

func NewVideoHandle() VideoHandle {
	return VideoHandle{
		videoSdk: microscrm.NewSourceVideoApi(),
	}
}

func (s VideoHandle) Detail(c *gin.Context) {
	var request sourceRequest.VideoDetailRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := s.videoSdk.Detail(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s VideoHandle) Update(c *gin.Context) {
	var request sourceRequest.VideoUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.videoSdk.Update(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

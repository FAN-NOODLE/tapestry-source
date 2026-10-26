package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/source"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type BrandHandle struct {
	handle.BaseHandle
	sourceSdk *microScrm.SourceApi
}

func NewBrandHandle() BrandHandle {
	return BrandHandle{
		sourceSdk: microScrm.NewScource(config.CfgData.Service.MicroScrm.Domain),
	}
}

func (h BrandHandle) BrandSearchList(c *gin.Context) {
	var request source.BrandListReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	res, err := h.sourceSdk.BrandSearchList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h BrandHandle) BrandGetCacheInfo(c *gin.Context) {
	var request source.BrandCacheReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	res, err := h.sourceSdk.BrandGetCacheInfo(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h BrandHandle) BrandCollectSearchList(c *gin.Context) {
	var request source.BrandListReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	res, err := h.sourceSdk.BrandCollectSearchList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

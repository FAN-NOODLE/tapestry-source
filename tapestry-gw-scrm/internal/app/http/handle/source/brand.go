package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request/source"
	"gw-scrm/internal/pkg/common"
	errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type BrandHandle struct {
	handle.BaseHandle
	SourceBrandApi microscrm.SourceBrandApi
}

func NewBrandHandle() BrandHandle {
	return BrandHandle{
		SourceBrandApi: microscrm.NewSourceBrandApi(),
	}
}

func (h BrandHandle) BrandSearchList(c *gin.Context) {
	var request source.BrandListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.AdminId = adminId

	res, err := h.SourceBrandApi.BrandSearchList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

//
//func (h BrandHandle) BrandLinkSave(c *gin.Context) {
//	var request source.RequestLinkSaveReq
//	if err := c.ShouldBind(&request); err != nil {
//		h.ResError(c, errs.ParamsErr)
//		return
//	}
//
//	enterpriseID, adminId := common.GetEnterprise(c)
//	request.EnterpriseId = enterpriseID
//	request.AdminId = adminId
//
//	res, err := h.SourceBrandApi.BrandLinkSave(ctx.Wrap(c), request)
//	if err != nil {
//		h.ResError(c, err)
//		return
//	}
//	h.ResSuccess(c, res)
//	return
//}

func (h BrandHandle) TalkieImageSearch(c *gin.Context) {
	var request source.BrandImageAiSearchReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.AdminId = adminId

	res, err := h.SourceBrandApi.TalkieImageSearch(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

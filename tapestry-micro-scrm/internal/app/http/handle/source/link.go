package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/modules/source"
)

type LinkHandle struct {
	handle.BaseHandle
	linkService sourceService.LinkService
}

func NewLinkHandle() LinkHandle {
	return LinkHandle{
		linkService: sourceService.NewLinkService(),
	}
}

func (h LinkHandle) Detail(c *gin.Context) {
	var request sourceRequest.DetailRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.linkService.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h LinkHandle) Update(c *gin.Context) {
	var request sourceRequest.LinkUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.linkService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	_ = sourceService.NewImageService().UpdateRelationGoods(ctx.Wrap(c), request.ID, uint64(request.EnterpriseID), request.RelationGoods, source.TypeLink)

	h.ResOk(c)
	return
}

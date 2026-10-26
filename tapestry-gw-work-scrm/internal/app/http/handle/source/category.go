package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/source"
	"gw-work-scrm/internal/pkg/errs"
	"gw-work-scrm/internal/pkg/sdk/microscrm"
)

type CategoryHandle struct {
	handle.BaseHandle
	sourceCategorySdk microScrm.SourceCategoryApi
}

func NewCategoryHandle() CategoryHandle {
	return CategoryHandle{
		sourceCategorySdk: microScrm.NewSourceCategoryApi(),
	}
}

func (s CategoryHandle) List(c *gin.Context) {
	var request source.CategoryListRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)
	memberInfo := s.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.OriginH5 = true
	res, err := s.sourceCategorySdk.List(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

package tags

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	tagsReq "gw-scrm/internal/app/http/request/tags"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type ImpressionHandle struct {
	handle.BaseHandle
}

func NewImpressionHandle() ImpressionHandle {
	return ImpressionHandle{}
}

func (h ImpressionHandle) GroupList(c *gin.Context) {
	var request tagsReq.GroupSearchList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.ImpressionTagGroupList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

func (h ImpressionHandle) TagList(c *gin.Context) {
	var request tagsReq.TagList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.ImpressionTagList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

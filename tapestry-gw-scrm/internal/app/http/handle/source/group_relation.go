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

type GroupRelationHandle struct {
	handle.BaseHandle
	groupRelationSdk microscrm.SourceGroupRelationApi
}

func NewGroupRelationHandle() GroupRelationHandle {
	return GroupRelationHandle{
		groupRelationSdk: microscrm.NewSourceGroupRelationApi(),
	}
}

// List 素材分组归属门店列表
func (s GroupRelationHandle) List(c *gin.Context) {
	var request source.GroupRelationListRequest

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	res, err := s.groupRelationSdk.List(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// Save 素材分组归属门店保存
func (s GroupRelationHandle) Save(c *gin.Context) {
	var request source.GroupRelationSaveRequest
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	res, err := s.groupRelationSdk.Save(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// GetListForStore List 根据门店获取素材分组归属门店列表
func (s GroupRelationHandle) GetListForStore(c *gin.Context) {
	var request source.GroupRelationFindRequest

	err := c.BindJSON(&request)
	if err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	res, err := s.groupRelationSdk.GetListForStore(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

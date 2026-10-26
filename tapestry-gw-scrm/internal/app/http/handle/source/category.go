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

type CategoryHandle struct {
	handle.BaseHandle
	sourceCategorySdk microscrm.SourceCategoryApi
}

func NewCategoryHandle() CategoryHandle {
	return CategoryHandle{
		sourceCategorySdk: microscrm.NewSourceCategoryApi(),
	}
}

func (s CategoryHandle) List(c *gin.Context) {
	var request source.CategoryList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.sourceCategorySdk.List(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s CategoryHandle) Create(c *gin.Context) {
	var request source.CategoryCreate
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.sourceCategorySdk.Create(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s CategoryHandle) Update(c *gin.Context) {
	var request source.CategoryUpdate
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.sourceCategorySdk.Update(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s CategoryHandle) SourceRelation(c *gin.Context) {
	var request source.Category
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.sourceCategorySdk.SourceRelation(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s CategoryHandle) Delete(c *gin.Context) {
	var request source.Category
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	res, err := s.sourceCategorySdk.Delete(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

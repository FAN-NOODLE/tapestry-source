package sensitivehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/sensitive"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/uid"
)

var service *sensitive.Service

// sensitiveList 敏感词列表
func sensitiveList(c *gin.Context) {
	req := &sensitive.ListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	count, entity, err := service.SelectList(ctx.Wrap(c), *req)
	res := sensitive.ResponseList{}
	if entity != nil {
		res.List = *entity
	} else {
		res.List = []sensitive.Detail{}
	}
	res.Total = count
	res.Page = req.Page
	res.Size = req.Size

	ginplus.ResSuccess(c, res)
}

// sensitiveSave 敏感词保存
func sensitiveSave(c *gin.Context) {
	req := &sensitive.SaveRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if req.SensitiveType == sensitive.TypeText && req.SensitiveText == "" {
		ginplus.ResError(c, errs.SensitiveTextEmptyErr)
		return
	}

	if req.SensitiveType == sensitive.TypeCart {
		req.SensitiveText = "名片"
	}

	cc := ctx.Wrap(c)
	req.ID = uid.ID(cc)
	req.TriggerNum = 0

	errAdd := service.SensitiveSave(cc, *req)
	if errAdd != nil {
		ginplus.ResError(c, errAdd)
		return
	}
	ginplus.ResSuccess(c, []struct{}{})
}

// sensitiveDelete 敏感词删除
func sensitiveDelete(c *gin.Context) {
	req := &sensitive.DeleteRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	errDelete := service.SensitiveDelete(ctx.Wrap(c), *req)
	if errDelete != nil {
		ginplus.ResError(c, errDelete)
		return
	}
	ginplus.ResSuccess(c, []struct{}{})
}

// sensitiveDetail 敏感词详情
func sensitiveDetail(c *gin.Context) {
	req := &sensitive.DetailListRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	detail, err := service.SensitiveDetail(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

// sensitiveDetailList 敏感词明细列表
func sensitiveDetailList(c *gin.Context) {
	req := &sensitive.DetailListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	count, entity, err := service.SensitiveDetailList(ctx.Wrap(c), *req)
	res := sensitive.DetailListPageResponse{}
	if entity != nil {
		res.List = *entity
	} else {
		res.List = []sensitive.DetailListResponse{}
	}
	res.Total = count
	res.Page = req.Page
	res.Size = req.Size

	ginplus.ResSuccess(c, res)
}

package miniProgramHttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/miniProgram"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

var s *miniProgram.Service

// add 新增小程序
func add(c *gin.Context) {
	params := &miniProgram.MiniProgram{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	params.Type = miniProgram.TypeManual
	id, err := s.Add(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &miniProgram.ID{ID: id})
	return
}

// authorize 维护官方小程序
func authorize(c *gin.Context) {
	params := &miniProgram.MiniProgram{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	params.Type = miniProgram.TypeAuthorize

	id, err := s.Authorize(ctx.Wrap(c), params)

	if err != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &miniProgram.ID{ID: id})
	return
}

// list 小程序列表
func list(c *gin.Context) {
	params := &miniProgram.RequestPrograms{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err := s.List(ctx.Wrap(c), params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// del 删除
func del(c *gin.Context) {
	cc := ctx.Wrap(c)
	params := &struct {
		ID int64 `json:"id"  validate:"required" required:"小程序ID不能为空"`
	}{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if ok := s.IsAuthOperate(cc, params.ID); !ok {
		response.Notice(c, errs.AuthorizeMiniProgramCannotOperate)
		return
	}

	rows, err := s.Del(cc, params.ID)

	if err != nil {
		response.Notice(c, errs.SourceGroupDelErr)
		return
	}

	response.Success(c, &miniProgram.Rows{Rows: rows})
	return
}

// edit 编辑
func edit(c *gin.Context) {
	cc := ctx.Wrap(c)
	params := &miniProgram.MiniProgram{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if params.ID == 0 {
		response.Notice(c, errs.SourceGroupIDEmpty)
		return
	}

	if ok := s.IsAuthOperate(cc, params.ID); !ok {
		response.Notice(c, errs.AuthorizeMiniProgramCannotOperate)
		return
	}

	id, err := s.Edit(cc, params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &miniProgram.ID{ID: id})
	return
}

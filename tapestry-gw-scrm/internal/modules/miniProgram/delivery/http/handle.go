package miniProgramHttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errorIgnore"
	scrm_errs "gw-scrm/internal/pkg/errs"
	response "gw-scrm/internal/pkg/resoponse"
	"gw-scrm/internal/pkg/validate_params"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
)

// add 新增小程序
func add(c *gin.Context) {
	params := &microScrm.MiniProgram{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	vo, err := miniProgramSDK.Add(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

// list 小程序列表
func list(c *gin.Context) {
	params := &microScrm.RequestPrograms{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	total, list, err := miniProgramSDK.List(ctx.Wrap(c), params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
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
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	vo, err := miniProgramSDK.Del(cc, params.ID)

	if err != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

// edit 编辑
func edit(c *gin.Context) {
	cc := ctx.Wrap(c)
	params := &microScrm.MiniProgram{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	vo, err := miniProgramSDK.Edit(cc, params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

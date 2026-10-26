package sourceGroupHttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/sourceGroup/delivery"
	"gw-scrm/internal/pkg/assist"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errorIgnore"
	scrm_errs "gw-scrm/internal/pkg/errs"
	response "gw-scrm/internal/pkg/resoponse"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/validate_params"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
)

// list 列表
func list(c *gin.Context) {
	params := &microScrm.SourceGroup{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	res, err := microscrm.ApiSource{}.SourceTree(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.Success(c, res)
	return
}

// del 删除
func del(c *gin.Context) {
	params := &struct {
		ID           int64 `json:"id" validate:"required" required:"素材分组IDs不能为空"`
		EnterpriseID int64 `json:"enterprise_id"`
	}{}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	_, err := microscrm.ApiSource{}.SourceDelete(ctx.Wrap(c), params)
	if err != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)

		return
	}
	response.Success(c, struct {
	}{})
	return
}

// add 新增
func add(c *gin.Context) {
	var receive delivery.NewSourceGroup
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	receive.EnterpriseID = tokenInfo.EnterpriseId
	res, err := microscrm.ApiSource{}.AddSource(ctx.Wrap(c), receive)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, res)
	return
}

// edit
func edit(c *gin.Context) {

	params := &struct {
		ID        int64  `json:"id" db:"id"  validate:"required" required:"素材分组ID不能为空"`
		GroupName string `json:"group_name" db:"group_name" validate:"required" required:"素材名称不能为空"`
	}{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	cc := ctx.Wrap(c)
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	sourceGroupParams := &microScrm.SourceGroup{
		ID:           params.ID,
		GroupName:    params.GroupName,
		UpdatedAt:    assist.Timestamp(),
		EnterpriseID: tokenInfo.EnterpriseId,
	}

	vo, err := sourceGroupSDK.Edit(cc, sourceGroupParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

func move(c *gin.Context) {
	var info delivery.MoveReceive
	if err := c.ShouldBindJSON(&info); err != nil {
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	info.EnterpriseID = tokenInfo.EnterpriseId

	_, err := microscrm.ApiSource{}.MoveBatch(ctx.Wrap(c), info)
	if err != nil {
		ginplus.ResError(c, err)
	}
	response.Success(c, struct {
	}{})
	return
}

func updateIgnoreIndex(c *gin.Context) {
	params := &struct {
		ID           int64 `json:"id" db:"id"  validate:"required" required:"素材分组ID不能为空"`
		EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
		IgnoreIndex  int   `json:"ignore_index"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	params.EnterpriseID = tokenInfo.EnterpriseId
	_, err := microscrm.ApiSource{}.UpdateIgnoreIndex(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.Success(c, struct{}{})
	return
}

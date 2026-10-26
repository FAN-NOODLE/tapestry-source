package sourceHttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/source"
	"gw-scrm/internal/pkg/assist"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errorIgnore"
	scrm_errs "gw-scrm/internal/pkg/errs"
	response "gw-scrm/internal/pkg/resoponse"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/tencent/tms"
	"gw-scrm/internal/pkg/validate_params"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
	"net/http"
)

var tokenInfo *common.TokenInfo

// SourcesImageList 图片素材List接口
func SourcesImageList(c *gin.Context) {
	params := &microScrm.RequestList{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceImageList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
}

// SourcesFileList 文件素材List接口
func SourcesFileList(c *gin.Context) {
	params := &microScrm.RequestList{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceFileList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, list)
}

// SourceImageAdd 图片素材新增
func SourceImageAdd(c *gin.Context) {
	params := &microScrm.SourcesImageRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	res, err := sourceSDK.SourcesImageAdd(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, res)
	return
}

// SourceFileAdd 文件素材新增
func SourceFileAdd(c *gin.Context) {
	params := &microScrm.SourcesFileRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	vo, err := sourceSDK.SourcesFileAdd(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.Success(c, vo)
	return
}

// sourceMiniProgramSave 小程序素材-新增
func sourceMiniProgramSave(c *gin.Context) {
	params := &microScrm.RequestSourceMini{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	apiResponse, ok, err := tms.Verification(params.Title)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	if !ok {
		resp := &ginplus.APIResponse{
			Code:      ginplus.ParamsErrCode,
			Message:   "存在非法关键字",
			RequestID: ctx.Wrap(c).ID(),
			Data:      apiResponse,
		}
		ginplus.ResJSON(c, http.StatusOK, resp)
		return
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId

	vo := &microScrm.ID{}
	if params.ID > 0 {
		vo, err = sourceSDK.SourceMiniProgramUpdate(ctx.Wrap(c), params)
	} else {
		voMini, errAdd := sourceSDK.SourceMiniProgramAdd(ctx.Wrap(c), params)
		if errAdd != nil {
			response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, errAdd), errAdd)
			return
		}
		response.Success(c, voMini)
		return
	}

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

// sourceMiniProgramList 素材小程序--列表
func sourceMiniProgramList(c *gin.Context) {
	params := &microScrm.RequestSourceSearch{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceMiniProgramList(ctx.Wrap(c), params)
	//

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// sourceMiniProgramDetail 素材小程序--详情
func sourceMiniProgramDetail(c *gin.Context) {
	params := &microScrm.ID{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	detail, err := sourceSDK.SourceMiniProgramDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.Success(c, detail)
	return
}

// sourceVideoSave 视频素材-新增、编辑
func sourceVideoSave(c *gin.Context) {
	params := &microScrm.RequestSourceVideo{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId

	vo := &microScrm.ID{}
	var err error
	if params.ID > 0 {
		vo, err = sourceSDK.SourceVideoUpdate(ctx.Wrap(c), params)
	} else {
		voVideo, errAdd := sourceSDK.SourceVideoAdd(ctx.Wrap(c), params)
		if errAdd != nil {
			response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, errAdd), errAdd)
			return
		}
		response.Success(c, voVideo)
		return
	}

	if err != nil {
		ginplus.ResError(c, err)
	}

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, vo)
	return
}

// sourceVideoList 视频素材--列表
func sourceVideoList(c *gin.Context) {
	params := &microScrm.RequestSourceSearch{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceVideoList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	type SourceVideo struct {
		microScrm.SourceVideo
		HSize string `json:"h_size"`
	}

	var resList []SourceVideo

	for _, v := range *list {
		tmp := SourceVideo{
			SourceVideo: v,
			HSize:       assist.FormatFileSize(v.Size),
		}
		resList = append(resList, tmp)
	}

	response.SuccessList(c, total, params.Page, params.Size, resList)
	return
}

// sourceVideoDetail 视频素材--详情
func sourceVideoDetail(c *gin.Context) {
	params := &microScrm.ID{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	detail, err := sourceSDK.SourceVideoDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	if detail == nil {
		response.Success(c, nil)
	}

	type SourceVideo struct {
		microScrm.SourceVideo
		HSize string `json:"h_size"`
	}

	resDetail := SourceVideo{
		SourceVideo: *detail,
		HSize:       assist.FormatFileSize(detail.Size),
	}

	response.Success(c, resDetail)
	return
}

// sourceTextSave 文本素材新增
func sourceTextSave(c *gin.Context) {
	params := &microScrm.RequestTextSave{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	apiResponse, ok, err := tms.Verification(params.Content)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	if !ok {
		resp := &ginplus.APIResponse{
			Code:      ginplus.ParamsErrCode,
			Message:   "存在非法关键字",
			RequestID: ctx.Wrap(c).ID(),
			Data:      apiResponse,
		}
		ginplus.ResJSON(c, http.StatusOK, resp)
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	params.AdminId = tokenInfo.AdminId
	res, err := sourceSDK.SourceTextSave(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}
	response.Success(c, res)
	return
}

// sourceLinkSave 链接素材新增
func sourceLinkSave(c *gin.Context) {
	params := &microScrm.RequestLinkSave{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	apiResponse, ok, err := tms.Verification(fmt.Sprintf("%s%s", params.Title, params.Description))
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	if !ok {
		resp := &ginplus.APIResponse{
			Code:      ginplus.ParamsErrCode,
			Message:   "存在非法关键字",
			RequestID: ctx.Wrap(c).ID(),
			Data:      apiResponse,
		}
		ginplus.ResJSON(c, http.StatusOK, resp)
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	params.AdminId = tokenInfo.AdminId
	res, err := sourceSDK.SourceLinkSave(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, res)
	return
}

// sourceTextList 文本素材列表
func sourceTextList(c *gin.Context) {
	params := &microScrm.RequestTextList{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceTextList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// sourceLinkList 链接素材列表
func sourceLinkList(c *gin.Context) {
	params := &microScrm.RequestLinkList{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.AdminID = tokenInfo.AdminId
	total, list, err := sourceSDK.SourceLinkList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// sourceTextDetail 文本素材详情
func sourceTextDetail(c *gin.Context) {
	params := &microScrm.RequestSourceDetail{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	detail, err := sourceSDK.SourceTextDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// sourceTextDetail 文本素材详情
func sourceLinkDetail(c *gin.Context) {
	params := &microScrm.RequestSourceDetail{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	detail, err := sourceSDK.SourceLinkDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(scrm_errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// 删除素材
func sourceDelete(c *gin.Context) {
	params := &microScrm.RequestSourceDelete{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	cc := ctx.Wrap(c)
	vo, err := sourceSDK.SourceDelete(cc, params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.Success(c, vo)
	return
}

// 移动素材分组
func moveGroup(c *gin.Context) {
	params := &microScrm.RequestGroupMove{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	cc := ctx.Wrap(c)
	vo, err := sourceSDK.MoveGroup(cc, params)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	response.Success(c, vo)
	return
}

// 获取分组名称
func getGroupName(c *gin.Context) {
	receive := source.GetGroupName{}
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	receive.EnterpriseID = tokenInfo.EnterpriseId
	res, err := microscrm.ApiSource{}.GetGroupName(ctx.Wrap(c), receive)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

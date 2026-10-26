package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/source"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"gw-work-scrm/internal/pkg/sdk/work"
)

type SourceApiHandle struct {
	handle.BaseHandle
	sourceSdk *microScrm.SourceApi
}

func NewSourceHandle() SourceApiHandle {
	return SourceApiHandle{
		sourceSdk: microScrm.NewScource(config.CfgData.Service.MicroScrm.Domain),
	}
}

// Credential 导购获取上传的凭证
func (h SourceApiHandle) Credential(c *gin.Context) {
	res, err := h.sourceSdk.GetUploadCredential(ctx.Wrap(c))
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

// UploadMedia 导购上传素材
func (h SourceApiHandle) UploadMedia(c *gin.Context) {
	var request source.UploadMediaRequest

	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.UserId = memberInfo.Userid
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.sourceSdk.UploadMedia(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

// CheckAuth 验证导购是否有上传素材的权限
func (h SourceApiHandle) CheckAuth(c *gin.Context) {
	var request source.CheckAuth
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.UserId = memberInfo.Userid

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := work.CheckSource(ctx.Wrap(c), &request)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

// Compress 压缩结果
func (h SourceApiHandle) Compress(c *gin.Context) {
	var request source.CompressReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.sourceSdk.Compress(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

// CheckFirst 验证是否用户首次打开素材交互提示
func (h SourceApiHandle) CheckFirst(c *gin.Context) {
	var request source.CheckAuth
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.UserId = memberInfo.Userid
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.sourceSdk.CheckFirst(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h SourceApiHandle) ReportTimes(c *gin.Context) {
	var request source.ReportTimesReq
	request.EnterpriseID = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.sourceSdk.ReportTimes(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h SourceApiHandle) CollectEdit(c *gin.Context) {
	var req source.CollectReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.sourceSdk.CollectEdit(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

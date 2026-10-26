package invite

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	inviteSvc "micro-scrm/internal/app/service/invite"
)

type TemplateHandle struct {
	handle.BaseHandle
	templateService *inviteSvc.InviteTemplateService
}

func NewTemplateHandle() TemplateHandle {
	return TemplateHandle{
		templateService: inviteSvc.NewInviteTemplateService(),
	}
}

// Save Create or Save 创建或保存邀请函模板
func (h TemplateHandle) Save(c *gin.Context) {

	var reqSave inviteRequest.TemplateSaveReq
	if err := c.ShouldBind(&reqSave); err != nil {
		h.ResError(c, err)
		return
	}

	var res inviteRequest.TemplateSaveRes
	if reqSave.ID > 0 {
		res.ID = reqSave.ID
		if err := h.templateService.Update(h.WarpContext(c), reqSave.EnterpriseId, reqSave); err != nil {
			h.ResError(c, err)
			return
		}
	} else {
		createId, err := h.templateService.Create(h.WarpContext(c), reqSave.EnterpriseId, reqSave)
		if err != nil {
			h.ResError(c, err)
			return
		}
		res.ID = createId
	}
	h.ResSuccess(c, res)

	return

}

// List 邀请函模板列表
func (h TemplateHandle) List(c *gin.Context) {

	var reqList inviteRequest.TemplateListReq
	if err := c.ShouldBind(&reqList); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.templateService.List(h.WarpContext(c), reqList.EnterpriseId, reqList)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, reqList.Page, reqList.Size, result)
	return

}

// Info 邀请函模板信息
func (h TemplateHandle) Info(c *gin.Context) {

	var reqInfo inviteRequest.TemplateInfoReq
	if err := c.ShouldBind(&reqInfo); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.templateService.Info(h.WarpContext(c), reqInfo.EnterpriseId, reqInfo)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return

}

// Delete 删除邀请函模板
func (h TemplateHandle) Delete(c *gin.Context) {

	var req inviteRequest.TemplateDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.templateService.Delete(h.WarpContext(c), req.EnterpriseId, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

// ConfigList 邀请函模板配置列表
func (h TemplateHandle) ConfigList(c *gin.Context) {

	var reqList inviteRequest.TemplateConfigListReq
	if err := c.ShouldBind(&reqList); err != nil {
		h.ResError(c, err)
		return
	}

	result, _, err := h.templateService.ConfigList(h.WarpContext(c), reqList.EnterpriseId, reqList)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return

}

// ConfigDetail 邀请函模板配置详情
func (h TemplateHandle) ConfigDetail(c *gin.Context) {

	var reqDetail inviteRequest.TemplateConfigDetailReq
	if err := c.ShouldBind(&reqDetail); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.templateService.ConfigDetail(h.WarpContext(c), reqDetail.EnterpriseId, reqDetail)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return

}

// UpdateStatus 修改邀请函模板上下架
func (h TemplateHandle) UpdateStatus(c *gin.Context) {

	var req inviteRequest.TemplateUpdateStatusReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.templateService.UpdateStatus(h.WarpContext(c), req.EnterpriseId, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

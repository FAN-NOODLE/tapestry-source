package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

type MedalRecordHandle struct {
	handle.BaseHandle
	recordSvc *excitationSvc.MedalRecordService
}

func NewMedalRecordHandle() MedalRecordHandle {
	return MedalRecordHandle{
		recordSvc: excitationSvc.NewMedalRecordService(),
	}
}

// List 奖章记录列表
func (h MedalRecordHandle) List(c *gin.Context) {
	var request excitationReq.MedalRecordList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}

	resp, total, err := h.recordSvc.GetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, resp)

	return
}

// Reissue 补发奖章
func (h MedalRecordHandle) Reissue(c *gin.Context) {
	var request excitationReq.MedalRecordReissue
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	resp, err := h.recordSvc.Reissue(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

// Remove 删除领取奖章记录
func (h MedalRecordHandle) Remove(c *gin.Context) {
	var request struct {
		EnterpriseId int64 `json:"enterprise_id" binding:"required"`
		Id           int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.recordSvc.Remove(h.WarpContext(c), request.EnterpriseId, request.Id)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}

// AlreadyGot 已获得奖章
func (h MedalRecordHandle) AlreadyGot(c *gin.Context) {
	var request excitationReq.MedalRecordAlreadyGot
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	receivedNum, err := h.recordSvc.AlreadyGot(h.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err != nil {
		return
	}

	h.ResSuccess(c, receivedNum)

	return
}

// MemberReceived 领取奖章
func (h MedalRecordHandle) MemberReceived(c *gin.Context) {
	var request excitationReq.MedalRecordMemberReceived
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.recordSvc.MemberReceived(h.WarpContext(c), request.EnterpriseId, request.RecordId, request.Mobile)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")

	return
}

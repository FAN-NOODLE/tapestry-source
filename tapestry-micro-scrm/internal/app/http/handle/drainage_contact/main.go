package drainageContactHandle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/drainage_contact"
	"micro-scrm/internal/app/service/drainage_contacts"
)

type Handle struct {
	handle.BaseHandle
	mainSvc *drainageContactsSvc.MainSvc
}

func NewHandle() Handle {
	return Handle{
		mainSvc: drainageContactsSvc.NewMainSvc(),
	}
}

func (h Handle) GroupSave(c *gin.Context) {
	var req drainageContactGroupReq.GroupSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ContactGroupSave(h.WarpContext(c), req)
	if err != nil {
		// 判断是否是40098未实名错误，将未实名员工信息放入data中返回，前端判断unverified_users有值则弹窗拦截
		if e, ok := err.(*errorx.Error); ok && e.Code == 40098 {
			h.ResSuccess(c, res)
			return
		}
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) ZipUrl(c *gin.Context) {
	var req drainageContactGroupReq.ZipUrlReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ZipUrl(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GroupDetail(c *gin.Context) {
	var req drainageContactGroupReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ContactGroupDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GroupDetailList(c *gin.Context) {
	var req drainageContactGroupReq.DetailListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ContactGroupDetailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GroupList(c *gin.Context) {
	var req drainageContactGroupReq.GroupListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mainSvc.ContactGroupList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) GroupDelete(c *gin.Context) {
	var req drainageContactGroupReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.ContactGroupDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

//
//func (h Handle) ActivityDelete(c *gin.Context) {
//	var req drainageContactGroupReq.InfoReq
//	if err := c.ShouldBindJSON(&req); err != nil {
//		h.ResError(c, err)
//		return
//	}
//	err := h.mainSvc.ActivityDelete(h.WarpContext(c), req)
//	if err != nil {
//		h.ResError(c, err)
//		return
//	}
//	h.ResOk(c)
//}

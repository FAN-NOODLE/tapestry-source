package notification

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	request2 "micro-scrm/internal/app/http/request"
	emailReq "micro-scrm/internal/app/http/request/email"
	request "micro-scrm/internal/app/http/request/notification"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service/email"
	memberSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/app/service/notification"
	"micro-scrm/internal/pkg/helper"
)

type Handle struct {
	handle.BaseHandle
	infoSvc *notificationSvc.InfoSvc
}

func NewHandle() Handle {
	return Handle{
		infoSvc: notificationSvc.NewInfoSvc(),
	}
}

func (h Handle) Save(c *gin.Context) {
	var req request.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	detail, err := memberSvc.NewMembersService().GetInfoByMobile(ctx, req.EnterpriseId, req.Mobile)
	if detail.ID == 0 {
		h.ResError(c, errors.New("该账号非导购身份"))
		return
	}
	if err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.infoSvc.GetManagerInfo(h.WarpContext(c), req.AdminManager)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if info.IsRo != request.IsRo && info.IsFT != 1 {
		h.ResError(c, errors.New("该账号无新建部门传讯权限"))
		return
	}
	err = h.infoSvc.Save(ctx, req, detail.ID, detail.Name)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) List(c *gin.Context) {
	var req request.SearchListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.infoSvc.GetManagerInfo(h.WarpContext(c), req.AdminManager)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if info.IsRo != request.IsRo && req.IsAll == request.IsAll {
		h.ResError(c, errors.New("该账号无RO权限"))
		return
	}
	total, list, err := h.infoSvc.List(h.WarpContext(c), req, info)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) Info(c *gin.Context) {
	var req request.SearchId
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := h.infoSvc.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, detail)
}

func (h Handle) Departments(c *gin.Context) {
	var req request2.EnterpriseBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	departmentLists, err := h.infoSvc.GetDepartments(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, departmentLists)
}

func (h Handle) Copy(c *gin.Context) {
	var req request.SearchId
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	newId, err := h.infoSvc.Copy(ctx, req.EnterpriseId, req.Id, req.AdminID)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"id": newId})
}

func (h Handle) Recall(c *gin.Context) {
	var req request.SearchId
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	err := h.infoSvc.Recall(ctx, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) Check(c *gin.Context) {
	var req request.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	info, err := h.infoSvc.GetManagerInfo(ctx, req.AdminManager)
	if info.IsRo != microScrmModel.NotificationIsRo {
		h.ResError(c, errors.New("该账号无审核权限"))
		return
	}
	err = h.infoSvc.Check(ctx, req, info)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) Manager(c *gin.Context) {
	var req request.AdminManager
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.infoSvc.GetManagerInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h Handle) Cancel(c *gin.Context) {
	var req request.SearchId
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	info, err := h.infoSvc.IsRoRole(h.WarpContext(c), req.AdminManager)
	if !info {
		h.ResError(c, errors.New("该账号无审核权限"))
		return
	}
	err = h.infoSvc.Cancel(ctx, req.EnterpriseId, req.Id)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) SendFileByEmail(c *gin.Context) {
	var req request.SendFileByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	var fileList []emailReq.FileInfo
	var content string
	title, _ := h.infoSvc.GetTitle(h.WarpContext(c), req)
	if req.Attachment == request.IsAttachment {
		fileList = []emailReq.FileInfo{{
			Url:  req.FileUrl,
			Name: req.Rename,
		}}
		content = fmt.Sprintf("您申请的【%s】已生成，烦请查阅附件，谢谢～", req.Rename)
	} else {
		emailData, _ := h.infoSvc.GetEmailData(h.WarpContext(c), req)
		content = emailData.Content
	}

	err := email.NewSendEmail().Send(emailReq.SendMailReq{
		Title:     fmt.Sprintf("【%s】传讯%s下载", title, helper.RuleThree(req.Attachment == request.IsAttachment, "附件", "内容")),
		Content:   content,
		ToAddress: []string{req.Email},
		FileList:  fileList,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) ListReadNum(c *gin.Context) {
	var req request.AdminManager
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	managerInfo, err := h.infoSvc.GetManagerInfo(h.WarpContext(c), req)

	res, err := h.infoSvc.ListReadNum(h.WarpContext(c), req.EnterpriseId, req.MemberId, managerInfo)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

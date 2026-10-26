package storeVisitHandle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	storeVisitReq "gw-work-scrm/internal/app/http/request/store_visit"
	storeVisitSdk "gw-work-scrm/internal/pkg/sdk/microscrm/store_visit"
)

type StoreVisitHandle struct {
	handle.BaseHandle
	storeVisitSdk *storeVisitSdk.StoreVisit
	Language      string
}

func NewStoreVisitHandle() *StoreVisitHandle {
	return &StoreVisitHandle{
		storeVisitSdk: storeVisitSdk.NewStoreVisit(),
	}
}

func (h StoreVisitHandle) SavePlan(c *gin.Context) {
	var req storeVisitReq.SavePlanReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.SavePlan(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) GetStoreManager(c *gin.Context) {
	var req storeVisitReq.StoreManagerReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.GetStoreManager(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) GetStoreMM(c *gin.Context) {
	var req storeVisitReq.StoreManagerReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.GetStoreMM(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) StoreVisitList(c *gin.Context) {
	var req storeVisitReq.PlanListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.List(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) DeptMemberList(c *gin.Context) {
	var req storeVisitReq.BaseReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.DeptMemberList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) GetEmailContact(c *gin.Context) {
	var req storeVisitReq.BaseReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.EmailContact(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) SaveEmail(c *gin.Context) {
	var req storeVisitReq.EmailContactReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.storeVisitSdk.SaveEmail(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) DelEmail(c *gin.Context) {
	var req storeVisitReq.EmailContactReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.storeVisitSdk.DelEmail(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) SearchMembers(c *gin.Context) {
	var req storeVisitReq.SearchMemberReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	list, err := h.storeVisitSdk.SearchMember(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h StoreVisitHandle) AddCopyMember(c *gin.Context) {
	var req storeVisitReq.AddCopyMemberReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	_, err := h.storeVisitSdk.AddCopyMember(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) LoginUserInfo(c *gin.Context) {
	var req storeVisitReq.BaseReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.LoginUserInfo(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) StartVisit(c *gin.Context) {
	var req storeVisitReq.StartVisitReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.StartVisit(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) SaveShowRecordContent(c *gin.Context) {
	var req storeVisitReq.SaveShowRecordContentReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.SaveShowRecordContent(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) SaveShowRecordMedia(c *gin.Context) {
	var req storeVisitReq.SaveShowRecordMediaReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	info, err := h.storeVisitSdk.SaveShowRecordMedia(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h StoreVisitHandle) DelShowRecord(c *gin.Context) {
	var req storeVisitReq.DelShowRecordReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.DelShowRecord(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
func (h StoreVisitHandle) DelShowRecordMedia(c *gin.Context) {
	var req storeVisitReq.DelShowRecordMediaReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.DelShowRecordMedia(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) ShowRecordList(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.ShowRecordList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) StoreVisitInfo(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.StoreVisitInfo(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) CompleteVisit(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.CompleteVisit(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) SendCompleteReport(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.SendCompleteReport(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) SaveVisitInfo(c *gin.Context) {
	var req storeVisitReq.SaveVisitInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.SaveVisitInfo(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) SaveReportPermission(c *gin.Context) {
	var req storeVisitReq.SaveReportPermissionReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.SaveReportPermission(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) SaveComment(c *gin.Context) {
	var req storeVisitReq.SaveCommentReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.SaveComment(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) CommentList(c *gin.Context) {
	var req storeVisitReq.CommentListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.CommentList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) DelComment(c *gin.Context) {
	var req storeVisitReq.DelCommentReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.DelComment(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) DashboardTree(c *gin.Context) {
	var req storeVisitReq.DashboardTreeReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.DashboardTree(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) ExportDashboardList(c *gin.Context) {
	var req storeVisitReq.DashboardExportEmail
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.ExportDashboardList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) SendPdfToEmail(c *gin.Context) {
	var req storeVisitReq.SendPdfToEmailReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	_, err := h.storeVisitSdk.SendPdfToEmail(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h StoreVisitHandle) GetConfig(c *gin.Context) {
	var req storeVisitReq.VisitConfigReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.storeVisitSdk.GetConfig(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitHandle) SaveLanguageConfig(c *gin.Context) {
	var req storeVisitReq.BaseReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.Language = c.Request.Header.Get("Accept-Language")
	res, err := h.storeVisitSdk.SaveLanguageConfig(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

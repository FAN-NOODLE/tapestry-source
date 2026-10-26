package storeVisitHandle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	storeVisitSvc "micro-scrm/internal/app/service/store_visit"
	"micro-scrm/internal/pkg/helper"
	"os"
)

type Handle struct {
	handle.BaseHandle
	mainSvc         *storeVisitSvc.MainSvc
	deptMemberSvc   *storeVisitSvc.DeptMemberSvc
	showRecordSvc   *storeVisitSvc.ShowRecordSvc
	visitCommentSvc *storeVisitSvc.CommentSvc
	mainAdminSvc    *storeVisitSvc.MainAdminSvc
}

func NewHandle() Handle {
	return Handle{
		mainSvc:         storeVisitSvc.NewMainSvc(),
		deptMemberSvc:   storeVisitSvc.NewDeptMemberSvc(),
		showRecordSvc:   storeVisitSvc.NewShowRecordSvc(),
		visitCommentSvc: storeVisitSvc.NewCommentSvc(),
		mainAdminSvc:    storeVisitSvc.NewMainAdminSvc(),
	}
}

func (h Handle) SavePlan(c *gin.Context) {
	var req storeVisitReq.SavePlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.SavePlan(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GetStoreMM(c *gin.Context) {
	var req storeVisitReq.StoreManagerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.GetStoreMM(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GetStoreManager(c *gin.Context) {
	var req storeVisitReq.StoreManagerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.GetStoreManager(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) StoreVisitList(c *gin.Context) {
	var req storeVisitReq.PlanListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mainSvc.StoreVisitList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) StoreVisitInfo(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.StoreVisitInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) DeptMemberList(c *gin.Context) {
	var req storeVisitReq.BaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.deptMemberSvc.DeptMemberList(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) GetEmailContact(c *gin.Context) {
	var req storeVisitReq.BaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.deptMemberSvc.GetEmailContact(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) SaveEmail(c *gin.Context) {
	var req storeVisitReq.EmailContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.deptMemberSvc.SaveEmail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) DelEmail(c *gin.Context) {
	var req storeVisitReq.EmailContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.deptMemberSvc.DelEmail(h.WarpContext(c), req)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) SearchMembers(c *gin.Context) {
	var req storeVisitReq.SearchMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.deptMemberSvc.SearchMembers(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) AddCopyMember(c *gin.Context) {
	var req storeVisitReq.AddCopyMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.AddCopyMembers(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) LoginUserInfo(c *gin.Context) {
	var req storeVisitReq.BaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.deptMemberSvc.LoginUserInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) StartVisit(c *gin.Context) {
	var req storeVisitReq.StartVisitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	var err error
	if req.CheckIn && req.ID > 0 {
		//打卡
		err = h.mainSvc.VisitCheckIn(h.WarpContext(c), req)
	} else {
		//开始巡店/直接巡店
		req.ID, err = h.mainSvc.StartVisit(h.WarpContext(c), req)
	}
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, req.ID)
}

func (h Handle) SaveShowRecordContent(c *gin.Context) {
	var req storeVisitReq.SaveShowRecordContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.showRecordSvc.SaveShowRecordContent(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//更新pdf
	visitInfo, _ := storeVisitDao.NewMainDao().GetOne(h.WarpContext(c), req.EnterpriseId, "id = ?", req.VisitID)
	if visitInfo.Status == storeVisitModel.VisitComplete && (visitInfo.Overview == "" || visitInfo.Content == "") {
		storeVisitSvc.NewCreatePDFQueueSvc().CreatePDFProduce(h.WarpContext(c), storeVisitReq.CreatePdfReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: req.EnterpriseId},
			ID:                req.VisitID,
			Language:          req.Language,
		})
	}
	h.ResSuccess(c, res)
}

func (h Handle) SaveShowRecordMedia(c *gin.Context) {
	var req storeVisitReq.SaveShowRecordMediaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.showRecordSvc.SaveShowRecordMedia(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//更新pdf
	updatePDF(h.WarpContext(c), req.EnterpriseId, req.VisitID, req.Language)
	h.ResSuccess(c, info)
}

func (h Handle) DelShowRecord(c *gin.Context) {
	var req storeVisitReq.DelShowRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.showRecordSvc.DelShowRecord(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	recordInfo, _ := storeVisitDao.NewShowRecordDao().GetOne(h.WarpContext(c), req.EnterpriseId, "id = ?", req.ID)
	//更新pdf
	updatePDF(h.WarpContext(c), req.EnterpriseId, recordInfo.VisitID, req.Language)
	h.ResOk(c)
}
func (h Handle) DelShowRecordMedia(c *gin.Context) {
	var req storeVisitReq.DelShowRecordMediaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.showRecordSvc.DelShowRecordMedia(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	mediaInfo, _ := storeVisitDao.NewMediaDao().GetOneUnscoped(h.WarpContext(c), req.EnterpriseId, "id = ?", req.MediaID)
	//更新pdf
	updatePDF(h.WarpContext(c), req.EnterpriseId, mediaInfo.VisitID, req.Language)
	h.ResOk(c)
}

func (h Handle) ShowRecordList(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.showRecordSvc.ShowRecordList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) CompleteVisit(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.CompleteVisit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//同步创建pdf
	_ = storeVisitSvc.NewCreatePDFQueueSvc().CreatePDFFile(h.WarpContext(c), storeVisitReq.CreatePdfReq{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: req.EnterpriseId},
		ID:                req.ID,
		Language:          req.Language,
	})
	h.ResOk(c)
}

func (h Handle) SendCompleteReport(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.SendCompleteReport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) SaveVisitInfo(c *gin.Context) {
	var req storeVisitReq.SaveVisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.SaveVisitInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//更新pdf
	updatePDF(h.WarpContext(c), req.EnterpriseId, req.ID, req.Language)
	h.ResOk(c)
}

func (h Handle) SaveReportPermission(c *gin.Context) {
	var req storeVisitReq.SaveReportPermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.SaveReportPermission(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) SaveComment(c *gin.Context) {
	var req storeVisitReq.SaveCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.visitCommentSvc.SaveComment(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) CommentList(c *gin.Context) {
	var req storeVisitReq.CommentListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.visitCommentSvc.CommentList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) DelComment(c *gin.Context) {
	var req storeVisitReq.DelCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.visitCommentSvc.DelComment(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) AdminStoreVisitList(c *gin.Context) {
	var req storeVisitReq.AdminStoreVisitListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mainAdminSvc.AdminStoreVisitList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) AdminStoreVisitClose(c *gin.Context) {
	var req storeVisitReq.VisitInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainAdminSvc.AdminStoreVisitClose(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) DeptNameList(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.deptMemberSvc.DeptNameList(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) AdminDashboard(c *gin.Context) {
	var req storeVisitReq.AdminStoreDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, sumVisitNum, err := h.mainAdminSvc.AdminDashboard(h.WarpContext(c), req, false)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{
		"total": total,
		"list":  list,
		"page":  req.Page,
		"size":  req.Size,
		"num":   sumVisitNum,
	})

}
func (h Handle) DashboardTree(c *gin.Context) {
	var req storeVisitReq.DashboardTreeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.DashboardTree(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) SendPdfToEmail(c *gin.Context) {
	var req storeVisitReq.SendPdfToEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.SendPdfToEmail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) ExportDashboardList(c *gin.Context) {
	var req storeVisitReq.DashboardExportEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.ExportDashboardList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) AdminStoreVisitExport(c *gin.Context) {
	var req storeVisitReq.AdminStoreVisitListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	fileName, err := h.mainAdminSvc.AdminStoreVisitExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	url, err := helper.UploadFileToCos(c, fileName, true)
	if err != nil {
		h.ResError(c, err)
		return
	}
	defer os.Remove(fileName)
	h.ResSuccess(c, gin.H{"url": url})
}

func (h Handle) AdminDashboardExport(c *gin.Context) {
	var req storeVisitReq.AdminStoreDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	fileName, err := h.mainAdminSvc.AdminDashboardExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	url, err := helper.UploadFileToCos(c, fileName, true)
	if err != nil {
		h.ResError(c, err)
		return
	}
	defer os.Remove(fileName)
	h.ResSuccess(c, gin.H{"url": url})
}

func updatePDF(c ctx.Context, enterpriseId, ID int64, language string) {
	if ID > 0 {
		visitInfo, _ := storeVisitDao.NewMainDao().GetOne(c, enterpriseId, "id = ?", ID)
		if visitInfo.Status == storeVisitModel.VisitComplete {
			storeVisitSvc.NewCreatePDFQueueSvc().CreatePDFProduce(c, storeVisitReq.CreatePdfReq{
				EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: enterpriseId},
				ID:                ID,
				Language:          language,
			})
		}
	}
}
func (h Handle) GetConfig(c *gin.Context) {
	var req storeVisitReq.VisitConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := storeVisitSvc.NewConfigSvc().GetConfig(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) SaveLanguageConfig(c *gin.Context) {
	var req storeVisitReq.BaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	storeVisitSvc.NewConfigSvc().SaveLanguageConfig(h.WarpContext(c), req)

	h.ResOk(c)
}

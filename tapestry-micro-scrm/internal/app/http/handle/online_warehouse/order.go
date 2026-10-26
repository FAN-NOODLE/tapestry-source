package onlineWarehouse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/handle"
	emailReq "micro-scrm/internal/app/http/request/email"
	onlineWarehouseReq "micro-scrm/internal/app/http/request/online_warehouse"
	"micro-scrm/internal/app/job/online_warehouse"
	departmentService "micro-scrm/internal/app/service/department"
	"micro-scrm/internal/app/service/email"
	onlineWarehouseSvc "micro-scrm/internal/app/service/online_warehouse"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"os"
)

type Handle struct {
	handle.BaseHandle
	orderSvc *onlineWarehouseSvc.OrderSvc
}

func NewHandle() Handle {
	return Handle{
		orderSvc: onlineWarehouseSvc.NewOrderSvc(),
	}
}

func (h Handle) List(c *gin.Context) {
	var req onlineWarehouseReq.ListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	if req.IsAll {
		thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}
	list, total, err := h.orderSvc.List(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) Submit(c *gin.Context) {
	var req onlineWarehouseReq.SubmitOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.orderSvc.Submit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//增加发货的时候补全商品信息
	go online_warehouse.NewSupplementSkuInfo().Handle(h.WarpContext(c), req.EnterpriseId)
	h.ResOk(c)
}

func (h Handle) Info(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.orderSvc.Info(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h Handle) Edit(c *gin.Context) {
	var req onlineWarehouseReq.SubmitOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.orderSvc.Edit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) StoreList(c *gin.Context) {
	var req onlineWarehouseReq.StoreListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.orderSvc.StoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) StoreOrderList(c *gin.Context) {
	var req onlineWarehouseReq.StoreOrderListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.orderSvc.StoreOrderList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) Export(c *gin.Context) {
	var req onlineWarehouseReq.ListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	fileName, err := h.orderSvc.ExportFile(h.WarpContext(c), req)
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

func (h Handle) ImportOrder(c *gin.Context) {
	var req onlineWarehouseReq.ImportOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.orderSvc.ImportOrder(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	//增加发货的时候补全商品信息
	go online_warehouse.NewSupplementSkuInfo().Handle(h.WarpContext(c), req.EnterpriseId)
	h.ResSuccess(c, result)
}

func (h Handle) BatchTransfer(c *gin.Context) {
	var req onlineWarehouseReq.BatchTransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.orderSvc.BatchTransfer(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) Refund(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.orderSvc.Refund(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) TipsNum(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	if req.IsAll {
		thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}
	res, err := h.orderSvc.GetTipsNum(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) Create(c *gin.Context) {
	var req onlineWarehouseReq.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	id, err := h.orderSvc.Create(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"id": id})
}

func (h Handle) StoreRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.orderSvc.StoreRank(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) SaRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	req.ThirdIds = thirdIds
	res, err := h.orderSvc.SaRank(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) SendFileByEmail(c *gin.Context) {
	var req onlineWarehouseReq.ExportMailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	if req.IsAll {
		thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}
	fileName, err := h.orderSvc.ExportFile(h.WarpContext(c), req.ListReq)
	if err != nil {
		h.ResError(c, err)
		return
	}
	defer os.Remove(fileName)
	err = email.NewSendEmail().SendEmailByDiy(emailReq.SendMailReq{
		Title:     fmt.Sprintf("网仓订单附件下载"),
		Content:   "",
		ToAddress: []string{req.Email},
		FileList: []emailReq.FileInfo{
			{
				Url:       "",
				Name:      "网仓订单附件",
				LocalPath: fileName,
			},
		},
	}, config.VchangyiEmailConfig{
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		Username: "storeoperation@vchangyi.com",
		Password: "REDACTED_FOR_SOURCE_SHARING",
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) CustomerOrderInfo(c *gin.Context) {
	var req onlineWarehouseReq.OrderInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = helper.RuleThree(app.Cfg().Env == "prod", int64(3360610327843968), int64(3363603667487424))
	info, err := h.orderSvc.CustomerOrderInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h Handle) CustomerOrderEdit(c *gin.Context) {
	var req onlineWarehouseReq.SaveReceiveInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = helper.RuleThree(app.Cfg().Env == "prod", int64(3360610327843968), int64(3363603667487424))
	err := h.orderSvc.SaveReceiveInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) LeaveSa(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	memberIds, err := h.orderSvc.GetAllMemberId(cc, req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	memberList, err := membersDao.NewMembers().GetUnscopedLeaveMemberInfoByIds(cc, req.EnterpriseId, memberIds)
	if err != nil {
		h.ResError(c, err)
		return
	}
	leaveMemberIds := make([]int64, 0)
	for _, info := range memberList {
		leaveMemberIds = append(leaveMemberIds, info.ID)
	}
	h.ResSuccess(c, leaveMemberIds)
}

func (h Handle) DepList(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	list, err := departmentService.NewDepartmentsService().GetDepList(cc, req.EnterpriseId, thirdIds)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)

}

func (h Handle) SaList(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	var thirdIds []int
	var err error
	if len(req.CheckIds) == 0 {
		thirdIds, err = h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
	} else {
		thirdIds = req.CheckIds
	}
	list, err := departmentService.NewDepartmentsService().GetSaList(cc, req.EnterpriseId, thirdIds)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)

}

func (h Handle) Cancel(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.orderSvc.Cancel(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) TransferExport(c *gin.Context) {
	var req onlineWarehouseReq.TransferExportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	transferFile, err := h.orderSvc.TransferExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	url, err := helper.UploadFileToCos(c, transferFile, false)
	if err != nil {
		h.ResError(c, err)
		return
	}
	defer os.Remove(transferFile)
	h.ResSuccess(c, gin.H{"url": url})
}

func (h Handle) TransferImport(c *gin.Context) {
	var req onlineWarehouseReq.TransferImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.orderSvc.TransferImport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) PosterMedia(c *gin.Context) {
	var req onlineWarehouseReq.OrderInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	mediaId, err := h.orderSvc.GetMediaId(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"media_id": mediaId})
}

func (h Handle) CheckOrderSn(c *gin.Context) {
	var req onlineWarehouseReq.CheckOrderSnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	orderInfo, err := h.orderSvc.CheckOrderSn(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	orderTime := ""
	if orderInfo.ID > 0 {
		orderTime = helper.GetFormatDateByParams(orderInfo.CreatedAt)
	}
	h.ResSuccess(c, gin.H{"order_time": orderTime, "order_sn": orderInfo.OrderSn})
}

func (h Handle) Dashboard(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	if req.IsAll {
		thirdIds, err := h.orderSvc.GetAllThirdId(cc, req.EnterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}
	res, err := h.orderSvc.Dashboard(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) GoodsRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.orderSvc.GoodsRank(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

package onlineWarehouse

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	onlineWarehouseReq "gw-work-scrm/internal/app/http/request/online_warehouse"
	"gw-work-scrm/internal/pkg/sdk/microscrm/online_warehouse"
)

type OrderHandle struct {
	handle.BaseHandle
	onlineWarehouseSdk *online_warehouse.Order
}

func NewOrderHandle() *OrderHandle {
	return &OrderHandle{
		onlineWarehouseSdk: online_warehouse.NewOrder(),
	}
}

func (h *OrderHandle) List(c *gin.Context) {
	var req onlineWarehouseReq.ListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	if memberInfo.Type == 4 { //店长
		if len(req.ThirdIds) == 0 && len(req.MemberIds) == 0 && len(req.LeaveMemberIds) == 0 {
			req.IsAll = true
		} else { // 如果有筛选查询，店长身份的话，那么这里的查询就不能有店长
			req.MemberId = 0
		}
	}
	req.IsH5 = true
	list, err := h.onlineWarehouseSdk.List(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) TipsNum(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if memberInfo.Type == 4 { //店长
		if len(req.ThirdIds) == 0 && len(req.MemberIds) == 0 && len(req.LeaveMemberIds) == 0 {
			req.IsAll = true
		} else { // 如果有筛选查询，店长身份的话，那么这里的查询就不能有店长
			req.MemberId = 0
		}
	}
	req.IsH5 = true
	res, err := h.onlineWarehouseSdk.TipsNum(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) Info(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	req.Origin = 1
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.onlineWarehouseSdk.Info(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) Create(c *gin.Context) {
	var req onlineWarehouseReq.CreateReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.onlineWarehouseSdk.Create(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) DcList(c *gin.Context) {
	var req onlineWarehouseReq.DcListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.onlineWarehouseSdk.DcList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) InventoryList(c *gin.Context) {
	var req onlineWarehouseReq.InventoryListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.onlineWarehouseSdk.InventoryList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) StoreRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	req.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.StoreRank(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) SaRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.SaRank(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) ExportMail(c *gin.Context) {
	var req onlineWarehouseReq.ExportMailReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.IsH5 = true
	cc := ctx.Wrap(c)
	if memberInfo.Type == 4 { //店长
		if len(req.ThirdIds) == 0 && len(req.MemberIds) == 0 {
			req.IsAll = true
		} else { // 如果有筛选查询，店长身份的话，那么这里的查询就不能有店长
			req.MemberId = 0
		}
	}
	list, err := h.onlineWarehouseSdk.ExportMail(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) CustomerInfo(c *gin.Context) {
	var req onlineWarehouseReq.OrderInfoReq
	req.EnterpriseId = 1
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.CustomerInfo(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) CustomerEdit(c *gin.Context) {
	var req onlineWarehouseReq.SaveReceiveInfoReq
	req.EnterpriseId = 1
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.CustomerEdit(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) LeaveSa(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	list, err := h.onlineWarehouseSdk.LeaveSa(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) SaList(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	list, err := h.onlineWarehouseSdk.SaList(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) DepList(c *gin.Context) {
	var req onlineWarehouseReq.LeaveSaReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	list, err := h.onlineWarehouseSdk.DepList(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) Cancel(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.Cancel(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) PosterMedia(c *gin.Context) {
	var req onlineWarehouseReq.OrderInfoReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.PosterMedia(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h *OrderHandle) CheckOrder(c *gin.Context) {
	var req onlineWarehouseReq.CheckOrderReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.CheckOrderSn(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)

}

func (h *OrderHandle) Dashboard(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if memberInfo.Type == 4 { //店长
		if len(req.ThirdIds) == 0 && len(req.MemberIds) == 0 && len(req.LeaveMemberIds) == 0 {
			req.IsAll = true
		} else { // 如果有筛选查询，店长身份的话，那么这里的查询就不能有店长
			req.MemberId = 0
		}
	}
	req.IsH5 = true
	res, err := h.onlineWarehouseSdk.Dashboard(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *OrderHandle) GoodsRank(c *gin.Context) {
	var req onlineWarehouseReq.DashboardReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.onlineWarehouseSdk.GoodsRank(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

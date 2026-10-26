package reserve

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/app/http/common"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/reserve"
	microScrmReserve "gw-work-scrm/internal/pkg/sdk/microscrm/reserve"
	"gw-work-scrm/internal/pkg/validate_params"
)

type OrderHandle struct {
	handle.BaseHandle
	ReserveSdk *microScrmReserve.OrderApi
}

func NewOrderHandle() OrderHandle {
	return OrderHandle{
		BaseHandle: handle.BaseHandle{},
		ReserveSdk: microScrmReserve.NewOrderApi(),
	}
}

// CreateOrder 创建邀约单
func (h OrderHandle) CreateOrder(c *gin.Context) {
	var request reserve.OrderInfoReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		errs := err.(ginplus.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "UserUnionid" && e.Tag() == "required" {
				h.ResError(c, errorx.New("该用户身份不是微信用户"))
			}
		}
		return
	}
	order, err := h.ReserveSdk.CreateOrder(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, order)
}

// OrderList 顾客详情查询邀约单列表
func (h OrderHandle) OrderList(c *gin.Context) {
	var request reserve.OrderListReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	request.MemberIds = append(request.MemberIds, memberInfo.ID)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if memberInfo.Type == 4 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(ctx.Wrap(c), request.EnterpriseId, request.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		request.StoreIds = thirdIds
		request.MemberId = memberInfo.ID // 店长需要带上自己的查询邀约单
		request.IsSM = 1
	}
	request.OrderSort = 1
	lists, err := h.ReserveSdk.OrderLists(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, lists)
}

// TimeList 导购的时间段邀约数
func (h OrderHandle) TimeList(c *gin.Context) {
	var request reserve.OrderTimeListReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	request.MemberId = memberInfo.ID
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	lists, err := h.ReserveSdk.TimeList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, lists)
}

func (h OrderHandle) ServiceList(c *gin.Context) {
	var request reserve.ReserveServiceListByMemberReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.Status = 1
	lists, err := h.ReserveSdk.MemberServiceLists(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, lists)
}

// OrderShare 邀约单分享上报
func (h OrderHandle) OrderShare(c *gin.Context) {
	var request reserve.OrderBaseReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.ReserveSdk.OrderShare(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, info)
}

// MemberDetail 用户的会员信息
func (h OrderHandle) MemberDetail(c *gin.Context) {
	var request reserve.MemberDetailReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := h.ReserveSdk.MemberDetail(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, detail)
}

// OrderDelete 取消分析小程序卡片，删除创建的邀约单
func (h OrderHandle) OrderDelete(c *gin.Context) {
	var request reserve.OrderBaseReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.ReserveSdk.OrderDelete(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, info)
}

// ChangeStatus 导购确认订单是否到店购买,设置新客老客
func (h OrderHandle) ChangeStatus(c *gin.Context) {
	var request reserve.ChangeStatus
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.ReserveSdk.OrderChangStatus(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

// CheckReserveTime SW&&ks验证同一时间段的邀约限制
func (h OrderHandle) CheckReserveTime(c *gin.Context) {
	var request reserve.CheckReserveTime
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.ReserveSdk.CheckReserveTime(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

// VerifyShare 验证导购是否可以分享邀约单
func (h OrderHandle) VerifyShare(c *gin.Context) {
	var request reserve.OrderAssignReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.ReserveSdk.VerifyShare(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h OrderHandle) OrderShow(c *gin.Context) {
	var request reserve.OrderShowReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.ReserveSdk.OrderShow(ctx.Wrap(c), request)
	if err != nil {
		return
	}
	h.ResSuccess(c, info)
}

func (h OrderHandle) CustomerReserveLists(c *gin.Context) {
	var request reserve.CustomerReserveListsReq
	request.EnterpriseId = h.EnterpriseID(c)
	if !validate_params.ValidateParams(c, &request) {
		return
	}
	lists, err := h.ReserveSdk.CustomerReserveLists(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, lists)
}

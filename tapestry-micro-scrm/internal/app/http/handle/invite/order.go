package invite

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	excitationModels "micro-scrm/internal/app/models/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	inviteSvc "micro-scrm/internal/app/service/invite"
	membersSvc "micro-scrm/internal/app/service/member"
	reserveSvc "micro-scrm/internal/app/service/reserve"
)

type OrderHandle struct {
	handle.BaseHandle
	templateService *inviteSvc.InviteTemplateService
	reserveSvc      *reserveSvc.OrderService
	membersService  *membersSvc.MembersService
}

func NewOrderHandle() OrderHandle {
	return OrderHandle{
		templateService: inviteSvc.NewInviteTemplateService(),
		reserveSvc:      reserveSvc.NewOrderService(),
		membersService:  membersSvc.NewMembersService(),
	}
}

// CreateRecord 发送邀约单记录并生成订单
func (h OrderHandle) CreateRecord(c *gin.Context) {
	var request inviteRequest.RecordInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.templateService.CreateRecord(h.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	result := inviteRequest.RecordResult{
		InviteRecordId: info.ID,
	}

	go func() {
		errs := func() (err error) {

			templateInfo, err := h.templateService.GetRecordInfo(h.WarpContext(c), request.EnterpriseId, result.InviteRecordId)
			if err != nil {
				return err
			}
			// 用户信息
			// 导购门店信息
			memberDetail, err := h.membersService.GetDetail(h.WarpContext(c), request.EnterpriseId, request.Mobile)
			if err != nil {
				return err
			}

			var orderInfo reserveRequest.OrderInfoReq
			orderInfo.EnterpriseId = request.EnterpriseId
			orderInfo.ReserveId = templateInfo.ReserveServiceID
			orderInfo.UserName = templateInfo.UserName
			orderInfo.UserUnionid = templateInfo.UserUnionid
			orderInfo.UserExternalId = templateInfo.UserOpenid

			res, err := h.reserveSvc.CreateInviteOrder(h.WarpContext(c), request.EnterpriseId, result.InviteRecordId, orderInfo.OrderInfo, memberDetail)
			if err != nil {
				return err
			}

			go reserveSvc.NewReserveLog().CreateLog(h.WarpContext(c), request.EnterpriseId, res.ID, "导购生成邀请函订单", reserveRequest.ReserveOperaType_BA)
			return nil
		}()
		if errs != nil {
			h.Log().Info(h.WarpContext(c)).Msgf("发送邀约单记录并生成订单错误日志%v", errs)
		}
	}()

	h.ResSuccess(c, result)

	return
}

// GetRecordInfo 用户查看邀请函信息
func (h OrderHandle) GetRecordInfo(c *gin.Context) {
	var request inviteRequest.GetRecordInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	var result inviteRequest.RecordDetailInfoRes
	info, err := h.templateService.GetRecordInfo(h.WarpContext(c), request.EnterpriseId, request.InviteRecordId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if request.UnionId != "" && info.UserUnionid != request.UnionId {
		h.ResError(c, errors.New("无权访问"))
		return
	}
	err = json.Unmarshal([]byte(info.TimeSlotStr), &result.TimeSlot)
	memberInfo, err := h.membersService.GetInfoByUserId(h.WarpContext(c), request.EnterpriseId, info.MemberUserid)
	if err != nil {
		h.ResError(c, err)
		return
	}
	info.MemberName = memberInfo.Name

	var orderReq inviteRequest.CreateOrderReq
	orderReq.ReserveOrderId = info.ReserveOrderId
	err = h.reserveSvc.UpdateInviteOrder(h.WarpContext(c), request.EnterpriseId, orderReq, 0)
	if err != nil {
		h.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(h.WarpContext(c), request.EnterpriseId, request.InviteRecordId, "顾客查看邀请函信息", reserveRequest.ReserveOperaType_Customer)

	orderInfo, err := h.reserveSvc.GetOrderInfo(h.WarpContext(c), request.EnterpriseId, info.ReserveOrderId)
	if err != nil {
		h.ResError(c, err)
		return
	}

	// 获取当前员工的门店编码
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(h.WarpContext(c), map[string]any{
		"enterprise_id":        request.EnterpriseId,
		"wechat_department_id": orderInfo.StoreID,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}

	info.ReserveTime = orderInfo.ReserveTime
	result.ReserveEndTime = orderInfo.ReserveEndTime
	if orderInfo.ReserveTime == 0 {
		info.ReserveOrderId = 0
	}
	info.StoreId = orderInfo.StoreID
	info.StoreName = storeDepartmentInfo.StoreName
	info.StoreCode = storeDepartmentInfo.ThirdStoreCode
	result.RecordDetailInfo = info
	h.ResSuccess(c, result)
	return
}

// CreateInviteOrder 用户确认邀请函订单
func (h OrderHandle) CreateInviteOrder(c *gin.Context) {

	var request inviteRequest.CreateOrderReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	templateInfo, err := h.templateService.GetRecordInfo(h.WarpContext(c), request.EnterpriseId, request.InviteRecordId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if templateInfo.UserUnionid != request.UnionId {
		h.ResError(c, errors.New("无权访问"))
		return
	}

	orderInfo, err := h.reserveSvc.GetOrderInfo(h.WarpContext(c), request.EnterpriseId, templateInfo.ReserveOrderId)

	if orderInfo.ReserveTime > 0 {
		h.ResError(c, errors.New("邀请函已确认"))
		return
	}

	request.ReserveOrderId = templateInfo.ReserveOrderId
	err = h.reserveSvc.UpdateInviteOrder(h.WarpContext(c), request.EnterpriseId, request, 1)
	if err != nil {
		h.ResError(c, err)
		return
	}

	go reserveSvc.NewReserveLog().CreateLog(h.WarpContext(c), request.EnterpriseId, templateInfo.ReserveOrderId, "顾客确认邀请函订单", reserveRequest.ReserveOperaType_Customer)

	var resData inviteRequest.CreateOrderRrs
	resData.Id = orderInfo.ID
	resData.ReserveNumber = orderInfo.ReserveNumber
	// 积分队列埋点,分享邀约才算
	excitationSvc.NewPointsProduceService().ProduceQueue(h.WarpContext(c), excitationReq.PointsProduceSave{
		MemberId:     orderInfo.MemberID,
		Type:         excitationModels.RuleCategoryId1,
		Value:        1,
		UniqueHash:   orderInfo.ReserveNumber + fmt.Sprintf("-%d-%d", orderInfo.MemberID, excitationModels.RuleCategoryId1),
		EnterpriseId: request.EnterpriseId,
	})
	h.ResSuccess(c, resData)

	return
}

// RecordOrderList 邀请数据列表
func (h OrderHandle) RecordOrderList(c *gin.Context) {

	var reqList inviteRequest.RecordOrderListReq
	if err := c.ShouldBind(&reqList); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.templateService.GetRecordOrderList(h.WarpContext(c), reqList.EnterpriseId, reqList)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, reqList.Page, reqList.Size, result)
	return

}

func (h OrderHandle) RecordDelete(c *gin.Context) {
	var request inviteRequest.GetRecordInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.templateService.RecordDelete(h.WarpContext(c), request.EnterpriseId, request.InviteRecordId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

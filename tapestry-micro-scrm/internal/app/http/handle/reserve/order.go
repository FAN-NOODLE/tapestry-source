package reserve

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/handle"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	"micro-scrm/internal/app/http/request/external"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	externalUserService "micro-scrm/internal/app/service/external_user"
	membersSvc "micro-scrm/internal/app/service/member"
	reserveSvc "micro-scrm/internal/app/service/reserve"
	"micro-scrm/internal/pkg/helper"
	"strconv"
)

type OrderHandle struct {
	handle.BaseHandle

	reserveSvc *reserveSvc.OrderService
}

func NewOrderHandle() OrderHandle {
	return OrderHandle{
		reserveSvc: reserveSvc.NewOrderService(),
	}
}

// CreateOrder 创建预约单
func (r OrderHandle) CreateOrder(c *gin.Context) {
	var request reserveRequest.OrderInfoReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.CreateOrder(r.WarpContext(c), request.EnterpriseId, request.OrderInfo)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, orderInfo.ID, "到店邀约", reserveRequest.ReserveOperaType_BA)
	result := reserveRequest.OrderBase{
		ID:            orderInfo.ID,
		ReserveNumber: orderInfo.ReserveNumber,
	}

	r.ResSuccess(c, result)
}

// EditOrder 编辑预约单
func (r OrderHandle) EditOrder(c *gin.Context) {
	var request reserveRequest.MemberEditOrder
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.Update(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.ID, "导购编辑预约单", reserveRequest.ReserveOperaType_BA)
	r.ResOk(c)
}

// CancelOrder 导购企微端取消预约单
func (r OrderHandle) CancelOrder(c *gin.Context) {
	var request reserveRequest.OrderBaseReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.CancelOrder(r.WarpContext(c), request.EnterpriseId, request.OrderBase)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "导购取消预约单", reserveRequest.ReserveOperaType_BA)
	r.ResOk(c)
}

// OrderList 预约单列表
func (r OrderHandle) OrderList(c *gin.Context) {
	var request reserveRequest.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	var result []reserveRequest.OrderListResult
	var total int64
	var err error
	if request.Search.IsSearch == 1 { // 针对于任务中心预约列表的筛选
		unionIdsAndMemberId, err := r.reserveSvc.GetSearchUnionid(r.WarpContext(c), request)
		if err != nil {
			r.ResError(c, err)
			return
		}
		if len(unionIdsAndMemberId) > 0 {
			result, total, err = r.reserveSvc.GetList(r.WarpContext(c), request, unionIdsAndMemberId)
		}
	} else {
		result, total, err = r.reserveSvc.GetList(r.WarpContext(c), request, nil)
	}
	if err != nil {
		r.ResError(c, err)
		return
	}
	if len(result) > 0 {
		var unionIds []string
		for _, info := range result {
			unionIds = append(unionIds, info.UserUnionid)
		}
		coachNameList, _ := externalUserService.NewWxUserCorpExternal().GetNickNameAndRemarkByUnionIds(r.WarpContext(c), request.EnterpriseId, request.MemberId, unionIds)
		var coachNameMap = make(map[string]external.CoachNameInfo)
		for _, info := range coachNameList {
			coachNameMap[info.Unionid] = info
		}
		for i, row := range result {
			if val, ok := coachNameMap[row.UserUnionid]; ok {
				result[i].UserName = helper.RuleThree(val.Remark != "", val.Remark, val.Nick)
			}
		}
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

// OrderListByCoach  预约单列表
func (r OrderHandle) OrderListByCoach(c *gin.Context) {
	var request reserveRequest.OrderListByCoachReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	var result = make([]reserveRequest.OrderListResult, 0)
	var total int64
	var err error
	cc := r.WarpContext(c)
	if request.IsSearch { // 针对于任务中心预约列表的筛选
		request.Search.MemberIds = request.MemberIds
		request.Search.EnterpriseID = request.EnterpriseId
		var externalRequest = customerSearchReq.ExternalSearchListBaseReq{
			ThirdIds:   request.ThirdIds,
			SearchList: request.Search,
		}
		unionIds, err := customerSvc.NewExternalSearch().GetExternalSearchCk(cc, externalRequest, "user_unionid")
		if err != nil {
			r.ResError(c, err)
			return
		}
		if len(unionIds) == 0 {
			r.ResList(c, 0, request.Page, request.Size, result)
			return
		}
		for _, v := range unionIds {
			// 使用类型断言将 interface{} 转换为 int64
			if val, ok := v.(string); ok {
				request.UnionIds = append(request.UnionIds, val)
			}
		}
	}
	result, total, err = r.reserveSvc.GetListByCoach(cc, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

// OrderDetail 预约单详情
func (r OrderHandle) OrderDetail(c *gin.Context) {
	var request reserveRequest.OrderDetailReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	result, err := r.reserveSvc.GetOrderDetail(r.WarpContext(c), request.EnterpriseId, request.ID, request.EnterpriseName)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResSuccess(c, result)
	return
}

// UserOrderInfo 微信用户端获取预约单信息
func (r OrderHandle) UserOrderInfo(c *gin.Context) {
	var request reserveRequest.UserOrderBase
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.GetOrderInfo(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	if orderInfo.UserUnionid != request.UnionId {
		r.ResError(c, errors.New("无权访问"))
		return
	}
	if orderInfo.IsOpen != 1 {
		_ = r.reserveSvc.OpenOrder(r.WarpContext(c), request.EnterpriseId, orderInfo)
	}
	serveInfo, _ := reserveSvc.NewReserveServiceService().GetInfoById(r.WarpContext(c), request.EnterpriseId, orderInfo.ReserveID)
	serviceScore, _ := reserveSvc.NewServiceScoreService().ScoreInfo(r.WarpContext(c), request.EnterpriseId, reserveRequest.UserOrderBase{
		OrderBase: reserveRequest.OrderBase{
			ID:            orderInfo.ID,
			ReserveNumber: orderInfo.ReserveNumber,
		},
		EnterpriseBase: reserveRequest.EnterpriseBase{
			EnterpriseId: request.EnterpriseId,
		},
		UnionId: orderInfo.UserUnionid,
	})
	result := reserveRequest.OrderInfoRes{
		ReserveOrderModel: orderInfo,
		ServeName:         serveInfo.ReserveName,
		ServeImg:          serveInfo.ReserveImg,
		ServeTime:         serveInfo.ReserveTime,
		ServiceScore:      serviceScore.ServiceScore,
		StoreScore:        serviceScore.StoreScore,
		Remark:            serviceScore.Remark,
	}
	r.ResSuccess(c, result)
}

// UserCancelOrder  用户微信端取消预约单
func (r OrderHandle) UserCancelOrder(c *gin.Context) {
	var request reserveRequest.UserOrderBase
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.CancelOrder(r.WarpContext(c), request.EnterpriseId, request.OrderBase)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "用户取消预约单", reserveRequest.ReserveOperaType_Customer)
	r.ResOk(c)
}

// UserEditOrder 用户编辑订单
func (r OrderHandle) UserEditOrder(c *gin.Context) {
	var request reserveRequest.UserEditOrder
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.UserEditOrder(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	if request.Status == 1 { // 用户确认预约单操作
		go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "用户确认预约单信息", reserveRequest.ReserveOperaType_Customer)
	}
	if request.Status == 2 { // 用户取消预约单操作
		go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "用户取消预约单", reserveRequest.ReserveOperaType_Customer)
	}
	r.ResOk(c)
}

// UserMakeScore 用户评分
func (r OrderHandle) UserMakeScore(c *gin.Context) {
	var request reserveRequest.UserScoreReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.GetOrderInfo(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	if orderInfo.UserUnionid != request.UnionId {
		r.ResError(c, errors.New("无权访问"))
		return
	}
	_, err = reserveSvc.NewServiceScoreService().UserMakeScore(r.WarpContext(c), request.EnterpriseId, request, orderInfo)
	if err != nil {
		r.ResError(c, errors.New("评分失败"))
	} else {
		go reserveSvc.NewReserveLog().CreateScoreLog(r.WarpContext(c), request, "用户填写邀约评价", 2)
		go r.reserveSvc.UserCommentOrder(r.WarpContext(c), request.EnterpriseId, request.UserOrderBase)
	}
	r.ResOk(c)
}

// OrderAssign 分配预约单
func (r OrderHandle) OrderAssign(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.OrderAssign(r.WarpContext(c), request.ID, request.EnterpriseId, request.Mobile)

	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "店长分配预约单给导购", reserveRequest.ReserveOperaType_BA)
	r.ResOk(c)
	return
}

// GetReserveTimeList 获取预约时间段
func (r OrderHandle) GetReserveTimeList(c *gin.Context) {
	var request reserveRequest.OrderTimeListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	if request.Mobile == "" && request.MemberId == 0 {
		r.ResError(c, errors.New("导购手机号不能为空"))
		return
	}
	result, err := r.reserveSvc.GetReserveTimeList(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	res := make(map[int]int64, 0)
	if len(result) > 0 {
		for _, v := range result {
			res[v.ReserveTime] = v.Num
		}
	}
	r.ResSuccess(c, res)
	return
}

// ReserveOrderList 顾客受邀单列表
func (r OrderHandle) ReserveOrderList(c *gin.Context) {
	var request reserveRequest.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(r.WarpContext(c), request.EnterpriseId, "store_name")
	if err != nil {
		r.ResError(c, err)
		return
	}
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}

	nums, result, total, err := r.reserveSvc.GetReserveOrderList(r.WarpContext(c), request, isCoach)
	if err != nil {
		r.ResError(c, err)
		return
	}
	if total == 0 {
		result = []reserveRequest.OrderDetailResBaseStr{}
	}
	type Res struct {
		Nums  interface{} `json:"nums"`
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
	}

	res := Res{
		Nums:  nums,
		Total: total,
		Page:  request.Page,
		Size:  request.Size,
		List:  result,
	}
	r.ResSuccess(c, res)
	return
}

func (r OrderHandle) RankList(c *gin.Context) {
	var request reserveRequest.OrderRankList
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	result, err := r.reserveSvc.GetRankList(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResSuccess(c, result)
}

// Share 分享预约单
func (r OrderHandle) Share(c *gin.Context) {
	var request reserveRequest.OrderBaseReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.Share(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResOk(c)
	return
}

// ScoreInfo 用户端的评分
func (r OrderHandle) ScoreInfo(c *gin.Context) {
	var request reserveRequest.UserOrderBase
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.GetOrderInfo(r.WarpContext(c), request.EnterpriseId, request.ID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	if orderInfo.UserUnionid != request.UnionId {
		r.ResError(c, errors.New("无权访问"))
		return
	}
	result, err := reserveSvc.NewServiceScoreService().ScoreInfo(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	serviceInfo, err := reserveSvc.NewReserveServiceService().GetInfoById(r.WarpContext(c), request.EnterpriseId, orderInfo.ReserveID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	res := reserveRequest.ScoreInfoRes{
		ReserveOrderModel: orderInfo,
		IsScore:           1,
		Remark:            result.Remark,
		ServiceScore:      result.ServiceScore,
		StoreScore:        result.StoreScore,
		FlowScore:         result.FlowScore,
		ServeName:         serviceInfo.ReserveName,
	}
	if result.ID == 0 {
		res.IsScore = 0
	}
	r.ResSuccess(c, res)
	return
}

func (r OrderHandle) OrderInfo(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.GetOrderPreloadDetail(r.WarpContext(c), request.EnterpriseId, request.ID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	serveInfo := orderInfo.ReserveService
	members, err := membersSvc.NewMembersService().GetInfoByMobile(r.WarpContext(c), request.EnterpriseId, request.Mobile) // 获取当前用户的身份
	isOwn := 0
	if orderInfo.MemberID == members.ID {
		isOwn = 1
	}
	memberInfo, err := membersSvc.NewMembersService().GetInfoById(r.WarpContext(c), request.EnterpriseId, orderInfo.MemberID) // 获取订单用户的身份
	if err != nil {
		r.ResError(c, err)
		return
	}
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(r.WarpContext(c), microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: request.EnterpriseId,
		},
		WechatDepartmentID: uint(orderInfo.StoreID),
	})
	if err != nil {
		r.ResError(c, err)
		return
	}
	var serviceScore, storeScore int32
	var remark string
	if orderInfo.IsComment == 1 {
		scoreInfo, _ := reserveSvc.NewServiceScoreService().ScoreInfo(r.WarpContext(c), request.EnterpriseId, reserveRequest.UserOrderBase{
			OrderBase: reserveRequest.OrderBase{
				ID:            orderInfo.ID,
				ReserveNumber: orderInfo.ReserveNumber,
			},
			EnterpriseBase: reserveRequest.EnterpriseBase{
				EnterpriseId: request.EnterpriseId,
			},
			UnionId: orderInfo.UserUnionid,
		})
		if scoreInfo.ID > 0 {
			serviceScore = scoreInfo.ServiceScore
			storeScore = scoreInfo.StoreScore
			remark = scoreInfo.Remark
		}
	}

	customerStatusStr := ""
	switch orderInfo.CustomerStatus {
	case 1:
		customerStatusStr = "新客"
	case 2:
		customerStatusStr = "老客"
	case 4:
		customerStatusStr = "沉睡客"
	}
	if orderInfo.IsVip == 1 {
		customerStatusStr += ",VIP客户"
	}
	relationOrderNum := len(orderInfo.ReserveOrderSales)
	result := reserveRequest.MemberOrderInfoRes{
		ReserveOrderModel: orderInfo,
		ServeImg:          serveInfo.ReserveImg,
		ServeName:         serveInfo.ReserveName,
		ServeTime:         serveInfo.ReserveTime,
		Mobile:            memberInfo.Mobile,
		UcDepartmentId:    storeDepartmentInfo.UcDepartmentID,
		ServiceScore:      serviceScore,
		StoreScore:        storeScore,
		Remark:            remark,
		IsOwn:             isOwn,
		CustomerStatusStr: customerStatusStr,
		RelationOrderNum:  relationOrderNum,
	}
	r.ResSuccess(c, result)
	return
}

// Finish 完成预约
func (r OrderHandle) Finish(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	err := r.reserveSvc.Finish(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}

	r.createLog(c, request.OrderBaseReq, "完成预约", reserveRequest.ReserveOperaType_Customer)
	r.ResOk(c)
	return
}

func (r OrderHandle) createLog(c *gin.Context, req reserveRequest.OrderBaseReq, eventName string, operaType int) {
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), req.EnterpriseId, req.OrderBase.ID, eventName, operaType)
}

// UserSignIn SW端用户增加签到按钮
func (r OrderHandle) UserSignIn(c *gin.Context) {
	var request reserveRequest.UserSignIn
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.UserSignIn(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "用户确认签到", reserveRequest.ReserveOperaType_Customer)
	r.ResOk(c)
}

// ChangeStatus 导购确认订单是否到店购买,设置新客老客
func (r OrderHandle) ChangeStatus(c *gin.Context) {
	var request reserveRequest.ChangeStatus
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	orderInfo, err := r.reserveSvc.GetOrderInfo(r.WarpContext(c), request.EnterpriseId, request.ID)
	if err != nil {
		r.ResError(c, err)
		return
	}
	members, err := membersSvc.NewMembersService().GetInfoByMobile(r.WarpContext(c), request.EnterpriseId, request.Mobile) // 获取当前用户的身份

	if orderInfo.MemberID != members.ID {
		r.ResError(c, errors.New("非本人邀约单，暂不可操作哦~"))
		return
	}

	err = r.reserveSvc.ChangeStatus(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	var addPointsType int64
	if request.CustomerStatus > 0 {
		addPointsType = excitationModels.RuleCategoryId17
	}
	if request.IsStore > 0 {
		addPointsType = excitationModels.RuleCategoryId18
	}
	if addPointsType > 0 {
		// 积分队列埋点
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			Mobile:       request.Mobile,
			Type:         addPointsType,
			Value:        1,
			UniqueHash:   strconv.FormatInt(request.ID, 10) + fmt.Sprintf("-%s-%d", request.Mobile, addPointsType),
			EnterpriseId: request.EnterpriseId,
		})
	}
	r.ResOk(c)
}

// ConfirmStore Coach导购确认用户到店
func (r OrderHandle) ConfirmStore(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	cc := r.WarpContext(c)

	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(cc, request.EnterpriseId, "store_name")
	if err != nil {
		r.ResError(c, err)
		return
	}
	isCoach := true
	if config.Name != "coach" {
		isCoach = false
	}
	err = r.reserveSvc.ConfirmStore(cc, request.EnterpriseId, request, isCoach)
	if err != nil {
		r.ResError(c, err)
		return
	}

	// 积分队列埋点
	excitationSvc.NewPointsProduceService().ProduceQueue(cc, excitationReq.PointsProduceSave{
		Mobile:       request.Mobile,
		Type:         excitationModels.RuleCategoryId6,
		Value:        1,
		UniqueHash:   request.ReserveNumber + fmt.Sprintf("-%s-%d", request.Mobile, excitationModels.RuleCategoryId6),
		EnterpriseId: request.EnterpriseId,
	})

	r.ResOk(c)
}

// CommentShare 导购发送邀约评价
func (r OrderHandle) CommentShare(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.CommentShare(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	go reserveSvc.NewReserveLog().CreateShareScoreLog(r.WarpContext(c), request.EnterpriseId, request.OrderBase.ID, "发送邀约评价", 1)
	r.ResOk(c)
}

func (r OrderHandle) UserEditScore(c *gin.Context) {
	var request reserveRequest.UserScoreReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	err := reserveSvc.NewServiceScoreService().UserEditScore(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
	}

	r.ResOk(c)
}

func (r OrderHandle) VerifyShare(c *gin.Context) {
	var request reserveRequest.OrderAssignReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.VerifyShare(r.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResOk(c)
}

func (r OrderHandle) Delete(c *gin.Context) {
	var request reserveRequest.OrderBaseReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	// err := r.reserveSvc.DeleteOrder(r.WarpContext(c), request.EnterpriseId, request.OrderBase)
	// if err != nil {
	// 	r.ResError(c, err)
	// 	return
	// }
	r.ResOk(c)
}

// OrderExportList 管理后台预约单列表导出
func (r OrderHandle) OrderExportList(c *gin.Context) {
	var request reserveRequest.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	var result []reserveRequest.OrderExportListResult
	var total int64
	var err error

	result, total, err = r.reserveSvc.GetExportList(r.WarpContext(c), request)

	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

func (r OrderHandle) CheckReserveTime(c *gin.Context) {
	var request reserveRequest.CheckReserveTime
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(r.WarpContext(c), request.EnterpriseId, "store_name")
	if err != nil {
		r.ResError(c, err)
		return
	}
	check, msg, err := r.reserveSvc.CheckReserveTime(r.WarpContext(c), request, resConfig.Name)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResSuccess(c, gin.H{"check": check, "msg": msg})
}

// CustomerReserveLists 获取用户的所有邀约单列表
func (r OrderHandle) CustomerReserveLists(c *gin.Context) {
	var request reserveRequest.CustomerReserveListsReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	result, total, err := r.reserveSvc.GetCustomerList(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

func (r OrderHandle) Show(c *gin.Context) {
	var request reserveRequest.OrderShowReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	err := r.reserveSvc.ShowOrder(r.WarpContext(c), request.EnterpriseId, request.OrderBase)
	if err != nil {
		r.ResError(c, err)
		return
	}
	// 积分队列埋点
	excitationSvc.NewPointsProduceService().ProduceQueue(r.WarpContext(c), excitationReq.PointsProduceSave{
		MemberId:     request.MemberId,
		Type:         excitationModels.RuleCategoryId1,
		Value:        1,
		UniqueHash:   request.ReserveNumber + fmt.Sprintf("-%d-%d", request.MemberId, excitationModels.RuleCategoryId1),
		EnterpriseId: request.EnterpriseId,
	})
	r.ResOk(c)
}

func (r OrderHandle) ReserveOrderListNew(c *gin.Context) {
	var request reserveRequest.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	result, total, err := r.reserveSvc.GetReserveOrderListNew(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

func (r OrderHandle) OrderListCk(c *gin.Context) {
	var request reserveRequest.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}
	var result []reserveRequest.OrderListResult
	var total int64
	var err error
	result, total, err = r.reserveSvc.GetListCK(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

func (r OrderHandle) ExportOrderRecord(c *gin.Context) {
	var request reserveRequest.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		r.ResError(c, err)
		return
	}

	result, total, err := r.reserveSvc.ExportOrderRecord(r.WarpContext(c), request)
	if err != nil {
		r.ResError(c, err)
		return
	}
	r.ResList(c, total, request.Page, request.Size, result)
	return
}

package reserveSvc

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	customerReq "micro-scrm/internal/app/http/request/customer"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	memberReq "micro-scrm/internal/app/http/request/member"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

type OrderService struct {
	service.Base
	orderDao          *externalUserDao.ReservesOrder
	ReserveServiceDao *externalUserDao.ReserveServiceDao
	ReserveLogDao     *externalUserDao.ReservesLog
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderDao:          externalUserDao.NewReservesOrder(),
		ReserveServiceDao: externalUserDao.NewReserveServiceDao(),
		ReserveLogDao:     externalUserDao.NewReserveLog(),
	}
}

func (s *OrderService) GetList(c app.Ctx, request reserveRequest.OrderListReq, unionids []customerReq.UserUnionIdAndMemberId) ([]reserveRequest.OrderListResult, int64, error) {
	res, total, err := s.orderDao.GetList(c, request.EnterpriseId, request.Page, request.Size, request, unionids)
	if err != nil {
		return []reserveRequest.OrderListResult{}, total, err
	}
	// 查询服务信息
	// var listReq reserveRequest.ReserveServiceListReq
	// serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, request.EnterpriseId, listReq)
	// if errs != nil {
	// 	return []reserveRequest.OrderListResult{}, total, errs
	// }
	// serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	// for _, item := range serviceList {
	// 	serviceListMap[item.ID] = item
	// }
	result := make([]reserveRequest.OrderListResult, 0)
	var userUnionIds []string
	for _, v := range res {
		userUnionIds = append(userUnionIds, v.UserUnionid)
	}
	many, err := externalUserDao.NewWxUserCorpExternal().GetMany(c, request.EnterpriseId, map[string]interface{}{
		"unionid": userUnionIds,
	})
	userCode := make(map[string]string, 0)
	for _, model := range many {
		userCode[model.Unionid] = model.UnionCode
	}
	for _, v := range res {
		item := reserveRequest.OrderListResult{
			ReserveOrderModel: v,
		}
		item.ReserveName = v.ReserveService.ReserveName
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		item.StatusStr = getStatusMap(v.Status)
		switch v.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
		case 2:
			item.CustomerStatusStr = "老客"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		}
		if v.IsVip == 1 {
			item.CustomerStatusStr += ",VIP客户"
		}
		if value, ok := userCode[item.UserUnionid]; ok {
			item.UnionCode = value
		}
		item.RelationOrderNum = len(v.ReserveOrderSales)
		result = append(result, item)
	}
	return result, total, err
}

func (s *OrderService) GetListByCoach(c app.Ctx, request reserveRequest.OrderListByCoachReq) ([]reserveRequest.OrderListResult, int64, error) {
	res, total, err := s.orderDao.GetListByCoach(c, request)
	if err != nil || len(res) == 0 {
		return []reserveRequest.OrderListResult{}, total, err
	}
	result := make([]reserveRequest.OrderListResult, 0)
	var userUnionIds []string
	for _, v := range res {
		userUnionIds = append(userUnionIds, v.UserUnionid)
	}
	many, err := externalSearchDao.NewExternalExtraCk().GetList(c, request.EnterpriseId, "user_unionid in ?", userUnionIds)
	externalMap := make(map[string]string, 0)
	for _, model := range many {
		externalMap[model.UserUnionid] = model.UnionCode
	}
	for _, v := range res {
		tmpRow := reserveRequest.OrderListResult{
			ReserveOrderModel: v,
		}
		tmpRow.ReserveName = v.ReserveService.ReserveName
		if value, ok := externalMap[tmpRow.UserUnionid]; ok {
			tmpRow.UnionCode = value
		}
		tmpRow.RelationOrderNum = len(v.ReserveOrderSales)
		result = append(result, tmpRow)
	}
	return result, total, err
}

func (s *OrderService) GetExportList(c app.Ctx, request reserveRequest.OrderListReq) ([]reserveRequest.OrderExportListResult, int64, error) {
	var total int64
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, request.EnterpriseId, "store_name")
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}

	res, total, err := s.orderDao.GetExportList(c, request.EnterpriseId, request.Page, request.Size, request, isCoach)
	if err != nil {
		return nil, total, err
	}

	var result []reserveRequest.OrderExportListResult
	// 查询服务信息
	var listReq reserveRequest.ReserveServiceListReq
	serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, request.EnterpriseId, listReq)
	if errs != nil {
		return result, total, errs
	}
	serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	for _, item := range serviceList {
		serviceListMap[item.ID] = item
	}
	for _, v := range res {
		item := reserveRequest.OrderExportListResult{}
		item.OrderListRes = v
		item.ReserveName = serviceListMap[v.ReserveID].ReserveName
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.CommentTime == 0 {
			item.CommentTimeStr = "--"
		} else {
			item.CommentTimeStr = helper.GetFormatDateByParams(int64(v.CommentTime))
		}

		if v.ReserveOrderTime == 0 {
			item.ReserveOrderTimeStr = "--"
		} else {
			item.ReserveOrderTimeStr = helper.GetFormatDateByParams(int64(v.ReserveOrderTime))
		}
		if v.IsOpen == 1 {
			item.IsOpenStr = "是"
		} else {
			item.IsOpenStr = "否"
		}
		if v.IsCommit == 1 {
			item.IsCommitStr = "是"
		} else {
			item.IsCommitStr = "否"
		}
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		if v.IsSign == 1 {
			item.IsSignInStr = "已签到"
		} else {
			item.IsSignInStr = "未签到"
		}
		if v.IsComment == 1 {
			item.IsCommentStr = "是"
		} else {
			item.IsCommentStr = "否"
		}
		isStoreStatus := "-"
		switch v.IsStore {
		case 1:
			isStoreStatus = "1"
			item.IsStoreStr = "是"
		case 2:
			isStoreStatus = "0"
			item.IsStoreStr = "否"
		}
		if isCoach {
			item.IsStoreStr = isStoreStatus
		}
		if v.IsPay == 1 {
			item.IsPayStr = "是"
		} else {
			item.IsPay = 0
			item.IsPayStr = "否"
		}
		switch v.OrderOrigin {
		case 1:
			item.OrderOriginStr = "到店邀约"
		case 2:
			item.OrderOriginStr = "邀请函"
		}
		customerStatus := "-"
		switch v.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
			customerStatus = "1"
		case 2:
			item.CustomerStatusStr = "老客"
			customerStatus = "0"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		default:
			item.CustomerStatusStr = "--"

		}
		if v.IsVip == 1 {
			item.CustomerStatusStr += ",VIP客"
		}
		if isCoach {
			item.CustomerStatusStr = customerStatus
		}
		item.StatusStr = getStatusMap(v.Status)
		if config.Name == "sw" {
			if v.Status == 0 {
				item.StatusStr = "邀约生效中"
			}
			if v.Status == 4 {
				item.StatusStr = "邀约匹配失败"
			}
		}
		item.ServiceScoreStr = fmt.Sprintf("%d星", item.ServiceScore)
		item.StoreScoreStr = fmt.Sprintf("%d星", item.StoreScore)
		// item.StatusStr = getStatusMap(v.Status)
		var reserveOrderPriceFloat decimal.Decimal
		relationOrderNumStr := "--"
		relationOrderNum := len(v.ReserveOrderSales)
		if relationOrderNum > 0 {
			relationOrderNumStr = fmt.Sprintf("%d单", relationOrderNum)
			for _, sale := range v.ReserveOrderSales {
				reserveOrderPriceFloat = decimal.NewFromInt(int64(sale.OrderPrice)).Div(decimal.NewFromInt(100)).Add(reserveOrderPriceFloat)
			}
		} else {
			if item.ReserveOrderPrice > 0 {
				reserveOrderPriceFloat = decimal.NewFromInt(int64(item.ReserveOrderPrice)).Div(decimal.NewFromInt(100))
			}
		}
		item.RelationOrderNum = relationOrderNum
		item.RelationOrderNumStr = relationOrderNumStr
		priceFloat, _ := reserveOrderPriceFloat.Float64()
		item.OrderPrice = priceFloat
		result = append(result, item)
	}

	return result, total, nil
}

func getStatusMap(status int32) string {
	statusStr := ""
	// 预约单的状态,0进行中,1已完成,2待完成,3已取消,4已过期（coach有这个状态）
	statusArrMap := make(map[int32]reserveRequest.StatusList, 0)
	statusArrMap[microScrmModel.RUNNING] = reserveRequest.StatusList{Status: microScrmModel.RUNNING, StatusStr: "进行中"}
	statusArrMap[microScrmModel.FINISH] = reserveRequest.StatusList{Status: microScrmModel.FINISH, StatusStr: "邀约成功"}
	statusArrMap[microScrmModel.WAITFINISH] = reserveRequest.StatusList{Status: microScrmModel.WAITFINISH, StatusStr: "邀约完成"}
	statusArrMap[microScrmModel.CANCEL] = reserveRequest.StatusList{Status: microScrmModel.CANCEL, StatusStr: "已取消"}
	statusArrMap[microScrmModel.OVERTIME] = reserveRequest.StatusList{Status: microScrmModel.OVERTIME, StatusStr: "已过期"}
	if _, ok := statusArrMap[status]; ok {
		statusStr = statusArrMap[status].StatusStr
	}
	return statusStr
}

func getOrderNumber() string {
	str := time.Now().Format("06")
	str += time.Now().Format("01")
	str += time.Now().Format("02")
	str += time.Now().Format("15")
	str += time.Now().Format("04")
	str += time.Now().Format("05")
	// ms := time.Now().Nanosecond() / 1e6
	ms := helper.Random(10000, 99999)
	orderNumber := fmt.Sprintf("%s%03d", str, ms)
	return orderNumber
}

// CreateOrder 创建预约单
func (s *OrderService) CreateOrder(c app.Ctx, enterpriseId int64, orderInfo reserveRequest.OrderInfo) (data microScrmModel.ReserveOrderModel, err error) {
	// 获取导购详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		Mobile: orderInfo.Mobile,
	})
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, errors.New("查询导购信息异常")
	}
	relation, err := externalUserDao.NewUserRelation().GetOne(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberInfo.ID,
		"unionid":       orderInfo.UserUnionid,
		"deleted_at":    0,
	})
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, err
	}
	if relation.RelationStatus != 0 {
		return microScrmModel.ReserveOrderModel{}, errors.New("抱歉，当前您已被用户删除好友，请添加客户后再进行操作哦~")
	}
	nowTime := time.Now().Unix()
	var InfoModel = microScrmModel.ReserveOrderModel{
		BaseModel: models.BaseModel{
			ID:        app.UID(),
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		ReserveNumber:  getOrderNumber(),
		ReserveID:      orderInfo.ReserveId,
		ReserveTime:    orderInfo.ReserveTime,
		StoreID:        orderInfo.StoreId,
		StoreName:      orderInfo.StoreName,
		StoreAddress:   orderInfo.StoreAddress,
		ReserveType:    0,
		UserName:       orderInfo.UserName,
		UserUnionid:    orderInfo.UserUnionid,
		UserExternalID: orderInfo.UserExternalId,
		MemberID:       memberInfo.ID,
		MemberUserID:   memberInfo.Userid,
		MemberName:     memberInfo.Name,
		ReserveRemark:  orderInfo.ReserveRemark,
		IsComment:      microScrmModel.COMMENT_NO,
		IsConfirm:      microScrmModel.CONFIRM_NO,
		OrderOrigin:    microScrmModel.ORDER_ORIGIN_RESERVE,
		IsPay:          microScrmModel.PAY_FAIL,
		InviteNum:      0,
		ReserveEndTime: orderInfo.ReserveEndTime,
	}

	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err == nil { // 如果报错就跳过
		if config.Name == "sw" { // sw的企业创建邀约时，需要确定顾客类型
			detail, _ := customerSvc.NewCustomer().GetCustomerDetail(c, enterpriseId, map[string]interface{}{
				"unionid": orderInfo.UserUnionid,
			})
			if detail.ID > 0 { // 查到用户信息
				InfoModel.CustomerStatus = int32(detail.CustomerType)
				InfoModel.IsVip = int32(detail.VipType)
			}
		}
	}

	data, err = s.orderDao.CreateOrder(c, enterpriseId, InfoModel)
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, err

	}
	return data, nil
}

func (s *OrderService) CancelOrder(c app.Ctx, enterpriseId int64, orderBase reserveRequest.OrderBase) error {
	return s.orderDao.CancelOrder(c, enterpriseId, orderBase)
}

func (s *OrderService) GetOrderInfo(c app.Ctx, enterpriseId, Id int64) (microScrmModel.ReserveOrderModel, error) {
	return s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{
		ID: Id,
	})
}

func (s *OrderService) GetOrderPreloadDetail(c app.Ctx, enterpriseId, Id int64) (microScrmModel.ReserveOrderModel, error) {
	detail, err := s.orderDao.GetOrderPreLoadDetail(c, enterpriseId, reserveRequest.OrderBase{
		ID: Id,
	})
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, err
	}
	// 查询客户实际昵称
	userInfo, _ := externalUserDao.NewWxUserCorpExternal().FindOne(c, enterpriseId, detail.UserExternalID)
	//优先备注,然后才是昵称
	remark, _ := externalUserDao.NewUserRelation().GetRemarkByUnionIdAndMemberId(c, enterpriseId, userInfo.Unionid, detail.MemberID)
	detail.UserName = helper.RuleThree(remark != "", remark, userInfo.Nick)
	return detail, nil

}

func (s *OrderService) OpenOrder(c app.Ctx, enterpriseId int64, info microScrmModel.ReserveOrderModel) error {
	return s.orderDao.OpenOrder(c, enterpriseId, info.ID)
}

func (s *OrderService) UserCommitOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserOrderBase) error {
	return s.orderDao.UserCommitOrder(c, enterpriseId, request)
}

func (s *OrderService) UserEditOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserEditOrder) error {
	return s.orderDao.UserEditOrder(c, enterpriseId, request)
}

// Update 编辑预约单
func (s *OrderService) Update(c app.Ctx, enterpriseId int64, request reserveRequest.MemberEditOrder) error {
	// 需要更新的数据
	updateData := map[string]interface{}{
		"reserve_time":     request.ReserveTime,
		"reserve_end_time": request.ReserveEndTime,
		"reserve_id":       request.ReserveId,
		"store_address":    request.StoreAddress,
		"reserve_remark":   request.ReserveRemark,
		"is_commit":        0,
	}
	// 查询是否存在
	orderInfo, err := s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{ID: request.ID})
	if err != nil || orderInfo.ID == 0 {
		return errorx.NewErrorMsg(1002, "查询记录异常")
	}

	// 获取员工详情信息
	if request.Mobile != "" {
		memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
			MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
			Mobile:         request.Mobile,
		})
		if err != nil || memberInfo.ID == 0 {
			return errorx.NewErrorMsg(1001, "查询导购信息异常")
		}

		if memberInfo.ID != orderInfo.MemberID { // 邀约单的导购ID与会员ID身份不一致，则认为是邀约单分配调整,则需要验证邀约单用户和导购的好友关系
			one, err := externalUserDao.NewUserRelation().GetOne(c, map[string]interface{}{
				"enterprise_id":   enterpriseId,
				"relation_status": 0,
				"deleted_at":      0,
				"member_id":       memberInfo.ID,
				"unionid":         orderInfo.UserUnionid,
			})
			if err != nil {
				return errors.New("该导购和用户非好友关系")
			}
			if one.ID == 0 {
				return errors.New("该导购和用户非好友关系")
			}

		}
		updateData["member_id"] = memberInfo.ID
		updateData["member_name"] = memberInfo.Name
		updateData["member_user_id"] = memberInfo.Userid
		updateData["reserve_type"] = 1
	}

	// 推送模板消息
	if request.ReserveId > 0 || request.ReserveTime > 0 {
		go func() {
			// 获取服务名称
			reserveServiceRes, err := externalUserDao.NewReserveServiceDao().GetInfoById(c, enterpriseId, request.ReserveId)
			if err != nil {
				s.Log().Error(c).Msgf("获取服务名称,异常：%v", err)
				return
			}
			if reserveServiceRes.ID == 0 {
				s.Log().Error(c).Msgf("获取服务名称-记录不存在")
				return
			}

			// 查询用户的openid
			openidRelationInfo, err := externalUserDao.NewOpenidRelation().GetInfo(c, enterpriseId, map[string]interface{}{
				"unionid": orderInfo.UserUnionid,
			})
			if err != nil {
				s.Log().Error(c).Msgf("查询用户的openid,异常：%v", err)
				return
			}
			if openidRelationInfo.WxOpenid != "" {
				_, _ = microwork.NewSubscribeMessage().Send(c, microwork.SubscribeMessageSendReq{
					EnterpriseId: enterpriseId,
					OpenId:       openidRelationInfo.WxOpenid,
					Type:         2,
					Data: microwork.SubscribeMessageSendDataReq{
						ReserveOrderId: orderInfo.ID,
						ServiceName:    reserveServiceRes.ReserveName,
						StoreName:      orderInfo.StoreName,
						ReserveTime:    request.ReserveTime,
						ReserveEndTime: request.ReserveEndTime,
					},
				})
			}
		}()
	}

	// 更新预约单数据
	err = s.orderDao.Update(c, enterpriseId, map[string]interface{}{"id": request.ID, "member_id": orderInfo.MemberID}, updateData)
	return err
}

func (s *OrderService) OrderAssign(c app.Ctx, id, enterpriseId int64, mobile string) error {
	// 获取员工详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Mobile:         mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}

	// 查询是否存在
	orderInfo, err := s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{ID: id})
	if err != nil || orderInfo.ID == 0 {
		return errors.New("预约单记录异常")
	}

	one, err := externalUserDao.NewUserRelation().GetOne(c, map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"relation_status": 0,
		"deleted_at":      0,
		"member_id":       memberInfo.ID,
		"unionid":         orderInfo.UserUnionid,
	})
	if err != nil {
		return errors.New("该导购和用户给好友关系")
	}
	if one.ID == 0 {
		return errors.New("该导购和用户给好友关系")
	}

	// 更新导购信息
	err = s.orderDao.Update(c, enterpriseId, map[string]interface{}{
		"id":        id,
		"member_id": orderInfo.MemberID,
	}, map[string]interface{}{
		"member_id":      memberInfo.ID,
		"member_name":    memberInfo.Name,
		"member_user_id": memberInfo.Userid,
		"reserve_type":   1,
	})

	return err
}

func (s *OrderService) GetReserveTimeList(c app.Ctx, enterpriseId int64, request reserveRequest.OrderTimeListReq) ([]reserveRequest.OrderTimeInfo, error) {
	startTime, endTime := helper.PeriodOfTimestamp()
	endTime = endTime + (30 * 86400)
	memberId := request.MemberId

	if request.Mobile != "" && request.MemberId == 0 {
		memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
			MicroBaseModel: models.MicroBaseModel{
				EnterpriseID: enterpriseId,
			},
			Mobile: request.Mobile,
		})
		if err != nil {
			return nil, errors.New("查询导购信息异常")
		}
		memberId = memberInfo.ID
	}

	return s.orderDao.GetReserveTimeList(c, enterpriseId, memberId, startTime, endTime)
}

func (s *OrderService) GetReserveOrderList(c app.Ctx, request reserveRequest.ReserveOrderListReq, isCoach bool) (reserveRequest.ReserveSum, []reserveRequest.OrderDetailResBaseStr, int64, error) {
	var lists []reserveRequest.OrderDetailResBaseStr
	nums, reserveOrderList, total, err := s.orderDao.GetReserveOrderList(c, request.EnterpriseId, request.Page, request.Size, request, isCoach)
	if err != nil {
		return reserveRequest.ReserveSum{}, lists, 0, err
	}
	reserveNumbers := make([]string, 0)
	for _, v := range reserveOrderList {
		item := reserveRequest.OrderDetailResBaseStr{}
		item.ReserveOrderModel = v
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.CommentTime == 0 {
			item.CommentTimeStr = "--"
		} else {
			item.CommentTimeStr = helper.GetFormatDateByParams(int64(v.CommentTime))
		}
		if v.ReserveOrderTime == 0 {
			item.ReserveOrderTimeStr = "--"
		} else {
			item.ReserveOrderTimeStr = helper.GetFormatDateByParams(int64(v.ReserveOrderTime))
		}

		if v.IsOpen == 1 {
			item.IsOpenStr = "是"
		} else {
			item.IsOpenStr = "否"
		}
		if v.IsCommit == 1 {
			item.IsCommitStr = "是"
		} else {
			item.IsCommitStr = "否"
		}
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		if v.IsSign == 1 {
			item.IsSignInStr = "已签到"
		} else {
			item.IsSignInStr = "未签到"
		}
		if v.IsComment == 1 {
			item.IsCommentStr = "是"
		} else {
			item.IsCommentStr = "否"
		}
		if v.IsStore == 1 {
			item.IsStoreStr = "是"
		} else {
			item.IsStoreStr = "否"
		}
		if v.IsPay == 1 {
			item.IsPayStr = "是"
		} else {
			item.IsPayStr = "否"
		}
		item.StatusStr = getStatusMap(v.Status)
		reserveNumbers = append(reserveNumbers, v.ReserveNumber)
		lists = append(lists, item)
	}
	salesInfoMap := make(map[string]reserveRequest.ReserveSalesNum, 0)
	salesInfo, err2 := externalUserDao.NewReservesOrderSales().GetReserveSalesCount(c, request.EnterpriseId, reserveNumbers, request.ReserveOrderStartTime, request.ReserveOrderEndTime)
	if err2 != nil {
		return nums, lists, total, err
	}
	for _, value := range salesInfo {
		salesInfoMap[value.ReserveNumber] = value
	}
	for i, item := range lists {
		if value, ok := salesInfoMap[item.ReserveNumber]; ok {
			item.RelationOrderNum = value.Num
			item.ReserveOrderPrice = value.ReserveOrderPrice
		}
		lists[i] = item
	}
	return nums, lists, total, err
}

// GetOrderDetail 预约单详情
func (s *OrderService) GetOrderDetail(c app.Ctx, enterpriseId, Id int64, enterpriseName string) (orderDetailInfo reserveRequest.OrderDetailRes, err error) {
	orderDetail, err := s.orderDao.GetOrderPreLoadDetail(c, enterpriseId, reserveRequest.OrderBase{
		ID: Id,
	})
	if err != nil {
		return orderDetailInfo, err
	}
	orderDetailInfo = reserveRequest.OrderDetailRes{
		ReserveOrderModel: orderDetail,
	}

	// 查询服务信息
	var listReq reserveRequest.ReserveServiceListReq
	// listReq.ReserveIds = append(listReq.ReserveIds, orderDetail.ReserveID)

	serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, enterpriseId, listReq)
	if errs != nil {
		return orderDetailInfo, errs
	}
	serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	for _, item := range serviceList {
		serviceListMap[item.ID] = item
	}

	if _, ok := serviceListMap[orderDetail.ReserveID]; ok {
		orderDetailInfo.ReserveName = serviceListMap[orderDetail.ReserveID].ReserveName
		orderDetailInfo.ReserveTimeStr = helper.GetFormatDateByParams(int64(orderDetail.ReserveTime))
		orderDetailInfo.CreateTimeStr = helper.GetFormatDateByParams(orderDetail.CreatedAt)
		orderDetailInfo.ReserveOrderPrice = orderDetail.ReserveOrderPrice
	}
	orderDetailInfo.IsOpenStr = "否"
	if orderDetail.IsOpen == 1 {
		orderDetailInfo.IsOpenStr = "是"
	}
	orderDetailInfo.IsCommitStr = "否"
	if orderDetail.IsCommit == 1 {
		orderDetailInfo.IsCommitStr = "是"
	}
	orderDetailInfo.IsCommentStr = "否"
	if orderDetail.IsComment == 1 {
		orderDetailInfo.IsCommentStr = "是"
	}
	orderDetailInfo.IsPayStr = "否"
	if orderDetail.IsPay == 1 {
		orderDetailInfo.IsPayStr = "是"
	}
	orderDetailInfo.IsSignInStr = "未签到"
	if orderDetail.IsSign == 1 {
		orderDetailInfo.IsSignInStr = "已签到"
	}
	orderDetailInfo.StatusStr = getStatusMap(orderDetail.Status)
	scoreInfo, err := NewServiceScoreService().ServiceScoreDao.GetScoreInfo(c, enterpriseId, reserveRequest.UserOrderBase{
		OrderBase: reserveRequest.OrderBase{
			ID:            Id,
			ReserveNumber: orderDetail.ReserveNumber,
		},
		UnionId: orderDetail.UserUnionid,
	})
	if err != nil {
		return reserveRequest.OrderDetailRes{}, err
	}
	orderDetailInfo.Remark = scoreInfo.Remark
	orderDetailInfo.ServiceScore = scoreInfo.ServiceScore
	orderDetailInfo.StoreScore = scoreInfo.StoreScore
	// 查询操作日志
	var logTypes []int // 订单操作日志
	logTypes = append(logTypes, 1, 2)
	logModelList, err := NewReserveLog().ListLog(c, enterpriseId, Id, logTypes)
	if err != nil {
		return reserveRequest.OrderDetailRes{}, err
	}

	// var logList []reserveRequest.LogList
	var logListArrs [][]reserveRequest.LogListArr
	for _, v := range logModelList {
		var logInfo reserveRequest.LogInfo
		helper.Unmarshal(v.LogInfo, &logInfo)

		var logListArr []reserveRequest.LogListArr

		logListArr = append(logListArr, reserveRequest.LogListArr{
			Key:   "操作事件",
			Value: logInfo.Event,
		})
		logListArr = append(logListArr, reserveRequest.LogListArr{
			Key:   "操作人",
			Value: logInfo.OpreateName,
		})
		logListArr = append(logListArr, reserveRequest.LogListArr{
			Key:   "操作时间",
			Value: helper.GetFormatDateByParams(v.CreatedAt),
		})

		reserveName := ""

		if v.LogType == 1 {

			var logOrderInfo reserveRequest.LogOrderInfo
			helper.Unmarshal(v.OrderInfo, &logOrderInfo)
			if _, ok := serviceListMap[logOrderInfo.ReserveId]; ok {
				reserveName = serviceListMap[logOrderInfo.ReserveId].ReserveName
			}

			logOrderInfo = reserveRequest.LogOrderInfo{
				ReserveId:      logOrderInfo.ReserveId,
				ReserveTime:    logOrderInfo.ReserveTime,
				StoreId:        logOrderInfo.StoreId,
				StoreName:      logOrderInfo.StoreName,
				StoreAddress:   logOrderInfo.StoreAddress,
				UserName:       logOrderInfo.UserName,
				UserUnionid:    logOrderInfo.UserUnionid,
				UserExternalId: logOrderInfo.UserExternalId,
				ReserveRemark:  logOrderInfo.ReserveRemark,
				UserRemark:     logOrderInfo.UserRemark,
				ReserveTimeStr: helper.GetFormatDateByParams(int64(logOrderInfo.ReserveTime)),
				ReserveName:    reserveName,
			}
			// logList = append(logList, reserveRequest.LogList{
			//	LogInfo:      logInfo,
			//	LogOrderInfo: logOrderInfo,
			//	CreateAt:     v.CreatedAt,
			//	CreateAtStr:  helper.GetFormatDateByParams(v.CreatedAt),
			// })

			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "服务门店",
				Value: logOrderInfo.StoreName,
			})
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "门店地址",
				Value: logOrderInfo.StoreAddress,
			})
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "服务名称",
				Value: logOrderInfo.ReserveName,
			})
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "服务时间",
				Value: logOrderInfo.ReserveTimeStr,
			})
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "客户名称",
				Value: logOrderInfo.UserName,
			})
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "客户备注",
				Value: logOrderInfo.UserRemark,
			})

			logListArrs = append(logListArrs, logListArr)
		}
		if v.LogType == 2 {

			var logScoreInfo reserveRequest.LogScoreInfo
			helper.Unmarshal(v.OrderInfo, &logScoreInfo)

			logScoreInfo = reserveRequest.LogScoreInfo{

				ReserveNumber:  logScoreInfo.ReserveNumber,
				ReserveOrderID: logScoreInfo.ReserveOrderID,
				UserName:       logScoreInfo.UserName,
				UserUnionid:    logScoreInfo.UserUnionid,
				UserExternalID: logScoreInfo.UserExternalID,
				MemberID:       logScoreInfo.MemberID,
				MemberName:     logScoreInfo.MemberName,
				MemberUserID:   logScoreInfo.MemberUserID,
				StoreID:        logScoreInfo.StoreID,
				StoreName:      logScoreInfo.StoreName,
				ServiceScore:   logScoreInfo.ServiceScore,
				StoreScore:     logScoreInfo.StoreScore,
				Remark:         logScoreInfo.Remark,
				FlowScore:      logScoreInfo.FlowScore,
				ReserveTime:    logScoreInfo.ReserveTime,
			}
			// logList = append(logList, reserveRequest.LogList{
			//	LogInfo:      logInfo,
			//	LogScoreInfo: logScoreInfo,
			//	CreateAt:     v.CreatedAt,
			//	CreateAtStr:  helper.GetFormatDateByParams(v.CreatedAt),
			// })
			if enterpriseName != "coach" {
				if logScoreInfo.UserName != "" {
					logListArr = append(logListArr, reserveRequest.LogListArr{
						Key:   "评论的用户",
						Value: logScoreInfo.UserName,
					})
				}
				if logScoreInfo.StoreName != "" {
					logListArr = append(logListArr, reserveRequest.LogListArr{
						Key:   "服务门店名称",
						Value: logScoreInfo.StoreName,
					})
				}
			}

			switch enterpriseName {
			case "coach":
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "活动安排",
					Value: strconv.Itoa(int(logScoreInfo.ServiceScore)),
				})
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "整体服务",
					Value: strconv.Itoa(int(logScoreInfo.StoreScore)),
				})
			case "sw":
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "活动安排",
					Value: strconv.Itoa(int(logScoreInfo.ServiceScore)),
				})
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "整体服务",
					Value: strconv.Itoa(int(logScoreInfo.StoreScore)),
				})
			case "ks":
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "店员服务",
					Value: strconv.Itoa(int(logScoreInfo.ServiceScore)) + "分",
				})
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "店铺环境",
					Value: strconv.Itoa(int(logScoreInfo.StoreScore)) + "分",
				})
				logListArr = append(logListArr, reserveRequest.LogListArr{
					Key:   "服务流程",
					Value: strconv.Itoa(int(logScoreInfo.FlowScore)) + "分",
				})
			}
			logListArr = append(logListArr, reserveRequest.LogListArr{
				Key:   "顾客评价",
				Value: logScoreInfo.Remark,
			})
			logListArrs = append(logListArrs, logListArr)
		}
	}

	// orderDetailInfo.LogList = logList
	orderDetailInfo.LogListArr = logListArrs

	// 查询客户信息
	userInfo, _ := externalUserDao.NewWxUserCorpExternal().FindOne(c, enterpriseId, orderDetailInfo.UserExternalID)
	//优先备注,然后才是昵称
	remark, _ := externalUserDao.NewUserRelation().GetRemarkByUnionIdAndMemberId(c, enterpriseId, userInfo.Unionid, orderDetail.MemberID)
	orderDetailInfo.UserName = helper.RuleThree(remark != "", remark, userInfo.Nick)
	orderDetailInfo.Avatar = userInfo.Avatar
	orderDetailInfo.RelationOrderNum = len(orderDetailInfo.ReserveOrderSales)
	return orderDetailInfo, nil
}

// GetTodayOrderList 获取今天进行中的订单
func (s *OrderService) GetTodayOrderList(c app.Ctx, enterpriseId, startTime, endTime int64) ([]reserveRequest.PushMemberInfo, error) {
	return s.orderDao.GetTodayOrderList(c, enterpriseId, startTime, endTime)
}

// GetPreTodayOrderList 获取昨天进行中(待完成)的订单
func (s *OrderService) GetPreTodayOrderList(c app.Ctx, enterpriseId int64, startTime int64, endTime int64, status int) ([]reserveRequest.PushMemberInfo, error) {
	return s.orderDao.GetPreTodayOrderList(c, enterpriseId, startTime, endTime, status)
}

func (s *OrderService) GetRankList(c app.Ctx, enterpriseId int64, request reserveRequest.OrderRankList) ([]reserveRequest.OrderRankRes, error) {
	result, err := s.orderDao.GetRankBaseList(c, enterpriseId, request)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []reserveRequest.OrderRankRes{}, nil
	}
	index := (request.Page - 1) * request.Size
	rankData := make(map[int64]uint, 0)
	var memberIds []int64
	var storeIds []int
	for _, v := range result {
		if _, ok := rankData[v.MemberID]; !ok {
			rankData[v.MemberID] = uint(v.StoreId)
		}
		memberIds = append(memberIds, v.MemberID)
		storeIds = append(storeIds, int(v.StoreId))
	}
	res, err := microWorkDao.NewStoreDepartment().GetStoreName(c, enterpriseId, storeIds)
	if err != nil {
		return nil, err
	}
	storeNameArr := make(map[uint]string, 0)
	for _, v := range res {
		storeNameArr[v.WechatDepartmentID] = v.StoreName
	}
	priceRes, err := s.orderDao.GetRankPrice(c, enterpriseId, memberIds, request.StartTime, request.EndTime)
	if err != nil {
		return nil, err
	}
	priceArr := make(map[int64]float64, 0)
	for _, v := range priceRes {
		priceArr[v.MemberID] = v.Price / 100
	}
	for k, v := range result {
		result[k].Index = index + k + 1
		result[k].Price = priceArr[v.MemberID]
		result[k].DepartmentName = storeNameArr[rankData[v.MemberID]]
	}
	return result, nil
}

// WriteOff 同步安可诚的订单，更新订单状态
func (s *OrderService) WriteOff(c app.Ctx, data reserveRequest.WriteOffReq) error {
	// 查询订单是否被核销过
	info, err := s.orderDao.GetOrderInfo(c, data.EnterpriseId, microScrmModel.ReserveOrderModel{ReserveOrderSn: data.ReserveOrderSn})
	if err != nil || info.ID > 0 {
		return errorx.NewErrorMsg(1000, "该订单已被核销")
	}

	// 查询预约单详情
	orderInfo, err := s.orderDao.GetOrderInfo(c, data.EnterpriseId, reserveRequest.OrderBase{
		ID:            data.ID,
		ReserveNumber: data.ReserveNumber,
	})
	if orderInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "预约单信息不存在")
	}
	// 核销预约单,更新订单编号和订单金额
	err = s.orderDao.Update(c, data.EnterpriseId, map[string]interface{}{
		"id":             data.ID,
		"reserve_number": data.ReserveNumber,
	}, map[string]interface{}{
		"reserve_order_sn":    data.ReserveOrderSn,
		"reserve_order_price": data.ReserveOrderPrice * 100,
		"status":              microScrmModel.FINISH,
		"reserve_order_time":  data.ReserveOrderTime,
		"is_pay":              1,
	})

	return err
}

// Share 分享订单
func (s *OrderService) Share(c app.Ctx, data reserveRequest.OrderBaseReq) error {
	// 查询预约单详情
	orderInfo, err := s.orderDao.GetOrderInfo(c, data.EnterpriseId, reserveRequest.OrderBase{
		ID:            data.ID,
		ReserveNumber: data.ReserveNumber,
	})

	if orderInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "预约单信息不存在")
	}

	// // 邀约次数+1
	err = s.orderDao.Increment(c, data.EnterpriseId, map[string]interface{}{
		"id":             data.ID,
		"reserve_number": data.ReserveNumber,
	}, "invite_num", 1)

	// 积分队列埋点,分享邀约才算
	excitationSvc.NewPointsProduceService().ProduceQueue(c, excitationReq.PointsProduceSave{
		MemberId:     orderInfo.MemberID,
		Type:         excitationModels.RuleCategoryId1,
		Value:        1,
		UniqueHash:   orderInfo.ReserveNumber + fmt.Sprintf("-%d-%d", orderInfo.MemberID, excitationModels.RuleCategoryId1),
		EnterpriseId: data.EnterpriseId,
	})

	return err
}

// Finish 完成预约单
func (s *OrderService) Finish(c app.Ctx, data reserveRequest.OrderAssignReq) error {
	// 查询预约单详情
	orderInfo, err := s.orderDao.GetOrderInfo(c, data.EnterpriseId, reserveRequest.OrderBase{
		ID:            data.ID,
		ReserveNumber: data.ReserveNumber,
	})

	if orderInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "预约单信息不存在")
	}

	if orderInfo.Status == microScrmModel.FINISH {
		return errorx.NewErrorMsg(1002, "此订单已完成")
	}

	// 查询员工信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: data.EnterpriseId},
		Mobile:         data.Mobile,
	})

	if memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1003, "操作人信息不存在")
	}

	if memberInfo.ID != orderInfo.MemberID {
		return errorx.NewErrorMsg(1004, "无操作权限")
	}

	// 更新预约单状态为已完成
	err = s.orderDao.Update(c, data.EnterpriseId, map[string]interface{}{
		"id":             data.ID,
		"reserve_number": data.ReserveNumber,
	}, map[string]interface{}{
		"status": microScrmModel.FINISH,
	})

	return err
}

func (s *OrderService) UpdatePreTodayOrders(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) error {
	return s.orderDao.UpdatePreTodayOrders(c, enterpriseId, startTime, endTime)
}

func (s *OrderService) UserSignIn(c app.Ctx, enterpriseId int64, request reserveRequest.UserSignIn) error {
	return s.orderDao.UserSignIn(c, enterpriseId, request)
}

// ChangeStatus  导购确认订单是否是到店购买,设置新客老客
func (s *OrderService) ChangeStatus(c app.Ctx, enterpriseId int64, request reserveRequest.ChangeStatus) error {
	return s.orderDao.ChangeStatus(c, enterpriseId, request)
}

// ConfirmStore 导购确认用户到店
func (s *OrderService) ConfirmStore(c ctx.Context, enterpriseId int64, request reserveRequest.OrderAssignReq, isCoach bool) error {
	// 获取员工详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	// 查询是否存在
	orderInfo, err := s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{ID: request.ID})
	if err != nil || orderInfo.ID == 0 {
		return errors.New("预约单记录异常")
	}
	if orderInfo.MemberID != memberInfo.ID {
		return errors.New("抱歉，非本人预约单暂不可操作")
	}

	err = s.orderDao.ConfirmStore(c, enterpriseId, request)
	if !isCoach { // ks&&sw 导购确认顾客到店，自动发送邀请评价的订阅消息
		go func() {
			reserveServiceRes, err := externalUserDao.NewReserveServiceDao().GetInfoById(c, enterpriseId, orderInfo.ReserveID)
			if err != nil {
				s.Log().Error(c).Msgf("获取服务名称,异常：%v", err)
				return
			}
			if reserveServiceRes.ID == 0 {
				s.Log().Error(c).Msgf("获取服务名称-记录不存在")
				return
			}
			// 查询用户的openid
			openidRelationInfo, err := externalUserDao.NewOpenidRelation().GetInfo(c, enterpriseId, map[string]interface{}{
				"unionid": orderInfo.UserUnionid,
			})
			if err != nil {
				s.Log().Error(c).Msgf("查询用户的openid,异常：%v", err)
				return
			}
			if openidRelationInfo.WxOpenid != "" {
				_, _ = microwork.NewSubscribeMessage().Send(c, microwork.SubscribeMessageSendReq{
					EnterpriseId: enterpriseId,
					OpenId:       openidRelationInfo.WxOpenid,
					Type:         3,
					Data: microwork.SubscribeMessageSendDataReq{
						ReserveOrderId: orderInfo.ID,
						ServiceName:    reserveServiceRes.ReserveName,
						StoreName:      orderInfo.StoreName,
						ReserveTime:    orderInfo.ReserveTime,
						ReserveEndTime: orderInfo.ReserveEndTime,
					},
				})
			}
		}()
	}

	return err
}

// UserCommentOrder 用户评论订单
func (s *OrderService) UserCommentOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserOrderBase) error {
	return s.orderDao.UserCommentOrder(c, enterpriseId, request)
}

func (s *OrderService) CommentShare(c ctx.Context, enterpriseId int64, request reserveRequest.OrderAssignReq) error {
	// 获取员工详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}

	// 查询是否存在
	orderInfo, err := s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{ID: request.ID})
	if err != nil || orderInfo.ID == 0 {
		return errors.New("预约单记录异常")
	}
	if orderInfo.MemberID != memberInfo.ID {
		return errors.New("抱歉，非本人预约单暂不可操作")
	}
	go func() {
		reserveServiceRes, err := externalUserDao.NewReserveServiceDao().GetInfoById(c, enterpriseId, orderInfo.ReserveID)
		if err != nil {
			s.Log().Error(c).Msgf("获取服务名称,异常：%v", err)
			return
		}
		if reserveServiceRes.ID == 0 {
			s.Log().Error(c).Msgf("获取服务名称-记录不存在")
			return
		}
		// 查询用户的openid
		openidRelationInfo, err := externalUserDao.NewOpenidRelation().GetInfo(c, enterpriseId, map[string]interface{}{
			"unionid": orderInfo.UserUnionid,
		})
		if err != nil {
			s.Log().Error(c).Msgf("查询用户的openid,异常：%v", err)
			return
		}
		if openidRelationInfo.WxOpenid != "" {
			_, _ = microwork.NewSubscribeMessage().Send(c, microwork.SubscribeMessageSendReq{
				EnterpriseId: enterpriseId,
				OpenId:       openidRelationInfo.WxOpenid,
				Type:         3,
				Data: microwork.SubscribeMessageSendDataReq{
					ReserveOrderId: orderInfo.ID,
					ServiceName:    reserveServiceRes.ReserveName,
					StoreName:      orderInfo.StoreName,
					ReserveTime:    orderInfo.ReserveTime,
					ReserveEndTime: orderInfo.ReserveEndTime,
				},
			})
		}
	}()

	return nil
}

// CreateInviteOrder 创建邀请函订单
func (s *OrderService) CreateInviteOrder(c app.Ctx, enterpriseId, inviteRecordId int64, orderInfo reserveRequest.OrderInfo, memberInfo memberReq.MemberDetailRes) (data microScrmModel.ReserveOrderModel, err error) {

	var InfoModel = microScrmModel.ReserveOrderModel{
		BaseModel:      models.CreateBase(app.UIDContext(c)),
		ReserveNumber:  getOrderNumber(),
		ReserveID:      orderInfo.ReserveId,
		ReserveTime:    orderInfo.ReserveTime,
		StoreID:        memberInfo.DepartmentThirdId,
		StoreName:      memberInfo.DepartmentName,
		StoreAddress:   memberInfo.Address,
		UserName:       orderInfo.UserName,
		UserUnionid:    orderInfo.UserUnionid,
		UserExternalID: orderInfo.UserExternalId,
		MemberID:       memberInfo.MembersId,
		MemberUserID:   memberInfo.MembersUserid,
		MemberName:     memberInfo.MembersName,
		ReserveRemark:  orderInfo.ReserveRemark,
		OrderOrigin:    microScrmModel.ORDER_ORIGIN_INVITE,
		IsComment:      microScrmModel.COMMENT_NO,
		IsConfirm:      microScrmModel.CONFIRM_NO,
		IsPay:          microScrmModel.PAY_FAIL,
	}

	data, err = s.orderDao.CreateOrder(c, enterpriseId, InfoModel)
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, err
	}

	// 更新邀请函记录订单
	cond := map[string]interface{}{"id": inviteRecordId}
	updateData := map[string]interface{}{"reserve_order_id": data.ID}
	err = externalUserDao.NewInviteRecordDao().Update(c, enterpriseId, cond, updateData)
	if err != nil {
		return microScrmModel.ReserveOrderModel{}, err
	}
	return data, nil
}

// UpdateInviteOrder 用户确认邀请函订单
func (s *OrderService) UpdateInviteOrder(c app.Ctx, enterpriseId int64, request inviteRequest.CreateOrderReq, isConfirm int) error {

	var updateOrder = reserveRequest.UserConfirmOrder{
		Id:             request.ReserveOrderId,
		ReserveTime:    request.ReserveTime,
		ReserveEndTime: request.ReserveEndTime,
	}
	return s.orderDao.UserConfirmOrder(c, enterpriseId, updateOrder, isConfirm)

}

// GetSearchUnionid  任务中心端邀约列表的筛选unionid
func (s *OrderService) GetSearchUnionid(c ctx.Context, request reserveRequest.OrderListReq) ([]customerReq.UserUnionIdAndMemberId, error) {
	memberIds := make([]int64, 0)
	if request.IsSM == 1 { // 判断是否为店长
		thirdIds := make([]int, 0)
		for _, v := range request.StoreIds {
			thirdIds = append(thirdIds, int(v))
		}
		lists, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, thirdIds, nil)
		if err != nil {
			return nil, err
		}

		for _, info := range lists {
			memberIds = append(memberIds, info.MembersId)
		}
	} else {
		memberIds = append(memberIds, request.MemberIds...)
	}
	unionids, err := s.orderDao.GetSearchUnionid(c, request, memberIds)
	if err != nil {
		return nil, err
	}
	userIds, err := membersDao.NewMembers().GetUserIdsById(c, request.EnterpriseId, memberIds)
	if err != nil {
		return nil, err
	}
	SearchUnionIds, err := externalUserDao.NewUserRelation().GetUnionIdsBySearchListsByCk(c, request.EnterpriseId, request.Search, unionids, userIds, request.MemberId)
	if err != nil {
		return nil, err
	}
	return SearchUnionIds, err

}

// GetOrderListsByStatus 查询coach昨天的邀约单数据给到安客诚(acxiom)
func (s *OrderService) GetOrderListsByStatus(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) ([]microScrmModel.ReserveOrderModel, error) {
	return s.orderDao.GetOrderListsByStatus(c, enterpriseId, startTime, endTime)
}

// VerifyShare  验证导购是否可以分享邀约单
func (s *OrderService) VerifyShare(c ctx.Context, enterpriseId int64, request reserveRequest.OrderAssignReq) error {
	// 获取员工详情信息
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}

	// 查询是否存在
	orderInfo, err := s.orderDao.GetOrderInfo(c, enterpriseId, reserveRequest.OrderBase{ID: request.ID})
	if err != nil || orderInfo.ID == 0 {
		return errors.New("预约单记录异常")
	}
	if orderInfo.MemberID != memberInfo.ID {
		return errors.New("抱歉，非本人预约单暂不可操作")
	}
	return err
}

// GetSuccessfulOrders 门店邀约数量(Successful_Appointment),当天日期-6的邀约成功记录
func (s *OrderService) GetSuccessfulOrders(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) ([]microScrmModel.ReserveOrderModel, error) {
	return s.orderDao.GetSuccessfulOrders(c, enterpriseId, startTime, endTime)
}

// GetStoreAppointmentLists Monthly的Store成功邀约文件，每个财月结束后上传(STORE_APPOINTMENTQTY)
func (s *OrderService) GetStoreAppointmentLists(c ctx.Context, enterpriseId int64, startTime int64, endTime int64) ([]reserveRequest.AcxiomStoreAppointmentQTYInfo, error) {
	return s.orderDao.GetStoreAppointmentLists(c, enterpriseId, startTime, endTime)
}

// GetStoreReportDataLists Monthly的Store成功邀约文件，每个财月结束后上传(Report_DATA)
func (s *OrderService) GetStoreReportDataLists(c ctx.Context, enterpriseId int64, startTime int64, endTime int64) ([]reserveRequest.AcxiomStoreReportDataInfo, error) {
	return s.orderDao.GetStoreReportDataLists(c, enterpriseId, startTime, endTime)
}

func (s *OrderService) DeleteOrder(c ctx.Context, enterpriseId int64, orderBase reserveRequest.OrderBase) error {
	return s.orderDao.DeleteOrder(c, enterpriseId, orderBase)
}

// CheckReserveTime 验证邀约时间段的上限
func (s *OrderService) CheckReserveTime(c ctx.Context, request reserveRequest.CheckReserveTime, name string) (bool, string, error) {
	memberCount, storeCount, err := s.orderDao.GetReserveTimeCount(c, request)
	if err != nil {
		return false, "", err
	}
	if name == "sw" {
		if memberCount >= 15 || storeCount >= 30 {
			return false, "该时段邀约人数已达个人预约上限，建议重新选择时段", err
		}
	}
	if name == "ks" {
		if memberCount >= 5 || storeCount >= 100 {
			return false, "该时段邀约人数已达个人预约上限，建议重新选择时段", err
		}
	}
	return true, "", err
}

// GetReserveTimeOrderLists 查询sw昨天的邀约单数据给到安客诚(acxiom)
func (s *OrderService) GetReserveTimeOrderLists(context ctx.Context, enterpriseId int64, startTime int64, endTime int64) ([]microScrmModel.ReserveOrderModel, error) {
	return s.orderDao.GetReserveTimeOrderLists(context, enterpriseId, startTime, endTime)
}

// GetCustomerList 获取导购的所有邀约单
func (s *OrderService) GetCustomerList(c ctx.Context, request reserveRequest.CustomerReserveListsReq) ([]reserveRequest.OrderListResult, int64, error) {
	res, total, err := s.orderDao.GetCustomerList(c, request.EnterpriseId, request.Page, request.Size, request)
	if err != nil {
		return nil, total, err
	}
	// 查询服务信息
	var result []reserveRequest.OrderListResult
	for _, v := range res {
		item := reserveRequest.OrderListResult{
			ReserveOrderModel: v,
		}
		item.ReserveName = v.ReserveService.ReserveName
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		item.StatusStr = getStatusMap(v.Status)
		switch v.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
		case 2:
			item.CustomerStatusStr = "老客"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		}
		if v.IsVip == 1 {
			item.CustomerStatusStr += ",VIP客户"
		}
		item.RelationOrderNum = len(v.ReserveOrderSales)
		result = append(result, item)
	}

	return result, total, nil
}

// ShowOrder 导购真正分享了邀约单才会展示
func (s *OrderService) ShowOrder(c ctx.Context, enterpriseId int64, orderBase reserveRequest.OrderBase) error {
	return s.orderDao.Update(c, enterpriseId, map[string]interface{}{
		"id":             orderBase.ID,
		"reserve_number": orderBase.ReserveNumber,
	}, map[string]interface{}{"deleted_at": 0})
}

func (s *OrderService) GetReserveOrderListNew(c app.Ctx, request reserveRequest.ReserveOrderListReq) ([]reserveRequest.OrderDetailResBaseStr, int64, error) {
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, request.EnterpriseId, "store_name")
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}
	var lists []reserveRequest.OrderDetailResBaseStr
	reserveOrderList, total, err := s.orderDao.GetReserveOrderListNew(c, request.EnterpriseId, request.Page, request.Size, request)
	if err != nil {
		return lists, 0, err
	}
	// 查询服务信息
	// var listReq reserveRequest.ReserveServiceListReq
	// serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, request.EnterpriseId, listReq)
	// if errs != nil {
	// 	return lists, total, errs
	// }
	// serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	// for _, item := range serviceList {
	// 	serviceListMap[item.ID] = item
	// }
	for _, v := range reserveOrderList {
		item := reserveRequest.OrderDetailResBaseStr{}
		item.ReserveOrderModel = v
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.CommentTime == 0 {
			item.CommentTimeStr = "--"
		} else {
			item.CommentTimeStr = helper.GetFormatDateByParams(int64(v.CommentTime))
		}
		if v.ReserveOrderTime == 0 {
			item.ReserveOrderTimeStr = "--"
		} else {
			item.ReserveOrderTimeStr = helper.GetFormatDateByParams(int64(v.ReserveOrderTime))
		}
		if v.IsOpen == 1 {
			item.IsOpenStr = "是"
		} else {
			item.IsOpenStr = "否"
		}
		if v.IsCommit == 1 {
			item.IsCommitStr = "是"
		} else {
			item.IsCommitStr = "否"
		}
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		if v.IsSign == 1 {
			item.IsSignInStr = "已签到"
		} else {
			item.IsSignInStr = "未签到"
		}
		if v.IsComment == 1 {
			item.IsCommentStr = "是"
		} else {
			item.IsCommentStr = "否"
		}
		isStoreStatus := "-"
		switch v.IsStore {
		case 1:
			isStoreStatus = "1"
			item.IsStoreStr = "是"
		case 2:
			isStoreStatus = "0"
			item.IsStoreStr = "否"
		}
		if isCoach {
			item.IsStoreStr = isStoreStatus
		}

		if v.IsPay == 1 {
			item.IsPayStr = "是"
		} else {
			item.IsPayStr = "否"
		}
		switch item.OrderOrigin {
		case 1:
			item.OrderOriginStr = "到店邀约"
		case 2:
			item.OrderOriginStr = "邀请函"
		}
		customerStatus := "-"
		switch item.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
			customerStatus = "1"
		case 2:
			item.CustomerStatusStr = "老客"
			customerStatus = "0"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		default:
			item.CustomerStatusStr = "--"

		}
		if item.IsVip == 1 {
			item.CustomerStatusStr += ",VIP客"
		}
		if isCoach {
			item.CustomerStatusStr = customerStatus
		}

		item.StatusStr = getStatusMap(item.Status)
		if config.Name == "sw" {
			if item.Status == 0 {
				item.StatusStr = "邀约生效中"
			}
			if item.Status == 4 {
				item.StatusStr = "邀约匹配失败"
			}
		}
		orderNum := len(v.ReserveOrderSales)
		var reserveOrderPriceFloat decimal.Decimal
		relationOrderNumStr := "--"
		if orderNum > 0 {
			relationOrderNumStr = fmt.Sprintf("%d单", orderNum)
			for _, sale := range v.ReserveOrderSales {
				reserveOrderPriceFloat = decimal.NewFromInt(int64(sale.OrderPrice)).Div(decimal.NewFromInt(100)).Add(reserveOrderPriceFloat)
			}
		} else {
			if item.ReserveOrderPrice > 0 {
				reserveOrderPriceFloat = decimal.NewFromInt(int64(item.ReserveOrderPrice)).Div(decimal.NewFromInt(100))
			}
		}
		priceFloat, _ := reserveOrderPriceFloat.Float64()
		item.OrderPrice = priceFloat
		item.RelationOrderNumStr = relationOrderNumStr
		item.RelationOrderNum = orderNum
		lists = append(lists, item)
	}
	return lists, total, err
}

func (s *OrderService) GetListCK(c app.Ctx, request reserveRequest.OrderListReq) ([]reserveRequest.OrderListResult, int64, error) {
	res, total, err := s.orderDao.GetListCk(c, request.EnterpriseId, request.Page, request.Size, request)
	if err != nil {
		return nil, total, err
	}
	// 查询服务信息
	var listReq reserveRequest.ReserveServiceListReq
	serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, request.EnterpriseId, listReq)
	if errs != nil {
		return nil, total, errs
	}
	serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	for _, item := range serviceList {
		serviceListMap[item.ID] = item
	}
	var result []reserveRequest.OrderListResult
	// orderListMap := make(map[string]reserveRequest.OrderListResult, 0)
	var reserveNumberSlice []string
	for _, v := range res {
		item := reserveRequest.OrderListResult{
			ReserveOrderModel: v,
		}
		if _, ok := serviceListMap[v.ReserveID]; ok {
			item.ReserveName = serviceListMap[v.ReserveID].ReserveName
		}
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
		item.CreateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
		if v.ReserveType == 1 {
			item.ReserveTypeStr = "分配预约单"
		} else {
			item.ReserveTypeStr = "默认预约单"
		}
		item.StatusStr = getStatusMap(v.Status)
		switch v.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
		case 2:
			item.CustomerStatusStr = "老客"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		}
		if v.IsVip == 1 {
			item.CustomerStatusStr += ",VIP客户"
		}
		// orderListMap[v.ReserveNumber] = item
		reserveNumberSlice = append(reserveNumberSlice, v.ReserveNumber)
		result = append(result, item)
	}
	salesInfoMap := make(map[string]reserveRequest.ReserveSalesNum, 0)
	salesInfo, err2 := externalUserDao.NewReservesOrderSales().GetReserveSalesCount(c, request.EnterpriseId, reserveNumberSlice, request.ReserveOrderStartTime, request.ReserveOrderEndTime)
	if err2 != nil {
		return result, total, err2
	}
	for _, value := range salesInfo {
		salesInfoMap[value.ReserveNumber] = value
	}
	for i, item := range result {
		if value, ok := salesInfoMap[item.ReserveNumber]; ok {
			item.RelationOrderNum = value.Num
			item.ReserveOrderPrice = value.ReserveOrderPrice
		}
		result[i] = item
	}
	return result, total, nil
}

func (s *OrderService) ExportOrderRecord(c app.Ctx, request reserveRequest.ReserveOrderListReq) ([]reserveRequest.ExportOrderRecord, int64, error) {
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, request.EnterpriseId, "store_name")
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}
	lists := make([]reserveRequest.ExportOrderRecord, 0)

	saleGoodsInfo, total, err := externalUserDao.NewReservesOrderSales().GetReserveSaleGoods(c, request.EnterpriseId, request)
	if err != nil {
		return lists, total, err
	}
	// 查询服务信息
	// var listReq reserveRequest.ReserveServiceListReq
	// serviceList, _, errs := s.ReserveServiceDao.GetListAll(c, request.EnterpriseId, listReq)
	// if errs != nil {
	// 	return lists, total, errs
	// }
	// serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	// for _, item := range serviceList {
	// 	serviceListMap[item.ID] = item
	// }
	for _, orderInfo := range saleGoodsInfo {

		item := reserveRequest.ExportOrderRecord{
			ReserveNumber:  orderInfo.ReserveNumber,
			ReserveTime:    orderInfo.ReserveTime,
			StoreName:      orderInfo.StoreName,
			Status:         orderInfo.Status,
			UserName:       orderInfo.UserName,
			MemberID:       orderInfo.MemberID,
			MemberName:     orderInfo.MemberName,
			IsOpen:         orderInfo.IsOpen,
			IsCommit:       orderInfo.IsCommit,
			InviteNum:      orderInfo.InviteNum,
			IsStore:        orderInfo.IsStore,
			IsPay:          orderInfo.IsPay,
			OrderOrigin:    orderInfo.OrderOrigin,
			CustomerStatus: orderInfo.CustomerStatus,
		}
		item.ReserveTimeStr = helper.GetFormatDateByParams(int64(orderInfo.ReserveTime))
		if orderInfo.IsOpen == 1 {
			item.IsOpenStr = "是"
		} else {
			item.IsOpenStr = "否"
		}
		if orderInfo.IsCommit == 1 {
			item.IsCommitStr = "是"
		} else {
			item.IsCommitStr = "否"
		}
		isStoreStatus := "-"
		switch orderInfo.IsStore {
		case 1:
			isStoreStatus = "1"
			item.IsStoreStr = "是"
		case 2:
			isStoreStatus = "0"
			item.IsStoreStr = "否"
		}
		if isCoach {
			item.IsStoreStr = isStoreStatus
		}
		if orderInfo.IsPay == 1 {
			item.IsPayStr = "是"
		} else {
			item.IsPay = 0
			item.IsPayStr = "否"
		}
		switch orderInfo.OrderOrigin {
		case 1:
			item.OrderOriginStr = "到店邀约"
		case 2:
			item.OrderOriginStr = "邀请函"
		}
		customerStatus := "-"
		switch orderInfo.CustomerStatus {
		case 1:
			item.CustomerStatusStr = "新客"
			customerStatus = "1"
		case 2:
			item.CustomerStatusStr = "老客"
			customerStatus = "0"
		case 4:
			item.CustomerStatusStr = "沉睡客"
		default:
			item.CustomerStatusStr = "--"
		}
		if isCoach {
			item.CustomerStatusStr = customerStatus
		}
		item.StatusStr = getStatusMap(orderInfo.Status)
		if config.Name == "sw" {
			if orderInfo.Status == 0 {
				item.StatusStr = "邀约生效中"
			}
			if orderInfo.Status == 4 {
				item.StatusStr = "邀约匹配失败"
			}
		}
		var GoodsPriceFloat decimal.Decimal
		GoodsPriceFloat = decimal.NewFromInt(int64(orderInfo.GoodsPrice)).Div(decimal.NewFromInt(100))
		priceFloat, _ := GoodsPriceFloat.Float64()
		item.RelationOrderID = orderInfo.RelationOrderID
		item.GoodsName = orderInfo.GoodsName
		item.Style = orderInfo.Style
		item.Color = orderInfo.Color
		item.OrderPrice = priceFloat
		item.PayTime = helper.GetFormatDateByParams(int64(orderInfo.PayTimeInt))
		item.StoreCode = orderInfo.StoreCode
		lists = append(lists, item)
	}

	return lists, total, err
}

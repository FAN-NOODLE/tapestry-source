package user

import (
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/go/cy-sdk-go/exception"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/pkg/common"
)

// Service 用户服务
type OrderService struct {
	customerRepo customer.Repository
}

func (o *OrderService) SyncOrder(c ctx.Context, request *OrderSyncData, userService *Service, event common.DdcEvent, marketService *marketing.Service) error {
	o.processRequestData(request)
	appid := request.User.AppID
	if appid == "" {
		appid, _ = userService.enterpriseRepo.QueryOneApp(c, request.User.EnterpriseID, customer.CustomerTypeShop)
	}
	source := &SourceUser{
		EnterpriseID: request.User.EnterpriseID,
		SourceUserBase: customer.SourceUserBase{
			AppID:        appid,
			OpenID:       request.User.OpenId,
			UnionID:      request.User.UnionID,
			CustomerType: customer.CustomerTypeShop,
		},
		EnterpriseUserBase: customer.EnterpriseUserBase{
			Nick:           request.User.Nickname,
			Avatar:         request.User.Avatar,
			Birthday:       request.User.Birthday,
			Sex:            int8(request.User.Sex),
			Mobile:         request.User.Mobile,
			DeliveryMobile: request.Info.ShopOrderDeliveryEntity.Mobile,
		},
		EnterpriseUserLocation: customer.EnterpriseUserLocation{
			Country:  request.User.Country,
			City:     request.User.City,
			Province: request.User.Province,
		},
		CreateAt: time.Now().Unix(),
	}

	if _, err := userService.MergeUser(c, source); err != nil {
		return err
	}

	var goodsSpuIds []int64
	for _, goodsData := range request.Goods {
		goodsSpuIds = append(goodsSpuIds, goodsData.SpuID)
	}

	var marketParam = marketing.MarketTrigRequest{
		EnterpriseID:  request.User.EnterpriseID,
		UserUnionID:   request.User.UnionID,
		Condition:     4,
		OrderID:       request.Info.ShopOrderEntity.OrderNo,
		AddressName:   request.Info.ShopOrderDeliveryEntity.Name,
		UserName:      source.Nick,
		Mobile:        source.Mobile,
		Point:         request.Info.ShopOrderEntity.TotalAmount,
		OrderGoodsSpu: goodsSpuIds,
	}
	//如果用户授权了手机号，就给用户发短信，如果没有就给收货地址的手机号发短信

	if marketParam.Mobile == "" {
		marketParam.Mobile = request.Info.ShopOrderDeliveryEntity.Mobile
	}
	switch event {
	case common.DdcEventEnum.OrderCreate:
		marketParam.Condition = 4
		marketParam.Point = int64(request.Info.ShopOrderEntity.FinalAmount)
		log.Info(c).Msgf("下单触发营销策略消息推送数据，marketParam:%+v", marketParam)
		err := marketService.MarketTrigBatch(c, &marketParam)
		if err != nil {
			exception.SendError(err, map[string]interface{}{"下单触发营销策略消息推送数据": request})
		}
	case common.DdcEventEnum.OrderPaySuccess:
		marketParam.Condition = 5
		marketParam.Point = int64(request.Info.ShopOrderEntity.CashFee)
		log.Info(c).Msgf("支付订单触发营销策略消息推送数据，marketParam:%+v", marketParam)
		err := marketService.MarketTrigBatch(c, &marketParam)
		if err != nil {
			exception.SendError(err, map[string]interface{}{"支付订单触发营销策略消息推送数据": request})
		}

		eventDesc := "实付金额：￥" + strconv.FormatFloat(float64(request.Info.ShopOrderEntity.CashFee)/100, 'f', 2, 32)
		err = userService.AddUserEvent(c, request.Info.ShopOrderEntity.EnterpriseID, source.EuID, common.EventEnum.PayOrder.String(), eventDesc, request.Info.ShopOrderEntity.OrderNo, "", int64(request.Info.ShopOrderEntity.PayAt))
		if err != nil {
			log.Info(c).Msgf("添加支付事件失败， err:%v", err)
		}
	case common.DdcEventEnum.OrderRefund,
		common.DdcEventEnum.OrderRefundPartialOverdue,
		common.DdcEventEnum.OrderRefundPartial:
	}
	return nil
}

// processRequestData 处理请求数据
func (o *OrderService) processRequestData(request *OrderSyncData) {
	request.Info.ThirdPaymentNo = request.GetPayInfo().ThirdPaymentNo
	for _, v := range request.Coupons {
		request.Info.CouponFee += v.Value
	}
}

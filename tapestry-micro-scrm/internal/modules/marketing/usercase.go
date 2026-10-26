package marketing

import (
	"encoding/json"
	"math"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"sync"
	"time"

	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/marketing"
	marketing2 "micro-scrm/internal/orm/service/marketing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRepository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/pkg/cdp"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/microactivities"
	"micro-scrm/internal/pkg/sdk/microgoods"
	"micro-scrm/internal/pkg/sdk/shop"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

const (
	StrategyTypePublic             int64 = 1 // 公共营销策略
	StrategyTypePrivate            int64 = 2 // 内置营销策略
	ConditionPay                   int   = 5 // 支付条件状态
	ConditionsInnerConditionMoney        = 2 // 金额条件限制
	ConditionsInnerConditionGoods        = 3 // 商品条件限制
	ConditionsConditions                 = 1 // 不限条件
	ConditionsConditionsVisit            = 2 // 访问触达
	ConditionsConditionsJoinCart         = 3 // 加入购物车触达
	ConditionsConditionsPlaceOrder       = 4 // 提交订单触达
	ConditionsConditionsPayOrder         = 5 // 支付订单触达

)

// Service 应用服务
type Service struct {
	repo            Repository
	customerRepo    customer.Repository
	esRepo          EsInterface
	SmsService      *sms.Service
	DrainageService *drainage.DrainageService
}

type UserInterface interface {
}

type EsInterface interface {
	QueryUserList(c ctx.Context, cond *customer.ListRequest) (int64, []*customer.EnterpriseUser, error)
	QueryUserCountByGroup(c ctx.Context, cond *customer.ListCond) (map[string]*customer.TagUserCount, error)
}

func (s Service) CloseMarketingInnerStatus(c ctx.Context, EnterpriseId int64) error {
	request := &MarketingDataDetailRequest{
		EnterpriseId: EnterpriseId,
	}
	err := s.repo.CloseMarketingInnerStatus(c, request)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) GetSmsCount(c ctx.Context, smsId int64) (int64, error) {
	total, err := s.repo.GeSmsCount(c, smsId)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// SimpleStrategyList 策略列表
func (s Service) SimpleStrategyList(cc ctx.Context, request ListRequest) ([]*MarketingStrategy, []int64, error) {

	list, err := s.repo.SimpleStrategyList(cc, request)

	if err != nil {
		return nil, nil, err
	}
	var smsIds []int64
	if len(list) > 0 {
		channelsMap := make(map[int64][]ChannelInfo)
		conditionMap := make(map[int64]*StrategyConditionItem)
		var ids []int64
		for _, v := range list {
			ids = append(ids, v.ID)
		}
		// 需要查策略对应的渠道
		if request.IsChannels == 1 {
			req := StrategyChannelListRequest{
				EnterpriseID:    request.EnterpriseID,
				StrategyIds:     ids,
				StrategyIdCount: len(ids),
			}
			channels, cErr := s.repo.ListChannelByStrategyIds(cc, req)

			if cErr == nil && len(*channels) > 0 {
				for _, vItem := range *channels {
					channelsMap[vItem.StrategyID] = append(channelsMap[vItem.StrategyID], *vItem)
					if ChannelTypeSMS == vItem.ChannelType {
						smsIds = append(smsIds, vItem.JuniorID)
					}
				}
			}

		}
		// 查询条件
		conditions, cErr := s.repo.ListConditionByStrategyIds(cc, ids)
		if cErr == nil && len(*conditions) > 0 {
			for _, vItem := range *conditions {
				conditionMap[vItem.StrategyID] = vItem
			}
		}
		// 处理list
		for _, item := range list {
			if channel, ok := channelsMap[item.ID]; ok {
				item.Channels = channel
			}

			if condition, yes := conditionMap[item.ID]; yes {
				item.Condition = condition.Condition
				item.InnerCondition = condition.InnerCondition
			}

		}
	}
	return list, smsIds, nil
}

// List 策略列表
func (s Service) List(c ctx.Context, request ListRequest) (int64, *[]MarketingStrategy) {
	offset := (request.Page - 1) * request.Size
	limit := request.Size
	option := &repos.FindPageOption{
		TplPrefix: "t_marketings/select",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MarketingStrategy{}
		},
	}
	if request.StrategyName != "" {
		var build strings.Builder
		build.WriteString("%")
		build.WriteString(request.StrategyName)
		build.WriteString("%")
		request.StrategyName = build.String()
	}
	total, list, err := s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request, option)
	if err != nil {
		return 0, nil
	}
	// 不存在内置策略则创建
	if total == 0 && request.StrategyType == 2 {
		// 获取默认短信模板
		var smsTemplate = sms.SmsTemplate{
			EnterpriseId: request.EnterpriseID,
		}
		SmsTemplateId, _ := s.SmsService.CheckDefaultSmsTemplate(c, smsTemplate, false)

		req := MarketingStrategy{
			StrategyName:   "加粉短信",
			StrategyType:   2,
			StrategyMethod: 1,
			Condition:      5,
			StrategyPolicy: 1,
			Status:         1,
			IsFinished:     1,
			InnerCondition: 1, // 内置营销策略默认 1 不设条件
			EnterpriseID:   request.EnterpriseID,
		}
		req.Channels = []ChannelInfo{}
		channel := ChannelInfo{
			ChannelType: 1,
			JuniorID:    SmsTemplateId,
			Order:       1,
		}
		req.Channels = append(req.Channels, channel)
		_, addErr := s.Add(c, req)
		if addErr != nil {
			return 0, nil
		}
		total, list, err = s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request, option)
		if err != nil {
			return 0, nil
		}
	}
	marketings := list.(*[]MarketingStrategy)
	var marketingList []MarketingStrategy
	for _, value := range *marketings {
		// 如果是立即执行则结束时间等于开始时间
		if value.StrategyMethod == 2 {
			value.EndTime = value.StartTime
		}
		marketingList = append(marketingList, *MarketingInfo(&value))
	}
	return total, &marketingList
}

// makeUpData 组装参数
func makeUpData(c ctx.Context, request *MarketingStrategy) (*MarketingStrategies, *MarketingStrategyCondition) {
	currentTime := time.Now().Unix()
	if request.ID == 0 {
		request.ID = uid.ID(ctx.New())
		request.Status = DisableStrategy
		request.IsFinished = NotFinished
		if request.StrategyType == CustomerStrategy {
			request.Status = EnableStrategy
			request.IsFinished = Finished
		}
		request.StrategyStatus = StrategyStatusInit
		request.CreatedAt = currentTime
	}
	strategyInfo := &MarketingStrategies{
		ID:              request.ID,
		StrategyName:    request.StrategyName,
		EnterpriseID:    request.EnterpriseID,
		Mobile:          request.Mobile,
		StartTime:       request.StartTime,
		EndTime:         request.EndTime,
		CustomerGroupID: request.CustomerGroupID,
		StrategyType:    request.StrategyType,
		StrategyMethod:  request.StrategyMethod,
		StrategyPolicy:  request.StrategyPolicy,
		StrategyStatus:  request.StrategyStatus,
		Status:          request.Status,
		IsFinished:      request.IsFinished,
		IsNotice:        request.IsNotice,
		AdminID:         request.AdminID,
		HistoryStatus:   request.HistoryStatus,
		CreatedAt:       request.CreatedAt,
		UpdatedAt:       currentTime,
		TriggerCount:    request.TriggerCount,
	}
	// 策略触发条件
	/** 根据时间类型，计算间隔秒 **/
	switch request.TimeType {
	case TimeTypeSecond:
		request.Time = request.TimeInterval
	case TimeTypeMin:
		request.Time = request.TimeInterval * 60
	case TimeTypeHour:
		request.Time = request.TimeInterval * 3600
	case TimeTypeDay:
		request.Time = request.TimeInterval * 86400
	}
	strategyCondition := &MarketingStrategyCondition{
		ID:             request.ID,
		Condition:      request.Condition,
		AssignType:     request.AssignType,
		RelationType:   request.RelationType,
		BelongType:     request.BelongType,
		InnerCondition: request.InnerCondition,
		StartPoint:     request.StartPoint,
		EndPoint:       request.EndPoint,
		TimeInterval:   request.TimeInterval,
		TimeType:       request.TimeType,
		Time:           request.Time,
		CreatedAt:      currentTime,
		UpdatedAt:      currentTime,
	}

	return strategyInfo, strategyCondition
}

//Add 策略添加
func (s Service) Add(c ctx.Context, request MarketingStrategy) (int64, error) {
	strategyInfo, strategyCondition := makeUpData(c, &request)
	request.ID = strategyInfo.ID
	currentTime := time.Now().Unix()
	if request.Condition == StrategyConditionTypeUnlimited || request.Condition == StrategyConditionTypeAddCart || request.Condition == StrategyConditionTypeCreateOrder || request.Condition == StrategyConditionTypePay {
		request.AssignType = AssignTypeAll
	}

	group := threading.NewRoutineGroup()

	// 策略主体
	err := s.repo.MarketingAdd(c, strategyInfo)
	if err != nil {
		return 0, err
	}
	// 触发人数
	recordCountErr := orm.Service.Marketing.RecordTriggerCount(marketing.MarketingStrategyTriggerCountModel{StrategyID: strategyInfo.ID, Count: int(request.TriggerCount)})
	if recordCountErr != nil {
		return 0, recordCountErr
	}
	// 策略条件
	group.RunSafe(func() {
		_ = s.repo.MarketingStrategyConditionAdd(c, strategyCondition)
	})

	// 策略触发条件关联
	group.RunSafe(func() {
		if len(request.Relations) > 0 {
			relationInfos := []RelationInfo{}
			for _, val := range request.Relations {
				val.ID = uid.ID(ctx.New())
				val.StrategyID = request.ID
				val.CreatedAt = currentTime
				val.UpdatedAt = currentTime
				relationInfos = append(relationInfos, val)
			}
			_ = s.repo.MarketingStrategyTriggerRelationAdd(c, relationInfos)
		}
	})

	// 策略触发渠道
	group.RunSafe(func() {
		if len(request.Channels) > 0 {
			channelInfos := []ChannelInfo{}
			for _, val := range request.Channels {
				contentJson, _ := json.Marshal(val.Content)
				val.TextContent = string(contentJson)
				val.ID = uid.ID(ctx.New())
				val.StrategyID = request.ID
				val.CreatedAt = currentTime
				val.UpdatedAt = currentTime
				channelInfos = append(channelInfos, val)
			}
			err = s.repo.MarketingStrategyChannelAdd(c, channelInfos)
		}
	})
	if err != nil {
		return 0, err
	}
	group.Wait()
	if err != nil {
		return 0, err
	}
	return request.ID, nil
}

//Update策略修改
func (s Service) Update(c ctx.Context, request MarketingStrategy) (int64, error) {
	currentTime := time.Now().Unix()
	if request.ID == 0 {
		return 0, errs.ParamsErr
	}
	resp, markerErr := s.repo.MarketingDetail(c, request.ID)

	if markerErr != nil {
		return 0, markerErr
	}
	// 如果投放中或者已结束不能编辑
	if resp.StrategyStatus == StrategyStatusStart {
		return 0, errs.MarketingLaunchZero
	}
	if resp.StrategyStatus == StrategyStatusEnd {
		return 0, errs.MarketingEndZero
	}

	// 删除原记录
	req := DeleteRequest{
		ID:           request.ID,
		EnterpriseID: request.EnterpriseID,
	}

	// 组装策略信息
	strategyInfo, strategyCondition := makeUpData(c, &request)

	relationInfos := []RelationInfo{}
	if len(request.Relations) > 0 {
		for _, val := range request.Relations {
			val.ID = uid.ID(ctx.New())
			val.StrategyID = request.ID
			val.CreatedAt = currentTime
			val.UpdatedAt = currentTime
			relationInfos = append(relationInfos, val)
		}
	}

	channelInfos := []ChannelInfo{}
	msCoupons := []MarketingStrategyCoupons{}
	if len(request.Channels) > 0 {
		for _, val := range request.Channels {
			contentJson, jsonErr := json.Marshal(val.Content)
			if jsonErr != nil {
				return 0, errs.JsonParseErr
			}
			val.TextContent = string(contentJson)

			val.ID = uid.ID(ctx.New())
			val.StrategyID = request.ID
			val.CreatedAt = currentTime
			val.UpdatedAt = currentTime
			channelInfos = append(channelInfos, val)

		}
	}

	err := s.repo.TxUpdateMarketingStrategy(c, req, strategyInfo, strategyCondition, relationInfos, channelInfos, msCoupons)
	if err != nil {
		return 0, err
	} else {
		return request.ID, nil
	}
}

//UpdateStatus 更新状态
func (s Service) UpdateStatus(c ctx.Context, request *UpdateStatusRequest) (int64, error) {
	resp, err := s.repo.MarketingDetail(c, request.ID)
	if err != nil {
		return 0, err
	}
	// 内置营销策略或者客户策略点击开始，判断短信模板是否可用
	if resp.StrategyType == BuiltInPolicy {
		channels, mErr := s.repo.MarketingDetailChannels(c, MarketingDetailRequest{
			ID:           resp.ID,
			EnterpriseID: request.EnterpriseID,
		})
		if mErr != nil {
			return 0, mErr
		}
		if len(*channels) > 0 {
			// 查询短信模板
			status, fErr := s.SmsService.FindSmsTemplateStatusByID(c, (*channels)[0].JuniorID)
			if fErr != nil {
				return 0, fErr
			}
			if status != 5 {
				return 0, errs.TemplateStatusNomarl
			}
		}
	}
	if resp.StrategyType == CustomerStrategy && request.StrategyStatus == StrategyStatusStart {
		channels, mErr := s.repo.MarketingDetailChannels(c, MarketingDetailRequest{
			ID:           resp.ID,
			EnterpriseID: request.EnterpriseID,
		})
		if mErr != nil {
			return 0, mErr
		}
		if len(*channels) > 0 && (*channels)[0].ChannelType == ChannelTypeSMS {
			// 查询短信模板
			status, fErr := s.SmsService.FindSmsTemplateStatusByID(c, (*channels)[0].JuniorID)
			if fErr != nil {
				return 0, fErr
			}
			if status != sms.TemplateStatusHasValue {
				return 0, errs.TemplateStatusNomarl
			}
		}
	}
	currentTime := time.Now().Unix()
	if resp.StrategyMethod == MethodImmediately && request.StrategyStatus == StrategyStatusStart {
		request.StartTime = currentTime
		request.EndTime = currentTime + (3 * 60)
	}
	request.UpdatedAt = currentTime
	//  如果修改状态为投放中且策略之前为异常暂停，需恢复字段
	if request.StrategyStatus == StrategyStatusStart && resp.AbnormalPause == AbnormalPauseYes {
		request.AbnormalPause = AbnormalPauseNo
	}
	err = s.repo.MarketingUpdateStatus(c, request)
	if err != nil {
		return 0, err
	}
	return request.ID, nil

}

// Detail 策略详情
func (s Service) Detail(c ctx.Context, request MarketingDetailRequest) (*MarketingStrategy, error) {
	resp, err := s.repo.MarketingDetail(c, request.ID)
	if err != nil {
		return nil, err
	}

	cres, cerr := s.repo.MarketingDetailChannels(c, request)
	if cerr != nil {
		return nil, cerr
	}
	resp.Channels = *cres
	// 渠道数据整理
	if len(*cres) > 0 {
		resp, _ = s.channelsHandler(c, resp)
	}

	rres, rerr := s.repo.MarketingDetailRelations(c, request)
	if rerr != nil {
		return nil, rerr
	}
	resp.Relations = *rres
	// 关联数据整理
	if resp.Condition == 2 || resp.InnerCondition == 3 && len(*rres) > 0 {
		resp, err = s.relationsHandler(c, resp)
		if err != nil {
			return nil, err
		}
	}
	// 转换类型
	for key, val := range resp.Channels {
		contentJson := make(map[string]interface{})
		if val.TextContent != "" {
			jsonError := json.Unmarshal([]byte(val.TextContent), &contentJson)
			if jsonError != nil {
				return nil, errs.JsonParseErr
			}
		}
		resp.Channels[key].Content = contentJson
		resp.Channels[key].TextContent = ""
	}

	// 新增结束时间
	if resp.StartTime > 0 {
		currentTime := time.Now().Unix()
		resp.StrategyEndTime = currentTime
		if resp.EndTime < currentTime {
			resp.StrategyEndTime = resp.EndTime
		}
	}

	return resp, nil
}

//channelsHandler 渠道信息处理
func (s Service) channelsHandler(c ctx.Context, resp *MarketingStrategy) (*MarketingStrategy, error) {
	for k, val := range resp.Channels {
		// 短信信息
		if val.ChannelType == 1 {
			var smsTemplate = sms.SmsTemplate{
				Id:           val.JuniorID,
				EnterpriseId: resp.EnterpriseID,
			}
			smsTemp, _ := s.SmsService.Repo.SmsTemplateFindOne(c, smsTemplate)
			if smsTemp != nil {
				templateName := strings.Split(smsTemp.TemplateName, "-")
				resp.Channels[k].ChannelDetail.Name = templateName[0]
				resp.Channels[k].ChannelDetail.MessageStr = smsTemp.MessageStr
			}
		}
	}
	return resp, nil
}

//relationsHandler 关联信息处理
func (s Service) relationsHandler(c ctx.Context, resp *MarketingStrategy) (*MarketingStrategy, error) {
	var RelationIDList []int64
	for _, v := range resp.Relations {
		relationID := v.RelationID
		RelationIDList = append(RelationIDList, relationID)
	}
	if len(RelationIDList) != 0 {
		if resp.InnerCondition == 3 || resp.RelationType == 2 {
			//商品详情
			var request = microgoods.SpuListRequest{
				EnterpriseID: resp.EnterpriseID,
				SpuIds:       RelationIDList,
				Page:         1,
				Size:         len(RelationIDList),
			}
			relationResp, rerr := getGoods(ctx.Wrap(c), &request)
			if rerr != nil {
				return nil, rerr
			}
			for key, val := range resp.Relations {
				if relationResp[val.RelationID] != nil {
					resp.Relations[key].RelationDetail = *relationResp[val.RelationID]
				}
			}
		}
		switch resp.RelationType {
		case 3:
			//直播
			var request = shop.LiveListRequest{
				EnterpriseID: resp.EnterpriseID,
				Ids:          RelationIDList,
			}
			relationResp, rerr := getLives(ctx.Wrap(c), &request)
			if rerr != nil {
				return nil, rerr
			}
			for key, val := range resp.Relations {
				if relationResp[val.RelationID] != nil {
					resp.Relations[key].RelationDetail = *relationResp[val.RelationID]
				}
			}
			break
		case 4:
			//活动页
			var request = shop.ActivityPageListRequest{
				EnterpriseID: resp.EnterpriseID,
				Ids:          RelationIDList,
			}
			relationResp, rerr := getActivityPages(ctx.Wrap(c), &request)
			if rerr != nil {
				return nil, rerr
			}
			for key, val := range resp.Relations {
				if relationResp[val.RelationID] != nil {
					resp.Relations[key].RelationDetail = *relationResp[val.RelationID]
				}
			}
			break
		case 5:
			//活动
			var request = microactivities.ActivitiesByIdsRequest{
				EnterpriseID: resp.EnterpriseID,
				RuleIds:      RelationIDList,
				Type:         resp.BelongType,
			}
			relationResp, rerr := getActivities(ctx.Wrap(c), &request)
			if rerr != nil {
				return nil, rerr
			}
			for key, val := range resp.Relations {
				if relationResp[val.RelationID] != nil {
					resp.Relations[key].RelationDetail = *relationResp[val.RelationID]
				}
			}
			break
		}
	}
	return resp, nil
}

//getGoods 获取商品
func getGoods(c ctx.Context, request *microgoods.SpuListRequest) (map[int64]*RelationDetail, error) {
	response, err := microgoods.SpuList(c, request)
	relationMapping := make(map[int64]*RelationDetail)
	if err != nil {
		return nil, err
	}
	for _, v := range response.Data.List {
		relationMapping[v.ID] = &RelationDetail{
			ImgUrl: v.SpuSource.URL,
			Name:   v.SpuName,
		}
	}
	return relationMapping, nil
}

//getLives 获取直播
func getLives(c ctx.Context, request *shop.LiveListRequest) (map[int64]*RelationDetail, error) {
	response, err := shop.LiveList(c, request)
	relationMapping := make(map[int64]*RelationDetail)
	if err != nil {
		return nil, err
	}
	for _, v := range response.Data.List {
		relationMapping[v.ID] = &RelationDetail{
			Name: v.LiveTitle,
		}
	}
	return relationMapping, nil
}

//getActivityPages 获取活动页
func getActivityPages(c ctx.Context, request *shop.ActivityPageListRequest) (map[int64]*RelationDetail, error) {
	response, err := shop.ActivityPageList(c, request)
	relationMapping := make(map[int64]*RelationDetail)
	if err != nil {
		return nil, err
	}
	for _, v := range response.Data.List {
		relationMapping[v.ID] = &RelationDetail{
			Name: v.TemplateName,
		}
	}
	return relationMapping, nil
}

//getActivities 获取活动
func getActivities(c ctx.Context, request *microactivities.ActivitiesByIdsRequest) (map[int64]*RelationDetail, error) {
	response, err := microactivities.Activities(c, request)
	relationMapping := make(map[int64]*RelationDetail)
	if err != nil {
		return nil, err
	}
	for _, v := range response.Data {
		relationMapping[v.ActivityRuleID] = &RelationDetail{
			Name: v.Name,
		}
	}
	return relationMapping, nil
}

//Delete 策略删除
func (s Service) Delete(c ctx.Context, request DeleteRequest) error {
	// 查询策略状态
	resp, err := s.repo.MarketingDetail(c, request.ID)
	if err != nil {
		return err
	}
	if resp.StrategyStatus != StrategyStatusInit && resp.StrategyStatus != StrategyStatusEnd {
		return errs.MarketingCannotDelete
	}
	request.DeletedAt = time.Now().Unix()
	err = s.repo.MarketingDelete(c, request)
	if err != nil {
		return err
	}
	return nil
}

// GetMarketingStrategyDetails 营销策略营销数据
func (s Service) GetMarketingStrategyDetails(c ctx.Context, request MarketingDataDetailRequest) (*cdp.StrategyDataDetails, error) {
	marketingStrategyResult, err := orm.Service.Marketing.GetStrategyResult(marketing2.MarketingResultQueryParam{
		EnterpriseID: request.EnterpriseId,
		StrategyID:   request.ID,
		StartTime:    int64(request.BeginTime),
		EndTime:      int64(request.EndTime),
	})

	if err != nil {
		return nil, err
	}

	// 用了参数结构
	var sdkResponse cdp.StrategyDataDetails
	sdkResponse.Total = &cdp.StrategyItem{
		SendMemberNumber:  0,
		ReachMemberNumber: 0,
	}
	for _, result := range marketingStrategyResult {
		sdkResponse.Data = append(sdkResponse.Data, &cdp.StrategyItem{
			SendMemberNumber:  result.Total,
			ReachMemberNumber: result.NumSuccess,
			Types:             int(result.ChannelType),
		})

		sdkResponse.Total.SendMemberNumber = sdkResponse.Total.SendMemberNumber + result.Total
		sdkResponse.Total.ReachMemberNumber = sdkResponse.Total.ReachMemberNumber + result.NumSuccess
	}

	return &sdkResponse, nil
}

func (s Service) MarketTrigCondition(params *MarketTrigRequest, TrigRecord MarketTrigDetailResponse) bool {
	// 访问首页
	if params.VisitTemplate == "index" && TrigRecord.ConditionsRelationType == 1 {
		return true
	}
	// 访问商品页面
	if params.VisitTemplate == "good-detail" && TrigRecord.ConditionsRelationType == 2 {
		//	1 全部商品 2 指定商品
		if TrigRecord.ConditionsAssignType == 2 {
			if strings.Contains(TrigRecord.RelationIDs, params.VisitDataId) {
				return true
			} else {
				return false
			}
		}
		return true
	}
	// 访问直播页面
	if params.VisitTemplate == "live-detail" && TrigRecord.ConditionsRelationType == 3 {
		//	1 全部商品 2 指定商品
		if TrigRecord.ConditionsAssignType == 2 {
			if strings.Contains(TrigRecord.RelationIDs, params.VisitDataId) {
				return true
			} else {
				return false
			}
		}
		return true
	}
	// 访问活动模板页面
	if params.VisitTemplate == "activity-template" && TrigRecord.ConditionsRelationType == 4 {
		//	1 全部商品 2 指定商品
		if TrigRecord.ConditionsAssignType == 2 {
			if strings.Contains(TrigRecord.RelationIDs, params.VisitDataId) {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

func (s Service) MarketTrigCoditionOne(params *MarketTrigRequest, resStrategy interface{}) bool {

	resByre, _ := json.Marshal(resStrategy)
	var TrigRecord MarketTrigDetailResponse
	_ = json.Unmarshal(resByre, &TrigRecord)

	switch TrigRecord.ConditionsConditions {
	// 不限条件
	case 1:
		if params.Condition != ConditionsConditions {
			return false
		} else {
			return true
		}
	// 访问触达（关联类型；1=首页，2=商品，3=直播，4=活动页，5=活动）
	case 2:
		if params.Condition != ConditionsConditionsVisit {
			return false
		}
		return s.MarketTrigCondition(params, TrigRecord)
	// 加购触发
	case 3:
		if params.Condition != ConditionsConditionsJoinCart {
			return false
		} else {
			return true
		}
	// 下单
	case 4:
		if params.Condition != ConditionsConditionsPlaceOrder {
			return false
		} else {
			return true
		}
	// 支付
	case 5:
		if params.Condition != ConditionsConditionsPayOrder {
			return false
		} else {
			return true
		}
	}
	return false
}

func (s Service) MarketTrigIsPass(c ctx.Context, user *customer.EnterpriseUser, params *MarketTrigRequest, ResStrategyList *MarketTrigResponse) error {
	newTime := time.Now().Unix()
	for _, MarketTrigRecode := range *ResStrategyList {
		// 判断是否在策略执行时间内
		if newTime > MarketTrigRecode.StartTime && newTime < MarketTrigRecode.EndTime {
			var isPass bool
			isPass = false
			// 客户策略
			if MarketTrigRecode.StrategyType == StrategyTypePublic {
				// 策略触发方式
				isPass = s.MarketTrigCoditionOne(params, *MarketTrigRecode)
			}
			// 内置营销策略
			if MarketTrigRecode.StrategyType == StrategyTypePrivate {
				// 是否是支付触发
				if params.Condition != ConditionPay {
					continue
				}
				// 判断策略限制条件是否开启
				if MarketTrigRecode.ConditionsInnerCondition == ConditionsInnerConditionMoney {
					// 购买金额限制条件
					if params.Point < MarketTrigRecode.ConditionsStartPoint || params.Point > MarketTrigRecode.ConditionsEndPoint {
						continue
					}
				}
				// 购买指定商品限制条件
				if MarketTrigRecode.ConditionsInnerCondition == ConditionsInnerConditionGoods {
					var isAdopt = false
					for _, goods := range params.OrderGoodsSpu {
						SpuId := strconv.FormatInt(goods, 10)
						if strings.Contains(MarketTrigRecode.RelationIDs, SpuId) == true {
							isAdopt = true
							break
						}
					}
					if isAdopt != true {
						continue
					}
				}

				externalUser := user.GetOnlyExternal()

				if externalUser != nil {
					if b, _ := s.customerRepo.HasNormalRelation(c, externalUser.ID, params.EnterpriseID, 0); !b {
						continue
					}
				}

				isPass = true
			}
			if isPass == false {
				continue
			}

			if MarketTrigRecode.ChannelType == ChannelTypeSMS && params.Mobile == "" {
				continue
			}
			//保存数据
			rec := &MarketStrategyRecord{
				ID:              uid.ID(c),
				EnterpriseID:    params.EnterpriseID,
				UserUnionID:     params.UserUnionID,
				StrategyID:      MarketTrigRecode.StrategyID,
				ChannelID:       MarketTrigRecode.ChannelID,
				RelationID:      MarketTrigRecode.RelationID,
				ChannelType:     MarketTrigRecode.ChannelType,
				Mobile:          params.Mobile,
				AddressUserName: params.AddressName,
				//UserName:        params.UserName,
				UserName:  "test",
				StartTime: MarketTrigRecode.StartTime,
				OrderID:   params.OrderID,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
				Status:    1,
			}
			// 策略发送时间等于当前时间+策略需要延迟的时间
			rec.StartTime = newTime + MarketTrigRecode.ConditionsTime
			saveErr := s.AddMarketStrategyRecord(c, rec)
			if saveErr != nil {
				return saveErr
			}

			if err := s.customerRepo.UpdateEnterpriseUserAttr(c, params.EnterpriseID, user.ID, &customer.EnterpriseUserAttr{MarketCount: 1}); err != nil {
				log.Err(c, err).Msgf("增加用户影响次数失败: %d", user.ID)
			}
		}
	}
	return nil
}

//策略触发
func (s Service) MarketTrigBatch(c ctx.Context, params *MarketTrigRequest) error {
	user, err := s.customerRepo.QueryEnterpriseUser(c, params.EnterpriseID, &customer.SourceUserBase{UnionID: params.UserUnionID})
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	params.StrategyStatus = StrategyStatusStart
	params.GroupIds = user.Groups
	params.Count = len(user.Groups)
	ResStrategyList, err := s.repo.MarketTrigList(c, params)
	if errorIgnore.Ignore(err) != nil {
		return err
	}

	//增加判断逻辑，如果用户没有授权手机号，以用户最近一笔订单的手机号去发短信
	if params.Mobile == "" {
		if user.Mobile != "" {
			params.Mobile = user.Mobile
		} else {
			params.Mobile = user.DeliveryMobile
		}
	}
	if len(*ResStrategyList) > 0 && params.UserUnionID != "" {
		err := s.MarketTrigIsPass(c, user, params, ResStrategyList)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) AddMarketStrategyRecord(ctx ctx.Context, params *MarketStrategyRecord) error {
	err := s.repo.AddMarketStrategyRecord(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

// MarketingCron 定时任务，执行可执行记录
func (s Service) MarketingCron(c ctx.Context, enterpriseID int64, domain string, currentTime int64, smsService *sms.Service) {
	// ==================查询当前企业正在进行中的策略=====================================
	listStrategy, smsIds, err := s.SimpleStrategyList(c, ListRequest{
		EnterpriseID:        enterpriseID,
		IsChannels:          1,
		StrategyStatusMulti: []int8{StrategyStatusStart, StrategyStatusInit},
	})
	if err != nil {
		log.Error(c).Msgf("策略列表查询,错误:", err)
		return
	}
	// ==================查询营销策略关联的短信模板是否都是已赋值状态===========================
	var smsTemplateList map[int64]*sms.SmsTemplate
	if len(smsIds) > 0 {
		smsTemplateList, err = smsService.TemplateListByTemplateIds(c, enterpriseID, smsIds, sms.TemplateStatusHasValue)
		if err != nil {
			log.Error(c).Msgf("根据模板ID 查询模板,错误:", err)
			return
		}
	}

	log.Info(c).Msgf("企业【%d】预计有【%d】条策略需要执行", enterpriseID, len(listStrategy))
	group := sync.WaitGroup{}

	for _, vStrategy := range listStrategy {
		if vStrategy.ID == 0 || (len(vStrategy.Channels) == 1 && vStrategy.Channels[0].ChannelType == ChannelTypeCoupon) {
			continue
		}
		// 判断客户策略
		if vStrategy.StrategyType == CustomerStrategy {
			if vStrategy.StrategyMethod == StrategyMethodALongTime { // 结束时间大于当前时间修改策略结束状态
				if vStrategy.EndTime < currentTime && vStrategy.StrategyStatus != StrategyStatusEnd {
					if err = s.repo.MarketingUpdateStatus(c, &UpdateStatusRequest{
						ID:             vStrategy.ID,
						StrategyStatus: StrategyStatusEnd,
						UpdatedAt:      currentTime,
					}); err != nil {
						log.Error(c).Msgf("更新客户策略,错误:", err)
					}
					continue
				}
			}
			if vStrategy.StrategyStatus != StrategyStatusStart {
				log.Info(c).Msgf("企业【%d】客户策略【%d】不是投放中不需执行", enterpriseID, vStrategy.ID)
				continue
			}
		} else if vStrategy.StrategyType == BuiltInPolicy {
			if vStrategy.Status != EnableStrategy {
				log.Info(c).Msgf("企业【%d】内置策略【%d】不是启用状态不需执行", enterpriseID, vStrategy.ID)
				continue
			}
		}

		if currentTime < vStrategy.StartTime {
			log.Info(c).Msgf("企业【%d】策略【%d】还未开始或者已解决", enterpriseID, vStrategy.ID)
			continue
		}

		if vStrategy.StrategyMethod == StrategyMethodALongTime { // 长期策略
			// 非不限策略不需要执行，非内置
			if StrategyConditionTypeUnlimited != vStrategy.Condition && vStrategy.StrategyType == CustomerStrategy {
				log.Info(c).Msgf("企业【%d】非内置策略【%d】非长期不限，无需在此执行", vStrategy.EnterpriseID, vStrategy.ID)
				continue
			}

			// 非不限策略不需要执行,内置
			if StrategyConditionTypeUnlimited != vStrategy.InnerCondition && vStrategy.StrategyType == BuiltInPolicy {
				log.Info(c).Msgf("企业【%d】内置策略【%d】非长期不限，无需在此执行", vStrategy.EnterpriseID, vStrategy.ID)
				continue
			}
			// 不在策略执行执行时间范围内不需要执行
			if vStrategy.EndTime < currentTime {
				log.Info(c).Msgf("企业【%d】策略【%d】未开始或者已结束", vStrategy.EnterpriseID, vStrategy.ID)
				continue
			}
		}

		log.Info(c).Msgf("开始执行企业【%d】策略【%d】", enterpriseID, vStrategy.ID)

		group.Add(1)
		go func(currentStrategy *MarketingStrategy, currentDomain string) {
			defer group.Done()
			// 先初始化策略人群
			if err = s.Perform(c, *currentStrategy, smsTemplateList, currentTime); err != nil {
				log.Error(c).Msgf("先初始化策略人群,错误", err)
				return
			}
			// 查询此策略可触达的人数
			listTargetPeople, err := s.repo.MarketingCron(c, &RecordCronRequest{
				CurrentTime: time.Now().Unix(),
				StrategyID:  currentStrategy.ID,
			})
			if err != nil {
				log.Error(c).Msgf("查询此策略可触达的人数,错误", err)
				return
			}

			var massHairEntity []MarketTrigExecRequest
			for _, vPeople := range listTargetPeople {
				// 判断策略处理的类型
				switch vPeople.ChannelType {
				// 短信触发
				case ChannelTypeSMS:
					sendReq := sms.SendRequest{
						UserId:         vPeople.UserID,
						EnterpriseId:   vPeople.EnterpriseID,
						TemplateId:     vPeople.ChannelID,
						Phone:          vPeople.Mobile,
						ChannelType:    int64(vPeople.ChannelType),
						ChannelValueId: 0,
						Vars: sms.VarsInfo{
							UserName: vPeople.UserName,
							//UserNick:    vPeople.UserName,
							//AddressName: vPeople.AddressName,
							//OrderId:     vPeople.OrderID,
						},
					}
					request := &MarketStrategyRecord{
						ID:           vPeople.ID,
						UpdatedAt:    currentTime,
						Status:       ExecuteStatusFail,
						StrategyID:   currentStrategy.ID,
						EnterpriseID: enterpriseID,
						Mobile:       vPeople.Mobile,
					}

					if err = helper.WithRetry(2, time.Second*2, func() error {
						return smsService.SmsSend(c, sendReq)
					}); err != nil {
						log.Err(c, err).Msgf("发送营销-短信失败，策略id: %d,手机号：%s - 用户id%d", currentStrategy.ID, vPeople.Mobile, vPeople.UserID)
					} else {
						request.Status = ExecuteStatusSuccess
						request.TouchedAt = time.Now().Unix()
						if err := s.customerRepo.UpdateEnterpriseUserAttr(c, enterpriseID, vPeople.UserID, &customer.EnterpriseUserAttr{MarketCount: 1, LastMarketTime: time.Now().Unix()}); err != nil {
							log.Err(c, err).Msgf("增加用户影响次数失败-短信: %d - %d", enterpriseID, vPeople.UserID)
						}
					}

					if err = helper.WithRetry(2, time.Second*2, func() error {
						return s.repo.UpdateMarketStrategyRecord(c, request)
					}); err != nil {
						log.Error(c).Msgf("更新策略记录参数,错误", err)
					}
					break
				// 企业群发
				case ChannelTypeMassHair:
					massHairEntity = append(massHairEntity, *vPeople)
					break
				}
			}
			// 需要更新的策略ID
			if StrategyMethodALongTime != currentStrategy.StrategyMethod {
				if err := helper.WithRetry(2, time.Second*2, func() error {
					return s.repo.UpdateStrategyStatusByIds(c, StrategyStatusEnd, []int64{currentStrategy.ID})
				}); err != nil {
					log.Error(c).Msgf("通过策略IDs修改策略状态,错误", err)
				}
			}

			if len(massHairEntity) > 0 {
				if currentDomain == "" {
					log.Error(c).Msgf("无效的ucdomain", err)
					return
				}
				s.MassHair(c, enterpriseID, currentStrategy.ID, massHairEntity, currentDomain)
			}
		}(vStrategy, domain)
	}
	group.Wait()
}

// MarketTrigChannelExec 策略执行
func (s Service) MarketTrigChannelExec(c ctx.Context, params *MarketTrigExecRequest) (bool, error) {
	var sendError error
	// 发送消息
	switch params.ChannelType {
	// 短信
	case 1:
		vars := sms.VarsInfo{
			UserName: params.UserName,
			//AddressName: params.AddressName,
			//OrderId:     params.OrderID,
		}
		smsSendRequest := sms.SendRequest{
			UserId:         params.UserID,
			EnterpriseId:   params.EnterpriseID,
			TemplateId:     params.ChannelID,
			Phone:          params.Mobile,
			ChannelType:    int64(params.ChannelType),
			ChannelValueId: 0,
			Vars:           vars,
		}
		sendError = s.SmsService.SmsSend(c, smsSendRequest)
		if sendError != nil {
			return false, nil
		}
	// 企微群发
	case 2:
	// 自动发券
	case 3:
	//	公众号
	case 4:
	}

	return true, nil
}

func (s Service) UpdateStrategyTime(c ctx.Context, request *UpdateStatusRequest) (int64, error) {
	currentTime := time.Now().Unix()
	request.UpdatedAt = currentTime
	err := s.repo.UpdateStrategyTime(c, request)
	if err != nil {
		return 0, err
	}
	return request.ID, nil
}

//tagGroupList 标签列表
func (s Service) TagGroupList(c ctx.Context, userGroup UserGroup) ([]UserGroup, error) {
	list, err := s.repo.GetUserGroup(c, userGroup)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return nil, nil
	}
	// 查询客户标签列表
	return list, nil
}

// Perform 初始化策略触达人群记录
func (s Service) Perform(c ctx.Context, params MarketingStrategy, smsTemplate map[int64]*sms.SmsTemplate, currentTime int64) error {
	// 已初始化策略无需执行
	if HistoryStatusYes == params.HistoryStatus {
		log.Info(c).Msgf("企业【%d】策略【%d】属于客户策略-长期-不限-已执行，无需执行", params.EnterpriseID, params.ID)
		return nil
	}

	for _, channel := range params.Channels {
		req := &customer.ListRequest{
			ListCond: customer.ListCond{
				EnterpriseID: params.EnterpriseID,
			},
			Page: 1,
			Size: MaxSendNum / 10,
		}
		req.CustomerGroupCond.Include = true
		req.CustomerGroupCond.GroupIDList = []int64{params.CustomerGroupID}
		if channel.ChannelType == ChannelTypeSMS {
			req.CustomerSource = []int{customer.CustomerTypeMobile, customer.CustomerTypeDeliveryMobile}
		} else {
			req.CustomerType = customer.CustomerTypeExternal
		}

		total, users, err := s.esRepo.QueryUserList(c, req)
		if err != nil {
			return err
		}

		if total == 0 {
			return nil
		}

		for len(users) > 0 {
			var insertData []*MarketStrategyRecord
			for _, val := range users {
				rec := &MarketStrategyRecord{
					ID:           uid.ID(c),
					EnterpriseID: params.EnterpriseID,
					StrategyID:   params.ID,
					UserId:       val.ID,
					ChannelID:    channel.JuniorID,
					ChannelType:  channel.ChannelType,
					Mobile:       val.Mobile,
					StartTime:    params.StartTime,
					Status:       1,
					CreatedAt:    time.Now().Unix(),
					UpdatedAt:    time.Now().Unix(),
					UserName:     val.Name,
				}

				for _, unionID := range val.WxUnion {
					rec.UserUnionID = unionID
					break
				}

				// 无手机号
				if channel.ChannelType == ChannelTypeSMS && val.Mobile == "" {
					rec.Mobile = val.DeliveryMobile
					if rec.Mobile == "" {
						continue
					}
				} else if channel.ChannelType == ChannelTypeMassHair {
					for _, eu := range val.WxExternal {
						rec.ExternalUserid = eu.ThirdID
						break
					}
					if rec.ExternalUserid == "" {
						continue
					}
				}

				insertData = append(insertData, rec)
			}
			if len(insertData) > 0 {
				log.Info(c).Msgf("【营销策略短信目标用户---】%v", insertData)
				aErr := s.repo.BatchAddMarketStrategyRecord(c, insertData)
				if aErr != nil {
					log.Error(c).Msgf("批量插入数据错误，错误信息：%+v", aErr)
				}
			}

			req.Page++
			if (req.Page-1)*req.Size > int(total) {
				return nil
			}
			total, users, err = s.esRepo.QueryUserList(c, req)
			if err != nil {
				return err
			}
		}
	}

	if params.Condition == StrategyConditionTypeUnlimited {
		uErr := s.repo.UpdateStrategyHistoryStatusByIds(c, HistoryStatusYes, []int64{params.ID})
		if uErr != nil {
			<-time.After(2 * time.Second)
			_ = s.repo.UpdateStrategyHistoryStatusByIds(c, HistoryStatusYes, []int64{params.ID})
		}
	}

	return nil
}

// GetGroupMsgResult 获取群发消息发送结果
func (s Service) GetGroupMsgResult(c ctx.Context) error {
	// 获取待发送列表
	list, err := s.repo.SelectMsgIds(c)
	if err != nil || len(list) == 0 {
		log.Error(ctx.Wrap(c)).Msgf("获取群发消息发送结果:%+v", err)
		return err
	}
	enterpriseService := enterprise.New(enterpriseRepository.New())

	for _, v := range list {

		externalReq := work.MsgResultInfoRequest{
			QyGwRequest: work.QyGwRequest{},
			MsgID:       v.MsgID,
		}
		enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), enterprise.Corp{
			EnterpriseID: v.EnterpriseID,
		})
		if enterpriseErr != nil {
			log.Error(ctx.Wrap(c)).Msgf("查询企业信息失败,错误信息:%+v", enterpriseErr)

			return enterpriseErr
		}

		workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
		res, csaErr := workImpl.GetGroupMsgResult(c, externalReq)
		if csaErr != nil {
			log.Error(ctx.Wrap(c)).Msgf("调用企业微信,获取企业群发消息发送结果:%+v", csaErr)

			continue
		}

		for _, val := range res {
			if val.Status == 0 {
				continue
			}
			sendResult := QyMsgSendResult{
				EnterpriseID:   v.EnterpriseID,
				ExternalUserID: val.ExternalUserID,
				UserID:         val.UserID,
				MsgID:          v.MsgID,
				Status:         val.Status,
				SendTime:       val.SendTime,
			}

			// 更新群发发送记录表
			updateRes, err := s.repo.UpdateQyMsgSendResult(c, &sendResult)
			if err != nil {
				log.Error(ctx.Wrap(c)).Msgf("更新群发发送结果错误:%+v", err)
			}
			if updateRes != 0 {
				// 更新策略发送记录表
				switch val.Status {
				case 1:
					sendResult.Status = 2
				case 2, 3:
					sendResult.Status = 3
				default:
					continue
				}
				_, err = s.repo.UpdateMarketStrategyRecordByEcuserid(c, &sendResult)
				if err != nil {
					log.Error(ctx.Wrap(c)).Msgf("更新策略群发发送结果错误:%+v", err)
				}
			}

			// 更新营销时间和营销次数
			if val.Status == 1 {
				euid, _ := s.repo.GetQyMsgSendLogWithIdr(c, v.EnterpriseID, v.MsgID, val.ExternalUserID)
				if euid > 0 {
					if err := s.customerRepo.UpdateEnterpriseUserAttr(c, v.EnterpriseID, euid, &customer.EnterpriseUserAttr{MarketCount: 1, LastMarketTime: time.Now().Unix()}); err != nil {
						log.Err(c, err).Msgf("增加用户影响次数失败-企微: %d - %d", v.EnterpriseID, euid)
					}
				}

			}
		}
	}

	return nil
}

// MassHair 执行群发
func (s Service) MassHair(c ctx.Context, enterpriseID int64, strategyID int64, params []MarketTrigExecRequest, domain string) {
	currentTime := time.Now().Unix()
	var massHairLog, sliceMassHairLog []QyMsgSendLog
	var recordsID, sliceRecordsID []int64
	var externalUserid, sliceExternalUserid []string
	for _, item := range params {
		externalUserid = append(externalUserid, item.ExternalUserid)
		recordsID = append(recordsID, item.ID)
		// 因QyMsgSendLog表没有用，删除了，因此记录id直接变成0
		massHairLog = append(massHairLog, QyMsgSendLog{
			ExternalUserid: item.ExternalUserid,
			StrategyID:     strategyID,
			ChatType:       ChatTypeUser,
			EnterpriseID:   item.EnterpriseID,
			ChannelID:      item.ChannelID,
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		})
	}

	offset := MaxSendNum
	maxPage := int(math.Ceil(float64(len(externalUserid)) / float64(offset)))
	for page := 1; page <= maxPage; page++ {
		start := (page - 1) * offset
		end := page * offset
		// 到最后一页
		if page == maxPage {
			end = len(externalUserid)
		}

		sliceExternalUserid = externalUserid[start:end]
		sliceRecordsID = recordsID[start:end]
		sliceMassHairLog = massHairLog[start:end]

		log.Info(ctx.Wrap(c)).Msgf("sliceExternalUserid 数量: %d,sliceRecordsID 数量:%d,sliceMassHairLog:%d", len(sliceExternalUserid), len(sliceRecordsID), len(sliceMassHairLog))
		returnRep, sendErr := s.SendQyMsg(c, enterpriseID, strategyID, domain, sliceExternalUserid)
		if sendErr != nil {
			log.Info(c).Msgf("sendErr: %+v", sendErr)
		}
		s.UpdateMassHairRecords(c, sliceRecordsID, returnRep, sliceMassHairLog)
	}
}

// SendQyMsg 执行企业群发
func (s Service) SendQyMsg(c ctx.Context, enterpriseID int64, strategyID int64, domain string, externalUserid []string) (work.QyMassHairResponse, error) {
	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	var massHairReq = work.QyMassHairRequest{}
	massHairReq.ChatType = ChatTypeUser
	massHairReq.ExternalUserid = externalUserid
	// 查询策略发送内容
	var channelReq = ChannelInfo{
		StrategyID:  strategyID,
		ChannelType: ChannelTypeMassHair,
	}
	channelInfo, channelErr := s.repo.MarketTrigChannelDetailByWhere(c, channelReq)
	if errorIgnore.Ignore(channelErr) != nil || channelInfo == nil {
		log.Info(c).Msgf("channelErr: %+v", channelErr)
	}
	var Content ContentEntity
	_ = json.Unmarshal([]byte(channelInfo.TextContent), &Content)
	// 图片
	if Content.ContentType == MaterialTypePicture {
		massHairReq.Image = work.Image{
			PicUrl:  Content.Picture.Url,
			MediaId: "",
		}
		// 获取media_id
		ImgMediaId := redis.Client().Get(c, Content.Picture.Url).Val()
		if ImgMediaId != "" {
			massHairReq.Image.MediaId = ImgMediaId
		} else {
			cacheMediaRequest := drainage.CacheMediaRequest{
				EnterpriseId: enterpriseID,
				ImgUrl:       Content.Picture.Url,
				FileType:     "1",
			}
			cacheMediaResponse, err := s.DrainageService.CacheMedia(c, cacheMediaRequest, domain)
			if err != nil || cacheMediaResponse.Result == false {
				log.Info(c).Msgf("cacheMediaResponseErr: %+v", err)
			}
			massHairReq.Image.MediaId = cacheMediaResponse.MediaId
		}
	}
	// 链接
	if Content.ContentType == MaterialTypeLink {
		massHairReq.Link = work.Link{
			Title:    Content.Link.Title,
			WxPicUrl: Content.Link.CoverUrl,
			Desc:     Content.Link.Desc,
			Url:      Content.Link.LinkUrl,
		}
	}
	// 小程序
	if Content.ContentType == MaterialTypeApplet {
		massHairReq.MiniProgram = work.MiniProgram{
			Title:      Content.Applet.Title,
			PicMediaId: "",
			Page:       Content.Applet.Page,
			Appid:      Content.Applet.AppID,
		}
		// 将图片素材转换成企业微信的素材ID
		if Content.Applet.CoverUrl != "" {
			//获取media_id
			ImgMediaId := redis.Client().Get(c, Content.Applet.CoverUrl).Val()
			if ImgMediaId != "" {
				massHairReq.MiniProgram.PicMediaId = ImgMediaId
			} else {
				cacheMediaRequest := drainage.CacheMediaRequest{
					EnterpriseId: enterpriseID,
					ImgUrl:       Content.Applet.CoverUrl,
					FileType:     "1",
				}
				cacheMediaResponse, err := s.DrainageService.CacheMedia(c, cacheMediaRequest, domain)
				if err != nil || cacheMediaResponse.Result == false {
					log.Info(c).Msgf("cacheMediaResponseErr: %+v", err)
				}
				massHairReq.MiniProgram.PicMediaId = cacheMediaResponse.MediaId
			}
		}
	}
	if errorIgnore.Ignore(channelErr) != nil {
		log.Info(c).Msgf("channelErr: %+v", channelErr)
	}
	massHairReq.Text.Content = Content.Text
	returnRep, sendErr := workImpl.QwMassHair(ctx.Wrap(c), massHairReq)
	return *returnRep, sendErr
}

// IsUserGroupInMarketing 客户分组是否被营销策略引用
func (s Service) IsUserGroupInMarketing(c ctx.Context, groupID int64) (bool, error) {
	ret, err := s.repo.IsUserGroupUserInMarket(c, groupID)
	if err != nil {
		log.Info(c).Msgf("查询客户分组引用的营销策略数据失败, err:%s, groupID:%d", err.Error(), groupID)
		return false, err
	}
	return ret, nil
}

// UpdateMassHairRecords 企业群发
func (s Service) UpdateMassHairRecords(c ctx.Context, recordsID []int64, returnRep work.QyMassHairResponse, massHairLog []QyMsgSendLog) {
	currentTime := time.Now().Unix()
	// 更新记录表状态
	updateStatus := ExecuteStatusIng
	ResultStatus := int64(0)
	if returnRep.MsgID == "" {
		updateStatus = ExecuteStatusFail
		ResultStatus = ResultStatusFail
	}
	request := BatchUpdateRecordsRequest{
		IDs:       recordsID,
		UpdatedAt: currentTime,
		Status:    updateStatus,
		MsgID:     returnRep.MsgID,
	}
	err := s.repo.BatchUpdateMarketStrategyRecord(c, &request)
	if err != nil {
		log.Info(c).Msgf("BatchUpdateMarketStrategyRecordErr: %+v", err)
	}
	// 添加发送记录表
	for _, values := range massHairLog {
		values.MsgID = returnRep.MsgID
		resultReq := QyMsgSendResult{
			ID:             uid.ID(c),
			EnterpriseID:   values.EnterpriseID,
			SendLogID:      values.ID,
			ExternalUserID: values.ExternalUserid,
			MsgID:          returnRep.MsgID,
			CreatedAt:      currentTime,
			Status:         ResultStatus,
		}
		AddQySendMsgResultErr := s.repo.AddQySendMsgResult(c, &resultReq)
		if AddQySendMsgResultErr != nil {
			log.Info(c).Msgf("AddQySendMsgResultErr: %+v", AddQySendMsgResultErr)
		}
	}
}

func (s Service) GetSpuDetail(c ctx.Context, spuId int64) (SpuDetailRes, error) {
	spuRes, err := s.repo.GetSpuDetail(c, spuId)
	if err != nil {
		log.Error(c).Msgf("GetSpuDetail失败：%+v", err)
		return spuRes, err
	}
	return spuRes, err
}
func (s Service) GetWxChannelList(c ctx.Context) ([]*WxChannelInfo, error) {
	wxChannelList, err := s.repo.GetWxChannelList(c)
	if err != nil {
		return nil, err
	}
	return wxChannelList, err
}

//DealWaitStartStrategy 定时任务开始没有开始策略
func (s Service) DealWaitStartStrategy(c ctx.Context, enterpriseID int64) bool {
	nowTime := time.Now().Unix()
	var ids []int64
	// 查询所有策略列表
	marketingList, _ := s.repo.SelectMarketingWaitStartList(c, enterpriseID)

	for _, value := range *marketingList {
		// 如果开始时间大于当前时间跳过
		if value.StartTime > nowTime {
			continue
		}
		ids = append(ids, value.ID)
	}
	if len(ids) > 0 {
		// 更新策略为已开始
		_ = s.repo.UpdateStrategyStatusByIds(c, StrategyStatusStart, ids)
	}
	return true
}

//DealWaitEndStrategy 定时任务结束没有结束策略
func (s Service) DealWaitEndStrategy(c ctx.Context, enterpriseID int64) bool {
	nowTime := time.Now().Unix()
	var ids []int64
	// 查询所有策略列表
	marketingList, _ := s.repo.SelectMarketingWaitEndList(c, enterpriseID)

	for _, value := range *marketingList {
		// 如果开始时间大于当前时间跳过
		if value.EndTime > nowTime {
			continue
		}
		ids = append(ids, value.ID)
	}
	if len(ids) > 0 {
		// 更新策略为已结束
		_ = s.repo.UpdateStrategyStatusByIds(c, StrategyStatusEnd, ids)
	}
	return true
}

/** CheckGoOnStrategyTemplate
检测投放中的策略模版是否为审核成功且已赋值
1.获取企业正在投放中的短信策略
2.查询策略中的模版ID状态是否为已赋值
3.如果短信模版状态不等于已赋值，暂停该策略，标识该策略为异常停止
4.若策略异常手机号不为空，给该手机号发送异常提醒
*/
func (s *Service) CheckGoOnStrategyTemplate(c ctx.Context, enterpriseID int64, smsService *sms.Service) {
	resp, err := s.repo.GoOnStrategyTemplateIDs(c, enterpriseID)
	if err != nil {
		return
	}
	if len(*resp) > 0 {
		// 查询运行中的策略模版ID
		var templateIDs, strategyIDs, channelStrategyIDs, remindStrategyID []int64
		strategyMap := make(map[int64]GetGoOnStrategyResp)
		for _, item := range *resp {
			if false == helper.InArray[int64](item.JuniorID, templateIDs) {
				templateIDs = append(templateIDs, item.JuniorID)
			}
			strategyMap[item.StrategyID] = *item
			channelStrategyIDs = append(channelStrategyIDs, item.StrategyID)
		}
		// 查询短信模版的状态
		smsTemplateList, err := smsService.SelectTemplateListByIDs(c, enterpriseID, templateIDs, channelStrategyIDs)
		if err != nil {
			log.Info(c).Msgf("查询短信的状态错误,err:%s", err.Error())
			return
		}
		if len(*smsTemplateList) > 0 {
			for _, item := range *smsTemplateList {
				if item.Status != sms.TemplateStatusHasValue {
					if data, ok := strategyMap[item.StrategyID]; ok {
						strategyIDs = append(strategyIDs, data.StrategyID)
						if data.Mobile != "" {
							remindStrategyID = append(remindStrategyID, data.StrategyID)
						}
					}
				}
			}
		}
		if len(remindStrategyID) == 0 {
			log.Info(c).Msgf("企业【%d】没有要发送提醒的策略", enterpriseID)
			return
		}
		group := threading.NewRoutineGroup()
		// 批量跟新策略为暂停
		group.RunSafe(func() {
			if len(strategyIDs) > 0 {
				_ = s.BatchUpdateStrategy(c, strategyIDs)
			}
		})
		// 发送短信
		group.RunSafe(func() {
			// 查询企业异常策略提醒模板
			smsTemplate, _ := smsService.GetDefaultTemplate(c, enterpriseID, sms.ReappearingName)
			if err != nil {
				return
			}
			if len(remindStrategyID) > 0 {
				for _, strategyID := range remindStrategyID {
					if strategyInfo, ok := strategyMap[strategyID]; ok {
						smsSendRequest := sms.SendRequest{
							EnterpriseId:   enterpriseID,
							TemplateId:     smsTemplate.Id,
							Phone:          strategyInfo.Mobile,
							ChannelType:    0,
							ChannelValueId: 0,
						}
						// 发送短信
						err = smsService.SmsSend(c, smsSendRequest)
						if err != nil {
							log.Info(c).Msgf("异常短信发送失败,err:%s", err.Error())
							continue
						}
					}
				}
			}
		})
		group.Wait()
	}
}

// 短信模版异常，批量更新策略状态为暂停
func (s *Service) BatchUpdateStrategy(c ctx.Context, strategyIDs []int64) error {
	err := s.repo.BatchUpdateStrategy(c, strategyIDs, StrategyStatusStop, AbnormalPauseYes)
	if err != nil {
		log.Info(c).Msgf("批量跟新策略为暂停,err:%s", err.Error())
		return err
	}
	return nil
}

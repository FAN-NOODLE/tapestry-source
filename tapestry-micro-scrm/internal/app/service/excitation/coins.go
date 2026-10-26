package excitationSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	"micro-scrm/internal/app/http/request"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	openReq "micro-scrm/internal/app/http/request/open"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"strconv"
	"strings"
	"time"
)

type CoinsService struct {
	service.Base
	intervalSeconds  int64 // 定时任务执行间隔，这里需要判断是否要触发延迟队列清空积分
	CoinsSettingsDao *excitationDao.CoinsSettings
}

func NewCoinsService() *CoinsService {
	return &CoinsService{
		intervalSeconds:  3600, // 一小时执行一次
		CoinsSettingsDao: excitationDao.NewCoinsSettings(),
	}
}

// AddOrder  下单
func (s *CoinsService) AddOrder(c ctx.Context, params openReq.AddOrderReq) error {
	mallOrderDao := excitationDao.NewMallOrder()
	memberDataDao := excitationDao.NewMemberData()
	if params.DeductionCoins <= 0 {
		return errors.New("下单金额异常")
	}
	orderInfo, err := mallOrderDao.GetOne(c, "enterprise_id=? AND order_num = ?", params.EnterpriseId, params.OrderNum)
	if err != nil {
		return err
	}
	if orderInfo.ID > 0 {
		return errors.New("订单号已存在")
	}
	memberId, _ := membersDao.NewMembers().GetMemberIdByUserId(c, params.EnterpriseId, params.UserId)
	if memberId == 0 {
		return errors.New("导购编码异常")
	}

	memberDataInfo, err := memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, memberId)
	if memberDataInfo.Coins < params.DeductionCoins {
		return errors.New("导购金币不足")
	}

	err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
		// 生成订单
		err = mallOrderDao.WithTx(tx).CreateData(c, excitationModels.ExcitationMallOrderModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			UserID:         params.UserId,
			MemberID:       memberId,
			OrderNum:       params.OrderNum,
			OrderCoins:     params.DeductionCoins,
			OrderStatus:    excitationModels.OrderStatusNormal,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("创建订单失败,原因:%v", err))
		}
		// 处理导购金币更新
		params.DeductionCoins, err = memberDataDao.SaveCoins(c, excitationReq.UpdateCoinsReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: params.EnterpriseId,
			},
			MemberId: memberId,
			Coins:    -params.DeductionCoins,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("下单更新积分失败,原因:%v", err))
		}

		err = excitationDao.NewCoinsRecord().WithTx(tx).CreateData(c, params.EnterpriseId, microScrmModel.ExcitationCoinsRecordModel{
			BaseModel:   models.CreateBase(app.UID()),
			MemberID:    memberId,
			Coins:       fmt.Sprintf("%0.1f", params.DeductionCoins),
			SourceType:  microScrmModel.ExcitationCoinsRecordAutoType_Mall,
			CoinsRemark: "兑换奖品",
			OrderNum:    params.OrderNum,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("下单积分扣除记录创建失败,原因:%v", err))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// CancelOrder 退单
func (s *CoinsService) CancelOrder(c ctx.Context, params openReq.CancelOrderReq) error {
	mallOrderDao := excitationDao.NewMallOrder()
	memberDataDao := excitationDao.NewMemberData()
	orderInfo, err := mallOrderDao.GetOne(c, "enterprise_id=? AND order_num = ?", params.EnterpriseId, params.OrderNum)
	if err != nil {
		return err
	}
	if orderInfo.ID == 0 || orderInfo.UserID != params.UserId {
		return errors.New("订单号或导购编码不正确")
	}
	if orderInfo.OrderStatus != excitationModels.OrderStatusNormal {
		return errors.New("订单状态不正确")
	}
	err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		err = mallOrderDao.WithTx(tx).UpdateData(c, map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"order_num":     orderInfo.OrderNum,
		}, map[string]interface{}{
			"order_status": excitationModels.OrderStatusCancel,
			"updated_at":   time.Now().Unix(),
		})
		if err != nil {
			return errors.New(fmt.Sprintf("退单失败,原因:%v", err))
		}
		// 更新导购金币余额
		orderInfo.OrderCoins, err = memberDataDao.SaveCoins(c, excitationReq.UpdateCoinsReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: params.EnterpriseId,
			},
			MemberId: orderInfo.MemberID,
			Coins:    orderInfo.OrderCoins,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("导购金币更新失败,原因:%v", err))
		}
		// 处理金币回归记录
		err = excitationDao.NewCoinsRecord().WithTx(tx).CreateData(c, params.EnterpriseId, microScrmModel.ExcitationCoinsRecordModel{
			BaseModel:   models.CreateBase(app.UID()),
			MemberID:    orderInfo.MemberID,
			Coins:       fmt.Sprintf("%0.1f", orderInfo.OrderCoins),
			SourceType:  microScrmModel.ExcitationCoinsRecordAutoType_Mall,
			CoinsRemark: "取消兑换订单",
			OrderNum:    orderInfo.OrderNum,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("下单积分扣除记录创建失败,原因:%v", err))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// OrderInfo  订单详情查询
func (s *CoinsService) OrderInfo(c ctx.Context, params openReq.OrderInfoReq) (res openResp.OrderInfoResp, err error) {

	orderInfo, err := excitationDao.NewMallOrder().GetOne(c, "enterprise_id=? AND order_num = ?", params.EnterpriseId, params.OrderNum)
	if err != nil {
		return res, err
	}
	if orderInfo.ID == 0 || orderInfo.UserID != params.UserId {
		return res, errors.New("订单号或导购编码不正确")
	}
	return openResp.OrderInfoResp{
		OrderStatus: orderInfo.OrderStatus,
		UserId:      orderInfo.UserID,
		OrderNum:    orderInfo.OrderNum,
		OrderCoins:  orderInfo.OrderCoins,
		CreatedAt:   orderInfo.CreatedAt,
		UpdatedAt:   orderInfo.UpdatedAt,
	}, nil
}

// GetCoinsList 小程序金币记录列表
func (s *CoinsService) GetCoinsList(c ctx.Context, params excitationReq.GetCoinsListReq) (resList []excitationResp.GetCoinsListResp, count int64, err error) {
	// 查询结果
	paramsWhere := excitationReq.CoinsRecordListReq{
		MemberSearchReq: request.MemberSearchReq{
			EnterpriseId: params.EnterpriseId,
			MemberIdArr: []int64{
				params.MemberId,
			},
		},
	}
	list, count, err := excitationDao.NewCoinsRecord().GetPagination(c, paramsWhere)
	if err != nil {
		return resList, 0, err
	}
	// 整理最终返回的数据
	for _, record := range list {
		resList = append(resList, excitationResp.GetCoinsListResp{
			Id:          record.ID,
			CoinsRemark: record.CoinsRemark,
			Coins:       helper.RuleThree(strings.HasPrefix(record.Coins, "-"), record.Coins, "+"+record.Coins),
			CreatedAt:   record.CreatedAt,
		})
	}
	return resList, count, nil
}

// GetCoinsSettingInfo 获取金币配置
func (s *CoinsService) GetCoinsSettingInfo(c ctx.Context, req excitationReq.GetCoinSettingInfoReq) (res excitationResp.CoinsSettingInfoResp, err error) {
	info, err := excitationDao.NewCoinsSettings().GetOne(c, "enterprise_id= ?", req.EnterpriseId)
	if err != nil {
		return excitationResp.CoinsSettingInfoResp{}, err
	}
	var expireList []excitationModels.ExcitationCoinsExpireRecordModel
	var editRecord []excitationModels.ExcitationCoinsSettingsEditRecordModel
	if req.HasRecord {
		// 查询金币过期记录
		expireList, _ = excitationDao.NewCoinsExpireRecord().GetMany(c, fmt.Sprintf("enterprise_id = %d", req.EnterpriseId))
		// 获取金币配置修改记录
		editRecord, _ = excitationDao.NewCoinsSettingsEditRecord().GetMany(c, fmt.Sprintf("enterprise_id = %d", req.EnterpriseId))
		if err != nil {
			return excitationResp.CoinsSettingInfoResp{}, err
		}
	}

	ratioSlice := strings.Split(info.RatioRemark, ":")
	res = excitationResp.CoinsSettingInfoResp{
		MallIsOpen:              info.MallIsOpen,
		RatioPoints:             helper.StringToInt(ratioSlice[0]),
		RatioCoins:              helper.StringToInt(ratioSlice[1]),
		CoinsExpireIsOpen:       info.ExpireIsOpen,
		CoinsExpireTime:         info.ExpireAt,
		CoinsExpireRecordList:   expireList,
		CoinsSettingsEditRecord: editRecord,
	}
	return
}

// SaveCoinsSetting 设置金币配置
func (s *CoinsService) SaveCoinsSetting(c ctx.Context, params excitationReq.SaveCoinsSettingsReq) error {
	coinsSettingsDao := excitationDao.NewCoinsSettings()
	info, err := coinsSettingsDao.GetOne(c, "enterprise_id= ?", params.EnterpriseId)
	if err != nil {
		return err
	}
	ratioRemark := fmt.Sprintf("%d:%d", params.RatioPoints, params.RatioCoins)
	ratio := helper.MathRounded(float64(params.RatioCoins)/float64(params.RatioPoints), 6)
	expireAt := helper.RuleThree(params.CoinsExpireIsOpen == excitationModels.CoinsExpireOpen, params.CoinsExpireTime, 0)
	if info.ID > 0 {
		err := coinsSettingsDao.UpdateData(c, fmt.Sprintf("enterprise_id = %d", params.EnterpriseId), map[string]interface{}{
			"ratio_remark":   ratioRemark,
			"ratio":          ratio,
			"mall_is_open":   params.MallIsOpen,
			"expire_is_open": params.CoinsExpireIsOpen,
			"expire_at":      expireAt,
		})
		if err != nil {
			return err
		}
	} else {
		err = coinsSettingsDao.CreateData(c, excitationModels.ExcitationCoinsSettingsModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			Ratio:          ratio,
			RatioRemark:    ratioRemark,
			MallIsOpen:     params.MallIsOpen,
			ExpireAt:       expireAt,
			ExpireIsOpen:   params.CoinsExpireIsOpen,
		})
		if err != nil {
			return err
		}
	}
	adminInfo, err := microBossDao.NewEnterpriseAdmin().GetInfo(c, microBossModel.EnterpriseAdminModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: params.EnterpriseId,
			BaseModel: models.BaseModel{
				ID: params.AdminId,
			},
		},
	})

	if info.RatioRemark != ratioRemark {
		ratioSlice := strings.Split(ratioRemark, ":")
		// 积分比例修改记录
		err = excitationDao.NewCoinsSettingsEditRecord().CreateData(c, excitationModels.ExcitationCoinsSettingsEditRecordModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			AdminID:        params.AdminId,
			AdminName:      adminInfo.Username,
			RatioRemark:    fmt.Sprintf("%s积分兑换%s金币", ratioSlice[0], ratioSlice[1]),
			MallIsOpen:     params.MallIsOpen,
			ExpireAt:       expireAt,
		})
		if err != nil {
			return err
		}
	}
	if expireAt > 0 && expireAt != info.ExpireAt {
		// 如果过期时间有调整则触发一次金币过期的任务（避免新设置的过期时间小于定时任务执行间隔）
		err = s.CoinsExpiredProduce(c, params.EnterpriseId)
		if err != nil {
			return err
		}
	}

	return nil
}

// ManualCoins 手动加金币
func (s *CoinsService) ManualCoins(c ctx.Context, enterpriseId int64, dataList []excitationReq.ManualCoinsRow) (err error) {
	var coinsRecordList []microScrmModel.ExcitationCoinsRecordModel
	memberDataDao := excitationDao.NewMemberData()
	err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
		for _, record := range dataList {
			if record.Coins == 0 {
				//忽略金币为0的数据
				continue
			}
			// 处理导购金币更新,考虑金币不够的情况返回真实扣减金币
			record.Coins, err = memberDataDao.WithTx(tx).SaveCoins(c, excitationReq.UpdateCoinsReq{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: enterpriseId,
				},
				MemberId: record.MemberId,
				Coins:    record.Coins,
			})
			if err != nil {
				log.Info(c).Msgf("手动更新导购金币数据失败，企业:%d,错误:%v", enterpriseId, err)
				return err
			}
			tmp := microScrmModel.ExcitationCoinsRecordModel{
				BaseModel:   models.CreateBase(app.UID()),
				MemberID:    record.MemberId,
				Coins:       fmt.Sprintf("%0.1f", record.Coins),
				SourceType:  microScrmModel.ExcitationCoinsRecordAutoType_Hand,
				CoinsRemark: record.CoinsRemark,
			}
			coinsRecordList = append(coinsRecordList, tmp)
		}
		err = excitationDao.NewCoinsRecord().WithTx(tx).CreateInBatchesData(c, enterpriseId, coinsRecordList)
		return err
	})
	if err != nil {
		log.Info(c).Msgf("批量写入金币记录失败，企业:%d,错误:%v", enterpriseId, err)
		return err
	}
	return nil
}

// GetRecordList 获取导购积分获取记录列表
func (s *CoinsService) GetRecordList(c ctx.Context, params excitationReq.CoinsRecordListReq) (resList []excitationResp.CoinsRecordList, count int64, err error) {
	// 处理通用导购部门等参数
	searchExists, memberIdArr := s.Base.GetSearchMemberIds(c, params.MemberSearchReq)
	if !searchExists {
		// 如果检索条件查不到用户直接返回空数据
		return resList, count, nil
	}
	params.MemberIdArr = memberIdArr
	// 查询结果
	list, count, err := excitationDao.NewCoinsRecord().GetPagination(c, params)
	if err != nil {
		return resList, 0, err
	}
	// 整理结果中的memberId用户查询导购信息
	var memberIds []int64
	for _, record := range list {
		memberIds = append(memberIds, record.MemberID)
	}
	memberIds = helper.ArrayUniq(memberIds)
	// 查询导购信息
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, params.EnterpriseId, memberIds)
	if err != nil {
		return nil, 0, err
	}

	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]excitationResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		memberMap[info.ID] = excitationResp.MemberInfo{
			MemberUserId:   info.Userid,
			MemberName:     info.Name,
			DepartmentName: dpName,
		}
	}
	changTypeStr := make(map[int]string, 2)
	changTypeStr[excitationReq.ChangeTypeInc] = "增加"
	changTypeStr[excitationReq.ChangeTypeDec] = "减少"
	// 整理最终返回的数据
	for _, record := range list {
		if _, ok := memberMap[record.MemberID]; ok {
			coinsFloat, err := strconv.ParseFloat(record.Coins, 64)
			if err != nil {
				return nil, 0, err
			}
			changeType := helper.RuleThree(coinsFloat > 0, excitationReq.ChangeTypeInc, excitationReq.ChangeTypeDec)
			resList = append(resList, excitationResp.CoinsRecordList{
				Id:             record.ID,
				CoinsRemark:    record.CoinsRemark,
				Coins:          coinsFloat,
				SourceType:     record.SourceType,
				CreatedAt:      record.CreatedAt,
				MemberInfo:     memberMap[record.MemberID],
				CreatedTimeStr: helper.GetFormatDateByParams(record.CreatedAt),
				SourceTypeStr:  getCoinsSourceTypeStr(record.SourceType),
				OrderNum:       record.OrderNum,
				ChangeType:     changeType,
				ChangeTypeStr:  changTypeStr[changeType],
			})
		}
	}
	return resList, count, nil
}

func getCoinsSourceTypeStr(SourceType int) string {
	var sourceTypeMap = make(map[int]string)
	sourceTypeMap[microScrmModel.ExcitationCoinsRecordAutoType_Points] = "积分同步"
	sourceTypeMap[microScrmModel.ExcitationCoinsRecordAutoType_Mall] = "兑礼中心"
	sourceTypeMap[microScrmModel.ExcitationCoinsRecordAutoType_Hand] = "手动补发"
	sourceTypeMap[microScrmModel.ExcitationCoinsRecordAutoType_Expire] = "金币过期"
	if _, ok := sourceTypeMap[SourceType]; ok {
		return sourceTypeMap[SourceType]
	}
	return ""
}

// CoinsExpiredProduce 金币过期生产者
func (s *CoinsService) CoinsExpiredProduce(c ctx.Context, enterpriseId int64) error {
	info, err := s.CoinsSettingsDao.GetOne(c, "enterprise_id= ?", enterpriseId)
	if err != nil {
		return err
	}
	expire, err := s.coinsExpiredCheck(c, info)
	if err != nil {
		return err
	}
	if !expire {
		// 不需要执行
		return nil
	}
	nowUnix := time.Now().Unix()
	if info.ExpireAt < nowUnix {
		// 过期时间已经过了
		return nil
	}
	delaySeconds := info.ExpireAt - nowUnix
	if delaySeconds > s.intervalSeconds {
		// 还不到需要执行延迟队列的时间
		return nil
	}
	// 写入kafka队列，消费时会自动延迟
	produceParams := excitationReq.CoinExpiredReq{
		EnterpriseId: enterpriseId,
		DelaySeconds: delaySeconds,
	}
	str, _ := json.Marshal(produceParams)
	err = pkgKafka.GetAppKafka().Producer(string(str), "coins_expired")
	if err != nil {
		log.Info(c).Msgf("金币过期写入执行队列失败,企业id:%d,过期时间:%d,错误原因：%+v", enterpriseId, info.ExpireAt, err)
	}
	return nil
}

// CoinsExpiredConsume 金币过期消费者
func (s *CoinsService) CoinsExpiredConsume(c ctx.Context, enterpriseId int64, delaySeconds int64) error {
	// 延迟时间
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	// 针对触发时的并发锁,避免多次执行定时任务触发多个延迟队列
	lockKey := fmt.Sprintf("lock:excitation:coin_expired:%d", enterpriseId)
	orderLock, lErr := lock.Mutex(
		c,
		lock.Key(lockKey),
		lock.TTL(10*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if lErr != nil {
		return lErr
	}
	if err := orderLock.Lock(); err != nil {
		return err
	}
	defer func() {
		_ = orderLock.Unlock()
	}()
	// 获取过期配置
	info, err := s.CoinsSettingsDao.GetOne(c, "enterprise_id= ?", enterpriseId)
	// 执行验证
	expire, err := s.coinsExpiredCheck(c, info)
	if err != nil {
		return err
	}
	if !expire {
		// 不需要执行
		return nil
	}
	if info.ExpireAt > time.Now().Unix() {
		// 如果过期时间大于当前时间，说明还没到时间（前面已经sleep了足够的时间，按道理应该是和当前时间相等或小于当前时间的）
		return nil
	}
	memberDataDao := excitationDao.NewMemberData()
	// 获取所有金币大于0的人
	memberDataList, err := memberDataDao.GetMany(c, "enterprise_id = ? AND coins > 0", enterpriseId)
	if err != nil {
		return err
	}
	var memberIds []int64
	for _, val := range memberDataList {
		memberIds = append(memberIds, val.MemberID)
	}
	err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
		var coinsRecordList []microScrmModel.ExcitationCoinsRecordModel
		for _, dataInfo := range memberDataList {
			// 处理导购金币更新
			dataInfo.Coins, err = memberDataDao.WithTx(tx).SaveCoins(c, excitationReq.UpdateCoinsReq{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: enterpriseId,
				},
				MemberId: dataInfo.MemberID,
				Coins:    -dataInfo.Coins,
			})
			// 组装金币明细数据
			tmp := microScrmModel.ExcitationCoinsRecordModel{
				BaseModel:   models.CreateBase(app.UID()),
				MemberID:    dataInfo.MemberID,
				Coins:       fmt.Sprintf("%0.1f", dataInfo.Coins),
				SourceType:  microScrmModel.ExcitationCoinsRecordAutoType_Expire,
				CoinsRemark: "金币过期",
			}
			coinsRecordList = append(coinsRecordList, tmp)
		}
		// 批量写入金币过期的金币记录
		err = excitationDao.NewCoinsRecord().WithTx(tx).CreateInBatchesData(c, enterpriseId, coinsRecordList)
		if err != nil {
			log.Info(c).Msgf("金币过期记录写入失败，企业:%d,错误:%v", enterpriseId, err)
			return err
		}
		err = excitationDao.NewCoinsExpireRecord().WithTx(tx).CreateData(c, excitationModels.ExcitationCoinsExpireRecordModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
			ExpiredAt:      info.ExpireAt,
		})
		if err != nil {
			return err
		}
		// 金币过期执行成功 自动设置为关闭
		err = s.CoinsSettingsDao.UpdateData(c, fmt.Sprintf("enterprise_id = %d", enterpriseId), map[string]interface{}{
			"expire_is_open": excitationModels.CoinsExpireClose,
			"expire_at":      0,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *CoinsService) coinsExpiredCheck(c ctx.Context, info excitationModels.ExcitationCoinsSettingsModel) (bool, error) {
	if info.ExpireIsOpen != excitationModels.CoinsExpireOpen || info.ExpireAt <= 0 {
		// 如果金币过期未设置直接跳过
		return false, nil
	}
	// 查询历史执行记录
	expireRecord, err := excitationDao.NewCoinsExpireRecord().GetMany(c, fmt.Sprintf("enterprise_id = %d", info.EnterpriseID))
	if err != nil {
		return false, err
	}
	for _, val := range expireRecord {
		if val.ExpiredAt >= info.ExpireAt {
			// 存在执行记录或者有执行记录的时间大于本次执行的时间说明已经执行过了
			return false, nil
		}
	}
	return true, nil
}

// CoinsInit 金币初始化
func (s *CoinsService) CoinsInit(c ctx.Context, enterpriseId int64) {
	var coinsRecordList []microScrmModel.ExcitationCoinsRecordModel
	memberDataList, err := excitationDao.NewMemberData().GetMany(c, "enterprise_id = ? AND point_value > 0", enterpriseId)
	updateData := map[string]interface{}{
		"coins":       gorm.Expr("point_value"),
		"total_coins": gorm.Expr("point_value"),
		"spend_coins": 0,
	}
	err = excitationDao.NewMemberData().UpdateData(c,
		map[string]interface{}{
			"enterprise_id": enterpriseId,
		}, updateData)
	for _, info := range memberDataList {
		tmp := microScrmModel.ExcitationCoinsRecordModel{
			BaseModel:   models.CreateBase(app.UID()),
			MemberID:    info.MemberID,
			Coins:       info.PointValue,
			SourceType:  microScrmModel.ExcitationCoinsRecordAutoType_Points,
			CoinsRemark: "金币初始化",
		}
		coinsRecordList = append(coinsRecordList, tmp)
	}
	err = excitationDao.NewCoinsRecord().CreateInBatchesData(c, enterpriseId, coinsRecordList)

	if err != nil {
		log.Info(c).Msgf("批量写入金币记录失败，企业:%d,错误:%v", enterpriseId, err)
	}
}

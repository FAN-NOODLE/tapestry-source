package customerSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	openReq "micro-scrm/internal/app/http/request/open"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	"micro-scrm/internal/app/models"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	tagsSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/workwechat"
	"strconv"
	"strings"
	"time"
)

type ExternalSearch struct {
	service.Base
	externalRelationshipCkDao        *externalSearchDao.ExternalRelationshipCk
	externalExtraCkDao               *externalSearchDao.ExternalExtraCk
	externalCustomerExtensionDao     *externalSearchDao.ExternalCustomerExtension
	externalRelationshipExtensionDao *externalSearchDao.ExternalRelationshipExtension
}

func NewExternalSearch() *ExternalSearch {
	return &ExternalSearch{
		externalRelationshipCkDao:        externalSearchDao.NewExternalRelationshipCk(),
		externalExtraCkDao:               externalSearchDao.NewExternalExtraCk(),
		externalCustomerExtensionDao:     externalSearchDao.NewExternalCustomerExtension(),
		externalRelationshipExtensionDao: externalSearchDao.NewExternalRelationshipExtension(),
	}
}
func (s *ExternalSearch) GetExternalSearchListCk(c app.Ctx, request customerSearchReq.ExternalSearchListBaseReq) ([]customerSearchResp.ExternalUserList, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserList, 0)
	// 判断当前员工类型 店员或店长
	memberIds, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	if err != nil {
		return list, 0, err
	}
	request.MemberIds = memberIds
	var saCode []string
	if request.IsRelationship {
		//顾客关系列表
		list, total, err = s.externalExtraCkDao.RelationShipListCk(c, request)
		for _, v := range list {
			saCode = append(saCode, v.UserId)
		}
	} else {
		//顾客列表
		list, total, err = s.externalExtraCkDao.ExternalSearchListCk(c, request)
		for _, v := range list {
			saCode = append(saCode, v.BaCode)
		}
	}
	if err != nil || total == 0 {
		return list, 0, err
	}
	// 整理keysa 的导购信息
	saMap := make(map[string]membersModel.MemberModel, 0)
	saLists, err := membersDao.NewMembers().GetManyByUserIds(c, request.EnterpriseID, saCode)
	if err != nil {
		return list, total, err
	}
	for _, v := range saLists {
		saMap[v.Userid] = v
	}
	for i, v := range list {
		if exist, ok := saMap[v.BaCode]; ok {
			list[i].KeySaName = exist.Name
		}
	}
	return list, total, err
}

func (s *ExternalSearch) GetExternalSearchCk(c app.Ctx, request customerSearchReq.ExternalSearchListBaseReq, pluckName string) ([]interface{}, error) {
	// 判断当前员工类型 店员或店长
	memberIds, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	if err != nil {
		return nil, err
	}
	request.MemberIds = memberIds
	return s.externalExtraCkDao.ExternalSearchPluckIds(c, request.SearchList, pluckName)
}

// GetMembersByManageType 根据导购或店长类型获取店员
func (s *ExternalSearch) GetMembersByManageType(c app.Ctx, enterpriseID, memberId int64, thirdIds []int, memberIds []int64) ([]int64, error) {
	if len(memberIds) > 0 {
		return memberIds, nil
	} else if len(thirdIds) > 0 {
		memberIds, err := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, enterpriseID, thirdIds, []int64{})
		if err != nil {
			return []int64{}, err
		}
		return helper.ArrayUniq(append(memberIds, memberId)), nil
	}
	return []int64{memberId}, nil
}

// SetReserveExtraData  顾客列表附加检索信息如表更新
func (s *ExternalSearch) SetReserveExtraData(c app.Ctx, enterpriseId int64) (total int, err error) {
	size := 3000
	var lastId int64
	externalExtraCk := externalSearchDao.NewExternalExtraCk()
	reserveOrm := microScrmDao.NewReservesOrder()
	for {
		// 只查当天的数据，因为目前存的是邀约相关的时间，检索的时候会自动判断是不是符合条件，所以不需要查30天的
		reserveList, err := reserveOrm.GetReserveByDays(c, enterpriseId, lastId, size)
		if err != nil {
			return 0, err
		}
		if len(reserveList) == 0 {
			break
		}
		// 邀约数据
		var uniqueMap = make(map[string]externalSearchDao.ExternalExtraData)
		for _, info := range reserveList {
			if existAt, ok := uniqueMap[info.UserUnionid]; ok {
				if info.CreatedAt > int64(existAt.ValidReserveAt) {
					existAt.ValidReserveAt = int(info.CreatedAt)
				}
				if info.ReserveOrderTime > existAt.SuccessReserveAt {
					existAt.SuccessReserveAt = info.ReserveOrderTime
				}
				uniqueMap[info.UserUnionid] = existAt
			} else {
				uniqueMap[info.UserUnionid] = externalSearchDao.ExternalExtraData{
					UserUnionid:      info.UserUnionid,
					ValidReserveAt:   int(info.CreatedAt),
					SuccessReserveAt: info.ReserveOrderTime,
				}
			}
		}
		// 业务上查询的是30天的数据,所以根据业务考虑更新频率
		err = externalExtraCk.SaveReserveDataByAll(c, enterpriseId, uniqueMap)
		if err != nil {
			return 0, err
		}
		s.saveReserveExtension(c, enterpriseId, uniqueMap)
		total += len(uniqueMap)
		// 获取最后一条记录的id
		lastId = reserveList[len(reserveList)-1].ID
	}
	return
}

// SetTaskExtraData  顾客列表附加检索信息入表更新
func (s *ExternalSearch) SetTaskExtraData(c app.Ctx, enterpriseId int64, days int) (total int, err error) {
	size := 3000
	externalExtraCk := externalSearchDao.NewExternalExtraCk()
	taskOrm := microScrmDao.NewPlanTaskConditionDataDao()
	// 处理任务数据, 这里任务分表后 要查询 30天内的所有新增的子任务找到主任务，循环查询触达情况
	planList, err := microScrmDao.NewPlanTaskDao().
		GetMany(c, enterpriseId, "task_start_time > ? ", time.Now().Unix()-int64(days*86400))
	if err != nil {
		return 0, err
	}
	var planIds []int64
	for _, planInfo := range planList {
		planIds = append(planIds, planInfo.PlanID)
	}
	planIds = helper.ArrayUniq(planIds)
	for _, planId := range planIds {
		var lastId int64
		for {
			// 任务的触达标识后就不会再变化了，只取触达时间为一天内的数据即可
			taskList, err := taskOrm.GetTaskByDays(c, enterpriseId, planId, lastId, size, 1, microScrmModel.TaskReach)
			if err != nil {
				return 0, err
			}
			if len(taskList) == 0 {
				break
			}
			// 任务数据
			var uniqueUnionIdMap = make(map[string]externalSearchDao.ExternalExtraData)
			for _, info := range taskList {
				if info.UnionId == "" {
					// 过滤unionid为空的顾客
					continue
				}
				if existAt, ok := uniqueUnionIdMap[info.UnionId]; ok {
					if info.FinishedAt > existAt.TaskFinishAt {
						uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraData{
							UserUnionid:  info.UnionId,
							TaskFinishAt: info.FinishedAt,
						}
					}
				} else {
					uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraData{
						UserUnionid:  info.UnionId,
						TaskFinishAt: info.FinishedAt,
					}
				}
			}
			// 更新任务相关数据
			err = externalExtraCk.SaveTaskData(c, enterpriseId, uniqueUnionIdMap, microScrmModel.TaskReach)
			if err != nil {
				return 0, err
			}
			s.saveTaskExtension(c, enterpriseId, uniqueUnionIdMap)
			total += len(uniqueUnionIdMap)
			// 获取最后一条记录的id
			lastId = taskList[len(taskList)-1].ID
		}
	}

	return
}

// BatchSaveCouponExtraData    顾客列表附加检索信息优惠券相关更新
func (s *ExternalSearch) BatchSaveCouponExtraData(c app.Ctx, req openReq.BatchUpdateCouponReq) error {
	var dataMap = make(map[string]openReq.UpdateCouponParams)
	for _, row := range req.DataList {
		if row.UnionId == "" {
			continue
		}
		dataMap[row.UnionId] = row
	}
	if len(dataMap) == 0 {
		return nil
	}
	err := externalSearchDao.NewExternalExtraCk().BatchSaveCouponData(c, req.EnterpriseId, dataMap)
	if err != nil {
		return err
	}
	s.saveCouponExtension(c, req.EnterpriseId, dataMap)
	// 数据合并
	_ = s.ExtraDataMerge(c, req.EnterpriseId, false)
	return nil
}

// ExtraDataMerge 数据合并
func (s *ExternalSearch) ExtraDataMerge(c app.Ctx, enterpriseId int64, force bool) error {
	unixTime := time.Now().Unix()
	couponUpdateTime := fmt.Sprintf("customer:extra:update-time-%d", enterpriseId)
	minuteTimesStr := s.Redis().Get(c, couponUpdateTime).Val()
	// 合并数据间隔一般需要大于5分钟
	if force || minuteTimesStr == "" {
		err := s.Redis().Set(c, couponUpdateTime, unixTime, time.Minute*10).Err()
		if err != nil {
			return err
		}
		// ck数据合并
		return externalSearchDao.NewExternalExtraCk().MergeData(c, enterpriseId)
	}
	return nil
}

// TruncateCdpTagTable  清空CDP表
func (s *ExternalSearch) TruncateCdpTagTable(c app.Ctx, enterpriseId int64) error {
	return externalSearchDao.NewCdpExternalTagCk().TruncateTable(c, enterpriseId)
}

// RelationshipDataMerge  数据合并
func (s *ExternalSearch) RelationshipDataMerge(c app.Ctx, enterpriseId int64, force bool) error {
	unixTime := time.Now().Unix()
	relationshipUpdateTime := fmt.Sprintf("customer:relationship:update-time-%d", enterpriseId)
	minuteTimesStr := s.Redis().Get(c, relationshipUpdateTime).Val()
	// 合并数据间隔一般需要大于5分钟
	if force || minuteTimesStr == "" {
		s.Redis().Set(c, relationshipUpdateTime, unixTime, time.Minute*10)
		// ck数据合并
		return externalSearchDao.NewExternalRelationshipCk().MergeData(c, enterpriseId)
	}
	return nil
}

// CdpTagDataMerge cdp_external_tag数据合并
func (s *ExternalSearch) CdpTagDataMerge(c app.Ctx, enterpriseId int64, force bool) error {
	unixTime := time.Now().Unix()
	relationshipUpdateTime := fmt.Sprintf("customer:cdpTag:update-time-%d", enterpriseId)
	minuteTimesStr := s.Redis().Get(c, relationshipUpdateTime).Val()
	// 合并数据间隔一般需要大于5分钟
	if force || minuteTimesStr == "" {
		s.Redis().Set(c, relationshipUpdateTime, unixTime, time.Minute*6)
		// ck数据合并
		return externalSearchDao.NewCdpExternalTagCk().MergeData(c, enterpriseId)
	}
	return nil
}

func (s *ExternalSearch) SetTaskHaveExtraData(c ctx.Context, enterpriseId int64) (total int, err error) {
	size := 3000
	externalExtraCk := externalSearchDao.NewExternalExtraCk()
	taskOrm := microScrmDao.NewPlanTaskConditionDataDao()
	// 处理任务数据, 这里任务分表后 要查询 当天新增了哪些子任务，循环查询拥有任务情况
	planList, err := microScrmDao.NewPlanTaskDao().
		GetMany(c, enterpriseId, "task_start_time > ? ", time.Now().Unix()-86400)
	if err != nil {
		return
	}
	if len(planList) == 0 {
		fmt.Print("【任务】顾客列表附加信息更新，未查询到新的任务")
	}
	var planIds []int64
	for _, planInfo := range planList {
		planIds = append(planIds, planInfo.PlanID)
	}
	planIds = helper.ArrayUniq(planIds)
	for _, planId := range planIds {
		var lastId int64
		for {
			// 任务的发送后created_at就不会再变化了，只取created_at为一天内的数据即可
			taskList, err := taskOrm.GetTaskByDays(c, enterpriseId, planId, lastId, size, 1, microScrmModel.TaskExist)
			if err != nil {
				return 0, err
			}
			if len(taskList) == 0 {
				break
			}
			// 任务数据
			var uniqueUnionIdMap = make(map[string]externalSearchDao.ExternalExtraData)
			for _, info := range taskList {
				if info.UnionId == "" {
					// 过滤unionid为空的顾客
					continue
				}
				if existAt, ok := uniqueUnionIdMap[info.UnionId]; ok {
					if info.CreatedAt > existAt.TaskCreatedAt {
						uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraData{
							UserUnionid:   info.UnionId,
							TaskCreatedAt: info.CreatedAt,
						}
					}
				} else {
					uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraData{
						UserUnionid:   info.UnionId,
						TaskCreatedAt: info.CreatedAt,
					}
				}
			}
			// 更新任务相关数据
			err = externalExtraCk.SaveTaskData(c, enterpriseId, uniqueUnionIdMap, microScrmModel.TaskExist)
			if err != nil {
				return 0, err
			}
			s.saveTaskExtension(c, enterpriseId, uniqueUnionIdMap)
			total += len(uniqueUnionIdMap)
			// 获取最后一条记录的id
			lastId = taskList[len(taskList)-1].ID
		}
	}
	return
}

// BatchSaveActiveData     更新顾客活跃沟通信息
func (s *ExternalSearch) BatchSaveActiveData(c app.Ctx, req openReq.BatchUpdateActiveReq) error {
	relationships, err := s.externalRelationshipCkDao.BatchSaveActiveData(c, req)
	if err != nil {
		return err
	}
	s.saveRelationshipExtension(c, req, relationships)
	// 数据合并
	_ = s.RelationshipDataMerge(c, req.EnterpriseId, false)
	return nil
}

// GetSearchOption  获取检索条件
func (s *ExternalSearch) GetSearchOption(c app.Ctx, request customerSearchReq.ReqChoice) (customerSearchReq.ChoiceOptionList, error) {
	choiceOptionList := customerSearchReq.ChoiceOptionList{}
	if len(request.ThirdIds) > 0 {
		// 员工筛选
		// 查询店长管理的门店下店员 GetMemberByThirdIds
		memberIds, err := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, request.EnterpriseID, request.ThirdIds, []int64{})
		if err != nil {
			return customerSearchReq.ChoiceOptionList{}, err
		}
		memberList, err := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND deleted_at = 0 AND id in ?", request.EnterpriseID, append(memberIds, request.MemberId))
		if err != nil {
			return customerSearchReq.ChoiceOptionList{}, err
		}
		for _, v := range memberList {
			choiceOptionList.MemberList = append(choiceOptionList.MemberList, customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name})
		}
	}
	// 会员卡等,卡等信息是coach的会员卡等信息(有两个筛选条件)
	var cardType []customerSearchReq.ChoiceOption
	cardType = append(cardType,
		customerSearchReq.ChoiceOption{
			Key:   "O",
			Value: "非会员",
		},
		customerSearchReq.ChoiceOption{
			Key:   "PROSPECT",
			Value: helper.GetCoachInfoByTime("品牌粉丝", "注册会员"),
		},
		customerSearchReq.ChoiceOption{
			Key:   "C",
			Value: helper.GetCoachInfoByTime("新晋会员", "已购会员"),
		},
		customerSearchReq.ChoiceOption{
			Key:   "SILVER",
			Value: "银卡会员",
		},
		customerSearchReq.ChoiceOption{
			Key:   "GOLD",
			Value: "金卡会员",
		},
		customerSearchReq.ChoiceOption{
			Key:   "DIAMOND",
			Value: "钻卡会员",
		},
	)
	choiceOptionList.CardType = cardType
	// 生日月份
	var birthMonth []customerSearchReq.ChoiceOption
	for i := 1; i <= 12; i++ {
		birthMonth = append(birthMonth,
			customerSearchReq.ChoiceOption{
				Key:   helper.EncodeToString(helper.Int64ToString(int64(i))),
				Value: helper.Int64ToString(int64(i)) + "月",
			})
	}
	choiceOptionList.BirthMonth = birthMonth
	// 顾问筛选
	var manageType []customerSearchReq.ChoiceList
	manageType = append(manageType,
		customerSearchReq.ChoiceList{
			Key:   1,
			Value: "专属",
		},
		customerSearchReq.ChoiceList{
			Key:   2,
			Value: "非专属",
		},
	)
	choiceOptionList.ManageType = manageType
	// 邀约筛选文案
	var reserveType []customerSearchReq.ChoiceList
	reserveType = append(reserveType,
		customerSearchReq.ChoiceList{
			Key:   1,
			Value: "邀约中",
		},
		customerSearchReq.ChoiceList{
			Key:   2,
			Value: "无邀约",
		},
		customerSearchReq.ChoiceList{
			Key:   3,
			Value: "邀约成功",
		})
	choiceOptionList.ReserveType = reserveType
	// 30天任务触达
	var taskType []customerSearchReq.ChoiceList
	taskType = append(taskType,
		customerSearchReq.ChoiceList{
			Key:   1,
			Value: "未触达",
		},
		customerSearchReq.ChoiceList{
			Key:   2,
			Value: "已触达",
		})
	choiceOptionList.TaskType = taskType
	// 近30天是否有任务
	var taskHave []customerSearchReq.ChoiceList
	taskHave = append(taskHave,
		customerSearchReq.ChoiceList{
			Key:   1,
			Value: "无任务",
		},
		customerSearchReq.ChoiceList{
			Key:   2,
			Value: "有任务",
		})
	choiceOptionList.TaskHave = taskHave
	choiceOptionList.TaskTypeValue = 7 //默认7天
	choiceOptionList.TaskHaveValue = 7 //默认7天
	configMap, _ := externalSearchDao.NewConfigDao().GetMapByKeys(c, request.EnterpriseID, []string{"task_type", "task_have"})
	if value, ok := configMap[customerSearchReq.TaskTypeKey]; ok {
		choiceOptionList.TaskTypeValue = helper.StringToInt(value)
	}
	if value, ok := configMap[customerSearchReq.TaskHaveKey]; ok {
		choiceOptionList.TaskHaveValue = helper.StringToInt(value)
	}
	// 优惠券筛选
	var couponType []customerSearchReq.ChoiceList
	couponType = append(couponType,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Has30DayCoupon,
			Value: "30天内临期券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.HasCoupon,
			Value: "待使用券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.BirthCoupon,
			Value: "生日礼券",
		},
	)
	if !helper.CheckCoachTime() {
		couponType = append(couponType,
			customerSearchReq.ChoiceList{
				Key:   customerSearchReq.RenewCoupon,
				Value: "续会礼券",
			},
			customerSearchReq.ChoiceList{
				Key:   customerSearchReq.UpgradeCoupon,
				Value: "升级礼券",
			})
	}
	choiceOptionList.CouponType = couponType

	// 历史消费时间(时间段选择，前端自己生成)

	// 近期未消费(时间段选择，前端自己生成)
	var notPurchasedDate []customerSearchReq.ChoiceList
	notPurchasedDate = append(notPurchasedDate,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt30D,
			Value: ">30D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt90D,
			Value: ">90D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt180D,
			Value: ">180D",
		})
	choiceOptionList.NotPurchasedDate = notPurchasedDate
	//历史购买金额LTV(金额区间选择，前端自己生成)

	//L12M总金额(金额区间选择，前端自己生成)

	//添加好友时间
	var relCreateTime []customerSearchReq.ChoiceList
	relCreateTime = append(relCreateTime,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeToday,
			Value: "今日",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeYesterday,
			Value: "昨日",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time2To7D,
			Value: "2-7D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time8To15D,
			Value: "8-15D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt15D,
			Value: ">15D",
		})
	choiceOptionList.RelCreateTime = relCreateTime
	//上次触达时间(和添加好友时间筛选项一样)
	var relInteractTime = relCreateTime
	choiceOptionList.RelInteractTime = relInteractTime

	//活跃沟通顾客
	var relActiveTime []customerSearchReq.ChoiceList
	relActiveTime = append(relActiveTime,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeLt7D,
			Value: "0-7D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time8To15D,
			Value: "8-15D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time16To30D,
			Value: "16-30D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt30D,
			Value: ">30D",
		})
	choiceOptionList.RelActiveTime = relActiveTime
	//单独沟通和勿扰标签
	chatTagList, err := tagsSvc.NewWorkTags().GetChatTagList(c, request.EnterpriseID)
	if err != nil {
		return customerSearchReq.ChoiceOptionList{}, err
	}
	chatTag := make([]customerSearchReq.ChoiceList, 0)
	for _, v := range chatTagList {
		temp := customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name}
		chatTag = append(chatTag, temp)
	}
	choiceOptionList.ChatTag = chatTag
	// 高意愿用户
	var highWilling []customerSearchReq.ChoiceOption
	highWilling = append(highWilling,
		customerSearchReq.ChoiceOption{
			Key:   "Y",
			Value: "是",
		},
		customerSearchReq.ChoiceOption{
			Key:   "N",
			Value: "否",
		},
	)
	choiceOptionList.HighWilling = highWilling

	// 价格敏感度
	var priceSensitivity []customerSearchReq.ChoiceOption
	priceSensitivity = append(priceSensitivity,
		customerSearchReq.ChoiceOption{
			Key:   "H",
			Value: "高",
		},
		customerSearchReq.ChoiceOption{
			Key:   "L",
			Value: "低",
		},
	)
	choiceOptionList.PriceSensitivity = priceSensitivity

	//商品偏好
	var categoryPreference []customerSearchReq.ChoiceOption
	categoryPreference = append(categoryPreference,
		customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.CategoryBD,
			Value: customerSearchReq.CategoryBD,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.CategoryXPJ,
			Value: customerSearchReq.CategoryXPJ,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.CategoryFS,
			Value: customerSearchReq.CategoryFS,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.CategoryXL,
			Value: customerSearchReq.CategoryXL,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.CategorySP,
			Value: customerSearchReq.CategorySP,
		})
	choiceOptionList.CategoryPreference = categoryPreference

	//顾客类型AGE
	var ageTag []customerSearchReq.ChoiceList
	ageTag = append(ageTag,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.AgeTagOne,
			Value: "Gen Z",
		}, customerSearchReq.ChoiceList{
			Key:   customerSearchReq.AgeTagTwo,
			Value: "YM",
		}, customerSearchReq.ChoiceList{
			Key:   customerSearchReq.AgeTagThree,
			Value: "OTHERS",
		})
	choiceOptionList.AgeTag = ageTag
	return choiceOptionList, err

}

func (s *ExternalSearch) GetExternalDetail(c app.Ctx, request customerSearchReq.ExternalDetailReq) (detail customerSearchResp.ExternalDetailResp, err error) {
	external, err := externalSearchDao.NewExternalExtraCk().GetOne(c, request.EnterpriseId, "openid=?", request.Openid)
	if err != nil {
		return customerSearchResp.ExternalDetailResp{}, err
	}
	if external.ID == 0 {
		//ck查不到去mysql查一次
		externalMysql, _ := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"openid": request.Openid}, request.EnterpriseId)
		if externalMysql.ID == 0 {
			externalMysql, err = FixExternal(c, request.EnterpriseId, request.Openid)
			if err != nil {
				return customerSearchResp.ExternalDetailResp{}, errors.New("找不到该用户")
			}
		}
		helper.CoverStruct(&detail, &externalMysql)
		detail.ID = externalMysql.ID
		detail.Birthday = ""
	} else {
		//合并结构体结果值
		helper.CoverStruct(&detail, &external)
		externalMysql, _ := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"openid": request.Openid}, request.EnterpriseId)
		detail.ID = external.ID
		detail.Unionid = external.UserUnionid
		detail.Birthday = ""
		if externalMysql.BirthYear != "" && externalMysql.BirthYear != "0" {
			year, _ := helper.DecodeString(externalMysql.BirthYear)
			month, _ := helper.DecodeString(externalMysql.BirthMonth)
			day, _ := helper.DecodeString(externalMysql.BirthDay)
			detail.Birthday = fmt.Sprintf("%s-%s-%s", year, month, day)
		}
		detail.CardLevel = externalMysql.CardLevel
		detail.UnionCode = externalMysql.UnionCode
	}

	var tmpCategoryPreference = make([]string, 0)
	for categoryName, categorySlice := range customerSearchReq.CdpCategoryMap {
		if len(helper.ArrayIntersect(categorySlice, detail.CategoryPreference)) > 0 {
			tmpCategoryPreference = append(tmpCategoryPreference, categoryName)
		}
	}
	detail.CategoryPreference = tmpCategoryPreference
	// 查询饼图
	cdpTagList, _ := externalSearchDao.NewCdpExternalTagCk().GetMany(c, request.EnterpriseId, "unionid = ? AND tag_type in ?",
		detail.Unionid, []string{externalSearchModel.CdpTagBuyCategoryPreference, externalSearchModel.CdpTagPurchasedMonth})
	if len(cdpTagList) > 0 {
		for _, tag := range cdpTagList {
			if tag.TagType == externalSearchModel.CdpTagBuyCategoryPreference {
				detail.CdpTagBuyCategoryPreference = s.parseCategoryPreferencePieChartData(tag.TagValue)
			}
			if tag.TagType == externalSearchModel.CdpTagPurchasedMonth {
				detail.CdpTagPurchasedMonth = s.parsePieChartData(tag.TagValue)
			}
		}
	}
	relationship, err := externalSearchDao.NewExternalRelationshipCk().GetOne(c, request.EnterpriseId, "unionid = ? AND member_id = ?", detail.Unionid, request.MemberID)
	if err != nil {
		return customerSearchResp.ExternalDetailResp{}, err
	}
	if relationship.ID > 0 {
		detail.ActiveTime = relationship.ActiveTime
		detail.InteractTime = relationship.InteractTime
	}
	mysqlRelation, _ := microScrmDao.NewUserRelation().GetOne(c, map[string]interface{}{
		"unionid":   detail.Unionid,
		"member_id": request.MemberID,
	})
	if mysqlRelation.ID > 0 && mysqlRelation.RelationStatus == 0 {
		helper.CoverStruct(&detail, &mysqlRelation)
		var existTagIds []int64
		tagList, err := tagsDao.NewWorkExternalUserTag().GetMany(c, map[string]interface{}{
			"enterprise_id": request.EnterpriseId,
			"relation_id":   mysqlRelation.ID,
		})
		for _, tag := range tagList {
			existTagIds = append(existTagIds, tag.ExternalTagID)
		}
		//单独沟通和勿扰标签
		chatTagList, err := tagsSvc.NewWorkTags().GetChatTagList(c, request.EnterpriseId)
		if err != nil {
			return customerSearchResp.ExternalDetailResp{}, err
		}
		chatTag := make([]customerSearchResp.ChatList, 0)
		for _, v := range chatTagList {
			temp := customerSearchResp.ChatList{Key: v.ID, Value: v.Name}
			if helper.InArray(v.ID, existTagIds) {
				temp.IsSelect = true
			}
			chatTag = append(chatTag, temp)
		}
		detail.ChatTag = chatTag
	}

	return detail, nil
}

func FixExternal(ctx ctx.Context, enterpriseId int64, openId string) (microScrmModel.WxUserCorpExternalModel, error) {
	external := microScrmModel.WxUserCorpExternalModel{}
	log.Info(ctx).Msgf("兜底修复外部联系人 - 开始%+v", openId)
	//mysql查不到则去企微查
	newClient, _ := workwechat.NewWorkWeChat(ctx, enterpriseId)
	externalDetail, err := newClient.GetWorkExternalUserInfo(openId)
	if err != nil {
		return external, err
	}
	corpID, err := enterprise.New(enterpriseMysql.New()).GetCorpID(ctx, enterpriseId)
	if err != nil {
		return external, err
	}
	log.Info(ctx).Msgf("兜底修复外部联系人 - 企微信息获取成功%+v", externalDetail)
	externalId := app.UID()
	external = microScrmModel.WxUserCorpExternalModel{
		BaseModel:    models.CreateBase(externalId),
		Appid:        corpID,
		Openid:       externalDetail.Base.ExternalUserID,
		Unionid:      externalDetail.Base.UnionID,
		Euid:         0,
		Nick:         externalDetail.Base.Name,
		Name:         externalDetail.Base.Name,
		Position:     externalDetail.Base.Position,
		Avatar:       externalDetail.Base.Avatar,
		CorpName:     externalDetail.Base.CorpName,
		CorpFullName: externalDetail.Base.CorpFullName,
		Type:         int32(externalDetail.Base.Type),
		Sex:          int32(externalDetail.Base.Gender),
	}
	err = microScrmDao.NewWxUserCorpExternal().SaveData(ctx, enterpriseId, external)
	if err != nil {
		return external, err
	}
	log.Info(ctx).Msgf("兜底修复外部联系人 - 顾客信息入库")
	if err != nil {
		return external, err
	}
	var memberUserIDs []string
	for _, user := range externalDetail.FollowUser {
		memberUserIDs = append(memberUserIDs, user.UserID)
	}
	memberList, err := membersDao.NewMembers().GetMany(ctx, "enterprise_id= ? AND userid in ? AND deleted_at=0", enterpriseId, memberUserIDs)
	if err != nil {
		return external, err
	}
	var addRelation []microScrmModel.WorkExternalUserRelationModel
	for _, relationV := range externalDetail.FollowUser {
		for _, member := range memberList {
			if member.Userid == relationV.UserID {
				addRelation = append(addRelation, microScrmModel.WorkExternalUserRelationModel{
					BaseModel:          models.CreateBase(app.UID()),
					EnterpriseID:       enterpriseId,
					WorkExternalUserID: externalId,
					MemberID:           member.ID,
					Remark:             relationV.Remark,
					Description:        relationV.Description,
					Createtime:         int(relationV.CreateTime),
					RemarkCorpName:     relationV.RemarkCorpName,
					RemarkMobiles:      strings.Join(relationV.RemarkMobiles, ","),
					AddWay:             int32(relationV.AddWay),
					State:              relationV.State,
					OperUserid:         relationV.OperUserID,
					RelationStatus:     0,
					UserID:             relationV.UserID,
					Unionid:            externalDetail.Base.UnionID,
				})
			}
		}
	}
	if len(addRelation) > 0 {
		_ = gormx.GetInstance().Model(microScrmModel.WorkExternalUserRelationModel{}).CreateInBatches(addRelation, len(addRelation))
		log.Info(ctx).Msgf("兜底修复外部联系人 - 关系修复结束，修复数量%+v", len(addRelation))
	}
	return external, err
}

func (s *ExternalSearch) parsePieChartData(str string) []customerSearchResp.PieChart {
	// 1. 将输入的字符串按 "|" 拆分成多个项
	items := strings.Split(str, "|")

	// 2. 用来存储结构体切片
	var pieCharts = make([]customerSearchResp.PieChart, 0)

	// 3. 用来计算总和
	var total int

	// 4. 解析每个项
	for _, item := range items {
		// 解析每个项，格式如 WHB:1 或 WLG:2
		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		// 计算总和
		total += value

		// 生成饼图数据结构
		pieCharts = append(pieCharts, customerSearchResp.PieChart{Name: name, Percent: float64(value)})
	}

	// 5. 计算前面部分的百分比
	var accumulatedPercent float64
	for i := 0; i < len(pieCharts)-1; i++ {
		pieCharts[i].Percent = helper.MathRounded((pieCharts[i].Percent/float64(total))*100, 2)
		accumulatedPercent += pieCharts[i].Percent
	}

	// 6. 计算最后一个部分的百分比，使得总和为100
	if len(pieCharts) > 0 {
		pieCharts[len(pieCharts)-1].Percent = helper.MathRounded(100-accumulatedPercent, 2)
	}

	return pieCharts
}

// parsePieChartData 用于根据输入的子分类数据计算各主分类的百分比
func (s *ExternalSearch) parseCategoryPreferencePieChartData(str string) []customerSearchResp.PieChart {
	// 1. 将输入的字符串按 "|" 拆分成多个项
	items := strings.Split(str, "|")

	// 2. 用来存储主分类和购买次数
	categoryCounts := make(map[string]int)

	// 3. 计算总和
	var total int

	// 4. 解析每个项
	for _, item := range items {
		// 解析每个项，格式如 WHB:1 或 WLG:2
		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		// 计算总和
		total += value

		// 5. 根据子分类标识查找主分类并累加购买次数
		for mainCategory, subCategories := range customerSearchReq.CdpCategoryMap {
			for _, sub := range subCategories {
				if name == sub {
					categoryCounts[mainCategory] += value
				}
			}
		}
	}

	// 6. 计算主分类的百分比并填充 PieChart
	var pieCharts []customerSearchResp.PieChart

	// 7. 计算百分比并生成饼图数据
	for _, category := range []string{
		customerSearchReq.CategoryBD,
		customerSearchReq.CategoryXPJ,
		customerSearchReq.CategoryFS,
		customerSearchReq.CategoryXL,
		customerSearchReq.CategorySP} {
		count, exists := categoryCounts[category]
		if exists {
			percent := helper.MathRounded((float64(count)/float64(total))*100, 2)
			pieCharts = append(pieCharts, customerSearchResp.PieChart{Name: category, Percent: percent})
		}
	}

	// 8. 计算百分比，避免浮动，总和为100
	var calculatedPercent float64
	for i := 0; i < len(pieCharts)-1; i++ {
		calculatedPercent += pieCharts[i].Percent
	}

	// 调整最后一个百分比，以确保总和为100
	if len(pieCharts) > 0 {
		pieCharts[len(pieCharts)-1].Percent = helper.MathRounded(100-calculatedPercent, 2)
	}

	return pieCharts
}

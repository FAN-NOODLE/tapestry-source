package customerSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	openReq "micro-scrm/internal/app/http/request/open"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	workTagSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

type ExternalSearchSw struct {
	service.Base
	externalRelationshipCkDao *microScrmDao.UserRelation
	externalExtraCkDao        *externalSearchDao.ExternalExtraCkSw
}

func NewExternalSearchSw() *ExternalSearchSw {
	return &ExternalSearchSw{
		externalRelationshipCkDao: microScrmDao.NewUserRelation(),
		externalExtraCkDao:        externalSearchDao.NewExternalExtraCkSw(),
	}
}

// GetSearchOption  获取检索条件
func (s *ExternalSearchSw) GetSearchOption(c app.Ctx, request customerSearchReq.ReqChoice) (customerSearchReq.ChoiceOptionSwList, error) {
	choiceOptionList := customerSearchReq.ChoiceOptionSwList{}
	if len(request.ThirdIds) > 0 {
		// 员工筛选
		// 查询店长管理的门店下店员 GetMemberByThirdIds
		memberIds, err := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, request.EnterpriseID, request.ThirdIds, []int64{})
		if err != nil {
			return customerSearchReq.ChoiceOptionSwList{}, err
		}
		memberList, err := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND deleted_at = 0 AND id in ?", request.EnterpriseID, append(memberIds, request.MemberId))
		if err != nil {
			return customerSearchReq.ChoiceOptionSwList{}, err
		}
		for _, v := range memberList {
			choiceOptionList.MemberList = append(choiceOptionList.MemberList, customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name})
		}
	}
	// 顾客类型
	var customerType []customerSearchReq.ChoiceList
	customerType = append(customerType,
		customerSearchReq.ChoiceList{
			Key:   1,
			Value: "新客",
		},
		customerSearchReq.ChoiceList{
			Key:   2,
			Value: "老客",
		},
		customerSearchReq.ChoiceList{
			Key:   3,
			Value: "VIP客",
		},
		customerSearchReq.ChoiceList{
			Key:   4,
			Value: "沉睡客",
		},
	)
	choiceOptionList.CustomerType = customerType
	// 会员卡等,卡等信息是
	var cardType []customerSearchReq.ChoiceOption
	cardType = append(cardType,
		customerSearchReq.ChoiceOption{
			Key:   "B",
			Value: "钻石「履」行家",
		},
		customerSearchReq.ChoiceOption{
			Key:   "C",
			Value: "金牌「履」行家",
		},
		customerSearchReq.ChoiceOption{
			Key:   "W",
			Value: "新星「履」行家",
		},
		customerSearchReq.ChoiceOption{
			Key:   "O",
			Value: "非会员",
		},
	)

	choiceOptionList.CardType = cardType
	// 生日月份
	var birthMonth []customerSearchReq.ChoiceOption
	for i := 1; i <= 12; i++ {
		birthMonth = append(birthMonth,
			customerSearchReq.ChoiceOption{
				Key:   helper.EncodeToString(helper.RuleThree(i < 10, "0", "") + helper.Int64ToString(int64(i))),
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
	// 优惠券筛选
	var couponType []customerSearchReq.ChoiceList
	couponType = append(couponType,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwHas30DayCoupon,
			Value: "30天内临期券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwHasCoupon,
			Value: "待使用券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwBirthCoupon,
			Value: "生日礼券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwNewMemberCoupon,
			Value: "新人礼券",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwRebateCoupon,
			Value: "回购礼券",
		}, customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwPromotionalCoupon,
			Value: "活动礼券",
		})
	choiceOptionList.CouponType = couponType
	//添加好友时间
	var relCreateTime []customerSearchReq.ChoiceList
	relCreateTime = append(relCreateTime,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeLt1D,
			Value: "0D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time1To30D,
			Value: "1-30D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time31To60D,
			Value: "31-60D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time61To90D,
			Value: "61-90D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt90D,
			Value: ">90D",
		})
	choiceOptionList.RelCreateTime = relCreateTime
	//上次触达时间
	var relInteractTime []customerSearchReq.ChoiceList
	relInteractTime = append(relInteractTime,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeLt30D,
			Value: "0-30D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time31To60D,
			Value: "31-60D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.Time61To90D,
			Value: "61-90D",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.TimeGt90D,
			Value: ">90D",
		})
	choiceOptionList.RelInteractTime = relInteractTime
	//商品偏好
	var categoryPreference []customerSearchReq.ChoiceOption
	categoryPreference = append(categoryPreference,
		customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryDX,
			Value: customerSearchReq.SwCategoryDX,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryCX,
			Value: customerSearchReq.SwCategoryCX,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryLX,
			Value: customerSearchReq.SwCategoryLX,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryGGX,
			Value: customerSearchReq.SwCategoryGGX,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryYDX,
			Value: customerSearchReq.SwCategoryYDX,
		}, customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwCategoryPDX,
			Value: customerSearchReq.SwCategoryPDX,
		})
	choiceOptionList.CategoryPreference = categoryPreference

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

	// 年消费金额
	//改为开始和结束时间（前端自己处理)
	// 价格偏好
	var priceTag []customerSearchReq.ChoiceOption
	priceTag = append(priceTag,
		customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwPriceTag1,
			Value: customerSearchReq.SwPriceTag1,
		},
		customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwPriceTag2,
			Value: customerSearchReq.SwPriceTag2,
		},
		customerSearchReq.ChoiceOption{
			Key:   customerSearchReq.SwPriceTag3,
			Value: customerSearchReq.SwPriceTag3,
		})
	choiceOptionList.PriceTag = priceTag

	// 价格偏好
	var ageTag []customerSearchReq.ChoiceOption
	ageTag = append(ageTag,
		customerSearchReq.ChoiceOption{
			Key:   "GENZ",
			Value: "GENZ",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Younger millennials",
			Value: "Younger millennials",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Old millennia",
			Value: "Old millennia",
		},
		customerSearchReq.ChoiceOption{
			Key:   "GENX",
			Value: "GENX",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Boomer+",
			Value: "Boomer+",
		},
		customerSearchReq.ChoiceOption{
			Key:   "unknow",
			Value: "unknow",
		})
	choiceOptionList.AgeTag = ageTag
	// SW at消费金额
	var atPurchasedYear []customerSearchReq.ChoiceOption
	atPurchasedYear = append(atPurchasedYear,
		customerSearchReq.ChoiceOption{
			Key:   "2024",
			Value: "2024",
		},
		customerSearchReq.ChoiceOption{
			Key:   "2025",
			Value: "2025",
		})
	choiceOptionList.AtPurchasedYear = atPurchasedYear

	var atPurchasedQuarter []customerSearchReq.ChoiceOption
	atPurchasedQuarter = append(atPurchasedQuarter,
		customerSearchReq.ChoiceOption{
			Key:   "Q1",
			Value: "Q1",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Q2",
			Value: "Q2",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Q3",
			Value: "Q3",
		},
		customerSearchReq.ChoiceOption{
			Key:   "Q4",
			Value: "Q4",
		})
	choiceOptionList.AtPurchasedQuarter = atPurchasedQuarter

	var atPurchasedAmount []customerSearchReq.ChoiceList
	atPurchasedAmount = append(atPurchasedAmount,
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwAtPurchasedAmountGt8000,
			Value: ">8k",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwAtPurchasedAmount6000To8000,
			Value: "6-8K(含6K)",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwAtPurchasedAmount4000To6000,
			Value: "4-6K(含4K)",
		},
		customerSearchReq.ChoiceList{
			Key:   customerSearchReq.SwAtPurchasedAmountLt4000,
			Value: "<4K",
		},
	)
	choiceOptionList.AtPurchasedAmount = atPurchasedAmount

	// 客户企微标签
	var tagList []customerSearchReq.ChoiceList
	tagLists, _ := workTagSvc.NewWorkTags().GetChoiceTagList(c, request.EnterpriseID)
	for _, v := range tagLists {
		temp := customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name}
		tagList = append(tagList, temp)
	}
	// 印象标签
	var impressionTagList []customerSearchReq.ChoiceList
	var chatTagList []customerSearchReq.ChoiceList
	var impressionTagReq tagsReq.ReqList
	impressionTagReq.EnterpriseID = request.EnterpriseID
	impressionTagLists, _ := workTagSvc.NewImpression().GetList(c, impressionTagReq)

	for _, v := range impressionTagLists {
		if helper.InArray(v.Name, []string{"勿扰", "愿意被触达"}) {
			tempChat := customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name}
			chatTagList = append(chatTagList, tempChat)
		} else {
			temp := customerSearchReq.ChoiceList{Key: v.ID, Value: v.Name}
			impressionTagList = append(impressionTagList, temp)
		}
	}

	choiceOptionList.TagList = tagList
	choiceOptionList.ImpressionTagList = impressionTagList
	choiceOptionList.ChatTag = chatTagList
	return choiceOptionList, nil
}

func (s *ExternalSearchSw) GetExternalSearchListCk(c app.Ctx, request customerSearchReq.ExternalSearchListBaseSwReq) ([]customerSearchResp.ExternalUserListSw, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserListSw, 0)
	// 判断当前员工类型 店员或店长
	memberIds, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	if err != nil {
		return list, 0, err
	}
	request.MemberIds = memberIds
	if request.IsRelationship {
		//顾客关系列表
		list, total, err = s.externalExtraCkDao.RelationShipListCk(c, request)
	} else {
		//顾客列表
		list, total, err = s.externalExtraCkDao.ExternalSearchListCk(c, request)
	}
	if err != nil || total == 0 {
		return list, 0, err
	}
	var storeCodes []string
	for i, row := range list {
		if row.LastPurchaseStoreCode != "" {
			storeCodes = append(storeCodes, row.LastPurchaseStoreCode)
		}
		// 生日解密
		list[i].Birthday, _ = helper.DecodeString(row.Birthday)
	}
	if len(storeCodes) > 0 {
		storeDepartmentInfoMap, _ := s.GetStoreDepartmentInfoMap(c, request.EnterpriseID, storeCodes)
		for i, row := range list {
			if row.LastPurchaseStoreCode != "" {
				lastPurchaseStoreCode := strings.ToUpper(row.LastPurchaseStoreCode)
				if _, ok := storeDepartmentInfoMap[lastPurchaseStoreCode]; ok {
					list[i].LastPurchaseStoreName = storeDepartmentInfoMap[lastPurchaseStoreCode].StoreName
				}
			}
		}
	}

	return list, total, err
}

// GetStoreDepartmentInfoMap 查询门店信息-sw
func (s *ExternalSearchSw) GetStoreDepartmentInfoMap(c app.Ctx, enterpriseID int64, storeCodes []string) (map[string]microWorkModel.StoreDepartmentModel, error) {

	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetStoreInfoByStoreCode(c, enterpriseID, storeCodes)
	var storeDepartmentInfoMap = make(map[string]microWorkModel.StoreDepartmentModel, 0)
	if err != nil {
		return storeDepartmentInfoMap, err
	}
	for _, item := range storeDepartmentInfo {
		storeDepartmentInfoMap[strings.ToUpper(item.ThirdStoreCode)] = item
	}

	return storeDepartmentInfoMap, nil
}

// GetMembersByManageType 根据导购或店长类型获取店员
func (s *ExternalSearchSw) GetMembersByManageType(c app.Ctx, enterpriseID, memberId int64, thirdIds []int, memberIds []int64) ([]int64, error) {
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

func (s *ExternalSearchSw) GetExternalSearchCk(c app.Ctx, request customerSearchReq.ExternalSearchListBaseSwReq, pluckName string) ([]interface{}, error) {
	// 判断当前员工类型 店员或店长
	memberIds, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	if err != nil {
		return nil, err
	}
	request.MemberIds = memberIds
	return s.externalExtraCkDao.ExternalSearchPluckIds(c, request.SearchListSw, pluckName)
}

// SetAllUserExtraData   顾客列表附加检索信息用户补全
func (s *ExternalSearchSw) SetAllUserExtraData(c app.Ctx, enterpriseId int64) error {
	return s.externalExtraCkDao.SetAllUserExtraData(c, enterpriseId)
}

// SetReserveExtraData  顾客列表附加检索信息如表更新
func (s *ExternalSearchSw) SetReserveExtraData(c app.Ctx, enterpriseId int64) (total int, err error) {
	size := 3000
	var lastId int64
	externalExtraCk := externalSearchDao.NewExternalExtraCkSw()
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
		var uniqueMap = make(map[string]externalSearchDao.ExternalExtraDataSw)
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
				uniqueMap[info.UserUnionid] = externalSearchDao.ExternalExtraDataSw{
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
		total += len(uniqueMap)
		// 获取最后一条记录的id
		lastId = reserveList[len(reserveList)-1].ID
	}
	return
}

// SetTaskExtraData  顾客列表附加检索信息入表更新
func (s *ExternalSearchSw) SetTaskExtraData(c app.Ctx, enterpriseId int64, days int) (total int, err error) {
	size := 3000
	externalExtraCk := externalSearchDao.NewExternalExtraCkSw()
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
			var uniqueUnionIdMap = make(map[string]externalSearchDao.ExternalExtraDataSw)
			for _, info := range taskList {
				if info.UnionId == "" {
					// 过滤unionid为空的顾客
					continue
				}
				if existAt, ok := uniqueUnionIdMap[info.UnionId]; ok {
					if info.FinishedAt > existAt.TaskFinishAt {
						uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraDataSw{
							UserUnionid:  info.UnionId,
							TaskFinishAt: info.FinishedAt,
						}
					}
				} else {
					uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraDataSw{
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
			total += len(uniqueUnionIdMap)
			// 获取最后一条记录的id
			lastId = taskList[len(taskList)-1].ID
		}
	}

	return
}

func (s *ExternalSearchSw) SetTaskHaveExtraData(c ctx.Context, enterpriseId int64) (total int, err error) {
	size := 3000
	externalExtraCk := externalSearchDao.NewExternalExtraCkSw()
	taskOrm := microScrmDao.NewPlanTaskConditionDataDao()
	// 处理任务数据, 这里任务分表后 要查询 当天新增了哪些子任务，循环查询拥有任务情况
	planList, err := microScrmDao.NewPlanTaskDao().
		GetMany(c, enterpriseId, "task_start_time > ? AND task_category = 0 ", time.Now().Unix()-86400)
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
			var uniqueUnionIdMap = make(map[string]externalSearchDao.ExternalExtraDataSw)
			for _, info := range taskList {
				if info.UnionId == "" {
					// 过滤unionid为空的顾客
					continue
				}
				if existAt, ok := uniqueUnionIdMap[info.UnionId]; ok {
					if info.CreatedAt > existAt.TaskCreatedAt {
						uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraDataSw{
							UserUnionid:   info.UnionId,
							TaskCreatedAt: info.CreatedAt,
						}
					}
				} else {
					uniqueUnionIdMap[info.UnionId] = externalSearchDao.ExternalExtraDataSw{
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
			total += len(uniqueUnionIdMap)
			// 获取最后一条记录的id
			lastId = taskList[len(taskList)-1].ID
		}
	}
	return
}

func (s *ExternalSearchSw) SetInteractTimeExtraData(c ctx.Context, enterpriseId int64) (total int, err error) {
	externalExtraCk := externalSearchDao.NewExternalExtraCkSw()
	var uniqueUnionIdMap = make(map[string]externalSearchDao.ExternalExtraDataSw, 0)
	externalList, _ := microScrmDao.NewWxUserCorpExternal().GetMany(c, enterpriseId, "customer_type_updated_at > ?", time.Now().Unix()-86400)
	for _, external := range externalList {
		uniqueUnionIdMap[external.Unionid] = externalSearchDao.ExternalExtraDataSw{
			UserUnionid:      external.Unionid,
			LastPurchaseDate: external.LastPurchaseDate,
		}
	}
	err = externalExtraCk.SaveTaskData(c, enterpriseId, uniqueUnionIdMap, microScrmModel.InteractTime)
	total = len(uniqueUnionIdMap)
	return
}

// ExtraDataMerge 数据合并
func (s *ExternalSearchSw) ExtraDataMerge(c app.Ctx, enterpriseId int64, force bool) error {
	unixTime := time.Now().Unix()
	couponUpdateTime := fmt.Sprintf("customer:extra:update-time-%d", enterpriseId)
	minuteTimesStr := s.Redis().Get(c, couponUpdateTime).Val()
	// 合并数据间隔一般需要大于5分钟
	if force || minuteTimesStr == "" {
		err := s.Redis().Set(c, couponUpdateTime, unixTime, time.Minute*6).Err()
		if err != nil {
			return err
		}
		// ck数据合并
		return externalSearchDao.NewExternalExtraCk().MergeData(c, enterpriseId)
	}
	return nil
}

// BatchSaveCouponExtraData    顾客列表附加检索信息优惠券相关更新
func (s *ExternalSearchSw) BatchSaveCouponExtraData(c app.Ctx, req openReq.BatchUpdateCouponSwReq) error {
	var dataMap = make(map[string]openReq.UpdateCouponSwParams)
	for _, row := range req.DataList {
		if row.UnionId == "" {
			continue
		}
		dataMap[row.UnionId] = row
	}
	if len(dataMap) == 0 {
		return nil
	}
	err := externalSearchDao.NewExternalExtraCkSw().BatchSaveCouponData(c, req.EnterpriseId, dataMap)
	if err != nil {
		return err
	}
	// 数据合并
	_ = s.ExtraDataMerge(c, req.EnterpriseId, false)
	return nil
}

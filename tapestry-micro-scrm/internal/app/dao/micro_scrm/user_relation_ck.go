package microScrmDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/http/request"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	customerReq "micro-scrm/internal/app/http/request/customer"
	memberReq "micro-scrm/internal/app/http/request/member"
	activityQRCodeResp "micro-scrm/internal/app/http/response/activity_qrcode"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strings"
	"time"
)

const ValidDays = 30        // 检索的有效期
const TaskHaveValidDays = 5 // 检索是否有任务的有效期
// GetUserListCk 根据条件筛选导购和外部联系人列表
func (d *UserRelation) GetUserListCk(c ctx.Context, request customerReq.ExternalUserListsReq) (int64, []workExternalModel.ExternalUserList, error) {
	var (
		total     int64
		err       error
		selectStr string
	)
	result := make([]workExternalModel.ExternalUserList, 0)
	selectStr = "work_external_user_relation.id as id,"
	if request.IsGroup == true {
		selectStr = "max(work_external_user_relation.id) as id,"
		if request.IsSelf == 1 { // 顾客列表筛选之后，点击只看好友按钮
			selectStr = "work_external_user_relation.work_external_user_id as id,"
		}

	}
	groupBy := "work_external_user_relation.id"
	orderBy := "id"
	if request.IsGroup == true {
		orderBy = "work_external_user_id"
		groupBy = "work_external_user_relation.work_external_user_id"
	}

	query := d.searchBuilderQuery(c, request.EnterpriseID, request.SearchList, request.Userids, request.MemberId)
	if request.MemberType == membersModel.MembersTypeSM { // 店长的情况
		if request.IsSelf == 1 { // 仅看我的条件
			if len(request.MemberList) == 0 { // 没有筛选导购的情况
				query.Where("work_external_user_relation.member_id = ?", request.MemberId)
			} else {
				query.Where("work_external_user_relation.member_id in (?)", request.MemberList)
			}
		} else {
			if len(request.MemberList) > 0 { // 有筛选导购的情况
				query.Where("work_external_user_relation.member_id in (?)", request.MemberList)
			} else {
				query.Where("work_external_user_relation.member_id in (?)", request.MemberLists)
			}
		}
	} else {
		query.Where("work_external_user_relation.member_id = ?", request.MemberId)
	}

	if request.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.nick LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
	}
	if request.IsGroup == true {
		if request.IsSelf != 1 {
			newQuery := query.Session(&gorm.Session{NewDB: true}).WithContext(c).Distinct("work_external_user_relation.work_external_user_id")
			newQuery.Count(&total)
		}
	} else {
		newQuery := query.Session(&gorm.Session{NewDB: true}).WithContext(c).Distinct("work_external_user_relation.id")
		newQuery.Count(&total)
	}
	var relationIds []int64
	query.Select(strings.TrimRight(selectStr, ","))
	if request.IsSelf != 1 {
		query.Offset((request.Page - 1) * request.Size).Limit(request.Size)
	}
	query.Order(fmt.Sprintf("%s %s", orderBy, "asc")).
		Group(groupBy).Pluck("id", &relationIds)
	err = query.Error
	if len(relationIds) == 0 {
		return total, result, err
	}
	if request.IsSelf == 1 {
		userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
		newQuery := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
			Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").Where(map[string]interface{}{
			"work_external_user_id": relationIds,
			"member_id":             request.MemberId,
			"relation_status":       0,
		})
		newQuery.Count(&total)
	}
	result, err = d.GetUserListByRelationIdCk(c, request.EnterpriseID, relationIds, request.IsSelf, request.MemberId)
	if err != nil {
		return total, result, err
	}
	result, err = d.setExtraData(c, request.EnterpriseID, result)
	if err != nil {
		return total, result, err
	}
	// 对结果进行排序
	if request.IsGroup == true {
		sort.Slice(result, func(i, j int) bool {
			return result[i].WorkExternalUserID < result[j].WorkExternalUserID
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Id < result[j].Id
		})
	}
	return total, result, err
}

func (d *UserRelation) setExtraData(c app.Ctx, enterpriseId int64, result []workExternalModel.ExternalUserList) ([]workExternalModel.ExternalUserList, error) {
	validUnix := int(time.Now().Unix() - ValidDays*86400)
	var unionidArr []string
	for _, data := range result {
		unionidArr = append(unionidArr, data.Unionid)
	}
	var extraList []externalSearchModel.ExternalExtraModel
	err := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", enterpriseId)).
		Where("user_unionid in ?", unionidArr).
		Find(&extraList).Error
	if err != nil {
		return nil, err
	}
	for i, data := range result {
		for _, extra := range extraList {
			if data.Unionid == extra.UserUnionid {
				if extra.SuccessReserveAt > validUnix {
					result[i].SuccessReserveAt = extra.SuccessReserveAt
				}
				if extra.ValidReserveAt > validUnix {
					result[i].ValidReserveAt = extra.ValidReserveAt
				}
				if extra.TaskFinishAt > validUnix {
					result[i].TaskFinishAt = extra.TaskFinishAt
				}
				result[i].CouponStatus = extra.CouponStatus
				result[i].CouponBy30daysStatus = extra.CouponBy30daysStatus
			}
		}
	}
	return result, err
}

// SwCustomersCk 顾客榜单查询
func (d *UserRelation) SwCustomersCk(c app.Ctx, request memberReq.SwCustomersListReq, membersID []int64) ([]memberReq.SwCustomersListReqRes, int64, error) {

	var (
		total  int64
		result []memberReq.SwCustomersListReqRes
		err    error
	)

	offset := (request.Page - 1) * request.Size
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseId)

	selectStr := "count(IF(c.customer_type = 1,1,null)) as new_customers," +
		"count(IF(c.customer_type = 2,1,null)) as old_customers," +
		"count(IF(c.customer_type = 3,1,null)) as sleep_customers," +
		"count(IF(c.vip_type = 1,1,null)) as vip_customers," +
		"work_external_user_relation.member_id"

	query := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
		Joins(" LEFT JOIN "+userCorpExternalTable+" c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   request.EnterpriseId,
			"work_external_user_relation.deleted_at":      0,
			"c.deleted_at":                                0,
			"work_external_user_relation.relation_status": 0,
		}).
		Where("c.unionid <> '' ").
		Where("c.customer_type_updated_at BETWEEN ? AND ?", startTime, endTime).
		Group("work_external_user_relation.member_id")

	if len(membersID) > 0 {
		query.Where("work_external_user_relation.member_id in ?", membersID)
		query.Count(&total)
	}

	query.Select(selectStr).
		Order("new_customers desc").
		Order("old_customers desc").
		Order("vip_customers desc").
		Order("work_external_user_relation.member_id asc")

	err = query.Offset(offset).Limit(request.Size).Find(&result).Error

	return result, total, err
}

// ExternalStatistics 顾客榜单查询
func (d *UserRelation) ExternalStatistics(c app.Ctx, enterpriseId int64) ([]memberReq.ExternalStatisticsRes, error) {

	var result []memberReq.ExternalStatisticsRes

	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)

	selectStr := "count(IF(c.customer_type = 1,1,null)) as new_customers," +
		"count(IF(c.customer_type = 2,1,null)) as old_customers," +
		"count(IF(c.customer_type = 3,1,null)) as sleep_customers," +
		"count(IF(c.vip_type = 1,1,null)) as vip_customers," +
		"work_external_user_relation.member_id"

	err := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
		Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   enterpriseId,
			"work_external_user_relation.deleted_at":      0,
			"c.deleted_at":                                0,
			"work_external_user_relation.relation_status": 0,
		}).
		Where("c.unionid <> '' ").
		Group("work_external_user_relation.member_id").
		Select(selectStr).Find(&result).Error

	return result, err
}

func (d *UserRelation) GetRemarkByExternalUser(c app.Ctx, enterpriseId int64, workExternalUserIDs []int64) (map[int64][]string, error) {
	result := make(map[int64][]string)

	var userRemark []struct {
		WorkExternalUserId int64  `json:"work_external_user_id"`
		Remark             string `json:"remark"`
	}

	// 取出关系中的备注
	err := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
		Select("work_external_user_id,remark").
		Where(map[string]any{
			"enterprise_id":         enterpriseId,
			"work_external_user_id": workExternalUserIDs,
		}).Scan(&userRemark).Error

	for _, user := range userRemark {
		result[user.WorkExternalUserId] = append(result[user.WorkExternalUserId], user.Remark)
	}

	return result, err
}

// 组装筛选条件的sqlbuilder
func (d *UserRelation) searchBuilderQuery(c ctx.Context, enterpriseId int64, request request.SearchList, userIds []string, memberId int64) *gorm.DB {
	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)
	query := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
		Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id": enterpriseId,
			"work_external_user_relation.deleted_at":    0,
		})
	if request.IsUnion < workExternalModel.IS_UNION_FALSE {
		query.Where("work_external_user_relation.relation_status=?", 0)
	} else {
		query.Where("work_external_user_relation.relation_status!=?", 0)
	}
	query.Where("c.unionid <> '' ")
	// 1【专属】指ks_sa所属顾问是自己的顾客
	// 2【非专属】指所属顾问不是自己的顾客
	if request.ManageType == workExternalModel.ManageIsSA {
		query.Where("c.ba_code in (?)", userIds)
	} else if request.ManageType == workExternalModel.ManageNotSA {
		query.Not("c.ba_code  in (?)", userIds)
	}
	// 企微标签
	if len(request.TagList) > 0 || len(request.ChatTag) > 0 {
		query.Joins(" LEFT JOIN micro_scrm_ck.work_external_user_tag b ON work_external_user_relation.id = b.relation_id").Where("b.deleted_at = 0")
	}
	if len(request.TagList) > 0 {
		query.Where(" b.external_tag_id in (?)  ", request.TagList)
	}
	if len(request.ChatTag) > 0 { // 顾客沟通意愿标签，之查看自己打的标签数据
		query.Where(" b.external_tag_id in (?)  ", request.ChatTag)
		if memberId > 0 {
			query.Where("work_external_user_relation.member_id = ?", memberId)
		}
	}
	// 印象标签
	if len(request.ImpressionTagList) > 0 {
		query.Joins(" LEFT JOIN micro_scrm_ck.impression_tag_relation d ON c.id = d.external_user_id").
			Where("d.deleted_at = 0").Where(" d.tag_id in (?) and d.user_id in (?) ", request.ImpressionTagList, userIds)
	}
	// 邀约和任务的检索
	extraTableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", enterpriseId)
	if request.ReserveType > 0 || request.TaskType > 0 || request.CouponType > 0 || request.TaskHave > 0 {
		validTime := time.Now().Unix() - ValidDays*86400                 // 任务和邀约查询的数据有效时间都是30天内的
		taskHaveValidTime := time.Now().Unix() - TaskHaveValidDays*86400 // 任务和邀约查询的数据有效时间都是30天内的
		query.Joins(" ANY LEFT JOIN " + extraTableName + " extra ON extra.user_unionid = c.unionid")
		switch request.ReserveType {
		case customerReq.HasValidReserve:
			query.Where("extra.valid_reserve_at > ?", validTime)
		case customerReq.NoReserve:
			query.Where("extra.valid_reserve_at < ? AND extra.success_reserve_at < ?", validTime, validTime)
		case customerReq.HasSuccessReserve:
			query.Where("extra.success_reserve_at > ?", validTime)
		}
		switch request.TaskType {
		case customerReq.NoFinishTask:
			query.Where("extra.task_finish_at < ?", validTime)
		case customerReq.HasFinishTask:
			query.Where("extra.task_finish_at > ?", validTime)
		}
		switch request.TaskHave {
		case customerReq.TaskHave:
			query.Where("extra.task_created_at < ?", taskHaveValidTime)
		case customerReq.TaskNoHave:
			query.Where("extra.task_created_at > ?", taskHaveValidTime)
		}
		switch request.CouponType {
		case customerReq.Has30DayCoupon:
			query.Where("extra.coupon_by_30days_status = ?", externalSearchModel.HasCoupon)
		case customerReq.HasCoupon:
			query.Where("extra.coupon_status = ?", externalSearchModel.HasCoupon)
		}
	}
	// 近一年消费金额
	if len(request.ConsumeType) > 0 {
		consumeSqlStr := d.setConsumeQuery(request.ConsumeType)
		if consumeSqlStr != "" {
			query.Where("(" + consumeSqlStr + ")")
		}

	}
	// 最近消费开始时间
	if request.RecentConsumeStartTime > 0 {
		query.Where("c.last_purchase_date > ?", request.RecentConsumeStartTime)
	}
	// 最近消费结束时间
	if request.RecentConsumeEndTime > 0 {
		query.Where("c.last_purchase_date < ?", request.RecentConsumeEndTime)
	}
	// clt类型
	if len(request.CltType) > 0 {
		query.Where("c.clt_type in (?)", request.CltType)
	}
	// 会员卡等
	if len(request.CardType) > 0 {
		var cardLevel []string
		if len(request.CardType) == 1 {
			switch request.CardType[0] {
			case customerReq.HIGH_GRADE:
				// 优质客 clt类型有值
				query.Where("c.clt_type <> '' ")
			case customerReq.POTENTIAL:
				query.Where("c.card_level = ? and c.amount_total = 0", customerReq.PROSPECT)
			case "O":
				cardLevel = append(cardLevel, "O", "")
				query.Where("c.card_level in ?", cardLevel)
			default:
				query.Where("c.card_level = ?", request.CardType[0])
			}
		} else {
			highGrade := false
			potential := false
			for _, v := range request.CardType {
				if v == customerReq.HIGH_GRADE {
					highGrade = true
				} else if v == customerReq.POTENTIAL {
					potential = true
				} else {
					if v == "O" { // 非会员 或为空
						cardLevel = append(cardLevel, "O", "")
					} else {
						cardLevel = append(cardLevel, v)
					}
				}
			}
			if len(cardLevel) > 0 {
				if !highGrade && !potential {
					query.Where("c.card_level in (?)", cardLevel)
				}
				if highGrade && potential {
					query.Where("c.card_level in (?) or c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
				}
				if highGrade && !potential {
					query.Where("c.card_level in (?) or c.clt_type <> ''", cardLevel)
				}
				if potential && !highGrade {
					query.Where("c.card_level in (?) or (c.card_level = ? and c.amount_total = 0)", cardLevel, customerReq.PROSPECT)
				}
			} else {
				if highGrade && potential {
					query.Where("c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
				}
			}
		}
	}

	// 会员code
	if len(request.UnionCodes) > 0 {
		query.Where("c.union_code in (?)", request.UnionCodes)
	}

	// sw顾客类型
	if len(request.CustomerType) > 0 {
		customerType := request.CustomerType
		// vip客
		if helper.InArray(workExternalModel.VipTypeChoice, request.CustomerType) {
			query.Where("c.vip_type = (?)", workExternalModel.VipType)

			item := []int{workExternalModel.VipTypeChoice}
			customerType = helper.DelItems(request.CustomerType, item)
		}
		if len(customerType) > 0 {
			query.Where("c.customer_type in (?)", customerType)
		}
	}

	// sw会员最近购买时间
	if request.LastPurchaseStartTime > 0 && request.LastPurchaseEndTime > 0 {
		query.Where("c.last_purchase_date >= (?) and c.last_purchase_date < (?) ", request.LastPurchaseStartTime, request.LastPurchaseEndTime)
	}

	// sw会员生日月份时间
	if request.BirthdayStart > 0 && request.BirthdayEnd > 0 {
		birthdayArr := helper.GetBirthdayArr(request.BirthdayStart, request.BirthdayEnd)
		query.Where("c.birth_month in (?) ", birthdayArr)
	}

	// coach的AgeTag 顾客类型标签
	if request.AgeTag > 0 {
		ageTagMap := make(map[int]string, 0)
		ageTagMap[workExternalModel.AgeTagOne] = "Z世代"
		ageTagMap[workExternalModel.AgeTagTwo] = "千禧一代"
		query.Where("c.age_tag = ? ", ageTagMap[request.AgeTag])
	}

	return query
}

// 整合拼接近一年消费金额查询条件
func (d *UserRelation) setConsumeQuery(consumeTypeSlice []int) string {
	sort.Ints(consumeTypeSlice)
	// 整合后的sql
	var consumeSqlStr string
	// 当前的开始金额
	var startCondition int
	// 当前的结束金额
	var endCondition int
	// 以一个条件中的开始金额
	var nextStartCondition int
	// 是否是第一次生成sql
	firstSqlSign := true

	for i, val := range consumeTypeSlice {
		// 是否当前条件需要生成sql的标识
		rowSqlSign := false
		tmpSlice := customerReq.ConsumeOptionMap[val]
		if len(tmpSlice) < 2 {
			// 数据异常忽略
			continue
		}
		if startCondition <= 0 {
			// 如果起始金额为0则获取当前起始金额
			startCondition = tmpSlice[0]
		}
		// 结束金额赋值
		endCondition = tmpSlice[1]
		if i < len(consumeTypeSlice)-1 {
			// 如果for循环中后面还有条件数据则判断后面的条件是否跟当前是连续的
			nextTmpSlice := customerReq.ConsumeOptionMap[consumeTypeSlice[i+1]]
			nextStartCondition = nextTmpSlice[0]
			if nextStartCondition-1 > endCondition {
				// 后面的条件和当前条件是不连续的则需要生成sql
				rowSqlSign = true
			}
		} else {
			// 后面没有条件了需要生成sql
			rowSqlSign = true
		}
		if rowSqlSign {
			// 需要生成sql
			if !firstSqlSign {
				// 如果不是第一次执行则需要加上 OR(因为条件之间是并集)
				consumeSqlStr += " OR "
			}
			// 金额要乘100，因为数据库存的是分
			consumeSqlStr += fmt.Sprintf("c.amount_total BETWEEN %d AND %d", startCondition*100, endCondition*100)
			// 生成过sql拼接则修改首次sql标识
			firstSqlSign = false
			// 生成完毕重置开始金额
			startCondition = 0
		}
	}
	return consumeSqlStr
}

// GetUnionIdsBySearchListsByCk  筛选出有关标签的unionid
func (d *UserRelation) GetUnionIdsBySearchListsByCk(c app.Ctx, enterpriseId int64, searchList request.SearchList, unionids []string, userIds []string, memberId int64) (data []customerReq.UserUnionIdAndMemberId, err error) {
	selectStr := "unionid,member_id"
	query := d.searchBuilderQuery(c, enterpriseId, searchList, userIds, memberId)
	var relationIds []int64
	err = query.Where("work_external_user_relation.unionid in ?", unionids).Pluck("work_external_user_relation.id", &relationIds).Error
	err = gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").Where(map[string]interface{}{
		"id": relationIds,
	}).Select(selectStr).Find(&data).Error
	return

}

// GetUserTotalCk 获取全选的好友关系的人数
func (d *UserRelation) GetUserTotalCk(c ctx.Context, request customerReq.ExternalUserListBaseReq) (int64, error) {
	var (
		total int64
		err   error
	)
	query := d.searchBuilderQuery(c, request.EnterpriseID, request.SearchList, request.Userids, request.MemberId)
	if len(request.MemberList) > 0 {
		query.Where("work_external_user_relation.member_id in (?) ", request.MemberList)
	} else {
		query.Where("work_external_user_relation.member_id = ?", request.MemberId)
	}
	if request.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
	}
	if request.IsGroup == true {
		newQuery := query.WithContext(c).Distinct("work_external_user_relation.work_external_user_id")
		newQuery.Count(&total)
	} else {
		newQuery := query.WithContext(c).Distinct("work_external_user_relation.id")
		newQuery.Count(&total)
	}
	if len(request.MemberList) > 0 {
		var relationIds []int64
		err = query.Select("work_external_user_relation.work_external_user_id as id").Pluck("id", &relationIds).Error
		if err != nil {
			return total, err
		}
		userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
		newQuery := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
			Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").Where(map[string]interface{}{
			"work_external_user_id": relationIds,
			"member_id":             request.MemberId,
			"relation_status":       0,
		})
		newQuery.Count(&total)
	}

	return total, err
}

func (d *UserRelation) GetExternalUserIdsCk(c ctx.Context, request customerReq.ExternalUserListBaseReq, userIds []string) ([]string, error) {
	var (
		data        []string
		relationIds []int64
		err         error
	)
	query := d.searchBuilderQuery(c, request.EnterpriseID, request.SearchList, userIds, request.MemberId)
	if len(request.MemberList) > 0 {
		query.Where("work_external_user_relation.member_id in (?) ", request.MemberList)
	} else {
		query.Where("work_external_user_relation.member_id = ?", request.MemberId)
	}
	if request.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
	}
	if len(request.FilterExternalUserId) > 0 {
		query.Not("c.id in ?", request.FilterExternalUserId)
	}
	err = query.Pluck("c.id", &relationIds).Limit(10000).Error
	if len(relationIds) > 0 {
		userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
		newQuery := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
			Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").Where(map[string]interface{}{
			"work_external_user_id": relationIds,
			"member_id":             request.MemberId,
			"relation_status":       0,
		})
		err = newQuery.Pluck("c.openid", &data).Error
	}
	return data, err
}

// GetUserListByRelationIdCk 通过relationId查找有关的外部联系人列表数据
func (d *UserRelation) GetUserListByRelationIdCk(c ctx.Context, enterpriseId int64, relationIds []int64, isSelf int, memberId int64) ([]workExternalModel.ExternalUserList, error) {
	var result []workExternalModel.ExternalUserList
	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)
	query := gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").
		Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id")
	if isSelf == 1 { // 只看好友的逻辑
		query.Where(map[string]interface{}{
			"work_external_user_relation.work_external_user_id": relationIds,
			"work_external_user_relation.member_id":             memberId,
			"work_external_user_relation.relation_status":       0,
		})
	} else {
		query.Where(map[string]interface{}{
			"work_external_user_relation.id": relationIds,
		})
	}
	selectStr := "work_external_user_relation.id as id,work_external_user_relation.work_external_user_id as work_external_user_id," +
		"work_external_user_relation.enterprise_id as enterprise_id," +
		"work_external_user_relation.remark as remark," +
		"work_external_user_relation.member_id as member_id," +
		"work_external_user_relation.is_key_sa as is_key_sa," +
		"work_external_user_relation.created_at as created_at," +
		"(c.ba_code) as sa_user_id," +
		"c.unionid as unionid," +
		"c.openid as openid," +
		"c.euid as euid," +
		"c.nick as nick," +
		"c.name as name," +
		"c.avatar as avatar," +
		"c.birthday as birthday," +
		"c.union_code as union_code," +
		"c.mobile as mobile," +
		"c.card_level as card_level," +
		"c.clt_type as clt_type," +
		"c.amount_total as amount_total," +
		"c.customer_code as customer_code," +
		"c.last_purchase_date as last_purchase_date," +
		"c.last_purchase_store_code as last_purchase_store_code," +
		"c.upgrade_need_amount as upgrade_need_amount"
	err := query.Select(selectStr).Find(&result).Error
	return result, err
}

// GetOnInheritSelectAllNum 在职分配全选的人数
func (d *UserRelation) GetOnInheritSelectAllNum(c ctx.Context, req customerReq.OnInheritReq, externalUserIds []int64, memberIds []int64, userIds []string) (int64, error) {
	var total int64
	query := d.searchBuilderQuery(c, req.EnterpriseID, req.SearchList, userIds, req.MemberId)
	if req.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}
	if len(req.MemberList) > 0 {
		query.Where("work_external_user_relation.member_id in (?)", req.MemberList)
	}
	if len(memberIds) > 0 {
		query.Where("work_external_user_relation.member_id in (?)", memberIds)
	}
	err := query.Not("work_external_user_relation.work_external_user_id in ? ", externalUserIds).Count(&total).Error
	return total, err
}

// GetOnInheritMemberIds  获取在职分配全选条件下的所有memberId
func (d *UserRelation) GetOnInheritMemberIds(c ctx.Context, request customerReq.OnInheritReq, userIds []string) ([]int64, error) {
	var (
		result    []int64
		err       error
		selectStr string
	)
	selectStr = "work_external_user_relation.member_id as member_id,"
	groupBy := "work_external_user_relation.member_id"
	query := d.searchBuilderQuery(c, request.EnterpriseID, request.SearchList, userIds, request.MemberId)
	if request.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
	}
	if len(request.RelationId) > 0 {
		query.Where("work_external_user_relation.id in (?)", request.RelationId)
	}
	if len(request.MemberList) > 0 {
		query.Where("work_external_user_relation.member_id in (?)", request.MemberList)
	}
	if len(request.FilterRelationId) > 0 {
		query.Not("work_external_user_relation.id in (?)", request.FilterRelationId)
	}

	query.Select(strings.TrimRight(selectStr, ","))
	err = query.Group(groupBy).Pluck("member_id", &result).Error
	if err != nil {
		return result, err
	}
	return result, err
}

// GetOnInheritSelectList 在职分配全选的人员列表,最多返回10000外部联系人
func (d *UserRelation) GetOnInheritSelectList(c ctx.Context, req customerReq.OnInheritReq, externalUserIds []int64, memberIds []int64, userIds []string) ([]customerReq.OnInheritExternalInfoRes, error) {
	var data []customerReq.OnInheritExternalInfoRes
	query := d.searchBuilderQuery(c, req.EnterpriseID, req.SearchList, userIds, req.MemberId)
	if req.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}
	if len(req.MemberList) > 0 {
		memberIds = helper.ArrayIntersect(memberIds, req.MemberList)
	}
	if len(memberIds) > 0 {
		query.Where("work_external_user_relation.member_id in (?)", memberIds)
	}
	query.Not("work_external_user_relation.work_external_user_id in ? ", externalUserIds)
	query.Select("work_external_user_relation.enterprise_id,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id,work_external_user_relation.user_id")
	err := query.Find(&data).Limit(10000).Error
	return data, err
}

// GetPlanUnionIdsBySearchListsByCk 筛选出与任务有关的数据
func (d *UserRelation) GetPlanUnionIdsBySearchListsByCk(c app.Ctx, enterpriseId int64, searchList request.SearchList, userId string, unionids []string) (data []string, err error) {
	query := d.searchBuilderQuery(c, enterpriseId, searchList, nil, 0)
	var relationIds []int64
	err = query.Where("work_external_user_relation.user_id=?", userId).Where("work_external_user_relation.unionid in ?", unionids).Pluck("work_external_user_relation.id", &relationIds).Error
	err = gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_relation").Where(map[string]interface{}{
		"id": relationIds,
	}).Pluck("work_external_user_relation.unionid", &data).Error
	return

}

// GetExtraData 获取数据的
func (d *UserRelation) GetExtraData(c ctx.Context, enterpriseId int64, unionIds []string) ([]externalSearchModel.ExternalExtraModel, error) {
	var extraList []externalSearchModel.ExternalExtraModel
	err := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", enterpriseId)).
		Where("user_unionid in ?", unionIds).
		Find(&extraList).Error
	if err != nil {
		return nil, err
	}
	return extraList, err
}

// GetActivityQrCodeStatistic  获取活动二维码统计数据的
func (d *UserRelation) GetActivityQrCodeStatistic(c ctx.Context, req activityQRCodeReq.StatisticsReq, stateSlice []string, statePageSlice []string) (activityQRCodeResp.TotalStatistics, []activityQRCodeResp.TotalStatistics, error) {
	var statisticsTotal activityQRCodeResp.TotalStatistics
	var statisticsList []activityQRCodeResp.TotalStatistics
	query := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_relation weur").
		Joins(fmt.Sprintf("LEFT JOIN (select card_level,unionid from micro_scrm_enterprise_%d_ck.wx_user_corp_external where card_level!='') wuce ON wuce.unionid=weur.unionid", req.EnterpriseId)).
		Where("weur.enterprise_id = ?", req.EnterpriseId)
	if req.StartTime > 0 {
		query.Where("weur.createtime > ?", req.StartTime)
	}
	if req.EndTime > 0 {
		query.Where("weur.createtime < ?", req.EndTime)
	}
	selectStr := "count(weur.id) as add_relation_total," +
		"count(DISTINCT weur.unionid) as add_customer_total," +
		"countIf(weur.relation_status > 0) as del_relation_total," +
		"uniqIf(wuce.unionid, wuce.card_level != '' AND wuce.card_level IS NOT NULL) as add_vip_total"
	query.Session(&gorm.Session{}).Where("weur.state in ?", stateSlice).Select(selectStr).
		Scan(&statisticsTotal)
	err := query.Where("weur.state in ?", statePageSlice).
		Select("weur.state as state," + selectStr).
		Group("weur.state").
		Scan(&statisticsList).Error
	return statisticsTotal, statisticsList, err
}

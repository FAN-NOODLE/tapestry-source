package externalSearchDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	membersDao "micro-scrm/internal/app/dao/members"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	openReq "micro-scrm/internal/app/http/request/open"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	swCrmCustomerModel "micro-scrm/internal/app/models/sw_crm_customer"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

type ExternalExtraCkSw struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewExternalExtraCkSw() *ExternalExtraCkSw {
	return &ExternalExtraCkSw{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      externalSearchModel.ExternalExtraSwModel{}.TableName(),
		localTableName: externalSearchModel.ExternalExtraSwModel{}.TableName() + "_local",
	}
}

type ExternalExtraDataSw struct {
	UserUnionid      string `json:"user_unionid"`
	ValidReserveAt   int    `json:"valid_reserve_at"`
	SuccessReserveAt int    `json:"success_reserve_at"`
	TaskFinishAt     int    `json:"task_finish_at"`
	TaskCreatedAt    int64  `json:"task_created_at"`
	LastPurchaseDate int    `json:"last_purchase_date"`
}

func (d *ExternalExtraCkSw) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (info externalSearchModel.ExternalExtraSwModel, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&info).Error
	return
}

func (d *ExternalExtraCkSw) GetList(c app.Ctx, enterpriseId, query interface{}, args ...interface{}) ([]externalSearchModel.ExternalExtraSwModel, error) {
	list := make([]externalSearchModel.ExternalExtraSwModel, 0)
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&list).Error
	return list, err
}
func (d *ExternalExtraCkSw) SetAllUserExtraData(c app.Ctx, enterpriseId int64) error {
	sql := fmt.Sprintf(`INSERT INTO  micro_scrm_enterprise_%d_ck.external_extra_local (
		id,
		user_unionid,
		created_at,
        updated_at
) 
SELECT
		id,
		unionid as user_unionid,
		created_at,
        updated_at
	FROM
	   micro_scrm_enterprise_%d_ck.wx_user_corp_external WHERE unionid!= '' AND id NOT IN (select id from micro_scrm_enterprise_%d_ck.external_extra)`,
		enterpriseId, enterpriseId, enterpriseId)
	err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
	if err != nil {
		return err
	}
	fmt.Printf("SW-extra缺少用户添加成功")
	return nil
}

// SaveReserveDataByAll  更新邀约相关数据
func (d *ExternalExtraCkSw) SaveReserveDataByAll(c app.Ctx, enterpriseId int64, uniqueMap map[string]ExternalExtraDataSw) error {
	var list []externalSearchModel.ExternalExtraSwModel
	var reserveUnionid []string
	for _, data := range uniqueMap {
		reserveUnionid = append(reserveUnionid, data.UserUnionid)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", reserveUnionid).Find(&list).Error
	if err != nil {
		return err
	}
	var newData []externalSearchModel.ExternalExtraSwModel // 待更新的数据
	// 循环查出来的非邀约数据，让有效数据合并
	for _, oldInfo := range list {
		if newInfo, ok := uniqueMap[oldInfo.UserUnionid]; ok {
			// 存在则更新老数据中需要更新的部分重新插入
			oldInfo.ValidReserveAt = helper.RuleThree(newInfo.ValidReserveAt > oldInfo.ValidReserveAt, newInfo.ValidReserveAt, oldInfo.ValidReserveAt)
			oldInfo.SuccessReserveAt = helper.RuleThree(newInfo.SuccessReserveAt > oldInfo.SuccessReserveAt, newInfo.SuccessReserveAt, oldInfo.SuccessReserveAt)
			oldInfo.InteractTime = helper.RuleThree(newInfo.ValidReserveAt > oldInfo.InteractTime, newInfo.ValidReserveAt, oldInfo.InteractTime)
			newData = append(newData, oldInfo)
		}
	}
	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		fmt.Printf("【邀约】顾客列表附加信息更新，邀约相关，更新条数:%d\n", len(newData))
	}
	return nil
}

// SaveTaskData  更新任务相关数据
func (d *ExternalExtraCkSw) SaveTaskData(c app.Ctx, enterpriseId int64, externalExtraDataMap map[string]ExternalExtraDataSw, queryType int) error {
	unixTime := time.Now().Unix()
	var list []externalSearchModel.ExternalExtraSwModel
	var taskUnionid []string
	for unionid := range externalExtraDataMap {
		taskUnionid = append(taskUnionid, unionid)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", taskUnionid).Find(&list).Error
	if err != nil {
		return err
	}
	var newData []externalSearchModel.ExternalExtraSwModel // 待新增的数据
	// 整理提交过来的数据，提交过来的数据有可能是重复的，所以
	for _, oldInfo := range list {
		if newInfo, ok := externalExtraDataMap[oldInfo.UserUnionid]; ok {
			// 存在老数据则需要更新则覆盖老数据重新插入
			if queryType == microScrmModel.TaskReach {
				oldInfo.TaskFinishAt = helper.RuleThree(newInfo.TaskFinishAt > oldInfo.TaskFinishAt, newInfo.TaskFinishAt, oldInfo.TaskFinishAt)
			}
			if queryType == microScrmModel.TaskExist {
				oldInfo.TaskCreatedAt = helper.RuleThree(newInfo.TaskCreatedAt > oldInfo.TaskCreatedAt, newInfo.TaskCreatedAt, oldInfo.TaskCreatedAt)
				oldInfo.InteractTime = helper.RuleThree(int(newInfo.TaskCreatedAt) > oldInfo.InteractTime, int(newInfo.TaskCreatedAt), oldInfo.InteractTime)
			}
			if queryType == microScrmModel.InteractTime {
				oldInfo.InteractTime = helper.RuleThree(newInfo.LastPurchaseDate > oldInfo.InteractTime, newInfo.LastPurchaseDate, oldInfo.InteractTime)

			}
			oldInfo.UpdatedAt = unixTime
			newData = append(newData, oldInfo)
		}
	}
	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		fmt.Printf("【任务】顾客列表附加信息更新，任务相关，插入条数:%d\n", len(newData))
	}

	return nil
}

// BatchSaveCouponData   批量更新顾客优惠券状态相关数据
func (d *ExternalExtraCkSw) BatchSaveCouponData(c app.Ctx, enterpriseId int64, dataMap map[string]openReq.UpdateCouponSwParams) error {
	unixTime := time.Now().Unix()
	var list []externalSearchModel.ExternalExtraSwModel
	var unionidSlice []string
	for unionid := range dataMap {
		unionidSlice = append(unionidSlice, unionid)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", unionidSlice).Find(&list).Error
	if err != nil {
		return err
	}
	// 整理已经存在的数据map,以Unionid为key
	var existMap = make(map[string]externalSearchModel.ExternalExtraSwModel)
	for _, exist := range list {
		existMap[exist.UserUnionid] = exist
	}
	var newData []externalSearchModel.ExternalExtraSwModel // 待新增的数据
	// 整理提交过来的数据
	for _, row := range dataMap {
		if oldInfo, ok := existMap[row.UnionId]; ok {
			if row.HasCoupon > 0 {
				oldInfo.CouponStatus = helper.RuleThree(row.HasCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasCouponBy30days > 0 {
				oldInfo.CouponBy30daysStatus = helper.RuleThree(row.HasCouponBy30days == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasNewMemberCoupon > 0 {
				oldInfo.NewMemberCouponStatus = helper.RuleThree(row.HasNewMemberCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasBirthdayCoupon > 0 {
				oldInfo.BirthdayCouponStatus = helper.RuleThree(row.HasBirthdayCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasRebateCoupon > 0 {
				oldInfo.RebateCouponStatus = helper.RuleThree(row.HasRebateCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasPromotionalCoupon > 0 {
				oldInfo.PromotionalCouponStatus = helper.RuleThree(row.HasPromotionalCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			oldInfo.CouponUpdatedAt = row.UpdatedAt
			oldInfo.UpdatedAt = unixTime
			newData = append(newData, oldInfo)
		}
	}
	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(tableName)
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		log.Info(c).Msgf("【优惠券】顾客列表附加信息更新，任务相关，插入条数:%d\n", len(newData))
	}

	return nil
}

// ExternalSearchListCk 根据条件筛选导购和外部联系人列表
func (d *ExternalExtraCkSw) ExternalSearchListCk(c ctx.Context, request customerSearchReq.ExternalSearchListBaseSwReq) ([]customerSearchResp.ExternalUserListSw, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserListSw, 0)
	query := d.GetExternalQuery(c, request.SearchListSw)
	err = query.Count(&total).Error
	if err != nil || request.OnlyCount || total == 0 {
		return list, total, err
	}
	if request.Size == 0 {
		request.Size = 10000 //控制一下最大获取数量
	}
	if !request.GetAll {
		if request.Page == 0 {
			request.Page = 1
		}
		query = query.Offset((request.Page - 1) * request.Size)
	}

	externalField := "id,ba_code,unionid,openid,nick,name,avatar,birthday," +
		"card_level,amount_total,vip_type,last_purchase_date,union_code,last_purchase_store_code"
	var externalList = make([]microScrmModel.WxUserCorpExternalModel, 0)
	err = query.Select(externalField).Order("id desc").Limit(request.Size).Find(&externalList).Error
	if err != nil {
		return list, total, err
	}

	var unionidSlice []string
	var externalIds []int64
	for _, row := range externalList {
		unionidSlice = append(unionidSlice, row.Unionid)
		externalIds = append(externalIds, row.ID)
		tmp := customerSearchResp.ExternalUserListSw{}
		helper.CoverStruct(&tmp, &row)
		tmp.ID = row.ID
		list = append(list, tmp)
	}
	//获取附加信息
	var extraList []externalSearchModel.ExternalExtraSwModel
	err = gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)).
		Where("user_unionid in ?", unionidSlice).
		Select("id,coupon_status,coupon_by_30days_status,valid_reserve_at,success_reserve_at,task_finish_at").
		Find(&extraList).Error
	var extraMap = make(map[int64]externalSearchModel.ExternalExtraSwModel)
	for _, row := range extraList {
		extraMap[row.ID] = row
	}

	var relationShipList []microScrmModel.WorkExternalUserRelationModel
	//获取备注
	err = gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_relation").
		Where("enterprise_id = ?", request.EnterpriseID).
		Where("member_id = ?", request.MemberId).
		Where("work_external_user_id in ?", externalIds).
		Select("work_external_user_id,remark,relation_status").
		Find(&relationShipList).Error
	var relationShipMap = make(map[int64]microScrmModel.WorkExternalUserRelationModel)
	for _, row := range relationShipList {
		relationShipMap[row.WorkExternalUserID] = row
	}
	var nowUnix = time.Now().Unix()
	for i, row := range list {
		//将结果提取到一个临时结构体
		list[i].SelfRelationStatus = externalSearchModel.RelationDelAll
		if exist, ok := relationShipMap[row.ID]; ok {
			list[i].Remark = exist.Remark
			list[i].SelfRelationStatus = int64(exist.RelationStatus)
		}
		if exist, ok := extraMap[row.ID]; ok {
			list[i].ValidReserveAt = exist.ValidReserveAt
			list[i].SuccessReserveAt = exist.SuccessReserveAt
			list[i].TaskFinishAt = exist.TaskFinishAt
			if exist.TaskFinishAt > 0 && nowUnix-int64(exist.TaskFinishAt) < 86400*30 {
				list[i].TaskFinishAtTag = true
			}
			if exist.ValidReserveAt > 0 && nowUnix-int64(exist.ValidReserveAt) < 86400*customerSearchReq.ValidDays {
				list[i].ValidReserveAtTag = true
			}
			list[i].CouponStatus = exist.CouponStatus
			list[i].CouponBy30daysStatus = exist.CouponBy30daysStatus
		}
	}
	return list, total, err
}

// GetAllByIdMap 查询添加导购信息
func (d *ExternalExtraCkSw) GetAllByIdMap(c app.Ctx, enterpriseID int64, memberIds []int64) (map[int64]membersModel.MemberModel, error) {

	var memberNameListMap = make(map[int64]membersModel.MemberModel, 0)
	memberInfoLists, err := membersDao.NewMembers().GetAllById(c, enterpriseID, memberIds)
	if err != nil {
		return memberNameListMap, err
	}
	for _, item := range memberInfoLists {
		memberNameListMap[item.ID] = item
	}
	return memberNameListMap, nil
}

// RelationShipListCk   根据条件筛选导购和外部联系人关系列表
func (d *ExternalExtraCkSw) RelationShipListCk(c ctx.Context, request customerSearchReq.ExternalSearchListBaseSwReq) ([]customerSearchResp.ExternalUserListSw, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserListSw, 0)
	query := d.GetRelationShipQuery(c, request.SearchListSw)
	if !request.GetAll {
		err = query.Count(&total).Error
		if err != nil || request.OnlyCount || total == 0 {
			return list, total, err
		}
		if request.Page > 0 && request.Size > 0 {
			query = query.Offset((request.Page - 1) * request.Size).Limit(request.Size)
		}
	}
	relationShipField := "id,remark,user_id,relation_status,work_external_user_id,unionid,member_id,created_at"
	var relationShipList []microScrmModel.WorkExternalUserRelationModel
	err = query.Select(relationShipField).Order("id desc").Find(&relationShipList).Error
	if err != nil {
		return list, total, err
	}
	var unionidSlice []string
	var memberMap = make(map[int64]membersModel.MemberModel)
	if request.IsInherit {
		var memberIds []int64
		for _, row := range relationShipList {
			memberIds = append(memberIds, row.MemberID)
		}
		memberMap, _ = d.GetAllByIdMap(c, request.EnterpriseID, memberIds)
	}
	for _, row := range relationShipList {
		unionidSlice = append(unionidSlice, row.Unionid)
		tmp := customerSearchResp.ExternalUserListSw{}
		helper.CoverStruct(&tmp, &row)
		tmp.ID = row.ID
		tmp.CreatedAt = row.CreatedAt
		if memberInfo, ok := memberMap[row.MemberID]; ok {
			tmp.MemberName = memberInfo.Name
		}
		list = append(list, tmp)
	}
	var externalList = make([]microScrmModel.WxUserCorpExternalModel, 0)
	externalField := "ba_code,unionid,openid,nick,name,avatar,birthday," +
		"card_level,amount_total,vip_type,last_purchase_date,last_purchase_store_code"
	//获取外部联系人信息
	err = gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)).
		Where("unionid in ?", unionidSlice).
		Select(externalField).
		Find(&externalList).Error
	var externalMap = make(map[string]microScrmModel.WxUserCorpExternalModel)
	for _, row := range externalList {
		externalMap[row.Unionid] = row
	}
	//获取附加信息
	var extraList []externalSearchModel.ExternalExtraSwModel
	err = gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)).
		Where("user_unionid in ?", unionidSlice).
		Select("user_unionid,valid_reserve_at,success_reserve_at,task_finish_at").
		Find(&extraList).Error
	var extraMap = make(map[string]externalSearchModel.ExternalExtraSwModel)
	for _, row := range extraList {
		extraMap[row.UserUnionid] = row
	}
	for i, row := range list {
		//将结果提取到一个临时结构体
		if exist, ok := externalMap[row.Unionid]; ok {
			helper.CoverStruct(&row, &exist)
			list[i] = row
		}
		if exist, ok := extraMap[row.Unionid]; ok {
			list[i].ValidReserveAt = exist.ValidReserveAt
			list[i].SuccessReserveAt = exist.SuccessReserveAt
			list[i].TaskFinishAt = exist.TaskFinishAt
		}
	}
	return list, total, err
}

// ExternalSearchPluckIds   根据条件筛选导购和外部联系人
func (d *ExternalExtraCkSw) ExternalSearchPluckIds(c ctx.Context, searchReq customerSearchReq.SearchListSw, pluckName string) ([]interface{}, error) {
	query := d.GetExternalQuery(c, searchReq)
	var pluckIds = make([]interface{}, 0)
	err := query.Group(pluckName).Pluck(pluckName, &pluckIds).Error
	return pluckIds, err
}

// GetExternalQuery 顾客列表查询组装
func (d *ExternalExtraCkSw) GetExternalQuery(c ctx.Context, request customerSearchReq.SearchListSw) *gorm.DB {
	searchQueryResp := d.SearchExternalOptionQuery(request)
	relationshipTable := "micro_scrm_ck.work_external_user_relation"
	externalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
	//顾客列表
	query := gormx.GetCkInstance().WithContext(c).Table(externalTable).
		Where("type = ?", externalSearchModel.WeChat) // 只查微信用户
	if len(searchQueryResp.External.SqlSlice) > 0 {
		query = query.Where(searchQueryResp.External.SqlConn, searchQueryResp.External.SqlSlice...)
	}
	if len(searchQueryResp.RelationShip.SqlSlice) > 0 {
		query = query.Where(fmt.Sprintf("id in (SELECT work_external_user_id FROM %s WHERE %s)", relationshipTable, searchQueryResp.RelationShip.SqlConn),
			searchQueryResp.RelationShip.SqlSlice...)
	}
	return query
}

// GetRelationShipQuery  顾客关系列表查询组装
func (d *ExternalExtraCkSw) GetRelationShipQuery(c ctx.Context, request customerSearchReq.SearchListSw) *gorm.DB {
	searchQueryResp := d.SearchExternalOptionQuery(request)
	relationshipTable := "micro_scrm_ck.work_external_user_relation"
	externalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
	//顾客关系列表
	query := gormx.GetCkInstance().WithContext(c).Table(relationshipTable).
		Where("enterprise_id = ? AND unionid != ?", request.EnterpriseID, "") // 只查微信用户
	if len(searchQueryResp.RelationShip.SqlSlice) > 0 {
		query = query.Where(searchQueryResp.RelationShip.SqlConn, searchQueryResp.RelationShip.SqlSlice...)
	}
	if len(searchQueryResp.External.SqlSlice) > 0 {
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT id FROM %s WHERE %s)", externalTable, searchQueryResp.External.SqlConn),
			searchQueryResp.External.SqlSlice...)
	}
	return query
}

// SearchExternalOptionQuery 组装筛选条件的sqlbuilder
func (d *ExternalExtraCkSw) SearchExternalOptionQuery(request customerSearchReq.SearchListSw) customerSearchResp.SearchQuerySwResp {
	unixTime := time.Now().Unix()
	externalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", request.EnterpriseID)
	relationshipTable := "micro_scrm_ck.work_external_user_relation"

	if len(request.MemberIds) == 0 { //如果没有圈定员工范围，则范围就是自己
		request.MemberIds = []int64{request.MemberId}
	}
	var relationshipSlice []interface{}
	var relationshipConn = "enterprise_id = ? AND deleted_at = 0"
	relationshipSlice = append(relationshipSlice, request.EnterpriseID)
	var externalExtraSlice []interface{}
	var externalExtraConn = "deleted_at = 0"
	var externalSlice []interface{}
	var externalConn = "deleted_at = 0"

	//好友非好友，如果没有选择好友关系则说明本身已经有顾客范围了，不需要根据员工去圈定
	if request.RelationshipType == customerSearchReq.IsUnionFalse {
		externalConn += " AND id NOT IN (SELECT work_external_user_id FROM " + relationshipTable + " WHERE member_id in (?) AND relation_status = ? )"
		externalExtraSlice = append(externalExtraSlice, request.MemberIds, 0)
	} else if request.RelationshipType == customerSearchReq.IsUnionTure {
		relationshipConn += " AND member_id IN (?) AND relation_status = ? "
		relationshipSlice = append(relationshipSlice, request.MemberIds, 0)
	}

	if len(request.UnionIds) > 0 {
		externalConn += " AND user_unionid  in ?"
		externalSlice = append(externalSlice, request.UnionIds)
	}
	if len(request.UnionCodes) > 0 {
		externalConn += " AND union_code  in ?"
		externalSlice = append(externalSlice, request.UnionCodes)
	}
	if len(request.FilterExternalUserId) > 0 {
		externalConn += " AND id not in ?"
		externalSlice = append(externalSlice, request.FilterExternalUserId)
	}
	if len(request.FilterRelationId) > 0 {
		relationshipConn += " AND id not in ?"
		relationshipSlice = append(relationshipSlice, request.FilterRelationId)
	}
	if len(request.ExternalUserIds) > 0 {
		externalConn += " AND id in ?"
		externalSlice = append(externalSlice, request.ExternalUserIds)
	}
	if len(request.RelationId) > 0 {
		relationshipConn += " AND id in ?"
		relationshipSlice = append(relationshipSlice, request.RelationId)
	}
	if request.Name != "" {
		relationshipConn += " AND (LOWER(remark) LIKE LOWER(?) OR work_external_user_id GLOBAL IN (SELECT id FROM " + externalTable + " WHERE LOWER(nick) LIKE LOWER(?) ))"
		relationshipSlice = append(relationshipSlice, "%"+request.Name+"%", "%"+request.Name+"%")
	}
	// 会员卡等
	if len(request.CardType) > 0 {
		var cardLevel []string
		for _, v := range request.CardType {
			if v == "O" { // 非会员 或为空
				cardLevel = append(cardLevel, "O", "")
			} else if v == "B" { // 钻石「履」行家 有两个
				cardLevel = append(cardLevel, "B", "Z")
			} else {
				cardLevel = append(cardLevel, v)
			}
		}
		if len(cardLevel) > 0 {
			externalConn += " AND card_level  in ?"
			externalSlice = append(externalSlice, cardLevel)
		}
	}
	if len(request.CustomerType) > 0 {
		customerType := request.CustomerType
		// vip客
		if helper.InArray(microScrmModel.VipTypeChoice, request.CustomerType) {
			externalConn += " AND vip_type = ?"
			externalSlice = append(externalSlice, microScrmModel.VipType)

			item := []int{microScrmModel.VipTypeChoice}
			customerType = helper.DelItems(request.CustomerType, item)
		}
		if len(customerType) > 0 {
			externalConn += " AND customer_type in (?)"
			externalSlice = append(externalSlice, customerType)
		}
	}
	//生日月份
	if len(request.BirthMonth) > 0 {
		externalConn += " AND birth_month  in ?"
		externalSlice = append(externalSlice, request.BirthMonth)
	}

	// 1【专属】指key_sa顾问是自己的顾客
	// 2【非专属】指key_sa顾问不是自己的顾客
	if request.ManageType == customerSearchReq.IsKeySA {
		relationshipConn += " AND is_key_sa = ?"
		relationshipSlice = append(relationshipSlice, 1)
	} else if request.ManageType == customerSearchReq.NotKeySA {
		relationshipConn += " AND is_key_sa = ?"
		relationshipSlice = append(relationshipSlice, 0)
	}

	// 邀约和任务的检索
	// 邀约查询的数据有效时间是30天内的
	validTime := time.Now().Unix() - customerSearchReq.ValidDays*86400
	if request.ReserveType > 0 {
		switch request.ReserveType {
		case customerSearchReq.HasValidReserve:
			externalExtraConn += " AND valid_reserve_at > ?"
			externalExtraSlice = append(externalExtraSlice, validTime)
		case customerSearchReq.NoReserve:
			externalExtraConn += " AND valid_reserve_at < ? AND success_reserve_at < ?"
			externalExtraSlice = append(externalExtraSlice, validTime, validTime)
		case customerSearchReq.HasSuccessReserve:
			externalExtraConn += " AND success_reserve_at > ?"
			externalExtraSlice = append(externalExtraSlice, validTime)
		}
	}
	// 任务触达查询的数据有效时间是30天内的
	if request.TaskType > 0 {
		switch request.TaskType {
		case customerSearchReq.NoFinishTask:
			externalExtraConn += " AND task_finish_at < ?"
			externalExtraSlice = append(externalExtraSlice, validTime)
		case customerSearchReq.HasFinishTask:
			externalExtraConn += " AND task_finish_at > ?"
			externalExtraSlice = append(externalExtraSlice, validTime)
		}
	}
	//有无任务查询的是30天的
	taskHaveValidTime := time.Now().Unix() - 30
	if request.TaskHave > 0 {
		switch request.TaskHave {
		case customerSearchReq.TaskNoHave:
			externalExtraConn += " AND task_created_at < ?"
			externalExtraSlice = append(externalExtraSlice, taskHaveValidTime)
		case customerSearchReq.TaskHave:
			externalExtraConn += " AND task_created_at > ?"
			externalExtraSlice = append(externalExtraSlice, taskHaveValidTime)
		}
	}
	//添加企微好友时间
	if len(request.RelCreateTime) > 0 {
		var relAddTimeConn string
		var relAddTimeSlice []interface{}
		relAddTimeConn += " AND ("
		for _, v := range request.RelCreateTime {
			switch v {
			case customerSearchReq.TimeLt1D:
				relAddTimeConn += " createtime >= ? OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-86400)
			case customerSearchReq.Time1To30D:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-30*86400, unixTime-86400)
			case customerSearchReq.Time31To60D:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-60*86400, unixTime-30*86400)
			case customerSearchReq.Time61To90D:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-90*86400, unixTime-60*86400)
			case customerSearchReq.TimeGt90D:
				relAddTimeConn += " createtime <= ? OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-90*86400)
			}
		}
		relAddTimeConn = strings.TrimSuffix(relAddTimeConn, "OR")
		relAddTimeConn += ") "
		relationshipConn += relAddTimeConn
		relationshipSlice = append(relationshipSlice, relAddTimeSlice...)
	}
	//上次触达时间
	if len(request.RelInteractTime) > 0 {
		externalExtraConn += " AND ("
		for _, v := range request.RelInteractTime {
			switch v {
			case customerSearchReq.TimeLt30D:
				externalExtraConn += " interact_time >= ? OR"
				externalExtraSlice = append(externalExtraSlice, unixTime-30*86400)
			case customerSearchReq.Time31To60D:
				externalExtraConn += " (interact_time BETWEEN ? AND ?) OR"
				externalExtraSlice = append(externalExtraSlice, unixTime-60*86400, unixTime-30*86400)
			case customerSearchReq.Time61To90D:
				externalExtraConn += " (interact_time BETWEEN ? AND ?) OR"
				externalExtraSlice = append(externalExtraSlice, unixTime-90*86400, unixTime-60*86400)
			case customerSearchReq.TimeGt90D:
				externalExtraConn += " interact_time <= ? OR"
				externalExtraSlice = append(externalExtraSlice, unixTime-90*86400)
			}
		}
		externalExtraConn = strings.TrimSuffix(externalExtraConn, "OR")
		externalExtraConn += ") "
	}
	//历史消费时间
	if request.PurchasedDate.StartTime > 0 || request.PurchasedDate.EndTime > 0 {
		externalExtraConn += " AND arrayExists(x ->"
		if request.PurchasedDate.StartTime > 0 {
			externalExtraConn += "x >= ? AND"
			externalExtraSlice = append(externalExtraSlice, request.PurchasedDate.StartTime)
		}
		if request.PurchasedDate.EndTime > 0 {
			externalExtraConn += " x <= ?"
			externalExtraSlice = append(externalExtraSlice, request.PurchasedDate.EndTime)
		} else {
			externalExtraConn = strings.TrimSuffix(externalExtraConn, "AND")
		}
		externalExtraConn += ",purchase_date)"
	}

	// 近期未消费时间
	if request.NotPurchasedDate > 0 {
		var tmpTime int64
		switch request.NotPurchasedDate {
		case customerSearchReq.TimeGt30D:
			tmpTime = unixTime - 30*86400
		case customerSearchReq.TimeGt90D:
			tmpTime = unixTime - 90*86400
		case customerSearchReq.TimeGt180D:
			tmpTime = unixTime - 180*86400
		}
		//近期未消费时间范围如果有开始时间则以开始时间为准，如果没有开始时间则以结束时间为准
		externalConn += " AND last_purchase_date < ?"
		externalSlice = append(externalSlice, tmpTime)
	}

	//优惠券筛选
	if len(request.CouponType) > 0 {
		externalExtraConn += " AND ("
		for _, v := range request.CouponType {
			switch v {
			case customerSearchReq.SwHas30DayCoupon:
				externalExtraConn += " coupon_by_30days_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.SwHasCoupon:
				externalExtraConn += " coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.SwBirthCoupon:
				externalExtraConn += " birthday_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.SwNewMemberCoupon:
				externalExtraConn += " new_member_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.SwRebateCoupon:
				externalExtraConn += " rebate_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.SwPromotionalCoupon:
				externalExtraConn += " promotional_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			}
		}
		externalExtraConn = strings.TrimSuffix(externalExtraConn, "OR")
		externalExtraConn += ") "

	}
	//近12个月消费
	if request.R12PurchaseAmount.StartNumber > 0 {
		externalExtraConn += " AND r12_purchase_amount >= ?"
		externalExtraSlice = append(externalExtraSlice, request.R12PurchaseAmount.StartNumber)
	}
	//近12个月消费
	if request.R12PurchaseAmount.EndNumber > 0 {
		externalExtraConn += " AND r12_purchase_amount <= ?"
		externalExtraSlice = append(externalExtraSlice, request.R12PurchaseAmount.EndNumber)
	}
	//价格标签
	if len(request.PriceTag) > 0 {
		externalExtraConn += " AND price_tag in ?"
		externalExtraSlice = append(externalExtraSlice, request.PriceTag)
	}

	//年龄标签
	if len(request.AgeTag) > 0 {
		externalExtraConn += " AND age_tag in ?"
		externalExtraSlice = append(externalExtraSlice, request.AgeTag)
	}
	//季度消费
	if request.AtPurchasedYear != "" && request.AtPurchasedQuarter != "" && len(request.AtPurchasedAmount) > 0 {
		var atValueConn string
		var atValueSlice []interface{}
		for _, v := range request.AtPurchasedAmount {
			switch v {
			case customerSearchReq.SwAtPurchasedAmountGt8000:
				atValueConn += " toFloat64(tag_value) >= ? OR"
				atValueSlice = append(atValueSlice, 8000)
			case customerSearchReq.SwAtPurchasedAmount6000To8000:
				atValueConn += " (toFloat64(tag_value) >= ? AND toFloat64(tag_value) < ?) OR"
				atValueSlice = append(atValueSlice, 6000, 8000)
			case customerSearchReq.SwAtPurchasedAmount4000To6000:
				atValueConn += " (toFloat64(tag_value) >= ? AND toFloat64(tag_value) < ?) OR"
				atValueSlice = append(atValueSlice, 4000, 6000)
			case customerSearchReq.SwAtPurchasedAmountLt4000:
				atValueConn += " toFloat64(tag_value) < ? OR"
				atValueSlice = append(atValueSlice, 4000)
			}
		}
		atValueConn = strings.TrimSuffix(atValueConn, "OR")
		externalConn += fmt.Sprintf(" AND unionid in (select unionid from micro_scrm_enterprise_%d_ck.external_tag where tag_type = ? AND tag_key = ? AND (%s) AND deleted_at = 0)", request.EnterpriseID, atValueConn)
		externalSlice = append(externalSlice, externalSearchModel.SWAtAmount, fmt.Sprintf("%s%s", request.AtPurchasedYear, request.AtPurchasedQuarter))
		externalSlice = append(externalSlice, atValueSlice...)

	}

	//商品偏好
	if len(request.CategoryPreference) > 0 {
		externalExtraConn += " AND ("
		for _, tag := range request.CategoryPreference {
			externalExtraConn += " has(category_preference,?) OR"
			externalExtraSlice = append(externalExtraSlice, tag)
		}
		externalExtraConn = strings.TrimSuffix(externalExtraConn, "OR")
		externalExtraConn += ") "
	}

	// 顾客沟通意愿标签，只查看自己打的标签数据
	if request.ChatTag > 0 && request.UserId != "" {
		externalConn += " AND id in (select external_user_id from micro_scrm_ck.impression_tag_relation where tag_id = ? AND user_id = ? AND deleted_at = 0)"
		externalSlice = append(externalSlice, request.ChatTag, request.UserId)
	}

	// 顾客沟通意愿标签，只查看自己打的标签数据
	if len(request.ImpressionTagList) > 0 {
		externalConn += " AND id in (select external_user_id from micro_scrm_ck.impression_tag_relation where tag_id in ? AND deleted_at = 0)"
		externalSlice = append(externalSlice, request.ImpressionTagList)
	}

	// 顾客沟通意愿标签，只查看自己打的标签数据
	if len(request.TagList) > 0 {
		relationshipConn += " AND id in (select relation_id from micro_scrm_ck.work_external_user_tag where external_tag_id in ? AND deleted_at = 0)"
		relationshipSlice = append(relationshipSlice, request.TagList)
	}

	if len(externalExtraSlice) > 0 {
		externalExtraTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)
		externalConn += fmt.Sprintf(" AND id in (select id from %s where %s)", externalExtraTable, externalExtraConn)
		externalSlice = append(externalSlice, externalExtraSlice...)
	}
	return customerSearchResp.SearchQuerySwResp{
		External: customerSearchResp.BuildSql{
			SqlConn:  externalConn,
			SqlSlice: externalSlice,
		},
		RelationShip: customerSearchResp.BuildSql{
			SqlConn:  relationshipConn,
			SqlSlice: relationshipSlice,
		},
	}
}

// GetTaskReachOrReserveAtExtraByDay   获取任务触达和邀约时间最新的更新信息
func (d *ExternalExtraCkSw) GetTaskReachOrReserveAtExtraByDay(c app.Ctx, enterpriseId int64, days int64) ([]externalSearchModel.ExternalExtraSwModel, error) {
	var list = make([]externalSearchModel.ExternalExtraSwModel, 0)
	unixTime := time.Now().Unix() - 86400*days
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("task_finish_at > ? or valid_reserve_at > ?", unixTime, unixTime).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// SavePurchaseData 保存订单信息
func (d *ExternalExtraCkSw) SavePurchaseData(c app.Ctx, enterpriseId int64, dataList []swCrmCustomerModel.CrmCustomerPurchaseModel) error {
	var list []externalSearchModel.ExternalExtraSwModel
	var unionIds []string
	for _, data := range dataList {
		unionIds = append(unionIds, data.UnionID)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", unionIds).Find(&list).Error
	if err != nil {
		return err
	}
	var newData []externalSearchModel.ExternalExtraSwModel // 待更新的数据
	// 循环查出来的非邀约数据，让有效数据合并
	for _, oldInfo := range list {
		for _, newInfo := range dataList {
			if oldInfo.UserUnionid == newInfo.UnionID {
				if newInfo.CategoryCode != "" {
					oldInfo.CategoryPreference = helper.ArrayUniq(append(oldInfo.CategoryPreference, strings.Split(newInfo.CategoryCode, "|")...))
				}
				if newInfo.PurchaseTime != "" {
					oldInfo.PurchaseDate = helper.ArrayUniq(append(oldInfo.PurchaseDate, helper.GetDateToTimes(newInfo.PurchaseTime)))
				}
				newData = append(newData, oldInfo)
			}
		}
	}
	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		fmt.Printf("【SW】顾客订单数据落表，更新条数:%d\n", len(newData))
	}
	return nil
}

// SaveExtensionData  保存用户附加信息
func (d *ExternalExtraCkSw) SaveExtensionData(c app.Ctx, enterpriseId int64, dataList []swCrmCustomerModel.CrmCustomerExtensionModel) error {
	var list []externalSearchModel.ExternalExtraSwModel
	var unionIds []string
	var dataMap = make(map[string]swCrmCustomerModel.CrmCustomerExtensionModel)
	for _, data := range dataList {
		unionIds = append(unionIds, data.UnionID)
		dataMap[data.UnionID] = data
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", unionIds).Find(&list).Error
	if err != nil {
		return err
	}
	// 整理已经存在的数据map,以Unionid为key
	var existMap = make(map[string]externalSearchModel.ExternalExtraSwModel)
	for _, exist := range list {
		existMap[exist.UserUnionid] = exist
	}
	var newData []externalSearchModel.ExternalExtraSwModel // 待更新的数据
	// 循环查出来的非邀约数据，让有效数据合并
	for _, newInfo := range dataList {
		if oldInfo, ok := existMap[newInfo.UnionID]; ok {
			if newInfo.PriceTag != "" {
				oldInfo.PriceTag = newInfo.PriceTag
			}
			if newInfo.AgeTag != "" {
				oldInfo.AgeTag = newInfo.AgeTag
			}
			if newInfo.R12PurchaseAmount > 0 {
				oldInfo.R12PurchaseAmount = int(newInfo.R12PurchaseAmount)
			}
			newData = append(newData, oldInfo)
		}
	}
	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		fmt.Printf("【SW】顾客标签附加信息更新，更新条数:%d\n", len(newData))
	}
	return nil
}

// MergeData  ck数据合并
func (d *ExternalExtraCkSw) MergeData(c app.Ctx, enterpriseId int64) error {
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	return gormx.GetCkInstance().WithContext(c).Debug().Exec(fmt.Sprintf("OPTIMIZE TABLE %s FINAL DEDUPLICATE", tableName)).Error
}

func (d *ExternalExtraCkSw) GetExternalIdsByOpenIds(c app.Ctx, enterpriseId int64, openIds []string) ([]int64, error) {
	var ids = make([]int64, 0)
	externalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)
	//顾客列表
	err := gormx.GetCkInstance().WithContext(c).Table(externalTable).
		Where("openid in ?", openIds).Pluck("id", &ids).Error
	return ids, err
}

package externalSearchDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	openReq "micro-scrm/internal/app/http/request/open"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

type ExternalExtraCk struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewExternalExtraCk() *ExternalExtraCk {
	return &ExternalExtraCk{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      externalSearchModel.ExternalExtraModel{}.TableName(),
		localTableName: externalSearchModel.ExternalExtraModel{}.TableName() + "_local",
	}
}

type ExternalExtraData struct {
	UserUnionid      string `json:"user_unionid"`
	ValidReserveAt   int    `json:"valid_reserve_at"`
	SuccessReserveAt int    `json:"success_reserve_at"`
	TaskFinishAt     int    `json:"task_finish_at"`
	TaskCreatedAt    int64  `json:"task_created_at"`
}

func (d *ExternalExtraCk) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (info externalSearchModel.ExternalExtraModel, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&info).Error
	return
}

func (d *ExternalExtraCk) GetList(c app.Ctx, enterpriseId, query interface{}, args ...interface{}) ([]externalSearchModel.ExternalExtraModel, error) {
	list := make([]externalSearchModel.ExternalExtraModel, 0)
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&list).Error
	return list, err
}

// SaveReserveDataByAll  更新邀约相关数据
func (d *ExternalExtraCk) SaveReserveDataByAll(c app.Ctx, enterpriseId int64, uniqueMap map[string]ExternalExtraData) error {
	var list []externalSearchModel.ExternalExtraModel
	var reserveUnionid []string
	for _, data := range uniqueMap {
		reserveUnionid = append(reserveUnionid, data.UserUnionid)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", reserveUnionid).Find(&list).Error
	if err != nil {
		return err
	}
	var newData []externalSearchModel.ExternalExtraModel // 待更新的数据
	// 循环查出来的非邀约数据，让有效数据合并
	for _, oldInfo := range list {
		if newInfo, ok := uniqueMap[oldInfo.UserUnionid]; ok {
			// 存在则更新老数据中需要更新的部分重新插入
			oldInfo.ValidReserveAt = helper.RuleThree(newInfo.ValidReserveAt > oldInfo.ValidReserveAt, newInfo.ValidReserveAt, oldInfo.ValidReserveAt)
			oldInfo.SuccessReserveAt = helper.RuleThree(newInfo.SuccessReserveAt > oldInfo.SuccessReserveAt, newInfo.SuccessReserveAt, oldInfo.SuccessReserveAt)
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
func (d *ExternalExtraCk) SaveTaskData(c app.Ctx, enterpriseId int64, externalExtraDataMap map[string]ExternalExtraData, queryType int) error {
	unixTime := time.Now().Unix()
	var list []externalSearchModel.ExternalExtraModel
	var taskUnionid []string
	for unionid := range externalExtraDataMap {
		taskUnionid = append(taskUnionid, unionid)
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("user_unionid in (?)", taskUnionid).Find(&list).Error
	if err != nil {
		return err
	}
	var newData []externalSearchModel.ExternalExtraModel // 待新增的数据
	// 整理提交过来的数据，提交过来的数据有可能是重复的，所以
	for _, oldInfo := range list {
		if newInfo, ok := externalExtraDataMap[oldInfo.UserUnionid]; ok {
			// 存在老数据则需要更新则覆盖老数据重新插入
			if queryType == microScrmModel.TaskReach {
				oldInfo.TaskFinishAt = helper.RuleThree(newInfo.TaskFinishAt > oldInfo.TaskFinishAt, newInfo.TaskFinishAt, oldInfo.TaskFinishAt)
			}
			if queryType == microScrmModel.TaskExist {
				oldInfo.TaskCreatedAt = helper.RuleThree(newInfo.TaskCreatedAt > oldInfo.TaskCreatedAt, newInfo.TaskCreatedAt, oldInfo.TaskCreatedAt)
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
func (d *ExternalExtraCk) BatchSaveCouponData(c app.Ctx, enterpriseId int64, dataMap map[string]openReq.UpdateCouponParams) error {
	unixTime := time.Now().Unix()
	var list []externalSearchModel.ExternalExtraModel
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
	var existMap = make(map[string]externalSearchModel.ExternalExtraModel)
	for _, exist := range list {
		existMap[exist.UserUnionid] = exist
	}
	var newData []externalSearchModel.ExternalExtraModel // 待新增的数据
	// 整理提交过来的数据
	for _, row := range dataMap {
		if oldInfo, ok := existMap[row.UnionId]; ok {
			if row.HasCoupon > 0 {
				oldInfo.CouponStatus = helper.RuleThree(row.HasCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasCouponBy30days > 0 {
				oldInfo.CouponBy30daysStatus = helper.RuleThree(row.HasCouponBy30days == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasBirthdayCoupon > 0 {
				oldInfo.BirthdayCouponStatus = helper.RuleThree(row.HasBirthdayCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasWelcomeCoupon > 0 {
				oldInfo.WelcomeCouponStatus = helper.RuleThree(row.HasWelcomeCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasRenewMembershipCoupon > 0 {
				oldInfo.RenewMembershipCouponStatus = helper.RuleThree(row.HasRenewMembershipCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
			}
			if row.HasUpgradeCoupon > 0 {
				oldInfo.UpgradeCouponStatus = helper.RuleThree(row.HasUpgradeCoupon == openReq.HasCoupon, externalSearchModel.HasCoupon, externalSearchModel.NoCoupon)
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

// MergeData  ck数据合并
func (d *ExternalExtraCk) MergeData(c app.Ctx, enterpriseId int64) error {
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	return gormx.GetCkInstance().WithContext(c).Debug().Exec(fmt.Sprintf("OPTIMIZE TABLE %s FINAL DEDUPLICATE", tableName)).Error
}

// UpdateExternalExtraCk    每3分钟将更新的数据写入物化视图数据
func (d *ExternalExtraCk) UpdateExternalExtraCk(c app.Ctx, enterpriseId int64, intervalInSeconds int64) error {
	unixTime := time.Now().Unix() - intervalInSeconds
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	// 拼接一个extra表的部分查询语句
	sql := fmt.Sprintf(`INSERT INTO %s (
		id,
		user_unionid,
		valid_reserve_at,
		success_reserve_at,
		task_finish_at,
		task_created_at,
		coupon_status,
		coupon_by_30days_status,
		coupon_updated_at,
		created_at,
		updated_at,
		deleted_at,
		birthday_coupon_status,
		renew_membership_coupon_status,
		welcome_coupon_status,
		upgrade_coupon_status,
		openid,
		nick,
		name,
		position,
		avatar,
		type,
		sex,
		birthday,
		union_code,
		point,
		grade,
		birth_year,
		birth_month,
		birth_day,
		mobile,
		card_level,
		clt_type,
		amount_total,
		history_amount_total,
		customer_code,
		customer_id,
		last_purchase_date,
		last_purchase_store_code,
		upgrade_need_amount,
		customer_type,
		customer_type_updated_at,
		vip_type,
		customer_deleted_at,
		age_tag,
		category_preference,
		purchased_date,
		high_willing,
		price_sensitivity,
		ba_code
)
SELECT
		wuce.id AS id,
		any(wuce.unionid) AS user_unionid,
		any(extra.valid_reserve_at) AS valid_reserve_at,
		any(extra.success_reserve_at) AS success_reserve_at,
		any(extra.task_finish_at) AS task_finish_at,
		any(extra.task_created_at) AS task_created_at,
		any(extra.coupon_status) AS coupon_status,
		any(extra.coupon_by_30days_status) AS coupon_by_30days_status,
		any(extra.coupon_updated_at) AS coupon_updated_at,
		any(wuce.created_at) AS created_at,
		GREATEST(max(wuce.updated_at),max(extra.updated_at)) AS updated_at,
		any(wuce.deleted_at) AS deleted_at,
		any(extra.birthday_coupon_status) AS birthday_coupon_status,
		any(extra.renew_membership_coupon_status) AS renew_membership_coupon_status,
		any(extra.welcome_coupon_status) AS welcome_coupon_status,
		any(extra.upgrade_coupon_status) AS upgrade_coupon_status,
		any(wuce.openid) AS openid,
		any(wuce.nick) AS nick,
		any(wuce.name) AS name,
		any(wuce.position) AS position,
		any(wuce.avatar) AS avatar,
		any(wuce.type) AS type,
		any(wuce.sex) AS sex,
		any(wuce.birthday) AS birthday,
		any(wuce.union_code) AS union_code,
		any(wuce.point) AS point,
		any(wuce.grade) AS grade,
		any(wuce.birth_year) AS birth_year,
		any(wuce.birth_month) AS birth_month,
		any(wuce.birth_day) AS birth_day,
		any(wuce.mobile) AS mobile,
		any(wuce.card_level) AS card_level,
		any(wuce.clt_type) AS clt_type,
		any(extra.amount_total) AS amount_total,
		any(extra.history_amount_total) AS history_amount_total,
		any(wuce.customer_code) AS customer_code,
		any(wuce.customer_id) AS customer_id,
		any(extra.last_purchase_date) AS last_purchase_date,
		any(extra.last_purchase_store_code) AS last_purchase_store_code,
		any(wuce.upgrade_need_amount) AS upgrade_need_amount,
		any(wuce.customer_type) AS customer_type,
		any(wuce.customer_type_updated_at) AS customer_type_updated_at,
		any(wuce.vip_type) AS vip_type,
		any(wuce.customer_deleted_at) AS customer_deleted_at,
		any(wuce.age_tag) AS age_tag,
		any(extra.category_preference) AS category_preference,
		any(extra.purchased_date) AS purchased_date,
		any(extra.high_willing) AS high_willing,
		any(extra.price_sensitivity) AS price_sensitivity,
		any(wuce.ba_code) AS ba_code
	FROM
	   (SELECT * FROM micro_scrm_enterprise_%d_ck.wx_user_corp_external WHERE updated_at >= %d) wuce
	   GLOBAL LEFT JOIN (SELECT * FROM micro_scrm_enterprise_%d_ck.external_extra  WHERE id IN (SELECT id
        FROM micro_scrm_enterprise_%d_ck.wx_user_corp_external 
        WHERE updated_at >= %d)) extra ON extra.id = wuce.id
	GROUP BY
	   wuce.id`, tableName, enterpriseId, unixTime, enterpriseId, enterpriseId, unixTime)
	err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
	if err != nil {
		return err
	}
	fmt.Printf("顾客列表大宽表更新成功，更新时间范围：%d分\n", intervalInSeconds/60)
	return nil
}

// FixExternalExtra  修复数据同步脚本可能遗漏的数据
func (d *ExternalExtraCk) FixExternalExtra(c app.Ctx, enterpriseId int64) error {
	sql := fmt.Sprintf(`INSERT INTO  micro_scrm_enterprise_%d_ck.external_extra_local (
		id,
		user_unionid,
		created_at,
		updated_at,
		deleted_at,
		openid,
		nick,
		name,
		position,
		avatar,
		type,
		sex,
		birthday,
		union_code,
		point,
		grade,
		birth_year,
		birth_month,
		birth_day,
		mobile,
		card_level,
		clt_type,
		customer_code,
		customer_id,
		upgrade_need_amount,
		customer_type,
		customer_type_updated_at,
		vip_type,
		customer_deleted_at,
		age_tag,
		ba_code
) 
SELECT
		wuce.id AS id,
		wuce.unionid AS user_unionid,
		wuce.created_at AS created_at,
		wuce.updated_at AS updated_at,
		wuce.deleted_at AS deleted_at,
		wuce.openid AS openid,
		wuce.nick AS nick,
		wuce.name AS name,
		wuce.position AS position,
		wuce.avatar AS avatar,
		wuce.type AS type,
		wuce.sex AS sex,
		wuce.birthday AS birthday,
		wuce.union_code AS union_code,
		wuce.point AS point,
		wuce.grade AS grade,
		wuce.birth_year AS birth_year,
		wuce.birth_month AS birth_month,
		wuce.birth_day AS birth_day,
		wuce.mobile AS mobile,
		wuce.card_level AS card_level,
		wuce.clt_type AS clt_type,
		wuce.customer_code AS customer_code,
		wuce.customer_id AS customer_id,
		wuce.upgrade_need_amount AS upgrade_need_amount,
		wuce.customer_type AS customer_type,
		wuce.customer_type_updated_at AS customer_type_updated_at,
		wuce.vip_type AS vip_type,
		wuce.customer_deleted_at AS customer_deleted_at,
		wuce.age_tag AS age_tag,
		wuce.ba_code AS ba_code
	FROM
	   (SELECT * FROM micro_scrm_enterprise_%d_ck.wx_user_corp_external WHERE id NOT IN (select id from micro_scrm_enterprise_%d_ck.external_extra)) wuce`,
		enterpriseId, enterpriseId, enterpriseId)
	err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
	if err != nil {
		return err
	}
	fmt.Printf("补偿修复成功")
	return nil
}

// ExternalSearchListCk 根据条件筛选导购和外部联系人列表
func (d *ExternalExtraCk) ExternalSearchListCk(c ctx.Context, request customerSearchReq.ExternalSearchListBaseReq) ([]customerSearchResp.ExternalUserList, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserList, 0)
	query := d.GetExternalQuery(c, request.SearchList)
	if request.OnlyCount || !request.GetAll {
		err = query.Count(&total).Error
		if err != nil || request.OnlyCount || total == 0 {
			return list, total, err
		}
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

	externalExtraField := "id,ba_code,user_unionid,openid,nick,name,avatar,birthday," +
		"coupon_status,coupon_by_30days_status,valid_reserve_at,success_reserve_at,task_finish_at," +
		"card_level,amount_total,vip_type,last_purchase_date,clt_type,union_code"
	var extraList = make([]externalSearchModel.ExternalExtraModel, 0)
	err = query.Select(externalExtraField).Order("id desc").Limit(request.Size).Find(&extraList).Error
	if err != nil {
		return list, total, err
	}
	var externalIds []int64
	for _, row := range extraList {
		externalIds = append(externalIds, row.ID)
		tmp := customerSearchResp.ExternalUserList{}
		helper.CoverStruct(&tmp, &row)
		tmp.Unionid = row.UserUnionid
		tmp.ID = row.ID
		list = append(list, tmp)
	}
	if request.IsTaskList {
		return list, total, err
	}
	var relationShipList []externalSearchModel.ExternalRelationshipModel
	//获取备注和本人最近一次联系时间
	err = gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_relationship", request.EnterpriseID)).
		Where("member_id = ?", request.MemberId).
		Where("work_external_user_id in ?", externalIds).
		Select("remark,relation_status,work_external_user_id,unionid,active_time,interact_time").
		Find(&relationShipList).Error
	var relationShipMap = make(map[int64]externalSearchModel.ExternalRelationshipModel)
	for _, row := range relationShipList {
		relationShipMap[row.WorkExternalUserID] = row
	}
	var nowUnix = time.Now().Unix()
	for i, row := range list {
		//将结果提取到一个临时结构体
		list[i].SelfRelationStatus = externalSearchModel.RelationDelAll
		if exist, ok := relationShipMap[row.ID]; ok {
			list[i].Remark = exist.Remark
			list[i].ActiveTime = exist.ActiveTime
			list[i].InteractTime = exist.InteractTime
			list[i].SelfRelationStatus = exist.RelationStatus
		}
		if list[i].TaskFinishAt > 0 {
			if nowUnix-int64(list[i].TaskFinishAt) < 86400*customerSearchReq.TaskHaveValidDays {
				list[i].TaskFinishAtTag = true
			} else {
				list[i].TaskFinishAt = 0
			}
		}
		if list[i].ValidReserveAt > 0 {
			if nowUnix-int64(list[i].ValidReserveAt) < 86400*customerSearchReq.ValidDays {
				list[i].ValidReserveAtTag = true
			} else {
				list[i].ValidReserveAt = 0
			}
		}
	}
	return list, total, err
}

// GetExtensionBackfillBatch 按主键游标读取客户扩展字段，用于一次性回填企业 MySQL 扩展表。
// FINAL 用于避免 ReplacingMergeTree 尚未合并时同一 id 返回多个版本。
func (d *ExternalExtraCk) GetExtensionBackfillBatch(c app.Ctx, enterpriseId, lastId int64, limit int) ([]externalSearchModel.ExternalExtraModel, error) {
	list := make([]externalSearchModel.ExternalExtraModel, 0)
	if limit <= 0 {
		return list, fmt.Errorf("invalid extension backfill limit: %d", limit)
	}
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)
	sql := fmt.Sprintf(`SELECT
		id,user_unionid,valid_reserve_at,success_reserve_at,task_finish_at,task_created_at,
		coupon_status,coupon_by_30days_status,birthday_coupon_status,
		renew_membership_coupon_status,welcome_coupon_status,upgrade_coupon_status,
		coupon_updated_at,created_at,updated_at
	FROM %s FINAL
	WHERE id > ?
	ORDER BY id ASC
	LIMIT %d`, tableName, limit)
	err := gormx.GetCkInstance().WithContext(c).Raw(sql, lastId).Scan(&list).Error
	return list, err
}

// RelationShipListCk   根据条件筛选导购和外部联系人关系列表
func (d *ExternalExtraCk) RelationShipListCk(c ctx.Context, request customerSearchReq.ExternalSearchListBaseReq) ([]customerSearchResp.ExternalUserList, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]customerSearchResp.ExternalUserList, 0)
	query := d.GetRelationShipQuery(c, request.SearchList)
	if !request.GetAll {
		err = query.Count(&total).Error
		if err != nil || request.OnlyCount || total == 0 {
			return list, total, err
		}
		if request.Page > 0 && request.Size > 0 {
			query = query.Offset((request.Page - 1) * request.Size).Limit(request.Size)
		}
	}
	relationShipField := "id,remark,user_id,relation_status,work_external_user_id,unionid,active_time,interact_time"
	var relationShipList []externalSearchModel.ExternalRelationshipModel
	err = query.Select(relationShipField).Order("id desc").Find(&relationShipList).Error
	if err != nil {
		return list, total, err
	}
	var externalIds []int64
	for _, row := range relationShipList {
		externalIds = append(externalIds, row.WorkExternalUserID)
		tmp := customerSearchResp.ExternalUserList{}
		helper.CoverStruct(&tmp, &row)
		tmp.ID = row.ID
		list = append(list, tmp)
	}
	var extraList = make([]externalSearchModel.ExternalExtraModel, 0)
	externalExtraField := "ba_code,user_unionid,openid,nick,name,avatar,birthday," +
		"coupon_status,coupon_by_30days_status,valid_reserve_at,success_reserve_at,task_finish_at," +
		"card_level,amount_total,vip_type,last_purchase_date"
	//获取外部联系人信息
	err = gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)).
		Where("id in ?", externalIds).
		Select(externalExtraField).
		Find(&extraList).Error
	var extraMap = make(map[string]externalSearchModel.ExternalExtraModel)
	for _, row := range extraList {
		extraMap[row.UserUnionid] = row
	}
	for i, row := range list {
		//将结果提取到一个临时结构体
		if exist, ok := extraMap[row.Unionid]; ok {
			helper.CoverStruct(&row, &exist)
			list[i] = row
		}
	}
	return list, total, err
}

// ExternalSearchPluckIds   根据条件筛选导购和外部联系人
func (d *ExternalExtraCk) ExternalSearchPluckIds(c ctx.Context, searchReq customerSearchReq.SearchList, pluckName string) ([]interface{}, error) {
	query := d.GetExternalQuery(c, searchReq)
	var pluckIds = make([]interface{}, 0)
	err := query.Group(pluckName).Pluck(pluckName, &pluckIds).Error
	return pluckIds, err
}

// GetExternalQuery 顾客列表查询组装
func (d *ExternalExtraCk) GetExternalQuery(c ctx.Context, request customerSearchReq.SearchList) *gorm.DB {
	searchQueryResp := d.SearchExternalOptionQuery(request)
	relationshipTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_relationship", request.EnterpriseID)
	externalExtraTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)
	//顾客列表
	query := gormx.GetCkInstance().WithContext(c).Table(externalExtraTable).
		Where("type = ?", externalSearchModel.WeChat) // 只查微信用户
	if len(searchQueryResp.ExternalExtra.SqlSlice) > 0 {
		query = query.Where(searchQueryResp.ExternalExtra.SqlConn, searchQueryResp.ExternalExtra.SqlSlice...)
	}
	if len(searchQueryResp.RelationShip.SqlSlice) > 0 {
		query = query.Where(fmt.Sprintf("id in (SELECT work_external_user_id FROM %s WHERE %s)", relationshipTable, searchQueryResp.RelationShip.SqlConn),
			searchQueryResp.RelationShip.SqlSlice...)
	}
	return query
}

// GetRelationShipQuery  顾客关系列表查询组装
func (d *ExternalExtraCk) GetRelationShipQuery(c ctx.Context, request customerSearchReq.SearchList) *gorm.DB {
	searchQueryResp := d.SearchExternalOptionQuery(request)
	relationshipTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_relationship", request.EnterpriseID)
	externalExtraTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)
	//顾客关系列表
	query := gormx.GetCkInstance().WithContext(c).Table(relationshipTable).
		Where("unionid != ?", "") // 只查微信用户
	if len(searchQueryResp.RelationShip.SqlSlice) > 0 {
		query = query.Where(searchQueryResp.RelationShip.SqlConn, searchQueryResp.RelationShip.SqlSlice...)
	}
	if len(searchQueryResp.ExternalExtra.SqlSlice) > 0 {
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT id FROM %s WHERE %s)", externalExtraTable, searchQueryResp.ExternalExtra.SqlConn),
			searchQueryResp.ExternalExtra.SqlSlice...)
	}
	return query
}

// SearchExternalOptionQuery 组装筛选条件的sqlbuilder
func (d *ExternalExtraCk) SearchExternalOptionQuery(request customerSearchReq.SearchList) customerSearchResp.SearchQueryResp {
	unixTime := time.Now().Unix()
	relationshipTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_relationship", request.EnterpriseID)
	externalExtraTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", request.EnterpriseID)

	var relationshipSlice []interface{}
	var relationshipConn = "deleted_at = 0"
	var externalExtraSlice []interface{}
	var externalExtraConn = "deleted_at = 0"

	if len(request.MemberIds) == 0 { //如果没有圈定员工范围，则范围就是自己
		request.MemberIds = []int64{request.MemberId}
	}
	//好友非好友，如果没有选择好友关系则说明本身已经有顾客范围了，不需要根据员工去圈定
	if request.RelationshipType == customerSearchReq.IsUnionFalse {
		externalExtraConn += " AND id NOT IN (SELECT work_external_user_id FROM " + relationshipTable + " WHERE member_id in (?) AND relation_status = ? )"
		externalExtraSlice = append(externalExtraSlice, request.MemberIds, 0)
	} else if request.RelationshipType == customerSearchReq.IsUnionTure {
		relationshipConn += " AND member_id IN (?) AND relation_status = ? "
		relationshipSlice = append(relationshipSlice, request.MemberIds, 0)
	}
	if request.IsSelf {
		//仅看我
		relationshipConn += " AND member_id = ?"
		relationshipSlice = append(relationshipSlice, request.MemberId)
	}
	if len(request.UnionIds) > 0 {
		externalExtraConn += " AND user_unionid  in ?"
		externalExtraSlice = append(externalExtraSlice, request.UnionIds)
	}
	if len(request.UnionCodes) > 0 {
		externalExtraConn += " AND union_code  in ?"
		externalExtraSlice = append(externalExtraSlice, request.UnionCodes)
	}
	if len(request.FilterExternalUserId) > 0 {
		externalExtraConn += " AND id not in ?"
		externalExtraSlice = append(externalExtraSlice, request.FilterExternalUserId)
	}
	if len(request.FilterRelationId) > 0 {
		relationshipConn += " AND id not in ?"
		relationshipSlice = append(relationshipSlice, request.FilterRelationId)
	}
	if len(request.ExternalUserIds) > 0 {
		externalExtraConn += " AND id in ?"
		externalExtraSlice = append(externalExtraSlice, request.ExternalUserIds)
	}
	if len(request.RelationId) > 0 {
		relationshipConn += " AND id in ?"
		relationshipSlice = append(relationshipSlice, request.RelationId)
	}
	if request.Name != "" {
		relationshipConn += " AND ((member_id = ? AND LOWER(remark) LIKE LOWER(?)) OR work_external_user_id GLOBAL IN (SELECT id FROM " + externalExtraTable + " WHERE LOWER(nick) LIKE LOWER(?) ))"
		relationshipSlice = append(relationshipSlice, request.MemberId, "%"+request.Name+"%", "%"+request.Name+"%")
	}
	// 会员卡等
	if len(request.CardType) > 0 {
		var cardLevel []string
		for _, v := range request.CardType {
			if v == "O" { // 非会员 或为空
				cardLevel = append(cardLevel, "O", "")
			} else {
				cardLevel = append(cardLevel, v)
			}
		}
		if len(cardLevel) > 0 {
			externalExtraConn += " AND card_level  in ?"
			externalExtraSlice = append(externalExtraSlice, cardLevel)
		}
	}
	//生日月份
	if len(request.BirthMonth) > 0 {
		externalExtraConn += " AND birth_month  in ?"
		externalExtraSlice = append(externalExtraSlice, request.BirthMonth)
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
	//有无任务查询的是7天的
	taskHaveValidTime := time.Now().Unix() - customerSearchReq.TaskHaveValidDays*86400
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
	//优惠券筛选
	if len(request.CouponType) > 0 {
		externalExtraConn += " AND ("
		for _, v := range request.CouponType {
			switch v {
			case customerSearchReq.Has30DayCoupon:
				externalExtraConn += " coupon_by_30days_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.HasCoupon:
				externalExtraConn += " coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.BirthCoupon:
				externalExtraConn += " birthday_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.WelcomeCoupon:
				externalExtraConn += " welcome_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.RenewCoupon:
				externalExtraConn += " renew_membership_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
			case customerSearchReq.UpgradeCoupon:
				externalExtraConn += " upgrade_coupon_status = ? OR"
				externalExtraSlice = append(externalExtraSlice, externalSearchModel.HasCoupon)
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
		externalExtraConn += ",purchased_date)"
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
		externalExtraConn += " AND last_purchase_date < ?"
		externalExtraSlice = append(externalExtraSlice, tmpTime)
	}
	//历史消费总金额开始
	if request.HistoryAmountTotal.StartNumber > 0 {
		externalExtraConn += " AND history_amount_total >= ?"
		externalExtraSlice = append(externalExtraSlice, request.HistoryAmountTotal.StartNumber)
	}
	//历史消费总金额结束
	if request.HistoryAmountTotal.EndNumber > 0 {
		externalExtraConn += " AND history_amount_total <= ?"
		externalExtraSlice = append(externalExtraSlice, request.HistoryAmountTotal.EndNumber)
	}
	//L12M总金额开始
	if request.AmountTotal.StartNumber > 0 {
		externalExtraConn += " AND amount_total >= ?"
		externalExtraSlice = append(externalExtraSlice, request.AmountTotal.StartNumber)
	}
	//L12M总金额结束
	if request.AmountTotal.EndNumber > 0 {
		externalExtraConn += " AND amount_total <= ?"
		externalExtraSlice = append(externalExtraSlice, request.AmountTotal.EndNumber)
	}
	//添加企微好友时间这里有本人和企业两种模式，如果是本人则可以直接查，如果是企业则需要单独查符合条件的人之后重新in
	if len(request.RelCreateTime) > 0 {
		var relAddTimeConn string
		var relAddTimeSlice []interface{}
		relAddTimeConn += " AND ("
		for _, v := range request.RelCreateTime {
			switch v {
			case customerSearchReq.TimeToday:
				todayZeroTimestamp, _ := helper.PeriodOfTimestamp()
				relAddTimeConn += " createtime > ? OR"
				relAddTimeSlice = append(relAddTimeSlice, todayZeroTimestamp)
			case customerSearchReq.TimeYesterday:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				yesterdayStart, yesterdayEnd := helper.YesterdayPeriodOfTimestamp(1)
				relAddTimeSlice = append(relAddTimeSlice, yesterdayStart, yesterdayEnd)
			case customerSearchReq.Time2To7D:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-7*86400, unixTime-2*86400)
			case customerSearchReq.Time8To15D:
				relAddTimeConn += " (createtime BETWEEN ? AND ?) OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-15*86400, unixTime-8*86400)
			case customerSearchReq.TimeGt15D:
				relAddTimeConn += " createtime < ? OR"
				relAddTimeSlice = append(relAddTimeSlice, unixTime-15*86400)
			}
		}
		relAddTimeConn = strings.TrimSuffix(relAddTimeConn, "OR")
		relAddTimeConn += ") "
		if request.RelCreateTimeType == customerSearchReq.Self {
			//如果是本人则直接查
			relationshipConn += relAddTimeConn
			relationshipSlice = append(relationshipSlice, relAddTimeSlice...)
			relationshipConn += " AND member_id = ?"
			relationshipSlice = append(relationshipSlice, request.MemberId)
		} else if request.RelCreateTimeType == customerSearchReq.Enterprise {
			//如果是企业则需要单独查出相关的unionid 单独in
			relationshipConn += fmt.Sprintf(" AND work_external_user_id IN (SELECT work_external_user_id FROM %s WHERE deleted_at = 0 %s)", relationshipTable, relAddTimeConn)
			relationshipSlice = append(relationshipSlice, relAddTimeSlice...)
		}

	}
	//上次触达时间分为本人和企业两种模式，如果是本人则可以直接查，如果是企业则需要单独查符合条件的人之后重新in
	if len(request.RelInteractTime) > 0 {
		var interactTimeConn string
		var interactTimeSlice []interface{}
		interactTimeConn += " AND ("
		for _, v := range request.RelInteractTime {
			switch v {
			case customerSearchReq.TimeToday:
				todayZeroTimestamp, _ := helper.PeriodOfTimestamp()
				interactTimeConn += " interact_time > ? OR"
				interactTimeSlice = append(interactTimeSlice, todayZeroTimestamp)
			case customerSearchReq.TimeYesterday:
				interactTimeConn += " (interact_time BETWEEN ? AND ?) OR"
				yesterdayStart, yesterdayEnd := helper.YesterdayPeriodOfTimestamp(1)
				interactTimeSlice = append(interactTimeSlice, yesterdayStart, yesterdayEnd)
			case customerSearchReq.Time2To7D:
				interactTimeConn += " (interact_time BETWEEN ? AND ?) OR"
				interactTimeSlice = append(interactTimeSlice, unixTime-7*86400, unixTime-2*86400)
			case customerSearchReq.Time8To15D:
				interactTimeConn += " (interact_time BETWEEN ? AND ?) OR"
				interactTimeSlice = append(interactTimeSlice, unixTime-15*86400, unixTime-8*86400)
			case customerSearchReq.TimeGt15D:
				interactTimeConn += " interact_time < ? OR"
				interactTimeSlice = append(interactTimeSlice, unixTime-15*86400)
			}
		}
		interactTimeConn = strings.TrimSuffix(interactTimeConn, "OR")
		interactTimeConn += ") "
		if request.RelInteractTimeType == customerSearchReq.Self {
			//如果是本人则直接查
			relationshipConn += interactTimeConn
			relationshipSlice = append(relationshipSlice, interactTimeSlice...)
			relationshipConn += " AND member_id = ?"
			relationshipSlice = append(relationshipSlice, request.MemberId)
		} else if request.RelInteractTimeType == customerSearchReq.Enterprise {
			//如果是企业则需要单独查出相关的unionid 单独in
			relationshipConn += fmt.Sprintf(" AND work_external_user_id IN (SELECT work_external_user_id FROM %s WHERE deleted_at = 0 %s)", relationshipTable, interactTimeConn)
			relationshipSlice = append(relationshipSlice, interactTimeSlice...)
		}

	}

	//活跃沟通时间分为本人和企业两种模式，如果是本人则可以直接查，如果是企业则需要单独查符合条件的人之后重新in
	if len(request.RelActiveTime) > 0 {
		var activeTimeConn string
		var activeTimeSlice []interface{}
		activeTimeConn += " AND ("
		for _, v := range request.RelActiveTime {
			switch v {
			case customerSearchReq.TimeLt7D:
				activeTimeConn += " active_time >= ? OR"
				activeTimeSlice = append(activeTimeSlice, unixTime-7*86400)
			case customerSearchReq.Time8To15D:
				activeTimeConn += " (active_time BETWEEN ? AND ?) OR"
				activeTimeSlice = append(activeTimeSlice, unixTime-15*86400, unixTime-8*86400)
			case customerSearchReq.Time16To30D:
				activeTimeConn += " (active_time BETWEEN ? AND ?) OR"
				activeTimeSlice = append(activeTimeSlice, unixTime-30*86400, unixTime-16*86400)
			case customerSearchReq.TimeGt30D:
				activeTimeConn += " active_time <= ? OR"
				activeTimeSlice = append(activeTimeSlice, unixTime-30*86400)
			}
		}
		activeTimeConn = strings.TrimSuffix(activeTimeConn, "OR")
		activeTimeConn += ") "
		if request.RelActiveTimeType == customerSearchReq.Self {
			//如果是本人则直接查
			relationshipConn += activeTimeConn
			relationshipSlice = append(relationshipSlice, activeTimeSlice...)
			relationshipConn += " AND member_id = ?"
			relationshipSlice = append(relationshipSlice, request.MemberId)
		} else if request.RelActiveTimeType == customerSearchReq.Enterprise {
			//如果是企业则需要单独查出相关的unionid 单独in
			relationshipConn += fmt.Sprintf(" AND work_external_user_id IN (SELECT work_external_user_id FROM %s WHERE deleted_at = 0 %s)", relationshipTable, activeTimeConn)
			relationshipSlice = append(relationshipSlice, activeTimeSlice...)
		}

	}
	// 顾客沟通意愿标签，只查看自己打的标签数据
	if request.ChatTag > 0 {
		relationshipConn += " AND has(external_tag_ids,?) AND member_id = ?"
		relationshipSlice = append(relationshipSlice, request.ChatTag, request.MemberId)
	}
	//高意愿用户
	if request.HighWilling != "" {
		externalExtraConn += " AND high_willing = ?"
		externalExtraSlice = append(externalExtraSlice, request.HighWilling)
	}
	//价格敏感度
	if request.PriceSensitivity != "" {
		externalExtraConn += " AND price_sensitivity = ?"
		externalExtraSlice = append(externalExtraSlice, request.PriceSensitivity)
	}
	//商品偏好
	if len(request.CategoryPreference) > 0 {
		externalExtraConn += " AND ("
		for _, v := range request.CategoryPreference {
			if categorySlice, ok := customerSearchReq.CdpCategoryMap[v]; ok {
				for _, tag := range categorySlice {
					externalExtraConn += " has(category_preference,?) OR"
					externalExtraSlice = append(externalExtraSlice, tag)
				}
			}
		}
		externalExtraConn = strings.TrimSuffix(externalExtraConn, "OR")
		externalExtraConn += ") "
	}
	// 任务触达查询的数据有效时间是30天内的
	// coach的AgeTag 顾客类型标签
	if request.AgeTag > 0 {
		externalExtraConn += " AND age_tag = ?"
		externalExtraSlice = append(externalExtraSlice, customerSearchReq.AgeMap[request.AgeTag])
	}
	return customerSearchResp.SearchQueryResp{
		ExternalExtra: customerSearchResp.BuildSql{
			SqlConn:  externalExtraConn,
			SqlSlice: externalExtraSlice,
		},
		RelationShip: customerSearchResp.BuildSql{
			SqlConn:  relationshipConn,
			SqlSlice: relationshipSlice,
		},
	}
}

// UpdateExternalExtraCkByCdpTag     每天全量更新CDP标签数据
func (d *ExternalExtraCk) UpdateExternalExtraCkByCdpTag(c app.Ctx, enterpriseId int64, pageSize int) error {
	//获取各个表更新数据相关的重要id
	var maxId int64
	extraTableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)
	cdpTagTableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, "cdp_external_tag")
	for {
		// 拼接一个extra表的部分查询语句
		extraPageSql := fmt.Sprintf(" FROM %s WHERE id > %d order by id asc LIMIT %d",
			extraTableName, maxId, pageSize)
		sql := fmt.Sprintf(`INSERT INTO %s (
		id,
		user_unionid,
		valid_reserve_at,
		success_reserve_at,
		task_finish_at,
		task_created_at,
		coupon_status,
		coupon_by_30days_status,
		coupon_updated_at,
		created_at,
		updated_at,
		deleted_at,
		birthday_coupon_status,
		renew_membership_coupon_status,
		welcome_coupon_status,
		upgrade_coupon_status,
		openid,
		nick,
		name,
		position,
		avatar,
		type,
		sex,
		birthday,
		union_code,
		point,
		grade,
		birth_year,
		birth_month,
		birth_day,
		mobile,
		card_level,
		clt_type,
		amount_total,
		history_amount_total,
		customer_code,
		customer_id,
		last_purchase_date,
		last_purchase_store_code,
		upgrade_need_amount,
		customer_type,
		customer_type_updated_at,
		vip_type,
		customer_deleted_at,
		age_tag,
		category_preference,
		purchased_date,
		external_impression_tag_ids,
		high_willing,
		price_sensitivity,
		ba_code
) SELECT
		extra.id AS id,
		any(extra.user_unionid) AS user_unionid,
		any(extra.valid_reserve_at) AS valid_reserve_at,
		any(extra.success_reserve_at) AS success_reserve_at,
		any(extra.task_finish_at) AS task_finish_at,
		any(extra.task_created_at) AS task_created_at,
		any(extra.coupon_status) AS coupon_status,
		any(extra.coupon_by_30days_status) AS coupon_by_30days_status,
		any(extra.coupon_updated_at) AS coupon_updated_at,
		any(extra.created_at) AS created_at,
		any(extra.updated_at) AS updated_at,
		any(extra.deleted_at) AS deleted_at,
		any(extra.birthday_coupon_status) AS birthday_coupon_status,
		any(extra.renew_membership_coupon_status) AS renew_membership_coupon_status,
		any(extra.welcome_coupon_status) AS welcome_coupon_status,
		any(extra.upgrade_coupon_status) AS upgrade_coupon_status,
		any(extra.openid) AS openid,
		any(extra.nick) AS nick,
		any(extra.name) AS name,
		any(extra.position) AS position,
		any(extra.avatar) AS avatar,
		any(extra.type) AS type,
		any(extra.sex) AS sex,
		any(extra.birthday) AS birthday,
		any(extra.union_code) AS union_code,
		any(extra.point) AS point,
		any(extra.grade) AS grade,
		any(extra.birth_year) AS birth_year,
		any(extra.birth_month) AS birth_month,
		any(extra.birth_day) AS birth_day,
		any(extra.mobile) AS mobile,
		any(extra.card_level) AS card_level,
		any(extra.clt_type) AS clt_type,
		anyIf(toInt64OrZero(tag.tag_value),tag.tag_type = '`+externalSearchModel.CdpTagAmountTotal+`') AS amount_total,
		anyIf(toInt64OrZero(tag.tag_value), tag.tag_type = '`+externalSearchModel.CdpTagHistoryAmountTotal+`') AS history_amount_total,
		any(extra.customer_code) AS customer_code,
		any(extra.customer_id) AS customer_id,
		anyIf(toInt64OrZero(tag.tag_value), tag.tag_type = '`+externalSearchModel.CdpTagLastPurchaseDate+`') AS last_purchase_date,
		anyIf(tag.tag_value, tag.tag_type = '`+externalSearchModel.CdpTagLastPurchaseStoreCode+`') AS last_purchase_store_code,
		any(extra.upgrade_need_amount) AS upgrade_need_amount,
		any(extra.customer_type) AS customer_type,
		any(extra.customer_type_updated_at) AS customer_type_updated_at,
		any(extra.vip_type) AS vip_type,
		any(extra.customer_deleted_at) AS customer_deleted_at,
		any(extra.age_tag) AS age_tag,
		anyIf(splitByChar('|', tag.tag_value), tag.tag_type = '`+externalSearchModel.CdpTagCategoryPreference+`') AS category_preference,
		anyIf(arrayMap(x -> toInt64OrZero(x), splitByChar('|', tag.tag_value)), tag.tag_type = '`+externalSearchModel.CdpTagPurchasedDate+`') AS purchased_date,
		any(extra.external_impression_tag_ids) AS external_impression_tag_ids,
		anyIf(tag.tag_value, tag.tag_type = '`+externalSearchModel.CdpTagHighWilling+`') AS high_willing,
		anyIf(tag.tag_value, tag.tag_type = '`+externalSearchModel.CdpTagPriceSensitivity+`') AS price_sensitivity,
		any(extra.ba_code) AS ba_code
	FROM
	   (SELECT * %s) extra
	   GLOBAL LEFT JOIN (SELECT * FROM %s WHERE unionid IN (SELECT user_unionid %s)) tag ON extra.user_unionid = tag.unionid
	GROUP BY
	   extra.id`,
			extraTableName+"_local", extraPageSql, cdpTagTableName, extraPageSql)
		err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
		if err != nil {
			return err
		}
		// 获取每pageSize条数据中的最大id,每次循环更新最大id,将数据join关联更新范围控制在pageSize条数据
		_ = gormx.GetCkInstance().WithContext(c).
			Raw(fmt.Sprintf("SELECT MAX(id) FROM (select id from micro_scrm_enterprise_%d_ck.%s where id > ? order by id asc limit ?)",
				enterpriseId, d.tableName), maxId, pageSize).
			Find(&maxId).Error
		if maxId == 0 {
			// 如果没查出来数据就结束循环
			return nil
		}
	}
}

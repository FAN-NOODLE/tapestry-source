package externalSearchDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type ExternalRelationshipCk struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewExternalRelationshipCk() *ExternalRelationshipCk {
	return &ExternalRelationshipCk{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      externalSearchModel.ExternalRelationshipModel{}.TableName(),
		localTableName: externalSearchModel.ExternalRelationshipModel{}.TableName() + "_local",
	}
}

func (d *ExternalRelationshipCk) BatchSaveActiveData(c app.Ctx, req openReq.BatchUpdateActiveReq) ([]externalSearchModel.ExternalRelationshipModel, error) {
	unixTime := time.Now().Unix()
	var list []externalSearchModel.ExternalRelationshipModel
	var whereIn [][]interface{}
	for _, v := range req.DataList {
		whereIn = append(whereIn, []interface{}{
			v.UnionId, v.UserId,
		})
	}
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseId, d.tableName)).
		Where("(unionid,user_id) in ?", whereIn).Find(&list).Error
	if err != nil {
		return nil, err
	}
	// 整理已经存在的数据map,以Unionid为key
	var existMap = make(map[string]externalSearchModel.ExternalRelationshipModel)
	for _, exist := range list {
		existMap[exist.Unionid+exist.UserID] = exist
	}
	var newData []externalSearchModel.ExternalRelationshipModel // 待新增的数据
	// 整理提交过来的数据，提交过来的数据有可能是重复的，所以
	for _, row := range req.DataList {
		if oldInfo, ok := existMap[row.UnionId+row.UserId]; ok {
			if row.ActiveTime > 0 {
				if row.ActiveTime < oldInfo.ActiveTime {
					continue
				}
				oldInfo.ActiveTime = row.ActiveTime
			}
			if row.InteractTime > 0 {
				if row.InteractTime < oldInfo.InteractTime {
					continue
				}
				oldInfo.InteractTime = row.InteractTime
			}
			oldInfo.UpdatedAt = unixTime
			newData = append(newData, oldInfo)
		}
	}
	if len(newData) > 0 {
		// 批量插入数据
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseId, d.localTableName)
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(tableName)
			}).Create(&newData).Error
		})
		if err != nil {
			return nil, err
		}
		log.Info(c).Msgf("【优惠券】顾客列表附加信息更新，任务相关，插入条数:%d\n", len(newData))
	}

	return list, nil
}

// MergeData  ck数据合并
func (d *ExternalRelationshipCk) MergeData(c app.Ctx, enterpriseId int64) error {
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	return gormx.GetCkInstance().WithContext(c).Debug().Exec(fmt.Sprintf("OPTIMIZE TABLE %s FINAL DEDUPLICATE", tableName)).Error
}

// UpdateExternalRelationshipCk     每3分钟将更新的数据写入数据
func (d *ExternalRelationshipCk) UpdateExternalRelationshipCk(c app.Ctx, enterpriseId int64, intervalInSeconds int64) error {
	unixTime := time.Now().Unix() - intervalInSeconds
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	relationTableName := "micro_scrm_ck.work_external_user_relation"
	tagTableName := "micro_scrm_ck.work_external_user_tag"
	// 拼接一个extra表的部分查询语句
	relationIdsSql := fmt.Sprintf(`WITH relation_ids AS (
    SELECT id
    FROM (
        SELECT id
        FROM %s 
        WHERE updated_at >= %d 
        UNION ALL 
        SELECT relation_id as id
        FROM %s 
        WHERE updated_at >= %d
    ) AS extra 
    GROUP BY id
)`, relationTableName, unixTime, tagTableName, unixTime)

	sql := fmt.Sprintf(`INSERT INTO %s (
		id,
		work_external_user_id,
		unionid,
		member_id,
		user_id,
		created_at,
		updated_at,
		deleted_at,
		remark,
		description,
		createtime,
		remark_mobiles,
		add_way,
		state,
		relation_status,
		member_deleted_time,
		external_deleted_time,
		is_key_sa,
		interact_time,
		active_time,
		external_tag_ids
) 
%s
SELECT
		weur.id AS id,
		any(weur.work_external_user_id) AS work_external_user_id,
		any(weur.unionid) AS unionid,
		any(weur.member_id) AS member_id,
		any(weur.user_id) AS user_id,
		any(weur.created_at) AS created_at,
		any(weur.updated_at) AS updated_at,
		any(weur.deleted_at) AS deleted_at,
		any(weur.remark) AS remark,
		any(weur.description) AS description,
		any(weur.createtime) AS createtime,
		any(weur.remark_mobiles) AS remark_mobiles,
		any(weur.add_way) AS add_way,
		any(weur.state) AS state,
		any(weur.relation_status) AS relation_status,
		any(weur.member_deleted_time) AS member_deleted_time,
		any(weur.external_deleted_time) AS external_deleted_time,
		any(weur.is_key_sa) AS is_key_sa,
		any(relationship.interact_time) AS interact_time,
		any(relationship.active_time) AS active_time,
		arrayFilter(x -> x != 0, arrayDistinct(groupArray(weut.external_tag_id))) AS external_tag_ids
	FROM
	   (SELECT * FROM micro_scrm_ck.work_external_user_relation WHERE enterprise_id=%d AND id IN (SELECT id FROM relation_ids)) weur
	   GLOBAL LEFT JOIN (SELECT * FROM micro_scrm_enterprise_%d_ck.external_relationship WHERE id IN (SELECT id FROM relation_ids)) relationship ON relationship.id = weur.id
	   GLOBAL LEFT JOIN (SELECT * FROM micro_scrm_ck.work_external_user_tag WHERE deleted_at=0 AND enterprise_id=%d AND relation_id IN (SELECT id FROM relation_ids)) weut ON weur.id = weut.relation_id
	GROUP BY
	   weur.id`, tableName, relationIdsSql, enterpriseId, enterpriseId, enterpriseId)
	err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
	if err != nil {
		return err
	}
	fmt.Printf("顾客关系大宽表更新成功，更新时间范围：%d分\n", intervalInSeconds/60)
	return nil
}

func (d *ExternalRelationshipCk) appendIds(allData externalSearchModel.MaterializeMainData, idsData []externalSearchModel.MaterializeMainId) externalSearchModel.MaterializeMainData {
	for _, data := range idsData {
		if data.RelationId > 0 {
			allData.RelationIds = append(allData.RelationIds, data.RelationId)
		}
	}
	allData.RelationIds = helper.ArrayUniq(allData.RelationIds)
	return allData
}

func (d *ExternalRelationshipCk) GetList(c app.Ctx, enterpriseId, memberId int64, workExternalUserIds []int64) ([]externalSearchModel.ExternalRelationshipModel, error) {
	list := make([]externalSearchModel.ExternalRelationshipModel, 0)
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("work_external_user_id in ? and member_id=?", workExternalUserIds, memberId).Find(&list).Error
	return list, err
}

func (d *ExternalRelationshipCk) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (info externalSearchModel.ExternalRelationshipModel, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&info).Error
	return
}

// GetExtensionBackfillBatch 按主键游标读取关系扩展字段，用于一次性回填企业 MySQL 扩展表。
// FINAL 用于避免 ReplacingMergeTree 尚未合并时同一 id 返回多个版本。
func (d *ExternalRelationshipCk) GetExtensionBackfillBatch(c app.Ctx, enterpriseId, lastId int64, limit int) ([]externalSearchModel.ExternalRelationshipModel, error) {
	list := make([]externalSearchModel.ExternalRelationshipModel, 0)
	if limit <= 0 {
		return list, fmt.Errorf("invalid extension backfill limit: %d", limit)
	}
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)
	sql := fmt.Sprintf(`SELECT id,unionid,user_id,interact_time,active_time,created_at,updated_at
	FROM %s FINAL
	WHERE id > ?
	ORDER BY id ASC
	LIMIT %d`, tableName, limit)
	err := gormx.GetCkInstance().WithContext(c).Raw(sql, lastId).Scan(&list).Error
	return list, err
}

// FixRelationShip   修复数据同步脚本可能遗漏的数据
func (d *ExternalRelationshipCk) FixRelationShip(c app.Ctx, enterpriseId int64) error {
	sql := fmt.Sprintf(`INSERT INTO  micro_scrm_enterprise_%d_ck.external_relationship_local (
		id,
		work_external_user_id,
		unionid,
		member_id,
		user_id,
		created_at,
		updated_at,
		deleted_at,
		remark,
		description,
		createtime,
		remark_mobiles,
		add_way,
		state,
		relation_status,
		member_deleted_time,
		external_deleted_time,
		is_key_sa,
		interact_time,
		active_time,
		external_tag_ids
) 
SELECT
    weur.id AS id,
    weur.work_external_user_id AS work_external_user_id,
    weur.unionid AS unionid,
    weur.member_id AS member_id,
    weur.user_id AS user_id,
    weur.created_at AS created_at,
    weur.updated_at AS updated_at,
    weur.deleted_at AS deleted_at,
    weur.remark AS remark,
    weur.description AS description,
    weur.createtime AS createtime,
    weur.remark_mobiles AS remark_mobiles,
    weur.add_way AS add_way,
    weur.state AS state,
    weur.relation_status AS relation_status,
    weur.member_deleted_time AS member_deleted_time,
    weur.external_deleted_time AS external_deleted_time,
    weur.is_key_sa AS is_key_sa,
    0,
    0,
    '[]' 
	FROM
	   (SELECT * FROM micro_scrm_ck.work_external_user_relation WHERE enterprise_id=%d AND id NOT IN (select id from micro_scrm_enterprise_%d_ck.external_relationship)) weur`,
		enterpriseId, enterpriseId, enterpriseId)
	err := gormx.GetCkInstance().WithContext(c).Exec(sql).Error
	if err != nil {
		return err
	}
	fmt.Printf("关系补偿修复成功")
	return nil
}

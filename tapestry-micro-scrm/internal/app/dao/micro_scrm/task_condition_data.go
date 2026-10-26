package microScrmDao

import (
	"context"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
	"time"
)

type PlanTaskConditionDataDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskConditionDataDao() *PlanTaskConditionDataDao {
	return &PlanTaskConditionDataDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.TaskConditionDataModel{}.TableName(),
	}
}

// BranchSharedTable 分表
func (d *PlanTaskConditionDataDao) BranchSharedTable(table string, enterpriseID, planId int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseID, table, planId)
		return tx.Table(tableName)
	}
}

// InsertConditionData
func (d *PlanTaskConditionDataDao) InsertConditionData(c app.Ctx, enterpriseID int64, planId int64, conditionModels []microScrmModel.TaskConditionDataModel) error {
	return d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseID, planId)).Create(&conditionModels).Error
}

// CreateConditionData
func (d *PlanTaskConditionDataDao) CreateConditionData(c app.Ctx, enterpriseID int64, planId int64, conditionModels microScrmModel.TaskConditionDataModel) error {
	return d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseID, planId)).Create(&conditionModels).Error
}

// GetCustomerLists 获取sacode&&unionid|customerid的顾客列表
func (d *PlanTaskConditionDataDao) GetCustomerLists(c ctx.Context, request planReq.CustomerListsReq, planId int64, relationUnionIds, searchNickUnionid []string) (total int64, result []microScrmModel.TaskConditionDataModel, err error) {
	offset := (request.Page - 1) * request.Size
	size := request.Size
	orderStr := "is_touch desc,is_finished asc,union_id asc"
	if request.TaskCategory == microScrmModel.PLANTASKMOMENT {
		orderStr = "union_id asc"
	}
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, request.EnterpriseId, planId)).Where("task_allot_log_id=?", request.Id).Where("is_match=?", 1).Scopes(func(db *gorm.DB) *gorm.DB {
		if request.FilterType > 0 {
			if request.FilterType == 1 {
				db.Where("is_finished=?", request.FilterType)
			} else {
				db.Where("is_finished!=?", 1)
			}
		}
		if request.IsUnion > 0 {
			if request.IsUnion == 1 {
				db.Where("union_id in (?)", relationUnionIds)
			} else {
				db.Where("union_id not in (?)", relationUnionIds)
			}
		}
		if len(searchNickUnionid) > 0 {
			db.Where("union_id in (?)", searchNickUnionid)
		}
		//if request.Keywords != "" {
		//	keywords := strings.Replace(request.Keywords, "%", "\\%", -1)
		//	db.Where("(remark like ? or user_name like ? )", "%"+keywords+"%", "%"+keywords+"%")
		//}
		return db
	}).Group("union_id,customer_id").
		Count(&total).Order(orderStr).
		Offset(offset).
		Limit(size).
		Find(&result).Error

	return
}

func (d *PlanTaskConditionDataDao) UpdateData(c app.Ctx, enterpriseId, planId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *PlanTaskConditionDataDao) GetOne(c app.Ctx, enterpriseId, planId int64, query interface{}) (microScrmModel.TaskConditionDataModel, error) {
	var result microScrmModel.TaskConditionDataModel
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query).Find(&result).Error
	return result, err
}

func (d *PlanTaskConditionDataDao) GetCount(c ctx.Context, enterpriseId, planId int64, query interface{}, args ...interface{}) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query, args...).Count(&total).Error
	return
}

func (d *PlanTaskConditionDataDao) GetGroupCount(c ctx.Context, enterpriseId, planId int64, query interface{}, args ...interface{}) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query, args...).Group("union_id,customer_id").Count(&total).Error
	return
}

// GetConditionTagsByUnionIds 获取unionid的标签数据
func (d *PlanTaskConditionDataDao) GetConditionTagsByUnionIds(c ctx.Context, enterpriseId, planId, tid int64, unionIds []string) (data []planReq.ConditionDataTagsList, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).
		Select("id,tag_name,union_id,ww_tag_id").
		Where(map[string]interface{}{
			"task_allot_log_id": tid,
			"union_id":          unionIds,
		}).Find(&data).Error
	return
}

func (d *PlanTaskConditionDataDao) GetListsByPage(c ctx.Context, enterpriseId, planId int64, page, size int, ids []int64) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			db.Where("task_allot_log_id in ?", ids)
		}
		return db
	}).Offset(offset).Limit(size).Find(&result).Error
	return result, err
}

// GetListsByTaskId 获取MA任务结束后未标记联系的顾客信息给MA推送回执
func (d *PlanTaskConditionDataDao) GetListsByTaskId(c ctx.Context, enterpriseId, planId, taskId int64, page, size int) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where("task_id=? and is_finished=0", taskId).Offset(offset).Limit(size).Find(&result).Error
	return result, err
}

// GetMany 获取列表数据
func (d *PlanTaskConditionDataDao) GetMany(c app.Ctx, enterpriseId, planId int64, query interface{}, args ...interface{}) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query, args...).Find(&result).Error
	return result, err
}

// GetOpenIds  获取openid列表
func (d *PlanTaskConditionDataDao) GetOpenIds(c app.Ctx, enterpriseId, planId int64, query interface{}, args ...interface{}) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).
		Where(query, args...).Pluck("work_external_user_id", &result).Error
	return result, err
}

// GetDistinctUserCount 获取去重客户数
func (d *PlanTaskConditionDataDao) GetDistinctUserCount(c ctx.Context, enterpriseId, planId int64, query interface{}, args ...interface{}) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query, args...).Select("count(distinct(CONCAT(union_id,customer_id)))").Count(&total).Error
	return
}

// GetNoFinishedPlanIdsByTaskId 查询出已完成的企业群发的沟通任务有未触达用户的主任务id
func (d *PlanTaskConditionDataDao) GetNoFinishedPlanIdsByTaskId(c ctx.Context, enterpriseId, planId int64, userid string, taskIds []int64, startTime, endTime int64) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(map[string]interface{}{
		"user_id":     userid,
		"task_id":     taskIds,
		"is_finished": 0,
	}).Where("created_at BETWEEN ? AND ? ", startTime, endTime).Count(&total).Error
	return total, err
}

func (d *PlanTaskConditionDataDao) GetNoFinishedCountByIds(c ctx.Context, enterpriseId, planId int64, ids []int64) ([]planReq.TaskNoFinishedCount, error) {
	var result []planReq.TaskNoFinishedCount

	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Select("task_allot_log_id as id ,count(DISTINCT(union_id)) as num").Where(map[string]interface{}{
		"task_allot_log_id": ids,
		"is_finished":       0,
	}).Group("task_allot_log_id").Scan(&result).Error
	return result, err
}

func (d *PlanTaskConditionDataDao) GetNoTouchLists(c ctx.Context, enterpriseId, planId int64, taskId int64) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(map[string]interface{}{
		"task_id":     taskId,
		"is_touch":    microScrmModel.TouchStatusLimit,
		"is_finished": 0,
	}).Find(&result).Error
	return result, err
}

// CreateShardTable 根据主任务创建task_condition_data分表
func (d *PlanTaskConditionDataDao) CreateShardTable(c ctx.Context, enterpriseId int64, planId int64) error {
	tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName, planId)
	createTableSql := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  task_allot_log_id int(10) NOT NULL,
  tag_name varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT 'CampaignTool标签名称',
  ww_tag_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企微标签id',
  union_id char(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '顾客的UnionID',
  work_external_user_id char(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '外部联系人user_id,可能空',
  is_union tinyint(2) NOT NULL DEFAULT '0' COMMENT '关系状态；0没有好友关系、1有好友关系',
  mobile varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机',
  mars_id varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'BA的MarsCode',
  user_id varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '员工的企微user_id',
  condition_id bigint(20) NOT NULL DEFAULT '0' COMMENT '条件对应的condition表的ID',
  task_id bigint(20) NOT NULL DEFAULT '0' COMMENT '条件对应的子任务的ID',
  plan_id bigint(20) NOT NULL DEFAULT '0' COMMENT '条件对应的子任务的ID',
  created_at int(10) NOT NULL DEFAULT '0',
  updated_at int(10) NOT NULL DEFAULT '0',
  deleted_at int(10) NOT NULL DEFAULT '0',
  customer_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户的customerid信息',
  is_finished tinyint(2) NOT NULL DEFAULT '0' COMMENT '判断任务是否标记完成',
  finished_at int(10) NOT NULL DEFAULT '0' COMMENT '标记完成的时间',
  remark varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '联系说明',
  is_contact tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否联系好友',
  is_call tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否拨打电话',
  is_add tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否点击添加好友',
  is_match tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否匹配成功',
  user_name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户的昵称',
  is_touch tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否触达,0为否，1为是',
  group_chat_id bigint(20) NOT NULL DEFAULT '0' COMMENT '社群id',
  customer_status tinyint(1) NOT NULL DEFAULT 0 COMMENT '顾客沟通类型，0为群发沟通，1为单独沟通',
  forward tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否转发给顾客，0为否，1为是',
  PRIMARY KEY (id) USING BTREE,
  KEY task_allot_log_id_index (task_allot_log_id) USING BTREE,
  KEY task_id_index (task_id) USING BTREE,
  KEY created_at_index (created_at) USING BTREE,
  KEY user_id_index (work_external_user_id,user_id),
  KEY plan_id_index (plan_id) USING BTREE,
  KEY delete_at_index (deleted_at) USING BTREE,
  KEY group_chat_id_index (group_chat_id),
  KEY union_id_index (union_id) USING BTREE,
  KEY index_0 (user_id,task_id,is_finished,created_at),
  KEY customer_id_index (customer_id) USING BTREE,
  KEY index_2 (union_id,deleted_at,customer_id) USING BTREE,
  KEY index_1 (plan_id,deleted_at,union_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC ENCRYPTION='Y' COMMENT='计划任务条件推送数据';`

	res := d.dao.WithContext(c).Exec(createTableSql)

	return res.Error
}

// AlterShardTableColumn 根据主任务创建task_condition_data分表增加字段
func (d *PlanTaskConditionDataDao) AlterShardTableColumn(c ctx.Context, enterpriseId int64, planId int64) error {
	tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName, planId)
	alterTableSql := `ALTER TABLE ` + tableName +
		` ADD COLUMN customer_status tinyint(1) NOT NULL DEFAULT 0 COMMENT '顾客沟通类型，0为群发沟通，1为单独沟通' AFTER group_chat_id,
ADD COLUMN forward tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否转发给顾客，0为否，1为是' AFTER customer_status;`
	res := d.dao.WithContext(c).Exec(alterTableSql)

	return res.Error
}

func (d *PlanTaskConditionDataDao) InsertDataToShard(c ctx.Context, enterpriseId int64, planId int64) error {
	insertTableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName, planId)
	selectTableName := fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName)
	sql := `insert into ` + insertTableName + ` select * from ` + selectTableName + ` where plan_id=?`
	res := d.dao.WithContext(c).Exec(sql, planId)
	return res.Error
}

func (d *PlanTaskConditionDataDao) GetListsByPageNew(c ctx.Context, enterpriseId, planId int64, page int, size int, ids []int64) ([]microScrmModel.TaskConditionDataModel, int64, error) {
	var result []microScrmModel.TaskConditionDataModel
	var total int64
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			db.Where("task_allot_log_id in ?", ids)
		}
		return db
	}).Count(&total).Offset(offset).Limit(size).Find(&result).Error
	return result, total, err
}

func (d *PlanTaskConditionDataDao) GetTaskByDays(c ctx.Context, enterpriseId int64, planId, id int64, size, days, queryType int) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	seconds := int64(days * 86400) // 需要查询的时间范围

	query := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).
		Select("id,union_id,finished_at,created_at")
	if queryType == microScrmModel.TaskReach {
		query.Where("finished_at > ?", time.Now().Unix()-seconds)
	}
	if queryType == microScrmModel.TaskExist {
		query.Where("created_at > ?", time.Now().Unix()-seconds)
	}
	err := query.Where("id > ?", id).
		Order("id asc").
		Limit(size).Find(&result).Error

	return result, err
}

// GetManyByPagination 获取列表分页数据
func (d *PlanTaskConditionDataDao) GetManyByPagination(c app.Ctx, enterpriseId, planId int64, page, size int, query interface{}, args ...interface{}) ([]microScrmModel.TaskConditionDataModel, error) {
	var result []microScrmModel.TaskConditionDataModel
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, enterpriseId, planId)).Where(query, args...).Offset(offset).Limit(size).Find(&result).Error
	return result, err
}

// GetSearchExternalLists  获取任务条件下的所有好友
func (d *PlanTaskConditionDataDao) GetSearchExternalLists(c ctx.Context, request planReq.ExternalSearchListsReq, planId int64) (result []microScrmModel.TaskConditionDataModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, request.EnterpriseId, planId)).Where("task_allot_log_id=?", request.Id).Where("is_match=?", 1).Scopes(func(db *gorm.DB) *gorm.DB {
		if request.FilterType > 0 {
			if request.FilterType == 1 {
				db.Where("is_finished=?", request.FilterType)
			} else {
				db.Where("is_finished!=?", 1)
			}
		}
		if request.Keywords != "" {
			resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, request.EnterpriseId, "store_name")
			if resConfig.Name == "sw" {
				keywords := strings.Replace(request.Keywords, "%", "\\%", -1)
				db.Where("(remark like ? or user_name like ? )", "%"+keywords+"%", "%"+keywords+"%")
			}
		}
		return db
	}).Where("customer_status=?", request.CustomerStatus).Group("union_id,customer_id").Find(&result).Error
	return
}

// GetSearchCustomerLists 获取任务条件下的所有好友
func (d *PlanTaskConditionDataDao) GetSearchCustomerLists(c ctx.Context, request planReq.CustomerSearchListsReq, planId int64, isUnion int, relationUnionIds []string) (result []microScrmModel.TaskConditionDataModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, request.EnterpriseId, planId)).Where("task_allot_log_id=?", request.Id).Where("is_match=?", 1).Scopes(func(db *gorm.DB) *gorm.DB {
		if request.FilterType > 0 {
			if request.FilterType == 1 {
				db.Where("is_finished=?", request.FilterType)
			} else {
				db.Where("is_finished!=?", 1)
			}
		}
		if isUnion > 0 {
			if isUnion == 1 {
				db.Where("union_id in (?)", relationUnionIds)
			} else {
				db.Where("union_id not in (?)", relationUnionIds)
			}
		}
		if request.Keywords != "" {
			keywords := strings.Replace(request.Keywords, "%", "\\%", -1)
			db.Where("(remark like ? or user_name like ? )", "%"+keywords+"%", "%"+keywords+"%")
		}
		return db
	}).Where("customer_status=?", request.CustomerStatus).Group("union_id,customer_id").Find(&result).Error
	return
}

func (d *PlanTaskConditionDataDao) GetSearchExternalListsByUnionIds(c ctx.Context, request planReq.ExternalSearchListsReq, planId int64, unionIds []string) (total int64, result []microScrmModel.TaskConditionDataModel, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, request.EnterpriseId, planId)).Where("task_allot_log_id=?", request.Id).Where("is_match=?", 1).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(unionIds) > 0 {
			db.Where("union_id in ?", unionIds)
		}
		return db
	}).Where("customer_status=?", request.CustomerStatus).Group("union_id,customer_id").Count(&total).Offset(offset).Limit(request.Size).Order("is_finished asc,id asc").Find(&result).Error
	return
}

func (d *PlanTaskConditionDataDao) GetSearchCustomerListsByUnionIds(c ctx.Context, request planReq.CustomerSearchListsReq, planId int64, unionIds []string) (total int64, result []microScrmModel.TaskConditionDataModel, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, request.EnterpriseId, planId)).Where("task_allot_log_id=?", request.Id).Where("is_match=?", 1).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(unionIds) > 0 {
			db.Where("union_id in ?", unionIds)
		}
		return db
	}).Where("customer_status=?", request.CustomerStatus).Group("union_id,customer_id").Count(&total).Offset(offset).Limit(request.Size).Find(&result).Error
	return
}

// CustomerMove 任务列表的顾客数据移动
func (d *PlanTaskConditionDataDao) CustomerMove(c ctx.Context, planId int64, req planReq.CustomerMoveReq) error {
	return d.dao.WithContext(c).Scopes(d.BranchSharedTable(d.tableName, req.EnterpriseId, planId)).Where(map[string]interface{}{
		"task_allot_log_id": req.Id,
		"union_id":          req.UnionIds,
	}).Update("customer_status", req.CustomerStatus).Error
}

func (d *PlanTaskConditionDataDao) BatchUpdateByUnionIds(c context.Context, enterpriseId, planId int64, mappings []microScrmModel.UnionidMapModel) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	// 获取当前批次的最小和最大ID
	minID := mappings[0].ID
	maxID := mappings[len(mappings)-1].ID

	query := `
        UPDATE ` + fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName, planId) + ` AS t1
        INNER JOIN ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, "unionid_map") + ` AS t2 ON t1.work_external_user_id = t2.openid
        SET t1.union_id = t2.new_unionid
        WHERE t2.id >= ? AND t2.id <= ?`

	db := d.dao.WithContext(c).Exec(query, minID, maxID)
	return db.RowsAffected, db.Error
}

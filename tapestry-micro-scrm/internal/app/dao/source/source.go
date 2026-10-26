package source

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/source"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type Dao struct {
	dao       dao.Dao
	tableName string
}

func NewSourceDao() *Dao {
	return &Dao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.SourceModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *Dao) WithTx(tx *gorm.DB) *Dao {
	d.dao.DB = tx
	return d
}

func (d *Dao) CreateData(c app.Ctx, data sourceModels.SourceModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *Dao) SaveData(c app.Ctx, data sourceModels.SourceModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Unscoped().Save(&data).Error
}

func (d *Dao) GetOneUnscoped(c app.Ctx, enterpriseID, sourceID int64) (sourceModels.SourceModel, error) {
	var result sourceModels.SourceModel
	err := d.dao.WithContext(c).Table(d.tableName).Unscoped().
		Where("enterprise_id = ? AND id = ?", enterpriseID, sourceID).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetOne(c app.Ctx, query interface{}, args ...interface{}) (sourceModels.SourceModel, error) {
	var result sourceModels.SourceModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Dao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]sourceModels.SourceModel, error) {
	var result []sourceModels.SourceModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Dao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

func (d *Dao) ForceDelete(c app.Ctx, query interface{}, args ...interface{}) error {
	var model sourceModels.SourceModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Unscoped().Delete(&model).Error
	return err
}

func (d *Dao) GetOneById(c app.Ctx, sourceId, enterpriseId uint64) error {
	var result source.ImageModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"id":            sourceId,
		"enterprise_id": enterpriseId,
	}).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return err
}

func (d *Dao) GetSourceImageDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.ImageCategoryRelationDetail, error) {
	var result []source.ImageCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.id AS source_id,source.group_id,source.merge_id," +
			"si.name AS name, si.url AS url," +
			"sg.group_name," +
			"IFNULL(group_concat(sc.name), ' ') category_names, IFNULL(group_concat(sc.id), ' ') category_ids").
		Joins("LEFT JOIN source_image si ON si.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON si.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"si.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceVideoDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.VideoCategoryRelationDetail, error) {
	var result []source.VideoCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.`id` AS source_id, " +
			"source.`group_id`, sg.`group_name`, " +
			"sv.`name` AS name, sv.url AS url, sv.cover_url," +
			"sv.`status` AS status, sv.file_id, sv.ext, sv.media_id, sv.media_expire_time, sv.media_status, sv.size, " +
			"IFNULL(group_concat(sc.`name`), ' ') category_names, IFNULL(group_concat(sc.`id`), ' ') category_ids").
		Joins("LEFT JOIN source_video sv ON sv.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON sv.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"sv.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceFileDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.FileCategoryRelationDetail, error) {
	var result []source.FileCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.`id` AS source_id, " +
			"source.`group_id`, sg.`group_name`, " +
			"sf.`name` AS name, sf.url AS url, " +
			"sf.`attach_name`, sf.attach_id, sf.file_type, sf.ext, sf.media_id, sf.media_expire_time, sf.media_status, sf.size, " +
			"IFNULL(group_concat(sc.`name`), ' ') category_names, IFNULL(group_concat(sc.`id`), ' ') category_ids").
		Joins("LEFT JOIN source_file sf ON sf.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON sf.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"sf.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceLinkDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.LinkCategoryRelationDetail, error) {
	var result []source.LinkCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.`id` AS source_id, " +
			"source.`group_id`, sg.`group_name`, " +
			"sl.`title`, sl.description, sl.`cover_url`, sl.link_type, sl.link_url, sl.author, sl.content, " +
			"IFNULL(group_concat(sc.`name`), ' ') category_names, IFNULL(group_concat(sc.`id`), ' ') category_ids").
		Joins("LEFT JOIN source_link sl ON sl.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON sl.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"sl.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceMiniProgramDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.MiniProgramCategoryRelationDetail, error) {
	var result []source.MiniProgramCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.`id` AS source_id, " +
			"source.`group_id`, sg.`group_name`, " +
			"smp.`mini_program_id`, smp.title, smp.`cover_url`, smp.path, " +
			"IFNULL(group_concat(sc.`name`), ' ') category_names, IFNULL(group_concat(sc.`id`), ' ') category_ids").
		Joins("LEFT JOIN source_mini_program smp ON smp.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON smp.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"smp.deleted_at":    0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceTextDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.TextCategoryRelationDetail, error) {
	var result []source.TextCategoryRelationDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.`id` AS source_id," +
			"source.`group_id`, sg.`group_name`, " +
			"st.`content`, " +
			"IFNULL(group_concat(sc.`name`), ' ') category_names, IFNULL(group_concat(sc.`id`), ' ') category_ids").
		Joins("LEFT JOIN source_text st ON st.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON st.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"st.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) GetSourceFriendDetail(c app.Ctx, sourceId, enterpriseId uint64) ([]source.FriendDetail, error) {
	var result []source.FriendDetail
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("source.start_time,source.end_time, source.id AS source_id,source.group_id,source.merge_id," +
			"sf.name AS name, sf.original_url,sf.text_content,sf.video_cover_url,sg.group_name," +
			"IFNULL(group_concat(sc.name), ' ') category_names, IFNULL(group_concat(sc.id), ' ') category_ids").
		Joins("LEFT JOIN source_friend sf ON sf.id = source.id").
		Joins("LEFT JOIN source_group sg ON sg.id = source.group_id").
		Joins("LEFT JOIN source_category_relation scr ON sf.id = scr.source_id").
		Joins("LEFT JOIN source_category  sc ON sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source.id":         sourceId,
			"scr.enterprise_id": enterpriseId,
			"source.deleted_at": 0,
			"sf.deleted_at":     0,
			"scr.deleted_at":    0,
			"sc.deleted_at":     0,
		}).
		Find(&result).Error
	return result, err
}

func (d *Dao) UpdateStatus(c app.Ctx, query, column map[string]interface{}, ids []int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).
		Where("id in (?)", ids).
		Updates(column).Error

}

// SelectEndList 查询已过期在启用的素材
func (d *Dao) SelectEndList(c app.Ctx, enterpriseId, nowTime int64) (result []int64, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("id").
		Where("enterprise_id", enterpriseId).
		Where("deleted_at", 0).
		Where("status", 0).
		Where("end_time <> ?", 0).
		Where("end_time <= ?", nowTime).
		Find(&result).Error
	return

}

// GetNewSourceCountForMember 获取起止时间内的素材关联门店下导购查看素材权限
func (d *Dao) GetNewSourceCountForMember(c app.Ctx, enterpriseId, startTime, endTime int64, ignoreGroupIds []int64) (sourceList []source.NewSourceList, err error) {

	// 查询门店上新素材
	query := d.dao.WithContext(c).Table(d.tableName).
		Select(fmt.Sprintf(
			"count(DISTINCT CASE WHEN source.type IN (%d, %d) AND source.merge_id > 0 THEN source.merge_id ELSE source.id END) AS id_count, MAX(source.start_time) AS start_time, sgr.third_id",
			source.TypeImg,
			source.TypeFriend,
		)).
		Joins("LEFT JOIN source_group_relation sgr ON source.group_id = sgr.source_group_id")
	if len(ignoreGroupIds) > 0 {
		query.Where("sgr.source_group_id not in ?", ignoreGroupIds)
	}
	err = query.Where(map[string]interface{}{
		"source.enterprise_id": enterpriseId,
		"source.deleted_at":    0,
		"source.status":        source.StatusEnable,
	}).
		Where("source.start_time >= ?", startTime).
		Where("source.start_time < ?", endTime).
		Where(map[string]interface{}{
			"sgr.enterprise_id": enterpriseId,
			"sgr.deleted_at":    0,
		}).
		Where("sgr.third_id > 0").
		Group("sgr.third_id").
		Find(&sourceList).Error
	return
}

// GetNewSourceCount 查询上新素材数量
func (d *Dao) GetNewSourceCount(c app.Ctx, enterpriseId, startTime, endTime int64, thirdIds []int) (res []source.NewSourceCountNum, err error) {

	// 查询门店上新素材
	err = d.dao.WithContext(c).Table(d.tableName).
		Select(fmt.Sprintf(
			"count(DISTINCT CASE WHEN source.type IN (%d, %d) AND source.merge_id > 0 THEN source.merge_id ELSE source.id END) AS count_num, source.type",
			source.TypeImg,
			source.TypeFriend,
		)).
		Joins("LEFT JOIN source_group_relation sgr ON source.group_id = sgr.source_group_id").
		Where(map[string]interface{}{
			"source.enterprise_id": enterpriseId,
			"source.deleted_at":    0,
			"source.status":        source.StatusEnable,
			"sgr.deleted_at":       0,
		}).
		Where("sgr.third_id in ?", thirdIds).
		Where("source.start_time >= ?", startTime).
		Where("source.start_time < ?", endTime).
		Group("sgr.third_id").
		Group("source.type").
		Find(&res).Error

	return

}

func (d *Dao) Increment(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *Dao) Decrement(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s - ?", column), num)).Error
	return err
}

func (d *Dao) GetSourceIdByCategoryIds(c app.Ctx, enterpriseId int64, sourceType int, categoryIds, sourceIds []int64) ([]int64, error) {
	var res []int64
	err := d.dao.WithContext(c).Table(sourceModels.CategoryRelationModel{}.TableName()).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"category_id":   categoryIds,
		"source_id":     sourceIds,
		"type":          sourceType,
	}).Pluck("source_id", &res).Error
	return res, err

}

// GetTopNum 获取素材的置顶数量
func (d *Dao) GetTopNum(c app.Ctx, enterpriseId int64, sourceType int, sourceGroupIds []int64) (int64, error) {
	var total int64
	target := false
	for _, v := range sourceGroupIds {
		if v == SOURCE_OTHER_GROUP {
			target = true
		}
	}

	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"type":          sourceType,
		"is_top":        sourceModels.TopYes,
	}).Scopes(func(db *gorm.DB) *gorm.DB {
		if target {
			db.Where("group_id=?", SOURCE_OTHER_GROUP)
		} else {
			db.Where("group_id <> ?", SOURCE_OTHER_GROUP)
		}
		return db
	}).Count(&total).Error
	return total, err
}

func (d *Dao) CleanTop(c app.Ctx, sourceIds []int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"id": sourceIds,
	}).Updates(map[string]interface{}{
		"is_top": sourceModels.TopNo,
	}).Error

}

func (d *Dao) GetSourceCollectNum(c app.Ctx, enterpriseId int64, sourceType int32, sourceIds []int64) ([]source.SourceCollectNumInfo, error) {
	var res []source.SourceCollectNumInfo
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"id":            sourceIds,
		"type":          sourceType,
	}).Select("id,merge_id,collect_num").Scan(&res).Error
	return res, err
}

func (d *Dao) GetSourceCollectNumByMergeId(c app.Ctx, enterpriseId int64, sourceType int32, mergeId []int64) ([]source.SourceCollectNumInfo, error) {
	var res []source.SourceCollectNumInfo
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"merge_id":      mergeId,
		"type":          sourceType,
	}).Select("id,merge_id,collect_num").Group("merge_id").Scan(&res).Error
	return res, err
}

func (d *Dao) GetExpireSourceIds(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]int64, error) {
	var res []int64
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"id":            sourceIds,
		"status":        source.StatusEnable,
	}).Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("group_id <> ?", SOURCE_OTHER_GROUP)
		if len(sourceGroupIds) > 0 {
			db.Where("group_id in?", sourceGroupIds)
		}
		return db
	}).Select("id").Scan(&res).Error
	return res, err
}

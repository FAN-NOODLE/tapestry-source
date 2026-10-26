package source

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceReq "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	sourceModels "micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/orm/model/micro_scrm/source"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

const (
	SOURCE_ALL_GROUP   = 0  // 全部分组
	SOURCE_UN_GROUP    = -1 // 未分组
	SOURCE_OTHER_GROUP = -2 // 回收站
)

const linkKOSExistsCondition = `EXISTS (
	SELECT 1 FROM source_link_kos slk
	WHERE slk.enterprise_id = source_link.enterprise_id
	  AND slk.source_id = s.id
	  AND slk.source_link_id = source_link.id
	  AND slk.deleted_at = 0
)`

const linkKOSSyncFlagSelect = linkKOSExistsCondition + ` AS is_kos_sync`

func linkKOSSourceOriginFilter(sourceOrigin int) string {
	switch sourceOrigin {
	case sourceReq.SourceOriginManual:
		return "NOT " + linkKOSExistsCondition
	case sourceReq.SourceOriginKOSSync:
		return linkKOSExistsCondition
	default:
		return ""
	}
}

type SourceGroupDao struct {
	dao.Dao
	tableName string
}

func NewSourceGroupDao() *SourceGroupDao {
	return &SourceGroupDao{
		Dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.SourceGroupModel{}.TableName(),
	}
}

// WithTx 切换数据库连接，用于事务。
func (dao *SourceGroupDao) WithTx(tx *gorm.DB) *SourceGroupDao {
	dao.DB = tx
	return dao
}

func (dao *SourceGroupDao) GetOneByName(c app.Ctx, enterpriseID, parentID int64, groupName string) (sourceModels.SourceGroupModel, error) {
	var result sourceModels.SourceGroupModel
	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("enterprise_id = ? AND parent_id = ? AND group_name = ? AND deleted_at = 0", enterpriseID, parentID, groupName).
		Find(&result).Error
	return result, err
}

func (dao *SourceGroupDao) GetMaxOrders(c app.Ctx, enterpriseID, parentID int64) (int, error) {
	var maxOrders int
	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("enterprise_id = ? AND parent_id = ? AND deleted_at = 0", enterpriseID, parentID).
		Select("COALESCE(MAX(orders), 0)").Scan(&maxOrders).Error
	return maxOrders, err
}

func (dao *SourceGroupDao) CreateData(c app.Ctx, data sourceModels.SourceGroupModel) error {
	return dao.DB.WithContext(c).Table(dao.tableName).Create(&data).Error
}

func (dao *SourceGroupDao) GetChildIDs(c app.Ctx, enterpriseId int64, sourceGroupIds []int64) (error, []source.SourceGroupModel) {
	var result []source.SourceGroupModel
	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("deleted_at = ? ", 0).
		Where("enterprise_id = ? ", enterpriseId).
		Where("parent_id  IN (?)", sourceGroupIds).
		Find(&result).Error

	return err, result
}

func (dao *SourceGroupDao) GetAllList(c app.Ctx, enterpriseId int64) (error, []source.SourceGroupModel) {
	var result []source.SourceGroupModel
	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("deleted_at = ? ", 0).
		Where("enterprise_id = ? ", enterpriseId).
		Order("orders asc").
		Find(&result).Error

	return err, result
}

// queryHandle 处理sql查询
func queryHandle(query *gorm.DB, sourceGroupIds []int64, categoryIds []int64, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64) {
	if len(sourceGroupIds) > 0 {
		var target bool
		for _, v := range sourceGroupIds {
			if v == SOURCE_OTHER_GROUP {
				target = true
			}
		}
		if target == true { // 存在
			query = query.Where(" s.group_id = ?", SOURCE_OTHER_GROUP)
		} else {
			query = query.Where(" s.group_id in (?) AND s.group_id <> ?", sourceGroupIds, SOURCE_OTHER_GROUP)
		}
	} else {
		query = query.Where(" s.group_id <> ?", SOURCE_OTHER_GROUP)
	}
	if len(categoryIds) > 0 {
		query = query.Where(" scr.category_id in (?)", categoryIds).Where("scr.deleted_at=0")
	}

	if status != 2 {
		query = query.Where(" s.status = ?", status)
	}
	if justMe == 1 {
		query = query.Where(" s.admin_id = ?", adminID)
	}
	if isUse == 1 {
		query = query.Where(" s.start_time <= ? AND s.end_time >= ?", time.Now().Unix(), time.Now().Unix())
	}
	if createdStart > 0 && createdEnd > 0 {
		query = query.Where(" s.created_at >= ? AND s.created_at <= ?", createdStart, createdEnd)
	}

}

// filterGoodsRelation 商品关联过滤 - 已关联商品
func (dao *SourceGroupDao) filterGoodsRelationLinked(c app.Ctx, goodsSourceIds []int64) ([]int64, bool) {
	var allGoodsSourceIds []int64
	dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
		Where("deleted_at = 0").
		Pluck("source_id", &allGoodsSourceIds)

	if len(allGoodsSourceIds) == 0 {
		return nil, false
	}

	// 与现有goodsSourceIds取交集
	if len(goodsSourceIds) > 0 {
		goodsSourceIds = helper.ArrayIntersect(goodsSourceIds, allGoodsSourceIds)
	} else {
		goodsSourceIds = allGoodsSourceIds
	}
	goodsSourceIds = helper.ArrayUniq(goodsSourceIds)

	return goodsSourceIds, len(goodsSourceIds) > 0
}

// filterGoodsRelationUnlink 商品关联过滤 - 未关联商品
func (dao *SourceGroupDao) filterGoodsRelationUnlink(c app.Ctx) []int64 {
	var allGoodsSourceIds []int64
	dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
		Where("deleted_at = 0").
		Pluck("source_id", &allGoodsSourceIds)
	return allGoodsSourceIds
}

// applyGoodsUnlinkFilter 应用未关联商品过滤到query
func (dao *SourceGroupDao) applyGoodsUnlinkFilter(query *gorm.DB, c app.Ctx, allGoodsSourceIds []int64, isMergeImage bool) *gorm.DB {
	if len(allGoodsSourceIds) == 0 {
		return query
	}

	if isMergeImage {
		// 图片/朋友圈类型：需要排除merge_id
		var goodsMergeIds []int64
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", allGoodsSourceIds).
			Where("merge_id > 0").
			Scan(&goodsMergeIds)
		query = query.Where(" s.id NOT IN (?) AND (s.merge_id = 0 OR s.merge_id NOT IN (?)) ", allGoodsSourceIds, goodsMergeIds)
	} else {
		// 链接/视频类型：只排除source_id
		query = query.Where(" s.id NOT IN (?) ", allGoodsSourceIds)
	}

	return query
}

func (dao *SourceGroupDao) GetTextConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, page, size int) (error, []source.ResponseConsoleTextList, int64) {

	var (
		total                   int64
		ResponseConsoleTextList []source.ResponseConsoleTextList
		err                     error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.TextModel{}).
		Joins("LEFT JOIN source s ON source_text.id=s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source_text.enterprise_id": enterpriseId,
			"source_text.deleted_at":    0,
			"s.enterprise_id":           enterpriseId,
			"s.deleted_at":              0,
			"s.type":                    sourceModels.TypeText,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		query = query.Where(" source_text.content LIKE ?", "%"+keyWord+"%")
	}
	if order != "" && sort != "" {
		query = query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query = query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	query.Distinct("s.id").Count(&total)

	err = query.
		Select("source_text.id,source_text.enterprise_id,source_text.content," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at," +
			"scr.`type` as source_type," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names,s.is_top").
		Offset((page - 1) * size).Limit(size).Find(&ResponseConsoleTextList).Error

	return err, ResponseConsoleTextList, total
}

func (dao *SourceGroupDao) GetImageConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, productName, skuItem, skuColor string, goodsRelation int, page, size int) (error, []source.ResponseConsoleImageList, int64) {

	var (
		total                    int64
		responseConsoleImageList []source.ResponseConsoleImageList
		err                      error
	)
	var goodsSourceIds []int64
	var results [][]int64

	// 查询函数：如果查询为空，直接返回空结果
	if productName != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("name LIKE ? AND deleted_at = 0", "%"+productName+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleImageList{}, 0
		}
		results = append(results, ids)
	}

	if skuItem != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("sku LIKE ? AND deleted_at = 0", "%"+skuItem+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleImageList{}, 0
		}
		results = append(results, ids)
	}

	if skuColor != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("color LIKE ? AND deleted_at = 0", "%"+skuColor+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleImageList{}, 0
		}
		results = append(results, ids)
	}

	// 计算交集
	if len(results) == 1 {
		goodsSourceIds = results[0]
	} else if len(results) > 1 {
		goodsSourceIds = results[0]
		for i := 1; i < len(results); i++ {
			goodsSourceIds = helper.ArrayIntersect(goodsSourceIds, results[i])
			if len(goodsSourceIds) == 0 {
				break
			}
		}
	}

	goodsSourceIds = helper.ArrayUniq(goodsSourceIds)

	// 如果用户输入了商品搜索条件，但交集为空，直接返回空结果
	if len(results) > 0 && len(goodsSourceIds) == 0 {
		return err, []source.ResponseConsoleImageList{}, 0
	}

	// 商品关联过滤
	if goodsRelation == sourceReq.GoodsRelationLinked { // 已关联商品
		var hasResult bool
		goodsSourceIds, hasResult = dao.filterGoodsRelationLinked(c, goodsSourceIds)
		if !hasResult {
			return err, []source.ResponseConsoleImageList{}, 0
		}
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.ImageModel{}).
		Joins("LEFT JOIN source s ON source_image.id=s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source_image.enterprise_id": enterpriseId,
			"source_image.deleted_at":    0,
			"s.enterprise_id":            enterpriseId,
			"s.deleted_at":               0,
			"s.type":                     sourceModels.TypeImg,
			"s.sort":                     0,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		var sourceImageIds []int64
		// 先查询出source_image.id的集合
		dao.DB.WithContext(c).
			Model(sourceGroup.ImageModel{}).Select("id").
			Where("enterprise_id = ?", enterpriseId).
			Where("name LIKE ?", "%"+keyWord+"%").
			Scan(&sourceImageIds)

		var sourceMergeIds []int64
		// 在查询出merge_id
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", sourceImageIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query.Where("source_image.name LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+keyWord+"%", sourceImageIds, sourceMergeIds)
	}
	if len(goodsSourceIds) > 0 {
		var sourceMergeIds []int64
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", goodsSourceIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)
		query = query.Where(" s.id IN (?) OR s.merge_id IN (?) ", goodsSourceIds, sourceMergeIds)
	}

	// 未关联商品过滤
	if goodsRelation == sourceReq.GoodsRelationUnlink {
		allGoodsSourceIds := dao.filterGoodsRelationUnlink(c)
		query = dao.applyGoodsUnlinkFilter(query, c, allGoodsSourceIds, true)
	}

	if order != "" && sort != "" {
		query.Group("s.id").Order(fmt.Sprintf("s.is_top desc,s.%s %s", order, sort))
	} else {
		query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	dao.QueryClone(c, query).Distinct("s.id").Count(&total)

	err = query.
		Select("source_image.id,source_image.enterprise_id,source_image.name,source_image.url,source_image.size,source_image.ext,source_image.media_id,source_image.media_expire_time,source_image.media_status," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at,s.merge_id,s.sort," +
			"scr.`type` as source_type,s.is_top," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).Limit(size).Find(&responseConsoleImageList).Error

	return err, responseConsoleImageList, total
}

func (dao *SourceGroupDao) GetLinkConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order, sort, productName, skuItem, skuColor string, goodsRelation, sourceOrigin int, page, size int) (error, []source.ResponseConsoleLinkList, int64) {

	var (
		total                   int64
		ResponseConsoleLinkList []source.ResponseConsoleLinkList
		err                     error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.LinkModel{}).
		Joins("LEFT JOIN source s ON source_link.id=s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Joins("LEFT JOIN source_relation_goods srg on s.id = srg.source_id AND srg.deleted_at = 0").
		Where(map[string]interface{}{
			"source_link.enterprise_id": enterpriseId,
			"source_link.deleted_at":    0,
			"s.enterprise_id":           enterpriseId,
			"s.deleted_at":              0,
			"s.type":                    sourceModels.TypeLink,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)
	if filter := linkKOSSourceOriginFilter(sourceOrigin); filter != "" {
		query = query.Where(filter)
	}

	if keyWord != "" {
		query = query.Where(" source_link.title LIKE ? OR  source_link.description LIKE ? ", "%"+keyWord+"%", "%"+keyWord+"%")
	}

	// 商品搜索：直接在JOIN后的表上过滤
	if productName != "" {
		query = query.Where("srg.name LIKE ?", "%"+productName+"%")
	}
	if skuItem != "" {
		query = query.Where("srg.sku LIKE ?", "%"+skuItem+"%")
	}
	if skuColor != "" {
		query = query.Where("srg.color LIKE ?", "%"+skuColor+"%")
	}

	// 商品关联过滤
	if goodsRelation == sourceReq.GoodsRelationLinked { // 已关联商品
		query = query.Where("srg.source_id IS NOT NULL")
	} else if goodsRelation == sourceReq.GoodsRelationUnlink { // 未关联商品
		query = query.Where("srg.source_id IS NULL")
	}

	if order != "" && sort != "" {
		query = query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query = query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	query.Distinct("s.id").Count(&total)

	err = query.
		Select("source_link.id,source_link.enterprise_id,source_link.title,source_link.description,source_link.cover_url,source_link.link_type,source_link.link_url,source_link.author,source_link.content," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at," +
			"scr.`type` as source_type,s.is_top," +
			linkKOSSyncFlagSelect + "," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).Limit(size).Find(&ResponseConsoleLinkList).Error

	return err, ResponseConsoleLinkList, total
}

func (dao *SourceGroupDao) GetMiniProgramConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, page, size int) (error, []source.ResponseConsoleMiniProgramList, int64) {

	var (
		total                          int64
		ResponseConsoleMiniProgramList []source.ResponseConsoleMiniProgramList
		err                            error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.MiniProgramModel{}).
		Joins("LEFT JOIN source s ON source_mini_program.id=s.id").
		Joins("LEFT JOIN `mini_program` mp ON mp.id = source_mini_program.mini_program_id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source_mini_program.enterprise_id": enterpriseId,
			"source_mini_program.deleted_at":    0,
			"s.enterprise_id":                   enterpriseId,
			"s.deleted_at":                      0,
			"s.type":                            sourceModels.TypeMiniProgram,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		query = query.Where(" source_mini_program.title LIKE ? OR mp.app_name LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
	}
	if order != "" && sort != "" {
		query = query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query = query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	query.Distinct("s.id").Count(&total)

	err = query.
		Select("source_mini_program.id,source_mini_program.enterprise_id," +
			"IFNULL(source_mini_program.mini_program_id,\"\") mini_program_id," +
			"IFNULL(source_mini_program.title, \"\") title," +
			"IFNULL(source_mini_program.path, \"\") path," +
			"IFNULL(source_mini_program.cover_url, \"\") cover_url," +
			"IFNULL(mp.origin_id, \"\") origin_id," +
			"IFNULL(mp.app_id, \"\") app_id," +
			"IFNULL(mp.logo, \"\") logo," +
			"IFNULL(mp.app_name, \"\") app_name," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at," +
			"scr.`type` as source_type,s.is_top," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).Limit(size).Find(&ResponseConsoleMiniProgramList).Error

	return err, ResponseConsoleMiniProgramList, total
}

func (dao *SourceGroupDao) GetVideoConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, productName, skuItem, skuColor string, goodsRelation int, page, size int) (error, []source.ResponseConsoleVideoList, int64) {

	var (
		total                    int64
		ResponseConsoleVideoList []source.ResponseConsoleVideoList
		err                      error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.VideoModel{}).
		Joins("LEFT JOIN source s ON source_video.id=s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Joins("LEFT JOIN source_relation_goods srg on s.id = srg.source_id AND srg.deleted_at = 0").
		Where(map[string]interface{}{
			"source_video.enterprise_id": enterpriseId,
			"source_video.deleted_at":    0,
			"s.enterprise_id":            enterpriseId,
			"s.deleted_at":               0,
			"s.type":                     sourceModels.TypeVideo,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		query = query.Where(" source_video.name LIKE ?  ", "%"+keyWord+"%")
	}

	// 商品搜索：直接在JOIN后的表上过滤
	if productName != "" {
		query = query.Where("srg.name LIKE ?", "%"+productName+"%")
	}
	if skuItem != "" {
		query = query.Where("srg.sku LIKE ?", "%"+skuItem+"%")
	}
	if skuColor != "" {
		query = query.Where("srg.color LIKE ?", "%"+skuColor+"%")
	}

	// 商品关联过滤
	if goodsRelation == sourceReq.GoodsRelationLinked { // 已关联商品
		query = query.Where("srg.source_id IS NOT NULL")
	} else if goodsRelation == sourceReq.GoodsRelationUnlink { // 未关联商品
		query = query.Where("srg.source_id IS NULL")
	}

	if order != "" && sort != "" {
		query = query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query = query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	query.Distinct("s.id").Count(&total)

	err = query.
		Select("source_video.id,source_video.enterprise_id,source_video.name,source_video.cover_url,source_video.url,source_video.status,source_video.file_id,source_video.ext,source_video.media_id,source_video.media_expire_time,source_video.media_status,source_video.size," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at," +
			"scr.`type` as source_type,s.is_top," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).Limit(size).Find(&ResponseConsoleVideoList).Error

	return err, ResponseConsoleVideoList, total
}

func (dao *SourceGroupDao) GetFileConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, page, size int) (error, []source.ResponseConsoleFileList, int64) {

	var (
		total                   int64
		ResponseConsoleFileList []source.ResponseConsoleFileList
		err                     error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.FileModel{}).
		Joins("LEFT JOIN source s ON source_file.id=s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source_file.enterprise_id": enterpriseId,
			"source_file.deleted_at":    0,
			"s.enterprise_id":           enterpriseId,
			"s.deleted_at":              0,
			"s.type":                    sourceModels.TypeFile,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		query = query.Where(" source_file.name LIKE ? ", "%"+keyWord+"%")
	}
	if order != "" && sort != "" {
		query = query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query = query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	query.Distinct("s.id").Count(&total)

	err = query.
		Select("source_file.id,source_file.enterprise_id,source_file.name,source_file.url,source_file.attach_name,source_file.attach_id,source_file.file_type,source_file.ext,source_file.media_id,source_file.media_expire_time,source_file.media_status,source_file.size," +
			"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at," +
			"scr.`type` as source_type,s.is_top," +
			"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
			"IFNULL( GROUP_CONCAT( CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).Limit(size).Find(&ResponseConsoleFileList).Error

	return err, ResponseConsoleFileList, total
}

func (dao *SourceGroupDao) GetFriendConsoleList(c app.Ctx, enterpriseId int64, sourceGroupIds []int64, categoryIds []int64, keyWord string, status int, justMe int, adminID int64, isUse int, createdStart int64, createdEnd int64, order string, sort string, productName, skuItem, skuColor string, goodsRelation int, page, size int) (error, []source.ResponseConsoleFriendList, int64) {

	var (
		total               int64
		responseConsoleList []source.ResponseConsoleFriendList
		err                 error
	)
	var goodsSourceIds []int64
	var results [][]int64

	// 查询函数：如果查询为空，直接返回空结果
	if productName != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("name LIKE ? AND deleted_at = 0", "%"+productName+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleFriendList{}, 0
		}
		results = append(results, ids)
	}

	if skuItem != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("sku LIKE ? AND deleted_at = 0", "%"+skuItem+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleFriendList{}, 0
		}
		results = append(results, ids)
	}

	if skuColor != "" {
		var ids []int64
		dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
			Where("color LIKE ? AND deleted_at = 0", "%"+skuColor+"%").
			Pluck("source_id", &ids)
		if len(ids) == 0 {
			return err, []source.ResponseConsoleFriendList{}, 0
		}
		results = append(results, ids)
	}

	// 计算交集
	if len(results) == 1 {
		goodsSourceIds = results[0]
	} else if len(results) > 1 {
		goodsSourceIds = results[0]
		for i := 1; i < len(results); i++ {
			goodsSourceIds = helper.ArrayIntersect(goodsSourceIds, results[i])
			if len(goodsSourceIds) == 0 {
				break
			}
		}
	}

	goodsSourceIds = helper.ArrayUniq(goodsSourceIds)

	// 如果用户输入了商品搜索条件，但交集为空，直接返回空结果
	if len(results) > 0 && len(goodsSourceIds) == 0 {
		return err, []source.ResponseConsoleFriendList{}, 0
	}

	// 商品关联过滤
	if goodsRelation == sourceReq.GoodsRelationLinked { // 已关联商品
		var hasResult bool
		goodsSourceIds, hasResult = dao.filterGoodsRelationLinked(c, goodsSourceIds)
		if !hasResult {
			return err, []source.ResponseConsoleFriendList{}, 0
		}
	}

	query := dao.DB.WithContext(c).Model(sourceModels.FriendModel{}).
		Joins("LEFT JOIN source s ON source_friend.id = s.id").
		Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
		Joins("LEFT JOIN source_category sc on sc.id = scr.category_id").
		Where(map[string]interface{}{
			"source_friend.enterprise_id": enterpriseId,
			"source_friend.deleted_at":    0,
			"s.enterprise_id":             enterpriseId,
			"s.deleted_at":                0,
			"s.type":                      sourceModels.TypeFriend,
			"s.sort":                      0,
		})

	queryHandle(query, sourceGroupIds, categoryIds, status, justMe, adminID, isUse, createdStart, createdEnd)

	if keyWord != "" {
		var sourceImageIds []int64
		// 先查询出source_friend.id的集合
		dao.DB.WithContext(c).
			Model(sourceModels.FriendModel{}).Select("id").
			Where("enterprise_id = ?", enterpriseId).
			Where("name LIKE ?", "%"+keyWord+"%").
			Scan(&sourceImageIds)

		var sourceMergeIds []int64
		// 在查询出merge_id
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", sourceImageIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query.Where("source_friend.name LIKE ? or source_friend.text_content LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+keyWord+"%", "%"+keyWord+"%", sourceImageIds, sourceMergeIds)

	}
	if len(goodsSourceIds) > 0 {
		var sourceMergeIds []int64
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", goodsSourceIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)
		query = query.Where(" s.id IN (?) OR s.merge_id IN (?) ", goodsSourceIds, sourceMergeIds)
	}

	// 未关联商品过滤
	if goodsRelation == sourceReq.GoodsRelationUnlink {
		allGoodsSourceIds := dao.filterGoodsRelationUnlink(c)
		query = dao.applyGoodsUnlinkFilter(query, c, allGoodsSourceIds, true)
	}

	if order != "" && sort != "" {
		query.Group("s.id").Order("s.is_top desc,s." + order + " " + sort)
	} else {
		query.Group("s.id").Order("s.is_top desc,s.created_at DESC")
	}

	dao.QueryClone(c, query).Distinct("s.id").Count(&total)

	err = query.Select("source_friend.id,source_friend.enterprise_id,source_friend.text_content,source_friend.name,source_friend.video_cover_url,source_friend.original_url,source_friend.compression_url,source_friend.ext,source_friend.size,source_friend.type," +
		"s.admin_id,s.group_id,s.type,s.status as source_status,s.start_time,s.end_time,s.created_at,s.merge_id,s.sort," +
		"scr.`type` as source_type,s.is_top," +
		"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.id ELSE NULL END ), \" \" ) AS category_ids," +
		"IFNULL( GROUP_CONCAT( DISTINCT CASE WHEN scr.deleted_at = 0 THEN sc.`name` ELSE NULL END), \"\" ) AS category_names").
		Offset((page - 1) * size).
		Limit(size).
		Find(&responseConsoleList).Error

	return err, responseConsoleList, total
}

func (dao *SourceGroupDao) GetIgnoreIndexGroupIds(c app.Ctx, enterpriseId int64) ([]int64, error) {
	var result []int64
	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("deleted_at = ? ", 0).
		Where("enterprise_id = ? ", enterpriseId).
		Where("ignore_index = ?", source.IgnoreIndex).
		Pluck("id", &result).Error

	return result, err
}

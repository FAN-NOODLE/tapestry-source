package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	sourceModels "micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/orm/model/micro_scrm/source"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type LinkDao struct {
	dao       dao.Dao
	tableName string
}

func NewLinkDao() *LinkDao {
	return &LinkDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.LinkModel{}.TableName(),
	}
}

// WithTx 切换数据库连接，用于事务。
func (d *LinkDao) WithTx(tx *gorm.DB) *LinkDao {
	d.dao.DB = tx
	return d
}

func (d *LinkDao) SaveData(c app.Ctx, data sourceModels.LinkModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Unscoped().Save(&data).Error
}

func (d *LinkDao) GetOneUnscoped(c app.Ctx, enterpriseID, sourceID int64) (sourceModels.LinkModel, error) {
	var result sourceModels.LinkModel
	err := d.dao.WithContext(c).Table(d.tableName).Unscoped().
		Where("enterprise_id = ? AND id = ?", enterpriseID, sourceID).
		Find(&result).Error
	return result, err
}

func (d *LinkDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.LinkModel, error) {
	var result []sourceModels.LinkModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *LinkDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.LinkModel, error) {
	var result []sourceModels.LinkModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *LinkDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

func (d *LinkDao) OpenSearch(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.LinkModel, error) {
	var result []sourceModels.LinkModel
	err := d.dao.WithContext(c).
		Joins("left join source_relation_goods srg on srg.source_id = source_link.id").
		Where("enterprise_id = ?", enterpriseId).
		Where(query, args...).Find(&result).Error
	return result, err
}

func (d *LinkDao) GetListByOpen(c app.Ctx, request openReq.SourceOpenReq) ([]sourceModels.LinkModel, int64, error) {
	var (
		result []sourceModels.LinkModel
		total  int64
	)
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? AND updated_at > ?", request.EnterpriseID, request.UpdatedAt).
		Unscoped().
		Count(&total).Preload("Source").Order("updated_at desc").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&result).Error
	return result, total, err
}

func (dao *GroupRelationDao) GetLinkOpenList(c app.Ctx, request openReq.SearchListRequest, sourceGroupIds []int64) ([]openResp.LinkListResponse, int64, error) {

	var (
		total            int64
		ResponseLinkList []openResp.LinkListResponse
		err              error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.LinkModel{}).
		Joins("LEFT JOIN source s ON source_link.id=s.id").
		Joins("LEFT JOIN source_relation_goods srg on source_link.id = srg.source_id").
		Where(map[string]interface{}{
			"source_link.enterprise_id": request.EnterpriseID,
			"source_link.deleted_at":    0,
			"s.enterprise_id":           request.EnterpriseID,
			"s.deleted_at":              0,
			"srg.deleted_at":            0,
			"s.status":                  sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(request.ProductSku) > 0 {
		query = query.Where(" srg.sku LIKE ?", "%"+request.ProductSku+"%")
	}
	if len(request.ProductColor) > 0 {
		query = query.Where(" srg.color LIKE ? ", "%"+request.ProductColor+"%")
	}
	if len(request.ProductName) > 0 {
		query = query.Where(" srg.name LIKE ? ", "%"+request.ProductName+"%")
	}

	query.Group("source_link.id").Order("source_link.created_at DESC").Count(&total)

	err = query.
		Select("source_link.id,source_link.title,source_link.description,source_link.cover_url,source_link.link_url,s.created_at,s.updated_at").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseLinkList).Error

	return ResponseLinkList, total, err
}

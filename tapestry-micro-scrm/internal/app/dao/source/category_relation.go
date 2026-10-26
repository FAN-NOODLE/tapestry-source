package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type CategoryRelationDao struct {
	dao       dao.Dao
	tableName string
}

func NewCategoryRelationDao() *CategoryRelationDao {
	return &CategoryRelationDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: source.CategoryRelationModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *CategoryRelationDao) WithTx(tx *gorm.DB) *CategoryRelationDao {
	d.dao.DB = tx
	return d
}

func (d *CategoryRelationDao) Create(c app.Ctx, data source.CategoryRelationModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

// BatchCreate 批量插入
func (d *CategoryRelationDao) BatchCreate(c app.Ctx, data []source.CategoryRelationModel) error {
	return d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
}

// Delete 逻辑删除
func (d *CategoryRelationDao) Delete(c app.Ctx, query map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Update("deleted_at", time.Now().Unix()).Error
}

// ForceDelete 强制删除
func (d *CategoryRelationDao) ForceDelete(c app.Ctx, query map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Delete(&source.CategoryRelationModel{}).Error
}

// GetSourceIdsByCategoryId 通过分类id获取关联的资源Id集合
func (d *CategoryRelationDao) GetSourceIdsByCategoryId(c app.Ctx, categoryId int64) ([]int64, error) {
	var sourceIds []int64
	err := d.dao.WithContext(c).Table(d.tableName).Where("category_id", categoryId).Pluck("source_id", &sourceIds).Error
	return sourceIds, err
}

func (d *CategoryRelationDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]source.CategoryRelationModel, error) {
	var result []source.CategoryRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ?", enterpriseId).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *CategoryRelationDao) GetManySourceIds(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var sourceIDs []int64
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and deleted_at=0", enterpriseId).
		Where(query, args...).
		Pluck("source_id", &sourceIDs).Error
	return sourceIDs, err
}

func (d *CategoryRelationDao) GetGetSourceIdsByCategoryIdsBySourceId(c app.Ctx, enterpriseId int64, sourceIds []int64) ([]int64, error) {
	var categoryIds []int64
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"source_id":     sourceIds,
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	}).Pluck("category_id", &categoryIds).Error
	return categoryIds, err
}

func (d *CategoryRelationDao) GetOpenCategoriesBySourceIds(c app.Ctx, enterpriseId int64, sourceIds []int64) (map[int64][]openResp.SourceCategory, error) {
	result := make(map[int64][]openResp.SourceCategory)
	if len(sourceIds) == 0 {
		return result, nil
	}

	var rows []struct {
		SourceID int64  `json:"source_id"`
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
	}
	err := d.dao.WithContext(c).
		Table(d.tableName+" AS scr").
		Select("DISTINCT scr.source_id, sc.id, sc.name").
		Joins("LEFT JOIN source_category AS sc ON sc.id = scr.category_id").
		Where("scr.enterprise_id = ? AND scr.source_id IN ? AND scr.deleted_at = 0 AND sc.enterprise_id = ? AND sc.deleted_at = 0", enterpriseId, sourceIds, enterpriseId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.SourceID] = append(result[row.SourceID], openResp.SourceCategory{
			ID:   row.ID,
			Name: row.Name,
		})
	}
	return result, nil
}

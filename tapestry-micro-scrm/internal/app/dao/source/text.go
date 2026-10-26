package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type TextDao struct {
	dao       dao.Dao
	tableName string
}

func NewTextDao() *TextDao {
	return &TextDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.TextModel{}.TableName(),
	}
}

func (d *TextDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.TextModel, error) {
	var result []sourceModels.TextModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *TextDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.TextModel, error) {
	var result []sourceModels.TextModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *TextDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

// CreateText H5端上传的文本素材插入
func (d *TextDao) CreateText(c app.Ctx, info sourceRequest.SourcesInfo, request sourceRequest.UploadMediaRequest, adminId int64) error {
	err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 插入素材主表，素材主表id和分表id一致
		sourceModelInfo := sourceModels.SourceModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel: models.BaseModel{
					ID: info.SourceId,
				},
				EnterpriseID: int64(int(request.EnterpriseID)),
			},
			GroupID:   request.GroupId,
			Type:      request.SourceType,
			AdminID:   adminId,
			StartTime: int(request.StartTime),
			EndTime:   int(request.EndTime),
		}
		if err := tx.Table(sourceModels.SourceModel{}.TableName()).Create(&sourceModelInfo).Error; err != nil {
			return err
		}
		// 分组标签关联表插入
		if len(request.CategoryIds) > 0 {
			var categoryRelations []sourceModels.CategoryRelationModel
			for _, id := range request.CategoryIds {
				categoryRelation := sourceModels.CategoryRelationModel{
					MicroBaseModel: models.MicroBaseModel{
						BaseModel: models.BaseModel{
							ID: app.UIDContext(c),
						},
						EnterpriseID: int64(int(request.EnterpriseID)),
					},
					CategoryID: id,
					SourceID:   uint64(info.SourceId),
					Type:       request.SourceType,
				}
				categoryRelations = append(categoryRelations, categoryRelation)
			}
			if err := tx.Table(sourceModels.CategoryRelationModel{}.TableName()).Create(&categoryRelations).Error; err != nil {
				return err
			}
		}
		// 插入文本素材表
		sourceText := sourceModels.TextModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel: models.BaseModel{
					ID: info.SourceId,
				},
				EnterpriseID: int64(int(request.EnterpriseID)),
			},
			Content: request.Content,
		}
		if err := tx.Table(sourceModels.TextModel{}.TableName()).Create(&sourceText).Error; err != nil {
			return err
		}
		return nil
	})
	return err

}

func (d *TextDao) GetListByOpen(c app.Ctx, request openReq.SourceOpenReq) ([]sourceModels.TextModel, int64, error) {
	var (
		result []sourceModels.TextModel
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

package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
	helper2 "micro-scrm/internal/pkg/helper"
)

type ImageDao struct {
	dao       dao.Dao
	tableName string
}

func NewImageDao() *ImageDao {
	return &ImageDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.ImageModel{}.TableName(),
	}
}

func (d *ImageDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.ImageModel, error) {
	var result []sourceModels.ImageModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ImageDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.ImageModel, error) {
	var result []sourceModels.ImageModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ImageDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

// CreateImage H5端图片的素材上传
func (d *ImageDao) CreateImage(c app.Ctx, imageInfos []sourceRequest.SourcesMediaInfo, request sourceRequest.UploadMediaRequest, adminId int64) error {
	// 判断是否为组合展示,组合展示需要新增合并id
	var mergeId int64
	if request.ShowType == 2 {
		mergeId = app.UID()
	}

	for _, v := range imageInfos {
		err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
			// 插入素材主表，素材主表id和分表id一致
			info := v
			sourceModelInfo := sourceModels.SourceModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				GroupID:   request.GroupId,
				Type:      request.SourceType,
				StartTime: int(request.StartTime),
				EndTime:   int(request.EndTime),
				MemberID:  request.MemberId,
				IsDeal:    sourceModels.Deal,
				AdminID:   adminId,
				MergeID:   mergeId,
				Sort:      helper2.RuleThree(request.ShowType == 1, 0, info.Sort),
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
			// 插入图片素材表
			sourceImage := sourceModels.ImageModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				Name: info.Name + "." + info.Ext,
				URL:  info.Url,
				Ext:  info.Ext,
				Size: info.Size,
			}
			if err := tx.Table(sourceModels.ImageModel{}.TableName()).Create(&sourceImage).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateCompressImage 要压缩的图片素材先默认插入数据库中 状态为压缩进行中
func (d *ImageDao) CreateCompressImage(c app.Ctx, imageInfos []sourceRequest.SourcesMediaInfo, request sourceRequest.UploadMediaRequest, compressId int64, adminId int64) error {
	// 判断是否为组合展示,组合展示需要新增合并id
	var mergeId int64
	if request.ShowType == 2 {
		mergeId = app.UID()
	}

	for _, v := range imageInfos {
		err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
			// 插入素材主表，素材主表id和分表id一致
			info := v
			sourceModelInfo := sourceModels.SourceModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				GroupID:    request.GroupId,
				Type:       request.SourceType,
				StartTime:  int(request.StartTime),
				EndTime:    int(request.EndTime),
				MemberID:   request.MemberId,
				IsDeal:     sourceModels.DealING,
				CompressID: compressId,
				AdminID:    adminId,
				MergeID:    mergeId,
				Sort:       helper2.RuleThree(request.ShowType == 1, 0, info.Sort),
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
								ID: app.UID(),
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
			// 插入图片素材表
			sourceImage := sourceModels.ImageModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				Name: info.Name + ".jpg",
				URL:  info.Url,
				Ext:  info.Ext,
				Size: info.Size,
			}
			if err := tx.Table(sourceModels.ImageModel{}.TableName()).Create(&sourceImage).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCompress 获取压缩结果
func (d *ImageDao) GetCompress(c ctx.Context, request sourceRequest.CompressReq) (result []sourceRequest.CompressRes, err error) {
	err = d.dao.WithContext(c).Table(sourceModels.SourceModel{}.TableName()).Joins("left join " + d.tableName + " as sm on sm.id=source.id").
		Where(map[string]interface{}{
			"source.member_id":   request.MemberId,
			"source.compress_id": request.CompressId,
		}).Select("sm.url,sm.name,source.is_deal").Scan(&result).Error
	return
}

func (d *ImageDao) GetManyLimit(c app.Ctx, enterpriseId int64, limit int, query interface{}, args ...interface{}) ([]sourceModels.ImageModel, error) {
	var result []sourceModels.ImageModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Order("id asc").Limit(limit).Find(&result).Error
	return result, err
}

func (d *ImageDao) GetListByOpen(c app.Ctx, request openReq.SourceOpenReq) ([]sourceModels.ImageModel, int64, error) {
	var (
		result []sourceModels.ImageModel
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

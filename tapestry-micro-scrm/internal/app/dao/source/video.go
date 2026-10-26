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
	"strconv"
)

type VideoDao struct {
	dao       dao.Dao
	tableName string
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.VideoModel{}.TableName(),
	}
}

func (d *VideoDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.VideoModel, error) {
	var result []sourceModels.VideoModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *VideoDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.VideoModel, error) {
	var result []sourceModels.VideoModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *VideoDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

// CreateVideo H5视频素材上传
func (d *VideoDao) CreateVideo(c app.Ctx, videoInfos []sourceRequest.SourcesMediaInfo, request sourceRequest.UploadMediaRequest, adminId int64) error {
	for _, v := range videoInfos {
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
			// 插入视频素材表
			size := strconv.FormatInt(int64(info.Size), 10)
			sourceImage := sourceModels.VideoModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				Name:     info.Name + "." + info.Ext,
				URL:      info.Url,
				Ext:      info.Ext,
				Size:     size,
				CoverURL: info.CoverUrl,
			}
			if err := tx.Table(sourceModels.VideoModel{}.TableName()).Create(&sourceImage).Error; err != nil {
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

// CreateCompressVideo 要压缩的视频素材先默认插入数据库中 状态为压缩进行中
func (d *VideoDao) CreateCompressVideo(c app.Ctx, videoInfos []sourceRequest.SourcesMediaInfo, request sourceRequest.UploadMediaRequest, compressId int64, adminId int64) error {
	for _, v := range videoInfos {
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
			}

			if err := tx.Table(sourceModels.SourceModel{}.TableName()).Create(&sourceModelInfo).Error; err != nil {
				return err
			}
			// 分组标签关联表插入
			var categoryRelations []sourceModels.CategoryRelationModel
			if len(request.CategoryIds) > 0 {
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
			// 插入视频素材表
			size := strconv.FormatInt(int64(info.Size), 10)
			sourceVideo := sourceModels.VideoModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID: info.SourceId,
					},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				Name:     info.Name + ".mp4",
				URL:      info.Url,
				Ext:      info.Ext,
				Size:     size,
				CoverURL: info.CoverUrl,
			}
			if err := tx.Table(sourceModels.VideoModel{}.TableName()).Create(&sourceVideo).Error; err != nil {
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

// GetCompress 视频的压缩结果
func (d *VideoDao) GetCompress(c ctx.Context, request sourceRequest.CompressReq) (result []sourceRequest.CompressRes, err error) {
	err = d.dao.WithContext(c).Table(sourceModels.SourceModel{}.TableName()).Joins("left join " + d.tableName + " as sv on sv.id=source.id").
		Where(map[string]interface{}{
			"source.member_id":   request.MemberId,
			"source.compress_id": request.CompressId,
		}).Select("sv.cover_url as url,sv.name,source.is_deal").Scan(&result).Error
	return
}

func (d *VideoDao) GetManyLimit(c app.Ctx, enterpriseId int64, limit int, query interface{}, args ...interface{}) ([]sourceModels.VideoModel, error) {
	var result []sourceModels.VideoModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Order("id asc").Limit(limit).Find(&result).Error
	return result, err
}

func (d *VideoDao) GetListByOpen(c app.Ctx, request openReq.SourceOpenReq) ([]sourceModels.VideoModel, int64, error) {
	var (
		result []sourceModels.VideoModel
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

package source

import (
	"strings"

	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
)

type FriendDao struct {
	dao       dao.Dao
	tableName string
}

func NewFriendDao() *FriendDao {
	return &FriendDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.FriendModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *FriendDao) WithTx(tx *gorm.DB) *FriendDao {
	d.dao.DB = tx
	return d
}

func (d *FriendDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.FriendModel, error) {
	var result []sourceModels.FriendModel
	err := d.dao.WithContext(c).Where("enterprise_id = ?", enterpriseId).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *FriendDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.FriendModel, error) {
	var result []sourceModels.FriendModel
	err := d.dao.WithContext(c).Where("enterprise_id = ?", enterpriseId).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *FriendDao) Update(c app.Ctx, query interface{}, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

func (d *FriendDao) ForceDelete(c app.Ctx, query interface{}, args ...interface{}) error {
	var model sourceModels.FriendModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Unscoped().Delete(&model).Error
	return err
}

// CreateFriend 素材上传
func (d *FriendDao) CreateFriend(c app.Ctx, mediaList []sourceRequest.SourcesMediaInfo, request sourceRequest.UploadMediaRequest, compressId, adminId, mergeId int64) error {

	// COS图片压缩参数（与CompressFile保持一致）
	const imageCosCompressParams = "?imageMogr2/thumbnail/1440x/quality/80/format/jpg"

	err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for _, info := range mediaList {
			// 判断是否为图片
			isImage := helper.InArray(strings.ToUpper(info.Ext), []string{"PNG", "JPG", "JPEG", "HEIF"})

			// 图片：is_deal直接标记为已完成（COS参数实时处理，无需异步压缩）
			// 视频：根据IsOrigin判断，非原图需要异步压缩
			isDeal := int32(sourceModels.Deal)
			if !isImage && info.IsOrigin != 1 {
				isDeal = sourceModels.DealING
			}

			// 插入素材主表，素材主表id和分表id一致
			sourceModelInfo := sourceModels.SourceModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel:    models.BaseModel{ID: info.SourceId},
					EnterpriseID: int64(int(request.EnterpriseID)),
				},
				GroupID:    request.GroupId,
				Type:       request.SourceType,
				StartTime:  int(request.StartTime),
				EndTime:    int(request.EndTime),
				MemberID:   request.MemberId,
				IsDeal:     isDeal,
				AdminID:    adminId,
				CompressID: compressId,
				MergeID:    mergeId,
				Sort:       info.Sort,
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
							BaseModel:    models.BaseModel{ID: app.UID()},
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

			// 图片拼COS压缩URL，视频留空（异步压缩后填充）
			compressionURL := ""
			if isImage {
				compressionURL = info.Url + imageCosCompressParams
			}

			// 插入素材表
			sourceFriend := sourceModels.FriendModel{
				MicroBaseModel: models.CreateMicroBase(int64(request.EnterpriseID), info.SourceId),
				TextContent:    request.Content,
				Name:           info.Name,
				VideoCoverURL:  info.CoverUrl,
				OriginalURL:    info.Url,
				CompressionURL: compressionURL,
				Ext:            info.Ext,
				Size:           info.Size,
				Type:           int32(request.FileType),
			}
			if err := tx.Table(sourceModels.FriendModel{}.TableName()).Create(&sourceFriend).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (d *FriendDao) GetListByOpen(c app.Ctx, request openReq.SourceOpenReq) ([]sourceModels.FriendModel, int64, error) {
	var (
		result []sourceModels.FriendModel
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

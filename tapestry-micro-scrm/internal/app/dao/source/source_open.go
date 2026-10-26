package source

import (
	"micro-scrm/internal/app"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	sourceModels "micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/orm/model/micro_scrm/source"
	"time"

	"gorm.io/gorm"
)

// GetTextOpenList open API文本列表（带权限、有效期过滤）
func (dao *GroupRelationDao) GetTextOpenList(c app.Ctx, request openReq.SearchListRequest, sourceGroupIds []int64) ([]openResp.TextListResponse, int64, error) {
	var (
		total    int64
		listResp []openResp.TextListResponse
		err      error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.TextModel{}).
		Joins("LEFT JOIN source s ON source_text.id=s.id").
		Where(map[string]interface{}{
			"source_text.enterprise_id": request.EnterpriseID,
			"source_text.deleted_at":    0,
			"s.enterprise_id":           request.EnterpriseID,
			"s.deleted_at":              0,
			"s.status":                  sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(request.ProductSku) > 0 {
		query = query.Where("source_text.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.sku LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductSku+"%")
	}
	if len(request.ProductColor) > 0 {
		query = query.Where("source_text.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.color LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductColor+"%")
	}
	if len(request.ProductName) > 0 {
		query = query.Where("source_text.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.name LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductName+"%")
	}

	query.Group("source_text.id").Order("source_text.created_at DESC").Count(&total)

	err = query.
		Select("source_text.id,source_text.content," +
			"s.group_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&listResp).Error

	return listResp, total, err
}

// GetImageOpenList open API图片列表（带权限、有效期过滤，只取封面图sort=0）
func (dao *GroupRelationDao) GetImageOpenList(c app.Ctx, request openReq.SearchListRequest, sourceGroupIds []int64) ([]openResp.ImageListResponse, int64, error) {
	var (
		total    int64
		listResp []openResp.ImageListResponse
		err      error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.ImageModel{}).
		Joins("LEFT JOIN source s ON source_image.id=s.id").
		Where(map[string]interface{}{
			"source_image.enterprise_id": request.EnterpriseID,
			"source_image.deleted_at":    0,
			"s.enterprise_id":            request.EnterpriseID,
			"s.deleted_at":               0,
			"s.status":                   sourceModels.StatusEnable,
			"s.sort":                     0,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(request.ProductSku) > 0 {
		query = query.Where("source_image.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.sku LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductSku+"%")
	}
	if len(request.ProductColor) > 0 {
		query = query.Where("source_image.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.color LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductColor+"%")
	}
	if len(request.ProductName) > 0 {
		query = query.Where("source_image.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.name LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductName+"%")
	}

	dao.QueryClone(c, query).Count(&total)

	err = query.
		Select("source_image.id,source_image.name,source_image.url,source_image.ext,source_image.size," +
			"s.group_id,s.merge_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Group("source_image.id").
		Order("source_image.created_at DESC,source_image.id DESC").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&listResp).Error

	return listResp, total, err
}

// GetVideoOpenList open API视频列表（带权限、有效期过滤）
func (dao *GroupRelationDao) GetVideoOpenList(c app.Ctx, request openReq.SearchListRequest, sourceGroupIds []int64) ([]openResp.VideoListResponse, int64, error) {
	var (
		total    int64
		listResp []openResp.VideoListResponse
		err      error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.VideoModel{}).
		Joins("LEFT JOIN source s ON source_video.id=s.id").
		Where(map[string]interface{}{
			"source_video.enterprise_id": request.EnterpriseID,
			"source_video.deleted_at":    0,
			"s.enterprise_id":            request.EnterpriseID,
			"s.deleted_at":               0,
			"s.status":                   sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(request.ProductSku) > 0 {
		query = query.Where("source_video.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.sku LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductSku+"%")
	}
	if len(request.ProductColor) > 0 {
		query = query.Where("source_video.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.color LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductColor+"%")
	}
	if len(request.ProductName) > 0 {
		query = query.Where("source_video.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.name LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductName+"%")
	}

	query.Group("source_video.id").Order("source_video.created_at DESC").Count(&total)

	err = query.
		Select("source_video.id,source_video.name,source_video.cover_url,source_video.url,source_video.ext,source_video.size," +
			"s.group_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&listResp).Error

	return listResp, total, err
}

// GetFriendOpenList open API朋友圈列表（带权限、有效期过滤，只取封面sort=0）
func (dao *GroupRelationDao) GetFriendOpenList(c app.Ctx, request openReq.SearchListRequest, sourceGroupIds []int64) ([]openResp.FriendListResponse, int64, error) {
	var (
		total    int64
		listResp []openResp.FriendListResponse
		err      error
	)

	query := dao.DB.WithContext(c).Model(sourceModels.FriendModel{}).
		Joins("LEFT JOIN source s ON source_friend.id=s.id").
		Where(map[string]interface{}{
			"source_friend.enterprise_id": request.EnterpriseID,
			"source_friend.deleted_at":    0,
			"s.enterprise_id":             request.EnterpriseID,
			"s.deleted_at":                0,
			"s.status":                    sourceModels.StatusEnable,
			"s.type":                      sourceModels.TypeFriend,
			"s.sort":                      0,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(request.ProductSku) > 0 {
		query = query.Where("source_friend.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.sku LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductSku+"%")
	}
	if len(request.ProductColor) > 0 {
		query = query.Where("source_friend.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.color LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductColor+"%")
	}
	if len(request.ProductName) > 0 {
		query = query.Where("source_friend.id IN (SELECT srg.source_id FROM source_relation_goods srg WHERE srg.name LIKE ? AND srg.deleted_at = 0)", "%"+request.ProductName+"%")
	}
	query.Group("source_friend.id").Order("source_friend.created_at DESC").Count(&total)

	err = query.
		Select("source_friend.id,source_friend.text_content,source_friend.name,source_friend.video_cover_url,source_friend.type," +
			"s.group_id,s.merge_id,s.created_at,s.updated_at").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&listResp).Error

	return listResp, total, err
}

// GetSourceIdsByProduct 按商品条件一次性查出所有匹配的 source_id + type
func (dao *GroupRelationDao) GetSourceIdsByProduct(c app.Ctx, enterpriseId int64, productSku, productColor, productName string) ([]sourceModels.IdWithType, error) {
	var result []sourceModels.IdWithType
	query := dao.DB.WithContext(c).Model(sourceModels.RelationGoodsModel{}).
		Select("DISTINCT source_id, type").
		Where("enterprise_id = ? AND deleted_at = 0", enterpriseId)

	if productSku != "" {
		query = query.Where("sku LIKE ?", "%"+productSku+"%")
	}
	if productColor != "" {
		query = query.Where("color LIKE ?", "%"+productColor+"%")
	}
	if productName != "" {
		query = query.Where("name LIKE ?", "%"+productName+"%")
	}

	err := query.Find(&result).Error
	return result, err
}

// GetTextOpenAll 按source_id集合查询文本素材（不分页）
func (dao *GroupRelationDao) GetTextOpenAll(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]openResp.TextListResponse, error) {
	var listResp []openResp.TextListResponse
	if len(sourceIds) == 0 {
		return listResp, nil
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.TextModel{}).
		Joins("LEFT JOIN source s ON source_text.id=s.id").
		Where(map[string]interface{}{
			"source_text.enterprise_id": enterpriseId,
			"source_text.deleted_at":    0,
			"s.enterprise_id":           enterpriseId,
			"s.deleted_at":              0,
			"s.status":                  sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)
	query = query.Where("source_text.id IN (?)", sourceIds)

	err := query.
		Select("source_text.id,source_text.content," +
			"s.group_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Group("source_text.id").
		Order("source_text.created_at DESC").
		Find(&listResp).Error

	return listResp, err
}

// GetImageOpenAll 按source_id集合查询图片素材（不分页，只取封面图sort=0）
func (dao *GroupRelationDao) GetImageOpenAll(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]openResp.ImageListResponse, error) {
	var listResp []openResp.ImageListResponse
	if len(sourceIds) == 0 {
		return listResp, nil
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.ImageModel{}).
		Joins("LEFT JOIN source s ON source_image.id=s.id").
		Where(map[string]interface{}{
			"source_image.enterprise_id": enterpriseId,
			"source_image.deleted_at":    0,
			"s.enterprise_id":            enterpriseId,
			"s.deleted_at":               0,
			"s.status":                   sourceModels.StatusEnable,
			"s.sort":                     0,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)
	query = query.Where("source_image.id IN (?) OR s.merge_id IN (?)", sourceIds, sourceIds)

	err := query.
		Select("source_image.id,source_image.name,source_image.url,source_image.ext,source_image.size," +
			"s.group_id,s.merge_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Group("source_image.id").
		Order("source_image.created_at DESC,source_image.id DESC").
		Find(&listResp).Error

	return listResp, err
}

// GetLinkOpenAll 按source_id集合查询链接素材（不分页）
func (dao *GroupRelationDao) GetLinkOpenAll(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]openResp.LinkListResponse, error) {
	var listResp []openResp.LinkListResponse
	if len(sourceIds) == 0 {
		return listResp, nil
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.LinkModel{}).
		Joins("LEFT JOIN source s ON source_link.id=s.id").
		Where(map[string]interface{}{
			"source_link.enterprise_id": enterpriseId,
			"source_link.deleted_at":    0,
			"s.enterprise_id":           enterpriseId,
			"s.deleted_at":              0,
			"s.status":                  sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)
	query = query.Where("source_link.id IN (?)", sourceIds)

	err := query.
		Select("source_link.id,source_link.title,source_link.description,source_link.cover_url,source_link.link_url,s.group_id,s.created_at,s.updated_at").
		Group("source_link.id").
		Order("source_link.created_at DESC").
		Find(&listResp).Error

	return listResp, err
}

// GetVideoOpenAll 按source_id集合查询视频素材（不分页）
func (dao *GroupRelationDao) GetVideoOpenAll(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]openResp.VideoListResponse, error) {
	var listResp []openResp.VideoListResponse
	if len(sourceIds) == 0 {
		return listResp, nil
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.VideoModel{}).
		Joins("LEFT JOIN source s ON source_video.id=s.id").
		Where(map[string]interface{}{
			"source_video.enterprise_id": enterpriseId,
			"source_video.deleted_at":    0,
			"s.enterprise_id":            enterpriseId,
			"s.deleted_at":               0,
			"s.status":                   sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)
	query = query.Where("source_video.id IN (?)", sourceIds)

	err := query.
		Select("source_video.id,source_video.name,source_video.cover_url,source_video.url,source_video.ext,source_video.size," +
			"s.group_id,s.start_time,s.end_time,s.created_at,s.updated_at").
		Group("source_video.id").
		Order("source_video.created_at DESC").
		Find(&listResp).Error

	return listResp, err
}

// GetFriendOpenAll 按source_id集合查询朋友圈素材（不分页，只取封面sort=0）
func (dao *GroupRelationDao) GetFriendOpenAll(c app.Ctx, enterpriseId int64, sourceIds []int64, sourceGroupIds []int64) ([]openResp.FriendListResponse, error) {
	var listResp []openResp.FriendListResponse
	if len(sourceIds) == 0 {
		return listResp, nil
	}

	query := dao.DB.WithContext(c).Model(sourceModels.FriendModel{}).
		Joins("LEFT JOIN source s ON source_friend.id=s.id").
		Where(map[string]interface{}{
			"source_friend.enterprise_id": enterpriseId,
			"source_friend.deleted_at":    0,
			"s.enterprise_id":             enterpriseId,
			"s.deleted_at":                0,
			"s.status":                    sourceModels.StatusEnable,
			"s.type":                      sourceModels.TypeFriend,
			"s.sort":                      0,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)
	query = query.Where("source_friend.id IN (?) OR s.merge_id IN (?)", sourceIds, sourceIds)

	err := query.
		Select("source_friend.id,source_friend.text_content,source_friend.name,source_friend.video_cover_url,source_friend.original_url,source_friend.type," +
			"s.group_id,s.merge_id,s.created_at,s.updated_at").
		Group("source_friend.id").
		Order("source_friend.created_at DESC").
		Find(&listResp).Error

	return listResp, err
}

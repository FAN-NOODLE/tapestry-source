package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/orm/model/micro_scrm/source"
	"time"
)

type CollectRelationDao struct {
	dao.Repo[source.CollectRelationModel]
	model     *source.CollectRelationModel
	tableName string
}

func NewCollectRelationDao() *CollectRelationDao {
	return &CollectRelationDao{
		Repo: dao.Repo[source.CollectRelationModel]{
			DB: app.DB(),
		},
		model:     new(source.CollectRelationModel),
		tableName: source.CollectRelationModel{}.TableName(),
	}
}

func (d CollectRelationDao) GetCollectSourceIds(c app.Ctx, enterpriseId, memberId int64, sourceIDs []int64, sourceType int) ([]int64, error) {
	var sourceIds []int64
	var err error
	if sourceType < source.TypeFriend {
		err = d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
			"type":          sourceType,
			"is_collect":    source.SourceCollectYes,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
			if len(sourceIDs) > 0 {
				db.Where("source_id in (?)", sourceIDs)
			}
			return db
		}).Order("created_at desc").Pluck("source_id", &sourceIds).Error
	} else {
		err = d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
			"type":          sourceType,
			"is_collect":    source.SourceCollectYes,
		}).Order("created_at desc").Pluck("merge_id", &sourceIds).Error
	}

	return sourceIds, err
}

// GetTextSource 获取文本素材
func (d CollectRelationDao) GetTextSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseTextList, int64, error) {
	var (
		total            int64
		ResponseTextList []source.ResponseTextList
		err              error
	)

	query := d.WithContext(c).Model(source.TextModel{}).
		Joins("LEFT JOIN source s ON source_text.id=s.id").
		Where(map[string]interface{}{
			"source_text.enterprise_id": req.EnterpriseId,
			"source_text.deleted_at":    0,
			"s.enterprise_id":           req.EnterpriseId,
			"s.deleted_at":              0,
			"s.status":                  source.StatusEnable,
		})

	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where(" source_text.id in  ?  ", sourceIds)
	}
	if len(req.Keyword) > 0 {
		query = query.Where(" source_text.content LIKE ?  ", "%"+req.Keyword+"%")
	}
	query.Group("source_text.id").Order("source_text.created_at DESC").Count(&total)

	err = query.
		Select("source_text.id,source_text.enterprise_id,source_text.content," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.copy_times,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&ResponseTextList).Error

	return ResponseTextList, total, err

}

// GetImageSource 获取图片素材
func (d CollectRelationDao) GetImageSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseImgList, int64, error) {
	var (
		total           int64
		responseImgList []source.ResponseImgList
		err             error
	)

	query := d.WithContext(c).Model(source.ImageModel{}).
		Joins("LEFT JOIN source s ON source_image.id=s.id").
		Where(map[string]interface{}{
			"source_image.enterprise_id": req.EnterpriseId,
			"source_image.deleted_at":    0,
			"s.enterprise_id":            req.EnterpriseId,
			"s.deleted_at":               0,
			"s.status":                   source.StatusEnable,
			"s.sort":                     0, // 只查询封面图
		})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where("source_image.id in ?", sourceIds)
	}
	if req.Keyword != "" {
		var sourceImageIds []int64
		// 先查询出source_image.id的集合
		d.WithContext(c).
			Model(source.ImageModel{}).Select("id").
			Where("enterprise_id = ?", req.EnterpriseId).
			Where("name LIKE ?", "%"+req.Keyword+"%").
			Scan(&sourceImageIds)

		var sourceMergeIds []int64
		// 在查询出merge_id
		d.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", sourceImageIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query = query.Where("source_image.name LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+req.Keyword+"%", sourceImageIds, sourceMergeIds)
	}

	d.QueryClone(c, query).Count(&total)

	err = query.
		Select("source_image.id,source_image.enterprise_id,source_image.name,source_image.url,source_image.size,source_image.ext,source_image.media_id,source_image.media_expire_time,source_image.media_status," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.merge_id,s.sort,s.download_times,s.forwarding_times,s.collect_num").
		Group("source_image.id").
		Order("source_image.created_at DESC,source_image.id DESC").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&responseImgList).Error

	return responseImgList, total, err
}

// GetLinkSource 获取链接素材
func (d CollectRelationDao) GetLinkSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseLinkList, int64, error) {
	var (
		total            int64
		ResponseLinkList []source.ResponseLinkList
		err              error
	)

	query := d.WithContext(c).Model(sourceGroup.LinkModel{}).
		Joins("LEFT JOIN source s ON source_link.id=s.id").
		Where(map[string]interface{}{
			"source_link.enterprise_id": req.EnterpriseId,
			"source_link.deleted_at":    0,
			"s.enterprise_id":           req.EnterpriseId,
			"s.deleted_at":              0,
			"s.status":                  source.StatusEnable,
		})

	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where(" source_link.id in  ?  ", sourceIds)
	}
	if len(req.Keyword) > 0 {
		query = query.Where(" source_link.title LIKE ?  OR source_link.description LIKE ? ", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	query.Group("source_link.id").Order("source_link.created_at DESC").Count(&total)

	err = query.
		Select("source_link.id,source_link.enterprise_id,source_link.title,source_link.description,source_link.cover_url,source_link.link_type,source_link.link_url,source_link.author,source_link.content," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&ResponseLinkList).Error

	return ResponseLinkList, total, err
}

// GetMiniProgramSource 获取小程序素材
func (d CollectRelationDao) GetMiniProgramSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseMiniProgramList, int64, error) {

	var (
		total                   int64
		ResponseMiniProgramList []source.ResponseMiniProgramList
		err                     error
	)

	query := d.WithContext(c).Model(sourceGroup.MiniProgramModel{}).
		Joins("LEFT JOIN source s ON source_mini_program.id=s.id").
		Joins("LEFT JOIN `mini_program` mp ON mp.id = source_mini_program.mini_program_id").
		Where(map[string]interface{}{
			"source_mini_program.enterprise_id": req.EnterpriseId,
			"source_mini_program.deleted_at":    0,
			"s.enterprise_id":                   req.EnterpriseId,
			"s.deleted_at":                      0,
			"s.status":                          source.StatusEnable,
		})

	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where(" source_mini_program.id in  ?  ", sourceIds)
	}

	if len(req.Keyword) > 0 {
		query = query.Where(" source_mini_program.title LIKE ? OR mp.app_name LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	query.Group("s.id").Order("s.created_at DESC").Count(&total)

	err = query.
		Select("source_mini_program.id,source_mini_program.enterprise_id,source_mini_program.mini_program_id,source_mini_program.title,source_mini_program.path,source_mini_program.cover_url," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id," +
			"mp.origin_id as origin_id,mp.app_id as app_id,mp.logo as logo,mp.app_name as app_name,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&ResponseMiniProgramList).Error

	return ResponseMiniProgramList, total, err
}

//	GetVideoSource 获取视频素材
func (d CollectRelationDao) GetVideoSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseVideoList, int64, error) {

	var (
		total             int64
		ResponseVideoList []source.ResponseVideoList
		err               error
	)

	query := d.WithContext(c).Model(sourceGroup.VideoModel{}).
		Joins("LEFT JOIN source s ON source_video.id=s.id").
		Where(map[string]interface{}{
			"source_video.enterprise_id": req.EnterpriseId,
			"source_video.deleted_at":    0,
			"s.enterprise_id":            req.EnterpriseId,
			"s.deleted_at":               0,
			"s.status":                   source.StatusEnable,
		})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where(" source_video.id in  ?  ", sourceIds)
	}
	if len(req.Keyword) > 0 {
		query = query.Where(" source_video.name LIKE ?  ", "%"+req.Keyword+"%")
	}

	query.Group("source_video.id").Order("source_video.created_at DESC").Count(&total)

	err = query.
		Select("source_video.id,source_video.enterprise_id,source_video.name,source_video.cover_url,source_video.url,source_video.status,source_video.file_id,source_video.ext,source_video.media_id,source_video.media_expire_time,source_video.media_status,source_video.size," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.download_times,s.forwarding_times,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&ResponseVideoList).Error

	return ResponseVideoList, total, err
}

//	GetFileSource 获取文件素材
func (d CollectRelationDao) GetFileSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseFileList, int64, error) {

	var (
		total            int64
		ResponseFileList []source.ResponseFileList
		err              error
	)
	query := d.WithContext(c).Model(sourceGroup.FileModel{}).
		Joins("LEFT JOIN source s ON source_file.id=s.id").
		Where(map[string]interface{}{
			"source_file.enterprise_id": req.EnterpriseId,
			"source_file.deleted_at":    0,
			"s.enterprise_id":           req.EnterpriseId,
			"s.deleted_at":              0,
			"s.status":                  source.StatusEnable,
		})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where(" source_file.id in  ?  ", sourceIds)
	}
	if len(req.Keyword) > 0 {
		query = query.Where(" source_file.name LIKE ?  ", "%"+req.Keyword+"%")
	}

	query.Group("source_file.id").Order("source_file.created_at DESC").Count(&total)

	err = query.
		Select("source_file.id,source_file.enterprise_id,source_file.name,source_file.url,source_file.attach_name,source_file.attach_id,source_file.file_type,source_file.ext,source_file.media_id,source_file.media_expire_time,source_file.media_status,source_file.size," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&ResponseFileList).Error

	return ResponseFileList, total, err
}

// 	GetFriendSource 获取朋友圈素材
func (d CollectRelationDao) GetFriendSource(c ctx.Context, req sourceRequest.CollectListReq, sourceIds []int64) ([]source.ResponseFriendList, int64, error) {

	var (
		total              int64
		responseFriendList []source.ResponseFriendList
		err                error
	)

	query := d.WithContext(c).Model(source.FriendModel{}).
		Joins("LEFT JOIN source s ON source_friend.id=s.id").
		Where(map[string]interface{}{
			"source_friend.enterprise_id": req.EnterpriseId,
			"source_friend.deleted_at":    0,
			"s.enterprise_id":             req.EnterpriseId,
			"s.deleted_at":                0,
			"s.status":                    source.StatusEnable,
			"s.type":                      source.TypeFriend,
			"s.sort":                      0,
		})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIds) > 0 {
		query = query.Where("source_friend.id in  ?  ", sourceIds)
	}
	if req.Keyword != "" {
		var sourceImageIds []int64
		// 先查询出source_friend.id的集合
		d.WithContext(c).
			Model(source.FriendModel{}).Select("id").
			Where("enterprise_id = ?", req.EnterpriseId).
			Where("name LIKE ?", "%"+req.Keyword+"%").
			Scan(&sourceImageIds)

		var sourceMergeIds []int64
		// 在查询出merge_id
		d.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", sourceImageIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query.Where("source_friend.name LIKE ? or source_friend.text_content LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+req.Keyword+"%", "%"+req.Keyword+"%", sourceImageIds, sourceMergeIds)
	}

	query.Group("source_friend.id").Order("source_friend.created_at DESC").Count(&total)

	err = query.
		Select("source_friend.id,source_friend.enterprise_id,source_friend.text_content,source_friend.name,source_friend.video_cover_url,source_friend.original_url,source_friend.compression_url,source_friend.ext,source_friend.size,source_friend.type," +
			"s.group_id,s.type,s.status as source_status,s.admin_id,s.member_id,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.merge_id,s.sort,s.download_times,s.forwarding_times,s.copy_times,s.collect_num").
		Offset((req.Page - 1) * req.Size).Limit(req.Size).
		Find(&responseFriendList).Error

	return responseFriendList, total, err
}

func (d CollectRelationDao) GetImageMergeCollectNum(c ctx.Context, enterpriseId int64) ([]source.ImgMergeCollect, error) {
	var res []source.ImgMergeCollect
	err := d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"type":          source.TypeImg,
		"is_collect":    source.SourceCollectYes,
	}).Select("count(DISTINCT member_id) as num,merge_id").Group("merge_id").Scan(&res).Error
	return res, err
}

func (d CollectRelationDao) GetImageMergeCollectNumByMergeId(c ctx.Context, enterpriseId int64, mergeId int64) (int64, error) {
	var total int64
	err := d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"type":          source.TypeImg,
		"merge_id":      mergeId,
		"is_collect":    source.SourceCollectYes,
	}).Select("count(DISTINCT member_id) as num").Count(&total).Error
	return total, err
}

func (d CollectRelationDao) GetCollectAllSourceIds(c app.Ctx, enterpriseId, memberId int64) ([]int64, error) {
	var sourceIds []int64
	err := d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"is_collect":    source.SourceCollectYes,
		"member_id":     memberId,
	}).Scopes(func(scope *gorm.DB) *gorm.DB {
		return scope.Where("source_id >0")
	}).Pluck("source_id", &sourceIds).Error
	return sourceIds, err
}

func (d CollectRelationDao) GetFriendSourceId(c app.Ctx, enterpriseId, memberId int64) ([]int64, error) {
	var mergeIds []int64
	var sourceIds []int64
	err := d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"is_collect":    source.SourceCollectYes,
		"member_id":     memberId,
		"type":          source.TypeFriend,
	}).Pluck("merge_id", &mergeIds).Error
	if err != nil || len(mergeIds) == 0 {
		return sourceIds, err
	}
	err = d.WithContext(c).Model(&source.SourceModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"type":          source.TypeFriend,
		"merge_id":      mergeIds,
		"sort":          0,
		"status":        source.StatusEnable,
	}).Pluck("id", &sourceIds).Error
	return sourceIds, err
}

func (d CollectRelationDao) GetCollectProduct(c app.Ctx, enterpriseId, memberId int64, keywords string) ([]source.CollectRelationModel, error) {
	var res []source.CollectRelationModel
	err := d.WithContext(c).Model(&d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
		"type":          source.TypeProduct,
		"is_collect":    source.SourceCollectYes,
	}).Scopes(func(db *gorm.DB) *gorm.DB {
		if keywords != "" {
			db.Where("product_code like ?", "%"+keywords+"%")
		}
		return db
	}).Find(&res).Error
	return res, err
}

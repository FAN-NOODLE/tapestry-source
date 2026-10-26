package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	sourceModels "micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/orm/model/micro_scrm/source"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type GroupRelationDao struct {
	dao.Dao
	tableName string
}

func NewGroupRelationDao() *GroupRelationDao {
	return &GroupRelationDao{
		Dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.GroupRelationModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (dao *GroupRelationDao) WithTx(tx *gorm.DB) *GroupRelationDao {
	dao.DB = tx
	return dao
}
func (dao *GroupRelationDao) CreateData(c app.Ctx, data source.GroupRelationList) error {
	err := dao.WithContext(c).Table(dao.tableName).Create(&data).Error
	return err
}

func (dao *GroupRelationDao) GetList(c app.Ctx, enterpriseId int64, sourceGroupId int64, page, size int) (error, []source.GroupRelationList, int64) {
	query := dao.DB.WithContext(c).Table(dao.tableName).
		Select("source_group_relation.id,source_group_relation.enterprise_id,source_group_relation.third_id,source_group_relation.source_group_id,d.`name` AS department_name").
		Joins("LEFT JOIN departments d ON source_group_relation.third_id=d.third_id").
		Where("source_group_relation.enterprise_id=? AND source_group_relation.deleted_at=? AND source_group_relation.source_group_id=?", enterpriseId, 0, sourceGroupId).
		Where("d.enterprise_id=? AND d.deleted_at=? ", enterpriseId, 0)

	var (
		total             int64
		groupRelationList []source.GroupRelationList
		err               error
	)

	err = query.Count(&total).Order("source_group_relation.department_id ASC").
		// Offset((page - 1) * size).Limit(size).
		Find(&groupRelationList).Error

	return err, groupRelationList, total
}

func (dao *GroupRelationDao) GetListForId(c app.Ctx, enterpriseId int64, sourceGroupId int64, thirdIds []int) (error, []source.GroupRelationListAll) {

	query := dao.DB.WithContext(c).Model(sourceGroup.GroupModel{}).
		Select("source_group.id,source_group.group_name,source_group.enterprise_id,source_group.parent_id,s.third_id,s.source_group_id,d.`name` AS department_name").
		Joins("LEFT JOIN source_group_relation s ON source_group.id=s.source_group_id").
		Joins("LEFT JOIN departments d ON s.third_id=d.third_id").
		Where("source_group.enterprise_id=? AND source_group.deleted_at=?", enterpriseId, 0).
		Where("s.enterprise_id=? AND s.deleted_at=? AND s.third_id in ?", enterpriseId, 0, thirdIds).
		Where("d.enterprise_id=? AND d.deleted_at=?", enterpriseId, 0)

	// if sourceGroupId == 0 {
	//	query = query.Where(" source_group.parent_id=0")
	// } else {
	//	query = query.Where(" source_group.parent_id=?", sourceGroupId)
	// }

	var (
		GroupRelationList []source.GroupRelationListAll
		err               error
	)

	err = query.Order("source_group.id ASC").Find(&GroupRelationList).Error

	return err, GroupRelationList
}

func (dao *GroupRelationDao) Save(c app.Ctx, groupRelation []*source.GroupRelationModel) error {
	if err := dao.DB.WithContext(c).Table(dao.tableName).Create(groupRelation).Error; err != nil {
		return err
	}
	return nil
}

func (dao *GroupRelationDao) DeleteForDepartmentIds(c app.Ctx, thirdIds []int) error {
	return dao.DB.WithContext(c).Table(dao.tableName).
		Where("deleted_at = 0").
		Where("third_id in (?)", thirdIds).
		Update("deleted_at", time.Now().Unix()).Error
}

func (dao *GroupRelationDao) DeleteForSourceGroupIds(c app.Ctx, sourceGroupIds []int64) error {
	return dao.DB.WithContext(c).Table(dao.tableName).
		Where("deleted_at = 0").
		Where("source_group_id in (?)", sourceGroupIds).
		Update("deleted_at", time.Now().Unix()).Error
}

// querySelectHandle 处理小程序个素材查询sql
func querySelectHandle(query *gorm.DB, sourceGroupIds []int64, categoryIds []int64) {

	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", time.Now().Unix(), time.Now().Unix())

	if len(sourceGroupIds) > 0 {
		query = query.Where(" sgr.source_group_id in (?)", sourceGroupIds)
	}
	if len(categoryIds) > 0 {
		query = query.Where(" scr.category_id in (?)  ", categoryIds)
	}
}

func (dao *GroupRelationDao) GetTextSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64) ([]source.ResponseTextList, int64, error) {

	var (
		total            int64
		ResponseTextList []source.ResponseTextList
		err              error
	)

	if len(request.CategoryIds) > 0 && len(sourceIDs) == 0 {
		return ResponseTextList, total, nil
	}

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
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where(" source_text.id in  ?  ", sourceIDs)
	}
	if len(request.Keyword) > 0 {
		query = query.Where(" source_text.content LIKE ?  ", "%"+request.Keyword+"%")
	}

	query.Group("source_text.id").Order("s.is_top desc,s.start_time DESC,source_text.created_at DESC").Count(&total)

	err = query.
		Select("source_text.id,source_text.enterprise_id,source_text.content," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.copy_times,s.collect_num,s.is_top").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseTextList).Error

	return ResponseTextList, total, err
}

func (dao *GroupRelationDao) GetImageSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64, collectSearch bool) ([]source.ResponseImgList, int64, error) {

	var (
		total           int64
		responseImgList []source.ResponseImgList
		err             error
	)

	if len(request.CategoryIds) > 0 && len(sourceIDs) == 0 {
		return responseImgList, total, nil
	}

	query := dao.DB.WithContext(c).Model(sourceGroup.ImageModel{}).
		Joins("LEFT JOIN source s ON source_image.id=s.id").
		Where(map[string]interface{}{
			"source_image.enterprise_id": request.EnterpriseID,
			"source_image.deleted_at":    0,
			"s.enterprise_id":            request.EnterpriseID,
			"s.deleted_at":               0,
			"s.status":                   sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			db.Where("s.group_id in (?)", sourceGroupIds)
		}

		if !collectSearch {
			db.Where("s.sort=0") // 只查询封面图
		}
		return db
	})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where("s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where("source_image.id in ?", sourceIDs)
	}
	if request.Keyword != "" {
		var sourceImageIds []int64
		// 先查询出source_image.id的集合
		dao.DB.WithContext(c).
			Model(sourceGroup.ImageModel{}).Select("id").
			Where("enterprise_id = ?", request.EnterpriseID).
			Where("name LIKE ?", "%"+request.Keyword+"%").
			Scan(&sourceImageIds)

		// 匹配商品名称和SKU
		var keywordGoodsIds []int64
		dao.DB.WithContext(c).
			Model(sourceModels.RelationGoodsModel{}).Select("source_id").
			Where("deleted_at = 0").
			Where("name LIKE ? OR sku_code LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%").
			Scan(&keywordGoodsIds)

		allSourceIds := append(sourceImageIds, keywordGoodsIds...)

		var sourceMergeIds []int64
		// 在查询出merge_id
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", allSourceIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query = query.Where("source_image.name LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+request.Keyword+"%", allSourceIds, sourceMergeIds)
	}

	dao.QueryClone(c, query).Count(&total)

	err = query.
		Select("source_image.id,source_image.enterprise_id,source_image.name,source_image.url,source_image.size,source_image.ext,source_image.media_id,source_image.media_expire_time,source_image.media_status," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.merge_id,s.sort,s.download_times,s.forwarding_times,s.collect_num,s.is_top").
		Group("source_image.id").
		Order("s.is_top desc,s.start_time DESC,source_image.created_at DESC,source_image.id DESC").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&responseImgList).Error

	return responseImgList, total, err
}

func (dao *GroupRelationDao) GetLinkSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64) ([]source.ResponseLinkList, int64, error) {

	var (
		total            int64
		ResponseLinkList []source.ResponseLinkList
		err              error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.LinkModel{}).
		Joins("LEFT JOIN source s ON source_link.id=s.id").
		Where(map[string]interface{}{
			"source_link.enterprise_id": request.EnterpriseID,
			"source_link.deleted_at":    0,
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
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where(" source_link.id in  ?  ", sourceIDs)
	}
	if len(request.Keyword) > 0 {
		// 匹配商品名称和SKU
		var keywordGoodsIds []int64
		dao.DB.WithContext(c).
			Model(sourceModels.RelationGoodsModel{}).Select("source_id").
			Where(map[string]interface{}{
				"enterprise_id": request.EnterpriseID,
				"type":          sourceModels.TypeLink,
				"deleted_at":    0,
			}).
			Where("name LIKE ? OR sku_code LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%").
			Scan(&keywordGoodsIds)
		query = query.Where(" source_link.title LIKE ? OR  source_link.description LIKE ? OR source_link.id in ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%", keywordGoodsIds)

	}

	dao.QueryClone(c, query).Count(&total)

	err = query.
		Select("source_link.id,source_link.enterprise_id,source_link.title,source_link.description,source_link.cover_url,source_link.link_type,source_link.link_url,source_link.author,source_link.content," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.collect_num,s.is_top").
		Order("s.is_top desc,s.start_time DESC,source_link.created_at DESC").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseLinkList).Error

	return ResponseLinkList, total, err
}

func (dao *GroupRelationDao) GetMiniProgramSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64) ([]source.ResponseMiniProgramList, int64, error) {

	var (
		total                   int64
		ResponseMiniProgramList []source.ResponseMiniProgramList
		err                     error
	)

	query := dao.DB.WithContext(c).Model(sourceGroup.MiniProgramModel{}).
		Joins("LEFT JOIN source s ON source_mini_program.id=s.id").
		Joins("LEFT JOIN `mini_program` mp ON mp.id = source_mini_program.mini_program_id").
		Where(map[string]interface{}{
			"source_mini_program.enterprise_id": request.EnterpriseID,
			"source_mini_program.deleted_at":    0,
			"s.enterprise_id":                   request.EnterpriseID,
			"s.deleted_at":                      0,
			"s.status":                          sourceModels.StatusEnable,
		}).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(sourceGroupIds) > 0 {
			return db.Where("s.group_id in (?)", sourceGroupIds)
		}
		return db
	})

	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where(" source_mini_program.id in  ?  ", sourceIDs)
	}

	if len(request.Keyword) > 0 {
		query = query.Where(" source_mini_program.title LIKE ? OR mp.app_name LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%")
	}

	query.Group("s.id").Order("s.is_top desc,s.start_time DESC,s.created_at DESC").Count(&total)

	err = query.
		Select("source_mini_program.id,source_mini_program.enterprise_id,source_mini_program.mini_program_id,source_mini_program.title,source_mini_program.path,source_mini_program.cover_url," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id," +
			"mp.origin_id as origin_id,mp.app_id as app_id,mp.logo as logo,mp.app_name as app_name,s.collect_num,s.is_top").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseMiniProgramList).Error

	return ResponseMiniProgramList, total, err

	// query := dao.DB.WithContext(c).Model(sourceGroup.MiniProgramModel{}).
	//	Joins("LEFT JOIN source s ON source_mini_program.id=s.id").
	//	Joins("LEFT JOIN `mini_program` mp ON mp.id = source_mini_program.mini_program_id").
	//	Joins("LEFT JOIN source_group_relation sgr ON sgr.source_group_id=s.group_id").
	//	Joins("LEFT JOIN source_category_relation scr ON scr.source_id=s.id").
	//	Where(map[string]interface{}{
	//		"source_mini_program.enterprise_id": request.EnterpriseID,
	//		"source_mini_program.deleted_at":    0,
	//		"s.enterprise_id":                   request.EnterpriseID,
	//		"s.deleted_at":                      0,
	//		"s.status":                          sourceModels.StatusEnable,
	//		"sgr.enterprise_id":                 request.EnterpriseID,
	//		"sgr.deleted_at":                    0,
	//		"sgr.third_id":                      request.ThirdIds,
	//	})
	//
	// querySelectHandle(query, request.GroupIDList, request.CategoryIds)
	//
	// if len(request.Keyword) > 0 {
	//	query = query.Where(" source_mini_program.title LIKE ? OR mp.app_name LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%")
	// }
	// var (
	//	total                   int64
	//	ResponseMiniProgramList []source.ResponseMiniProgramList
	//	err                     error
	// )
	// query.Group("s.id").Order("s.created_at DESC").Distinct("s.id").Count(&total)
	//
	// err = query.
	//	Select("source_mini_program.id,source_mini_program.enterprise_id,source_mini_program.mini_program_id,source_mini_program.title,source_mini_program.path,source_mini_program.cover_url," +
	//		"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,sgr.third_id,s.is_deal,sgr.source_group_id,IFNULL(group_concat(scr.category_id), \" \") AS category_ids," +
	//		"mp.origin_id as origin_id,mp.app_id as app_id,mp.logo as logo,mp.app_name as app_name").
	//	Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&ResponseMiniProgramList).Error
	//
	// return err, ResponseMiniProgramList, total
}

func (dao *GroupRelationDao) GetVideoSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64) ([]source.ResponseVideoList, int64, error) {

	var (
		total             int64
		ResponseVideoList []source.ResponseVideoList
		err               error
	)

	if len(request.CategoryIds) > 0 && len(sourceIDs) == 0 {
		return ResponseVideoList, total, nil
	}

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
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where(" source_video.id in  ?  ", sourceIDs)
	}
	if request.Keyword != "" {
		// 匹配商品名称和SKU
		var keywordGoodsIds []int64
		dao.DB.WithContext(c).
			Model(sourceModels.RelationGoodsModel{}).Select("source_id").
			Where("deleted_at = 0").
			Where("name LIKE ? OR sku_code LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%").
			Scan(&keywordGoodsIds)

		query = query.Where(" source_video.name LIKE ? OR source_video.id IN (?)", "%"+request.Keyword+"%", keywordGoodsIds)
	}

	query.Group("source_video.id").Order("s.is_top desc,s.start_time DESC,source_video.created_at DESC").Count(&total)

	err = query.
		Select("source_video.id,source_video.enterprise_id,source_video.name,source_video.cover_url,source_video.url,source_video.status,source_video.file_id,source_video.ext,source_video.media_id,source_video.media_expire_time,source_video.media_status,source_video.size," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.download_times,s.forwarding_times,s.collect_num,s.is_top").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseVideoList).Error

	return ResponseVideoList, total, err
}

func (dao *GroupRelationDao) GetFileSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds []int64) ([]source.ResponseFileList, int64, error) {

	var (
		total            int64
		ResponseFileList []source.ResponseFileList
		err              error
	)
	query := dao.DB.WithContext(c).Model(sourceGroup.FileModel{}).
		Joins("LEFT JOIN source s ON source_file.id=s.id").
		Where(map[string]interface{}{
			"source_file.enterprise_id": request.EnterpriseID,
			"source_file.deleted_at":    0,
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
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where(" source_file.id in  ?  ", sourceIDs)
	}
	if len(request.Keyword) > 0 {
		query = query.Where(" source_file.name LIKE ?  ", "%"+request.Keyword+"%")
	}

	query.Group("source_file.id").Order("s.is_top desc,s.start_time DESC,source_file.created_at DESC").Count(&total)

	err = query.
		Select("source_file.id,source_file.enterprise_id,source_file.name,source_file.url,source_file.attach_name,source_file.attach_id,source_file.file_type,source_file.ext,source_file.media_id,source_file.media_expire_time,source_file.media_status,source_file.size," +
			"s.group_id,s.type,s.status as source_status,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.collect_num,s.is_top").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&ResponseFileList).Error

	return ResponseFileList, total, err
}

func (dao *GroupRelationDao) GetFriendSelectList(c app.Ctx, request sourceRequest.RequestSelectList, sourceIDs, sourceGroupIds, collectSourceIDs []int64) ([]source.ResponseFriendList, int64, error) {

	var (
		total              int64
		responseFriendList []source.ResponseFriendList
		err                error
	)

	if len(request.CategoryIds) > 0 && len(sourceIDs) == 0 {
		return responseFriendList, total, nil
	}

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
			db.Where("s.group_id in (?)", sourceGroupIds)
		}
		if len(collectSourceIDs) > 0 {
			db.Where("s.merge_id in (?)", collectSourceIDs)
		}
		return db
	})
	var nowTime = time.Now().Unix()
	// 有效期内
	query = query.Where(" s.start_time <= ? AND s.end_time >= ?", nowTime, nowTime)

	if len(sourceIDs) > 0 {
		query = query.Where("source_friend.id in  ?  ", sourceIDs)
	}
	if request.Keyword != "" {
		var sourceFriendIds []int64
		// 先查询出source_friend.id的集合
		dao.DB.WithContext(c).
			Model(sourceModels.FriendModel{}).Select("id").
			Where("enterprise_id = ?", request.EnterpriseID).
			Where("name LIKE ?", "%"+request.Keyword+"%").
			Scan(&sourceFriendIds)

		// 匹配商品名称和SKU
		var keywordGoodsIds []int64
		dao.DB.WithContext(c).
			Model(sourceModels.RelationGoodsModel{}).Select("source_id").
			Where("deleted_at = 0").
			Where("name LIKE ? OR sku_code LIKE ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%").
			Scan(&keywordGoodsIds)

		allSourceIds := append(sourceFriendIds, keywordGoodsIds...)

		var sourceMergeIds []int64
		// 在查询出merge_id
		dao.DB.WithContext(c).
			Model(sourceGroup.SourceModel{}).Select("merge_id").
			Where("id in (?)", allSourceIds).
			Where("merge_id > 0").
			Scan(&sourceMergeIds)

		query = query.Where("source_friend.name LIKE ? or source_friend.text_content LIKE ? or s.id IN (?) or s.merge_id IN (?)", "%"+request.Keyword+"%", "%"+request.Keyword+"%", allSourceIds, sourceMergeIds)
	}

	query.Group("source_friend.id").Order("s.is_top desc,s.start_time DESC,source_friend.created_at DESC").Count(&total)

	err = query.
		Select("source_friend.id,source_friend.enterprise_id,source_friend.text_content,source_friend.name,source_friend.video_cover_url,source_friend.original_url,source_friend.compression_url,source_friend.ext,source_friend.size,source_friend.type," +
			"s.group_id,s.type,s.status as source_status,s.admin_id,s.member_id,s.created_at,s.start_time,s.end_time,s.is_deal,s.group_id as source_group_id,s.merge_id,s.sort,s.download_times,s.forwarding_times,s.copy_times,s.collect_num,s.is_top").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&responseFriendList).Error

	return responseFriendList, total, err
}

func (dao *GroupRelationDao) GetManySourceGroupIds(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var sourceGroupIds []int64

	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("enterprise_id = ? and deleted_at=0", enterpriseId).
		Where(query, args...).
		Pluck("source_group_id", &sourceGroupIds).Error

	return sourceGroupIds, err
}

func (dao *GroupRelationDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]source.GroupRelationModel, error) {
	var result []source.GroupRelationModel

	err := dao.DB.WithContext(c).Table(dao.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where(query, args...).
		Find(&result).Error

	return result, err
}

package source

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	source2 "micro-scrm/internal/app/dao/source"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/modules/drainage"
	drainageMysql "micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/source"
	"micro-scrm/internal/pkg/assist"
)

// Service  应用服务
type Service struct {
	repo Repository
}

var DrainageService DrainageInterface

type DrainageInterface interface {
	MakeMediaBySourceUrl(ctx ctx.Context, enterpriseID int64, sourceUrl string, sourceType string, sourceName string, domain string) (string, error)
}

// SourceFileList 文件素材列表
func (s Service) SourceFileList(c ctx.Context, params *RequestList) (int64, *[]ResponseFileList, error) {
	total, list, err := s.repo.SourcesFileList(c, params)
	if err != nil {
		return 0, nil, err
	}

	var resultList []ResponseFileList
	resultList = *list

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}
	return total, &resultList, nil
}

// SourceImageAdd 图片素材新增
func (s Service) SourceImageAdd(c ctx.Context, params *SourcesImageRequest) (int64, error) {
	params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseID, params.AdminID)
	id, err := s.repo.AddSourcesImage(c, params)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

// SourceFileAdd 文件素材新增
func (s Service) SourceFileAdd(c ctx.Context, params *SourcesFileRequest) (int64, error) {
	params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseID, params.AdminID)
	_, err := s.repo.AddSourcesFile(c, params)
	if err != nil {
		return 0, nil
	}
	return params.ID, nil
}

// SourceImageList 图片素材列表
func (s Service) SourceImageList(c ctx.Context, params *RequestList) (int64, *[]ResponseImgList, error) {
	total, list, err := s.repo.SourcesImageList(c, params)
	if err != nil {
		return 0, nil, err
	}

	var resultList []ResponseImgList
	resultList = *list

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}

	return total, &resultList, nil
}

// TextSave 文本素材保存
func (s Service) TextSave(c ctx.Context, params *RequestTextSave) (int64, error) {

	if params.ID == 0 {
		params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseId, params.AdminId)
		params.ID = int64(assist.IdGenerate(c))
		params.CreatedAt = assist.Timestamp()
		params.UpdatedAt = 0
	} else {
		params.UpdatedAt = assist.Timestamp()
	}
	params.Type = 1

	_, err := s.repo.TextSave(c, *params)

	if err != nil {
		return 0, err
	}

	return params.ID, nil
}

// LinkSave 链接素材保存
func (s Service) LinkSave(c ctx.Context, params *RequestLinkSave) (int64, error) {

	if params.ID == 0 {
		params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseId, params.AdminId)
		params.ID = int64(assist.IdGenerate(c))
		params.CreatedAt = assist.Timestamp()
		params.UpdatedAt = 0
	} else {
		params.UpdatedAt = assist.Timestamp()
	}
	params.Type = 3

	_, err := s.repo.LinkSave(c, *params)

	if err != nil {
		return 0, err
	}

	// 处理关联商品
	if len(params.RelationGoods) > 0 {
		relationGoodsDao := source2.NewRelationGoodsDao()
		// 批量插入新的关联商品
		var goodsList []sourceModels.RelationGoodsModel
		for _, item := range params.RelationGoods {
			goods := sourceModels.RelationGoodsModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID:        app.UID(),
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
					},
					EnterpriseID: params.EnterpriseId,
				},
				SourceID: params.ID,
				Type:     sourceModels.TypeLink, // 链接素材类型
				UPC:      item.Upc,
				SKU:      item.Sku,
				Name:     item.Name,
				Color:    item.Color,
				SkuCode:  fmt.Sprintf("%s %s", item.Sku, item.Color),
				Url:      item.Url,
				AddType:  item.AddType,
				Score:    item.Score,
			}
			goodsList = append(goodsList, goods)
		}

		if len(goodsList) > 0 {
			err := relationGoodsDao.CreateInBatchesData(c, goodsList)
			if err != nil {
				return 0, err
			}
		}
	}
	return params.ID, nil
}

// TextList 文本素材列表
func (s Service) TextList(c ctx.Context, params *RequestTextList) (int64, *ResponseTextList, error) {
	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-text-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &ResponseTextList{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if err != nil {
		return 0, nil, err
	}
	resList := list.(*ResponseTextList)

	var resultList ResponseTextList
	resultList = *resList

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}

	return total, &resultList, nil
}

// LinkList 链接素材列表
func (s Service) LinkList(c ctx.Context, params *RequestLinkList) (int64, *ResponseLinkList, error) {

	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-link-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &ResponseLinkList{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if err != nil {
		return 0, nil, err
	}
	resList := list.(*ResponseLinkList)

	var resultList ResponseLinkList
	resultList = *resList

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}

	return total, &resultList, nil
}

// SourceTextDetail 文本素材详情
func (s Service) SourceTextDetail(c ctx.Context, params *RequestSourceDetail) (*ResponseTextDetail, error) {
	res, err := s.repo.TextDetail(c, params)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// SourceLinkDetail 链接素材详情
func (s Service) SourceLinkDetail(c ctx.Context, params *RequestSourceDetail) (*ResponseLinkDetail, error) {
	res, err := s.repo.LinkDetail(c, params)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// SourceMiniProgramAdd 小程序素材-新增
func (s Service) SourceMiniProgramAdd(c ctx.Context, params *RequestSourceMini) (int64, error) {
	params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseID, params.AdminID)
	id, err := s.repo.SourceMiniProgramAdd(c, params)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// SourceVideoUpdate 小程序素材-编辑
func (s Service) SourceMiniProgramUpdate(c ctx.Context, params *RequestSourceMini) (int64, error) {
	id, err := s.repo.SourceMiniProgramUpdate(c, params)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// SourceMiniProgramList 小程序素材--列表
func (s Service) SourceMiniProgramList(c ctx.Context, search *RequestSourceSearch) (int64, *[]SourceMini, error) {
	total, list, err := s.repo.SourceMiniProgramList(c, search)
	if err != nil {
		return 0, nil, err
	}

	var resultList []SourceMini
	resultList = *list

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}

	return total, &resultList, err
}

// SourceMiniProgramDetail 小程序素材--详情
func (s Service) SourceMiniProgramDetail(cc ctx.Context, id int64) (*SourceMini, error) {
	entityFactory := func(ctx.Context) interface{} {
		return &SourceMini{}
	}

	option := &repos.DetailOption{
		TplPrefix: "t_source/source-mini-detail",
		Fn:        entityFactory,
	}
	vo, err := s.repo.DetailPlus(cc, id, option)

	if err != nil {
		return nil, err
	}

	if vo == nil {
		return nil, nil
	}

	return vo.(*SourceMini), nil
}

// SourceVideoAdd 视频素材-新增
func (s Service) SourceVideoAdd(c ctx.Context, params *RequestSourceVideo) (int64, error) {
	params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseID, params.AdminID)
	id, err := s.repo.SourceVideoAdd(c, params)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// SourceVideoUpdate 视频素材-编辑
func (s Service) SourceVideoUpdate(c ctx.Context, params *RequestSourceVideo) (int64, error) {
	id, err := s.repo.SourceVideoUpdate(c, params)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// SourceVideoDetail 视频素材--详情
func (s Service) SourceVideoDetail(cc ctx.Context, id int64) (*SourceVideo, error) {
	entityFactory := func(ctx.Context) interface{} {
		return &SourceVideo{}
	}

	option := &repos.DetailOption{
		TplPrefix: "t_source/source-video-detail",
		Fn:        entityFactory,
	}
	vo, err := s.repo.DetailPlus(cc, id, option)

	if err != nil {
		return nil, err
	}

	if vo == nil {
		return nil, nil
	}

	return vo.(*SourceVideo), nil
}

// SourceVideoList 视频素材--列表
func (s Service) SourceVideoList(c ctx.Context, search *RequestSourceSearch) (int64, *[]SourceVideo, error) {
	total, list, err := s.repo.SourceVideoList(c, search)
	if err != nil {
		return 0, nil, err
	}

	var resultList []SourceVideo
	resultList = *list

	for k, v := range resultList {
		res, _ := boss.GetAdminDetail(c, map[string]interface{}{"admin_id": v.AdminID})
		resultList[k].AdminName = res.Data.Username
	}

	return total, &resultList, err
}

// SourceDelete 删除素材
func (s Service) SourceDelete(cc ctx.Context, data *RequestSourceDelete) (int64, error) {
	if len(data.IDs) > 0 {
		return s.repo.SourceDelete(cc, data)
	}
	return 0, nil
}

// MoveGroup 移动分组
func (s Service) MoveGroup(cc ctx.Context, data *RequestGroupMove) (int64, error) {
	return s.repo.MoveGroup(cc, data)
}

// UpdateGroup 更新素材分组
func (s Service) UpdateGroup(cc ctx.Context, data *RequestGroupMove) (int64, error) {
	return s.repo.UpdateGroup(cc, data)
}

func (s Service) UpdateSourceMedia(cc ctx.Context, req UpdateSourceMediaRequest) error {
	return s.repo.UpdateSourceMedia(cc, req)
}

// MakeMediaForSource 生成素材
func (s Service) MakeMediaForSource(cc ctx.Context, enterpriseID int64, sourceID int64, sourceType string, sourceUrl string, sourceName string, domain string) UpdateSourceMediaRequest {
	reqUpdateSource := UpdateSourceMediaRequest{
		ID:           sourceID,
		MediaStatus:  1,
		UpdateAt:     time.Now().Unix(),
		SourceType:   sourceType,
		EnterpriseID: enterpriseID,
	}

	DrainageService = drainage.NewDrainage(drainageMysql.NewDrainage())
	mediaID, err := DrainageService.MakeMediaBySourceUrl(cc, enterpriseID, sourceUrl, sourceType, sourceName, domain)
	if err != nil {
		log.Info(cc).Msgf("企业[%d]生成企业微信临时素材出错，错误信息%+v", enterpriseID, err)
		reqUpdateSource.MediaStatus = 2
		reqUpdateSource.MediaFailReason = err.Error()
		return reqUpdateSource
	}

	reqUpdateSource.MediaExpireTime = time.Now().Unix() + 3*24*60*60
	reqUpdateSource.MediaID = mediaID
	return reqUpdateSource

}

// SourceMediaDetail 临时素材信息
func (s Service) SourceMediaDetail(ctx ctx.Context, enterpriseID int64, sourceID int64, sourceType string) (*MediaItemDetail, error) {
	return s.repo.SourceMediaDetail(ctx, enterpriseID, sourceID, sourceType)
}

func (s Service) SourceLevel(groupID int64) int {
	res, _ := orm.Service.Source.GetGroupID(source.GroupModel{})
	newMap := make(map[int64]source.GroupModel, 0)
	for _, v := range res {
		newMap[v.ID] = v
	}

	return sadfas(newMap, groupID, 1)
}

func sadfas(infos map[int64]source.GroupModel, groupID int64, i int) int {
	if infos[groupID].ParentId == 0 {
		return i
	} else {
		i = i + 1
		sadfas(infos, infos[groupID].ParentId, i)
	}
	return i
}

// SourceFriendAdd 朋友圈素材-新增
func (s Service) SourceFriendAdd(c ctx.Context, params *RequestSourceFriend) error {
	params.MemberID = resolveAdminUploaderMemberID(c, params.EnterpriseID, params.AdminID)
	params.Content = helper.CleanUserInput(params.Content)

	err := gormx.GetInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
		if params.Type == 1 && len(params.Images) > 0 {
			for _, image := range params.Images {
				sourceId := app.UID()
				err := tx.Table(sourceModels.SourceModel{}.TableName()).Create(sourceModels.SourceModel{
					MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: sourceId}, EnterpriseID: params.EnterpriseID},
					GroupID:        params.GroupID,
					Type:           params.Type,
					AdminID:        params.AdminID,
					MemberID:       params.MemberID,
					StartTime:      params.StartTime,
					EndTime:        params.EndTime,
				}).Error

				if err != nil {
					return err
				}

				var categoryRelationData []sourceModels.CategoryRelationModel

				for _, categoryId := range params.CategoryIds {
					categoryRelationData = append(categoryRelationData, sourceModels.CategoryRelationModel{
						MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: app.UID()}, EnterpriseID: params.EnterpriseID},
						CategoryID:     uint64(categoryId),
						SourceID:       uint64(sourceId),
						Type:           TypeFile,
					})
				}

				err = tx.Table(sourceModels.CategoryRelationModel{}.TableName()).CreateInBatches(categoryRelationData, len(categoryRelationData)).Error
				if err != nil {
					return err
				}

				err = tx.Table(sourceModels.FriendModel{}.TableName()).Create(sourceModels.FriendModel{
					MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: sourceId}, EnterpriseID: params.EnterpriseID},
					TextContent:    params.Content,
					Name:           image.Name,
					OriginalURL:    image.URL,
					Ext:            image.Ext,
					Size:           int(image.Size),
					Type:           params.Type,
				}).Error

				if err != nil {
					return err
				}
			}
		}

		if params.Type == 2 {
			sourceId := app.UID()
			err := tx.Table(sourceModels.SourceModel{}.TableName()).Create(sourceModels.SourceModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel:    models.BaseModel{ID: sourceId},
					EnterpriseID: params.EnterpriseID,
				},
				GroupID:   params.GroupID,
				Type:      params.Type,
				AdminID:   params.AdminID,
				MemberID:  params.MemberID,
				StartTime: params.StartTime,
				EndTime:   params.EndTime,
			}).Error

			if err != nil {
				return err
			}

			var categoryRelationData []sourceModels.CategoryRelationModel

			for _, categoryId := range params.CategoryIds {
				categoryRelationData = append(categoryRelationData, sourceModels.CategoryRelationModel{
					MicroBaseModel: models.MicroBaseModel{
						BaseModel:    models.BaseModel{ID: app.UID()},
						EnterpriseID: params.EnterpriseID,
					},
					CategoryID: uint64(categoryId),
					SourceID:   uint64(sourceId),
					Type:       TypeFile,
				})
			}

			err = tx.Table(sourceModels.CategoryRelationModel{}.TableName()).CreateInBatches(categoryRelationData, len(categoryRelationData)).Error
			if err != nil {
				return err
			}

			err = tx.Table(sourceModels.FriendModel{}.TableName()).Create(sourceModels.FriendModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel:    models.BaseModel{ID: sourceId},
					EnterpriseID: params.EnterpriseID,
				},
				TextContent:   params.Content,
				Name:          params.Name,
				VideoCoverURL: params.CoverUrl,
				OriginalURL:   params.Url,
				Ext:           params.Ext,
				Size:          int(params.Size),
				Type:          params.Type,
			}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

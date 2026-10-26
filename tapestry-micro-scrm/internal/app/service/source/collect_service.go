package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models"
	sourceModel "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"time"
)

type CollectService struct {
	collectDao *sourceDao.CollectRelationDao
}

func NewCollectService() *CollectService {
	return &CollectService{
		collectDao: sourceDao.NewCollectRelationDao(),
	}
}

func (s *CollectService) CollectSource(c ctx.Context, req sourceRequest.CollectReq) (sourceResp.CollectNumResp, error) {
	// 查询现有记录
	var collectResInfo sourceResp.CollectNumResp
	collectSourceList := req.SourceInfo
	if len(collectSourceList) == 0 {
		return collectResInfo, nil
	}
	isReturnCollectNum := false
	if len(collectSourceList) == 1 {
		isReturnCollectNum = true
	}
	for _, sourceInfo := range collectSourceList {
		whereMap := map[string]interface{}{
			"enterprise_id": req.EnterpriseId,
			"member_id":     req.MemberId,
			"type":          sourceInfo.Type,
		}
		if sourceInfo.Type < sourceModel.TypeFriend {
			whereMap["source_id"] = sourceInfo.SourceId
		}
		if sourceInfo.Type == sourceModel.TypeFriend {
			whereMap["merge_id"] = sourceInfo.MergeId
		}
		if sourceInfo.Type == sourceModel.TypeProduct {
			whereMap["product_code"] = sourceInfo.ProductCode
		}
		list, err := s.collectDao.GetMany(c, whereMap)
		if err != nil {
			return collectResInfo, err
		}
		// 如果是取消收藏且记录不存在，直接返回
		if req.IsCollect <= sourceModel.SourceCollectNo && len(list) == 0 {
			return collectResInfo, nil
		}
		// 如果是新增收藏且记录不存在，先创建记录
		sourceIds := make([]int64, 0)
		mergeIds := make([]int64, 0)
		productCodes := make([]string, 0)
		if req.IsCollect == sourceModel.SourceCollectYes && len(list) == 0 {
			sourceImageMerge := make(map[int64]int64)
			if sourceInfo.Type < sourceModel.TypeFriend {
				sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]interface{}{
					"enterprise_id": req.EnterpriseId,
					"id":            sourceInfo.SourceId,
					"type":          sourceInfo.Type,
				})
				if err != nil {
					return collectResInfo, err
				}
				if len(sourceList) == 0 {
					return collectResInfo, nil
				}
				for _, source := range sourceList {
					sourceIds = append(sourceIds, source.ID)
					if sourceInfo.Type == sourceModel.TypeImg {
						sourceImageMerge[source.ID] = source.MergeID
					}
				}
			}
			batchInsertData := make([]sourceModel.CollectRelationModel, 0)
			if len(sourceIds) > 0 {
				for _, value := range sourceIds {
					id := app.UID()
					mergeId := int64(0)
					if item, ok := sourceImageMerge[value]; ok {
						mergeId = item
					}
					insertData := sourceModel.CollectRelationModel{
						MicroBaseModel: models.MicroBaseModel{
							BaseModel: models.BaseModel{
								ID:        id,
								CreatedAt: time.Now().Unix(),
								UpdatedAt: time.Now().Unix(),
								DeletedAt: 0,
							},
							EnterpriseID: req.EnterpriseId,
						},
						MemberId:  req.MemberId,
						Type:      sourceInfo.Type,
						SourceID:  value,
						IsCollect: req.IsCollect,
						MergeId:   mergeId,
					}
					batchInsertData = append(batchInsertData, insertData)
				}
			}
			if len(sourceInfo.ProductCode) > 0 {
				for _, value := range sourceInfo.ProductCode {
					id := app.UID()
					insertData := sourceModel.CollectRelationModel{
						MicroBaseModel: models.MicroBaseModel{
							BaseModel: models.BaseModel{
								ID:        id,
								CreatedAt: time.Now().Unix(),
								UpdatedAt: time.Now().Unix(),
								DeletedAt: 0,
							},
							EnterpriseID: req.EnterpriseId,
						},
						MemberId:    req.MemberId,
						Type:        sourceInfo.Type,
						ProductCode: value,
						IsCollect:   req.IsCollect,
					}
					batchInsertData = append(batchInsertData, insertData)
					productCodes = append(productCodes, value)
				}
			}
			if len(sourceInfo.MergeId) > 0 {
				for _, value := range sourceInfo.MergeId {
					id := app.UID()
					insertData := sourceModel.CollectRelationModel{
						MicroBaseModel: models.MicroBaseModel{
							BaseModel: models.BaseModel{
								ID:        id,
								CreatedAt: time.Now().Unix(),
								UpdatedAt: time.Now().Unix(),
								DeletedAt: 0,
							},
							EnterpriseID: req.EnterpriseId,
						},
						MemberId:  req.MemberId,
						Type:      sourceInfo.Type,
						MergeId:   value,
						IsCollect: req.IsCollect,
					}
					batchInsertData = append(batchInsertData, insertData)
					mergeIds = append(mergeIds, value)
				}
			}
			if len(batchInsertData) == 0 {
				return collectResInfo, nil
			}
			if err := s.collectDao.CreateInBatchesData(c, batchInsertData); err != nil {
				return collectResInfo, err
			}
		} else {
			sourceIds = sourceIds[:0]
			mergeIds = mergeIds[:0]
			productCodes = productCodes[:0]
			isCollect := helper.RuleThree(req.IsCollect == sourceModel.SourceCollectYes, sourceModel.SourceCollectYes, sourceModel.SourceCollectNo)
			for _, one := range list {
				if isCollect == int(one.IsCollect) {
					continue
				}
				// 更新收藏状态
				err = s.collectDao.UpdateData(c, map[string]interface{}{
					"id": one.ID,
				}, map[string]interface{}{
					"is_collect": isCollect,
				})
				if err != nil {
					return collectResInfo, err
				}
				sourceIds = append(sourceIds, one.SourceID)
				mergeIds = append(mergeIds, one.MergeId)
				productCodes = append(productCodes, one.ProductCode)
			}
		}
		if sourceInfo.Type < sourceModel.TypeFriend {
			if req.IsCollect == sourceModel.SourceCollectYes {
				err = sourceDao.NewSourceDao().Increment(c, map[string]interface{}{
					"id": sourceIds,
				}, "collect_num", 1)
			} else {
				err = sourceDao.NewSourceDao().Decrement(c, map[string]interface{}{
					"id": sourceIds,
				}, "collect_num", 1)
			}
			if isReturnCollectNum {
				collectNumList, err := sourceDao.NewSourceDao().GetSourceCollectNum(c, req.EnterpriseId, sourceInfo.Type, sourceIds)
				if err != nil {
					return collectResInfo, err
				}
				var imageMergeId int64
				for _, value := range collectNumList {
					collectResInfo.CollectNum = value.CollectNum
					collectResInfo.SourceId = value.Id
					if value.MergeId > 0 && sourceInfo.Type == sourceModel.TypeImg {
						imageMergeId = value.MergeId
					}
				}
				collectResInfo.SourceType = sourceInfo.Type
				sourceCollectNum := make([]sourceResp.SourceCollectNumResp, 0)
				collectResInfo.SourceList = sourceCollectNum
				if imageMergeId > 0 {
					total, err := sourceDao.NewCollectRelationDao().GetImageMergeCollectNumByMergeId(c, req.EnterpriseId, imageMergeId)
					if err != nil {
						return collectResInfo, err
					}
					for _, value := range collectNumList {
						item := sourceResp.SourceCollectNumResp{
							CollectNum: value.CollectNum,
							SourceId:   value.Id,
						}
						sourceCollectNum = append(sourceCollectNum, item)
					}
					collectResInfo.CollectNum = total
					collectResInfo.MergeId = imageMergeId
					collectResInfo.SourceList = sourceCollectNum
				}

			}

		}
		if sourceInfo.Type == sourceModel.TypeFriend {
			if req.IsCollect == sourceModel.SourceCollectYes {
				err = sourceDao.NewSourceDao().Increment(c, map[string]interface{}{
					"merge_id": mergeIds,
					"type":     sourceInfo.Type,
				}, "collect_num", 1)
			} else {
				err = sourceDao.NewSourceDao().Decrement(c, map[string]interface{}{
					"merge_id": mergeIds,
					"type":     sourceInfo.Type,
				}, "collect_num", 1)
			}
			if isReturnCollectNum {
				collectNumList, err := sourceDao.NewSourceDao().GetSourceCollectNumByMergeId(c, req.EnterpriseId, sourceInfo.Type, mergeIds)
				if err != nil {
					return collectResInfo, err
				}
				for _, value := range collectNumList {
					collectResInfo.CollectNum = value.CollectNum
					collectResInfo.SourceId = value.Id
				}
				collectResInfo.SourceType = sourceInfo.Type

			}

		}
		if sourceInfo.Type == sourceModel.TypeProduct {
			// 增加转发次数
			if req.IsCollect == sourceModel.SourceCollectYes {
				err = sourceDao.NewBaozunCollectDao().Increment(c, map[string]interface{}{
					"product_code": productCodes,
				}, "collect_num", 1)
			} else {
				err = sourceDao.NewBaozunCollectDao().Decrement(c, map[string]interface{}{
					"product_code": productCodes,
				}, "collect_num", 1)
			}
			if isReturnCollectNum {
				collectNumList, err := sourceDao.NewBaozunCollectDao().GetCollectNumByProductCode(c, req.EnterpriseId, productCodes)
				if err != nil {
					return collectResInfo, err
				}
				for _, value := range collectNumList {
					collectResInfo.CollectNum = value.CollectNum
					collectResInfo.ProductCode = value.ProductCode
				}
				collectResInfo.SourceType = sourceInfo.Type
			}

		}
		if err != nil {
			return collectResInfo, err
		}
	}
	return collectResInfo, nil
}

// List 收藏列表
func (s *CollectService) List(c ctx.Context, req sourceRequest.CollectListReq) ([]interface{}, int, error) {
	var textSourceId, imageSourceId, linkSourceId, miniProgramSourceId, videoSourceId, fileSourceId, friendSourceId, productSourceId []int64
	result := make([]interface{}, 0)
	whereMap := map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"member_id":     req.MemberId,
		"is_collect":    sourceModel.SourceCollectYes,
	}
	if req.Type >= sourceModel.TypeText {
		whereMap["type"] = req.Type
	}
	many, err := s.collectDao.GetMany(c, whereMap)
	if err != nil {
		return result, 0, err
	}
	for _, one := range many {
		switch one.Type {
		case sourceModel.TypeText:
			textSourceId = append(textSourceId, one.SourceID)
		case sourceModel.TypeImg:
			imageSourceId = append(imageSourceId, one.SourceID)
		case sourceModel.TypeLink:
			linkSourceId = append(linkSourceId, one.SourceID)
		case sourceModel.TypeMiniProgram:
			miniProgramSourceId = append(miniProgramSourceId, one.SourceID)
		case sourceModel.TypeVideo:
			videoSourceId = append(videoSourceId, one.SourceID)
		case sourceModel.TypeFile:
			fileSourceId = append(fileSourceId, one.SourceID)
		case sourceModel.TypeFriend:
			friendSourceId = append(friendSourceId, one.SourceID)
		case sourceModel.TypeProduct:
			productSourceId = append(productSourceId, one.SourceID)
		default:

		}
	}
	if len(req.CategoryIds) > 0 {
		textSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeText, req.CategoryIds, textSourceId)
		imageSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeImg, req.CategoryIds, imageSourceId)
		linkSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeLink, req.CategoryIds, linkSourceId)
		miniProgramSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeMiniProgram, req.CategoryIds, miniProgramSourceId)
		videoSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeVideo, req.CategoryIds, videoSourceId)
		fileSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeFile, req.CategoryIds, fileSourceId)
		friendSourceId, err = s.getSourceIdByCategory(c, req.EnterpriseId, sourceModel.TypeFriend, req.CategoryIds, friendSourceId)
		if len(textSourceId)+len(imageSourceId)+len(linkSourceId)+len(miniProgramSourceId)+len(videoSourceId)+len(fileSourceId)+len(friendSourceId) == 0 {
			return result, 0, err
		}
	}
	textSourceMap := make(map[int64]sourceModel.ResponseTextList)
	if len(textSourceId) > 0 {
		textSourceList, _, err := s.collectDao.GetTextSource(c, req, textSourceId)
		if err != nil {
			return result, 0, err
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(textSourceList))
		for i := range textSourceList {
			identityItems = append(identityItems, &textSourceList[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range textSourceList {
			textSourceMap[one.ID] = one
		}

	}
	imageSourceMap := make(map[int64]sourceResp.ImageSelectList)
	if len(imageSourceId) > 0 {
		imageSourceList, total, err := s.collectDao.GetImageSource(c, req, imageSourceId)
		if err != nil {
			return result, 0, err
		}

		resp := make([]sourceResp.ImageSelectList, 0)
		for _, list := range imageSourceList {
			resp = append(resp, sourceResp.ImageSelectList{
				ResponseSelectList: list.ResponseSelectList,
				Name:               list.Name,
				URL:                list.URL,
				Size:               list.Size,
				MediaID:            list.MediaID,
				MediaExpireTime:    list.MediaExpireTime,
				MediaStatus:        list.MediaStatus,
				DepartmentId:       list.DepartmentId,
				SourceGroupID:      list.SourceGroupID,
				MergeId:            list.MergeId,
				MergeImage:         make([]sourceResp.MergeImage, 0),
				DownloadTimes:      list.DownloadTimes,
				ForwardingTimes:    list.ForwardingTimes,
			})
		}

		if total > 0 {
			var mergeIds []int64
			for _, v := range resp {
				if v.MergeId > 0 {
					mergeIds = append(mergeIds, v.MergeId)
				}
			}

			sourceMap := make(map[int64][]sourceResp.MergeImage)
			// 查询组合图
			if len(mergeIds) > 0 {
				sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": req.EnterpriseId, "merge_id": mergeIds})
				if err != nil {
					return nil, 0, err
				}

				var sourceIds []int64
				for _, val := range sourceList {
					sourceIds = append(sourceIds, val.ID)
				}

				sourceImageMap := make(map[int64]sourceModel.ImageModel)
				sourceImageList, err := sourceDao.NewImageDao().GetMany(c, req.EnterpriseId, map[string]any{"id": sourceIds})
				if err != nil {
					return nil, 0, err
				}

				for _, image := range sourceImageList {
					sourceImageMap[image.ID] = image
				}

				for _, val := range sourceList {
					tmpImage := sourceResp.MergeImage{
						Id:              sourceImageMap[val.ID].ID,
						Name:            sourceImageMap[val.ID].Name,
						URL:             sourceImageMap[val.ID].URL,
						Sort:            val.Sort,
						DownloadTimes:   val.DownloadTimes,
						ForwardingTimes: val.ForwardingTimes,
						CollectNum:      val.CollectNum,
					}

					sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
				}
			}

			for k, v := range resp {
				// 组装合并图片
				if sourceImage, ok := sourceMap[v.MergeId]; ok {
					sort.Slice(sourceImage, func(i, j int) bool {
						return sourceImage[i].Sort < sourceImage[j].Sort
					})
					resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
				}
			}
		}

		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(resp))
		for i := range resp {
			identityItems = append(identityItems, &resp[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range resp {
			imageSourceMap[one.ID] = one
		}

	}
	linkSourceMap := make(map[int64]sourceModel.ResponseLinkList)
	if len(linkSourceId) > 0 {
		linkSourceList, _, err := s.collectDao.GetLinkSource(c, req, linkSourceId)
		if err != nil {
			return result, 0, err
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(linkSourceList))
		for i := range linkSourceList {
			identityItems = append(identityItems, &linkSourceList[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range linkSourceList {
			linkSourceMap[one.ID] = one
		}

	}
	miniProgramSourceMap := make(map[int64]sourceModel.ResponseMiniProgramList)
	if len(miniProgramSourceId) > 0 {
		miniProgramSourceList, _, err := s.collectDao.GetMiniProgramSource(c, req, miniProgramSourceId)
		if err != nil {
			return result, 0, err
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(miniProgramSourceList))
		for i := range miniProgramSourceList {
			identityItems = append(identityItems, &miniProgramSourceList[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range miniProgramSourceList {
			miniProgramSourceMap[one.ID] = one
		}
	}
	videoSourceMap := make(map[int64]sourceModel.ResponseVideoList)
	if len(videoSourceId) > 0 {
		videoSourceList, _, err := s.collectDao.GetVideoSource(c, req, videoSourceId)
		if err != nil {
			return result, 0, err
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(videoSourceList))
		for i := range videoSourceList {
			identityItems = append(identityItems, &videoSourceList[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range videoSourceList {
			videoSourceMap[one.ID] = one
		}
	}
	fileSourceMap := make(map[int64]sourceModel.ResponseFileList)
	if len(fileSourceId) > 0 {
		fileSourceList, _, err := s.collectDao.GetFileSource(c, req, fileSourceId)
		if err != nil {
			return result, 0, err
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(fileSourceList))
		for i := range fileSourceList {
			identityItems = append(identityItems, &fileSourceList[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range fileSourceList {
			fileSourceMap[one.ID] = one
		}
	}
	friendSourceMap := make(map[int64]sourceResp.FriendSelectList)
	if len(friendSourceId) > 0 {
		friendSourceList, total, err := s.collectDao.GetFriendSource(c, req, friendSourceId)
		if err != nil {
			return result, 0, err
		}

		resp := make([]sourceResp.FriendSelectList, 0)
		for _, list := range friendSourceList {
			resp = append(resp, sourceResp.FriendSelectList{
				ResponseSelectList: list.ResponseSelectList,
				TextContent:        list.TextContent,
				Name:               list.Name,
				VideoCoverURL:      list.VideoCoverURL,
				OriginalURL:        list.OriginalURL,
				CompressionURL:     list.CompressionURL,
				Ext:                list.Ext,
				Size:               list.Size,
				Type:               list.Type,
				Sort:               list.Sort,
				MergeId:            list.MergeId,
				DownloadTimes:      list.DownloadTimes,
				ForwardingTimes:    list.ForwardingTimes,
				CopyTimes:          list.CopyTimes,
				MergeImage:         make([]sourceResp.MergeImage, 0),
				CollectNum:         list.CollectNum,
			})
		}

		if total > 0 {
			var mergeIds []int64

			for _, v := range resp {
				if v.MergeId > 0 {
					mergeIds = append(mergeIds, v.MergeId)
				}
			}

			sourceMap := make(map[int64][]sourceResp.MergeImage)
			// 查询组合图
			if len(mergeIds) > 0 {
				sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": req.EnterpriseId, "merge_id": mergeIds})
				if err != nil {
					return nil, 0, err
				}

				var sourceIds []int64
				for _, val := range sourceList {
					sourceIds = append(sourceIds, val.ID)
				}

				sourceFriendMap := make(map[int64]sourceModel.FriendModel)
				sourceFriendList, err := sourceDao.NewFriendDao().GetMany(c, req.EnterpriseId, map[string]any{"id": sourceIds})
				if err != nil {
					return nil, 0, err
				}

				for _, friend := range sourceFriendList {
					sourceFriendMap[friend.ID] = friend
				}

				for _, val := range sourceList {
					tmp := sourceResp.MergeImage{
						Id:              sourceFriendMap[val.ID].ID,
						Name:            sourceFriendMap[val.ID].Name,
						URL:             sourceFriendMap[val.ID].OriginalURL,
						Sort:            val.Sort,
						DownloadTimes:   val.DownloadTimes,
						ForwardingTimes: val.ForwardingTimes,
					}
					sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmp)
				}
			}

			for k, v := range resp {
				// 组装合并图片
				if sourceImage, ok := sourceMap[v.MergeId]; ok {
					resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
				}
			}
		}
		identityItems := make([]*sourceModel.ResponseSelectList, 0, len(resp))
		for i := range resp {
			identityItems = append(identityItems, &resp[i].ResponseSelectList)
		}
		enrichUploaderIdentity(c, req.EnterpriseId, identityItems)
		for _, one := range resp {
			friendSourceMap[one.ID] = one
		}
	}

	for _, one := range many {
		switch one.Type {
		case sourceModel.TypeText:
			if item, ok := textSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeImg:
			if item, ok := imageSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeLink:
			if item, ok := linkSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeMiniProgram:
			if item, ok := miniProgramSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeVideo:
			if item, ok := videoSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeFile:
			if item, ok := fileSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		case sourceModel.TypeFriend:
			if item, ok := friendSourceMap[one.SourceID]; ok {
				result = append(result, item)
			}
		default:

		}
	}

	return result, len(result), nil
}

func (s *CollectService) getSourceIdByCategory(c ctx.Context, enterpriseId int64, sourceType int, categoryIds, sourceIds []int64) ([]int64, error) {
	return sourceDao.NewSourceDao().GetSourceIdByCategoryIds(c, enterpriseId, sourceType, categoryIds, sourceIds)
}

package source

import (
	"github.com/jinzhu/copier"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models/source"
	sourceModel "micro-scrm/internal/app/models/source"
	departmentService "micro-scrm/internal/app/service/department"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"time"
)

type GroupRelationService struct {
	groupRelationDao    *sourceDao.GroupRelationDao
	sourceGroupDao      *sourceDao.SourceGroupDao
	categoryRelationDao *sourceDao.CategoryRelationDao
}

func NewGroupRelationService() GroupRelationService {
	return GroupRelationService{
		groupRelationDao:    sourceDao.NewGroupRelationDao(),
		sourceGroupDao:      sourceDao.NewSourceGroupDao(),
		categoryRelationDao: sourceDao.NewCategoryRelationDao(),
	}
}

// GroupRelationGetList 查看素材分组归属门店列表
func (s GroupRelationService) GroupRelationGetList(c app.Ctx, request sourceRequest.GroupRelationListRequest) (error, []sourceModel.GroupRelationList, int64) {

	return s.groupRelationDao.GetList(c, request.EnterpriseID, request.SourceGroupId, request.Page, request.Size)
}

// GroupRelationSave 素材分组归属门店保存
func (s GroupRelationService) GroupRelationSave(c app.Ctx, request sourceRequest.GroupRelationSaveRequest) error {

	var groupRelation []*source.GroupRelationModel

	groupIds := request.SourceGroupIds
	thirdIds := request.ThirdIds
	enterpriseId := request.EnterpriseID
	for _, groupId := range groupIds {
		for _, thirdId := range thirdIds {
			temp := &source.GroupRelationModel{
				ID:            app.UIDContext(c),
				CreatedAt:     time.Now().Unix(),
				EnterpriseID:  enterpriseId,
				ThirdId:       thirdId,
				SourceGroupID: groupId,
			}
			groupRelation = append(groupRelation, temp)
		}
	}

	if err := s.groupRelationDao.DeleteForSourceGroupIds(c, groupIds); err != nil {
		return err
	}
	if len(groupRelation) == 0 {
		return nil
	}
	if err := s.groupRelationDao.Save(c, groupRelation); err != nil {
		return err
	}
	return nil
}

func initialGroupThirdIDs(parentID int64, allDepartmentThirdIDs []int, parentRelations []source.GroupRelationModel) []int {
	thirdIDSet := make(map[int]struct{})
	if parentID == 0 {
		for _, thirdID := range allDepartmentThirdIDs {
			thirdIDSet[thirdID] = struct{}{}
		}
	} else {
		for _, relation := range parentRelations {
			thirdIDSet[relation.ThirdId] = struct{}{}
		}
	}
	thirdIDs := make([]int, 0, len(thirdIDSet))
	for thirdID := range thirdIDSet {
		thirdIDs = append(thirdIDs, thirdID)
	}
	sort.Ints(thirdIDs)
	return thirdIDs
}

// GroupRelationSaveInitial applies the product rule used only when a group is
// created: root groups are visible to all departments and child groups copy
// their direct parent's current permissions. Later edits remain independent.
func (s GroupRelationService) GroupRelationSaveInitial(c app.Ctx, enterpriseID, groupID, parentID int64) error {
	request := sourceRequest.GroupRelationSaveRequest{
		EnterpriseID:   enterpriseID,
		SourceGroupIds: []int64{groupID},
	}
	if parentID == 0 {
		return s.GroupRelationSaveAll(c, request)
	}
	parentRelations, err := s.groupRelationDao.GetMany(c, enterpriseID, map[string]interface{}{
		"source_group_id": parentID,
		"deleted_at":      0,
	})
	if err != nil {
		return err
	}
	request.ThirdIds = initialGroupThirdIDs(parentID, nil, parentRelations)
	return s.GroupRelationSave(c, request)
}

// GroupRelationSaveAll 素材分组归属门店保存
func (s GroupRelationService) GroupRelationSaveAll(c app.Ctx, request sourceRequest.GroupRelationSaveRequest) error {
	err, res := departmentService.NewDepartmentsService().GetAllDepartments(c, request.EnterpriseID)

	if err != nil {
		return err
	}
	allDepartmentThirdIDs := make([]int, 0, len(res))
	for _, department := range res {
		allDepartmentThirdIDs = append(allDepartmentThirdIDs, department.ThirdID)
	}
	request.ThirdIds = initialGroupThirdIDs(0, allDepartmentThirdIDs, nil)
	return s.GroupRelationSave(c, request)
}

// SourceListForId 分组树结构
func (s GroupRelationService) SourceListForId(c app.Ctx, request sourceRequest.GroupRelationFindRequest) (error, []*sourceModel.SourceGroupTreeList) {

	err, resAllList := s.sourceGroupDao.GetAllList(c, request.EnterpriseID)
	err, res := s.groupRelationDao.GetListForId(c, request.EnterpriseID, request.SourceGroupId, request.ThirdIds)

	if err != nil {
		return err, nil
	}

	treeMap := make(map[int64]*sourceModel.SourceGroupTreeList, 0)
	treeList := make([]*sourceModel.SourceGroupTreeList, 0)
	newList := make([]*sourceModel.SourceGroupTreeList, 0)
	copier.Copy(&treeList, &resAllList)
	for _, v := range treeList {
		treeMap[v.ID] = v
	}

	treeMapTemp := make(map[int64]*sourceModel.GroupRelationTreeList, 0)
	treeListTemp := make([]*sourceModel.GroupRelationTreeList, 0)
	copier.Copy(&treeListTemp, &res)

	for _, v := range treeListTemp {
		treeMapTemp[v.ID] = v
	}

	for _, sourceList := range treeList {

		if sourceList.ParentID == 0 {
			if _, ok := treeMapTemp[sourceList.ID]; ok {
				sourceList.IsShow = 1
			}
			newList = append(newList, sourceList)
		} else {
			if val, ok := treeMap[sourceList.ParentID]; ok {
				if val.Children == nil {
					val.Children = make([]*sourceModel.SourceGroupTreeList, 0)
				}
				if _, ok1 := treeMapTemp[sourceList.ID]; ok1 {
					sourceList.IsShow = 1
				}
				val.Children = append(val.Children, sourceList)
			}
		}
	}

	return nil, newList

}

func (s GroupRelationService) SourceListForId1(c app.Ctx, request sourceRequest.GroupRelationFindRequest) (error, []*sourceModel.GroupRelationTreeList) {

	err, res := s.groupRelationDao.GetListForId(c, request.EnterpriseID, request.SourceGroupId, request.ThirdIds)

	if err != nil || res == nil {
		return err, nil
	}
	treeMap := make(map[int64]*sourceModel.GroupRelationTreeList, 0)
	treeList := make([]*sourceModel.GroupRelationTreeList, 0)
	newList := make([]*sourceModel.GroupRelationTreeList, 0)
	copier.Copy(&treeList, &res)

	for _, v := range treeList {
		treeMap[v.ID] = v
	}

	temp := false
	for _, sourceList := range treeList {
		if sourceList.ParentId == 0 {
			temp = true
			break
		}
	}

	if temp == true {
		for _, sourceList := range treeList {

			if sourceList.ParentId == 0 {
				newList = append(newList, sourceList)
			} else {
				if val, ok := treeMap[sourceList.ParentId]; ok {
					if val.Children == nil {
						val.Children = make([]*sourceModel.GroupRelationTreeList, 0)
					}
					val.Children = append(val.Children, sourceList)
				}
			}
		}
	} else {
		for _, sourceList := range treeList {

			if sourceList.ParentId == treeList[0].ParentId {
				newList = append(newList, sourceList)
			} else {
				if val, ok := treeMap[sourceList.ParentId]; ok {
					if val.Children == nil {
						val.Children = make([]*sourceModel.GroupRelationTreeList, 0)
					}
					val.Children = append(val.Children, sourceList)
				}
			}
		}
	}

	return nil, newList
}

// TextSelectList 查看文本素材分组归属门店列表
func (s GroupRelationService) TextSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceModel.ResponseTextList, total int64, err error) {

	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (len(collectSourceIDs) == 0 && collectSearch) {
		return
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeText,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	res, total, err = s.groupRelationDao.GetTextSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
		res[k].CollectAt = collectTimeMap[v.ID]
	}
	if collectSearch {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CollectAt > res[j].CollectAt
		})
	}
	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(res))
	for i := range res {
		identityItems = append(identityItems, &res[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return res, total, nil
}

// ImageSelectList 查看图片素材分组归属门店列表
func (s GroupRelationService) ImageSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) ([]sourceResp.ImageSelectList, int64, error) {
	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	var err error
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeImg,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	imgList, total, err := s.groupRelationDao.GetImageSelectList(c, request, sourceIDs, sourceGroupIds, collectSearch)
	resp := make([]sourceResp.ImageSelectList, 0)
	for _, list := range imgList {
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
		imgMergeCollect, err := sourceDao.NewCollectRelationDao().GetImageMergeCollectNum(c, request.EnterpriseID)
		if err != nil {
			return nil, 0, err
		}
		imgMergeCollectMap := make(map[int64]int, 0)
		for _, v := range imgMergeCollect {
			imgMergeCollectMap[v.MergeId] = v.Num
		}
		var mergeIds []int64
		// 先遍历resp找出所有mergeId
		for i := 0; i < len(resp); {
			v := resp[i]
			if v.MergeId > 0 {
				if helper.InArray(v.MergeId, mergeIds) {
					// 如果mergeId已存在，从resp中删除该元素
					resp = append(resp[:i], resp[i+1:]...)
					continue // 不增加i，因为删除了当前元素
				}
				mergeIds = append(mergeIds, v.MergeId)
			}
			i++
		}

		sourceMap := make(map[int64][]sourceResp.MergeImage)
		// 查询组合图
		if len(mergeIds) > 0 {
			sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
			if err != nil {
				return nil, 0, err
			}

			var sourceIds []int64
			for _, val := range sourceList {
				sourceIds = append(sourceIds, val.ID)
			}

			sourceImageMap := make(map[int64]source.ImageModel)
			sourceImageList, err := sourceDao.NewImageDao().GetMany(c, request.EnterpriseID, map[string]any{"id": sourceIds})
			if err != nil {
				return nil, 0, err
			}

			for _, image := range sourceImageList {
				sourceImageMap[image.ID] = image
			}

			for _, val := range sourceList {
				isCollect := sourceModel.SourceCollectNo
				if helper.InArray(val.ID, collectSourceIDs) {
					isCollect = sourceModel.SourceCollectYes
				}
				if collectSearch && isCollect == sourceModel.SourceCollectNo { // 我的收藏页面，组合图片的收藏中只收藏了其中几个图片，则为收藏的不展示
					continue
				}
				tmpImage := sourceResp.MergeImage{
					Id:              sourceImageMap[val.ID].ID,
					Name:            sourceImageMap[val.ID].Name,
					URL:             sourceImageMap[val.ID].URL,
					Sort:            val.Sort,
					DownloadTimes:   val.DownloadTimes,
					ForwardingTimes: val.ForwardingTimes,
					IsCollect:       isCollect,
					CollectNum:      val.CollectNum,
				}

				sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
			}
		}

		for k, v := range resp {
			// 组装合并图片
			if sourceImage, ok := sourceMap[v.MergeId]; ok {
				resp[k].CollectNum = imgMergeCollectMap[v.MergeId]
				resp[k].IsCollect = sourceModel.SourceCollectYes
				for _, image := range sourceImage {
					if image.IsCollect == sourceModel.SourceCollectNo { // 如果有一个没收藏，就不是全部收藏了
						resp[k].IsCollect = sourceModel.SourceCollectNo
						break
					}
				}
				sort.Slice(sourceImage, func(i, j int) bool {
					return sourceImage[i].Sort < sourceImage[j].Sort
				})
				resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
				resp[k].URL = sourceImage[0].URL
			} else {
				if helper.InArray(v.ID, collectSourceIDs) {
					resp[k].IsCollect = sourceModel.SourceCollectYes
				}
			}
			resp[k].CollectAt = collectTimeMap[v.ID]

		}
		if collectSearch {
			sort.Slice(resp, func(i, j int) bool {
				return resp[i].CollectAt > resp[j].CollectAt
			})
			total = int64(len(resp))
		}
	}

	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(resp))
	for i := range resp {
		identityItems = append(identityItems, &resp[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return resp, total, err
}

func (s GroupRelationService) LinkSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceModel.ResponseLinkList, total int64, err error) {

	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeLink,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	res, total, err = s.groupRelationDao.GetLinkSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
		res[k].CollectAt = collectTimeMap[v.ID]
	}
	if collectSearch {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CollectAt > res[j].CollectAt
		})
	}
	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(res))
	for i := range res {
		identityItems = append(identityItems, &res[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return res, total, nil
}

func (s GroupRelationService) MiniProgramSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceModel.ResponseMiniProgramList, total int64, err error) {

	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeMiniProgram,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	res, total, err = s.groupRelationDao.GetMiniProgramSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
		res[k].CollectAt = collectTimeMap[v.ID]
	}
	if collectSearch {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CollectAt > res[j].CollectAt
		})
	}
	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(res))
	for i := range res {
		identityItems = append(identityItems, &res[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return res, total, nil
}

func (s GroupRelationService) VideoSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceModel.ResponseVideoList, total int64, err error) {
	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeVideo,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	res, total, err = s.groupRelationDao.GetVideoSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
		res[k].CollectAt = collectTimeMap[v.ID]
	}
	if collectSearch {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CollectAt > res[j].CollectAt
		})
	}
	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(res))
	for i := range res {
		identityItems = append(identityItems, &res[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return res, total, nil
}

func (s GroupRelationService) FileSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceModel.ResponseFileList, total int64, err error) {
	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeFile,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.SourceID] = v.UpdatedAt
	}
	res, total, err = s.groupRelationDao.GetFileSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
		res[k].CollectAt = collectTimeMap[v.ID]
	}
	if collectSearch {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CollectAt > res[j].CollectAt
		})
	}
	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(res))
	for i := range res {
		identityItems = append(identityItems, &res[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return res, total, nil
}

func (s GroupRelationService) FriendSelectList(c app.Ctx, request sourceRequest.RequestSelectList, collectSearch bool) (res []sourceResp.FriendSelectList, total int64, err error) {
	var sourceIDs []int64
	var sourceGroupIds []int64
	var collectSourceIDs []int64
	var mergeIds []int64
	if collectSearch {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getCollectSelectParams(c, request)
		mergeIds = collectSourceIDs
	} else {
		sourceIDs, sourceGroupIds, collectSourceIDs, err = s.getSelectParams(c, request)
	}
	if err != nil || (collectSearch && len(collectSourceIDs) == 0) {
		return nil, 0, err
	}
	many, err := sourceDao.NewCollectRelationDao().GetMany(c, map[string]any{
		"enterprise_id": request.EnterpriseID,
		"member_id":     request.MemberID,
		"type":          sourceModel.TypeFriend,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, 0, err
	}
	collectTimeMap := make(map[int64]int64, 0)
	for _, v := range many {
		collectTimeMap[v.MergeId] = v.UpdatedAt
	}
	friendList, total, err := s.groupRelationDao.GetFriendSelectList(c, request, sourceIDs, sourceGroupIds, mergeIds)

	resp := make([]sourceResp.FriendSelectList, 0)
	for _, list := range friendList {
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
			sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
			if err != nil {
				return nil, 0, err
			}

			var sourceIds []int64
			for _, val := range sourceList {
				sourceIds = append(sourceIds, val.ID)
			}

			sourceFriendMap := make(map[int64]source.FriendModel)
			sourceFriendList, err := sourceDao.NewFriendDao().GetMany(c, request.EnterpriseID, map[string]any{"id": sourceIds})
			if err != nil {
				return nil, 0, err
			}

			for _, friend := range sourceFriendList {
				if friend.OriginalURL != "" && friend.Type == sourceModel.FriendTypeImg {
					sourceFriendMap[friend.ID] = friend
				}
			}

			for _, val := range sourceList {
				if _, ok := sourceFriendMap[val.ID]; !ok {
					continue
				}
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
			if helper.InArray(v.MergeId, collectSourceIDs) {
				resp[k].IsCollect = sourceModel.SourceCollectYes
			}
			resp[k].CollectAt = collectTimeMap[v.MergeId]
		}
		if collectSearch {
			sort.Slice(resp, func(i, j int) bool {
				return resp[i].CollectAt > resp[j].CollectAt
			})
		}
	}

	identityItems := make([]*sourceModel.ResponseSelectList, 0, len(resp))
	for i := range resp {
		identityItems = append(identityItems, &resp[i].ResponseSelectList)
	}
	enrichUploaderIdentity(c, request.EnterpriseID, identityItems)
	return resp, total, err
}

// 获取查询参数
func (s GroupRelationService) getSelectParams(c app.Ctx, request sourceRequest.RequestSelectList) ([]int64, []int64, []int64, error) {
	var (
		err              error
		sourceIDs        []int64
		sourceGroupIds   []int64
		collectSourceIDs []int64
	)

	// 所属门店下 分组的sourceGroupIDs
	var whereGroupRelation = map[string]interface{}{
		"third_id": request.ThirdIds,
	}
	if len(request.GroupIDList) > 0 {
		whereGroupRelation = map[string]interface{}{
			"third_id":        request.ThirdIds,
			"source_group_id": request.GroupIDList,
		}
	}
	sourceGroupIds, err = s.groupRelationDao.GetManySourceGroupIds(c, request.EnterpriseID, whereGroupRelation)
	if err != nil {
		return sourceIDs, sourceGroupIds, collectSourceIDs, err
	}
	if request.IgnoreIndex {
		//首页需要过滤相关分组
		ignoreIndexGroupIds, _ := s.sourceGroupDao.GetIgnoreIndexGroupIds(c, request.EnterpriseID)
		if len(ignoreIndexGroupIds) > 0 {
			sourceGroupIds = helper.ArrayDiff(sourceGroupIds, ignoreIndexGroupIds)
		}
	}
	// 分类下的sourceIDs
	if len(request.CategoryIds) > 0 {
		var whereCategoryIds = map[string]interface{}{
			"category_id": request.CategoryIds,
			"type":        request.Type,
		}
		sourceIDs, err = s.categoryRelationDao.GetManySourceIds(c, request.EnterpriseID, whereCategoryIds)
		if err != nil {
			return sourceIDs, sourceGroupIds, collectSourceIDs, err
		}
	}
	collectSourceIDs, err = sourceDao.NewCollectRelationDao().GetCollectSourceIds(c, request.EnterpriseID, request.MemberID, sourceIDs, request.Type)
	if err != nil {
		return sourceIDs, sourceGroupIds, collectSourceIDs, err
	}
	return sourceIDs, sourceGroupIds, collectSourceIDs, nil

}

// 获取查询参数
func (s GroupRelationService) getCollectSelectParams(c app.Ctx, request sourceRequest.RequestSelectList) ([]int64, []int64, []int64, error) {
	var (
		err              error
		sourceIDs        []int64
		sourceGroupIds   []int64
		collectSourceIDs []int64
	)
	// 获取用户所有收藏素材
	collectSourceIDs, err = sourceDao.NewCollectRelationDao().GetCollectSourceIds(c, request.EnterpriseID, request.MemberID, []int64{}, request.Type)
	if err != nil {
		return sourceIDs, sourceGroupIds, collectSourceIDs, err
	}
	if request.Type < sourceModel.TypeFriend {
		sourceIDs = collectSourceIDs
	}
	// 所属门店下 分组的sourceGroupIDs
	var whereGroupRelation = map[string]interface{}{
		"third_id": request.ThirdIds,
	}
	sourceGroupIds, err = s.groupRelationDao.GetManySourceGroupIds(c, request.EnterpriseID, whereGroupRelation)
	if err != nil {
		return sourceIDs, sourceGroupIds, collectSourceIDs, err
	}
	// 分类下的sourceIDs
	if len(request.CategoryIds) > 0 {
		var whereCategoryIds = map[string]interface{}{
			"category_id": request.CategoryIds,
			"type":        request.Type,
		}
		categorySourceIDs, err := s.categoryRelationDao.GetManySourceIds(c, request.EnterpriseID, whereCategoryIds)
		if err != nil {
			return categorySourceIDs, sourceGroupIds, collectSourceIDs, err
		}
		if len(categorySourceIDs) == 0 && !helper.InArray(request.Type, []int{
			sourceModel.TypeLink,
			sourceModel.TypeMiniProgram,
			sourceModel.TypeFile,
		}) {
			return categorySourceIDs, sourceGroupIds, collectSourceIDs, err
		}
		if request.Type < sourceModel.TypeFriend {
			sourceIDs, err = sourceDao.NewCollectRelationDao().GetCollectSourceIds(c, request.EnterpriseID, request.MemberID, categorySourceIDs, request.Type)
			if err != nil {
				return sourceIDs, sourceGroupIds, collectSourceIDs, err
			}
		} else {
			sourceIDs = categorySourceIDs
		}
	}
	return sourceIDs, sourceGroupIds, collectSourceIDs, nil

}

func (s GroupRelationService) CollectTextSelectList(c ctx.Context, request sourceRequest.RequestSelectList) (res []sourceModel.ResponseTextList, total int64, err error) {
	sourceIDs, sourceGroupIds, collectSourceIDs, err := s.getCollectSelectParams(c, request)
	if err != nil || len(collectSourceIDs) == 0 {
		return
	}
	res, total, err = s.groupRelationDao.GetTextSelectList(c, request, sourceIDs, sourceGroupIds)
	if err != nil {
		return
	}
	for k, v := range res {
		if helper.InArray(v.ID, collectSourceIDs) {
			res[k].IsCollect = sourceModel.SourceCollectYes
		}
	}
	return res, total, nil
}

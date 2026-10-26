package activityQRCodeSvc

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	activityQRCodeDao "micro-scrm/internal/app/dao/activity_qrcode"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	activityQRCodeResp "micro-scrm/internal/app/http/response/activity_qrcode"
	"micro-scrm/internal/app/models"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"sort"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type MainSvc struct {
	service.Base
	mainDao        *activityQRCodeDao.MainDao
	memberScopeDao *activityQRCodeDao.MemberScopeDao
	pointDao       *activityQRCodeDao.PointDao
	categoryDao    *activityQRCodeDao.CategoryDao
	TagGroupName   string
}

func NewMainSvc() *MainSvc {
	return &MainSvc{
		mainDao:        activityQRCodeDao.NewMainDao(),
		memberScopeDao: activityQRCodeDao.NewMemberScopeDao(),
		pointDao:       activityQRCodeDao.NewPointDao(),
		categoryDao:    activityQRCodeDao.NewCategoryDao(),
		TagGroupName:   "活动码_活动名称",
	}
}

// ActivitySave 活动保存
func (s *MainSvc) ActivitySave(c ctx.Context, req activityQRCodeReq.ActivitySaveReq) (int64, error) {
	nowUnix := time.Now().Unix()
	var activityID int64
	// 获取提交上来的投放渠道id集合
	var categoryIds []string
	if len(req.CategoryList) > 0 {
		for _, info := range req.CategoryList {
			categoryIds = append(categoryIds, helper.Int64ToString(info.ID))
		}
	}
	var tagId int64
	if req.AddNameTag > 0 {
		tagId = s.addActivityNameTag(c, req)
	}
	if req.ID > 0 {
		// 更新操作
		err := s.mainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
			"name":             req.Name,
			"alias_name":       req.AliasName,
			"banner_url":       req.BannerUrl,
			"time_type":        req.TimeType,
			"start_time":       req.StartTime,
			"end_time":         req.EndTime,
			"skip_verify":      req.SkipVerify,
			"add_name_tag":     req.AddNameTag,
			"add_category_tag": req.AddCategoryTag,
			"category_ids":     strings.Join(categoryIds, ","),
			"tag_ids":          helper.IntSliceToString(req.TagIds, ","),
			"name_tag_id":      tagId,
			"updated_at":       nowUnix,
		})
		activityID = req.ID
		if err != nil {
			return activityID, err
		}
	} else {
		// 创建操作
		activityID = app.UID()
		if err := s.mainDao.CreateData(c, req.EnterpriseId, activityQRCodeModel.ActivityQRCodeMainModel{
			BaseModel:      models.CreateBase(activityID),
			Name:           req.Name,
			AliasName:      req.AliasName,
			BannerUrl:      req.BannerUrl,
			TimeType:       req.TimeType,
			StartTime:      req.StartTime,
			EndTime:        req.EndTime,
			SkipVerify:     req.SkipVerify,
			AddNameTag:     req.AddNameTag,
			AddCategoryTag: req.AddCategoryTag,
			CategoryIds:    strings.Join(categoryIds, ","),
			TagIds:         helper.IntSliceToString(req.TagIds, ","),
			NameTagId:      tagId,
		}); err != nil {
			return activityID, err
		}
	}

	// 删除原有成员范围数据
	if err := s.memberScopeDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("activity_id = %d", req.ID), map[string]any{
		"deleted_at": nowUnix,
		"updated_at": nowUnix,
	}); err != nil {
		return activityID, err
	}

	// 批量插入新的成员范围数据
	var memberScopeList []activityQRCodeModel.ActivityQRCodeMemberScopeModel
	for _, memberID := range req.ScopeMemberIds {
		memberScopeList = append(memberScopeList, activityQRCodeModel.ActivityQRCodeMemberScopeModel{
			BaseModel:  models.CreateBase(app.UID()),
			ActivityID: activityID,
			MemberId:   memberID,
		})
	}
	for _, areaID := range req.ScopeAreaIds {
		memberScopeList = append(memberScopeList, activityQRCodeModel.ActivityQRCodeMemberScopeModel{
			BaseModel:  models.CreateBase(app.UID()),
			ActivityID: activityID,
			AreaId:     areaID,
		})
	}

	if len(memberScopeList) > 0 {
		if err := s.memberScopeDao.CreateInBatchesData(c, req.EnterpriseId, memberScopeList); err != nil {
			return activityID, err
		}
	}

	// 处理提交上来的渠道和点位
	var addStr []string
	for _, category := range req.CategoryList {
		//渠道id和点位名拼接唯一字符串
		addStr = append(addStr, fmt.Sprintf("%d-%s", category.ID, ""))
		for _, point := range category.Points {
			addStr = append(addStr, fmt.Sprintf("%d-%s", category.ID, point))
		}
	}

	var originalPoints = make([]activityQRCodeModel.ActivityQRCodePointModel, 0)
	var originalStr []string
	if req.ID > 0 {
		//查询原始点位
		originalPoints, _ = s.pointDao.GetMany(c, req.EnterpriseId, "activity_id = ?", activityID)
		for _, originalPoint := range originalPoints {
			//渠道id和点位名拼接唯一字符串
			originalStr = append(originalStr, fmt.Sprintf("%d-%s", originalPoint.CategoryID, originalPoint.Name))
		}
		//差集获取需要新增的和需要删除的
		tmpStr := addStr
		addStr = helper.GenericArrayDiff(addStr, originalStr)
		needDelStr := helper.GenericArrayDiff(originalStr, tmpStr)
		var delPointIds []int64
		if len(needDelStr) > 0 {
			for _, point := range originalPoints {
				for _, str := range needDelStr {
					tmpSlice := strings.Split(str, "-")
					if len(tmpSlice) == 2 && point.CategoryID == int64(helper.StringToInt(tmpSlice[0])) && point.Name == tmpSlice[1] {
						delPointIds = append(delPointIds, point.ID)
					}
				}
			}
		}
		if len(delPointIds) > 0 {
			//处理需要删除的
			_ = s.pointDao.UpdateData(c, req.EnterpriseId, map[string]any{"id": delPointIds}, map[string]any{
				"deleted_at": nowUnix,
				"updated_at": nowUnix,
			})
		}
		//过期二维码
		s.expiredContactWay(c, req.EnterpriseId, activityID)
		//删除人员范围缓存
		app.Redis().Del(c, fmt.Sprintf("activity-qrcode:member-scope:%d-%d", req.EnterpriseId, activityID))
	}
	if len(addStr) > 0 {
		var addPoints = make([]activityQRCodeModel.ActivityQRCodePointModel, 0)

		for _, str := range addStr {
			pointId := app.UID()
			tmpSlice := strings.Split(str, "-")
			tmpModel := activityQRCodeModel.ActivityQRCodePointModel{
				BaseModel:  models.CreateBase(pointId),
				ActivityID: activityID,
				Name:       tmpSlice[1],
				IsPoint:    helper.RuleThree(tmpSlice[1] == "", 0, 1),
				CategoryID: int64(helper.StringToInt(tmpSlice[0])),
			}
			//调用外部接口获取太阳码url
			qrCodeInfo, err := microwork.NewEcHub().GetQrCodeAndLink(c, microwork.SunQRCodeReq{
				EnterpriseID: req.EnterpriseId,
				Path:         "webPackage/pages/addCS/addCS",
				Scene:        helper.Int64ToString(pointId),
			})
			if err != nil {
				return activityID, err
			}
			tmpModel.CodeURL = qrCodeInfo.CodeUrl
			tmpModel.LinkURL = qrCodeInfo.Link
			addPoints = append(addPoints, tmpModel)
		}

		err := s.pointDao.CreateInBatchesData(c, req.EnterpriseId, addPoints)
		if err != nil {
			return activityID, err
		}
	}

	welcomeDao := activityQRCodeDao.NewWelcomeDao()
	if req.ID > 0 {
		_ = welcomeDao.UpdateData(c, req.EnterpriseId, map[string]any{
			"rel_type": activityQRCodeModel.WelcomeTypeActivity,
			"rel_id":   activityID,
		}, map[string]any{
			"updated_at": nowUnix,
			"deleted_at": nowUnix,
		})
	}
	if len(req.WelcomeList) > 0 {
		var welcomeModelList = make([]activityQRCodeModel.ActivityQRCodeWelcomeModel, 0)
		for _, info := range req.WelcomeList {
			contentStr, _ := json.Marshal(info.Content)
			welcomeModelList = append(welcomeModelList, activityQRCodeModel.ActivityQRCodeWelcomeModel{
				BaseModel: models.CreateBase(app.UID()),
				RelType:   activityQRCodeModel.WelcomeTypeActivity,
				RelID:     activityID,
				Content:   string(contentStr),
				Type:      info.Type,
			})
		}
		_ = welcomeDao.CreateInBatchesData(c, req.EnterpriseId, welcomeModelList)

	}
	return activityID, nil
}

func (s *MainSvc) addActivityNameTag(c ctx.Context, req activityQRCodeReq.ActivitySaveReq) int64 {
	tagGroupInfo, _ := tagsDao.NewWorkExternalTagGroup().GetOne(c, "enterprise_id=? AND group_name= ?", req.EnterpriseId, s.TagGroupName)
	if tagGroupInfo.ID > 0 {
		tagInfo, _ := tagsDao.NewWorkExternalTag().GetOne(c, "enterprise_id=? AND group_id=? AND name=?", req.EnterpriseId, tagGroupInfo.ID, req.Name)
		if tagInfo.ID > 0 {
			return tagInfo.ID
		}
	}
	addTagParams := work.AddCorpTagRequest{
		GroupName: s.TagGroupName,
	}
	addTagParams.TagList = append(addTagParams.TagList, work.AddCorpTag{
		Name: req.Name,
	})
	workImpl := work.NewQyWorkImpl().SetEnterprise(req.EnterpriseId)
	result, err := workImpl.AddCorpTag(c, addTagParams)
	if err == nil && len(result.TagGroup.TagList) > 0 {
		if tagGroupInfo.ID == 0 {
			tagGroupInfo.ID = app.UID()
			tagGroupModel := workModels.WorkExternalTagGroupModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, tagGroupInfo.ID),
				GroupID:        result.TagGroup.GroupID,
				GroupName:      result.TagGroup.GroupName,
				Type:           workModels.WorkExternalTagTypeQy,
				CreateTime:     int(result.TagGroup.CreateTime),
			}
			_ = tagsDao.NewWorkExternalTagGroup().CreateData(c, tagGroupModel)
		}
		for _, tag := range result.TagGroup.TagList {
			tagId := app.UID()
			tagModel := workModels.WorkExternalTagModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, tagId),
				GroupID:        tagGroupInfo.ID,
				TagID:          tag.ID,
				Name:           tag.Name,
				Type:           workModels.WorkExternalTagTypeQy,
				CreateTime:     int(tag.CreateTime),
			}
			_ = tagsDao.NewWorkExternalTag().CreateData(c, tagModel)
			return tagId
		}

	}
	return 0
}

// ActivityDetail  活动详情
func (s *MainSvc) ActivityDetail(c ctx.Context, req activityQRCodeReq.ActivityInfoReq) (activityQRCodeResp.ActivityInfoResp, error) {
	var resp = activityQRCodeResp.ActivityInfoResp{
		TagList:         []activityQRCodeResp.Info{},
		ScopeMemberList: []activityQRCodeResp.Info{},
		ScopeAreaList:   []activityQRCodeResp.Info{},
		CategoryList:    []activityQRCodeResp.CategoryInfo{},
		WelcomeList:     []activityQRCodeResp.ActivityWelcomeInfo{},
	}
	if req.State != "" {
		pointInfo, err := s.pointDao.GetOne(c, req.EnterpriseId, "id = ?", req.State)
		if err != nil {
			return activityQRCodeResp.ActivityInfoResp{}, err
		}
		req.ID = pointInfo.ActivityID
	}
	// 查询活动主信息
	activity, err := s.mainDao.GetOne(c, req.EnterpriseId, "id = ?", req.ID)
	if err != nil {
		return activityQRCodeResp.ActivityInfoResp{}, err
	}
	helper.CoverStruct(&resp, &activity)
	resp.ID = activity.ID
	// 查询成员范围
	memberScopes, err := s.memberScopeDao.GetMany(c, req.EnterpriseId, "activity_id = ?", req.ID)
	if err != nil {
		return activityQRCodeResp.ActivityInfoResp{}, err
	}

	// 查询点位信息
	points, err := s.pointDao.GetMany(c, req.EnterpriseId, "activity_id = ?", req.ID)
	if err != nil {
		return activityQRCodeResp.ActivityInfoResp{}, err
	}

	// 转换成员范围
	var areaIds []int
	var memberIds []int64
	for _, ms := range memberScopes {
		if ms.AreaId > 0 {
			areaIds = append(areaIds, ms.AreaId)
		}
		if ms.MemberId > 0 {
			memberIds = append(memberIds, ms.MemberId)
		}
	}
	memberList, err := membersDao.NewMembers().GetUnscopedMemberInfoByIds(c, req.EnterpriseId, memberIds)
	if err != nil {
		return activityQRCodeResp.ActivityInfoResp{}, err
	}
	for _, model := range memberList {
		resp.ScopeMemberList = append(resp.ScopeMemberList, activityQRCodeResp.Info{
			ID:   model.ID,
			Name: model.Name,
		})
	}
	areaList, err := dashboardDao.NewArea().GetMany(c, req.EnterpriseId, "id in ?", areaIds)
	if err != nil {
		return activityQRCodeResp.ActivityInfoResp{}, err
	}
	for _, model := range areaList {
		resp.ScopeAreaList = append(resp.ScopeAreaList, activityQRCodeResp.Info{
			ID:   model.ID,
			Name: model.Name,
		})
	}
	// 转换点位信息
	pointResp := make(map[int64][]activityQRCodeResp.PointInfo)
	var categorySort = make([]int64, 0)
	var categoryPointMap = make(map[int64]activityQRCodeModel.ActivityQRCodePointModel)
	for _, pointModel := range points {
		categorySort = append(categorySort, pointModel.CategoryID)
		if pointModel.Name == "" {
			categoryPointMap[pointModel.CategoryID] = pointModel
			continue
		}
		pointResp[pointModel.CategoryID] = append(pointResp[pointModel.CategoryID], activityQRCodeResp.PointInfo{
			ID:       pointModel.ID,
			Name:     pointModel.Name,
			CodeUrl:  pointModel.CodeURL,
			LinkUrl:  pointModel.LinkURL,
			QueryStr: fmt.Sprintf("scene=%d", pointModel.ID),
		})
	}
	categorySort = helper.ArrayUniq(categorySort)

	categoryList, _ := s.categoryDao.GetMany(c, req.EnterpriseId, "id in ?", categorySort)
	var categoryMap = make(map[int64]string)
	for _, model := range categoryList {
		categoryMap[model.ID] = model.Name
	}
	for _, categoryID := range categorySort {
		tmp := activityQRCodeResp.CategoryInfo{
			ID:        categoryID,
			LinkUrl:   categoryPointMap[categoryID].LinkURL,
			CodeUrl:   categoryPointMap[categoryID].CodeURL,
			QueryStr:  fmt.Sprintf("scene=%d", categoryPointMap[categoryID].ID),
			Name:      categoryMap[categoryID],
			PointList: []activityQRCodeResp.PointInfo{},
		}
		if pointInfoList, ok := pointResp[categoryID]; ok {
			sort.Slice(pointInfoList, func(i, j int) bool {
				return pointInfoList[i].ID < pointInfoList[j].ID
			})
			tmp.PointList = pointInfoList
		}
		resp.CategoryList = append(resp.CategoryList, tmp)

	}
	if req.State == "" && activity.TagIds != "" {
		var tagIds []int64
		for _, tagIdStr := range strings.Split(activity.TagIds, ",") {
			tagIds = append(tagIds, int64(helper.StringToInt(tagIdStr)))
		}
		//转换标签
		tagList, err := tagsDao.NewWorkExternalTag().GetMany(c, "enterprise_id = ? AND id in ?", req.EnterpriseId, tagIds)
		if err != nil {
			return activityQRCodeResp.ActivityInfoResp{}, err
		}
		for _, tagInfo := range tagList {
			resp.TagList = append(resp.TagList, activityQRCodeResp.Info{
				ID:   tagInfo.ID,
				Name: tagInfo.Name,
			})
		}
	}
	if req.State == "" {
		welcomeList, _ := activityQRCodeDao.NewWelcomeDao().GetMany(c, req.EnterpriseId, "rel_type = ? AND rel_id = ?", activityQRCodeModel.WelcomeTypeActivity, req.ID)
		for _, welcome := range welcomeList {
			var contentJson []map[string]interface{}
			_ = json.Unmarshal([]byte(welcome.Content), &contentJson)
			resp.WelcomeList = append(resp.WelcomeList, activityQRCodeResp.ActivityWelcomeInfo{
				Type:    welcome.Type,
				Content: contentJson,
			})
		}
	}

	return resp, nil
}

// ActivityList 活动列表
func (s *MainSvc) ActivityList(c ctx.Context, req activityQRCodeReq.ActivityListReq) ([]activityQRCodeResp.ActivityListResp, int64, error) {
	var respList = make([]activityQRCodeResp.ActivityListResp, 0)
	activityList, total, err := s.mainDao.GetActivityList(c, req)
	if err != nil || len(activityList) == 0 {
		return nil, 0, err
	}
	var tagIds []int64
	var categoryIds []int64
	var activityIds []int64
	for _, model := range activityList {
		if model.TagIds != "" {
			tagStrSlice := strings.Split(model.TagIds, ",")
			for _, tagStr := range tagStrSlice {
				tagIds = append(tagIds, int64(helper.StringToInt(tagStr)))
			}
		}
		if model.CategoryIds != "" {
			categoryStrSlice := strings.Split(model.CategoryIds, ",")
			for _, categoryStr := range categoryStrSlice {
				categoryIds = append(categoryIds, int64(helper.StringToInt(categoryStr)))
			}
		}
		activityIds = append(activityIds, model.ID)
	}
	welcomeList, _ := activityQRCodeDao.NewWelcomeDao().GetMany(c, req.EnterpriseId, "rel_type = ? AND rel_id in ?", activityQRCodeModel.WelcomeTypeActivity, activityIds)
	var welcomeMap = make(map[int64]bool)
	for _, model := range welcomeList {
		welcomeMap[model.RelID] = true
	}
	var tagMap = make(map[int64]string)
	if len(tagIds) > 0 {
		tagList, _ := tagsDao.NewWorkExternalTag().GetMany(c, "enterprise_id = ? AND id in ?", req.EnterpriseId, tagIds)
		for _, tagInfo := range tagList {
			tagMap[tagInfo.ID] = tagInfo.Name
		}
	}
	categoryList, _ := s.categoryDao.GetMany(c, req.EnterpriseId, "id in ?", categoryIds)
	var categoryMap = make(map[int64]string)
	for _, model := range categoryList {
		categoryMap[model.ID] = model.Name
	}
	for _, model := range activityList {
		var tmp = activityQRCodeResp.ActivityListResp{
			TagNames:      []string{},
			CategoryNames: []string{},
		}
		helper.CoverStruct(&tmp, &model)
		tmp.ID = model.ID
		tmp.CreatedAt = model.CreatedAt
		if model.TagIds != "" {
			tagStrSlice := strings.Split(model.TagIds, ",")
			for _, tagIdStr := range tagStrSlice {
				if tagName, ok := tagMap[int64(helper.StringToInt(tagIdStr))]; ok {
					tmp.TagNames = append(tmp.TagNames, tagName)
				}
			}
		}
		if model.CategoryIds != "" {
			categoryStrSlice := strings.Split(model.CategoryIds, ",")
			for _, categoryStr := range categoryStrSlice {
				if categoryName, ok := categoryMap[int64(helper.StringToInt(categoryStr))]; ok {
					tmp.CategoryNames = append(tmp.CategoryNames, categoryName)
				}
			}
		}
		if _, ok := welcomeMap[model.ID]; ok {
			tmp.HasWelcome = 1
		}
		respList = append(respList, tmp)
	}

	return respList, total, nil
}

func (s *MainSvc) ActivityCategoryList(c ctx.Context, req activityQRCodeReq.ActivityCategoryListReq) ([]activityQRCodeModel.ActivityQRCodeCategoryModel, error) {
	var categoryIds = make([]int64, 0)
	if req.AddWelcome {
		categoryIds, _ = activityQRCodeDao.NewWelcomeDao().GetAllCategoryIds(c, req.EnterpriseId, "rel_type = ?", activityQRCodeModel.WelcomeTypeCategory)
	}
	if len(categoryIds) > 0 {
		return s.categoryDao.GetMany(c, req.EnterpriseId, "id not in ? AND deleted_at=0", categoryIds)
	}
	return s.categoryDao.GetMany(c, req.EnterpriseId, "deleted_at=0")
}

func (s *MainSvc) ActivityDelete(c ctx.Context, req activityQRCodeReq.InfoReq) error {
	nowUnix := time.Now().Unix()
	// 查询活动主信息
	_ = s.mainDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
		"updated_at": nowUnix,
		"deleted_at": nowUnix,
	})
	_ = s.memberScopeDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("activity_id = %d", req.ID), map[string]interface{}{
		"updated_at": nowUnix,
		"deleted_at": nowUnix,
	})
	_ = s.pointDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("activity_id = %d", req.ID), map[string]interface{}{
		"updated_at": nowUnix,
		"deleted_at": nowUnix,
	})
	wayList, _ := activityQRCodeDao.NewContactWayDao().GetMany(c, req.EnterpriseId, "activity_id = ?", req.ID)
	go func() {
		if len(wayList) > 0 {
			for _, wayInfo := range wayList {
				_, _ = qyWork.NewQyWork(req.EnterpriseId).DeleteContactWay(c, qyWork.ContactDeleteRequest{
					ConfigId: wayInfo.ConfigId,
				})
			}
		}
	}()
	_ = activityQRCodeDao.NewContactWayDao().UpdateData(c, req.EnterpriseId, fmt.Sprintf("activity_id = %d", req.ID), map[string]interface{}{
		"updated_at": nowUnix,
		"deleted_at": nowUnix,
	})

	return nil
}

func (s *MainSvc) expiredContactWay(c ctx.Context, enterpriseId, activityId int64) {
	//如果活动被修改则所有二维码立刻过期
	_ = activityQRCodeDao.NewContactWayDao().UpdateData(c, enterpriseId, map[string]any{
		"activity_id": activityId,
	}, map[string]any{
		"updated_at": 0,
	})
}

func (s *MainSvc) WelcomeSave(c ctx.Context, req activityQRCodeReq.WelcomeInfoReq) (int64, error) {
	nowUnix := time.Now().Unix()
	var welcomeId int64
	contentStr, _ := json.Marshal(req.Content)
	if req.ID > 0 {
		welcomeDao := activityQRCodeDao.NewWelcomeDao()
		_ = welcomeDao.UpdateData(c, req.EnterpriseId, map[string]any{
			"id": req.ID,
		}, map[string]any{
			"rel_id":     req.CategoryID,
			"content":    string(contentStr),
			"updated_at": nowUnix,
		})
		welcomeId = req.ID
	} else {
		welcomeId = app.UID()
		_ = activityQRCodeDao.NewWelcomeDao().CreateData(c, req.EnterpriseId, activityQRCodeModel.ActivityQRCodeWelcomeModel{
			BaseModel: models.CreateBase(welcomeId),
			RelType:   activityQRCodeModel.WelcomeTypeCategory,
			RelID:     req.CategoryID,
			Content:   string(contentStr),
		})
	}
	return welcomeId, nil

}

func (s *MainSvc) WelcomeDetail(c ctx.Context, req activityQRCodeReq.InfoReq) (activityQRCodeResp.WelcomeInfo, error) {
	welcome, err := activityQRCodeDao.NewWelcomeDao().GetOne(c, req.EnterpriseId, "id = ?", req.ID)
	if err != nil {
		return activityQRCodeResp.WelcomeInfo{}, err
	}
	categoryInfo, _ := s.categoryDao.GetOne(c, req.EnterpriseId, "id = ?", welcome.RelID)

	var contentJson []map[string]interface{}
	_ = json.Unmarshal([]byte(welcome.Content), &contentJson)

	return activityQRCodeResp.WelcomeInfo{
		ID:           welcome.ID,
		CategoryId:   welcome.RelID,
		CategoryName: categoryInfo.Name,
		Content:      contentJson,
	}, nil
}

func (s *MainSvc) WelcomeList(c ctx.Context, req activityQRCodeReq.WelcomeListReq) ([]activityQRCodeResp.WelcomeListResp, int64, error) {
	var resList = make([]activityQRCodeResp.WelcomeListResp, 0)
	welcomeList, total, err := activityQRCodeDao.NewWelcomeDao().GetWelcomeList(c, req)
	if err != nil || len(welcomeList) == 0 {
		return resList, 0, err
	}
	categoryList, _ := s.categoryDao.GetMany(c, req.EnterpriseId, "deleted_at=0")
	var categoryMap = make(map[int64]string)
	for _, category := range categoryList {
		categoryMap[category.ID] = category.Name
	}
	for _, welcome := range welcomeList {
		var contentJson []map[string]interface{}
		_ = json.Unmarshal([]byte(welcome.Content), &contentJson)
		if categoryName, ok := categoryMap[welcome.RelID]; ok {
			resList = append(resList, activityQRCodeResp.WelcomeListResp{
				ID:           welcome.ID,
				CategoryName: categoryName,
				Content:      contentJson,
				CreatedAt:    welcome.CreatedAt,
			})
		}
	}
	return resList, total, err
}

func (s *MainSvc) WelcomeDelete(c ctx.Context, req activityQRCodeReq.InfoReq) error {
	nowUnix := time.Now().Unix()
	return activityQRCodeDao.NewWelcomeDao().UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.ID), map[string]interface{}{
		"updated_at": nowUnix,
		"deleted_at": nowUnix,
	})
}

func (s *MainSvc) Statistics(c ctx.Context, req activityQRCodeReq.StatisticsReq) (resp activityQRCodeResp.StatisticsResp, err error) {
	resp = activityQRCodeResp.StatisticsResp{
		TotalStatistics: activityQRCodeResp.TotalStatistics{
			AddRelationTotal: 0,
			AddCustomerTotal: 0,
			DelRelationTotal: 0,
			AddVipTotal:      0,
		},
		List:  make([]activityQRCodeResp.PointStatistics, 0),
		Total: 0,
		Page:  req.Page,
		Size:  req.Size,
	}
	//获取相关活动
	activityList, err := s.mainDao.GetStatisticsList(c, req)
	if err != nil || len(activityList) == 0 {
		return resp, nil
	}
	var activityIds []int64
	var activityNameMap = make(map[int64]string)
	for _, model := range activityList {
		activityIds = append(activityIds, model.ID)
		activityNameMap[model.ID] = model.Name
	}
	//获取相关点位
	statisticList, pointIds, err := s.pointDao.GetStatisticList(c, req, activityIds)
	if err != nil {
		return activityQRCodeResp.StatisticsResp{}, err
	}
	var categoryIds []int64
	var statePageSlice []string
	for _, pointInfo := range statisticList {
		categoryIds = append(categoryIds, pointInfo.CategoryID)
		statePageSlice = append(statePageSlice, helper.Int64ToString(pointInfo.ID))
	}
	//获取对应投放渠道
	categoryList, err := s.categoryDao.GetMany(c, req.EnterpriseId, "id in ?", helper.ArrayUniq(categoryIds))
	if err != nil {
		return activityQRCodeResp.StatisticsResp{}, err
	}
	var categoryNameMap = make(map[int64]string)
	for _, model := range categoryList {
		categoryNameMap[model.ID] = model.Name
	}
	resp.Total = len(pointIds)
	var stateSlice []string
	for _, id := range pointIds {
		stateSlice = append(stateSlice, helper.Int64ToString(id))
	}
	//获取按state分组统计
	totalStatistics, statisticsList, err := microScrmDao.NewUserRelation().GetActivityQrCodeStatistic(c, req, stateSlice, statePageSlice)
	if err != nil {
		return activityQRCodeResp.StatisticsResp{}, err
	}
	resp.TotalStatistics = totalStatistics
	var statisticMap = make(map[string]activityQRCodeResp.TotalStatistics)
	for _, statistic := range statisticsList {
		statisticMap[statistic.State] = statistic
	}
	for _, pointInfo := range statisticList {
		tmp := activityQRCodeResp.PointStatistics{
			ActivityName: "-",
			CategoryName: "-",
			PointName:    pointInfo.Name,
		}
		if activityName, ok1 := activityNameMap[pointInfo.ActivityID]; ok1 {
			tmp.ActivityName = activityName
		}
		if categoryName, ok2 := categoryNameMap[pointInfo.CategoryID]; ok2 {
			tmp.CategoryName = categoryName
		}
		if statistic, ok := statisticMap[helper.Int64ToString(pointInfo.ID)]; ok {
			tmp.TotalStatistics.AddRelationTotal = statistic.AddRelationTotal
			tmp.TotalStatistics.AddCustomerTotal = statistic.AddCustomerTotal
			tmp.TotalStatistics.DelRelationTotal = statistic.DelRelationTotal
			tmp.TotalStatistics.AddVipTotal = statistic.AddVipTotal
		}
		resp.List = append(resp.List, tmp)
	}

	return resp, nil
}

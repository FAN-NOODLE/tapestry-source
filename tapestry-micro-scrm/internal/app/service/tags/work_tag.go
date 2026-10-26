package tagsSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	tagModels "micro-scrm/internal/app/models/tag"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/app/service"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/sdk/work"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type WorkTags struct {
	service.Base
	workTagsDao     *tagsDao.WorkTags
	qyWorkImpl      *work.QyWorkImpl
	userRelationDao *microScrmDao.UserRelation
}

func NewWorkTags() *WorkTags {
	return &WorkTags{
		workTagsDao:     tagsDao.NewWorkTagsDao(),
		qyWorkImpl:      work.NewQyWorkImpl(),
		userRelationDao: microScrmDao.NewUserRelation(),
	}
}

func (t WorkTags) GetTagUserList(c app.Ctx, req tagsReq.InfoReq) ([]tagsReq.UserList, int64, error) {
	userList, total, err := microScrmDao.NewWxUserCorpExternalCk().GetTagList(c, req)
	if err != nil {
		return nil, 0, err
	}

	var users []tagsReq.UserList
	for _, item := range userList {
		var user tagsReq.UserList
		user.Mobile = helper.HandleMobile(item.Mobile)
		user.Name = item.Name
		user.UnionID = item.Unionid
		users = append(users, user)
	}
	return users, total, nil
}

// SyncThirdTagsToWorkTags 同步三方标签到企微标签
func (t WorkTags) SyncThirdTagsToWorkTags(c app.Ctx, enterpriseID int64, thirdTagsInfo []tagsReq.ThirdTagsInfo) error {
	groupTags := make(map[string][]tagsReq.ThirdTagsInfo, 0)
	tagsNameMap := make(map[string][]string, 0)
	thirdIdMap := make(map[string][]string, 0)
	for _, info := range thirdTagsInfo {
		groupTags[info.TagGroup] = append(groupTags[info.TagGroup], info)
		tagsNameMap[info.TagGroup] = append(tagsNameMap[info.TagGroup], info.TagName)
		thirdIdMap[info.TagGroup] = append(thirdIdMap[info.TagGroup], info.ID)
	}
	workImpl := t.qyWorkImpl.SetEnterprise(enterpriseID)
	for groupName, value := range groupTags {
		groupInfo, err := t.workTagsDao.WorkTagsGroupInfoByName(c, enterpriseID, groupName)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		tagsList, err := t.workTagsDao.GetWorkTagsByName(c, enterpriseID, groupInfo.ID, tagsNameMap[groupName])
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		updateIds := false
		tagsNameUpdateThirdIDs := make(map[string]string, 0)
		if len(tagsList) > 0 {
			updateIds = true
			for _, v := range tagsList {
				tagsNameUpdateThirdIDs[v.Name] = v.ThirdTagID
			}
		}
		tagsIdList, err := t.workTagsDao.GetWorkTagsByThirdIds(c, enterpriseID, groupInfo.ID, thirdIdMap[groupName])
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		updateTagNames := false
		tagsThirdIDsUpdateName := make(map[string]string, 0)
		if len(tagsIdList) > 0 {
			updateTagNames = true
			for _, v := range tagsList {
				tagsThirdIDsUpdateName[v.Name] = v.ThirdTagID
			}
		}
		var addCorpTagParam work.AddCorpTagRequest
		addCorpTagParam.GroupID = groupInfo.GroupID
		addCorpTagParam.GroupName = groupName
		tagsNameThirdIdMap := make(map[string]string, 0)
		tagsDescriptions := make(map[string]string, 0)
		for k, info := range value {
			tagsNameAdd := true
			tagThirdAdd := true
			if _, ok := tagsNameUpdateThirdIDs[info.TagName]; ok {
				tagsNameAdd = false
				tagsNameUpdateThirdIDs[info.TagName] = info.ID
			}
			if _, ok := tagsThirdIDsUpdateName[info.ID]; ok {
				tagThirdAdd = false
				tagsThirdIDsUpdateName[info.ID] = info.TagName
			}
			var addCorpTag work.AddCorpTag
			addCorpTag.Name = info.TagName
			addCorpTag.Order = k + 1
			if tagThirdAdd && tagsNameAdd {
				addCorpTagParam.TagList = append(addCorpTagParam.TagList, addCorpTag)
				tagsNameThirdIdMap[info.TagName] = info.ID
				tagsDescriptions[info.TagName] = info.TagDescription
			}
		}
		if len(addCorpTagParam.TagList) > 0 {
			Result, err := workImpl.AddCorpTag(c, addCorpTagParam)
			if err != nil {
				return err
			}
			groupID := groupInfo.ID
			if groupID == 0 {
				group := Result.TagGroup
				groupModel := workModels.WorkExternalTagGroupModel{
					GroupName:      group.GroupName,
					GroupID:        group.GroupID,
					CreateTime:     int(group.CreateTime),
					IsThirdTag:     1,
					MicroBaseModel: models.CreateMicroBase(enterpriseID, app.UIDContext(c)),
				}
				groupID, err = t.workTagsDao.InsertWorkTagGroup(c, groupModel)
				if err != nil {
					return err
				}
			}
			tags := Result.TagGroup.TagList
			var InsertTagsInfo []workModels.WorkExternalTagModel
			for _, tag := range tags {
				tagInfo := workModels.WorkExternalTagModel{
					TagID:               tag.ID,
					GroupID:             groupID,
					Name:                tag.Name,
					Order:               tag.Order,
					CreateTime:          int(tag.CreateTime),
					ThirdTagID:          tagsNameThirdIdMap[tag.Name],
					ThirdTagDescription: tagsDescriptions[tag.Name],
					Type:                workModels.WorkExternalTagTypeQy,
					MicroBaseModel:      models.CreateMicroBase(enterpriseID, app.UIDContext(c)),
				}
				InsertTagsInfo = append(InsertTagsInfo, tagInfo)
			}

			if err = t.workTagsDao.InsertWorkTags(c, InsertTagsInfo); err != nil {
				return err
			}
		}
		if updateIds {
			for k, v := range tagsNameUpdateThirdIDs {
				if err = t.workTagsDao.UpdateWorkTagsThirdByName(c, enterpriseID, groupInfo.ID, k, v); err != nil {
					return err
				}
			}
		}
		if updateTagNames {
			for k, v := range tagsThirdIDsUpdateName {
				if err = t.workTagsDao.UpdateWorkTagsNameByThird(c, enterpriseID, groupInfo.ID, k, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// SyncExternalWorkTags 第三方的用户和标签同步 1为新增标签 2为删除标签
func (t WorkTags) SyncExternalWorkTags(c app.Ctx, enterpriseID int64, data []tagsReq.ExternalUserThirdTag, status int) error {
	thirdIdTOTagIDMap := make(map[string][]tagsReq.ExternalUserThirdTag, 0)
	for _, item := range data {
		thirdIdTOTagIDMap[item.ThirdTagID] = append(thirdIdTOTagIDMap[item.ThirdTagID], item)
	}
	workImpl := t.qyWorkImpl.SetEnterprise(enterpriseID)
	for k, value := range thirdIdTOTagIDMap {
		uniqueUnionIdMap := make(map[string]string, 0)
		res, err := t.workTagsDao.GetTagInfoByThirdID(c, enterpriseID, k)
		if err != nil {
			continue
		}
		for _, v := range value {
			if _, ok := uniqueUnionIdMap[v.UnionId]; ok {
				continue
			}
			uniqueUnionIdMap[v.UnionId] = k

			relationInfo, err := t.userRelationDao.GetUserRelationByUnionId(c, enterpriseID, v.UnionId)
			if err != nil || relationInfo.ID == 0 {
				continue
			}

			externalUserTagInfo, err := t.workTagsDao.GetExternalUserTagInfo(c, enterpriseID, res.ID, relationInfo.ID)
			if err != nil || externalUserTagInfo.ID > 0 {
				continue
			}

			externalUser, err := t.workTagsDao.GetExternalIdByUnionId(c, enterpriseID, v.UnionId)
			if err != nil || externalUser.ID == 0 {
				continue
			}

			var userTag work.UserTagHandlerRequest
			userTag.Userid = relationInfo.UserID
			userTag.ExternalUserid = externalUser.Openid
			if status == 1 {
				userTag.AddTag = append(userTag.AddTag, res.TagID)
			}
			if status == 2 {
				userTag.RemoveTag = append(userTag.RemoveTag, res.TagID)
			}
			_, err = workImpl.UserTagHandler(c, userTag)
			if err != nil {
				continue
			}

			if status == 1 {
				externalUserTag := microScrmModel.WorkExternalUserTagModel{
					ExternalTagID:  res.ID,
					RelationID:     relationInfo.ID,
					MicroBaseModel: models.CreateMicroBase(enterpriseID, app.UIDContext(c)),
				}
				if err = t.workTagsDao.AddExternalUserTag(c, externalUserTag); err != nil {
					continue
				}
			}
			if status == 2 {
				if err = t.workTagsDao.DelExternalUserTag(c, enterpriseID, res.ID, relationInfo.ID); err != nil {
					continue
				}
			}
		}
	}
	return nil
}

func (t WorkTags) UserTagsAndTagGroup(c app.Ctx, req tagsReq.UserTagsReq) (info []tagsReq.UserTagsGroupList, err error) {
	var userTagsList []tagsReq.UserTagsRes
	if len(req.MemberIds) == 0 {
		req.MemberIds = append(req.MemberIds, req.MemberId)
	}
	userTagsList, _, err = t.workTagsDao.UserTagsAndTagGroup(c, req.EnterpriseID, req.WorkExternalUserId, req.MemberIds)

	if err != nil {
		return nil, err
	}
	// configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
	// 	"enterprise_id": req.EnterpriseID,
	// 	"key":           "only_third_tag",
	// })
	configData, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, req.EnterpriseID, "only_third_tag")
	var onlyThirdTag int
	if err == nil && configData.Value == 1 {
		onlyThirdTag = 1
	}
	var userTagsGroupList []tagsReq.UserTagsGroupList
	userGroupInfoMap := make(map[int64]tagsReq.GroupInfoRes)
	userTagsMap := make(map[int64][]tagsReq.TagsInfo, 0)

	for _, v := range userTagsList {
		if onlyThirdTag == 1 && v.IsThirdTag != 1 {
			continue
		}
		if _, ok := userGroupInfoMap[v.GroupId]; !ok {
			userGroupInfoMap[v.GroupId] = tagsReq.GroupInfoRes{
				// RelationId:    v.RelationId,
				ExternalTagId: v.ExternalTagId,
				Type:          v.Type,
				GroupBaseInfo: tagsReq.GroupBaseInfo{
					GroupId:   v.GroupId,
					GroupName: v.GroupName,
				},
			}
		}
		userTagsMap[v.GroupId] = append(userTagsMap[v.GroupId], tagsReq.TagsInfo{
			TagId:               v.TagId,
			Name:                v.Name,
			ThirdTagDescription: v.ThirdTagDescription,
		})
	}

	for k := range userGroupInfoMap {
		temp := tagsReq.UserTagsGroupList{
			GroupInfoRes: userGroupInfoMap[k],
			TagsList:     userTagsMap[k],
		}
		userTagsGroupList = append(userTagsGroupList, temp)
	}

	// 标签排序
	sort.Slice(userTagsGroupList, func(i, j int) bool {
		return userTagsGroupList[i].GroupId > userTagsGroupList[j].GroupId
	})
	for _, v := range userTagsGroupList {
		sort.Slice(v.TagsList, func(i, j int) bool {
			return v.TagsList[i].TagId > v.TagsList[j].TagId
		})
	}

	return userTagsGroupList, nil
}

func (t WorkTags) TagsAuth(c app.Ctx, req tagsReq.InfoBaseReq) (int64, error) {
	return t.workTagsDao.TagsAuth(c, req.EnterpriseID)
}

func (t WorkTags) GetTagList(c app.Ctx, enterpriseID int64) (data []workModels.WorkExternalTagModel, err error) {
	return t.workTagsDao.GetExternalTagList(c, enterpriseID)
}

type UserTagHandler struct {
	Status     int
	RelationId int64
	AddTagIds  []int64
	DelTagIds  []int64
	work.UserTagHandlerRequest
}

// BatchMarking 批量打标
// 禁用并发调用,必须同步执行
func (t WorkTags) BatchMarking(c app.Ctx, title string, enterpriseId int64, externalUserTags []tagsReq.ExternalUserTag, status int) error {
	if len(externalUserTags) == 0 {
		return errors.New("打标关系不能为空")
	}

	// 防止并发执行,或者已有执行中的任务了，还继续执行
	isLock := t.Redis().SetNX(c, "qy_batch_marking_lock", time.Now().Unix(), 0).Val()
	if !isLock {
		return errors.New("禁止同时执行")
	}

	defer func() {
		t.Redis().Del(c, "qy_batch_marking_lock")
	}()

	// 取出所有 unionIdUserId 和 unionIds 的集合做一次性查询
	unionIdUserIdSlice := make([][]any, 0)
	unionIds := make([]string, 0)
	for _, externalUserTag := range externalUserTags {
		unionIdUserIdSlice = append(unionIdUserIdSlice, []any{externalUserTag.UnionId, externalUserTag.UserId})
		unionIds = append(unionIds, externalUserTag.UnionId)
	}

	unionIdUserIdMap := make(map[string]int64)
	// 取出所有 外部联系人和员工的关系
	result, err := t.userRelationDao.GetAllByUnionIdUserIdMap(c, enterpriseId, unionIdUserIdSlice)
	if err != nil {
		return err
	}

	// 将unionId和UserId拼接作为key
	for _, val := range result {
		unionIdUserIdMap[val.Unionid+val.UserID] = val.ID
	}

	// 获取 unionIds 对应的openId
	wxUserCorpExternal, err := microScrmDao.NewWxUserCorpExternal().GetMany(c, enterpriseId, map[string]any{"unionid": unionIds})
	if err != nil {
		return err
	}

	// 匹配 unionId 对应的 Openid
	unionIdMap := make(map[string]string)
	for _, external := range wxUserCorpExternal {
		unionIdMap[external.Unionid] = external.Openid
	}

	// 判断限速-不满足则阻塞等待
	requestQyApiLimit(c, "qy_api_hour_count", len(externalUserTags))

	NewTaskNotice().MarkSendNotice(c, true, TaskMessageData{
		Title:         fmt.Sprintf("%s-批量打标-开始", title),
		EnterpriseId:  enterpriseId,
		NeedMarkCount: len(externalUserTags),
	})

	// 打标成功的数量
	var successTagNum int64

	useMemory, useCpu, useTime := helper.FuncExecTime(func() {
		successTagNum = t.handleMark(c, enterpriseId, externalUserTags, unionIdUserIdMap, unionIdMap, status)
	})

	NewTaskNotice().MarkSendNotice(c, false, TaskMessageData{
		Title:            fmt.Sprintf("%s-批量打标-结束", title),
		EnterpriseId:     enterpriseId,
		NeedMarkCount:    len(externalUserTags),
		SuccessMarkCount: successTagNum,
		UseMemory:        useMemory,
		UseCpu:           useCpu,
		UseTime:          useTime,
	})

	return nil
}

func (t WorkTags) handleMark(
	c app.Ctx,
	enterpriseId int64,
	externalUserTags []tagsReq.ExternalUserTag,
	unionIdUserIdMap map[string]int64,
	unionIdMap map[string]string,
	status int,
) int64 {

	// 打标成功的数量
	var successTagNum int64

	// 在方法内部，通过方法参数中传递过来的企业 ID，初始化相应的 QyWorkImpl 和外部标签 DAO 实例
	workImpl := t.qyWorkImpl.SetEnterprise(enterpriseId)

	limitChan := make(chan interface{}, 2100)

	threading.GoSafe(func() error {
		resultChan := make(chan UserTagHandler, len(externalUserTags))
		wg := &sync.WaitGroup{}
		pool, _ := ants.NewPool(30)

		for _, externalUserTag := range externalUserTags {
			wg.Add(1)
			_ = pool.Submit(func(externalUserTagItem tagsReq.ExternalUserTag) func() {
				return func() {
					defer wg.Done()
					// 判断 unionId和UserId的key 是否存在关系，不存在则没有好友关系
					relationId, ok := unionIdUserIdMap[externalUserTagItem.UnionId+externalUserTagItem.UserId]
					if !ok {
						app.Logger().Error(c).Msgf("根据企业ID和UnionID获取用户关联信息异常,UnionID = %s UserId = %s", externalUserTagItem.UnionId, externalUserTagItem.UserId)
						return
					}

					// 从 unionIdMap 中取出对应的 unionId 的 openId
					openId, ok := unionIdMap[externalUserTagItem.UnionId]
					if !ok {
						app.Logger().Error(c).Msgf("获取当前UnionId对应的ExternalUserid异常,UnionID = %s", externalUserTagItem.UnionId)
						return
					}

					userTagRequest := UserTagHandler{
						Status:     status,
						RelationId: relationId,
						UserTagHandlerRequest: work.UserTagHandlerRequest{
							Userid:         externalUserTagItem.UserId,
							ExternalUserid: openId,
						},
					}

					tagIds, _ := t.workTagsDao.GetExternalTagIds(c, enterpriseId, []string{externalUserTagItem.TagId})

					if status == 1 {
						userTagRequest.AddTag = append(userTagRequest.AddTag, externalUserTagItem.TagId)
						userTagRequest.AddTagIds = tagIds
					}

					if status == 2 {
						userTagRequest.RemoveTag = append(userTagRequest.RemoveTag, externalUserTagItem.TagId)
						userTagRequest.DelTagIds = tagIds
					}

					resultChan <- userTagRequest
				}
			}(externalUserTag))
		}

		wg.Wait()
		close(resultChan)

		for userTagRequest := range resultChan {
			limitChan <- userTagRequest
		}

		close(limitChan)
		return nil
	})

	wg := &sync.WaitGroup{}
	pool, _ := ants.NewPool(200)

	// 限频的维度,1小时频率为120000次,即1分钟频率为2100次,即1秒频率为35次
	chanResult := helper.ChannelRateLimit(limitChan, time.Second, 35)

	for result := range chanResult {
		resultTmp := result.(UserTagHandler)
		wg.Add(1)
		_ = pool.Submit(func(res UserTagHandler) func() {
			return func() {
				defer wg.Done()
				// 构造 UserTagHandlerRequest 请求对象
				userTagRequest := work.UserTagHandlerRequest{
					Userid:         res.Userid,
					ExternalUserid: res.ExternalUserid,
				}
				userTagRequest.AddTag = append(userTagRequest.AddTag, res.AddTag...)
				userTagRequest.RemoveTag = append(userTagRequest.RemoveTag, res.RemoveTag...)

				// 调用企业微信 API 实现标签成员增加和删除操作
				if _, err := workImpl.UserTagHandler(c, userTagRequest); err != nil {
					return
				}

				// 打标成功计次+1
				atomic.AddInt64(&successTagNum, 1)

				// 根据 status 进行相应的数据库操作
				if res.Status == 1 && len(res.AddTagIds) > 0 {
					var externalUserTagData []microScrmModel.WorkExternalUserTagModel
					for _, tagId := range res.AddTagIds {
						externalUserTagData = append(externalUserTagData, microScrmModel.WorkExternalUserTagModel{
							ExternalTagID:  tagId,
							RelationID:     res.RelationId,
							MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
						})
					}
					err := t.workTagsDao.BatchCreate(c, externalUserTagData)
					if err != nil {
						return
					}
				}

				if res.Status == 2 && len(res.DelTagIds) > 0 {
					err := t.workTagsDao.DelExternalUserTags(c, enterpriseId, res.RelationId, res.DelTagIds)
					if err != nil {
						return
					}
				}

				return
			}
		}(resultTmp))
	}

	wg.Wait()

	return successTagNum
}

func requestQyApiLimit(c app.Ctx, hourKey string, runCount int) {
	// 获取 redis 客户端
	redisClient := app.Redis()
	for {
		// 企业微信调用的频次限制
		var qyApiHourLimit = 120000

		// 判断限频的key值是否存在，若不存在则创建一个
		hourKeyExists := redisClient.Exists(c, hourKey).Val()
		if hourKeyExists == 0 {
			redisClient.SetNX(c, hourKey, qyApiHourLimit-runCount, time.Hour)
			break
		}

		// 获取请求企微接口剩余频次
		// 判断是否是够本次执行数量
		hourCount, _ := redisClient.Get(c, hourKey).Int()
		if hourCount > runCount {
			redisClient.IncrBy(c, hourKey, int64(-runCount))
			break
		}

		// 若当前剩余限额大于或等于请求量，则限制频次在小时维度上
		hourTime := redisClient.TTL(c, hourKey).Val()
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("企微频次限制,剩余频次%d,执行数量%d,等待时间为%s", hourCount, runCount, hourTime))
		<-time.After(hourTime)
		notice.NewWechat().SendTextMessage(c, "企微频次恢复,开始执行任务")
	}
}

// CreateTaskTag 创建任务标签
func (t WorkTags) CreateTaskTag(c app.Ctx, enterpriseId int64, groupName string, tagName string, endTime, tagType int) (string, int64, error) {
	var workTagId string
	var id int64
	groupInfo, err := t.workTagsDao.WorkTagsGroupInfoByName(c, enterpriseId, groupName)
	if err != nil {
		return workTagId, id, err
	}
	if groupInfo.ID == 0 {
		return workTagId, id, errors.New(groupName + "组不存在")
	}
	tagsList, err := t.workTagsDao.GetWorkTagsByName(c, enterpriseId, groupInfo.ID, []string{tagName})
	if err != nil {
		return workTagId, id, err
	}
	if len(tagsList) > 0 {
		return workTagId, id, errors.New(tagName + "标签已存在")
	}
	workImpl := t.qyWorkImpl.SetEnterprise(enterpriseId)
	var addCorpTagParam work.AddCorpTagRequest
	addCorpTagParam.GroupID = groupInfo.GroupID
	addCorpTagParam.GroupName = groupName
	var addCorpTag work.AddCorpTag
	addCorpTag.Name = tagName
	addCorpTagParam.TagList = append(addCorpTagParam.TagList, addCorpTag)
	result, err := workImpl.AddCorpTag(c, addCorpTagParam)
	if err != nil {
		return workTagId, id, errors.New(tagName + "标签已存在")
	}
	tags := result.TagGroup.TagList
	var insertTagsInfo []workModels.WorkExternalTagModel
	for _, tag := range tags {
		workTagId = tag.ID
		id = app.UID()
		tagInfo := workModels.WorkExternalTagModel{
			TagID:          tag.ID,
			GroupID:        groupInfo.ID,
			Name:           tag.Name,
			Order:          tag.Order,
			CreateTime:     int(tag.CreateTime),
			ExpiredTime:    endTime,
			Type:           int32(tagType),
			MicroBaseModel: models.CreateMicroBase(enterpriseId, id),
		}
		insertTagsInfo = append(insertTagsInfo, tagInfo)
	}
	if err = t.workTagsDao.InsertWorkTags(c, insertTagsInfo); err != nil {
		return workTagId, id, err
	}
	return workTagId, id, err
}

// GetChoiceTagList 获取小程序端筛选条件的标签列表
func (t WorkTags) GetChoiceTagList(c app.Ctx, enterpriseID int64) (data []workModels.WorkExternalTagModel, err error) {
	qyTagRedisKey := fmt.Sprintf("tags:qy-tag:%d", enterpriseID)
	qyTagJson, _ := t.Redis().Get(c, qyTagRedisKey).Bytes()
	if qyTagJson != nil {
		//有缓存直接拿缓存
		_ = json.Unmarshal(qyTagJson, &data)
	}
	if data == nil || len(data) == 0 {
		data, err = t.workTagsDao.GetChoiceExternalTagList(c, enterpriseID)
		if len(data) > 0 {
			// 写缓存
			qyTagByte, _ := json.Marshal(data)
			t.Redis().Set(c, qyTagRedisKey, string(qyTagByte), time.Minute)
		}
	}
	return
}

func (t WorkTags) GetChatTagList(c ctx.Context, enterpriseID int64) (data []workModels.WorkExternalTagModel, err error) {
	chatTagRedisKey := fmt.Sprintf("tags:chat-tag:%d", enterpriseID)
	chatTagJson, _ := t.Redis().Get(c, chatTagRedisKey).Bytes()
	if chatTagJson != nil {
		//有缓存直接拿缓存
		_ = json.Unmarshal(chatTagJson, &data)
	}
	if data == nil || len(data) == 0 {
		data, err = t.workTagsDao.GetMany(c, map[string]interface{}{
			"enterprise_id": enterpriseID,
			"third_tag_id":  workModels.ChatTagThirdType,
			"deleted_at":    0,
		})
		if len(data) > 0 {
			// 写缓存
			chatTagByte, _ := json.Marshal(data)
			t.Redis().Set(c, chatTagRedisKey, string(chatTagByte), time.Minute)
		}
	}
	return
}

func (t WorkTags) GetUserChatTags(c app.Ctx, enterpriseID int64, memberId int64, externalUserId int64) (tagModels.UserTagsGroupList, error) {
	var result tagModels.UserTagsGroupList
	many, err := t.workTagsDao.GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"third_tag_id":  workModels.ChatTagThirdType,
		"deleted_at":    0,
	})
	if err != nil {
		return result, err
	}

	tagIds := make([]int64, 0)
	tagRelation := make(map[int64]string, 0)
	var groupId int64
	for _, model := range many {
		groupId = model.GroupID
		tagIds = append(tagIds, model.ID)
		tagRelation[model.ID] = model.Name
		result.ImpressionTagsList = append(result.ImpressionTagsList, tagModels.ImpressionTagsInfo{
			TagId:  model.ID,
			Name:   model.Name,
			IsHave: 0,
		})
	}
	groupInfo, err := tagsDao.NewWorkExternalTagGroup().GetOne(c, map[string]interface{}{
		"id":         groupId,
		"deleted_at": 0,
	})
	if err != nil {
		return result, err
	}
	result.GroupId = groupId
	result.GroupName = groupInfo.GroupName
	one, err := microScrmDao.NewUserRelation().GetOne(c, map[string]interface{}{
		"enterprise_id":         enterpriseID,
		"member_id":             memberId,
		"work_external_user_id": externalUserId,
	})
	if err != nil || one.ID == 0 {
		return result, err
	}
	getMany, err := tagsDao.NewWorkExternalUserTag().GetMany(c, map[string]interface{}{
		"relation_id":     one.ID,
		"external_tag_id": tagIds,
		"deleted_at":      0,
	})
	if err != nil {
		return result, err
	}
	for _, value := range getMany {
		for i, v := range result.ImpressionTagsList {
			if value.ExternalTagID == v.TagId {
				result.ImpressionTagsList[i].IsHave = 1
			}
		}
	}
	return result, err
}

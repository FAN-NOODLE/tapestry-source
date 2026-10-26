package http

import (
	"context"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagsSvc "micro-scrm/internal/app/service/tags"
	userGroupSvc "micro-scrm/internal/app/service/user_group"
	"micro-scrm/internal/modules/customer"
	helper2 "micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	response "micro-scrm/internal/pkg/resoponse"
	"sort"
	"strconv"

	"micro-scrm/internal/modules/user"

	"micro-scrm/internal/events"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/pkg/common"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"

	// "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/errs"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	EventData "micro-scrm/internal/events/entity"
)

var tagService *tag.Service
var tagGroupService *tag.GroupService
var marketService *marketing.Service
var enterpriseService *enterprise.Service
var userService *user.Service

// work标签服务系列
var (
	workExternalTagService      *tag.WorkExternalTagService
	workExternalTagGroupService *tag.WorkExternalTagGroupService
	// 客户分组服务
	userGroupService *tag.UserGroupService
)

// 标签分组列表
func tagList(c *gin.Context) {
	var request entity.SearchTagRequest
	_ = c.BindJSON(&request)
	//
	// if request.Status != 3 {
	//	ginplus.ResError(c,errors.New("Parameter status error,please confirm!"))
	//	return
	// }
	var tagGroupList *entity.TagGroupList
	var err error

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if request.Keyword == "" {
		request.SearchType = 0
	}
	// if request.Keyword == "" {
	//	tagGroupList, err = tagGroupService.TagGroupList(wrapContext(c), enterpriseID.(int64), request.Type)
	// } else {
	tagGroupList, err = workExternalTagService.SearchTag(wrapContext(c), request, enterpriseID.(int64))
	// }

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, tagGroupList)
}

// createGroup 创建手动标签组
func createGroup(c *gin.Context) {
	var request entity.CreateGroupRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	enterpriseId := common.EnterpriseID(cc)
	if enterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	var workTagIds []string
	if request.GroupId > 0 {
		params := entity.CreateManualTagParam{
			EnterpriseID: enterpriseId,
			TagGroupID:   request.GroupId,
			TagName:      request.TagNames[0],
			ExpiredTime:  request.ExpiredTime,
		}
		workTagIds, err = workExternalTagService.CreateTag(cc, params)
	} else {
		params := entity.CreateManualTagGroupParams{
			EnterpriseID:       enterpriseId,
			CreateGroupRequest: request,
		}
		workTagIds, err = workExternalTagGroupService.AddTagGroup(cc, params)
	}

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if len(workTagIds) > 0 && request.FileUrl != "" {
		threading.GoSafe(func() error {
			_ = syncMarkingTag(cc, workTagIds, enterpriseId, request.FileUrl)
			return nil
		})
	}

	ginplus.ResOK(c)
}

// 创建手动标签接口
func createManualTag(c *gin.Context) {
	var request entity.CreateTagRequest
	helper.Must(c.BindJSON(&request))
	if len(request.TagName) > 30 {
		ginplus.ResError(c, errors.New("超出范围"))
		return
	}

	// 获取企业ID
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	// 为指定的标签组创建标签
	_, err := workExternalTagService.CreateTag(wrapContext(c), entity.CreateManualTagParam{
		EnterpriseID: enterpriseID.(int64),
		TagGroupID:   request.TagGroupID,
		TagName:      request.TagName,
	})

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, nil)
}

// userCreateGroup 标签分组创建
func userCreateGroup(c *gin.Context) {
	var CreateUserGroupBodyParams entity.CreateUserGroupBody
	err := c.ShouldBindJSON(&CreateUserGroupBodyParams)
	if err != nil {
		ginplus.ResError(c, err)
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	groupID, createErr := userGroupService.AddUserGroup(wrapContext(c), enterpriseID.(int64), CreateUserGroupBodyParams)
	if createErr != nil {
		ginplus.ResError(c, createErr)
		return
	}

	if CreateUserGroupBodyParams.UserGroupType == 1 {
		go func() {
			err := userGroupSvc.NewGroupSvc().TagsHandle(wrapContext(c), enterpriseID.(int64), groupID)
			if err != nil {
				log.Info(wrapContext(c)).Msgf("创建分组同步标签类型的客户分组失败::,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID)
			}
		}()
	}
	if CreateUserGroupBodyParams.UserGroupType == 3 {
		go func() {
			err := userGroupSvc.NewGroupSvc().DepartmentHandle(wrapContext(c), enterpriseID.(int64), groupID)
			if err != nil {
				log.Info(wrapContext(c)).Msgf("创建分组同步门店分组类型的客户分组失败::,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID)
			}
		}()
	}
	if CreateUserGroupBodyParams.UserGroupType == 2 {
		go func() {
			err := userGroupService.SyncGroupImportData(wrapContext(c), enterpriseID.(int64), CreateUserGroupBodyParams.FileID, groupID)
			if err != nil {
				return
			} else {
				err = userGroupSvc.NewGroupSvc().ImportDataHandle(wrapContext(c), enterpriseID.(int64), CreateUserGroupBodyParams.FileID, groupID)

			}

		}()
	}
	ginplus.ResSuccess(c, entity.CreateUserGroupResponse{
		GroupID: groupID,
	})
}

func userGroupExport(c *gin.Context) {
	var ExportReqParams entity.ExportReq
	err := c.ShouldBindJSON(&ExportReqParams)
	if err != nil {
		ginplus.ResError(c, err)
	}
	var result entity.ExportList
	result.Total = ExportReqParams.Total
	if ExportReqParams.Total == 0 {
		Total, totalErr := userGroupService.GetUserGroupLogCount(ExportReqParams.EnterpriseID, ExportReqParams.UserGroupID)
		if totalErr != nil {
			ginplus.ResError(c, totalErr)
		}
		result.Total = Total
	}
	list, listErr := userGroupService.GetUserGroupLogList(ExportReqParams.EnterpriseID, ExportReqParams.UserGroupID, ExportReqParams.Page, ExportReqParams.PageSize)
	if listErr != nil {
		ginplus.ResError(c, listErr)
	}
	result.ExportData = list
	ginplus.ResSuccess(c, &result)
	return

}

// userGroupFileDetection
func userGroupFileDetection(c *gin.Context) {
	var fileDetection entity.FileDetection
	if err := c.ShouldBindJSON(&fileDetection); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if fileDetection.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	detection, err := userGroupService.GetFileNameDetection(fileDetection.EnterpriseID, fileDetection.FileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var result entity.CountInfo
	result.Count = detection
	ginplus.ResSuccess(c, &result)
	return
}

// userGroupFileSaveInfo
func userGroupFileSaveInfo(c *gin.Context) {
	var request entity.DownLoadFileInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空"))
		return
	}
	var idInfo entity.IDInfo
	ID, err := userGroupService.SaveUploadInfo(request, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	idInfo.ID = ID
	ginplus.ResSuccess(c, &idInfo)
	return
}

// userGroupFileUpdate
func userGroupFileUpdate(c *gin.Context) {
	var request entity.UpdateIdInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空! "))
		return
	}
	if err := userGroupService.GetUpdateInfo(request.EnterpriseID, request.ID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

// userGroupEdit 客户分组编辑
func userGroupEdit(c *gin.Context) {
	var CreateUserGroupBodyParams entity.CreateUserGroupBody
	err := c.BindJSON(&CreateUserGroupBodyParams)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if CreateUserGroupBodyParams.ID == 0 {
		ginplus.ResError(c, errs.UserGroupIDIsEmpty)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	ctx1 := wrapContext(c)
	// 删除客户分组详情，删除客户分组对应的客户分组groups
	events.ScrmEvent.DeleteCustomerGroup.Dispatch(EventData.Data{
		Ctx:             ctx1,
		EnterpriseID:    enterpriseID.(int64),
		CustomerGroupID: CreateUserGroupBodyParams.ID,
	})

	// 跟新客户分组详情
	groupID, createErr := userGroupService.EditCustomerGroup(ctx1, enterpriseID.(int64), CreateUserGroupBodyParams)
	if createErr != nil {
		ginplus.ResError(c, createErr)
		return
	}
	if CreateUserGroupBodyParams.UserGroupType == 1 {
		go func() {
			err := userGroupSvc.NewGroupSvc().TagsHandle(wrapContext(c), enterpriseID.(int64), groupID)
			if err != nil {
				log.Info(wrapContext(c)).Msgf("创建分组同步标签类型的客户分组失败::,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID)
			}
		}()
	}
	if CreateUserGroupBodyParams.UserGroupType == 3 {
		go func() {
			err := userGroupSvc.NewGroupSvc().DepartmentHandle(wrapContext(c), enterpriseID.(int64), groupID)
			if err != nil {
				log.Info(wrapContext(c)).Msgf("创建分组同步门店分组类型的客户分组失败::,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID)
			}
		}()
	}
	ginplus.ResSuccess(c, entity.CreateUserGroupResponse{
		GroupID: groupID,
	})
}

// userGroupDelete 客户分组删除
func userGroupDelete(c *gin.Context) {
	var CreateUserGroupBodyParams entity.CreateUserGroupBody
	err := c.BindJSON(&CreateUserGroupBodyParams)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if CreateUserGroupBodyParams.ID == 0 {
		ginplus.ResError(c, errs.UserGroupIDIsEmpty)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	// 不依赖CDP,删除客户分组
	err = userGroupService.DeleteCustomerGroup(wrapContext(c), enterpriseID.(int64), CreateUserGroupBodyParams.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, struct {
		GroupID int64 `json:"group_id"`
	}{
		GroupID: CreateUserGroupBodyParams.ID,
	})
}

func userGroupDetail(c *gin.Context) {
	var CreateUserGroupBodyParams entity.CreateUserGroupBody
	err := c.BindJSON(&CreateUserGroupBodyParams)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if CreateUserGroupBodyParams.ID == 0 {
		ginplus.ResError(c, errs.UserGroupIDIsEmpty)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	Res, err := userGroupService.UserGroupDetail(wrapContext(c), enterpriseID.(int64), CreateUserGroupBodyParams.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, Res)
}

// userGroupList 客户分组列表
func userGroupList(c *gin.Context) {
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	res, err := userGroupService.GetUserGroupList(wrapContext(c), enterpriseID.(int64))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	result := entity.UserGroupList{
		Aipl:   []entity.UserGroupModel{},
		Rfm:    []entity.UserGroupModel{},
		Custom: make([]entity.CustomUserGroupModel, len(res)),
	}
	for index, data := range res {
		count := 0
		cond := &customer.ListCond{EnterpriseID: enterpriseID.(int64)}
		cond.CustomerGroupCond.Include = true
		cond.CustomerGroupCond.GroupIDList = []int64{data.ID}
		totalRes, _ := userService.QueryUserCount(wrapContext(c), cond)
		if totalRes != nil {
			count = int(totalRes.MemberCount)
		}
		result.Custom[index] = entity.CustomUserGroupModel{
			GroupID:   data.ID,
			GroupName: data.GroupName,
			Value:     strconv.Itoa(count),
		}
	}

	ginplus.ResSuccess(c, result)
}

// 更新自动标签设置
func updateAutoTag(c *gin.Context) {
	var request entity.AutoTagUpdateRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Info(wrapContext(c)).Msgf("updateAutoTag参数解析异常：%s", err.Error())
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	err = tagGroupService.UpdateTagGroup(wrapContext(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	go func() {
		_ = tagService.SyncTagData(wrapContext(c))
	}()
	ginplus.ResSuccess(c, nil)
}

// 更新标签组同步到企业微信配置
func updateTagGroupSync(c *gin.Context) {
	var request entity.UpdateTagGroupSyncRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = tagGroupService.UpdateTagGroupSync(wrapContext(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// 客户分组预计算人数
func userGroupPrecalculate(c *gin.Context) {
	var request entity.UserGroupPrecalculateRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	count, countErr := userGroupService.PreCalculateCount(wrapContext(c), request, enterpriseID.(int64))
	if countErr != nil {
		ginplus.ResError(c, countErr)
		return
	}
	ginplus.ResSuccess(c, entity.UserGroupPrecalculateResponse{
		UserCount: count,
	})
}

// tagGroupInfo 标签组详情
func tagGroupInfo(c *gin.Context) {
	var request entity.TagGroupInfoRequest
	requestErr := c.BindJSON(&request)
	if requestErr != nil {
		ginplus.ResError(c, requestErr)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	params := entity.TagParam{
		EnterpriseID: enterpriseID.(int64),
		TagGroupInfoRequest: entity.TagGroupInfoRequest{
			TagGroupID:  request.TagGroupID,
			TagGroupIds: request.TagGroupIds,
			Order:       request.Order,
			Page:        request.Page,
			Size:        request.Size,
			NeedGoods:   request.NeedGoods,
			NeedCoupon:  request.NeedCoupon,
			TagKeywords: request.TagKeywords,
		},
	}

	// 通过标签组ID，获取组内标签
	tagInfoList, total, tagErr := workExternalTagService.QueryTagListByGroupId(wrapContext(c), params)
	if tagErr != nil {
		ginplus.ResError(c, tagErr)
		return
	}

	var tagIds []int64
	for _, tagInfo := range tagInfoList {
		tagIds = append(tagIds, tagInfo.TagId)
	}

	tagCount := make(map[int64]int64)
	if len(tagIds) > 0 && !request.IgnoreTagCount {
		// 取出标签对应的人数
		tagCount, tagErr = tagsDao.NewWorkExternalUserTag().GetTagCount(ctx.Wrap(c), enterpriseID.(int64), tagIds)
		if tagErr != nil {
			ginplus.ResError(c, tagErr)
			return
		}
	}

	// 循环查询每个标签对应的人数
	for k, v := range tagInfoList {
		if countVal, ok := tagCount[v.TagId]; ok {
			v.UseNum = int(countVal)
		}
		if v.UseNum <= 0 {
			v.IsDelete = 0
		} else {
			v.IsDelete = 1
		}
		tagInfoList[k] = v
	}
	if params.Order > 0 {
		sort.Slice(tagInfoList, func(i, j int) bool {
			if params.Order == 1 {
				return tagInfoList[i].UseNum > tagInfoList[j].UseNum
			}
			if params.Order == 2 {
				return tagInfoList[i].UseNum < tagInfoList[j].UseNum
			}
			return false
		})
	}

	// ginplus.ResSuccess(c, tagInfoList)

	response.SuccessList(c, total, params.Page, params.Size, tagInfoList)
}

// 标签组排序
func sortTagGroup(c *gin.Context) {
	var request entity.SortTagGroupRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	err = tagGroupService.SortTagGroup(wrapContext(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// 标签搜索
func searchTag(c *gin.Context) {
	var request entity.SearchTagRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	res, err := tagGroupService.SearchTag(wrapContext(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// 删除手动标签组
func deleteManualTagGroup(c *gin.Context) {
	var request entity.DeleteManualTagGroupRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	err = workExternalTagGroupService.DeleteTagGroup(wrapContext(c), entity.DeleteManualTagGroupRequestParam{
		EnterpriseID: enterpriseID.(int64),
		TagGroupID:   request.TagGroupID,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, nil)
}

// 更新手动标签
func updateManualTag(c *gin.Context) {
	var request entity.UpdateManualTagRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	// 获取enterprise_id
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	updateErr := workExternalTagService.UpdateManualTag(wrapContext(c), entity.UpdateManualTagParam{
		EnterpriseID: enterpriseID.(int64),
		TagID:        request.TagID,
		NewTagName:   request.TagName,
	})

	if updateErr != nil {
		ginplus.ResError(c, updateErr)
		return
	}

	ginplus.ResSuccess(c, nil)
}

// 删除手动标签
func deleteManualTag(c *gin.Context) {
	var request entity.DeleteManualTagRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	// 获取enterprise_id
	EnterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	enterpriseID := EnterpriseID.(int64)
	if !isExist || enterpriseID <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	// 初始化删除标签参数
	params := entity.TagOperationParam{
		EnterpriseID: enterpriseID,
		TagId:        request.TagID,
	}

	// 删除标签
	delErr := workExternalTagService.DeleteTag(wrapContext(c), params)
	if delErr != nil {
		ginplus.ResError(c, delErr)
		return
	}

	ginplus.ResSuccess(c, nil)
}

func TagSync(c *gin.Context) {
	_ = tagService.SyncTagData(wrapContext(c))
	_ = tagService.SyncUserGroup(wrapContext(c))
}

func getTagList(c *gin.Context) {
	var request entity.GetTagListRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	var list []*entity.GetTagListItem
	var total int64
	if request.Type == 1 {
		total, list, err = tagService.GetTagList(wrapContext(c), &request)
	} else {
		total, list, err = tagService.GetUserGroupList(wrapContext(c), &request)
	}

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var ret = map[string]interface{}{
		"list":  list,
		"page":  request.Page,
		"size":  request.Size,
		"total": total,
	}

	ginplus.ResSuccess(c, ret)
}

func EnterpriseInit(c *gin.Context) error {
	err := tagGroupService.TagEnterpriseInit(wrapContext(c), nil)
	if err != nil {
		return err
	}
	return nil
}

func wrapContext(c *gin.Context) ctx.Context {
	cc := ctx.Wrap(c)
	enterpriseID := common.EnterpriseID(cc)
	enterpriseInfo, err := enterpriseService.Detail(cc, enterprise.Corp{EnterpriseID: enterpriseID})
	if err != nil {
		log.Info(cc).Msgf("获取企业信息失败, err:%s, enterpriseID:%d", err.Error(), enterpriseID)
		return cc
	}
	cc = ctx.Wrap(context.WithValue(cc, common.ContextKeyEnum.EnterpriseHash, enterpriseInfo.Hash))
	return cc
}

func customerManualList(c *gin.Context) {
	req := &entity.CustomerManualTagRequest{}
	err := c.BindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	ctxContent := wrapContext(c)
	tagGroupList, err := tagGroupService.TagGroupList(ctxContent, req.EnterpriseID, 0)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	userTagListReturn := tagService.CustomerTags(ctxContent, req.EnterpriseID, req.RelationID, tagGroupList)
	ginplus.ResSuccess(c, userTagListReturn)
}

// SyncMarkingTag 异步打标
func syncMarkingTag(c ctx.Context, workTagIds []string, enterpriseId int64, fileUrl string) error {
	mapData, err := helper2.ReadRemoteUrlData(c, fileUrl, 0)
	if err != nil {
		return err
	}

	// 获取文件里面的成员信息
	var unionIds []string
	var customerIds []string
	var saCode []string

	for key, data := range mapData {
		// 优先取unionId 如果不存在则取customerId
		unionIdVal, unionIdOk := data["unionId"]
		if unionIdOk && unionIdVal.(string) != "" {
			unionIds = append(unionIds, unionIdVal.(string))
		} else {
			if val, ok := data["customerId"]; ok {
				customerIds = append(customerIds, val.(string))
			}
		}

		if val, ok := data["saCode"]; ok {
			saCode = append(saCode, val.(string))
		}

		// 判断 取出的unionIds和customerIds是否等于调用企微的阈值，或是否为最后一条
		if len(unionIds)+len(customerIds) == 2000 || key+1 == len(mapData) {
			externalResult, _ := microScrmDao.NewWxUserCorpExternal().GetMany(c, enterpriseId, "customer_code in (?) or union_code in (?)", customerIds, customerIds)
			var externalUserIds []int64
			for _, external := range externalResult {
				externalUserIds = append(externalUserIds, external.ID)
			}
			if len(externalUserIds) == 0 && len(unionIds) == 0 {
				continue
			}
			// 查询与当前 UnionIds 有关系的员工
			userRelation, relationErr := microScrmDao.NewUserRelation().GetInfoByUnionIdsOrCustomers(c, enterpriseId, unionIds, externalUserIds)
			if relationErr != nil {
				log.Err(c, relationErr).Msgf("查询与当前UnionIds有关系的员工")
				continue
			}

			var externalUserTags []tagsReq.ExternalUserTag
			for _, tagId := range workTagIds {
				for _, relation := range userRelation {
					tmpUserTags := tagsReq.ExternalUserTag{TagId: tagId}
					// 判断是否上传用户id,如果用户id不在关系中则跳过
					if len(saCode) > 0 && !helper2.InArray(relation.UserID, saCode) {
						continue
					}
					tmpUserTags.UserId = relation.UserID
					tmpUserTags.UnionId = relation.Unionid
					externalUserTags = append(externalUserTags, tmpUserTags)
				}
			}

			err = tagsSvc.NewWorkTags().BatchMarking(c, "客户标签", enterpriseId, externalUserTags, 1)
			if err != nil {
				log.Err(c, err).Msgf("批量打标异常")
				notice.NewWechat().SendTextMessage(c, fmt.Sprintf("客户标签-批量打标异常,错误:[%+v]", err))
			}

			// 清空已打标的数据
			unionIds, customerIds, saCode = []string{}, []string{}, []string{}
		}
	}

	return nil
}

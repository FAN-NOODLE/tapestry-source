package moment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/uid"
	"strconv"
	"strings"
	"time"
)

type TopicService struct {
	repo Repository
}

type MemberService struct {
	repo Repository
}

type TopicMemberService struct {
	repo Repository
}

// MemberGetInfo 获取朋友圈形象
func (s *MemberService) MemberGetInfo(c ctx.Context, memberRequest MomentMemberGetRequest,
	enterpriseService *enterprise.Service) (*MomentMemberGetResponse, error) {
	momentMember, error := s.repo.MemberFindOne(c, memberRequest)
	if error != nil {
		return nil, error
	}
	domain := config.CfgData.Service.WebH5.Domain

	corpRequest := enterprise.Corp{
		EnterpriseID: memberRequest.EnterpriseId,
	}
	corp, error := enterpriseService.Detail(c, corpRequest)
	if error != nil {
		return nil, error
	}
	corpId := corp.CorpID
	fmt.Println(corpId)

	req := &boss.EnterpriseDetailRequest{
		EnterpriseID: memberRequest.EnterpriseId,
	}
	enterpriseData, _ := boss.EnterpriseDetail(c, req)
	var EnterpriseHashValue = ""
	if enterpriseData != nil {
		EnterpriseHashValue = corp.Hash
	} else {
		EnterpriseHashValue = "public"
	}

	link := ""
	link += domain
	link += "/friend/"
	link += EnterpriseHashValue
	link += "/friend-circle-list?mark="
	link += memberRequest.Mark

	/*var link string = ""
	link += domain
	link += "/"
	link += corpId
	link += "/crs-manu-customer/friend-circle-list?mark="
	link += memberRequest.Mark
	link += "&hash="
	link += corp.Hash*/

	momentMemberGetResponse := &MomentMemberGetResponse{
		ID:           momentMember.ID,
		MemberId:     momentMember.MemberId,
		Name:         momentMember.Name,
		Sign:         momentMember.Sign,
		Mark:         momentMember.Mark,
		Avatar:       momentMember.Avatar,
		EnterpriseId: momentMember.EnterpriseId,
		CreatedAt:    momentMember.CreatedAt,
		UpdatedAt:    momentMember.UpdatedAt,
		Background:   momentMember.Background,
		DeletedAt:    momentMember.DeletedAt,
		Link:         link,
	}
	//羽西同步用户查看朋友圈足迹
	//if isCan, authErr := auth.BrandAuth(c, auth.BrandAuthConfig{
	//	EnterpriseID: momentMember.EnterpriseId,
	//	Business:     auth.BUSINESS_USER_VIEW_MOMENT_EVENT,
	//}); isCan && authErr == nil {
	//
	//}
	//end羽西同步用户足迹
	eModel, err := orm.Service.WorkExternalUser.GetWxUserCorpExternalInfoByUnionID(momentMember.EnterpriseId, memberRequest.UnionId)
	if eModel.Euid != 0 {
		model := micro_scrm.UserEventModel{
			UserId: eModel.Euid,
			Event:  common.EventEnum.ViewMoment.String(),
			Object: strconv.FormatInt(momentMember.MemberId, 10),
			Extra:  "",
		}
		model.ID = uid.ID(c)
		model.EnterpriseID = momentMember.EnterpriseId
		err = orm.Service.UserEvent.CreateModel(model)
		if err != nil {
			return nil, err
		}
	}

	return momentMemberGetResponse, nil
}

// MemberEditInfo 编辑朋友圈形象
func (s *MemberService) MemberEditInfo(c ctx.Context,
	momentMemberEditRequest MomentMemberEditRequest) (int64, error) {

	momentMemberRequest := MomentMemberGetRequest{
		Mark:         momentMemberEditRequest.Mark,
		EnterpriseId: momentMemberEditRequest.EnterpriseId,
	}

	oldMomentMember, queryError := s.repo.MemberFindOne(c, momentMemberRequest)
	if queryError != nil {
		return 0, queryError
	}

	newMomentMember := MomentMember{
		Name:         momentMemberEditRequest.Name,
		Sign:         momentMemberEditRequest.Sign,
		Mark:         oldMomentMember.Mark,
		Avatar:       momentMemberEditRequest.Avatar,
		Background:   momentMemberEditRequest.Background,
		ID:           oldMomentMember.ID,
		MemberId:     oldMomentMember.MemberId,
		EnterpriseId: oldMomentMember.EnterpriseId,
		CreatedAt:    oldMomentMember.CreatedAt,
		UpdatedAt:    time.Now().Unix(),
	}

	return s.repo.MemberUpdateOne(c, newMomentMember)
}

// TopicEdit 编辑朋友圈
func (s *MemberService) TopicEdit(c ctx.Context, topicEditRequest TopicEditRequest,
	topicService *TopicService, topicMemberService *TopicMemberService, memberService *MemberService, originalMemberService *members.Service) (int64, error) {
	topicDeleteRequest := TopicDeleteRequest{
		ID: topicEditRequest.ID,
	}

	_, delError := s.TopicDelete(c, topicDeleteRequest)
	if delError != nil {
		return 0, delError
	}

	topicCreateRequest := TopicCreateRequest{
		Name:            topicEditRequest.Name,
		PublishType:     topicEditRequest.PublishType,
		PublishTime:     topicEditRequest.PublishTime,
		EnterpriseId:    topicEditRequest.EnterpriseId,
		CreatedMemberId: topicEditRequest.CreatedMemberId,
		MemberIds:       topicEditRequest.MemberIds,
		DepartmentIds:   topicEditRequest.DepartmentIds,
		Content:         topicEditRequest.Content,
	}

	newId, createErr := s.TopicCreate(c, topicCreateRequest, topicService, topicMemberService, memberService, originalMemberService)
	if createErr != nil {
		return 0, createErr
	}

	return newId, nil
}

// getInsertMemberMap 组装moment_members数据
func (s *MemberService) getInsertMemberMap(c ctx.Context, topicCreateRequest TopicCreateRequest, currentTime int64) (map[int64]MomentMember, int64) {
	var memberCount int64 = 0
	insertMemberMap := make(map[int64]MomentMember)
	for _, memberVal := range topicCreateRequest.MemberIds {
		hashStr := s.GetMD5Encode(strconv.FormatInt((topicCreateRequest.EnterpriseId + memberVal + time.Now().Unix()), 10))
		insertMemberMap[memberVal] = MomentMember{
			ID:           uid.ID(c),
			MemberId:     memberVal,
			Name:         "",
			Sign:         "",
			Mark:         hashStr,
			Avatar:       "",
			Background:   "",
			EnterpriseId: topicCreateRequest.EnterpriseId,
			CreatedAt:    currentTime,
			UpdatedAt:    currentTime,
			DeletedAt:    0,
		}
		memberCount++
	}
	return insertMemberMap, memberCount
}

// getInsertTopic 组装moment_topic数据
func (s *MemberService) getInsertTopic(c ctx.Context, topicCreateRequest TopicCreateRequest, currentTime int64, memberCount int64) (MomentTopic, error) {
	insertTopic := MomentTopic{}
	contentJson, jsonErr := json.Marshal(topicCreateRequest.Content)
	if jsonErr != nil {
		return insertTopic, errs.JsonParseErr
	}
	insertTopic.ID = uid.ID(c)
	insertTopic.Name = topicCreateRequest.Name
	insertTopic.PublishType = topicCreateRequest.PublishType
	insertTopic.PublishTime = topicCreateRequest.PublishTime
	insertTopic.Content = string(contentJson)
	insertTopic.CreatedMemberId = topicCreateRequest.CreatedMemberId
	insertTopic.MemberCount = memberCount
	insertTopic.EnterpriseId = topicCreateRequest.EnterpriseId
	insertTopic.CreatedAt = currentTime
	insertTopic.UpdatedAt = currentTime
	insertTopic.DeletedAt = 0
	return insertTopic, nil
}

// updateMemberTime 更改moment_member表的更新时间
func (s *MemberService) updateMemberTime(c ctx.Context, topicCreateRequest TopicCreateRequest) (bool, error) {
	batchUpdateMemberTime := BatchUpdateMemberTime{}
	batchUpdateMemberTime.DeletedAt = 0
	batchUpdateMemberTime.EnterpriseId = topicCreateRequest.EnterpriseId
	batchUpdateMemberTime.UpdatedAt = time.Now().Unix()
	for _, MomentMemberId := range topicCreateRequest.MemberIds {
		batchUpdateMemberTime.MemberIDs = append(batchUpdateMemberTime.MemberIDs, MomentMemberId)
	}

	ok, err := s.repo.BatchUpdateMemberTime(c, batchUpdateMemberTime)
	if err != nil || ok == false {
		return false, errs.UpdateMemberTimeErr
	} else {
		return true, nil
	}
}

// insertTopicMembers 组装topicMember数据
func (s *MemberService) insertTopicMembers(c ctx.Context, topicCreateRequest TopicCreateRequest,
	insertMemberMap map[int64]MomentMember, insertTopic MomentTopic, currentTime int64) []*MomentTopicMember {
	insertTopicMembers := make([]*MomentTopicMember, 0)
	for _, insertMemberStruct := range insertMemberMap {
		insertTopicMember := MomentTopicMember{
			ID:             uid.ID(c),
			MomentTopicId:  insertTopic.ID,
			MomentMemberId: insertMemberStruct.MemberId,
			EnterpriseId:   topicCreateRequest.EnterpriseId,
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
			DeletedAt:      0,
		}

		insertTopicMembers = append(insertTopicMembers, &insertTopicMember)
	}
	return insertTopicMembers
}

// TopicCreate 发布朋友圈
func (s *MemberService) TopicCreate(c ctx.Context, topicCreateRequest TopicCreateRequest,
	topicService *TopicService, topicMemberService *TopicMemberService, memberService *MemberService, originalMemberService *members.Service) (int64, error) {

	currentTime := time.Now().Unix()
	// 定时发布时间不能小于当前时间
	if topicCreateRequest.PublishType == 2 && topicCreateRequest.PublishTime < currentTime {
		return 0, errs.PublishTimeErr
	}

	// 组装moment_members数据
	insertMemberMap, memberCount := s.getInsertMemberMap(c, topicCreateRequest, currentTime)

	existQuery := MemberBatchGetByMemberIds{}
	existQuery.EnterpriseId = topicCreateRequest.EnterpriseId
	insertMemberSlice := make([]*MomentMember, 0)
	for _, insertMember := range insertMemberMap {
		tmpMember := insertMember
		insertMemberSlice = append(insertMemberSlice, &tmpMember)
		existQuery.MemberIds = append(existQuery.MemberIds, tmpMember.MemberId)
	}

	// 去除已经存在的moment_members数据
	existMembers, _ := topicMemberService.MemberBatchGetByMemberIds(c, existQuery)
	BatchGetMemberByMemberId := members.BatchGetMemberByMemberId{}
	BatchGetMemberByMemberId.DeletedAt = 0
	BatchGetMemberByMemberId.EnterpriseID = topicCreateRequest.EnterpriseId
	notExistMembers := make([]*MomentMember, 0)
	notExistMemberIds := make([]int64, 0)
	for _, insertMemberVal := range insertMemberSlice {
		isExist := false
		for _, existMember := range *existMembers {
			if existMember.MemberId == insertMemberVal.MemberId {
				isExist = true
				break
			}
		}
		if isExist == false {
			notExistMembers = append(notExistMembers, insertMemberVal)
			notExistMemberIds = append(notExistMemberIds, insertMemberVal.MemberId)
			BatchGetMemberByMemberId.IDs = append(BatchGetMemberByMemberId.IDs, insertMemberVal.MemberId)
		}
		isExist = false
	}

	// 组装moment_topics数据
	insertTopic, err := s.getInsertTopic(c, topicCreateRequest, currentTime, memberCount)
	if err != nil {
		return 0, err
	}

	// 组装moment_topic_members数据
	insertTopicMembers := s.insertTopicMembers(c, topicCreateRequest, insertMemberMap, insertTopic, currentTime)

	// 新增moment_members数据
	_, insertTopicError := topicService.AddMomentTopic(c, insertTopic)
	if insertTopicError != nil {
		return 0, errs.TopicInsertErr
	}

	// 新增moment_members数据
	if len(notExistMembers) > 0 {
		membersInfo, getMemberError := originalMemberService.GetMemberByMemberIds(c, &BatchGetMemberByMemberId)
		if getMemberError == nil && len(*membersInfo) != 0 {
			for _, notExistMemberVal := range notExistMembers {
				for _, membersInfoVal := range *membersInfo {
					if membersInfoVal.ID == notExistMemberVal.MemberId {
						notExistMemberVal.Name = membersInfoVal.Name
						notExistMemberVal.Avatar = membersInfoVal.Avatar
						break
					}
				}
			}
		}
		_, insertMemberError := memberService.AddMomentMembers(c, notExistMembers)
		if insertMemberError != nil {
			return 0, errs.MemberInsertErr
		}
	}

	// 新增moment_topic_members数据
	_, insertTopicMemberError := topicMemberService.AddMomentTopicMembers(c, insertTopicMembers)
	if insertTopicMemberError != nil {
		return 0, errs.TopicMemberInsertErr
	}

	// 更新moment_member的更新时间
	ok, err := s.updateMemberTime(c, topicCreateRequest)
	if !ok || err != nil {
		return 0, err
	}

	return insertTopic.ID, nil
}

// GetDeeplyDps 根据third_ids最终查询后的所有部门
func (s *MemberService) GetDeeplyDps(c ctx.Context, memberRequest TopicCreateRequest,
	departmentService *department.Service) (*[]int64, error) {

	finalDps := make([]int64, 0)

	batchDept := &department.BatchGetDepartment{}
	batchDept.DeletedAt = 0

	for _, depId := range memberRequest.DepartmentIds {
		batchDept.ThirdIDs = append(batchDept.ThirdIDs, depId)
		finalDps = append(finalDps, depId)
	}

	for {
		dps, queryErr := departmentService.BatchGetDepartmentByThirdIds(c, batchDept)
		if queryErr != nil {
			break
		}
		batchDept.ThirdIDs = batchDept.ThirdIDs[0:0]
		for _, dept := range *dps {
			finalDps = append(finalDps, dept.ThirdID)
			batchDept.ThirdIDs = append(batchDept.ThirdIDs, dept.ThirdID)
		}
	}

	return &finalDps, nil
}

// AddMomentTopic 新增Topic
func (t *TopicService) AddMomentTopic(c ctx.Context, momentTopic MomentTopic) (int64, error) {
	return t.repo.MomentTopicAddOne(c, momentTopic)
}

// AddMomentMembers 新增moment_members
func (m *MemberService) AddMomentMembers(c ctx.Context, momentMembers []*MomentMember) (int64, error) {
	return m.repo.MomentMemberAddMany(c, momentMembers)
}

// AddMomentTopicMembers 新增moment_topic_members
func (t *TopicMemberService) AddMomentTopicMembers(c ctx.Context, momentTopicMembers []*MomentTopicMember) (int64, error) {
	return t.repo.MomentTopicMemberAddMany(c, momentTopicMembers)
}

// 获取32位MD5加密字符串
func (s *MemberService) GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// TopicDisplay 一个朋友圈展示
func (s *MemberService) TopicDisplay(c ctx.Context, topicDisplayRequest TopicDisplayRequest) (*TopicDisplayResponse, error) {
	momentTopicQuery := MomentTopic{
		ID: topicDisplayRequest.ID,
	}
	momentTopic, momentTopicFindError := s.repo.MomentTopicFindOne(c, &momentTopicQuery)

	if momentTopicFindError != nil {
		return nil, errs.TopicDisplayErr
	}

	contentJson := make(map[string]interface{})
	jsonError := json.Unmarshal([]byte(momentTopic.Content), &contentJson)
	if jsonError != nil {
		return nil, errs.JsonParseErr
	}

	offset := (topicDisplayRequest.Page - 1) * topicDisplayRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(topicDisplayRequest.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_moments/select-topic-display",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]TopicDisplayExt{}
		},
	}

	total, list, err := s.repo.FindPage(c, pageVo, topicDisplayRequest, option)

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	resList := list.(*[]TopicDisplayExt)
	topicDisplayResponse := &TopicDisplayResponse{
		ID:              momentTopic.ID,
		Name:            momentTopic.Name,
		PublishType:     momentTopic.PublishType,
		PublishTime:     momentTopic.PublishTime,
		Content:         contentJson,
		CreatedMemberId: momentTopic.CreatedMemberId,
		MemberCount:     momentTopic.MemberCount,
		CreatedAt:       momentTopic.CreatedAt,
		UpdatedAt:       momentTopic.UpdatedAt,
		EnterpriseId:    momentTopic.EnterpriseId,
		Page:            topicDisplayRequest.Page,
		Size:            topicDisplayRequest.Size,
		Total:           total,
		List:            *resList,
	}

	return topicDisplayResponse, nil
}

// TopicDelete 删除朋友圈
func (s *MemberService) TopicDelete(c ctx.Context, topicDeleteRequest TopicDeleteRequest) (*TopicDeleteResponse, error) {
	response := &TopicDeleteResponse{
		Result: true,
	}
	// 查询该朋友圈是否存在
	topicDisplayRequest := TopicDisplayRequest{
		ID:   topicDeleteRequest.ID,
		Page: 1,
		Size: 1,
	}
	topicDisplay, displayError := s.TopicDisplay(c, topicDisplayRequest)
	if displayError != nil {
		return nil, errs.TopicDisplayErr
	}

	currentTime := time.Now().Unix()

	// 删除朋友圈
	topic := &MomentTopic{
		ID:        topicDisplay.ID,
		DeletedAt: currentTime,
	}

	_, topicUpdateError := s.repo.MomentTopicUpdateOne(c, topic)
	if topicUpdateError != nil {
		return nil, errs.TopicDeleteErr
	}

	// 删除中间表
	topicMemberBatchDelete := &TopicMemberBatchDelete{
		DeletedAt:     time.Now().Unix(),
		MomentTopicId: topicDeleteRequest.ID,
	}

	_, memberDeleteError := s.repo.MomentTopicMemberBatchDelete(c, topicMemberBatchDelete)
	if memberDeleteError != nil {
		return nil, errs.MemberDeleteError
	}

	return response, nil
}

// ContentList 按内容搜索
func (s *MemberService) ContentList(c ctx.Context, contentListRequest ContentListRequest) (*ContentListResponse, error) {

	defaultResponse := &ContentListResponse{
		List:      make([]TopicContent, 0),
		Total:     0,
		UpdatedAt: 0,
		Page:      contentListRequest.Page,
		Size:      contentListRequest.Size,
	}

	contentListRequest.CurrentTime = time.Now().Unix()
	offset := (contentListRequest.Page - 1) * contentListRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(contentListRequest.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_moments/topic-list-content",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]TopicContent{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, contentListRequest, option)

	if err != nil || total == 0 {
		return defaultResponse, nil
	}

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	resList := list.(*[]TopicContent)
	resFormat := make([]TopicContent, 0)
	for _, res := range *resList {
		tmpRes := res
		if res.PublishType == 1 {
			tmpRes.Status = 1
		} else {
			if res.PublishTime < time.Now().Unix() {
				tmpRes.Status = 1
			} else {
				tmpRes.Status = 2
			}
		}
		contentJson := make(map[string]interface{})
		jsonError := json.Unmarshal([]byte(tmpRes.Content), &contentJson)
		if jsonError != nil {
			return nil, errs.JsonParseErr
		}
		tmpRes.ContentJson = contentJson
		resFormat = append(resFormat, tmpRes)
	}

	momentListResponse := &ContentListResponse{
		List:      resFormat,
		Total:     total,
		UpdatedAt: resFormat[0].UpdatedAt,
		Page:      contentListRequest.Page,
		Size:      contentListRequest.Size,
	}

	return momentListResponse, nil
}

// MemberList 按员工搜索
func (s *MemberService) MemberList(c ctx.Context, memberListRequest MemberListRequest,
	departmentService *department.Service, originalMemberService *members.Service) (*MemberListResponse, error) {

	defaultResponse := &MemberListResponse{
		List:  make([]TopicMember, 0),
		Total: 0,
		Page:  memberListRequest.Page,
		Size:  memberListRequest.Size,
	}

	if memberListRequest.DepartmentId > 0 {
		memberList, queryError := originalMemberService.GetMemberByDpId(c, memberListRequest.DepartmentId, memberListRequest.EnterpriseId)
		if queryError == nil || memberList != nil {
			for _, memberRow := range *memberList {
				memberListRequest.MemberIds = append(memberListRequest.MemberIds, memberRow.MemberID)
			}
		}
	}

	if memberListRequest.MemberId > 0 {
		if len(memberListRequest.MemberIds) > 0 {
			memberExist := false
			for _, memberId := range memberListRequest.MemberIds {
				if memberId == memberListRequest.MemberId {
					// 搜索条件应该只有一个员工
					memberExist = true
					memberListRequest.MemberIds = make([]interface{}, 0)
					memberListRequest.MemberIds = append(memberListRequest.MemberIds, memberListRequest.MemberId)
				}
			}
			if memberExist == false {
				return defaultResponse, nil
			}
		} else {
			memberListRequest.MemberIds = append(memberListRequest.MemberIds, memberListRequest.MemberId)
		}
	}

	offset := (memberListRequest.Page - 1) * memberListRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(memberListRequest.Size)}

	option := &repos.FindPageOption{
		TplPrefix: "t_moments/topic-list-member",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]TopicMember{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, memberListRequest, option)

	if total == 0 {
		return defaultResponse, nil
	}

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	resList := list.(*[]TopicMember)

	// 根据员工获取部门信息数组、部门信息字符串、员工姓名
	batchGetDpnameByMembers := &department.BatchGetDpnameByMembers{}
	batchGetDpnameByMembers.DeletedAt = 0
	batchGetDpnameByMembers.EnterpriseId = memberListRequest.EnterpriseId
	for _, res := range *resList {
		batchGetDpnameByMembers.MemberIds = append(batchGetDpnameByMembers.MemberIds, res.MemberId)
	}

	membersInfo, getInfoErr := departmentService.DepartnameByMembers(c, batchGetDpnameByMembers)

	if getInfoErr != nil {
		return nil, getInfoErr
	}

	returnList := make([]TopicMember, 0)
	for _, res := range *resList {
		tmpRes := TopicMember{}
		tmpRes.MemberId = res.MemberId
		tmpRes.TopicCount = res.TopicCount
		tmpRes.UpdatedAt = res.UpdatedAt
		tmpRes.Mark = res.Mark
		tmpStr := ""
		for _, memberInfo := range *membersInfo {
			if memberInfo.MemberId == res.MemberId {
				tmpStr += ";" + memberInfo.DepartmentStr
				tmpRes.DepartmentId = memberInfo.DepartmentId
				tmpRes.DepartmentName = memberInfo.DepartmentName
				tmpRes.UserName = memberInfo.UserName
			}
		}
		tmpStr = strings.Replace(tmpStr, ";", "", 1)
		tmpRes.DepartmentStr = tmpStr
		returnList = append(returnList, tmpRes)
	}

	memberListResponse := &MemberListResponse{
		List:      returnList,
		Total:     total,
		Page:      memberListRequest.Page,
		Size:      memberListRequest.Size,
		UpdatedAt: returnList[0].UpdatedAt,
	}

	return memberListResponse, nil
}

// MemberTopic 一个员工的朋友圈列表
func (s *MemberService) MemberTopic(c ctx.Context, request MemberTopicRequest) (*MemberTopicResponse, error) {
	request.PublishTime = time.Now().Unix()
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(request.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_moments/member-topic-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberTopicExt{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, request, option)

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	resList := list.(*[]MemberTopicExt)
	returnList := make([]MemberTopicExt, 0)
	for _, res := range *resList {
		contentJson := make(map[string]interface{})
		jsonError := json.Unmarshal([]byte(res.Content), &contentJson)
		if jsonError != nil {
			return nil, errs.JsonParseErr
		}
		res.ContentJson = contentJson
		returnList = append(returnList, res)
	}

	response := &MemberTopicResponse{
		List:  returnList,
		Total: total,
		Page:  request.Page,
		Size:  request.Size,
	}
	return response, nil
}

// MemberBatchGetByMemberIds 批量获取朋友圈用户信息
func (t *TopicMemberService) MemberBatchGetByMemberIds(c ctx.Context, request MemberBatchGetByMemberIds) (*[]MomentMember, error) {
	return t.repo.MemberBatchGetByMemberIds(c, request)
}

// TopicMemberDelete 一个员工的朋友圈删除
func (s *MemberService) TopicMemberDelete(c ctx.Context, request TopicMemberDeleteRequest) (*TopicMemberDeleteResponse, error) {
	// 获取moment_topics
	topicCond := &MomentTopic{
		ID: request.MomentTopicId,
	}
	topic, queryErr := s.repo.MomentTopicFindOne(c, topicCond)
	if queryErr != nil {
		return nil, errs.TopicDisplayErr
	}
	if topic.MemberCount == 1 {
		topic.MemberCount = 0
	} else {
		topic.MemberCount -= 1
	}

	// 获取moment_topics.member_count-1
	_, updateErr := s.repo.TopicUpdateMemberCount(c, *topic)
	if updateErr != nil {
		return nil, updateErr
	}

	if topic.MemberCount == 0 {
		topicDeleteRequest := TopicDeleteRequest{
			ID: request.MomentTopicId,
		}
		// 删除moment_topic
		_, delError := s.TopicDelete(c, topicDeleteRequest)
		if delError != nil {
			return nil, delError
		}
	} else {
		// 删除moment_topic_members表ID
		request.DeletedAt = time.Now().Unix()
		_, deleteErr := s.repo.TopicMemberDeleteOne(c, request)
		if deleteErr != nil {
			return nil, deleteErr
		}
	}

	response := &TopicMemberDeleteResponse{Result: true}
	return response, nil
}

// TopicMemberEdit 一个员工的朋友圈编辑
func (s *MemberService) TopicMemberEdit(c ctx.Context, request TopicMemberEditRequest,
	topicService *TopicService, topicMemberService *TopicMemberService, memberService *MemberService, originalMemberService *members.Service) (*TopicMemberEditResponse, error) {
	topicMemberDeleteRequest := TopicMemberDeleteRequest{
		MemberId:            request.EditMemberId,
		EnterpriseId:        request.EnterpriseId,
		MomentTopicId:       request.MomentTopicId,
		MomentTopicMemberId: request.MomentTopicMemberId,
	}

	_, delErr := s.TopicMemberDelete(c, topicMemberDeleteRequest)
	if delErr != nil {
		return nil, delErr
	}

	topicCreateRequest := TopicCreateRequest{
		Name:            request.Name,
		PublishType:     request.PublishType,
		PublishTime:     request.PublishTime,
		EnterpriseId:    request.EnterpriseId,
		CreatedMemberId: request.CreatedMemberId,
		EditMemberId:    request.EditMemberId,
		MemberIds:       request.MemberIds,
		DepartmentIds:   request.DepartmentIds,
		Content:         request.Content,
	}

	lastInsertId, insertErr := s.TopicCreate(c, topicCreateRequest, topicService, topicMemberService, memberService, originalMemberService)
	if insertErr != nil {
		return nil, insertErr
	}

	response := &TopicMemberEditResponse{
		MomentTopicId: lastInsertId,
	}
	return response, nil
}

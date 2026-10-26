package http

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gorm.io/gorm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"sync"
	"time"

	"github.com/shakinm/xlsReader/xls/structure"

	"github.com/shakinm/xlsReader/xls"

	"micro-scrm/internal/modules/department"

	"micro-scrm/internal/modules/customer"

	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/sdk/work"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/chat"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/socialDraw"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
)

var groupChatService socialDraw.Service
var drainageService drainage.DrainageService
var memberService members.Service
var chatService chat.GroupChatService
var departmentService department.Service
var customerService customer.Service
var groupChatMemberService chat.GroupChatMemberService

// workGroupChatList 获取群列表数据
func workGroupChatList(c *gin.Context) {
	var req socialDraw.WorkGroupListRequest
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}

	// 运营者名称+所属部门 -> 运营者userid
	if req.Operator != "" || req.DepartmentID != 0 {
		req.OperatorIds = groupChatService.GetWorkGroupChatListFilters(ctx.Wrap(c), req)
		if len(req.OperatorIds) == 0 {
			ginplus.ResSuccess(c, []*socialDraw.WorkGroupChatResponse{})
			return
		}
	}
	chatList, err := groupChatService.GetGroupChatList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if chatList == nil {
		ginplus.ResSuccess(c, chatList)
		return
	}

	var (
		wg         sync.WaitGroup
		staticsMap map[int64]*socialDraw.WorkGroupChatStatistic
		staticsReq = socialDraw.WorkGroupChatStaticsRequest{
			StartAt: req.StartAt,
			EndAt:   req.EndAt,
		}
		operatorIds   []string
		ownerIds      []int64
		channelIds    []int64
		ownerMap      = make(map[int64]members.Member)
		operatorMap   = make(map[string]*members.Member)
		departmentMap = make(map[int64]string)
		channelMap    = make(map[int64]string)
	)
	// 群主，所属渠道，群运营者+所属部门
	for _, v := range chatList {
		if v.OperatorMemberID != "" {
			operatorIds = append(operatorIds, v.OperatorMemberID)
		}
		if v.OwnerMemberID != 0 {
			ownerIds = append(ownerIds, v.OwnerMemberID)
		}
		if v.ChannelID != 0 {
			channelIds = append(channelIds, v.ChannelID)
		}
		staticsReq.GroupChatID = append(staticsReq.GroupChatID, v.ID)
	}

	// 群运营者
	if len(operatorIds) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var departmentIds []int64
			operators, _ := memberService.GetMembersByUserIdsOrMobiles(ctx.Wrap(c), req.EnterpriseID, members.MemberByUserID, operatorIds)
			for _, v := range operators {
				operatorMap[v.Userid] = v
				departmentIds = append(departmentIds, v.Department...)
			}
			departments, _ := departmentService.GetDepartmentsByIds(ctx.Wrap(c), req.EnterpriseID, departmentIds)
			for _, v := range departments {
				departmentMap[v.ID] = v.Name
			}
		}()
	}

	// 群主
	wg.Add(1)
	go func() {
		defer wg.Done()
		var ownerReq = &members.BatchGetMemberByMemberId{EnterpriseID: req.EnterpriseID}
		for _, id := range ownerIds {
			ownerReq.IDs = append(ownerReq.IDs, id)
		}
		owners, _ := memberService.GetMemberByMemberIds(ctx.Wrap(c), ownerReq)
		for _, v := range *owners {
			ownerMap[v.ID] = v
		}
	}()

	// 渠道
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, channels, _ := drainageService.SimpleChannelsList(ctx.Wrap(c), drainage.ContactChannelListRequest{EnterpriseID: req.EnterpriseID})
		for _, v := range *channels {
			channelMap[v.ID] = v.Name
		}
	}()

	// 社群统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		staticsMap, _ = groupChatService.GetWorkGroupChatStatics(ctx.Wrap(c), staticsReq)
	}()

	wg.Wait()

	for _, v := range chatList {
		if statics, ok := staticsMap[v.ID]; ok {
			v.WorkGroupChatStatistic = *statics
		}
		if operator, ok := operatorMap[v.OperatorMemberID]; ok {
			v.Operator = operator.Name
			for _, dpId := range operator.Department {
				if dpName, ok := departmentMap[dpId]; ok {
					v.OperatorDepartments = append(v.OperatorDepartments, dpName)
				}
			}
		}
		if channelName, ok := channelMap[v.ChannelID]; ok {
			v.ChannelName = channelName
		}
		if ownerMember, ok := ownerMap[v.OwnerMemberID]; ok {
			v.Owner = ownerMember.Name
		}
	}

	ginplus.ResSuccess(c, chatList)
	return
}

// workGroupChatInfo 社群详情
func workGroupChatInfo(c *gin.Context) {
	var req chat.WorkGroupChatReq
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}
	if req.ChatID == "" {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	baseChat, err := chatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{
		EnterpriseID: req.EnterpriseID,
		ChatID:       req.ChatID,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	info := &chat.WorkGroupChatBase{WorkGroupChat: *baseChat}

	getWorkGroupChatExtInfo(ctx.Wrap(c), req.EnterpriseID, info)

	ginplus.ResSuccess(c, info)
	return
}

// getWorkGroupChatExtInfo 补全群详情相关信息
func getWorkGroupChatExtInfo(ctx ctx.Context, enterpriseId int64, info *chat.WorkGroupChatBase) {
	var (
		wg        sync.WaitGroup
		userIds   []string
		memberMap = make(map[string]*members.Member)
	)
	_ = json.Unmarshal(info.TmpAdmins, &userIds)
	if info.OperatorMemberID != "" {
		userIds = append(userIds, info.OperatorMemberID)
	}

	// 运营者 + 所属部门 + 群管理员
	if len(userIds) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			memberInfos, _ := memberService.GetMembersByUserIdsOrMobiles(ctx, enterpriseId, members.MemberByUserID, userIds)
			for _, v := range memberInfos {
				memberMap[v.Userid] = v
			}
			if info.OperatorMemberID != "" && memberMap[info.OperatorMemberID] != nil {
				info.Operator = memberMap[info.OperatorMemberID].Name
				dps, _ := departmentService.GetDepartmentsByIds(ctx, enterpriseId, memberMap[info.OperatorMemberID].Department)
				for _, v := range dps {
					info.OperatorDepartments = append(info.OperatorDepartments, v.Name)
				}
			}
			for userId, m := range memberMap {
				if userId == info.OperatorMemberID {
					continue
				}
				info.Admins = append(info.Admins, m.Name)
			}
		}()
	}

	// 渠道信息
	if info.ChannelID > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if channelInfo, err := drainageService.SimpleChannelsDetail(ctx, drainage.ContactChannels{ID: info.ChannelID, EnterpriseID: enterpriseId}); err == nil {
				info.ChannelName = channelInfo.Name
			}
		}()
	}

	// 群主信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		if ownerMembers, err := memberService.GetMemberByMemberIds(ctx, &members.BatchGetMemberByMemberId{EnterpriseID: enterpriseId, IDs: []interface{}{info.OwnerMemberID}}); err == nil && len(*ownerMembers) > 0 {
			info.Owner = (*ownerMembers)[0].Name
			info.UserID = (*ownerMembers)[0].Userid
		}
	}()

	// 群类型信息
	if info.ChatTypeID > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if typeInfo, _ := groupChatService.GetGroupChatTypeInfo(ctx, socialDraw.WorkGroupChatType{
				ID: info.ChatTypeID, EnterpriseID: enterpriseId,
			}); typeInfo != nil {
				info.ChatTypeName = typeInfo.Name
			}

		}()
	}
	wg.Wait()

	return
}

// workGroupChatSetting update settings of work group chat
func workGroupChatSetting(c *gin.Context) {
	var req socialDraw.WorkGroupChatSettingRequest
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}
	if req.ChatID == "" {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if req.ChannelID > 0 {
		channelInfo, _ := drainageService.ChannelDetail(ctx.Wrap(c), &drainage.ContactChannels{
			ID:           req.ChannelID,
			EnterpriseID: req.EnterpriseID,
		})
		if channelInfo == nil {
			ginplus.ResError(c, errs.ParamsErr)
			return
		}
		req.ChannelName = channelInfo.Name
	}
	if req.ChatTypeID > 0 {
		typeInfo, _ := groupChatService.GetGroupChatTypeInfo(ctx.Wrap(c), socialDraw.WorkGroupChatType{
			ID:           req.ChatTypeID,
			EnterpriseID: req.EnterpriseID,
		})
		if typeInfo == nil {
			ginplus.ResError(c, errs.ParamsErr)
			return
		}
	}
	if req.OperatorID != "" {
		memberInfo, _ := memberService.GetMemberByUserId(ctx.Wrap(c), req.EnterpriseID, req.OperatorID)
		if memberInfo == nil {
			ginplus.ResError(c, errs.ParamsErr)
			return
		}
		req.Operator = memberInfo.Name
		req.DepartmentID = memberInfo.Department
		dByte, _ := json.Marshal(req.DepartmentID)
		req.DepartmentIds = string(dByte)

		if len(memberInfo.Department) > 0 {
			departments, _ := departmentService.GetDepartmentsByIds(ctx.Wrap(c), req.EnterpriseID, memberInfo.Department)
			for _, v := range departments {
				req.Departments = append(req.Departments, v.Name)
			}
			dpByte, _ := json.Marshal(req.Departments)
			req.DepartmentStr = string(dpByte)
		}
	}
	err := groupChatService.WorkGroupChatSetting(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

// workGroupChatBatchSetting 社群批量设置
func workGroupChatBatchSetting(c *gin.Context) {
	var req socialDraw.WorkGroupChatBatchSettingRequest
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	settingData, _ := redis.Client().Get(ctx.Wrap(c), req.SettingKey).Bytes()
	if len(settingData) == 0 {
		ginplus.ResError(c, errs.GroupChatImportAgain)
		return
	}
	var settings []*socialDraw.WorkGroupChatImportSetting
	if err := json.Unmarshal(settingData, &settings); err != nil {
		ginplus.ResError(c, errs.GroupChatImportAgain)
		return
	}

	for _, v := range settings {
		if err := groupChatService.WorkGroupChatSetting(ctx.Wrap(c), v.WorkGroupChatSettingRequest); err != nil {
			log.Error(ctx.Wrap(c)).Msgf("更新群【%s】设置失败, %+v", v.ChatID, err)
		}
	}

	ginplus.ResSuccess(c, nil)
	return
}

// workChatGroupTypeList 获取群类型列表
func workChatGroupTypeList(c *gin.Context) {
	var req socialDraw.WorkGroupChatType
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}

	types, err := groupChatService.GetGroupChatTypeList(ctx.Wrap(c), req, true)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, types)
	return
}

// workChatGroupTypeManage 群类型管理
func workChatGroupTypeManage(c *gin.Context) {
	var req socialDraw.WorkGroupChatType
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}
	// 检测数据类型是否重复设置
	existInfo, _ := groupChatService.GetGroupChatTypeInfo(ctx.Wrap(c), req)
	if req.Opt != 3 {
		if existInfo != nil {
			ginplus.ResError(c, errs.GroupChatTypeExisted)
			return
		}
	}
	// 删除时有被使用，不允许删除
	if req.Opt == 3 {
		used, _ := groupChatService.CheckWorkGroupChatTypeUsed(ctx.Wrap(c), req.EnterpriseID, req.ID)
		if used {
			ginplus.ResError(c, errs.GroupChatTypeNotDeleted)
			return
		}
	}
	typeInfo, err := groupChatService.ManageGroupChatType(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, typeInfo)
	return
}

// workGroupChatInfo h5社群详情
func workH5GroupChatInfo(c *gin.Context) {
	var req socialDraw.GroupInfoReq
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}
	if req.ChatID == "" {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	cc := ctx.Wrap(c)
	resp := &socialDraw.WorkGroupChatInfo{}

	chatInfo, err := chatService.GetGroupInfo(cc, chat.WorkGroupChat{ChatID: req.ChatID, EnterpriseID: req.EnterpriseID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	chatRes := chat.WorkGroupChatBase{WorkGroupChat: *chatInfo}

	getWorkGroupChatExtInfo(cc, req.EnterpriseID, &chatRes)

	ntime := time.Now().AddDate(0, 0, -30)
	monthSt := time.Date(ntime.Year(), ntime.Month(), ntime.Day(), 0, 0, 0, 0, ntime.Location()).Unix()

	req.MonthTS = monthSt
	req.CID = chatRes.ID
	resp.WorkGroupChatBase = chatRes

	var data socialDraw.WorkGroupChatStatisticInfo

	wg := threading.NewRoutineGroup()

	// 查询群累计人数
	wg.RunSafe(func() {
		var count int
		gormx.GetCkInstance().WithContext(c).
			Table(fmt.Sprintf("%s.%s", "micro_scrm_ck", "work_group_chat_member")).
			Distinct("count(DISTINCT userid) as count").
			Where("group_chat_id", chatInfo.ID).
			Where("enterprise_id = ?", req.EnterpriseID).
			Where("status = 1").
			Where("deleted_at = 0").
			Group("group_chat_id").
			Scan(&count)

		data.ValidUv = count
	})

	// 群新增人数
	wg.RunSafe(func() {
		var count int
		gormx.GetCkInstance().WithContext(c).
			Table(fmt.Sprintf("%s.%s", "micro_scrm_ck", "work_group_chat_member")).
			Distinct("count(DISTINCT userid) as count").
			Where("group_chat_id", chatInfo.ID).
			Where("join_time >= ? and join_time <= ?", req.STS, req.ETS).
			Where("enterprise_id = ?", req.EnterpriseID).
			Where("status = 1").
			Where("deleted_at = 0").
			Group("group_chat_id").
			Scan(&count)

		data.AddUv = count
	})

	// 退群人数
	wg.RunSafe(func() {
		var count int
		// 获取群成员id集合
		gormx.GetCkInstance().WithContext(c).
			Table(fmt.Sprintf("%s.%s", "micro_scrm_ck", "work_group_chat_member")).
			Distinct("count(DISTINCT userid) as count").
			Where("group_chat_id", chatInfo.ID).
			Where("enterprise_id = ?", req.EnterpriseID).
			Where("status = 2").
			Where("deleted_at = 0").
			Group("group_chat_id").
			Scan(&count)

		data.AllLossUv = count
	})

	// 会员人数
	wg.RunSafe(func() {
		var count int64

		gormx.GetCkInstance().WithContext(c).
			Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseID, "wx_user_corp_external")).
			Where("union_code <> '' or customer_code <> ''").
			Where("openid in (?)", func() *gorm.DB {
				return gormx.GetCkInstance().WithContext(c).
					Table(fmt.Sprintf("%s.%s", "micro_scrm_ck", "work_group_chat_member")).
					Distinct("userid").
					Where("group_chat_id", chatInfo.ID).
					Where("enterprise_id = ?", req.EnterpriseID).
					Where("status = 1").
					Where("deleted_at = 0")
			}).
			Count(&count)

		data.MemberUv = int(count)
	})

	wg.Wait()

	// 非会员人数
	data.GeneralUv = data.ValidUv - data.MemberUv

	// 退群率
	if data.AllLossUv == 0 {
		data.LossUvFee = "0%"
	} else {
		data.LossUvFee = mathrounded(float64(data.AllLossUv), float64(data.ValidUv))
	}

	// 会员率
	if data.MemberUv == 0 {
		data.MemberCardUv = "0%"
	} else {
		data.MemberCardUv = mathrounded(float64(data.MemberUv), float64(data.MemberUv+data.GeneralUv))
	}

	resp.WorkGroupChatStatisticInfo = data

	nowDate := time.Now().Format("2006-01-02")
	theTime := nowDate + " 04:00:00"
	resp.UpdateTime = theTime
	ginplus.ResSuccess(c, resp)
	return
}
func mathrounded(params1, params2 float64) string {
	if params1 == 0 {
		return "0%"
	}
	params3 := helper.MathRounded((params1/params2)*100, 2)
	if params3 == float64(int64(params3)) {
		return fmt.Sprintf("%.0f", params3) + "%"
	}
	return fmt.Sprintf("%.2f", params3) + "%"
}

// workGroupChatFilters 自定义筛选列表
func workGroupChatFilters(c *gin.Context) {
	var req socialDraw.WorkGroupFilterRequest
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.FiledName == "" {
		req.FiledName = "channels"
	}
	data, err := groupChatService.GetWorkGroupFilters(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
	return
}

// getWorkGroupChatMembers 获取社群成员明细
func getWorkGroupChatMembers(c *gin.Context) {
	var (
		req  socialDraw.WorkGroupChatMembersRequest
		resp socialDraw.PaginationResponse
	)
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	req.HandleOffset()

	data, totalCount, err := groupChatService.GetWorkGroupChatMembers(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var (
		wg             sync.WaitGroup
		userIds        []string
		externalUsers  = make(map[string]*socialDraw.WxUserCorpExternal)
		invitorMap     = make(map[string]string)
		relatedMembers = make(map[string][]string) // openid -> []memberName
	)
	for _, v := range data {
		invitorMap[v.UserID] = v.Name
		userIds = append(userIds, v.UserID)
	}

	// 邀请人信息
	for _, v := range data {
		if invitor, ok := invitorMap[v.InviteUserID]; ok {
			v.InviteUser = invitor
		}
	}

	// 会员信息 + 已添加导购
	wg.Add(1)
	go func() {
		defer wg.Done()
		externalUsers, _ = groupChatService.GetWorkGroupExternalUsers(ctx.Wrap(c), req.EnterpriseID, userIds)
		if req.MemberType > socialDraw.ChatMemberAll {
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		externalUserRelatedMembers, _ := getExternalUserRelatedMembers(ctx.Wrap(c), req.EnterpriseID, userIds)
		for openId, mems := range externalUserRelatedMembers {
			for _, m := range mems {
				relatedMembers[openId] = append(relatedMembers[openId], m.Name)
			}
		}
	}()

	wg.Wait()

	for _, v := range data {
		if extUser, ok := externalUsers[v.UserID]; ok {
			if extUser.UnionCode != "" {
				v.VipStatus = 1
				v.VipLevel = extUser.Grade
				v.VipCode = extUser.UnionCode
				v.Mobile = extUser.Mobile
			}
		}
		if guiders, ok := relatedMembers[v.UserID]; ok {
			v.Guiders = guiders
		}
	}

	resp.Items = data
	resp.TotalCount = totalCount

	ginplus.ResSuccess(c, resp)
	return
}

// getExternalUserRelatedMembers 获取外部联系人关联的员工信息
func getExternalUserRelatedMembers(ctx ctx.Context, enterpriseId int64, externalUserIds []string) (map[string][]*members.Member, error) {
	externalUsers, err := groupChatService.GetWorkGroupExternalUsers(ctx, enterpriseId, externalUserIds) // openID -> corpInfo
	if err != nil {
		return nil, err
	}
	var (
		corpExternalId []int64
		memberParams   = &members.BatchGetMemberByMemberId{
			EnterpriseID: enterpriseId,
			DeletedAt:    0,
		}
		memberMap       = make(map[int64]members.Member)                 // memberID -> memberInfo
		corpExternalMap = make(map[int64]*socialDraw.WxUserCorpExternal) // corpID -> corpInfo
	)
	for _, eu := range externalUsers {
		corpExternalMap[eu.ID] = eu
		corpExternalId = append(corpExternalId, eu.ID)
	}
	relatedMembers, err := memberService.GetExternalUserRelatedMembers(ctx, enterpriseId, corpExternalId) // corpID -> []memberID
	if err != nil {
		return nil, err
	}
	for _, memberIds := range relatedMembers {
		for _, memberId := range memberIds {
			memberParams.IDs = append(memberParams.IDs, memberId)
		}
	}
	mems, err := memberService.GetMemberByMemberIds(ctx, memberParams) // members
	if err != nil {
		return nil, err
	}
	for _, member := range *mems {
		memberMap[member.ID] = member
	}

	var (
		externalIdRelatedMembers     = make(map[int64][]*members.Member)  // corpID -> []members
		externalOpenIdRelatedMembers = make(map[string][]*members.Member) // corpOpenID -> []members
	)
	for corpId, memberIds := range relatedMembers {
		for _, memId := range memberIds {
			if memInfo, ok := memberMap[memId]; ok {
				externalIdRelatedMembers[corpId] = append(externalIdRelatedMembers[corpId], &memInfo)
			}
		}
	}
	for corpId, mems := range externalIdRelatedMembers {
		if corpInfo, ok := corpExternalMap[corpId]; ok {
			externalOpenIdRelatedMembers[corpInfo.OpenID] = mems
		}
	}
	return externalOpenIdRelatedMembers, nil
}

// getWorkGroupChatDataGrid 获取统计数据
func getWorkGroupChatDataGrid(c *gin.Context) {
	var (
		req  socialDraw.GroupInfoReq
		resp socialDraw.WorkGroupChatDataGridResponse
	)
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EnterpriseID <= 0 {
		ginplus.ResError(c, errs.EnterpriseIdEmpty)
		return
	}
	if req.ChatID == "" {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	baseChat, err := chatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{ChatID: req.ChatID, EnterpriseID: req.EnterpriseID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	chatInfo := &chat.WorkGroupChatBase{WorkGroupChat: *baseChat}
	req.CID = chatInfo.ID

	var (
		wg      sync.WaitGroup
		statics = &socialDraw.WorkGroupChatStatistic{}
		charts  = &socialDraw.DataGridChart{}
	)

	// 统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		if data, err := groupChatService.GroupUserUvStatistics(ctx.Wrap(c), req); err == nil {
			statics.AddUv = data.AddUv
			statics.AllLossUv = data.AllLossUv
			statics.ValidUv = data.ValidUv
			statics.NewLossUv = data.NewLossUv
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if data, err := groupChatService.GroupMemberStatistics(ctx.Wrap(c), req); err == nil {
			statics.MemberUv = data.MemberUv
			statics.GeneralUv = data.GeneralUv
			statics.WhiteGlodDown = data.WhiteGlodDown
		}
	}()

	// 饼图
	wg.Add(1)
	go func() {
		defer wg.Done()
		if data, err := groupChatService.NewGroupUserChat(ctx.Wrap(c), req); err == nil {
			charts.New = data
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if data, err := groupChatService.ValidGroupUserChat(ctx.Wrap(c), req); err == nil {
			charts.Valid = data
		}
	}()

	wg.Wait()

	resp.Target = statics
	resp.Chart = charts

	ginplus.ResSuccess(c, resp)
	return
}

// getMemberList
func getMemberList(c *gin.Context) {
	var (
		req  socialDraw.WorkGroupChatMembersRequest
		resp socialDraw.PaginationResponse
	)
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.ChatID == "" {
		ginplus.ResError(c, fmt.Errorf("缺少参数"))
		return
	}
	req.HandleOffset()
	baseChat, err := chatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{ChatID: req.ChatID, EnterpriseID: req.EnterpriseID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	chatInfo := &chat.WorkGroupChatBase{WorkGroupChat: *baseChat}
	req.CID = chatInfo.ID
	req.MarsID, _ = groupChatService.GetMarsID(ctx.Wrap(c), req.EnterpriseID, req.UserId)
	if req.MarsID == "" && req.ListType == 0 {
		ginplus.ResSuccess(c, resp)
		return
	}
	data, total, err := groupChatService.GetWorkGroupMembersList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp.TotalCount = total
	resp.Items = data
	ginplus.ResSuccess(c, resp)
	return
}

// fileUploadExchangeData 上传文件交换上传数据
func fileUploadExchangeData(c *gin.Context) {
	enterpriseIdStr := c.PostForm("enterprise_id")
	enterpriseID, _ := strconv.ParseInt(enterpriseIdStr, 10, 64)
	fh, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var (
		ft               string
		importedSettings []*socialDraw.WorkGroupChatImportSetting
	)
	fileData, err := common.ReadExcelFile(fh, &ft)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if ft == ".xlsx" {
		importedSettings = readXlsxData(fileData.(*xlsx.File).Sheets)
	} else {
		importedSettings = readXlsData(*fileData.(*xls.Workbook))
	}

	resp := readGroupChatSettings(ctx.Wrap(c), enterpriseID, importedSettings)
	cacheKey := fmt.Sprintf("work_group_setting:%d", time.Now().Unix())
	if resp.Success > 0 {
		cacheBody, _ := json.Marshal(resp.SuccessItems)
		redis.Client().SetEX(ctx.Wrap(c), cacheKey, string(cacheBody), time.Minute*10)
		resp.SettingKey = cacheKey
	}

	ginplus.ResSuccess(c, resp)
	return
}

func readXlsxData(sheets []*xlsx.Sheet) []*socialDraw.WorkGroupChatImportSetting {
	var settings []*socialDraw.WorkGroupChatImportSetting
	for index, sheet := range sheets {
		if index > 0 {
			continue
		}
		for ri, row := range sheet.Rows {
			if ri < 2 {
				continue
			}
			if len(row.Cells) == 0 {
				continue
			}
			if row.Cells[0].Value == "" {
				continue
			}
			var setting = &socialDraw.WorkGroupChatImportSetting{}
			setting.LineNum = ri + 1
			setting.ChatID = row.Cells[0].String() // 群主键ID
			if len(row.Cells) < 7 {
				setting.Errors = append(setting.Errors, "请按模板格式上传数据")
				settings = append(settings, setting)
				continue
			}
			setting.ChatName = row.Cells[1].String()
			setting.ChannelName = row.Cells[2].String()
			setting.ChatType = row.Cells[3].String()
			setting.Operator = row.Cells[4].String()
			setting.Mobile = row.Cells[5].String()
			setting.ChatTarget = row.Cells[6].String()
			settings = append(settings, setting)
		}
	}
	return settings
}

func readXlsData(wb xls.Workbook) []*socialDraw.WorkGroupChatImportSetting {
	var settings []*socialDraw.WorkGroupChatImportSetting
	for index, sheet := range wb.GetSheets() {
		if index > 0 {
			continue
		}
		for ri, row := range sheet.GetRows() {
			if ri < 2 {
				continue
			}
			if len(row.GetCols()) == 0 {
				continue
			}
			if cell, err := row.GetCol(0); err != nil || cell.GetString() == "" {
				continue
			}
			getColStr := func(c structure.CellData, err error) string {
				if err != nil {
					return ""
				}
				return c.GetString()
			}
			var setting = &socialDraw.WorkGroupChatImportSetting{}
			setting.LineNum = ri + 1
			setting.ChatID = getColStr(row.GetCol(0)) // 群主键ID
			if len(row.GetCols()) < 7 {
				setting.Errors = append(setting.Errors, "请按模板格式上传数据")
				settings = append(settings, setting)
				continue
			}
			setting.ChatName = getColStr(row.GetCol(1))
			setting.ChannelName = getColStr(row.GetCol(2))
			setting.ChatType = getColStr(row.GetCol(3))
			setting.Operator = getColStr(row.GetCol(4))
			setting.Mobile = getColStr(row.GetCol(5))
			setting.ChatTarget = getColStr(row.GetCol(6))
			settings = append(settings, setting)
		}
	}
	return settings
}

// readGroupChatSettings 读取数据
func readGroupChatSettings(ctx ctx.Context, enterpriseId int64, totalSettings []*socialDraw.WorkGroupChatImportSetting) *socialDraw.BatchImportResponse {
	var (
		chatIds      []string
		channelNames []string
		chatTypes    []string
		mobiles      []string
		repeatNum    int
		settings     []*socialDraw.WorkGroupChatImportSetting
		settingMap   = make(map[string]string)
	)
	for i := 0; i <= len(totalSettings)-1; i++ {
		setting := totalSettings[i]
		if _, ok := settingMap[setting.ChatID]; ok {
			repeatNum++
			continue
		}
		settingMap[setting.ChatID] = setting.ChatID

		chatIds = append(chatIds, setting.ChatID)
		channelNames = append(channelNames, setting.ChannelName)
		chatTypes = append(chatTypes, setting.ChatType)
		mobiles = append(mobiles, setting.Mobile)

		settings = append(settings, setting)
	}

	// 获取校验数据
	var (
		wg          sync.WaitGroup
		chatMap     = make(map[string]string)
		typeMap     = make(map[string]int64)
		channelMap  = make(map[string]int64)
		operatorMap = make(map[string]*members.Member)
		dpMap       = make(map[int64]string)
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatList, _ := groupChatService.GetGroupChatList(ctx, socialDraw.WorkGroupListRequest{
			EnterpriseID: enterpriseId,
			ChatIds:      chatIds,
			StartAt:      0,
			EndAt:        time.Now().Unix() + 86400,
		})
		for _, v := range chatList {
			chatMap[v.ChatID] = v.ChatID
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		types, _ := groupChatService.GetGroupChatTypeList(ctx, socialDraw.WorkGroupChatType{
			EnterpriseID: enterpriseId,
			Names:        chatTypes,
		}, false)
		for _, v := range types {
			typeMap[v.Name] = v.ID
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, channels, _ := drainageService.SimpleChannelsList(ctx, drainage.ContactChannelListRequest{
			EnterpriseID: enterpriseId,
			ChannelNames: channelNames,
			Page:         1,
			Size:         10000,
		})
		for _, v := range *channels {
			channelMap[v.Name] = v.ID
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var dpIds []int64
		operators, _ := memberService.GetMembersByUserIdsOrMobiles(ctx, enterpriseId, members.MemberByMobile, mobiles)
		for _, v := range operators {
			dpIds = append(dpIds, v.Department...)
			operatorMap[v.Mobile] = v
		}
		dps, _ := departmentService.GetDepartmentsByIds(ctx, enterpriseId, dpIds)
		for _, v := range dps {
			dpMap[v.ID] = v.Name
		}
	}()

	wg.Wait()

	// 检测数据
	var (
		success, errors []*socialDraw.WorkGroupChatImportSetting
		resp            = &socialDraw.BatchImportResponse{}
	)
	for _, v := range settings {
		if len(v.Errors) > 0 { // 模板格式问题，不再继续判断其他错误
			errors = append(errors, v)
			continue
		}
		if _, ok := chatMap[v.ChatID]; !ok {
			v.Errors = append(v.Errors, "未查询到群ID")
		}
		if _, ok := channelMap[v.ChannelName]; !ok && v.ChannelName != "" {
			v.Errors = append(v.Errors, "所属渠道不存在")
		}
		if _, ok := typeMap[v.ChatType]; !ok && v.ChatType != "" {
			v.Errors = append(v.Errors, "群类型不存在")
		}
		if _, ok := operatorMap[v.Mobile]; !ok {
			v.Errors = append(v.Errors, "群运营者手机号不匹配")
		}
		if len([]rune(v.ChatTarget)) > 50 {
			v.Errors = append(v.Errors, "群目标字数超过50字符")
		}
		if len(v.Errors) > 0 {
			errors = append(errors, v)
			continue
		}
		v.ChannelID = channelMap[v.ChannelName]
		v.ChatTypeID = typeMap[v.ChatType]
		v.EnterpriseID = enterpriseId
		if _, ok := operatorMap[v.Mobile]; ok {
			v.Operator = operatorMap[v.Mobile].Name
			v.OperatorID = operatorMap[v.Mobile].Userid
			for _, dpId := range operatorMap[v.Mobile].Department {
				v.Departments = append(v.Departments, dpMap[dpId])
			}
			dpBody, _ := json.Marshal(v.Departments)
			v.DepartmentStr = string(dpBody)
		}
		success = append(success, v)
	}
	resp.Success = len(success)
	resp.Errors = len(errors)
	resp.Repeated = repeatNum

	resp.SuccessItems = success
	resp.ErrorItems = errors

	return resp
}

// getGroupChatMembersInfo 获取社群成员明细-h5
func getGroupChatMembersInfo(c *gin.Context) {
	var (
		req  socialDraw.WorkGroupChatMembersRequest
		resp socialDraw.PaginationResponse
	)
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	req.HandleOffset()
	baseChat, err := chatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{ChatID: req.ChatID, EnterpriseID: req.EnterpriseID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	chat := &chat.WorkGroupChatBase{WorkGroupChat: *baseChat}
	req.CID = chat.ID
	data, totalCount, err := groupChatService.GetWorkGroupChatInfo(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var (
		userIds    []string
		membersMap = make(map[string]*members.Member)
	)

	for _, v := range data {
		v.Time = time.Unix(v.TmpTime, 0).Format("2006-01-02 15:04")
		if v.MemberID != 0 {
			userIds = append(userIds, v.UserID)
		}
	}
	if len(userIds) > 0 {
		mems, _ := memberService.GetMembersByUserIdsOrMobiles(ctx.Wrap(c), req.EnterpriseID, members.MemberByUserID, userIds)
		for _, m := range mems {
			membersMap[m.Userid] = m
		}
		for _, v := range data {
			if member, ok := membersMap[v.UserID]; ok {
				v.Avatar = member.Avatar
			}
		}
	}
	resp.Items = data
	resp.TotalCount = totalCount
	ginplus.ResSuccess(c, resp)
	return
}

// getLastDataUpdateTime 最后一次数据更新时间
func getLastDataUpdateTime(c *gin.Context) {
	var req struct {
		EnterpriseID int64 `json:"enterprise_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	timestamp := groupChatService.GetLastUpdateTime(ctx.Wrap(c), req.EnterpriseID)

	n := time.Now()
	clockEight := time.Date(n.Year(), n.Month(), n.Day(), 8, 0, 0, 0, time.Local)
	if n.After(clockEight) { // 当天8点之后，有更新展示更新时间，无更新展示定时任务时间
		if timestamp > clockEight.Unix() {
			ginplus.ResSuccess(c, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))
			return
		}
		ginplus.ResSuccess(c, n.Format("2006-01-02 08:00:00"))
		return
	}
	// 当天8点之前，昨日8点之后有更新，展示更新时间，否则展示昨日定时任务时间
	if timestamp > clockEight.Unix()-86400 {
		ginplus.ResSuccess(c, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))
		return
	}
	ginplus.ResSuccess(c, time.Unix(clockEight.Unix()-86400, 0).Format("2006-01-02 15:04:05"))
	return
}

// BatchSyncGroupChats 批量同步群信息
func BatchSyncGroupChats(c *gin.Context) {
	var req socialDraw.WorkGroupListRequest
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if req.EndAt <= 0 {
		req.EndAt = time.Now().Unix()
	}
	if syncAt := groupChatService.GetWorkGroupSyncFlag(ctx.Wrap(c), req.EnterpriseID); syncAt > 0 {
		ginplus.ResSuccess(c, "正在同步群信息，请稍后刷新查询")
		return
	}

	groupChatService.SetLastUpdateTime(ctx.Wrap(c), req.EnterpriseID, time.Now().Unix())
	groupChats, err := groupChatService.GetGroupChatList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var corpDomain []string
	if err := db.DB().Select(ctx.Wrap(c), &corpDomain, `select domain from enterprise_corp where enterprise_id=?`, req.EnterpriseID); err != nil {
		ginplus.ResError(c, err)
		return
	}

	go func() {
		groupChatService.SetOrClearWorkGroupSyncFlag(ctx.Wrap(c), req.EnterpriseID, 0)
		for _, v := range groupChats {
			log.Info(ctx.Wrap(c)).Msgf("--------------- sync chat [%s] --------------------", v.ChatID)
			if err := manualSyncGroupChatInfo(ctx.Wrap(c), req.EnterpriseID, corpDomain[0], v.ChatID, new(work.QyWorkImpl)); err != nil {
				log.Error(ctx.Wrap(c)).Msgf("update work group chat [%+v] error %+v", v.ChatID, err)
				continue
			}
		}
		groupChatService.SetOrClearWorkGroupSyncFlag(ctx.Wrap(c), req.EnterpriseID, 30)
	}()
	ginplus.ResSuccess(c, "正在同步群信息，请稍后刷新查询")
	return
}

func manualSyncGroupChatInfo(ctx ctx.Context, enterpriseId int64, corpDomain, chatId string, workImpl *work.QyWorkImpl) error {
	chatInfo, err := customerService.FindWorkGroupChatDetail(ctx, customer.WorkGroupChat{EnterpriseID: enterpriseId, ChatID: chatId})
	if err != nil {
		return err
	}

	workImpl = workImpl.SetDomain(corpDomain)
	req := work.ExternalGroupRequest{
		ChatID:      chatId,
		NeedName:    1,
		QyGwRequest: work.QyGwRequest{},
	}
	if req.ChatID == "" {
		return nil
	}
	chatDetail, err := workImpl.ExternalGroupChatGet(ctx, req)
	if err != nil {
		return err
	}
	var (
		adminList     []string
		chatMemberNum = len(chatDetail.MemberList)
		workGroupType = customer.WorkGroupTypeDissolve
	)
	if chatMemberNum > 0 {
		workGroupType = customer.WorkGroupTypeExists
	}
	for _, a := range chatDetail.AdminList {
		adminList = append(adminList, a.Userid)
	}
	admin := "[]"
	if len(adminList) > 0 {
		a1, _ := json.Marshal(adminList)
		admin = string(a1)
	}
	// 更新信息
	chatInfo.Name = chatDetail.Name
	chatInfo.Notice = chatDetail.Notice
	chatInfo.UpdatedAt = time.Now().Unix()
	chatInfo.Status = workGroupType
	chatInfo.MemberNum = chatMemberNum
	chatInfo.AdminList = admin
	chatInfo.FormatNotice()
	if workGroupType == customer.WorkGroupTypeDissolve && chatInfo.DissolveTime == 0 {
		chatInfo.DissolveTime = chatInfo.UpdatedAt
	}
	memInfo, _ := memberService.GetMemberByUserId(ctx, enterpriseId, chatDetail.Owner)
	if memInfo != nil {
		chatInfo.OwnerMemberID = memInfo.ID
	}
	if err := customerService.UpdateWorkGroupChat(ctx, *chatInfo); err != nil {
		log.Error(ctx).Msgf("更新客户群信息失败，更新数据:%+v,错误信息：%+v", chatInfo, err)
	}
	if _, err := groupChatMemberService.BatchUpdateExternalWorkGroupChatMembersInfo(ctx, enterpriseId, chatInfo.ID, chatDetail.MemberList); err != nil {
		log.Error(ctx).Msgf("update chat member info error %+v", err)
	}
	return nil
}

// getActivityInfo 获取活动详情
func getActivityInfo(c *gin.Context) {
	var (
		req socialDraw.ActivityReq
	)
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	baseChat, err := chatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{ChatID: req.ChatID, EnterpriseID: req.EnterpriseID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	chat := &chat.WorkGroupChatBase{WorkGroupChat: *baseChat}
	req.Cid = chat.ID
	data, err := groupChatService.GetWorkGroupUserActivityInfo(ctx.Wrap(c), req.EnterpriseID, req.UserID, req.Cid)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
	return
}

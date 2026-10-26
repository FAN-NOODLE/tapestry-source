package sidebarGroupChat

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk"
)

// Service 应用服务
type SidebarGroupChatService struct {
	Repo           Repository
	enterpriseRepo enterprise.Repository
	customerRepo   customer.Repository
}

// GroupChatInfo获取社群侧边栏基础信息
func (s *SidebarGroupChatService) GroupChatInfo(c ctx.Context, request GroupChatInfoRequest) (*GroupChatInfoResponse, error) {

	var groupChatInfoRequest = GroupChatInfoRequest{
		GroupChatId:  request.GroupChatId,
		EnterpriseId: request.EnterpriseId,
	}
	groupChatInfo, groupChatInfoErr := s.Repo.GroupChatInfo(c, groupChatInfoRequest)

	if groupChatInfo == nil {
		return nil, errs.GroupChartNotFound
	}

	if groupChatInfoErr != nil {
		return nil, groupChatInfoErr
	}

	// 社群基础信息
	groupChatInfoResponse := GroupChatInfoResponse{
		Name:          groupChatInfo.Name,
		CreateTime:    groupChatInfo.CreateTime,
		GroupImageUrl: "https://img-crs-test.vchangyi.com/goods16062044997940.png",
	}

	s.GroupChatMemberStatistics(c, groupChatInfoRequest, &groupChatInfoResponse)

	// 获取群主信息
	groupChatLeaderInfo, _ := s.Repo.GroupChatLeaderInfo(c, groupChatInfoRequest)
	groupChatInfoResponse.GroupLeaderName = groupChatLeaderInfo.Name
	//群成员数量
	groupChatMember, _ := s.Repo.GroupChatMemberUnionid(c, groupChatInfoRequest)
	groupChatInfoResponse.GroupMemberNum = int64(len(*groupChatMember))
	uidList := []string{}
	//遍历所有群成员，获取所有群成员unionid
	for _, value := range *groupChatMember {
		uidList = append(uidList, value.Unionid)
	}
	uidStr := strings.Join(uidList, "','")
	groupCorpUser, _ := s.Repo.GroupChatCorpUser(c, uidStr, request.EnterpriseId)
	vipNum := 0
	gradeMap := make(map[string]int)
	unionidMap := make(map[string]int)
	//所有会员等级
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	params := gin.H{
		"enterprise_id": request.EnterpriseId,
	}
	// 查询用户信息
	url := config.CfgData.Service.MicroWork.Domain + "/service/v1/data-lake/query-tier-list"
	resp, err := client.Post(c, url, clients.JSONBody(params))
	if err != nil {
		return nil, err
	}
	var resultData resultData
	_, _ = resp.Read(clients.JSON(&resultData))
	//if err != nil {
	//	return nil, err
	//}
	if resultData.Status == http.StatusOK {
		for _, gradeTex := range resultData.Data {
			gradeMap[gradeTex.TierName] = 0
		}
		//return nil, err
	}

	//遍历所有群成员，获取所有群成员unionid,统计会员数量
	if groupCorpUser != nil {
		for _, value := range *groupCorpUser {
			unionidMap[value.Unionid] = 1
			if value.Grade != "" {
				vipNum++
				_, ok := gradeMap[value.Grade]
				if ok {
					gradeMap[value.Grade] += 1
				} else {
					gradeMap[value.Grade] = 1
				}
			}
		}
	}
	//兼容通过分享二维码进群用户
	if len(uidList) != len(unionidMap) {
		group := threading.NewRoutineGroup()
		for _, v := range uidList {
			if _, exist := unionidMap[v]; !exist {
				group.RunSafe(func() {
					params := gin.H{
						"enterprise_id": request.EnterpriseId,
						"type":          "unionId",
						"value":         v,
					}
					userGrade, err := GetUserGrade(c, params)
					if err != nil {
						log.Error(ctx.Wrap(c)).Msgf("群客户画像---DL查询会员等级（GetUserGrade）返回错误信息，%s", err)
						return
					}
					gradeMap[userGrade] += 1
					vipNum++
				})
			}
		}
		group.Wait()
	}
	groupChatInfoResponse.GroupVipMemberNum = vipNum
	if int64(vipNum) > groupChatInfoResponse.GroupMemberNum {
		groupChatInfoResponse.GroupVipMemberNum = int(groupChatInfoResponse.GroupMemberNum)
	}
	groupChatInfoResponse.GroupVipGroup = gradeMap
	return &groupChatInfoResponse, nil
}

// GroupChatInfo获取社群侧边栏基础信息
func (s *SidebarGroupChatService) GroupChatMemberStatistics(c ctx.Context, groupChatInfoRequest GroupChatInfoRequest, response *GroupChatInfoResponse) bool {
	// 当前时间
	now := time.Now()
	dayTime := now.Format("2006-01-02")

	// 查询社群用户列表
	groupMemberList, _ := s.GroupChatMemberList(ctx.Wrap(c), groupChatInfoRequest)

	for _, value := range *groupMemberList {
		if value.Sex == 1 {
			response.MaleTotal++
		}
		if value.Sex == 2 {
			response.GenderFemaleTotal++
		}
		if value.Sex == int64(0) {
			response.UnknownGenderTotal++
		}
		if value.Status == 1 {
			if value.JoinDataTime == dayTime {
				response.TodayJoinUser++
			}
		}
		if value.Status == 2 {
			if value.LeaveDataTime == dayTime {
				response.TodayLeaveUser++
			}
		}
		if value.TodayMessageTotal != int64(0) {
			response.TodayActiveUser++
		}
		if value.TodayDealOrderTotal != int64(0) {
			response.TodayDealUser++
		}
		response.TodayDealOrder += value.TodayDealOrderTotal
		response.TodayDealOrderAmount += value.TodayDealOrderAmount
		response.TodayRefundOrder += value.TodayRefundOrderTotal
		response.TodayRefundAmount += value.TodayRefundAmount
	}
	return true
}

// GroupChatInfo获取社群侧边栏基础信息
func (s *SidebarGroupChatService) GroupChatMemberList(c ctx.Context, request GroupChatInfoRequest) (*[]WorkGroupChatMember, error) {
	// 获取第三方模板信息
	var groupChatInfoRequest = GroupChatInfoRequest{
		GroupChatId:  request.GroupChatId,
		EnterpriseId: request.EnterpriseId,
	}
	chatMemberList, err := s.Repo.GroupChatMember(c, groupChatInfoRequest)
	if err != nil {
		return nil, err
	}
	if chatMemberList == nil {
		return nil, errs.GroupChartNotFound
	}

	var dataMap []WorkGroupChatMember
	for _, value := range chatMemberList {
		user, err := s.customerRepo.QueryEnterpriseUser(c, request.EnterpriseId, &customer.SourceUserBase{UnionID: value.Unionid, CustomerType: customer.CustomerTypeExternal})
		if err == nil {
			if user != nil {
				nick := user.Nick
				if nick == "" {
					nick = user.Name
				}
				value.NickName = nick
				value.Avatar = user.Avatar
				value.Sex = int64(user.Sex)
				for _, eu := range user.WxExternal {
					value.ExternalUserId = eu.ThirdID
					break
				}
				for _, eu := range user.WxShop {
					value.ShopUserId = eu.ID
					break
				}
			}
		}
		/** 群用户消息内容 **/
		MessageTotal, TodayMessageTotal, LastChatAt, err := s.GroupChatMessage(ctx.Wrap(c), *value)
		if err != nil {
			return nil, err
		}
		// 绑定聊天总数
		value.MessageTotal = MessageTotal
		// 绑定今日聊天总数
		value.TodayMessageTotal = TodayMessageTotal
		// 绑定最新发送时间
		value.LastChatAt = LastChatAt

		dataMap = append(dataMap, *value)
	}
	return &dataMap, nil
}

// GroupChatMessage 群用户消息统计
func (s *SidebarGroupChatService) GroupChatMessage(c ctx.Context, value WorkGroupChatMember) (int64, int64, int64, error) {
	// 用户聊天总数
	var MessageTotal int64 = 0
	// 今日用户聊天总数
	var TodayMessageTotal int64 = 0
	// 用户最后聊天时间
	var LastChatAt int64 = 0
	// 当前时间
	now := time.Now()
	dayTime := now.Format("2006-01-02")

	// 查询用户聊天记录
	var chatMsgRequest = ChatMsgRequest{
		RoomId:   value.ChatId,
		FromUser: value.ExternalUserId,
	}
	ChatMemberMsgList, ChatMemberMsgErr := s.Repo.GroupChatMemberMsg(c, chatMsgRequest)
	if ChatMemberMsgErr != nil {
		return 0, 0, 0, ChatMemberMsgErr
	}

	for key, msgInfo := range *ChatMemberMsgList {
		// 用户聊天总数
		if msgInfo.FromUser == value.ExternalUserId {
			MessageTotal++
		}
		// 用户今日聊天数量
		if msgInfo.SendTime == dayTime {
			TodayMessageTotal++
		}
		if key == 0 {
			LastChatAt = msgInfo.MsgAt
		}
	}
	return MessageTotal, TodayMessageTotal, LastChatAt, nil
}

// TodayMessageMemberList 群成员活跃排行
type TodayMessageMemberList []WorkGroupChatMember

func (t TodayMessageMemberList) Len() int { return len(t) }

func (t TodayMessageMemberList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TodayMessageMemberList) Less(i, j int) bool {

	return t[i].TodayMessageTotal < t[j].TodayMessageTotal
}

// TodayDealOrderMemberList 群成员订单排行
type TodayDealOrderMemberList []WorkGroupChatMember

func (t TodayDealOrderMemberList) Len() int { return len(t) }

func (t TodayDealOrderMemberList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TodayDealOrderMemberList) Less(i, j int) bool {

	return t[i].TodayDealOrderTotal < t[j].TodayDealOrderTotal
}

// GroupChatUserList 群用户消息统计
func (s *SidebarGroupChatService) GroupChatUserList(c ctx.Context, groupChatInfoRequest GroupChatInfoRequest) (GroupChatUserListResponse, error) {
	// 查询社群用户基础信息
	groupMemberList, _ := s.GroupChatMemberList(ctx.Wrap(c), groupChatInfoRequest)

	if groupChatInfoRequest.Type > 4 || groupChatInfoRequest.Type < 1 {
		groupChatInfoRequest.Type = 1
	}

	var list []WorkGroupChatMember
	if groupChatInfoRequest.Type == 1 {
		// 社群成员列表
		for _, value := range *groupMemberList {
			if value.Status == 1 {
				list = append(list, value)
			}
		}
	}
	if groupChatInfoRequest.Type == 2 {
		// 活跃成员列表
		for _, value := range *groupMemberList {
			if value.MessageTotal != 0 {
				list = append(list, value)
			}
		}
		sort.Sort(sort.Reverse(TodayMessageMemberList(list)))
	}

	if groupChatInfoRequest.Type == 3 {
		// 成交成员列表
		for _, value := range *groupMemberList {
			if value.DealOrderTotal != 0 {
				list = append(list, value)
			}
		}
		sort.Sort(sort.Reverse(TodayDealOrderMemberList(list)))
	}
	if groupChatInfoRequest.Type == 4 {
		list = s.toYX(groupMemberList, groupChatInfoRequest)
	}
	//羽西业务
	//if isCan, authErr := auth.BrandAuth(c, auth.BrandAuthConfig{
	//	EnterpriseID: groupChatInfoRequest.EnterpriseId,
	//	Business:     auth.BUSINESS_USER_JOIN_FISSION_EVENT,
	//}); isCan && authErr == nil {
	//
	//}
	//end

	total := int64(len(list))
	size := groupChatInfoRequest.Size
	page := groupChatInfoRequest.Page
	min := size * (page - 1)
	max := page * size
	if max > total {
		max = total
	}
	if total == 1 {
		min = 0
	}

	var cmp GroupChatUserListResponse
	cmp.Total = total
	cmp.Page = groupChatInfoRequest.Page
	cmp.Size = groupChatInfoRequest.Size
	if list == nil {
		cmp.List = make([]WorkGroupChatMember, 0)
		return cmp, nil
	}
	cmp.List = list[min:max]
	return cmp, nil
}

//toYX 兼容羽西新增群成员参与活动排行榜
func (s *SidebarGroupChatService) toYX(groupMemberList *[]WorkGroupChatMember, groupChatInfoRequest GroupChatInfoRequest) (list []WorkGroupChatMember) {
	unionIds := make([]string, 0)
	exUserIds := make(map[string]string)
	for _, value := range *groupMemberList {
		if value.Status != 1 || value.Unionid == "" {
			continue
		}
		unionIds = append(unionIds, value.Unionid)
		exUserIds[value.Unionid] = value.ExternalUserId
	}
	var userIds []string
	f := gormx.GetInstance()
	f.Select("id").Where("union_id in ?", unionIds).Model(&fission.FissionMpUserModel{}).Pluck("id", &userIds)
	f.Raw("SELECT u.user_id,mp.union_id as `unionid`,mp.nick_name,mp.img_url as `avatar`,count( u.user_id ) AS join_num "+
		"FROM fission_user AS u LEFT JOIN fission_mp_user AS mp ON u.user_id = mp.id"+
		" WHERE u.user_id IN ? and u.enterprise_id=?"+
		" GROUP BY u.user_id ORDER BY join_num DESC", userIds, groupChatInfoRequest.EnterpriseId).Scan(&list)
	for key, _ := range list {
		list[key].ExternalUserId = exUserIds[list[key].Unionid]
	}
	return
}

// GroupUsersAddEmployees 用户添加员工列表
func (s *SidebarGroupChatService) GroupUsersAddEmployees(c ctx.Context, groupUsersAddEmployeesRequest GroupUsersAddEmployeesRequest) (GroupUsersAddEmployeesResponse, error) {
	appid, err := s.enterpriseRepo.QueryOneApp(c, groupUsersAddEmployeesRequest.EnterpriseId, customer.CustomerTypeExternal)
	if err != nil {
		return GroupUsersAddEmployeesResponse{}, err
	}

	user, err := s.customerRepo.QueryEnterpriseUser(c, groupUsersAddEmployeesRequest.EnterpriseId, &customer.SourceUserBase{
		AppID:        appid,
		OpenID:       groupUsersAddEmployeesRequest.ExternalUserId,
		CustomerType: customer.CustomerTypeExternal,
	})
	if err != nil {
		return GroupUsersAddEmployeesResponse{}, err
	}

	eu, ok := user.WxExternal[appid]
	if !ok {
		return GroupUsersAddEmployeesResponse{}, nil
	}

	// 查询群用户订单相关
	employeesList, _ := s.Repo.UsersAddEmployees(c, groupUsersAddEmployeesRequest.EnterpriseId, eu.ID, user.ID)

	var list []GroupUsersAddEmployees
	for _, value := range *employeesList {
		if value.CustomerRelationGuideId == 0 {
			value.IsBindingMember = 1
		} else {
			value.IsBindingMember = 2
		}
		list = append(list, value)
	}

	total := int64(len(list))
	size := groupUsersAddEmployeesRequest.Size
	page := groupUsersAddEmployeesRequest.Page
	min := size * (page - 1)
	max := page * size
	if max > total {
		max = total
	}
	if total == 1 {
		min = 0
	}

	var cmp GroupUsersAddEmployeesResponse
	cmp.Total = total
	cmp.Page = groupUsersAddEmployeesRequest.Page
	cmp.Size = groupUsersAddEmployeesRequest.Size
	if list == nil {
		cmp.List = make([]GroupUsersAddEmployees, 0)
		return cmp, nil
	}
	cmp.List = list[min:max]
	return cmp, nil
}

// GroupUsersAddEmployees 用户参加裂变活动列表
func (s *SidebarGroupChatService) GroupUsersJoinActivity(c ctx.Context, groupUsersJoinActivityRequest GroupUsersJoinActivityRequest) (GroupUsersJoinActivityResponse, error) {

	var list []GroupUsersJoinActivity
	var userId int
	db := gormx.GetInstance()
	db.Raw("select id from fission_mp_user where union_id=? and enterprise_id=?", groupUsersJoinActivityRequest.UnionID, groupUsersJoinActivityRequest.EnterpriseId).Scan(&userId)
	db.Raw("SELECT a.`name`, a.created_at, CASE WHEN a.start_time <= ? AND ?>="+
		" a.end_time THEN '进行中' ELSE '已结束' end as `status` FROM fission_user AS u LEFT JOIN fission_activity AS a ON"+
		" u.activity_id = a.id where u.user_id=? and u.enterprise_id=?", time.Now().Unix(), time.Now().Unix(), userId, groupUsersJoinActivityRequest.EnterpriseId).Scan(&list)

	total := int64(len(list))
	size := groupUsersJoinActivityRequest.Size
	page := groupUsersJoinActivityRequest.Page
	min := size * (page - 1)
	max := page * size
	if max > total {
		max = total
	}
	if total == 1 {
		min = 0
	}

	var cmp GroupUsersJoinActivityResponse
	cmp.Total = total
	cmp.Page = groupUsersJoinActivityRequest.Page
	cmp.Size = groupUsersJoinActivityRequest.Size
	if list == nil {
		cmp.List = make([]GroupUsersJoinActivity, 0)
		return cmp, nil
	}
	cmp.List = list[min:max]
	return cmp, nil
}

// dl api返回结果
type resultData struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    []data `json:"data"`
}
type data struct {
	Category string `json:"category"`
	Remark   string `json:"remark"`
	SortVal  int    `json:"sortVal"`
	TierCode string `json:"tierCode"`
	TierName string `json:"tierName"`
}

type gradeResultData struct {
	ID      string    `json:"id"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    gradeData `json:"data"`
}
type gradeData struct {
	IsUpdateGrade  string `json:"isUpdateGrade"`
	Grade          string `json:"grade"`
	GradeEndDate   string `json:"gradeEndDate"`
	GradeCode      string `json:"gradeCode"`
	GradeStartDate string `json:"gradeStartDate"`
}

func GetUserGrade(ctx ctx.Context, params gin.H) (string, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, config.CfgData.Service.MicroWork.Domain+"/service/v1/data-lake/query-member-grade", clients.JSONBody(params))
	if err != nil {
		return "", err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return "", rErr
	}
	var response = gradeResultData{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return "", jErr
	}
	if response.Status != 200 {
		return "", errorx.NewErrorMsg(response.Status, response.Message)
	}
	return response.Data.Grade, nil
}

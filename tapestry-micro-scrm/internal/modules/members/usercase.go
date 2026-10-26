package members

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	memberReq "micro-scrm/internal/app/http/request/member"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	helper2 "micro-scrm/internal/pkg/helper"
	"strconv"
	"time"

	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
	"micro-scrm/internal/logic/smart_shopping_guide"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// Service 应用服务
type Service struct {
	repo       Repository
	Enterprise EnterpriseInterface
}

// GetMemberIDByUserid 当数据为空的时候返回 nil
func (s Service) GetMemberIDByUserid(c ctx.Context, userID string, enterpriseID int64) (*Member, error) {
	return s.repo.GetMemberIDByUserid(c, userID, enterpriseID)
}

func (s Service) GetMemberId(c ctx.Context, enterpriseID int64, QyUserID string) (int64, error) {
	member := &Member{
		Userid:       QyUserID,
		EnterpriseID: enterpriseID,
		DeletedAt:    0,
	}
	res, err := s.repo.FindMemberByUserid(c, *member)
	if errorIgnore.Ignore(err) != nil {
		return 0, err
	}
	if res == nil {
		return 0, nil
	}
	return res.ID, nil
}

// GetMemberByUserId 通过员工ID获取员工详情
func (s Service) GetMemberByUserId(c ctx.Context, enterpriseID int64, QyUserID string) (*Member, error) {
	memberInfo, err := s.repo.FindMemberByUserid(c, Member{
		Userid:       QyUserID,
		EnterpriseID: enterpriseID,
		DeletedAt:    0,
	})
	if memberInfo == nil {
		return nil, err
	}
	dps, _ := s.repo.GetDepartmentsById(c, memberInfo.ID, memberInfo.EnterpriseID)
	for _, v := range dps {
		memberInfo.Department = append(memberInfo.Department, v.DepartmentID)
	}
	return memberInfo, nil
}

// EnterpriseInterface 企业服务接口
type EnterpriseInterface interface {
	GetInfo(c ctx.Context, enterpriseID int64) (string, string, string, error)
}

// Add 成员添加
func (s Service) Add(c ctx.Context, member *Member) typex.ID {
	res, err := s.repo.Add(c, member)
	helper.Must(err)
	return typex.ID(res)
}

// GetMemberUserID 获取员工userid
func (s Service) GetMemberUserID(c ctx.Context, memberID int64) (string, error) {
	memberInfo, err := s.Detail(c, memberID)
	if err != nil {
		return "", err
	}
	return memberInfo.Userid, nil
}

// Detail 成员详情
func (s Service) Detail(c ctx.Context, id int64) (*Member, error) {
	member, err := s.repo.Detail(c, id)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	memberData := member.(*Member)
	memberData = MemberInfo(*memberData)
	return memberData, nil
}

// Page 成员列表（分页）
func (s Service) Page(c ctx.Context, request MemberListRequest) (int64, *[]Member) {
	offset := (request.Page - 1) * request.Size
	limit := request.Size
	total, list, err := s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request)
	if err != nil {
		return 0, nil
	}
	members := list.(*[]Member)
	var memberList []Member
	for _, value := range *members {
		memberList = append(memberList, *MemberInfo(value))
	}
	return total, &memberList
}

// All 成员列表（全部）
func (s Service) All(c ctx.Context, enterpriseID int64) ([]Member, error) {
	return s.repo.FindAll(c, enterpriseID)
}

// MemberDepartmentRelation 成员组织关系入库
func (s Service) MemberDepartmentRelation(c ctx.Context, md []*MemberDepartment) (int64, error) {
	return s.repo.MemberDepartmentRelation(c, md)
}

// FindMemberByUserid 通过用户ID查找成员
func (s Service) FindMemberByUserid(c ctx.Context, m Member) *Member {
	res, err := s.repo.FindMemberByUserid(c, m)
	if err != nil {
		return nil
	}
	return res
}

// CreateMemberOnLock 使用锁查找人员不存在则创建
func (s Service) CreateMemberOnLock(c ctx.Context, m Member, workSdk *work.QyWorkImpl) (*Member, error) {
	lockKey := fmt.Sprintf(
		"lock:member_%s_%s",
		strconv.FormatInt(m.EnterpriseID, 10),
		m.Userid,
	)
	memberLock, err := lock.Mutex(
		c,
		lock.Key(lockKey),
		lock.TTL(5*time.Second),
		lock.Retry(
			lock.LimitRetry(
				lock.ExponentialBackoff(500*time.Millisecond, 2*time.Second),
				4,
			),
		),
	)

	if err != nil {
		return nil, err
	}

	log.Debug(c).Msgf("加锁准备开始时间:%d", time.Now().UnixNano())
	if err = memberLock.Lock(); err != nil {
		log.Debug(c).Msgf("加锁失败，错误原因:%s, 时间是：%d", err.Error(), time.Now().UnixNano())
		return nil, err
	}
	defer func() {
		memberLock.Unlock()
		log.Debug(c).Msgf("当前解锁时间是:%d", time.Now().UnixNano())
	}()

	// 查询人员是否存在
	member := s.FindMemberByUserid(c, m)
	if member != nil {
		return member, nil
	}

	// 则调用企业微信查询该人员信息
	memberDetail, err := workSdk.WorkMemberDetail(c, work.MemberDetailRequest{
		QyGwRequest: work.QyGwRequest{ReqPL: work.QWCommunalContact},
		Userid:      m.Userid})
	if err != nil {
		return nil, err
	}

	insertMember := TranslateMember(*memberDetail)
	insertMember.ID = uid.ID(ctx.New())
	insertMember.CreatedAt = time.Now().Unix()
	insertMember.UpdatedAt = time.Now().Unix()
	insertMember.EnterpriseID = m.EnterpriseID
	if err = s.BatchInsertMembers(c, []*Member{insertMember}); err != nil {
		return nil, err
	}

	for departmentId := range memberDetail.Department {
		// 部门关系同步
		memberDepartmentDao := membersDao.NewMemberDepartments()
		// 查询员工和部门关系数据关系是否存在
		memberDepartmentInfo, tmpErr := memberDepartmentDao.GetInfo(c, membersModel.MemberDepartmentsModel{
			MicroBaseModel: models.MicroBaseModel{EnterpriseID: m.EnterpriseID},
			MemberID:       insertMember.ID,
			DepartmentID:   int64(departmentId),
		})
		if tmpErr != nil {
			log.Error(c).Msgf("查询员工和部门关系数据,异常：%v", tmpErr)
			return nil, tmpErr
		}

		if memberDepartmentInfo.ID > 0 {
			// 更新员工和部门关系数据
			tmpErr = memberDepartmentDao.EnableByMemberIdAndDepartmentId(c, m.EnterpriseID, insertMember.ID, int64(departmentId))
			if tmpErr != nil {
				log.Error(c).Msgf("更新员工和部门关系数据,异常：%v", tmpErr)
				return nil, tmpErr
			}
			return nil, tmpErr
		}

		// 添加员工和部门关系数据
		tmpErr = memberDepartmentDao.Create(c, membersModel.MemberDepartmentsModel{
			MicroBaseModel: models.CreateMicroBase(m.EnterpriseID, app.UID()),
			MemberID:       insertMember.ID,
			DepartmentID:   int64(departmentId),
		})
		if tmpErr != nil {
			log.Error(c).Msgf("添加员工和部门关系数据,异常：%v", tmpErr)
			return nil, tmpErr
		}
	}

	return insertMember, nil
}

// FindMemberByUseridStatus 通过用户id查找成员 过滤状态
func (s Service) FindMemberByUseridStatus(c ctx.Context, m Member) Member {
	res, err := s.repo.FindMemberByUseridStatus(c, m)
	if err != nil {
		log.Error(c).Msgf("通过用户id查找成员 过滤状态错误:" + err.Error())

		return res
	}
	return res
}

// Update 更新成员信息
func (s Service) Update(c ctx.Context, m Member) int64 {
	id, err := s.repo.Update(c, m)
	helper.Must(err)
	return id
}

// MemberDepartmentBatchInsert 批量新增关系数据
func (s Service) MemberDepartmentBatchInsert(c ctx.Context, md []*MemberDepartment) error {
	return s.repo.MemberDepartmentBatchInsert(c, md)
}

// BatchInsertMembers 批量新增成员
func (s Service) BatchInsertMembers(c ctx.Context, md []*Member) error {
	_, err := s.repo.Add(c, md)
	return err
}

// GetMembersByUserIdsOrMobiles 通过userid或者手机号批量获取员工信息
func (s Service) GetMembersByUserIdsOrMobiles(c ctx.Context, enterpriseID int64, byType int, items []string) ([]*Member, error) {
	var (
		err     error
		members []*Member
	)
	if byType == MemberByUserID {
		members, err = s.repo.GetMembersByUserIds(c, enterpriseID, items)
		if err != nil {
			return nil, err
		}
	} else if byType == MemberByMobile {
		members, err = s.repo.GetMembersByMobiles(c, enterpriseID, items)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	var memberIds []int64
	for _, v := range members {
		memberIds = append(memberIds, v.ID)
	}

	memberDpMap, _ := s.repo.GetMembersDepartments(c, enterpriseID, memberIds)
	for _, v := range members {
		if dps, ok := memberDpMap[v.ID]; ok {
			v.Department = dps
		}
	}
	return members, nil
}

// GetMembersByMarsIds 通过mars_id批量获取员工信息
func (s Service) GetMembersByMarsIds(c ctx.Context, enterpriseId int64, marsIds []string) ([]*Member, error) {
	totalMembers, err := s.repo.GetMembersByMarsIds(c, enterpriseId, marsIds)
	if err != nil {
		return nil, err
	}
	var (
		marsUserIds    []string
		marsMembers    []*Member
		existMemberMap = make(map[string]string)
	)
	for _, v := range totalMembers {
		marsUserIds = append(marsUserIds, v.Userid)
	}
	existMembers, _ := s.repo.GetMembersByUserIds(c, enterpriseId, marsUserIds)
	for _, v := range existMembers {
		existMemberMap[v.Userid] = v.Userid
	}
	for _, mem := range totalMembers {
		if _, ok := existMemberMap[mem.Userid]; ok {
			marsMembers = append(marsMembers, mem)
		}
	}
	return marsMembers, nil
}

func (s Service) GetMembersByDepartIds(c ctx.Context, batchDept *BatchGetMember) (*[]MemberDepartment, error) {
	return s.repo.GetMembersByDepartIds(c, batchDept)
}

// MomentList 选择员工
func (s Service) MomentList(c ctx.Context, momentListRequest *MomentListRequest) (*MomentListResponse, error) {
	offset := (momentListRequest.Page - 1) * momentListRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(momentListRequest.Size)}

	option := &repos.FindPageOption{
		TplPrefix: "t_members/select-member-moment",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MomentList{}
		},
	}

	total, list, err := s.repo.FindPage(c, pageVo, momentListRequest, option)

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	resList := list.(*[]MomentList)
	momentListResponse := &MomentListResponse{
		NoDpList: make([]MomentList, 0),
		List:     *resList,
		Total:    total,
		Page:     momentListRequest.Page,
		Size:     momentListRequest.Size,
	}

	if momentListRequest.GetNoDp > 0 {
		searchCond := map[string]interface{}{
			"mem_username":  momentListRequest.Search.MemUsername,
			"enterprise_id": momentListRequest.EnterpriseId,
			"userid":        momentListRequest.Search.Userid,
			"mobile":        momentListRequest.Search.Mobile,
		}
		momentList := &[]MomentList{}
		memberOption := &repos.ListOption{
			TplPrefix: "t_members/select-no-department-member",
			Fn: func(context ctx.Context, i int64) interface{} {
				return momentList
			},
		}

		_, err := s.repo.List(c, searchCond, memberOption)

		if err == nil && len(*momentList) > 0 {
			momentListResponse.NoDpList = *momentList
		}

	}

	return momentListResponse, nil
}

// GetMemberByMemberIds 通过MemberIds批量获取员工
func (s Service) GetMemberByMemberIds(c ctx.Context, batchMemberIDs *BatchGetMemberByMemberId) (*[]Member, error) {
	return s.repo.GetMemberByMemberIds(c, batchMemberIDs)
}

// GetMemberByDpId 通过部门获取员工
func (s Service) GetMemberByDpId(c ctx.Context, deptId int64, enterprise_id int64) (*[]MemberDepartment, error) {
	return s.repo.GetMemberByDpId(c, deptId, enterprise_id)
}

// SimpleMemberList 查询member表信息
func (s Service) SimpleMemberList(c ctx.Context, request *MemberListRequest) (int64, *[]Member, error) {
	offset := (request.Page - 1) * request.Size
	limit := request.Size
	listOption := &repos.FindPageOption{
		TplPrefix: "t_members/select-simple-member",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]Member{}
		},
	}

	// 分页查询
	total, list, err := s.repo.FindPage(
		c,
		repos.PageVO{Offset: offset, Limit: limit},
		request,
		listOption,
	)

	if err != nil {
		return 0, nil, err
	}

	members := list.(*[]Member)
	var memberList []Member
	for _, value := range *members {
		memberList = append(memberList, *MemberInfo(value))
	}
	return total, &memberList, nil
}

func (s Service) FixByChatId(c ctx.Context, chatId string, enterpriseID int64) error {
	var workImpl = new(work.QyWorkImpl)
	enterpriseService := enterprise.New(enterpriseMysql.New())
	domain, _, _, err := enterpriseService.GetInfo(c, enterpriseID)
	if err != nil {
		log.Info(c).Msgf("获取企业信息失败", err)
		return err
	}
	workImpl.SetDomain(domain)
	err = s.GetGroupChatAndInsert(c, enterpriseID, chatId, workImpl, true)
	if err != nil {
		log.Info(c).Msgf("群聊获取数据失败,err:%s, ChatID:%v", err.Error(), chatId)
	}
	return nil
}

func (s Service) FixAllMemberChat(c ctx.Context, enterpriseID int64, domain, corpID string) error {
	// 查询企业所有员工
	allMembers, err := s.All(c, enterpriseID)
	if err != nil {
		log.Info(c).Msgf("获取企业所有成员失败:%s, enterpriseID:%d", err.Error(), enterpriseID)
		return err
	}
	var workImpl = new(work.QyWorkImpl)
	workImpl.SetDomain(domain)
	var processMemberMap = make(map[string]int)
	for i := range allMembers {
		if _, ok := processMemberMap[allMembers[i].Userid]; ok {
			continue
		}
		processMemberMap[allMembers[i].Userid] = 1
		s.getAllMemberChatList(c, enterpriseID, &(allMembers[i]), workImpl)
		time.Sleep(time.Millisecond * 20)
	}
	return nil
}

func (s Service) getAllMemberChatList(c ctx.Context, enterpriseID int64, member *Member, workImpl *work.QyWorkImpl) error {
	log.Info(c).Msgf("================ FixMemberChat:%s =====================", member.Userid)
	var request = work.ExternalGroupChatListRequestV2{}
	var chatList []*work.ExternalGroupChatItem
	request.OwnerFilter.UserIDList = []string{member.Userid}
	request.Limit = 1000
	// request.Cursor = 0
	for {
		list, err := workImpl.ExternalGroupChatListUpdate(c, request)
		if err != nil {
			log.Info(c).Msgf("获取群聊列表失败,err:%s, request:%v", err.Error(), request)
			break
		}
		chatList = append(chatList, list...)
		if len(list) < request.Limit {
			break
		}
	}

	for _, item := range chatList {
		err := s.GetGroupChatAndInsert(c, enterpriseID, item.ChatID, workImpl, false)
		if err != nil {
			log.Info(c).Msgf("群聊获取出入数据失败,err:%s, ChatID:%v", err.Error(), item.ChatID)
		}
	}
	return nil
}

func (s Service) GetGroupChatAndInsert(c ctx.Context, enterpriseID int64, ChatID string, workImpl *work.QyWorkImpl, updateFlag bool) error {
	req := work.ExternalGroupRequest{
		ChatID:      ChatID,
		NeedName:    1,
		QyGwRequest: work.QyGwRequest{},
	}
	if req.ChatID == "" {
		log.Info(c).Msgf("群聊chatid为空")
		return nil
	}
	workChatDetail, err := workImpl.ExternalGroupChatGet(c, req)
	if err != nil {
		return err
	}
	var adminList []string
	for _, a := range workChatDetail.AdminList {
		adminList = append(adminList, a.Userid)
	}
	admin := "[]"
	if len(adminList) > 0 {
		a1, _ := json.Marshal(adminList)
		admin = string(a1)
	}

	// 查询群是否已经在数据库中存在
	chatReq := customer.WorkGroupChat{
		ChatID:       ChatID,
		EnterpriseID: enterpriseID,
	}
	customerRepo := customerMysql.NewCustomer()
	custemService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	chat, cErr := custemService.FindWorkGroupChatDetail(c, chatReq)
	if cErr != nil {
		log.Info(c).Msgf("查询客户群失败，错误信息：%+v,查询条件：%+v", cErr, chatReq)
		return cErr
	}

	var memberUserIDs []string
	var existsInQwMembers []*customer.WorkGroupChatMember
	// 群存在
	workGroupType := customer.WorkGroupTypeDissolve
	memberUserIDs = append(memberUserIDs, workChatDetail.Owner)
	/*
	*人员列表
	 */
	chatMemberlogs := map[string]work.ExternalGroupMemberList{}
	// =====================客户群新增成员====================
	var chatMemberNum = len(workChatDetail.MemberList)
	if chatMemberNum > 0 {
		for _, chatMemberItem := range workChatDetail.MemberList {
			existsInQwMembers = append(existsInQwMembers, customer.TransferWorkGroupChatMember(&chatMemberItem))
			if chatMemberItem.Type == customer.WorkChatGroupMemberTypeQw {
				memberUserIDs = append(memberUserIDs, chatMemberItem.UserID)
			}
			chatMemberlogs[chatMemberItem.UserID] = chatMemberItem
		}
		workGroupType = customer.WorkGroupTypeExists
	}
	// =====================查询群里的内部群员（企业微信成员）====================
	membersMapFromDb := s.SearchChatMemberInner(c, memberUserIDs, enterpriseID)
	// =====================群信息处理=========================================
	var workGroupChatID int64
	var hasExistWorkGroupMembersFromDb map[string]customer.WorkGroupChatMember
	// 新增群信息
	if chat == nil {
		insertGroupChat := customer.TransferWorkGroupChat(workChatDetail)
		insertGroupChat.ID = uid.ID(c)
		insertGroupChat.CreatedAt = time.Now().Unix()
		insertGroupChat.UpdatedAt = time.Now().Unix()
		insertGroupChat.EnterpriseID = enterpriseID
		insertGroupChat.Type = customer.WorkGroupChatInnerGroup
		insertGroupChat.Status = workGroupType
		insertGroupChat.MemberNum = chatMemberNum
		insertGroupChat.AdminList = admin
		insertGroupChat.FormatNotice()
		if owner, ok := membersMapFromDb[workChatDetail.Owner]; ok {
			insertGroupChat.OwnerMemberID = owner.ID
			// insertGroupChat.OwnerMemberID = owner.Userid
		}
		_, aErr := customerRepo.AddWorkGroupChat(c, *insertGroupChat)
		if aErr != nil {
			return aErr
		}
		workGroupChatID = insertGroupChat.ID
		chat = insertGroupChat
	} else {
		if updateFlag {
			// 更新群信息
			chat.Name = workChatDetail.Name
			chat.Notice = workChatDetail.Notice
			chat.UpdatedAt = time.Now().Unix()
			chat.Status = workGroupType
			chat.MemberNum = chatMemberNum
			chat.AdminList = admin
			chat.FormatNotice()
			if workGroupType == customer.WorkGroupTypeDissolve && chat.DissolveTime == 0 {
				chat.DissolveTime = chat.UpdatedAt
			}
			if owner, ok := membersMapFromDb[workChatDetail.Owner]; ok {
				chat.OwnerMemberID = owner.ID
			}
			uErr := customerRepo.UpdateWorkGroupChat(c, *chat)
			if uErr != nil {
				log.Error(c).Msgf("更新客户群信息失败，更新数据:%+v,错误信息：%+v", chat, uErr)
			}
			workGroupChatID = chat.ID
			hasExistWorkGroupMembersFromDb, _ = searchWorkChatMemberFromDb(c, workGroupChatID, enterpriseID, custemService)
		} else {
			return errors.New("群数据已存在")
		}
	}

	// ====================群成员处理========================================
	// 新增群成员
	var batchChatMembers []customer.WorkGroupChatMember
	// 需要添加的群成员
	batchAddchatMemberlogs := []group_chat.WorkGroupChatMemberLog{}
	// 已经处理的群成员
	hasDealChatMembers := make(map[string]bool)
	for _, chatMemberItem := range existsInQwMembers {
		hasDealChatMembers[chatMemberItem.UserID] = true
		// 是企业内部成员
		if chatMemberItem.Type == customer.WorkChatGroupMemberTypeQw {
			if qwMember, ok := membersMapFromDb[chatMemberItem.UserID]; ok {
				chatMemberItem.MemberID = qwMember.ID
			}
		}
		chatMemberItem.GroupChatID = workGroupChatID
		chatMemberItem.EnterpriseID = enterpriseID
		chatMemberItem.Status = 1
		// 查询客户群成员是否存在于数据库
		if _, ok := hasExistWorkGroupMembersFromDb[chatMemberItem.UserID]; !ok {
			chatMemberItem.ID = uid.ID(c)
			chatMemberItem.CreatedAt = time.Now().Unix()
			chatMemberItem.UpdatedAt = time.Now().Unix()
			chatMemberItem.Status = customer.WorkGroupChatMemberJoinStatus
			batchChatMembers = append(batchChatMembers, *chatMemberItem)
			// 2021-08-27添加
			if _, ok := chatMemberlogs[chatMemberItem.UserID]; ok {
				batchAddchatMemberlogs = append(batchAddchatMemberlogs, group_chat.WorkGroupChatMemberLog{
					MicroScrmBaseModel: model.MicroScrmBaseModel{
						BaseModel: model.BaseModel{
							ID: uid.ID(c),
						},
						EnterpriseID: enterpriseID,
					},
					GroupChatId:     chatMemberItem.GroupChatID,
					ExternalUserid:  chatMemberItem.UserID,
					UnionId:         chatMemberItem.UnionID,
					OperationStatus: 1,
					OperationTime:   int(time.Now().Unix()),
					MemberType:      chatMemberItem.Type,
					JoinScene:       chatMemberItem.JoinScene,
					GroupNickname:   chatMemberlogs[chatMemberItem.UserID].GroupNickname,
					Name:            chatMemberlogs[chatMemberItem.UserID].Name,
					InvitorUserid:   chatMemberItem.InviteUserID,
				})
			}
		}
	}
	// 批量新增客户群成员
	if len(batchChatMembers) > 0 {
		_, aErr := custemService.AddWorkGroupChatMember(c, batchChatMembers)
		if aErr != nil {
			log.Error(c).Msgf("新增客户群成员数据失败，错误信息：%+v", aErr)
		} else {
			s.batchAddChatEvent(c, batchChatMembers, chat, enterpriseID, common.EventEnum.JoinWorkChatGroup, workChatDetail.CreateTime, custemService)
		}
		if len(batchAddchatMemberlogs) > 0 {
			orm.Service.PlanTask.CreatememberLos(batchAddchatMemberlogs)
		}
	}

	// ====================退群群成员处理========================================

	var hasDeletedMembers []int64
	// 需要删除的群成员
	batchDeletechatMemberlogs := []group_chat.WorkGroupChatMemberLog{}
	var deleteMembers []customer.WorkGroupChatMember
	for key, value := range hasExistWorkGroupMembersFromDb {
		if _, ok := hasDealChatMembers[value.UserID]; !ok {
			hasDeletedMembers = append(hasDeletedMembers, value.ID)
			deleteMembers = append(deleteMembers, hasExistWorkGroupMembersFromDb[key])
			// 2021-08-27添加删除日志
			batchDeletechatMemberlogs = append(batchDeletechatMemberlogs, group_chat.WorkGroupChatMemberLog{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID: uid.ID(c),
					},
					EnterpriseID: enterpriseID,
				},
				GroupChatId:     workGroupChatID,
				ExternalUserid:  value.UserID,
				UnionId:         value.UnionID,
				OperationStatus: 2,
				OperationTime:   int(time.Now().Unix()),
				MemberType:      value.Type,
				JoinScene:       value.JoinScene,
				GroupNickname:   chatMemberlogs[value.UserID].GroupNickname,
				Name:            chatMemberlogs[value.UserID].Name,
				InvitorUserid:   value.InviteUserID,
			})
		}
	}
	if len(hasDeletedMembers) > 0 {
		num, dErr := custemService.BatchUpdateWorkGroupChatMembers(c, hasDeletedMembers, customer.WorkGroupChatMember{
			Status:       customer.WorkGroupChatMemberLeaveStatus,
			LeaveTime:    time.Now().Unix(),
			EnterpriseID: enterpriseID,
			UpdatedAt:    time.Now().Unix(),
		})

		if dErr != nil {
			log.Error(c).Msgf("批量删除已退群成员数据失败，错误信息：%+v", dErr)
		} else if num > 0 {
			s.batchAddChatEvent(c, deleteMembers, chat, chat.EnterpriseID, common.EventEnum.LeaveWorkChatGroup, workChatDetail.CreateTime, custemService)
		}
		if len(batchDeletechatMemberlogs) > 0 {
			orm.Service.PlanTask.CreatememberLos(batchDeletechatMemberlogs)
		}
	}
	return nil
}

func (s Service) batchAddChatEvent(ctx ctx.Context, members []customer.WorkGroupChatMember, chat *customer.WorkGroupChat, enterpriseID int64, event common.Event, eventTime int64, customerService *customer.Service) error {
	log.Info(ctx).Msgf("批量设置用户群聊事件，members:%v, chat:%v, event:%s, enterpriseID:%d", members, *chat, event.String(), enterpriseID)
	for _, v := range members {
		if v.UnionID == "" {
			continue
		}
		userEntity, err := customerService.QueryEnterpriseUser(ctx, enterpriseID, &customer.SourceUserBase{UnionID: v.UnionID})
		if err != nil {
			log.Info(ctx).Msgf("添加用户群聊事件：获取用户数据失败：err:%s, unionID:%s, enterpriseID:%d", err.Error(), v.UnionID, enterpriseID)
			continue
		}
		UserService := user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())
		err = UserService.AddUserEvent(ctx, enterpriseID, userEntity.ID, event.String(), chat.Name, chat.ChatID, chat, eventTime)
		if err != nil {
			log.Info(ctx).Msgf("添加用户事件失败：err:%s, userInfo:%+v, chat:%+v", err.Error(), userEntity, chat)
			continue
		}
	}
	return nil
}

// searchWorkChatMemberFromDb 从数据库查询客户群成员信息
func searchWorkChatMemberFromDb(ctx ctx.Context, workGroupChatID int64, enterpriseID int64, customerService *customer.Service) (map[string]customer.WorkGroupChatMember, map[string]*customer.WorkGroupChatMember) {
	workGroupMembersMap := make(map[string]customer.WorkGroupChatMember)
	leavedMembersMap := make(map[string]*customer.WorkGroupChatMember)
	// 查询已有的群成员
	batchReq := customer.BatchWorkGroupMemberRequest{
		Size:         0,
		EnterpriseID: enterpriseID,
		GroupChatID:  workGroupChatID,
	}
	total, workGroupMembers, _ := customerService.BatchFindWorkGroupChatMembers(ctx, batchReq)
	const (
		WorkChatGroupMemberStatusJoin = 1 // 企微群状态：加群
	)
	if total > 0 {
		for i, v := range workGroupMembers {
			if v.Status == WorkChatGroupMemberStatusJoin {
				workGroupMembersMap[v.UserID] = v
			} else {
				leavedMembersMap[v.UserID] = &workGroupMembers[i]
			}
		}
	}
	return workGroupMembersMap, leavedMembersMap
}

// SearchChatMemberInner 查询群里的内部成员
func (s Service) SearchChatMemberInner(ctx ctx.Context, memberUserIDs []string, enterpriseID int64) map[string]Member {
	var membersMap = make(map[string]Member)
	if len(memberUserIDs) > 0 {
		reqMembers := &MemberListRequest{
			UserIDs:      memberUserIDs,
			UserIDCount:  len(memberUserIDs),
			Page:         1,
			Size:         len(memberUserIDs),
			EnterpriseID: enterpriseID,
		}
		total, memberData, _ := s.SimpleMemberList(ctx, reqMembers)
		if total > 0 {
			for _, v := range *memberData {
				membersMap[v.Userid] = v
			}
		}
	}
	return membersMap
}

// migrateMemberChat 迁移成员群聊数据
func (s Service) migrateMemberChat(c ctx.Context, member *Member, workImpl *work.QyWorkImpl, corpID string, customerService *customer.Service) error {
	log.Info(c).Msgf("================ migrateMemberChat:%s =====================", member.Userid)
	var request = work.ExternalGroupChatListRequest{}
	var chatList []*work.ExternalGroupChatItem
	request.OwnerFilter.UserIDList = []string{member.Userid}
	request.Limit = 1000
	request.Offset = 0
	for {
		list, err := workImpl.ExternalGroupChatList(c, request)
		if err != nil {
			log.Info(c).Msgf("获取群聊列表失败,err:%s, request:%v", err.Error(), request)
			break
		}
		chatList = append(chatList, list...)
		if len(list) < request.Limit {
			break
		}
		request.Offset += request.Limit
	}

	var migrationList []*WorkExternalUserMigrationEntity
	now := time.Now().Unix()
	for _, v := range chatList {
		migrationList = append(migrationList, &WorkExternalUserMigrationEntity{
			EnterpriseID: member.EnterpriseID,
			CorpID:       corpID,
			Event:        "change_external_chat",
			ChatID:       v.ChatID,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	_, err := s.AddWorkExternalUserMigration(c, migrationList)
	if err != nil {
		log.Info(c).Msgf("添加外部联系人迁移记录失败, err:%s", err.Error())
		return err
	}
	return nil
}

// AddWorkExternalUserMigration 添加外部联系人迁移记录
func (s Service) AddWorkExternalUserMigration(c ctx.Context, data []*WorkExternalUserMigrationEntity) (int64, error) {
	index := 0
	length := 1000
	var total int64
	var totalData = len(data)
	if totalData > length {
		for {
			end := index + length
			if end > totalData {
				end = totalData
			}
			dataT := data[index:end]
			count, err := s.repo.AddExternalMigration(c, dataT)
			if err != nil {
				log.Info(c).Msgf("添加企微迁移数据失败, err:%s", err.Error())
				return total, err
			}
			total += count
			index = end
			if index >= totalData {
				break
			}
		}
		return total, nil
	} else {
		return s.repo.AddExternalMigration(c, data)
	}
}

// batchUpdateExternalMigration 批量更新迁移记录状态v
func (s Service) batchUpdateExternalMigration(c ctx.Context, list []*WorkExternalUserMigrationEntity) error {
	if len(list) == 0 {
		return nil
	}
	var idList []int64
	for _, v := range list {
		idList = append(idList, v.ID)
	}
	_, err := s.repo.BatchUpdateExternalMigration(c, idList, 1)
	if err != nil {
		log.Info(c).Msgf("更新迁移记录数据状态失败, err:%s", err.Error())
		return err
	}
	return nil
}

// AuthInfoByCode 通过code返回授权信息
func (s Service) AuthInfoByCode(c ctx.Context, req AuthInfoRequest) (*AuthInfoResponse, error) {
	// 请求企业微信
	domain, _, _, err := s.Enterprise.GetInfo(c, req.EnterpriseID)
	if err != nil {
		return nil, err
	}
	workRequest := work.UserInfoRequest{
		Code: req.Code,
		QyGwRequest: work.QyGwRequest{
			ReqPL: work.QWGadgetIdentifier,
		},
	}
	if req.LoginType == "miniprogram" {
		workRequest.ReqPL = work.QWCrsStoreAssistantSelf
	}

	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	authInfoResult, err := workImpl.AuthUserInfo(c, workRequest)
	if err != nil {
		return nil, err
	}

	if authInfoResult == nil {
		return nil, fmt.Errorf("获取用户身份异常")
	}

	ar := new(AuthInfoResponse)
	if authInfoResult.UserId != "" {

		memberInfo, _ := membersDao.NewMembers().GetOneByUserId(c, req.EnterpriseID, authInfoResult.UserId)
		ar.Member = AuthInfoMember{
			UserId:   authInfoResult.UserId,
			DeviceId: authInfoResult.DeviceId,
			MemberId: memberInfo.ID,
			Name:     memberInfo.Name,
		}
	}

	if authInfoResult.OpenId != "" {
		ar.External = AuthInfoExternal{
			OpenId:         authInfoResult.OpenId,
			DeviceId:       authInfoResult.DeviceId,
			ExternalUserid: authInfoResult.ExternalUserid,
		}
	}

	return ar, err
}

// MemberList  选择员工
func (s Service) MemberList(c ctx.Context, momentListRequest *MomentListRequest) (MomentListResponse, error) {
	momentListResponse := MomentListResponse{
		Page: momentListRequest.Page,
		Size: momentListRequest.Size,
	}
	// 根据BA查询BA所在的部门
	var dpIds []int64
	if momentListRequest.Search.DpId != 0 {
		dpIds = append(dpIds, momentListRequest.Search.DpId)
	}
	memberLists, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByNameAndDepartmentId(c, momentListRequest.EnterpriseId, memberReq.MemberSearch{
		Name:   momentListRequest.Search.MemUsername,
		Userid: momentListRequest.Search.Userid,
		Mobile: momentListRequest.Search.Mobile,
	}, dpIds)
	if err != nil {
		return momentListResponse, err
	}
	momentListResponse.List = make([]MomentList, 0)
	uniqueMemberIds := make([]int64, 0)
	for _, value := range memberLists {

		if value.Status == membersModel.MembersStatusEnable && !helper2.InArray(value.MembersId, uniqueMemberIds) {
			momentListResponse.List = append(momentListResponse.List, MomentList{
				UserId:             value.MembersUserid,
				Name:               value.MembersName,
				MemberId:           value.MembersId,
				MemSubscribeStatus: value.Status,
				DeptName:           value.DepartmentName,
				Mobile:             value.Mobile,
			})
			uniqueMemberIds = append(uniqueMemberIds, value.MembersId)
		}
	}
	momentListResponse.Total = int64(len(momentListResponse.List))

	return momentListResponse, nil
}

// 获取同步更新通讯录同步记录
func (s Service) GetSyncContact(ctx ctx.Context, enterpriseID int64) bool {
	return s.repo.GetSyncContact(ctx, enterpriseID)
}

// 手动同步更新通讯录同步记录
func (s Service) SyncContact(ctx ctx.Context, enterpriseID int64, adminID int64) error {
	return s.repo.CommitSyncMembers(ctx, enterpriseID, adminID)
}

// BatchRefreshDepartments 批量插入部门
func (s Service) BatchRefreshDepartments(ctx ctx.Context, corp *enterprise.Corp, departs map[int64]*department.Department) error {
	return s.repo.BatchRefreshDepartments(ctx, corp.EnterpriseID, departs)

}

// BatchInsertMemberDepartment 批量插入员工部门
func (s Service) BatchInsertMemberDepartment(ctx ctx.Context, corp *enterprise.Corp, membersWork []work.MemberDetail, dbDeparts map[int64]department.Department) error {
	return s.repo.BatchInsertMemberDepartment(ctx, corp.EnterpriseID, membersWork, true, dbDeparts)
}

func (s Service) FixAllGroupChatAndMember(c ctx.Context, enterpriseID int64, domain, corpID string) error {
	log.Info(c).Msgf("================ FixGroupAndMember:%d =====================", enterpriseID)
	var request = work.ExternalGroupChatListRequestV2{}
	var chatList []*work.ExternalGroupChatItem
	var workImpl = new(work.QyWorkImpl)
	workImpl.SetDomain(domain)
	request.Limit = 1000
	for {
		list, err := workImpl.ExternalGroupChatListUpdate(c, request)
		if err != nil {
			log.Info(c).Msgf("获取群聊列表失败,err:%s, request:%v", err.Error(), request)
			break
		}
		chatList = append(chatList, list...)
		if len(list) < request.Limit {
			break
		}

	}

	for _, item := range chatList {
		// todo 调试此方法
		err := s.GetGroupChatAndInsert(c, enterpriseID, item.ChatID, workImpl, true)
		if err != nil {
			log.Info(c).Msgf("群聊获取出入数据失败,err:%s, ChatID:%v", err.Error(), item.ChatID)
		}
	}
	return nil
}

// GetMembersInfoByDepartmentIds 多部门选择员工
func (s Service) GetMembersInfoByDepartmentIds(c ctx.Context, enterpriseId int64, departmentIds []int64, memType string) ([]*SimpleMember, error) {
	// 获取所以BA的userId
	var baUserIds []string
	if memType == MemberTypeBA {
		userIds, err := smart_shopping_guide.NewSmartShoppingGuideLogic(enterpriseId).GetBaUserIDList(c, []uint8{})
		if err != nil {
			return nil, err
		}
		baUserIds = userIds
		if len(userIds) == 0 {
			return nil, nil
		}
	}
	return s.repo.GetMembersInfoByDepartmentIds(c, enterpriseId, departmentIds, baUserIds)
}

// 根据员工部门，查询员工的user_Id集合
func (s Service) GetUserIdsByDepId(c ctx.Context, deptId int64, enterpriseID int64) ([]string, error) {
	return s.repo.GetUserIdsByDepId(c, deptId, enterpriseID)
}

// GetExternalUserRelatedMembers 获取外部联系人关联的员工关系
func (s Service) GetExternalUserRelatedMembers(c ctx.Context, enterpriseId int64, workExternalUserId []int64) (map[int64][]int64, error) {
	return s.repo.GetExternalUserRelatedMembers(c, enterpriseId, workExternalUserId)
}

// Code2Session 通过code返回授权信息
func (s Service) Code2Session(c ctx.Context, req Code2SessionRequest) (*work.Code2SessionReq, error) {
	workRequest := work.Code2SessionRequest{
		JsCode:    req.JsCode,
		GrantType: "authorization_code",
		QyGwRequest: work.QyGwRequest{
			ReqPL: work.QWCrsStoreAssistantSelf,
		},
	}

	workImpl := new(work.QyWorkImpl).SetDomain(req.Domain)
	authInfo, err := workImpl.Code2Session(c, workRequest)
	if err != nil {
		return nil, err
	}
	if authInfo.Userid != "" {
		memberInfo, _ := membersDao.NewMembers().GetOne(c, map[string]any{
			"enterprise_id": req.EnterpriseID,
			"userid":        authInfo.Userid,
		})
		authInfo.Mobile = memberInfo.Mobile
		authInfo.Name = memberInfo.Name
		authInfo.MemberId = memberInfo.ID
		if memberInfo.Type == membersModel.MembersTypeSM {
			authInfo.IsManager = 1
		}

		// 获取当前员工的部门关系
		thirdIds, _ := departmentsDao.NewMemberDepartments().GetThirdsByMemberIds(c, req.EnterpriseID, memberInfo.ID)
		if len(thirdIds) > 0 {
			for _, departId := range thirdIds {
				storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
					MicroBaseModel:     models.MicroBaseModel{EnterpriseID: req.EnterpriseID},
					WechatDepartmentID: uint(departId),
				})
				if storeDepartmentInfo.StoreName != "" || storeDepartmentInfo.ThirdStoreCode != "" {
					authInfo.StoreCode = storeDepartmentInfo.ThirdStoreCode
					authInfo.StoreName = storeDepartmentInfo.StoreName
					authInfo.StoreCity = storeDepartmentInfo.City
					authInfo.StorePhone = storeDepartmentInfo.Phone
					authInfo.StoreInventoryName = storeDepartmentInfo.InventoryName
					break
				}
			}
		}
	}

	return authInfo, err
}

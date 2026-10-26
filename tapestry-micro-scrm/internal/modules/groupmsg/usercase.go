package groupmsg

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	workDao "micro-scrm/internal/app/dao/work"
	"micro-scrm/internal/pkg/sdk/work"
	"time"

	"micro-scrm/internal/modules/ws_trace"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/uid"
)

type Service struct {
	repo Repository
}

var WorkApi WorkApiInterface
var MemberService MemberInterface
var EnterpriseService EnterpriseInterface
var WsTraceService *ws_trace.Service

// SendMsgTemplate 发送群发消息
func (s *Service) SendMsgTemplate(c ctx.Context, content string, taskMember *ChatGroupMsgTaskMember, domain string, sendType int) (err error) {
	if taskMember.TaskStatus == MsgTaskMemberTaskStatusDone {
		return nil
	}

	defer func() {
		if err != nil {
			_ = s.setTaskMemberTaskStatus(c, taskMember, MsgTaskMemberTaskStatusFailed, err.Error())
			log.Error(c).Err(err).Interface("taskmember", taskMember).Msg("SendMsgTemplate")
		}
	}()
	var userID string
	if userID, err = s.getMemberUserID(c, taskMember.MemberId); err != nil {
		return err
	}

	chatType := ChatTypeUser
	if sendType == MsgTaskTaskSendTypeGroup {
		chatType = ChatTypeGroup
	}

	taskMember.MsgID, err = WorkApi.SendGroupMsg(c, userID, content, chatType, domain, WsTraceService, taskMember.MemberId, taskMember.MsgTaskId, taskMember.EnterpriseId)
	if err != nil {
		return err
	}

	taskMember.Status = MsgTaskMemberStatusNotSend
	return s.setTaskMemberTaskStatus(c, taskMember, MsgTaskMemberTaskStatusDone, "")
}

// setTaskMemberTaskStatus 设置任务成员状态
func (s *Service) setTaskMemberTaskStatus(c ctx.Context, taskMember *ChatGroupMsgTaskMember, taskStatus int, msg string) error {
	if msg != "" {
		taskMember.FailMsg = msg
	}
	taskMember.TaskStatus = taskStatus
	err := s.UpdateTaskMember(c, taskMember)
	if err != nil {
		return err
	}
	return nil
}

// taskLock 任务处理锁
func (s *Service) taskLock(c ctx.Context, enterpriseID int64, taskID int64) (lock.Locker, error) {
	lockKey := fmt.Sprintf("%s_GroupMsgTaskLock_%d_%d", config.CfgData.App, enterpriseID, taskID)
	mutex, err := lock.Mutex(c,
		lock.Key(lockKey),
		lock.TTL(MsgTaskProcessLockTime*time.Second),
		lock.Retry(lock.NoRetry()))
	if err != nil {
		log.Info(c).Msgf("初始化锁失败, err:%s", err.Error())
		return nil, err
	}
	err = mutex.Lock()
	return mutex, err
}

// getEnterpriseDomain 获取企业domain
func (s *Service) getEnterpriseDomain(c ctx.Context, enterpriseID int64) (string, error) {
	domain, _, _, err := EnterpriseService.GetInfo(c, enterpriseID)
	if err != nil {
		return "", err
	}
	if domain == "" {
		return "", errs.EnterpriseNotFund
	}
	return domain, nil
}

// processTask 处理任务逻辑
func (s *Service) processTask(c ctx.Context, taskInfo *ChatGroupMsgTask) error {
	// 任务上锁
	taskLocker, err := s.taskLock(c, taskInfo.EnterpriseId, taskInfo.ID)
	if err != nil {
		return err
	}
	defer taskLocker.Unlock()

	memberList, err := s.QueryTaskMembers(c, taskInfo.EnterpriseId, taskInfo.ID, nil, nil)
	if err != nil {
		return err
	}
	if len(memberList) == 0 {
		return nil
	}
	domain, err := s.getEnterpriseDomain(c, taskInfo.EnterpriseId)
	if err != nil {
		log.Info(c).Msgf("获取企业domain失败, err:%s", err.Error())
		return err
	}
	var taskStatus = MsgTaskTaskStatusDone
	content, err := WorkApi.SetMediaID(c, taskInfo.Content, taskInfo.EnterpriseId, domain)
	if err != nil {
		return err
	}
	for _, member := range memberList {
		err = s.SendMsgTemplate(c, content, member, domain, taskInfo.SendType)
		// 失败的任务定时任务重新跑一次以后
		if err != nil {
			if taskInfo.TaskStatus == MsgTaskTaskStatusFailed {
				taskStatus = MsgTaskTaskStatusFailedDone
			} else {
				taskStatus = MsgTaskTaskStatusFailed
			}
		} else {
			taskInfo.Status = NotSend
		}
	}
	taskInfo.TaskStatus = taskStatus
	if taskInfo.Status != NotSend {
		taskInfo.Status = SendFail
	}
	return s.UpdateMsgTask(c, taskInfo)
}

// QueryTaskMembers 查询任务成员列表
func (s *Service) QueryTaskMembers(c ctx.Context, enterpriseID int64, taskID int64, status []int64, taskStatus []int64) ([]*ChatGroupMsgTaskMember, error) {
	return s.repo.GetMsgTaskMembers(c, taskID, enterpriseID, status, taskStatus)
}

// QueryTaskMembers4Page 分页查询任务成员列表
func (s *Service) QueryTaskMembers4Page(c ctx.Context, enterpriseID int64, taskID int64, status []int64, taskStatus []int64, page int64, size int64) (int64, *[]*ChatGroupMsgTaskMember, error) {
	return s.repo.GetMsgTaskMembers4Page(c, taskID, enterpriseID, status, taskStatus, page, size)
}

// UpdateTaskMember 更新任务成员信息
func (s *Service) UpdateTaskMember(c ctx.Context, taskMember *ChatGroupMsgTaskMember) error {
	if taskMember == nil {
		return nil
	}
	taskMember.UpdatedAt = time.Now().Unix()
	_, err := s.repo.UpdateMsgTaskMember(c, taskMember)
	if err != nil {
		log.Info(c).Msgf("更新群发任务人员信息失败, err:%s, member:%v", err.Error(), taskMember)
		return err
	}
	return nil
}

// UpdateMsgTask 更新群发消息任务信息
func (s *Service) UpdateMsgTask(c ctx.Context, task *ChatGroupMsgTask) error {
	if task == nil {
		return nil
	}
	log.Debug(c).Msgf("UpdateMsgTask:%v", *task)
	task.UpdatedAt = time.Now().Unix()
	_, err := s.repo.UpdateMsgTask(c, task)
	if err != nil {
		log.Info(c).Msgf("更新群发任务信息失败, err:%s, member:%v", err.Error(), task)
		return err
	}
	return nil
}

// UpdateMsgTask 更新群发消息任务状态
func (s *Service) UpdateMsgTaskStatus(c ctx.Context, task *ChatGroupMsgTask) error {
	if task == nil {
		return nil
	}
	_, err := s.repo.UpdateMsgTaskStatus(c, task)
	if err != nil {
		log.Info(c).Msgf("更新群发任务信息失败, err:%s, member:%v", err.Error(), task)
		return err
	}
	return nil
}

// getMemberUserID 获取员工userid
func (s *Service) getMemberUserID(c ctx.Context, memberID int64) (string, error) {
	// todo 加缓存
	userID, err := MemberService.GetMemberUserID(c, memberID)
	if err != nil {
		log.Info(c).Msgf("获取员工userid返回错误")
		return "", err
	}
	// todo 判断userid为空
	return userID, nil
}

// GroupSendAdd 新增群发任务
func (s *Service) GroupSendAdd(c ctx.Context, request *GroupMsgAddRequest) error {
	var (
		err         error
		nowTime     = time.Now().Unix()
		TaskMembers []*ChatGroupMsgTaskMember
		contentJson []byte
		taskId      = uid.ID(c)
		task        *ChatGroupMsgTask
	)
	request.ID = taskId
	request.CreatedAt = nowTime
	contentJson, err = json.Marshal(request.Content)
	if err != nil {
		return err
	}
	//立刻发送
	if AtOnce == request.Type {
		request.SendTime = nowTime
	}
	for _, memberId := range request.MemberIds {
		member := &ChatGroupMsgTaskMember{
			MemberId:     memberId,
			EnterpriseId: request.EnterpriseId,
			MsgTaskId:    taskId,
			ID:           uid.ID(c),
			CreatedAt:    nowTime,
		}
		TaskMembers = append(TaskMembers, member)
	}
	task = &ChatGroupMsgTask{
		ID:           taskId,
		Content:      string(contentJson),
		CreatedAt:    nowTime,
		Name:         request.Name,
		Type:         request.Type,
		SendType:     request.SendType,
		EnterpriseId: request.EnterpriseId,
		SendTime:     request.SendTime,
	}
	return s.repo.GroupSendAdd(c, task, TaskMembers)
}

// GroupSendEdit新增群发任务
func (s *Service) GroupSendEdit(c ctx.Context, request *GroupMsgEditRequest) error {
	var (
		err         error
		nowTime     = time.Now().Unix()
		TaskMembers []*ChatGroupMsgTaskMember
		contentJson []byte
		task        *ChatGroupMsgTask
	)
	contentJson, err = json.Marshal(request.Content)
	if err != nil {
		return err
	}
	//重新插入执行人
	for _, memberId := range request.MemberIds {
		member := &ChatGroupMsgTaskMember{
			MemberId:     memberId,
			EnterpriseId: request.EnterpriseId,
			MsgTaskId:    request.ID,
			ID:           uid.ID(c),
			CreatedAt:    nowTime,
		}
		TaskMembers = append(TaskMembers, member)
	}
	task = &ChatGroupMsgTask{
		ID:           request.ID,
		Content:      string(contentJson),
		UpdatedAt:    nowTime,
		Name:         request.Name,
		Type:         request.Type,
		SendType:     request.SendType,
		EnterpriseId: request.EnterpriseId,
		SendTime:     request.SendTime,
	}

	return s.repo.GroupSendEdit(c, task, TaskMembers)
}

// MemberList 任务明细-执行群主下拉列表
func (s *Service) MemberList(c ctx.Context, request *MemberListRequest) ([]*Member, error) {

	var (
		members []*Member
		err     error
	)
	members, err = s.repo.MemberList(c, request)
	if nil != errorIgnore.Ignore(err) {
		return members, err
	}
	return members, nil
}

// TaskInfo 群发任务详情
func (s *Service) TaskInfo(c ctx.Context, request *TaskIfoRequest) (*TaskInfoResponse, error) {
	var (
		task              *ChatGroupMsgTask
		err               error
		members           []*Member
		memberListRequest *MemberListRequest
		response          *TaskInfoResponse
		content           = make(map[string]interface{})
	)
	task, err = s.repo.GetMsgTaskInfo(c, request.ID)
	if nil != err {
		return response, err
	}
	if nil == task {
		return response, err
	}
	err = json.Unmarshal([]byte(task.Content), &content)
	if nil != err {
		return response, err
	}
	response = &TaskInfoResponse{
		ID:           task.ID,
		EnterpriseId: task.EnterpriseId,
		Type:         task.Type,
		SendType:     task.SendType,
		Name:         task.Name,
		Status:       task.Status,
		TaskStatus:   task.Status,
		FailMsg:      task.FailMsg,
		SendTime:     task.SendTime,
		LastSendMsg:  task.LastSendMsg,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
		DeletedAt:    task.DeletedAt,
		Content:      content,
	}
	memberListRequest = &MemberListRequest{
		ID:           request.ID,
		EnterpriseID: request.EnterpriseID,
	}
	members, err = s.repo.MemberList(c, memberListRequest)
	if nil != err {
		return response, err
	}
	response.Members = members
	return response, nil
}

func (s *Service) TaskDetail(c ctx.Context, request TaskDetailRequest) (*TaskDetailResponse, error) {

	var (
		err        error
		taskDetail = &TaskDetailResponse{}
		total      int64
		list       []*TaskDetailItem
	)
	total, list, err = s.repo.TaskDetailList(c, request)
	if nil != err {
		return taskDetail, err
	}
	taskDetail.List = list
	taskDetail.Total = total
	taskDetail.Page = request.Page
	taskDetail.Size = request.Size
	return taskDetail, nil
}

func (s *Service) TaskDetailMember(c ctx.Context, request TaskDetailMemberRequest) (*TaskDetailMemberResponse, error) {
	var (
		err        error
		taskDetail = &TaskDetailMemberResponse{}
		total      int64
		list       []*TaskDetailMemberItem
	)
	total, list, err = s.repo.TaskDetailMemberList(c, request)
	if nil != err {
		return taskDetail, err
	}
	taskDetail.List = list
	taskDetail.Total = total
	taskDetail.Page = request.Page
	taskDetail.Size = request.Size
	return taskDetail, nil
}

func (s *Service) TaskMemberSendList(c ctx.Context, request TaskMemberSendListReq) (TaskMemberSendListResp, error) {
	var result TaskMemberSendListResp

	userId, err := s.getMemberUserID(c, request.MemberId)
	if err != nil || userId == "" {
		return result, errorx.New("获取员工信息异常")
	}

	if request.MsgId == "" {
		return TaskMemberSendListResp{
			Page: request.Page,
			Size: request.Size,
		}, nil
	}

	workSdk := new(work.QyWorkImpl).SetEnterprise(request.EnterpriseID)

	var limit = 1000
	workReq := work.GroupMsgSendResultRequest{
		MsgID:  request.MsgId,
		UserID: userId,
		Limit:  limit,
	}

	var responseList []*work.GroupMsgSendResultItem
	// 循环调用企微接口查询msgId下的所有发送结果
	for {
		workResp, err := workSdk.GetGroupMsgSendResult(c, workReq)
		if err != nil {
			return result, err
		}
		if workResp.ErrCode != 0 {
			return result, errorx.NewErrorMsg(workResp.ErrCode, workResp.ErrMsg)
		}
		responseList = append(responseList, workResp.SendList...)
		if workResp.NextCursor == "" {
			break
		}
		workReq.Cursor = workResp.NextCursor
	}

	var resultList []TaskMemberSendList
	for _, v := range responseList {
		resultList = append(resultList, TaskMemberSendList{
			UserId:         v.UserID,
			ChatId:         v.ChatID,
			Status:         v.Status,
			SendTime:       v.SendTime,
			ExternalUserID: v.ExternalUserID,
		})
	}

	var finalList []TaskMemberSendList
	// 导出数据标识
	if request.Export == 1 {
		finalList = resultList
	} else {
		// 计算分页
		if request.Page > len(resultList) || request.Page == 0 {
			request.Page = 1
		}
		if request.Size == 0 {
			request.Size = 10
		}
		page := request.Page
		size := request.Size

		// 最终分页数据
		if request.Size > len(resultList) {
			finalList = resultList[page-1:]
		} else {
			finalList = resultList[page-1 : size]
		}
		result.Page = page
		result.Size = size
	}

	// 获取企微返回的外部联系人id(openid),群组id(chat_id)
	var openIds []string
	var chatIds []string
	for _, record := range finalList {
		if record.ExternalUserID != "" {
			openIds = append(openIds, record.ExternalUserID)
		}
		if record.ChatId != "" {
			chatIds = append(chatIds, record.ChatId)
		}
	}

	// 按外部联系人id对应姓名
	userListMap := make(map[string]string)
	if len(openIds) > 0 {
		// 获取分页数据中的外部联系信息
		corpExternalList, err := microScrmDao.NewWxUserCorpExternal().GetMany(c, request.EnterpriseID,
			map[string]interface{}{
				"openid": openIds,
			})
		if err != nil {
			return result, err
		}
		if len(corpExternalList) == 0 {
			return result, nil
		}
		for _, user := range corpExternalList {
			userListMap[user.Openid] = user.Name
		}
	}

	// 按群id对应群信息
	chatListMap := make(map[string]string)
	if len(chatIds) > 0 {
		// 获取分页数据中的群信息
		chatList, err := workDao.NewGroupChat().GetMany(c,
			map[string]interface{}{
				"enterprise_id": request.EnterpriseID,
				"chat_id":       chatIds,
			})
		if err != nil {
			return result, err
		}
		if len(chatList) == 0 {
			return result, nil
		}
		for _, chat := range chatList {
			chatListMap[chat.ChatID] = chat.Name
		}
	}

	var resList []TaskMemberSendList
	for _, record := range finalList {
		resList = append(resList, TaskMemberSendList{
			UserId:         record.UserId,
			ChatId:         record.ChatId,
			ChatName:       chatListMap[record.ChatId],
			UserName:       userListMap[record.ExternalUserID],
			Status:         record.Status,
			SendTime:       record.SendTime,
			ExternalUserID: record.ExternalUserID,
		})
	}

	result.Total = int64(len(resultList))
	result.List = resList

	return result, nil
}

// MsgTaskList 群发任务列表
func (s *Service) MsgTaskList(c ctx.Context, request *MsgTaskRequest) (*MsgTaskListResponse, error) {
	response := &MsgTaskListResponse{
		Page:  request.Page,
		Size:  request.Size,
		Total: 0,
		List:  make([]ChatGroupMsgTaskList, 0),
	}

	offset := (request.Page - 1) * request.Size
	limit := request.Size
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_group_msg_task/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]ChatGroupMsgTaskList{}
		},
	}

	total, list, err := s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request, option)
	if err != nil || total == 0 {
		return response, nil
	}

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	resList := list.(*[]ChatGroupMsgTaskList)
	returnList := make([]ChatGroupMsgTaskList, 0)

	for _, res := range *resList {
		contentJson := make(map[string]interface{})
		jsonError := json.Unmarshal([]byte(res.Content), &contentJson)
		if jsonError != nil {
			return nil, errs.JsonParseErr
		}
		res.ContentObj = contentJson
		returnList = append(returnList, res)
	}

	response.Total = total
	response.List = returnList

	return response, nil
}

// SendGroupMsgCron 定时执行群发消息
func (s *Service) SendGroupMsgCron(c ctx.Context, enterpriseID int64) {
	ev := log.Info(c).Int64("enterprise_id", enterpriseID)
	list, err := s.repo.ByEnterpriseIdGetMsgTaskList(c, enterpriseID)
	ev = ev.Err(err).Interface("list", list)
	defer func() {
		ev.Msg("SendGroupMsgCron")
	}()
	// 当前时间
	for _, task := range list {
		if err = s.processTask(c, task); err != nil {
			ev = ev.AnErr(fmt.Sprintf("error-%d", task.ID), err)
		}
	}
}

// CheckEnterpriseTaskSendResult 检查企业任务发送状态
func (s *Service) CheckEnterpriseTaskSendResult(c ctx.Context, enterpriseID int64) error {
	log.Info(c).Msgf("检查企业下群发任务发送状态, enterpriseID:%d", enterpriseID)
	// 查询企业任务列表
	if enterpriseID == 0 {
		return nil
	}
	// 处理任务数据
	var page int64 = 1
	var size int64 = 20
	var index int64 = 0
	for {
		total, taskList, err := s.GetUnfinishedTaskList(c, enterpriseID, page, size, index)
		if err != nil {
			return err
		}
		if total == 0 || len(taskList) == 0 {
			break
		}

		wg := threading.NewRoutineGroup()
		for _, v := range taskList {
			if index == 0 {
				index = v.ID
			} else {
				if v.ID < index {
					index = v.ID
				}
			}
			tmpValue := v
			wg.RunSafe(func() {
				err = s.CheckTaskSendResult(c, tmpValue)
				if err != nil {
					log.Err(c, err).Interface("检查处理任务发送结果异常", tmpValue)
				}
			})
		}
		wg.Wait()
		if page*size >= total {
			break
		}
	}
	return nil
}

// CheckTaskSendResult 检查处理任务发送结果
func (s *Service) CheckTaskSendResult(c ctx.Context, taskInfo *ChatGroupMsgTask) error {
	log.Info(c).Msgf("检查处理任务发送结果, CheckTaskSendResult, id:%d", taskInfo.ID)
	if !s.checkThisTime(taskInfo) {
		return nil
	}
	log.Info(c).Msgf("检查处理任务发送结果, CheckTaskSendResult-process, task:%v", taskInfo)
	var statusCond []int64
	if taskInfo.Status == AtOnceNotStart || taskInfo.Status == AlreadySend {
		return nil
	}
	if taskInfo.Status == PartNotSend {
		statusCond = []int64{MsgTaskMemberStatusNotSend}
	}
	members, err := s.QueryTaskMembers(c, taskInfo.EnterpriseId, taskInfo.ID, statusCond, nil)
	if err != nil {
		log.Info(c).Msgf("查询任务成员列表失败, err:%s", err.Error())
		return err
	}
	// 未查询到任务成员
	if len(members) == 0 {
		// 如果是部分发送状态，则转任务的status为已发送
		taskInfo.Status = AlreadySend
		err = s.UpdateMsgTask(c, taskInfo)
		if err != nil {
			return err
		}
		return nil
	}
	domain, err := s.getEnterpriseDomain(c, taskInfo.EnterpriseId)
	if err != nil {
		return err
	}
	var sendOver = true
	var partSend = false
	for _, member := range members {
		err = s.checkTaskMemberSendResult(c, member, domain, taskInfo.SendType)
		if member.Status == MsgTaskMemberStatusNotSend {
			sendOver = false
		}
		if member.Status == MsgTaskMemberStatusSend {
			partSend = true
		}
		if member.SendTime > int64(taskInfo.LastSendMsg) {
			taskInfo.LastSendMsg = int(member.SendTime)
		}
	}
	if sendOver {
		taskInfo.Status = AlreadySend
	} else if partSend && taskInfo.Status == NotSend {
		taskInfo.Status = PartNotSend
	}
	taskInfo.LatestCheckTime = time.Now().Unix()
	err = s.UpdateMsgTask(c, taskInfo)
	if err != nil {
		return err
	}
	return nil
}

// checkThisTime 判断本次是否执行查询当前任务
func (s *Service) checkThisTime(taskInfo *ChatGroupMsgTask) bool {
	if taskInfo == nil {
		return false
	}
	now := time.Now().Unix()
	var sendTime int64
	switch taskInfo.Type {
	case AtOnce:
		sendTime = taskInfo.CreatedAt
	case Timing:
		sendTime = taskInfo.SendTime
	default:
		sendTime = taskInfo.CreatedAt
	}
	//48小时内的任务每次都执行
	if (now - sendTime) < 3600*48 {
		return true
	}
	// 超过48小时的任务每12小时执行一次
	if (now - taskInfo.LatestCheckTime) > 3600*12 {
		return true
	}
	return false
}

// checkTaskMemberSendResult 检查任务成员发送结果
func (s *Service) checkTaskMemberSendResult(c ctx.Context, taskMember *ChatGroupMsgTaskMember, domain string, sendType int) error {
	if taskMember.Status == MsgTaskMemberStatusSend {
		return nil
	}
	sendList, err := s.GetGroupMsgSendResultList(c, taskMember, domain)
	if err != nil {
		return err
	}
	if len(sendList) == 0 {
		return nil
	}

	// 群发总数
	var memberSendTime int64
	for _, v := range sendList {
		v.EnterpriseId = taskMember.EnterpriseId
		v.MsgTaskId = taskMember.MsgTaskId
		v.MsgTaskMemberId = taskMember.MemberId
		v.ID = uid.ID(c)
		v.CreatedAt = time.Now().Unix()
		if v.SendTime > taskMember.SendTime {
			taskMember.SendTime = v.SendTime
		}

		// 发送状态：0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败
		switch v.Status {
		case 0:
			taskMember.NotSendCount++
		case 1:
			taskMember.SuccessCount++
			memberSendTime = v.SendTime
		case 2:
			fallthrough
		case 3:
			taskMember.FailedCount++
		}
	}

	// 群发插入记录
	if sendType == MsgTaskTaskSendTypeGroup {
		_, err = s.repo.AddSendRecord(c, sendList)
		if err != nil {
			log.Info(c).Msgf("添加群发发送记录失败, err:%s", err.Error())
			return err
		}
	}

	if memberSendTime > 0 {
		taskMember.Status = MsgTaskMemberStatusSend
		taskMember.SendTime = memberSendTime
	}

	err = s.UpdateTaskMember(c, taskMember)
	if err != nil {
		return err
	}
	return nil
}

// GetGroupMsgSendResultList 获取群发消息发送结果
func (s *Service) GetGroupMsgSendResultList(c ctx.Context, taskMember *ChatGroupMsgTaskMember, domain string) ([]*ChatGroupMsgSendRecord, error) {
	if taskMember == nil {
		return nil, nil
	}
	userID, err := s.getMemberUserID(c, taskMember.MemberId)
	if err != nil {
		return nil, err
	}
	list, err := WorkApi.GroupMsgSendResult(c, userID, taskMember.MsgID, domain)
	if err != nil {
		log.Info(c).Msgf("获取消息发送结果错误, err:%s", err.Error())
		return nil, err
	}
	return list, nil
}

// GetUnfinishedTaskList 获取未结束的任务列表
func (s *Service) GetUnfinishedTaskList(c ctx.Context, enterpriseID int64, page int64, size int64, index int64) (int64, []*ChatGroupMsgTask, error) {
	// 查询未发送和部分未发送的
	total, list, err := s.repo.GetTaskList(c, enterpriseID, []int64{NotSend, PartNotSend}, page, size, index)
	if err != nil {
		log.Info(c).Msgf("查询任务失败, err:%s, enterpriseID:%d", err.Error(), enterpriseID)
		return 0, nil, err
	}
	if list == nil {
		return total, nil, nil
	}
	return total, *list, nil
}

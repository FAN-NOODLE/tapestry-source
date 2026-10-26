package mysql

import (
	"database/sql"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
	"time"
)

// MemberRepo 定义成员存储
type Repo struct {
	*repos.SQLRepo
}

// UpdateMsgTask 更新群发任务信息
func (r *Repo) UpdateMsgTask(c ctx.Context, entity *groupmsg.ChatGroupMsgTask) (int64, error) {
	return db.Update(c, r.SQLRepo, "t_chat_group_msg_task/update", entity)
}

// UpdateMsgTask 更新群发任务状态
func (r *Repo) UpdateMsgTaskStatus(c ctx.Context, entity *groupmsg.ChatGroupMsgTask) (int64, error) {
	return db.Update(c, r.SQLRepo, "t_chat_group_msg_task/update-status", entity)
}

// UpdateMsgTaskMember 更新群发任务员工信息
func (r *Repo) UpdateMsgTaskMember(c ctx.Context, entity *groupmsg.ChatGroupMsgTaskMember) (int64, error) {
	return db.Update(c, r.SQLRepo, "t_chat_group_msg_task_member/update", entity)
}

// GroupSendAdd 新增群发任务
func (r *Repo) GroupSendAdd(c ctx.Context, request *groupmsg.ChatGroupMsgTask, members []*groupmsg.ChatGroupMsgTaskMember) error {
	var (
		err     error
		cashErr error
	)
	err = db.DB().Tx(c, func(tx *sqlplus.Tx) error {
		err, cashErr = fx.RunSafe(func() error {
			if groupmsg.AtOnce == request.Type {
				request.TaskStatus = groupmsg.MsgTaskTaskStatusDoing
				request.Status = groupmsg.Sending
			} else {
				request.TaskStatus = groupmsg.MsgTaskTaskStatusUndo
				request.Status = groupmsg.AtOnceNotStart
			}
			err = groupSendTaskAdd(c, tx, request)

			if nil != err {
				log.Info(c).Msgf("新增群发任务失败,错误是:%s", err.Error())
				return err
			}
			for _, member := range members {
				err = groupSendMemberAdd(c, tx, member)
				if nil != err {
					log.Info(c).Msgf("新增群发任务成员失败,错误是:%s", err.Error())
					return err
				}
			}
			return nil
		})
		if nil != cashErr {
			return cashErr
		}
		if nil != err {
			return err
		}
		return nil
	})
	return err
}

// GroupSendEdit 修改群发
func (r *Repo) GroupSendEdit(c ctx.Context, request *groupmsg.ChatGroupMsgTask, members []*groupmsg.ChatGroupMsgTaskMember) error {
	var (
		err        error
		cashErr    error
		task       *groupmsg.ChatGroupMsgTask
		oldMembers []*groupmsg.ChatGroupMsgTaskMember
		nowTime    = time.Now().Unix()
		temp       []*groupmsg.ChatGroupMsgTaskMember
	)
	//查找确定任务存在
	task, err = r.GetMsgTaskInfo(c, request.ID)
	if nil != errorIgnore.Ignore(err) {
		return err
	}
	if nil == task {
		return errors.New("当前群发消息任务不存在")
	}
	err = db.DB().Tx(c, func(tx *sqlplus.Tx) error {
		err, cashErr = fx.RunSafe(func() error {
			//由定时任务修改为立即执行
			if groupmsg.AtOnce == request.Type {
				request.TaskStatus = groupmsg.MsgTaskTaskStatusDoing
				request.Status = groupmsg.Sending
				request.SendTime = nowTime
			} else {
				request.Status = task.Status
				request.TaskStatus = task.TaskStatus
			}
			err = groupSendTaskEdit(c, tx, request)
			if nil != err {
				log.Info(c).Msgf("新增群发任务失败,错误是:%s", err.Error())
				return err
			}
			//查找出之前所有的执行人
			oldMembers, err = r.GetMsgTaskMembers(c, request.ID, request.EnterpriseId, nil, nil)
			//求新旧执行人的交集
			temp = same(oldMembers, members)
			//求旧的执行人和temp的差集  是需要删除掉的
			for _, item := range difference(temp, oldMembers) {
				request.DeletedAt = nowTime
				err = deleteMember(c, tx, request, item.MemberId)
				if nil != err {
					log.Info(c).Msgf("删除不需要的执行人错误,错误是:%s", err.Error())
					return err
				}
			}
			request.DeletedAt = 0
			//求新的执行人和temp的差集  是需要新增的
			for _, item := range difference(temp, members) {
				err = groupSendMemberAdd(c, tx, item)
				if nil != err {
					log.Info(c).Msgf("新增群发任务执行人失败,错误是:%s", err.Error())
					return err
				}
			}
			return nil
		})
		if nil != cashErr {
			return cashErr
		}
		if nil != err {
			return err
		}
		return nil
	})
	return err
}

// slice求交集
func same(slice1 []*groupmsg.ChatGroupMsgTaskMember, slice2 []*groupmsg.ChatGroupMsgTaskMember) (result []*groupmsg.ChatGroupMsgTaskMember) {
	temp := make(map[int64]interface{})
	for _, item := range slice1 {
		temp[item.MemberId] = true
	}

	for _, item := range slice2 {
		if _, ok := temp[item.MemberId]; ok {
			result = append(result, item)
		}
	}
	return result
}

// slice求差集
func difference(slice1 []*groupmsg.ChatGroupMsgTaskMember, slice2 []*groupmsg.ChatGroupMsgTaskMember) (result []*groupmsg.ChatGroupMsgTaskMember) {
	temp := make(map[int64]interface{})
	for _, item := range slice1 {
		temp[item.MemberId] = true
	}
	for _, item := range slice2 {
		if _, ok := temp[item.MemberId]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// groupSendMemberAdd 新增群发关联成员
func groupSendMemberAdd(c ctx.Context, tx *sqlplus.Tx, request *groupmsg.ChatGroupMsgTaskMember) error {
	var (
		err  error
		stmt *sqlplus.TemplateStmt
	)

	stmt, err = tx.Template("t_chat_group_msg_task_member/store")
	if nil != err {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if nil != err {
		return err
	}
	return nil
}

// groupSendTaskAdd 新增群发任务
func groupSendTaskAdd(c ctx.Context, tx *sqlplus.Tx, request *groupmsg.ChatGroupMsgTask) error {
	var (
		err  error
		stmt *sqlplus.TemplateStmt
	)
	stmt, err = tx.Template("t_chat_group_msg_task/store")
	if nil != err {
		return err
	}

	_, err = stmt.NamedExec(c, request)
	if nil != err {
		return err
	}
	return nil
}

// groupSendTaskAdd 新增群发任务
func deleteMember(c ctx.Context, tx *sqlplus.Tx, request *groupmsg.ChatGroupMsgTask, memberId int64) error {
	var (
		err  error
		stmt *sqlplus.TemplateStmt
	)
	stmt, err = tx.Template("t_chat_group_msg_task_member/delete")
	if nil != err {
		return err
	}

	_, err = stmt.Exec(c, request.DeletedAt, request.EnterpriseId, request.ID, memberId)
	if nil != err {
		return err
	}
	return nil
}

// groupSendTaskAdd 修改群发任务
func groupSendTaskEdit(c ctx.Context, tx *sqlplus.Tx, request *groupmsg.ChatGroupMsgTask) error {
	var (
		err  error
		stmt *sqlplus.TemplateStmt
	)
	stmt, err = tx.Template("t_chat_group_msg_task/edit")
	if nil != err {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if nil != err {
		return err
	}
	return nil
}

// MemberList任务明细-执行群主下拉列表
func (r *Repo) MemberList(c ctx.Context, request *groupmsg.MemberListRequest) ([]*groupmsg.Member, error) {

	var (
		err     error
		members = make([]*groupmsg.Member, 0)
		stmt    *sqlplus.TemplateStmt
	)

	stmt, err = r.DB.Template("t_chat_group_msg_task_member/member-list", r.Cond2Args(request))
	if nil != err {
		log.Info(c).Msgf("获取任务明细-执行群主下拉列表错误,错误是:%s", err.Error())
		return members, err
	}
	err = stmt.NamedSelect(c, &members, request)
	if nil != err {
		log.Info(c).Msgf("获取任务明细-执行群主下拉列表错误,错误是:%s", err.Error())
		return members, err
	}
	return members, nil
}

// TaskDetailList
func (r *Repo) TaskDetailList(c ctx.Context, request groupmsg.TaskDetailRequest) (int64, []*groupmsg.TaskDetailItem, error) {
	var (
		cmp    = &[]*groupmsg.TaskDetailItem{}
		result interface{}
		err    error
		total  int64
		list   = make([]*groupmsg.TaskDetailItem, 0)
	)
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(request.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_group_msg_task/detail-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, result, err = r.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	list = *result.(*[]*groupmsg.TaskDetailItem)
	return total, list, nil
}

func (r *Repo) TaskDetailMemberList(c ctx.Context, request groupmsg.TaskDetailMemberRequest) (int64, []*groupmsg.TaskDetailMemberItem, error) {
	var (
		cmp    = &[]*groupmsg.TaskDetailMemberItem{}
		result interface{}
		err    error
		total  int64
		list   = make([]*groupmsg.TaskDetailMemberItem, 0)
	)
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: offset, Limit: request.Size}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_group_msg_task_member/detail-member-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, result, err = r.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	list = *result.(*[]*groupmsg.TaskDetailMemberItem)
	return total, list, nil
}

// GetMsgTaskInfo 获取群发任务详情
func (r *Repo) GetMsgTaskInfo(c ctx.Context, taskID int64) (*groupmsg.ChatGroupMsgTask, error) {
	var task groupmsg.ChatGroupMsgTask
	err := db.Get(c, r.SQLRepo, "t_chat_group_msg_task/detail", &task, taskID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}

// GetMsgTaskMembers 查询群发任务员工
func (r *Repo) GetMsgTaskMembers(c ctx.Context, taskID int64, enterpriseID int64, status []int64, taskStatus []int64) ([]*groupmsg.ChatGroupMsgTaskMember, error) {
	var list = make([]*groupmsg.ChatGroupMsgTaskMember, 0)
	var mapping = make(map[string]interface{})
	mapping["taskID"] = taskID
	mapping["enterpriseID"] = enterpriseID
	mapping["status"] = status
	mapping["taskStatus"] = taskStatus
	err := db.NamedSelect(c, r.SQLRepo, "t_chat_group_msg_task_member/query-member-list", &list, mapping, &mapping)
	if err != nil {
		return list, err
	}
	return list, nil
}

// GetMsgTaskMembers4Page 查询群发任务员工(分页)
func (r *Repo) GetMsgTaskMembers4Page(c ctx.Context, taskID int64, enterpriseID int64, status []int64, taskStatus []int64, page int64, size int64) (int64, *[]*groupmsg.ChatGroupMsgTaskMember, error) {
	var resList = &[]*groupmsg.ChatGroupMsgTaskMember{}
	var mapping = make(map[string]interface{})
	mapping["taskID"] = taskID
	mapping["enterpriseID"] = enterpriseID
	mapping["status"] = status
	mapping["taskStatus"] = taskStatus
	total, list, err := db.FindPage(c, r.SQLRepo, "t_chat_group_msg_task_member/query-member-list-page", page, size, mapping, resList)
	if err != nil {
		return 0, nil, err
	}
	resList = list.(*[]*groupmsg.ChatGroupMsgTaskMember)
	return total, resList, nil
}

// AddSendRecord 添加发送记录
func (r *Repo) AddSendRecord(c ctx.Context, records []*groupmsg.ChatGroupMsgSendRecord) (int64, error) {
	return db.Add(c, r.SQLRepo, "t_chat_group_msg_send_record/insert", records)
}

// GetWorkList 获取企业任务列表
func (r *Repo) GetTaskList(c ctx.Context, enterpriseID int64, statusList []int64, page int64, size int64, index int64) (int64, *[]*groupmsg.ChatGroupMsgTask, error) {
	var resList = &[]*groupmsg.ChatGroupMsgTask{}
	var mapping = make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	mapping["status_list"] = statusList
	mapping["index"] = index
	total, list, err := db.FindPage(c, r.SQLRepo, "t_chat_group_msg_task/task-list", page, size, mapping, resList)
	if err != nil {
		return 0, nil, err
	}
	resList = list.(*[]*groupmsg.ChatGroupMsgTask)
	return total, resList, nil
}

// ByEnterpriseIdGetMsgTaskList 通过企业Id查询企业下定制执行群发任务列表
func (r *Repo) ByEnterpriseIdGetMsgTaskList(c ctx.Context, enterpriseID int64) ([]*groupmsg.ChatGroupMsgTask, error) {
	var list = make([]*groupmsg.ChatGroupMsgTask, 0)
	err := db.Query(c, r.SQLRepo, "t_chat_group_msg_task/by-enterprise-id-get-msg-task", &list, enterpriseID)
	return list, err
}

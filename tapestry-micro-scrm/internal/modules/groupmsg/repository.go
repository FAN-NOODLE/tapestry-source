package groupmsg

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 数据库操作
//
//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=groupmsg
type Repository interface {
	repos.IRepository
	UpdateMsgTask(c ctx.Context, entity *ChatGroupMsgTask) (int64, error)
	UpdateMsgTaskStatus(c ctx.Context, entity *ChatGroupMsgTask) (int64, error)
	UpdateMsgTaskMember(c ctx.Context, entity *ChatGroupMsgTaskMember) (int64, error)
	//GroupSendAdd 新增群发
	GroupSendAdd(c ctx.Context, task *ChatGroupMsgTask, members []*ChatGroupMsgTaskMember) error
	//GroupSendEdit 编辑群发
	GroupSendEdit(c ctx.Context, task *ChatGroupMsgTask, members []*ChatGroupMsgTaskMember) error
	//MemberList 任务明细-执行群主下拉列表
	MemberList(c ctx.Context, request *MemberListRequest) ([]*Member, error)
	TaskDetailList(c ctx.Context, request TaskDetailRequest) (int64, []*TaskDetailItem, error)
	TaskDetailMemberList(c ctx.Context, request TaskDetailMemberRequest) (int64, []*TaskDetailMemberItem, error)
	GetMsgTaskInfo(c ctx.Context, taskID int64) (*ChatGroupMsgTask, error)
	ByEnterpriseIdGetMsgTaskList(c ctx.Context, enterpriseID int64) ([]*ChatGroupMsgTask, error)
	// GetMsgTaskMembers 查询任务下的
	GetMsgTaskMembers(c ctx.Context, taskID int64, enterpriseID int64, status []int64, taskStatus []int64) ([]*ChatGroupMsgTaskMember, error)
	// GetMsgTaskMembers4Page 分页查询任务下的
	GetMsgTaskMembers4Page(c ctx.Context, taskID int64, enterpriseID int64, status []int64, taskStatus []int64, page int64, size int64) (int64, *[]*ChatGroupMsgTaskMember, error)
	GetTaskList(c ctx.Context, enterpriseID int64, statusList []int64, page int64, size int64, index int64) (int64, *[]*ChatGroupMsgTask, error)
	// AddSendRecord 添加发送记录数据
	AddSendRecord(c ctx.Context, records []*ChatGroupMsgSendRecord) (int64, error)
}

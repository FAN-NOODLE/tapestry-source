package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/sdk/work"
)

// 群聊天
type GroupChatMemberService struct {
	repo GroupChatMemberRepository
}

// BatchUpdateWorkGroupChatMembers 批量更新
func (s GroupChatMemberService) BatchUpdateWorkGroupChatMembers(c ctx.Context, ids []int64, update WorkGroupChatMember) (int64, error) {
	where := make(map[string]interface{})
	where["ids"] = ids
	return s.repo.BatchUpdateWorkGroupChatMembers(c, where, update)
}

// BatchUpdateExternalWorkGroupChatMembersInfo 批量更新外部联系人群成员信息
func (s GroupChatMemberService) BatchUpdateExternalWorkGroupChatMembersInfo(c ctx.Context, EnterpriseId int64,
	groupChatId int64, chatMembers []work.ExternalGroupMemberList) (int64, error) {
	return s.repo.BatchUpdateExternalWorkGroupChatMembersInfo(c, EnterpriseId, groupChatId, chatMembers)
}

package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// 群聊天
type GroupChatService struct {
	repo GroupChatRepository
}

// GetGroupInfo 获取群信息
func (s *GroupChatService) GetGroupInfo(c ctx.Context, groupChat WorkGroupChat) (*WorkGroupChat, error) {
	ex, err := s.repo.FindWorkGroupChatDetail(c, groupChat)
	if err != nil {
		return nil, err
	}
	return ex, nil
}

package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/chat"
	"micro-scrm/internal/pkg/db"
)

func NewChatMsg() chat.MsgRepository {
	return &ChatMsgRepo{repos.NewSQLRepo("t_chat_msg", db.DB())}
}

func NewChatMember() chat.MemberRepository {
	return &ChatMemberRepo{repos.NewSQLRepo("t_chat_member", db.DB())}
}

func NewChatReceiveUser() chat.ReceiveUserRepository {
	return &ReceiveUserRepo{repos.NewSQLRepo("t_chat_receive_user", db.DB())}
}

func NewChatObject() chat.ObjectRepository {
	return &ChatObjectRepo{repos.NewSQLRepo("t_chat_object", db.DB())}
}

func NewChatMemberToGroup() chat.MemberToGroupRepository {
	return &MemberToGroupRepo{repos.NewSQLRepo("t_chat_member_to_group", db.DB())}
}

func NewChatMsgMedia() chat.MsgMediaRepository {
	return &MsgMediaRepo{repos.NewSQLRepo("t_chat_msg_media", db.DB())}
}

func NewWorkGroupChat() chat.GroupChatRepository {
	return &GroupChatRepo{repos.NewSQLRepo("t_work_group_chat", db.DB())}
}

func NewWorkGroupChatMember() chat.GroupChatMemberRepository {
	return &GroupChatMemberRepo{repos.NewSQLRepo("t_work_group_chat_member", db.DB())}
}

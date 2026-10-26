package chat

func NewMsg(repository MsgRepository) *Service {
	return &Service{
		repo: repository,
	}
}

func NewMember(repository MemberRepository) *MemberService {
	return &MemberService{
		repo: repository,
	}
}

func NewReceive(repository ReceiveUserRepository) *ReceiveUserService {
	return &ReceiveUserService{
		repo: repository,
	}
}

func NewObject(repository ObjectRepository) *ObjectService {
	return &ObjectService{
		repo: repository,
	}
}

func NewMemberToGroup(repository MemberToGroupRepository) *MemberToGroupService {
	return &MemberToGroupService{
		repo: repository,
	}
}

func NewMsgMedia(repository MsgMediaRepository) *MsgMediaService {
	return &MsgMediaService{
		repo: repository,
	}
}

func NewGroupChat(repository GroupChatRepository) *GroupChatService {
	return &GroupChatService{
		repo: repository,
	}
}

func NewGroupChatMember(repository GroupChatMemberRepository) *GroupChatMemberService {
	return &GroupChatMemberService{
		repo: repository,
	}
}

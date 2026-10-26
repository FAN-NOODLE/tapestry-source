package moment

func NewTopic(repository Repository) *TopicService {
	return &TopicService{
		repo: repository,
	}
}

func NewMember(repository Repository) *MemberService {
	return &MemberService{
		repo: repository,
	}
}

func NewTopicMember(repository Repository) *TopicMemberService {
	return &TopicMemberService{
		repo: repository,
	}
}

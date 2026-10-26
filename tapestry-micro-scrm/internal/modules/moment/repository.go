package moment

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository
	MemberFindOne(ctx ctx.Context, memberRequest MomentMemberGetRequest) (*MomentMember, error)
	MemberUpdateOne(ctx ctx.Context, momentMember MomentMember) (int64, error)
	MomentTopicAddOne(c ctx.Context, momentTopic MomentTopic) (int64, error)
	MomentMemberAddMany(c ctx.Context, momentMembers []*MomentMember) (int64, error)
	MomentTopicMemberAddMany(c ctx.Context, momentTopicMembers []*MomentTopicMember) (int64, error)
	MomentTopicFindOne(c ctx.Context, momentTopic *MomentTopic) (*MomentTopic, error)
	MomentTopicUpdateOne(c ctx.Context, momentTopic *MomentTopic) (int64, error)
	MomentTopicMemberBatchDelete(c ctx.Context, topicMemberBatchDelete *TopicMemberBatchDelete) (int64, error)
	MemberBatchGetByMemberIds(c ctx.Context, request MemberBatchGetByMemberIds) (*[]MomentMember, error)
	TopicMemberDeleteOne(c ctx.Context, request TopicMemberDeleteRequest) (bool, error)
	TopicUpdateMemberCount(c ctx.Context, request MomentTopic) (bool, error)
	BatchUpdateMemberTime(c ctx.Context, request BatchUpdateMemberTime) (bool, error)
}

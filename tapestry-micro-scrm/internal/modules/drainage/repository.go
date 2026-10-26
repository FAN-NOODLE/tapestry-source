package drainage

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
)

type Repository interface {
	repos.IRepository
	ChannelsCreate(ctx ctx.Context, Request ContactChannels) (int64, error)
	ChannelsLabelCreate(ctx ctx.Context, req []ContactChannelLabels) error
	ChannelsLabelsList(ctx ctx.Context, req ContactChannelLabelsListRequest) (int64, *[]ContactChannelLabels, error)
	ChannelsList(ctx ctx.Context, req ContactChannelListRequest) (int64, *[]ContactChannelsList, error)
	ChannelDetail(ctx ctx.Context, req ContactChannels) (*ContactChannels, error)
	ChannelEdit(ctx ctx.Context, req CreateContactChannelLabelsRequest) error

	ContactAdd(c ctx.Context, request Contacts) error
	ContactMembersAdd(c ctx.Context, ages []*ContactMembers) error
	ContactUpdate(c ctx.Context, ages Contacts) error
	ContactMembersDelete(c ctx.Context, ages Contacts) error
	ContactSelect(c ctx.Context, request ContactDetailsRequest) (*Contacts, error)
	ContactMemberSelect(c ctx.Context, request ContactDetailsRequest) ([]MemberData, error)
	ContactDepartmentSelect(c ctx.Context, request ContactDetailsRequest) ([]DepartmentData, error)
	ContactDelete(c ctx.Context, request ContactDeleteRequest) error

	InsertWelcomeMember(c ctx.Context, request []*WelcomeMember) (int64, error)
	SelectWelcomeOne(c ctx.Context, request WelcomeDetailRequest) (*Welcome, error)
	SelectWelcomeDefault(c ctx.Context, request WelcomeDetailRequest) (*Welcome, error)
	DeleteWelcomeOne(c ctx.Context, request WelcomeDeleteRequest) (int64, error)
	DeleteWelcomeMemberList(c ctx.Context, request WelcomeDeleteRequest) (int64, error)

	GetWorkExternalUserRelationList(c ctx.Context, req WorkExternalUserRelation) (*[]WorkExternalUserRelation, error)
	GetDrainageContactChannelsList(c ctx.Context, req ContactChannelsData) (*[]ContactChannelsData, error)
	SumAddMemberNumber(c ctx.Context, contactsID int64) error
	SelectLabelsByChannelID(c ctx.Context, request *ContactChannelLabelsListRequest) ([]*ContactChannelLabels, error)
	GetWorkExternalUserRelation(c ctx.Context, req customer.ExternalUserRelation) (*customer.ExternalUserRelation, error)
	GetSynChannelsList(c ctx.Context, req Channel) (*[]Channel, error)
	GetSynContactList(c ctx.Context, req ContactQrCode) (*[]ContactQrCode, error)
	GetSynMemberId(c ctx.Context, req UcMember) (*[]UcMember, error)
	GetSynContacts(ctx ctx.Context, req Contacts) (*Contacts, error)
	GetSynChannelsInfo(ctx ctx.Context, req ContactChannels) (*ContactChannels, error)

	WelcomeChannelList(ctx ctx.Context, req WelcomeValidateRequest) (*[]WelcomeChannelsList, error)
	CountWelcomeByChannelID(ctx ctx.Context, req *WelcomeListRequest) (int64, error)
	UpdateWelcomeMsg(ctx ctx.Context, req *Welcome) (int64, error)

	// ExistRelationCount 查询留存好友总数
	ExistRelationCount(c ctx.Context, req *StatisticsAllReq) (int64, error)
	ExistRelationCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error)
	// IncreaseCustomerCount 新增好友总数
	IncreaseCustomerCount(c ctx.Context, req *StatisticsAllReq) (int64, error)
	IncreaseCustomerCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error)
	// LossCustomerCount 流失好友总数
	LossCustomerCount(c ctx.Context, req *StatisticsAllReq) (int64, error)
	LossCustomerCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error)
	// MemberDelCustomerCountOrChannel 员工删除客户总数
	MemberDelCustomerCountOrChannel(c ctx.Context, req *StatisticsChannelReq, corpId string) (int64, error)
	// CustomerDelMemberCount 好友删员工
	CustomerDelMemberCount(c ctx.Context, req *StatisticsAllReq) (int64, error)
	CustomerDelMemberCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error)
	// GetActivityChannelIDs 查询所有活动的关联渠道
	GetActivityChannelIDs(c ctx.Context, enterpriseID int64) ([]ChannelID, error)
	// GetWelcomeChannelIDs 查询欢迎语的关联渠道
	GetWelcomeChannelIDs(c ctx.Context, enterpriseID int64) ([]ChannelID, error)
	// FixDrainageWelcomeContentRep 查询欢迎语的关联渠道
	FixDrainageWelcomeContentRep(c ctx.Context, enterpriseID int64) error
}

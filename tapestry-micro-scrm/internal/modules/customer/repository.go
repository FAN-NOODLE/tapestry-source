// Package callback 对接Open 平台
package customer

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/sdk/work"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository

	// QueryEnterpriseUser 获取聚合用户
	QueryEnterpriseUser(cc ctx.Context, enterpriseID int64, source *SourceUserBase) (*EnterpriseUser, error)
	// UpdateEnterpriseUserAttr 更新用户属性
	UpdateEnterpriseUserAttr(cc ctx.Context, enterpriseID int64, id int64, attr *EnterpriseUserAttr) error
	// SelectEUidBYDepartment  通过department_id查询外部联系人的euid
	SelectEUidBYDepartment(cc ctx.Context, enterpriseID int64, departmentIDs []int64, page repos.PageVO) (int64, *[]*EUidList, error)
	// UpdateEnterpriseUserBase 更新用户基本属性
	UpdateEnterpriseUserBase(cc ctx.Context, enterpriseID int64, id int64, base *EnterpriseUserBase) error
	// BelongToMembers 客户所属员工列表
	BelongToMembers(cc ctx.Context, externalUserIDs []int64, page repos.PageVO) (int64, *ResponseMemberList, error)

	// BelongToChats 客户加入所属员工群列表
	BelongToChats(cc ctx.Context, enterpriseID int64, externalUserIDs []string, MemberIDs []int64, status int, page repos.PageVO) (int64, *BelongToGroupChatList, error)

	FindExternal(ctx ctx.Context, user WorkExternalUser) (*WorkExternalUser, error)
	FindExternalRelation(ctx ctx.Context, r WorkUserMemberRelation) (*WorkUserMemberRelation, error)
	AddExternalRelation(ctx ctx.Context, r []WorkUserMemberRelation) (int64, error)
	AddExternalRelationLog(ctx ctx.Context, r []WorkUserMemberRelation) (int64, error)
	UpdateExternalRelation(ctx ctx.Context, r WorkUserMemberRelation) (int64, error)
	BatchDeleteUpdateExternalReleation(c ctx.Context, enterpriseID int64, lastUpdateAt int64) (int64, error)
	BatchTagGroupByName(ctx ctx.Context, g []string, enterpriseID int64) ([]WorkExternalTagGroup, error)
	AddTagGroup(ctx ctx.Context, t WorkExternalTagGroup) (int64, error)
	FindWorkGroupChatDetail(ctx ctx.Context, req WorkGroupChat) (*WorkGroupChat, error)
	AddWorkGroupChat(ctx ctx.Context, req WorkGroupChat) (int64, error)
	UpdateWorkGroupChat(ctx ctx.Context, req WorkGroupChat) error
	FindWorkGroupChatMember(ctx ctx.Context, req WorkGroupChatMember) (*WorkGroupChatMember, error)
	AddWorkGroupChatMember(ctx ctx.Context, req []WorkGroupChatMember) (int64, error)
	BatchFindWorGroupChatMember(ctx ctx.Context, req BatchWorkGroupMemberRequest) (int64, []WorkGroupChatMember, error)
	BatchDeleteWorkGroupChatMembers(ctx ctx.Context, ids []int64) error
	BatchUpdateWorkGroupChatMembers(ctx ctx.Context, where map[string]interface{}, update WorkGroupChatMember) (int64, error)
	MemberMultiDepartments(ctx ctx.Context, memberIds []int64) (map[int64]MemberDepartsList, error)
	FindTagGroupFirst(ctx ctx.Context, req WorkExternalTagGroup) (*WorkExternalTagGroup, error)
	UpdateTagGroup(ctx ctx.Context, req WorkExternalTagGroup) error
	AddCustomerCallbackRecord(c ctx.Context, req *work.MemberExternalCallback) (int64, error)
	UpdateCustomerCallbackRecord(ctx ctx.Context, req work.MemberExternalCallback) (int64, error)
	QueryCustomerCallbackRecord(cc ctx.Context, page int, size int, cond work.MemberExternalCallback) ([]*work.MemberExternalCallback, error)
	CallbackRecordDetail(c ctx.Context, req *work.MemberExternalCallback) (*work.MemberExternalCallback, error)
	DelCallbackRecord(c ctx.Context, page int, size int, corpID string) ([]*work.MemberExternalCallback, error)
	LastRecordDetail(c ctx.Context, req *work.MemberExternalCallback) (*work.MemberExternalCallback, error)
	BatchExternalRelations(c ctx.Context, req BatchWorkUserMemberRelationRequest) ([]BatchWorkUserMemberRelationResponse, error)
	GetExternalTagsOfficial(c ctx.Context, enterpriseID int64, tagIDList []string) ([]*WorkExternalTag, error)
	GetExternalUserTagsByTagId(ctx ctx.Context, enterpriseId, relationId int64, tagIds []int64) ([]int64, error)
	AddExternalTag(c ctx.Context, data []*WorkExternalTag) (int64, error)
	GetExternalTagGroups(c ctx.Context, enterpriseID int64, tagGroupIDList []string) ([]*WorkExternalTagGroup, error)
	AddExternalTagGroup(c ctx.Context, data []*WorkExternalTagGroup) (int64, error)
	AddExternalUserTags(c ctx.Context, data []*WorkExternalUserTag) (int64, error)
	QueryExternalUserTags(c ctx.Context, relationID int64) ([]*WorkExternalUserTag, error)
	GetExternalUserRelation(c ctx.Context, enterpriseID int64, workExternalUserID int64, memberUserID string) (*WorkUserMemberRelation, error)
	GetExternalTagsByID(c ctx.Context, idList []int64) ([]*WorkExternalTag, error)
	DeleteExternalUserTags(c ctx.Context, idList []int64) (int64, error)
	QueryCustomerTags(c ctx.Context, enterpriseID int64, memberID int64, tags []string) ([]*WorkExternalTag, error)
	ExternalListRemark(ctx ctx.Context, enterpriseID int64, memberID int64, externalIDList []int64) ([]ExternalRemark, error)
	GetCustomerOfficialAccounts(c ctx.Context, enterpriseID int64, wxOpen map[string]*ThirdUserID) ([]*DetailOfficialAccount, error)
	// 更新企微群聊成员信息
	UpdateWorkChatMemberInfo(c ctx.Context, entity *WorkGroupChatMember) (int64, error)
	// 客户备注信息获取
	GetUserRemark(ctx ctx.Context, id ReturnUID) (*RemarkInfo, error)
	// 客户备注信息保存
	RemarkSave(ctx ctx.Context, info RemarkInfo) (int64, error)
	// 设置外部联系人企业标签
	SetExternalTags(ctx ctx.Context, info SetExternalTags) (int64, error)
	FindCustomerRelation(ctx ctx.Context, memberID int64, workExternalUserID int64, enterpriseID int64) (*WorkCustomerRelationResponse, error)
	BelongToChatsCount(cc ctx.Context, enterpriseID int64, externalUserIDs map[string]*ThirdUserID, memberIDs []int64, status int) (int64, error)
	FindExternalBindMembers(cc ctx.Context, params *ExternalRelationByUserIdRequest) (int64, *WorkUserMemberRelationMemberIds, error)
	// 判断外部联系人是否存在正常的关系
	HasNormalRelation(context ctx.Context, workExternalUserID int64, enterpriseId, channelId int64) (bool, error)
	// 查询在开始时间之前，客户是否添加过这些员工
	ExistsRelationByStartTime(c ctx.Context, req *ExistsRelationByTimeReq) (bool, error)
	// 获取用户和群直接最早添加的关系
	FindFirstJoinTimeInGroups(context ctx.Context, externalUserID string, groupIDs []int64, enterpriseID int64) (int64, error)

	//获取EterpriseUser总数
	AllEterpriseUserForEsList(cc ctx.Context, enterpriseId int64, beginTime, endTime string, pageSize int, elist chan<- []EnterpriseUserEntity)
	UpdateEsResult(cc ctx.Context, esIndex, dtUpdate string, ids []string)
	// 获取指定userid的外部联系人user_id
	GetExternalUserID(ctx ctx.Context, userId int64, enterpriseId int64) (int64, string, error)
	// 根据tagid查询tag下的所有人
	GetExternalUserUsedTag(ctx ctx.Context, enterpriseId int64, tagID int64) ([]WorkExternalRelationTag, error)
	// 根据relation_id，删除用户的tag
	DeleteUserTagsByUserID(c ctx.Context, enterpriseID int64, externalUserID int64)

	//全量计算重复用户
	RepeatCustomerCountJob(ctx ctx.Context, enterpriseID int64) error
	//增量计算重复用户
	RepeatCustomerCount(ctx ctx.Context, enterpriseID int64, userExternalId int64) error

	//全量计算流失用户添加关系时间
	LossCustomerAddRelationTimeJob(ctx ctx.Context, enterpriseId int64, corpID string) error
	//增量计算流失用户添加关系时间
	LossCustomerAddRelationTime(ctx ctx.Context, req *work.MemberExternalCallback) work.LossCustomerAddRelationTimeChannelID

	//根据user_id查员工信息
	FindMemberByUserid(c ctx.Context, enterpriseId int64, userId string) (LossMember, error)
}

type EsRepository interface {
	CreateIndex(ctx.Context, string, string) error
	AddEsIndexDoc(ctx.Context, string, string, string) error
	AddDocsByBulk(ctx.Context, string) (map[string]map[string]int, error)
	DeleteDoc(cc ctx.Context, index, id string) error
	QueryUserList(ctx.Context, *ListRequest) (int64, []*EnterpriseUser, error)
	QueryByUnionID(ctx.Context, *CondUnionID) (int64, []*EnterpriseUser, error)
	QueryUserCount(ctx.Context, *ListCond) (*TagUserCount, error)
	SearchDoc(cc ctx.Context, index, body string) (string, error)
	CountDoc(cc ctx.Context, index, body string) (string, error)
	ExistsIndex(cc ctx.Context, index string) bool
	DelsIndex(cc ctx.Context, index string) bool
}

type WorkUserRepository interface {
	repos.IRepository
	QueryExternalRelationRemark(c ctx.Context, relationID int64, enterpriseID int64) (*WorkExternalUserRemark, error)
	QueryExternalRelationRemarks(ctx.Context, int64, []int64) ([]*WorkExternalUserRemark, error)
	UpdateExternalRelationRemark(ctx.Context, int64, string) error
	QueryWorkUserByUserIDs(c ctx.Context, enterpriseID int64, IDs []int64) (*[]*WorkExternalUser, error)
	// 获取企微员工和外部联系人的关系
	QueryExternalRelationInfo(c ctx.Context, params QueryRelationParams) (*ExternalUserRelationEntity, error)
}

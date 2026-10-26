package members

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/pkg/sdk/work"
)

//Repository 数据库操作
type Repository interface {
	repos.IRepository
	FindAll(ctx ctx.Context, enterpriseID int64) ([]Member, error)
	MemberDepartmentRelation(ctx ctx.Context, md []*MemberDepartment) (int64, error)
	FindMemberByUserid(c ctx.Context, m Member) (*Member, error)
	FindMemberByUseridStatus(c ctx.Context, m Member) (Member, error)
	DeleteAllMemberDepartmentRelation(ctx ctx.Context, req interface{}) error
	MemberDepartmentBatchInsert(ctx ctx.Context, relation []*MemberDepartment) error
	BatchDeleteMembersByIds(ctx ctx.Context, m *BatchDeleteMember) error
	GetMembersByDepartIds(c ctx.Context, batchDept *BatchGetMember) (*[]MemberDepartment, error)
	GetMemberByMemberIds(c ctx.Context, batchDept *BatchGetMemberByMemberId) (*[]Member, error)
	GetMemberByDpId(c ctx.Context, deptId int64, enterprise_id int64) (*[]MemberDepartment, error)
	GetNoDepartmentMember(c ctx.Context, momentListRequest *MomentListRequest) (*[]MomentList, error)
	AddExternalMigration(c ctx.Context, data []*WorkExternalUserMigrationEntity) (int64, error)
	GetExternalMigration(c ctx.Context, page int, size int) ([]*WorkExternalUserMigrationEntity, error)
	UpdateExternalMigration(c ctx.Context, data *WorkExternalUserMigrationEntity) (int64, error)
	BatchUpdateExternalMigration(c ctx.Context, idList []int64, status int) (int64, error)
	GetMemberIDByUserid(c ctx.Context, userID string, enterpriseID int64) (*Member, error)

	GetDepartmentsById(c ctx.Context, memberId, enterpriseID int64) ([]*MemberDepartment, error)
	GetMembersByUserIds(c ctx.Context, enterpriseID int64, userIds []string) ([]*Member, error)
	GetMembersByMarsIds(c ctx.Context, enterpriseID int64, userIds []string) ([]*Member, error)
	GetMembersByMobiles(c ctx.Context, enterpriseID int64, mobiles []string) ([]*Member, error)
	GetMembersDepartments(c ctx.Context, enterpriseID int64, memberIds []int64) (map[int64][]int64, error)
	GetMembersInfoByDepartmentIds(c ctx.Context, enterpriseId int64, departmentIds []int64, baUserIds []string) ([]*SimpleMember, error)
	GetExternalUserRelatedMembers(c ctx.Context, enterpriseId int64, workExternalUserId []int64) (map[int64][]int64, error)

	BatchRefreshDepartments(c ctx.Context, enterpriseID int64, departs map[int64]*department.Department) error
	BatchInsertMemberDepartment(ctx ctx.Context, enterpriseID int64, relation []work.MemberDetail, bRefresh bool, dbDeparts map[int64]department.Department) error
	CommitSyncMembers(c ctx.Context, enterpriseID int64, adminID int64) error
	GetSyncContact(c ctx.Context, enterpriseID int64) bool
	GetUserIdsByDepId(c ctx.Context, deptId int64, enterpriseID int64) ([]string, error)
}

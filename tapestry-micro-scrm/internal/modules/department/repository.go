package department

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//Repository 存储方法
type Repository interface {
	repos.IRepository
	FindDeptByThirdID(c ctx.Context, department Department) (*Department, error)
	SelectList(c ctx.Context, DepartmentListRequest *DepartmentListRequest) (*[]Department, error)
	FindAll(c ctx.Context, enterpriseID int64) ([]Department, error)
	FindUpdateRecordRecent(c ctx.Context, enterpriseID int64) (*DepartmentUpdateRecord, error)
	BatchDeleteDepartmentByIds(c ctx.Context, deleted *BatchDeleteDepartment) error
	FindMinRootDeptId(c ctx.Context, enterpriseID int64, rootDeptId int64) (*Department, error)
	BatchGetDepartmentByThirdIds(c ctx.Context, thirdIds *BatchGetDepartment) (*[]Department, error)
	DepartnameByMembers(c ctx.Context, batchGetDpnameByMembers *BatchGetDpnameByMembers) (*[]TopicMember, error)
	GetDepartmentsByIds(c ctx.Context, enterpriseID int64, dpIds []int64) ([]*Department, error)
	GetChildrenForDepartment(c ctx.Context, enterpriseId int64) (map[int64][]int64, error)
}

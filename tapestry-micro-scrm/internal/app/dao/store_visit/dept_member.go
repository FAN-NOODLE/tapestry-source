package storeVisitDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type DeptMemberDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitDeptMemberModel]
	model     *storeVisitModel.StoreVisitDeptMemberModel
	tableName string
}

func NewDeptMemberDao() *DeptMemberDao {
	return &DeptMemberDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitDeptMemberModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitDeptMemberModel),
		tableName: storeVisitModel.StoreVisitDeptMemberModel{}.TableName(),
	}
}

// GetMemberIdsByDeptName 根据部门名称获取成员id集合
func (d DeptMemberDao) GetMemberIdsByDeptName(c ctx.Context, enterpriseId int64, deptName string) ([]int64, error) {
	var memberIds []int64
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where("dept_name =?", deptName).
		Pluck("member_id", &memberIds).Error
	return memberIds, err
}

// GetDeptNameList 获取部门名称集合
func (d DeptMemberDao) GetDeptNameList(c ctx.Context, enterpriseId int64) ([]string, error) {
	var deptNames []string
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).Group("dept_name").
		Where("dept_name <> ''").
		Pluck("dept_name", &deptNames).Error
	return deptNames, err
}

// GetAllInfo 获取所有配置的员工信息
func (d DeptMemberDao) GetAllInfo(c ctx.Context, enterpriseId int64, query any, args ...any) ([]storeVisitModel.StoreVisitDeptMemberModel, error) {
	var list []storeVisitModel.StoreVisitDeptMemberModel
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where(query, args...).
		Preload("Members.MemberDepartments.Department").
		Find(&list).Error
	return list, err
}

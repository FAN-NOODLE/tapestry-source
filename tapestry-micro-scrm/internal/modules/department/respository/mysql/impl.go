package mysql

import (
	"database/sql"
	"encoding/json"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"github.com/jmoiron/sqlx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/department"
)

//DepartmentRepo 组织存储初始化
type DepartmentRepo struct {
	*repos.SQLRepo
}

// FindDeptByThirdID 根据企业微信组织ID查询组织信息
func (d *DepartmentRepo) FindDeptByThirdID(c ctx.Context, params department.Department) (*department.Department, error) {
	stmt, err := d.DB.Template("t_departments/select-dept-thirdid")
	if err != nil {
		return nil, err
	}
	dept := &department.Department{}
	sErr := stmt.Get(c, dept, params.DeletedAt, params.EnterpriseID, params.ThirdID)
	if sErr != nil {
		return nil, err
	}
	return dept, nil
}

// SelectList 组织列表
func (d *DepartmentRepo) SelectList(c ctx.Context, departmentListRequest *department.DepartmentListRequest) (*[]department.Department, error) {
	stmt, err := d.DB.Template("t_departments/select-list", d.Cond2Args(departmentListRequest))
	if err != nil {
		return nil, err
	}
	var dept []department.Department
	err = stmt.NamedSelect(c, &dept, departmentListRequest)
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// FindAll 查询全量的数据
func (d *DepartmentRepo) FindAll(c ctx.Context, enterpriseID int64) ([]department.Department, error) {
	var list []department.Department
	stmt, err := d.DB.Template("t_departments/select-all")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Select(c, &list, enterpriseID, 0)
	if errStmt != nil {
		return nil, errStmt
	}
	return list, nil
}

// FindUpdateRecordRecent 查询最近同步通讯录记录
func (d *DepartmentRepo) FindUpdateRecordRecent(c ctx.Context, enterpriseID int64) (*department.DepartmentUpdateRecord, error) {
	query := "SELECT admin_id, commit_at as created_at, succ_at as updated_at FROM work_member_sync where enterprise_id=?"
	dept := &department.DepartmentUpdateRecord{}

	if err := d.DB.DB.Get(c, dept, query, enterpriseID); err != nil {
		return nil, err
	}
	return dept, nil
}

// BatchDeleteDepartmentByIds 批量删除组织
func (d *DepartmentRepo) BatchDeleteDepartmentByIds(c ctx.Context, deleted *department.BatchDeleteDepartment) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = deleted.IDs
	stmt, err := d.DB.Template("t_departments/delete-batch", &mapping)
	if err != nil {
		return err
	}
	_, sErr := stmt.NamedExec(c, deleted)
	if sErr != nil {
		return sErr
	}

	return nil
}

// BatchInsertDepartment 批量新增组织
func (d *DepartmentRepo) BatchInsertDepartment(c ctx.Context, de []*department.Department) error {
	stmt, err := d.DB.Template("t_department_update_records/insert")
	if err != nil {
		return err
	}
	_, sErr := stmt.NamedExec(c, de)
	if sErr != nil {
		return sErr
	}
	return nil
}

// FindMinRootDeptId 查询企业根组织
func (d *DepartmentRepo) FindMinRootDeptId(c ctx.Context, enterpriseID int64, rootDeptId int64) (*department.Department, error) {
	stmt, err := d.DB.Template("t_departments/root-department")
	if err != nil {
		return nil, err
	}
	dept := &department.Department{}
	err = stmt.Get(c, dept, enterpriseID, enterpriseID, rootDeptId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return dept, nil
}

func (d *DepartmentRepo) BatchGetDepartmentByThirdIds(c ctx.Context, thirdIds *department.BatchGetDepartment) (*[]department.Department, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = thirdIds.ThirdIDs
	stmt, stmtErr := d.DB.Template("t_departments/batch-get-departments-by-thirds", &mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	members := &[]department.Department{}
	res, sErr := stmt.NamedQuery(c, mapping)
	if sErr != nil {
		return nil, sErr
	}
	scanErr := sqlplus.ScanAll(res, members)
	if scanErr != nil {
		return nil, scanErr
	}

	return members, nil
}

func (d *DepartmentRepo) DepartnameByMembers(c ctx.Context, batchGetDpnameByMembers *department.BatchGetDpnameByMembers) (*[]department.TopicMember, error) {
	mapping := make(map[string]interface{})
	mapping["MemberIds"] = batchGetDpnameByMembers.MemberIds
	mapping["enterprise_id"] = batchGetDpnameByMembers.EnterpriseId
	stmt, stmtErr := d.DB.Template("t_departments/dpname-by-members", &mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	memberInfos := &[]department.TopicMember{}
	res, sErr := stmt.NamedQuery(c, mapping)
	if sErr != nil {
		return nil, sErr
	}
	scanErr := sqlplus.ScanAll(res, memberInfos)
	if scanErr != nil {
		return nil, scanErr
	}

	return memberInfos, nil
}

// GetDepartmentsByIds 按部门ID获取部门列表
func (d *DepartmentRepo) GetDepartmentsByIds(c ctx.Context, enterpriseID int64, dpIds []int64) ([]*department.Department, error) {
	stmt, err := d.DB.Template("t_departments/select-dept-by-ids")
	if err != nil {
		return nil, err
	}
	var data []*department.Department
	query, args, err := sqlx.In(stmt.SQLRaw(), enterpriseID, dpIds)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Select(c, &data, query, args...); err != nil {
		return nil, err
	}
	return data, nil
}

// GetChildrenForDepartment 获取部门，子部门，子子部门
func (d *DepartmentRepo) GetChildrenForDepartment(c ctx.Context, enterpriseId int64) (map[int64][]int64, error) {
	var (
		dps   []*department.ChildrenForDepartment
		dpMap = make(map[int64][]int64)
	)
	if err := d.DB.Select(c, &dps, `select id,if(isnull(child_departments), '[]', child_departments) child_departments from departments where enterprise_id=?`, enterpriseId); err != nil {
		return nil, err
	}
	for _, v := range dps {
		if err := json.Unmarshal(v.TmpIds, &v.Children); err != nil {
			log.Debug(c).Msgf("unmarshal tmp department_ids to children error %+v ", err)
			continue
		}
		dpMap[v.DepartmentID] = v.Children
	}
	return dpMap, nil
}

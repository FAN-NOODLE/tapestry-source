package department

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/sdk/boss"
	"strings"
)

// Service 应用服务
type Service struct {
	repo Repository
}

// Add 新增
func (d *Service) Add(ctx ctx.Context, dept *Department) typex.ID {
	res, err := d.repo.Add(ctx, dept)
	if err != nil {
		return 0
	}
	return typex.ID(res)
}

// Detail 详情
func (d *Service) Detail(ctx ctx.Context, id int64) *Department {
	dept, err := d.repo.Detail(ctx, id)
	if err != nil {
		return nil
	}
	department := dept.(Department)
	return &department
}

// List 组织列表
func (d *Service) List(ctx ctx.Context, DepartmentListRequest *DepartmentListRequest) *[]Department {
	// pid等于0 查询根组织
	if DepartmentListRequest.Pid == 0 {
		// 每个企业展示的组织架构跟目录都不同,根据配置获取
		rootDeptId := int64(1)
		req := &boss.EnterpriseConfigRequest{
			EnterpriseId: DepartmentListRequest.EnterpriseID,
			Key:          boss.EnterpriseConfigRootDepartmentIdKey,
		}
		configData, err := boss.EnterpriseConfig(ctx, req)
		if err != nil {
			log.Error(ctx).Msgf("获取自定义企业配置中root组织架构部门third_id失败", err)
		} else {
			if len(configData.Data) > 0 {
				rootDeptId = configData.Data[0].Value
			}
		}

		var departments []Department
		minDepart, err := d.repo.FindMinRootDeptId(ctx, DepartmentListRequest.EnterpriseID, rootDeptId)
		if err != nil {
			log.Error(ctx).Msgf("查询部门失败", err)
			return nil
		}
		if minDepart.ID > 0 {
			departments = append(departments, *minDepart)
			return &departments
		}
	}
	dept, err := d.repo.SelectList(ctx, DepartmentListRequest)
	if err != nil {
		log.Error(ctx).Msgf("查询部门失败", err)
		return nil
	}
	return dept
}

// Del 删除组织
func (d *Service) Del(ctx ctx.Context, id int64) (int64, error) {
	return d.repo.Del(ctx, id)
}

// FindDeptByThirdID 获取组织信息
func (d Service) FindDeptByThirdID(c ctx.Context, department Department) *Department {
	dept, err := d.repo.FindDeptByThirdID(c, department)
	if err != nil {
		return nil
	}
	return dept
}

// Update 根据主键更新
func (d Service) Update(ctx ctx.Context, department *Department) (int64, error) {
	id, err := d.repo.Update(ctx, department)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindAll 获取组织列表（全部）
func (d Service) FindAll(ctx ctx.Context, enterpriseID int64) ([]Department, error) {
	return d.repo.FindAll(ctx, enterpriseID)
}

// FindUpdateRecordRecent 获取更新记录
func (d *Service) FindUpdateRecordRecent(ctx ctx.Context, enterpriseID int64) *DepartmentUpdateRecord {
	res, err := d.repo.FindUpdateRecordRecent(ctx, enterpriseID)
	if err != nil {
		return nil
	}
	return res
}

// BatchDeleteDepartmentByIds 批量删除组织
func (d *Service) BatchDeleteDepartmentByIds(ctx ctx.Context, delete *BatchDeleteDepartment) error {
	return d.repo.BatchDeleteDepartmentByIds(ctx, delete)
}

// BatchInsertDepartment 批量新增组织
func (d *Service) BatchInsertDepartment(ctx ctx.Context, de []*Department) error {
	_, err := d.repo.Add(ctx, de)
	return err
}

func (d Service) BatchGetDepartmentByThirdIds(c ctx.Context, m *BatchGetDepartment) (*[]Department, error) {
	dps, err := d.repo.BatchGetDepartmentByThirdIds(c, m)
	if err != nil {
		return nil, err
	}
	return dps, nil
}

// DepartmentMomentList 选择部门
func (d Service) DepartmentMomentList(c ctx.Context, momentListRequest *MomentListRequest) (*[]MomentDepartment, error) {
	pageVo := repos.PageVO{Offset: 0, Limit: 1000}
	option := &repos.FindPageOption{
		TplPrefix: "t_departments/select-dept-moment",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MomentDepartment{}
		},
	}
	_, list, err := d.repo.FindPage(c, pageVo, momentListRequest, option)
	if errorIgnore.Ignore(err) != nil {
		return &[]MomentDepartment{}, err
	}
	resList := list.(*[]MomentDepartment)
	if len(*resList) > 0 {
		var result []MomentDepartment
		for _, v := range *resList {
			m := &v
			var storeDepartmentName = v.DpName
			if strings.Contains(storeDepartmentName, "OCF") == true { // 奥莱
				m.StoreType = 2
			} else if strings.Contains(storeDepartmentName, "OC") == true { // 零售
				m.StoreType = 1
			}
			result = append(result, *m)
		}
		return &result, nil
	} else {
		return resList, nil
	}
}

func (d Service) DepartnameByMembers(c ctx.Context, batchGetDpnameByMembers *BatchGetDpnameByMembers) (*[]TopicMember, error) {
	return d.repo.DepartnameByMembers(c, batchGetDpnameByMembers)
}

// GetDepartmentsByIds 根据部门ID获取部门列表
func (d Service) GetDepartmentsByIds(c ctx.Context, enterpriseID int64, ids []int64) ([]*Department, error) {
	return d.repo.GetDepartmentsByIds(c, enterpriseID, ids)
}

// GetChildrenForDepartment 获取指定企业部门下所有部门ID列表
func (d Service) GetChildrenForDepartment(c ctx.Context, enterpriseId int64) (map[int64][]int64, error) {
	return d.repo.GetChildrenForDepartment(c, enterpriseId)
}

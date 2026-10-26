package structure

import "micro-scrm/internal/orm/model/micro_scrm"

// GetMemberInfoByDepartmentIDs 根据部门ID，获取部门内所有的员工
func (s *StructureService) GetMemberInfoByDepartmentIDs(enterpriseID int64, ids []int64) (data []micro_scrm.MemberDepartmentsModel, err error) {
	err = s.Preload("Department", "deleted_at = ?", 0).
		Preload("Member", "deleted_at = ?", 0).
		Where("deleted_at = ?", 0).
		Where("department_id IN ?", ids).
		Find(&data).Error

	return
}

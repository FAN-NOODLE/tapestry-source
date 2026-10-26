package structure

import "micro-scrm/internal/orm/model/micro_scrm"

// GetDepartmentDetail 精确查询部门详情
func (s *StructureService) GetDepartmentDetail(param micro_scrm.DepartmentsModel) (data micro_scrm.DepartmentsModel, err error) {
	err = s.DB.Where(&param).First(&data).Error

	return
}

// GetDepartmentAll 获取企业的部门信息
func (s *StructureService) GetDepartmentAll(enterpriseID int64) (data []micro_scrm.DepartmentsModel, err error) {
	err = s.DB.Where("enterprise_id = ?", enterpriseID).Find(&data).Error

	return
}

// GetDepartmentInfoByThirdIDList 根据third_id集合，获取部门信息
func (s *StructureService) GetDepartmentInfoByThirdIDList(enterpriseID int64, thirdIDList []int) (data []micro_scrm.DepartmentsModel, err error) {
	err = s.DB.Where("enterprise_id = ?", enterpriseID).
		Where("third_id IN ?", thirdIDList).
		Find(&data).Error

	return
}

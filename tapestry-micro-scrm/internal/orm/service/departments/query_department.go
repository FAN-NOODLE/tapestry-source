package departments

import (
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

// GetDepartmentsChildrenIds 根据部门ID，获取所有下级部门id集合
func (d *DepartmentsService) GetDepartmentsChildrenIds(enterpriseID, departmentId int64) []int64 {
	var department micro_scrm.DepartmentsModel
	var departmentIds []int64
	if departmentId > 0 {
		gormx.GetInstance().Debug().Where("id", departmentId).Where("enterprise_id", enterpriseID).First(&department)

		// 所有要返回的部门id集合
		departmentIds = append(departmentIds, departmentId)

		var dpIds []int64
		dpIds = append(dpIds, int64(department.ThirdID))
		for {
			var departments []micro_scrm.DepartmentsModel
			gormx.GetInstance().Where("enterprise_id", enterpriseID).Where("third_parent_id in (?)", dpIds).Find(&departments)
			if len(departments) <= 0 {
				break
			}
			var tmpDpIds []int64
			for _, depart := range departments {
				departmentIds = append(departmentIds, depart.ID)
				tmpDpIds = append(tmpDpIds, int64(depart.ThirdID))
			}
			dpIds = tmpDpIds
		}
	}

	return departmentIds
}

func (d *DepartmentsService) GetUserDepartmentInfo(enterpriseID int64, userID string) (data micro_scrm.MemberDepartmentsModel, err error) {
	mainTable := data.TableName()

	err = d.Model(&data).Joins("Department").
		Joins("Member").
		Where(fmt.Sprintf("%s.enterprise_id = ? AND %s.deleted_at = ?", mainTable, mainTable), enterpriseID, 0).
		Where("Member.userid = ?", userID).
		Take(&data).Error

	return
}

func (d *DepartmentsService) GetDepartmentInfoBYIds(enterpriseID int64, departmentId []int64) (data []micro_scrm.DepartmentsModel, err error) {
	err = gormx.GetInstance().Debug().Where(departmentId).Where("enterprise_id", enterpriseID).Find(&data).Error
	return
}

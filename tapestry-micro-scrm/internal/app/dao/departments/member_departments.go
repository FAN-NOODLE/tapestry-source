package departmentsDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	departmentsModels "micro-scrm/internal/app/models/departments"
	"micro-scrm/internal/pkg/db/gormx"
)

type MemberDepartments struct {
	dao       dao.Dao
	tableName string
}

func NewMemberDepartments() *MemberDepartments {
	return &MemberDepartments{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: departmentsModels.MemberDepartmentsModel{}.TableName(),
	}
}

// GetMemberIdsByThirdIds 根据企微部门id查询部门id
func (d *MemberDepartments) GetMemberIdsByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) (memberIds []int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"departments.third_id":             thirdIds,
			"departments.deleted_at":           0,
		}).
		Joins("left join departments on member_departments.department_id = departments.id").
		Pluck("member_departments.member_id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetThirdsByMemberIds 查询导购所在门店
func (d *MemberDepartments) GetThirdsByMemberIds(c app.Ctx, enterpriseId, memberId int64) (thirdIds []int, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"member_departments.member_id":     memberId,
			"departments.deleted_at":           0,
		}).
		Joins("left join departments on member_departments.department_id = departments.id").
		Pluck("departments.third_id", &thirdIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return

}

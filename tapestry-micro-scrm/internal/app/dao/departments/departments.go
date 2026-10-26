package departmentsDao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/http/request/department"
	departmentsModels "micro-scrm/internal/app/models/departments"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type Departments struct {
	dao       dao.Dao
	tableName string
}

func NewDepartments() *Departments {
	return &Departments{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: departmentsModels.DepartmentsModel{}.TableName(),
	}
}

func (d *Departments) GetOne(c app.Ctx, query interface{}) (departmentsModels.DepartmentsModel, error) {
	var result departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

func (d *Departments) GetInfo(c app.Ctx, data departmentsModels.DepartmentsModel) (departmentsModels.DepartmentsModel, error) {
	var result departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(data).Where("deleted_at = 0").First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *Departments) GetAll(c app.Ctx, data departmentsModels.DepartmentsModel) ([]departmentsModels.DepartmentsModel, error) {
	var result []departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(data).Where("deleted_at = 0").Find(&result).Error
	return result, err
}

func (d *Departments) Create(c app.Ctx, data departmentsModels.DepartmentsModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *Departments) Update(c app.Ctx, data interface{}, cond map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).UpdateColumns(data).Error
	return err
}

func (d *Departments) RemoveById(c app.Ctx, id, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"id":            id,
			"enterprise_id": enterpriseId,
		}).
		Update("deleted_at", time.Now().Unix()).Error
}

func (d *Departments) GetAllDepartments(c app.Ctx, enterpriseId int64) (error, []*departmentsModels.DepartmentsModel) {
	var result []*departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Find(&result).Error

	return err, result
}

// GetDepartmentsThirdIds 通过企业微信中的组织ID获取部门信息
func (d *Departments) GetDepartmentsThirdIds(c app.Ctx, enterpriseId int64, thirdIds []string) ([]departmentsModels.DepartmentsModel, error) {
	var result []departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).
		Where("third_id IN ?", thirdIds).
		Find(&result).Error

	return result, err
}

func (d *Departments) GetDepartmentsMembersByDepartmentId(c app.Ctx, enterpriseId, departmentId int64) ([]departmentsModels.DepartmentsMembers, error) {
	var result []departmentsModels.DepartmentsMembers
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("departments.id as department_id,members.userid,members.id as member_id,members.name as member_name,members.position").
		Joins("left join member_departments on member_departments.department_id = departments.id").
		Joins("left join members on members.id = member_departments.member_id").
		Where("departments.enterprise_id", enterpriseId).
		Where("departments.id", departmentId).
		Where("members.deleted_at = 0 and member_departments.deleted_at = 0 and members.status = 1").
		Find(&result).Error
	return result, err
}

func (d *Departments) GetDepartmentsMemberByMemberIds(c app.Ctx, enterpriseId int64, MemberIds []int64) ([]departmentsModels.DepartmentsMembersInfo, error) {
	var result []departmentsModels.DepartmentsMembersInfo
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("departments.id as department_id,member_departments.member_id,departments.name as department_name,departments.third_id").
		Joins("left join member_departments on member_departments.department_id = departments.id").
		Where("member_departments.member_id", MemberIds).
		Where("departments.enterprise_id", enterpriseId).
		Where("member_departments.deleted_at = 0").
		Scan(&result).Error
	return result, err
}

func (d *Departments) GetDepartmentsNameByMemberId(c app.Ctx, enterpriseId int64, MemberId int64) (string, error) {
	var result string
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("departments.name as department_name").
		Joins("left join member_departments on member_departments.department_id = departments.id").
		Joins("left join members on members.id = member_departments.member_id").
		Where("departments.enterprise_id", enterpriseId).
		Where("members.id", MemberId).
		Where("member_departments.deleted_at", 0).Group("member_departments.id").
		Find(&result).Error
	return result, err
}

// GetDepartmentsInfoThirdIds 通过企业微信中的组织ID获取部门信息
func (d *Departments) GetDepartmentsInfoThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) ([]departmentsModels.DepartmentsModel, error) {
	var result []departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).
		Where("third_id IN ?", thirdIds).
		Find(&result).Error

	return result, err
}

// GetSaStoreLists 获取coach的sa 店铺列表，零售，奥莱： 1为零售，2为奥莱
func (d *Departments) GetSaStoreLists(c app.Ctx, enterpriseId int64, storeType int) (total int64, lists []department.NewMomentDepartments, err error) {
	// selectStr := fmt.Sprintf("`id`,`name` as dp_name,`third_id`,`third_parent_id`,( CASE `third_parent_id` WHEN 0 THEN 0 ELSE 1 END ) AS is_child_department,`order`,%d as `store_type`", storeType)
	selectStr := fmt.Sprintf("`id`,`name`,`third_id`,true as `is_dp`,0 AS leaf,1 as `mem_subscribe_status`,%d as `store_type`", storeType)
	err = d.dao.WithContext(c).Table(d.tableName).
		Select(selectStr).
		Where("enterprise_id=? and deleted_at=0", enterpriseId).Scopes(func(db *gorm.DB) *gorm.DB {
		switch storeType {
		case 0:
			db.Where("name Like ?", "OC%")
		case 1:
			db.Where("name REGEXP ?", "^OC[0-9]")
		case 2:
			db.Where("name REGEXP ?", "^OCF[0-9]")

		}
		return db
	}).Count(&total).Find(&lists).Error
	return
}

// GetDepartmentIdsByThirdIds 根据企微部门id查询部门id
func (d *Departments) GetDepartmentIdsByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int64) (departmentIds []int64, err error) {
	err = d.dao.WithContext(c).Model(departmentsModels.DepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Where("third_id in ?", thirdIds).
		Pluck("id", &departmentIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *Departments) GetMemberDataByMemberIds(c app.Ctx, enterpriseId int64, memberIds []int64) ([]departmentsModels.DepartmentsMembersData, error) {
	var result []departmentsModels.DepartmentsMembersData
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("departments.third_id,departments.name as third_name, members.userid,members.mobile, members.id as member_id, members.name as member_name").
		Joins("left join member_departments on member_departments.department_id = departments.id").
		Joins("left join members on members.id = member_departments.member_id").
		Where(map[string]interface{}{
			"departments.enterprise_id":        enterpriseId,
			"departments.deleted_at":           0,
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"members.enterprise_id":            enterpriseId,
			"members.deleted_at":               0,
			"members.id":                       memberIds,
		}).
		Find(&result).Error
	return result, err
}

// GetSaStoreListsByThirdIds sw&&ks 根据peopleHUB部门id查询出所在的门店类型，奥莱还是零售，
func (d *Departments) GetSaStoreListsByThirdIds(c app.Ctx, enterpriseId int64, storeType int, thirdIds []int64) (total int64, lists []department.NewMomentDepartments, err error) {
	// selectStr := fmt.Sprintf("`id`,`name` as dp_name,`third_id`,`third_parent_id`,( CASE `third_parent_id` WHEN 0 THEN 0 ELSE 1 END ) AS is_child_department,`order`,%d as `store_type`", storeType)
	selectStr := fmt.Sprintf("`id`,`name`,`third_id`,true as `is_dp`,0 AS leaf,1 as `mem_subscribe_status`,%d as `store_type`", storeType)
	err = d.dao.WithContext(c).Table(d.tableName).
		Select(selectStr).
		Where("enterprise_id=? and deleted_at=0", enterpriseId).Where(map[string]interface{}{
		"third_id": thirdIds,
	}).Count(&total).Find(&lists).Error
	return
}

func (d *Departments) DeleteData(c app.Ctx, query interface{}, args ...interface{}) error {
	var model departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Delete(&model).Error
	return err
}

func (d *Departments) GetListByThirdParentId(c app.Ctx, enterpriseId int64, thirdParentId int) ([]departmentsModels.DepartmentsModel, error) {
	var model []departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Model(&departmentsModels.DepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Where("third_parent_id = ?", thirdParentId).
		Find(&model).Error
	return model, err
}

func (d *Departments) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]departmentsModels.DepartmentsModel, error) {
	var result []departmentsModels.DepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Departments) GetMemberDataByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) ([]departmentsModels.DepartmentsMembersData, error) {
	var result []departmentsModels.DepartmentsMembersData
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("departments.third_id,departments.name as third_name, members.userid,members.mobile, members.id as member_id, members.name as member_name").
		Joins("left join member_departments on member_departments.department_id = departments.id").
		Joins("left join members on members.id = member_departments.member_id").
		Where(map[string]interface{}{
			"departments.enterprise_id":        enterpriseId,
			"departments.deleted_at":           0,
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"members.enterprise_id":            enterpriseId,
			"members.deleted_at":               0,
			"departments.third_id":             thirdIds,
		}).
		Find(&result).Error
	return result, err
}

package membersDao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type MemberDepartments struct {
	dao       dao.Dao
	tableName string
}

func NewMemberDepartments() *MemberDepartments {
	return &MemberDepartments{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: membersModel.MemberDepartmentsModel{}.TableName(),
	}
}

func (d *MemberDepartments) GetInfo(c app.Ctx, data membersModel.MemberDepartmentsModel) (membersModel.MemberDepartmentsModel, error) {
	var result membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(data).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *MemberDepartments) GetAllByMemberId(c app.Ctx, enterpriseId, memberId int64) ([]membersModel.MemberDepartmentsModel, error) {
	var result []membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).Unscoped().WithContext(c).
		Where(map[string]interface{}{
			"member_id":     memberId,
			"enterprise_id": enterpriseId,
		}).Find(&result).Error
	return result, err
}

func (d *MemberDepartments) GetAllByDepartmentId(c app.Ctx, enterpriseId, departmentId int64) ([]membersModel.MemberDepartmentsModel, error) {
	var result []membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"department_id": departmentId,
		}).Find(&result).Error
	return result, err
}

func (d *MemberDepartments) Create(c app.Ctx, data membersModel.MemberDepartmentsModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *MemberDepartments) BatchesCreate(c app.Ctx, data []membersModel.MemberDepartmentsModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *MemberDepartments) EnableByMemberIdAndDepartmentId(c app.Ctx, enterpriseId, memberId, departmentId int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).Unscoped().WithContext(c).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
		"department_id": departmentId,
	}).Updates(map[string]interface{}{
		"updated_at": time.Now().Unix(),
		"deleted_at": 0,
	}).Error
	return err
}

func (d *MemberDepartments) RemoveByMemberId(c app.Ctx, memberId, enterpriseId int64) error {
	now := time.Now().Unix()
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
		}).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func (d *MemberDepartments) EnableByMemberIdAndDepIds(c app.Ctx, memberId, enterpriseId int64, depIds []int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
		}).Where("department_id IN ?", depIds).
		Update("deleted_at", 0).Error
}

func (d *MemberDepartments) RemoveByMemberIdAndDepIds(c app.Ctx, memberId, enterpriseId int64, depIds []int64) error {
	now := time.Now().Unix()
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
		}).Where("department_id IN ?", depIds).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func (d *MemberDepartments) RemoveByDepartmentId(c app.Ctx, departmentId, enterpriseId int64) error {
	now := time.Now().Unix()
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"department_id": departmentId,
		}).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func (d *MemberDepartments) GetOneByMemberID(c app.Ctx, enterpriseId int64, memberId int64) (error, membersModel.MemberDepartmentsModel) {
	var result membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"member_id":     memberId,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err, result
}

// SelectEUidBYDepartmentCount 根据部门id获取enterprise_user表的id
func (d *MemberDepartments) SelectEUidBYDepartmentCount(c app.Ctx, enterpriseId int64, departmentIds []int64) (total int64, err error) {
	RawSql := fmt.Sprintf("SELECT count(1) AS count FROM ( SELECT IFNULL(WUSER.euid,0) as euid FROM member_departments AS MD "+
		"LEFT JOIN members AS M ON MD.member_id = M.id LEFT JOIN work_external_user_relation AS WR ON WR.member_id=M.id "+
		"LEFT JOIN micro_scrm_enterprise_%d.wx_user_corp_external AS WUSER ON WUSER.id = WR.work_external_user_id "+
		"WHERE MD.enterprise_id = ? AND MD.department_id IN ? AND WR.relation_status=0 "+
		"AND WUSER.deleted_at=0 GROUP BY WUSER.euid ) as a", enterpriseId)
	err = d.dao.WithContext(c).Raw(RawSql, enterpriseId, departmentIds).Count(&total).Error
	return
}

// SelectEUidBYDepartment 根据部门id获取enterprise_user表的id
func (d *MemberDepartments) SelectEUidBYDepartment(c app.Ctx, enterpriseId int64, departmentIds []int64, page, size int) (data []memberReq.EUidList, err error) {
	RawSql := fmt.Sprintf("SELECT IFNULL(WUSER.euid,0) as euid FROM member_departments AS MD "+
		"LEFT JOIN members AS M ON MD.member_id = M.id LEFT JOIN work_external_user_relation AS WR ON WR.member_id=M.id "+
		"LEFT JOIN micro_scrm_enterprise_%d.wx_user_corp_external AS WUSER ON WUSER.id = WR.work_external_user_id "+
		"WHERE MD.enterprise_id = ? AND MD.department_id IN ? AND WR.relation_status=0 "+
		"AND WUSER.deleted_at=0 GROUP BY WUSER.euid ", enterpriseId)
	err = d.dao.WithContext(c).Raw(RawSql, enterpriseId, departmentIds).Offset(page).Limit(size).Find(&data).Error
	return
}

func (d *MemberDepartments) GetListsByMemberID(c app.Ctx, enterpriseId int64, memberId int64) (error, []membersModel.MemberDepartmentsModel) {
	var result []membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"member_id":     memberId,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Find(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err, result
}

// GetMemberIdsByThirdIds 根据third_id获取导购信息
func (d *MemberDepartments) GetMemberIdsByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int, memberIds []int64) ([]int64, error) {
	var result []int64
	query := d.dao.WithContext(c).Model(membersModel.MemberDepartmentsModel{}).
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where(map[string]interface{}{
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"departments.deleted_at":           0,
		})
	if len(memberIds) > 0 {
		query = query.Where("member_id in (?)", memberIds)
	}
	if len(thirdIds) > 0 {
		query = query.Where("departments.third_id in (?)", thirdIds)
	}

	err := query.Pluck("member_id", &result).Error

	return result, err

}

// GetMemberIds 根据department表id获取member表id集合
func (d *MemberDepartments) GetMemberIds(c app.Ctx, enterpriseId int64, departmentIds []int64) ([]int64, error) {
	var memberIds []int64
	err := d.dao.WithContext(c).Model(membersModel.MemberDepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Where("department_id IN ?", departmentIds).
		Pluck("member_id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return memberIds, err
}

// GetMembersAndDepartmentsByMemberIds 根据memberId集合获取导购信息和部门信息
func (d *MemberDepartments) GetMembersAndDepartmentsByMemberIds(c app.Ctx, enterpriseId int64, ids []int64) (list []membersModel.MemberDepartmentsModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberDepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Where("member_id IN ?", ids).
		Unscoped().
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Department", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Find(&list).Error

	return
}

func (d *MemberDepartments) GetInfoByMemberDepartment(c app.Ctx, enterpriseId int64, memberId, departmentId int64) (membersModel.MemberDepartmentsModel, error) {
	var result membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Unscoped().
		Where("enterprise_id = ?", enterpriseId).
		Where("member_id = ? and department_id = ?", memberId, departmentId).
		Find(&result).Error
	return result, err
}

func (d *MemberDepartments) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(data).Error
	return err
}

func (d *MemberDepartments) GetMemberListByDpIds(c app.Ctx, params memberReq.MemberListByThirdIdsReq) ([]membersModel.MemberDepartmentsModel, error) {
	var result []membersModel.MemberDepartmentsModel
	err := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id", params.EnterpriseId).
		Where("department_id in ?", params.DpIds).
		Group("member_id").
		Preload("Member").
		Find(&result).Error
	return result, err
}

// GetMembersByDepartmentIds  根据department_id集合获取导购信息
func (d *MemberDepartments) GetMembersByDepartmentIds(c app.Ctx, enterpriseId int64, ids []int64) (list []membersModel.MemberDepartmentsModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberDepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Where("department_id IN ?", ids).
		Preload("Member").
		Find(&list).Error
	return
}

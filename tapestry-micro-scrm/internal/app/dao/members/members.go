package membersDao

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	openReq "micro-scrm/internal/app/http/request/open"
	membersResp "micro-scrm/internal/app/http/response/members"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type Members struct {
	dao       dao.Dao
	tableName string
	model     *membersModel.MemberModel
}

func NewMembers() *Members {
	return &Members{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: membersModel.MemberModel{}.TableName(),
		model:     new(membersModel.MemberModel),
	}
}

var StoreManagerPosition = []string{
	"Assistant SM",
	"Store Supervisor",
	"Store Manager",
	"General Manager",
	"Retail Management Trainee SUP",
	"Lead Supervisor I",
	"Assistant Store Manager I",
}
var StoreVisitManagerPosition = []string{
	"Assistant SM",
	"Store Supervisor",
	"Store Manager",
	"General Manager",
	//"Retail Management Trainee SA",
}
var StoreManagerPositionMap = map[string]string{
	"Assistant SM":     "ASM",
	"Store Supervisor": "SUP",
	"Store Manager":    "SM",
	"General Manager":  "GM",
}

var ConversationArchiveTag = []string{
	"Assistant SM",
	"Clienteling Leader",
	"Retail Management Trainee SA",
	"Retail Management Trainee SSA",
	"SALES",
	"Sales Assistant",
	"Sales Associate",
	"Seasonal Sales Associate",
	"Senior Sales",
	"Sr Sales Associate",
	"Store Manager",
	"Store Supervisor",
}

func (d *Members) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Members) GetUnscopedMany(c app.Ctx, query interface{}, args ...interface{}) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Unscoped().Find(&result).Error
	return result, err
}

// GetManyLimit 分页批量获取导购列表
func (d *Members) GetManyLimit(c app.Ctx, enterpriseId int64, pageLimit int, beginID int64) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? ", enterpriseId).
		Where("status = ? ", membersModel.MembersStatusEnable).
		Where("id > ? ", beginID).
		// Where("id = ? ", 3405232142238400).
		Order("id asc").
		Limit(pageLimit).
		Find(&result).Error
	return result, err
}

func (d *Members) GetAllTotal(c app.Ctx, query interface{}) (int64, error) {
	var result int64
	err := d.dao.WithContext(c).Model(d.model).Where(query).Count(&result).Error
	return result, err
}

func (d *Members) GetMemberAndDepartmentList(c app.Ctx, params openReq.StoreWorkListReq) ([]membersModel.MemberModel, int64, error) {
	var result []membersModel.MemberModel
	var total int64
	offset := (params.Page - 1) * params.Size
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id", params.EnterpriseId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.UserId != "" {
				db.Unscoped().Where("userid", params.UserId)
			}
			if params.StartTime > 0 || params.EndTime > 0 {
				db.Where("updated_at BETWEEN ? AND ?", params.StartTime, params.EndTime)
			}
			return db
		}).
		Count(&total).
		Offset(offset).
		Limit(params.Size).
		Preload("MemberDepartments.Department").
		Find(&result).Error
	return result, total, err
}

func (d *Members) GetHkStoreMemberAndDepartmentList(c app.Ctx, params openReq.HkStoreWorkListReq, departmentIds []int) ([]membersModel.MemberModel, int64, error) {
	var result []membersModel.MemberModel
	var total int64
	offset := (params.Page - 1) * params.Size
	query := d.dao.WithContext(c).Model(d.model).
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ?", params.EnterpriseId).
		Where("member_departments.deleted_at = 0").
		Where("departments.deleted_at = 0").
		Where("departments.third_id in ?", departmentIds).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.UserId != "" {
				db.Unscoped().Where("members.userid = ?", params.UserId)
			}
			if params.StartTime > 0 || params.EndTime > 0 {
				db.Where("members.updated_at BETWEEN ? AND ?", params.StartTime, params.EndTime)
			}
			return db
		})
	if err := query.Session(&gorm.Session{}).Distinct("members.id").Count(&total).Error; err != nil {
		return result, total, err
	}
	err := query.Group("members.id").
		Offset(offset).
		Limit(params.Size).
		Preload("MemberDepartments.Department").
		Find(&result).Error
	return result, total, err
}

func (d *Members) GetOne(c app.Ctx, query interface{}) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

func (d *Members) GetOneByUnscoped(c app.Ctx, query interface{}) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Table(d.tableName).Unscoped().Where(query).Find(&result).Error
	return result, err
}

func (d *Members) GetMemberInfo(c app.Ctx, data interface{}) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).
		Where(data).
		Preload("MemberDepartments.Department").
		Order("created_at desc").
		Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *Members) GetAllById(c app.Ctx, enterpriseId int64, id []int64) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("id in (?)", id).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"status":        1,
		}).Preload("MemberDepartments").Find(&result).Error
	return result, err
}

func (d *Members) Create(c app.Ctx, data membersModel.MemberModel) error {
	return d.dao.WithContext(c).Model(d.model).Create(&data).Error
}

func (d *Members) UpdateDataByUserId(c app.Ctx, data interface{}, userId string, enterpriseId int64) error {
	err := d.dao.WithContext(c).Model(d.model).Unscoped().Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"userid":        userId,
	}).UpdateColumns(data).Error
	return err
}

// UpdateStatusByUserId 更新member表员工状态
func (d *Members) UpdateStatusByUserId(c app.Ctx, userId string, enterpriseId, status int64) error {
	return d.dao.WithContext(c).Model(d.model).Unscoped().WithContext(c).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"userid":        userId,
			"deleted_at":    0,
		}).
		Update("status", status).Error
}

func (d *Members) RemoveByUserId(c app.Ctx, userId string, enterpriseId, status int64) error {
	return d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"userid":        userId,
		}).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().Unix(),
		}).Error
}

// GetMemberInfoAndDepartmentsByUserId 获取指定员工的部门和员工信息
func (d *Members) GetMemberInfoAndDepartmentsByUserId(c app.Ctx, enterpriseId int64, userIds []string) ([]membersModel.MemberDepartment, error) {
	var result []membersModel.MemberDepartment
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ? and members.userid IN(?)", enterpriseId, userIds).
		Where("member_departments.deleted_at = 0 and members.deleted_at = 0").
		Group("member_departments.id desc").
		Find(&result).Error
	return result, err
}

func (d *Members) GetMemberInfoAndDepartmentsByMobile(c app.Ctx, enterpriseId int64, mobile string) (membersModel.MemberDepartment, error) {
	var result membersModel.MemberDepartment
	err := d.dao.WithContext(c).Table(membersModel.MemberDepartmentsModel{}.TableName()).
		Select("members.enterprise_id,members.id AS members_id,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN members ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ? and members.mobile = ?", enterpriseId, mobile).
		Where("member_departments.deleted_at = 0 and members.deleted_at = 0").
		Group("member_departments.id").
		Find(&result).Error
	return result, err
}

func (d *Members) GetMobilesByMemberIds(c app.Ctx, enterpriseId int64, memberIds []int64) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Model(d.model).Where(map[string]interface{}{
		"id":            memberIds,
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	}).Pluck("mobile", &result).Error
	return result, err
}

func (d *Members) GetMemberByMobiles(c app.Ctx, enterpriseId int64, mobiles []string) ([]membersModel.MemberInfo, error) {

	var result []membersModel.MemberInfo
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"mobile":        mobiles,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Find(&result).Error
	return result, err
}

// GetUnscopedMemberByMobiles queries uploader identities, including members
// that have been soft deleted. Other member queries keep their original scope.
func (d *Members) GetUnscopedMemberByMobiles(c app.Ctx, enterpriseId int64, mobiles []string) ([]membersModel.MemberInfo, error) {
	var result []membersModel.MemberInfo
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"mobile":        mobiles,
			"enterprise_id": enterpriseId,
		}).Unscoped().Find(&result).Error
	return result, err
}

// GetMemberByMobile 根据手机号查询导购信息
func (d *Members) GetMemberByMobile(c app.Ctx, enterpriseId int64, mobiles string) (membersModel.MemberInfo, error) {

	var result membersModel.MemberInfo
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"mobile":        mobiles,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Find(&result).Error
	return result, err
}

func (d *Members) GetMemberByDepartmentId(c app.Ctx, enterpriseId int64, departments []int64) ([]membersModel.MemberDepartment, error) {

	var result []membersModel.MemberDepartment
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where(map[string]interface{}{
			"members.enterprise_id": enterpriseId,
			"members.status":        membersModel.MembersStatusEnable,
		})

	if len(departments) > 0 {
		query = query.Where("member_departments.department_id in (?)", departments)
	}

	err := query.Order("member_departments.deleted_at asc").Group("members.id").
		Find(&result).Error

	return result, err

}

// GetMemberByThirdIds 根据third_id获取导购信息
func (d *Members) GetMemberByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int, memberIds []int64) ([]memberReq.MemberStaffRes, error) {
	var result []memberReq.MemberStaffRes
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("members.id AS members_id,members.position,members.avatar AS members_avatar,members.name AS members_name,members.userid AS members_userid,departments.third_id AS store_id,departments.name As store_name,members.mobile").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where(map[string]interface{}{
			"members.enterprise_id":         enterpriseId,
			"members.status":                membersModel.MembersStatusEnable,
			"member_departments.deleted_at": 0,
			"departments.deleted_at":        0,
			"members.deleted_at":            0,
		})
	if len(memberIds) > 0 {
		query = query.Where("members.id in (?)", memberIds)
	}
	if len(thirdIds) > 0 {
		query = query.Where("departments.third_id in (?)", thirdIds)
	}

	err := query.Group("members.id").
		Find(&result).Error

	return result, err

}

func (d *Members) GetMemberInfoAndDepartmentsByMemberId(c app.Ctx, enterpriseId, id int64) (membersModel.MemberDepartment, error) {
	var result membersModel.MemberDepartment
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.mobile,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ? and members.id = ? ", enterpriseId, id).
		Where("member_departments.deleted_at = 0 and members.deleted_at=0").
		Group("member_departments.id").
		Find(&result).Error
	return result, err
}

// GetMemberInfoAndAllDepartmentsByMobile 根据导购手机号，获取导购的信息
func (d *Members) GetMemberInfoAndAllDepartmentsByMobile(c app.Ctx, enterpriseId int64, mobile string) ([]membersModel.MemberDepartment, error) {
	var result []membersModel.MemberDepartment
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ? and members.mobile = ?", enterpriseId, mobile).
		Where("member_departments.deleted_at = 0 and members.deleted_at=0").
		Find(&result).Error
	return result, err
}

func (d *Members) GetMemberAndDepartmentsByMemberId(c app.Ctx, enterpriseId, id int64) ([]membersModel.MemberDepartment, error) {
	var result []membersModel.MemberDepartment
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.mobile,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ? and members.deleted_at = 0 and members.id = ? ", enterpriseId, id).
		Where("member_departments.deleted_at = 0 and members.deleted_at=0").
		Group("member_departments.id").
		Find(&result).Error
	return result, err
}

func (d *Members) GetMemberIdsByUserIds(c app.Ctx, enterpriseId int64, userIds []string) (memberIds []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("userid in ?", userIds).
		Pluck("id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *Members) GetMemberIdsByName(c app.Ctx, enterpriseId int64, name string) (memberIds []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("name like ?", "%"+name+"%").
		Pluck("id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *Members) GetMemberIdsByUserIdStr(c app.Ctx, enterpriseId int64, userIdStr string) (memberIds []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("userid like ?", "%"+userIdStr+"%").
		Pluck("id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *Members) GetManyByUserIds(c app.Ctx, enterpriseId int64, userIds []string) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? and userid in ?", enterpriseId, userIds).Find(&result).Error
	return result, err
}

// GetByEmployeeIDs 按企微 userid 或 open_userid 批量查找上传人，包含已离职员工。
func (d *Members) GetByEmployeeIDs(c app.Ctx, enterpriseID int64, employeeIDs []string) ([]membersModel.MemberModel, error) {
	if len(employeeIDs) == 0 {
		return nil, nil
	}
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("id, enterprise_id, userid, open_userid, mobile, deleted_at").
		Where("enterprise_id = ? AND (userid IN ? OR open_userid IN ?)", enterpriseID, employeeIDs, employeeIDs).
		Order("id ASC").Find(&result).Error
	return result, err
}

func (d *Members) GetOneByUserId(c app.Ctx, enterpriseId int64, userId string) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? and userid = ?", enterpriseId, userId).Find(&result).Error
	return result, err
}

func (d *Members) GetUserId(c app.Ctx, enterpriseId int64, userId string) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? and userid like ?", enterpriseId, "%"+userId+"%").Find(&result).Error
	return result, err
}

func (d *Members) GetMemberCount(c app.Ctx, enterpriseId int64) int64 {
	var count int64
	d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? and status = ?", enterpriseId, membersModel.MembersStatusEnable).Count(&count)
	return count
}

// GetUserIdsByName 根据导购昵称匹配导购的userid
func (d *Members) GetUserIdsByName(c app.Ctx, enterpriseId int64, name string) (userIds []string, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("name like ?", "%"+name+"%").
		Pluck("userid", &userIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetMemberInfoAndDepartmentsByNameAndDepartmentId 根据导购昵称和部门id搜索导购信息
func (d *Members) GetMemberInfoAndDepartmentsByNameAndDepartmentId(c app.Ctx, enterpriseId int64, search memberReq.MemberSearch, departmentId []int64) ([]membersModel.MemberDepartment, error) {
	var result []membersModel.MemberDepartment
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.status,members.mobile,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("member_departments.deleted_at=0").
		Where("members.enterprise_id = ?", enterpriseId)
	if search.Name != "" {
		query = query.Where("members.name like ? ", "%"+search.Name+"%")
	}
	if search.Mobile != "" {
		query = query.Where("members.mobile = ? ", search.Mobile)
	}
	if search.Userid != "" {
		query = query.Where("members.userid = ? ", search.Userid)
	}
	if len(departmentId) > 0 {
		query = query.Where("member_departments.department_id in ?", departmentId)
	}
	err := query.Group("member_departments.id desc").
		Find(&result).Error
	return result, err
}

// GetMemberIdsByUserid 根据导购userid模糊匹配导购的memberId
func (d *Members) GetMemberIdsByUserid(c app.Ctx, enterpriseId int64, userid string) (memberIds []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("userid like ?", "%"+userid+"%").
		Pluck("id", &memberIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *Members) GetUserIdsByMemberIDs(c app.Ctx, enterpriseId int64, memberIds []int64) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"id":            memberIds,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Pluck("userid", &result).Error

	return result, err
}

func (d *Members) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *Members) UpdateDataTable(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(data).Error
	return err
}

func (d *Members) GetMemberIdByUserId(c app.Ctx, enterpriseId int64, userid string) (int64, error) {
	var result int64
	err := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"userid":        userid,
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		}).Pluck("id", &result).Error

	return result, err
}

func (d *Members) GetUserIdsById(c app.Ctx, enterpriseId int64, ids []int64) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"id":            ids,
			"status":        1,
			"deleted_at":    0,
		}).Pluck("userid", &result).Error

	return result, err
}

// GetMemberInfoAndDepartmentByMemberIds preload方式根据memberId集合获取导购信息
func (d *Members) GetMemberInfoAndDepartmentByMemberIds(c app.Ctx, enterpriseId int64, ids []int64) (list []membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("id IN ?", ids).
		Unscoped().
		Preload("MemberDepartments.Department", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Find(&list).Error
	return
}

// GetMemberInfoAndDepartmentByUserIds preload方式根据userid集合获取导购信息
func (d *Members) GetMemberInfoAndDepartmentByUserIds(c app.Ctx, enterpriseId int64, userIds []string) (list []membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("userid IN ?", userIds).
		Unscoped().
		Preload("MemberDepartments.Department", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Find(&list).Error
	return
}

// GetAllMemberInfo  preload方式获取导购信息
func (d *Members) GetAllMemberInfo(c app.Ctx, enterpriseId int64) (list []membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberModel{}).
		Where("enterprise_id = ?", enterpriseId).
		Unscoped().
		Preload("MemberDepartments.Department", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Find(&list).Error
	return
}

func (d *Members) GetUnscopedMemberInfoAndDepartmentsByUserId(c app.Ctx, enterpriseId int64, userIds []string) ([]membersModel.MemberDepartment, error) {
	var result []membersModel.MemberDepartment
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("members.enterprise_id,members.id AS members_id,members.name AS members_name,members.userid AS members_userid,departments.id AS department_id,departments.name AS department_name,departments.third_id AS department_third_id").
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("members.enterprise_id = ?", enterpriseId)
	if len(userIds) > 0 {
		query = query.Where(" members.userid IN(?) ", userIds)
	}
	err := query.Group("member_departments.id").Order("member_departments.deleted_at asc").
		Find(&result).Error
	return result, err
}

func (d *Members) GetMemberInfoById(c app.Ctx, enterpriseId int64, id int64) (info membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberModel{}).
		Where("enterprise_id = ? AND id = ?", enterpriseId, id).
		Unscoped().
		Find(&info).Error
	return
}

func (d *Members) GetUnscopedMemberInfoByIds(c app.Ctx, enterpriseId int64, ids []int64) (list []membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberModel{}).
		Where("enterprise_id = ? AND id in ?", enterpriseId, ids).
		Unscoped().
		Find(&list).Error
	return
}

func (d *Members) GetUnscopedLeaveMemberInfoByIds(c app.Ctx, enterpriseId int64, ids []int64) (list []membersModel.MemberModel, err error) {
	err = d.dao.WithContext(c).Model(membersModel.MemberModel{}).
		Where("enterprise_id = ? AND id in ? and status=?", enterpriseId, ids, membersModel.MembersStatusLeave).
		Unscoped().
		Find(&list).Error
	return
}

func (d *Members) GetStoreVisitManagerByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) ([]membersModel.StoreManager, error) {
	var result []membersModel.StoreManager
	err := d.dao.WithContext(c).Table(d.tableName).
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where(map[string]interface{}{
			"departments.enterprise_id":        enterpriseId,
			"departments.deleted_at":           0,
			"member_departments.enterprise_id": enterpriseId,
			"member_departments.deleted_at":    0,
			"members.enterprise_id":            enterpriseId,
			"members.deleted_at":               0,
			"departments.third_id":             thirdIds,
			"members.position":                 StoreVisitManagerPosition,
		}).Select("members.id as member_id,members.name as member_name,members.position as position,departments.third_id,departments.name as store_code").
		Scan(&result).Error
	return result, err
}

func (d *Members) GetMemberInfoByNameOrEmail(c app.Ctx, enterpriseId int64, str string) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ?", enterpriseId).
		Where("deleted_at = 0").
		Where("(name like ? or email like ?)", "%"+str+"%", "%"+str+"%").
		Preload("MemberDepartments.Department").
		Order("created_at desc").
		Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

// GetMemberForShenCe   获取需要推送给神策的用户
func (d *Members) GetMemberForShenCe(c app.Ctx, enterpriseId, startTime, endTime int64) ([]membersResp.MemberForShenCeResp, error) {
	var result []membersResp.MemberForShenCeResp
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("members.id,members.userid,members.name,members.mobile,members.position,members.type,members.gender,members.avatar,svm.dept_name as sv_department,svm.is_store_visit_manager as sv_manager").
		Joins(fmt.Sprintf("LEFT JOIN micro_scrm_enterprise_%d.store_visit_dept_member svm ON svm.member_id = members.id", enterpriseId)).
		Where("members.enterprise_id = ? and ((members.updated_at >= ? and members.updated_at <= ?) or (svm.updated_time >= ? and svm.updated_time <= ?))",
			enterpriseId, startTime, endTime, helper.GetFormatDateByParams(startTime), helper.GetFormatDateByParams(endTime)).
		Where("members.deleted_at = 0").
		Scan(&result).Error
	return result, err
}

// GetMemberIdsByThirdIds  根据third_id获取导购id集合
func (d *Members) GetMemberIdsByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) ([]int64, error) {
	var result []int64
	query := d.dao.WithContext(c).Table(d.tableName).
		Joins("LEFT JOIN member_departments ON member_departments.member_id = members.id").
		Joins("LEFT JOIN departments ON departments.id = member_departments.department_id").
		Where("departments.third_id in (?)", thirdIds).
		Where(map[string]interface{}{
			"members.enterprise_id":         enterpriseId,
			"member_departments.deleted_at": 0,
			"departments.deleted_at":        0,
			"members.deleted_at":            0,
		})
	err := query.Group("members.id").
		Pluck("members.id", &result).Error
	return result, err

}

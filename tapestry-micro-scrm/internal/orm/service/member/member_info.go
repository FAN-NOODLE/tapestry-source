package member

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/service"
)

// QueryAllMemberInfo 查询员工的信息，所有符合条件的
func (m *MemberService) QueryAllMemberInfo(param micro_scrm.MemberModel) (result []micro_scrm.MemberModel, err error) {
	err = m.Scopes(service.ScopeCond_EQDelete).
		Where(&param).Find(&result).Error

	return
}

// QueryMemberInfo 查询员工信息
func (m *MemberService) QueryMemberInfo(param micro_scrm.MemberModel) (result micro_scrm.MemberModel, err error) {
	err = m.Scopes(service.ScopeCond_EQDelete).Where(&param).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

// GetMemberInfoByUserID 根据企业微信员工的userID，获取对应到本地数据库表member中的信息
func (m *MemberService) GetMemberInfoByUserID(enterpriseID int64, userID string) (data micro_scrm.MemberModel) {
	m.Model(&data).
		Scopes(service.ScopeCond_EQDelete).
		Where("userid = ?", userID).
		Where("enterprise_id = ?", enterpriseID).
		First(&data)

	return
}

// GetUserIdByMobile 根据企业微信员工的手机号获取userID
func (m *MemberService) GetUserIdByMobile(enterpriseID int64, mobile string) (data micro_scrm.MemberModel) {
	m.Model(&data).
		Scopes(service.ScopeCond_EQDelete).
		Where("mobile = ?", mobile).
		Where("enterprise_id = ?", enterpriseID).
		Where("status=1").
		First(&data)

	return
}

type Department struct {
	Userid string `json:"userid"`
	DpId   int64  `json:"dp_id"`
	DpName string `json:"dp_name"`
}

// GetMemberDepartment 查询员工对应的部门
func (m *MemberService) GetMemberDepartment(userId string, enterpriseID int64) (department Department) {
	m.Model(&micro_scrm.MemberModel{}).
		Select("members.userid,members.name,departments.name as dp_name,departments.id as dp_id").
		Joins("left join member_departments on member_departments.member_id = members.id").
		Joins("left join departments on departments.id = member_departments.department_id").
		Where("members.userid", userId).
		Where("members.enterprise_id", enterpriseID).
		Where("member_departments.deleted_at", 0).
		First(&department)
	return
}

func (m MemberService) GetMemberLsit(enterpriseID int64) (memberList []micro_scrm.MemberModel, err error) {
	err = m.Model(&micro_scrm.MemberModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("deleted_at = ?", 0).
		Find(&memberList).Error
	return
}

// GetMemberListPage 分页查询member信息
func (m MemberService) GetMemberListPage(enterpriseID int64, offset, limit int) (memberList []micro_scrm.MemberModel, err error) {
	err = m.Model(&micro_scrm.MemberModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("deleted_at = ?", 0).Limit(limit).Offset(offset).
		Find(&memberList).Error
	return
}

// GetMemberCount 获取member总数
func (m MemberService) GetMemberCount(enterpriseID int64) (count int64, err error) {
	err = m.Model(&micro_scrm.MemberModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("deleted_at = ?", 0).
		Count(&count).Error
	return
}

// GetMembersByUserId 根据userid 获取member_id
func (m MemberService) GetMembersByUserId(enterpriseID int64, userIds []string) (memberList []*micro_scrm.MemberModel, err error) {
	err = m.Model(&micro_scrm.MemberModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("userid in ? ", userIds).
		Where("deleted_at = ?", 0).
		Group("userid").
		Find(&memberList).Error
	return memberList, err
}

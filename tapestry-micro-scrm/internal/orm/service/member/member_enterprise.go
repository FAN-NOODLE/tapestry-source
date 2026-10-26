package member

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// 品牌分库表的员工信息数据

func (m *MemberService) InsertEnterpriseMemberEnterprise(enterpriseID int64, data micro_scrm_enterprise.MembersModel) (err error) {
	err = m.Scopes(model.BranchTable(&micro_scrm_enterprise.MembersModel{}, enterpriseID)).Create(&data).Error
	return
}

// IsExistEmp_Enterprise 查询对应的员工是否存在
func (m *MemberService) IsExistEmpEnterprise(enterpriseID int64, userID string) (isExist bool, data micro_scrm_enterprise.MembersModel, err error) {
	if err = m.Scopes(model.BranchTable(&micro_scrm_enterprise.MembersModel{}, enterpriseID)).
		Where("userid = ?", userID).Take(&data).Error; err == nil {
		isExist = true
	}

	return
}

// UpdateEmp_Enterprise 更新数据
func (m *MemberService) UpdateEmpEnterprise(enterpriseID int64, data micro_scrm_enterprise.MembersModel) error {
	return m.Scopes(model.BranchTable(&micro_scrm_enterprise.MembersModel{}, enterpriseID)).Save(data).Error
}

// TakeEmp_Enterprise 获取企业员工信息
func (m *MemberService) TakeEmpEnterprise(enterpriseID int64, param micro_scrm_enterprise.MembersModel) (data micro_scrm_enterprise.MembersModel, err error) {
	err = m.Scopes(model.BranchTable(&micro_scrm_enterprise.MembersModel{}, enterpriseID)).
		Where(&param).
		Take(&data).Error

	return
}

/**
*查询BA
 */
func (m *MemberService) FindBas(enterpriseID int64, BAUserIds []string) (data []int64, err error) {
	err = m.Model(&micro_scrm_enterprise.MembersModel{}).Select("id").
		Where("enterprise_id=? AND deleted_at=?", enterpriseID, 0).
		Where("userid IN (?)", BAUserIds).Find(&data).Error
	return
}

func (m *MemberService) FindBainfos(enterpriseID int64, userids []string) (data []*micro_scrm_enterprise.MemberInfo, err error) {
	err = m.Model(&micro_scrm_enterprise.MembersModel{}).
		Select("members.id menber_id,members.userid user_id").
		Where("members.enterprise_id=? AND members.userid IN (?) AND members.deleted_at=?", enterpriseID, userids, 0).
		Find(&data).Error
	return
}
func (m *MemberService) GetMemenbersBa(enterpriseID int64, baUserIds []string, deptId int64, userName string, page, size int64) (bas []DeptMemberList, err error) {
	query := m.Model(&micro_scrm.MemberModel{}).
		Select("members.id AS member_id,members.userid AS user_id,members.name AS user_name,members.`status` AS mem_subscribe_status,d.`name` AS dept_name,mobile").
		Joins("LEFT JOIN member_departments md ON members.id = md.member_id AND md.deleted_at=? AND md.enterprise_id=?", 0, enterpriseID).
		Joins("LEFT JOIN departments d ON md.department_id=d.id AND d.deleted_at = ? AND d.enterprise_id=?", 0, enterpriseID).
		Where("members.enterprise_id=? AND members.deleted_at=? AND md.id IS NOT NULL AND d.id IS NOT NULL AND members.userid IN (?)", enterpriseID, 0, baUserIds)
	if deptId > 0 {
		query.Where("d.id=?", deptId)
	}
	if userName != "" {
		query.Where("members.name  LIKE ?", "%"+userName+"%")
	}
	err = query.Group("members.id, members.userid, members.name, members.status, d.name").Order("members.id DESC").Limit(int(size)).Offset(int((page - 1) * size)).
		Find(&bas).Error
	return
}
func (m *MemberService) GetMemenbersBaTotal(enterpriseID int64, baUserIds []string, deptId int64, userName string) (total int64, err error) {
	var (
		params   = []interface{}{enterpriseID, enterpriseID, enterpriseID}
		w        = ""
		userWher = ""
	)
	for k, value := range baUserIds {
		if k == 0 {
			userWher = "?"
		} else {
			userWher += ",?"
		}
		params = append(params, value)
	}
	if deptId > 0 {
		w += " AND d.id=?"
		params = append(params, deptId)
	}
	if userName != "" {
		w += " AND members.name  LIKE ? "
		params = append(params, "%"+userName+"%")
	}
	sql := `
	SELECT t.member_id FROM (
		SELECT members.id AS member_id,members.userid AS user_id,members.name AS user_name,members.status AS mem_subscribe_status FROM members LEFT JOIN member_departments md 
    	ON members.id = md.member_id AND md.deleted_at = 0 AND md.enterprise_id = ?
	LEFT JOIN departments d ON md.department_id = d.id AND d.deleted_at = 0  AND d.enterprise_id = ?
	 WHERE members.enterprise_id = ? AND members.deleted_at = 0 AND md.id IS NOT NULL AND d.id IS NOT NULL
	 AND  members.userid IN (` + userWher + `)  ` + w + `
	 GROUP BY members.id,members.userid,members.name,members.status,d.name
	) t   
	`
	tabsql := "SELECT COUNT(1) FROM ( " + sql + ") t1"
	err = m.Raw(tabsql, params...).Count(&total).Error
	return
}

func (m *MemberService) GetMemenbersBaNoDept(enterpriseID int64, baUserIds []string, deptId int64, userName string) (bas []DeptMemberList, err error) {
	query := m.Model(&micro_scrm.MemberModel{}).
		Select("members.id AS member_id,members.userid AS user_id,members.name AS user_name,members.`status` AS mem_subscribe_status,mobile").
		Where("members.enterprise_id=? AND members.deleted_at=? AND members.userid IN (?)", enterpriseID, 0, baUserIds).
		Where("members.id NOT IN (?)", m.Model(&micro_scrm.MemberDepartmentsModel{}).Select("DISTINCT member_id").Where("enterprise_id=? AND deleted_at=? ", enterpriseID, 0))
	if userName != "" {
		query.Where("members.name  LIKE ?", "%"+userName+"%")
	}
	err = query.Group("members.id, members.userid, members.name, members.status").
		Find(&bas).Error
	return
}

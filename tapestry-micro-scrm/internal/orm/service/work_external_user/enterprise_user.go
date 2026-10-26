package work_external_user

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// UpdateExternalUserUserGroup 更新外部联系人的分组ID（主要为对于，新增客户分组时，把新增的客户分组同步到enterprise_user.groups）
func (w *WorkExternalUserService) UpdateExternalUserUserGroup(enterpriseID, userGroupID int64, externalUserIDList []int64) (err error) {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id IN ?", externalUserIDList).
		Update("groups", gorm.Expr(fmt.Sprintf("IF(IFNULL(groups, \"\") = \"\", JSON_ARRAY(%d),JSON_ARRAY_APPEND(groups, '$',%d))", userGroupID, userGroupID))).Error
}

// 更新外部联系人的分组ID（主要对于，更新外部联系人的enterprise_user.groups。整个替换）
func (w *WorkExternalUserService) ReplaceExternalUserGroups(enterpriseID, externalUserID int64, customerGroupIDList []int64) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id = ?", func() int64 {
			data := micro_scrm_enterprise.WxUserCorpExternalModel{}
			w.Scopes(model.BranchTable(&data, enterpriseID)).Where("id = ?", externalUserID).First(&data)
			return data.Euid
		}()).
		Update("groups", gorm.Expr("JSON_ARRAY(?)", customerGroupIDList)).Error
}

// 跟新外部联系人的分组ID （主要对于，更新外部联系人的enterprise_user.groups。单个替换）
func (w *WorkExternalUserService) UpdateExternalUserGroupsOne(enterpriseID, ID int64, customerGroupIDList []int64) (err error) {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id = ?", ID).
		Update("groups", func(idList []int64) clause.Expr {
			if 0 == len(idList) {
				return gorm.Expr("?", "")
			} else {
				return gorm.Expr("JSON_ARRAY(?)", idList)
			}
		}(customerGroupIDList)).Error
}

// UpdateExternalUserTag 更新外部联系人的标签
func (w *WorkExternalUserService) UpdateExternalUserTag(enterpriseID, externalUserID int64, tagIDList []int64) (err error) {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id = ?", func() int64 {
			data := micro_scrm_enterprise.WxUserCorpExternalModel{}
			w.Scopes(model.BranchTable(&data, enterpriseID)).Where("id = ?", externalUserID).First(&data)
			return data.Euid
		}()).
		Update("tags", func(idList []int64) clause.Expr {
			if 0 == len(idList) {
				return gorm.Expr("?", "")
			} else {
				return gorm.Expr("JSON_ARRAY(?)", idList)
			}
		}(tagIDList)).Error
}

// 根据客户分组获取外部联系人信息
func (w *WorkExternalUserService) GetExternalUserByCustomerGroupID(enterpriseID, customerGroupID int64) (dataList []micro_scrm_enterprise.EnterpriseUserModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("groups LIKE ?", fmt.Sprintf("%s%d%s", "%", customerGroupID, "%")).
		Find(&dataList).Error

	return
}

func (w *WorkExternalUserService) UpdateEnterpriseUserByunionId(enterpriseID, euid int64, memberCode, adminBa, gradeCodeName string) (err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id=?", euid).Updates(map[string]interface{}{
		"grade":      gradeCodeName,
		"union_code": memberCode,
		"ba_code":    adminBa,
	}).Error
	return
}

// UpdateEnterpriseUserById 更新字段
func (w *WorkExternalUserService) UpdateEnterpriseUserById(enterpriseID int64, ID int64, unionCode, birthday, userName, mobile string) (err error) {
	userMap := make(map[string]interface{})
	userMap["union_code"] = unionCode
	//userMap["name"] = userName
	//userMap["mobile"] = mobile
	if birthday != "" {
		userMap["birthday"] = birthday
	}
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
		Where("id = ?", ID).
		Updates(userMap).Error
	return
}

package work_external_user

import (
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/db/gormx"
)

// 根据用户的关系ID，获取用户的uid
func (w *WorkExternalUserService) GetExternalUserIDByRelationID(enterpriseID, relationID int64) (externalUserID int64, err error) {
	data := work_external_user.WorkExternalUserRelationModel{}
	err = w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "id": relationID}).
		First(&data).Error
	externalUserID = data.WorkExternalUserID

	return
}

// GetCustomerRelationBA 获取和顾客相关的BA
func (w *WorkExternalUserService) GetCustomerRelationBA(enterpriseID, workExternalUserID int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).Model(&work_external_user.WorkExternalUserRelationModel{}).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "work_external_user_id": workExternalUserID}).
		Where("relation_status = ?", 0).
		Preload("MemberInfo").Find(&data).Error

	return
}

// GetAllRelation 获取和顾客相关的BA
func (w *WorkExternalUserService) GetAllRelation(enterpriseID int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).Model(&work_external_user.WorkExternalUserRelationModel{}).
		Where(map[string]interface{}{"enterprise_id": enterpriseID}).
		Where("relation_status = ?", 0).
		Find(&data).Error

	return
}

// 根据外部联系人的unionid和员工userid是否有绑定关系
func (w *WorkExternalUserService) IsExternalRelation(enterpriseID, uId, mId int64) (rowsAffected int64, err error) {
	data := work_external_user.WorkExternalUserRelationModel{}
	result := w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "member_id": mId, "work_external_user_id": uId}).
		First(&data)
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

// 根据外部联系人的unionid和员工userid是否有绑定关系
func (w *WorkExternalUserService) IsExternalRelationBak(unionId, userId string, enterpriseID int64) int64 {
	var rid, uid, mid int64
	db := gormx.GetInstance()
	sql := fmt.Sprintf("select id from micro_scrm_enterprise_%d.wx_user_corp_external where unionid=?", enterpriseID)
	db.Raw(sql, unionId).Scan(&uid)
	db.Raw("select id from members where userid=? and enterprise_id=?", userId, enterpriseID).Scan(&mid)
	db.Raw("select id from work_external_user_relation where work_external_user_id=? and"+
		" member_id=? and relation_status=0 and enterprise_id=?", uid, mid, enterpriseID).Scan(&rid)
	if rid != 0 {
		rid = 1
	}
	return rid
}

//GetExternalUserCountByMemberId 根据memberId 获取关系外部联系人数
func (w *WorkExternalUserService) GetExternalUserCountByMemberId(enterpriseID,mId int64) (count int64, err error) {
	data := work_external_user.WorkExternalUserRelationModel{}
	result := w.Scopes(service.ScopeCond_EQDelete).Model(&data).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "member_id": mId }).
		Count(&count)
	err = result.Error

	return
}


// GetExternalUserInfoByMemberId 根据memberId 获取关系信息 默认取1000个数据
func (w *WorkExternalUserService) GetExternalUserInfoByMemberId(enterpriseID, mId int64,offset,limit int ) (data []*work_external_user.WorkExternalUserRelationModel, err error) {
	//data := work_external_user.WorkExternalUserRelationModel{}
	result := w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "member_id": mId }).
		Offset(offset).Limit(limit).
		Find(&data)
	err = result.Error

	return
}
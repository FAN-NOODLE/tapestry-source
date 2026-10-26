package data_compass

import (
	"fmt"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (d DataCompassService) GetNewGoodFriendsSingal(memberId int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	err = d.Where(service.ScopeCond_EQDelete).
		Where("member_id = ?", memberId).
		Find(&data).Error
	return
}

// GetNewGoodFriendsCount 获取新增好友数据
func (d DataCompassService) GetNewGoodFriendsCount(enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (count int64, err error) {
	scope := d.Model(&work_external_user.WorkExternalUserRelationModel{})
	if startTime > 0 && endTime > 0 && endTime > startTime {
		scope.Where("createtime > ?", startTime).
			Where("createtime < ?", endTime)
	}
	builder := scope.Where("deleted_at = ?", 0)
	if searchAll != 1 {
		builder.Where("member_id in (?)", memberId)
	}
	err = builder.Where("unionid != ''").
		Where("enterprise_id = ?", enterpriseId).
		Distinct("work_external_user_id").Count(&count).Error

	return
}

func (d DataCompassService) GetNewGoodFriendsAllBymemberId(enterpriseID, memberId int64, startTime int64, endTime int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "

	scope := d.Model(&work_external_user.WorkExternalUserRelationModel{}).
		Joins("LEFT JOIN " + userexternalTable + " b ON work_external_user_relation.work_external_user_id=b.id ")
	if startTime > 0 && endTime > 0 && endTime > startTime {
		scope.Where("work_external_user_relation.created_at >= ?", startTime).
			Where("work_external_user_relation.created_at <= ?", endTime)
	}
	err = scope.Where("work_external_user_relation.deleted_at = ? AND b.id is NOT NULL", 0).
		Where("work_external_user_relation.member_id =? ", memberId).
		Distinct("work_external_user_relation.work_external_user_id").Find(&data).Error
	return
}

//获取累计好友数据
func (d DataCompassService) GetGoodFriendsAll(enterpriseId int64, memberId []string, endTime, searchAll int64) (total int64, err error) {
	scope := d.Model(&work_external_user.WorkExternalUserRelationModel{})
	if endTime > 0 {
		scope.Where("work_external_user_relation.createtime < ?", endTime)
	}
	builder := scope.Where("work_external_user_relation.deleted_at = ?", 0).
		Where("enterprise_id = ?", enterpriseId).
		Where("work_external_user_relation.relation_status = ?", 0).
		Where("work_external_user_relation.unionid <> ''")
	if searchAll != 1 {
		err = builder.Where("work_external_user_relation.member_id in (?)", memberId).Count(&total).Error
	} else {
		err = builder.Count(&total).Error
	}

	return
}

//会员绑定数
func (d DataCompassService) GetFriendsVIP(enterpriseId int64, startTime int64, endTime int64, memberIds []string) (data int64, err error) {
	//scope := d.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseId)).
	//	Model(&micro_scrm_enterprise.WxUserCorpExternalModel{})

	//err = scope.Where("deleted_at = ?", 0).
	//	Where("id in (?)", workExternalIds).
	//	Where("union_code <> ?", "").
	//	Count(&data).Error
	enterpriseIdStr := fmt.Sprintf("%d", enterpriseId)
	sql := "SELECT  count(1) as count FROM  work_external_user_relation as u force index(good_friend_index) , " +
		"micro_scrm_enterprise_" + enterpriseIdStr + ".wx_user_corp_external as wx WHERE  " +
		"u.work_external_user_id = wx.id  and u.createtime > ? and u.createtime < ?   AND u.deleted_at = 0  " +
		"AND u.member_id IN (?)   AND u.enterprise_id = ?  and wx.union_code <>'' AND wx.type = ? "
	err = d.Raw(sql, startTime, endTime, memberIds, enterpriseId, micro_scrm_enterprise.WECHATt_USER_TYPE).Scan(&data).Error
	return
}

type Res struct {
	Count int64 `gorm:"column:count"`
}

//历史会员绑定数
func (d DataCompassService) GetFriendsVIPAll(enterpriseId int64, endTime int64, memberIds []string, searchAll int64) (data int64, err error) {
	//scope := d.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseId)).
	//	Model(&micro_scrm_enterprise.WxUserCorpExternalModel{})

	//err = scope.Where("deleted_at = ?", 0).
	//	Where("id in (?)", workExternalIds).
	//	Where("union_code <> ?", "").
	//	Count(&data).Error
	enterpriseIdStr := fmt.Sprintf("%d", enterpriseId)
	if searchAll != 1 {
		sql := "SELECT  count(DISTINCT work_external_user_id) as count FROM  " +
			"work_external_user_relation as u , " +
			"micro_scrm_enterprise_" + enterpriseIdStr + ".wx_user_corp_external as wx WHERE  " +
			"u.work_external_user_id = wx.id  and u.createtime < ? AND u.deleted_at = 0 AND  " +
			"u.relation_status = 0 AND u.member_id IN (?)   AND u.enterprise_id = ?  and wx.union_code <>'' AND wx.type = ?"
		err = d.Raw(sql, endTime, memberIds, enterpriseId, micro_scrm_enterprise.WECHATt_USER_TYPE).Scan(&data).Error
	} else {
		sql := "SELECT  count(DISTINCT work_external_user_id) as count FROM  " +
			"work_external_user_relation as u , " +
			"micro_scrm_enterprise_" + enterpriseIdStr + ".wx_user_corp_external as wx WHERE  " +
			"u.work_external_user_id = wx.id  and u.createtime < ? AND u.deleted_at = 0 AND  " +
			"u.relation_status = 0 AND u.enterprise_id = ?  and wx.union_code <>'' AND wx.type = ?"
		err = d.Raw(sql, endTime, enterpriseId, micro_scrm_enterprise.WECHATt_USER_TYPE).Scan(&data).Error
	}

	return
}

//获取指定id的unionIds
func (d DataCompassService) GetMemberId(enterpriseId int64, memberUserId []string) (data []micro_scrm_enterprise.MembersModel, err error) {
	err = d.Model(&micro_scrm_enterprise.MembersModel{}).
		Where("deleted_at = ?", 0).
		Where("enterprise_id = ?", enterpriseId).
		Where("userid in (?)", memberUserId).
		Find(&data).Error
	return
}

//获取指定id的unionIds
func (d DataCompassService) GeUnionIds(enterpriseId int64, Ids []int64) (data []micro_scrm_enterprise.WxUserCorpExternalModel, err error) {

	err = d.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseId)).
		Where("deleted_at = ?", 0).
		Where("id in (?)", Ids).
		Distinct("unionid").Find(&data).Error
	return
}

//是否入群
func (d DataCompassService) GetJoinCroupCount(enterpriseId int64, baGroupIds []int64, memberIds []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	//err = d.Model(&micro_scrm.WorkGroupChatMember{}).
	//	Where("deleted_at = ?", 0).
	//	Where("enterprise_id = ?", enterpriseId).
	//	Where("group_chat_id in (?)", baGroupIds).
	//	Where("unionid in (?)", unionids).
	//	Where("status = ?", 1).
	//	Distinct("unionid").Count(&data).Error
	enterpriseIdStr := fmt.Sprintf("%d", enterpriseId)
	if searchAll != 1 {
		sql := "SELECT  COUNT( DISTINCT work_external_user_id ) as count " +
			"FROM `work_external_user_relation` UR " +
			"WHERE  UR.enterprise_id = " + enterpriseIdStr + "  " +
			"AND UR.deleted_at = 0  AND UR.member_id IN (?) AND UR.relation_status = 0  AND " +
			"UR.createtime >= ? AND UR.createtime <= ? AND UR.unionid IN " +
			"(SELECT DISTINCT unionid FROM (SELECT unionid, status, Max( created_at ) " +
			"FROM `work_group_chat_member` WHERE enterprise_id = " + enterpriseIdStr + " AND type = 2  " +
			"AND deleted_at = 0  AND group_chat_id IN (?) AND unionid <> '' " +
			"GROUP BY `unionid`  HAVING status = 1) F )"
		err = d.Raw(sql, memberIds, startTime, endTime, baGroupIds).Find(&data).Error
	} else {
		sql := "SELECT  COUNT( DISTINCT work_external_user_id ) as count " +
			"FROM `work_external_user_relation` UR " +
			"WHERE  UR.enterprise_id = " + enterpriseIdStr + "  " +
			"AND UR.deleted_at = 0  AND UR.relation_status = 0  AND " +
			"UR.createtime >= ? AND UR.createtime <= ? AND UR.unionid <> ''AND UR.unionid IN " +
			"(SELECT DISTINCT unionid FROM (SELECT unionid, status, Max( created_at ) " +
			"FROM `work_group_chat_member` WHERE enterprise_id = " + enterpriseIdStr + " AND type = 2  " +
			"AND deleted_at = 0  AND group_chat_id IN (?) AND unionid <> '' " +
			"GROUP BY `unionid`  HAVING status = 1) F )"
		err = d.Raw(sql, startTime, endTime, baGroupIds).Find(&data).Error
	}

	return
}

//是否入群
func (d DataCompassService) GetJoinCroupCountAll(enterpriseId int64, baGroupIds []int64, memberIds []string, endTime, searchAll int64) (data int64, err error) {
	//err = d.Model(&micro_scrm.WorkGroupChatMember{}).
	//	Where("deleted_at = ?", 0).
	//	Where("enterprise_id = ?", enterpriseId).
	//	Where("group_chat_id in (?)", baGroupIds).
	//	Where("unionid in (?)", unionids).
	//	Where("status = ?", 1).
	//	Distinct("unionid").Count(&data).Error
	enterpriseIdStr := fmt.Sprintf("%d", enterpriseId)
	if searchAll != 1 {
		sql := "SELECT COUNT(1) as count FROM ( SELECT  DISTINCT work_external_user_id FROM `work_external_user_relation` UR " +
			"WHERE  UR.enterprise_id = " + enterpriseIdStr + "  " +
			"AND UR.deleted_at = 0  AND UR.member_id IN (?) AND UR.relation_status = 0  AND UR.createtime <= ? " +
			"AND UR.unionid <> '' AND UR.unionid IN (SELECT DISTINCT unionid FROM " +
			"(SELECT unionid, status, Max( created_at ) FROM `work_group_chat_member` WHERE " +
			"enterprise_id = " + enterpriseIdStr + " AND type = 2  " +
			"AND deleted_at = 0  AND group_chat_id IN (?) AND unionid <> '' AND status = 1 " +
			"GROUP BY `unionid` ) F ) )as a"
		err = d.Raw(sql, memberIds, endTime, baGroupIds).Find(&data).Error
	} else {
		sql := "SELECT COUNT(1) as count FROM ( SELECT  DISTINCT work_external_user_id FROM `work_external_user_relation` UR " +
			"WHERE  UR.enterprise_id = " + enterpriseIdStr + "  " +
			"AND UR.deleted_at = 0  AND UR.relation_status = 0  AND UR.createtime <= ? " +
			"AND UR.unionid <> '' AND UR.unionid IN (SELECT DISTINCT unionid FROM " +
			"(SELECT unionid, status, Max( created_at ) FROM `work_group_chat_member` WHERE " +
			"enterprise_id = " + enterpriseIdStr + " AND type = 2  " +
			"AND deleted_at = 0  AND group_chat_id IN (?) AND unionid <> '' AND status = 1 " +
			"GROUP BY `unionid` ) F ) ) as a"
		err = d.Raw(sql, endTime, baGroupIds).Find(&data).Error
	}

	return
}

// 根据unionIDS 查询 是否入群
func (d DataCompassService) GetJoinCroupByUnoinId(enterpriseId int64, unionids []string) (data int64, err error) {
	err = d.Model(&micro_scrm.WorkGroupChatMember{}).
		Where("deleted_at = ? AND enterprise_id=?", 0, enterpriseId).
		Where("unionid in (?)", unionids).
		Where("status = ?", 1).
		Distinct("unionid").Count(&data).Error
	return
}

//获取流失好友数据
func (d DataCompassService) GetDrainFriends(enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (data int64, err error) {

	scope := d.Model(&work_external_user.WorkExternalUserRelationModel{})
	if endTime > 0 && startTime > 0 && endTime > startTime {
		whereStr := fmt.Sprintf("((work_external_user_relation.member_deleted_time>0 and work_external_user_relation.member_deleted_time <= %d) or ("+
			"work_external_user_relation.external_deleted_time>0 and work_external_user_relation.external_deleted_time <= %d))", endTime, endTime)
		scope.Where(whereStr)
	}
	if searchAll != 1 {
		scope.Where("work_external_user_relation.member_id in (?)", memberId)
	}
	err = scope.Where("work_external_user_relation.deleted_at = ?", 0).
		Where("work_external_user_relation.enterprise_id = ?", enterpriseId).
		Where("work_external_user_relation.unionid  <> '' ").
		Distinct("work_external_user_relation.unionid").Count(&data).Error
	return
}

//获取流失好友数据
func (d DataCompassService) GetDrainFriendsBymemberId(enterpriseID, memberId int64, endTime int64) (data int64, err error) {
	scope := d.Model(&work_external_user.WorkExternalUserRelationModel{})
	if endTime > 0 {
		scope.Where("created_at < ?", endTime)
	}
	err = scope.Where("deleted_at = ? AND enterprise_id=?", 0, enterpriseID).
		Where("member_id in (?)", memberId).
		Where("relation_status  <> 0").
		Distinct("work_external_user_id").Count(&data).Error
	return
}

/**
*根据导购查询累计好友数
 */

func (d DataCompassService) GetLeijiFriendsAllBymemberId(enterpriseID, memberId int64, endTime int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "

	scope := d.Model(&work_external_user.WorkExternalUserRelationLogModel{}).
		Joins("LEFT JOIN " + userexternalTable + " b ON work_external_user_relation_log.work_external_user_id=b.id ")
	if endTime > 0 {
		scope.Where("work_external_user_relation_log.created_at <= ?", endTime)
	}
	err = scope.Where("work_external_user_relation_log.deleted_at = ? AND b.id is NOT NULL", 0).
		Where("work_external_user_relation_log.member_id = ?", memberId).
		Where("work_external_user_relation_log.relation_status = ?", 0).
		Distinct("work_external_user_relation_log.work_external_user_id").Find(&data).Error
	return
}

/**
*好友入群数量
 */

func (d DataCompassService) GetHaoyouruqunBymemberId(enterpriseID, memberId int64, startTime, endTime int64, groupIds []int64) (data []work_external_user.WorkExternalUserRelationModel, err error) {
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "

	scope := d.Model(&work_external_user.WorkExternalUserRelationLogModel{}).
		Joins("LEFT JOIN "+userexternalTable+" b ON work_external_user_relation_log.work_external_user_id=b.id ").
		Joins("LEFT JOIN  work_group_chat_member c ON b.unionid=c.unionid  AND c.enterprise_id=? AND c.status=?", enterpriseID, 1).
		Where("work_external_user_relation_log.created_at <= ?", endTime).
		Where("c.group_chat_id IN (?)", groupIds)
	err = scope.Where("work_external_user_relation_log.deleted_at = ? AND c.id is NOT NULL ", 0).
		Where("work_external_user_relation_log.member_id = ?", memberId).
		Distinct("work_external_user_relation_log.work_external_user_id").Find(&data).Error
	return
}

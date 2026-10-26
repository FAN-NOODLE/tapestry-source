package data_compass

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// GetBAGroupInfo 获取ba的群信息
func (b *DataCompassService) GetBAGroupInfo(enterpriseID, startTime, endTime int64, baInfo []string) (result []model.GroupUserInfo, err error) {
	var groupChatMember group_chat.WorkGroupChatMemberModel
	var groupChat group_chat.WorkGroupChatModel
	err = b.Model(&groupChatMember).
		Select("userid AS user_id,group_chat_id").
		Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
		Where("status = ? ", 1).
		Where("userid IN (?)", baInfo).
		Where("group_chat_id IN (?)", b.Model(&groupChat).
			Select("id").
			Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
			Where(" type = ?", 1).
			Where("status = ?", 1).
			Where("create_time between ? AND ?", startTime, endTime)).
		Group("group_chat_id").
		Find(&result).Error
	return
}

// GetGroupCustomCount 获取群成员统计
func (b *DataCompassService) GetGroupCustomCount(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64, baInfo []string, searchAll int64) (count int64, err error) {
	// 调整为clickhouse查询
	// db.ClickHouse().Get(ctx)
	// var result []struct {
	//	Count int64 `json:"count" db:"count"`
	// }
	// ckSql := "SELECT count(DISTINCT union_id) as count FROM micro_scrm.group_chat_member_log_local WHERE enterprise_id= ? and type= ? and status = ? and opera_time<= ? and group_chat_id in (?) and union_id != ''"
	// err = db.ClickHouse().Get(ctx, &result, ckSql, enterpriseID, micro_scrm.WGCML_MEMBER_TYPE, micro_scrm.WGCML_OPERATION_STATUS, endTime, groupChatIngo)
	// if err != nil {
	//	return 0, err
	// }
	// return result[0].Count, nil

	var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	builder := b.Model(&groupChatMemberLog).Debug().
		Select("count( DISTINCT work_group_chat_member_log.union_id)").
		Joins("LEFT JOIN work_external_user_relation ON work_external_user_relation.unionid=work_group_chat_member_log.union_id").
		Where("work_group_chat_member_log.enterprise_id = ? AND work_group_chat_member_log.deleted_at = 0", enterpriseID)
	if searchAll != 1 {
		builder.Where("work_external_user_relation.user_id IN(?)", baInfo)
	}
	err = builder.Where("work_group_chat_member_log.member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
		Where("work_group_chat_member_log.operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
		Where("work_group_chat_member_log.created_at <= ?", endTime).
		Where("work_group_chat_member_log.group_chat_id IN (?)", groupChatIngo).
		Where("work_group_chat_member_log.union_id != '' ").
		Where("work_external_user_relation.enterprise_id = ? AND work_external_user_relation.deleted_at = 0", enterpriseID).
		Where("work_external_user_relation.relation_status = 0 ").
		Find(&count).Error

	return count, err
}

// GetGroupCustomerNumberCount 获取群成员统计
func (b *DataCompassService) GetGroupCustomerNumberCount(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64) (count int64, err error) {
	// 调整为clickhouse查询
	// db.ClickHouse().Get(ctx)
	// var result []struct {
	//	Count int64 `json:"count" db:"count"`
	// }
	// ckSql := "SELECT count(DISTINCT union_id) as count FROM micro_scrm.group_chat_member_log_local WHERE enterprise_id= ? and type= ? and status = ? and opera_time<= ? and group_chat_id in (?) and union_id != ''"
	// err = db.ClickHouse().Get(ctx, &result, ckSql, enterpriseID, micro_scrm.WGCML_MEMBER_TYPE, micro_scrm.WGCML_OPERATION_STATUS, endTime, groupChatIngo)
	// if err != nil {
	//	return 0, err
	// }
	// return result[0].Count, nil

	var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	err = b.Model(&groupChatMemberLog).Debug().
		// Select("count( DISTINCT work_group_chat_member_log.union_id)").
		Where("work_group_chat_member_log.enterprise_id = ? AND work_group_chat_member_log.deleted_at = 0", enterpriseID).
		Where("work_group_chat_member_log.member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
		Where("work_group_chat_member_log.operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
		Where("work_group_chat_member_log.created_at <= ?", endTime).
		Where("work_group_chat_member_log.group_chat_id IN (?)", groupChatIngo).
		Where("work_group_chat_member_log.union_id != '' ").
		Distinct("work_group_chat_member_log.union_id").
		Count(&count).Error

	return count, err
}

// GetContactPerson 获取群成员信息
func (b *DataCompassService) GetContactPerson(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64) (result []model.UnionInfo, err error) {

	// var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	// err = b.Model(&groupChatMemberLog).
	//	Select("union_id").
	//	//Select("union_id,operation_status,Max(created_at)").
	//	Where("enterprise_id = ? AND deleted_at = 0", enterpriseID).
	//	Where("member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
	//	Where("operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
	//	Where("created_at <= ?", endTime).
	//	//Where("external_userid LIKE ?", micro_scrm.EXTERNAL_CONTACTS_STRING).
	//	Where("group_chat_id IN (?)", groupChatIngo).
	//	Where("union_id != '' ").
	//	Group("union_id").
	//	//Having("operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
	//	Find(&result).Error

	// return

	// 调整为clickhouse的查询
	if len(groupChatIngo) < 1 {
		return nil, nil
	}
	ckSql := "SELECT unionid FROM micro_scrm_ck.work_group_chat_member_local WHERE enterprise_id= ? and type= ? and status = ? and created_at <= ? and group_chat_id in (?) and unionid != ''"
	err = db.ClickHouse().Get(ctx, &result, ckSql, enterpriseID, micro_scrm.WGCML_MEMBER_TYPE, micro_scrm.WGCML_OPERATION_STATUS, endTime, groupChatIngo)
	return result, err

	// var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	// err = b.Model(&groupChatMemberLog).
	//	Select("union_id").
	//	//Select("union_id,operation_status,Max(created_at)").
	//	Where("enterprise_id = ? AND deleted_at = 0", enterpriseID).
	//	Where("member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
	//	Where("operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
	//	Where("created_at <= ?", endTime).
	//	//Where("external_userid LIKE ?", micro_scrm.EXTERNAL_CONTACTS_STRING).
	//	Where("group_chat_id IN (?)", groupChatIngo).
	//	Where("union_id != '' ").
	//	Group("union_id").
	//	//Having("operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
	//	Find(&result).Error
}

// GetContactPersonCount  获取群成员数量
func (b *DataCompassService) GetContactPersonCount(enterpriseID, endTime int64, groupChatIngo []int64) (result []model.UnionInfo, total int64, err error) {
	var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	model1 := b.Model(&groupChatMemberLog)
	model1.
		Where("enterprise_id = ? AND deleted_at = 0", enterpriseID).
		Where("member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
		Where("created_at <= ?", endTime).
		Where("group_chat_id IN (?)", groupChatIngo).
		Where("union_id != '' ").
		Group("union_id")
	model1.Count(&total)
	err = model1.Find(&result).Error

	return
}

func (b *DataCompassService) GetBAGroupInfo2(ctx ctx.Context, enterpriseID, endTime int64, baInfo []string, searchAll int64) (result []model.GroupUserInfo, err error) {
	var groupChat group_chat.WorkGroupChatModel

	// 列表查询ck和mysql的差距不大
	// var groupIds []int64
	// err = b.Model(&groupChat).
	//	Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
	//	Where("type = ?", 1).
	//	Where("status = ?", 1).
	//	Where("create_time <= ?", endTime).
	//	Pluck("id", &groupIds).Error
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	//	return nil, err
	// }
	// ckSql := "SELECT group_chat_id,userid as user_id FROM micro_scrm.work_group_chat_member_local WHERE status = 1 and group_chat_id in (?) "
	// if searchAll != 1 {
	//	ckSql += "and userid in (?)"
	//	err = db.ClickHouse().Get(ctx, &result, ckSql, groupIds, baInfo)
	// } else {
	//	err = db.ClickHouse().Get(ctx, &result, ckSql, groupIds)
	// }
	//
	// if err != nil {
	//	return nil, err
	// }
	// return

	var groupChatMember group_chat.WorkGroupChatMemberModel

	builder := b.Model(&groupChatMember).Debug().
		Select("work_group_chat_member.userid AS user_id,work_group_chat_member.group_chat_id").
		Where("work_group_chat_member.deleted_at = 0").
		Where("work_group_chat_member.status = ?", 1)
	if searchAll != 1 {
		builder.Where("work_group_chat_member.userid IN (?)", baInfo)
	}
	err = builder.Where("work_group_chat_member.group_chat_id IN (?)", b.Model(&groupChat).
		Select("id").
		Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
		Where("type = ?", 1).
		Where("status = ?", 1).
		Where("create_time <= ?", endTime)).
		Group("group_chat_id").
		Find(&result).Error

	return
}

// GetContactEnterPerson 获取进群成员信息
func (b *DataCompassService) GetContactEnterPerson(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatIngo []int64) (count int64, err error) {
	// 调整为clickhouse的查询
	// var result []struct {
	//	Count int64 `json:"count" db:"count"`
	// }
	// ckSql := "SELECT count(union_id) as count FROM micro_scrm.group_chat_member_log_local WHERE enterprise_id= ? and type= ? and status = ? and opera_time BETWEEN ? AND ? and group_chat_id in (?) and union_id != ''"
	// err = db.ClickHouse().Get(ctx, &result, ckSql, enterpriseID, micro_scrm.WGCML_MEMBER_TYPE, micro_scrm.WGCML_OPERATION_STATUS, startTime, endTime, groupChatIngo)
	// if err != nil {
	//	return 0, err
	// }
	// return result[0].Count, nil

	var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	err = b.Model(&groupChatMemberLog).
		// Joins("left join members on members.userid = work_group_chat_member_log.invitor_userid").
		// Select("work_group_chat_member_log.union_id").
		// Where("members.status", membersModel.MembersStatusEnable).
		Where("work_group_chat_member_log.enterprise_id = ? AND work_group_chat_member_log.deleted_at = 0 ", enterpriseID).
		Where("work_group_chat_member_log.member_type = ?", micro_scrm.WGCML_MEMBER_TYPE).
		Where("work_group_chat_member_log.operation_status = ? ", micro_scrm.WGCML_OPERATION_STATUS).
		// Where("external_userid LIKE ? ", micro_scrm.EXTERNAL_CONTACTS_STRING).
		Where("work_group_chat_member_log.created_at BETWEEN ? AND ?", startTime, endTime).
		Where("work_group_chat_member_log.group_chat_id IN (?)", groupChatIngo).
		Where("work_group_chat_member_log.union_id != '' ").
		Group("union_id").
		Count(&count).Error
	return
}

// GetContactLeavePerson 获取退群群成员信息
func (b *DataCompassService) GetContactLeavePerson(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatIngo []int64) (count int64, err error) {
	// 调整为clickhouse的查询
	// var result []struct {
	//	Count int64 `json:"count" db:"count"`
	// }
	// ckSql := "SELECT count(unionid) as count FROM micro_scrm.work_group_chat_member_local WHERE enterprise_id= ? and type= ? and status = ? and created_at BETWEEN ? AND ? and group_chat_id in (?) and unionid != ''"
	// err = db.ClickHouse().Get(ctx, &result, ckSql, enterpriseID, micro_scrm.WGCML_MEMBER_TYPE, 2, startTime, endTime, groupChatIngo)
	// if err != nil {
	//	return 0, err
	// }
	// return result[0].Count, nil

	var groupChatMemberLog micro_scrm.WorkGroupChatMemberLog
	err = b.Model(&groupChatMemberLog).
		Select("union_id").
		Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
		Where("member_type = ? ", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status = ? ", 2).
		Where("created_at BETWEEN ? AND ? ", startTime, endTime).
		Where("group_chat_id IN (?) ", groupChatIngo).
		Where("union_id != '' ").
		Group("union_id").
		Count(&count).Error
	return
}

// GetMemberIDList 获取ba的memberID
func (b *DataCompassService) GetMemberIDList(enterpriseID int64, baInfo []string) (result []int64, err error) {
	var member micro_scrm.MemberModel
	b.Model(&member).
		Select("id").
		Where("enterprise_id = ? AND deleted_at = 0 ").
		Where("userid IN (?)", baInfo).
		Find(&result)
	return
}

// GetRelationInfo 获取绑定关系信息
func (b *DataCompassService) GetRelationInfo(enterpriseID int64, unionList []string, baInfo []string) (result []string, err error) {
	var workExternalUserRelation work_external_user.WorkExternalUserRelationModel
	var wxUserCorpExternal micro_scrm_enterprise.WxUserCorpExternalModel
	var member micro_scrm.MemberModel
	// err = b.Scopes(model.BranchTable(&wxUserCorpExternal, enterpriseID)).
	//	Select("unionid").
	//	Where("deleted_at = 0").
	//	Where("unionid IN (?)", unionList).
	//	Where("EXISTS (?)", b.Model(&workExternalUserRelation).
	//		Select("work_external_user_id").
	//		Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
	//		Where("work_external_user_id = wx_user_corp_external.id").
	//		Where(" relation_status = 0").
	//		Where("EXISTS (?)", b.Model(&member).
	//			Select("id").
	//			Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
	//			Where("id = work_external_user_relation.member_id").
	//			Where("userid IN (?) ", baInfo))).
	//	Find(&result).Error

	err = b.Scopes(model.BranchTable(&wxUserCorpExternal, enterpriseID)).
		Select("DISTINCT wx_user_corp_external.unionid").
		Joins("LEFT JOIN "+workExternalUserRelation.TableName()+"  work_external_user_relation ON wx_user_corp_external.id = work_external_user_relation.work_external_user_id").
		Joins("LEFT JOIN "+member.TableName()+"  members ON work_external_user_relation.member_id=members.id").
		Where("wx_user_corp_external.deleted_at = 0").
		Where("wx_user_corp_external.unionid IN (?)", unionList).
		Where("work_external_user_relation.enterprise_id = ? AND work_external_user_relation.deleted_at = 0 ", enterpriseID).
		Where("work_external_user_relation.relation_status = 0 ").
		Where("members.enterprise_id = ? AND members.deleted_at = 0 ", enterpriseID).
		Where("members.userid IN (?) ", baInfo).
		Find(&result).Error

	return
}

// GetUnionCodeInfo 获取是会员的人员id
func (b *DataCompassService) GetUnionCodeInfo(enterpriseID int64, unionList []string) (result []string, err error) {
	var wxUserCorpExternal micro_scrm_enterprise.WxUserCorpExternalModel
	err = b.Scopes(model.BranchTable(&wxUserCorpExternal, enterpriseID)).
		Select("DISTINCT unionid").
		Where("deleted_at = 0").
		Where("union_code != '' AND unionid IN (?)", unionList).
		Find(&result).Error
	return
}

// GetContactActivityPerson 获取活动群成员信息
func (b *DataCompassService) GetContactActivityPerson(enterpriseID, startTime, endTime int64, activityID []int64) (result []string, err error) {
	var fissionActivity fission.FissionActivityModel
	var fissionAssRecords fission.FissionAssistantRecordsModel
	var wxUserCorpExternal micro_scrm_enterprise.WxUserCorpExternalModel
	err = b.Model(&fissionAssRecords).
		Select("UC.unionid").
		Joins("INNER JOIN "+fissionActivity.TableName()+" activity ON activity.id = fission_assistant_records.activity_id AND activity.deleted_at = 0 AND activity.enterprise_id = ? ", enterpriseID).
		Joins("INNER JOIN "+wxUserCorpExternal.AllTableName(enterpriseID)+" UC ON UC.id = fission_assistant_records.assistant_user_id AND UC.deleted_at = 0 AND ISNULL(UC.unionid) != 1").
		Where("fission_assistant_records.enterprise_id = ? AND fission_assistant_records.deleted_at = 0", enterpriseID).
		Where("activity.fission_type = ? ", 2).
		Where("fission_assistant_records.created_at BETWEEN ? AND ? ", startTime, endTime).
		Where("fission_assistant_records.activity_id IN (?)", activityID).
		Group("UC.unionid").
		Find(&result).Error
	return
}

func (b *DataCompassService) GetContactActivityPerson2(enterpriseID, startTime, endTime int64, activityID []int64) (result []string, err error) {
	var fissionUser fission.FissionUserModel
	var fissionMp fission.FissionMpUserModel
	err = b.Model(&fissionUser).
		Select("mp.union_id").
		Joins("LEFT JOIN "+fissionMp.TableName()+" mp ON mp.id = fission_user.user_id AND "+
			"mp.deleted_at = 0 AND mp.enterprise_id = ?", enterpriseID).
		Where("fission_user.enterprise_id = ? AND fission_user.deleted_at = 0 ", enterpriseID).
		Where("fission_user.activity_id IN (?)", activityID).
		Where("fission_user.created_at BETWEEN ? AND ? ", startTime, endTime).
		Group("mp.union_id").
		Find(&result).Error
	return
}

// GetFissionActivityInfo 获取裂变活动的群
func (b *DataCompassService) GetFissionActivityInfo(enterpriseID, endTime int64, baInfo []string) (result []model.FissionChatGroup, err error) {
	var fissionChatGroup fission.FissionChatGroupModel
	var groupChatMember group_chat.WorkGroupChatMemberModel
	var member micro_scrm.MemberModel
	err = b.Model(&fissionChatGroup).
		Select("chat_id,activity_id").
		Joins("LEFT JOIN "+groupChatMember.TableName()+"  member ON member.group_chat_id = fission_chat_group.chat_id "+
			"AND member.deleted_at = 0 AND member.enterprise_id = ?", enterpriseID).
		Where("fission_chat_group.enterprise_id = ? AND fission_chat_group.deleted_at = 0 ", enterpriseID).
		Where("fission_chat_group.created_at <= ? ", endTime).
		Where("member.member_id IN  (?) ", b.Model(&member).
			Select("id").
			Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
			Where("userid IN (?) ", baInfo)).
		Find(&result).Error
	return
}
func (b *DataCompassService) GetIntoActivityGroupInfo(enterpriseID, startTime, endTime int64, activityID []int64) (result []int64, err error) {
	var fissionAssistabtRecords fission.FissionAssistantRecordsModel
	err = b.Model(&fissionAssistabtRecords).
		Select("DISTINCT assistant_user_id").
		Where("enterprise_id = ? AND deleted_at = 0 ", enterpriseID).
		Where("join_chat_group_status = ? ", 1).
		Where("created_at BETWEEN ? AND ? ", startTime, endTime).
		Where("activity_id IN (?)", activityID).
		Pluck("assistant_user_id", &result).Error
	return
}

func (b *DataCompassService) GetMemberInfo(enterpriseID int64, helpUnion []string) (result int64) {
	var userExternal micro_scrm_enterprise.WxUserCorpExternalModel
	b.Scopes(model.BranchTable(&userExternal, enterpriseID)).
		Select("count(DISTINCT id)").
		Where("union_code != '' AND deleted_at = 0 ").
		Where("unionid IN (?)", helpUnion).
		Find(&result)
	return
}

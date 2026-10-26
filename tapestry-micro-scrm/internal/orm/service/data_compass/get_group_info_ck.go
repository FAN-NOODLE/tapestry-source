package data_compass

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

// GetBAGroupInfoCk 获取外部群信息
func (b *DataCompassService) GetBAGroupInfoCk(ctx ctx.Context, enterpriseID, endTime int64, baInfo []string, searchAll int64) (result []model.GroupUserInfo, err error) {

	//原生sql

	//SELECT
	//work_group_chat_member.userid AS user_id,
	//	work_group_chat_member.group_chat_id
	//FROM
	//`work_group_chat_member`
	//WHERE
	//( work_group_chat_member.deleted_at = 0 )
	//AND ( work_group_chat_member.STATUS = 1 )
	//AND (
	//	work_group_chat_member.group_chat_id IN ( SELECT `id` FROM `work_group_chat` WHERE ( enterprise_id = 3363603667487424 AND deleted_at = 0 ) AND type = 1 AND STATUS = 1 AND create_time <= 1677167999 ))
	//GROUP BY
	//`group_chat_id`

	// ck查询
	groupChatIdSql := gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat").
		Select("id").
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("type", 1).
		Where("status", 1).
		Where("create_time <= ?", endTime).
		Group("id")

	chatMemberSql := gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member").
		Select("any(work_group_chat_member.userid) AS user_id,work_group_chat_member.group_chat_id").
		Where("group_chat_id in(?)", groupChatIdSql).
		Where("status", 1).
		Where("deleted_at", 0).
		Group("group_chat_id")

	if searchAll != 1 {
		chatMemberSql.Where("userid IN (?)", baInfo)
	}

	chatMemberSql.Find(&result)

	return
}

// GetGroupCustomCountCk 获取群成员统计
func (b *DataCompassService) GetGroupCustomCountCk(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64, baInfo []string, searchAll int64) (count int64, err error) {

	//原始sql 需要把left join改成 in的方式
	//SELECT
	//count( DISTINCT work_group_chat_member_log.union_id )
	//FROM
	//`work_group_chat_member_log`
	//LEFT JOIN work_external_user_relation ON work_external_user_relation.unionid = work_group_chat_member_log.union_id
	//WHERE
	//( work_group_chat_member_log.enterprise_id = 3363603667487424 AND work_group_chat_member_log.deleted_at = 0 )
	//AND work_group_chat_member_log.member_type = 2
	//AND work_group_chat_member_log.operation_status = 1
	//AND work_group_chat_member_log.created_at <= 1677167999
	//AND work_group_chat_member_log.group_chat_id IN ( 3416916650425024 )
	//AND work_group_chat_member_log.union_id != ''
	//AND ( work_external_user_relation.enterprise_id = 3363603667487424 AND work_external_user_relation.deleted_at = 0 )
	//AND work_external_user_relation.relation_status = 0

	unionIdSql := gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_external_user_relation").
		Select("unionid").
		Where("enterprise_id", enterpriseID).
		Where("unionid != ?", "").
		Where("deleted_at", 0).
		Where("relation_status", 0).
		Group("unionid")
	if searchAll != 1 {
		unionIdSql.Where("user_id IN(?)", baInfo)
	}

	err = gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member_log").
		Select("count( DISTINCT union_id)").
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("member_type", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status", micro_scrm.WGCML_OPERATION_STATUS).
		Where("created_at <= ?", endTime).
		Where("group_chat_id in ?", groupChatIngo).
		Where("union_id in(?)", unionIdSql).
		Find(&count).Error

	return count, err

}

//	GetGroupCustomerNumberCountCk 获取累计群成员信息
func (b *DataCompassService) GetGroupCustomerNumberCountCk(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64) (count int64, err error) {

	//原生sql
	//SELECT
	//COUNT(
	//	DISTINCT ( `work_group_chat_member_log`.`union_id` ))
	//FROM
	//`work_group_chat_member_log`
	//WHERE
	//( work_group_chat_member_log.enterprise_id = 3363603667487424 AND work_group_chat_member_log.deleted_at = 0 )
	//AND work_group_chat_member_log.member_type = 2
	//AND work_group_chat_member_log.operation_status = 1
	//AND work_group_chat_member_log.created_at <= 1677513599
	//AND work_group_chat_member_log.group_chat_id IN ( 3418855331173165, 3418855331173176, 3397458890123968 )
	//AND work_group_chat_member_log.union_id != ''

	err = gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member_log").
		Select("count( DISTINCT union_id)").
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("member_type", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status", micro_scrm.WGCML_OPERATION_STATUS).
		Where("created_at <= ?", endTime).
		Where("group_chat_id in ?", groupChatIngo).
		Where("union_id != ?", "").
		Find(&count).Error

	return count, err
}

// GetContactEnterPersonCk 获取进群成员信息
func (d *DataCompassService) GetContactEnterPersonCk(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatIngo []int64) (count int64, err error) {

	//原生sql
	//SELECT
	//count(*)
	//FROM
	//`work_group_chat_member_log`
	//WHERE
	//( work_group_chat_member_log.enterprise_id = 3363603667487424 AND work_group_chat_member_log.deleted_at = 0 )
	//AND work_group_chat_member_log.member_type = 2
	//AND work_group_chat_member_log.operation_status = 1
	//AND ( work_group_chat_member_log.created_at BETWEEN 1677427200 AND 1677513599 )
	//AND work_group_chat_member_log.group_chat_id IN (3432300738318272)
	//AND work_group_chat_member_log.union_id != ''
	//GROUP BY
	//`union_id`

	err = gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member_log").
		Select("count(DISTINCT union_id)").
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("member_type", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status", micro_scrm.WGCML_OPERATION_STATUS).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Where("group_chat_id in ?", groupChatIngo).
		Where("union_id != ?", "").
		Find(&count).Error

	return count, err
}

// GetContactLeavePersonCk 获取退群群成员信息
func (b *DataCompassService) GetContactLeavePersonCk(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatIngo []int64) (count int64, err error) {

	//原生sql
	//SELECT
	//COUNT( `union_id` )
	//FROM
	//`work_group_chat_member_log`
	//WHERE
	//( enterprise_id = 3363603667487424 AND deleted_at = 0 )
	//AND member_type = 2
	//AND operation_status = 2
	//AND ( created_at BETWEEN 1677427200 AND 1677513599 )
	//AND group_chat_id IN (3426106005045184)
	//AND union_id != ''
	//GROUP BY
	//`union_id`

	err = gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member_log").
		Select("union_id").
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("member_type", micro_scrm.WGCML_MEMBER_TYPE).
		Where("operation_status", 2).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Where("group_chat_id in ?", groupChatIngo).
		Where("union_id != ?", "").
		Group("union_id").
		Count(&count).Error
	return count, err
}

// GetFissionActivityInfoCk 获取裂变活动的群
func (b *DataCompassService) GetFissionActivityInfoCk(ctx ctx.Context, enterpriseID, endTime int64, baInfo []string, searchAll int64) (result []model.FissionChatGroup, err error) {

	var res []model.FissionChatGroup
	return res, nil
	//原生sql
	//SELECT
	//chat_id,
	//	activity_id
	//FROM
	//`fission_chat_group`
	//LEFT JOIN work_group_chat_member member ON member.group_chat_id = fission_chat_group.chat_id
	//AND member.deleted_at = 0
	//AND member.enterprise_id = 3363603667487424
	//WHERE
	//( fission_chat_group.enterprise_id = 3363603667487424 AND fission_chat_group.deleted_at = 0 )
	//AND fission_chat_group.created_at <= 1677513599
	//AND member.member_id IN (
	//	SELECT
	//`id`
	//FROM
	//`members`
	//WHERE
	//( enterprise_id = 3363603667487424 AND deleted_at = 0 )
	//AND userid IN ( NULL )
	//)

	//1.查members member_id 2.查work_group_chat_member group_chat_id 3.fission_chat_group chat_id和activity_id
	//memberIdSql := gormx.GetCkInstance().Table("micro_scrm_ck.members").
	//	Select("id").
	//	Where("userid != ?", "").
	//	Where("enterprise_id", enterpriseID).
	//	Where("deleted_at", 0)
	//if searchAll != 1 {
	//	memberIdSql.Where("userid IN(?)", baInfo)
	//}
	//memberIdSql.Group("id")
	//memberIdSql.Group("userid")
	//
	//groupChatMemberSql := gormx.GetCkInstance().Table("micro_scrm_ck.work_group_chat_member").
	//	Select("group_chat_id").
	//	Where("enterprise_id", enterpriseID).
	//	Where("status > ?", 0).
	//	Where("deleted_at", 0).
	//	Where("member_id in(?)", memberIdSql).
	//	Group("group_chat_id").
	//	Group("enterprise_id")
	//
	//err = gormx.GetCkInstance().Table("micro_scrm_ck.fission_chat_group").
	//	Select("chat_id,activity_id").
	//	Where("enterprise_id", enterpriseID).
	//	Where("deleted_at", 0).
	//	Where("created_at <= ?", endTime).
	//	Where("chat_id in(?)", groupChatMemberSql).
	//	Find(&result).Error
	//return
}

// GetContactPersonCk 获取群成员信息
func (b *DataCompassService) GetContactPersonCk(ctx ctx.Context, enterpriseID, endTime int64, groupChatIngo []int64) (result []model.UnionInfo, err error) {
	if len(groupChatIngo) < 1 {
		return nil, nil
	}
	err = gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_group_chat_member").
		Select("unionid").
		Where("enterprise_id", enterpriseID).
		Where("status", micro_scrm.WGCML_OPERATION_STATUS).
		Where("deleted_at", 0).
		Where("type", micro_scrm.WGCML_MEMBER_TYPE).
		Where("created_at <= ?", endTime).
		Where("group_chat_id in(?)", groupChatIngo).
		Where("unionid != ?", "").
		Group("unionid").
		Find(&result).Error

	return result, err
}

// GetContactActivityPersonCk 获取活动群成员信息 社区数据不查询
func (b *DataCompassService) GetContactActivityPersonCk(ctx ctx.Context, enterpriseID, startTime, endTime int64, activityID []int64) (result []string, err error) {

	var res []string
	return res, nil
}

// GetContactActivityPerson2Ck 社区数据不查询
func (d *DataCompassService) GetContactActivityPerson2Ck(ctx ctx.Context, enterpriseID, startTime, endTime int64, activityID []int64) (result []string, err error) {
	var res []string
	return res, nil
}

// GetIntoActivityGroupInfoCk 社区数据不查询
func (d *DataCompassService) GetIntoActivityGroupInfoCk(ctx ctx.Context, enterpriseID, startTime, endTime int64, activityID []int64) (result []int64, err error) {
	var res []int64
	return res, nil
}

// GetMemberInfoCk 社区数据不查询
func (d *DataCompassService) GetMemberInfoCk(ctx ctx.Context, enterpriseID int64, helpUnion []string) (result int64) {
	var res int64
	return res
}

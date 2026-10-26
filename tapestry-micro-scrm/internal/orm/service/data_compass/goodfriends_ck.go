package data_compass

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db/gormx"
)

var scrmDatabases = "micro_scrm_ck"

func (d DataCompassService) GetScrmTableCk(table string) string {
	return fmt.Sprintf("%s.%s", scrmDatabases, table)
}

//GetMemberIdCk 获取指定id的unionIds
func (d DataCompassService) GetMemberIdCk(c ctx.Context, enterpriseId int64, memberUserId []string) (data []micro_scrm_enterprise.MembersModel, err error) {

	fmt.Println(d.GetScrmTableCk("members"))
	err = gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("members")).
		Select("id").
		Where("enterprise_id = ?", enterpriseId).
		Where("deleted_at = ?", 0).
		Where("userid in (?)", memberUserId).
		Find(&data).Error
	return
}

//GetNewGoodFriendsCountCk 获取好友新增数
func (d DataCompassService) GetNewGoodFriendsCountCk(c ctx.Context, enterpriseId int64, memberId []string, startTime int64, endTime int64, searchAll int64) (count int64, err error) {

	scope := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_external_user_relation")).
		Select("COUNT(DISTINCT(`work_external_user_id`))").
		Where("enterprise_id", enterpriseId).
		Where("unionid != ?", "").
		Where("deleted_at", 0)

	if startTime > 0 && endTime > 0 && endTime > startTime {
		scope.Where("createtime > ?", startTime).
			Where("createtime < ?", endTime)
	}

	if searchAll != 1 {
		scope.Where("member_id in (?)", memberId)
	}
	err = scope.Count(&count).Error
	return
}

//GetJoinCroupCountCk 是否入群(获取新增好友入群人数)
func (d DataCompassService) GetJoinCroupCountCk(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberIds []string, startTime int64, endTime, searchAll int64) (data int64, err error) {

	groupChatMemberSql := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_group_chat_member")).
		Select("DISTINCT unionid").
		Where("enterprise_id", enterpriseId).
		Where("type", 2).
		Where("deleted_at", 0).
		Where("group_chat_id in(?)", baGroupIds).
		Where("unionid != ?", "").
		Where("status", 1)

	relationSql := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_external_user_relation")).
		Select("COUNT(DISTINCT work_external_user_id)").
		Where("enterprise_id", enterpriseId).
		Where("deleted_at", 0)

	if searchAll != 1 {
		relationSql.Where("member_id in(?)", memberIds)
	}

	err = relationSql.Where("relation_status", 0).
		Where("createtime >= ?", startTime).
		Where("createtime <= ?", endTime).
		Where("unionid != ?", "").
		Where("unionid in(?)", groupChatMemberSql).
		Find(&data).Error
	return
}

//GetGoodFriendsAllCk 获取累计好友数据
func (d DataCompassService) GetGoodFriendsAllCk(c ctx.Context, enterpriseId int64, memberId []string, endTime, searchAll int64) (total int64, err error) {

	scope := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_external_user_relation")).Select("count(*)")
	if endTime > 0 {
		scope.Where("createtime < ?", endTime)
	}

	builder := scope.Where("deleted_at", 0).
		Where("enterprise_id", enterpriseId).
		Where("relation_status", 0).
		Where("unionid != ?", "")
	if searchAll != 1 {
		builder.Where("member_id in (?)", memberId)
	}
	err = builder.Find(&total).Error
	return
}

//GetDrainFriendsCk 获取流失好友数据
func (d DataCompassService) GetDrainFriendsCk(c ctx.Context, enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (data int64, err error) {

	scope := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_external_user_relation")).Select("count(DISTINCT unionid)")
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
		Find(&data).Error
	return
}

//GetJoinCroupCountAllCk 通过unionid获取入群人数
func (d DataCompassService) GetJoinCroupCountAllCk(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberIds []string, endTime, searchAll int64) (data int64, err error) {

	groupChatMemberSql := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_group_chat_member")).
		Select("DISTINCT unionid").
		Where("enterprise_id", enterpriseId).
		Where("type", 2).
		Where("deleted_at", 0).
		Where("group_chat_id in(?)", baGroupIds).
		Where("unionid != ?", "").
		Where("status", 1)

	relationSql := gormx.GetCkInstance().WithContext(c).Table(d.GetScrmTableCk("work_external_user_relation")).
		Select("COUNT(DISTINCT work_external_user_id)").
		Where("enterprise_id", enterpriseId).
		Where("deleted_at", 0)

	if searchAll != 1 {
		relationSql.Where("member_id in(?)", memberIds)
	}

	err = relationSql.Where("relation_status", 0).
		Where("createtime <= ?", endTime).
		Where("unionid != ?", "").
		Where("unionid in(?)", groupChatMemberSql).
		Find(&data).Error
	return
}

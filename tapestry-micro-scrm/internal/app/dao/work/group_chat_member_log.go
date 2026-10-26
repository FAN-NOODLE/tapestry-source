package workDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	workGroupModel "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkGroupChatMemberLog struct {
	dao       dao.Dao
	tableName string
}

func NewWorkGroupChatMemberLog() *WorkGroupChatMemberLog {
	return &WorkGroupChatMemberLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workGroupModel.WorkGroupChatMemberLogModel{}.TableName(),
	}
}

// GetJoinGroupUserCounts 获取加群人数
func (d *WorkGroupChatMemberLog) GetJoinGroupUserCounts(c ctx.Context, enterpriseId, startTime, entTime int64, groupChatIds []int64) (data int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"group_chat_id":    groupChatIds,
		"operation_status": workGroupModel.JOINGROUPSTATUS,
		"member_type":      workGroupModel.MEMBERTYPEWECHAT,
	}).Where("created_at < ? AND union_id!='' ", entTime).Distinct("union_id").Count(&data).Error
	return
}

// GetChatFriendsRelationCount 获取群成员是好友数
func (d *WorkGroupChatMemberLog) GetChatFriendsRelationCount(c ctx.Context, enterpriseId, startTime, entTime int64, groupChatIds []int64, memberIds []int64) (data int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Joins("LEFT JOIN (?) as b on b.unionid=work_group_chat_member_log.union_id", d.dao.WithContext(c).Table(microScrmModel.WorkExternalUserRelationModel{}.TableName()).Where(map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"member_id":       memberIds,
		"relation_status": 0,
	})).Where("work_group_chat_member_log.created_at < ? AND work_group_chat_member_log.union_id!='' ", entTime).Where(map[string]interface{}{
		"work_group_chat_member_log.enterprise_id":    enterpriseId,
		"work_group_chat_member_log.group_chat_id":    groupChatIds,
		"work_group_chat_member_log.operation_status": workGroupModel.JOINGROUPSTATUS,
		"work_group_chat_member_log.member_type":      workGroupModel.MEMBERTYPEWECHAT,
	}).Distinct("work_group_chat_member_log.union_id").Count(&data).Error
	return
}

// GetLeaveGroupCounts 获取离群的人数
func (d *WorkGroupChatMemberLog) GetLeaveGroupCounts(c ctx.Context, enterpriseId, startTime, entTime int64, groupChatIds []int64) (data int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"group_chat_id":    groupChatIds,
		"operation_status": workGroupModel.LEAVEGROUPSTATUS,
		"member_type":      workGroupModel.MEMBERTYPEWECHAT,
	}).Where("created_at BETWEEN ? AND ? AND union_id!='' ", startTime, entTime).Distinct("union_id").Count(&data).Error
	return
}

// GetJoinGroupCounts 获取入群的人数
func (d *WorkGroupChatMemberLog) GetJoinGroupCounts(c ctx.Context, enterpriseId, startTime, entTime int64, groupChatIds []int64) (data int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"group_chat_id":    groupChatIds,
		"operation_status": workGroupModel.JOINGROUPSTATUS,
		"member_type":      workGroupModel.MEMBERTYPEWECHAT,
	}).Where("created_at BETWEEN ? AND ? AND union_id!='' ", startTime, entTime).Distinct("union_id").Count(&data).Error
	return
}

// GetChatNewFriendsRelationCount 获得新增好友的入群统计
func (d *WorkGroupChatMemberLog) GetChatNewFriendsRelationCount(c ctx.Context, enterpriseId, startTime, entTime int64, groupChatIds []int64, memberIds []int64) (total int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Joins("LEFT JOIN (?) as b on b.unionid=work_group_chat_member_log.union_id", d.dao.WithContext(c).Table(microScrmModel.WorkExternalUserRelationModel{}.TableName()).Where(map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"member_id":       memberIds,
		"relation_status": 0,
	}).Where("created_at between ? and ? ", startTime, entTime)).Where("work_group_chat_member_log.created_at BETWEEN ? AND ? AND work_group_chat_member_log.union_id!='' ", startTime, entTime).Where(map[string]interface{}{
		"work_group_chat_member_log.enterprise_id":    enterpriseId,
		"work_group_chat_member_log.group_chat_id":    groupChatIds,
		"work_group_chat_member_log.operation_status": workGroupModel.JOINGROUPSTATUS,
		"work_group_chat_member_log.member_type":      workGroupModel.MEMBERTYPEWECHAT,
	}).Distinct("work_group_chat_member_log.union_id").Count(&total).Error
	return
}

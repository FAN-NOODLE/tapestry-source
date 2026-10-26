package workDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app/dao"
	dataCompassReq "micro-scrm/internal/app/http/request/datacompass"
	workGroupModel "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkGroupChatMemberCk struct {
	dao       dao.Dao
	database  string
	tableName string
}

type groupChatMemberCount struct {
	GroupChatId int64 `json:"group_chat_id"`
	Count       int64 `json:"count"`
}

func NewWorkGroupChatMemberCk() *WorkGroupChatMemberCk {
	return &WorkGroupChatMemberCk{
		dao:       dao.Dao{DB: gormx.GetCkInstance()},
		database:  "micro_scrm_ck",
		tableName: workGroupModel.WorkGroupChatMemberModel{}.TableName(),
	}
}

// SelectTotalCount 查询总数
func (d *WorkGroupChatMemberCk) SelectTotalCount(c ctx.Context, chatGroupIds []int64, req dataCompassReq.ChatGroupReq) (map[int64]int64, error) {
	var groupChatMember []groupChatMemberCount

	result := make(map[int64]int64)

	err := d.dao.WithContext(c).
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("%s.%s", d.database, d.tableName))
		}).
		Select("group_chat_id,count(DISTINCT userid) as count").
		Where("group_chat_id", chatGroupIds).
		Where("join_time <= ?", req.EndTime).
		Where("enterprise_id = ?", req.EnterpriseID).
		Where("status = 1").
		Where("deleted_at = 0").
		Group("group_chat_id").
		Scan(&groupChatMember).Error

	if err != nil {
		return nil, err
	}

	for _, group := range groupChatMember {
		result[group.GroupChatId] = group.Count
	}

	return result, nil
}

// SelectAddCount 查询新增群成员数
func (d *WorkGroupChatMemberCk) SelectAddCount(c ctx.Context, chatGroupIds []int64, req dataCompassReq.ChatGroupReq) (map[int64]int64, error) {
	var groupChatMember []groupChatMemberCount

	result := make(map[int64]int64)

	err := d.dao.WithContext(c).
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("%s.%s", d.database, d.tableName))
		}).
		Select("group_chat_id,count(DISTINCT userid) as count").
		Where("group_chat_id", chatGroupIds).
		Where("join_time >= ? and join_time <= ?", req.StartTime, req.EndTime).
		Where("enterprise_id = ?", req.EnterpriseID).
		Where("status = 1").
		Where("deleted_at = 0").
		Group("group_chat_id").
		Scan(&groupChatMember).Error

	if err != nil {
		return nil, err
	}

	for _, group := range groupChatMember {
		result[group.GroupChatId] = group.Count
	}

	return result, nil
}

// SelectLossCount 查询流失群成员数
func (d *WorkGroupChatMemberCk) SelectLossCount(c ctx.Context, chatGroupIds []int64, req dataCompassReq.ChatGroupReq) (map[int64]int64, error) {
	var groupChatMember []groupChatMemberCount

	result := make(map[int64]int64)

	err := d.dao.WithContext(c).
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("%s.%s", d.database, d.tableName))
		}).
		Select("group_chat_id,count(DISTINCT userid) as count").
		Where("group_chat_id", chatGroupIds).
		Where("leave_time >= ? and leave_time <= ?", req.StartTime, req.EndTime).
		Where("enterprise_id = ?", req.EnterpriseID).
		Where("status = 2").
		Where("deleted_at = 0").
		Group("group_chat_id").
		Scan(&groupChatMember).Error

	if err != nil {
		return nil, err
	}

	for _, group := range groupChatMember {
		result[group.GroupChatId] = group.Count
	}

	return result, nil
}

// SelectMemberCount 获取群中导购数量
func (d *WorkGroupChatMemberCk) SelectMemberCount(c ctx.Context, chatGroupIds []int64, req dataCompassReq.ChatGroupReq) (map[int64]int64, error) {
	var groupChatMember []groupChatMemberCount

	result := make(map[int64]int64)

	err := d.dao.WithContext(c).
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("%s.%s", d.database, d.tableName))
		}).
		Select("group_chat_id,count(DISTINCT member_id) as count").
		Where("group_chat_id", chatGroupIds).
		Where("join_time >= ? and join_time <= ?", req.StartTime, req.EndTime).
		Where("enterprise_id = ?", req.EnterpriseID).
		Where("status = 1").
		Where("member_id > 0").
		Where("deleted_at = 0").
		Group("group_chat_id").
		Scan(&groupChatMember).Error

	if err != nil {
		return nil, err
	}

	for _, group := range groupChatMember {
		result[group.GroupChatId] = group.Count
	}

	return result, nil
}

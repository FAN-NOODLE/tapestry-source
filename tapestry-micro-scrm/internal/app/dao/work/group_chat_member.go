package workDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao"
	workGroupModel "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkGroupChatMember struct {
	dao       dao.Dao
	tableName string
}

func NewWorkGroupChatMember() *WorkGroupChatMember {
	return &WorkGroupChatMember{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workGroupModel.WorkGroupChatMemberModel{}.TableName(),
	}
}

func (d *WorkGroupChatMember) GetChatUnionIds(c ctx.Context, enterpriseId, groupChatId, startTime, endTime int64, status int) ([]string, error) {
	var data []string
	err := d.dao.WithContext(c).Table(d.tableName).Where("leave_time between ? and ?", startTime, endTime).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"group_chat_id": groupChatId,
		"status":        status,
	}).Group("unionid").Pluck("unionid", &data).Error
	return data, err
}

// GetChatIdsByMemberIds 根据memberIds 查询时间内的群id
func (d *WorkGroupChatMember) GetChatIdsByMemberIds(c ctx.Context, enterpriseId, startTime, entTime int64, memberIds []int64) (data []int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberIds,
		"type":          1,
		"leave_time":    0,
	}).Where("created_at < ? ", entTime).Group("group_chat_id").Pluck("group_chat_id", &data).Error
	return
}

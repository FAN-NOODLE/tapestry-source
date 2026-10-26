package workGroupDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	workGroupModel "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkGroupChat struct {
	dao       dao.Dao
	tableName string
}

func NewWorkGroupChat() *WorkGroupChat {
	return &WorkGroupChat{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workGroupModel.GroupChatModel{}.TableName(),
	}
}

func (d *WorkGroupChat) GetUserGroupChatNum(c app.Ctx, enterpriseId int64, memberIds []int64) ([]workGroupModel.WorkGroupChatNumList, error) {
	var result []workGroupModel.WorkGroupChatNumList
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("owner_member_id in ?", memberIds).
		Where("enterprise_id = ?", enterpriseId).
		Select("count(*) as num,owner_member_id").Find(&result).Error

	return result, err
}

func (d *WorkGroupChat) GetUserGroupChatNumForUserId(c app.Ctx, enterpriseId int64, userIds []string) ([]workGroupModel.WorkGroupChatNumList, error) {
	var result []workGroupModel.WorkGroupChatNumList
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("count(*) as num,work_group_chat.owner_member_id,m.userid").
		Joins("LEFT JOIN members m ON m.id = work_group_chat.owner_member_id").
		Where("m.userid in ?", userIds).
		Where("m.status = 1").
		Where("work_group_chat.status = 1").
		Where("work_group_chat.type = 1").
		Where("work_group_chat.enterprise_id = ?", enterpriseId).
		Group("work_group_chat.owner_member_id").
		Find(&result).Error

	return result, err
}

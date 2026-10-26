package group_chat

import (
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/service"
)

// GetCustomerGroupInfo 根据顾客的union_id，获取顾客的有效群信息
func (g *GroupChatService) GetCustomerGroupInfo(param group_chat.WorkGroupChatMemberModel) (data []group_chat.WorkGroupChatMemberModel, err error) {
	err = g.Scopes(service.ScopeCond_EQDelete).
		Where(&param).
		Where("leave_time = ?", 0).
		Preload("GroupChat").
		Find(&data).Error

	return
}

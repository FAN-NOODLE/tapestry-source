package group_chat

import (
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"time"
)

// UpdateWorkGroupChatDismiss 更新群状态
func (g *GroupChatService) UpdateWorkGroupChatDismiss(enterpriseID int64, ChatID string) error {
	return g.Model(&group_chat.WorkGroupChatModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"chat_id":       ChatID,
	}).Updates(map[string]interface{}{
		"status":        group_chat.WGC_STATUS_NO,
		"dissolve_time": time.Now().Unix(),
	}).Error
}

package microScrmModel

import "micro-scrm/internal/app/models"

type TaskAllotLogMsgModel struct {
	models.BaseModel
	TaskAllotLogId int64  `gorm:"column:task_allot_log_id;default:0;NOT NULL" json:"task_allot_log_id"`
	MsgId          string `gorm:"column:msg_id;default:'';NOT NULL" json:"msg_id"`  // 群发消息msg_id
	IsSend         int    `gorm:"column:is_send;default:0;NOT NULL" json:"is_send"` // 是否完成发送，0否，1是
}

func (TaskAllotLogMsgModel) TableName() string {
	return "task_allot_log_msg"
}

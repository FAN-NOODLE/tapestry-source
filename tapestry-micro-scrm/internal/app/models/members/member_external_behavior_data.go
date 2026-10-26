package membersModel

import "micro-scrm/internal/app/models"

type MemberExternalBehaviorDataModel struct {
	models.MicroBaseModel

	MemberID            string `json:"member_id" gorm:"column:member_id;type:varchar(50);not null;index:member_id-index;comment:members表id"`
	Userid              uint64 `json:"userid" gorm:"column:userid;type:bigint(20) unsigned;not null;index:userid-index;default:0;comment:员工企微账号"`
	Date                string `json:"date" gorm:"column:date;type:varchar(50);not null;comment:日期"`
	StatTime            uint   `json:"stat_time" gorm:"column:stat_time;type:int(10) unsigned;not null;index:stat_time-index;default:0;comment:开始时间"`
	NewApplyCnt         uint   `json:"new_apply_cnt" gorm:"column:new_apply_cnt;type:int(10) unsigned;not null;default:0;comment:发起申请数"`
	NewContactCnt       uint   `json:"new_contact_cnt" gorm:"column:new_contact_cnt;type:int(10) unsigned;not null;default:0;comment:新增客户数"`
	ChatCnt             uint   `json:"chat_cnt" gorm:"column:chat_cnt;type:int(10) unsigned;not null;default:0;comment:聊天总数"`
	MessageCnt          uint   `json:"message_cnt" gorm:"column:message_cnt;type:int(10) unsigned;not null;default:0;comment:发送消息数"`
	ReplyPercentage     string `json:"reply_percentage" gorm:"column:reply_percentage;type:varchar(50);not null;comment:已回复聊天占比"`
	AvgReplyTime        uint   `json:"avg_reply_time" gorm:"column:avg_reply_time;type:int(10) unsigned;not null;default:0;comment:平均首次回复时长(分钟)"`
	NegativeFeedbackCnt uint   `json:"negative_feedback_cnt" gorm:"column:negative_feedback_cnt;type:int(10) unsigned;not null;default:0;comment:删除/拉黑成员的客户数"`
}

func (MemberExternalBehaviorDataModel) TableName() string {
	return "member_external_behavior_data"
}

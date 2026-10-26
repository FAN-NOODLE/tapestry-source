package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationAchievementModel struct {
	models.MicroBaseModel
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:enterprise_id_member_id_type_deleted_at_index;default:0;comment:导购member表id"`
	Type         int64  `json:"type" gorm:"column:type;type:tinyint(1) unsigned;not null;index:enterprise_id_member_id_type_deleted_at_index;default:0;comment:触发类型:1:邀约数(每次) 2:邀约到店率(每周) 3:邀约成功率(每月) 4:邀约成功(每次) 5:邀约占比(每周) 6:邀约到店确认(每次 )7:任务完成(每次) 8:任务完成率(每月) 9:任务触达(每人) 10: 新增好友数(每人) 11:客户流失数(每人) 12:聊天总数(每月) 13:已回复聊天占比(每月) 14:完成播种按(次) 15:积分等比兑换（每分） 16:线上销售业绩率（每周）"`
	Value        int64  `json:"value" gorm:"column:value;type:int(10) unsigned;not null;default:0;comment:任务值(百分比值存百分比整数）"`
	UniqueHash   string `json:"unique_hash" gorm:"column:unique_hash;type:varchar(32);uniqueIndex:unique_hash_index;comment:重复验证，主要用于避免定时任务重复写入"`
	SourceRemark string `json:"source_remark" gorm:"column:source_remark;type:varchar(255);comment:记录数据的特殊来源，比如培训相关"`
}

// TableName returns the table name of the ExcitationAchievement model
func (ExcitationAchievementModel) TableName() string {
	return "excitation_achievement"
}

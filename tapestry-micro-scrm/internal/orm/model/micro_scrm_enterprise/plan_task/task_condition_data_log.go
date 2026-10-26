package plan_task

import "micro-scrm/internal/orm/model"

const (
	TCDL_IsSuccess_No = 1 + iota
	TCDL_IsSuccess_Yes
)

const (
	TCDL_IsUnion_No = iota
	TCDL_IsUnion_Yes
)

const (
	TCDL_IsHasTag_No = iota
	TCDL_IsHasTag_Yes
)

type TaskConditionDataLogModel struct {
	model.BranchModel

	TagName            string `gorm:"column:tag_name;NOT NULL" json:"tag_name"`                           // CampaignTool标签名称
	WwTagId            string `gorm:"column:ww_tag_id;NOT NULL" json:"ww_tag_id"`                         // 企微标签id
	UnionId            string `gorm:"column:union_id;NOT NULL" json:"union_id"`                           // 顾客的UnionID
	WorkExternalUserId string `gorm:"column:work_external_user_id;NOT NULL" json:"work_external_user_id"` // 外部联系人user_id,可能空
	IsUnion            int    `gorm:"column:is_union;default:0;NOT NULL" json:"is_union"`                 // 关系状态；0没有好友关系、1有好友关系
	Mobile             string `gorm:"column:mobile" json:"mobile"`                                        // 手机
	MarsId             string `gorm:"column:mars_id;NOT NULL" json:"mars_id"`                             // BA的MarsCode
	UserId             string `gorm:"column:user_id;NOT NULL" json:"user_id"`                             // 员工的企微user_id
	IsSuccess          int    `gorm:"column:is_success;default:1" json:"is_success"`                      // 1 未成功，2成功
	TagNum             int    `gorm:"column:tag_num;default:0" json:"tag_num"`                            // 请求企微标签接口的次数
	MemberId           int64  `gorm:"column:member_id;default:0" json:"member_id"`                        // micro_scrm.member.id
	TagId              int64  `gorm:"column:tag_id;default:0" json:"tag_id"`                              // micro_scrm.work_external_tag.id
	RelationId         int64  `gorm:"column:relation_id;default:0" json:"relation_id"`                    // micro_scrm.work_external_user_relation.id
	IsHasTag           int    `gorm:"column:is_has_tag;default:0" json:"is_has_tag"`                      // 0 之前没有该标签，1有
}

func (TaskConditionDataLogModel) TableName() string {
	return "task_condition_data_log"
}

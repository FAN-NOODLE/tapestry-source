package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

const (
	SMCT_IsMember_NO  = 0
	SMCT_IsMember_YES = 1
)

type SyncMemberCampaignToolModel struct {
	model.BranchModel

	UnionID      string `json:"union_id" gorm:"column:union_id;type:varchar(64);not null;index:union_id_index"`
	IsMember     uint8  `json:"is_member" gorm:"column:is_member;type:tinyint(1) unsigned;not null;default:0;comment:是否是会员"` // 是否是会员
	GradeCode    string `json:"grade_code" gorm:"column:grade_code;type:varchar(20);not null;default:"`
	MemberCode   string `json:"member_code" gorm:"column:member_code;type:varchar(20);not null;default:"`
	AdminBa      string `json:"admin_ba" gorm:"column:admin_ba;type:varchar(20);not null;default:"`
	MemberChange uint8  `json:"member_change" gorm:"column:member_change;type:tinyint(1) unsigned;not null;default:0;comment:是否是会员"` // "会员信息是否改变:0 未改变 1改变"
}

func (SyncMemberCampaignToolModel) TableName() string {
	return "sync_member_campaign_tool"
}

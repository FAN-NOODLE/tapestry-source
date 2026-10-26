package excitationReq

type MedalList struct {
	Name            string `json:"name"`
	Source          int    `json:"source"`
	RuleDimensionId int    `json:"rule_dimension_id"`
	RuleCategoryId  int    `json:"rule_category_id"`
	OtherRewards    int    `json:"other_rewards"`
	StartTime       int    `json:"start_time"`
	EndTime         int    `json:"end_time"`
	RuleStartTime   int    `json:"rule_start_time"`
	RuleEndTime     int    `json:"rule_end_time"`
	Status          int    `json:"status"`
	IsEffective     int    `json:"is_effective"`
	EnterpriseId    int64  `json:"enterprise_id"`
	Page            int    `json:"page"`
	Size            int    `json:"size"`
}

type MedalCreate struct {
	EnterpriseId    int64        `json:"enterprise_id"`
	Name            string       `json:"name" binding:"required" comment:"奖章规则名称"`
	Source          int32        `json:"source" binding:"required,oneof=1 2" comment:"奖章来源 1手动上传,2规则统计"`
	RuleType        int32        `json:"rule_type" binding:"required,oneof=1 2 3" comment:"规则有效期 1长期有效,2指定时间,3仅生效一次"`
	RuleStartTime   int          `json:"rule_start_time" comment:"规则有效期-开始"`
	RuleEndTime     int          `json:"rule_end_time" comment:"规则有效期-结束"`
	MedaType        int32        `json:"meda_type" binding:"required,oneof=1 2" comment:"奖章类型 1普通奖章,2隐藏奖章"`
	MedaDescription string       `json:"meda_description" binding:"required" comment:"奖章描述"`
	OtherRewards    int32        `json:"other_rewards" binding:"required,oneof=1 2 3 4" comment:"其它奖励 1无其它奖励,2仅奖励积分,3仅奖励实物奖品,4奖励积分+实物奖品"`
	Integral        int32        `json:"integral" comment:"积分"`
	ItemName        string       `json:"item_name" comment:"物品名称"`
	RuleDimensionID int          `json:"rule_dimension_id" comment:"规则维度-规则表主键id"`
	RuleCategoryID  int          `json:"rule_category_id" comment:"规则类目-规则表主键id"`
	InviteTop       int32        `json:"invite_top" comment:"邀约业绩TOP"`
	UserIds         []string     `json:"user_ids" binding:"required,dive,required" comment:"员工id集合"`
	Styles          []MedalStyle `json:"styles" binding:"required,dive,required"`
}

type MedalUpdate struct {
	Id int64 `json:"id" binding:"required"`
	MedalCreate
}

type MedalDetail struct {
	EnterpriseId int64 `json:"enterprise_id"`
	Id           int64 `json:"id" binding:"required"`
}

type MedalStyle struct {
	Id          int64  `json:"id"`
	MedalNumber string `json:"medal_number" comment:"奖章ID"`
	TimeValue   int32  `json:"time_value" comment:"时间值"`
	TimeName    string `json:"time_name" comment:"时间名称"`
	MedalName   string `json:"medal_name" comment:"奖章名称"`
	LightIcon   string `json:"light_icon" comment:"点亮图标"`
	UnlitIcon   string `json:"unlit_icon" comment:"未点亮图标"`
}

type MedalRuleList struct {
	EnterpriseId int64 `json:"enterprise_id"`
	Pid          int   `json:"pid"`
	Page         int   `json:"page"`
	Size         int   `json:"size"`
}

type MedalStatus struct {
	Id           int64 `json:"id" binding:"required"`
	EnterpriseId int64 `json:"enterprise_id"`
	Status       int   `json:"status" binding:"required" comment:"1 启用 2 禁用"`
}

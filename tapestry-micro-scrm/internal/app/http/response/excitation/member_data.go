package excitationResp

type MemberDataListResp struct {
	Id             int64   `json:"id"`
	MemberID       int64   `json:"member_id"`
	MemberUserid   string  `json:"member_userid"`
	MemberName     string  `json:"member_name"`
	ThirdName      string  `json:"third_name"`
	ThirdId        int     `json:"third_id"`
	SpendCoins     float64 `json:"spend_coins" gorm:"comment:总花费金币"`
	TotalCoins     float64 `json:"total_coins" gorm:"comment:总获得金币"`
	Coins          float64 `json:"coins" gorm:"comment:金币"`
	PointValue     string  `json:"point_value" gorm:"comment:积分(整数)"`
	PointUpdatedAt int     `json:"point_updated_at" gorm:"comment:积分更新时间"`
	Level          string  `json:"level" gorm:"comment:等级名称"`
	LevelID        int64   `json:"level_id" gorm:"comment:level_rule表Id"`
	LevelUpdatedAt int     `json:"level_updated_at" gorm:"comment:等级更新时间"`
	MedalNum       int     `json:"medal_num" gorm:"comment:奖章数量(整数)"`
	MedalUpdatedAt int     `json:"medal_updated_at" gorm:"comment:奖章更新时间"`
	MaxLevel       string  `json:"max_level"`
	MaxLevelTime   int64   `json:"max_level_time"`
}

type MemberDataIndexResp struct {
	MemberId          int64          `json:"member_id"`
	Avatar            string         `json:"avatar"`            //头像
	Name              string         `json:"name"`              //姓名
	Level             string         `json:"level"`             //等级
	LevelIcon         string         `json:"level_icon"`        //等级图标
	PointValue        string         `json:"point_value"`       //积分
	MedalNum          int            `json:"medal_num"`         //奖章数量
	InviteValue       string         `json:"invite_value"`      //邀约能力指标
	InviteRuleValue   int            `json:"invite_rule_value"` //邀约能力指标条件值
	InviteUpdatedAt   int            `json:"invite_updated_at"`
	InviteRuleStr     string         `json:"invite_rule_str"`
	TaskValue         string         `json:"task_value"`      //任务能力指标
	TaskRuleValue     int            `json:"task_rule_value"` //任务能力指标条件值
	TaskUpdatedAt     int            `json:"task_updated_at"`
	TaskRuleStr       string         `json:"task_rule_str"`
	InteractValue     string         `json:"interact_value"`      //互动能力指标
	InteractRuleValue int            `json:"interact_rule_value"` //互动能力指标条件值
	InteractUpdatedAt int            `json:"interact_updated_at"`
	InteractRuleStr   string         `json:"interact_rule_str"`
	TrainValue        string         `json:"train_value"`      //培训能力指标
	TrainRuleValue    int            `json:"train_rule_value"` //培训能力指标条件值
	TrainUpdatedAt    int            `json:"train_updated_at"`
	TrainRuleStr      string         `json:"train_rule_str"`
	OnlineValue       string         `json:"online_value"`      //线上能力指标
	OnlineRuleValue   int            `json:"online_rule_value"` //线上能力指标条件值
	OnlineUpdatedAt   int            `json:"online_updated_at"`
	OnlineRuleStr     string         `json:"online_rule_str"`
	MaxLevel          string         `json:"max_level"`       //下个级别名称
	MaxLevelValue     int            `json:"max_level_value"` //下个级别值
	IsLastLevel       int            `json:"is_last_level"`   //1为最后一个等级
	RelationRule      []RelationRule `json:"invite_rule"`     //指标规则
	Coins             float64        `json:"coins"`           //金币
}

type RelationRule struct {
	IncentiveConfigType int    `json:"incentive_config_type"`
	RuleCondition       int    `json:"rule_condition"`
	RuleValue           int    `json:"rule_value"`
	RuleGive            string `json:"rule_give"`
}

type GradePointsResp struct {
	LevelName      string          `json:"level_name"`       //等级名称
	LevelUpdatedAt int             `json:"level_updated_at"` //升级时间
	Percentage     int64           `json:"percentage"`       //百分比
	LevelRule      []LevelRuleResp `json:"level_rule"`       //等级数据
}

type LevelRuleResp struct {
	Id        int64      `json:"id"`
	LevelName string     `json:"level_name"`
	LevelIcon string     `json:"level_icon"`
	MinValue  int        `json:"min_value"`
	MaxValue  int        `json:"max_value"`
	Selected  bool       `json:"selected"`
	Rule      []RuleResp `json:"rule"`
}

type RuleResp struct {
	IncentiveConfigType int    `json:"incentive_config_type"`
	RuleCondition       int    `json:"rule_condition"`
	RuleValue           int    `json:"rule_value"`
	RuleGive            string `json:"rule_give"`
}

type GetPointsRulesResp struct {
	EnterpriseId   int64 `json:"enterprise_id"`
	RuleCategoryId int64 `json:"rule_category_id"`
	Integral       int64 `json:"integral"`
	IntegralType   int64 `json:"integral_type"`
	RuleValue      int64 `json:"rule_value"`
	AstrictType    int64 `json:"astrict_type"` //积分上限类型 1.日，2周，3月，4年
	MaxIntegral    int64 `json:"max_integral"`
	NowIntegral    int64 `json:"now_integral"` //当前阶段已获得积分
}

type LevelRecordResp struct {
	MemberId  int64  `json:"member_id"`
	Level     string `json:"level"`
	CreatedAt int64  `json:"created_at"`
	IsMax     int64  `json:"is_max"`
	LevelIcon string `json:"level_icon"`
}

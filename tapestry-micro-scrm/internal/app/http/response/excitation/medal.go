package excitationResp

type MedalBase struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Source          int32  `json:"source"`
	RuleType        int32  `json:"rule_type"`
	RuleStartTime   int    `json:"rule_start_time"`
	RuleEndTime     int    `json:"rule_end_time"`
	MedaType        int32  `json:"meda_type"`
	MedaDescription string `json:"meda_description"`
	OtherRewards    int32  `json:"other_rewards"`
	Integral        int32  `json:"integral"`
	ItemName        string `json:"item_name"`
	RuleDimensionID int    `json:"rule_dimension_id"`
	RuleCategoryID  int    `json:"rule_category_id"`
	InviteTop       int32  `json:"invite_top"`
	CreateTime      int64  `json:"create_time"`
}

type MedalList struct {
	MedalBase
	IsEffective       int    `json:"is_effective"`
	Status            int32  `json:"status"`
	RuleDimensionName string `json:"medal_dimension_name"`
	RuleCategoryName  string `json:"rule_category_name"`
	ExistsRecord      bool   `json:"exists_record"`
}

type MedalDetail struct {
	MedalBase
	RuleDimensionName string             `json:"medal_dimension_name"`
	RuleCategoryName  string             `json:"rule_category_name"`
	Styles            []MedalDetailStyle `json:"styles"`
}

type MedalDetailStyle struct {
	Id        int64  `json:"id"`
	TimeValue int32  `json:"time_value"`
	TimeName  string `json:"time_name"`
	MedalName string `json:"medal_name"`
	LightIcon string `json:"light_icon"`
	UnlitIcon string `json:"unlit_icon"`
}

type MedalRuleList struct {
	Id          int64  `json:"id"`
	Pid         int    `json:"pid"`
	RuleName    string `json:"rule_name"`
	TimeType    int32  `json:"time_type"`
	Description string `json:"description"`
}

type MedalWall struct {
	Total   int64           `json:"total"`
	Gained  int64           `json:"gained" desc:"已获得奖章"`
	NowYear int64           `json:"now_year"`
	List    []YearMedalWall `json:"list"`
}
type YearMedalWall struct {
	Year     int64           `json:"year"`
	Ordinary []MedalWallList `json:"ordinary" desc:"普通奖章"`
	Rare     []MedalWallList `json:"rare" desc:"稀有勋章"`
}

type FiscalYear struct {
	Year      int64 `json:"year"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

type MedalWallList struct {
	RuleId     int64  `json:"rule_id"`
	IsShow     int32  `json:"is_show" desc:"是否已查看"`
	Icon       string `json:"icon" desc:"点亮图标"`
	RuleGained int64  `json:"rule_gained" desc:"已获得数量"`
	RuleTotal  int64  `json:"rule_total" desc:"总数"`
	RuleName   string `json:"rule_name" desc:"规则名称"`
	TimeIndex  int64  `json:"time_index" desc:"当前所属时间坐标"`
	GetTime    int64  `json:"get_time" desc:"奖章获取时间"`
}

type MedalWallDetail struct {
	RuleId          int64                  `json:"rule_id"`
	RuleName        string                 `json:"rule_name" desc:"规则名称"`
	RuleDescription string                 `json:"rule_description" desc:"规则描述"`
	TopNum          int32                  `json:"top_num" desc:"top榜数"`
	Integral        int32                  `json:"integral" desc:"积分"`
	ItemName        string                 `json:"item_name" desc:"物品名称"`
	MedalStyle      []MedalWallDetailStyle `json:"medal_style"`
}

type MedalWallDetailStyle struct {
	StyleId             int64   `json:"style_id" desc:"奖章id"`
	MedalName           string  `json:"medal_name" desc:"奖章名称"`
	LightIcon           string  `json:"light_icon" desc:"点亮图标"`
	UnlitIcon           string  `json:"unlit_icon" desc:"未点亮图标"`
	ReceiveTime         string  `json:"receive_time" desc:"获得时间"`
	CreateTime          int64   `json:"create_time" desc:""`
	ReceiveNumberPeople int64   `json:"receive_number_people" desc:"获得人数"`
	HasExceeded         float64 `json:"has_exceeded" desc:"已超过"`
}

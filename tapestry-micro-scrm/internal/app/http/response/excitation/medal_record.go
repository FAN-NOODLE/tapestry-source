package excitationResp

type MedalRecordList struct {
	Id                int64  `json:"id"`
	MemberId          int64  `json:"member_id"`
	MemberUserId      string `json:"member_user_id"`
	MemberName        string `json:"member_name"`
	RuleName          string `json:"rule_name"`
	MedalName         string `json:"medal_name"`
	Source            int32  `json:"source"`
	RuleDimensionID   int    `json:"rule_dimension_id"`
	RuleCategoryID    int    `json:"rule_category_id"`
	RuleDimensionName string `json:"rule_dimension_name"`
	RuleCategoryName  string `json:"rule_category_name"`
	Integral          int32  `json:"integral"`
	ItemName          string `json:"item_name"`
	DepartmentName    string `json:"department_name"`
	CreatedAt         int64  `json:"created_at"`
	Coins             string `json:"coins"`
	Ranking           string `json:"ranking"`
	RankingValue      string `json:"ranking_value"`
	Description       string `json:"description"`
	InviteTop         int    `json:"invite_top"`
}

type MedalRecordAlreadyGot struct {
	Number       int64                        `json:"number"`
	MedalList    []MedalRecordAlreadyGotMedal `json:"medal_list"`
	NewMedalSign int64                        `json:"new_medal_sign"`
}

type MedalRecordAlreadyGotMedal struct {
	MedalName    string `json:"medal_name"`
	RecordId     int64  `json:"record_id"`
	RuleName     string `json:"rule_name"`
	ReceivedTime int64  `json:"received_time"`
	LightIcon    string `json:"light_icon"`
}

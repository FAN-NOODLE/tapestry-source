package excitationResp

type IntegralRuleList struct {
	Id           int64  `json:"id"`
	EnterpriseId int64  `json:"enterprise_id"`
	Type         int64  `json:"type"`
	TypeName     string `json:"type_name"`
	Status       int64  `json:"status"`
	RelationId   int64  `json:"relation_id"`
	RuleName     string `json:"rule_name"`
	RuleNum      int64  `json:"rule_num"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type Configuration struct {
	Id   int64             `json:"id"`
	Name string            `json:"name"`
	Rule ConfigurationRule `json:"rule"`
}

type ConfigurationRule struct {
	RuleIf   []RuleIfValue `json:"rule_if"`
	RuleIf1  []RuleIfValue `json:"rule_if1"`
	Selected int64         `json:"selected"`
}

type RuleIfValue struct {
	Value int64  `json:"value"`
	Name  string `json:"name"`
}

//IntegralRule 返回的规则信息节结构
type IntegralRule struct {
	ID           int64                `json:"id"`
	EnterpriseID int64                `json:"enterprise_id"`
	Type         int64                `json:"type"`
	RuleName     string               `json:"rule_name"`
	Detail       []IntegralRuleDetail `json:"rule_detail"`
}

//IntegralRuleDetail 返回的规则信息节结构
type IntegralRuleDetail struct {
	ID             int64                     `json:"id"`
	EnterpriseID   int64                     `json:"enterprise_id"`
	RuleID         int64                     `json:"rule_id"`
	RuleCategoryID int64                     `json:"rule_category_id"`
	ValidityType   int64                     `json:"validity_type"`
	StartTime      int64                     `json:"start_time"`
	EndTime        int64                     `json:"end_time"`
	IsAstrict      int64                     `json:"is_astrict"`
	AstrictType    int64                     `json:"astrict_type"`
	MaxIntegral    int64                     `json:"max_integral"`
	IsDiffStore    int64                     `json:"is_diff_store"`
	StoreType      int64                     `json:"store_type"`
	DetailValue    []IntegralRuleDetailValue `json:"detail_value"`
}

//IntegralRuleDetailValue 规则值
type IntegralRuleDetailValue struct {
	Integral     int64 `json:"integral"`
	IntegralType int64 `json:"integral_type"`
	RuleIf       int64 `json:"rule_if"`
	RuleIf1      int64 `json:"rule_if1"`
	RuleValue    int64 `json:"rule_value"`
	StoreType    int64 `json:"store_type"`
}

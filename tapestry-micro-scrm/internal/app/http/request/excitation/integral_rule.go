package excitationReq

//RuleList 列表
type RuleList struct {
	RuleName     string `json:"rule_name"`
	Type         int    `json:"type"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	EnterpriseId int64  `json:"enterprise_id"  binding:"required"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

//RuleCreate 添加规则
type RuleCreate struct {
	EnterpriseId int64              `json:"enterprise_id" binding:"required"`
	Type         int64              `json:"type" binding:"required" comment:"维度类型"`
	RuleName     string             `json:"rule_name" binding:"required" comment:"规则名称"`
	RuleDetail   []RuleDetailCreate `json:"rule_detail"`
}

//RuleDetailCreate 添加规则
type RuleDetailCreate struct {
	EnterpriseId   int64             `json:"enterprise_id"`
	RuleID         int64             `json:"rule_id"  binding:"required"`
	RuleCategoryID int64             `json:"rule_category_id" binding:"required" comment:"规则类型（1:邀约数每次 2:邀约到店率每周 3:邀约成功率每月 4:邀约成功每次 5:邀约占比每周 6:邀约到店确认每次 7:任务完成每次 8:任务完成率每月 9:任务触达每人 10: 新增好友数每人 11:客户流失数每人 12:聊天总数每月 13:已回复聊天占比每月 14:完成播种按次 15:积分等比兑换 16:线上销售业绩率每周)"`
	ValidityType   int64             `json:"validity_type" binding:"required" comment:"1.长期 2.指定时间"`
	StartTime      int64             `json:"start_time" comment:"开始时间"`
	EndTime        int64             `json:"end_time" comment:"结束时间"`
	IsAstrict      int64             `json:"is_astrict"`
	AstrictType    int64             `json:"astrict_type" comment:"积分上限类型 1.日，2周，3月，4年 "`
	MaxIntegral    int64             `json:"max_integral" comment:"结束时间"`
	IsDiffStore    int64             `json:"is_diff_store" comment:"是否区分门店：0 否 1 区分奥莱零售"`
	RuleValue      []RuleDetailValue `json:"rule_value"`
}

type RuleDetailValue struct {
	Integral     int64 `json:"integral" binding:"required" comment:"奖励积分"`
	IntegralType int64 `json:"integral_type" binding:"required" comment:"1.加积分 2.减积分"`
	RuleIf       int64 `json:"rule_if" comment:"1:次数条件 2:高于3：达到"`
	RuleIf1      int64 `json:"rule_if1" comment:"0.默认无规则 1：数值2：平均值"`
	RuleValue    int64 `json:"rule_value" comment:"条件值"`
	StoreType    int64 `json:"store_type" comment:"0 默认 1 奥莱 2 零售"`
}

//RuleDetail 详情
type RuleDetail struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required"`
	Id           int64 `json:"id" binding:"required" comment:"规则ID"`
}

//RuleDestroy 删除
type RuleDestroy struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required"`
	Id           int64 `json:"id" binding:"required" comment:"规则ID"`
}

//RuleUpdate 编辑
type RuleUpdate struct {
	ID           int64              `json:"id" binding:"required"`
	EnterpriseId int64              `json:"enterprise_id" binding:"required"`
	RuleName     string             `json:"rule_name" binding:"required" comment:"规则名称"`
	RuleDetail   []RuleDetailCreate `json:"rule_detail"`
}

//ParameterConfiguration 配置信息
type ParameterConfiguration struct {
	Type int `json:"type"`
}

// RuleMany 列表
type RuleMany struct {
	EnterpriseId int64   `json:"enterprise_id"  binding:"required"`
	Types        []int32 `json:"types"`
}

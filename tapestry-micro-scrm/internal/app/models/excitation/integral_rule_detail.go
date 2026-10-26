package excitationModels

import "micro-scrm/internal/app/models"

const (
	RuleCategoryId1  = iota + 1 // 邀约数每次
	RuleCategoryId2             // 邀约到店率每周
	RuleCategoryId3             // 邀约成功率每月（往前推一月）
	RuleCategoryId4             // 邀约成功每次
	RuleCategoryId5             // 邀约占比每月
	RuleCategoryId6             // 邀约到店确认每次
	RuleCategoryId7             // 任务完成每次
	RuleCategoryId8             // 任务完成率每月
	RuleCategoryId9             // 任务触达每人
	RuleCategoryId10            // 新增好友数每人
	RuleCategoryId11            // 客户流失数每月
	RuleCategoryId12            // 聊天总数每月
	RuleCategoryId13            // 已回复聊天占比每月
	RuleCategoryId14            // 完成播种按次
	RuleCategoryId15            // 积分等比兑换
	RuleCategoryId16            // 线上销售业绩率每月
	RuleCategoryId17            // 点击新老客按钮
	RuleCategoryId18            // 点击邀约的是否到店购买
	RuleCategoryId19            // 线上销售金额
)
const (
	RuleIfTimes = iota + 1 // 次数条件
	RuleIfGT               // 高于
	RuleIfEgt              // 达到
)
const (
	RuleIf1Value   = iota + 1 // 数值（包含百分比）
	RuleIf1Average            // 平均值
)
const (
	IntegralTypeInc = iota + 1 // 加积分
	IntegralTypeDec            // 减积分
)

const (
	ValidityTypeLong     = iota + 1 // 长期
	ValidityTypeDeadline            // 限时
)

const (
	NoAstrict = iota + 1 //无上限
	IsAstrict            //有上限
)

const (
	Astrict_Day = iota + 1
	Astrict_Week
	Astrict_Month
	Astrict_Year
)

const (
	INTEGRAL_STORE_ALL = iota
	INTEGRAL_STORE_OLE
	INTEGRAL_STORE_REATIL
)

type ExcitationIntegralRuleDetailModel struct {
	models.MicroBaseModel

	ID             int64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) auto_increment"`
	EnterpriseID   int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);default:0"`
	RuleID         int64 `json:"rule_id" gorm:"column:rule_id;type:bigint(20);not null;comment:规则表id"`
	RuleCategoryID int64 `json:"rule_category_id" gorm:"column:rule_category_id;type:tinyint(4);comment:规则类型（1:邀约数每次 2:邀约到店率每周 3:邀约成功率每月 4:邀约成功每次 5:邀约占比每周 6:邀约到店确认每次 7:任务完成每次 8:任务完成率每月 9:任务触达每人 10: 新增好友数每人 11:客户流失数每人 12:聊天总数每月 13:已回复聊天占比每月 14:完成播种按次 15:积分等比兑换 16:线上销售业绩率每周)"`
	ValidityType   int64 `json:"validity_type" gorm:"column:validity_type;type:tinyint(4);default:0;comment:1.长期 2.指定时间"`
	StartTime      int64 `json:"start_time" gorm:"column:start_time;type:int(11);comment:开始时间"`
	EndTime        int64 `json:"end_time" gorm:"column:end_time;type:int(11);comment:结束时间"`
	IsDiffStore    int64 `json:"is_diff_store" gorm:"column:is_diff_store;type:tinyint(4);default:0;comment:是否区分门店：0通用规则 1区分奥莱零售'"`
	StoreType      int64 `json:"store_type" gorm:"column:store_type;type:tinyint(4);default:0;comment:是否区分门店：0通用规则 1奥莱2 零售'"`
	Integral       int64 `json:"integral" gorm:"column:integral;type:int(11);default:0;comment:奖励积分"`
	IntegralType   int64 `json:"integral_type" gorm:"column:integral_type;type:tinyint(4);default:0;comment:1.加积分 2.减积分"`
	RuleIf         int64 `json:"rule_if" gorm:"column:rule_if;type:tinyint(4);default:0;comment:1:次数条件 2:高于3：达到"`
	RuleIf1        int64 `json:"rule_if1" gorm:"column:rule_if1;type:tinyint(4);comment:0.默认无规则 1：数值2：平均值 "`
	RuleValue      int64 `json:"rule_value" gorm:"column:rule_value;type:int(11);default:0;comment:条件值"`
	IsAstrict      int64 `json:"is_astrict" gorm:"column:is_astrict;type:tinyint(4);default:1;comment:积分上限 1.无上限 2.有上限"`
	AstrictType    int64 `json:"astrict_type" gorm:"column:astrict_type;type:tinyint(4);default:1;comment:积分上限类型 1.日，2周，3月，4年"`

	MaxIntegral int64 `json:"max_integral" gorm:"column:max_integral;type:int(11);default:0;comment:最大上线积分"`

	// 关联规则主表配置
	Rule ExcitationIntegralRuleModel `gorm:"foreignKey:RuleID;references:ID"`
}

func (ExcitationIntegralRuleDetailModel) TableName() string {
	return "excitation_integral_rule_detail"
}

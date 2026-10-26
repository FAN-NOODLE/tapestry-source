package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationMedalRuleModel struct {
	models.MicroBaseModel

	RuleName    string `json:"rule_name" gorm:"column:rule_name;type:varchar(100);not null;comment:规则名称"`
	ShortName   string `json:"short_name" gorm:"column:short_name;type:varchar(100);not null;comment:短名"`
	Description string `json:"description" gorm:"column:description;type:varchar(100);not null;comment:描述"`
	TimeType    int32  `json:"time_type" gorm:"column:time_type;type:tinyint(1);not null;default:0;comment:时间类型 0无,1月,2季度,3年"`
	Pid         int    `json:"pid" gorm:"column:pid;type:int(11);not null;default:1;comment:父级id"`
	Status      int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:1;comment:状态 1启用,2禁用"`
}

// TableName returns the table name of the WjMedalRule model
func (ExcitationMedalRuleModel) TableName() string {
	return "excitation_medal_rule"
}

// 状态 1启用,2禁用
const (
	MedalRuleStatusEnable = iota + 1
	MedalRuleStatusDisable
)

// 时间类型 0无,1月,2季度,3年
const (
	MedalRuleTimeTypeMonth = iota + 1
	MedalRuleTimeTypeQuarter
	MedalRuleTimeTypeYear
)

// 邀约业绩月榜
// 邀约业绩季度榜
// 邀约业绩年度榜
// 销售业绩月榜
// 销售业绩季度榜
// 销售业绩年度榜
// 线上销售月榜
// 线上销售季度榜
// 线上销售年度榜
// 任务完成奖章
// 任务触达奖章
// 平均首次回复时长
// 好友人数
// 互动达人
// 播种达人
// 好学生
// 考试狂人
// 学习委员

// invite_month
// invite_season
// invite_year
// sale_month
// sale_season
// sale_year
// online_sale_month
// online_sale_season
// online_sale_year
// task_finish
// task_touch
// response_time
// friend_number
// interactive
// sown
// good_student
// exam_maniac
// study_member

const (
	MedalRuleInviteMonthlyList string = "invite_month"
	MedalRuleInviteSeasonList         = "invite_season"
	MedalRuleInviteYearList           = "invite_year"

	MedalRuleTaskFinishList = "task_finish"
	MedalRuleTaskTouchList  = "task_touch"

	MedalRuleResponseTimeList = "response_time"
	MedalRuleFriendNumberList = "friend_number"
	MedalRuleInteractiveList  = "interactive"
	MedalRuleSownList         = "sown"

	MedalRuleSaleMonth        = "sale_month"
	MedalRuleSaleSeason       = "sale_season"
	MedalRuleSaleYear         = "sale_year"
	MedalRuleOnlineSaleMonth  = "online_sale_month"
	MedalRuleOnlineSaleSeason = "online_sale_season"
	MedalRuleOnlineSaleYear   = "online_sale_year"

	MedalRuleGoodStudent = "good_student"
	MedalRuleExamManiac  = "exam_maniac"
	MedalRuleStudyMember = "study_member"

	MedalStudyCoursesList = "study_courses_list" //行走的百科全书-学课数量榜单
	MedalStudyTimeList    = "study_time_list"    //高质量专家学员-学习时长榜单
	MedalNewCoursesList   = "new_courses_list"   //高觉悟三好学生-新品课程完成榜单
)

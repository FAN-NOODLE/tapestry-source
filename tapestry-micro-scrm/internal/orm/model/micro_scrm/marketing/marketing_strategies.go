package marketing

import "micro-scrm/internal/orm/model"

// 营销策略表模型
type MarketingStrategiesModel struct {
	model.BaseModel

	StrategyName    string `gorm:"column:strategy_name;type:varchar(64);not null;comment:策略名称"`
	CustomerGroupID int64  `gorm:"column:customer_group_id;type:bigint(20);not null;index:customer_group_id_index;default:0;comment:客户分组ID"`
	AdminID         int64  `gorm:"column:admin_id;type:bigint(20);not null;index:admin_id_index;default:0;comment:操作人ID"`
	StrategyType    int32  `gorm:"column:strategy_type;type:tinyint(1);not null;default:1;comment:策略类型: 1=客户策略，2=内置策略"`
	StrategyMethod  int32  `gorm:"column:strategy_method;type:tinyint(1);not null;default:1;comment:执行方式: 1=长期执行，2=立即执行，3=定时执行"`
	StartTime       int    `gorm:"column:start_time;type:int(11);not null;default:0;comment:执行开始时间"`
	EndTime         int    `gorm:"column:end_time;type:int(11);not null;default:0;comment:执行结束时间"`
	StrategyPolicy  int32  `gorm:"column:strategy_policy;type:tinyint(1);not null;default:1;comment:触达规则：1=通知一个渠道，2=通知所有渠道"`
	StrategyStatus  int32  `gorm:"column:strategy_status;type:tinyint(1);not null;default:1;comment:状态: 1=未开始，2=投放中，3=暂停中，4=已结束"`
	Status          int32  `gorm:"column:status;type:tinyint(1);not null;default:1;comment:状态: 1=禁用，2=启用"`
	IsFinished      int32  `gorm:"column:is_finished;type:tinyint(1);not null;default:1;comment:是否完成: 1=未完成，2=已完成"`
	SettingDay      int32  `gorm:"column:setting_day;type:tinyint(1);not null;default:7;comment:设置默认天数"`
	IsNotice        int32  `gorm:"column:is_notice;type:tinyint(1);not null;default:1;comment:是否已通知发送短信提醒"`
	HistoryStatus   int32  `gorm:"column:history_status;type:tinyint(1);not null;default:0;comment:长期不限存量数据是否已执行，1=执行 0=未执行"`
	Mobile          string `gorm:"column:mobile;type:varchar(20);not null;comment:联系方式"`
	AbnormalPause   int32  `gorm:"column:abnormal_pause;type:tinyint(1) unsigned;not null;default:0;comment:是否为异常情况暂停的策略，1:是,0:否"`

	MarketingStrategyRecords      []MarketingStrategyRecordsModel    `gorm:"foreignKey:StrategyID;references:ID" desc:"营销策略记录表"`
	MarketingStrategyTriggerCount MarketingStrategyTriggerCountModel `gorm:"foreignKey:StrategyID;references:ID" desc:"营销策略预触达人数"`
}

func (MarketingStrategiesModel) TableName() string {
	return "marketing_strategies"
}

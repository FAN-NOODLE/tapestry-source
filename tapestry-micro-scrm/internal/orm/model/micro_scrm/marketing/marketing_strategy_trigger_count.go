package marketing

type MarketingStrategyTriggerCountModel struct {
	StrategyID int64 `gorm:"primaryKey;column:strategy_id;type:bigint(20);comment:策略id"`
	Count      int   `gorm:"column:count;type:int(10);not null;default:0;comment:可触达人数"`
}

// TableName returns the table name of the struct_name model
func (MarketingStrategyTriggerCountModel) TableName() string {
	return "marketing_strategy_trigger_count"
}

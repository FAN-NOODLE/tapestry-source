package marketing

// 策略执行的结果
type StrategyExecResult struct {
	EnterpriseID int64 `gorm:"column:enterprise_id"`
	StrategyID   int64 `gorm:"column:strategy_id"`
	Total        int   `gorm:"column:total"`
	NumSuccess   int   `gorm:"column:num_success"`
	NumFail      int   `gorm:"num_fail"`
	ChannelType  int8  `gorm:"channel_type"`
}

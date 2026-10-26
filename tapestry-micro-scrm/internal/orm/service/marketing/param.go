package marketing

// 营销策略下载参数
type MarketingExportParam struct {
	EnterpriseID, ID   int64
	StartTime, EndTime int
}

// 营销策略结果查询参数
type MarketingResultQueryParam struct {
	EnterpriseID int64
	StrategyID   int64
	StartTime    int64
	EndTime      int64
}

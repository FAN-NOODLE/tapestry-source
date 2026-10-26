package customerReq

type CalcRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Day          int64 `json:"day" desc:"0 今日 1昨天 近n天"`
}

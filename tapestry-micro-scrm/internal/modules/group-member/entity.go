package groupmember

type GroupParam struct {
	Day          int64 `json:"day" desc:"开始时间"`
	EndAt        int64 `json:"end_at" desc:"结束时间" db:"end_at"`
	StartAt      int64 `json:"start_at" desc:"开始时间" db:"start_at"`
	EnterpriseId int64 `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
}

type AsRes struct {
	RelationNum int64 `json:"relation_num" db:"relation_num"`
}

type GroupRes struct {
	GroupNum     int64 `json:"group_num"`
	GroupNumThan int64 `json:"group_num_than"`
	ReUv         int64 `json:"re_group_uv"`
	ReUvThan     int64 `json:"re_group_uv_than"`
	AddUv        int64 `json:"add_group_uv"`
	AddUvThan    int64 `json:"add_group_uv_than"`
	LossUv       int64 `json:"loss_group_uv"`
	LossUvThan   int64 `json:"loss_group_uv_than"`
}

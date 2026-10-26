package source

type VideoDetailRequest struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type VideoUpdateRequest struct {
	ID            int64        `json:"id"`
	EnterpriseID  int64        `json:"enterprise_id"`
	AdminId       int64        `json:"admin_id"`
	StartTime     int64        `json:"start_time"`
	EndTime       int64        `json:"end_time"`
	CategoryIds   []int64      `json:"category_ids"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}

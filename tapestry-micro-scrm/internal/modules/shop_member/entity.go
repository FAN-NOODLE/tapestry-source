package shop_member

// DDCRequestBase
type DDCRequestBase struct {
	Method   string `json:"method"`
	Platform string `json:"platform"`
}

// SyncShopMembersRequest 同步商城会员接口请求参数
type SyncShopMemberRequest struct {
	DDCRequestBase
	ShopMemberInfo ShopMemberEntity `json:"data"`
}

// ShopMemberEntity 商城会员结构体
type ShopMemberEntity struct {
	ID                     int64  `json:"id" db:"id" desc:"会员ID"`
	EnterpriseID           int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UnionID                string `json:"union_id" db:"union_id" desc:"用户unionID"`
	Number                 string `json:"number" db:"number" desc:"编号"`
	LevelID                int64  `json:"level_id" db:"level_id" desc:"等级id"`
	AvailableIntegralTotal int64  `json:"available_integral_total" db:"available_integral_total" desc:"可用积分总数量"`
	YearIntegralTotal      int64  `json:"year_integral_total" db:"year_integral_total" desc:"本年积分总数量"`
	HistoryIntegralTotal   int64  `json:"history_integral_total" db:"history_integral_total" desc:"历史积分总数量"`
	QrCodeURL              string `json:"qr_code_url" db:"qr_code_url" desc:"二维码链接"`
	BarCodeURL             string `json:"bar_code_url" db:"bar_code_url" desc:"条形码链接"`
	OpenTime               int64  `json:"open_time" db:"open_time" desc:"当前等级开通时间"`
	OverdueTime            int64  `json:"overdue_time" db:"overdue_time" desc:"当前等级过期时间"`
	CreatedAt              int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt              int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt              int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// SyncShopMemberLevelRequest 同步商城会员等级接口请求参数
type SyncShopMemberLevelRequest struct {
	DDCRequestBase
	ShopMemberLevelInfo ShopMemberLevelEntity `json:"data"`
}

// ShopMemberLevelEntity 商城会员等级结构体
type ShopMemberLevelEntity struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	LevelID      int64  `json:"level_id" db:"level_id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Level        int8   `json:"level" db:"level" desc:"等级"`
	LevelName    string `json:"level_name" db:"level_name" desc:"等级名称"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MemberLevelListRequest 会员等级列表请求结构体
type MemberLevelListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// MemberLevelListResponse 会员等级列表返回结构体
type ConfigGuideListResponse struct {
	List []*ShopMemberLevelEntity `json:"list"`
}

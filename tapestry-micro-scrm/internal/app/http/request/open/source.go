package openReq

// SearchListRequest 对外链接列表请求
type SearchListRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	UserId       string `json:"user_id"`
	ThirdIds     []int  `json:"third_ids"`
	ProductSku   string `json:"product_sku"`
	ProductColor string `json:"product_color"`
	ProductName  string `json:"product_name"`
}

// SourceTreeRequest 素材分组树请求
type SourceTreeRequest struct {
	EnterpriseID int64 `json:"enterprise_id"` // 企业ID
}

// SourceOpenReq 素材开放接口请求
type SourceOpenReq struct {
	EnterpriseID      int64 `json:"enterprise_id"` // 企业ID
	Type              int64 `json:"type"`
	UpdatedAt         int64 `json:"updated_at"`         // 更新时间戳，获取该时间戳之后更新的数据
	Page              int   `json:"page"`               // 页码
	Size              int   `json:"size"`               // 每页数量
	IncludeCategories bool  `json:"include_categories"` // 是否返回素材标签
}

// SourceListRequest 素材列表请求
// UserId为空时查询全企业素材，非空时查询该用户有权限的素材
type SourceListRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	UpdatedAt    int64  `json:"updated_at"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

// SearchAllRequest 合并检索请求（通过商品条件一次性检索所有素材类型）
// ProductSku / ProductColor / ProductName 至少填一个
type SearchAllRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	ThirdIds     []int  `json:"third_ids"`
	ProductSku   string `json:"product_sku"`
	ProductColor string `json:"product_color"`
	ProductName  string `json:"product_name"`
}

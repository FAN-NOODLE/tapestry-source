package microgoods

// SpuListRequest spu列表请求参数
type SpuListRequest struct {
	EnterpriseID int64   `json:"enterprise_id"`
	SpuIds       []int64 `json:"spu_ids"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
}

// SpuListResponse spu列表返回格式
type SpuListResponse struct {
	ResponseBase
	Data struct {
		List  []SpuInfo `json:"list"`
		Total int       `json:"total"`
		Page  int       `json:"page"`
		Size  int       `json:"size"`
	} `json:"data"`
}

// SpuInfo spu基本信息
type SpuInfo struct {
	ID           int64  `json:"id"`
	SpuName      string `json:"spu_name"`
	EnterpriseID int64  `json:"enterprise_id"`
	SpuSource    struct {
		URL string `json:"url"`
	} `json:"spu_source"`
}

package baiqiuSdk

// BatchSearchRequest 批量搜索请求
type BatchSearchRequest struct {
	Products []ProductQuery `json:"products"`
}

// ProductQuery 产品查询项
type ProductQuery struct {
	ProductCode string `json:"product_code"`
	Color       string `json:"color"`
}

// BatchSearchResponse 批量搜索响应
type BatchSearchResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    *BatchData `json:"data"`
}

// BatchData 批量查询数据
type BatchData struct {
	TotalQuery    int        `json:"total_query"`
	FoundCount    int        `json:"found_count"`
	NotFoundCount int        `json:"not_found_count"`
	Products      []Product  `json:"products"`
	NotFound      []NotFound `json:"not_found"`
}

// Product 产品信息
type Product struct {
	Query       string `json:"query"`
	Key         string `json:"key"`
	ProductCode string `json:"product_code"`
	Color       string `json:"color"`
	ImagePath   string `json:"image_path"`
}

// NotFound 未找到的产品
type NotFound struct {
	Query  string `json:"query"`
	Reason string `json:"reason"`
}

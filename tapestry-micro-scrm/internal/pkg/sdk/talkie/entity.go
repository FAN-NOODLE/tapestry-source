package talkie

// TokenRequest Token获取请求
type TokenRequest struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

// TokenResponse Token获取响应
type TokenResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AccessToken string `json:"accessToken"`
		Expires     int    `json:"expires"`
	} `json:"data"`
	Detail *TokenErrorDetail `json:"detail,omitempty"`
}

// TokenErrorDetail Token错误详情
type TokenErrorDetail struct {
	Msg        string `json:"msg"`
	StatusCode string `json:"statusCode"`
}

// ImageSearchRequest 图片搜索请求
type ImageSearchRequest struct {
	ImageURL   string `json:"imageUrl"`             // 图片URL地址
	CategoryID int    `json:"categoryId,omitempty"` // 商品类目ID（可选）
	Num        int    `json:"num"`                  // 每个区域返回产品数量，默认5，最多10
}

// ImageSearchResponse 图片搜索响应
type ImageSearchResponse struct {
	Success bool                      `json:"success"`
	Data    []ImageSearchRegionResult `json:"data"`
	Detail  *ImageSearchErrorDetail   `json:"detail,omitempty"`
}

// ImageSearchErrorDetail 图片搜索错误详情
type ImageSearchErrorDetail struct {
	Msg        string `json:"msg"`
	StatusCode string `json:"statusCode"`
}

// ImageSearchRegionResult 图片主体区域搜索结果
type ImageSearchRegionResult struct {
	Region  string               `json:"region"`  // 主体区域坐标，格式：x1,x2,y1,y2
	PicList []ImageSearchPicInfo `json:"picList"` // 相似产品列表
}

// ImageSearchPicInfo 相似图片产品信息
type ImageSearchPicInfo struct {
	ProductID  string  `json:"productId"`  // 商品ID [款色]，示例：65129_GDMEL
	PicName    string  `json:"picName"`    // 图片名称
	CategoryID int     `json:"categoryId"` // 商品类目
	Score      float64 `json:"score"`      // 相似度分数（0-1）
}

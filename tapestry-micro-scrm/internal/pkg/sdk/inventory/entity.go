package inventory

type ParamsReq struct {
	Country     string     `json:"country"`
	LocCode     []string   `json:"locCode,omitempty"`
	SkuItem     []string   `json:"skuItem,omitempty"`
	UpcCode     []string   `json:"upcCode,omitempty"`
	InvLookType string     `json:"invLookType,omitempty"`
	SummaryType string     `json:"summaryType,omitempty"`
	SortOption  SortOption `json:"sortOption"`
	// ExpectedSegment carries caller context and is not sent upstream; business
	// assembly filters usable DCL rows without treating other segments as invalid.
	ExpectedSegment string `json:"-"`
}

type SortOption struct {
	OnHand    int `json:"onHand,omitempty"`
	InTransit int `json:"inTransit,omitempty"`
}

// Response  三方库存接口返回结果固定参数
type Response struct {
	DetailList []RespDetail `json:"detailList"`
	Header     *RespHeader  `json:"header"`
	Error      interface{}  `json:"error"`
}
type RespDetail struct {
	Country       string `json:"country"`
	UpcCode       string `json:"upcCode"`
	InTransit     int    `json:"inTransit"`
	LocCode       string `json:"locCode"`
	OnHand        int    `json:"onHand"`
	OnHold        int    `json:"onHold"`
	TotalOnHand   int    `json:"totalOnHand"`
	Segment       string `json:"segment"`
	SkuAttr1      string `json:"skuAttr1"`
	SkuAttr2      any    `json:"skuAttr2"`
	OwnStoreStock int    `json:"ownStoreStock"`
}

type RespHeader struct {
	InTransitSum   int    `json:"inTransitSum,omitempty"`
	OnHandSum      int    `json:"onHandSum,omitempty"`
	OnHoldSum      int    `json:"onHoldSum,omitempty"`
	ResultCode     string `json:"resultCode"`
	ResultMessage  string `json:"resultMessage"`
	TotalCount     int    `json:"totalCount,omitempty"`
	TotalOnHandSum int    `json:"totalOnHandSum,omitempty"`
}

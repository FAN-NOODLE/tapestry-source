package inventorySw

type ParamsReq struct {
	SearchOtherStore bool   `json:"search_other_store"`
	StoreCode        string `json:"store_code"`
	StkName          string `json:"stk_name,omitempty"` //款号名称
	StkCode          string `json:"stk_code,omitempty"` //款号
	Plu              string `json:"plu,omitempty"`      //stk_code*color*size
}

// Response  三方库存接口返回结果固定参数
type Response struct {
	Data []RespDetail `json:"data"`
	Msg  string       `json:"msg"`
	Code int          `json:"code"`
}
type RespDetail struct {
	StoreCode string  `json:"store_code"`
	StoreName string  `json:"store_name"`
	StkCode   string  `json:"stk_code"` //款号
	StkName   string  `json:"stk_name"` //款号名称
	Color     string  `json:"color"`
	Size      string  `json:"size"`
	Plu       string  `json:"plu"`      //stk_code*color*size
	ListPrc   float64 `json:"list_prc"` //零售标价
	RtlPrc    float64 `json:"rtl_prc"`  //优先取该门店价格本上的标价,如果没有取商品目录上零售标价
	Qty       int     `json:"qty"`      //在手库存
	TraQty    float64 `json:"tra_qty"`  //在途库存
	PicUrl    string  `json:"pic_url"`
	StkAttri  string  `json:"stk_attri"`  //门店类型名称
	StoreTel  string  `json:"store_tel"`  //门店电话
	StoreCity string  `json:"store_city"` //门店城市
}

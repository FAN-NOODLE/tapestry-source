package inventoryResp

type ListResult struct {
	List            []SkuInfo `json:"list"`
	InventoryStatus string    `json:"inventory_status"`
	FailedUPCs      []string  `json:"failed_upcs"`
	DCLStatus       string    `json:"dcl_status"`
}

type SkuInfo struct {
	SkuItem  string `json:"sku_item"`
	ImageUrl string `json:"image_url"`
	Color    string `json:"color"`
	Style    any    `json:"style"`
	UpcInfo
}

type UpcInfo struct {
	UpcCode        string `json:"upc_code"`
	InTransit      int    `json:"in_transit"`
	LocCode        string `json:"loc_code"`
	ThirdStoreCode string `json:"third_store_code"`
	OnHand         int    `json:"on_hand"`
	InventoryName  string `json:"inventory_name"`
	StoreName      string `json:"store_name"`
	City           string `json:"city"`
	Phone          string `json:"phone"`
	DclOnHand      int    `json:"dcl_on_hand"`
}

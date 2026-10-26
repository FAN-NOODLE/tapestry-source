package inventoryReq

import (
	"micro-scrm/internal/app/http/request"
	inventorySw "micro-scrm/internal/pkg/sdk/inventory_sw"
)

type ListReq struct {
	EnterpriseId int64    `json:"enterprise_id" binding:"required"`
	SkuItem      string   `json:"sku_item"`
	UpcCode      string   `json:"upc_code"`
	BaStoreCode  string   `json:"store_code"`     //登录用户门店标识
	StoreCodeArr []string `json:"store_code_arr"` //前端传参门店标识
	MemberId     int64    `json:"member_id"`      // 网仓导购获取大仓的库存
	Scan         bool     `json:"scan"`           // 是否扫码查询
}

type DcListReq struct {
	ListReq
	MemberId int64 `json:"member_id"` // 网仓导购获取大仓的库存
}

type KeywordsReq struct {
	request.EnterpriseBaseReq
	Keyword  string `json:"keyword"`
	MemberId int64  `json:"member_id" binding:"required"`
}

type DelKeywordsReq struct {
	request.EnterpriseBaseReq
	Ids      []int64 `json:"ids"`
	MemberId int64   `json:"member_id" binding:"required"`
	IsAll    int     `json:"is_all" comment:"0:非全部删除，1:sku_item全部删除,2:upc_code全部删除"`
}

type ListSwReq struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required"`
	MemberId     int64 `json:"member_id"` // 网仓导购获取大仓的库存
	Scan         bool  `json:"scan"`      // 是否扫码查询
	inventorySw.ParamsReq
}

type DelKeywordsSwReq struct {
	request.EnterpriseBaseReq
	Ids      []int64 `json:"ids"`
	MemberId int64   `json:"member_id" binding:"required"`
	DelType  string  `json:"del_type"`
}

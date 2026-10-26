package microScrm

// ActivityGoodsListRequest GoodsList 处理器请求
type ActivityGoodsListRequest struct {
	EnterpriseID   int64  `json:"enterprise_id"` // EnterpriseID 企业ID
	MemberUserId   string `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	ActivityType   int    `json:"activity_type"`   // ActivityType 活动类型
	ActivityStatus int    `json:"activity_status"` // ActivityStatus 活动状态 0:草稿 1：待上架 2：已上架 3：已下架 4：已结束 5：暂停 6：定时上架 7：定时下架
	GoodsName      string `json:"goods_name"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

// ActivityGoodsListItem GoodsListResponse 单条item
type ActivityGoodsListItem struct {
	SkuName       string            `json:"sku_name" db:"sku_name"`
	ActivityID    int64             `json:"shop_activity_id" db:"shop_activity_id"` // ActivityID 商城活动id
	ActivityType  int               `json:"activity_type"`
	ActivityPrice int               `json:"activity_price"`             // 活动价钱（单位：分）
	SellPrice     int               `json:"sell_price" db:"sell_price"` //  商品售卖价（单位：分）
	SkuAttribute  string            `json:"sku_attribute" db:"sku_attribute"`
	Attribute     []AttributeEntity `json:"attributes"` // 商品属性
	NeedNum       int               `json:"need_num"`   // NeedNum 需要数量
	SkuID         int64             `json:"sku_id" db:"sku_id"`
	SpuID         int64             `json:"spu_id" db:"spu_id"`
	SkuPicture    string            `json:"image" db:"sku_picture"`
	ShareAppId    string            `json:"share_app_id"`
	SharePath     string            `json:"share_path"`
}

// ActivityGoodsListResponse GoodsList 处理器响应
type ActivityGoodsListResponse struct {
	List  []*ActivityGoodsListItem `json:"list"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
	Total int64                    `json:"total"`
}

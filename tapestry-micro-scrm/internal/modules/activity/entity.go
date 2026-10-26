package activity

import (
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/pkg/common"
)

// 活动商品状态
const (
	SKuActivityUpStatus = iota + 1
	SKuActivityDownStatus
)

// 活动状态
const (
	DraftStatus = iota
	PendingStatus
	NormalStatus
	DownStatus
	FinishStatus
	StopStatus
	TimingUpStatus
	TimingDownStatus
)

// ShopActivityEntity
type ShopActivityEntity struct {
	ID             int64  `db:"id"`
	ActivityID     int64  `json:"id" db:"activity_id"`                  // ActivityID 商城活动id
	Source         string `json:"source" db:"source"`                   // Source 活动数据来源：crs
	Name           string `json:"name" db:"name"`                       // Name 活动名称
	StartAt        int64  `json:"start_at" db:"start_at"`               // StartAt 活动开始时间
	EndAt          int64  `json:"end_at" db:"end_at"`                   // EndAt 活动结束时间
	NeedNum        int    `json:"need_num" db:"need_num"`               // NeedNum 需要数量
	ActivityType   int    `json:"activity_type" db:"activity_type"`     // ActivityType 活动类型
	ActivityStatus int    `json:"activity_status" db:"activity_status"` // ActivityStatus 活动状态 0:草稿 1：待上架 2：已上架 3：已下架 4：已结束 5：暂停 6：定时上架 7：定时下架
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`     // EnterpriseID 企业id
	CreatedAt      int64  `json:"created_at" db:"created_at"`           // CreatedAt 创建时间
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`           // UpdatedAt 修改时间
	DeletedAt      int64  `json:"deleted_at" db:"deleted_at"`           // DeletedAt 删除时间
}

// ShopActivitySku ...
type ShopActivitySku struct {
	ID                int64  `json:"id" db:"id"`
	ShopActivityID    int64  `json:"shop_activity_id" db:"shop_activity_id"`       // ShopActivityID 商城活动主表ID
	EnterpriseID      int64  `json:"enterprise_id" db:"enterprise_id"`             // EnterpriseID 企业ID
	SpuID             int64  `json:"spu_id" db:"spu_id"`                           // SpuID 商品spu_id
	SkuID             int64  `json:"sku_id" db:"sku_id"`                           // SkuID 商品sku_id
	SkuStatus         int    `json:"sku_status" db:"sku_status"`                   // SkuStatus 状态 1：启用 2：禁用
	ActivityStockType int    `json:"activity_stock_type" db:"activity_stock_type"` // ActivityStockType 库存类型 1：使用活动库存 2：使用商品库存
	ActivityStockSet  int    `json:"activity_stock_set" db:"activity_stock_set"`   // ActivityStockSet 活动设置库存
	SkuPicture        string `json:"sku_picture" db:"sku_picture"`                 // SkuPicture sku图像
	SellPrice         int    `json:"sell_price" db:"sell_price"`                   // SellPrice 生成活动时候的商品销售价格
	ActivityPrice     int    `json:"activity_price" db:"activity_price"`           // ActivityPrice 活动价格
	CreatedAt         int64  `json:"created_at" db:"created_at"`                   // CreatedAt 创建时间
	UpdatedAt         int64  `json:"updated_at" db:"updated_at"`                   // UpdatedAt 更新时间
	DeletedAt         int64  `json:"deleted_at" db:"deleted_at"`                   // DeletedAt 删除时间
}

// SyncRequest 活动数据同步请求参数结构体
type SyncRequest struct {
	Method   common.DdcEvent `json:"method"`
	Platform string          `json:"platform"`
	Data     struct {
		ShopActivityEntity
		SkuList []*ShopActivitySku `json:"sku_list"`
	} `json:"data"`
}

// GoodsListResponse GoodsList 处理器请求
type GoodsListRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"` // EnterpriseID 企业ID
	MemberUserId   string `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	ActivityType   int    `json:"activity_type" db:"activity_type"`     // ActivityType 活动类型
	ActivityStatus int    `json:"activity_status" db:"activity_status"` // ActivityStatus 活动状态 0:草稿 1：待上架 2：已上架 3：已下架 4：已结束 5：暂停 6：定时上架 7：定时下架
	GoodsName      string `json:"goods_name"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

// GoodsListItem GoodsListResponse 单条item
type GoodsListItem struct {
	SkuName       string                  `json:"sku_name" db:"sku_name"`
	ActivityID    int64                   `json:"shop_activity_id" db:"shop_activity_id"` // ActivityID 商城活动id
	ActivityType  int                     `json:"activity_type"`
	ActivityPrice int                     `json:"activity_price" db:"activity_price"` // 活动价钱（单位：分）
	SellPrice     int                     `json:"sell_price" db:"sell_price"`         //  商品售卖价（单位：分）
	SkuAttribute  string                  `json:"sku_attribute" db:"sku_attribute"`
	Attribute     []goods.AttributeEntity `json:"attributes"` // 商品属性
	NeedNum       int                     `json:"need_num"`   // NeedNum 需要数量
	SkuID         int64                   `json:"sku_id" db:"sku_id"`
	SpuID         int64                   `json:"spu_id" db:"spu_id"`
	SkuPicture    string                  `json:"image" db:"sku_picture"`
	ShareAppId    string                  `json:"share_app_id"`
	SharePath     string                  `json:"share_path"`
}

// GoodsListResponse GoodsList 处理器响应
type GoodsListResponse struct {
	List  []*GoodsListItem `json:"list"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int64            `json:"total"`
}

type SkuListByActivityIdsReq struct {
	ActivityIds  []int64
	EnterpriseID int64
	GoodsName    string
	Status       int //`sku_status` 状态 1：启用 2：禁用',
	Page         int
	Size         int
}

package goods

import (
	"encoding/json"
	"micro-scrm/internal/pkg/common"
)

// SKU 状态
const (
	SkuPendingStatus = 1
	SkuNormalStatus  = 2
	SkuDownStatus    = 3
)

// SPU 状态
const (
	SPuPendingStatus = 0
	SPuNormalStatus  = 1
	SPuDownStatus    = 2
)

// SkuEntity ...
type SkuEntity struct {
	ID           int64   `db:"id"`
	SkuID        int64   `json:"id" db:"sku_id"`                   // SkuID SKU ID
	SpuID        int64   `json:"spu_id" db:"spu_id"`               // SpuID SPU ID
	SourceID     int64   `json:"source_id" db:"source_id"`         // SourceID 门店Id
	MaterialID   int64   `json:"material_id" db:"material_id"`     // MaterialID 物料ID
	SkuName      string  `json:"sku_name" db:"sku_name"`           // SkuName sku名称
	CardTitle    string  `json:"card_title" db:"card_title"`       // CardTitle 分享卡片标题
	PosterTitle  string  `json:"poster_title" db:"poster_title"`   // PosterTitle 海报标题
	SellPoint    string  `json:"sell_point" db:"sell_point"`       // SellPoint 卖点
	ShortTitle   string  `json:"short_title" db:"short_title"`     // ShortTitle 商品短标题
	Sales        int     `json:"sales" db:"sales"`                 // Sales 销量
	Stock        int     `json:"stock" db:"stock"`                 // Stock 商品库存
	GuidePrice   int     `json:"guide_price" db:"guide_price"`     // GuidePrice 总部指导价
	SignPrice    int     `json:"sign_price" db:"sign_price"`       // SignPrice 商品标牌价（单位：分）
	SellPrice    int     `json:"sell_price" db:"sell_price"`       // SellPrice 商品售卖价（单位：分）
	SaleType     int     `json:"sale_type" db:"sale_type"`         // SaleType 商品售卖类型（1:快递 2:自提 3:快递+自提）
	QrcodeURL    string  `json:"qrcode_url" db:"qrcode_url"`       // QrcodeURL 二维码地址
	Status       int     `json:"status" db:"status"`               // Status 商品状态 1待上架 2上架中 3已下架
	GoodsWeight  float32 `json:"goods_weight" db:"goods_weight"`   // GoodsWeight 商品重量
	GoodsVolume  float32 `json:"goods_volume" db:"goods_volume"`   // GoodsVolume 商品体积
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id"` // EnterpriseID 企业ID
	SkuAttribute string  `json:"sku_attribute" db:"sku_attribute"` // SkuAttribute SKU商品属性值集合
	SkuPicture   string  `json:"sku_picture" db:"sku_picture"`     // SkuPicture sku图像
	CreateTime   int     `json:"created_at" db:"create_time"`      // CreateTime sku创建事件
	CreatedAt    int64   `db:"created_at"`                         // CreatedAt 创建时间
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`       // UpdatedAt 更新时间
	DeletedAt    int64   `json:"deleted_at" db:"deleted_at"`       // DeletedAt 删除时间
}

// SpuEntity ...
type SpuEntity struct {
	ID               int64  `json:"id" db:"id"`
	SpuID            int64  `json:"spu_id" db:"spu_id"`                       // SpuID spu_id
	SourceID         int64  `json:"source_id" db:"source_id"`                 // SourceID 门店Id
	SpuName          string `json:"spu_name" db:"spu_name"`                   // SpuName 商品名称
	SpuSn            string `json:"spu_sn" db:"spu_sn"`                       // SpuSn 商品编号
	ShortTitle       string `json:"short_title" db:"short_title"`             // ShortTitle 商品短标题
	SellPoint        string `json:"sell_point" db:"sell_point"`               // SellPoint 卖点
	Sales            int    `json:"sales" db:"sales"`                         // Sales 销量
	Status           int    `json:"is_up" db:"status"`                        // Status 商品状态 1待上架 2上架中 3已下架
	CategoryID       int64  `json:"category_id" db:"category_id"`             // CategoryID 分类ID
	CategoryName     string `json:"category_name" db:"category_name"`         // CategoryName 分类名称
	CreateTime       int    `json:"created_at" db:"create_time"`              // CreateTime 商品创建时间
	EnterpriseID     int64  `json:"enterprise_id" db:"enterprise_id"`         // EnterpriseID 企业ID
	Picture          string `json:"picture" db:"picture"`                     // Picture spu图像
	DefaultAttribute string `json:"default_attribute" db:"default_attribute"` // DefaultAttribute SPU商品默认属性值集合
	CreatedAt        int64  `db:"created_at"`                                 // CreatedAt 创建时间
	UpdatedAt        int64  `json:"updated_at" db:"updated_at"`               // UpdatedAt 更新时间
	DeletedAt        int64  `json:"deleted_at" db:"deleted_at"`               // DeletedAt 删除时间
}

// SyncRequest 同步数据请求参数结构
type SyncRequest struct {
	Method   common.DdcEvent `json:"method"`
	Platform string          `json:"platform"`
	Data     struct {
		SpuEntity
		Attributes        []*AttributeItem `json:"attributes"`
		DefaultAttributes struct {
			Attributes []*SkuAttribute `json:"attributes"`
			Pictures   struct {
				Image string `json:"image"`
			} `json:"pictures"`
			EnterpriseId int64 `json:"enterprise_id"`
		} `json:"default_attributes"`
		Sku struct {
			List []*SyncSkuItem `json:"list"`
		} `json:"sku"`
	} `json:"data"`
}

// SyncChangeRequest DDC同步商品变更事件请求结构体
type SyncChangeRequest struct {
	Method   common.DdcEvent `json:"method"`
	Platform string          `json:"platform"`
	Data     struct {
		List []*SyncChangeItem `json:"list"`
	} `json:"data"`
}

// SyncDeleteRequest 删除商品同步事件请求
type SyncDeleteRequest struct {
	Method   common.DdcEvent `json:"method"`
	Platform string          `json:"platform"`
	Data     struct {
		EnterpriseID int64   `json:"enterprise_id"`
		SpuID        int64   `json:"spu_id"`
		SkuIDs       []int64 `json:"sku_ids"`
		IsDeleteSpu  bool    `json:"is_delete_spu"`
	} `json:"data"`
}

// SyncChangeItem DDC同步商品变更事件请求结构体
type SyncChangeItem struct {
	SpuID     int64            `json:"spu_id"`
	SpuStatus int              `json:"spu_status"`
	SkuList   []*SyncChangeSku `json:"sku"`
}

// SyncChangeSku DDC同步商品变更事件请求结构体-sku信息
type SyncChangeSku struct {
	SpuID     int64 `json:"spu_id"`
	SourceID  int64 `json:"source_id"`
	ID        int64 `json:"id"`
	Status    int   `json:"status"`
	Stock     int   `json:"stock"`
	Price     int   `json:"price"`
	SignPrice int   `json:"sign_price"`
}

// SyncSkuItem 同步请求sku结构
type SyncSkuItem struct {
	SkuEntity
	Attributes []*SkuAttribute `json:"attributes"`
	Pictures   struct {
		Image string `json:"image"`
	} `json:"pictures"`
}

// AttributeItem 商品属性定义
type AttributeItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SkuAttribute sku属性值
type SkuAttribute struct {
	AttributeID    int64 `json:"attribute_id"`
	AttributeName  string
	AttributeValue string `json:"attribute_value"`
}

// StoreAttribute 数据库中存储的规格数据格式
type StoreAttribute struct {
	AttributeName      string `json:"attribute_name"`
	AttributeItemValue string `json:"attribute_item_value"`
}

// ListReq 商品请求实体
type ListReq struct {
	SkuID        int64  `json:"sku_id" db:"sku_id" desc:"sku_id"`
	SkuName      string `json:"sku_name" db:"sku_name" desc:"商品名称"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberUserId string `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	Status       uint
}

// SkuItem SPU单条数据
type SkuItem struct {
	ID           int64             `json:"id" db:"id"`
	SkuID        int64             `json:"sku_id" db:"sku_id"`               // SkuID SKU ID
	SpuID        int64             `json:"spu_id" db:"spu_id"`               // SpuID SPU ID
	SourceID     int64             `json:"source_id" db:"source_id"`         // SourceID 门店Id
	MaterialID   int64             `json:"material_id" db:"material_id"`     // MaterialID 物料ID
	SkuName      string            `json:"sku_name" db:"sku_name"`           // SkuName sku名称
	CardTitle    string            `json:"card_title" db:"card_title"`       // CardTitle 分享卡片标题
	PosterTitle  string            `json:"poster_title" db:"poster_title"`   // PosterTitle 海报标题
	SellPoint    string            `json:"sell_point" db:"sell_point"`       // SellPoint 卖点
	ShortTitle   string            `json:"short_title" db:"short_title"`     // ShortTitle 商品短标题
	Sales        int               `json:"sales" db:"sales"`                 // Sales 销量
	Stock        int               `json:"stock" db:"stock"`                 // Stock 商品库存
	GuidePrice   int               `json:"guide_price" db:"guide_price"`     // GuidePrice 总部指导价
	SignPrice    int               `json:"sign_price" db:"sign_price"`       // SignPrice 商品标牌价（单位：分）
	SellPrice    int               `json:"sell_price" db:"sell_price"`       // SellPrice 商品售卖价（单位：分）
	SaleType     int               `json:"sale_type" db:"sale_type"`         // SaleType 商品售卖类型（1:快递 2:自提 3:快递+自提）
	QrcodeURL    string            `json:"qrcode_url" db:"qrcode_url"`       // QrcodeURL 二维码地址
	Status       int               `json:"status" db:"status"`               // Status 商品状态 1待上架 2上架中 3已下架
	GoodsWeight  float32           `json:"goods_weight" db:"goods_weight"`   // GoodsWeight 商品重量
	GoodsVolume  float32           `json:"goods_volume" db:"goods_volume"`   // GoodsVolume 商品体积
	EnterpriseID int64             `json:"enterprise_id" db:"enterprise_id"` // EnterpriseID 企业ID
	SkuAttribute string            `json:"sku_attribute" db:"sku_attribute"` // SkuAttribute SKU商品属性值集合
	SkuPicture   string            `json:"sku_picture" db:"sku_picture"`     // SkuPicture sku图像
	CreateTime   int               `json:"create_time" db:"create_time"`     // CreateTime sku创建事件
	CreatedAt    int64             `json:"created_at" db:"created_at"`       // CreatedAt 创建时间
	UpdatedAt    int64             `json:"updated_at" db:"updated_at"`       // UpdatedAt 更新时间
	DeletedAt    int64             `json:"deleted_at" db:"deleted_at"`       // DeletedAt 删除时间
	Attribute    []AttributeEntity `json:"attributes"`
	ShareAppId   string            `json:"share_app_id"`
	SharePath    string            `json:"share_path"`
}

// Attribute 解析属性
func (s *SkuItem) SetAttributes() (err error) {
	var attr []AttributeEntity
	if s.SkuAttribute == "" {
		return
	}

	err = json.Unmarshal([]byte(s.SkuAttribute), &attr)
	s.Attribute = attr
	return
}

type AttributeEntity struct {
	Name  string `json:"attribute_name"`
	Value string `json:"attribute_item_value"`
}

// SpuListRes 响应结构体
type SpuListRes struct {
	List  []*SkuItem `json:"list"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int64      `json:"total"`
}

// SPU 列表请求结构体
type SpuListRequest struct {
	SpuEntity
	SpuIDs []int64
	Page   int `json:"page"`
	Size   int `json:"size"`
}

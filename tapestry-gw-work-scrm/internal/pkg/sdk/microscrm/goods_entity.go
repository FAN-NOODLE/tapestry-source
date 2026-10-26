package microScrm

// GoodsListReq 商品请求实体
type GoodsListReq struct {
	SkuID        int64  `json:"sku_id" db:"sku_id" desc:"sku_id"`
	SkuName      string `json:"sku_name" db:"sku_name" desc:"商品名称"`
	MemberUserId string `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

// GoodsSpuListReq 商品请求实体
type GoodsSpuListReq struct {
	SpuIds       []int64 `json:"spu_ids"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
}

// SpuListRes 响应结构体
type SpuListRes struct {
	List  []*SkuItem `json:"list"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int64      `json:"total"`
}

// SkuItem SPU单条数据
type SkuItem struct {
	SkuEntity
	Attribute []AttributeEntity `json:"attributes"`
}

// SkuEntity ...
type SkuEntity struct {
	ID           int64   `json:"id" db:"id"`
	SkuID        int64   `json:"sku_id" db:"sku_id"`               // SkuID SKU ID
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
	CreateTime   int     `json:"create_time" db:"create_time"`     // CreateTime sku创建事件
	CreatedAt    int64   `json:"created_at" db:"created_at"`       // CreatedAt 创建时间
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`       // UpdatedAt 更新时间
	DeletedAt    int64   `json:"deleted_at" db:"deleted_at"`       // DeletedAt 删除时间
	ShareAppId   string  `json:"share_app_id"`
	SharePath    string  `json:"share_path"`
}

// AttributeEntity 商品信息
type AttributeEntity struct {
	Name  string `json:"attribute_name"`
	Value string `json:"attribute_item_value"`
}

package marketingPromotion

const (
	/**
	 * 1 商品 2 优惠券 6 好友助力-付邮 7 好友助力-买赠  12 限时购
	 * 13 助力抽奖 14 拼团 15 砍价
	 */
	MarketingGoods         = 1
	MarketingCoupon        = 2
	MarketingHelperPost    = 6
	MarketingHelperBuyGive = 7
	MarketingSecKill       = 12
	MarketingLottery       = 13
	MarketingGroup         = 14
	MarketingBargain       = 15

	OperatorTypeSend     = 1
	OperatorTypeDownload = 2
	ImgFileType          = "1"
)

type Test struct {
	FileUrl string `json:"file_url" desc:"file_url"`
	Type    int    `json:"type" desc:"type"`
}

// GetPosterRequest 生成海报
type GetPosterRequest struct {
	Type           int    `json:"type" desc:"二维码类型"`
	OperatorType   int    `json:"operator_type" desc:"操作类型 1 发送 2 下载"`
	MemUserID      string `json:"member_userid" desc:"导购userid"`
	MemberName     string `json:"member_name" desc:"员工姓名"`
	MemberAvatar   string `json:"member_avatar" desc:"员工头像"`
	EnterpriseID   int64  `json:"enterprise_id" desc:"企业ID"`
	EnterpriseHash string `json:"enterprise_hash" desc:"企业ID"`
	SunCode        string `json:"sun_code" desc:"太阳码"`
	Params         string `json:"params" desc:"params海报参数"`
	FileTmpUrl     string `json:"file_tmp_url" desc:"海报临时地址"`
	ShareImg       string `json:"share_img" desc:"分享图片地址"`
	IsTest         int    `json:"is_test" desc:"测试"`
	FileName       string `json:"file_name" desc:"文件名称"`
}

// GetPosterResponse 海报返回
type GetPosterResponse struct {
	Object string `json:"object" desc:"object"`
}

// WxQrcodeRequest 生成小程序太阳码请求参数
type WxQrcodeRequest struct {
	EnterpriseID int64     `json:"enterprise_id"`
	Path         string    `json:"path"`
	Params       CodeParam `json:"params"`
}

// GoodsParams 商品海报参数
type GoodsParams struct {
	SkuName    string `json:"sku_name" desc:"sku_name"`
	SkuID      int64  `json:"sku_id" desc:"sku_id"`
	SpuID      int64  `json:"spu_id" desc:"spu_id"`
	SkuPicture string `json:"sku_picture" desc:"sku_picture"`
	SellPrice  int    `json:"sell_price" desc:"sell_price"`
	SignPrice  int    `json:"sign_price" desc:"sign_price"`
}

// CouponParams 优惠券海报参数
type CouponParams struct {
	CouponID   int64  `json:"coupon_id" desc:"coupon_id"`
	Value      int    `json:"coupon_value" desc:"coupon_value"`
	CouponMini int    `json:"coupon_mininumn" desc:"coupon_mininumn"`
	StartTime  int64  `json:"start_time" desc:"start_time"`
	EndTime    int64  `json:"end_time" desc:"end_time"`
	LogoUrl    string `json:"logo_url" desc:"logo_url"`
	IsAll      int    `json:"is_all" desc:"is_all"`
	Condition  string `json:"condition" desc:"condition"`
}

// ActivityParams 活动海报参数
type ActivityParams struct {
	GoodsParams
	Type           int    `json:"activity_type" desc:"activity_type"`
	Image          string `json:"image" desc:"image"`
	ActivityPrice  int    `json:"activity_price" desc:"activity_price"`
	ShopActivityID int64  `json:"shop_activity_id" desc:"shop_activity_id"`
	AppName        string `json:"app_name" desc:"app_name"`
	AppLogo        string `json:"app_logo" desc:"app_logo"`
}

// 活动icon结构体
type ActivityIcon struct {
	X          int    `json:"x" desc:"x"`
	Y          int    `json:"y" desc:"y"`
	TopIcon    string `json:"top_icon" desc:"top_icon"`
	BottomIcon string `json:"bottom_icon" desc:"bottom_icon"`
}

type CodeParam interface {
	SetShareImg(string)
	GetShareImg() string
}

type GoodsCodeParam struct {
	SpuID    int64  `json:"spu_id"`
	SkuID    int64  `json:"sku_id"`
	MemUID   string `json:"mem_uid"`
	StoreID  int64  `json:"store_id"`
	ShareUid string `json:"share_uid"`
	Source   string `json:"source"`
}

func (g *GoodsCodeParam) SetShareImg(_ string) {
}

func (g *GoodsCodeParam) GetShareImg() string {
	return ""
}

type CouponCodeParam struct {
	CouponID int64  `json:"id"`
	StoreID  int64  `json:"store_id"`
	MemUID   string `json:"mem_uid"`
	ShareUid string `json:"share_uid"`
	Source   string `json:"source"`
	ShareImg string `json:"share_img"`
}

func (c *CouponCodeParam) SetShareImg(s string) {
	c.ShareImg = s
}

func (c *CouponCodeParam) GetShareImg() string {
	return c.ShareImg
}

type GroupCodeParam struct {
	SpuID    int64  `json:"spu_id"`
	SkuID    int64  `json:"sku_id"`
	MemUID   string `json:"mem_uid"`
	StoreID  int64  `json:"store_id"`
	Type     int    `json:"type"`
	RuleID   int64  `json:"rule_id"`
	ShareUid string `json:"share_uid"`
	Source   string `json:"source"`
}

func (g *GroupCodeParam) SetShareImg(_ string) {
}

func (g *GroupCodeParam) GetShareImg() string {
	return ""
}

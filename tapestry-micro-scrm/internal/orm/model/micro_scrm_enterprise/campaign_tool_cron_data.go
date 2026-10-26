package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type CampaignToolCronDataModel struct {
	model.BranchModel

	MemberCode          string `json:"member_code" gorm:"column:member_code;type:varchar(64);not null;comment:会员卡号"`                                   // 会员卡号
	UnionID             string `json:"union_id" gorm:"column:union_id;type:varchar(500);not null;comment:微信union_id"`                                  // 微信union_id
	IsContact           string `json:"is_contact" gorm:"column:is_contact;type:varchar(64);not null;comment:可沟通"`                                      // 可沟通
	Gc                  string `json:"gc" gorm:"column:gc;type:varchar(64);not null;comment:当前会员等级是GC"`                                                // 当前会员等级是GC
	My                  string `json:"my" gorm:"column:my;type:varchar(64);not null;comment:当前会员等级是名媛"`                                                // 当前会员等级是名媛
	JinKa               string `json:"jin_ka" gorm:"column:jin_ka;type:varchar(64);not null;comment:当前会员等级是金卡"`                                        // 当前会员等级是金卡
	BaiJinKa            string `json:"bai_jin_ka" gorm:"column:bai_jin_ka;type:varchar(64);not null;comment:当前会员等级是白金卡"`                               // 当前会员等级是白金卡
	ZuanKa              string `json:"zuan_ka" gorm:"column:zuan_ka;type:varchar(64);not null;comment:当前会员等级是钻卡"`                                      // 当前会员等级是钻卡
	HeiZuanKa           string `json:"hei_zuan_ka" gorm:"column:hei_zuan_ka;type:varchar(64);not null;comment:当前会员等级是黑钻卡"`                             // 当前会员等级是黑钻卡
	GuiTaiCode          string `json:"gui_tai_code" gorm:"column:gui_tai_code;type:varchar(64);not null;comment:管理柜台编号"`                               // 管理柜台编号
	PinPaiBuyLast       string `json:"pin_pai_buy_last" gorm:"column:pin_pai_buy_last;type:varchar(200);not null;comment:品牌最后一笔订单购买产品"`                // 品牌最后一笔订单购买产品
	OnlineBuyLast       string `json:"online_buy_last" gorm:"column:online_buy_last;type:varchar(200);not null;comment:线上渠道最后一笔订单购买产品"`                // 线上渠道最后一笔订单购买产品
	OfflineBuyLast      string `json:"offline_buy_last" gorm:"column:offline_buy_last;type:varchar(200);not null;comment:线下渠道最后一笔订单购买产品"`              // 线下渠道最后一笔订单购买产品
	TianMaoBuyLast      string `json:"tian_mao_buy_last" gorm:"column:tian_mao_buy_last;type:varchar(200);not null;comment:天猫渠道最后一笔订单购买产品"`            // 天猫渠道最后一笔订单购买产品
	JingDongBuyLast     string `json:"jing_dong_buy_last" gorm:"column:jing_dong_buy_last;type:varchar(200);not null;comment:京东渠道最后一笔订单购买产品"`          // 京东渠道最后一笔订单购买产品
	MiniBuyLast         string `json:"mini_buy_last" gorm:"column:mini_buy_last;type:varchar(200);not null;comment:mini program渠道最后一笔订单购买产品"`          // mini program渠道最后一笔订单购买产品
	PinPaiBuyLastTime   string `json:"pin_pai_buy_last_time" gorm:"column:pin_pai_buy_last_time;type:varchar(64);not null;comment:品牌最后一笔订单距今时长"`       // 品牌最后一笔订单距今时长
	OnlineBuyLastTime   string `json:"online_buy_last_time" gorm:"column:online_buy_last_time;type:varchar(64);not null;comment:线上渠道最后一笔订单距今时长"`       // 线上渠道最后一笔订单距今时长
	OfflineBuyTime      string `json:"offline_buy_time" gorm:"column:offline_buy_time;type:varchar(64);not null;comment:线下渠道最后一笔订单距今时长"`               // 线下渠道最后一笔订单距今时长
	TianMaoBuyLastTime  string `json:"tian_mao_buy_last_time" gorm:"column:tian_mao_buy_last_time;type:varchar(64);not null;comment:天猫渠道最后一笔订单距今时长"`   // 天猫渠道最后一笔订单距今时长
	JingDongBuyLastTime string `json:"jing_dong_buy_last_time" gorm:"column:jing_dong_buy_last_time;type:varchar(64);not null;comment:京东渠道最后一笔订单距今时长"` // 京东渠道最后一笔订单距今时长
	MiniBuyLastTime     string `json:"mini_buy_last_time" gorm:"column:mini_buy_last_time;type:varchar(64);not null;comment:mini program渠道最后一笔订单距今时长"` // mini program渠道最后一笔订单距今时长
	R12mBuyProductList  string `json:"r12m_buy_product_list" gorm:"column:r12m_buy_product_list;type:text;not null;comment:R12M品牌购买的产品"`               // R12M品牌购买的产品
	YtdOnlineAus        string `json:"ytd_online_aus" gorm:"column:ytd_online_aus;type:text;not null;comment:YTD线上渠道AUS"`                              // YTD线上渠道AUS
	YtdOfflineAus       string `json:"ytd_offline_aus" gorm:"column:ytd_offline_aus;type:text;not null;comment:YTD线下渠道AUS"`                            // YTD线下渠道AUS
	MainBuyProduct      string `json:"main_buy_product" gorm:"column:main_buy_product;type:varchar(200);not null;comment:主要购买产品"`                      // 主要购买产品
	CronFileName        string `json:"cron_file_name" gorm:"column:cron_file_name;type:varchar(64);not null;comment:sftp文件名"`                          // sftp文件名
	CreateTime          string `json:"create_time" gorm:"column:create_time;type:varchar(14);not null;comment:campaign-tool推送时间"`                      // campaign-tool推送时间
	IsPull              int32  `json:"is_pull" gorm:"column:is_pull;type:tinyint(2);not null;default:0;comment:是否拉取1拉取，0没有拉取"`                         // 是否拉取1拉取，0没有拉取
	TagName             string `json:"tag_name" gorm:"column:tag_name;type:varchar(50);not null;comment:标签名"`                                          // 标签名
}

func (CampaignToolCronDataModel) TableName() string {
	return "campaign_tool_cron_data"
}

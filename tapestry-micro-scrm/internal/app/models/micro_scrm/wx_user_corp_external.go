package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	VipTypeChoice = 3 // VIP客筛选 3

	VipType = 1 // VIP客 1

	CustomerTypeNew   = 1 // 新客
	CustomerTypeOld   = 2 // 老客
	CustomerTypeSleep = 4 // 沉睡客

	AgeTagOne = 1 // Z时代
	AgeTagTwo = 2 // 千禧一代
)

type WxUserCorpExternalModel struct {
	models.BaseModel
	Appid                 string `json:"appid" gorm:"column:appid;type:char(64);not null;uniqueIndex:uin;comment:对应的是corp_id"`
	Openid                string `json:"openid" gorm:"column:openid;type:char(64);not null;uniqueIndex:uin;comment:对应的是external_user_id"`
	Unionid               string `json:"unionid" gorm:"column:unionid;type:char(64);not null;index:unionid_idx;comment:用户unionid"`
	Euid                  int64  `json:"euid" gorm:"column:euid;type:bigint(20);not null;index:euid;default:0;comment:聚合id"`
	Nick                  string `json:"nick" gorm:"column:nick;type:varchar(255);not null;comment:昵称"`
	Name                  string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:姓名"`
	Position              string `json:"position" gorm:"column:position;type:varchar(255);not null;comment:用户职称"`
	Avatar                string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;comment:头像"`
	CorpName              string `json:"corp_name" gorm:"column:corp_name;type:varchar(255);not null;comment:企业简称"`
	CorpFullName          string `json:"corp_full_name" gorm:"column:corp_full_name;type:varchar(255);not null;comment:企业全称"`
	Type                  int32  `json:"type" gorm:"column:type;type:smallint(4);not null;default:0;comment:外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户"`
	Sex                   int32  `json:"sex" gorm:"column:sex;type:smallint(4);not null;default:0;comment:外部联系人性别 0-未知 1-男性 2-女性"`
	ExternalProfile       string `json:"external_profile" gorm:"column:external_profile;type:text;not null;comment:外部联系人的自定义展示信息,JSON"`
	Birthday              string `json:"birthday" gorm:"column:birthday;type:varchar(20);comment:生日"`
	UnionCode             string `json:"union_code" gorm:"column:union_code;type:varchar(20);comment:会员卡号"`
	Point                 string `json:"point" gorm:"column:point;type:varchar(10);comment:会员积分"`
	Grade                 string `json:"grade" gorm:"column:grade;type:varchar(10);comment:会员等级"`
	StoreCode             string `json:"store_code" gorm:"column:store_code;type:varchar(50);not null;comment:门店编号"`
	BaCode                string `json:"ba_code" gorm:"column:ba_code;type:varchar(50);not null;comment:BA编号"`
	BirthYear             string `json:"birth_year" gorm:"column:birth_year;type:varchar(50);not null;comment:生日的年份"`
	BirthMonth            string `json:"birth_month" gorm:"column:birth_month;type:varchar(50);not null;comment:生日的月份"`
	BirthDay              string `json:"birth_day" gorm:"column:birth_day;type:varchar(50);not null;comment:生日的天"`
	Mobile                string `json:"mobile" gorm:"column:mobile;type:varchar(15);not null;comment:企业微信手机号"`
	CardLevel             string `json:"card_level" gorm:"column:card_level;type:varchar(100);not null;comment:会员等级"`
	CltType               string `json:"clt_type" gorm:"column:clt_type;type:varchar(100);not null;comment:CLT类型"`
	AmountTotal           int64  `json:"amount_total" gorm:"column:amount_total;type:bigint(20);not null;default:0;comment:总消费金额：分"`
	CustomerCode          string `json:"customer_code" gorm:"column:customer_code;type:varchar(255);not null;comment:cdp的用户customer_code"`
	CustomerID            string `json:"customer_id" gorm:"column:customer_id;type:varchar(10);not null;comment:cdp的customer表id"`
	LastPurchaseDate      int    `json:"last_purchase_date" gorm:"column:last_purchase_date;type:int(11);not null;default:0;comment:最近一次订单消费时间"`
	LastPurchaseStoreCode string `json:"last_purchase_store_code" gorm:"column:last_purchase_store_code;type:varchar(20);not null;comment:最近一次订单购买门店code"`
	UpgradeNeedAmount     int    `json:"upgrade_need_amount" gorm:"column:upgrade_need_amount;type:int(11);not null;default:0;comment:升级下一等级所需金额：分"`
	CustomerType          int    `json:"customer_type" gorm:"column:customer_type;type:tinyint(5);not null;default:0;comment:顾客类型：1新客，2老客，4沉睡客"`
	CustomerTypeUpdatedAt int    `json:"customer_type_updated_at" gorm:"column:customer_type_updated_at;type:int(11);not null;default:0;comment:顾客类型更新时间"`
	VipType               int    `json:"vip_type" gorm:"column:vip_type;type:tinyint(5);not null;default:0;comment:vip类型:1vip"`
	AgeTag                string `json:"age_tag" gorm:"column:age_tag;type:varchar(255);not null;comment:coach顾客类型的tag"`
}

// TableName returns the table name of the struct_name model
func (WxUserCorpExternalModel) TableName() string {
	return "wx_user_corp_external"
}

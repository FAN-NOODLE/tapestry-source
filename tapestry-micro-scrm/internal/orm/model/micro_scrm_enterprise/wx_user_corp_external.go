package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

const (
	WECHATt_USER_TYPE = iota + 1
	ENTERPRISE_WECHATt_USER_TYPE
)

type WxUserCorpExternalModel struct {
	model.BranchModel

	Appid           string `gorm:"column:appid;type:char(64);not null;uniqueIndex:uin;comment:对应的是corp_id"`
	Openid          string `gorm:"column:openid;type:char(64);not null;uniqueIndex:uin;comment:对应的是external_user_id"`
	Unionid         string `gorm:"column:unionid;type:char(64);not null;comment:用户unionid"`
	Euid            int64  `gorm:"column:euid;type:bigint(20);not null;index:euid;default:0;comment:聚合id"`
	Nick            string `gorm:"column:nick;type:varchar(255);not null;comment:昵称"`
	Name            string `gorm:"column:name;type:varchar(255);not null;comment:姓名"`
	Position        string `gorm:"column:position;type:varchar(255);not null;comment:用户职称"`
	Avatar          string `gorm:"column:avatar;type:varchar(255);not null;comment:头像"`
	Birthday        string `gorm:"column:birthday;type:varchar(50);not null;comment:生日"`
	UnionCode       string `gorm:"column:union_code;type:varchar(50);not null;comment:会员卡号"`
	Point           string `gorm:"column:point;type:varchar(10);not null;comment:会员积分"`
	Grade           string `gorm:"column:grade;type:varchar(20);not null;comment:会员等级"`
	StoreCode       string `gorm:"column:store_code;type:varchar(50);not null;comment:门店编号"`
	BaCode          string `gorm:"column:ba_code;type:varchar(50);not null;comment:BA编号"`
	CorpName        string `gorm:"column:corp_name;type:varchar(255);not null;comment:企业简称"`
	CorpFullName    string `gorm:"column:corp_full_name;type:varchar(255);not null;comment:企业全称"`
	Type            int32  `gorm:"column:type;type:smallint(4);not null;default:0;comment:外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户"`
	Sex             int32  `gorm:"column:sex;type:smallint(4);not null;default:0;comment:外部联系人性别 0-未知 1-男性 2-女性"`
	ExternalProfile string `gorm:"column:external_profile;type:text;not null;comment:外部联系人的自定义展示信息,JSON"`
	BirthYear       string `gorm:"column:birth_year;type:int(10);not null;comment:生日的年份" json:"birth_year"`
	BirthMonth      string `gorm:"column:birth_month;type:int(10);not null;comment:生日的月份" json:"birth_month"`
	BirthDay        string `gorm:"column:birth_day;type:int(10);not null;comment:生日的天" json:"birth_day"`
	CustomerCode    string `json:"customer_code" gorm:"column:customer_code;type:varchar(255);not null;comment:cdp的用户customer_code"`
}

type WxUserCorpExternalModelInfo struct {
	model.BranchModel
	Appid           string `gorm:"column:appid;type:char(64);not null;uniqueIndex:uin;comment:对应的是corp_id"`
	Openid          string `gorm:"column:openid;type:char(64);not null;uniqueIndex:uin;comment:对应的是external_user_id"`
	Unionid         string `gorm:"column:unionid;type:char(64);not null;comment:用户unionid"`
	Euid            int64  `gorm:"column:euid;type:bigint(20);not null;index:euid;default:0;comment:聚合id"`
	Nick            string `gorm:"column:nick;type:varchar(255);not null;comment:昵称"`
	Name            string `gorm:"column:name;type:varchar(255);not null;comment:姓名"`
	Position        string `gorm:"column:position;type:varchar(255);not null;comment:用户职称"`
	Avatar          string `gorm:"column:avatar;type:varchar(255);not null;comment:头像"`
	Birthday        string `gorm:"column:birthday;type:varchar(50);not null;comment:生日"`
	UnionCode       string `gorm:"column:union_code;type:varchar(50);not null;comment:会员卡号"`
	Point           string `gorm:"column:point;type:varchar(10);not null;comment:会员积分"`
	Grade           string `gorm:"column:grade;type:varchar(20);not null;comment:会员等级"`
	StoreCode       string `gorm:"column:store_code;type:varchar(50);not null;comment:门店编号"`
	BaCode          string `gorm:"column:ba_code;type:varchar(50);not null;comment:BA编号"`
	CorpName        string `gorm:"column:corp_name;type:varchar(255);not null;comment:企业简称"`
	CorpFullName    string `gorm:"column:corp_full_name;type:varchar(255);not null;comment:企业全称"`
	Type            int32  `gorm:"column:type;type:smallint(4);not null;default:0;comment:外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户"`
	Sex             int32  `gorm:"column:sex;type:smallint(4);not null;default:0;comment:外部联系人性别 0-未知 1-男性 2-女性"`
	ExternalProfile string `gorm:"column:external_profile;type:text;not null;comment:外部联系人的自定义展示信息,JSON"`
}

// TableName returns the table name of the struct_name model
func (WxUserCorpExternalModel) TableName() string {
	return "wx_user_corp_external"
}
func (w WxUserCorpExternalModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", w.DBPrefix(), enterpriseID, w.TableName())
}

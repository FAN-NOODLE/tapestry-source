package microScrmModel

import (
	"fmt"
	"micro-scrm/internal/config"
	"time"

	"gorm.io/gorm"
)

type EnterpriseUserModel struct {
	ID                     int64     `gorm:"primaryKey;column:id;type:bigint(20)" json:"id"`
	CreatedAt              int64     `gorm:"column:created_at;type:int(10);not null;default:0" json:"created_at"`
	UpdatedAt              int64     `gorm:"column:updated_at;type:int(10);not null;default:0" json:"updated_at"`
	WxUnion                string    `gorm:"column:wx_union;type:text;not null;comment:关联用户列表"`
	WxShop                 string    `gorm:"column:wx_shop;type:text;comment:商城信息集合"`
	WxExternal             string    `gorm:"column:wx_external;type:text;comment:外部联系人集合"`
	WxOpen                 string    `gorm:"column:wx_open;type:text;comment:公众号用户信息集合"`
	Nick                   string    `gorm:"column:nick;type:varchar(255);not null;comment:昵称"`
	Name                   string    `gorm:"column:name;type:varchar(255);not null;comment:姓名"`
	Tags                   string    `gorm:"column:tags;type:text;comment:标签列表"`
	Groups                 string    `gorm:"column:groups;type:text;comment:rfm/aipl分组"`
	Avatar                 string    `gorm:"column:avatar;type:varchar(255);not null;comment:头像"`
	Mobile                 string    `gorm:"column:mobile;type:varchar(20);comment:手机"`
	Birthday               string    `gorm:"column:birthday;type:varchar(50);not null;comment:生日"`
	UnionCode              string    `gorm:"column:union_code;type:varchar(50);not null;comment:会员卡号"`
	Point                  string    `gorm:"column:point;type:varchar(10);not null;comment:会员积分"`
	Grade                  string    `gorm:"column:grade;type:varchar(20);not null;comment:会员等级"`
	StoreCode              string    `gorm:"column:store_code;type:varchar(50);not null;comment:门店编号"`
	BaCode                 string    `gorm:"column:ba_code;type:varchar(50);not null;comment:BA编号"`
	DeliveryMobile         string    `gorm:"column:delivery_mobile;type:varchar(20);comment:最近收货人手机"`
	Sex                    int32     `gorm:"column:sex;type:smallint(4);not null;default:0;comment:用户的性别，值为1时是男性，值为2时是女性，值为0时是未知"`
	Language               string    `gorm:"column:language;type:varchar(20);comment:语言"`
	Country                string    `gorm:"column:country;type:varchar(255);not null;comment:国家"`
	Province               string    `gorm:"column:province;type:varchar(255);not null;comment:省份"`
	City                   string    `gorm:"column:city;type:varchar(255);not null;comment:城市"`
	Remark                 string    `gorm:"column:remark;type:text;comment:备注"`
	Version                int       `gorm:"column:version;type:int(11);not null;default:0"`
	PayOrder               int       `gorm:"column:pay_order;type:int(11);not null;default:0;comment:已支付订单数"`
	TotalFee               int       `gorm:"column:total_fee;type:int(11);not null;default:0;comment:实付金额总数"`
	RefundOrder            int       `gorm:"column:refund_order;type:int(11);not null;default:0;comment:退款订单数"`
	RefundFee              int       `gorm:"column:refund_fee;type:int(11);not null;default:0;comment:退款金额数"`
	PerCustomerTransaction int       `gorm:"column:per_customer_transaction;type:int(11);not null;default:0;comment:客单价"`
	LatestMarketTime       int       `gorm:"column:latest_market_time;type:int(11);not null;default:0;comment:最新营销时间"`
	MarketCount            int       `gorm:"column:market_count;type:int(11);not null;default:0;comment:营销次数"`
	BindTo                 int64     `gorm:"column:bind_to;type:bigint(20);not null;default:0;comment:绑定对象"`
	LevelID                int64     `gorm:"column:level_id;type:bigint(20);not null;default:0"`
	EsUpdate               time.Time `gorm:"autoUpdateTime,column:es_update;type:timestamp;index:es"`
	DtUpdate               time.Time `gorm:"column:dt_update;type:timestamp;index:dt;default:CURRENT_TIMESTAMP"`
}

// TableName 企业用户聚合表
func (EnterpriseUserModel) TableName() string {
	return "enterprise_user"
}

// BranchTable 切表
func (m EnterpriseUserModel) BranchTable(enterpriseID int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s", config.CfgData.DB.DBPrefix, enterpriseID, m.TableName())
		return tx.Table(tableName)
	}
}

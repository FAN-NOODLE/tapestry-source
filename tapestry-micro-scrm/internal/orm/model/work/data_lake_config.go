package work

import "micro-scrm/internal/orm/model"

const (
	DLC_DELETED = iota
)

type DataLakeConfigModel struct {
	model.MicroScrmBaseModel
	Vendor       string `gorm:"column:vendor;NOT NULL" json:"vendor"`               // 第三方名称
	BrandCode    string `gorm:"column:brand_code;NOT NULL" json:"brand_code"`       // 品牌编号
	ProgramCode  string `gorm:"column:program_code;NOT NULL" json:"program_code"`   // 俱乐部编号
	StoreCode    string `gorm:"column:store_code;NOT NULL" json:"store_code"`       // 门店号
	SellerName   string `gorm:"column:seller_name;NOT NULL" json:"seller_name"`     // 小程序账号主体
	ClientId     string `gorm:"column:client_id;NOT NULL" json:"client_id"`         // httpBase 认证账号
	ClientSecret string `gorm:"column:client_secret;NOT NULL" json:"client_secret"` // httpBase 认证密码
	AuthUser     string `gorm:"column:auth_user;NOT NULL" json:"auth_user"`         // 获取token账号
	AuthPassword string `gorm:"column:auth_password;NOT NULL" json:"auth_password"` // 获取token密码
}

func (DataLakeConfigModel) TableName() string {
	return "data_lake_config"
}

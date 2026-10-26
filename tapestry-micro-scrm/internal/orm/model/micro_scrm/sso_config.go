package micro_scrm

import "micro-scrm/internal/orm/model"

// SSOConfigModel  企业SSO配置模型
type SSOConfigModel struct {
	model.MicroScrmBaseModel
	TenantID     string `json:"tenant_id" gorm:"column:tenant_id;type:varchar(255);not null;comment:租户ID"`
	ClientID     string `json:"client_id" gorm:"column:client_id;type:varchar(255);not null;comment:客户端ID"`
	ClientSecret string `json:"client_secret" gorm:"column:client_secret;type:varchar(512);not null;comment:客户端密钥"`
	RedirectURI  string `json:"redirect_uri" gorm:"column:redirect_uri;type:varchar(512);not null;comment:回调地址"`
	Scope        string `json:"scope" gorm:"column:scope;type:varchar(255);default:'User.Read';comment:授权范围"`
}

// TableName 设置表名
func (SSOConfigModel) TableName() string {
	return "sso_config"
}

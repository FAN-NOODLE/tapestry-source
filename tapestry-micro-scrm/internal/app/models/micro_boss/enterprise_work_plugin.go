package microWorkModel

type EnterpriseWorkPluginModel struct {
	ID            uint64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned"`
	EnterpriseID  uint64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;index:enterprise_work_plugin_enterprise_id_IDX;default:0;comment:企业id"`
	CorpID        string `json:"corp_id" gorm:"column:corp_id;type:varchar(255);not null;comment:企业corpid"`
	AgentID       string `json:"agent_id" gorm:"column:agent_id;type:varchar(255);not null;comment:应用id"`
	Secret        string `json:"secret" gorm:"column:secret;type:varchar(255);not null;comment:应用secret"`
	Identifier    string `json:"identifier" gorm:"column:identifier;type:varchar(255);not null;comment:应用标识"`
	CheckFile     int32  `json:"check_file" gorm:"column:check_file;type:tinyint(4);not null;default:2;comment:域名校验文件是否上传 1 是 2否"`
	FileURL       string `json:"file_url" gorm:"column:file_url;type:varchar(255);not null;comment:证书文件"`
	PlPluginID    string `json:"pl_plugin_id" gorm:"column:pl_plugin_id;type:varchar(255);not null;comment:uc返回id"`
	MiniAppID     string `json:"mini_app_id" gorm:"column:mini_app_id;type:varchar(255);not null;comment:应用关联小程序appid"`
	UpdatedAt     uint   `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0"`
	CreatedAt     uint   `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0"`
	DeletedAt     uint   `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;index:enterprise_work_plugin_enterprise_id_IDX;default:0"`
	MiniAppSecret string `json:"mini_app_secret" gorm:"column:mini_app_secret;type:varchar(255);not null;comment:应用关联小程序的secret"`
}

// TableName returns the table name of the EnterpriseWorkPlugin model
func (EnterpriseWorkPluginModel) TableName() string {
	return "enterprise_work_plugin"
}

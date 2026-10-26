package micro_scrm

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/pkg/sftp"
)

const (
	SftpTypePlan      = iota + 1 // 1 计划任务
	SftpTypeBa                   // 2 ba业绩看板
	SftpTypeCommunity            // 3 社群营销数据
	SftpTypeAcxiom               // 4 Sftp 对接Acxiom信息，（Coach,SW）
	SftpTypeWedo                 // 5 Sftp 对接Wedo信息，（KS）
	SftpTypeDEL                  // 6 Sftp 对接DEL信息 （Coach）
	SftpTypeKUNXIN               // 7 Sftp 对接DEL信息上传数据给到KUNXIN （Coach）
	SftpTypeCdpTag               // 8 Sftp CDP全量标签 （Coach）
	SftpTypeWorkday              // 9 Sftp 对接Workday员工主数据（Coach）
	SftpTypeDELALI = 12 // 12 DataFabric 阿里云 SFTP（Coach）

)

type BrandConfigSftpModel struct {
	model.MicroScrmBaseModel

	BrandCode  string `json:"brand_code" gorm:"column:brand_code;type:varchar(20);not null;comment:品牌简称"` // 品牌简称
	BrandName  string `json:"brand_name" gorm:"column:brand_name;type:varchar(50);not null;comment:品牌名称"` // 品牌名称
	SftpType   int    `json:"sftp_type"  gorm:"column:sftp_type;default:0"`                               // 配置类型:1 计划任务 2 ba业绩看板 3 社群营销数据
	SftpConfig string `json:"sftp_config" gorm:"column:sftp_config;type:text;comment:SFTP配置"`             // SFTP配置
	SftpPath   string `json:"sftp_path" gorm:"column:sftp_path;type:varchar(255);not null;default:根目录配置"`
}

func (BrandConfigSftpModel) TableName() string {
	return "brand_config_sftp"
}

type BrandConfigSftpDetail struct {
	BrandConfigSftpModel

	ConfigDetail sftp.SSHConfig
	PathConfig   SFTPPathConfig
}

type SFTPPathConfig struct {
	DataFabric       string   `json:"data_fabric" comment:"DataFabric上传基础目录，不含日期子目录"`
	MemberSync       string   `json:"member_sync"`
	PlanTask         []string `json:"plan_task"`
	TagLog           string   `json:"tag_log"`
	BaKanBan         string   `json:"ba_kan_ban"`
	UploadFile       string   `json:"upload_file"`
	CustomerLog      string   `json:"customer_log"`
	Acxiom           []string `json:"acxiom" comment:"coach和安客诚的sftp对接"`
	CdpExternalUser  []string `json:"cdp_external_user" comment:"安客诚cdp外部联系人会员信息"`
	ChangYiReport    []string `json:"chang_yi_report" comment:"vchangyi将report数据上传sftp给安客诚"`
	OrderAchievement string   `json:"order_achievement" comment:"sa订单的sftp对接"`
	SwCrmPurchase    []string `json:"sw_crm_purchase" comment:"SW安客诚会员订单及附加信息"`
}

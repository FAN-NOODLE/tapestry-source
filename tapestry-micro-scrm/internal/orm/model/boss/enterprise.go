package boss

type EnterpriseModel struct {
	ID                int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned;comment:企业ID-自增"`
	Hash              string `json:"hash" gorm:"column:hash;type:varchar(50);not null;uniqueIndex:hash_unique;comment:企业HASH"`
	Name              string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:企业名称"`
	Contact           string `json:"contact" gorm:"column:contact;type:varchar(40);not null;comment:企业联系人"`
	Mobile            string `json:"mobile" gorm:"column:mobile;type:varchar(20);not null;comment:企业联系手机号"`
	Address           string `json:"address" gorm:"column:address;type:varchar(255);not null;comment:详细地址"`
	Logo              string `json:"logo" gorm:"column:logo;type:varchar(255);not null;comment:企业LOGO"`
	CorpID            string `json:"corp_id" gorm:"column:corp_id;type:varchar(100);not null;comment:企业Corpid"`
	Domain            string `json:"domain" gorm:"column:domain;type:varchar(32);not null;comment:企业域名"`
	SeparateSwitch    uint32 `json:"separate_switch" gorm:"column:separate_switch;type:tinyint(4) unsigned;not null;default:1;comment:分账开关：1未开启 2开启"`
	IndustryID        uint64 `json:"industry_id" gorm:"column:industry_id;type:bigint(20) unsigned;not null;default:0;comment:行业ID"`
	LockTime          uint   `json:"lock_time" gorm:"column:lock_time;type:int(11) unsigned;not null;default:1440;comment:时间分 每单锁定这个时间,到了之后进行分账"`
	Proportion        uint   `json:"proportion" gorm:"column:proportion;type:int(11) unsigned;not null;default:0;comment:分账比例 精度小数点后两位, 实际百分比为 值/10000"`
	StartFee          uint   `json:"start_fee" gorm:"column:start_fee;type:int(11) unsigned;not null;default:0;comment:分账起始金额, 大于此金额分账"`
	CrsStatus         uint32 `json:"crs_status" gorm:"column:crs_status;type:tinyint(3) unsigned;not null;default:1;comment:CRS小程序状态（1:正常，2:禁用）"`
	StoreStatus       uint32 `json:"store_status" gorm:"column:store_status;type:tinyint(3) unsigned;not null;default:1;comment:商城小程序状态（1:正常，2:禁用）"`
	ContactStatus     uint   `json:"contact_status" gorm:"column:contact_status;type:int(10) unsigned;not null;default:1;comment:企业微信通讯录应用状态（1:正常，2:禁用, 3 自建应用）"`
	ExternalSecret    string `json:"external_secret" gorm:"column:external_secret;type:varchar(255);not null;comment:客户联系secret"`
	Status            uint32 `json:"status" gorm:"column:status;type:tinyint(3) unsigned;not null;default:1;comment:状态（1:正常，2:禁用）"`
	BusinessSource    uint32 `json:"business_source" gorm:"column:business_source;type:tinyint(3) unsigned;not null;default:0;comment:业务来源; 0:-,1:企业微信,10:畅移（直签）,11:畅移（代理）,12:畅移（试用）20:畅鲟（直签）21:畅鲟（代理）22:畅鲟（试用）"`
	Version           uint32 `json:"version" gorm:"column:version;type:tinyint(4) unsigned;not null;default:1;comment:1:标准版（旧版）,2:中心商城版(有门店), 3:中心商城, 4:直营连锁门店"`
	VersionTag        string `json:"version_tag" gorm:"column:version_tag;type:varchar(16);not null;comment:版本标签;正式版,试用版"`
	ErpStatus         uint   `json:"erp_status" gorm:"column:erp_status;type:int(10) unsigned;not null;default:1;comment:ERP接口调用权限 1：可以 2：不可以"`
	CrmStatus         uint32 `json:"crm_status" gorm:"column:crm_status;type:tinyint(3) unsigned;not null;default:1;comment:CRM接口调用权限 1：可以 2：不可以"`
	ReturnAddress     string `json:"return_address" gorm:"column:return_address;type:varchar(255);not null;comment:退货地址"`
	ServiceMerchantID uint64 `json:"service_merchant_id" gorm:"column:service_merchant_id;type:bigint(20) unsigned;not null;default:1;comment:服务商ID"`
	OpenPlatformID    uint64 `json:"open_platform_id" gorm:"column:open_platform_id;type:bigint(20) unsigned;not null;default:0;comment:企业关联的开放平台id"`
	AccountLimit      uint   `json:"account_limit" gorm:"column:account_limit;type:int(10) unsigned;not null;default:0;comment:最大允许创建账户数量"`
	AccountStartAt    uint   `json:"account_start_at" gorm:"column:account_start_at;type:int(10) unsigned;not null;default:0;comment:账户有效期开始日期"`
	AccountExpireAt   uint   `json:"account_expire_at" gorm:"column:account_expire_at;type:int(10) unsigned;not null;default:0;comment:账户有效期失效日期"`
	IsSplitFlow       uint32 `json:"is_split_flow" gorm:"column:is_split_flow;type:tinyint(1) unsigned;not null;default:0;comment:是否切流"`
	AvailableSystem   string `json:"available_system" gorm:"column:available_system;type:text;not null;comment:企业可用系统,json字符存储[1-CRS后台 2-店长小程序 3-天擎CDP]"`
	EditAt            uint   `json:"edit_at" gorm:"column:edit_at;type:int(10) unsigned;not null;default:0;comment:最后操作时间"`
	CreatedAt         uint   `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt         uint   `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:更新时间"`
	DeletedAt         uint   `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;default:0;comment:删除时间"`
}

func (EnterpriseModel) TableName() string {
	return "enterprise"
}

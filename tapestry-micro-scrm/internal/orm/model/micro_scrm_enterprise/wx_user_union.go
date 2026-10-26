package micro_scrm_enterprise

type WxUserUnionModel struct {
	Unionid            string `gorm:"primaryKey;column:unionid;type:char(64);comment:微信开放平台unionid"`
	OpenPlatformID     string `gorm:"column:open_platform_id;type:char(64);not null;comment:unionid所属开放平台"`
	OpenPlatformSource int32  `gorm:"column:open_platform_source;type:smallint(4);not null;default:0;comment:unionid归属开放平台是否确认1确认 0当前配置"`
	Euid               int64  `gorm:"column:euid;type:bigint(20);not null;index:euid;default:0;comment:企业聚合用户id"`
}

// TableName 开放平台unionid映射聚合id
func (WxUserUnionModel) TableName() string {
	return "wx_user_union"
}

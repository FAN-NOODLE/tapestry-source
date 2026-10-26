package externalSearchModel

const (
	CdpTagAmountTotal           = "cdpl12mzje"       //一年总金额
	CdpTagHistoryAmountTotal    = "cdplsgmzjeltv"    //历史总金额
	CdpTagLastPurchaseDate      = "cdpzjycxfrq"      //最后一次消费日期
	CdpTagLastPurchaseStoreCode = "cdpzjycgmdp"      //最后一次消费店铺
	CdpTagCategoryPreference    = "cdplljgmtop2plph" //品类偏好
	CdpTagPurchasedDate         = "cdpr24lsxfsj"     //历史消费日期
	CdpTagHighWilling           = "cdpgyygk"         //高意愿
	CdpTagPriceSensitivity      = "cdpjgmgd"         //价格敏感度
	CdpTagBuyCategoryPreference = "cdpgmplphcs"      //购买品类偏好
	CdpTagPurchasedMonth        = "cdpgmsjfb"        //购买时间分布
	CdpTagManagerStoreCode      = "cmcc"             //管理柜台
)

type CdpExternalTag struct {
	Unionid  string `gorm:"column:unionid;type:varchar(255);not null;comment:外部联系人的unionid"`
	TagType  string `gorm:"column:tag_type;type:varchar(255);not null;comment:标签类型"`
	TagValue string `gorm:"column:tag_value;type:varchar(255);not null;comment:成员设置的外部联系人备注"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (CdpExternalTag) TableName() string {
	return "cdp_external_tag"
}

package microWorkModel

type EnterpriseConfigModel struct {
	ID           uint64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned;comment:id-自增"`
	EnterpriseID uint64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;index:enterprise_config_enterprise_id_index;default:0;comment:企业id"`
	Type         uint64 `json:"type" gorm:"column:type;type:bigint(20) unsigned;not null;default:0;comment:类别(1：好物圈设置，2：商品详情页设置，3：下单小程序授权，4：购物袋设置，5：直播设置，6：模板设置，7：第三方会员开关, 8:有数, 9:王者人生  10 积分名称  11 发券插件 注意:新增设置时开发人员需要维护该值)"`
	Key          string `json:"key" gorm:"column:key;type:varchar(100);not null;index:key_idx;comment:字段标示"`
	Value        uint   `json:"value" gorm:"column:value;type:int(10) unsigned;not null;default:0;comment:1=开启 2=关闭"`
	Desc         string `json:"desc" gorm:"column:desc;type:varchar(100);not null;default:0;comment:字段说明"`
	CreatedAt    uint   `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt    uint   `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:更新时间"`
	DeletedAt    uint   `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;default:0;comment:删除时间"`
	Name         string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:品牌标识"`
}

// TableName returns the table name of the EnterpriseConfig model
func (EnterpriseConfigModel) TableName() string {
	return "enterprise_config"
}

package microWorkModel

import "micro-scrm/internal/app/models"

// store_type 门店类别：1奥莱 2零售
const (
	StoreTypeOle = iota + 1
	StoreTypeRetail
)

const (
	StoreTypeOleStr    = "Outlet"
	StoreTypeRetailStr = "Retail"
	StoreTypeESA       = "ESA"
)

type StoreDepartmentModel struct {
	models.MicroBaseModel

	ThirdStoreCode     string `json:"third_store_code" gorm:"column:third_store_code;type:varchar(32);not null;uniqueIndex:unique_store_code;comment:三方门店id"`
	UcDepartmentID     string `json:"uc_department_id" gorm:"column:uc_department_id;type:varchar(64);not null;comment:uc部门id"`
	WechatDepartmentID uint   `json:"wechat_department_id" gorm:"column:wechat_department_id;type:int(10) unsigned;not null;default:0;comment:企微部门id"`
	StoreName          string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;default:0;comment:门店名称"`
	StoreType          int32  `json:"store_type" gorm:"column:store_type;type:tinyint(1);not null;default:0;comment:门店类别：1奥莱 2零售"`
	OperateName        string `json:"operate_name" gorm:"column:operate_name;type:varchar(255);not null;comment:操作人"`
	SyncStatus         int32  `json:"sync_status" gorm:"column:sync_status;type:tinyint(1);not null;default:0;comment:同步状态"`
	SyncAt             int    `json:"sync_at" gorm:"column:sync_at;type:int(10);not null;default:0;comment:上次同步时间"`
	Address            string `json:"address" gorm:"column:address;type:varchar(255);not null;default:'';comment:门店地址"`
	City               string `json:"city" gorm:"column:city;type:varchar(255);not null;comment:城市"`
	Phone              string `json:"phone" gorm:"column:phone;type:varchar(255);not null;comment:门店电话"`
	InventoryName      string `json:"inventory_name" gorm:"column:inventory_name;type:varchar(255);not null;comment:库存相关门店名称"`
	Email              string `json:"email" gorm:"column:email;type:varchar(255);not null;comment:门店邮箱"`
	ClosedAt           int64  `json:"closed_at" gorm:"column:closed_at;type:int(10);not null;default:0;comment:闭店时间"`
	ProvinceZh         string `json:"province_zh" gorm:"column:province_zh;type:varchar(255);not null;comment:省/直辖市(中文)"`
	CityZh             string `json:"city_zh" gorm:"column:city_zh;type:varchar(255);not null;comment:市/区(中文)"`
	Lat                string `json:"lat" gorm:"column:lat;type:varchar(255);not null;comment:纬度"`
	Lng                string `json:"lng" gorm:"column:lng;type:varchar(255);not null;comment:经度"`
}

// TableName 三方门店表（业务架构）
func (StoreDepartmentModel) TableName() string {
	return "store_department"
}

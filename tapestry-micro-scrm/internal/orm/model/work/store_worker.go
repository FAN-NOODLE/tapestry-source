package work

import "micro-scrm/internal/orm/model"

const (
	SW_TYPE_STORE = iota + 1
	SW_TYPE_GROUP
	SW_TYPE_ENTERPRISE
	SW_TYPE_ENTERPRISE_MANAGE
)

const (
	SW_WORKER_STATUS_ON = iota + 1
	SW_WORKER_STATUS_RESIGN
)

type StoreWorkerModel struct {
	model.MicroScrmBaseModel

	Type                  uint8  `json:"type" gorm:"column:type;type:tinyint(3) unsigned;not null;default:1;comment:导购类型：1、门店导购，2、社群导购，3、企业导购"`
	UserID                int64  `json:"user_id" gorm:"column:user_id;type:bigint(20) unsigned;not null;default:0;comment:用户id"`
	StoreID               int64  `json:"store_id" gorm:"column:store_id;type:bigint(20) unsigned;not null;default:0;comment:门店ID"`
	WorkerNo              string `json:"worker_no" gorm:"column:worker_no;type:varchar(120);not null;comment:员工编号"`
	Name                  string `json:"name" gorm:"column:name;type:varchar(120);not null;comment:名称"`
	Mobile                string `json:"mobile" gorm:"column:mobile;type:varchar(150);not null;comment:电话"`
	MemUID                string `json:"mem_uid" gorm:"column:mem_uid;type:varchar(150);not null;comment:员工uid"`
	MemUserid             string `json:"mem_userid" gorm:"column:mem_userid;type:varchar(150);not null;comment:员工企业微信userid"`
	ContactQrCodeID       int64  `json:"contact_qr_code_id" gorm:"column:contact_qr_code_id;type:bigint(20) unsigned;not null;default:0;comment:联系我二维码id"`
	Nickname              string `json:"nickname" gorm:"column:nickname;type:varchar(120);not null;comment:员工昵称"`
	Avatar                string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;comment:头像"`
	SmsMobile             string `json:"sms_mobile" gorm:"column:sms_mobile;type:varchar(20);not null;comment:接收短信手机号"`
	IsSendSms             uint8  `json:"is_send_sms" gorm:"column:is_send_sms;type:tinyint(3) unsigned;not null;default:2;comment:是否接收提醒短信1:发送 2：不发送"`
	QrCode                string `json:"qr_code" gorm:"column:qr_code;type:varchar(255);not null;comment:企业微信二维码链接"`
	StoreSunCode          string `json:"store_sun_code" gorm:"column:store_sun_code;type:varchar(255);not null;comment:太阳码链接"`
	IsBoss                uint8  `json:"is_boss" gorm:"column:is_boss;type:tinyint(3) unsigned;not null;default:0;comment:是否是店长：1=是, 2=否"`
	Source                uint8  `json:"source" gorm:"column:source;type:tinyint(3) unsigned;not null;default:0;comment:来源：1=企业微信"`
	WorkerStatus          uint8  `json:"worker_status" gorm:"column:worker_status;type:tinyint(3) unsigned;not null;default:0;comment:员工状态 1=在职，2=离职"`
	EntryTime             int64  `json:"entry_time" gorm:"column:entry_time;type:int(10) unsigned;not null;default:0;comment:入职时间"`
	OutTime               int64  `json:"out_time" gorm:"column:out_time;type:int(10) unsigned;not null;default:0;comment:离职时间"`
	RecentlyLoggedStoreID int64  `json:"recently_logged_store_id" gorm:"column:recently_logged_store_id;type:bigint(20) unsigned;not null;default:0;comment:最近登录门店id"`
	RecentlyLoggedTime    int64  `json:"recently_logged_time" gorm:"column:recently_logged_time;type:int(10) unsigned;not null;default:0;comment:最近登录门店时间"`
}

func (StoreWorkerModel) TableName() string {
	return "store_worker"
}

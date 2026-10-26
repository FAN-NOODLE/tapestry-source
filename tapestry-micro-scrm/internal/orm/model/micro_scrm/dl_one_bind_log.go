package micro_scrm

type DLOneBindLogModel struct {
	ID                 int    `json:"id" gorm:"primaryKey;column:id;type:int(20) auto_increment"`
	EnterpriseID       int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;default:0"`
	WorkExternalUserID string `json:"work_external_user_id" gorm:"column:work_external_user_id;type:char(64);not null"`
	UnionCode          string `json:"union_code" gorm:"column:union_code;type:varchar(50);not null;comment:会员卡号"`                            // 会员卡号
	UnionID            string `json:"union_id" gorm:"column:union_id;type:char(64);not null;comment:用户unionid"`                              // 用户unionid
	BaUserID           string `json:"ba_user_id" gorm:"column:ba_user_id;type:varchar(64);not null;comment:userid"`                          // userid
	BindType           int    `json:"bind_type" gorm:"column:bind_type;type:tinyint(2);not null;default:0;comment:1绑定，2解绑"`                  // 1绑操作
	BindStatus         int    `json:"bind_status" gorm:"column:bind_status;type:tinyint(2);not null;default:0;comment:0无意义；1 一绑；2 二绑；3专属BA"` // 1绑操作
	OpDate             string `json:"op_date" gorm:"column:op_date;type:char(20);not null"`
	BaCode             string `json:"ba_code" gorm:"column:ba_code;type:varchar(64);not null;comment:ba_code"`
	Status             int    `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:激活状态: 0 没有跑 1 成功"` // 激活状态: 0 没有跑 1 成功
	CreatedAt          int    `json:"created_at" gorm:"column:created_at;type:int(10);not null;default:0"`
	UpdatedAt          int    `json:"updated_at" gorm:"column:updated_at;type:int(10);not null;default:0"`
	DeletedAt          int    `json:"deleted_at" gorm:"column:deleted_at;type:int(10);not null;default:0"`
}

func (DLOneBindLogModel) TableName() string {
	return "dl_one_bind_log"
}

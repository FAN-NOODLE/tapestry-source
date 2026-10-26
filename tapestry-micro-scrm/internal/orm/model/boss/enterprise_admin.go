package boss

type EnterpriseAdminModel struct {
	Id                   int64  `gorm:"column:id;primary_key" json:"id"`                                        // 企业管理员ID-自增
	EnterpriseId         int64  `gorm:"column:enterprise_id;default:0;NOT NULL" json:"enterprise_id"`           // 企业ID
	Username             string `gorm:"column:username;NOT NULL" json:"username"`                               // 用户名
	Account              string `gorm:"column:account;NOT NULL" json:"account"`                                 // 账户
	Password             string `gorm:"column:password;NOT NULL" json:"password"`                               // 密码
	Salt                 string `gorm:"column:salt;NOT NULL" json:"salt"`                                       // 盐
	Mobile               string `gorm:"column:mobile;NOT NULL" json:"mobile"`                                   // 手机号码
	Email                string `gorm:"column:email;NOT NULL" json:"email"`                                     // 邮箱（目前作为超管账号)
	GlobalRoaming        string `gorm:"column:global_roaming;NOT NULL" json:"global_roaming"`                   // 国际区号
	Role                 int    `gorm:"column:role;default:0;NOT NULL" json:"role"`                             // 角色（1：系统管理员）
	IsSuper              int    `gorm:"column:is_super;default:2;NOT NULL" json:"is_super"`                     // 是否超管: 1-是 2-不是
	LastLoginTime        int    `gorm:"column:last_login_time;default:0;NOT NULL" json:"last_login_time"`       // 最后一次登录时间
	Status               int    `gorm:"column:status;default:1;NOT NULL" json:"status"`                         // 状态（1:正常，2:禁用）
	OperatorId           int    `gorm:"column:operator_id;default:0;NOT NULL" json:"operator_id"`               // 操作人ID
	MemberId             string `gorm:"column:member_id;NOT NULL" json:"member_id"`                             // UC用户ID
	CreatedAt            int    `gorm:"column:created_at;default:0;NOT NULL" json:"created_at"`                 // 创建时间
	UpdatedAt            int    `gorm:"column:updated_at;default:0;NOT NULL" json:"updated_at"`                 // 更新时间
	DeletedAt            int    `gorm:"column:deleted_at;default:0;NOT NULL" json:"deleted_at"`                 // 删除时间
	RangeType            int    `gorm:"column:range_type;default:1;NOT NULL" json:"range_type"`                 // 账号范围 1-不限 2-指定范围
	AdminAvailableSystem string `gorm:"column:admin_available_system;NOT NULL" json:"admin_available_system"`   // 账号可用系统,json字符存储[1-CRS后台 2-店长小程序 3-天擎CDP]
	PasswordExpireTime   int    `gorm:"column:password_expire_time;default:0;NOT NULL" json:"password_expire_time"` // 密码过期时间，如果过期了则需要重新修改密码，默认90天过期
}

func (EnterpriseAdminModel) TableName() string {
	return "enterprise_admin"
}

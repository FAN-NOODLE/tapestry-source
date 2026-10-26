package microWorkModel

import "micro-scrm/internal/app/models"

const (
	// EnterpriseAdminEnable 员工状态 1=正常
	// EnterpriseAdminDisable 员工状态 2=禁用
	EnterpriseAdminEnable = iota + 1
	EnterpriseAdminDisable
)

type EnterpriseAdminModel struct {
	models.MicroBaseModel
	Username             string `json:"username" gorm:"column:username;type:varchar(40);not null;comment:用户名"`
	Account              string `json:"account" gorm:"column:account;type:varchar(50);not null;index:enterprise_admin_account_index;comment:账户"`
	Password             string `json:"password" gorm:"column:password;type:varchar(32);not null;comment:密码"`
	Salt                 string `json:"salt" gorm:"column:salt;type:varchar(6);not null;comment:盐"`
	Mobile               string `json:"mobile" gorm:"column:mobile;type:varchar(32);not null;comment:手机号码"`
	Email                string `json:"email" gorm:"column:email;type:varchar(50);not null;comment:邮箱（目前作为超管账号)"`
	GlobalRoaming        string `json:"global_roaming" gorm:"column:global_roaming;type:varchar(32);not null;comment:国际区号"`
	Role                 uint   `json:"role" gorm:"column:role;type:int(10) unsigned;not null;default:0;comment:角色（1：系统管理员）"`
	IsSuper              int32  `json:"is_super" gorm:"column:is_super;type:tinyint(1);not null;default:2;comment:是否超管: 1-是 2-不是"`
	LastLoginTime        uint   `json:"last_login_time" gorm:"column:last_login_time;type:int(10) unsigned;not null;default:0;comment:最后一次登录时间"`
	Status               uint32 `json:"status" gorm:"column:status;type:tinyint(3) unsigned;not null;default:1;comment:状态（1:正常，2:禁用）"`
	OperatorID           uint   `json:"operator_id" gorm:"column:operator_id;type:int(10) unsigned;not null;default:0;comment:操作人ID"`
	MemberID             string `json:"member_id" gorm:"column:member_id;type:varchar(32);not null;comment:UC用户ID"`
	RangeType            uint32 `json:"range_type" gorm:"column:range_type;type:tinyint(4) unsigned;not null;default:1;comment:账号范围 1-不限 2-指定范围"`
	AdminAvailableSystem string `json:"admin_available_system" gorm:"column:admin_available_system;type:text;not null;comment:账号可用系统,json字符存储[1-CRS后台 2-店长小程序 3-天擎CDP]"`
	PasswordExpireTime   int    `json:"password_expire_time" gorm:"column:password_expire_time;type:int(10);not null;default:0;comment:密码过期时间，如果过期了则需要重新修改密码，默认90天过期"`
}

// TableName 企业管理员表
func (EnterpriseAdminModel) TableName() string {
	return "enterprise_admin"
}

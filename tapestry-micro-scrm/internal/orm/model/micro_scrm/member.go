package micro_scrm

import "micro-scrm/internal/orm/model"

type MemberModel struct {
	model.MicroScrmBaseModel

	Userid      string `gorm:"column:userid;type:varchar(64);not null;index:user_id_index;comment:userid"`
	Name        string `gorm:"column:name;type:varchar(128);not null;comment:企业微信姓名"`
	Mobile      string `gorm:"column:mobile;type:varchar(15);not null;index:mobile_index;comment:企业微信手机号"`
	Position    string `gorm:"column:position;type:varchar(255);not null;comment:职务"`
	Gender      int32  `gorm:"column:gender;type:tinyint(2);not null;default:0;comment:性别0=未知 1=男性 2=女性"`
	Email       string `gorm:"column:email;type:varchar(128);not null;index:email_index;comment:邮箱"`
	Avatar      string `gorm:"column:avatar;type:varchar(512);not null;comment:员工头像"`
	ThumbAvatar string `gorm:"column:thumb_avatar;type:varchar(128);not null;comment:缩略图"`
	Telephone   string `gorm:"column:telephone;type:varchar(32);not null;comment:座机"`
	Alias       string `gorm:"column:alias;type:varchar(64);not null;comment:别名"`
	Status      int32  `gorm:"column:status;type:tinyint(2);not null;default:0;comment:激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业"`
	QrCode      string `gorm:"column:qr_code;type:varchar(128);not null;comment:员工个人二维码"`
	Address     string `gorm:"column:address;type:varchar(255);not null;comment:地址"`
	OpenUserid  string `gorm:"column:open_userid;type:varchar(64);not null;default:open_userid"`
}

// TableName returns the table name of the struct_name model
func (MemberModel) TableName() string {
	return "members"
}

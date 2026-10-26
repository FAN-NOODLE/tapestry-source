package membersModel

import (
	"micro-scrm/internal/app/models"
)

const (
	MembersStatusCreate      = 0 // 0=刚创建
	MembersStatusEnable      = 1 // 1=已激活
	MembersStatusDisable     = 2 // 2=已禁用
	MembersStatusInactivated = 4 // 4=未激活
	MembersStatusLeave       = 5 // 5=退出企业
)

const (
	MembersTypeSA = 3 // 导购或店员
	MembersTypeSM = 4 // 店长或管理
)

type MemberModel struct {
	models.MicroBaseModel

	Userid         string `json:"userid" gorm:"primaryKey;column:userid;type:varchar(64);index:userid_index;comment:userid"`
	Name           string `json:"name" gorm:"column:name;type:varchar(128);not null;comment:企业微信姓名"`
	Mobile         string `json:"mobile" gorm:"column:mobile;type:varchar(15);not null;index:mobile_index;comment:企业微信手机号"`
	Position       string `json:"position" gorm:"column:position;type:varchar(255);not null;comment:职务"`
	Gender         int32  `json:"gender" gorm:"column:gender;type:tinyint(2);not null;default:0;comment:性别0=未知 1=男性 2=女性"`
	Email          string `json:"email" gorm:"column:email;type:varchar(128);not null;index:email_index;comment:邮箱"`
	Avatar         string `json:"avatar" gorm:"column:avatar;type:varchar(512);not null;comment:员工头像"`
	ThumbAvatar    string `json:"thumb_avatar" gorm:"column:thumb_avatar;type:varchar(128);not null;comment:缩略图"`
	Telephone      string `json:"telephone" gorm:"column:telephone;type:varchar(32);not null;comment:座机"`
	Alias          string `json:"alias" gorm:"column:alias;type:varchar(64);not null;comment:别名"`
	Status         int32  `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业"`
	QrCode         string `json:"qr_code" gorm:"column:qr_code;type:varchar(128);not null;comment:员工个人二维码"`
	Address        string `json:"address" gorm:"column:address;type:varchar(255);not null;comment:地址"`
	OpenUserid     string `json:"open_userid" gorm:"column:open_userid;type:varchar(64);not null;default:open_userid"`
	Type           uint32 `json:"type" gorm:"column:type;type:tinyint(3) unsigned;not null;default:3;comment:导购类型：1、门店导购，2、社群导购，3、企业导购，4、店长或企业管理"`
	QrPosterCodeID int32  `json:"qr_poster_code_id" gorm:"column:qr_poster_code_id;type:int(11);not null;default:0;comment:员工海报id"`
	QrPosterCode   string `json:"qr_poster_code" gorm:"column:qr_poster_code;type:varchar(200);not null;comment:员工海报"`

	MemberDepartments []MemberDepartmentsModel `json:"member_departments" gorm:"foreignKey:MemberID;references:ID"`
}

// TableName 企业微信员工
func (MemberModel) TableName() string {
	return "members"
}

type MemberDepartment struct {
	EnterpriseId      int64  `json:"enterprise_id"`
	MembersId         int64  `json:"members_id"`
	MembersName       string `json:"members_name"`
	MembersUserid     string `json:"members_userid"`
	DepartmentId      int64  `json:"department_id"`
	DepartmentThirdId int64  `json:"department_third_id"`
	DepartmentName    string `json:"department_name"`
	Mobile            string `json:"mobile"`
	Status            int    `json:"status"`
}

type MemberInfo struct {
	Id             int64  `json:"id"`
	EnterpriseId   int64  `json:"enterprise_id"`
	Name           string `json:"name"`
	Userid         string `json:"userid"`
	Mobile         string `json:"mobile"`
	Type           uint32 `json:"type"`
	Avatar         string `json:"avatar"`
	QrPosterCode   string `json:"qr_poster_code"`
	Status         int32  `json:"status"`
	DepartmentName string `json:"department_name"`
	ThirdId        int    `json:"third_id"`
	StoreId        int    `json:"store_id"`
}

type StoreManager struct {
	MemberId     int64  `json:"member_id"`
	ThirdId      int    `json:"third_id"`
	MemberName   string `json:"member_name"`
	StoreCode    string `json:"store_code"`
	Position     string `json:"position"`
	PositionAbbr string `json:"position_abbr"`
}

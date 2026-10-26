package micro_scrm_enterprise

import (
	"fmt"
	"micro-scrm/internal/orm/model"
)

type MembersModel struct {
	model.BranchModel

	Userid           string `json:"userid" gorm:"column:userid;type:varchar(64);not null;index:user_id_index;comment:userid"`
	Name             string `json:"name" gorm:"column:name;type:varchar(128);not null;comment:企业微信姓名"`
	Department       string `json:"department" gorm:"column:department;type:varchar(64);not null;comment:department"`
	Order            string `json:"order" gorm:"column:order;type:varchar(64);not null;comment:order"`
	Position         string `json:"position" gorm:"column:position;type:varchar(64);not null;comment:职务"`
	Mobile           string `json:"mobile" gorm:"column:mobile;type:varchar(15);not null;index:mobile_index;comment:企业微信手机号"`
	Gender           int    `json:"gender" gorm:"column:gender;type:tinyint(2);not null;default:0;comment:性别0=未知 1=男性 2=女性"`
	Email            string `json:"email" gorm:"column:email;type:varchar(128);not null;index:email_index;comment:邮箱"`
	IsLeaderInDept   string `json:"is_leader_in_dept" gorm:"column:is_leader_in_dept;type:varchar(64);not null;comment:is_leader_in_dept"`
	Avatar           string `json:"avatar" gorm:"column:avatar;type:varchar(512);not null;comment:员工头像"`
	ThumbAvatar      string `json:"thumb_avatar" gorm:"column:thumb_avatar;type:varchar(128);not null;comment:缩略图"`
	Telephone        string `json:"telephone" gorm:"column:telephone;type:varchar(32);not null;comment:座机"`
	Alias            string `json:"alias" gorm:"column:alias;type:varchar(64);not null;comment:别名"`
	Status           int8   `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业"`
	Address          string `json:"address" gorm:"column:address;type:varchar(255);not null;comment:地址"`
	HideMobile       string `json:"hide_mobile" gorm:"column:hide_mobile;type:varchar(32);not null;comment:座机"`
	EnglishName      string `json:"english_name" gorm:"column:english_name;type:varchar(64);not null;comment:english_name"`
	OpenUserid       string `json:"open_userid" gorm:"column:open_userid;type:varchar(64);not null;default:open_userid"`
	MainDepartment   int32  `json:"main_department" gorm:"column:main_department;type:int(10);not null;comment:main_department"`
	Extattr          string `json:"extattr" gorm:"column:extattr;type:text;not null;comment:extattr"`
	QrCode           string `json:"qr_code" gorm:"column:qr_code;type:varchar(128);not null;comment:员工个人二维码"`
	ExternalPosition string `json:"external_position" gorm:"column:external_position;type:varchar(64);not null;comment:external_position"`
	ExternalProfile  string `json:"external_profile" gorm:"column:external_profile;type:text;not null;comment:main_department"`
	SaleArea         string `json:"sale_area" gorm:"column:sale_area;type:varchar(50);not null;comment:柜台名称"`
	City             string `json:"city" gorm:"column:city;type:varchar(50);not null;comment:城市"`
	BrandCode        string `json:"brand_code" gorm:"column:brand_code;type:varchar(50);not null;comment:品牌代码"`
	CounterName      string `json:"counter_name" gorm:"column:counter_name;type:varchar(50);not null;comment:柜台名称"`
	MarsCode         string `json:"mars_code" gorm:"column:mars_code;type:varchar(50);not null;comment:mars_code"`
	MarsID           string `json:"mars_id" gorm:"column:mars_id;type:varchar(50);not null;comment:mars_id"`
}

// TableName returns the table name of the MembersModel model
func (MembersModel) TableName() string {
	return "members"
}

func (m MembersModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", m.DBPrefix(), enterpriseID, m.TableName())
}

type MemberInfo struct {
	MenberId int64  `json:"menber_id"`
	User_id  string `json:"user_id"`
}

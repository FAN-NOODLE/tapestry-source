package microScrm

// AuthInfoRequest 用户授权请求
type AuthInfoRequest struct {
	Code         string `json:"code" desc:"企业授权code"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业id"`
	LoginType    string `json:"login_type" desc:"登录方式"`
}

// AuthInfoResponse 企业微信授权信息
type AuthInfoResponse struct {
	Member   AuthInfoMember   `json:"member"`
	External AuthInfoExternal `json:"external"`
}

// AuthInfoMember 内部成员信息
type AuthInfoMember struct {
	UserId   string `json:"user_id" desc:"成员ID"`
	DeviceId string `json:"device_id" desc:"手机设备号"`
	Name     string `json:"name"`
	MemberId int64  `json:"member_id"`
}

// AuthInfoExternal 外部成员信息
type AuthInfoExternal struct {
	OpenId         string `json:"openid" desc:"非企业成员的标识"`
	DeviceId       string `json:"device_id" desc:"手机设备号"`
	ExternalUserid string `json:"external_userid" desc:"外部联系人id"`
}

type DetailUseridRequest struct {
	Userid       string `json:"userid" desc:"企业微信userid"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业id"`
}

// Member 成员
type Member struct {
	ID           int64  `json:"id" db:"id" desc:"员工ID"`
	Userid       string `json:"userid" db:"userid" desc:"企业微信userid"`
	Name         string `json:"name" db:"name" desc:"员工姓名"`
	Mobile       string `json:"mobile" db:"mobile" desc:"员工手机号"`
	Position     string `json:"position" db:"position" desc:"职位"`
	Gender       int    `json:"gender" db:"gender" desc:"性别"`
	Email        string `json:"email" db:"email" desc:"员工邮箱"`
	Avatar       string `json:"avatar" db:"avatar" desc:"头像"`
	ThumbAvatar  string `json:"thumb_avatar" db:"thumb_avatar" desc:""`
	Telephone    string `json:"telephone" db:"telephone" desc:"电话"`
	Alias        string `json:"alias" db:"alias" desc:"别名"`
	Status       int8   `json:"status" db:"status" desc:"状态,1=已激活，2=已禁用，4=未激活，5=退出企业"`
	QrCode       string `json:"qr_code" db:"qr_code" desc:"二维码"`
	Address      string `json:"address" db:"address" desc:"地址"`
	OpenUserid   string `json:"open_userid" db:"open_userid" desc:""`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	Forms        string `json:"forms" desc:"所属部门"`
	Type         int    `json:"type" db:"type" desc:"3、企业导购，4、店长或企业管理"`
	// 非member 成员字段
	IsLeaderInDept []int8  `json:"is_leader_in_dept" desc:"表示在所在的部门内是否为上级"`
	Department     []int64 `json:"department" desc:"成员所属部门id列表，仅返回该应用有查看权限的部门id"`
	GenderText     string  `json:"gender_text" desc:"性别"`
	StatusText     string  `json:"status_text" desc:"状态"`
}

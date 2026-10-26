package openResp

type StoreList struct {
	StoreName          string `json:"store_name" desc:"门店名称"`
	StoreCode          string `json:"store_code" desc:"门店编码"`
	InventoryName      string `json:"inventory_name"  desc:"英文名（库存相关名）"`
	Address            string `json:"address"  desc:"地址"`
	City               string `json:"city"  desc:"城市"`
	Phone              string `json:"phone"  desc:"电话"`
	Email              string `json:"email"  desc:"邮箱"`
	StoreType          int32  `json:"store_type"  desc:"门店类型 0未设置 1奥莱 2零售"`
	ClosedAt           int64  `json:"closed_at"`
	WechatDepartmentID uint   `json:"wechat_department_id"  desc:"企微部门id"`
}

type SaList struct {
	SaCode    string `json:"sa_code" desc:"导购编码"`
	StoreCode string `json:"store_code"`
	Name      string `json:"name" desc:"姓名"`
}

type SaListByCustomerList struct {
	Name      string `json:"name" desc:"姓名"`
	SaCode    string `json:"sa_code" desc:"导购编码"`
	CreatedAt int64  `json:"created_at" desc:"好友关系添加时间"`
}

type StoreWorkListInfo struct {
	MemberName int64  `json:"member_id"`
	Member     string `json:"member" desc:""`
}

type StoreWorkList struct {
	MembersUserid  string                     `json:"members_userid"`
	MembersName    string                     `json:"members_name"`
	Mobile         string                     `json:"mobile"`
	Position       string                     `json:"position"`
	IsStoreManager int                        `json:"is_store_manager"`
	QrCode         string                     `json:"qr_code"`
	Status         int32                      `json:"status"`
	Departments    []StoreWorkListDepartments `json:"departments"`
}

type StoreWorkListDepartments struct {
	Name      string `json:"name"`
	ThirdId   int    `json:"third_id"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
}

type CustomerStoreWorkListResp struct {
	Id           int64  `json:"id"`
	MemberUserId string `json:"user_id"`
	MemberName   string `json:"member_name"`
	Remark       string `json:"remark"`
	CreateTime   int    `json:"create_time"`
	// IsKeySa      int                        `json:"is_key_sa"`
	Departments []StoreWorkListDepartments `json:"departments"`
}

type GenQrCodeResp struct {
	StoreWorkList
}

type AreaList struct {
	ID          int       `json:"id"`
	LMSID       string    `json:"lms_id"`
	Type        int       `json:"type" `
	CustomGroup string    `json:"custom_group"`
	Name        string    `json:"name"`
	EnName      string    `json:"en_name"`
	ParentID    int       `json:"parent_id"`
	ThirdID     int       `json:"third_id"`
	ManagerList []Manager `json:"manager_list"`
	ClosedAt    int64     `json:"closed_at"`
	StoreCode   string    `json:"store_code"`
	UpdatedAt   int64     `json:"updated_at"`
	//获取子区域
	Children []AreaList `json:"children"`
}

type Manager struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Position string `json:"position"`
}

type ManagerInfo struct {
	UserId    string      `json:"user_id"`
	UserName  string      `json:"user_name"`
	Position  string      `json:"position"`
	StoreList []StoreInfo `json:"store_list"`
}

type StoreInfo struct {
	ThirdId   int    `json:"third_id"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
	EnName    string `json:"en_name"`
}

type HkStoreList struct {
	StoreName          string `json:"store_name" desc:"门店名称"`
	StoreCode          string `json:"store_code" desc:"门店编码"`
	Area               string `json:"area" desc:"hk/mo"`
	InventoryName      string `json:"inventory_name" desc:"英文名（库存相关名）"`
	Address            string `json:"address" desc:"地址"`
	City               string `json:"city" desc:"城市"`
	Phone              string `json:"phone" desc:"电话"`
	Email              string `json:"email" desc:"邮箱"`
	ClosedAt           int64  `json:"closed_at"`
	WechatDepartmentID uint   `json:"wechat_department_id" desc:"企微部门id"`
}

type HkStoreWorkList struct {
	MembersUserid  string                      `json:"members_userid"`
	MembersName    string                      `json:"members_name"`
	Mobile         string                      `json:"mobile"`
	Position       string                      `json:"position"`
	IsStoreManager int                         `json:"is_store_manager"`
	QrCode         string                      `json:"qr_code"`
	Status         int32                       `json:"status"`
	Departments    []HkStoreWorkListDepartment `json:"departments"`
}

type HkStoreWorkListDepartment struct {
	Name      string `json:"name"`
	ThirdId   int    `json:"third_id"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
}

type HkManagerInfo struct {
	UserId    string        `json:"user_id"`
	UserName  string        `json:"user_name"`
	Position  string        `json:"position"`
	StoreList []HkStoreInfo `json:"store_list"`
}

type HkStoreInfo struct {
	ThirdId   int    `json:"third_id"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
	Area      string `json:"area"`
	EnName    string `json:"en_name"`
}

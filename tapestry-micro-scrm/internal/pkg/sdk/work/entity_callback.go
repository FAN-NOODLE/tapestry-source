package work

// ContactCallback 企微通讯录回调基础信息
type ContactCallback struct {
	ToUserName   string `json:"toUserName"`
	FromUserName string `json:"fromUserName"`
	CreateTime   int    `json:"createTime"`
	MsgType      string `json:"msgType"`
	Event        string `json:"event"`
	ChangeType   string `json:"changeType"`
}

// ContactCallbackUser 企微通讯录回调-员工信息
// 文档地址 https://developer.work.weixin.qq.com/document/path/90970
type ContactCallbackUser struct {
	ContactCallback
	UserID         string `json:"userID"`
	NewUserID      string `json:"newUserID"`
	Name           string `json:"name"`
	Department     string `json:"department"`
	MainDepartment string `json:"mainDepartment"`
	IsLeaderInDept int    `json:"isLeaderInDept"`
	Position       string `json:"position"`
	Mobile         string `json:"mobile"`
	Gender         int32  `json:"gender"`
	Email          string `json:"email"`
	Status         string `json:"status"` // 激活状态：1=已激活 2=已禁用 4=未激活 已激活代表已激活企业微信或已关注微信插件（原企业号）5=成员退出
	Avatar         string `json:"avatar"`
	Alias          string `json:"alias"`
	Telephone      string `json:"telephone"`
	Address        string `json:"address"`
}

// ContactCallbackDepartment 企微通讯录回调-部门信息
// 文档地址 https://developer.work.weixin.qq.com/document/path/90971
type ContactCallbackDepartment struct {
	ContactCallback
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Order    float64 `json:"order"`
	ParentId string  `json:"parentId"`
}

// ContactCallbackTag 企微通讯录回调-标签信息
// 文档地址 https://developer.work.weixin.qq.com/document/path/90972
type ContactCallbackTag struct {
	ContactCallback
	TagId         int    `json:"tagId"`
	AddUserItems  string `json:"addUserItems"`
	DelUserItems  string `json:"delUserItems"`
	AddPartyItems string `json:"addPartyItems"`
	DelPartyItems string `json:"delPartyItems"`
}

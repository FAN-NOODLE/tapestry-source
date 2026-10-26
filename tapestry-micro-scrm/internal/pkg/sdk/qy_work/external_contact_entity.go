package qyWork

// SendWelcomeReq 发送欢迎语
type SendWelcomeReq struct {
	WelcomeCode string        `json:"welcome_code"`
	Text        Text          `json:"text"`
	Attachments []Attachments `json:"attachments"`
}

// ContactCreateRequest 添加联系我二维码
type ContactCreateRequest struct {
	Remark     string   `json:"remark"  desc:"联系方式的备注信息，用于助记，不超过30个字符"`
	SkipVerify bool     `json:"skip_verify"  desc:"外部客户添加时是否无需验证，默认为true"`
	Type       int      `json:"type"  desc:"联系方式类型，1-单人，2-多人"`
	Scene      int      `json:"scene" desc:"场景，1-在小程序中联系，2-通过二维码联系"`
	User       []string `json:"user" desc:"使用该联系方式的用户userID列表，在type为1时为必填，且只能有一个"`
	Party      []int    `json:"party"`
	State      string   `json:"state,omitempty" desc:"用于区分不同的添加渠道"`
}

// ContactUpdateRequest   更新联系我二维码
type ContactUpdateRequest struct {
	ConfigId   string   `json:"config_id"`
	SkipVerify bool     `json:"skip_verify"  desc:"外部客户添加时是否无需验证，默认为true"`
	User       []string `json:"user" desc:"使用该联系方式的用户userID列表，在type为1时为必填，且只能有一个"`
	Party      []int    `json:"party"`
	State      string   `json:"state,omitempty" desc:"用于区分不同的添加渠道"`
}

// ContactDeleteRequest   删除联系我二维码
type ContactDeleteRequest struct {
	ConfigId string `json:"config_id"`
}

// ContactListRequest    查询联系我二维码
type ContactListRequest struct {
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Limit     int    `json:"limit"`
	Cursor    string `json:"cursor"`
}

type ContactCreateResp struct {
	QyResponse
	ConfigId string `json:"config_id,omitempty"`
	QrCode   string `json:"qr_code,omitempty"`
}

type ContactList struct {
	ConfigId string `json:"config_id"`
}

type ContactListResp struct {
	QyResponse
	ContactWay []ContactList `json:"contact_way"`
	NextCursor string        `json:"next_cursor"`
}

type ContactGetResp struct {
	QyResponse
	ContactWay ContactUpdateRequest `json:"contact_way"`
}

// GetGroupMsgResultResp 群发消息发送结果
type GetGroupMsgResultResp struct {
	ExternalUserID string `json:"external_userid"`
	ChatID         string `json:"chat_id"`
	UserID         string `json:"userid"`
	Status         int64  `json:"status"`
	SendTime       int64  `json:"send_time"`
}

// AddMsgTemplateReq 添加企业群发消息
type AddMsgTemplateReq struct {
	ChatType       string        `json:"chat_type"`
	ExternalUserid []string      `json:"external_userid"`
	Sender         string        `json:"sender"`
	AllowSelect    bool          `json:"allow_select"`
	Text           Text          `json:"text"`
	Attachments    []Attachments `json:"attachments"`
}

type Text struct {
	Content string `json:"content"`
}

type Attachments struct {
	MsgType     string      `json:"msgtype"`
	Image       Image       `json:"image,omitempty"`
	Link        Link        `json:"link,omitempty"`
	MiniProgram MiniProgram `json:"miniprogram,omitempty"`
	Video       Video       `json:"video,omitempty"`
	File        File        `json:"file,omitempty"`
}

type Image struct {
	MediaId string `json:"media_id"`
	PicUrl  string `json:"pic_url"`
}

type Link struct {
	Title  string `json:"title"`
	Picurl string `json:"picurl"`
	Desc   string `json:"desc"`
	Url    string `json:"url"`
}

type MiniProgram struct {
	Title      string `json:"title"`
	PicMediaId string `json:"pic_media_id"`
	Appid      string `json:"appid"`
	Page       string `json:"page"`
}

type Video struct {
	MediaId string `json:"media_id"`
}

type File struct {
	MediaId string `json:"media_id"`
}

type AddMsgTemplateResp struct {
	QyResponse
	FailList []string `json:"fail_list"`
	Msgid    string   `json:"msgid"`
}

// MarkTagReq 编辑客户企业标签
type MarkTagReq struct {
	Userid         string   `json:"userid"`
	ExternalUserid string   `json:"external_userid"`
	AddTag         []string `json:"add_tag"`
	RemoveTag      []string `json:"remove_tag"`
}

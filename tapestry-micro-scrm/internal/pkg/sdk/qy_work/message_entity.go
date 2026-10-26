package qyWork

// SendMessageRequest 发送应用消息的基本请求体
type SendMessageRequest struct {
	EnterpriseID int64               `json:"enterprise_id"`
	MemUserIds   []string            `json:"mem_user_ids"`
	PartyIds     []string            `json:"partyIds"`
	Title        string              `json:"title"`
	Content      []map[string]string `json:"content"`
	Page         string              `json:"page"`
}

type SendTemplateCardMessageRequest struct {
	SendMessageRequest
	CardType   string     `json:"card_type"`
	JumpList   []JumpItem `json:"jump_list"`
	CardAction CardAction `json:"card_action"`
}

type SendMessageResp struct {
	QyResponse
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	MsgId        string `json:"msgid"`
}

type SendMessageBase struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
}

// SendMiniMessageReq 文档地址 https://developer.work.weixin.qq.com/document/path/90236#%E5%B0%8F%E7%A8%8B%E5%BA%8F%E9%80%9A%E7%9F%A5%E6%B6%88%E6%81%AF
type SendMiniMessageReq struct {
	SendMessageBase
	MiniProgramNotice      MiniProgramNotice `json:"miniprogram_notice"`
	EnableIdTrans          int               `json:"enable_id_trans"`
	EnableDuplicateCheck   int               `json:"enable_duplicate_check"`
	DuplicateCheckInterval int               `json:"duplicate_check_interval"`
}
type MiniProgramNotice struct {
	Appid             string        `json:"appid"`
	Page              string        `json:"page"`
	Title             string        `json:"title"`
	Description       string        `json:"description,omitempty"`
	EmphasisFirstItem bool          `json:"emphasis_first_item"`
	ContentItem       []ContentItem `json:"content_item"`
}

type ContentItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SendTemplateCardMessageReq 文档地址 https://developer.work.weixin.qq.com/document/path/90236#%E6%A8%A1%E6%9D%BF%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
type SendTemplateCardMessageReq struct {
	SendMessageBase
	AgentId                int          `json:"agentid"`
	TemplateCard           TemplateCard `json:"template_card"`
	EnableIdTrans          int          `json:"enable_id_trans"`
	EnableDuplicateCheck   int          `json:"enable_duplicate_check"`
	DuplicateCheckInterval int          `json:"duplicate_check_interval"`
}

type TemplateCard struct {
	CardType              string               `json:"card_type,omitempty"`
	Source                *Source              `json:"source,omitempty"`
	ActionMenu            *ActionMenu          `json:"action_menu,omitempty"`
	TaskID                string               `json:"task_id,omitempty"`
	MainTitle             *MainTitle           `json:"main_title,omitempty"`
	QuoteArea             *QuoteArea           `json:"quote_area,omitempty"`
	EmphasisContent       *EmphasisContent     `json:"emphasis_content,omitempty"`
	SubTitleText          string               `json:"sub_title_text,omitempty"`
	HorizontalContentList *[]HorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              *[]JumpItem          `json:"jump_list,omitempty"`
	CardAction            *CardAction          `json:"card_action,omitempty"`
}

type Source struct {
	IconURL   string `json:"icon_url"`
	Desc      string `json:"desc"`
	DescColor int    `json:"desc_color"`
}

type ActionMenu struct {
	Desc       string       `json:"desc"`
	ActionList []ActionItem `json:"action_list"`
}

type ActionItem struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type MainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type QuoteArea struct {
	Type      int    `json:"type"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	QuoteText string `json:"quote_text"`
}

type EmphasisContent struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type HorizontalContent struct {
	KeyName string `json:"keyname"`
	Value   string `json:"value"`
	Type    int    `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	MediaID string `json:"media_id,omitempty"`
	UserID  string `json:"userid,omitempty"`
}

type JumpItem struct {
	Type     int    `json:"type"`
	Title    string `json:"title"`
	URL      string `json:"url,omitempty"`
	AppID    string `json:"appid,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
}

type CardAction struct {
	Type     int    `json:"type"`
	URL      string `json:"url,omitempty"`
	AppID    string `json:"appid,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
}

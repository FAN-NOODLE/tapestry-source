package microScrm

// SourceListRequest 素材列表请求参数
type SourceListRequest struct {
	Page           int     `json:"page"`
	Size           int     `json:"size"`
	Keyword        string  `json:"keyword"`
	Type           string  `json:"type"  validate:"required" required:"素材类型必填（file、image、text、video、link、mini-program）" desc:"素材类型 file、image、text、video、link、mini-program"`
	GroupID        int64   `json:"group_id"`
	MemberID       int64   `json:"member_id" db:"member_id" desc:"ws员工Id"`
	MemberUserId   string  `json:"member_user_id" db:"member_user_id" desc:"企业员工Id"`
	GroupChatID    string  `json:"group_chat_id" db:"group_chat_id" desc:"ws群id"`
	ExternalUserid string  `json:"external_userid" db:"external_userid" desc:"ws外部联系人openid"`
	CategoryIds    []int64 `json:"category_ids"`
	RequestBase
}

// SourceGroupListRequest 素材分组列表
type SourceGroupListRequest struct {
	RequestBase
}

// SourceMediaRequest 生成临时素材
type SourceMediaRequest struct {
	SourceID    int64  `json:"source_id"`
	SourceType  string `json:"source_type"`
	SourceUrl   string `json:"source_url"`
	SourceName  string `json:"source_name"`
	ResourceKey string `json:"resource_key"`
	RequestBase
}

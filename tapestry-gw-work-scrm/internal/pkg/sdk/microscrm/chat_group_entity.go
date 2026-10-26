package microScrm

type GroupChatStaticsRequest struct {
	ChatID   string `json:"chat_id" db:"chat_id" desc:"群ID" validate:"required" required:"群ID不能为空"`
	STS      int64  `json:"sts" db:"sts" desc:"开始时间"`
	ETS      int64  `json:"ets" db:"ets" desc:"开始时间"`
	Friend   int    `json:"friend" db:"friend" desc:"0全部 1好友 2非好友"`
	MemberID int64  `json:"member_id"`
	UserID   string `json:"user_id"`
	ListType int    `json:"list_type"`
	RequestBase
	ListBase
}

type ActivityInfo struct {
	ChatID string `json:"chat_id" db:"chat_id" validate:"required" required:"群ID不能为空"`
	UserID string `json:"user_id" db:"user_id" validate:"required" required:"用户ID不能为空"`
	RequestBase
}

type GroupChatBaseResponse struct {
	ResponseBase
	Data interface{} `json:"data"`
}

package groupmsg

type TaskMemberSendListResp struct {
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	List  []TaskMemberSendList `json:"list"`
	Total int64                `json:"total"`
}

type TaskMemberSendList struct {
	ChatId         string `json:"chat_id"`
	ChatName       string `json:"chat_name"`
	UserId         string `json:"user_id"`
	UserName       string `json:"user_name"`
	ExternalUserID string `json:"external_userid"`
	Status         int    `json:"status"`
	SendTime       int64  `json:"send_time"`
}

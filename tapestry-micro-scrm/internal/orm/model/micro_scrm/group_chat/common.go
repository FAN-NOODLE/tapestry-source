package group_chat

type GroupChatPoeple struct {
	ID            int64  `json:"id"`
	OwnerMemberID string `json:"owner_member_id"`
	UserId        string `json:"user_id"`
	Number        int64  `json:"number"`
	ChatId        string `json:"chat_id"`
}

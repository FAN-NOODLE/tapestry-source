package dataCompassReq

type ChatGroupReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	PageNum      int   `json:"page_num"`
	PageSize     int   `json:"page_size"`
	StartTime    int64 `json:"start_time"`
	EndTime      int64 `json:"end_time"`
}

// ChatGroupResp 社群数据明细
type ChatGroupResp struct {
	GroupChatId               int64  `json:"group_chat_id"`                // 群编号
	GroupName                 string `json:"group_name"`                   // 群组名称
	CumulativeGroupNumber     int64  `json:"cumulative_group_number"`      // 累计群成员数
	NumberOfNewcomers         int64  `json:"number_of_newcomers"`          // 新入群人数
	LeaveGroupNumber          int64  `json:"leave_group_number"`           // 退群人数
	LeaveGroupPercentage      string `json:"leave_group_percentage"`       // 退群率
	AddFriendNumber           int64  `json:"add_friend_number"`            // 加好友人数
	AddFriendPercentage       string `json:"add_friend_percentage"`        // 加好友率
	WxMemberBindingNumber     int64  `json:"wx_member_binding_number"`     // 微信会员绑定人数
	WxMemberBindingPercentage string `json:"wx_member_binding_percentage"` // 微信会员绑定率
	ParticipateNumber         int64  `json:"participate_number"`           // 群活动参与人数
	ParticipatePercentage     string `json:"participate_percentage"`       // 群活动参与率
	IntoGroupNumber           int64  `json:"into_group_number"`            // 活动进群人数
	IntoGroupPercentage       string `json:"into_group_percentage"`        // 活动进群率
	BANumber                  int64  `json:"ba_number"`                    // 导购人数
}

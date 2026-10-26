package dataCompassResp

// ChatGroupResp 社群数据明细
type ChatGroupResp struct {
	GroupChatId           int64  `json:"group_chat_id"`           // 群编号
	GroupName             string `json:"group_name"`              // 群组名称
	CumulativeGroupNumber int64  `json:"cumulative_group_number"` // 累计群成员数
	NumberOfNewcomers     int64  `json:"number_of_newcomers"`     // 新入群人数
	LeaveGroupNumber      int64  `json:"leave_group_number"`      // 退群人数
	LeaveGroupPercentage  string `json:"leave_group_percentage"`  // 退群率
	BANumber              int64  `json:"ba_number"`               // 导购人数
}

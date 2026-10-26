package customerReq

type BaseReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
}

type CompassReq struct {
	BaseReq
	Mobile        string  `json:"mobile"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
	MemberIds     []int64 `json:"member_ids"`
	DepartmentIds []int64 `json:"department_ids"`
	SearchAll     int     `json:"search_all" comment:"是否是查询全部,0为否,1为是"`
	IsSM          int     `json:"is_sm" comment:"判断是否为店长身份,0为否,1为是"`
	ThirdIds      []int   `json:"third_ids" comment:"店长所管理的门店"`
}

type ChatGroupCompassRes struct {
	CumulativeGroupNumber int64  `json:"cumulative_group_number" comment:"累计群成员数"`
	GroupNumber           int64  `json:"group_number" comment:"群总数"`
	AddFriendNumber       int64  `json:"add_friend_number" comment:"群成员添加好友数"`
	AddFriendPer          string `json:"add_friend_per" comment:"群成员添加好友率"`
	LeaveGroupNumber      int64  `json:"leave_group_number" comment:"退群人数"`
	LeaveGroupPer         string `json:"leave_group_per" comment:"退群率"`
	JoinGroupNumber       int64  `json:"join_group_number" comment:"入群人数"`
}

type FriendsCompassRes struct {
	AddFriendNumber   int64  `json:"add_friend_number" comment:"新增好友数"`
	AllFriendsNumber  int64  `json:"all_friends_number" comment:"累计好友数"`
	LossFriendsNumber int64  `json:"loss_friends_number" comment:"好友流失人数"`
	JoinGroupNumber   int64  `json:"join_group_number" comment:"入群人数"`
	NewJoinGroupPer   string `json:"new_join_group_per" comment:"新增好友入群率"`
	JoinGroupPer      string `json:"join_group_per" comment:"好友入群率"`
	LossFriendsPer    string `json:"loss_friends_per" comment:"好友流失率"`
	SaNumber          int    `json:"sa_number"`
}

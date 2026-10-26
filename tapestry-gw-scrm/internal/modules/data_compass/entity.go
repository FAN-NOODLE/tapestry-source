package data_compass

/**
*数据罗盘
 */

type CommunityDataDetailExportRequest struct {
	EnterpriseID  int64    `json:"enterprise_id"`
	ShoppingGuide []string `json:"shopping_guide"`
	UserIds       []string `json:"user_ids"`
	StartTime     int64    `json:"start_time"`
	EndTime       int64    `json:"end_time"`
}
type CommunityDataDetail struct {
	GroupChatId               int64  `json:"group_chat_id"` //群编号
	GroupName                 string `json:"group_name"`
	CumulativeGroupNumber     int64  `json:"cumulative_group_number"`      //累计群成员数
	NumberOfNewcomers         int64  `json:"number_of_newcomers"`          //新入群人数
	LeaveGroupNumber          int64  `json:"leave_group_number"`           //退群人数
	LeaveGroupPercentage      string `json:"leave_group_percentage"`       //退群率
	AddFriendNumber           int64  `json:"add_friend_number"`            // 加好友人数
	AddFriendPercentage       string `json:"add_friend_percentage"`        //加好友率
	WxMemberBindingNumber     int64  `json:"wx_member_binding_number"`     //微信会员绑定人数
	WxMemberBindingPercentage string `json:"wx_member_binding_percentage"` //微信会员绑定率
	ParticipateNumber         int64  `json:"participate_number"`           //群活动参与人数
	ParticipatePercentage     string `json:"participate_percentage"`       // 群活动参与率
	IntoGroupNumber           int64  `json:"into_group_number"`            // 活动进群人数
	IntoGroupPercentage       string `json:"into_group_percentage"`        // 活动进群率
	BANumber                  int64  `json:"ba_number"`                    // ba 导购人数
}

type CommunityDate struct {
	List []CommunityDataDetail `json:"list"`
}
type BuddyData struct {
	List []BuddyDataDetail `json:"list"`
}
type BuddyDataDetail struct {
	ShoppingGuideID                 string `json:"shopping_guide_id"`                   //导购id
	ShoppingGuideName               string `json:"shopping_guide_name"`                 //导购
	BelongingArea                   string `json:"belonging_area"`                      //所属区域
	BelongingStore                  string `json:"belonging_store"`                     // 所属门店
	NewAddBuddyNumber               int64  `json:"new_add_buddy_number"`                //新增好友数
	NewBuddyMemberBindingNumber     int64  `json:"new_buddy_member_binding_number"`     //新增好友会员绑定数
	NewBuddyMemberBindingPercentage string `json:"new_buddy_member_binding_percentage"` //新增好友会员绑定率
	NewBuddyIntoNumber              int64  `json:"new_buddy_into_number"`               //新增好友入群数
	NewBuddyIntoPercentage          string `json:"new_buddy_into_percentage"`           //新增好友入群率
	BuddyChurnNumber                int64  `json:"buddy_churn_number"`                  //好友流失人数
	BuddyChurnPercentage            string `json:"buddy_churn_percentage"`              //好友流失率
	GrandTotalNumber                int64  `json:"grand_total_number"`                  //累计好友数
	BuddyIntoPercenuNmber           int64  `json:"buddy_into_percen_number"`            //好友入群数
	BuddyIntoPercentage             string `json:"buddy_into_percentage"`
}

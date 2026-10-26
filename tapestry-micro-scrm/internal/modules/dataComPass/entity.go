package dataComPass

// 接收参数实体/返回模型实体

// DataRequest 通用查询参数
type DataRequest struct {
	SearchAll     int64    `json:"search_all" comment:"1为查询全部导购,0为选择导购"`
	EnterpriseID  int64    `json:"enterprise_id"`
	ShoppingGuide []string `json:"shopping_guide"`
	UserIds       []string `json:"user_ids"`
	StartTime     int64    `json:"start_time"`
	EndTime       int64    `json:"end_time"`
	Page
	Types
}

// Page 分页参数
type Page struct {
	PageSize int `json:"page_size"`
	PageNum  int `json:"page_num"`
}

// Types Type 类型
type Types struct {
	Type int `json:"type"`
}

// CommunityData 社群数据
type CommunityData struct {
	GroupNumber           string `json:"group_number"`            //群数量
	CumulativeGroupNumber string `json:"cumulative_group_number"` //累计群成员数
	NumberOfNewcomers     string `json:"number_of_newcomers"`     //新入群人数
	LeaveGroupNumber      string `json:"leave_group_number"`      //退群人数
	LeaveGroupPercentage  string `json:"leave_group_percentage"`  //退群率
	AddFriendNumber       string `json:"add_friend_number"`       //加好友人数
	AddFriendPercentage   string `json:"add_friend_percentage"`   //加好友率
	//WxMemberBindingNumber     string `json:"wx_member_binding_number"`     //微信会员绑定人数
	//WxMemberBindingPercentage string `json:"wx_member_binding_percentage"` //微信会员绑定率
}

// CommunityActivityData 社群活动数据
type CommunityActivityData struct {
	ActivityGroupNumber   string `json:"activity_group_number"`  //活动群数
	ParticipateNumber     string `json:"participate_number"`     //群活动参与人数
	ParticipatePercentage string `json:"participate_percentage"` // 群活动参与率
	IntoGroupNumber       string `json:"into_group_number"`      // 活动进群人数
	IntoGroupPercentage   string `json:"into_group_percentage"`  // 活动进群率
}

// BuddyData 好友数据
type BuddyData struct {
	NewAddBuddyNumber    string `json:"new_add_buddy_number"`   //新增好友数
	GrandTotalNumber     string `json:"grand_total_number"`     //累计好友数
	BuddyIntoNumber      string `json:"buddy_into_number"`      //好友入群数
	BuddyIntoPercentage  string `json:"buddy_into_percentage"`  //好友入群率
	BuddyChurnNumber     string `json:"buddy_churn_number"`     //好友流失人数
	BuddyChurnPercentage string `json:"buddy_churn_percentage"` //好友流失率
	//NewBuddyMemberBindingNumber string `json:"new_buddy_member_binding_number"` //新增好友会员绑定率
	NewBuddyIntoNumber string `json:"new_buddy_into_number"` //新增好友入群率
	BANumber           string `json:"ba_number"`             // ba 导购人数
	//WxMemberBindingNumber       string `json:"wx_member_binding_number"`        //微信会员绑定人数
	//WxMemberBindingPercentage   string `json:"wx_member_binding_percentage"`    //微信会员绑定率
}

// CommunityDataDetail 社群数据明细
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

// CommunityDataDetailList 社群数据明细列表
type CommunityDataDetailList struct {
	Count               int                    `json:"count"`
	CommunityDataDetail []*CommunityDataDetail `json:"community_data_detail"`
}

// BuddyDataDetail 好友数据明细
type BuddyDataDetail struct {
	ShoppingGuideID                 string `json:"shopping_guide_id"`                   //导购id
	ShoppingGuideName               string `json:"shopping_guide_name"`                 //导购
	ShoppingUserId                  string `json:"shopping_user_id"`                    //导购
	ShoppingMemberId                int64  `json:"shopping_member_id"`                  //导购的memberId
	BelongingArea                   string `json:"belonging_area"`                      //所属区域
	Avatar                          string `json:"avatar"`                              // ba图像
	BelongingStore                  string `json:"belonging_store"`                     // 所属门店
	NewAddBuddyNumber               int64  `json:"new_add_buddy_number"`                //新增好友数
	NewBuddyMemberBindingNumber     int64  `json:"new_buddy_member_binding_number"`     //新增好友会员绑定数
	NewBuddyMemberBindingPercentage string `json:"new_buddy_member_binding_percentage"` //新增好友会员绑定率
	NewBuddyIntoNumber              int64  `json:"new_buddy_into_number"`               //新增好友入群数
	NewBuddyIntoPercentage          string `json:"new_buddy_into_percentage"`           //新增好友入群率
	BuddyChurnNumber                int64  `json:"buddy_churn_number"`                  //好友流失人数
	BuddyChurnPercentage            string `json:"buddy_churn_percentage"`              //好友流失率
	GrandTotalNumber                int64  `json:"grand_total_number"`                  //累计好友数
	BuddyIntoPercenNumber           int64  `json:"buddy_into_percen_number"`            //好友入群数
	BuddyIntoPercentage             string `json:"buddy_into_percentage"`               //好友入群率
}

// BuddyDataDetailList 好友数据明细列表
type BuddyDataDetailList struct {
	Count           int                `json:"count"`
	BuddyDataDetail []*BuddyDataDetail `json:"buddy_data_detail"`
}

// PanelStatisticsData 统计面板数据集合
type PanelStatisticsData struct {
	CommunityData         CommunityData         `json:"community_data"`          //社群数据
	CommunityActivityData CommunityActivityData `json:"community_activity_data"` //社群活动数据
	BuddyData             BuddyData             `json:"buddy_data"`              //好友数据
}

// GetBaFirstDepartmentParam 根据智慧导购的导购的mem_uid，获取导购的一级部门
type GetBaFirstDepartmentParam struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required"`
	MemID        string `json:"mem_id" binding:"required"`
}

// RankDataCk ck查询导购添加好友排行榜
type RankDataCk struct {
	MemberId int64  `json:"member_id"`
	Name     string `json:"name"`    // 导购
	DpName   string `json:"dp_name"` // 导购门店
	Total    int64  `json:"total"`
}

type MemberInfo struct {
	MemberUserId   string `json:"member_user_id"`
	MemberName     string `json:"member_name"`
	DepartmentName string `json:"department_name"`
	Avatar         string `json:"avatar"`
}

package cdp

// UserGroupListResponse 用户分组列表返回数据类型
type UserGroupListResponse struct {
	CheckCodeResponse
	Data struct {
		Aipl          []UserGroupModel       `json:"aipl"`
		Rfm           []UserGroupModel       `json:"rfm"`
		Custom        []CustomUserGroupModel `json:"custom"`
		ValueUpdateAt int64                  `json:"value_update_at"`
	} `json:"data"`
}

// UserGroupModel 用户分组模型
type UserGroupModel struct {
	GroupID       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
	GroupType     int    `json:"group_type"`
	Value         string `json:"value"`
	ValueUpdateAt int64  `json:"value_update_at"`
}

// CustomUserGroupModel 自定义用户分组模型
type CustomUserGroupModel struct {
	GroupID       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
	Value         string `json:"value"`
	ValueUpdateAt int64  `json:"value_update_at"`
	NewAggeResult int    `json:"new_agge_result"`
	Conditions    []struct {
		IncludeType int `json:"include_type"`
		TagGroups   []struct {
			TagGroupID int64   `json:"tag_group_id"`
			TagIds     []int64 `json:"tag_ids"`
		} `json:"tag_groups"`
	} `json:"conditions"`
}

// CreateManualTagResponse 创建手动标签返回结构体
type CreateManualTagResponse struct {
	CheckCodeResponse
	Data struct {
		TagID int64 `json:"tag_id"`
	} `json:"data"`
}

// CreateManualTagRequest cdp创建标签请求参数
type CreateManualTagRequest struct {
	TagGroupID int64  `json:"tag_group_id"`
	TagName    string `json:"tag_name"`
}

// UserGroupPrecalculateCountResponse 用户分组与计算接口返回值
type UserGroupPrecalculateCountResponse struct {
	CheckCodeResponse
	Data struct {
		UserCount int `json:"user_count"`
	} `json:"data"`
}

// CreateManualTagGroupRequest 创建手动标签组接口请求参数
type CreateManualTagGroupRequest struct {
	Name     string   `json:"name"`
	TagNames []string `json:"tag_names"`
}

// CreateManualTagGroupResponse 创建手动标签组接口返回数据
type CreateManualTagGroupResponse struct {
	CheckCodeResponse
	Data struct {
		TagGroupID int64 `json:"tag_group_id"`
	} `json:"data"`
}

// CreateUserGroupResponse 创建用户分组请求参数
type CreateUserGroupResponse struct {
	CheckCodeResponse
	Data struct {
		GroupID int64 `json:"group_id"`
	} `json:"data"`
}

// TagGroupInfoRequest 标签组详情请求参数
type TagGroupInfoRequest struct {
	TagGroupID int64 `json:"tag_group_id"`
	Order      int   `json:"order"`
	Page       int   `json:"page_index"`
	Size       int   `json:"page_size"`
}

// TagGroupInfoResponse 标签组详情返回数据格式
type TagGroupInfoResponse struct {
	CheckCodeResponse
	Data TagGroup `json:"data"`
}

// TagGroup 标签组基本信息
type TagGroup struct {
	GroupID           int64     `json:"tag_group_id"`
	GroupName         string    `json:"name"`
	ValueUpdateAt     int       `json:"value_update_at"`
	SyncQW            int       `json:"sync_qw"`
	TagGroupFirstType int       `json:"tag_group_first_type"`
	TagGroupType      int       `json:"tag_group_type"`
	Description       string    `json:"description"`
	PageCount         int       `json:"page_count"`
	TagList           []TagInfo `json:"tags"`
}

// TagInfo 标签基本信息
type TagInfo struct {
	ID            int64       `json:"tag_id"`
	Name          string      `json:"name"`
	Value         string      `json:"value"`
	NewAggeResult int         `json:"new_agge_result"`
	Condition     interface{} `json:"ConditionsOf"`
}
type CdpStrategyRequest struct {
	SignTime   int64    `json:"sign_time"`
	StrategyId int64    `json:"strategy_id"`
	Day        int      `json:"day"`
	Sts        int      `json:"sts"`
	Ets        int      `json:"ets"`
	Indices    []string `json:"indices"`
	SignType   int      `json:"sign_type"`
	//EnterpriseId int `json:"enterprise_id"`
}
type StrategyDataDetails struct {
	CheckErrorResponse
	Data  []*StrategyItem `json:"data"`
	Total *StrategyItem   `json:"total"`
}

type StrategyItem struct {
	SendMemberNumber   int     `json:"send_member_number"`
	ReachMemberNumber  int     `json:"reach_member_number"`
	FansNumber         int     `json:"fans_number"`
	VisitorNumber      int     `json:"visitor_number"`
	AddOrderNumber     int     `json:"add_order_number"`
	PayOrder           int     `json:"pay_order"`
	FinalAmount        float64 `json:"final_amount"`
	RefundAmountNumber int     `json:"refund_amount_number"`
	RefundAmount       float64 `json:"refund_amount"`
	Types              int     `json:"type"`
}

// UserGroupDetailRequest 用户分组详情请求参数
type UserGroupDetailRequest struct {
	ID int64 `json:"id"`
}

type UserGroupDetailResponse struct {
	CheckCodeResponse
	Data UserGroupDetailEntity `json:"data"`
}

// UserGroupDetailEntity 创建/编辑客户分组请求参数
type UserGroupDetailEntity struct {
	Conditions []*UserGroupDetailCondition `json:"conditions"`
	Name       string                      `json:"user_group_name" desc:"分组名称"`
	ID         int64                       `json:"user_group_id" desc:"客户分组ID"`
	Value      string                      `json:"user_group_value" desc:"客户分组标签值"`
}

// UserGroupDetailCondition 用户分组详情参数
type UserGroupDetailCondition struct {
	TagGroups   []*UserGroupDetailTagGroup `json:"tag_groups" desc:"标签组列表"`
	IncludeType int                        `json:"include_type"`
}

// UserGroupDetailTagGroup 用户分组详情群组信息
type UserGroupDetailTagGroup struct {
	TagGroupID   int64                 `json:"tag_group_id" desc:"标签组id"`
	TagGroupName string                `json:"tag_group_name"`
	Tags         []*UserGroupDetailTag `json:"tags" desc:"标签列表"`
}

// UserGroupDetailTag 用户分组详情标签信息
type UserGroupDetailTag struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}

type CouponStrategyCdpRequest struct {
	SignTime   int64    `json:"sign_time"`
	StrategyId int64    `json:"strategy_id"`
	Sts        int      `json:"sts"`
	Ets        int      `json:"ets"`
	Indices    []string `json:"indices"`
	SignType   int      `json:"sign_type"`
	Gid        string   `json:"gid"`
}

type CouponStrategyItem struct {
	UsedNum    int     `json:"hx_num"`
	PayNum     int     `json:"pay_uv"`
	CashFee    float64 `json:"sf_fee"`
	RefundNum  int     `json:"re_uv"`
	RefundFee  float64 `json:"re_fee"`
	StockID    int     `json:"stock_id"`
	SendNum    int64   `json:"send_num"`
	TotalTouch int64   `json:"touch_num"`
	IsSendNum  int64   `json:"is_send_num"`
}

type UserGroupsByUidResult struct {
	CheckErrorResponse
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
	Data      struct {
		Groups []int64 `json:"groups"`
	} `json:"data"`
}

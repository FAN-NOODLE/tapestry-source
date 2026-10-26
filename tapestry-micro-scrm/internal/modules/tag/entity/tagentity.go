package entity

const (
	GroupTypeAuto   = 1
	GroupTypeManual = 2
)

type CustomerManualTagRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	RelationID   int64 `json:"relation_id" db:"relation_id" validate:"required" required:"客户关系ID不能为空"`
}

// TagInitRequest cdP 标签请求
type TagInitRequest struct {
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID"`
	CorpID       string `json:"corp_id" desc:"企业corpId"`
}

// TagInfo 标签基本信息
type TagInfo struct {
	ID            int64                  `json:"tag_id"`
	Name          string                 `json:"name"`
	Value         string                 `json:"value"`
	WxTagID       string                 `json:"qw_id"`
	NewAggeResult int                    `json:"new_agge_result"`
	GoodData      *TagGoodInfo           `json:"good_data"`
	CouponData    *TagCouponInfo         `json:"coupon_data"`
	Condition     map[string]interface{} `json:"ConditionsOf"`
	GroupID       int64
}

// CouponTagTypeEnum 优惠券标签类型
var CouponTagTypeEnum = struct {
	UserCount int // 按券使用次数
	Coupon    int // 按优惠券名称
}{
	UserCount: 1,
	Coupon:    2,
}

// TagConditionCoupon 券类标签条件
type TagConditionCoupon struct {
	Value    int    `json:"value"`
	CouponID string `json:"coupon_id"`
	Type     int    `json:"type"`
}

// TagCondition 标签配置结构体
type TagCondition struct {
	Name      string `json:"name"`
	Min       int    `json:"min"`
	Max       int    `json:"max"`
	DateType  int    `json:"date_type"`
	DateBegin string `json:"date_beg"`
	DateEnd   string `json:"date_end"`
	SpuID     int64  `json:"spu_id"`
	Value     int    `json:"value"`
}

// TagConditionEdit 标签配置更新结构体
type TagConditionEdit struct {
	TagCondition
	UpdateType int `json:"update_type"`
}

// TagGoodInfo 标签中商品显示信息
type TagGoodInfo struct {
	ImgURL string `json:"img_url"`
	Name   string `json:"name"`
}

// TagCouponInfo 标签中优惠券显示信息
type TagCouponInfo struct {
	Name    string `json:"coupon_name"`
	StockID string `json:"stock_id"`
}

// TagGroupList 标签组列表
type TagGroupList struct {
	AutoList   []TagGroup `json:"auto"`
	ManualList []TagGroup `json:"manual"`
}

type TagGroupBase struct {
	GroupID           int64  `json:"tag_group_id"`
	GroupName         string `json:"tag_group_name"`
	WxTagGroupID      string `json:"qw_id"`
	ValueUpdateAt     int    `json:"value_update_at"`
	SyncQW            int    `json:"sync_qw"`
	TagGroupFirstType int    `json:"tag_group_first_type"`
	TagGroupType      int    `json:"tag_group_type"`
	Description       string `json:"description"`
	PageCount         int    `json:"page_count"`

	IsDelete int `json:"is_delete"`
	Type     int // 标签分组类型，1-自动标签，2-手动标签
}

// TagGroup 标签组基本信息
type TagGroup struct {
	TagGroupBase
	TagList TagList `json:"tags"`
}

// CustomerTagResponse 客户标签列表
type CustomerTagResponse struct {
	ManualList []CustomerTagGroup `json:"manual"`
}

// CustomerTag 客户标签
type CustomerTagGroup struct {
	TagGroupBase
	TagList []*CustomerTag `json:"tags"`
}

// CustomerTag 客户标签列表
type CustomerTag struct {
	TagInfo
	IsChoice int `json:"is_choice"`
}

type TagList []*TagInfo

func (t TagList) Len() int { return len(t) }

func (t TagList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TagList) Less(i, j int) bool {
	if t[i].Condition == nil || t[j].Condition == nil {
		return true
	}
	if t[i].Condition["min"].(float64) == t[j].Condition["min"].(float64) {
		if int64(t[i].Condition["max"].(float64)) == -1 {
			return false
		} else if int64(t[j].Condition["max"].(float64)) == -1 {
			return true
		}
		return t[i].Condition["max"].(float64) < t[j].Condition["max"].(float64)
	}
	return t[i].Condition["min"].(float64) < t[j].Condition["min"].(float64)
}

// CreateUserGroupResponse 创建用户群组返回数据结构
type CreateUserGroupResponse struct {
	GroupID int64 `json:"group_id"`
}

// UserGroupList 用户分组列表
type UserGroupList struct {
	Aipl          []UserGroupModel       `json:"aipl"`
	Rfm           []UserGroupModel       `json:"rfm"`
	Custom        []CustomUserGroupModel `json:"custom"`
	ValueUpdateAt int64                  `json:"value_update_at"`
}

// UserGroupModel 用户分组模型
type UserGroupModel struct {
	GroupID       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
	GroupType     int    `json:"group_type"`
	Value         string `json:"value"`
	ValueUpdateAt int64  `json:"value_update_at"`
}

// CustomUserGroupModel 用户分组自定义模型
type CustomUserGroupModel struct {
	GroupID       int64                `json:"group_id"`
	GroupName     string               `json:"group_name"`
	Value         string               `json:"value"`
	ValueUpdateAt int64                `json:"value_update_at"`
	NewAggeResult int                  `json:"new_agge_result"`
	Conditions    []UserGroupCondition `json:"conditions"`
}

// UserGroupCondition 用户分组自定义模型条件
type UserGroupCondition struct {
	IncludeType int                `json:"include_type"`
	TagGroups   []UserGroupTagInfo `json:"tag_groups"`
}

// UserGroupTagInfo 用户分组自定义模型标签信息
type UserGroupTagInfo struct {
	TagGroupID int64   `json:"tag_group_id"`
	TagIds     []int64 `json:"tag_ids"`
}

// CreateUserGroupBody 创建/编辑客户分组请求参数
type CreateUserGroupBody struct {
	Conditions    []ConditionsParams `json:"conditions"`
	Name          string             `json:"name" desc:"分组名称"`
	ID            int64              `json:"id" desc:"客户分组ID"`
	UserGroupType int                `json:"user_group_type"`
	FileID        int64              `json:"file_id"`
	DepartmentIDs []int64            `json:"department_ids"`
}

// ConditionsParams 用户群组参数
type ConditionsParams struct {
	TagGroups   []TagGroupsParams `json:"tag_groups" desc:"标签组列表"`
	IncludeType int               `json:"include_type"`
}

// TagGroupsParams 用户群组
type TagGroupsParams struct {
	TagGroupID int64 `json:"tag_group_id" desc:"标签组id"`
	TagIds     []int `json:"tag_ids" desc:"标签列表"`
}

// UserGroupDetail 创建/编辑客户分组请求参数
type UserGroupDetail struct {
	Conditions    []*UserGroupDetailCondition `json:"conditions"`
	Name          string                      `json:"name" desc:"分组名称"`
	ID            int64                       `json:"id" desc:"客户分组ID"`
	Value         string                      `json:"value" desc:"客户分组标签值"`
	CanEdit       bool                        `json:"can_edit"`
	FileID        int64                       `json:"file_id"`
	IsExtract     int                         `json:"is_extract"`
	UserGroupType int8                        `json:"user_group_type"`
	Departments   []Departments
}

type Departments struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
	TagName string `json:"name"`
}

// AutoTagUpdateRequest 自动标签更新接口请求参数
type AutoTagUpdateRequest struct {
	TagGroupID     int64       `json:"tag_group_id" desc:"标签分组id"`
	TagGroupType   int         `json:"tag_group_type" desc:"标签组类型"`
	SyncWechatWork int         `json:"sync_qw" desc:"是否同步到企业微信：1同步，2不同步"`
	Conditions     interface{} `json:"conditions"`
}

// UpdateTagGroupSyncRequest 更新标签组同步企业微信标识接口请求参数格式
type UpdateTagGroupSyncRequest struct {
	TagGroupID int64 `json:"tag_group_id" desc:"标签组id"`
	SyncQW     int   `json:"sync_qw" desc:"是否同步到企业微信 1:是 2：否"`
}

// UserGroupPrecalculateRequest 用户组预计算请求参数
type UserGroupPrecalculateRequest struct {
	Conditions []UserGroupPrecalculateCondition `json:"conditions"`
}

// UserGroupPrecalculateResponse 用户组预计算返回数据
type UserGroupPrecalculateResponse struct {
	UserCount int `json:"user_count"`
}

// UserGroupPrecalculateCondition 客户分组预计算condition结构体
type UserGroupPrecalculateCondition struct {
	IncludeType int                 `json:"include_type"`
	TagGroups   []UserGroupTagGroup `json:"tag_groups"`
}

// UserGroupTagGroup 客户分组标签组组成结构体
type UserGroupTagGroup struct {
	TagGroupID int64   `json:"tag_group_id"`
	TagIds     []int64 `json:"tag_ids"`
}

// TagGroupInfoRequest 标签组详情接口请求参数
type TagGroupInfoRequest struct {
	TagGroupID     int64   `json:"tag_group_id" db:"tag_group_id"`
	TagGroupIds    []int64 `json:"tag_group_ids"`
	Order          int     `json:"order" db:"order"`
	Page           int     `json:"page"`
	Size           int     `json:"size"`
	NeedGoods      bool    `json:"need_goods"`
	NeedCoupon     bool    `json:"need_coupon"`
	TagKeywords    string  `json:"tag_keywords"`
	IgnoreTagCount bool    `json:"ignore_tag_count"`
}

// CreateTagGroupResponse 创建手动标签组返回数据结构体
type CreateTagGroupResponse struct {
	TagGroupID int64 `json:"tag_group_id"`
}

// TagGroupType 标签组类型
type TagGroupType int

// TagGroupTypeEnum 标签组类型枚举
var TagGroupTypeEnum = struct {
	UserChannel                TagGroupType // 1:用户渠道
	WechatWorkCustomerChannel  TagGroupType // 2:企微客户渠道
	WechatWorkCustomerStaff    TagGroupType // 3:企微客户所属员工
	WechatWorkUserTag          TagGroupType // 4:企微个人标签
	WechatOfficialAccountGroup TagGroupType // 5:公众号分组
	ChannelWithParam           TagGroupType // 6:带参渠道
	VisitActivation            TagGroupType // 7:访问活跃度
	VisitFrequency             TagGroupType // 8:访问次数
	VisitDepth                 TagGroupType // 9:访问深度
	VisitDuration              TagGroupType // 10:访问时长
	VisitTime                  TagGroupType // 11:访问时间
	PayTime                    TagGroupType // 12:支付时间
	OrderTime                  TagGroupType // 13:下单时间
	VisitGood                  TagGroupType // 14:浏览商品
	AddGoodToCart              TagGroupType // 15:加购未支付指定商品
	PayForGood                 TagGroupType // 16:支付指定商品
	PayFrequency               TagGroupType // 17:支付活跃度
	PerCustomerTransaction     TagGroupType // 18:客单价
	GMV                        TagGroupType // 19:成交总额
	RePurchaseCycle            TagGroupType // 20:复购周期
	PurchaseFrequency          TagGroupType // 21:购买频次
	RefundRate                 TagGroupType // 22:退款率
	AuthorizeGender            TagGroupType // 23:授权性别
	Region                     TagGroupType // 24:地区
	UseCoupon                  TagGroupType // 25:是否用券
	ReceiveCoupon              TagGroupType // 26:是否领券
}{
	UserChannel:                1,
	WechatWorkCustomerChannel:  2,
	WechatWorkCustomerStaff:    3,
	WechatWorkUserTag:          4,
	WechatOfficialAccountGroup: 5,
	ChannelWithParam:           6,
	VisitActivation:            7,
	VisitFrequency:             8,
	VisitDepth:                 9,
	VisitDuration:              10,
	VisitTime:                  11,
	PayTime:                    12,
	OrderTime:                  13,
	VisitGood:                  14,
	AddGoodToCart:              15,
	PayForGood:                 16,
	PayFrequency:               17,
	PerCustomerTransaction:     18,
	GMV:                        19,
	RePurchaseCycle:            20,
	PurchaseFrequency:          21,
	RefundRate:                 22,
	AuthorizeGender:            23,
	Region:                     24,
	UseCoupon:                  25,
	ReceiveCoupon:              26,
}

// SortTagGroupRequest 标签组排序接口请求参数数据结构
type SortTagGroupRequest struct {
	TagGroupFirstType int     `json:"tag_group_first_type"`
	TagGroupIds       []int64 `json:"tag_group_ids"`
}

// SearchTagRequest 标签搜索接口请求参数
type SearchTagRequest struct {
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	Keyword    string `json:"keyword"`
	SearchType int    `json:"search_type"`
}

// DeleteManualTagGroupRequest 删除手动标签组
type DeleteManualTagGroupRequest struct {
	TagGroupID int64 `json:"tag_group_id" binding:"required"`
}

type DeleteManualTagGroupRequestParam struct {
	EnterpriseID int64 `db:"enterprise_id"`
	TagGroupID   int64 `db:"tag_group_id"`
}

// UpdateManualTagRequest 更新手动标签信息
type UpdateManualTagRequest struct {
	TagID   int64  `json:"tag_id" binding:"required"`
	TagName string `json:"tag_name" binding:"required"`
}

// DeleteManualTagRequest 删除手动标签
type DeleteManualTagRequest struct {
	TagID int64 `json:"tag_id" binding:"required"`
}

// AutoTagUpdateBase 自动标签更新接口公共参数
type AutoTagUpdateBase struct {
	TagGroupID int64 `json:"tag_group_id" binding:"required"`
	SyncQW     int   `json:"sync_qw"`
}

// AutoTagUpdateCondBase 自动标签更新条件数据结构公共参数
type AutoTagUpdateCondBase struct {
	Name       string `json:"name" binding:"required"`
	TagID      int64  `json:"tag_id"`
	UpdateType int    `json:"update_type" binding:"required"`
}

// AutoTagUpdateNumberRange 自动标签更新请求参数：最小最大数字类型参数
type AutoTagUpdateNumberRange struct {
	AutoTagUpdateBase
	Conditions []struct {
		AutoTagUpdateCondBase
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"conditions" binding:"required"`
}

// AutoTagUpdateTimeRange 自动标签更新请求参数：最小最大数字类型参数
type AutoTagUpdateTimeRange struct {
	AutoTagUpdateBase
	Conditions []struct {
		AutoTagUpdateCondBase
		DateType int    `json:"date_type" binding:"required"`
		DateBeg  string `json:"date_beg"`
		DateEnd  string `json:"date_end"`
	} `json:"conditions" binding:"required"`
}

// CreateTagRequest 创建手动标签请求参数结构体
type CreateTagRequest struct {
	TagGroupID int64  `json:"tag_group_id"`
	TagName    string `json:"tag_name"`
}

// CreateTagResponse 创建标签返回参数
type CreateTagResponse struct {
	TagID int64 `json:"tag_id"`
}

// CreateGroupRequest 创建手动标签组请求参数结构体
type CreateGroupRequest struct {
	GroupId     int64    `json:"group_id"`
	GroupName   string   `json:"group_name" db:"group_name"`
	TagNames    []string `json:"tag_names" binding:"required"`
	ExpiredTime int64    `json:"expired_time" db:"expired_time"`
	FileUrl     string   `json:"file_url"`
}

// TagEntity 标签实体
type TagEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`
	TagID        int64  `db:"tag_id"`
	TagName      string `db:"tag_name"`
	TagGroupID   int64  `db:"tag_group_id"`
	CreatedAt    int    `db:"created_at"`
	UpdatedAt    int    `db:"updated_at"`
	DeletedAt    int    `db:"deleted_at"`
}

// UserTagEntity 用户标签实体
type UserTagEntity struct {
	UserID int64 `db:"user_id"`
	TagEntity
}

// TagGroupEntity 标签组实体
type TagGroupEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`
	GroupID      int64  `db:"group_id"`
	GroupName    string `db:"group_name"`
	Type         int    `db:"type"`
	CreatedAt    int    `db:"created_at"`
	UpdatedAt    int    `db:"updated_at"`
	DeletedAt    int    `db:"deleted_at"`
}

// UserTagGroupEntity 用户标签组实体
type UserTagGroupEntity struct {
	UserID int64 `db:"user_id"`
	UserGroupEntity
}

// GetTagListRequest 标签列表请求参数
type GetTagListRequest struct {
	Type         int    `json:"type"`
	Search       string `json:"search"`
	EnterpriseID int64
	Page         int `json:"page"`
	Size         int `json:"size"`
}

// GetTagListItem 标签列表数据格式
type GetTagListItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserGroupEntity 用户分组实体
type UserGroupEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"` // UserGroupID cdp客户分组id
	UserGroupID  int64  `db:"user_group_id"` // Name 客户分组名称
	Name         string `db:"name"`          // Type 客户分组类型：1-忠诚度模型AIPL,2-RFM模型，3-自定义
	Type         int    `db:"type"`
	CreatedAt    int    `db:"created_at"`
	UpdatedAt    int    `db:"updated_at"`
	DeletedAt    int    `db:"deleted_at"`
}

type WorkExternalTag struct {
	RelationID int64  `json:"relation_id" db:"relation_id"`
	TagName    string `json:"tag_name" db:"tag_name"`
	TagID      string `json:"tag_id" db:"tag_id"`
}

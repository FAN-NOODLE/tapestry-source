package drainage

import "gw-scrm/internal/pkg/sdk"

type ContactsListResponse struct {
	Name         string `json:"name" db:"name" desc:"name"`
	ChannelId    int64  `json:"channel_id"  db:"channel_id" desc:"渠道id"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Page         int    `json:"page" db:"page" desc:"页码"`
	Size         int    `json:"size"  db:"size" desc:"分页"`
}
type ContactCreateRequest struct {
	ID           int64    `json:"id" db:"id" desc:"id"`
	CodeCame     string   `json:"code_name" binding:"required"  db:"code_name" desc:"二维码名称"`
	ChannelId    int64    `json:"channel_id" binding:"required"  db:"channel_id" desc:"渠道id"`
	AutoFriend   int      `json:"auto_friend"  db:"auto_friend" desc:"是否自动开启"`
	MemberIds    []int64  `json:"member_ids" binding:"required"  db:"member_ids" desc:"员工id"`
	ThirdIds     []int    `json:"third_ids" desc:"部门Id"`
	EnterpriseId int64    `json:"enterprise_id" binding:"required" db:"enterprise_id" desc:"企业ID"`
	UserIds      []string `json:"user_ids" binding:"required" db:"user_ids" desc:"user_ids"`
	Type         int      `json:"type"  db:"type" desc:"类型"`
}
type ContactDeleteRequest struct {
	ID           []int64 `json:"ids" db:"ids" binding:"required" desc:"ids"`
	EnterpriseId int64   `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}
type ContactDetailsRequest struct {
	ID           int64 `json:"id" db:"id" binding:"required" desc:"id"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

type QrCodeRequest struct {
	EnterpriseId int64  `json:"enterprise_id"`
	ImageUrl     string `json:"image_url"`
	Width        int    `json:"width"`
	Weight       int    `json:"height"`
}

// Welcome drainage_welcomes 引流模块欢迎语表
type Welcome struct {
	ID           int64  `json:"id" db:"id" desc:"id"`
	Name         string `json:"name" db:"name" desc:"主题名称"`
	Content      string `json:"content" db:"content" desc:"内容"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	IsDefault    int64  `json:"is_default" db:"is_default" desc:"是否是默认欢迎语 1是 0不是"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// WelcomeMember drainage_welcome_members 引流模块欢迎语关联员工表
type WelcomeMember struct {
	ID           int64 `json:"id" db:"id" desc:"id"`
	WelcomeID    int64 `json:"welcome_id" db:"welcome_id" desc:"drainage_welcomes表ID"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberID     int64 `json:"member_id" db:"member_id" desc:"员工ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// WelcomeCreateRequest 创建欢迎语请求
type WelcomeCreateRequest struct {
	ID           int64                    `json:"id" db:"id" desc:"id"`
	Content      []map[string]interface{} `json:"content" db:"content" binding:"required" desc:"内容"`
	EnterpriseId int64                    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	IsDefault    int64                    `json:"is_default" db:"is_default" desc:"是否是默认欢迎语 1是 0不是"`
	ChannelId    int64                    `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	CreatedAt    int64                    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64                    `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64                    `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// WelcomeCreateResponse 创建欢迎语响应
type WelcomeCreateResponse struct {
	Result bool `json:"result" desc:"是否创建成功"`
}

// WelcomeDetailRequest 欢迎语详情请求
type WelcomeDetailRequest struct {
	ID           int64 `json:"id" db:"id" desc:"id"`
	IsDefault    int64 `json:"is_default" db:"is_default" desc:"是否是默认欢迎语 1是 0不是"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Page         int64 `json:"page" db:"page" desc:"page"`
	Size         int64 `json:"size" db:"size" desc:"size"`
}

// WelcomeDetailResponse 欢迎语详情响应
type WelcomeDetailResponse struct {
	ID           int64                    `json:"id" db:"id" desc:"id"`
	Content      []map[string]interface{} `json:"content" db:"content" desc:"内容"`
	EnterpriseId int64                    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	IsDefault    int64                    `json:"is_default" db:"is_default" desc:"是否是默认欢迎语 1是 0不是"`
	ChannelId    int64                    `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	Channel      string                   `json:"channel" desc:"渠道名称"`
	CreatedAt    int64                    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64                    `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64                    `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	Page         int64                    `json:"page" db:"page" desc:"page"`
	Size         int64                    `json:"size" db:"size" desc:"size"`
	Total        int64                    `json:"total" db:"total" desc:"total"`
}

// WelcomeEditRequest 编辑欢迎语请求
type WelcomeEditRequest struct {
	OldID int64 `json:"old_id" db:"old_id" binding:"required" desc:"old_id旧欢迎语的ID"`
	WelcomeCreateRequest
}

// WelcomeEditResponse 编辑欢迎语响应
type WelcomeEditResponse struct {
	Result bool `json:"result" desc:"是否编辑成功"`
}

// WelcomeDeleteRequest 删除欢迎语请求
type WelcomeDeleteRequest struct {
	ID        int64 `json:"id" db:"id" binding:"required" desc:"id"`
	DeletedAt int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// WelcomeDeleteResponse 删除欢迎语响应
type WelcomeDeleteResponse struct {
	Result bool `json:"result" desc:"是否删除成功"`
}

// WelcomeListRequest 列表欢迎语请求
type WelcomeListRequest struct {
	Page         int64  `json:"page" db:"page"  binding:"required" desc:"page"`
	Size         int64  `json:"size" db:"size"  binding:"required" desc:"size"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ChannelId    int64  `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	MemberID     int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemberIDStr  string `json:"member_id_str" db:"member_id_str" desc:"员工ID"`
	Content      string `json:"content" db:"content" desc:"内容"`
}

// WelcomeListResponse 列表欢迎语响应
type WelcomeListResponse struct {
	Page  int64               `json:"page" db:"page" desc:"page"`
	Size  int64               `json:"size" db:"size" desc:"size"`
	Total int64               `json:"total" db:"total" desc:"total"`
	List  []WelcomeListStruct `json:"list" db:"list" desc:"list"`
}

// WelcomeListStruct 列表欢迎语响应子结构体
type WelcomeListStruct struct {
	ID         int64                    `json:"id" db:"id" desc:"id"`
	ContentStr string                   `json:"content_str" db:"content_str" desc:"内容"`
	Content    []map[string]interface{} `json:"content" db:"content" desc:"内容"`
	CreatedAt  int64                    `json:"created_at" db:"created_at" desc:"创建时间"`
	ChannelId  int64                    `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	Channel    string                   `json:"channel" desc:"渠道名称"`
}

// WelcomeSendRequest 列表欢发送请求
type WelcomeSendRequest struct {
	MemberID     int64 `json:"member_id" db:"member_id" desc:"员工ID"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// WelcomeSendResponse 欢迎语发送响应
type WelcomeSendResponse struct {
	Result bool `json:"result" desc:"是否发送成功"`
}

// WelcomeValidateRequest 员工欢迎语是否已配置请求
type WelcomeValidateRequest struct {
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	AddChannel   bool  `json:"add_channel" desc:"是否添加渠道"`
}

// WelcomeValidateResponse 员工欢迎语是否已配置响应
type WelcomeValidateResponse struct {
	Page  int64               `json:"page" db:"page"`
	Size  int64               `json:"size" db:"size"`
	Total int64               `json:"total" db:"total"`
	List  []WelcomeMemberInfo `json:"member_ids" db:"member_ids" desc:"员工IDs"`
}

// WelcomeMemberInfo 员工欢迎语是否已配置响应子结构体
type WelcomeMemberInfo struct {
	MemberId int64  `json:"member_id" db:"member_id"`
	Name     string `json:"name" db:"name"`
}

type WelcomeAllRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	KeyWords     string `json:"key_words"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
}

type KdTypaAll struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
}

type WelcomeBase struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	ID           int64 `json:"id" required:"数据id不能为空"`
}

type WelcomeNewEditRequest struct {
	ID               int64                    `json:"id"`
	EnterpriseId     int64                    `json:"enterprise_id"`
	Content          []map[string]interface{} `json:"content"`
	ConditionContent map[string][]string      `json:"condition_content" comment:"条件集合json格式"`
}

// ChannelsDataRequest 渠道统计数据请求
type ChannelsDataRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	StartTime    int64 `json:"start_time" db:"start_time"`
	EndTime      int64 `json:"end_time" db:"end_time"`
	ChannelsId   int64 `json:"channels_id" db:"channels_id"`
	Page         int64 `json:"page" db:"page"`
	Size         int64 `json:"size" db:"size"`
	Type         int64 `json:"type" db:"type"`
}

// ChannelsCountData 渠道统计数据
type ChannelsCountData struct {
	sdk.ResponseBase
	Data *ChannelsCount `json:"data"`
}
type ChannelsCount struct {
	List  []*ChannelsDataListResponse `json:"list"`
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
	Total int64                       `json:"total"`
}
type ChannelsDataListResponse struct {
	TimeOrName              string `json:"time_or_name"`
	NewlyAddedCustomerCount int64  `json:"newly_added_customer_count"` // 净增客户数
	ExistRelationCustomer   int64  `json:"exist_relation_customer"`    // 留存好友数
	IncreaseCustomerCount   int64  `json:"increase_customer_count"`    // 新增好友总数
	LossCustomerCount       int64  `json:"loss_customer_count"`        // 流失好友数
	CustomerDelMember       int64  `json:"customer_del_member"`        // 好友删除员工数
	MemberDelCustomer       int64  `json:"member_del_customer"`        // 员工删除客户总数
}

// CacheMediaRequest 缓存微信图片素材请求
type CacheMediaRequest struct {
	ImgUrl       string `json:"img_url" db:"img_url" binding:"required" desc:"img_url"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	IsTest       int64  `json:"is_test"`
}

// CacheMediaResponse 缓存微信图片素材响应
type CacheMediaResponse struct {
	Result  bool   `json:"result"`
	MediaId string `json:"mediaId"`
}

type WelcomeChannelsList struct {
	ID           int64  `json:"id" db:"id" desc:"删除时间"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Channel      string `json:"channel"  db:"channel" desc:"渠道名称"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

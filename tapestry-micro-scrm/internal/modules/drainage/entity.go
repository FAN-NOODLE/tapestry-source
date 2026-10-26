package drainage

import (
	"micro-scrm/internal/orm/model/micro_scrm/welcome"
)

const (
	ActivityWelcome = 2
)

// Welcome drainage_welcomes 引流模块欢迎语表
type Welcome struct {
	ID               int64  `json:"id" db:"id" desc:"id"`
	Name             string `json:"name" db:"name" desc:"主题名称"`
	Content          string `json:"content" db:"content" desc:"内容"`
	EnterpriseId     int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	IsDefault        int64  `json:"is_default" db:"is_default" desc:"是否是默认欢迎语 1是 0不是"`
	ChannelId        int64  `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	Channel          string `json:"channel" desc:"渠道名称"`
	Type             int    `json:"type" desc:"欢迎语类型：1-系统欢迎语，2-裂变活动欢迎语"`
	CreatedAt        int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt        int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt        int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	ConditionContent string `json:"condition_content" db:"condition_content" desc:"条件集合"`
}

// WelcomeMember drainage_welcome_members 引流模块欢迎语关联员工
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
	IdNot        int64  `json:"id_not" db:"id" desc:"去除"`
	MemberID     int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemberIDStr  string `json:"member_id_str" db:"member_id_str" desc:"员工ID"`
	Content      string `json:"content" db:"content" desc:"内容"`
	Type         int64  `json:"type" db:"type" desc:"欢迎语类型：1-系统欢迎语，2-裂变活动欢迎语"`
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

// WelcomeValidateRequest 渠道欢迎语选择渠道列表
type WelcomeValidateRequest struct {
	EnterpriseId int64   `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	AddChannel   bool    `json:"add_channel" desc:"是否添加渠道"`
	Filter       []int64 `json:"filter_ids" desc:"过滤id"`
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

// 联系我们
type Contacts struct {
	ID                 int64            `json:"id" db:"id" desc:"id"`
	CodeName           string           `json:"code_name" db:"code_name" desc:"code_name"`
	Name               string           `json:"name" db:"name" desc:"name"`
	ChannelId          int64            `json:"channel_id"  db:"channel_id" desc:"渠道id"`
	AutoFriend         int              `json:"auto_friend"  db:"auto_friend" desc:"是否自动开启"`
	QrCodeUrl          string           `json:"qr_code_url"  db:"qr_code_url" desc:"二维码地址"`
	MemberCount        int              `json:"member_count"  db:"member_count" desc:"创建员工数量"`
	EnterpriseId       int64            `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UserId             []ContactMembers `json:"user_id"  db:"user_id" desc:"员工id"`
	CreatedAt          int64            `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt          int64            `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt          int64            `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	MemberDatas        []MemberData     `json:"member_data" db:"member_data" desc:"员工数据"`
	DepartmentData     []DepartmentData `json:"department_data" db:"department_data" desc:"部门数据"`
	ChannelName        string           `json:"channel_name"  db:"channel_name" desc:"渠道名称"`
	SumAddMemberNumber int              `json:"sum_add_member_number"  db:"sum_add_member_number" desc:"累加添加数量"`
	ConfigId           string           `json:"config_id"  db:"config_id" desc:"微信配置id"`
	Type               int              `json:"type"  db:"type" desc:"类型 0 多人共享码（原本的） 1 员工码+门店码"`
}

// ShowContacts 联系我详情页
type ShowContacts struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	ChannelName string         `json:"channel_name"`
	TagNameList []string       `json:"tag_name_list"`
	Type        int            `json:"type" desc:"类型 0 多人共享码（原本的） 1 员工码+门店码"`
	QrcodeList  []MemberQrcode `json:"qrcode_list"`
}

// MemberQrcode 联系我详情页子列表
type MemberQrcode struct {
	Name           string   `json:"name"`
	DeptMemberList []string `json:"dept_member_list"`
	AddNumber      int64    `json:"add_number" desc:"添加好友数量"`
	CreatedAt      int64    `json:"created_at" db:"created_at" desc:"创建时间"`
	QrCodeUrl      string   `json:"qr_code_url"  db:"qr_code_url" desc:"二维码地址"`
}

// 联系我员工信息
type MemberData struct {
	ID         int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemberName string `json:"member_name" db:"member_name" desc:"员工名称"`
	UserId     string `json:"user_id" db:"user_id" desc:"用户id"`
}

// 联系我部门信息
type DepartmentData struct {
	ThirdId        int    `json:"third_id" db:"third_id" desc:"部门id"`
	DepartmentName string `json:"department_name" db:"department_name" desc:"部门名称"`
}

// 联系我员工表
type ContactMembers struct {
	ID           int64 `json:"id" db:"id" desc:"id"`
	ContactId    int64 `json:"contact_id" db:"contact_id" desc:"联系我id"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	MemberId     int64 `json:"member_id" db:"member_id" desc:"用户id"`
	ThirdId      int   `json:"third_id" db:"third_id" desc:"企微部门id"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

const ContactTypeStoreCode = 2

// 联系我新增请求
type ContactCreateRequest struct {
	ID           int64         `json:"id" db:"id" desc:"id"`
	CodeCame     string        `json:"code_name" binding:"required"  db:"code_name" desc:"二维码名称"`
	ChannelId    int64         `json:"channel_id" binding:"required"  db:"channel_id" desc:"渠道id"`
	AutoFriend   int           `json:"auto_friend"   db:"auto_friend" desc:"是否自动开启"`
	MemberIds    []int64       `json:"member_ids" binding:"required"  db:"member_ids" desc:"员工id"`
	ThirdIds     []int         `json:"third_ids" desc:"部门Id"`
	EnterpriseId int64         `json:"enterprise_id" binding:"required" db:"enterprise_id" desc:"企业ID"`
	Type         int           `json:"type"  db:"type" desc:"类型 0 多人共享码（原本的） 1 员工码+门店码"`
	UserIds      []string      `json:"user_ids" binding:"required" db:"user_ids" desc:"user_ids"`
	State        string        `json:"state"`
	WelcomeList  []WelcomeInfo `json:"welcome_list"`
	Domain       string        `json:"domain"`
}

type WelcomeInfo struct {
	ID      int64                    `json:"id" db:"id" desc:"id"`
	Type    int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content []map[string]interface{} `json:"content" binding:"required" desc:"内容"`
}

// 联系我详情请求
type ContactDetailsRequest struct {
	ID           int64 `json:"id" db:"id" binding:"required" desc:"id"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// 联系我删除请求
type ContactDeleteRequest struct {
	ID           []int64 `json:"ids" db:"ids" binding:"required" desc:"ids"`
	EnterpriseId int64   `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// 联系我列表请求
type ContactsListResponse struct {
	IDs          []int64 `json:"ids"`
	IdsCount     int
	Name         string `json:"name" db:"name" desc:"name"`
	ChannelId    int64  `json:"channel_id"  db:"channel_id" desc:"渠道id"`
	EnterpriseId int64  `json:"enterprise_id" binding:"required"  db:"enterprise_id" desc:"企业ID"`
	Page         int    `json:"page" db:"page" desc:"页码"`
	Size         int    `json:"size"  db:"size" desc:"分页"`
}

// 联系我列表公共返回参数
type CommonList struct {
	List  []Contacts `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page" db:"page" desc:"页码"`
	Size  int        `json:"size"  db:"size" desc:"分页"`
}

// CacheMediaRequest 缓存微信图片素材请求
type CacheMediaRequest struct {
	ImgUrl       string `json:"img_url" db:"img_url" binding:"required" desc:"img_url"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	FileType     string `json:"file_type" desc:"图片1"`
	IsTest       int64  `json:"is_test"`
	FileName     string `json:"file_name"`
}

// CacheMediaResponse 缓存微信图片素材响应
type CacheMediaResponse struct {
	Result  bool   `json:"result"`
	MediaId string `json:"mediaId"`
}

type QrCodeRequest struct {
	ImageUrl string `json:"image_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type ChannelsDataRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	StartTime    int64 `json:"start_time" db:"start_time"`
	EndTime      int64 `json:"end_time" db:"end_time"`
	ChannelsId   int64 `json:"channels_id" db:"channels_id"`
	Page         int64 `json:"page" db:"page"`
	Size         int64 `json:"size" db:"size"`
	Type         int64 `json:"type" db:"type"`
}

// ChannelsDataResponse 渠道统计返回结构体
type ChannelsDataResponse struct {
	TotalCustomer            int64                       `json:"total_customer"`
	NewlyAddedCustomerCount  int64                       `json:"newly_added_customer_count"`
	LossCustomerCount        int64                       `json:"loss_customer_count"`
	IncreaseCustomerCount    int64                       `json:"increase_customer_count"`
	CustomerDelMember        int64                       `json:"customer_del_member"`
	MemberDelCustomer        int64                       `json:"member_del_customer"`
	ChannelsDataListResponse []*ChannelsDataListResponse `json:"list"`
	Page                     int64                       `json:"page"`
	Size                     int64                       `json:"size"`
	Total                    int                         `json:"total"`
}

type ChannelsDataListResponse struct {
	TimeOrName              string `json:"time_or_name"`
	TotalCustomer           int64  `json:"total_customer"`
	NewlyAddedCustomerCount int64  `json:"newly_added_customer_count"`
	LossCustomerCount       int64  `json:"loss_customer_count"`
	IncreaseCustomerCount   int64  `json:"increase_customer_count"`
	CustomerDelMember       int64  `json:"customer_del_member"`
	MemberDelCustomer       int64  `json:"member_del_customer"`
}

// WorkExternalUserRelation 外部联系人关系表
type WorkExternalUserRelation struct {
	ID                     int64  `json:"id" db:"id"`
	EnterpriseID           int64  `json:"enterprise_id" db:"enterprise_id"`
	WorkExternalUserID     int64  `json:"work_external_user_id" db:"work_external_user_id"`
	MemberID               int64  `json:"member_id" db:"member_id"`
	Remark                 string `json:"remark" db:"remark"`
	Description            string `json:"description" db:"description"`
	Createtime             int64  `json:"createtime" db:"createtime"`
	RemarkCorpName         string `json:"remark_corp_name" db:"remark_corp_name"`
	RemarkMobiles          string `json:"remark_mobiles" db:"remark_mobiles"`
	AddWay                 int64  `json:"add_way" db:"add_way"`
	State                  string `json:"state" db:"state"`
	OperUserid             string `json:"oper_userid" db:"oper_userid"`
	RelationStatus         int64  `json:"relation_status" db:"relation_status"`
	MemberDeletedTime      int64  `json:"member_deleted_time" db:"member_deleted_time"`
	ExternalDeletedTime    int64  `json:"external_deleted_time" db:"external_deleted_time"`
	CreatedAt              int64  `json:"created_at" db:"created_at"`
	UpdatedAt              int64  `json:"updated_at" db:"updated_at"`
	DeletedAt              int64  `json:"deleted_at" db:"deleted_at"`
	Time                   string `json:"time" db:"time"`
	StartTime              int64  `json:"start_time" db:"start_time"`
	EndTime                int64  `json:"end_time" db:"end_time"`
	ChannelsName           string `json:"channels_name" db:"channels_name"`
	ChannelId              int64  `json:"channels_id" db:"channels_id"`
	ExternalDeletedTimeStr string `json:"external_deleted_time_str" db:"external_deleted_time"`
	MemberDeletedTimeStr   string `json:"member_deleted_time_str" db:"member_deleted_time"`
}

// ContactChannelsData 渠道列表数据
type ContactChannelsData struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Name         string `json:"name" db:"name"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// StatisticsItemResponse 统计相应结构体
type StatisticsItemResponse struct {
	TimeOrName              string `json:"time_or_name"`
	NewlyAddedCustomerCount int64  `json:"newly_added_customer_count"` // 净增客户数
	ExistRelationCustomer   int64  `json:"exist_relation_customer"`    // 留存好友数
	IncreaseCustomerCount   int64  `json:"increase_customer_count"`    // 新增好友总数
	LossCustomerCount       int64  `json:"loss_customer_count"`        // 流失好友数
	CustomerDelMember       int64  `json:"customer_del_member"`        // 好友删除员工数
	MemberDelCustomer       int64  `json:"member_del_customer"`        // 员工删除客户总数
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	Page  int64                     `json:"page"`
	Size  int64                     `json:"size"`
	Total int64                     `json:"total"`
	List  []*StatisticsItemResponse `json:"list"`
}

// StatisticsAllReq 统计所有的渠道的结构体
type StatisticsAllReq struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	StartTime    int64 `json:"start_time" db:"start_time"`
	EndTime      int64 `json:"end_time" db:"end_time"`
}

// // StatisticsChannelReq 统计单个渠道渠道的结构体
type StatisticsChannelReq struct {
	ChannelsState int64 `json:"channels_state" db:"channels_state"`
	StatisticsAllReq
}

// ChannelID 渠道
type ChannelID struct {
	ChannelID int64 `json:"channel_id" db:"channel_id" desc:"channel_id"`
}

type WelcomeContent struct {
	Id      int64  `json:"id" db:"id" desc:"id"`
	Content string `json:"content" db:"content" desc:"content"`
}

type WelcomesAllResp struct {
	welcome.DrainageWelcomesModel
	ContentText []map[string]interface{} `json:"content_text"`
	TypeName    string                   `json:"type_name"`
}

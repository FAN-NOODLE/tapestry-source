package chat

import (
	"encoding/json"

	"micro-scrm/internal/pkg/common"
	chat "micro-scrm/internal/pkg/sdk/microchat"
)

// 用户类型:员工
const UserTypeMember = 1

// 用户类型:机器人
const UserTypeRobot = 2

// 用户类型:企业微信外部联系人
const UserTypeWorkExternalUser = 3

// 用户类型:微信外部联系人
const UserTypeExternalUser = 4

// 用户类型:群聊
const UserTypeGroup = 5

// 同意状态:默认同意
const AgreeStatusDefaultAgree = 1

// 同意状态:同意
const AgreeStatusAgree = 2

// 同意状态:不同意
const AgreeStatusDisAgree = 3

// 群类型: 内部群
const GroupTypeInside = 2

const MemberStatusNormal = 1

// FileType 包含文件的消息类型
var FileType = []string{
	chat.MsgTypeImage,
	chat.MsgTypeVoice,
	chat.MsgTypeVideo,
	chat.MsgTypeEmotion,
	chat.MsgTypeFile,
	chat.MsgTypeMixed,
	chat.MsgTypeMeetingVoiceCall,
	chat.MsgTypeVoipDocShare,
	chat.MsgTypeChatrecord,
}

// Msg 聊天消息表
type Msg struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	MsgId        string `json:"msg_id" db:"msg_id" desc:"企微消息id"`
	Action       string `json:"action" db:"action" desc:"消息动作"`
	FromUser     string `json:"from_user" db:"from_user" desc:"发送用户"`
	FromUserType int64  `json:"from_user_type" db:"from_user_type" desc:"发送用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人"`
	ToList       string `json:"to_list" db:"to_list" desc:"消息接收用户列表"`
	RoomId       string `json:"room_id" db:"room_id" desc:"群聊消息的群ID"`
	MsgAt        int64  `json:"msg_at" db:"msg_at" desc:"消息发送时间"`
	MsgTime      int64  `json:"msg_time" db:"msg_time" desc:"消息发送时间ms"`
	Type         string `json:"type" db:"type" desc:"类型"`
	Content      string `json:"content" db:"content" desc:"内容"`
	Seq          int64  `json:"seq" db:"seq" desc:"企业微信指针"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`

	// 处理过的数据
	ToListArr    []string               `json:"to_list_arr" desc:"消息接收用户列表"`
	ContentMap   map[string]interface{} `json:"content_map" desc:"内容decode"`
	FromUserInfo UserInfo               `json:"from_user_info" desc:"发送人信息"`
	ToUserInfo   UserInfo               `json:"to_user_info" desc:"接收人信息"`
}

// MsgPageListRequest 分页列表请求
type MsgPageListRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	Page         int64  `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size         int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	FromUser     string `json:"from_user" db:"from_user" desc:"发送方,员工" jpath:"page" validate:"required"`
	ToUser       string `json:"to_user" db:"to_user" desc:"接收方,客户" jpath:"page" validate:"required"`
	Type         string `json:"type" db:"type" desc:"消息类型" jpath:"type" validate:""`
	Keywords     string `json:"keywords" desc:"搜索关键字" jpath:"page" validate:""`
	StartAt      int64  `json:"start_at" desc:"开始聊天时间" jpath:"start_at" validate:""`
	EndAt        int64  `json:"end_at" desc:"开始聊天时间" jpath:"end_at" validate:""`
}

// MsgPageListResponse 分页消息响应
type MsgPageListResponse struct {
	common.PageType
	List []Msg `json:"list" desc:"列表数据"`
}

// GroupMsgListRequest 群消息列表请求
type GroupMsgListRequest struct {
	EnterpriseID int64  `json:"enterprise_id" jpath:"enterprise_id" validate:"required"`
	GroupId      string `json:"group_id" jpath:"group_id" validate:"required"`
	Page         int64  `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size         int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	Type         string `json:"type" db:"type" desc:"消息类型" jpath:"type" validate:""`
	Keywords     string `json:"keywords" desc:"搜索关键字" jpath:"page" validate:""`
	StartAt      int64  `json:"start_at" desc:"开始聊天时间" jpath:"start_at" validate:""`
	EndAt        int64  `json:"end_at" desc:"开始聊天时间" jpath:"end_at" validate:""`
}

// MsgListByIDRequest 根据ID找消息分页列表请求
type MsgListByIDRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	FromUser     string `json:"from_user" db:"from_user" desc:"发送方,员工" jpath:"page" validate:"required"`
	ToUser       string `json:"to_user" db:"to_user" desc:"接收方,客户" jpath:"page" validate:"required"`
	Size         int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	MsgID        int64  `json:"msg_id" db:"id" validate:"required" desc:"从哪个消息ID开始找"`
	Sort         string `json:"sort" validate:"required" desc:"排序 ASC:往下找,DESC:往上找"`
	IsNeedSelf   bool   `json:"is_need_self" validate:"" desc:"是否需要消息本身"`
}

// MsgListByIDResponse 根据ID找消息分页列表响应
type MsgListByIDResponse struct {
	List []Msg `json:"list" desc:"列表数据"`
}

// MemberRelation 聊天员工关系
type MemberRelation struct {
	ID         int64 `json:"id" db:"id" desc:"ID"`
	MemberID   int64 `json:"member_id" db:"member_id" desc:"员工ID"`
	LastChatAt int64 `json:"last_chat_at" db:"last_chat_at" desc:"最后聊天时间"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MemberInfo 聊天员工结构体
type MemberInfo struct {
	MemberID    int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	Userid      string `json:"userid" db:"userid" desc:"userid"`
	Name        string `json:"name" db:"name" desc:"企业微信姓名"`
	Mobile      string `json:"mobile" db:"mobile" desc:"企业微信手机号"`
	Position    string `json:"position" db:"position" desc:"职务"`
	Gender      int64  `json:"gender" db:"gender" desc:"性别0=未知 1=男性 2=女性"`
	Avatar      string `json:"avatar" db:"avatar" desc:"头像"`
	ThumbAvatar string `json:"thumb_avatar" db:"thumb_avatar" desc:"缩略图头像"`
	Alias       string `json:"alias" db:"alias" desc:"别名"`
	Status      int64  `json:"status" db:"status" desc:"状态 0:刚创建,1:已激活,2:已禁用,4:未激活,5:退出企业"`
	DeletedAt   int64  `json:"deleted_at" db:"deleted_at" desc:"员工删除时间,>0说明离职了"`
	LastChatAt  int64  `json:"last_chat_at" db:"last_chat_at" desc:"最后一次聊天时间"`
}

// MemberPageListRequest 员工分页列表请求
type MemberPageListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	Page         int64 `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size         int64 `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	DepartmentID int64 `json:"department_id" jpath:"page" validate:""`
	MemberID     int64 `json:"member_id" jpath:"member_id" validate:""`
}

// MemberPageListResponse 员工分页列表响应
type MemberPageListResponse struct {
	common.PageType
	List []MemberInfo `json:"list" desc:"列表数据"`
}

// MemberToMemberPageListRequest 员工的员工列表请求
type MemberToMemberPageListRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	MemberID       int64  `json:"member_id" desc:"员工ID" jpath:"member_id" validate:"required"`
	Page           int64  `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size           int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	ToName         string `json:"to_name" desc:"搜索,员工姓名" jpath:"to_name" validate:""`
	IsReply        int64  `json:"is_reply" desc:"是否回复,1:是,2:否" jpath:"is_reply" validate:""`
	ToDepartmentID int64  `json:"to_department_id" jpath:"to_department_id" validate:""`
	ToMemberID     int64  `json:"to_member_id" jpath:"to_member_id" validate:""`
}

// MemberToMemberInfo 员工的员工
type MemberToMemberInfo struct {
	Name   string `json:"name" db:"name"`
	Avatar string `json:"avatar" db:"avatar"`
	Gender int64  `json:"gender" db:"gender"`
	Userid string `json:"userid" db:"userid"`
	Status int64  `json:"status" db:"status"`
}

// MemberToMemberPageListResponse 员工的员工列表响应
type MemberToMemberPageListResponse struct {
	common.PageType
	List []MemberToMemberInfo `json:"list" desc:"列表数据"`
}

// MemberToExternal 员工的外部联系人
type MemberToExternalInfo struct {
	Name        string `json:"name" db:"name"`
	Avatar      string `json:"avatar" db:"avatar"`
	Gender      int64  `json:"gender" db:"gender"`
	Type        int64  `json:"type" db:"type"`
	Userid      string `json:"userid" db:"external_userid"`
	AgreeStatus int64  `json:"agree_status" db:"agree_status"`
	IsReply     int64  `json:"is_reply" db:"is_reply"`
}

// MemberToExternalPageListRequest 员工与外部联系人聊天列表请求
type MemberToExternalPageListRequest struct {
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	MemberID     int64  `json:"member_id" desc:"员工ID" jpath:"member_id" validate:"required"`
	Page         int64  `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size         int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	Name         string `json:"name" desc:"搜索,客户姓名" jpath:"name" validate:""`
	AgreeStatus  int64  `json:"agree_status" desc:"同意情况 1:默认同意 2:同意 3:不同意" jpath:"agree_status" validate:""`
	IsReply      int64  `json:"is_reply" desc:"是否回复,1:是,2:否" jpath:"is_reply" validate:""`
}

// MemberToExternalPageListResponse 员工与外部联系人聊天列表响应
type MemberToExternalPageListResponse struct {
	common.PageType
	List []MemberToExternalInfo `json:"list" desc:"列表数据"`
}

// MemberToGroupInfo 员工的群
type MemberToGroupInfo struct {
	ID            int64  `json:"id" db:"id" desc:"主键id"`
	ChatID        string `json:"chat_id" db:"chat_id" desc:"企业微信群聊id"`
	Name          string `json:"name" db:"name" desc:"群名"`
	OwnerMemberID int64  `json:"owner_member_id" db:"owner_member_id" desc:"群主"`
	IsOwner       int8   `json:"is_owner" desc:"是否是群主 1 是 0 否"`
}

// MemberToGroupPageListRequest 员工与群聊天分页列表请求
type MemberToGroupPageListRequest struct {
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID" jpath:"enterprise_id" validate:"required"`
	MemberID     int64  `json:"member_id" desc:"员工ID" jpath:"member_id" validate:"required"`
	Page         int64  `json:"page" desc:"页数" jpath:"page" validate:"required,min=0"`
	Size         int64  `json:"size" desc:"每页记录数" jpath:"page" validate:"required,min=0,max=100"`
	Name         string `json:"name" desc:"搜索,群聊名称" jpath:"name" validate:""`
}

// MemberGroupPageListResponse 员工的客户列表响应
type MemberToGroupPageListResponse struct {
	common.PageType
	List []MemberToGroupInfo `json:"list" desc:"列表数据"`
}

// ReceiveUser 消息接收人表
type ReceiveUser struct {
	ID          int64  `json:"id" db:"id"`
	ChatMsgId   int64  `json:"chat_msg_id" db:"chat_msg_id" desc:"消息表ID"`
	ToUser      string `json:"to_user" db:"to_user" desc:"接收用户企微userid"`
	UserChannel int64  `json:"user_channel" db:"user_channel" desc:"用户渠道 1:同企业,2:机器人,3:企业微信外部联系人,4:微信外部联系人"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// AgreeInfoRequest 会话存档同意信息
type AgreeInfoRequest struct {
	EnterpriseID int64                      `json:"enterprise_id"`
	List         []AgreeInfoRequestListInfo `json:"list" jpath:"list" validate:"required"`
}

// AgreeInfoRequestListInfo 会话存档同意信息列表信息
type AgreeInfoRequestListInfo struct {
	Userid          string `json:"userid" validate:"required"`
	ExteranalUserid string `json:"exteranal_userid" validate:"required"`
}

// Object 聊天对象结构
type Object struct {
	ID            int64  `json:"id" db:"id" desc:"ID"`
	FromUserid    string `json:"from_userid" db:"from_userid" desc:"发送方userid"`
	ToUserid      string `json:"to_userid" db:"to_userid" desc:"接收方userid"`
	FromObjectID  int64  `json:"from_object_id" db:"from_object_id" desc:"发送对象ID"`
	ToObjectID    int64  `json:"to_object_id" db:"to_object_id" desc:"接收对象ID"`
	FromChannel   int64  `json:"from_channel" db:"from_channel" desc:"发送用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人"`
	ToChannel     int64  `json:"to_channel" db:"to_channel" desc:"接收用户渠道"`
	LastChatAt    int64  `json:"last_chat_at" db:"last_chat_at" desc:"最后聊天时间"`
	AgreeStatus   int64  `json:"agree_status" db:"agree_status" desc:"同意状态 1:默认同意 2:同意 3:不同意"`
	AgreeStatusAt int64  `json:"agree_status_at" db:"agree_status_at" desc:"同意状态变更时间"`
	IsReply       int64  `json:"is_reply" db:"is_reply" desc:"是否回复,1:是2:否"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MemberToMember 员工与员工聊天表
type MemberToMember struct {
	ID           int64 `json:"id" db:"id"`
	FromMemberID int64 `json:"from_member_id" db:"from_member_id" desc:"发送人ID"`
	ToMemberID   int64 `json:"to_member_id" db:"to_member_id" desc:"接收人ID"`
	LastChatAt   int64 `json:"last_chat_at" db:"last_chat_at" desc:"最后聊天时间"`
	IsReply      int64 `json:"is_reply" db:"id_reply" desc:"是否回复,1:是2:否"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MemberToExternal 员工与外部联系人聊天表
type MemberToExternal struct {
	ID                 int64  `json:"id" db:"id"`
	Userid             string `json:"userid" db:"userid" desc:"员工userid"`
	MemberID           int64  `json:"member_id" db:"member_id" desc:"members:ID"`
	ExternalUserid     string `json:"external_userid" db:"external_userid" desc:"外部联系人userid"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id" desc:"外部联系人表内部:ID"`
	AgreeStatus        int64  `json:"agree_status" db:"agree_status" desc:"同意情况 1:默认同意 2:同意 3:不同意"`
	StatusChangeTime   int64  `json:"status_change_time" db:"status_change_time" desc:"状态修改时间"`
	LastChatAt         int64  `json:"last_chat_at" db:"last_chat_at" desc:"最后聊天时间"`
	IsReply            int64  `json:"is_reply" db:"id_reply" desc:"是否回复,1:是2:否"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MsgMedia 聊天记录关联媒体素材表
type MsgMedia struct {
	ID          int64  `json:"id" db:"id"`
	ChatMsgID   int64  `json:"chat_msg_id" db:"chat_msg_id" desc:"消息表ID"`
	Suffix      string `json:"suffix" desc:"文件后缀"`
	Name        string `json:"name" desc:"文件名称"`
	SdkFileID   string `json:"sdk_file_id" db:"sdk_file_id" desc:"资源文件ID"`
	Path        string `json:"path" db:"path" desc:"资源路径"`
	UserChannel int64  `json:"user_channel" db:"user_channel" desc:"用户渠道 1:同企业,2:机器人,3:企业微信外部联系人,4:微信外部联系人"`

	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// UserInfo 聊天用户信息
type UserInfo struct {
	ObjectID     int64  `json:"object_id" db:"object_id" desc:"员工ID,外部联系人ID"`
	UserID       string `json:"user_id" db:"userid" desc:"企微userid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Name         string `json:"name" db:"name" desc:"姓名"`
	Avatar       string `json:"avatar" db:"avatar" desc:"头像"`
	Remark       string `json:"remark" desc:"备注"`
	Status       int32  `json:"status" db:"status" desc:"备注"`
}

// GroupMemberInfo 群聊用户信息
type GroupMemberInfo struct {
	ChatID       string `json:"chat_id" db:"chat_id" desc:"企微群聊id"`
	Name         string `json:"name" db:"name" desc:"企微群聊名称"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	GroupChatID  int64  `json:"group_chat_id" db:"group_chat_id" desc:"群聊表id"`
	UserID       string `json:"userid" db:"userid" desc:"成员userid"`
	MemberID     int64  `json:"member_id" db:"member_id" desc:"成员表id"`
}

// MixMediaContent 混合媒体消息
type MixMediaContent struct {
	Items []*MixMediaItem `json:"item"`
}

// MixMediaItem 混合媒体消息
type MixMediaItem struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// ContentEntity 非混合消息
type ContentEntity struct {
	SdkFileid   string `json:"sdkfileid"`
	FileExt     string `json:"fileext"`
	EmotionType int32  `json:"type"`
	FileName    string `json:"filename"`
	FileType    string `json:"file_type"`
}

// 接收用户搜索返回值
type SearchReceiveResponse struct {
	ChatMsgId int64 `json:"chat_msg_id" db:"chat_msg_id" desc:"消息id"`
}

// WorkGroupChat 客户群信息
type WorkGroupChat struct {
	ID               int64  `json:"id" db:"id" desc:"主键"`
	EnterpriseID     int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Type             int    `json:"type" db:"type" desc:"群类型"`
	ChatID           string `json:"chat_id" db:"chat_id" desc:"群ID"`
	Name             string `json:"name" db:"name" desc:"群名称"`
	OwnerMemberID    int64  `json:"owner_member_id" db:"owner_member_id" desc:"群主ID"`
	Notice           string `json:"notice" db:"notice" desc:"群公告"`
	CreateTime       int64  `json:"create_time" db:"create_time" desc:"群创建时间"`
	Status           int    `json:"status" db:"status" desc:"群状态 1=群存在 2=群已解散"`
	DissolveTime     int64  `json:"dissolve_time" db:"dissolve_time" desc:"群解散时间"`
	CreatedAt        int64  `json:"created_at" db:"created_at"`
	UpdatedAt        int64  `json:"updated_at" db:"updated_at"`
	DeletedAt        int64  `json:"deleted_at" db:"deleted_at"`
	Owner            string `json:"owner"`
	MemberNum        int    `json:"member_num" db:"member_num"`
	ChatTarget       string `json:"chat_target" db:"chat_target" desc:"群目标"`
	ChannelID        int64  `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	ChatTypeID       int64  `json:"chat_type_id" db:"chat_type_id" desc:"群类型ID"`
	OperatorMemberID string `json:"operator_member_id" db:"operator_member_id" desc:"运营人员userid"`
	TmpAdmins        []byte `json:"-" db:"admin_list"`
}
type WorkGroupChatBase struct {
	WorkGroupChat
	ChannelName         string   `json:"channel_name" db:"channel_name" desc:"渠道名"`
	ChatTypeName        string   `json:"chat_type_name" db:"chat_type_name" desc:"群类型名称"`
	Owner               string   `json:"owner" desc:"群主"`
	UserID              string   `json:"user_id" db:"-" desc:"群主user_id"`
	Operator            string   `json:"operator" db:"operator" desc:"运营人员名称"`
	OperatorDepartments []string `json:"operator_departments" db:"-"`
	Admins              []string `json:"admins" db:"-" desc:"管理员列表"`
}

type WorkGroupChatReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ChatID       string `json:"chat_id" db:"chat_id" desc:"群ID"`
}

// WorkGroupChatMember 群成员
type WorkGroupChatMember struct {
	ID           int64  `json:"id" db:"id" desc:"主键"`
	GroupChatID  int64  `json:"group_chat_id" db:"group_chat_id" desc:"work_group_chat 主表ID"`
	Type         int    `json:"type" db:"type" desc:"成员类型:1 - 企业成员,2 - 外部联系人"`
	UserID       string `json:"userid" db:"userid" desc:"群成员id"`
	JoinTime     int64  `json:"join_time" db:"join_time" desc:"加群时间"`
	JoinScene    int    `json:"join_scene" db:"join_scene" desc:"入群方式"`
	LeaveTime    int64  `json:"leave_time" db:"leave_time" desc:"退群时间"`
	Status       int    `json:"status" db:"status" desc:"群状态1=入群 2=退群"`
	MemberID     int64  `json:"member_id" db:"member_id" desc:"成员ID"`
	UnionID      string `json:"unionid" db:"unionid" desc:"微信unionid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// GetSdkFileId 获取消息中的所有媒体id
func (msg *Msg) GetSdkFileId() []string {
	var sdkFileIdArr []string

	// 混合消息单独处理
	if msg.Type == chat.MsgTypeMixed {
		var mixedStruct chat.Mixed
		_ = json.Unmarshal([]byte(msg.Content), &mixedStruct)

		for _, v := range mixedStruct.Item {
			var itemContentMap map[string]interface{}
			_ = json.Unmarshal([]byte(v.Content), &itemContentMap)
			if sdkFileid, ok := itemContentMap["sdkfileid"]; ok {
				sdkFileIdArr = append(sdkFileIdArr, sdkFileid.(string))
			}
		}
	}

	// 观察微信文档基本直接转map取sdkfileid字段即可
	if sdkfileid, ok := msg.ContentMap["sdkfileid"]; ok {
		sdkFileIdArr = append(sdkFileIdArr, sdkfileid.(string))
	}

	return sdkFileIdArr
}

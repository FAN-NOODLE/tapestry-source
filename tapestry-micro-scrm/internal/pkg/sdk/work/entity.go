package work

import (
	"encoding/json"
	"micro-scrm/internal/pkg/sdk"
	"mime/multipart"
	"os"
	"unicode/utf8"
)

// QWExternalIdentifier 获取外部联系人固定使用的应用标识
const QWExternalIdentifier = "ExternalContact" // 外部联系人标识
const QWArchiveIdentifier = "scrm-session-archiving"
const QWCommunalContact = "communalContact"                // 通讯录标识
const QWCrsStoreAssistantSelf = "crs_store_assistant_self" // 小程序标识
const QWVM = "vm"                                          // VM陈列应用标识

const DrainageContactsTypeMultiple = 2 //默认多人

const DrainageContactsSceneQrcode = 2 //默认二维码联系

// QWGadgetIdentifier 侧边栏应用标识符
const QWGadgetIdentifier = "crs-manu-customer"

const (
	ExternalTagTypeOfficial = 1 // ExternalTagTypeOfficial 企业标签
	ExternalTagTypeCustomer = 2 // ExternalTagTypeCustomer 自定义标签
)

// STokenResponse 获取token返回结果
type STokenResponse struct {
	Data STokenResponseData `json:"data"`
}

// STokenResponseData token 里的data
type STokenResponseData struct {
	AccessToken string `json:"access_token"`
	Expire      int64  `json:"expire"`
}

// STokenRequest 获取token 请求
type STokenRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// QyWorkImpl 企业微信
type QyWorkImpl struct {
	domain         string
	workHttpClient *sdk.Client
}

// QyGwRequest 企业微信
type QyGwRequest struct {
	ReqPL     string `json:"reqPL"`
	ReqMethod string `json:"reqMethod"`
	ReqURL    string `json:"reqUrl"`
	Token     string `json:"token"`
}

// QyDecryptRequest 回调参数解密
type QyDecryptRequest struct {
	Timestamp     string `json:"timestamp"`
	Nonce         string `json:"nonce"`
	MsgSignature  string `json:"msgSignature"`
	EncryptString string `json:"encryptString"`
	Token         string `json:"token"`
	PlIdentifier  string `json:"plIdentifier"`
	CorpID        string `json:"corpid"`
}

// MemberExternalCallback 外部联系人回调解密数据
type MemberExternalCallback struct {
	ID              int64  `json:"id" db:"id"`
	ToUserName      string `json:"toUserName" db:"to_user_name" desc:"企业微信CorpID"`
	FromUserName    string `json:"fromUserName" db:"from_user_name" desc:"此事件该值固定为sys，表示该消息由系统生成"`
	CreateTime      int64  `json:"createTime" db:"create_time" desc:"消息创建时间 （整型）"`
	MsgType         string `json:"msgType" db:"msg_type" desc:"消息的类型，此时固定为event"`
	Event           string `json:"event" db:"event" desc:"事件的类型"`
	ChangeType      string `json:"changeType" db:"change_type" desc:"此时固定为add_external_contact"`
	UserID          string `json:"userID" db:"user_id" desc:"企业服务人员的UserID"`
	ExternalUserID  string `json:"externalUserID" db:"external_user_id" desc:"外部联系人的userid，注意不是企业成员的帐号"`
	State           string `json:"state" db:"state" desc:"添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道"`
	WelcomeCode     string `json:"welcomeCode" db:"welcome_code" desc:"欢迎语code，可用于发送欢迎语"`
	FailReason      string `json:"failReason" db:"fail_reason" desc:"接替失败的原因, customer_refused-客户拒绝， customer_limit_exceed-接替成员的客户数达到上限"`
	ChatID          string `json:"chatId" db:"chat_id" desc:"群ID"`
	Status          int    `json:"status" db:"status"`
	StatusTrace     int    `json:"status_trace" db:"status_trace"`
	ChannelID       int64  `json:"channel_id" db:"channel_id"`
	AddRelationTime int64  `json:"add_relation_time" db:"add_relation_time"`
}

type LossCustomerAddRelationTimeChannelID struct {
	ChannelID       int64 `json:"channel_id" db:"channel_id"`
	AddRelationTime int64 `json:"add_relation_time" db:"add_relation_time"`
}

// WorkDepartment 组织
type WorkDepartment struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	Parentid int64  `json:"parentid"`
	Order    int32  `json:"order"`
}

// QyResponse 企业微信返回结果
type QyResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// WorkDepartmentListRequest 组织列表
type WorkDepartmentListRequest struct {
	QyGwRequest
	ID interface{} `json:"id"`
}

// WorkMemberRequest 请求成员
type WorkMemberRequest struct {
	QyGwRequest
	DepartmentID int64 `json:"department_id" desc:"获取的部门id"`
	FetchChild   int   `json:"fetch_child" desc:"是否递归获取子部门下面的成员：1-递归获取，0-只获取本部门"`
}

// AddTagUsersRequest 添加标签成员请求
type AddTagUsersRequest struct {
	QyGwRequest
	TagID    int      `json:"tagid"`
	UserList []string `json:"userlist"`
}

// MemberItem 成员对象
type MemberItem struct {
	Userid     string  `json:"userid" desc:"userid"`
	Name       string  `json:"name" desc:"成员姓名"`
	Department []int64 `json:"department"`
	OpenUserid string  `json:"open_userid"`
}

// MemberDetailRequest 成员详情请求参数
type MemberDetailRequest struct {
	QyGwRequest
	Userid string `json:"userid"`
}

// MemberDetail 成员详情
type MemberDetail struct {
	Userid         string  `json:"userid"`
	Name           string  `json:"name"`
	Department     []int64 `json:"department"`
	Position       string  `json:"position"`
	Mobile         string  `json:"mobile"`
	Gender         string  `json:"gender"`
	Email          string  `json:"email"`
	IsLeaderInDept []int8  `json:"is_leader_in_dept"`
	Avatar         string  `json:"avatar"`
	ThumbAvatar    string  `json:"thumb_avatar"`
	Telephone      string  `json:"telephone"`
	Alias          string  `json:"alias"`
	Address        string  `json:"address"`
	OpenUserid     string  `json:"open_userid"`
	MainDepartment int32   `json:"main_department"`
	Status         int8    `json:"status"`
	QrCode         string  `json:"qr_code"`
	Extattr        Extattr `json:"extattr"`
}

type Attrs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Extattr struct {
	Attrs []Attrs `json:"attrs"`
}

// ContentMessage 消息内容
type MessageTypeTextCardContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

// TextData 文本消息content
type TextData struct {
	Content string `json:"content"`
}

type MiniprogramNoticeItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MiniprogramNoticeContent struct {
	Appid       string                  `json:"appid"`
	Page        string                  `json:"page"`
	Title       string                  `json:"title"`
	Description string                  `json:"description,omitempty"`
	ContentItem []MiniprogramNoticeItem `json:"content_item"`
}

// SendMessageRequest 发送消息
type SendMessageRequest struct {
	QyGwRequest
	Touser                 string                      `json:"touser,omitempty"`
	Toparty                string                      `json:"toparty,omitempty"`
	Totag                  string                      `json:"totag,omitempty"`
	Msgtype                string                      `json:"msgtype"`
	Agentid                int64                       `json:"agentid,omitempty"`
	EnableIdTrans          int                         `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int                         `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int                         `json:"duplicate_check_interval,omitempty"`
	TextCard               *MessageTypeTextCardContent `json:"textcard,omitempty"`
	Text                   *TextData                   `json:"text,omitempty"`
	MiniprogramNotice      *MiniprogramNoticeContent   `json:"miniprogram_notice,omitempty"`
}

type SendMessageResponse struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
	Msgid        string `json:"msgid"`
}

// ExternalContactRequest 外部联系人详情
type ExternalContactRequest struct {
	QyGwRequest
	ExternalUserID string `json:"external_userid"`
}

// ExternalContactResponse 外部联系人详情返回参数
type ExternalContactResponse struct {
	Base       ExternalContactBase  `json:"external_contact"`
	FollowUser []ExternalFollowUser `json:"follow_user"`
}

// BatchExternalContactRequest 批量获取外部联系人详情
type BatchExternalContactRequest struct {
	QyGwRequest
	MemberUserID string `json:"userid"`
	Cursor       string `json:"cursor"`
	Limit        int    `json:"limit"`
}

type WorkExternalContact struct {
	Base       ExternalContactBase `json:"external_contact"`
	FollowInfo ExternalFollowUser  `json:"follow_info"`
}

type BatchExternalContactResponse struct {
	QyResponse
	ExternalContactList []WorkExternalContact `json:"external_contact_list"`
	NextCursor          string                `json:"next_cursor"`
}

// 批量获取客户详情信息
// 注意：批量获取外部联系人时，tag_id，和获取客户详情的结构不一样
type BatchExternalFollowUser struct {
	UserID         string   `json:"userid" desc:"添加了此外部联系人的企业成员userid"`
	Remark         string   `json:"remark" desc:"该成员对此外部联系人的备注"`
	Description    string   `json:"description" desc:"该成员对此外部联系人的描述"`
	CreateTime     int64    `json:"createtime" desc:"该成员添加此外部联系人的时间"`
	TagIDList      []string `json:"tag_id"`
	RemarkCorpName string   `json:"remark_corp_name" desc:"该成员对此客户备注的企业名称"`
	RemarkMobiles  []string `json:"remark_mobiles" desc:"	该成员对此客户备注的手机号码"`
	AddWay         int      `json:"add_way" desc:"该成员添加此客户的来源，具体含义详见"`
	OperUserID     string   `json:"oper_userid" desc:"发起添加的userid，如果成员主动添加，为成员的userid；如果是客户主动添加，则为客户的外部联系人userid；如果是内部成员共享/管理员分配，则为对应的成员/管理员userid"`
}

// ExternalContactBase 外部联系人主信息
type ExternalContactBase struct {
	ExternalUserID  string          `json:"external_userid" desc:"外部联系人的userid"`
	Name            string          `json:"name" desc:"外部联系人的名称"`
	Avatar          string          `json:"avatar" desc:"外部联系人头像"`
	Type            int8            `json:"type" desc:"外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户"`
	Gender          int8            `json:"gender" desc:"外部联系人性别 0-未知 1-男性 2-女性"`
	UnionID         string          `json:"unionid" desc:"外部联系人在微信开放平台的唯一身份标识"`
	Position        string          `json:"position" desc:"外部联系人的职位"`
	CorpName        string          `json:"corp_name" desc:"外部联系人所在企业的简称"`
	CorpFullName    string          `json:"corp_full_name" desc:"外部联系人所在企业的主体名称"`
	ExternalProfile json.RawMessage `json:"external_profile" desc:"外部联系人的自定义展示信息"`
}

// ExternalFollowUser 外部联系人添加的成员
type ExternalFollowUser struct {
	UserID         string        `json:"userid" desc:"添加了此外部联系人的企业成员userid"`
	Remark         string        `json:"remark" desc:"该成员对此外部联系人的备注"`
	Description    string        `json:"description" desc:"该成员对此外部联系人的描述"`
	CreateTime     int64         `json:"createtime" desc:"该成员添加此外部联系人的时间"`
	Tags           []ExternalTag `json:"tags"`
	RemarkCorpName string        `json:"remark_corp_name" desc:"该成员对此客户备注的企业名称"`
	RemarkMobiles  []string      `json:"remark_mobiles" desc:"	该成员对此客户备注的手机号码"`
	AddWay         int           `json:"add_way" desc:"该成员添加此客户的来源，具体含义详见"`
	OperUserID     string        `json:"oper_user_id" desc:"发起添加的userid，如果成员主动添加，为成员的userid；如果是客户主动添加，则为客户的外部联系人userid；如果是内部成员共享/管理员分配，则为对应的成员/管理员userid"`
	State          string        `json:"state" desc:"企业自定义的state参数，用于区分客户具体是通过哪个「联系我」添加"`
}

// ExternalTag 外部联系人打的标签
type ExternalTag struct {
	TagID     string `json:"tag_id" desc:"该成员添加此外部联系人所打企业标签的id,仅企业设置（type为1）的标签返回"`
	GroupName string `json:"group_name" desc:"该成员添加此外部联系人所打标签的分组名称"`
	TagName   string `json:"tag_name" desc:"该成员添加此外部联系人所打标签名称"`
	Type      int8   `json:"type" desc:"该成员添加此外部联系人所打标签类型, 1-企业设置, 2-用户自定义"`
}

// ExternalGroupRequest 获取客户群详情
type ExternalGroupRequest struct {
	QyGwRequest
	ChatID   string `json:"chat_id" desc:"客户群ID"`
	NeedName int    `json:"need_name" desc:"是否需要返回群成员的名字group_chat.member_list.name。0-不返回；1-返回。默认不返回"`
}

// ExternalGroupResponse 客户群返回结果
type ExternalGroupResponse struct {
	ChatID     string                    `json:"chat_id"`
	Name       string                    `json:"name"`
	Owner      string                    `json:"owner"`
	CreateTime int64                     `json:"create_time"`
	Notice     string                    `json:"notice"`
	MemberList []ExternalGroupMemberList `json:"member_list"`
	AdminList  []ExternalGroupAdminList  `json:"admin_list"`
}

// ExternalGroupMemberList 客户群成员列表
type ExternalGroupMemberList struct {
	UserID        string `json:"userid"`
	GroupNickname string `json:"group_nickname"`
	Name          string `json:"name"`
	Type          int    `json:"type"`
	UnionID       string `json:"unionid"`
	JoinTime      int64  `json:"join_time"`
	JoinScene     int    `json:"join_scene"`
	Invitor       struct {
		UserID string `json:"userid"`
	} `json:"invitor"`
}

type ExternalGroupAdminList struct {
	Userid string `json:"userid"`
}

// EditExternalDescriptionRequest 修改外部联系人描述请求参数
type EditExternalDescriptionRequest struct {
	QyGwRequest
	UserID         string `json:"userid"`
	ExternalUserid string `json:"external_userid"`
	Remark         string `json:"remark"`
	Description    string `json:"description"`
}

// EditExternalDescriptionResponse 修改外部联系人描述返回值
type EditExternalDescriptionResponse struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

// 获取企业微信请求参数
type GetWxCodeRequest struct {
	QyGwRequest
	Type       int      `json:"type" db:"type"`
	Scene      int      `json:"scene" db:"scene"`
	User       []string `json:"user" db:"user"`
	Party      []int    `json:"party"`
	State      string   `json:"state" db:"state"`
	SkipVerify bool     `json:"skip_verify" db:"skip_verify"`
}

// 企业微信返回参数
type GetWxCodeResult struct {
	Errcode  int64  `json:"errcode" db:"errcode"`
	Errmsg   string `json:"errmsg" db:"errmsg"`
	ConfigId string `json:"config_id" db:"config_id"`
	QrCode   string `json:"qr_code" db:"qr_code"`
}

// 添加联系我二维码
type ContactCreateRequest struct {
	ID           int64    `json:"id" db:"id" desc:"id"`
	CodeCame     string   `json:"code_name" binding:"required"  db:"code_name" desc:"二维码名称"`
	ChannelId    int64    `json:"channel_id" binding:"required"  db:"channel_id" desc:"渠道id"`
	AutoFriend   int      `json:"auto_friend" binding:"required"  db:"auto_friend" desc:"是否自动开启"`
	MemberIds    []int64  `json:"member_ids" binding:"required"  db:"member_ids" desc:"员工id"`
	ThirdIds     []int    `json:"third_ids" desc:"部门Id"`
	EnterpriseId int64    `json:"enterprise_id" binding:"required" db:"enterprise_id" desc:"企业ID"`
	Type         int      `json:"type"  db:"type" desc:"类型"`
	UserIds      []string `json:"user_ids" binding:"required" db:"user_ids" desc:"user_ids"`
	State        string   `json:"state"`
}

// 修改联系我二维码
type UpdateWxCodeRequest struct {
	QyGwRequest
	User       []string `json:"user" db:"user"`
	Party      []int    `json:"party"`
	ConfigId   string   `json:"config_id" db:"config_id"`
	State      string   `json:"state" db:"state"`
	SkipVerify bool     `json:"skip_verify" db:"skip_verify"`
}

// CacheMediaRequest 缓存微信图片素材请求
type CacheMediaWorkRequest struct {
	QyGwRequest
	File         *multipart.File `json:"file"`
	FileOrg      string          `json:"fileOrg"`
	FileType     string          `json:"fileType"`
	FileName     string          `json:"filename"`
	EnterpriseID int64           `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ImgUrl       string          `json:"img_url"`
}

// CacheMediaResponse 缓存微信图片素材响应
type CacheMediaWorkResponse struct {
	Code string               `json:"code"`
	Msg  string               `json:"msg"`
	Data CacheMediaWorkStruct `json:"data"`
}

// CacheMediaWorkStruct 缓存微信图片素材响应子结构体
type CacheMediaWorkStruct struct {
	Type         string `json:"type"`
	MediaId      string `json:"mediaId"`
	ThumbMediaId string `json:"thumbMediaId"`
}

// SendWelcomeRequest 发送欢迎语请求
type SendWelcomeRequest struct {
	WelcomeCode string `json:"welcome_code"`
	QyGwRequest
	Text        Text                     `json:"text"`
	Attachments []map[string]interface{} `json:"attachments"`
}

// SendWelcomeMsgTextImage 发送文本和图片类欢迎语
type SendWelcomeMsgTextImage struct {
	QyGwRequest
	Text        `json:"text"`
	Image       `json:"image"`
	WelcomeCode string `json:"welcome_code"`
}

// SendWelcomeMsgTextLink 发送文本和链接类欢迎语
type SendWelcomeMsgTextLink struct {
	QyGwRequest
	Text        `json:"text"`
	Link        `json:"link"`
	WelcomeCode string `json:"welcome_code"`
}

// SendWelcomeMsgText 发送文本和链接类欢迎语
type SendWelcomeMsgText struct {
	QyGwRequest
	Text        `json:"text"`
	WelcomeCode string `json:"welcome_code"`
}

// SendWelcomeMsgTextMiniProgram 发送文本和小程序类欢迎语
type SendWelcomeMsgTextMiniProgram struct {
	QyGwRequest
	Text        `json:"text"`
	MiniProgram `json:"miniprogram"`
	WelcomeCode string `json:"welcome_code"`
}

// 视频
type Video struct {
	MediaId  string `json:"media_id"`
	CoverUrl string `json:"cover_url"`
	Url      string `json:"url"`
}

// 文件
type File struct {
	MediaId  string `json:"media_id"`
	CoverUrl string `json:"cover_url"`

	PicUrl string `json:"pic_url"`
	Title  string `json:"title"`
}

// Text 文本
type Text struct {
	Content string `json:"content"`
}

// Image 图片
type Image struct {
	MediaId string `json:"media_id"`
	PicUrl  string `json:"pic_url"`
}

// Link 链接
type Link struct {
	Title    string `json:"title"`
	WxPicUrl string `json:"picurl"`
	Desc     string `json:"desc"`
	Url      string `json:"url"`
	PicUrl   string `json:"pic_url"`
}

const (
	qyLinkTitleMaxBytes = 127
	qyLinkDescMaxBytes  = 511
)

func truncateUTF8ByBytes(value string, maxBytes int) string {
	if len(value) <= maxBytes {
		return value
	}
	end := maxBytes
	for end > 0 && !utf8.RuneStart(value[end]) {
		end--
	}
	return value[:end]
}

// LimitLinkForQy 在发送企微请求前限制链接标题和描述的 UTF-8 字节数。
func LimitLinkForQy(l Link) Link {
	l.Title = truncateUTF8ByBytes(l.Title, qyLinkTitleMaxBytes)
	l.Desc = truncateUTF8ByBytes(l.Desc, qyLinkDescMaxBytes)
	return l
}

// MiniProgram 小程序
type MiniProgram struct {
	Title      string `json:"title"`
	PicMediaId string `json:"pic_media_id"`
	Appid      string `json:"appid"`
	Page       string `json:"page"`
	CoverUrl   string `json:"cover_url"`
}

// UserTagHandlerRequest 用户标签处理请求
type UserTagHandlerRequest struct {
	QyGwRequest
	Userid         string   `json:"userid" desc:"添加外部联系人的userid"`
	ExternalUserid string   `json:"external_userid" desc:"外部联系人userid"`
	AddTag         []string `json:"add_tag" desc:"要标记的标签列表"`
	RemoveTag      []string `json:"remove_tag" desc:"要移除的标签列表"`
}

// CheckSingleAgreeRequest 获取会话同意情况
type CheckSingleAgreeRequest struct {
	QyGwRequest
	Info []CheckSingleAgreeInfoRequest `json:"info"`
}

// CheckSingleAgreeInfoRequest 检查会话同意信息请求
type CheckSingleAgreeInfoRequest struct {
	Userid          string `json:"userid"`
	Exteranalopenid string `json:"exteranalopenid"`
}

// ExternalGroupResponse 会话同意情况结果
type CheckSingleAgreeResponse struct {
	StatusChangeTime int64  `json:"status_change_time" desc:"同意状态改变的具体时间"`
	UserID           string `json:"userid" desc:"员工userid"`
	ExteranalOpenID  string `json:"exteranalopenid" desc:"外部联系人userid"`
	AgreeStatus      string `json:"agree_status" desc:"同意状态 同意:Agree,不同意:Disagree,默认同意:Default_Agree"`
}

// GetCorpTagListRequest 获取企业标签数据请求参数
type GetCorpTagListRequest struct {
	QyGwRequest
	TagID []string `json:"tag_id"`
}

// GetCorpTagListResponse 获取企业标签数据返回数据
type GetCorpTagListResponse struct {
	Code     int                   `json:"errcode"`
	Msg      string                `json:"errmsg"`
	TagGroup []*CorpTagGroupEntity `json:"tag_group"`
}

// CorpTagGroupEntity 企业标签组实体
type CorpTagGroupEntity struct {
	GroupID    string           `json:"group_id"`
	GroupName  string           `json:"group_name"`
	CreateTime int              `json:"create_time"`
	Order      int              `json:"order"`
	Delete     bool             `json:"delete"`
	Tag        []*CorpTagEntity `json:"tag"`
}

// CorpTagEntity 企业标签实体
type CorpTagEntity struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreateTime int    `json:"create_time"`
	Order      int    `json:"order"`
	Delete     bool   `json:"delete"`
}

// RobotInfoResponse 机器人信息
type RobotInfoResponse struct {
	Code int       `json:"errcode"`
	Msg  string    `json:"errmsg"`
	Data RobotInfo `json:"data"`
}

// RobotInfo 机器人信息
type RobotInfo struct {
	RobotId       string `json:"robot_id"`
	Name          string `json:"name"`
	CreatorUserID string `json:"creator_userid"`
}

// RobotInfoRequest 机器人信息
type RobotInfoRequest struct {
	QyGwRequest
	RobotID string `json:"robot_id"`
}

// InsideGroupRequest 获取内部群信息请求体
type InsideGroupRequest struct {
	QyGwRequest
	RoomID string `json:"roomid"`
}

// GroupChatInfoResponse 获取内部群信息返回结果
type InsideGroupRequestResponse struct {
	RoomName       string         `json:"roomname" desc:"roomid对应的群名称"`
	Creator        string         `json:"creator" desc:"roomid对应的群创建者，userid"`
	RoomCreateTime int64          `json:"room_create_time" desc:"roomid对应的群创建时间"`
	Notice         string         `json:"notice" desc:"roomid对应的群公告"`
	MemberList     []GroupMembers `json:"members" desc:"roomid对应的群成员列表"`
}

// GroupMembers  群成员
type GroupMembers struct {
	MemberID string `json:"memberid" desc:"roomid群成员的id，userid"`
	JoinTime int64  `json:"jointime" desc:"roomid群成员的入群时间"`
}

// QyMassHairRequest 获取企业标签数据请求参数
type QyMassHairRequest struct {
	QyGwRequest
	ChatType       string      `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string    `json:"external_userid" desc:"客户的外部联系人id列表"`
	Sender         string      `json:"sender" desc:"发送企业群发消息的成员userid"`
	Text           Text        `json:"text" desc:"消息文本内容"`
	Image          Image       `json:"image" desc:"图片"`
	Link           Link        `json:"link" desc:"图片"`
	MiniProgram    MiniProgram `json:"miniprogram" desc:"miniprogram"`
}

// QyMassHairResponse 获取企业标签数据请求返回
type QyMassHairResponse struct {
	QyResponse
	FailList []string `json:"fail_list" desc:"fail_list"`
	MsgID    string   `json:"msgid" desc:"msgid"`
}

// SendMassHairText
type SendMassHairText struct {
	QyGwRequest
	ChatType       string   `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string `json:"external_userid" desc:"客户的外部联系人id列表"`
	Text           `json:"text"`
}

// SendMassHairTextImage
type SendMassHairTextImage struct {
	QyGwRequest
	ChatType       string   `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string `json:"external_userid" desc:"客户的外部联系人id列表"`
	Text           `json:"text"`
	Image          `json:"image"`
}

// SendMassHairTextLink
type SendMassHairTextLink struct {
	QyGwRequest
	ChatType       string   `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string `json:"external_userid" desc:"客户的外部联系人id列表"`
	Text           `json:"text"`
	Link           Link `json:"link" desc:"图片"`
}

// SendMassHairMiniProgram
type SendMassHairMiniProgram struct {
	QyGwRequest
	ChatType       string   `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string `json:"external_userid" desc:"客户的外部联系人id列表"`
	Text           `json:"text"`
	MiniProgram    MiniProgram `json:"miniprogram" desc:"miniprogram"`
}

// MsgResultInfoRequest 群发消息发送结果请求参数
type MsgResultInfoRequest struct {
	QyGwRequest
	MsgID string `json:"msgid"`
}

// MsgResultInfoResponse 群发消息发送结果
type MsgResultInfoResponse struct {
	ExternalUserID string `json:"external_userid"`
	ChatID         string `json:"chat_id"`
	UserID         string `json:"userid"`
	Status         int64  `json:"status"`
	SendTime       int64  `json:"send_time"`
}

// ExternalContactListRequest 外部联系人列表请求参数
type ExternalContactListRequest struct {
	QyGwRequest
	UserID string `json:"userid"`
}

// ExternalContactListResponse 外部联系人列表返回值
type ExternalContactListResponse struct {
	QyResponse
	ExternalUserIDList []string `json:"external_userid"`
}

// ExternalGroupChatListRequest 客户群列表请求参数
type ExternalGroupChatListRequest struct {
	QyGwRequest
	StatusFilter int `json:"status_filter"`
	OwnerFilter  struct {
		UserIDList  []string `json:"userid_list"`
		PartyIDList []int    `json:"partyid_list"`
	} `json:"owner_filter"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ExternalGroupChatListRequest 客户群列表请求参数官网更新版
type ExternalGroupChatListRequestV2 struct {
	QyGwRequest
	StatusFilter int `json:"status_filter"`
	OwnerFilter  struct {
		UserIDList  []string `json:"userid_list"`
		PartyIDList []int    `json:"partyid_list"`
	} `json:"owner_filter"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// ExternalGroupChatListResponse 客户群列表返回结构
type ExternalGroupChatListResponse struct {
	QyResponse
	GroupChatList []*ExternalGroupChatItem `json:"group_chat_list"`
}

// ExternalGroupChat 客户群列表item
type ExternalGroupChatItem struct {
	ChatID string `json:"chat_id"`
	Status int    `json:"status"`
}

// UserInfoRequest 获取用户身份参数
type UserInfoRequest struct {
	QyGwRequest
	Code string `json:"code"`
}

// AuthInfoResponse 企业微信授权信息
type AuthInfoResponse struct {
	QyResponse
	*AuthAllUserInfo
}

// AuthAllUserInfo 授权所以成员信息
type AuthAllUserInfo struct {
	OpenId         string `json:"openid" desc:"非企业成员的标识"`
	ExternalUserid string `json:"external_userid" desc:"外部联系人id"`
	UserId         string `json:"UserId" desc:"成员ID"`
	DeviceId       string `json:"DeviceId" desc:"手机设备号"`
}

// AgentTicketRequest 应用 ticket 请求参数
type AgentTicketRequest struct {
	QyGwRequest
	Type string `json:"type"`
}

// TicketResponse JS ticket 响应数据
type TicketResponse struct {
	QyResponse
	Ticket  string `json:"ticket"`
	Expires int    `json:"expires_in"`
}

// SendMsgTemplateReq 发送群发消息请求结构体
type SendMsgTemplateReq struct {
	QyGwRequest
	ChatType       string                       `json:"chat_type" desc:"群发任务的类型"`
	ExternalUserid []string                     `json:"external_userid" desc:"客户的外部联系人id列表"`
	Sender         string                       `json:"sender" desc:"发送企业群发消息的成员userid"`
	Text           *Text                        `json:"text,omitempty" desc:"消息文本内容"`
	AllowSelect    bool                         `json:"allow_select"  desc:"是否可二次选择接收人"`
	Attachments    []SendMsgTemplateAttachments `json:"attachments"`
}
type SendMsgTemplateAttachments struct {
	MsgType     string       `json:"msgtype"`
	Image       *Image       `json:"image,omitempty"`
	Link        *Link        `json:"link,omitempty"`
	MiniProgram *MiniProgram `json:"miniprogram,omitempty"`
	Video       *Video       `json:"video,omitempty"`
	File        *File        `json:"file,omitempty"`
}

// GroupMsgSendResultRequest 查询群发消息请求参数结构体
type GroupMsgSendResultRequest struct {
	QyGwRequest
	MsgID  string `json:"msgid"`
	UserID string `json:"userid"`
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor,omitempty"`
}

// GroupMsgSendResultResponse 查询群发消息结果返回结构体
type GroupMsgSendResultResponse struct {
	QyResponse
	NextCursor string                    `json:"next_cursor"`
	SendList   []*GroupMsgSendResultItem `json:"send_list"`
}

// GroupMsgSendResultItem 查询群发消息结果
type GroupMsgSendResultItem struct {
	ExternalUserID string `json:"external_userid"`
	ChatID         string `json:"chat_id"`
	UserID         string `json:"userid"`
	Status         int    `json:"status"`
	SendTime       int64  `json:"send_time"`
}

type UploadMediaRequest struct {
	QyGwRequest
	File         *os.File `json:"img_url"`
	EnterpriseID int64    `json:"enterprise_id"`
}

type UploadMediaResponse struct {
	QyResponse
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

// MediaContent 素材文件结构体
type MediaContent struct {
	Text        *Text        `json:"text,omitempty" desc:"消息文本内容"`
	Image       *Image       `json:"image,omitempty" desc:"图片"`
	Link        *Link        `json:"link,omitempty" desc:"图片"`
	MiniProgram *MiniProgram `json:"miniprogram,omitempty" desc:"miniprogram"`
	Video       *Video       `json:"video,omitempty" desc:"视频"`
	Files       *File        `json:"files,omitempty" desc:"文件"`
}
type MediaContentData struct {
	Text    *Text          `json:"text,omitempty" desc:"消息文本内容"`
	Content []MediaContent `json:"content"`
}

// 获取指定的应用详情
type GetAgentRequest struct {
	QyGwRequest
	Agentid int64 `json:"agentid"`
}

// 获取指定的应用详情
type GetAgentResponse struct {
	QyResponse
	Agentid        int                 `json:"agentid"`
	AllowUserinfos AgentAllowUserinfos `json:"allow_userinfos"`
	AllowPartys    AgentAllowPartys    `json:"allow_partys"`
}

type AgentUser struct {
	Userid string `json:"userid"`
}
type AgentAllowUserinfos struct {
	User []AgentUser `json:"user"`
}
type AgentAllowPartys struct {
	Partyid []int `json:"partyid"`
}

type TransferResultRequst struct {
	QyGwRequest
	HandoverUserid string `json:"handover_userid" desc:"原添加成员的userid"`
	TakeoverUserid string `json:"takeover_userid" desc:"接替成员的userid"`
	Cursor         string `json:"cursor" desc:"分页查询的cursor，每个分页返回的数据不会超过1000条；不填或为空表示获取第一个分页"`
}

type TransferResultRespone struct {
	QyResponse
	Customer   []TransferResultCustomer `json:"customer" desc:"客户信息"`
	NextCursor string                   `json:"next_cursor" desc:"下个分页的起始cursor"`
}
type TransferResultCustomer struct {
	ExternalUserid string `json:"external_userid" desc:"转接客户的外部联系人userid"`
	Status         int    `json:"status" desc:"接替状态， 1-接替完毕 2-等待接替 3-客户拒绝 4-接替成员客户达到上限 5-无接替记录"`
	TakeoverTime   int64  `json:"takeover_time" desc:"接替客户的时间，如果是等待接替状态，则为未来的自动接替时间"`
}

// 分配在职成员的客户
// TransferCustomerAllotRequest 分配在职成员的客户请求
type TransferCustomerAllotRequest struct {
	QyGwRequest
	HandoverUserID     string   `json:"handover_userid"`
	TakeoverUserID     string   `json:"takeover_userid"`
	ExternalUserID     []string `json:"external_userid"`
	TransferSuccessMsg string   `json:"transfer_success_msg"`
}

// TransferCustomerAllotResponse 分配在职成员的客户请求响应
type TransferCustomerAllotResponse struct {
	QyResponse
	ExternalUserID []TransferCustomerAllotStatus `json:"external_userid"`
}

type TransferCustomerAllotStatus struct {
	ExternalUserID string `json:"external_user_id"`
	ErrCode        int    `json:"err_code"`
}

// ResignedTransferCustomerReq 分配离职成员的客户-请求参数
type ResignedTransferCustomerReq struct {
	QyGwRequest
	HandoverUserid string   `json:"handover_userid"` // 原跟进成员的userid
	TakeoverUserid string   `json:"takeover_userid"` // 接替成员的userid
	ExternalUserid []string `json:"external_userid"` // 客户的external_userid列表，最多一次转移100个客户
}

// ResignedTransferCustomerResp 分配离职成员的客户-返回值
type ResignedTransferCustomerResp struct {
	QyResponse
	Customer []struct {
		ExternalUserid string `json:"external_userid"`
		Errcode        int    `json:"errcode"`
	} `json:"customer"`
}

// ResignedTransferCustomerGroupReq 分配离职成员的客户群-请求参数
type ResignedTransferCustomerGroupReq struct {
	QyGwRequest
	ChatIdList []string `json:"chat_id_list"` // 需要转群主的客户群ID列表。取值范围： 1 ~ 100
	NewOwner   string   `json:"new_owner"`    // 新群主ID
}

// ResignedTransferCustomerGroupReps 分配离职成员的客户群-返回值
type ResignedTransferCustomerGroupReps struct {
	QyResponse
	FailedChatList []struct {
		ChatId  string `json:"chat_id"`
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	} `json:"failed_chat_list"`
}

// GetUnassignedListRequest 获取待分配的离职成员列表Request
type GetUnassignedListRequest struct {
	QyGwRequest
	Cursor   string `json:"cursor"`
	PageSize int    `json:"page_size"`
}

// GetUnassignedListResponse 获取待分配的离职成员列表Response
type GetUnassignedListResponse struct {
	QyResponse
	Info       []UnassignedList `json:"info"`
	IsLast     bool             `json:"is_last"`
	NextCursor string           `json:"next_cursor"`
}

type UnassignedList struct {
	HandoverUserid string `json:"handover_userid"`
	ExternalUserid string `json:"external_userid"`
	DimissionTime  int    `json:"dimission_time"`
}

type GetResignedTransferResultRequest struct {
	QyGwRequest
	HandoverUserid string `json:"handover_userid"` // 原跟进成员的userid
	TakeoverUserid string `json:"takeover_userid"` // 接替成员的userid
	Cursor         string `json:"cursor"`
}

type ResignedTransferResultResponse struct {
	QyResponse
	Customer []struct {
		ExternalUserid string `json:"external_userid"`
		Status         int    `json:"status"`
		TakeoverTime   int    `json:"takeover_time"`
	}
	NextCursor string `json:"next_cursor"`
}

type UserBehaviorDataReq struct {
	QyGwRequest
	UserId    []string `json:"userid"`
	PartyId   []string `json:"partyid"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
}

type UserBehaviorDataResp struct {
	QyResponse
	BehaviorData []BehaviorData `json:"behavior_data"`
}
type BehaviorData struct {
	StatTime            int     `json:"stat_time"`
	ChatCnt             int     `json:"chat_cnt"`
	MessageCnt          int     `json:"message_cnt"`
	ReplyPercentage     float64 `json:"reply_percentage"`
	AvgReplyTime        int     `json:"avg_reply_time"`
	NegativeFeedbackCnt int     `json:"negative_feedback_cnt"`
	NewApplyCnt         int     `json:"new_apply_cnt"`
	NewContactCnt       int     `json:"new_contact_cnt"`
}

// Code2SessionReq 小程序授权获取成员信息
type Code2SessionReq struct {
	Corpid             string `json:"corpid"`
	Userid             string `json:"userid" desc:"员工企微userid"`
	SessionKey         string `json:"session_key"`
	Mobile             string `json:"mobile"`
	MemberId           int64  `json:"member_id"`
	Name               string `json:"name"`
	IsManager          int    `json:"is_manager" desc:"是否企业管理，0否 1是"`
	StoreCode          string `json:"store_code"`
	StoreName          string `json:"store_name"`
	StoreCity          string `json:"store_city"`
	StorePhone         string `json:"store_phone"`
	StoreInventoryName string `json:"store_inventory_name"`
}

// UserInfoRequest 获取用户身份参数
type Code2SessionRequest struct {
	QyGwRequest
	JsCode    string `json:"js_code"`
	GrantType string `json:"grant_type"`
}

type MomentTaskImage struct {
	MsgType string `json:"msg_type,omitempty"`
	Image   []struct {
		MediaId string `json:"media_id"`
	}
}

type MomentTaskVideo struct {
	MsgType string `json:"msg_type,omitempty"`
	Video   []struct {
		MediaId string `json:"media_id"`
	} `json:"video"`
}

type MomentTaskLink struct {
	MsgType string `json:"msg_type,omitempty"`
	Link    []struct {
		Title   string `json:"title"`
		Url     string `json:"url"`
		MediaId string `json:"media_id"`
	} `json:"link"`
}

type AddMomentTaskAttachments struct {
	MsgType string       `json:"msgtype"`
	Image   *MomentImage `json:"image,omitempty"`
	Link    *MomentLink  `json:"link,omitempty"`
	Video   *MomentVideo `json:"video,omitempty"`
}

type AddMomentVisibleRange struct {
	SenderList *struct {
		UserList       []string `json:"user_list"`
		DepartmentList []int    `json:"department_list,omitempty"`
	} `json:"sender_list"`
	ExternalContactList *struct {
		TagList []string `json:"tag_list"`
	} `json:"external_contact_list,omitempty"`
}
type AddMomentTaskMediaContent struct {
	Text        *Text                      `json:"text,omitempty" desc:"消息文本内容"`
	Attachments []AddMomentTaskAttachments `json:"attachments,omitempty"`
}

type AddMomentTaskReq struct {
	QyGwRequest
	AddMomentTaskMediaContent
	VisibleRange AddMomentVisibleRange `json:"visible_range,omitempty"`
}

type AddMomentTaskResp struct {
	QyResponse
	Jobid string `json:"jobid" desc:"异步任务id"`
}

// MomentImage 朋友圈图片
type MomentImage struct {
	MediaId string `json:"media_id"`
}

// MomentLink 朋友圈链接
type MomentLink struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	MediaId string `json:"media_id"`
}

// LimitMomentLinkForQy 在发送企微朋友圈请求前限制链接标题的 UTF-8 字节数。
func LimitMomentLinkForQy(l MomentLink) MomentLink {
	l.Title = truncateUTF8ByBytes(l.Title, qyLinkTitleMaxBytes)
	return l
}

// MomentVideo 朋友圈视频
type MomentVideo struct {
	MediaId string `json:"media_id"`
}

// MomentTaskResultRequest 查询朋友圈的发表结果
type MomentTaskResultRequest struct {
	QyGwRequest
	JobId string `json:"jobid"`
}

type MomentTaskResultResp struct {
	QyResponse
	Status int              `json:"status"`
	Type   string           `json:"type"`
	Result MomentTaskResult `json:"result"`
}
type MomentTaskResult struct {
	QyResponse
	MomentId                   string                           `json:"moment_id"`
	InvalidSenderList          MomentInvalidSenderList          `json:"invalid_sender_list"`
	InvalidExternalContactList MomentInvalidExternalContactList `json:"invalid_external_contact_list"`
}
type MomentInvalidSenderList struct {
	UserList       []string `json:"user_list"`
	DepartmentList []int    `json:"department_list"`
}
type MomentInvalidExternalContactList struct {
	TagList []string `json:"tag_list"`
}

type GetMomentTaskRequest struct {
	QyGwRequest
	MomentId string  `json:"moment_id"`
	Cursor   *string `json:"cursor,omitempty"`
}

type GetMomentTaskResp struct {
	QyResponse
	NextCursor string              `json:"next_cursor"`
	TaskList   []GetMomentTaskInfo `json:"task_list"`
}

type GetMomentTaskInfo struct {
	UserId        string `json:"userid"`
	PublishStatus int    `json:"publish_status"`
}

type GetMomentSendTaskRequest struct {
	GetMomentTaskRequest
	UserId string `json:"userid"`
}

type GetMomentCommentsRequest struct {
	QyGwRequest
	MomentId string `json:"moment_id"`
	UserId   string `json:"userid"`
}

type GetMomentCommentResp struct {
	QyResponse
	CommentList []struct {
		UserId         string `json:"userid"`
		ExternalUserid string `json:"external_userid"`
		CreateTime     int64  `json:"create_time"`
	} `json:"comment_list"`
	LikeList []struct {
		UserId         string `json:"userid"`
		ExternalUserid string `json:"external_userid"`
		CreateTime     int64  `json:"create_time"`
	} `json:"like_list"`
}

type GetMomentSendTaskResp struct {
	QyResponse
	NextCursor   string `json:"next_cursor"`
	CustomerList []struct {
		ExternalUserId string `json:"external_userid"`
	} `json:"customer_list"`
}

type CancelMomentTaskRequest struct {
	QyGwRequest
	MomentId string `json:"moment_id"`
}

type CancelGroupMsgSendRequest struct {
	QyGwRequest
	MsgId string `json:"msgid"`
}

type SaTagListResp struct {
	QyGwRequest
	TagList []TagInfo `json:"taglist"`
}
type TagInfo struct {
	TagId   int    `json:"tagid"`
	TagName string `json:"tagname"`
}

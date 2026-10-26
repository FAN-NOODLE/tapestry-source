package fission

import (
	"time"

	"micro-scrm/internal/pkg/common"
)

const (
	StatusNotStarted = 1 + iota // 活动未开始
	StatusOffShelf              // 活动已下架
	StatusEnd                   // 活动已结束
	StatusOnGoing               // 活动进行中
	StatusDeleted               // 活动已删除
	StatusOnShelf    = 1        // 活动上架中
	IsMutualType     = 1        // 含有互斥活动
)

// 裂变类型
const (
	MemberFissionType = 1 + iota
	GroupFissionType
	VipFissionType
)

// 添加好友的状态
const (
	InitRelationStatus = iota
	AddMemberStatus
	DelRelationStatus
)

// 人员群的状态
const (
	InitGroupStatus = iota
	JoinGroupStatus
	LeaveGroupStatus
)

const (
	ExchangeWarn = "fission_warn" //兑换奖品提醒配置的key值
	WarnOn       = 1              //提醒开启
	WarnOff      = 0              //提醒关闭
	WarnTitle    = "兑换奖品提醒"       //兑奖提醒标题
	WarnText     = "成功兑换奖品"       //兑奖提醒文本
)

// 操作活动动作
const (
	ActionAdd    = 1 + iota // 添加活动
	ActionUpdate            // 编辑活动
)

const (
	CondTypeAnyOne        = 1 // 参与条件：1-添加任意员工
	CondTypeFissionMember = 2 // 参与条件：2-添加裂变员工
)

const (
	TypeSystem = 1 + iota // 系统联系我二维码
	TypeFissionDefault
	TypeFissionSend
)

// 对积分操作的行为
const (
	ActionScore    = "get_score"
	ActionReduce   = "reduce"
	ActionRecovery = "recovery"
)

const (
	PrizeStatusOK          = 0
	PrizeStatusNotStart    = 1
	PrizeStatusEmpty       = 2
	PrizeStatusDeleted     = 3
	PrizeStatusExchangeEnd = 4
)

const EnterpriseAuthTypeDefault = 1

const (
	ContactTypeActivityDefault = 2
	ContactTypeFissionUser     = 3
)

const ScrollListRobotNum = 20    // 弹幕数量
const ScrollListAddRobotTime = 5 // 弹幕列表增加假数据延迟时间(分钟)

// 欢迎语封面图 && 欢迎语描述
const (
	ActivityCover = "https://img-crs.vchangyi.com/goods16106094944490.png"
	ActivityStr   = "邀请好友助力，好礼兑不停"
)

// ActivityListRequest 裂变活动列表请求结构体
type ActivityListRequest struct {
	EnterpriseID int64              `json:"enterprise_id" desc:"企业Id"` //EnterpriseID企业ID
	Page         int64              `json:"page" desc:"页码"`            //Page页码
	Size         int64              `json:"size" desc:"数量"`            //Size数量
	Search       ActivityListSearch `json:"search"`
}

// ActivityListSearch 搜索字段
type ActivityListSearch struct {
	ActivityName string `json:"activity_name" desc:"活动名称"`             //ActivityName活动名称
	ChannelID    int64  `json:"channel_id" desc:"所属渠道"`                //ChannelID渠道ID
	FissionType  int    `json:"fission_type" desc:"裂变类型"`              //FissionType  裂变类型
	Status       int    `json:"status" desc:"活动状态"`                    //Status 活动状态
	TimeType     int    `json:"time_type" desc:"筛选时间类型:1-活动时间 2-创建时间"` //TimeType 筛选时间类型
	StartTime    int64  `json:"start_time" desc:"开始时间"`                //StartTime 开始时间
	EndTime      int64  `json:"end_time" desc:"结束时间"`                  //EndTime 结束时间
	NowTime      int64  `json:"now_time" desc:"当前时间"`                  //NowTime 当前时间
}

// ActivityListResponse 裂变活动返回结构体
type ActivityListResponse struct {
	common.PageType
	List []*ActivityEntity `json:"list"` //List 活动列表
}

// ExchangePrizeRecordsEntity 奖品兑换记录
type ExchangePrizeRecordsEntity struct {
	ID               int64  `db:"id" json:"id"`
	EnterpriseID     int64  `db:"enterprise_id" json:"enterprise_id"`     // EnterpriseID 企业id
	ActivityID       int64  `db:"activity_id" json:"activity_id"`         // ActivityID 活动id
	PrizeID          int64  `db:"prize_id" json:"prize_id"`               // PrizeID 奖品id
	FissionUserID    int64  `db:"fission_user_id" json:"fission_user_id"` // FissionUserID 裂变用户关系表id
	UserID           int64  `db:"user_id"`
	CostAssistantNum int    `db:"cost_assistant_num" json:"cost_assistant_num"` // CostAssistantNum 消耗的助力人数
	DeliverName      string `db:"deliver_name" json:"deliver_name"`             // DeliverName 收件人姓名
	DeliverMobile    string `db:"deliver_mobile" json:"deliver_mobile"`         // DeliverMobile 收件人手机号
	DeliverAddress   string `db:"deliver_address" json:"deliver_address"`       // DeliverAddress 收件人地址
	DeliverArea      string `db:"deliver_area" json:"deliver_area"`
	//DeliverAreaCode  int    `db:"deliver_area_code" json:"deliver_area_code"`
	DeliverAreaCode string `db:"deliver_area_code" json:"deliver_area_code"`
	Num             int    `db:"num"`
	CreatedAt       int64  `db:"created_at" json:"created_at"` // CreatedAt 创建时间
	UpdatedAt       int64  `db:"updated_at" json:"updated_at"` // UpdatedAt 更新时间
	DeletedAt       int64  `db:"deleted_at" json:"deleted_at"` // DeletedAt 删除时间
}

// ExchangeDetail 兑换详情
type ExchangeDetail struct {
	ExchangeInfo
	AddressInfo
}

// AddressInfo 地址信息
type AddressInfo struct {
	DeliverName   string `json:"deliver_name"`   //DeliverName 收件人姓名
	DeliverMobile string `json:"deliver_mobile"` //DeliverMobile 收件人电话
	Address       string `json:"address"`        //Address   收件人地址
	Area          string `json:"area"`           //Area 收货地址省市区
	AreaCode      string `json:"area_code"`      //AreaCode 收货地址省市区编码
}

// ExchangeInfo 奖品兑换信息
type ExchangeInfo struct {
	ID               int64  `json:"id" `
	ExchangeTime     int64  `json:"exchange_time" db:"exchange_time"` //ExchangeTime 兑换时间
	PrizePic         string `json:"prize_pic" db:"prize_pic"`         //PrizePic 奖品图片
	PrizeName        string `json:"prize_name" db:"prize_name"`       //PrizeName 精品名称
	CostAssistantNum int    `json:"cost_assistant_num"`               //CostAssistantNum 消耗的助力值
	Num              int    `json:"num" db:"num"`                     //Num 兑换数量
}

//UserPrizeRequest 用户奖品查询结构体
type UserPrizeRequest struct {
	EnterpriseID int64 `json:"enterprise_id"` //EnterpriseID 企业ID
	UserID       int64 `json:"user_id"`       //UserID  兑换人ID
	ActivityID   int64 `json:"activity_id"`   //ActivityID活动ID
	PrizeID      int64 `json:"prize_id"`      //PrizeID 奖品ID
}

//ExchangePrizeRequest用户兑奖信息查询结构体
type ExchangePrizeRequest struct {
	UserPrizeRequest
	AddressInfo
}

// ExchangeList 兑换列表
type ExchangeList struct {
	common.PageType
	List []*ExchangeInfo `json:"list"` //List 兑换记录列表
}

// PrizeList 奖品列表
type PrizeList struct {
	common.PageType
	List []*PrizeEntity `json:"list"` //List 奖品列表
}

// PrizeDetail 奖品详情-h5页面使用
type PrizeDetail struct {
	PrizeEntity
	ActivityDesc string `json:"activity_desc"` //ActivityDesc 活动描述
	Status       int    `json:"status"`        //Status 活动状态
}

// AssistantList 助力列表-h5页面
type AssistantList struct {
	common.PageType
	List            []*AssistantItem `json:"list"`              //List 助力记录列表
	CurAssistantNum int              `json:"cur_assistant_num"` //CurAssistantNum 当前助力人数
}

//AssistantItem 助力人信息
type AssistantItem struct {
	AssistantUserID     int64  `json:"assistant_user_id"`                                    //AssistantUserID 助力人ID
	Name                string `json:"name" desc:"助力好友名称"`                                   //Name 助力好友名称
	ImgUrl              string `json:"img_url" desc:"头像"`                                    //ImgUrl 助力好友头像
	AddMemberStatus     int    `json:"add_member_status" desc:"添加好友状态：0-未添加，1-已添加，2-已删除"`    //AddMemberStatus 好友添加状态
	JoinChatGroupStatus int    `json:"join_chat_group_status" desc:"入群状态：0-未入群，1-已入群，2-已退群"` //JoinChatGroupStatus 好友入群状态
	Score               int    `json:"score"`                                                //Score 计分
	IsVip               int    `json:"is_vip"`
}

//AssistantUserInfo 助力人信息
type AssistantUserInfo struct {
	Name   string `json:"name" desc:"助力好友名称"` //Name助力好友名称
	ImgUrl string `json:"img_url" desc:"头像"`  //ImgUrlName助力好友头像
}

type MpHomeInfo struct {
	ActivityInfo struct {
		Name            string `json:"name"`                                      //Name 活动名称
		FissionType     int    `json:"fission_type" desc:"裂变类型：1-添加员工，2-添加员工且入群"` //FissionType裂变类型
		StartTime       int64  `json:"start_time"`                                //StartTime开始时间
		EndTime         int64  `json:"end_time"`                                  //EndTime结束时间
		ExchangeEndTime int64  `json:"exchange_end_time"`                         //ExchangeEndTime兑换时间
		ActivityStatus  int    `json:"activity_status"`                           //ActivityStatus 活动状态
		Desc            string `json:"desc"`                                      //Desc 活动描述
		ChatQrCode      string `json:"chat_qr_code"`                              //ChatQrCode
		SeverTime       int64  `json:"sever_time"`                                //SeverTime
		MemberQrCode    string `json:"member_qr_code"`                            //MemberQrCode
		RgbColors       string `json:"rgb_colors"`                                //RgbColors
		Styles          string `json:"styles"`                                    //Styles
		BgImg           string `json:"bg_img"`
		ShareImg        string `json:"share_img"`
	} `json:"activity_info"` //ActivityInfo 活动信息
	UserInfo struct {
		ImgUrl            string `json:"img_url"`                    //ImgUrl  用户头像
		CurAssistantNum   int    `json:"cur_assistant_num"`          //CurAssistantNum 当前助力人数
		TotalAssistantNum int    `json:"total_assistant_num"`        // TotalAssistantNum 累计助力人数
		ContactQrCode     string `json:"contact_qr_code"`            //ContactQrCode 联系我二维码
		JoinOtherActivity bool   `json:"join_other_mutual_activity"` //JoinOtherActivity 是否参加了其它活动
		UnionID           string `json:"union_id"`                   //UnionID   用户unionID
	} `json:"user_info"` //UserInfo 用户信息
	ShopInfo ShopAuthInfo `json:"shop_info"` //ShopInfo
}

type ShopAuthInfo struct {
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// PrizeEntity 裂变奖品信息
type PrizeEntity struct {
	ID           int64  `db:"id" json:"id"`
	EnterpriseID int64  `db:"enterprise_id" json:"enterprise_id"`                                   // EnterpriseID 企业id
	ActivityID   int64  `db:"activity_id" json:"activity_id"`                                       // ActivityID 活动id
	Name         string `db:"name" json:"name" validate:"min=1,max=50" min:"活动名称不能为空" max:"最多50个字"` // Name 奖品名称
	Pic          string `db:"pic" json:"pic" validate:"required"`                                   // Pic 奖品图片链接
	Stocks       int    `db:"stocks" json:"stocks" validate:"min=0,max=999999"`                     // Stocks 奖品库存
	ExchangeNum  int    `db:"exchange_num" json:"exchange_num"`                                     // ExchangeNum 已兑换奖品数
	NeedNum      int    `db:"need_num" json:"need_num" validate:"min=0,max=999999"`                 // FissionNum 裂变人数
	LimitNum     int    `db:"limit_num" json:"limit_num" validate:"min=0,max=999999"`               // LimitNum 每人限兑换个数
	Sort         int64  `db:"sort" json:"sort"`                                                     // Sort 排序权重：升序
	CreatedAt    int64  `db:"created_at" json:"created_at"`                                         // CreatedAt 创建时间
	UpdatedAt    int64  `db:"updated_at" json:"updated_at"`                                         // UpdatedAt 更新时间
	DeletedAt    int64  `db:"deleted_at" json:"deleted_at"`                                         // DeletedAt 删除时间
}

//AssistantRecordInfo 助力记录信息
type AssistantRecordInfo struct {
	AssistantRecordsEntity
	ImgUrl   string `json:"img_url" db:"img_url"`     //ImgUrl 头像地址
	Nickname string `json:"nick_name" db:"nick_name"` //Nickname 昵称
}

// AssistantRecordsEntity 助力记录结构体
type AssistantRecordsEntity struct {
	ID                  int64  `db:"id" json:"id"`
	EnterpriseID        int64  `db:"enterprise_id" json:"enterprise_id"`                   // EnterpriseID 企业id
	ActivityID          int64  `db:"activity_id" json:"activity_id"`                       // ActivityID 活动id
	FissionUserID       int64  `db:"fission_user_id" json:"fission_user_id"`               // FissionUserID 裂变用户关系表id
	AssistantUserID     int64  `db:"assistant_user_id" json:"assistant_user_id"`           // AssistantUserID 助力用户id
	ExternalUserID      string `json:"external_user_id" db:"external_user_id"`             // ExternalUserID 企业微信外部联系人id
	MemberID            int64  `db:"member_id" json:"member_id"`                           // MemberID 添加的成员ID
	AddMemberStatus     int    `db:"add_member_status" json:"add_member_status"`           // AddMemberStatus 添加好友记录：0-未添加，1-已添加，2-已删除
	AddMemberAt         int64  `db:"add_member_at" json:"add_member_at"`                   //AddMemberAt好友添加时间
	JoinChatGroupStatus int    `db:"join_chat_group_status" json:"join_chat_group_status"` // JoinChatGroupStatus 入群状态：0-未入群，1-已入群，2-已退群
	JoinChatGroupAt     int64  `db:"join_chat_group_at" json:"join_chat_group_at"`         // JoinChatGroupAt 加群时间
	Score               int    `db:"score" json:"score"`                                   // Score 计分
	CreatedAt           int64  `db:"created_at" json:"created_at"`                         // CreatedAt 创建时间
	UpdatedAt           int64  `db:"updated_at" json:"updated_at"`                         // UpdatedAt 更新时间
	DeletedAt           int64  `db:"deleted_at" json:"deleted_at"`                         // DeletedAt 删除时间
	IsVip               int    `db:"is_vip" json:"is_vip"`
}

// HasScore 判断是否加分
func (a *AssistantRecordsEntity) HasScore() bool {
	return a.Score > 0
}

// GetScore 获取相应行为的分值
func (a *AssistantRecordsEntity) GetScore(activityType int) int {
	switch true {
	case activityType == MemberFissionType && a.AddMemberStatus == AddMemberStatus:
		return 1
	case activityType == GroupFissionType && a.AddMemberStatus == AddMemberStatus && a.JoinChatGroupStatus == JoinGroupStatus:
		return 1
	default:
		return 0
	}
}

// ActivityEntity 裂变活动实体
type ActivityEntity struct {
	ID                 int64  `db:"id" json:"id"`
	EnterpriseID       int64  `db:"enterprise_id" json:"enterprise_id"`                                                     // EnterpriseID 企业id
	Name               string `db:"name" json:"name" validate:"min=1,max=10" max:"活动名称最多为10个字"`                             // Name 活动名称
	Desc               string `db:"desc" json:"desc" validate:"required" required:"活动描述不能为空"`                               // Desc 活动描述
	Status             int    `db:"status" json:"status"`                                                                   // Status 活动状态：1-上架中，2-已下架
	ChannelID          int64  `db:"channel_id" json:"channel_id" validate:"required" required:"渠道id不能为空"`                   // ChannelID 渠道id
	StartTime          int64  `db:"start_time" json:"start_time" validate:"required" required:"活动开始时间不能为空"`                 // StartTime 活动开始时间
	EndTime            int64  `db:"end_time" json:"end_time" validate:"required" required:"活动结束时间不能为空"`                     // EndTime 活动结束时间
	ExchangeEndTime    int64  `db:"exchange_end_time" json:"exchange_end_time" validate:"required" required:"兑换奖品截至时间不能为空"` // ExchangeEndTime 兑换奖品截至时间
	FissionType        int    `db:"fission_type" json:"fission_type" validate:"required" required:"裂变类型不能为空"`               // FissionType 裂变类型：1-添加员工，2-添加员工且入群
	CondType           int    `db:"cond_type" json:"cond_type"`                                                             // CondType 参与条件：1-添加任意员工，2-添加裂变员工
	MutualType         int    `db:"mutual_type" json:"mutual_type"`                                                         // MutualType 互斥类型：1-互斥，2-不互斥
	DrainageWelcomesID int64  `db:"drainage_welcomes_id" json:"drainage_welcomes_id"`                                       // DrainageWelcomesID 欢迎语id
	ContactId          int64  `db:"contact_id" json:"contact_id"`                                                           // 活动默认联系我二维码id，drainage_contacts表id
	UserNum            int    `db:"user_num" json:"user_num"`                                                               // UserNum 发起人数
	FissionNum         int    `db:"fission_num" json:"fission_num"`                                                         // FissionNum 好友裂变数
	RunAwayNum         int    `db:"run_away_num" json:"run_away_num"`                                                       // RunAwayNum 流失人数
	JoinGroupNum       int    `db:"join_group_num" json:"join_group_num"`                                                   // JoinGroupNum 群成员裂变数
	LeaveGroupNum      int    `db:"leave_group_num" json:"leave_group_num"`                                                 // LeaveGroupNum 群成员退群数
	CreatedAt          int64  `db:"created_at" json:"created_at"`                                                           // CreatedAt 创建时间
	UpdatedAt          int64  `db:"updated_at" json:"updated_at"`                                                           // UpdatedAt 更新时间
	DeletedAt          int64  `db:"deleted_at" json:"deleted_at"`                                                           // DeletedAt 删除时间
	ChannelName        string `db:"channel_name" json:"channel_name"`                                                       //ChannelName 渠道名称
}

// RunTime 获取时间
func (r *ActivityEntity) RunTime() bool {
	nowTime := time.Now().Unix()
	if r.StartTime > nowTime || r.EndTime < nowTime {
		return false
	}

	return true
}

// IsGroupType 判断是群类型活动
func (r *ActivityEntity) IsGroupType() bool {
	return r.FissionType == GroupFissionType
}

// RunExchange 是否能兑换
func (r *ActivityEntity) RunExchange() bool {
	nowTime := time.Now().Unix()
	if r.ExchangeEndTime < nowTime || r.StartTime > nowTime {
		return false
	}

	return true
}

// ChatGroupEntity 活动社群实体
type ChatGroupEntity struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"` // EnterpriseID 企业id
	ActivityID   int64  `json:"activity_id" db:"activity_id"`     // ActivityID 活动id
	ChatID       int64  `json:"chat_id" db:"chat_id"`             // ChatID 群聊id
	QrCode       string `json:"qr_code" db:"qr_code"`             // MemberQrCode 群聊 二维码
	Sort         int    `json:"sort" db:"sort"`                   // Sort 群聊排序
	CreatedAt    int64  `json:"created_at" db:"created_at"`       // CreatedAt 创建时间
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`       // UpdatedAt 更新时间
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`       // DeletedAt 删除时间
}

// MembersEntity 裂变活动员工实体
type MembersEntity struct {
	ID           int64 `db:"id" json:"id"`
	EnterpriseID int64 `db:"enterprise_id" json:"enterprise_id"` // EnterpriseID 企业id
	ActivityID   int64 `db:"activity_id" json:"activity_id"`     // ActivityID 活动id
	MemberID     int64 `db:"member_id" json:"member_id"`         // MemberID 员工id
	CreatedAt    int64 `db:"created_at" json:"created_at"`       // CreatedAt 创建时间
	UpdatedAt    int64 `db:"updated_at" json:"updated_at"`       // UpdatedAt 更新时间
	DeletedAt    int64 `db:"deleted_at" json:"deleted_at"`       // DeletedAt 删除时间
}

// MutualRelationEntity 活动互斥关系表实体
type MutualRelationEntity struct {
	ID               int64 `db:"id" json:"id"`
	EnterpriseID     int64 `db:"enterprise_id" json:"enterprise_id"`           // EnterpriseID 企业id
	ActivityID       int64 `db:"activity_id" json:"activity_id"`               // ActivityID 活动id
	MutualActivityID int64 `db:"mutual_activity_id" json:"mutual_activity_id"` // MutualActivityID 互斥的活动id
	CreatedAt        int64 `db:"created_at" json:"created_at"`                 // CreatedAt 创建时间
	UpdatedAt        int64 `db:"updated_at" json:"updated_at"`                 // UpdatedAt 更新时间
	DeletedAt        int64 `db:"deleted_at" json:"deleted_at"`                 // DeletedAt 删除时间
}

// ActivityUserEntity 活动发起人结构体
type ActivityUserEntity struct {
	ID                int64  `db:"id" json:"id"`
	EnterpriseID      int64  `db:"enterprise_id" json:"enterprise_id"`             // EnterpriseID 企业id
	ActivityID        int64  `db:"activity_id" json:"activity_id"`                 // ActivityID 活动id
	UserID            int64  `db:"user_id" json:"user_id"`                         // UserID 用户id
	TotalAssistantNum int    `db:"total_assistant_num" json:"total_assistant_num"` // TotalAssistantNum 累计助力人数
	CurAssistantNum   int    `db:"cur_assistant_num" json:"cur_assistant_num"`     // CurAssistantNum 当前剩余助力人数
	ContactID         int64  `db:"contact_id" json:"contact_id"`                   // ContactID 联系我二维码id，drainage_contacts表id
	CreatedAt         int64  `db:"created_at" json:"created_at"`                   // CreatedAt 创建时间
	UpdatedAt         int64  `db:"updated_at" json:"updated_at"`                   // UpdatedAt 更新时间
	DeletedAt         int64  `db:"deleted_at" json:"deleted_at"`                   // DeletedAt 删除时间
	FissionNum        int    `db:"fission_num" json:"fission_num"`                 // FissionNum 好友裂变数
	RunAwayNum        int    `db:"run_away_num" json:"run_away_num"`               // RunAwayNum 流失人数
	JoinGroupNum      int    `db:"join_group_num" json:"join_group_num"`           // JoinGroupNum 群成员裂变数
	LeaveGroupNum     int    `db:"leave_group_num" json:"leave_group_num"`         // LeaveGroupNum 群成员退群数
	Nickname          string `json:"nick_name" db:"nick_name"`                     //Nickname昵称
	ImgUrl            string `json:"img_url" db:"img_url"`                         //ImgUrl头像地址
}

//ActivityDetailRequest 活动详情请求结构体
type ActivityDetailRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"  validate:"required" required:"企业ID不能为空"` //EnterpriseID企业ID
	ActivityID   int64  `json:"activity_id" db:"activity_id"  validate:"required" required:"活动ID不能为空"`     //ActivityID活动ID
	Page         int64  `json:"page" desc:"页码"`                                                            //Page 页码
	Size         int64  `json:"size" desc:"数量"`                                                            //Size 数量
	Nickname     string `json:"nick_name" desc:"发起人昵称"`                                                    //Nickname 昵称
}

//ActivityDetailResponse 活动详情响应结构体
type ActivityDetailResponse struct {
	Activity ActivityEntity `json:"activity"` //活动信息
}

//ActivityDetailListResponse 活动详情列表响应结构体
type ActivityDetailListResponse struct {
	List  *[]ActivityUserEntity `json:"list"`  //List 裂变发起人列表
	Page  int64                 `json:"page"`  //Page  页码
	Size  int64                 `json:"size"`  //Size 数量
	Total int64                 `json:"total"` //Total 总数
}

// ActivitySetRequest 设置裂变活动结构体
type ActivitySetRequest struct {
	EnterpriseID      int64                    `json:"enterprise_id" db:"enterprise_id"`
	Activity          *ActivityEntity          `json:"activity" validate:"required,dive" required:"活动内容不完整"`
	MemberIDs         []int64                  `json:"member_ids" validate:"min=1,max=100" min:"请选择裂变员工" max:"最多选择100个员工"`
	ChatGroups        []*ChatGroupEntity       `json:"group_list"`
	MutualActivityIDs []int64                  `json:"mutual_activity_ids"`
	WelcomeContent    []map[string]interface{} `json:"welcome_content" db:"content" desc:"内容" validate:"required" required:"欢迎语不能为空"`
	Prize             []*PrizeEntity           `json:"prize_list" desc:"活动奖品" validate:"min=1,dive,max=10" min:"最少添加一个奖品" max:"最多添加10个奖品"`
}

// MutualActivityListRequest 互斥活动下拉列表
type MutualActivityListRequest struct {
	Page                int64   `json:"page"`
	Size                int64   `json:"size"`
	EnterpriseID        int64   `json:"enterprise_id" db:"enterprise_id"`
	ActivityName        string  `json:"activity_name" desc:"活动名称"`
	MutualActivityIDs   []int64 `json:"mutual_activity_ids" desc:"互斥活动IDs"`
	MutualActivityCount int     `json:"mutual_activity_count" desc:"互斥活动IDs数量"`
	Filter              []int64 `json:"filter" desc:"过滤的社群"`
	FilterCount         int     `json:"filter_count" desc:"过滤的社群数量"`
	Action              string  `json:"action" desc:"操作类型 detail详情 list列表"`
}

// MutualActivityListResponse 裂变活动返回结构体
type MutualActivityListResponse struct {
	common.PageType
	List []*MutualActivityEntity `json:"list"`
}

// MutualActivityEntity 互斥活动结构体
type MutualActivityEntity struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`             // Name 活动名称
	StartTime int64  `db:"start_time" json:"start_time"` // StartTime 活动开始时间
	EndTime   int64  `db:"end_time" json:"end_time"`     // EndTime 活动结束时间
	Status    int    `db:"status" json:"status"`         // Status 活动状态：1-上架中，2-已下架
}

// ChatGroupListRequest 社群列表请求结构体
type ChatGroupListRequest struct {
	Page         int64   `json:"page"`
	Size         int64   `json:"size"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id"`
	GroupName    string  `json:"group_name" desc:"社群名称"`
	ChatIDs      []int64 `json:"chat_ids" desc:"社群IDs"`
	ChatIDCount  int     `json:"chat_id_count" desc:"社群IDs数量"`
	Filter       []int64 `json:"filter" desc:"过滤的社群"`
	FilterCount  int     `json:"filter_count" desc:"过滤的社群数量"`
}

//  ChatGroupListListResponse 社群列表返回结构体
type ChatGroupListListResponse struct {
	common.PageType
	List []*WorkChatGroupEntity `json:"list"`
}

// WorkChatGroupEntity 群组结构体
type WorkChatGroupEntity struct {
	GroupID         int64  `json:"id" db:"id" desc:"社群ID"`
	OwnerMemberID   int64  `json:"owner_member_id" db:"owner_member_id" desc:"社群ID"`
	GroupName       string `json:"name" db:"name" desc:"社群名称"`
	OwnerMemberName string `json:"owner_member_name" db:"owner_member" desc:"群主"`
	GroupNumber     int    `json:"member_num" db:"member_num" desc:"群内人数"`
	QrCode          string `json:"qr_code" db:"qr_code" desc:"群二维码"`
	Sort            int    `json:"sort" db:"sort" desc:"排序"`
}

// ActivityInfoRequest 活动详情请求结构体
type ActivityInfoRequest struct {
	ActivityID   int64 `json:"activity_id" desc:"activity_id"`
	EnterpriseID int64 `json:"enterprise_id" desc:"enterprise_id"`
}

// StyleConfigRequest 配色列表请求参数
type StyleConfigRequest struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required"`
	RgbColors    string `json:"rgb_colors"`
	Styles       string `json:"styles"`
	BgImg        string `json:"bg_img"`
	ShareImg     string `json:"share_img"`
	ID           int64  `json:"id"`
}

type DelStyleConfigRequest struct {
	EnterpriseID int64   `json:"enterprise_id" binding:"required"`
	ID           []int64 `json:"id"`
}

type StyleConfigResponse struct {
	ID        int64  `json:"id"`
	RgbColors string `json:"enterprise_id"`
	Styles    string `json:"styles"`
}

// ActivityInfoResponse 活动详情响应结构体
type ActivityInfoResponse struct {
	Activity       ActivityEntity           `json:"activity" desc:"活动信息"`
	Members        []*Member                `json:"member_list" desc:"员工列表"`
	ChatGroups     []WorkChatGroupEntity    `json:"group_list" desc:"社群列表"`
	MutualActivity []*MutualActivityEntity  `json:"mutual_activity_list" desc:"互斥活动"`
	WelcomeContent []map[string]interface{} `json:"welcome_content" desc:"欢迎语"`
	Prize          []*PrizeEntity           `json:"prize_list" desc:"活动奖品"`
}

// UpdateStatusRequest 更新活动状态请求结构体
type UpdateStatusRequest struct {
	EnterpriseID int64 `db:"fission_user_id" json:"enterprise_id"`
	ActivityID   int64 `json:"activity_id" desc:"activity_id"`
	Status       int   `db:"status" json:"status"` // Status 活动状态：1-上架中，2-已下架
}

// InviteDetailRequest发起人邀请明细请求结构体
type InviteDetailRequest struct {
	EnterpriseID int64 ` json:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	ActivityID   int64 ` json:"activity_id" validate:"required" required:"活动ID不能为空"` //活动id
	UserID       int64 ` json:"user_id"  validate:"required" required:"裂变用户ID不能为空"`  //裂变用户id
	Page         int64 `json:"page" desc:"页码" validate:"required" required:"Page不能为空"`
	Size         int64 `json:"size" desc:"数量" validate:"required" required:"Size不能为空"`
}

//InviteDetailResponse发起人邀请明细响应结构体
type InviteDetailResponse struct {
	List  []*AssistantRecordInfo `json:"list"`
	Page  int64                  `json:"page"`
	Size  int64                  `json:"size"`
	Total int64                  `json:"total"`
}

// Member员工列表
type Member struct {
	ID         int64  `db:"id" json:"id"`
	MemberID   int64  `db:"member_id" json:"member_id"`
	MemberName string `db:"name" json:"name"`
	UserID     string `db:"userid" json:"userid"`
}

// QwActivityListRequest 裂变活动请求结构体
type QwActivityListRequest struct {
	EnterpriseId   int64  `json:"enterprise_id" desc:"企业Id" validate:"required" required:"企业ID不能为空"`
	EnterpriseHash string `json:"enterprise_hash" desc:"企业hash" validate:"required" required:"企业hash不能为空"`
	SearchText     string `json:"search_text" desc:"搜索条件"`
	Status         int64  `json:"status" desc:"活动状态：1-上架中，2-已下架"`
	NowTime        int64  `json:"now_time"`
	Page           int64  `json:"page" validate:"required" required:"page不能为空"`
	Size           int64  `json:"size" validate:"required" required:"size不能为空"`
}

// QwActivityListResponse 裂变活动返回结构体
type QwActivityListResponse struct {
	List       *[]QwActivityList `json:"list" desc:"列表数据"`
	Statistics Statistics        `json:"statistics" desc:"统计数据"`
	Page       int64             `json:"page" desc:"页码"`
	Size       int64             `json:"size" desc:"单页条数"`
	Total      int64             `json:"total" desc:"总条数"`
}

// Statistics 统计数据
type Statistics struct {
	UserNum       int64 `json:"user_num" db:"user_num"  desc:"发起人数"`
	FissionNum    int64 `json:"fission_num" db:"fission_num" desc:"好友裂变数"`
	RunAwayNum    int64 `json:"run_away_num" db:"run_away_num" desc:"流失人数"`
	JoinGroupNum  int64 `json:"join_group_num" db:"join_group_num" desc:"群成员裂变数"`
	LeaveGroupNum int64 `json:"leave_group_num" db:"leave_group_num" desc:"群成员退群数"`
	UpdateAt      int64 `json:"update_at" db:"updated_at"`
}

// QwActivityList 列表数据
type QwActivityList struct {
	Id        int64  `json:"id" db:"id" desc:"活动ID"`
	Name      string `json:"name" db:"name"  desc:"活动名称"`
	Desc      string `json:"desc" db:"desc" desc:"活动描述"`
	Status    int    `json:"status" db:"status" desc:"活动状态"`
	HtmlUrl   string `json:"html_url" desc:"h5url地址"`
	ImgUrl    string `json:"img_url" desc:"图片地址"`
	StartTime int64  `json:"start_time" db:"start_time" desc:"活动开始时间"`
	EndTime   int64  `json:"end_time" db:"end_time" desc:"活动结束时间"`
}

//ExchangeDetailRequest 兑换记录请求结构体
type ExchangeDetailRequest struct {
	EnterpriseID  int64  `db:"enterprise_id" json:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	ActivityID    int64  `db:"activity_id" json:"activity_id"  desc:"活动id"` //活动id
	UserID        int64  `db:"user_id" json:"user_id" desc:"裂变用户id"`        //裂变用户id
	Page          int64  `json:"page" desc:"页码"`
	Size          int64  `json:"size" desc:"数量"`
	Nickname      string `json:"nick_name" desc:"用户昵称"`
	StartTime     int64  `json:"start_time" desc:"开始时间"`
	EndTime       int64  `json:"end_time" desc:"结束时间"`
	DeliverName   string `db:"deliver_name" json:"deliver_name" desc:"收件人姓名"`      // DeliverName 收件人姓名
	DeliverMobile string `db:"deliver_mobile" json:"deliver_mobile" desc:"收件人手机号"` // DeliverMobile 收件人手机号

}

//ExchangeRecordInfo 兑换记录信息
type ExchangeRecordInfo struct {
	ExchangePrizeRecordsEntity
	ExchangeInfo
	Nickname     string `json:"nick_name" db:"nick_name" desc:"用户昵称"`
	ImgUrl       string `json:"img_url" db:"img_url" desc:"头像"`
	ActivityName string `json:"activity_name" db:"activity_name" desc:"活动名称"`
}

//ExchangeDetailResponse 兑换记录响应结构体
type ExchangeDetailResponse struct {
	List  *[]ExchangeRecordInfo `json:"list"`
	Page  int64                 `json:"page"`
	Size  int64                 `json:"size"`
	Total int64                 `json:"total"`
}

// UpdateAssistantRecord 更新助力情况
type UpdateAssistantRecord struct {
	ID                  int64   `db:"id" json:"id"`
	IDS                 []int64 `json:"ids" db:"ids"`
	AddMemberStatus     int     `db:"add_member_status" json:"add_member_status"`           // AddMemberStatus 添加好友记录：0-未添加，1-已添加，2-已删除
	AddMemberAt         int64   `db:"add_member_at" json:"add_member_at"`                   //AddMemberAt好友添加时间
	JoinChatGroupStatus int     `db:"join_chat_group_status" json:"join_chat_group_status"` // JoinChatGroupStatus 入群状态：0-未入群，1-已入群，2-已退群
	JoinChatGroupAt     int64   `db:"join_chat_group_at" json:"join_chat_group_at"`         // JoinChatGroupAt 加群时间
	Score               int     `db:"score" json:"score"`                                   // Score 计分
	UpdateAt            int64   `json:"update_at" db:"update_at"`
}

func (a *UpdateAssistantRecord) GetScore(activityType int) int {
	record := AssistantRecordsEntity{
		ID:                  a.ID,
		AddMemberStatus:     a.AddMemberStatus,
		JoinChatGroupStatus: a.JoinChatGroupStatus,
	}

	return record.GetScore(activityType)
}

func (a *UpdateAssistantRecord) HasScore() bool {
	record := AssistantRecordsEntity{
		Score: a.Score,
	}

	return record.HasScore()
}

//ExchangeWarnRequest 配置兑换提醒请求结构体
type ExchangeWarnRequest struct {
	Status       int     `json:"status"`     //提醒开启状态 0-关闭 1-开启
	MemberIds    []int64 `json:"member_ids"` // 更多提醒员工的id
	EnterpriseID int64   `json:"enterprise_id"`
}

//ExchangeWarnRequest 配置兑换提醒响应结构体
type ExchangeWarnReponse struct {
	Status       int       `json:"status"`  //提醒开启状态 0-关闭 1-开启
	Members      []*Member `json:"members"` // 更多提醒员工信息
	EnterpriseID int64     `json:"enterprise_id"`
}

// MpUser 公众号授权用户信息
type MpUser struct {
	ID           int64  `db:"id" json:"id"`
	OpenId       string `db:"open_id" json:"open_id"`
	NickName     string `db:"nick_name" json:"nick_name"`
	Sex          int    `db:"sex" json:"sex"`
	Province     string `db:"province" json:"province"`
	City         string `db:"city" json:"city"`
	Country      string `db:"country" json:"country"`
	ImgUrl       string `db:"img_url" json:"img_url"`
	Privilege    string `db:"privilege" json:"privilege"`
	UnionId      string `db:"union_id" json:"union_id"`
	EnterpriseId int64  `db:"enterprise_id" json:"enterprise_id"`
	CreatedAt    int    `db:"created_at" json:"created_at"`
	UpdatedAt    int    `db:"updated_at" json:"updated_at"`
	DeletedAt    int    `db:"deleted_at" json:"deleted_at"`
}

// BatchDelete 批量删除
type BatchDelete struct {
	DeleteIDs []int64 `json:"ids" desc:"删除ids"`
	DeletedAt int64   `json:"deleted_at" desc:"删除时间"`
}

// ExistRelationReq 关系状态查询
type ExistRelationReq struct {
	CustomerID     int64   `json:"customer_id"`
	Members        []int64 `json:"members"`
	EnterpriseID   int64   `json:"enterprise_id"`
	RelationStatus int     `json:"relation_status"`
}

// Channel 渠道
type Channel struct {
	ChannelID int64 `json:"channel_id" db:"channel_id" desc:"channel_id"`
}

//UpdateStatisticsData 更新裂变统计数据
type UpdateStatisticsData struct {
	Statistics
	TotalAssistantNum int64 `json:"total_assistant_num" desc:"总助力人数"`
	CurAssistantNum   int64 `json:"cur_assistant_num" desc:"剩余助力人数"`
}

// MembersSet
type MembersSet struct {
	AddMembers []int64 `json:"add_members"`
	DelMembers []int64 `json:"del_members"`
}

// PrizeSet
type PrizeSet struct {
	AddPrizes    []*PrizeEntity `json:"add_prizes"`
	UpdatePrizes []*PrizeEntity `json:"update_prizes"`
	DelPrizes    []int64        `json:"del_prizes"`
}

// ChatGroupSet
type ChatGroupSet struct {
	AddGroups    []*ChatGroupEntity `json:"add_groups"`
	UpdateGroups []*ChatGroupEntity `json:"update_groups"`
	DelGroups    []int64            `json:"del_groups"`
}

// MutualSet
type MutualSet struct {
	AddMutual []int64 `json:"add_mutual"`
	DelMutual []int64 `json:"del_mutual"`
}

// ActivityByChatIDReq 通过群聊ID获取活动
type ActivityByChatIDReq struct {
	WithDel               bool  `json:"with_del"`
	ChatID                int64 `json:"chat_id"`
	EnterpriseID          int64 `json:"enterprise_id"`
	ActivityRunAt         int64 `db:"end_time" json:"end_time"`
	ActivityExchangeRunAt int64 `db:"exchange_end_time" json:"exchange_end_time" `
}

// WorkUserBasicInfo 外部联系人基本信息
type WorkUserBasicInfo struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	ImgURL string `json:"img_url"`
}

// AssistantRecordsWithTime 助力记录和入群时间
type AssistantRecordsWithTime struct {
	*AssistantRecordsEntity
	JoinTime int64 `json:"join_time" db:"join_time"`
}

// ScrollListRequest 滚动条数据列表请求结构体
type ScrollListRequest struct {
	ActivityID   int64 `json:"activity_id"`
	EnterpriseID int64 `json:"enterprise_id"`
	common.RequestPage
}

// Text 文本
type Text struct {
	Content string `json:"content"`
}

// Image 图片
type Image struct {
	MediaId string `json:"media"`
	PicUrl  string `json:"pic_url"`
}

// Link 链接
type Link struct {
	Title  string `json:"title"`
	PicUrl string `json:"pic_url"`
	Desc   string `json:"desc"`
	Url    string `json:"url"`
}

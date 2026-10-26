package members

import (
	"strconv"

	"micro-scrm/internal/pkg/sdk/work"
)

// ShopIdentifier 商城应用标识
const ShopIdentifier = "crs_shop"

// CrsManuCustomerIdentifier 侧边栏应用
const CrsManuCustomerIdentifier = "crs-manu-customer"

// LockForContactKey 同步通讯录加锁
const LockForContactKey = "sync_contract:%d"

const MemberTypeBA = "ba"

const (
	MemberByUserID = 0 // 通过userid获取
	MemberByMobile = 1 // 通过mobile获取
)

// Member 成员
type Member struct {
	ID           int64  `json:"id" db:"id" desc:"员工ID"`
	Userid       string `json:"userid" db:"userid" desc:"企业微信userid"`
	Name         string `json:"name" db:"name" desc:"员工姓名"`
	Mobile       string `json:"mobile" db:"mobile" desc:"员工手机号"`
	Position     string `json:"position" db:"position" desc:"职位"`
	Gender       int    `json:"gender" db:"gender" desc:"性别"`
	Email        string `json:"email" db:"email" desc:"员工邮箱"`
	Avatar       string `json:"avatar" db:"avatar" desc:"头像"`
	ThumbAvatar  string `json:"thumb_avatar" db:"thumb_avatar" desc:""`
	Telephone    string `json:"telephone" db:"telephone" desc:"电话"`
	Alias        string `json:"alias" db:"alias" desc:"别名"`
	Status       int8   `json:"status" db:"status" desc:"状态,1=已激活，2=已禁用，4=未激活，5=退出企业"`
	QrCode       string `json:"qr_code" db:"qr_code" desc:"二维码"`
	Address      string `json:"address" db:"address" desc:"地址"`
	OpenUserid   string `json:"open_userid" db:"open_userid" desc:""`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	Forms        string `json:"forms" desc:"所属部门"`
	Type         int    `json:"type" db:"type" desc:"3、企业导购，4、店长或企业管理"`
	// 非member 成员字段
	IsLeaderInDept []int8  `json:"is_leader_in_dept" desc:"表示在所在的部门内是否为上级"`
	Department     []int64 `json:"department" desc:"成员所属部门id列表，仅返回该应用有查看权限的部门id"`
	GenderText     string  `json:"gender_text" desc:"性别"`
	StatusText     string  `json:"status_text" desc:"状态"`
	// 分库特有字段
	MarsId         string `json:"mars_id" db:"mars_id" desc:"用户唯一标识，分库字段"`
	TempDepartment []byte `json:"-" db:"department" desc:"成员所属部门ID列表"`
	QrPosterCodeID int32  `json:"qr_poster_code_id" db:"qr_poster_code_id" desc:"员工海报id"`
	QrPosterCode   string `json:"qr_poster_code" db:"qr_poster_code" desc:"员工海报"`
}

// SimpleMember 基本企业成员信息
type SimpleMember struct {
	ID           int64  `json:"member_id" db:"id" desc:"员工ID"`
	Userid       string `json:"user_id" db:"userid" desc:"企业微信userid"`
	Name         string `json:"user_name" db:"name" desc:"员工姓名"`
	Position     string `json:"position" db:"position" desc:"职位"`
	Mobile       string `json:"mobile"  db:"mobile" desc:"员工手机号"`
	Status       int    `json:"mem_subscribe_status" db:"status" desc:"员工状态"`
	DepartmentId int64  `json:"department_id" db:"department_id" desc:"部门id"`
}

// MemberDetailRequest 成员详情请求
type MemberDetailRequest struct {
	ID           int64 `json:"id" desc:"员工ID"`
	EnterpriseID int64 `json:"enterprise_id" desc:"企业id"`
}

// BADetailRequest 成员详情请求
type BADetailRequest struct {
	UserId       string `json:"user_id" binding:"required"`
	UnionId      string `json:"union_id" binding:"required"`
	EnterpriseID int64  `json:"enterprise_id" binding:"required"`
}

// SyncContact 同步通讯录请求
type SyncContact struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" desc:"企业ID"`
	Domain       string `json:"domain" binding:"required" desc:"企业Domain"`
	AdminID      int64  `json:"admin_id" binding:"required" desc:"操作人ID"`
}

// MemberListRequest 列表请求数据
type MemberListRequest struct {
	ID           int64    `json:"id" db:"id" desc:"员工ID"`
	IDs          []int64  `json:"ids" desc:"多个员工ID"`
	IDCount      int      `json:"id_count" desc:"员工数"`
	UserIDs      []string `json:"user_ids" desc:"多个成员ID"`
	UserIDCount  int      `json:"user_id_count" desc:"成员数"`
	Userid       string   `json:"userid" db:"userid" desc:"企业微信userid"`
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DepartmentID int64    `json:"department_id" db:"department_id" desc:"部门id"`
	Status       int8     `json:"status" db:"status" desc:"激活状态1=激活 2=禁用 4=未关注"`
	Mobile       string   `json:"mobile" db:"mobile"`
	Name         string   `json:"name" db:"name" desc:"员工姓名"`
	Page         int      `json:"page" desc:"页码"`
	Size         int      `json:"size" desc:"分页"`
}

// MemberListResponse 返回数据
type MemberListResponse struct {
	MemberList []Member `json:"member_list"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Size       int      `json:"size"`
}

// MemberDepartment 成员组织关系
type MemberDepartment struct {
	DepartmentID      int64 `json:"department_id" db:"department_id" desc:"组织ID"`
	MemberID          int64 `json:"member_id" db:"member_id" desc:"成员ID"`
	EnterpriseID      int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt         int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt         int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt         int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	UserID            string
	DepartmentThirdID int64
}

// SyncContactResponse 同步通讯录返回结果
type SyncContactResponse struct {
	UpdateTime int64 `json:"updateTime" desc:"同步通讯录更新时间"`
}

// MemberInfo 整理用户数据
func MemberInfo(member Member) *Member {
	var genderArr = []string{0: "未定义", 1: "男", 2: "女"}
	var statusArr = []string{1: "已激活", 2: "已禁用", 4: "未激活", 5: "退出企业"}
	m := &Member{
		ID:           member.ID,
		Userid:       member.Userid,
		Name:         member.Name,
		Mobile:       member.Mobile,
		Position:     member.Position,
		Gender:       member.Gender,
		GenderText:   genderArr[member.Gender],
		Email:        member.Email,
		Avatar:       member.Avatar,
		ThumbAvatar:  member.ThumbAvatar,
		Telephone:    member.Telephone,
		Alias:        member.Alias,
		Status:       member.Status,
		StatusText:   statusArr[member.Status],
		QrCode:       member.QrCode,
		Address:      member.Address,
		OpenUserid:   member.OpenUserid,
		EnterpriseID: member.EnterpriseID,
		CreatedAt:    member.CreatedAt,
		UpdatedAt:    member.UpdatedAt,
		Forms:        member.Forms,
	}
	if m.Mobile != "" {
		m.Mobile = m.Mobile[:3] + "****" + m.Mobile[7:]
	}
	return m
}

// TranslateMember member 成员字段转换
func TranslateMember(md work.MemberDetail) *Member {
	gender, _ := strconv.Atoi(md.Gender)
	return &Member{
		Userid:         md.Userid,
		Name:           md.Name,
		Mobile:         md.Mobile,
		Address:        md.Address,
		Alias:          md.Alias,
		Avatar:         md.Avatar,
		Email:          md.Email,
		Status:         md.Status,
		QrCode:         md.QrCode,
		Position:       md.Position,
		Gender:         gender,
		Department:     md.Department,
		IsLeaderInDept: md.IsLeaderInDept,
	}
}

// BatchDeleteMember 批量删除成员
type BatchDeleteMember struct {
	DeletedAt int64 `db:"deleted_at"`
	IDs       []interface{}
}

// BatchGetMember 批量获取员工
type BatchGetMember struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DeletedAt    int64 `db:"deleted_at"`
	DepartIDs    []interface{}
}

// MomentListRequest 朋友圈创建选择员工请求
type MomentListRequest struct {
	GetNoDp      int64 `json:"get_no_dp" desc:"没有部门的员工"`
	Page         int64 `json:"page" desc:"第几页"`
	Size         int64 `json:"size" desc:"每页条数"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Search       struct {
		DpId        int64  `json:"dp_id" desc:"部门ID"`
		MemUsername string `json:"mem_username" desc:"名称"`
		Userid      string `json:"userid"`
		Mobile      string `json:"mobile"`
	} `json:"search" desc:"搜索条件"`
}

// MomentList 朋友圈创建选择员工返回
type MomentListResponse struct {
	List     []MomentList `json:"list"`
	NoDpList []MomentList `json:"no_dp_list"`
	Total    int64        `json:"total" desc:"total"`
	Page     int64        `json:"page" desc:"page"`
	Size     int64        `json:"size" desc:"size"`
}

// MomentList 朋友圈创建选择员工返回子结构体
type MomentList struct {
	UserId             string `json:"user_id" db:"user_id" desc:"members.userid"`
	Name               string `json:"user_name" db:"user_name" desc:"members.name"`
	MemberId           int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemSubscribeStatus int    `json:"mem_subscribe_status" db:"mem_subscribe_status" desc:"是不是子部门"`
	DeptName           string `json:"dept_name" db:"dept_name"`
	Mobile             string `json:"mobile" db:"mobile"`
}

// BatchGetMemberByMemberId 通过MemberIds批量获取员工
type BatchGetMemberByMemberId struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DeletedAt    int64 `db:"deleted_at"`
	IDs          []interface{}
}

// SyncMemberJob 同步任务消息
type SyncMemberJob struct {
	Job  string `json:"job"`
	ID   int64  `json:"id" db:"id" desc:"记录ID"`
	Data *SyncContact
}

// SyncMemberRelationJob 同步员工的外部联系人关系
type SyncMemberRelationJob struct {
	ID           int64  `json:"id" db:"id" desc:"员工的member表id"`
	Userid       string `json:"userid" db:"userid" desc:"员工的userid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	CorpID       string `json:"corp_id" db:"corp_id"`
}

// ExternalContactQueueData 外部联系人队列数据结构
type ExternalContactQueueData struct {
	Job  string `json:"job"`
	Data *work.MemberExternalCallback
}

// WorkExternalUserMigrationEntity 企微外部联系人迁移跟踪表结构
type WorkExternalUserMigrationEntity struct {
	ID             int64  `db:"id"`
	EnterpriseID   int64  `db:"enterprise_id"`
	CorpID         string `db:"corp_id"`
	UserID         string `db:"user_id"`
	ExternalUserID string `db:"external_user_id"`
	Event          string `db:"event"`
	ChatID         string `db:"chat_id"`
	Status         int    `db:"status"`
	CreatedAt      int64  `db:"created_at"`
	UpdatedAt      int64  `db:"updated_at"`
}

// AuthInfoRequest 成员详情请求
type AuthInfoRequest struct {
	Code         string `json:"code" desc:"企业授权code"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业id"`
	LoginType    string `json:"login_type" desc:"登录方式"`
}

// AuthInfoResponse 企业微信授权信息
type AuthInfoResponse struct {
	Member   AuthInfoMember   `json:"member"`
	External AuthInfoExternal `json:"external"`
}

// AuthInfoMember 内部成员信息
type AuthInfoMember struct {
	UserId   string `json:"user_id" desc:"成员ID"`
	DeviceId string `json:"device_id" desc:"手机设备号"`
	Name     string `json:"name"`
	MemberId int64  `json:"member_id"`
}

// AuthInfoExternal 外部成员信息
type AuthInfoExternal struct {
	OpenId         string `json:"openid" desc:"非企业成员的标识"`
	DeviceId       string `json:"device_id" desc:"手机设备号"`
	ExternalUserid string `json:"external_userid" desc:"外部联系人id"`
}

// BADetailResponse BA信息
type BADetailResponse struct {
	IsExist  int64    `json:"is_exist"`
	BADetail BADetail `json:"ba_detail"`
}

type BADetail struct {
	Name        string `json:"name,omitempty"`
	Position    string `json:"position,omitempty"`
	CounterName string `json:"counter_name,omitempty"`
	Marscode    string `json:"marscode,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	QrCode      string `json:"qr_code,omitempty"`
}

// Code2SessionRequest 成员详情请求
type Code2SessionRequest struct {
	JsCode       string `json:"js_code" desc:"企业授权code"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业id"`
	Domain       string `json:"domain" desc:"UC特有的企微id"`
}

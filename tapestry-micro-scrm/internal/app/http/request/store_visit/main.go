package storeVisitReq

import (
	"micro-scrm/internal/app/http/request"
)

type BaseReq struct {
	request.EnterpriseBaseReq
	MemberId int64  `json:"member_id"`
	Language string `json:"language"`
}

type PlanListReq struct {
	BaseReq
	DataType  int   `json:"data_type"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	ThirdId   int   `json:"third_id"`
	Status    int   `json:"status"`
	request.PageReq
}

const DataTypeCreator = 1
const DataTypeAll = 2

type VisitInfoReq struct {
	BaseReq
	ID int64 `json:"id"`
}

type SaveVisitInfoReq struct {
	VisitInfoReq
	ResetOverview int    `json:"reset_overview"`
	ResetContent  int    `json:"reset_content"`
	Overview      string `json:"overview"`
	Content       string `json:"content"`
}

type SaveReportPermissionReq struct {
	VisitInfoReq
	VisitorIsShow   int       `json:"visitor_is_show"`
	StoreIsShow     int       `json:"store_is_show"`
	AddVisitorList  []Visitor `json:"add_visitor_list"`
	AddStoreManager []int64   `json:"add_store_manager"`
	DelVisitorList  []Visitor `json:"del_visitor_list"`
}

type StoreManagerReq struct {
	BaseReq
	ThirdIds []int `json:"third_ids"`
}

type SavePlanReq struct {
	BaseReq
	StoreList      []StoreInfo `json:"store_list"`
	VisitorList    []Visitor   `json:"visitor_list"`
	SendMsgToStore int         `json:"send_msg_to_store"` // 是否通知店长
	Remark         string      `json:"remark"`            // 备注
}

type StoreInfo struct {
	Date         string  `json:"date"`       // 巡店日期
	TimeSpan     string  `json:"time_span"`  // 巡店时段
	UnixTime     int64   `json:"unix_time"`  //巡店预约开始的unix时间
	StoreCode    string  `json:"store_code"` // 门店编码
	ThirdId      int     `json:"third_id"`   //企微部门id
	StoreName    string  `json:"store_name"` // 门店名称
	StoreManager []int64 `json:"store_manager"`
}

type Visitor struct {
	MemberId int64  `json:"member_id"`
	Email    string `json:"email"`
}

type SearchMemberReq struct {
	BaseReq
	Str string `json:"str"`
}

type EmailContactReq struct {
	BaseReq
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type AddCopyMemberReq struct {
	BaseReq
	ID             int64     `json:"id"`
	CopyMemberList []Visitor `json:"copy_member_list"`
}

type StartVisitReq struct {
	VisitInfoReq
	Date     string `json:"date"`      // 巡店日期
	TimeSpan string `json:"time_span"` // 巡店时段
	UnixTime int64  `json:"unix_time"` //巡店预约开始的unix时间
	ThirdId  int    `json:"third_id"`
	CheckIn  bool   `json:"check_in"`
}

type SaveShowRecordContentReq struct {
	BaseReq
	ID       int64  `json:"id"`
	VisitID  int64  `json:"visit_id"`
	TypeName string `json:"type_name"`
	Content  string `json:"content"`
}
type SaveShowRecordMediaReq struct {
	BaseReq
	ID            int64  `json:"id"`
	VisitID       int64  `json:"visit_id"`
	TypeName      string `json:"type_name"`
	MediaType     int    `json:"media_type"`
	Url           string `json:"url"`
	VideoImageUrl string `json:"video_image_url"`
	VideoSeconds  int    `json:"video_second"`
}
type DelShowRecordReq struct {
	BaseReq
	ID int64 `json:"id"`
}

type DelShowRecordMediaReq struct {
	BaseReq
	MediaID int64 `json:"media_id"`
}

type AdminStoreVisitListReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	StoreName      string `json:"store_name"`
	DeptName       string `json:"dept_name"`
	Status         int    `json:"status"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	VisitStartTime int64  `json:"visit_start_time"`
	VisitEndTime   int64  `json:"visit_end_time"`
}

type AdminStoreDashboardReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	AreaId    int64  `json:"area_id"`
	DeptName  string `json:"dept_name"`
	Type      int    `json:"type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type SaveCommentReq struct {
	VisitInfoReq
	Content        string         `json:"content"`
	ReplyCommentId int64          `json:"reply_comment_id"`
	MediaList      []MediaInfoReq `json:"media_list"`
}

type CommentListReq struct {
	VisitInfoReq
	request.PageReq
}

type DelCommentReq struct {
	BaseReq
	CommentID int64 `json:"comment_id"`
}

type MediaInfoReq struct {
	MediaType     int    `json:"media_type"`
	Url           string `json:"url"`
	VideoImageUrl string `json:"video_image_url"`
	VideoSeconds  int    `json:"video_second"`
}

type DashboardTreeReq struct {
	BaseReq
	DataType  int     `json:"data_type"`
	StartTime int64   `json:"start_time"`
	EndTime   int64   `json:"end_time"`
	VisitIds  []int64 `json:"visit_ids"`
}

type DashboardExportEmail struct {
	DashboardTreeReq
	Email       string `json:"email"`
	SearchName  string `json:"search_name"`
	SearchLevel int    `json:"search_level"`
}

type CreatePdfReq struct {
	request.EnterpriseBaseReq
	ID       int64  `json:"id"`
	Language string `json:"language"`
}

type SendPdfToEmailReq struct {
	VisitInfoReq
	Email string `json:"email"`
}

type VisitConfigReq struct {
	BaseReq
	Key string `json:"key"`
}

const SendMsg = 1

const StoreVisitStatusVisiting = 1
const StoreVisitStatusVisited = 2

const (
	DashboardStore  = iota + 1 // 按照门店汇总
	DashboardSource            // 按照渠道汇总
	DashboardRegion            // 按照大区汇总
	DashboardArea              // 按照小区汇总
)

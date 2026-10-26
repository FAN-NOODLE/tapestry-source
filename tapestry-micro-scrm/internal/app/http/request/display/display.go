package display

import (
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
)

type AddTaskReq struct {
	request.EnterpriseBaseReq
	TaskInfoRes
}

type ClassInfo struct {
	ClassName string `json:"class_name"`
	Url       string `json:"url"`
}

type TemplateDetail struct {
	TemplateName string `json:"template_name"`
	IsRequired   bool   `json:"is_required"`
}

type StoreDepartmentInfo struct {
	DepartmentID       int                   `json:"department_id"`
	ThirdID            int                   `json:"third_id"`
	DpID               int                   `json:"dp_id"`
	DpThirdID          int                   `json:"dp_third_id"`
	DpThirdParentID    int                   `json:"dp_third_parent_id"`
	DpParentID         int                   `json:"dp_parent_id"`
	DpName             string                `json:"dp_name"`
	StoreType          int                   `json:"store_type"`
	StoreName          string                `json:"store_name"`
	IsChildDepartment  int                   `json:"is_child_department"`
	DpStoreName        string                `json:"dp_store_name"`
	ChildrenList       []StoreDepartmentInfo `json:"children_list"`
	MemSubscribeStatus int                   `json:"mem_subscribe_status"`
}

type SearchReq struct {
	request.EnterpriseBaseReq
	DisplayName    string `json:"display_name"`
	StartTime      int    `json:"start_time"`
	EndTime        int    `json:"end_time"`
	CreatedStartAt int    `json:"created_start_at"`
	CreatedEndAt   int    `json:"created_end_at"`
	Status         int32  `json:"status"`
	ShowType       int32  `json:"show_type"`
	request.PageReq
}

type SearchAllotBaseReq struct {
	request.EnterpriseBaseReq
	DisplayTaskId   int64  `json:"display_task_id"`
	DisplayName     string `json:"display_name"`
	OwnerMemberID   int64  `json:"owner_member_id"`
	SubmitStartTime int    `json:"submit_start_time"`
	SubmitEndTime   int    `json:"submit_end_time"`
	FinishStatus    int32  `json:"finish_status"`
	StoreId         int    `json:"store_id"`
	StoreName       string `json:"store_name"`
}

type SearchAllotReq struct {
	SearchAllotBaseReq
	request.PageReq
}

type TaskBaseReq struct {
	ID int64 `json:"id"`
	request.EnterpriseBaseReq
}

type TaskInfoRes struct {
	ID           int64                 `json:"id"`
	DisplayName  string                `json:"display_name"`
	StoreType    int32                 `json:"store_type" `
	Content      string                `json:"content" `
	ShowType     int32                 `json:"show_type"`
	Status       int32                 `json:"status"`
	StartTime    int                   `json:"start_time"`
	EndTime      int                   `json:"end_time"`
	CreatedAt    int64                 `json:"created_at"`
	ClassURL     []ClassInfo           `json:"class_url"`
	TemplateInfo []TemplateDetail      `json:"template_info"`
	StoreList    []StoreDepartmentInfo `json:"store_list"`
	ReviewType   int                   `json:"review_type"`
	ReviewFile   string                `json:"review_file"`
}

type StoreTaskSaveInfoReq struct {
	request.EnterpriseBaseReq
	ID           int64  `json:"id"`
	MemberId     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	MemberUserId string `json:"member_user_id"`
	SubmitType   int    `json:"submit_type"`
	StoreTaskSaveInfo
}

type StoreTaskSaveInfo struct {
	AllotID           int64       `json:"allot_id" comment:"任务分配表id"`
	OwnerMemberId     int64       `json:"owner_member_id"`
	OwnerMemberUserId string      `json:"owner_member_user_id"`
	OwnerMemberName   string      `json:"owner_member_name"`
	FinishStatus      int32       `json:"finish_status"`
	Grade             uint32      `json:"grade"`
	ReviewRemark      string      `json:"review_remark"`
	StoreCode         string      `json:"store_code"`
	StoreName         string      `json:"store_name"`
	SubmitTime        int64       `json:"submit_time"`
	ImageList         []ImageInfo `json:"image_list"`
}

type ImageInfo struct {
	ID         int64         `json:"id"`
	Identifier string        `json:"identifier"`
	Url        string        `json:"url"`
	Status     int           `json:"status"`
	Num        int           `json:"num"`
	VmNum      int           `json:"vm_num"`
	Comment    []CommentInfo `json:"comment"`
	IsLook     int           `json:"is_look"`
	IsDraft    int           `json:"is_draft"`
}

type CommentInfo struct {
	ID           int64      `json:"id"`
	MemberId     int64      `json:"member_id"`
	MemberUserId string     `json:"member_user_id"`
	MemberName   string     `json:"member_name"`
	Avatar       string     `json:"avatar"`
	StoreCode    string     `json:"store_code"`
	Content      RemarkInfo `json:"content"`
	ReplyName    string     `json:"reply_name" comment:"回复人名称"`
	CommentTime  int        `json:"comment_time"`
	IsVm         int        `json:"is_vm"`
	IsSelf       int        `json:"is_self"`
	Action       int        `json:"action"`
}

type RemarkInfo struct {
	ImageUrl []string `json:"image_url,omitempty"`
	Text     string   `json:"text,omitempty"`
}

type MemberTaskListReq struct {
	MemberTaskReq
	request.PageReq
	Name         string `json:"name"`
	FinishStatus int    `json:"finish_status"`
	Status       int    `json:"status"`
}

type MemberTaskReq struct {
	TaskBaseReq
	MemberId     int64  `json:"member_id"`
	MemberUserId string `json:"member_user_id"`
}

type MemberStoreTaskInfoReq struct {
	MemberTaskReq
	StoreId int  `json:"store_id"`
	IsAdmin bool `json:"is_admin"`
}

type StoreTaskDetail struct {
	TaskInfo  TaskInfoRes       `json:"task_info"`
	StoreInfo StoreTaskSaveInfo `json:"store_info"`
}

type MemberTaskListRes struct {
	AllotId      int64  `json:"allot_id"`
	Id           int64  `json:"id"`
	DisplayName  string `json:"display_name"`
	FinishStatus int    `json:"finish_status"`
	Status       int    `json:"status"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	Content      string `json:"content"`
	StoreId      int    `json:"store_id"`
	Grade        int    `json:"grade"`
	ClassUrl     string `json:"class_url"`
}

type StoreBaseInfo struct {
	StoreCode     string `json:"store_code"`
	StoreId       int    `json:"store_id"`
	StoreName     string `json:"store_name"`
	StoreType     int    `json:"store_type"`
	StoreCity     string `json:"store_city"`
	SpecialUserId string `json:"special_user_id"`
}

type LookReq struct {
	request.EnterpriseBaseReq
	AllotID   int64 `json:"allot_id"`
	ImageID   int64 `json:"image_id"`
	MemberId  int64 `json:"member_id"`
	CommentId int64 `json:"comment_id"`
}

type LookRes struct {
	AllotDataId int64 `json:"allot_data_id"`
	CommentId   int64 `json:"comment_id"`
}

type AuthReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type AuthInfo struct {
	IsVm      int `json:"is_vm"`
	IsManager int `json:"is_manager"`
	IsSpecial int `json:"is_special"`
}

type VmListRes struct {
	MemberId   int64  `json:"member_id"`
	UserId     string `json:"user_id"`
	MemberName string `json:"member_name"`
}

type StoreSearchReq struct {
	TaskBaseReq
	Name    string `json:"name"`
	StoreId int    `json:"store_id"`
}

type VmStoreBase struct {
	StoreId   int    `json:"store_id"`
	StoreName string `json:"store_name"`
}

type VmStore struct {
	StoreCity string        `json:"store_city"`
	StoreInfo []VmStoreBase `json:"store_info"`
}

type ListTips struct {
	AllNum         int64 `json:"all_num" comment:"全部的数量"`
	WaitNum        int64 `json:"wait_num" comment:"待提交中"`
	CheckNum       int64 `json:"check_num" comment:"审核中"`
	DoubleCheckNum int64 `json:"double_check_num" comment:"待复审"`
	FailNum        int64 `json:"fail_num" comment:"未通过"`
	SuccessNum     int64 `json:"success_num" comment:"已通过"`
}

type ReviewStatus struct {
	DisplayTaskId int64 `json:"display_task_id"`
	AllotListCount
}
type AllotListCount struct {
	FinishStatus int   `json:"finish_status"`
	Num          int64 `json:"num"`
}

type CheckListRes struct {
	TaskInfoRes
	ListTips
}

type ListTipsReq struct {
	TaskBaseReq
	MemberId     int64  `json:"member_id"`
	IsVm         int    `json:"is_vm" comment:"VM查询的身份信息"`
	MemberUserId string `json:"member_user_id"`
}
type StoreListInfoReq struct {
	MemberTaskReq
	Name         string `json:"name"`
	FinishStatus int    `json:"finish_status"`
}

type StoreListInfo struct {
	ID           int64  `json:"id" comment:"任务id"`
	AllotID      int64  `json:"allot_id" comment:"分配id"`
	StoreId      int    `json:"store_id"`
	StoreName    string `json:"store_name"`
	StoreCode    string `json:"store_code"`
	FinishStatus int32  `json:"finish_status"`
	Grade        uint32 `json:"grade"`
}

type StoreTaskInfo struct {
	ID                int64       `json:"id"`
	DisplayName       string      `json:"display_name"`
	AllotID           int64       `json:"allot_id" comment:"任务分配表id"`
	OwnerMemberId     int64       `json:"owner_member_id"`
	OwnerMemberUserId string      `json:"owner_member_user_id"`
	OwnerMemberName   string      `json:"owner_member_name"`
	FinishStatus      int32       `json:"finish_status"`
	Grade             uint32      `json:"grade"`
	ReviewRemark      string      `json:"review_remark"`
	StoreCode         string      `json:"store_code"`
	StoreName         string      `json:"store_name"`
	SubmitTime        int64       `json:"submit_time"`
	ReviewTime        int64       `json:"review_time"`
	ImageList         []ImageInfo `json:"image_list"`
	SpecialUserId     string      `json:"special_user_id"`
	SpecialName       string      `json:"special_name"`
	IsSpecialCheck    bool        `json:"is_special_check"`
	IsCorporateCheck  bool        `json:"is_corporate_check"`
}

type StoreTaskReviewBasReq struct {
	TaskBaseReq
	ReviewRemark string `json:"review_remark"`
	Grade        uint32 `json:"grade"`
	ReviewType   int    `json:"review_type"`
	StoreTaskSaveInfo
}

type StoreTaskReviewReq struct {
	StoreTaskReviewBasReq
	MemberId     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	MemberUserId string `json:"member_user_id"`
}

type StoreTaskReviewWebReq struct {
	StoreTaskReviewBasReq
	AdminId   int64  `json:"admin_id"`
	AdminName string `json:"admin_name"`
}

type SendNoticeReq struct {
	TaskBaseReq
	MemberId int64 `json:"member_id"`
	AllotID  int64 `json:"allot_id"`
}

type AllotListRes struct {
	AllotId         int64  `json:"allot_id"`
	DisplayTaskId   int64  `json:"display_task_id"`
	DisplayName     string `json:"display_name"`
	StoreId         int    `json:"store_id"`
	StoreCode       string `json:"store_code"`
	StoreName       string `json:"store_name"`
	StoreType       int32  `json:"store_type"`
	StoreCity       string `json:"store_city"`
	FinishStatus    int32  `json:"finish_status"`
	Grade           uint32 `json:"grade"`
	OwnerMemberId   int64  `json:"owner_member_id"`
	OwnerMemberName string `json:"owner_member_name"`
	SubmitTime      int64  `json:"submit_time"`
	ReviewTime      int64  `json:"review_time"`
	SubmitName      string `json:"submit_name"`
	SubmitMemberId  int64  `json:"submit_member_id"`
	ReviewName      string `json:"review_name"`
}

type AllotListExportRes struct {
	AllotListRes
	SubmitTimeStr   string `json:"submit_time_str"`
	ReviewTimeStr   string `json:"review_time_str"`
	GradeStr        string `json:"grade_str"`
	StoreTypeStr    string `json:"store_type_str"`
	FinishStatusStr string `json:"finish_status_str"`
}

type NoticeSettingReq struct {
	request.EnterpriseBaseReq
	Value uint `json:"value"`
}

type AllotSendNoticeInfo struct {
	DisplayTaskId int64 `json:"display_task_id"`
	StoreId       int   `json:"store_id"`
	FinishStatus  int   `json:"finish_status"`
}

type SaveImageReq struct {
	request.EnterpriseBaseReq
	AllotId     int64  `json:"allot_id"`
	StoreCode   string `json:"store_code"`
	StoreName   string `json:"store_name"`
	DisplayName string `json:"display_name"`
}

type SaveStoreImageReq struct {
	request.EnterpriseBaseReq
	DisplayTaskId int64 `json:"display_task_id"`
	AllotId       int64 `json:"allot_id"`
	ImageInfo
}

type DelImageReq struct {
	request.EnterpriseBaseReq
	Id int64 `json:"id"`
}

type SaveStoreCommentReq struct {
	request.EnterpriseBaseReq
	DisplayTaskId int64 `json:"display_task_id"`
	AllotId       int64 `json:"allot_id"`
	AllotDataId   int64 `json:"allot_data_id"`
	CommentInfo
}

type SimpleTaskInfo struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
}

type DashBoard struct {
	request.EnterpriseBaseReq
	Ids     []int64 `json:"ids" comment:"任务表id"`
	AreaIds []int64 `json:"area_ids" comment:"筛选的组织树id "`
}

type DashBoardInfo struct {
	ListTips
	WaitRate        string `json:"wait_rate" comment:"待提交的占比"`
	CheckRate       string `json:"check_rate" comment:"审核中的占比"`
	FailRate        string `json:"fail_rate" comment:"未通过的占比"`
	SuccessRate     string `json:"success_rate" comment:"已通过的占比"`
	DoubleCheckRate string `json:"double_check_rate" comment:"待复审的占比"`
}

type VmFailList struct {
	ID          int64        `json:"id"`
	DisplayName string       `json:"display_name"`
	List        []VmFailInfo `json:"list"`
}

type VmFailInfo struct {
	AllotId    int64  `json:"allot_id"`
	StoreName  string `json:"store_name"`
	StoreId    int    `json:"store_id"`
	ReviewTime int64  `json:"review_time"`
}

type DashBoardListReq struct {
	Ids          []int64 `json:"ids"`
	FinishStatus int     `json:"finish_status"`
	dashboardReq.ManagerAuthAreaIds
}

type StoreCount struct {
	StoreId int `json:"store_id"`
	Num     int `json:"num"`
}

type DashBoardListRes struct {
	dashboardResp.AreaInfo
	StoreNum int `json:"store_num" comment:"门店数"`
	TaskNum  int `json:"task_num" comment:"任务数"`
}

type StoreListReq struct {
	request.EnterpriseBaseReq
	Ids          []int `json:"ids"`
	StoreId      int   `json:"store_id"`
	FinishStatus int   `json:"finish_status"`
}

type StoreListRes struct {
	ID           int64  `json:"id"`
	DisplayName  string `json:"display_name"`
	StoreType    int32  `json:"store_type" `
	Content      string `json:"content" `
	ShowType     int32  `json:"show_type"`
	Status       int32  `json:"status"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	AllotId      int64  `json:"allot_id"`
	StoreId      int    `json:"store_id"`
	FinishStatus int32  `json:"finish_status"`
}

type StoreDashBoard struct {
	StoreId int `json:"store_id"`
	AllotListCount
}

type SpecialInfo struct {
	ID        int64  `json:"id"`
	UserId    string `json:"user_id"`
	StoreCode string `json:"store_code"`
	StoreType int32  `json:"store_type"`
	ThirdId   int64  `json:"third_id"`
	Num       int    `json:"num"`
}

type ReviewImportData struct {
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
	StoreCode string `json:"store_code"`
}

type ReviewData struct {
	ReviewImportData
	FailReason string `json:"fail_reason"`
}
type ReviewFileReq struct {
	request.EnterpriseBaseReq
	FileUrl string `json:"file_url"`
}

type ReviewImportResp struct {
	UploadNum    int    `json:"upload_num"`
	Success      int    `json:"success"`
	DuplicateNum int    `json:"duplicate_num"`
	FailedNum    int    `json:"failed_num"`
	FailUrl      string `json:"fail_url"`
}

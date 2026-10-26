package display

import (
	"gw-scrm/internal/app/http/request"
)

type AddTaskReq struct {
	request.EnterpriseBaseReq
	ID           int64                 `json:"id"`
	DisplayName  string                `json:"display_name"`
	StoreType    int32                 `json:"store_type" `
	Content      string                `json:"content" `
	ClassURL     []ClassInfo           `json:"class_url"`
	TemplateInfo []TemplateDetail      `json:"template_info" `
	ShowType     int32                 `json:"show_type"`
	Status       int32                 `json:"status"`
	StartTime    int                   `json:"start_time"`
	EndTime      int                   `json:"end_time"`
	StoreList    []StoreDepartmentInfo `json:"store_list"`
	ReviewType   int                   `json:"review_type"`
	ReviewFile   string                `json:"review_file"`
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

type TaskBaseReq struct {
	ID int64 `json:"id"`
	request.EnterpriseBaseReq
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

type NoticeSettingReq struct {
	request.EnterpriseBaseReq
	Value uint `json:"value"`
}

type AllotBaseReq struct {
	TaskBaseReq
	StoreId int  `json:"store_id" binding:"required"`
	IsAdmin bool `json:"is_admin"`
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

type StoreTaskReviewReq struct {
	TaskBaseReq
	AdminId      int64  `json:"admin_id"`
	AdminName    string `json:"admin_name"`
	ReviewRemark string `json:"review_remark"`
	Grade        uint32 `json:"grade"`
	ReviewType   int    `json:"review_type"`
	StoreTaskSaveInfo
}

type SaveImageReq struct {
	request.EnterpriseBaseReq
	AllotId     int64  `json:"allot_id"`
	StoreCode   string `json:"store_code"`
	StoreName   string `json:"store_name"`
	DisplayName string `json:"display_name"`
}

type CosAvatar struct {
	Url string `json:"url" binding:"required" required:"图片地址必填"`
}

type ReviewFileReq struct {
	request.EnterpriseBaseReq
	FileUrl string `json:"file_url"`
}

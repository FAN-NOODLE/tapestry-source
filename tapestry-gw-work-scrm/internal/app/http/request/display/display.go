package display

import (
	"gw-work-scrm/internal/app/http/request"
)

type TaskBaseReq struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" comment:"任务id"`
}

type MemberStoreTaskInfoReq struct {
	TaskBaseReq
	MemberId int64 `json:"member_id"`
	StoreId  int   `json:"store_id"`
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

type StoreTaskSaveInfoReq struct {
	TaskBaseReq
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
	ImageList         []ImageInfo `json:"image_list"`
}

type ImageInfo struct {
	ID         int64         `json:"id"`
	Identifier string        `json:"identifier"`
	Url        string        `json:"url"`
	Status     int           `json:"status" comment:"vm的审核状态，0为通过，1为不通过"`
	Comment    []CommentInfo `json:"comment"`
}

type CommentInfo struct {
	ID           int64      `json:"id"`
	MemberId     int64      `json:"member_id"`
	MemberUserId string     `json:"member_user_id"`
	MemberName   string     `json:"member_name"`
	Avatar       string     `json:"avatar"`
	ReplyName    string     `json:"reply_name" comment:"回复人名称"`
	CommentTime  int        `json:"comment_time"`
	StoreCode    string     `json:"store_code"`
	Content      RemarkInfo `json:"content"`
	IsVm         int        `json:"is_vm"`
	IsSelf       int        `json:"is_self"`
	Action       int        `json:"action"`
}

type RemarkInfo struct {
	ImageUrl []string `json:"image_url,omitempty"`
	Text     string   `json:"text,omitempty"`
}

type LookReq struct {
	request.EnterpriseBaseReq
	AllotID   int64 `json:"allot_id"`
	ImageID   int64 `json:"image_id"`
	MemberId  int64 `json:"member_id"`
	CommentId int64 `json:"comment_id"`
}

type AuthReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type StoreSearchReq struct {
	TaskBaseReq
	Name    string `json:"name"`
	StoreId int    `json:"store_id"`
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

type StoreTaskReviewReq struct {
	TaskBaseReq
	MemberId     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	MemberUserId string `json:"member_user_id"`
	ReviewRemark string `json:"review_remark"`
	Grade        uint32 `json:"grade"`
	ReviewType   int    `json:"review_type"`
	StoreTaskSaveInfo
}

type SendNoticeReq struct {
	TaskBaseReq
	AllotID  int64 `json:"allot_id"`
	MemberId int64 `json:"member_id"`
}

type WeLearnToken struct {
	LoginStr string `json:"login_str" comment:"登录的str"`
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

type DashBoard struct {
	request.EnterpriseBaseReq
	IDs     []int64 `json:"ids" comment:"任务表id"`
	AreaIds []int64 `json:"area_ids" comment:"筛选的组织树id"`
}

type DashBoardListReq struct {
	request.EnterpriseBaseReq
	Ids          []int64 `json:"ids"`
	FinishStatus int     `json:"finish_status"`
	AreaIds      []int64 `json:"area_ids"`
	TargetAreaId int64   `json:"target_area_id"`
}

type StoreListReq struct {
	request.EnterpriseBaseReq
	Ids          []int `json:"ids"`
	StoreId      int   `json:"store_id"`
	FinishStatus int   `json:"finish_status"`
}

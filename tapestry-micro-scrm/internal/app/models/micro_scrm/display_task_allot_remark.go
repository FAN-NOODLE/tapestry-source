package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	RemarkDel = iota + 1
	RemarkEdit
)

type DisplayTaskAllotRemarkModel struct {
	models.BaseModel
	DisplayTaskID int64  `json:"display_task_id" gorm:"column:display_task_id;type:bigint(20);not null;default:0;comment:任务表id"`
	AllotID       int64  `json:"allot_id" gorm:"column:allot_id;type:bigint(20);not null;default:0;comment:分配表id"`
	AllotDataID   int64  `json:"allot_data_id" gorm:"column:allot_data_id;type:bigint(20);not null;index:allot_data_id_index;default:0;comment:VM任务图片表id"`
	MemberID      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购信息"`
	MemberUserID  string `json:"member_user_id" gorm:"column:member_user_id;type:varchar(255);not null;comment:导购工号"`
	MemberName    string `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:导购名称"`
	Avatar        string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;comment:导购头像"`
	Content       string `json:"content" gorm:"column:content;type:text;not null;comment:备注内容"`
	IsVm          int    `json:"is_vm" gorm:"column:is_vm;type:tinyint;not null;default:0;comment:是否是VM备注"`
	CommentTime   int    `json:"comment_time" gorm:"column:comment_time;type:int(10);not null;index:comment_time_index;default:0;comment:备注时间"`
	StoreCode     string `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:门店code"`
	ReplyName     string `json:"reply_name" gorm:"column:reply_name;type:varchar(255);not null;comment:回复人名称"`
}

// TableName returns the table name of the DisplayTaskAllotReamrk model
func (DisplayTaskAllotRemarkModel) TableName() string {
	return "display_task_allot_remark"
}

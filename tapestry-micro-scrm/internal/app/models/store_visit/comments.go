package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitCommentModel struct {
	models.BaseModel
	VisitID       int64                    `json:"visit_id" gorm:"column:visit_id;type:bigint(20);not null;index:visit_id_index;comment:关联的巡店ID"`
	MemberID      int64                    `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;comment:评论用户ID"`
	Content       string                   `json:"content" gorm:"column:content;type:text;not null;comment:评论内容"`
	ParentID      int64                    `json:"parent_id" gorm:"column:parent_id;type:bigint(20);not null;default:0;index:parent_id_index;comment:父评论ID"`
	ReplyMemberID int64                    `json:"reply_member_id" gorm:"column:reply_member_id;type:bigint(20);comment:回复目标用户ID"`
	IsDel         int                      `json:"is_del" gorm:"column:is_del;type:tinyint(1);not null;default:0;comment:是否已删除 0否 1是"`
	MediaList     []StoreVisitMediaModel   `json:"media_list" gorm:"foreignKey:RelID;references:ID"`
	Children      []StoreVisitCommentModel `json:"children" gorm:"foreignKey:ParentID;references:ID"`
}

func (StoreVisitCommentModel) TableName() string {
	return "store_visit_comment"
}

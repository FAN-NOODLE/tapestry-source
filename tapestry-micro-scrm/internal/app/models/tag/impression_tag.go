package tagModels

import "micro-scrm/internal/app/models"

type ImpressionTagModel struct {
	GroupID      int64  `json:"group_id" gorm:"column:group_id;type:bigint(20);not null;comment:外部联系人印象标签组表id"`
	ThirdGroupID int64  `json:"third_group_id" gorm:"column:third_group_id;type:bigint(20);not null;index:group_id_index;comment:三方标签组id"`
	ThirdTagID   string `json:"third_tag_id" gorm:"column:third_tag_id;type:varchar(255);not null;index:tag_id_index;comment:三方标签id"`
	Name         string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:标签名称"`
	Description  string `json:"description" gorm:"column:description;type:varchar(255);not null;comment:标签描述"`

	models.MicroBaseModel

	UseNum int64 `json:"use_num" comment:"标签使用统计"`
}

type TagsAndGroupTagListRes struct {
	ImpressionTagModel
	ImpressionTagGroupModel
	GroupID int64 `json:"group_id"`
	TagID   int64 `json:"tag_id"`
	IsHave  int   `json:"is_have"`
}

type UserTagsRes struct {
	GroupInfoRes
	ImpressionTagsInfo
}

type UserTagsGroupList struct {
	GroupInfoRes
	//TagsList []TagsInfo `json:"tags_list"`
	ImpressionTagsList []ImpressionTagsInfo `json:"impression_tags_list"`
}

type ImpressionTagsInfo struct {
	TagId  int64  `json:"tag_id"`
	Name   string `json:"name"`
	IsHave int    `json:"is_have"`
}

type GroupBaseInfo struct {
	GroupId   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

type GroupInfoRes struct {
	GroupBaseInfo
}

// TableName 外部联系人印象标签表
func (ImpressionTagModel) TableName() string {
	return "impression_tag"
}

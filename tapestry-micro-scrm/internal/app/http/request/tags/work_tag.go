package tagsReq

type PageVO struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}
type InfoBaseReq struct {
	EnterpriseID int64 `json:"enterprise_id"  comment:"企业ID"`
}

type InfoReq struct {
	InfoBaseReq
	TagID   int64 `json:"tag_id"`
	TagType int   `json:"tag_type"`
	PageVO
}

type UserList struct {
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	UnionID string `json:"union_id"`
}

type TagListInfo struct {
	Total    int64      `json:"total"`
	UserList []UserList `json:"user_list"`
}

type ThirdTagsInfo struct {
	ID             string `json:"id"`
	TagGroup       string `json:"tag_group"`
	TagName        string `json:"tag_name"`
	ActionType     int    `json:"action_type"`
	TagDescription string `json:"tag_description"`
}

type UpdateTagsInfo struct {
	Name    string `json:"name"`
	ThirdId string `json:"third_id"`
}

type ExternalUserThirdTag struct {
	ThirdTagID string `json:"third_tag_id"`
	TagName    string `json:"tag_name"`
	Mobile     string `json:"mobile"`
	UnionId    string `json:"union_id"`
	ActionType int    `json:"action_type"`
}

type ExternalUserTag struct {
	TagId   string `json:"third_tag_id"`
	UnionId string `json:"union_id"`
	UserId  string `json:"user_id"`
}

type BatchMarkingTagReq struct {
	EnterpriseID        int64             `json:"enterprise_id"  comment:"企业ID"`
	TaskID              int64             `json:"task_id"  comment:"任务ID"`
	ExternalUserTagList []ExternalUserTag `json:"external_user_tag_list"`
}

type UserTagsReq struct {
	InfoBaseReq
	WorkExternalUserId int64    `json:"work_external_user_id"  comment:"外部联系人user_id"`
	MemberId           int64    `json:"member_id"  comment:"member_id"`
	MemberIds          []int64  `json:"member_ids" `
	EnterpriseUserId   int64    `json:"enterprise_user_id" db:"enterprise_user_id"`
	Userid             string   `json:"userid"`
	Userids            []string `json:"userids"`
}

type UserTagsRes struct {
	GroupInfoRes
	TagsInfo
}

type UserTagsGroupList struct {
	GroupInfoRes
	TagsList []TagsInfo `json:"tags_list"`
}

type TagsInfo struct {
	TagId               string `json:"tag_id"`
	Name                string `json:"name"`
	ThirdTagDescription string `json:"third_tag_description"`
}

type GroupBaseInfo struct {
	GroupId    int64  `json:"group_id"`
	GroupName  string `json:"group_name"`
	IsThirdTag int64  `json:"is_third_tag"`
}

type GroupInfoRes struct {
	// RelationId    int64 `json:"relation_id"`
	ExternalTagId int64 `json:"external_tag_id"`
	Type          int   `json:"type"`
	GroupBaseInfo
}

type ListReq struct {
	InfoBaseReq
}

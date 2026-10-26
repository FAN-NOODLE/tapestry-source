package microscrm

import "gw-scrm/internal/pkg/sdk"

// TagListResponse 标签列表返回参数
type TagListResponse struct {
	sdk.ResponseBase
	Data TagListData `json:"data"`
}

type TagListData struct {
	Auto   []TagGroups `json:"auto" desc:"自动标签组"`
	Manual []TagGroups `json:"manual" desc:"手动标签组"`
}

// TagGroups 标签组
type TagGroups struct {
	TagGroupId   int64  `json:"tag_group_id" desc:"标签组ID"`
	TagGroupName string `json:"tag_group_name" desc:"标签组名称"`
	Tags         []Tag  `json:"tags" desc:"标签组下的标签列表"`
}

// Tag 标签
type Tag struct {
	Name  string `json:"name" desc:"标签名称"`
	TagId int64  `json:"tag_id" desc:"标签ID"`
	Value string `json:"value" desc:"标签人数"`
}

type TagInfoResp struct {
	sdk.ResponseBase
	Data TagInfoData `json:"data"`
}

type TagInfoData struct {
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
	List  []TagInfoListData `json:"list"`
}

type TagInfoListData struct {
	WorkTagId   string `json:"work_tag_id"`
	TagId       int64  `json:"tag_id"`
	TagName     string `json:"tag_name"`
	UseNum      int    `json:"use_num"`
	IsDelete    int    `json:"is_delete"`
	CreateTime  int    `json:"create_time"`
	ExpiredTime int    `json:"expired_time"`
	TagGroupId  int64  `json:"tag_group_id"`
}

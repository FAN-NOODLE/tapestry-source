package tagsReq

type GroupSearchList struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Keyword      string `json:"keyword" comment:"搜索关键字"`
	SearchType   int    `json:"search_type" binding:"required,min=1" comment:"搜索类型" desc:"1.标签组,2.标签"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type TagList struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	TagGroupId   int64  `json:"tag_group_id" comment:"群组id"`
	Keyword      string `json:"keyword" comment:"搜索关键字"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=6,max=99999999" comment:"分页数"`
}

type CreateGroup struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	ThirdGroupId string `json:"third_group_id" comment:"三方标签组id"`
	GroupName    string `json:"group_name" comment:"标签组名称"`
}

type CreateTag struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	ThirdTagId   string `json:"third_tag_id" comment:"三方标签id"`
	Name         string `json:"name" comment:"标签名称"`
	Description  string `json:"description" comment:"标签描述"`
	GroupId      int64  `json:"group_id" comment:"分组id"`
}

type CreateGroupAndTag struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	ThirdGroupId string `json:"third_group_id" comment:"三方标签组id"`
	ThirdTagId   string `json:"third_tag_id" binding:"required" comment:"三方标签id"`
	GroupName    string `json:"group_name" binding:"required" comment:"标签组名称"`
	Name         string `json:"name" binding:"required" comment:"标签名称"`
	Description  string `json:"description" binding:"required" comment:"标签描述"`
}

type ReqList struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

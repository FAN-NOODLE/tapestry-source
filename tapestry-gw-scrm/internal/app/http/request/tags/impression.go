package tagsReq

type GroupSearchList struct {
	EnterpriseID int64  `json:"enterprise_id" comment:"企业ID"`
	Keyword      string `json:"keyword" comment:"搜索关键字"`
	SearchType   int    `json:"search_type" comment:"搜索类型" desc:"1.标签组,2.标签"`
	Page         int    `json:"page" comment:"页码"`
	Size         int    `json:"size" comment:"分页数"`
}

type TagList struct {
	EnterpriseID int64  `json:"enterprise_id" comment:"企业ID"`
	TagGroupId   int64  `json:"tag_group_id" comment:"群组id"`
	Keyword      string `json:"keyword" comment:"搜索关键字"`
	Page         int    `json:"page" comment:"页码"`
	Size         int    `json:"size" comment:"分页数"`
}

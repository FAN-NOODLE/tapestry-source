package entity

// CustomerGroupDetail 客户分组详情
type CustomerGroupDetail struct {
	EnterpriseID    int64            `json:"enterprise_json" desc:"企业ID"`
	ID              int64            `json:"id" desc:"客户分组ID"`
	Name            string           `json:"name" desc:"客户分组名称"`
	IncludeTagGroup []TagGroupDetail `json:"include_tag_group" desc:"包含的客户标签组"`
	FilterTagGroup  []TagGroupDetail `json:"filter_tag_group" desc:"过滤的客户标签组"`
	Type            int32            `json:"type" desc:"客户分组类型，分为三类：1-忠诚度模型AIPL,2-RFM模型，3-自定义"`
}

// TagGroupDetail 标签组详情
type TagGroupDetail struct {
	EnterpriseID int64       `json:"enterprise_json" desc:"企业ID"`
	ID           int64       `json:"id" desc:"标签组ID"`
	Name         string      `json:"name" desc:"标签组名称"`
	WwID         string      `json:"ww_id" desc:"企微标签组ID"`
	TagList      []TagDetail `json:"tag_list" desc:"标签组内的标签集合"`
	Order        int         `json:"order" desc:"标签组排序标志，同步企微order"`
}

// TagDetail 标签详情
type TagDetail struct {
	EnterpriseID int64  `json:"enterprise_json" desc:"企业ID"`
	ID           int64  `json:"id" desc:"标签ID"`
	Name         string `json:"name" desc:"标签名称"`
	WwID         string `json:"ww_id" desc:"企微标签ID"`
	Type         int32  `json:"type" desc:"标签类型，分为两类：1-企业设置, 2-用户自定义"`
	TagGroupID   int64  `json:"tag_group_id" desc:"标签分组ID"`
	Order        int    `json:"order" desc:"标签组内标签排序标志，同步企微order"`
}

// 客户分组
type CustomerGroup struct {
	ID   int64  `json:"id" desc:"客户分组ID"`
	Name string `json:"name" desc:"客户分组名称"`
	Type int32  `json:"type" desc:"客户分组类型，分为三类：1-忠诚度模型AIPL,2-RFM模型，3-自定义"`
}

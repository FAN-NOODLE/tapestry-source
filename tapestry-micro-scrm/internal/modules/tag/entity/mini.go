package entity

// 小程序请求参数

type MiniBaseParam struct {
	EnterpriseID int64 `json:"enterprise_id" desc:"企业ID"`
}

type MiniTagGroupDetail struct {
	ID        int64           `json:"id"`
	GroupID   string          `json:"group_id" desc:"企业微信的标签ID"`
	GroupName string          `json:"group_name"`
	Order     int             `json:"order"`
	TagList   []MiniTagDetail `json:"tag_list"`
}

type MiniTagDetail struct {
	ID    int64  `json:"id"`
	TagID string `json:"tag_id" desc:"企业微信标签ID"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

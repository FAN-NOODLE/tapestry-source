package customerResp

type UserList struct {
	Id        int64  `json:"id"`
	Avatar    string `json:"avatar"`
	Birthday  string `json:"birthday"`
	Grade     string `json:"grade"`
	Mobile    string `json:"mobile"`
	Name      string `json:"name"`
	Nick      string `json:"nick"`
	UnionCode string `json:"union_code"`
	UnionId   string `json:"union_id"`
}

// UserDetail 客户详情
type UserDetail struct {
	UserID        int64    `json:"user_id"`
	UnionID       string   `json:"union_id"`
	UserName      string   `json:"user_name"`
	Nickname      string   `json:"nick_name"`
	Avatar        string   `json:"avatar"`
	Gender        int      `json:"gender"`
	TagList       []string `json:"tag_list"`
	ImpressionTag []string `json:"impression_tag"`
	GroupList     []string `json:"group_list"`
	Birthday      string   `json:"birthday"`
	UnionCode     string   `json:"union_code"`
	Point         string   `json:"point"`
	Grade         string   `json:"grade"`
	StoreCode     string   `json:"store_code"`
	BACode        string   `json:"ba_code"`
	Mobile        string   `json:"mobile"`
}

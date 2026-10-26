package work

// DelCorpTagParam 删除标签参数
type DelCorpTagParam struct {
	TagID   []string
	GroupID []string
}

// 新增企业标签标签请求参数
type AddCorpTagRequest struct {
	QyGwRequest
	GroupID   string       `json:"group_id,omitempty"`
	GroupName string       `json:"group_name,omitempty"`
	Order     int          `json:"order,omitempty"`
	TagList   []AddCorpTag `json:"tag"`
	AgentID   int          `json:"agentid,omitempty"`
}

type AddCorpTag struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// 返回企业信息集合
type AddCorpTagResponse struct {
	QyResponse
	TagGroup TagGroupResponse `json:"tag_group"`
}

type TagGroupResponse struct {
	GroupID    string        `json:"group_id"`
	GroupName  string        `json:"group_name"`
	CreateTime int64         `json:"create_time"`
	Order      int           `json:"order"`
	TagList    []TagResponse `json:"tag"`
}

type TagResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
	Order      int    `json:"order"`
}

// DelCorpTagRequest 删除企业标签请求参数
type DelCorpTagRequest struct {
	QyGwRequest
	TagID   []string `json:"tag_id"`
	GroupID []string `json:"group_id"`
	AgentID int      `json:"agentid"`
}

// 更新标签请求参数
type UpdateTagRequest struct {
	QyGwRequest
	ID      string `json:"id"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	AgentID int64  `json:"agentid"`
}

type TaglistResponse struct {
	QyResponse
	Taglist []struct {
		Tagid   int    `json:"tagid"`
		Tagname string `json:"tagname"`
	} `json:"taglist"`
}

type TagCreateRequest struct {
	QyGwRequest
	Tagname string `json:"tagname"`
	Tagid   int    `json:"tagid"`
}

type TagCreateResponse struct {
	QyResponse
	Tagid int `json:"tagid"`
}

type TagUsersResponse struct {
	QyResponse
	Tagname   string                     `json:"tagname"`
	Userlist  []TagUsersUserlistResponse `json:"userlist"`
	Partylist []int                      `json:"partylist"`
}
type TagUsersUserlistResponse struct {
	Userid string `json:"userid"`
	Name   string `json:"name"`
}

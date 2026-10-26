package drainage

// ContactChannels 引流模块联系我渠道表
type ContactChannels struct {
	ID           int64                  `json:"id" db:"id" desc:"删除时间"`
	EnterpriseID int64                  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Name         string                 `json:"name" db:"name" desc:"渠道名称"`
	CreatedAt    int64                  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64                  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64                  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	Labels       []ContactChannelLabels `json:"labels"`
}

// ContactChannelLabels 引流模块联系我渠道标签表
type ContactChannelLabels struct {
	ID           int64  `json:"id" db:"id" desc:"删除时间"`
	ChannelID    int64  `json:"channel_id" db:"channel_id" desc:"渠道id"`
	LabelGroupID int64  `json:"label_group_id" db:"label_group_id" desc:"标签分组id"`
	LabelID      int64  `json:"label_id" db:"label_id" desc:"标签id"`
	WxTagGroupID string `json:"wx_tag_group_id" db:"wx_tag_group_id"`
	WxTagID      string `json:"wx_tag_id" db:"wx_tag_id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	LabelName    string `json:"label_name" db:"label_name" desc:"标签名称"`
}

// CreateContactChannelLabelsRequest 创建渠道请求参数
type CreateContactChannelLabelsRequest struct {
	ID           int64                  `json:"id"`
	EnterpriseID int64                  `json:"enterprise_id"`
	Name         string                 `json:"name"`
	Label        []ContactChannelLabels `json:"label"`
}

// ContactChannelListRequest 渠道列表请求
type ContactChannelListRequest struct {
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id"`
	ChannelName  string   `json:"channel_name"`
	ChannelNames []string `json:"-" desc:"渠道名字符串，支持按渠道名列表查询"`
	ChannelID    int64    `json:"channel_id"`
	Page         int      `json:"page"`
	Size         int      `json:"size"`
}

// ContactChannelDetailResponse 渠道详情返回结果
type ContactChannelDetailResponse struct {
	ChannelID      int64  `json:"channel_id" db:"channel_id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`
	ChannelName    string `json:"channel_name" db:"channel_name"`
	ChannelCreated int64  `json:"channel_created" db:"channel_created"`
	LabelGroupID   int64  `json:"label_group_id" db:"label_group_id"`
	LabelGroupName string `json:"label_group_name"`
	LabelID        int64  `json:"label_id" db:"label_group_id"`
	LabelName      string `json:"label_name"`
}

// ContactChannelLabelsListRequest 渠道标签列表
type ContactChannelLabelsListRequest struct {
	Page            int     `json:"page"`
	Size            int     `json:"size"`
	EnterpriseID    int64   `json:"enterprise_id" db:"enterprise_id"`
	ChannelID       int64   `json:"channel_id" db:"channel_id"`
	ChannelIds      []int64 `json:"channel_ids"`
	ChannelIdsCount int
}

// ContactChannels 引流模块联系我渠道表
type ContactChannelsList struct {
	ID            int64    `json:"id" db:"id" desc:"删除时间"`
	EnterpriseID  int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Name          string   `json:"name" db:"name" desc:"渠道名称"`
	CreatedAt     int64    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt     int64    `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt     int64    `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	ChannelLabels []string `json:"channel_labels"`
}

// WelcomeChannelsList
type WelcomeChannelsList struct {
	ID           int64  `json:"id" db:"id" desc:"删除时间"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Channel      string `json:"channel"  db:"channel" desc:"渠道名称"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// SimpleChannelList 渠道列表
type SimpleChannelList struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

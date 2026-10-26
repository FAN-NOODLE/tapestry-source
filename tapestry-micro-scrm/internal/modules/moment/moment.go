package moment

// MomentTopic MomentTopic表
type MomentTopic struct {
	ID              int64  `json:"id" db:"id" desc:"id"`
	Name            string `json:"name" db:"name" desc:"主题名称"`
	PublishType     int8   `json:"publish_type" db:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	PublishTime     int64  `json:"publish_time" db:"publish_time" desc:"定时发布时间"`
	Content         string `json:"content" db:"content" desc:"内容"`
	CreatedMemberId int64  `json:"created_member_id" db:"created_member_id" desc:"朋友圈创建员工ID"`
	MemberCount     int64  `json:"member_count" db:"member_count" desc:"员工数量"`
	EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt       int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt       int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MomentMember MomentMember表
type MomentMember struct {
	ID           int64  `json:"id" db:"id" binding:"required" desc:"id"`
	MemberId     int64  `json:"member_id" db:"member_id" binding:"required" desc:"员工id"`
	Name         string `json:"name" db:"name" binding:"required,max=15" desc:"名称"`
	Sign         string `json:"sign" db:"sign" binding:"required,max=30" desc:"签名"`
	Mark         string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	Avatar       string `json:"avatar" db:"avatar" binding:"required" desc:"签名"`
	Background   string `json:"background" db:"background" binding:"required" desc:"签名"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" binding:"required" desc:"企业ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MomentTopicMember MomentTopicMember表
type MomentTopicMember struct {
	ID             int64 `json:"id" db:"id" desc:"id"`
	MomentTopicId  int64 `json:"moment_topic_id" db:"moment_topic_id" desc:"circle_topics表ID"`
	MomentMemberId int64 `json:"moment_member_id" db:"moment_member_id" desc:"circle_members表ID"`
	EnterpriseId   int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt      int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt      int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt      int64 `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MomentTopicContent MomentTopicContent表的content
type MomentTopicContent struct {
	Pics    []string `json:"pics" desc:"图片"`
	Links   []string `json:"links" desc:"链接"`
	Videos  []string `json:"videos" desc:"视频"`
	Applets []string `json:"applets" desc:"小程序"`
	Text    string   `json:"text" desc:"内容"`
}

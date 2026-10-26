package moment

// MomentMemberGetRequest 获取朋友圈形象请求
type MomentMemberGetRequest struct {
	Mark         string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	UnionId      string `json:"union_id" db:"union_id"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// MomentMemberGetResponse 获取朋友圈形象响应
type MomentMemberGetResponse struct {
	ID           int64  `json:"id" db:"id" binding:"required" desc:"id"`
	MemberId     int64  `json:"member_id" db:"member_id" binding:"required" desc:"员工id"`
	Name         string `json:"name" db:"name" binding:"required,max=15" desc:"名称"`
	Sign         string `json:"sign" db:"sign" binding:"required,max=30" desc:"签名"`
	Mark         string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	Link         string `json:"link" db:"link" desc:"link"`
	Avatar       string `json:"avatar" db:"avatar" binding:"required" desc:"签名"`
	Background   string `json:"background" db:"background" binding:"required" desc:"签名"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" binding:"required" desc:"企业ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MomentMemberEditRequest 编辑朋友圈形象请求
type MomentMemberEditRequest struct {
	Name         string `json:"name" binding:"required" desc:"名称"`
	Sign         string `json:"sign" desc:"签名"`
	Mark         string `json:"mark" binding:"required" desc:"链接加密mark"`
	Avatar       string `json:"avatar" binding:"required" desc:"签名"`
	Background   string `json:"background" binding:"required" desc:"签名"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" binding:"required" desc:"企业ID"`
}

// MomentMemberEditResponse 编辑朋友圈形象响应
type MomentMemberEditResponse struct {
	AffectedNum int64 `json:"affected_num" desc:"受影响的行数"`
}

// TopicCreateRequest 朋友圈手动发布请求
type TopicCreateRequest struct {
	Name            string                 `json:"name" validate:"required,max=30" required:"主题名称不能为空" max:"最大长度30" desc:"主题名称"`
	PublishType     int8                   `json:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	PublishTime     int64                  `json:"publish_time" desc:"定时发布时间"`
	EnterpriseId    int64                  `json:"enterprise_id" desc:"企业ID"`
	CreatedMemberId int64                  `json:"created_member_id" desc:"朋友圈创建员工ID"`
	EditMemberId    int64                  `json:"edit_member_id" desc:"朋友圈创建员工ID"`
	MemberIds       []int64                `json:"member_ids" desc:"员工IDS"`
	DepartmentIds   []int64                `json:"department_ids" desc:"部门IDS"`
	Content         map[string]interface{} `json:"content" desc:"内容"`
}

// TopicCreateResponse 朋友圈手动发布返回
type TopicCreateResponse struct {
	MomentTopicId int64 `json:"moment_topic_id" desc:"最新发布朋友圈的ID"`
}

// TopicEditRequest 朋友圈编辑发布请求
type TopicEditRequest struct {
	ID int64 `json:"id" binding:"required" desc:"id"`
	TopicCreateRequest
}

// TopicEditResponse 朋友圈手动发布返回
type TopicEditResponse struct {
	MomentTopicId int64 `json:"moment_topic_id" desc:"最新发布朋友圈的ID"`
}

// TopicDisplayRequest 朋友圈展示请求
type TopicDisplayRequest struct {
	ID           int64 `json:"id" binding:"required" db:"id" desc:"id"`
	Page         int64 `json:"page" binding:"required" db:"page" desc:"page"`
	Size         int64 `json:"size" binding:"required" db:"size" desc:"size"`
	EditMemberId int64 `json:"edit_member_id" db:"moment_member_id" desc:"朋友圈创建员工ID"`
}

// TopicDisplayResponse 朋友圈展示返回
type TopicDisplayResponse struct {
	ID              int64                  `json:"id" db:"id" desc:"id"`
	Name            string                 `json:"name" db:"name" desc:"主题名称"`
	PublishType     int8                   `json:"publish_type" db:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	PublishTime     int64                  `json:"publish_time" db:"publish_time" desc:"定时发布时间"`
	Content         map[string]interface{} `json:"content" db:"content" desc:"内容"`
	CreatedMemberId int64                  `json:"created_member_id" db:"created_member_id" desc:"朋友圈创建员工ID"`
	MemberCount     int64                  `json:"member_count" db:"member_count" desc:"员工数量"`
	EnterpriseId    int64                  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt       int64                  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt       int64                  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt       int64                  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	Page            int64                  `json:"page" desc:"page"`
	Size            int64                  `json:"size" desc:"size"`
	Total           int64                  `json:"total" desc:"total"`
	List            []TopicDisplayExt      `json:"list" desc:"list"`
}

// TopicDisplayExt 朋友圈展示子结构体
type TopicDisplayExt struct {
	MomentTopicMemberId int64  `json:"moment_topic_member_id" db:"moment_topic_member_id"`
	MomentMemberId      int64  `json:"moment_member_id" db:"moment_member_id"`
	MemberId            int64  `json:"member_id" db:"member_id"`
	Name                string `json:"user_name" db:"user_name"`
	Userid              string `json:"userid" db:"userid"`
}

// TopicDeleteRequest 删除topic请求
type TopicDeleteRequest struct {
	ID int64 `json:"id" binding:"required" db:"id" desc:"id"`
}

// TopicDeleteResponse 删除topic响应
type TopicDeleteResponse struct {
	Result bool `json:"result" desc:"是否成功删除"`
}

// TopicMemberBatchDelete 批量删除
type TopicMemberBatchDelete struct {
	DeletedAt     int64 `json:"deleted_at" db:"deleted_at"`
	MomentTopicId int64 `json:"moment_topic_id" db:"moment_topic_id"`
}

// ContentListRequest 朋友圈按内容列表展示请求
type ContentListRequest struct {
	Name            string `json:"name" validate:"required,max=30" required:"主题名称不能为空" max:"最大长度30" desc:"主题名称"`
	EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedMemberId int64  `json:"created_member_id" db:"created_member_id" desc:"朋友圈创建员工ID"`
	CurrentTime     int64  `json:"current_time" db:"current_time" desc:"当前时间"`
	StartTime       int64  `json:"start_time" db:"start_time" desc:"开始时间"`
	EndTime         int64  `json:"end_time" db:"end_time" desc:"结束时间"`
	Status          int8   `json:"status" db:"status" desc:"1已发布2定时未发布"`
	Page            int64  `json:"page" desc:"page"`
	Size            int64  `json:"size" desc:"size"`
}

// ContentListResponse 朋友圈按内容列表展示响应
type ContentListResponse struct {
	Page      int64          `json:"page" desc:"page"`
	Size      int64          `json:"size" desc:"size"`
	Total     int64          `json:"total" desc:"total"`
	UpdatedAt int64          `json:"updated_at" db:"updated_at" desc:"更新时间"`
	List      []TopicContent `json:"list" desc:"list"`
}

// TopicContent 朋友圈按内容列表展示响应子结构体
type TopicContent struct {
	ID              int64                  `json:"id" db:"id" desc:"id"`
	Name            string                 `json:"name" db:"name" desc:"主题名称"`
	Content         string                 `json:"content" db:"content" desc:"内容"`
	ContentJson     map[string]interface{} `json:"content_json" db:"content_json" desc:"内容"`
	PublishType     int8                   `json:"publish_type" db:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	PublishTime     int64                  `json:"publish_time" db:"publish_time" desc:"定时发布时间"`
	CreatedMemberId int64                  `json:"created_member_id" db:"created_member_id" desc:"朋友圈创建员工ID"`
	MemberCount     int64                  `json:"member_count" db:"member_count" desc:"员工数量"`
	CreatedAt       int64                  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt       int64                  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	Status          int8                   `json:"status" db:"status"`
}

// MemberListRequest 朋友圈列表按员工搜索请求
type MemberListRequest struct {
	MemberId     int64         `json:"member_id" db:"member_id" desc:"员工id"`
	MemberIds    []interface{} `json:"member_ids" db:"member_ids" desc:"员工id"`
	EnterpriseId int64         `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DepartmentId int64         `json:"department_id" desc:"部门ID"`
	StartTime    int64         `json:"start_time" db:"start_time" desc:"开始时间"`
	EndTime      int64         `json:"end_time" db:"end_time" desc:"结束时间"`
	Page         int64         `json:"page" desc:"page"`
	Size         int64         `json:"size" desc:"size"`
}

// MemberListResponse 朋友圈列表按员工搜索响应
type MemberListResponse struct {
	Page      int64         `json:"page" desc:"page"`
	Size      int64         `json:"size" desc:"size"`
	Total     int64         `json:"total" desc:"total"`
	UpdatedAt int64         `json:"updated_at" desc:"更新时间"`
	List      []TopicMember `json:"list" desc:"list"`
}

// MomentTopicInsertRequest 发布朋友圈请求
type MomentTopicInsertRequest struct {
	Name            string             `json:"name" validate:"required,max=30" required:"主题名称不能为空" max:"最大长度30" desc:"主题名称"`
	PublishType     int8               `json:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	MemberIDs       []int64            `json:"member_ids" desc:"员工IDS"`
	EnterpriseId    int64              `json:"enterprise_id" desc:"企业ID"`
	DepartmentIds   []int64            `json:"department_ids" desc:"部门IDS"`
	CreatedMemberId int64              `json:"created_member_id" desc:"朋友圈创建员工ID"`
	Content         MomentTopicContent `json:"content" desc:"内容"`
}

// TopicMember 朋友圈列表按员工搜索子结构体
type TopicMember struct {
	MemberId       int64  `json:"member_id" db:"member_id" desc:"员工id"`
	Mark           string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	UserName       string `json:"user_name" db:"user_name" desc:"员工姓名"`
	DepartmentId   int64  `json:"department_id" db:"department_id" desc:"部门IDS"`
	DepartmentName string `json:"department_name" db:"department_name" desc:"部门名称"`
	DepartmentStr  string `json:"department_str" db:"department_str" desc:"部门名称"`
	TopicCount     int64  `json:"topic_count" db:"topic_count" desc:"朋友圈数"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// MemberTopicRequest 一个人的朋友圈列表请求
type MemberTopicRequest struct {
	Page         int64  `json:"page" db:"page" desc:"page"`
	Size         int64  `json:"size" db:"size" desc:"size"`
	PublishTime  int64  `json:"publish_time" db:"publish_time" desc:"定时发布时间"`
	Mark         string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// MemberTopicResponse 一个人的朋友圈列表返回
type MemberTopicResponse struct {
	Page  int64            `json:"page" db:"page" desc:"page"`
	Size  int64            `json:"size" db:"size" desc:"size"`
	Total int64            `json:"total" db:"total" desc:"total"`
	List  []MemberTopicExt `json:"list" desc:"list"`
}

// MemberTopicExt 一个人的朋友圈列表中List单元
type MemberTopicExt struct {
	ID                  int64                  `json:"moment_topic_id" db:"moment_topic_id" desc:"moment_topic_id"`
	Name                string                 `json:"name" db:"name" desc:"主题名称"`
	MemberId            int64                  `json:"member_id" db:"member_id" desc:"员工id"`
	MomentTopicMemberId int64                  `json:"moment_topic_member_id" db:"moment_topic_member_id" desc:"员工id"`
	PublishType         int8                   `json:"publish_type" db:"publish_type" desc:"发布类型 1立即发送 2定时发送"`
	PublishTime         int64                  `json:"publish_time" db:"publish_time" desc:"定时发布时间"`
	Content             string                 `json:"content" db:"content" desc:"内容"`
	ContentJson         map[string]interface{} `json:"content_json" db:"content_json" desc:"内容"`
	CreatedMemberId     int64                  `json:"created_member_id" db:"created_member_id" desc:"朋友圈创建员工ID"`
	MemberCount         int64                  `json:"member_count" db:"member_count" desc:"员工数量"`
	EnterpriseId        int64                  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt           int64                  `json:"created_at" db:"created_at" desc:"创建时间"`
	DeletedAt           int64                  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
	UserName            string                 `json:"user_name" db:"user_name" desc:"企业微信姓名"`
	Avatar              string                 `json:"user_avatar" db:"user_avatar" desc:"员工头像"`
	MomentUserName      string                 `json:"moment_user_name" db:"moment_user_name" desc:"朋友圈微信姓名"`
	MomentAvatar        string                 `json:"moment_user_avatar" db:"moment_user_avatar" desc:"朋友圈员工头像"`
}

// MemberBatchGetByMemberIds 批量获取朋友圈用户信息
type MemberBatchGetByMemberIds struct {
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DeletedAt    int64 `db:"deleted_at"`
	MemberIds    []interface{}
}

// TopicMemberDeleteRequest 删除员工朋友圈请求
type TopicMemberDeleteRequest struct {
	DeletedAt           int64 `json:"deleted_at" db:"deleted_at"`
	MemberId            int64 `json:"member_id" db:"member_id" desc:"企业ID"`
	EnterpriseId        int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MomentTopicId       int64 `json:"moment_topic_id" db:"moment_topic_id"`
	MomentTopicMemberId int64 `json:"moment_topic_member_id" db:"moment_topic_member_id"`
}

// TopicMemberDeleteResponse 删除员工朋友圈响应
type TopicMemberDeleteResponse struct {
	Result bool `json:"result"`
}

// TopicMemberEditRequest 编辑员工朋友圈请求
type TopicMemberEditRequest struct {
	TopicCreateRequest
	MomentTopicId       int64 `json:"moment_topic_id" db:"moment_topic_id"`
	MomentTopicMemberId int64 `json:"moment_topic_member_id" db:"moment_topic_member_id"`
}

// TopicMemberEditResponse 编辑员工朋友圈响应
type TopicMemberEditResponse struct {
	MomentTopicId int64 `json:"moment_topic_id" desc:"最新发布朋友圈的ID"`
}

// BatchUpdateMemberTime 批量修改moment_member更新时间
type BatchUpdateMemberTime struct {
	DeletedAt    int64 `db:"deleted_at"`
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"企业ID"`
	MemberIDs    []interface{}
}

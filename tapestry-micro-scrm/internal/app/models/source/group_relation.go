package source

type GroupRelationModel struct {
	ID            int64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned"`
	EnterpriseID  int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;default:0;comment:企业ID"`
	DepartmentID  int64 `json:"department_id" gorm:"column:department_id;type:bigint(20) unsigned;not null;index:department_id_idx;default:0;comment:关联部门/门店表id"`
	ThirdId       int   `json:"third_id" gorm:"column:third_id;type:int(10) unsigned;not null;index:third_id_idx;default:0;comment:关联企微部门id"`
	SourceGroupID int64 `json:"source_group_id" gorm:"column:source_group_id;type:bigint(20) unsigned;not null;index:source_group_id_idx;default:0;comment:关联素材分组id"`
	CreatedAt     int64 `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt     int64 `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:修改时间"`
	DeletedAt     int64 `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;default:0;comment:删除时间"`
}

type GroupRelationData struct {
	ID            int64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned"`
	EnterpriseID  int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;default:0;comment:企业ID"`
	DepartmentID  int64 `json:"department_id"`
	ThirdId       int   `json:"third_id"`
	SourceGroupID int64 `json:"source_group_id" gorm:"column:source_group_id;type:bigint(20) unsigned;not null;index:source_group_id_idx;default:0;comment:关联素材分组id"`
}

type GroupRelationList struct {
	GroupRelationData
	DepartmentName string `json:"department_name" gorm:"column:department_name;type:varchar(128);not null" `
}

type GroupRelationTreeList struct {
	GroupRelationData
	GroupName string                   `json:"group_name" `
	ParentId  int64                    `gorm:"column:parent_id;NOT NULL" json:"parent_id"` // 父级id
	Orders    int                      `gorm:"column:orders;NOT NULL" json:"orders"`       // 排序字段
	Children  []*GroupRelationTreeList `json:"children"`
}

type GroupRelationListAll struct {
	GroupRelationData
	GroupName string `json:"group_name" `
	ParentId  int64  `gorm:"column:parent_id;NOT NULL" json:"parent_id"` // 父级id
	Orders    int    `gorm:"column:orders;NOT NULL" json:"orders"`       // 排序字段
}

// ResponseSelectList 公有返回结构
type ResponseSelectList struct {
	ID            int64  `json:"id" db:"id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID       int64  `json:"group_id" db:"group_id"`
	Type          int    `json:"type" db:"type"`
	CreatedAt     int    `json:"created_at" db:"created_at"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	CategoryIds   string `json:"category_ids" db:"category_ids"`
	CategoryNames string `json:"category_names" db:"category_names"`
	SourceStatus  int    `json:"source_status" db:"source_status"`
	SourceType    int    `json:"source_type" db:"source_type"`
	IsNew         int    `json:"is_new"`
	ThirdId       int    `json:"third_id"`
	IsDeal        int32  `json:"is_deal"`
	AdminId       int64  `json:"admin_id"`
	MemberId      int64  `json:"member_id"`
	MemberName    string `json:"member_name"`
	MemberAvatar  string `json:"member_avatar"`
	IsCollect     int    `json:"is_collect"`
	CollectNum    int    `json:"collect_num"`
	IsTop         int    `json:"is_top"`
	CollectAt     int64  `json:"collect_at"`
}

type ResponseTextList struct {
	ResponseSelectList

	Content       string `json:"content" db:"content" desc:"文本内容"`
	Status        int    `json:"status" db:"status"`
	DepartmentId  int64  `json:"department_id"`
	SourceGroupId int64  `json:"source_group_id"`
	CopyTimes     int    `json:"copy_times"`
}

// ResponseImgList 图片类型返回值
type ResponseImgList struct {
	ResponseSelectList

	Name            string `json:"name" db:"name"`
	URL             string `json:"url" db:"url"`
	Size            int64  `json:"size"  db:"size"`
	MediaID         string `json:"media_id" db:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64  `json:"media_status" db:"media_status"`
	DepartmentId    int64  `json:"department_id"`
	SourceGroupID   int64  `json:"source_group_id" db:"source_group_id"`
	MergeId         int64  `json:"merge_id"`
	Sort            int32  `json:"sort"`
	DownloadTimes   int    `json:"download_times"`
	ForwardingTimes int    `json:"forwarding_times"`
}

// ResponseLinkList 链接素材列表返回结构
type ResponseLinkList struct {
	ResponseSelectList

	GroupName     string `json:"group_name" db:"group_name"`
	Title         string `json:"title" db:"title" desc:"标题"`
	Description   string `json:"description" db:"description" desc:"描述"`
	CoverUrl      string `json:"cover_url" db:"cover_url" desc:"描述"`
	LinkType      int    `json:"link_type" db:"link_type" desc:"链接类型"`
	LinkUrl       string `json:"link_url" db:"link_url" desc:"链接URL"`
	Author        string `json:"author" db:"author" desc:"作者"`
	Content       string `json:"content" db:"content" desc:"文本内容"`
	DepartmentId  int64  `json:"department_id"`
	SourceGroupID int64  `json:"source_group_id" db:"source_group_id"`
}

// ResponseMiniProgramList 小程序素材
type ResponseMiniProgramList struct {
	ResponseSelectList

	Title         string `json:"title" db:"title"`
	MiniProgramID int64  `json:"mini_program_id" db:"mini_program_id"`
	AppName       string `json:"app_name"  db:"app_name"`
	AppID         string `json:"app_id"  db:"app_id"`
	OriginID      string `json:"origin_id" db:"origin_id"`
	Logo          string `json:"logo" db:"logo"`
	Path          string `json:"path" db:"path"`
	CoverUrl      string `json:"cover_url" db:"cover_url"`
	DepartmentId  int64  `json:"department_id"`
	SourceGroupID int64  `json:"source_group_id" db:"source_group_id"`
}

// ResponseVideoList 视频素材
type ResponseVideoList struct {
	ResponseSelectList

	Name            string `json:"name"  db:"name"`
	CoverUrl        string `json:"cover_url" db:"cover_url"`
	Size            int64  `json:"size"  db:"size"`
	Ext             string `json:"ext"  db:"ext"`
	Url             string `json:"url"  db:"url"`
	Status          int    `json:"status"  db:"status"`
	MediaID         string `json:"media_id" db:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64  `json:"media_status" db:"media_status"`
	DepartmentId    int64  `json:"department_id"`
	SourceGroupID   int64  `json:"source_group_id" db:"source_group_id"`
	DownloadTimes   int    `json:"download_times"`
	ForwardingTimes int    `json:"forwarding_times"`
}

// ResponseFileList 文件类型List返回值
type ResponseFileList struct {
	ResponseSelectList

	Name            string `json:"name" db:"name"`
	URL             string `json:"url" db:"url"`
	Ext             string `json:"ext" db:"ext"`
	FileType        int    `json:"file_type" db:"file_type"`
	MediaID         string `json:"media_id" db:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64  `json:"media_status" db:"media_status"`
	Size            int    `json:"size" db:"size"`
	DepartmentId    int64  `json:"department_id"`
	SourceGroupID   int64  `json:"source_group_id" db:"source_group_id"`
}

type ResponseFriendList struct {
	ResponseSelectList

	TextContent     string `json:"text_content" comment:"文本内容"`
	Name            string `json:"name"  comment:"资源名称"`
	VideoCoverURL   string `json:"video_cover_url" comment:"视频封面地址"`
	OriginalURL     string `json:"original_url" comment:"原资源地址"`
	CompressionURL  string `json:"compression_url" comment:"压缩资源地址"`
	Ext             string `json:"ext" comment:"文件后缀"`
	Size            int32  `json:"size" comment:"文件大小"`
	Type            int32  `json:"type" comment:"类型 1:图片 2:视频"`
	Sort            int32  `json:"sort" comment:"排序"`
	MergeId         int64  `json:"merge_id"`
	DownloadTimes   int    `json:"download_times" comment:"下载次数"`
	ForwardingTimes int    `json:"forwarding_times" comment:"转发次数"`
	CopyTimes       int    `json:"copy_times" comment:"复制次数"`
	CollectNum      int    `json:"collect_num" comment:"收藏数"`
}

// TableName returns the table name of the SourceGroupRelation model
func (GroupRelationModel) TableName() string {
	return "source_group_relation"
}

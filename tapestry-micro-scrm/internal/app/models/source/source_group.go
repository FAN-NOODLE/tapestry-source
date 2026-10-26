package source

type SourceGroupModel struct {
	ID           int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	EnterpriseID int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;default:0;comment:企业id"`
	GroupName    string `json:"group_name" gorm:"column:group_name;type:varchar(255);not null"`
	IgnoreIndex  int    `json:"ignore_index" gorm:"column:ignore_index;type:tinyint(1);not null"`
	CreatedAt    int    `json:"created_at" gorm:"column:created_at;type:int(11);not null;default:0;comment:创建时间"`
	UpdatedAt    int    `json:"updated_at" gorm:"column:updated_at;type:int(11);not null;default:0;comment:更新时间"`
	DeletedAt    int    `json:"deleted_at" gorm:"column:deleted_at;type:int(11);not null;default:0;comment:删除时间"`
	ParentID     int64  `json:"parent_id" gorm:"column:parent_id;type:bigint(20);not null;comment:父级ID"`
	Orders       int    `json:"orders" gorm:"column:orders;type:int(11);not null;comment:分组排序"`
}

type SourceGroupTreeList struct {
	SourceGroupModel
	IsShow   int                    `json:"is_show"` // 是否展示 0 不展示 1展示
	Children []*SourceGroupTreeList `json:"children"`
}

// ResponseConsoleList 公有返回结构
type ResponseConsoleList struct {
	ID            int64    `json:"id"`
	EnterpriseID  int64    `json:"enterprise_id"`
	GroupID       int64    `json:"group_id"`
	Type          int      `json:"type"`
	AdminID       int64    `json:"admin_id"`
	AdminName     string   `json:"admin_name"`
	CreatedAt     int      `json:"created_at"`
	StartTime     int64    `json:"start_time"`
	EndTime       int64    `json:"end_time"`
	CategoryIds   string   `json:"category_ids"`
	CategoryNames string   `json:"category_names"`
	SourceStatus  int      `json:"source_status"`
	SourceType    int      `json:"source_type"`
	IsNew         int      `json:"is_new"`
	IsTop         int      `json:"is_top"`
	CanTop        bool     `json:"can_top"`
	SkuList       []string `json:"sku_list"`
}

// ResponseConsoleTextList 文本素材列表返回结构
type ResponseConsoleTextList struct {
	ResponseConsoleList

	Content string `json:"content" comment:"文本内容"`
}

// ResponseConsoleImageList 图片类型返回值
type ResponseConsoleImageList struct {
	ResponseConsoleList

	Name            string `json:"name"`
	URL             string `json:"url"`
	MediaID         string `json:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time"`
	MediaStatus     int64  `json:"media_status"`
	MergeId         int64  `json:"merge_id"`
	Sort            int32  `json:"sort"`
}

// ResponseConsoleLinkList 链接素材列表返回结构
type ResponseConsoleLinkList struct {
	ResponseConsoleList

	GroupName   string `json:"group_name"`
	Title       string `json:"title" comment:"标题"`
	Description string `json:"description" comment:"描述"`
	CoverUrl    string `json:"cover_url" comment:"描述"`
	LinkType    int    `json:"link_type" comment:"链接类型"`
	LinkUrl     string `json:"link_url" comment:"链接URL"`
	Author      string `json:"author" comment:"作者"`
	Content     string `json:"content" comment:"文本内容"`
	IsKOSSync   int    `json:"is_kos_sync"`
}

// ResponseConsoleMiniProgramList 小程序素材
type ResponseConsoleMiniProgramList struct {
	ResponseConsoleList

	Title         string `json:"title" `
	MiniProgramID int64  `json:"mini_program_id"`
	AppName       string `json:"app_name"`
	AppID         string `json:"app_id"`
	OriginID      string `json:"origin_id"`
	Logo          string `json:"logo"`
	Path          string `json:"path"`
	CoverUrl      string `json:"cover_url"`
}

// ResponseConsoleVideoList 视频素材
type ResponseConsoleVideoList struct {
	ResponseConsoleList

	Name            string `json:"name"`
	CoverUrl        string `json:"cover_url"`
	Size            int64  `json:"size"`
	Ext             string `json:"ext"`
	Url             string `json:"url"`
	Status          int    `json:"status"`
	MediaID         string `json:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time"`
	MediaStatus     int64  `json:"media_status"`
}

// ResponseConsoleFileList 文件类型List返回值
type ResponseConsoleFileList struct {
	ResponseConsoleList

	Name            string `json:"name"`
	URL             string `json:"url"`
	Ext             string `json:"ext"`
	FileType        int    `json:"file_type"`
	MediaID         string `json:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time"`
	MediaStatus     int64  `json:"media_status"`
	Size            int    `json:"size"`
}

// ResponseConsoleFriendList 朋友圈素材
type ResponseConsoleFriendList struct {
	ResponseConsoleList

	TextContent    string `json:"text_content" comment:"文本内容"`
	Name           string `json:"name"  comment:"资源名称"`
	VideoCoverURL  string `json:"video_cover_url" comment:"视频封面地址"`
	OriginalURL    string `json:"original_url" comment:"原资源地址"`
	CompressionURL string `json:"compression_url" comment:"压缩资源地址"`
	Ext            string `json:"ext" comment:"文件后缀"`
	Size           int32  `json:"size" comment:"文件大小"`
	Type           int32  `json:"type" comment:"类型 1:图片 2:视频"`
	Sort           int32  `json:"sort" comment:"排序"`
	MergeId        int64  `json:"merge_id"`
}

// TableName returns the table name of the SourceGroup model
func (s SourceGroupModel) TableName() string {
	return "source_group"
}

const IgnoreIndex = 1

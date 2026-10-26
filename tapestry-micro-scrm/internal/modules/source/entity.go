package source

import "micro-scrm/internal/app/http/request/source"

type ID struct {
	ID int64 `json:"id"`
}

// Source 素材主表
type Source struct {
	ID int64 `json:"id" db:"id"`
	// EnterpriseID 企业id
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	// GroupID 分组ID 0 未分组
	GroupID int64 `json:"group_id" db:"group_id"`
	// Type 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
	Type int `json:"type" db:"type"`
	// AdminID 管理员ID
	AdminID int64 `json:"admin_id" db:"admin_id"`
	// MemberID 对应企业微信员工ID
	MemberID int64 `json:"member_id" db:"member_id"`
	// 开始时间
	StartTime int64 `json:"start_time" db:"start_time"`
	// 结束时间
	EndTime int64 `json:"end_time" db:"end_time"`
	// CreatedAt 创建时间
	CreatedAt int64 `json:"created_at" db:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt int64 `json:"updated_at" db:"updated_at"`
	// DeletedAt 删除时间
	DeletedAt int64 `json:"deleted_at" db:"deleted_at"`
	// 合并组id
	MergeID int64 `json:"merge_id" db:"merge_id"`
	// 排序
	Sort int32 `json:"sort" db:"sort"`
}

// SourcesImage 图片素材
type SourcesImage struct {
	ID int64 `json:"id"`
	// EnterpriseID 企业id
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	// SourceID 素材主表id
	SourceID int64 `json:"source_id" db:"source_id"`
	// SourceName 图片名称
	Name string `json:"name" db:"name"`
	// SourceURL 图片url
	URL string `json:"url" db:"url"`
	// Ext 文件后缀
	Ext string `json:"ext" db:"ext"`
	// 文件大小
	Size int64 `json:"size" db:"size"`
	// CreatedAt 创建时间
	CreatedAt int64 `json:"created_at" db:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt int64 `json:"updated_at" db:"updated_at"`
	// DeletedAt 删除时间
	DeletedAt int64 `json:"deleted_at" db:"deleted_at"`
}

// SourcesFile 文件素材
type SourcesFile struct {
	ID int64 `json:"id"`
	// EnterpriseID 企业id
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	// SourceName 图片名称
	Name string `json:"name" db:"name"`
	// SourceURL 图片url
	URL string `json:"url" db:"url"`
	// AttachName 第三方文件名称
	AttachName string `json:"attach_name" db:"attach_name"`
	// AttachID 第三方file ID
	AttachID int64 `json:"attach_id" db:"attach_id"`
	// FileType 文件类型 1 word 2 excel 3 ppt 4 pdf 6 其他
	FileType int `json:"file_type" db:"file_type"`
	// FileExt 文件后缀
	Ext string `json:"ext" db:"ext"`
	// Size 文件大小
	Size int64 `json:"size" db:"size"`
	// CreatedAt 创建时间
	CreatedAt int64 `json:"created_at" db:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt int64 `json:"updated_at" db:"updated_at"`
	// DeletedAt 删除时间
	DeletedAt int64 `json:"deleted_at" db:"deleted_at"`
}

// RequestList 请求
type RequestList struct {
	Page            int     `json:"page" db:"page" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
	Size            int     `json:"size" db:"size" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
	EnterpriseId    int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID         int64   `json:"group_id" db:"group_id"`
	JustMe          int     `json:"just_me" db:"just_me"`
	AdminID         int64   `json:"admin_id" db:"admin_id"`
	Keyword         string  `json:"keyword" db:"keyword"`
	MediaExpireTime int64   `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64   `json:"media_status" db:"media_status"`
	CreateStart     int64   `json:"create_start" db:"create_start"`
	CreateEnd       int64   `json:"create_end" db:"create_end"`
	CategoryIds     []int64 `json:"category_ids"`
	GroupIDList     []int64 `json:"group_id_list" db:"group_id_list"`
	Order           string  `json:"order"`
	Sort            string  `json:"sort"`
	SourceStatus    int     `json:"source_status"`
	IsUse           int     `json:"is_use"`
	NowTime         int64   `json:"now_time"`
}

// SourcesImageRequest 图片素材请求
type SourcesImageRequest struct {
	EnterpriseID int64                     `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID      int64                     `json:"group_id" db:"group_id"`
	AdminID      int64                     `json:"admin_id" db:"admin_id"`
	MemberID     int64                     `json:"member_id" db:"member_id"`
	Images       []*SourcesImageRequestAdd `json:"images" validate:"required" required:"图片必传"`
	CategoryIds  []int64                   `json:"category_ids" required:"产品类别"`
	StartTime    int64                     `json:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64                     `json:"end_time" validate:"required" required:"结束时间必传"`
	ShowType     int                       `json:"show_type" comment:"展示类型:1单张,2合并"`
}

// SourcesImageRequestAdd 图片素材新建
type SourcesImageRequestAdd struct {
	ID            int64               `json:"id" db:"id"`
	EnterpriseID  int64               `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	Name          string              `json:"name" db:"name" validate:"required" required:"图片名称必传"`
	URL           string              `json:"url" db:"url" validate:"required" required:"图片地址必传"`
	Ext           string              `json:"ext" db:"ext"`
	Size          int64               `json:"size" db:"size"`
	Sort          int32               `json:"sort" db:"sort" comment:"排序"`
	RelationGoods []source.BrandGoods `json:"relation_goods"`
}

// SourcesImageRequestUpdate 图片素材操作
type SourcesImageRequestUpdate struct {
	IDs          []int64 `json:"ids" db:"ids"`
	Type         int     `json:"type"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
}

// SourcesFileRequest 文件素材请求
type SourcesFileRequest struct {
	ID           int64   `json:"id" db:"id"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	AdminID      int64   `json:"admin_id" db:"admin_id"`
	MemberID     int64   `json:"member_id" db:"member_id"`
	Name         string  `json:"name" db:"name" validate:"required" required:"文件名称必传"`
	URL          string  `json:"url" db:"url" validate:"required" required:"文件URL必传"`
	AttachName   string  `json:"attach_name" db:"attach_name"`
	AttachID     int64   `json:"attach_id" db:"attach_id"`
	FileType     int     `json:"file_type" db:"file_type"`
	Ext          string  `json:"ext" db:"ext" validate:"required" required:"文件后缀必传"`
	Size         int64   `json:"size" db:"size" validate:"required" required:"文件大小必传"`
	CategoryIds  []int64 `json:"category_ids" required:"产品类别"`
	StartTime    int64   `json:"start_time" db:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64   `json:"end_time" db:"end_time" validate:"required" required:"结束时间必传"`
}

// SourcesFileRequestUpdate 文件素材操作
type SourcesFileRequestUpdate struct {
	IDs          []int64 `json:"ids" db:"ids"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	Type         int     `json:"type"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	SourceName   string  `json:"source_name" db:"source_name"`
	SourceURL    string  `json:"source_url" db:"source_url"`
	AttachName   string  `json:"attach_name" db:"attach_name"`
	AttachID     int64   `json:"attach_id" db:"attach_id"`
	FileType     int     `json:"file_type" db:"file_type"`
	FileExt      string  `json:"file_ext" db:"file_ext"`
	Status       int     `json:"status" db:"status"`
}

// RequestTextList 文本素材列表请求结构
type RequestTextList struct {
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	Keyword      string  `json:"keyword" db:"content"`
	AdminID      int64   `json:"admin_id" db:"admin_id"`
	JustMe       int64   `json:"just_me" db:"admin_id"`
	CreatedStart int64   `json:"create_start" db:"create_start"`
	CreatedEnd   int64   `json:"create_end" db:"create_end"`
	Size         int     `json:"size" db:"size"`
	Page         int     `json:"page" db:"page"`
	GroupIDList  []int64 `json:"group_id_list" db:"group_id_list"`
	CategoryIds  []int64 `json:"category_ids"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// ResponseTextList 文本素材列表返回结构
type ResponseTextList []struct {
	ID            int64  `json:"id" db:"id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID       int64  `json:"group_id" db:"group_id"`
	GroupName     string `json:"group_name" db:"group_name"`
	Content       string `json:"content" db:"content" desc:"文本内容"`
	Type          int    `json:"type" db:"type"`
	AdminID       int64  `json:"admin_id" db:"admin_id"`
	CreatedAt     int    `json:"created_at" db:"created_at"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	AdminName     string `json:"admin_name" db:"admin_name"`
	CategoryIds   string `json:"category_ids" db:"category_ids"`
	CategoryNames string `json:"category_names" db:"category_names"`
	SourceStatus  int    `json:"source_status" db:"source_status"`
}

// RequestLinkList 链接素材列表请求结构
type RequestLinkList struct {
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	Keyword      string  `json:"keyword" db:"title"`
	AdminID      int64   `json:"admin_id" db:"admin_id"`
	JustMe       int64   `json:"just_me" db:"admin_id"`
	CreatedStart int     `json:"create_start" db:"create_start"`
	CreatedEnd   int     `json:"create_end" db:"create_end"`
	Size         int     `json:"size" db:"size"`
	Page         int     `json:"page" db:"page"`
	GroupIDList  []int64 `json:"group_id_list" db:"group_id_list"`
	CategoryIds  []int64 `json:"category_ids"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// ResponseLinkList 链接素材列表返回结构
type ResponseLinkList []*struct {
	ID            int64    `json:"id" db:"id"`
	EnterpriseID  int64    `json:"enterprise_id" db:"enterprise_id"`
	GroupID       int64    `json:"group_id" db:"group_id"`
	GroupName     string   `json:"group_name" db:"group_name"`
	Title         string   `json:"title" db:"title" desc:"标题"`
	Description   string   `json:"description" db:"description" desc:"描述"`
	CoverUrl      string   `json:"cover_url" db:"cover_url" desc:"描述"`
	LinkType      int      `json:"link_type" db:"link_type" desc:"链接类型"`
	LinkUrl       string   `json:"link_url" db:"link_url" desc:"链接URL"`
	Author        string   `json:"author" db:"author" desc:"作者"`
	Content       string   `json:"content" db:"content" desc:"文本内容"`
	Type          int      `json:"type" db:"type"`
	AdminID       int64    `json:"admin_id" db:"admin_id"`
	CreatedAt     int      `json:"created_at" db:"created_at"`
	StartTime     int64    `json:"start_time" db:"start_time"`
	EndTime       int64    `json:"end_time" db:"end_time"`
	AdminName     string   `json:"admin_name" db:"admin_name"`
	CategoryIds   string   `json:"category_ids" db:"category_ids"`
	CategoryNames string   `json:"category_names" db:"category_names"`
	SourceStatus  int      `json:"source_status" db:"source_status"`
	SkuList       []string `json:"sku_list"`
}

// ResponseImgList 图片类型返回值
type ResponseImgList struct {
	ID              int64    `json:"id" db:"id"`
	EnterpriseID    int64    `json:"enterprise_id" db:"enterprise_id"`
	GroupID         int64    `json:"group_id" db:"group_id"`
	Name            string   `json:"name" db:"name"`
	URL             string   `json:"url" db:"url"`
	Type            int      `json:"type" db:"type"`
	AdminID         int64    `json:"admin_id" db:"admin_id"`
	MediaID         string   `json:"media_id" db:"media_id"`
	MediaExpireTime int64    `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64    `json:"media_status" db:"media_status"`
	CreatedAt       int      `json:"created_at" db:"created_at"`
	StartTime       int64    `json:"start_time" db:"start_time"`
	EndTime         int64    `json:"end_time" db:"end_time"`
	AdminName       string   `json:"admin_name" db:"admin_name"`
	CategoryIds     string   `json:"category_ids" db:"category_ids"`
	CategoryNames   string   `json:"category_names" db:"category_names"`
	SourceStatus    int      `json:"source_status" db:"source_status"`
	SkuList         []string `json:"sku_list"`
}

// ResponseFileList 文件类型List返回值
type ResponseFileList struct {
	ID              int64  `json:"id" db:"id"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID         int64  `json:"group_id" db:"group_id"`
	Name            string `json:"name" db:"name"`
	URL             string `json:"url" db:"url"`
	Ext             string `json:"ext" db:"ext"`
	FileType        int    `json:"file_type" db:"file_type"`
	AdminID         int64  `json:"admin_id" db:"admin_id"`
	MediaID         string `json:"media_id" db:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64  `json:"media_status" db:"media_status"`
	Size            int    `json:"size" db:"size"`
	CreatedAt       int    `json:"created_at" db:"created_at"`
	StartTime       int64  `json:"start_time" db:"start_time"`
	EndTime         int64  `json:"end_time" db:"end_time"`
	AdminName       string `json:"admin_name" db:"admin_name"`
	CategoryIds     string `json:"category_ids" db:"category_ids"`
	CategoryNames   string `json:"category_names" db:"category_names"`
	SourceStatus    int    `json:"source_status" db:"source_status"`
}

// RequestTextSave 文本素材保存请求格式
type RequestTextSave struct {
	ID           int64   `json:"id" db:"id"`
	Type         int     `json:"type" db:"type"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	Content      string  `json:"content" required:"文本内容必传" db:"content"`
	EnterpriseId int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	AdminId      int64   `json:"admin_id" db:"admin_id"`
	MemberID     int64   `json:"member_id" db:"member_id"`
	CreatedAt    int64   `json:"created_at" db:"created_at"`
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`
	CategoryIds  []int64 `json:"category_ids" required:"产品类别"`
	StartTime    int64   `json:"start_time" db:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64   `json:"end_time" db:"end_time" validate:"required" required:"结束时间必传"`
	MergeId      int64   `json:"merge_id" db:"merge_id"`
	Sort         int64   `json:"sort" db:"sort"`
}

// ResponseTextSave 文本素材保存返回格式
type ResponseTextSave struct {
	ID int64 `json:"id" desc:"新增ID"`
}

// RequestLinkSave 链接素材保存请求格式
type RequestLinkSave struct {
	ID            int64               `json:"id" db:"id"`
	Type          int                 `json:"type" db:"type"`
	GroupID       int64               `json:"group_id" db:"group_id"`
	Title         string              `json:"title" required:"素材标题必传" db:"title"`
	Description   string              `json:"description" required:"素材描述必传" db:"description"`
	LinkType      int                 `json:"link_type" db:"link_type" validate:"required,max=2" required:"链接类型必传" max:"链接类型不合法"`
	CoverUrl      string              `json:"cover_url" required:"封面图URL必传" db:"cover_url"`
	LinkUrl       string              `json:"link_url" db:"link_url"`
	Author        string              `json:"author" db:"author"`
	Content       string              `json:"content" db:"content"`
	AdminId       int64               `json:"admin_id" db:"admin_id"`
	MemberID      int64               `json:"member_id" db:"member_id"`
	EnterpriseId  int64               `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	CreatedAt     int64               `json:"created_at" db:"created_at"`
	UpdatedAt     int64               `json:"updated_at" db:"updated_at"`
	CategoryIds   []int64             `json:"category_ids" required:"产品类别"`
	StartTime     int64               `json:"start_time" db:"start_time" validate:"required" required:"开始时间必传"`
	EndTime       int64               `json:"end_time" db:"end_time" validate:"required" required:"结束时间必传"`
	MergeId       int64               `json:"merge_id" db:"merge_id"`
	Sort          int64               `json:"sort" db:"sort"`
	RelationGoods []source.BrandGoods `json:"relation_goods"`
}

// ResponseLinkSave 链接素材保存返回格式
type ResponseLinkSave struct {
	ID int64 `json:"id" desc:"新增ID"`
}

// 素材详情 请求参数
type RequestSourceDetail struct {
	ID int64 `json:"id" db:"id" validate:"required" required:"素材ID不能为空"`
}

// ResponseTextDetail 文本素材详情 返回参数
type ResponseTextDetail struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	Content      string `json:"content" db:"content" desc:"文本内容"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
}

// ResponseLinkDetail 链接素材详情 返回参数
type ResponseLinkDetail struct {
	ID            int64               `json:"id" db:"id"`
	EnterpriseID  int64               `json:"enterprise_id" db:"enterprise_id"`
	GroupID       int64               `json:"group_id" db:"group_id"`
	GroupName     string              `json:"group_name" db:"group_name"`
	Title         string              `json:"title" db:"title" desc:"标题"`
	Description   string              `json:"description" db:"description" desc:"描述"`
	CoverUrl      string              `json:"cover_url" db:"cover_url" desc:"描述"`
	LinkType      int                 `json:"link_type" db:"link_type" desc:"链接类型"`
	LinkUrl       string              `json:"link_url" db:"link_url" desc:"链接URL"`
	Author        string              `json:"author" db:"author" desc:"作者"`
	Content       string              `json:"content" db:"content" desc:"文本内容"`
	Type          int                 `json:"type" db:"type"`
	AdminID       int64               `json:"admin_id" db:"admin_id"`
	CreatedAt     int                 `json:"created_at" db:"created_at"`
	StartTime     int64               `json:"start_time" db:"start_time"`
	EndTime       int64               `json:"end_time" db:"end_time"`
	RelationGoods []source.BrandGoods `json:"relation_goods"`
}

// RequestSourceMini 小程序素材新增
type RequestSourceMini struct {
	ID            int64   `json:"id" db:"id"`
	EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID       int64   `json:"group_id" db:"group_id"`
	Title         string  `json:"title" required:"卡片标题必传" db:"title"`
	MiniProgramID int64   `json:"mini_program_id" validate:"required" required:"小程序ID必传" db:"mini_program_id"`
	Path          string  `json:"path"  validate:"required" required:"链接类型必传" db:"path"`
	CoverUrl      string  `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
	AdminID       int64   `json:"admin_id" db:"admin_id"`
	MemberID      int64   `json:"member_id" db:"member_id"`
	CreatedAt     int64   `json:"created_at" db:"created_at"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
	CategoryIds   []int64 `json:"category_ids"`
	StartTime     int64   `json:"start_time" validate:"required" required:"开始时间必传"`
	EndTime       int64   `json:"end_time" validate:"required" required:"结束时间必传"`
}

// SourceMini 小程序素材
type SourceMini struct {
	ID            int64  `json:"id" db:"id"`
	Type          int    `json:"type" db:"type"`
	AdminID       int64  `json:"admin_id" db:"admin_id"`
	GroupID       int64  `json:"group_id" db:"group_id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	Title         string `json:"title" db:"title"`
	MiniProgramID int64  `json:"mini_program_id" db:"mini_program_id"`
	AppName       string `json:"app_name"  db:"app_name"`
	AppID         string `json:"app_id"  db:"app_id"`
	OriginID      string `json:"origin_id" db:"origin_id"`
	Logo          string `json:"logo" db:"logo"`
	Path          string `json:"path" db:"path"`
	CoverUrl      string `json:"cover_url" db:"cover_url"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	AdminName     string `json:"admin_name" db:"admin_name"`
	SourceStatus  int    `json:"source_status" db:"source_status"`
}

// RequestSourceVideo 视频素材新增
type RequestSourceVideo struct {
	ID            int64               `json:"id" db:"id"`
	EnterpriseID  int64               `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	GroupID       int64               `json:"group_id" db:"group_id"`
	CoverUrl      string              `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
	Name          string              `json:"name" validate:"required" required:"视频名称必传" db:"name"`
	Size          int64               `json:"size" validate:"required" required:"视频大小必传" db:"size"`
	Url           string              `json:"url"  validate:"required" required:"视频Url必传" db:"url"`
	Ext           string              `json:"ext" validate:"required" required:"视频类型必传" db:"ext"`
	AdminID       int64               `json:"admin_id" db:"admin_id"`
	MemberID      int64               `json:"member_id" db:"member_id"`
	CategoryIds   []int64             `json:"category_ids" required:"产品类别"`
	StartTime     int64               `json:"start_time" db:"start_time" validate:"required" required:"开始时间必传"`
	EndTime       int64               `json:"end_time" db:"end_time" validate:"required" required:"结束时间必传"`
	CreatedAt     int64               `json:"created_at" db:"created_at"`
	UpdatedAt     int64               `json:"updated_at" db:"updated_at"`
	MergeId       int64               `json:"merge_id" db:"merge_id"`
	Sort          int64               `json:"sort" db:"sort"`
	RelationGoods []source.BrandGoods `json:"relation_goods"`
}

// SourceVideo 视频素材
type SourceVideo struct {
	ID              int64    `json:"id" db:"id"`
	Type            int      `json:"type" db:"type"`
	GroupID         int64    `json:"group_id" db:"group_id"`
	AdminID         int64    `json:"admin_id" db:"admin_id"`
	EnterpriseID    int64    `json:"enterprise_id" db:"enterprise_id"`
	Name            string   `json:"name"  db:"name"`
	CoverUrl        string   `json:"cover_url" db:"cover_url"`
	Size            int64    `json:"size"  db:"size"`
	Ext             string   `json:"ext"  db:"ext"`
	Url             string   `json:"url"  db:"url"`
	Status          int      `json:"status"  db:"status"`
	MediaID         string   `json:"media_id" db:"media_id"`
	MediaExpireTime int64    `json:"media_expire_time" db:"media_expire_time"`
	MediaStatus     int64    `json:"media_status" db:"media_status"`
	CreatedAt       int64    `json:"created_at" db:"created_at"`
	UpdatedAt       int64    `json:"updated_at" db:"updated_at"`
	DeletedAt       int64    `json:"deleted_at" db:"deleted_at"`
	StartTime       int64    `json:"start_time" db:"start_time"`
	EndTime         int64    `json:"end_time" db:"end_time"`
	AdminName       string   `json:"admin_name" db:"admin_name"`
	CategoryIds     string   `json:"category_ids" db:"category_ids"`
	CategoryNames   string   `json:"category_names" db:"category_names"`
	SourceStatus    int      `json:"source_status" db:"source_status"`
	SkuList         []string `json:"sku_list"`
}

// RequestSourceSearch 素材列表 搜索参数
type RequestSourceSearch struct {
	Page            int     `json:"page"`
	Size            int     `json:"size"`
	JustMe          int     `json:"just_me"`
	CreateStart     int64   `json:"create_start" db:"create_start"`
	CreateEnd       int64   `json:"create_end" db:"create_end"`
	EnterpriseID    int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	MediaExpireTime int64   `json:"media_expire_time" db:"media_expire_time"`
	GroupID         int64   `json:"group_id" db:"group_id"`
	Keyword         string  `json:"keyword" db:"keyword"`
	AdminID         int64   `json:"admin_id" db:"admin_id"`
	MemberUserId    string  `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	ExternalUserid  string  `json:"external_userid" db:"external_userid" desc:"ws外部联系人openid"`
	MemberID        int64   `json:"member_id" db:"member_id" desc:"ws员工Id"`
	GroupChatID     string  `json:"group_chat_id" db:"group_chat_id" desc:"ws群id"`
	GroupIDList     []int64 `json:"group_id_list" db:"group_id_list"`
	CategoryIds     []int64 `json:"category_ids"`
	Order           string  `json:"order"`
	Sort            string  `json:"sort"`
	SourceStatus    int     `json:"source_status"`
}

// RequestSourceDelete 素材删除 请求参数
type RequestSourceDelete struct {
	IDs []int64 `json:"ids" db:"ids" validate:"required" required:"素材IDs不能为空"`
}

// ResponseSourceDelete 素材删除 返回参数
type ResponseSourceDelete struct {
	Rows int64 `json:"rows"`
}

// RequestGroupMove 移动分组 请求参数
type RequestGroupMove struct {
	GroupID      int64   `json:"group_id" db:"group_id" validate:"required" required:"分组ID不能为空"`
	EnterpriseID int64   `json:"enterprise_id"`
	SourceIDs    []int64 `json:"source_ids" db:"source_ids"`
}

// UpdateSourceMediaRequest 更新素材的临时素材media状态
type UpdateSourceMediaRequest struct {
	ID              int64  `json:"id" db:"id"`
	MediaID         string `json:"media_id" db:"media_id"`
	SourceType      string `json:"source_type"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	UpdateAt        int64  `json:"update_at" db:"updated_at"`
	MediaStatus     int    `json:"media_status" db:"media_status"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaFailReason string `json:"media_fail_reason" db:"media_fail_reason"`
}

// MakeMediaRequest 给某个素材生成临时素材ID
type MakeMediaRequest struct {
	SourceID     int64  `json:"source_id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	SourceName   string `json:"source_name"`
	ResourceKey  string `json:"resource_key"`
}

// MediaItemDetail 查询 source_image、source_video 、source_file 表中企业微信临时素材信息
type MediaItemDetail struct {
	SourceID        int64  `json:"source_id" db:"source_id"`
	MediaExpireTime int64  `json:"media_expire_time" db:"media_expire_time"`
	MediaID         string `json:"media_id" db:"media_id"`
	ResourceKey     string `json:"resource_key" db:"resource_key"`
}

// RequestSourceFriend 朋友圈素材
type RequestSourceFriend struct {
	EnterpriseID int64                     `json:"enterprise_id" validate:"required" required:"企业ID必传"`
	Content      string                    `json:"content" required:"文本内容必传"`
	GroupID      int64                     `json:"group_id"`
	Images       []*SourcesImageRequestAdd `json:"images" validate:"required" required:"图片必传"`
	CoverUrl     string                    `json:"cover_url" validate:"required" required:"封面图URL必传"`
	Name         string                    `json:"name" validate:"required" required:"名称必传"`
	Size         int64                     `json:"size" validate:"required" required:"大小必传"`
	Url          string                    `json:"url"  validate:"required" required:"Url必传"`
	Ext          string                    `json:"ext" validate:"required" required:"类型必传"`
	Type         int32                     `json:"type" comment:"类型 1:图片 2:视频"`
	AdminID      int64                     `json:"admin_id"`
	MemberID     int64                     `json:"member_id" db:"member_id"`
	CategoryIds  []int64                   `json:"category_ids"`
	StartTime    int                       `json:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int                       `json:"end_time" validate:"required" required:"结束时间必传"`
}

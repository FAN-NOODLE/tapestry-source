package source

//// 素材主表
//type Source struct {
//	ID int64 `json:"id" db:"id"`
//	// EnterpriseID 企业id
//	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
//	// GroupID 分组ID 0 未分组
//	GroupID int64 `json:"group_id" db:"group_id"`
//	// Type 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
//	Type int8 `json:"type" db:"type"`
//	// AdminID 管理员ID
//	AdminID int64 `json:"admin_id" db:"admin_id"`
//	// CreatedAt 创建时间
//	CreatedAt int64 `json:"created_at" db:"created_at"`
//	// UpdatedAt 更新时间
//	UpdatedAt int64 `json:"updated_at" db:"updated_at"`
//	// DeletedAt 删除时间
//	DeletedAt int64 `json:"deleted_at" db:"deleted_at"`
//}
//
//// 图片素材
//type SourcesImage struct {
//	ID int64 `json:"id"`
//	// EnterpriseID 企业id
//	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
//	// SourceID 素材主表id
//	SourceID int64 `json:"source_id" db:"source_id"`
//	// SourceName 图片名称
//	SourceName string `json:"source_name" db:"source_name"`
//	// SourceURL 图片url
//	SourceURL string `json:"source_url" db:"source_url"`
//	// CreatedAt 创建时间
//	CreatedAt int64 `json:"created_at" db:"created_at"`
//	// UpdatedAt 更新时间
//	UpdatedAt int64 `json:"updated_at" db:"updated_at"`
//	// DeletedAt 删除时间
//	DeletedAt int64 `json:"deleted_at" db:"deleted_at"`
//}
//
//// 文件素材
//type SourcesFile struct {
//	ID int64 `json:"id"`
//	// EnterpriseID 企业id
//	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
//	// SourceID 素材主表id
//	SourceID int64 `json:"source_id" db:"source_id"`
//	// SourceName 图片名称
//	SourceName string `json:"source_name" db:"source_name"`
//	// SourceURL 图片url
//	SourceURL string `json:"source_url" db:"source_url"`
//	// AttachName 第三方文件名称
//	AttachName string `json:"attach_name" db:"attach_name"`
//	// AttachID 第三方file ID
//	AttachID int64 `json:"attach_id" db:"attach_id"`
//	// FileType 文件类型 1 word 2 excel 3 ppt 4 pdf 5 txt 6 其他
//	FileType int `json:"file_type" db:"file_type"`
//	// FileExt 文件后缀
//	FileExt string `json:"file_ext" db:"file_ext"`
//	// Status 上传状态（1=转码中；2=正常）
//	Status int `json:"status" db:"status"`
//	// Size 文件大小
//	Size int `json:"size" db:"size"`
//	// CreatedAt 创建时间
//	CreatedAt int64 `json:"created_at"`
//	// UpdatedAt 更新时间
//	UpdatedAt int64 `json:"updated_at"`
//	// DeletedAt 删除时间
//	DeletedAt int64 `json:"deleted_at"`
//}
//
//// List请求
//type RequestList struct {
//	Page         int    `json:"page" db:"page" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
//	Size         int    `json:"size" db:"size" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
//	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Type         int    `json:"type" db:"type" validate:"required,max=6" required:"素材类型必传" max:"素材类型不合法"`
//	Keyword      string `json:"keyword" db:"keyword"`
//	StartAt      int    `json:"start_at" db:"start_at"`
//	EndAt        int    `json:"end_at" db:"end_at"`
//}
//
//// 图片素材请求
//type SourcesImageRequest struct {
//	EnterpriseID int64                    `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64                    `json:"group_id" db:"group_id"`
//	AdminID      int64                    `json:"admin_id" db:"admin_id"`
//	Images       []SourcesImageRequestAdd `json:"images"`
//}
//
//// 图片素材新建
//type SourcesImageRequestAdd *struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	SourceID     int64  `json:"source_id" db:"source_id"`
//	SourceName   string `json:"source_name" db:"source_name"`
//	SourceURL    string `json:"source_url" db:"source_url"`
//}
//
//// 图片素材操作
//type SourcesImageRequestUpdate struct {
//	IDs          []int64 `json:"ids" db:"ids"`
//	Type         int     `json:"type"`
//	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//}
//
//// 文件素材请求
//type SourcesFileRequest struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//	SourceName   string `json:"source_name" db:"source_name"`
//	SourceURL    string `json:"source_url" db:"source_url"`
//	AttachName   string `json:"attach_name" db:"attach_name"`
//	AttachID     int64  `json:"attach_id" db:"attach_id"`
//	FileType     int    `json:"file_type" db:"file_type"`
//	FileExt      string `json:"file_ext" db:"file_ext"`
//	Status       int    `json:"status" db:"status"`
//	Size         int    `json:"size" db:"size"`
//}
//
//// 文件素材操作
//type SourcesFileRequestUpdate struct {
//	IDs          []int64 `json:"ids" db:"ids"`
//	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	Type         int     `json:"type"`
//	GroupID      int64   `json:"group_id" db:"group_id"`
//	SourceName   string  `json:"source_name" db:"source_name"`
//	SourceURL    string  `json:"source_url" db:"source_url"`
//	AttachName   string  `json:"attach_name" db:"attach_name"`
//	AttachID     int64   `json:"attach_id" db:"attach_id"`
//	FileType     int     `json:"file_type" db:"file_type"`
//	FileExt      string  `json:"file_ext" db:"file_ext"`
//	Status       int     `json:"status" db:"status"`
//}
//
//// 小程序类型List返回值
//type ResponseMiniProgramList []struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Name         string `json:"name" db:"name"`
//	Type         int    `json:"type" db:"type"`
//	Content      string `json:"content" db:"content"`
//	//ContentObj   MiniProgramSubContent `json:"content_obj"`
//	ObjectID    int64  `json:"object_id" db:"object_id"`
//	ProgramName string `json:"program_name" db:"program_name"`
//	AdminID     int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt   int    `json:"created_at" db:"created_at"`
//	UpdatedAt   int    `json:"updated_at" db:"updated_at"`
//}
//
//// 文本素材列表请求结构
//type RequestTextList struct {
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Keyword      string `json:"keyword" db:"content"`
//	JustMe       int64  `json:"just_me" db:"admin_id"`
//	CreatedStart int    `json:"created_start" db:"created_at"`
//	CreatedEnd   int    `json:"created_end" db:"created_at"`
//}
//
//// 文本素材列表返回结构
//type ResponseTextList []struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	GroupName    string `json:"group_name" db:"group_name"`
//	Content      string `json:"content" db:"content" desc:"文本内容"`
//	Type         int    `json:"type" db:"type"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt    int    `json:"created_at" db:"created_at"`
//}
//
//// 链接素材列表请求结构
//type RequestLinkList struct {
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Keyword      string `json:"keyword" db:"content"`
//	JustMe       int64  `json:"just_me" db:"admin_id"`
//	CreatedStart int    `json:"created_start" db:"created_at"`
//	CreatedEnd   int    `json:"created_end" db:"created_at"`
//}
//
//// 链接素材列表返回结构
//type ResponseLinkList []struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	GroupName    string `json:"group_name" db:"group_name"`
//	Title        string `json:"title" db:"title" desc:"标题"`
//	Description  string `json:"description" db:"description" desc:"描述"`
//	CoverUrl     string `json:"cover_url" db:"cover_url" desc:"描述"`
//	LinkType     int    `json:"link_type" db:"link_type" desc:"链接类型"`
//	LinkUrl      string `json:"link_url" db:"link_url" desc:"链接URL"`
//	Author       string `json:"author" db:"author" desc:"作者"`
//	Content      string `json:"content" db:"content" desc:"文本内容"`
//	Type         int    `json:"type" db:"type"`
//	Creator      string `json:"creator" db:"creator"`
//	CreatedAt    int    `json:"created_at" db:"created_at"`
//}
//
//// 图片类型返回值
//type ResponseImgList struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	SourceName   string `json:"source_name" db:"source_name"`
//	SourceURL    string `json:"source_url" db:"source_url"`
//	Type         int    `json:"type" db:"type"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt    int    `json:"created_at" db:"created_at"`
//}
//
//// 文件类型List返回值
//type ResponseFileList struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	SourceName   string `json:"source_name" db:"source_name"`
//	SourceURL    string `json:"source_url" db:"source_url"`
//	FileExt      string `json:"file_ext"`
//	FileType     int    `json:"file_type"`
//	Creator      string `json:"creator" db:"creator"`
//	Status       int    `json:"status" db:"status"`
//	CreatedAt    int    `json:"created_at" db:"created_at"`
//}
//
//// RequestTextSave 文本素材保存请求格式
//type RequestTextSave struct {
//	ID           int64  `json:"id" db:"id"`
//	Type         int    `json:"type" db:"type" validate:"required,max=6" required:"素材类型必传" max:"素材类型不合法"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Content      string `json:"content" required:"文本内容必传" db:"content"`
//	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	AdminId      int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt    int64  `json:"created_at" db:"created_at"`
//	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
//}
//
////ResponseTextSave 文本素材保存返回格式
//type ResponseTextSave struct {
//	ID int64 `json:"id" desc:"新增ID"`
//}
//
//// RequestLinkSave 链接素材保存请求格式
//type RequestLinkSave struct {
//	ID           int64  `json:"id" db:"id"`
//	Type         int    `json:"type" db:"type" validate:"required,max=6" required:"素材类型必传" max:"素材类型不合法"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Title        string `json:"title" required:"素材标题必传" db:"title"`
//	Description  string `json:"description" required:"素材描述必传" db:"description"`
//	LinkType     int    `json:"link_type" db:"link_type" validate:"required,max=2" required:"链接类型必传" max:"链接类型不合法"`
//	CoverUrl     string `json:"cover_url" required:"封面图URL必传" db:"cover_url"`
//	LinkUrl      string `json:"link_url" db:"link_url"`
//	Author       string `json:"author" db:"author"`
//	Content      string `json:"content" db:"content"`
//	AdminId      int64  `json:"admin_id" db:"admin_id"`
//	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	CreatedAt    int64  `json:"created_at" db:"created_at"`
//	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
//}
//
//// ResponseLinkSave 链接素材保存返回格式
//type ResponseLinkSave struct {
//	ID int64 `json:"id" desc:"新增ID"`
//}
//
//// RequestSourceMini 小程序素材新增
//type RequestSourceMini struct {
//	ID            int64  `json:"id" db:"id"`
//	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID       int64  `json:"group_id" db:"group_id"`
//	Title         string `json:"title" required:"卡片标题必传" db:"title"`
//	MiniProgramID int64  `json:"mini_program_id" validate:"required" required:"小程序ID必传" db:"mini_program_id"`
//	Path          string `json:"path"  validate:"required" required:"链接类型必传" db:"path"`
//	CoverUrl      string `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
//	AdminID       int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt     int64  `json:"created_at" db:"created_at"`
//	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
//}
//
//// 小程序素材
//type SourceMini struct {
//	ID            int64  `json:"id" db:"id"`
//	Type          int    `json:"type" db:"type"`
//	AdminID       int64  `json:"admin_id" db:"admin_id"`
//	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
//	Title         string `json:"title" db:"title"`
//	MiniProgramID int64  `json:"mini_program_id" db:"mini_program_id"`
//	AppName       string `json:"app_name"  db:"app_name"`
//	Path          string `json:"path" db:"path"`
//	CoverUrl      string `json:"cover_url" db:"cover_url"`
//	CreatedAt     int64  `json:"created_at" db:"created_at"`
//	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
//	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
//}
//
//// RequestSourceVideo 视频素材新增
//type RequestSourceVideo struct {
//	ID           int64  `json:"id" db:"id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	CoverUrl     string `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
//	Name         string `json:"name" validate:"required" required:"视频名称必传" db:"name"`
//	Size         string `json:"size" validate:"required" required:"视频大小必传" db:"size"`
//	Url          string `json:"url"  validate:"required" required:"视频Url必传" db:"url"`
//	Ext          string `json:"ext" validate:"required" required:"视频类型必传" db:"ext"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//	CreatedAt    int64  `json:"created_at" db:"created_at"`
//	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
//}
//
//// SourceVideo 视频素材
//type SourceVideo struct {
//	ID           int64  `json:"id" db:"id"`
//	Type         int    `json:"type" db:"type"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
//	Name         string `json:"name"  db:"name"`
//	CoverUrl     string `json:"cover_url" db:"cover_url"`
//	Size         string `json:"size"  db:"size"`
//	Ext          string `json:"ext"  db:"ext"`
//	Url          string `json:"url"  db:"url"`
//	Status       int    `json:"status"  db:"status"`
//	CreatedAt    int64  `json:"created_at" db:"created_at"`
//	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
//	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
//}
//
//// 素材列表 搜索参数
//type RequestSourceSearch struct {
//	Page         int    `json:"page"`
//	Size         int    `json:"size"`
//	JustMe       int    `json:"just_me"`
//	CreateStart  int64  `json:"create_start" db:"create_start"`
//	CreateEnd    int64  `json:"create_end" db:"create_end"`
//	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
//	GroupID      int64  `json:"group_id" db:"group_id"`
//	Keyword      string `json:"keyword" db:"keyword"`
//	AdminID      int64  `json:"admin_id" db:"admin_id"`
//}


type GetGroupName struct {
	GroupID      int64 `json:"group_id"`
	EnterpriseID int64 `json:"enterprise_id"`
}
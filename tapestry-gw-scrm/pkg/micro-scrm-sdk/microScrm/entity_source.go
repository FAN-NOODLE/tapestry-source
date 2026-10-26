package microScrm

// 素材列表 搜索参数
type RequestSourceSearch struct {
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	JustMe       int    `json:"just_me"`
	CreateStart  int64  `json:"create_start" db:"create_start"`
	CreateEnd    int64  `json:"create_end" db:"create_end"`
	EnterpriseID int64  `json:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	Keyword      string `json:"keyword" db:"keyword"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
}

// 图片素材新增
type SourcesImageRequestAdd *struct {
	ID            int64        `json:"id" db:"id"`
	EnterpriseID  int64        `json:"enterprise_id" db:"enterprise_id"`
	Name          string       `json:"name" db:"name"`
	URL           string       `json:"url" db:"url"`
	Ext           string       `json:"ext" db:"ext"`
	Size          int          `json:"size" db:"size"`
	Sort          int32        `json:"sort" comment:"排序"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}

type BrandGoods struct {
	Name    string `json:"name"`
	Sku     string `json:"sku"`
	Upc     string `json:"upc"`
	Url     string `json:"url"`
	Color   string `json:"color"`
	AddType int    `json:"add_type"` // 1 手动  2 阿里图像识别
	Score   string `json:"score" comment:"置信度（AI识别）"`
}

// 图片素材请求
type SourcesImageRequest struct {
	EnterpriseID int64                    `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64                    `json:"group_id" db:"group_id"`
	AdminID      int64                    `json:"admin_id" db:"admin_id"`
	Images       []SourcesImageRequestAdd `json:"images"`
	CategoryIds  []int64                  `json:"category_ids" required:"产品类别"`
	StartTime    int64                    `json:"start_time" db:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64                    `json:"end_time" db:"end_time" validate:"required" required:"结束时间必传"`
	ShowType     int                      `json:"show_type" comment:"展示类型:1单张,2合并"`
}

// 文件素材请求
type SourcesFileRequest struct {
	ID           int64   `json:"id" db:"id"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	AdminID      int64   `json:"admin_id" db:"admin_id"`
	Name         string  `json:"name" db:"name" validate:"required" required:"文件名称必传"`
	URL          string  `json:"url" db:"url" validate:"required" required:"文件地址必传"`
	AttachName   string  `json:"attach_name" db:"attach_name"`
	AttachID     int64   `json:"attach_id" db:"attach_id"`
	FileType     int     `json:"file_type" db:"file_type"`
	Ext          string  `json:"ext" db:"ext" validate:"required" required:"文件后缀必传"`
	Size         int     `json:"size" db:"size" validate:"required" required:"文件大小必传"`
	CategoryIds  []int64 `json:"category_ids" required:"产品类别"`
	StartTime    int64   `json:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64   `json:"end_time" validate:"required" required:"结束时间必传"`
}

// List请求
type RequestList struct {
	Page         int    `json:"page" db:"page" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
	Size         int    `json:"size" db:"size" validate:"required,min=1" required:"分页必传" min:"分页大于1"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	JustMe       int    `json:"just_me" db:"just_me"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	Keyword      string `json:"keyword" db:"keyword"`
	CategoryId   int64  `json:"category_id"`
	CreateStart  int    `json:"create_start" db:"create_start"`
	CreateEnd    int    `json:"create_end" db:"create_end"`
}

// 图片类型返回值
type ResponseImgList []struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	Name         string `json:"name" db:"name"`
	URL          string `json:"url" db:"url"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedAt    int    `json:"created_at" db:"created_at"`
}

// 文件类型List返回值
type ResponseFileList []struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	Name         string `json:"name" db:"name"`
	URL          string `json:"url" db:"url"`
	Ext          string `json:"ext" db:"ext"`
	FileType     int    `json:"file_type" db:"file_type"`
	Creator      string `json:"creator" db:"creator"`
	CreatedAt    int    `json:"created_at" db:"created_at"`
}

// RequestTextSave 文本素材保存请求格式
type RequestTextSave struct {
	ID           int64   `json:"id" db:"id"`
	Type         int     `json:"type" db:"type"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	Content      string  `json:"content" required:"文本内容必传" db:"content"`
	EnterpriseId int64   `json:"enterprise_id" db:"enterprise_id"`
	AdminId      int64   `json:"admin_id" db:"admin_id"`
	CreatedAt    int64   `json:"created_at" db:"created_at"`
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`
	CategoryIds  []int64 `json:"category_ids" required:"产品类别"`
	StartTime    int64   `json:"start_time" validate:"required" required:"开始时间必传"`
	EndTime      int64   `json:"end_time" validate:"required" required:"结束时间必传"`
}

// ResponseTextSave 文本素材保存返回格式
type ResponseTextSave struct {
	ID int64 `json:"id" desc:"新增ID"`
}

// RequestLinkSave 链接素材保存请求格式
type RequestLinkSave struct {
	ID            int64        `json:"id" db:"id"`
	Type          int          `json:"type" db:"type"`
	GroupID       int64        `json:"group_id" db:"group_id"`
	Title         string       `json:"title" required:"素材标题必传" db:"title"`
	Description   string       `json:"description" required:"素材描述必传" db:"description"`
	LinkType      int          `json:"link_type" db:"link_type" validate:"required,max=2" required:"链接类型必传" max:"链接类型不合法"`
	CoverUrl      string       `json:"cover_url" required:"封面图URL必传" db:"cover_url"`
	LinkUrl       string       `json:"link_url" db:"link_url"`
	Author        string       `json:"author" db:"author"`
	Content       string       `json:"content" db:"content"`
	AdminId       int64        `json:"admin_id" db:"admin_id"`
	EnterpriseId  int64        `json:"enterprise_id" db:"enterprise_id"`
	CreatedAt     int64        `json:"created_at" db:"created_at"`
	UpdatedAt     int64        `json:"updated_at" db:"updated_at"`
	CategoryIds   []int64      `json:"category_ids"`
	StartTime     int64        `json:"start_time"`
	EndTime       int64        `json:"end_time"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}

// ResponseLinkSave 链接素材保存返回格式
type ResponseLinkSave struct {
	ID int64 `json:"id" desc:"新增ID"`
}

// 文本素材列表请求结构
type RequestTextList struct {
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64   `json:"group_id" db:"group_id"`
	Keyword      string  `json:"keyword" db:"content"`
	JustMe       int64   `json:"just_me" db:"admin_id"`
	AdminID      int64   `json:"admin_id" db:"admin_id"`
	CreatedStart int     `json:"create_start" db:"create_start"`
	CreatedEnd   int     `json:"create_end" db:"create_end"`
	Size         int     `json:"size" db:"size"`
	Page         int     `json:"page" db:"page"`
	CategoryIds  []int64 `json:"category_ids"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
}

// 文本素材列表返回结构
type ResponseTextList []struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	GroupName    string `json:"group_name" db:"group_name"`
	Content      string `json:"content" db:"content" desc:"文本内容"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedAt    int    `json:"created_at" db:"created_at"`
}

// 链接素材列表请求结构
type RequestLinkList struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	Keyword      string `json:"keyword" db:"content"`
	JustMe       int64  `json:"just_me" db:"admin_id"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedStart int    `json:"create_start" db:"create_start"`
	CreatedEnd   int    `json:"create_end" db:"create_end"`
	Size         int    `json:"size" db:"size"`
	Page         int    `json:"page" db:"page"`
}

// 链接素材列表返回结构
type ResponseLinkList []struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	GroupName    string `json:"group_name" db:"group_name"`
	Title        string `json:"title" db:"title" desc:"标题"`
	Description  string `json:"description" db:"description" desc:"描述"`
	CoverUrl     string `json:"cover_url" db:"cover_url" desc:"描述"`
	LinkType     int    `json:"link_type" db:"link_type" desc:"链接类型"`
	LinkUrl      string `json:"link_url" db:"link_url" desc:"链接URL"`
	Author       string `json:"author" db:"author" desc:"作者"`
	Content      string `json:"content" db:"content" desc:"文本内容"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedAt    int    `json:"created_at" db:"created_at"`
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
	CreatedAt    int    `json:"created_at" db:"created_at"`
}

// 链接素材详情 返回参数
type ResponseLinkDetail struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	GroupName    string `json:"group_name" db:"group_name"`
	Title        string `json:"title" db:"title" desc:"标题"`
	Description  string `json:"description" db:"description" desc:"描述"`
	CoverUrl     string `json:"cover_url" db:"cover_url" desc:"描述"`
	LinkType     int    `json:"link_type" db:"link_type" desc:"链接类型"`
	LinkUrl      string `json:"link_url" db:"link_url" desc:"链接URL"`
	Author       string `json:"author" db:"author" desc:"作者"`
	Content      string `json:"content" db:"content" desc:"文本内容"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	CreatedAt    int    `json:"created_at" db:"created_at"`
}

// 素材删除 请求参数
type RequestSourceDelete struct {
	IDs []int64 `json:"ids" db:"ids" validate:"required" required:"素材IDs不能为空"`
}

// 移动分组 请求参数
type RequestGroupMove struct {
	GroupID   int64   `json:"group_id" db:"group_id" validate:"required" required:"分组ID不能为空"`
	SourceIDs []int64 `json:"source_ids" db:"source_ids"`
}

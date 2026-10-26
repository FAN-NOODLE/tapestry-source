package source

// RequestTextList 文本素材列表请求结构
type RequestTextList struct {
	EnterpriseID int64   `json:"enterprise_id"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
	GroupID      int64   `json:"group_id"`
	Keyword      string  `json:"keyword"`
	JustMe       int64   `json:"just_me"`
	AdminID      int64   `json:"admin_id"`
	CreatedStart int     `json:"create_start"`
	CreatedEnd   int     `json:"create_end"`
	CategoryIds  []int64 `json:"category_ids"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// RequestList List请求
type RequestList struct {
	EnterpriseID int64   `json:"enterprise_id"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
	GroupID      int64   `json:"group_id"`
	JustMe       int     `json:"just_me"`
	AdminID      int64   `json:"admin_id"`
	Keyword      string  `json:"keyword"`
	CategoryIds  []int64 `json:"category_ids"`
	CreateStart  int     `json:"create_start"`
	CreateEnd    int     `json:"create_end"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// RequestLinkList 链接素材列表请求结构
type RequestLinkList struct {
	EnterpriseID int64   `json:"enterprise_id"`
	GroupID      int64   `json:"group_id"`
	Keyword      string  `json:"keyword"`
	JustMe       int64   `json:"just_me"`
	AdminID      int64   `json:"admin_id"`
	CreatedStart int     `json:"create_start"`
	CreatedEnd   int     `json:"create_end"`
	Size         int     `json:"size"`
	Page         int     `json:"page"`
	CategoryIds  []int64 `json:"category_ids"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// RequestSourceSearch 素材列表 搜索参数
type RequestSourceSearch struct {
	Page         int     `json:"page"`
	Size         int     `json:"size"`
	JustMe       int     `json:"just_me"`
	CreateStart  int64   `json:"create_start"`
	CreateEnd    int64   `json:"create_end"`
	EnterpriseID int64   `json:"enterprise_id"`
	GroupID      int64   `json:"group_id"`
	Keyword      string  `json:"keyword"`
	AdminID      int64   `json:"admin_id"`
	CategoryIds  []int64 `json:"category_ids"`
	Order        string  `json:"order"`
	Sort         string  `json:"sort"`
	SourceStatus int     `json:"source_status"`
}

// RequestDetail 素材详情请求结构
type RequestDetail struct {
	EnterpriseID int64 `json:"enterprise_id" comment:"企业ID"`
	SourceID     int64 `json:"source_id"`
	Type         int64 `json:"type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈)"`
}

// RequestUpdate 素材保存请求结构
type RequestUpdate struct {
	ID           int64   `json:"id"`
	EnterpriseID int64   `json:"enterprise_id"`
	AdminId      int64   `json:"admin_id"`
	GroupId      int64   `json:"group_id"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
	CategoryIds  []int64 `json:"category_ids"`
	SourceType   int64   `json:"source_type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈)"`

	// 文本素材
	Content string `json:"content" comment:"文本内容"`
	// (文件.视频)素材
	Name string `json:"name" comment:"文件名"`

	FileType    int    `json:"file_type" comment:"文件类型: 1图片,2视频"`
	TextContent string `json:"text_content" comment:"文本内容"`
	Sources     []struct {
		Id            int64        `json:"id"`
		Url           string       `json:"url"`
		IsOrigin      int          `json:"is_origin"`
		Name          string       `json:"name"`
		Ext           string       `json:"ext"`
		Size          int          `json:"size"`
		IsTranscode   int          `json:"is_transcode"`
		Sort          int          `json:"sort"`
		RelationGoods []BrandGoods `json:"relation_goods"`
	} `json:"sources"`

	// (小程序.链接)公用
	Title    string `json:"title" comment:"标题"`
	CoverUrl string `json:"cover_url" comment:"封面图"`

	// 小程序素材
	Path          string `json:"path" comment:"链接类型"`
	MiniProgramID int64  `json:"mini_program_id" comment:"小程序ID"`

	// 链接素材
	Description   string       `json:"description" comment:"链接描述"`
	LinkType      int          `json:"link_type" comment:"链接类型"`
	LinkUrl       string       `json:"link_url" comment:"跳转地址"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}

// RequestRemoveCategory 移除素材关联产品类别
type RequestRemoveCategory struct {
	SourceID     int64 `json:"source_id"`
	EnterpriseID int64 `json:"enterprise_id"`
	CategoryId   int64 `json:"category_id"`
}

type RequestUpdateStatus struct {
	IDs          []int64 `json:"ids" comment:"ids"`
	EnterpriseID int64   `json:"enterprise_id" comment:"企业ID"`
	SourceType   int     `json:"source_type" comment:"source_type"`
	SourceStatus int     `json:"source_status" comment:"source_status"`
}

// RequestConsoleList 后台管理各素材列表请求结构
type RequestConsoleList struct {
	EnterpriseID   int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	Keyword        string  `json:"keyword" db:"content"`
	UploaderName   string  `json:"uploader_name"`
	DepartmentName string  `json:"department_name"`
	AdminID        int64   `json:"admin_id" db:"admin_id"`
	JustMe         int     `json:"just_me" db:"just_me"`
	CreatedStart   int64   `json:"create_start" db:"create_start"`
	CreatedEnd     int64   `json:"create_end" db:"create_end"`
	Size           int     `json:"size" db:"size"`
	Page           int     `json:"page" db:"page"`
	GroupIDList    []int64 `json:"group_id_list" db:"group_id_list"`
	CategoryIds    []int64 `json:"category_ids"`
	Order          string  `json:"order"`
	Sort           string  `json:"sort"`
	SourceStatus   int     `json:"source_status"`
	SourceGroupID  int64   `json:"source_group_id" db:"source_group_id"`
	GroupIDs       []int64 `json:"group_ids" db:"group_ids"`
	Type           int     `json:"type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	IsUse          int     `json:"is_use"`
	SkuItem        string  `json:"sku_item"`
	ProductName    string  `json:"product_name"`
	SkuColor       string  `json:"sku_color"`
	GoodsRelation  int     `json:"goods_relation" comment:"商品关联过滤: 0-全部 1-已关联商品 2-未关联商品"`
	SourceOrigin   int     `json:"source_origin" comment:"素材来源: 0-全部 1-手动创建 2-CH同步"`
}

// RequestGetMediaInfo 获取上传media-id的参数
type RequestGetMediaInfo struct {
	EnterpriseID int64  `json:"enterprise_id" comment:"企业ID"`
	SourceID     int64  `json:"source_id"`
	SourceType   int64  `json:"source_type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	Url          string `json:"url" comment:"素材的cos地址"`
	Type         int64  `json:"type" comment:"0为临时素材，1为附件资源(用户朋友圈)"`
}

type SourceSetTop struct {
	EnterpriseID int64 `json:"enterprise_id" comment:"企业ID"`
	SourceID     int64 `json:"source_id"`
	SourceType   int32 `json:"source_type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	IsTop        int32 `json:"is_top" comment:"是否置顶,1是，2否"`
}

type BrandListReq struct {
	BrandCacheReq
	Size     int    `json:"size"`
	LastCode string `json:"last_code"`
	AdminId  int64  `json:"admin_id"`
}
type BrandCacheReq struct {
	EnterpriseId            int64  `json:"enterprise_id"`
	SkuItem                 string `json:"sku_item"`
	ProductName             string `json:"product_name"`
	SkuColor                string `json:"sku_color"`
	UpcCode                 string `json:"upc_code"`
	VariantCode             string `json:"variant_code"`
	StoreCode               string `json:"store_code"`
	Mobile                  string `json:"mobile"`
	IncludeNoneImageProduct bool   `json:"include_none_image_product"`
}

//// RequestLinkSaveReq  链接素材保存请求格式
//type RequestLinkSaveReq struct {
//	ID            int64        `json:"id" db:"id"`
//	Type          int          `json:"type" db:"type"`
//	GroupID       int64        `json:"group_id" db:"group_id"`
//	Title         string       `json:"title" required:"素材标题必传" db:"title"`
//	Description   string       `json:"description" required:"素材描述必传" db:"description"`
//	LinkType      int          `json:"link_type" db:"link_type" validate:"required,max=2" required:"链接类型必传" max:"链接类型不合法"`
//	CoverUrl      string       `json:"cover_url" required:"封面图URL必传" db:"cover_url"`
//	LinkUrl       string       `json:"link_url" db:"link_url"`
//	Author        string       `json:"author" db:"author"`
//	Content       string       `json:"content" db:"content"`
//	AdminId       int64        `json:"admin_id" db:"admin_id"`
//	EnterpriseId  int64        `json:"enterprise_id" db:"enterprise_id"`
//	CreatedAt     int64        `json:"created_at" db:"created_at"`
//	UpdatedAt     int64        `json:"updated_at" db:"updated_at"`
//	CategoryIds   []int64      `json:"category_ids"`
//	StartTime     int64        `json:"start_time"`
//	EndTime       int64        `json:"end_time"`
//	RelationGoods []BrandGoods `json:"relation_goods"`
//}

type BrandGoods struct {
	Name    string `json:"name"`
	Sku     string `json:"sku"`
	Upc     string `json:"upc"`
	Url     string `json:"url"`
	Color   string `json:"color"`
	AddType int8   `json:"add_type" comment:"添加方式 1 手动添加 2 AI识别"`
	Score   string `json:"score" comment:"置信度（AI识别）"`
}

type BrandImageAiSearchReq struct {
	AdminId      int64  `json:"admin_id" db:"admin_id"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	ImageUrl     string `json:"image_url"`
	Num          int    `json:"num"` //每个区域返回几个商品
}

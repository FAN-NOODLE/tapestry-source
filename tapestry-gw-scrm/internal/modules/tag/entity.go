package tag

// TAG_EXPORT_TYPE_AUTO 导出自动标签
const TAG_EXPORT_TYPE_AUTO = 1

// TAG_EXPORT_TYPE_MANUAL 导出手动标签
const TAG_EXPORT_TYPE_MANUAL = 2

// TAG_EXPORT_TYPE_ALL 导出自动、手动标签
const TAG_EXPORT_TYPE_ALL = 0

// Response 默认返回结果
type Response struct{}

// GoodsSpuRequest spu列表请求
type GoodsSpuRequest struct {
	Page         int    `json:"page" desc:"分页"`
	Size         int    `json:"size" desc:"单页查询数量"`
	SearchType   int    `json:"search_type" desc:"搜索方式（1:spuId搜索，2:商品名称搜索）"`
	SearchValue  string `json:"search" desc:"搜索内容"`
	EnterpriseId int64  `json:"enterprise_id" desc:"企业ID"`
}

// CouponRequest 优惠券请求
type CouponRequest struct {
	SearchIsAll  int    `json:"search_is_all" desc:"查询所有"`
	Page         int    `json:"page" desc:"分页"`
	Size         int    `json:"size" desc:"单页查询数量"`
	SearchText   string `json:"search_text" desc:"搜索条件"`
	EnterpriseId int64  `json:"enterprise_id" desc:"企业ID"`
}

// TagExportRequest 标签列表请求参数
type TagExportRequest struct {
	SearchType int64  `json:"search_type" desc:"检索类型 1:标签组名称 2:标签名称"`
	KeyWords   string `json:"key_words" desc:"检索关键字"`
	TagType    int64  `json:"tag_type" desc:"导出标签类型 1:自动标签 2:手动标签 0:全部"`
}

type UserGroupFileDetection struct {
	FileName     string `json:"file_name" desc:"文件名"`
	EnterpriseID int64  `json:"enterprise_id"`
}

type CountInfo struct {
	Count int64 `json:"id"`
}

type DownLoadFileInfo struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
	FileUrl      string `json:"file_url"`
}

type IDInfo struct {
	ID int64 `json:"id"`
}

type FileField struct {
	Unionid string `json:"unionid"`
}

type GetFileInfo struct {
	EnterpriseID int64       `json:"enterprise_id"`
	ID           int64       `json:"id"`
	FileName     string      `json:"file_name"`
	File         []FileField `json:"file"` //csv 文件
}

type NewIDInfo struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type ExportList struct {
	ExportData []ExportData
	Total      int64 `json:"total"`
}
type ExportData struct {
	UnionID   string `json:"union_id"`
	ImportRes string `json:"import_res"`
}

type ExportReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	UserGroupID  int64 `json:"user_group_id"`
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Total        int64 `json:"total"`
}

type PageVO struct {
	Page int `json:"page" comment:"页码"`
	Size int `json:"size" comment:"分页数"`
}

type InfoReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	TagID        int64 `json:"tag_id"`
	TagType      int   `json:"tag_type"`
	PageVO
}
type ExportInfo struct {
	Total    int64          `json:"total"`
	UserList []InfoUserList `json:"user_list"`
}

type InfoUserList struct {
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	UnionID string `json:"union_id"`
}

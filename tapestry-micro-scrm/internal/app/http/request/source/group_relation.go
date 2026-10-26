package source

// 商品关联过滤常量（引用models中的定义）
const (
	GoodsRelationLinked = 1
	GoodsRelationUnlink = 2
	SourceOriginAll     = 0
	SourceOriginManual  = 1
	SourceOriginKOSSync = 2
)

type GroupRelationListRequest struct {
	EnterpriseID  int64 `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	SourceGroupId int64 `json:"source_group_id" binding:"required" comment:"素材分组ID不能为空"`
	Page          int   `json:"page" comment:"页码"`
	Size          int   `json:"size" comment:"分页数"`
}

type GroupRelationSaveRequest struct {
	EnterpriseID   int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	SourceGroupIds []int64 `json:"source_group_ids" binding:"gt=0,dive,required,min=1" comment:"素材分组ID不能为空"`
	ThirdIds       []int   `json:"third_ids" binding:"gt=0,dive,required,min=1" comment:"门店ID不能为空"`
}

type GroupRelationFindRequest struct {
	EnterpriseID  int64  `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	SourceGroupId int64  `json:"source_group_id"  comment:"素材分组ID"`
	MemberID      int64  `json:"member_id" comment:"员工Id"`
	MemberUserId  string `json:"member_user_id" comment:"企业员工Id"`
	ThirdId       int    `json:"third_id"`
	ThirdIds      []int  `json:"third_ids"`
}

// RequestSelectList 各素材列表请求结构
type RequestSelectList struct {
	EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	SourceGroupID int64   `json:"source_group_id" db:"source_group_id"`
	Keyword       string  `json:"keyword" db:"keyword"`
	CreatedStart  int64   `json:"create_start" db:"create_start"`
	CreatedEnd    int64   `json:"create_end" db:"create_end"`
	Size          int     `json:"size" db:"size"`
	Page          int     `json:"page" db:"page"`
	GroupIDList   []int64 `json:"group_id_list" db:"group_id_list"`
	GroupIDs      []int64 `json:"group_ids" db:"group_ids"`
	CategoryIds   []int64 `json:"category_ids"`
	Mobile        string  `json:"mobile"`
	MemberID      int64   `json:"member_id" comment:"员工Id"`
	MemberUserId  string  `json:"member_user_id" comment:"企业员工Id"`
	Type          int     `json:"type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	ThirdId       int     `json:"third_id"`
	ThirdIds      []int   `json:"third_ids"`
	CollectSearch bool    `json:"collect_search" db:"collect_search" comment:"是否收藏查询"`
	IgnoreIndex   bool    `json:"ignore_index"`
	SkuItem       string  `json:"sku_item"`
	ProductName   string  `json:"product_name"`
}

// RequestChildList 文本素材列表请求结构
type RequestChildList struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	ParentID     int64 `json:"parent_id" db:"parent_id"`
}

// RequestConsoleList 后台管理各素材列表请求结构
type RequestConsoleList struct {
	EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	Keyword       string  `json:"keyword" db:"content"`
	AdminID       int64   `json:"admin_id" db:"admin_id"`
	JustMe        int     `json:"just_me" db:"just_me"`
	CreatedStart  int64   `json:"create_start" db:"create_start"`
	CreatedEnd    int64   `json:"create_end" db:"create_end"`
	Size          int     `json:"size" db:"size"`
	Page          int     `json:"page" db:"page"`
	GroupIDList   []int64 `json:"group_id_list" db:"group_id_list"`
	CategoryIds   []int64 `json:"category_ids"`
	Order         string  `json:"order"`
	Sort          string  `json:"sort"`
	SourceStatus  int     `json:"source_status"`
	SourceGroupID int64   `json:"source_group_id" db:"source_group_id"`
	GroupIDs      []int64 `json:"group_ids" db:"group_ids"`
	Type          int     `json:"type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	IsUse         int     `json:"is_use"`
	SkuItem       string  `json:"sku_item"`
	ProductName   string  `json:"product_name"`
	SkuColor      string  `json:"sku_color"`
	GoodsRelation int     `json:"goods_relation" comment:"商品关联过滤: 0-全部 1-已关联商品 2-未关联商品"`
	SourceOrigin  int     `json:"source_origin" comment:"素材来源: 0-全部 1-手动创建 2-CH同步"`
}

type RequestNewSource struct {
	EnterpriseId int64  `json:"enterprise_id" validate:"required" required:"企业ID必传"`
	Mobile       string `json:"mobile"`
	MemberID     int64  `json:"member_id" comment:"员工Id"`
	ThirdIds     []int  `json:"third_ids"`
}

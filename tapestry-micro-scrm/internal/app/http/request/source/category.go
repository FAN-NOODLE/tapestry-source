package source

type CategoryListRequest struct {
	Name         string `json:"name"`
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
	IsCollect    int    `json:"is_collect" comment:"我的收藏判断，0全部，1我的收藏的分类"`
	MemberId     int64  `json:"member_id" comment:"导购ID"`
	ThirdIds     []int  `json:"third_ids" comment:"第三方ID"`
	OriginH5     bool   `json:"origin_h5" comment:"是否是H5"`
}

type CategoryCreateRequest struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Name         string `json:"name" binding:"required" comment:"名称"`
}

type CategoryUpdateRequest struct {
	Id           int64  `json:"id" binding:"required" comment:"ID不能为空"`
	EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Name         string `json:"name" binding:"required" comment:"名称"`
}

type CategoryUpdateIgnoreIndexRequest struct {
	Id           int64 `json:"id" binding:"required" comment:"ID不能为空"`
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
	IgnoreIndex  int   `json:"ignore_index" comment:"首页忽略该目录 ，0否 1是"`
}

type CategoryRequest struct {
	Id           int64 `json:"id" binding:"required" comment:"ID"`
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

// RemoveCategoryRequest 移除素材关联产品类别
type RemoveCategoryRequest struct {
	SourceID     int64 `json:"source_id" binding:"required" comment:"素材Id"`
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
	CategoryId   int64 `json:"category_id" binding:"required" comment:"产品类别ID"`
}

package source

import "gw-work-scrm/internal/app/http/request"

type BaseRequest struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type Request struct {
	BaseRequest
	StartTime   int64    `json:"start_time" comment:"开始时间"`
	EndTime     int64    `json:"end_time" comment:"结束时间"`
	CategoryIds []uint64 `json:"category_ids" comment:"产品类别"`
	SourceType  int32    `json:"source_type" binding:"required" comment:"素材类型不正确"`
	GroupId     int64    `json:"group_id"`
}

type UploadMediaRequest struct {
	Request
	MemberId int64       `json:"member_id"`
	UserId   string      `json:"userid"`
	Content  string      `json:"content" comment:"文本内容"`
	Sources  []MediaInfo `json:"sources" comment:"资源"`
	Mobile   string      `json:"mobile"`
	ShowType int         `json:"show_type" comment:"展示类型:1单张,2合并"`
	FileType int         `json:"file_type" comment:"文件类型: 1图片,2视频"`
}

type MediaInfo struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Ext         string `json:"ext"`
	Size        int    `json:"size"`
	IsOrigin    int    `json:"is_origin"`
	Sort        int32  `json:"sort" comment:"排序"`
	IsTranscode int    `json:"is_transcode" comment:"是否转码,0为否,1为是"`
	CoverUrl    string `json:"cover_url" comment:"视频第一帧的封面图"`
}
type CheckAuth struct {
	BaseRequest
	UserId string `json:"userid"`
}

type CompressReq struct {
	BaseRequest
	SourceType int32 `json:"source_type"`
	MemberId   int64 `json:"member_id"`
	CompressId int64 `json:"compress_id"`
}

type ReportTimesReq struct {
	Ids          []int    `json:"ids"`
	Type         int      `json:"type" binding:"required" required:"类型"`
	EnterpriseID int64    `json:"enterprise_id" comment:"企业ID"`
	ResourceKeys []string `json:"resource_keys"`
}

type BrandListReq struct {
	BrandCacheReq
	Size     int    `json:"size"`
	LastCode string `json:"last_code"`
}
type BrandCacheReq struct {
	EnterpriseId            int64  `json:"enterprise_id"`
	SkuItem                 string `json:"sku_item"`
	ProductName             string `json:"product_name"`
	UpcCode                 string `json:"upc_code"`
	VariantCode             string `json:"variant_code"`
	StoreCode               string `json:"store_code"`
	Mobile                  string `json:"mobile"`
	IncludeNoneImageProduct bool   `json:"include_none_image_product"`
}

type CollectReq struct {
	request.EnterpriseBaseReq
	SourceInfo []CollectSourceInfoReq `json:"source_info"`
	MemberId   int64                  `json:"member_id" binding:"required"`
	IsCollect  int32                  `json:"is_collect" comment:"是否收藏,0取消收藏，1添加收藏"`
}
type CollectSourceInfoReq struct {
	Type        int32    `json:"type" required:"类型"`
	SourceId    []int64  `json:"source_id"`
	ProductCode []string `json:"product_code"`
}

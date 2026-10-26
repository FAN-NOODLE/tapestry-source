package sourceResp

import (
	source2 "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
)

type ImageList struct {
	source.ResponseConsoleList
	Name            string       `json:"name"`
	URL             string       `json:"url"`
	MediaID         string       `json:"media_id"`
	MediaExpireTime int64        `json:"media_expire_time"`
	MediaStatus     int64        `json:"media_status"`
	MergeId         int64        `json:"merge_id"`
	MergeImage      []MergeImage `json:"merge_image"`
}

type FriendList struct {
	source.ResponseConsoleList
	TextContent string       `json:"text_content"`
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	MergeId     int64        `json:"merge_id"`
	MergeImage  []MergeImage `json:"merge_image"`
}

type MergeImage struct {
	Id              int64                `json:"id"`
	Name            string               `json:"name"`
	URL             string               `json:"url"`
	Sort            int32                `json:"sort"`
	DownloadTimes   int                  `json:"download_times"`
	ForwardingTimes int                  `json:"forwarding_times"`
	IsCollect       int                  `json:"is_collect"`
	CollectNum      int                  `json:"collect_num"`
	RelationGoods   []source2.BrandGoods `json:"relation_goods"`
}

type ImageDetail struct {
	source.ImageCategoryRelationDetail
	RelationGoods []source2.BrandGoods `json:"relation_goods"`
	MergeImage    []MergeImage         `json:"merge_image"`
}

type FriendDetail struct {
	source.FriendDetail
	RelationGoods []source2.BrandGoods `json:"relation_goods"`
	MergeImage    []MergeImage         `json:"merge_image"`
}

type ImageSelectList struct {
	source.ResponseSelectList

	Name            string       `json:"name"`
	URL             string       `json:"url"`
	Size            int64        `json:"size"`
	MediaID         string       `json:"media_id"`
	MediaExpireTime int64        `json:"media_expire_time"`
	MediaStatus     int64        `json:"media_status"`
	DepartmentId    int64        `json:"department_id"`
	SourceGroupID   int64        `json:"source_group_id"`
	MergeId         int64        `json:"merge_id"`
	MergeImage      []MergeImage `json:"merge_image"`
	DownloadTimes   int          `json:"download_times" comment:"下载次数"`
	ForwardingTimes int          `json:"forwarding_times" comment:"转发次数"`
}

type FriendSelectList struct {
	source.ResponseSelectList

	TextContent     string       `json:"text_content" comment:"文本内容"`
	Name            string       `json:"name"  comment:"资源名称"`
	VideoCoverURL   string       `json:"video_cover_url" comment:"视频封面地址"`
	OriginalURL     string       `json:"original_url" comment:"原资源地址"`
	CompressionURL  string       `json:"compression_url" comment:"压缩资源地址"`
	Ext             string       `json:"ext" comment:"文件后缀"`
	Size            int32        `json:"size" comment:"文件大小"`
	Type            int32        `json:"type" comment:"类型 1:图片 2:视频"`
	Sort            int32        `json:"sort" comment:"排序"`
	MergeId         int64        `json:"merge_id"`
	MergeImage      []MergeImage `json:"merge_image"`
	DownloadTimes   int          `json:"download_times" comment:"下载次数"`
	ForwardingTimes int          `json:"forwarding_times" comment:"转发次数"`
	CopyTimes       int          `json:"copy_times" comment:"复制次数"`
	CollectNum      int          `json:"collect_num"`
}

type BrandProductInfo struct {
	ID                 string            `json:"id" db:"id"`
	EnterpriseID       int64             `json:"enterprise_id" db:"enterprise_id"`
	Color              string            `json:"color"`
	StyleCode          string            `json:"style_code"`
	Department         string            `json:"department"`
	SkuTagPrice        string            `json:"sku_tag_price"`
	SilhouetteDesc     string            `json:"silhouette_desc"`
	ShortDescriptionZH string            `json:"short_description_zh"`
	LongDescriptionZH  string            `json:"long_description_zh"`
	ProductName        string            `json:"product_name"`
	SkuItem            string            `json:"sku_item"`
	UpcCode            string            `json:"upc_code"`
	UpcList            []UpcListResp     `json:"upc_list"`
	ProductCode        string            `json:"product_code"`
	StorePrice         float64           `json:"store_price"`
	StoreCode          string            `json:"store_code"`
	BrandSourceList    []BrandSourceInfo `json:"brand_source_list"`
	IsCollect          int               `json:"is_collect"`
	CollectNum         int               `json:"collect_num"`
	CollectAt          int64             `json:"collect_at"`
}

type BrandSourceInfo struct {
	ResourceKey     string `json:"resource_key" db:"resource_key"`
	Sort            int    `json:"sort"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	Name            string `json:"name"`
	Type            int    `json:"type"`
	URL             string `json:"url"`
	Size            int64  `json:"size"`
	Ext             string `json:"ext"`
	IsMain          int    `json:"is_main"`
	MediaID         string `json:"media_id"`
	MediaExpireTime int64  `json:"media_expire_time"`
	MediaStatus     int    `json:"media_status"`
	DownloadTimes   int    `json:"download_times" comment:"下载次数"`
	ForwardingTimes int    `json:"forwarding_times" comment:"转发次数"`
}

type UpcListResp struct {
	UpcCode     string `json:"upc_code"`
	VariantCode string `json:"variant_code"`
	SkuNumber   string `json:"sku_number"`
}

type BrandImage struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Sort            int32  `json:"sort"`
	DownloadTimes   int    `json:"download_times"`
	ForwardingTimes int    `json:"forwarding_times"`
}

type BrandCacheResp struct {
	IsMain      int    `json:"is_main"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	ProductName string `json:"product_name"`
	SkuItem     string `json:"sku_item"`
	UpcCode     string `json:"upc_code"`
	Color       string `json:"color"`
	ResourceKey string `json:"resource_key"`
	Size        int64  `json:"size"`
	Ext         string `json:"ext"`
	SourceNum   int    `json:"source_num"`
}

type CollectNumResp struct {
	CollectNum  int64                  `json:"collect_num"`
	MergeId     int64                  `json:"merge_id"`
	ProductCode string                 `json:"product_code"`
	SourceId    int64                  `json:"source_id"`
	SourceList  []SourceCollectNumResp `json:"source_list"`
	SourceType  int32                  `json:"source_type"`
}

type SourceCollectNumResp struct {
	CollectNum int64 `json:"collect_num"`
	SourceId   int64 `json:"source_id"`
}

type SourceDetail struct {
	source.LinkCategoryRelationDetail
	RelationGoods []source2.BrandGoods `json:"relation_goods"`
}

type VideoDetail struct {
	source.VideoCategoryRelationDetail
	RelationGoods []source2.BrandGoods `json:"relation_goods"`
}

// TalkieImageSearchRegion Talkie图片搜索主体区域结果
type TalkieImageSearchRegion struct {
	Region  string                     `json:"region"`  // 主体区域坐标，格式：x1,x2,y1,y2
	PicList []TalkieImageSearchPicInfo `json:"picList"` // 相似产品列表
}

// TalkieImageSearchPicInfo Talkie图片搜索相似产品信息
type TalkieImageSearchPicInfo struct {
	ProductID  string  `json:"productId"`  // 商品ID [款色]，示例：65129_GDMEL
	PicName    string  `json:"picName"`    // 图片名称
	CategoryID int     `json:"categoryId"` // 商品类目
	Score      float64 `json:"score"`      // 相似度分数（0-1）
	Sku        string  `json:"sku"`        // SKU款号（不带颜色），示例：65129
	Name       string  `json:"name"`       // 商品名称（来自宝尊）
	Upc        string  `json:"upc"`        // 商品名称（来自宝尊）
	Url        string  `json:"url"`        // a0主图地址（来自宝尊）
	Color      string  `json:"color"`      // 商品颜色（来自宝尊）
}

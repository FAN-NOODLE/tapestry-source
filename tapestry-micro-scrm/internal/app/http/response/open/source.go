package openResp

import "micro-scrm/internal/app/models"

// LinkListResponse 对外链接列表响应
type LinkListResponse struct {
	ID          int64        `json:"id"`
	GroupId     int64        `json:"group_id"`
	Title       string       `json:"title"`       // 标题
	Description string       `json:"description"` // 描述
	CoverUrl    string       `json:"cover_url"`   // 封面图
	LinkUrl     string       `json:"link_url"`    // 链接URL
	CreatedAt   int64        `json:"created_at"`  // 创建时间
	UpdatedAt   int64        `json:"updated_at"`  // 更新时间
	BrandGoods  []BrandGoods `json:"brand_goods"`
}
type BrandGoods struct {
	Name  string `json:"name"`
	Sku   string `json:"sku"`
	Upc   string `json:"upc"`
	Url   string `json:"url"`
	Color string `json:"color"`
	// AddType int8   `json:"add_type"`
	// Score   string `json:"score"`
}

type SourceCategory struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// SourceListResponse  素材列表响应
type SourceListResponse struct {
	ID         int64             `json:"id"`
	GroupId    int64             `json:"group_id"`
	MergeId    int64             `json:"merge_id"`
	Name       string            `json:"name"`
	Content    string            `json:"content"`
	CoverURL   string            `json:"cover_url"`
	URL        string            `json:"url"`
	SourceURL  string            `json:"source_url"`
	Ext        string            `json:"ext"`
	Status     int32             `json:"status"`
	StartTime  int               `json:"start_time"`
	EndTime    int               `json:"end_time"`
	CreatedAt  int64             `json:"created_at"`
	UpdatedAt  int64             `json:"updated_at"`
	DeleteAt   models.DeletedAt  `json:"delete_at"`
	BrandGoods []BrandGoods      `json:"brand_goods"`
	Categories *[]SourceCategory `json:"categories,omitempty"`
}

// OpenMergeImage 组合图子图响应（open API 用）
type OpenMergeImage struct {
	Id         int64        `json:"id"`
	Name       string       `json:"name"`
	URL        string       `json:"url"`
	Sort       int32        `json:"sort"`
	BrandGoods []BrandGoods `json:"brand_goods"`
}

// TextListResponse 文本列表响应
type TextListResponse struct {
	ID         int64        `json:"id"`
	GroupID    int64        `json:"group_id"`
	Content    string       `json:"content"`
	CreatedAt  int64        `json:"created_at"`
	UpdatedAt  int64        `json:"updated_at"`
	BrandGoods []BrandGoods `json:"brand_goods"`
}

// ImageListResponse 图片列表响应
type ImageListResponse struct {
	ID         int64            `json:"id"`
	GroupID    int64            `json:"group_id"`
	MergeId    int64            `json:"merge_id"`
	Name       string           `json:"name"`
	URL        string           `json:"url"`
	CreatedAt  int64            `json:"created_at"`
	UpdatedAt  int64            `json:"updated_at"`
	MergeImage []OpenMergeImage `json:"merge_image"`
	BrandGoods []BrandGoods     `json:"brand_goods"`
}

// VideoListResponse 视频列表响应
type VideoListResponse struct {
	ID         int64        `json:"id"`
	GroupID    int64        `json:"group_id"`
	Name       string       `json:"name"`
	CoverUrl   string       `json:"cover_url"`
	URL        string       `json:"url"`
	CreatedAt  int64        `json:"created_at"`
	UpdatedAt  int64        `json:"updated_at"`
	BrandGoods []BrandGoods `json:"brand_goods"`
}

// FriendListResponse 朋友圈列表响应
type FriendListResponse struct {
	ID            int64            `json:"id"`
	GroupID       int64            `json:"group_id"`
	MergeId       int64            `json:"merge_id"`
	TextContent   string           `json:"text_content"`
	Name          string           `json:"name"`
	VideoCoverUrl string           `json:"video_cover_url"`
	OriginalUrl   string           `json:"original_url"`
	Type          int32            `json:"type"`
	CreatedAt     int64            `json:"created_at"`
	UpdatedAt     int64            `json:"updated_at"`
	MergeImage    []OpenMergeImage `json:"merge_image"`
	BrandGoods    []BrandGoods     `json:"brand_goods"`
}

// SearchAllResponse 合并检索响应（各素材类型独立 list）
type SearchAllResponse struct {
	TextList   []TextListResponse   `json:"text_list"`
	ImageList  []ImageListResponse  `json:"image_list"`
	LinkList   []LinkListResponse   `json:"link_list"`
	VideoList  []VideoListResponse  `json:"video_list"`
	FriendList []FriendListResponse `json:"friend_list"`
}

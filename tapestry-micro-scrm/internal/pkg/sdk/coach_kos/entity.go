package coachkos

import "encoding/json"

// PostResult 是 COACH KOS 帖子结果。
type PostResult struct {
	PostID           string          `json:"post_id"`
	SourceURL        string          `json:"source_url"`
	DataDate         string          `json:"data_date"`
	PublishedAt      string          `json:"published_at"`
	Name             string          `json:"name"`
	Content          string          `json:"content"`
	Hashtag          string          `json:"hashtag"`
	URL              string          `json:"url"`
	CoverURL         string          `json:"cover_url"`
	InfluencerID     string          `json:"influencer_id"`
	InfluencerName   string          `json:"influencer_name"`
	EmployeeID       string          `json:"employee_id"`
	EmployeeTag      string          `json:"employee_tag"`
	StoreID          string          `json:"store_id"`
	StoreName        string          `json:"store_name"`
	ProductID        []string        `json:"product_id"`
	ContentTag       string          `json:"content_tag"`
	LikeCount        int64           `json:"like_count"`
	CommentCount     int64           `json:"comment_count"`
	FavoriteCount    int64           `json:"favorite_count"`
	ShareCount       int64           `json:"share_count"`
	InteractionCount int64           `json:"interaction_count"`
	IsValid          string          `json:"is_valid"`
	CommentSentiment json.RawMessage `json:"comment_sentiment"`
	PriceMentioned   string          `json:"price_mentioned"`
	UpdatedAt        string          `json:"updated_at"`
	SourceUpdatedAt  string          `json:"source_updated_at"`
}

// ListRequest 是帖子结果列表请求。
type ListRequest struct {
	UpdatedAt string
	Page      int
	Size      int
}

// PageData 是帖子结果分页数据。
type PageData struct {
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	List  []PostResult `json:"list"`
}

type listResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    PageData `json:"data"`
}

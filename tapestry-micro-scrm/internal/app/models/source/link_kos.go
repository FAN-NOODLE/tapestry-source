package source

import "micro-scrm/internal/app/models"

type LinkKOSModel struct {
	models.MicroBaseModel
	SourceID         int64  `gorm:"column:source_id;type:bigint(20);not null"`
	SourceLinkID     int64  `gorm:"column:source_link_id;type:bigint(20);not null"`
	PostID           string `gorm:"column:post_id;type:varchar(128);not null"`
	SourceURL        string `gorm:"column:source_url;type:varchar(1000);not null"`
	DataDate         string `gorm:"column:data_date;type:varchar(32);not null"`
	PublishedAt      string `gorm:"column:published_at;type:varchar(32);not null"`
	Name             string `gorm:"column:name;type:varchar(255);not null"`
	Content          string `gorm:"column:content;type:text;not null"`
	Hashtag          string `gorm:"column:hashtag;type:text;not null"`
	URL              string `gorm:"column:url;type:varchar(1000);not null"`
	CoverURL         string `gorm:"column:cover_url;type:varchar(1000);not null"`
	InfluencerID     string `gorm:"column:influencer_id;type:varchar(128);not null"`
	InfluencerName   string `gorm:"column:influencer_name;type:varchar(255);not null"`
	EmployeeID       string `gorm:"column:employee_id;type:varchar(128);not null"`
	EmployeeTag      string `gorm:"column:employee_tag;type:varchar(255);not null"`
	StoreID          string `gorm:"column:store_id;type:varchar(128);not null"`
	StoreName        string `gorm:"column:store_name;type:varchar(255);not null"`
	ProductID        string `gorm:"column:product_id;type:text;not null"`
	ContentTag       string `gorm:"column:content_tag;type:text;not null"`
	LikeCount        int64  `gorm:"column:like_count;type:bigint(20);not null"`
	CommentCount     int64  `gorm:"column:comment_count;type:bigint(20);not null"`
	FavoriteCount    int64  `gorm:"column:favorite_count;type:bigint(20);not null"`
	ShareCount       int64  `gorm:"column:share_count;type:bigint(20);not null"`
	InteractionCount int64  `gorm:"column:interaction_count;type:bigint(20);not null"`
	IsValid          string `gorm:"column:is_valid;type:varchar(8);not null"`
	CommentSentiment string `gorm:"column:comment_sentiment;type:text;not null"`
	PriceMentioned   string `gorm:"column:price_mentioned;type:varchar(255);not null"`
	APIUpdatedAt     string `gorm:"column:api_updated_at;type:varchar(32);not null"`
	SourceUpdatedAt  string `gorm:"column:source_updated_at;type:varchar(32);not null"`
}

func (LinkKOSModel) TableName() string { return "source_link_kos" }

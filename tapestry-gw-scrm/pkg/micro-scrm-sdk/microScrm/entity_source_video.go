package microScrm

// RequestSourceVideo 视频素材新增
type RequestSourceVideo struct {
	ID            int64        `json:"id" db:"id"`
	EnterpriseID  int64        `json:"enterprise_id" db:"enterprise_id"`
	GroupID       int64        `json:"group_id" db:"group_id"`
	CoverUrl      string       `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
	Name          string       `json:"name" validate:"required" required:"视频名称必传" db:"name"`
	Size          int64        `json:"size" validate:"required" required:"视频大小必传" db:"size"`
	Url           string       `json:"url"  validate:"required" required:"视频Url必传" db:"url"`
	Ext           string       `json:"ext" validate:"required" required:"视频类型必传" db:"ext"`
	AdminID       int64        `json:"admin_id" db:"admin_id"`
	CreatedAt     int64        `json:"created_at" db:"created_at"`
	UpdatedAt     int64        `json:"updated_at" db:"updated_at"`
	CategoryIds   []int64      `json:"category_ids"`
	StartTime     int64        `json:"start_time"`
	EndTime       int64        `json:"end_time"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}

// SourceVideo 视频素材
type SourceVideo struct {
	ID           int64  `json:"id" db:"id"`
	Type         int    `json:"type" db:"type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	GroupID      int64  `json:"group_id" db:"group_id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Name         string `json:"name"  db:"name"`
	CoverUrl     string `json:"cover_url" db:"cover_url"`
	Size         int64  `json:"size"  db:"size"`
	Ext          string `json:"ext"  db:"ext"`
	Url          string `json:"url"  db:"url"`
	Status       int    `json:"status"  db:"status"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

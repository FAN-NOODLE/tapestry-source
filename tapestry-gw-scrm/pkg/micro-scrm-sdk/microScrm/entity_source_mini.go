package microScrm

// RequestSourceMini 小程序素材新增
type RequestSourceMini struct {
	ID            int64   `json:"id" db:"id"`
	EnterpriseID  int64   `json:"enterprise_id"`
	GroupID       int64   `json:"group_id" db:"group_id"`
	Title         string  `json:"title" required:"卡片标题必传" db:"title"`
	MiniProgramID int64   `json:"mini_program_id" validate:"required" required:"小程序ID必传" db:"mini_program_id"`
	Path          string  `json:"path"  validate:"required" required:"链接类型必传" db:"path"`
	CoverUrl      string  `json:"cover_url" validate:"required" required:"封面图URL必传" db:"cover_url"`
	AdminID       int64   `json:"admin_id" db:"admin_id"`
	CreatedAt     int64   `json:"created_at" db:"created_at"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
	CategoryIds   []int64 `json:"category_ids"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
}

// 小程序素材
type SourceMini struct {
	ID            int64   `json:"id" db:"id"`
	Type          int     `json:"type" db:"type"`
	AdminID       int64   `json:"admin_id" db:"admin_id"`
	GroupID       int64   `json:"group_id" db:"group_id"`
	EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id"`
	Title         string  `json:"title" db:"title"`
	MiniProgramID int64   `json:"mini_program_id" db:"mini_program_id"`
	AppName       string  `json:"app_name"  db:"app_name"`
	AppID         string  `json:"app_id"  db:"app_id"`
	OriginID      string  `json:"origin_id" db:"origin_id"`
	Logo          string  `json:"logo" db:"logo"`
	Path          string  `json:"path" db:"path"`
	CoverUrl      string  `json:"cover_url" db:"cover_url"`
	CreatedAt     int64   `json:"created_at" db:"created_at"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
	DeletedAt     int64   `json:"deleted_at" db:"deleted_at"`
	CategoryIds   []int64 `json:"category_ids"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
}

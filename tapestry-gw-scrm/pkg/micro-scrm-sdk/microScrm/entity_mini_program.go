package microScrm

type MiniProgram struct {
	ID           int64  `json:"id" db:"id"`
	AppName      string `json:"app_name" db:"app_name" validate:"required" required:"请输入小程序名称"`
	AppID        string `json:"app_id" db:"app_id" validate:"required" required:"请输入小程序APP ID"`
	Logo         string `json:"logo" db:"logo"`
	OriginID     string `json:"origin_id" db:"origin_id" validate:"required" required:"请输入小程序原始ID"`
	Type         int    `json:"type" db:"type"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
}

type ProgramList []struct {
	ID       int64  `json:"id" db:"id"`
	AppName  string `json:"app_name" db:"app_name"`
	Logo     string `json:"logo" db:"logo"`
	Type     int    `json:"type" db:"type"`
	AppID    string `json:"app_id" db:"app_id"`
	OriginID string `json:"origin_id" db:"origin_id"`
}

type RequestPrograms struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	Page         int   `json:"page"`
	Size         int   `json:"size"`
}

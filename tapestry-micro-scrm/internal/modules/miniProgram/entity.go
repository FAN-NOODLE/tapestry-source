package miniProgram

type MiniProgram struct {
	ID           int64  `json:"id" db:"id"`
	AppName      string `json:"app_name" db:"app_name" validate:"required" required:"请输入小程序名称"`
	AppID        string `json:"app_id" db:"app_id" validate:"required" required:"请输入小程序APP ID"`
	OriginID     string `json:"origin_id" db:"origin_id" validate:"required" required:"请输入小程序原始ID"`
	Logo         string `json:"logo" db:"logo"`
	Type         int    `json:"type" db:"type"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
}

type ProgramList []struct {
	ID       int64  `json:"id" db:"id"`
	AppName  string `json:"app_name" db:"app_name"`
	Type     int    `json:"type" db:"type"`
	AppID    string `json:"app_id" db:"app_id"`
	OriginID string `json:"origin_id" db:"origin_id"`
	Logo     string `json:"logo" db:"logo"`
}

type RequestPrograms struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	Page         int   `json:"page"`
	Size         int   `json:"size"`
}

type ID struct {
	ID int64 `json:"id"`
}

type Rows struct {
	Rows int64 `json:"rows"`
}

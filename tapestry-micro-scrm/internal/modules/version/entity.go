package version

// GetCurrentVersionInfoRequest.Type的值
// 0 表示 获取当前版本信息的内容
// 1 表示 获取当前本本信息的内容，并记录版本信息提示记录；如果用户已经看过提示信息
const (
	CurrentVersionInfoRequest_Type_Get = iota
	CurrentVersionInfoRequest_Type_Record
)

// GetCurrentVersionInfoRequest 获取当前版本信息的请求参数
type GetCurrentVersionInfoRequest struct {
	Type         uint8 `json:"type"`
	AdminID      int64 `json:"admin_id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type GetCurrentVersionInfoResponse struct {
	Type        uint8       `json:"type"`
	IsHasShow   bool        `json:"is_has_show"`
	VersionInfo VersionInfo `json:"version_info"`
}

type VersionInfo struct {
	Them   string          `json:"them"`
	Number string          `json:"number"`
	Date   string          `json:"date"`
	Detail []VersionDetail `json:"detail"`
}

// VersionDetail 用来解析version_record.version_info
type VersionDetail struct {
	Title  string   `json:"title"`
	Detail []string `json:"detail"`
}

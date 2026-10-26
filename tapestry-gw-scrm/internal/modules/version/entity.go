package version

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

package entity

// 客户分组列表数据库查询
type UserGroupInfo struct {
	ID        int64  `db:"user_group_id"`
	GroupName string `db:"group_name"`
	UseNum    int    `db:"use_num"`
}

// 客户分组标签信息
type UserGroupDetailTagInfo struct {
	UserGroupID   int64  `db:"user_group_id" json:"user_group_id"`
	UserGroupName string `db:"user_group_name" json:"user_group_name"`
	TagGroupID    int64  `db:"tag_group_id" json:"tag_group_id"`
	TagGroupName  string `db:"tag_group_name" json:"tag_group_name"`
	TagID         int64  `db:"tag_id" json:"tag_id"`
	TagName       string `db:"tag_name" json:"tag_name"`
	IncludeType   int8   `db:"include_type" json:"include_type"`
	UserGroupType int8   `db:"user_group_type" json:"user_group_type"`
	FileID        int64  `db:"file_id" json:"file_id"`
	DepartmentID  int64  `db:"department_id" json:"department_id"`
	IsExtract     int    `db:"is_extract" json:"is_extract"`
}

type ExportReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	UserGroupID  int64 `json:"user_group_id"`
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Total        int64 `json:"total"`
}

type FileDetection struct {
	FileName     string `json:"file_name"`
	EnterpriseID int64  `json:"enterprise_id"`
}

type CountInfo struct {
	Count int64 `json:"id"`
}

type ExportList struct {
	ExportData []ExportDataInfo
	Total      int64 `json:"total"`
}
type ExportDataInfo struct {
	UnionID   string `json:"union_id"`
	ImportRes string `json:"import_res"`
}

type DownLoadFileInfo struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
	FileUrl      string `json:"file_url"`
}

type IDInfo struct {
	ID int64 `json:"id"`
}

type GetFileInfo struct {
	EnterpriseID int64       `json:"enterprise_id"`
	ID           int64       `json:"id"`
	FileName     string      `json:"file_name"`
	File         []FileField `json:"file"` //csv 文件
}

type FileField struct {
	Unionid string `json:"unionid"`
}

type UpdateIdInfo struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type ImportData struct {
	ID int64 `db:"id" json:"id"`
}

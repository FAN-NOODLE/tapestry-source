package delivery

type GetFileInfo struct {
	TagName        string      `json:"tag_name"`
	EnterpriseID   int64       `json:"enterprise_id"`
	ID             int64       `json:"id"`
	FileName       string      `json:"file_name"`
	File           []FileField `json:"file"`             //csv 文件
	UnionIDEncrypt string      `json:"union_id_encrypt"` //union是否加密 1.加密过 2.未加密
}

type FileField struct {
	Ccode     string `json:"ccode"`
	Unionid   string `json:"unionid"`
	GradeCode string `json:"grade_code"`
	AdminBA   string `json:"admin_ba"`
}

type DownLoadFileInfo struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
	TagName      string `json:"tag_name"` //标签名称
	FileUrl      string `json:"file_url"`
	TagGroupName string `json:"tag_group_name"` //标签组
	BAType       string `json:"ba_type"`        //BA类型  1.普通BA 2.管理BA
	TaskType     int    `json:"task_type"`
}

type IDInfo struct {
	ID int64 `json:"id"`
}

type FileDetection struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
}

type CountInfo struct {
	Count int64 `json:"id"`
}

type NewIDInfo struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

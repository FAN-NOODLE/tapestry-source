package plan_task

import "micro-scrm/internal/orm/model"

// 计划任务数据文件记录表模型

// 上传方式
const (
	PTDF_FileType_WS = 1 + iota
	PTDF_FileType_CT
)

// 打标签的BA类型
const (
	PTDF_BaType_Relation = 1 + iota
	PTDF_BaType_Manager
)

// 数据的流程
const (
	PTDF_IsExtract_SFTPData = 1 + iota
	PTDF_IsExtract_PreData
	PTDF_IsExtract_Tag
	PTDF_IsExtract_SyncTag
)

type PlanTaskDownloadedFileModel struct {
	model.BranchModel

	FileName     string `json:"file_name" gorm:"column:file_name;type:varchar(300);not null;comment:文件名"`
	FileType     uint8  `json:"file_type" gorm:"column:file_type;type:tinyint(1);default:0;comment:上传的文件方式类型：1、wesocial上传；2、campaign_tool上传"`
	FileURL      string `json:"file_url" gorm:"column:file_url;type:varchar(300);not null;comment:文件在blog的地址"`
	GroupTagName string `json:"group_tag_name" gorm:"column:group_tag_name;type:varchar(100);not null;comment:标签组名称"`
	IsExtract    uint8  `json:"is_extract" gorm:"column:is_extract;type:tinyint(1);default:0;comment:数据文件是否已经提取完成：1: 文件数据提取完成；2、预处理数据完成；3、数据打标签完成"`
	TagName      string `json:"tag_name" gorm:"column:tag_name;type:varchar(100);not null;comment:标签名称"`
	BaType       uint8  `json:"ba_type" gorm:"column:ba_type;type:tinyint(1);default:0;comment:BA的类型。1、关系BA；2、管理BA"`
	IsShowTag    uint8  `json:"is_show_tag" gorm:"column:is_show_tag;type:tinyint(1);default:0;comment:BA的类型。0不展示，1展示"`
}

func (PlanTaskDownloadedFileModel) TableName() string {
	return "plan_task_downloaded_file"
}

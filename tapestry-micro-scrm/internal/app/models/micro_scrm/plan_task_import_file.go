package microScrmModel

import "micro-scrm/internal/app/models"

type TaskImportFileModel struct {
	models.BaseModel
	ID          int64  `json:"id" gorm:"primaryKey;column:id;type:int(11);autoIncrement"`
	FileName    string `json:"file_name" gorm:"column:file_name;type:varchar(300);not null;comment:文件名"`
	FileType    int32  `json:"file_type" gorm:"column:file_type;type:tinyint(1);default:0;comment:上传的文件方式类型：1、wesocial上传；2、campaign_tool上传"`
	FileURL     string `json:"file_url" gorm:"column:file_url;type:varchar(300);not null;comment:文件在cos的地址"`
	IsExtract   int32  `json:"is_extract" gorm:"column:is_extract;type:tinyint(1);default:0;comment:数据文件是否已经提取完成：1: 文件数据提取完成；2、预处理数据完成；3、数据打标签完成"`
	TaskType    int    `json:"task_type" gorm:"column:task_type;type:tinyint(1);default:3;comment:导入的任务类型，3为客户导入，4为客户匹配SaCode导入"`
	DownloadURL string `json:"download_url" gorm:"column:download_url;type:varchar(255);not null;comment:匹配完的结果URL地址"`
}

// TableName returns the table name of the TaskImportFile model
func (TaskImportFileModel) TableName() string {
	return "plan_task_import_file"
}

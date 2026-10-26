package crmSftpFileModels

import "micro-scrm/internal/app/models"

type CrmSftpDownloadedFileModel struct {
	models.MicroBaseModel
	FileName     string `json:"file_name" gorm:"column:file_name;type:varchar(300);not null;index:file_name;comment:文件名"`
	FileFullPath string `json:"file_full_path" gorm:"column:file_full_path;type:varchar(300);not null;comment:sftp路径"`
	FileType     int32  `json:"file_type" gorm:"column:file_type;type:tinyint(1);default:0;comment:文件类型：1、标签组；2、标签3、标签与客户关系；4专属导购 5 印象标签日志"`
	FileURL      string `json:"local_file_path" gorm:"column:file_url;type:varchar(300);not null;comment:解密后文件真实地址"`
	IsExtract    int32  `json:"is_extract" gorm:"column:is_extract;type:tinyint(1);default:0;comment:数据文件是否已经提取完成：1: 文件数据提取完成；2、预处理数据完成；3、数据打标签完成"`
}

// TableName returns the table name of the CrmSftpDownloadedFile model
func (CrmSftpDownloadedFileModel) TableName() string {
	return "crm_sftp_downloaded_file"
}

const (
	TAG_GROUP               = iota + 1 // 1=标签组
	TAG                                // 2=标签
	TAG_CUSTOMER_REL                   // 3=标签与客户关系
	KAY_SA                             // 4=专属导购
	Impress_Tag                        // 5=印象标签
	Member_Info                        // 6=coach的用户信息同步
	Member_Info_SW                     // 7=SW的用户信息同步
	Member_Info_KS                     // 8=KS的用户信息同步
	Customer_Purchase_SW               // 9=SW的订单信息同步
	Customer_Extension_SW              // 10=SW的会员标签信息同步
	Store_Master_LMS                   // 11=Coach LMS门店主数据
	Employee_Master_Workday            // 12=Coach Workday员工主数据
)

package micro_scrm

import "micro-scrm/internal/orm/model"

const (
	UGD_IncludeType_Include = 1 + iota
	UGD_IncludeType_Filter
)

type UserGroupDetailModel struct {
	model.MicroScrmBaseModel

	UserGroupID   int64                     `gorm:"column:user_group_id;type:bigint(20) unsigned;not null;default:0;comment:客户分组id"`
	UserGroup     UserGroupModel            `gorm:"foreignKey:UserGroupID"`
	TagGroupID    int64                     `gorm:"column:tag_group_id;type:bigint(20) unsigned;not null;default:0;comment:work_external_tag_group.id"`
	TagGroup      WorkExternalTagGroupModel `gorm:"foreignKey:TagGroupID"`
	TagID         int64                     `gorm:"column:tag_id;type:bigint(20) unsigned;not null;default:0;comment:work_external_tag_group.id"`
	Tag           WorkExternalTagModel      `gorm:"foreignKey:TagID"`
	IncludeType   int32                     `gorm:"column:include_type;type:tinyint(3) unsigned;not null;default:0;comment:1包含，2排除"`
	UserGroupType int                       `gorm:"column:user_group_type;type:tinyint(1) ;not null;default:0;comment:分组的维度；1标签分组，2客户导入，3门店"`
	FileID        int64                     `gorm:"column:file_id;type:int(11);not null;default:0;comment:客户导入的文件id"`
	DepartmentID  int64                     `gorm:"column:department_id;type:bigint(20);not null;default:0;comment:企业部门的门店id(部门id)"`
	IsExtract     int                       `gorm:"column:is_extract;type:tinyint(1);default:0;comment:名单导入的unionid是否成功入库，0:否,1:是"`
}

// TableName returns the table name of the struct_name model
func (UserGroupDetailModel) TableName() string {
	return "user_group_detail"
}

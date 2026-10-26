package externalSearchModel

type ExternalCrowdDataModel struct {
	CrowdId            int64 `json:"crowd_id" gorm:"column:crowd_id;type:bigint(10) unsigned;not null;default:0;comment:crowd表id"`
	WorkExternalUserId int64 `json:"work_external_user_id" gorm:"column:work_external_user_id;type:bigint(10);comment:外部联系人id"`
	IsContact          int   `json:"is_contact" gorm:"column:is_contact;type:bigint(10);comment:是否联系 1已联系"`
	CreatedAt          int64 `json:"created_at" gorm:"column:created_at;type:bigint(10);comment:插入时间"`
}

// TableName returns the table name of the ExcitationAchievement model
func (ExternalCrowdDataModel) TableName() string {
	return "external_crowd_data"
}

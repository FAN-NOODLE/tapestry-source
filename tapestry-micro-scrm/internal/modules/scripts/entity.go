package scripts

type Data struct {
	Id  int64 `db:"id" desc:"id"`
	Num int64 `db:"num" desc:"重复数"`
}

type ExternalUserRelationData struct {
	EnterpriseID       int64  `db:"enterprise_id" desc:"企业id"`
	WorkExternalUserID int64  `db:"work_external_user_id" desc:"外部联系人id"`
	Description        string `db:"description" desc:"员工给客户的备注"`
}

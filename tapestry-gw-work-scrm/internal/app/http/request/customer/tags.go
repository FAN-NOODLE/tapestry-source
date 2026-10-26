package customerReq

type CreateTags struct {
	EnterpriseID       int64   `json:"enterprise_id" comment:"企业ID"`
	EnterpriseUserID   int64   `json:"enterprise_user_id" comment:"客户ID"`
	AddTagIds          []int64 `json:"add_tag_ids" comment:"添加标签ID"`
	DelTagIds          []int64 `json:"del_tag_ids" comment:"删除标签ID"`
	OperateUserId      string  `json:"operate_user_id" comment:"操作人Id(企业微信userid)"`
	OperateName        string  `json:"operate_name" comment:"操作人姓名"`
	OperateMobile      string  `json:"operate_mobile" comment:"操作人手机号"`
	ExternalUserID     string  `json:"external_userid"  comment:"外部联系人userid不能为空"`
	WorkExternalUserId int64   `json:"work_external_user_id"comment:"客户ID(兼容处理micro-scrm的字段调整)"`
}

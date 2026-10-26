package customerReq

type CreateTags struct {
	EnterpriseID       int64   `json:"enterprise_id" binding:"required" comment:"企业ID"`
	WorkExternalUserId int64   `json:"work_external_user_id" binding:"required" comment:"客户ID"`
	AddTagIds          []int64 `json:"add_tag_ids"  comment:"添加标签ID"`
	DelTagIds          []int64 `json:"del_tag_ids"  comment:"删除标签ID"`
	OperateUserId      string  `json:"operate_user_id" comment:"操作人Id(企业微信userid)"`
	OperateName        string  `json:"operate_name" comment:"操作人姓名"`
	OperateMobile      string  `json:"operate_mobile" comment:"操作人手机号"`
}

type CreateTagOperateLogReq struct {
	EnterpriseID       int64    `json:"enterprise_id"`
	WorkExternalUserId int64    `json:"work_external_user_id"`
	ThirdTagIds        []string `json:"third_tag_ids"`
	OperateUserId      string   `json:"operate_user_id"`
	OperateName        string   `json:"operate_name"`
	OperateMobile      string   `json:"operate_mobile"`
	ActionType         int32    `json:"action_type"`
}

type OperateWorkTags struct {
	EnterpriseID       int64   `json:"enterprise_id" binding:"required" comment:"企业ID"`
	WorkExternalUserId int64   `json:"work_external_user_id" binding:"required" comment:"客户ID"`
	GroupId            int64   `json:"group_id" comment:"企业标签群组id"`
	AddTagIds          []int64 `json:"add_tag_ids"  comment:"添加标签ID"`
	DelTagIds          []int64 `json:"del_tag_ids"  comment:"删除标签ID"`
	OperateUserId      string  `json:"operate_user_id" comment:"操作人Id(企业微信userid)"`
	OperateMobile      string  `json:"operate_mobile" comment:"操作人手机号"`
	MemberId           int64   `json:"member_id"`
}

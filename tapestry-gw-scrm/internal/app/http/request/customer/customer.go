package customerReq

type UserList struct {
	EnterpriseID      int64    `json:"enterprise_id" comment:"企业ID"`
	ManageType        int      `json:"manage_type"`
	MemberList        []int64  `json:"member_list"`
	CardType          []string `json:"card_type"`
	ImpressionTagList []int64  `json:"impression_tag_list"`
	TagList           []int64  `json:"tag_list"`
	CltType           []string `json:"clt_type"`
	Mobile            string   `json:"mobile"`
	Name              string   `json:"name"`
	Remark            string   `json:"remark"`
	Page              int      `json:"page" comment:"页码"`
	Size              int      `json:"size" comment:"分页数"`
}

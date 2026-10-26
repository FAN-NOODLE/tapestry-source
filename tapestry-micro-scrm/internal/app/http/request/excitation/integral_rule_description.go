package excitationReq

//DescriptionCreateOrUpdateReq 添加规则说明
type DescriptionCreateOrUpdateReq struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required"`
	Content      string `json:"content" binding:"required" comment:"规则说明不能不空"`
}

type DescriptionReq struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required"`
}

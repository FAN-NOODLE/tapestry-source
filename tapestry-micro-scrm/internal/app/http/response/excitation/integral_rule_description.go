package excitationResp

//DescriptionResp 规则说明
type DescriptionResp struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required"`
	Content      string `json:"content" binding:"required" comment:"规则说明不能不空"`
}

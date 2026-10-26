package transfer_customer

// 在职继承业务实体

type OnTransfer_Common struct {
	EnterpriseID     int64 `json:"enterprise_id" binding:"required" desc:"企业ID"`
	HandOverMemberID int64 `json:"hand_over_member_id" binding:"required" desc:"原跟进成员的member_id"`
	TakeOverMemberID int64 `json:"take_over_member_id" binding:"required" desc:"接替成员的member_id"`
}

// OnTransfer_AllotCustomerRequest 分配客户的请求参数
type OnTransfer_AllotCustomerRequest struct {
	OnTransfer_Common
	ExternalUserIDList []int64 `json:"external_user_id_list" binding:"required" desc:"被接替的客户的ID"`
}

// OnTransfer_QueryCustomerStatusRequest 查询接替成员状态的请求参数
type OnTransfer_QueryCustomerStatusRequest struct {
	OnTransfer_Common
}

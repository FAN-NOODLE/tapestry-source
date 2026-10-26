package delivery

import "micro-scrm/internal/modules/customer"

func TransListParams(params *customer.RepeatListRequest) *customer.RepeatListParams {
	serviceParams := &customer.RepeatListParams{
		EnterpriseID: params.EnterpriseID,
		Page:         params.Page,
		Size:         params.Size,
		CustomerName: params.Search.CustomerName,
		DepartmentID: params.Search.DepartmentID,
		MemberID:     params.Search.MemberID,
	}
	return serviceParams
}
